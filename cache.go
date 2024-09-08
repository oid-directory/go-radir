package radir

/*
cache.go offers a generic, thread-safe, in-memory Registration/Registrant
caching subsystem.
*/

import (
	"encoding/gob"
	"os"
	"sync"
	"time"
)

/*
cachedRegistration contains a *[Registration] instance alongside an
expiry [time.Time] instance.

Instances of this type are stored within an instance of *[Cache] and
need not be managed directly by the user.
*/
type cachedRegistration struct {
	Value  *Registration
	Expiry time.Time
}

/*
cachedRegistrant contains a *[Registrant] instance alongside
an expiry [time.Time] instance.

Instances of this type are stored within an instance of *[Cache] and
need not be managed directly by the user.
*/
type cachedRegistrant struct {
	Value  *Registrant
	Expiry time.Time
}

/*
Expired returns a Boolean value indicative of whether the receiver
instance has expired. A value of true is returned if the instance
cannot be found.

Note that, unlike [Cache.Registration], use of this method will not
result in the deletion of an expired instance.
*/
func (r *Cache) RegistrationExpired(dn string) bool {
	return r.expired(dn, 0)
}

/*
Expired returns a Boolean value indicative of whether the receiver
instance has expired. A value of true is returned if the instance
cannot be found.

Note that, unlike [Cache.Registrant], use of this method will not
result in the deletion of an expired instance.
*/
func (r *Cache) RegistrantExpired(dn string) bool {
	return r.expired(dn, 1)
}

func (r *Cache) expired(dn string, t int) bool {
	if r.IsZero() {
		return false
	}

	if len(dn) == 0 {
		return false
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	key := lc(dn)

	if t == 0 {
		if item, _ := r.registrations[key]; item.Value != nil {
			return now().After(item.Expiry)
		}
		return true
	}

	if item, _ := r.registrants[key]; item.Value != nil {
		return now().After(item.Expiry)
	}

	return true
}

/*
Cache is a thread-safe, memory-based caching type, meant to store any
number of *[Registrant] and *[Registration] instances for the
purpose of reducing LDAP utilization. Caching is covered in the [RADUA]
and [RASCHEMA IDs].

Instances of either type are associated with their respective LDAP DNs,
which are queried immediately prior to an LDAP Search.

Unexpired instances found within a *[Cache] that match the queried DN are
returned instead of reaching out to the RA DSA.

Requesting cached instances that have since expired will result in their
immediate annihilation and a nil return. Under ordinary circumstances,
at this point the DUA could reach out to the directory in an attempt to
re-acquire the now-expired instance.

The [NewCache] function initializes and returns instances of this type.

Following initialization, an instance of *[Cache] may be written to file
using the [Cache.WriteRegistrations] and [Cache.WriteRegistrants] methods.

Conversely, the [Cache.LoadRegistrations] and [Cache.LoadRegistrants]
methods allow an unfrozen *[Cache] to be loaded from file.

The [Cache.Registration] and [Cache.Registrant] methods allow accessing
unexpired cached instances. Attempting to access a cached instance that
has since expired will result in its deletion unless the *[Cache] has
been frozen.

The [Cache.RegistrationLen] and [Cache.RegistrantLen] methods return the
current respective integer length of a *[Cache]. The [Cache.RegistrationCap]
and [Cache.RegistrantCap] return the maximum number of elements within a
*[Cache].

The [Cache.RegistrationExpired] and [Cache.RegistrantExpired] methods
safely allow expiration checks of specific cached instances without the
risk of deletion. The [Cache.RegistrationKeys] and [Cache.RegistrantKeys]
return string slices of cached element DNs, also without deletion risk.

The [Cache.Add] method allows the addition of elements into the *[Cache],
provided it is not frozen.

The [Cache.RemoveRegistration] and [Cache.RemoveRegistrant] methods allow
for the removal of select instances from an unfrozen *[Cache], regardless
of the expiration status.

The [Cache.TouchRegistration] and [Cache.TouchRegistrant] methods allow
for expired -- but as-of-yet undeleted -- cached instances to be "reborn"
or "resurrected" if the *[Cache] is not frozen.  A "touch" effectively
results in the reset of the respective "expiration timer" to the specified
duration.

The [Cache.Freeze] and [Cache.Thaw] methods impose read-only and read-write
policies respectively, influencing the ability for updates and amendments
to the *[Cache] to be recognized. The [Cache.Frozen] method offers a means
for checking the frozen state of a *[Cache]. A frozen *[Cache] will always
permit read operations. When first initialized, a *[Cache] is always in a
thawed state.

[Cache.Tidy] and [Cache.Flush] can be used to clean-up or outright purge
multiple instances from an unfrozen *[Cache]. The [Cache.IsZero] method
reveals whether the instance has been initialized or not. [Cache.Free]
destroys the *[Cache] -- regardless of freeze state.

[RADUA]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radua

[RASCHEMA I-Ds]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema
*/
type Cache struct {
	threshold     [2]int
	lock          *sync.Mutex
	frozen        bool
	registrations map[string]cachedRegistration
	registrants   map[string]cachedRegistrant
}

/*
NewCache returns a freshly initialized instance of *[Cache].

The registrationMax and registrantMax integer input values define the
maximum number of entries that will be cached respectively. Specifying
0 disables the respective threshold.

Attempts to exceed this threshold will silently disregard submissions
for NEW (uncached) instances, however previously cached instances will
still be refreshed.
*/
func NewCache(registrationMax, registrantMax int) *Cache {
	return newCache(registrationMax, registrantMax)
}

func newCache(registrationMax, registrantMax int) *Cache {
	var max [2]int
	if registrationMax < 0 {
		registrationMax = 0
	}
	if registrantMax < 0 {
		registrantMax = 0
	}
	max = [2]int{registrationMax, registrantMax}

	return &Cache{
		threshold:     max,
		lock:          &sync.Mutex{},
		registrations: make(map[string]cachedRegistration, max[0]),
		registrants:   make(map[string]cachedRegistrant, max[1]),
	}
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *Cache) IsZero() bool {
	return r == nil
}

/*
RegistrationLen returns the integer length of the receiver instance, thereby
revealing how many *[Registration] instances are cached. This does not take
expiration status into account.
*/
func (r *Cache) RegistrationLen() int {
	return len(r.registrations)
}

/*
RegistrationCap returns the maximum permitted number of *[Registration]
instances that may be cached. A value of zero (0) indicates no limits are
imposed upon caching requests of this form.
*/
func (r *Cache) RegistrationCap() int {
	return r.threshold[0]
}

/*
RegistrantLen returns the integer length of the receiver instance, thereby
revealing how many *[Registrant] instances are cached. This does
not take expiration status into account.
*/
func (r *Cache) RegistrantLen() int {
	return len(r.registrants)
}

/*
RegistrantCap returns the maximum permitted number of *[Registrant]
instances that may be cached. A value of zero (0) indicates no limits are
imposed upon caching requests of this form.
*/
func (r *Cache) RegistrantCap() int {
	return r.threshold[1]
}

/*
Registration returns an instance of *[Registration] following a search for
the input dn value within the receiver instance.

A nil return value can indicate any of the following:

  - Instance had expired and has since been purged, or has not yet been cached
  - Instance was found but was nil, indicating caching is disabled for the entry

Case is not significant in the matching process.
*/
func (r *Cache) Registration(dn string) *Registration {
	reg, _ := r.get(dn, 0).(*Registration)
	return reg
}

/*
Registrant returns a *[Registrant] instance following a search
for the input dn value within the receiver instance.

A nil return value can indicate any of the following:

  - Instance had expired and has since been purged, or has not yet been cached
  - Instance was found but was nil, indicating caching is disabled for the entry

Case is not significant in the matching process.
*/
func (r *Cache) Registrant(dn string) *Registrant {
	reg, _ := r.get(dn, 1).(*Registrant)
	return reg
}

func (r *Cache) get(dn string, t int) any {
	if r.IsZero() {
		return nil
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	key := lc(dn)

	if t == 0 {
		item, ok := r.registrations[key]
		if ok && item.Value != nil {
			if now().After(item.Expiry) {
				r.deleteRegistration(key)
				return nil
			}
		}

		return item.Value
	}

	item, ok := r.registrants[key]
	if ok && item.Value != nil {
		if now().After(item.Expiry) {
			r.deleteRegistrant(key)
			return nil
		}
	}

	return item.Value
}

/*
RegistrationKeys returns slices of cached element DNs, each representing
a *[Registration] instance present within the receiver.

Expiration status is not taken into account, nor are any expired elements
purged from the receiver as a result of using this method.
*/
func (r *Cache) RegistrationKeys() []string {
	return r.keys(0)
}

/*
RegistrantKeys returns slices of cached element DNs, each representing
a *[Registrant] instance present within the receiver.

Expiration status is not taken into account, nor are any expired elements
purged from the receiver as a result of using this method.
*/
func (r *Cache) RegistrantKeys() []string {
	return r.keys(1)
}

func (r *Cache) keys(t int) (keys []string) {
	if !r.IsZero() {

		appender := func(dn string) {
			if len(dn) > 0 {
				keys = append(keys, dn)
			}
		}

		r.lock.Lock()
		defer r.lock.Unlock()

		if t == 0 {
			for _, v := range r.registrations {
				appender(v.Value.DN())
			}
			return
		}

		for _, v := range r.registrants {
			appender(v.Value.DN())
		}
	}

	return
}

/*
TouchRegistration will refresh the targeted *[Registration] instance by
the input dn value, and replace its Expiry struct field with a fresh
[time.Time] instance based upon the input minutes value, which may be a
string or an int.

In addition to preserving instances past their original expiration time,
this method may be used to "resurrect" instances that have since expired
but have not yet been purged from the receiver instance.

This method has no effect if the targeted instance is not found, the receiver
is zero, or the minutes value is <= 0.
*/
func (r *Cache) TouchRegistration(dn string, minutes any) {
	r.touch(dn, minutes, 0)
}

/*
TouchRegistrant will refresh the targeted *[Registrant] instance
by the input dn value, and replace its Expiry struct field with a fresh
[time.Time] instance based upon the input minutes value, which may be a
string or an int.

In addition to preserving instances past their original expiration time,
this method may be used to "resurrect" instances that have since expired
but have not yet been purged from the receiver instance.

This method has no effect if the targeted instance is not found, the receiver
is zero, or the minutes value is <= 0.
*/
func (r *Cache) TouchRegistrant(dn string, minutes any) {
	r.touch(dn, minutes, 1)
}

func (r *Cache) touch(dn string, minutes any, t int) {
	if r.IsZero() || r.Frozen() {
		return
	}

	min := assertTTL(minutes)

	if len(dn) == 0 || min <= 0 {
		return
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	key := lc(dn)

	if t == 0 {
		if item, _ := r.registrations[key]; item.Value != nil {
			item.Expiry = newExpiry(min)
		}
		return
	}

	if item, _ := r.registrants[key]; item.Value != nil {
		item.Expiry = newExpiry(min)
	}
}

/*
Add assigns the input *[Registration] or *[Registrant] instance
to the receiver instance. The minutes input value (which may be a string
or an int) should correspond to one of the following states:

  - <=0 (entry default) indicates no caching for the indicated instance (always call LDAP)
  - All other positive values indicate a cached lifespan in minutes (cache and don't call LDAP for this entry until N minutes)

If the target instance is already cached, it shall be replaced with the
input instance, and will be subject to the new lifespan value. This will
achieve the same outcome as use of [Cache.Touch].

This method is meant for use either of the following scenarios:

  - Automatically, whereby an 'rATTL' or 'c-rATTL' value has been set within the RA DIT and is being observed following retrieval one or more LDAP entries to be marshaled
  - Manually, whereby an instance crafted by the user is being deliberately cached, whether or not LDAP is presently involved

Input instances may be cached at any point, whether modified or not.
*/
func (r *Cache) Add(instance, minutes any) {
	if r.IsZero() || r.Frozen() {
		return
	}

	min := assertTTL(minutes)

	if instance == nil {
		return
	}

	switch tv := instance.(type) {
	case *Registration:
		if len(tv.DN()) > 0 {
			r.cacheRegistration(tv, min)
		}
	case *Registrant:
		if len(tv.DN()) > 0 {
			r.cacheRegistrant(tv, min)
		}
	}
}

/*
RemoveRegistration deletes the specified *[Registration] instances from
the receiver instance.

Case is not significant in the matching process.
*/
func (r *Cache) RemoveRegistration(dn ...string) {
	r.remove(1, dn...)
}

/*
RemoveRegistrant deletes the specified *[Registrant] instances
from the receiver instance.

Case is not significant in the matching process.
*/
func (r *Cache) RemoveRegistrant(dn ...string) {
	r.remove(1, dn...)
}

func (r *Cache) remove(t int, keys ...string) {
	if r.IsZero() {
		return
	} else if len(keys) == 0 {
		return
	} else if r.Frozen() {
		return
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if t == 0 {
		r.deleteRegistration(keys...)
	} else {
		r.deleteRegistrant(keys...)
	}
}

func (r *Cache) registrationsFull() bool {
	return r.threshold[0] <= len(r.registrations) && r.threshold[0] != 0
}

func (r *Cache) cacheRegistration(reg *Registration, minutes any) {

	if r.Registration(reg.DN()) == nil {
		// reg is not presently cached.
		if r.registrationsFull() {
			// cannot cache: full house
			return
		}
	}

	ttl := assertTTL(minutes)
	if ttl <= 0 {
		return
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	key := lc(reg.DN())
	r.registrations[key] = cachedRegistration{
		Value:  reg,
		Expiry: newExpiry(ttl),
	}
}

func (r *Cache) registrantsFull() bool {
	return r.threshold[1] <= len(r.registrants) && r.threshold[1] != 0
}

func (r *Cache) cacheRegistrant(reg *Registrant, minutes any) {

	if r.Registrant(reg.DN()) == nil {
		// reg is not presently cached.
		if r.registrantsFull() {
			// cannot cache: full house
			return
		}
	}

	ttl := assertTTL(minutes)
	if ttl <= 0 {
		return
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	key := lc(reg.DN())
	r.registrants[key] = cachedRegistrant{
		Value:  reg,
		Expiry: newExpiry(ttl),
	}
}

/*
Freeze freezes the receiver instance, thereby preventing any subsequent
updates and clean-ups from proceeding. This means no expirations (removals)
of expired cache entries will occur, nor can any new elements be added to
the instance. Cached elements may still be accessed using conventional means.

See the [Cache.Thaw] method for a means of unfreezing the receiver instance.
See the [Cache.Frozen] method for a means of confirming a frozen or thawed
state.

Note that the receiver can still be freed (destroyed) by [Cache.Free]
while frozen.
*/
func (r *Cache) Freeze() {
	if !r.IsZero() && !r.Frozen() {
		r.lock.Lock()
		defer r.lock.Unlock()

		r.frozen = true
	}
}

/*
Thaw unfreezes the receiver instance, thereby allowing any subsequent
updates and clean-ups to proceed. This means that expirations (removals)
of expired cache entries will occur, and new elements may be added to
the instance.

See the [Cache.Freeze] method for a means of freezing the receiver instance.
See the [Cache.Frozen] method for a means of confirming a frozen state.
*/
func (r *Cache) Thaw() {
	if !r.IsZero() && r.Frozen() {
		r.lock.Lock()
		defer r.lock.Unlock()

		r.frozen = false
	}
}

/*
Frozen returns a Boolean value indicative of a frozen receiver state.
During a freeze state, no new elements may be added to the instance,
nor can expired elements be purged.

Note that the receiver can still be freed (destroyed) by [Cache.Free]
while frozen.
*/
func (r *Cache) Frozen() bool {
	return r.frozen
}

/*
Free frees (destroys) the *[Cache] instance, rendering it nil and unusable.

Note that this method is immutable, and will not honor any frozen state or
mutex lock.
*/
func (r *Cache) Free() {
	*r =  Cache{}
}

/*
Flush will purge all cached *[Registration] or *[Registrant] instances
from the receiver, regardless of expiration status. Following completion,
the receiver remains initialized and usable.
*/
func (r *Cache) Flush() {
	r.flush(true)
}

/*
Tidy scans for, and purges the receiver instance of, any *[Registration]
or *[Registrant] instance which have expired.
*/
func (r *Cache) Tidy() {
	r.flush(false)
}

func (r *Cache) flush(all bool) {
	if r.IsZero() || r.Frozen() {
		return
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	for k, v := range r.registrations {
		if now().After(v.Expiry) || all {
			r.deleteRegistration(k)
		}
	}

	for k, v := range r.registrants {
		if now().After(v.Expiry) || all {
			r.deleteRegistrant(k)
		}
	}
}

/*
newExpiry returns an instance of [time.Time] that defines when a given
item will be considered expired and needing refresh.
*/
func newExpiry(min int) time.Time {
	return now().Add(time.Duration(min) * time.Minute)
}

/*
WriteRegistrations returns an error following an attempt to write the
current contents of the *[Registration] cache to the filename indicated.
*/
func (r *Cache) WriteRegistrations(filename string) error {
	return r.write(filename, 0)
}

/*
WriteRegistrants returns an error following an attempt to write the current
contents of the *[Registrant] cache to the filename indicated.
*/
func (r *Cache) WriteRegistrants(filename string) error {
	return r.write(filename, 1)
}

func (r *Cache) write(filename string, t int) error {
	if r.IsZero() {
		return nil // nothing to write!
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	codec := gob.NewEncoder(file)

	if t == 0 {
		err = codec.Encode(r.registrations)
	} else {
		err = codec.Encode(r.registrants)
	}

	return err
}

/*
LoadRegistrations returns an error following an attempt to read the filename
indicated into the receiver's *[Registration] cache.
*/
func (r *Cache) LoadRegistrations(filename string) error {
	return r.load(filename, 0)
}

/*
LoadRegistrations returns an error following an attempt to read the filename
indicated into the receiver's *[Registrant] cache.

If the receiver instance is zero (uninitialized), a [NilCacheErr] is returned.
If the receiver is frozen, a [FrozenCacheErr] is returned.
*/
func (r *Cache) LoadRegistrants(filename string) error {
	return r.load(filename, 1)
}

func (r *Cache) load(filename string, t int) error {
	if r.IsZero() {
		return NilCacheErr
	} else if r.Frozen() {
		return FrozenCacheErr
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	codec := gob.NewDecoder(file)

	if t == 0 {
		err = codec.Decode(&r.registrations)
	} else {
		err = codec.Decode(&r.registrants)
	}

	return err
}

func (r *Cache) deleteRegistration(dn ...string) {
	if !r.Frozen() {
		for _, x := range dn {
			delete(r.registrations, lc(x))
		}
	}
}

func (r *Cache) deleteRegistrant(dn ...string) {
	if !r.Frozen() {
		for _, x := range dn {
			delete(r.registrants, lc(x))
		}
	}
}
