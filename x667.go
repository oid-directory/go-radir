package radir

/*
X667 implements [RASCHEMA § 2.5.5] and derives various concepts from
[ITU-T Rec. X.667].

	( 1.3.6.1.4.1.56521.101.2.5.5
	    NAME 'x667Context'
	    DESC 'X.667 contextual class'
	    SUP registration AUXILIARY
	    MUST registeredUUID )

Instances of this type need not be initialized by the user directly.

[RASCHEMA § 2.5.5]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.5
[ITU-T Rec. X.667]: https://www.itu.int/rec/T-REC-X.667
*/
type X667 struct {
	R_UUID string `ldap:"registeredUUID"` // RASCHEMA § 2.3.102

	r_DITProfile *DITProfile
	r_root       *registeredRoot
}

/*
X667 returns (and if needed, initializes) the embedded instance of *[X667].
*/
func (r *Registration) X667() *X667 {
	if r.IsZero() {
		return &X667{}
	}

	if r.R_X667.IsZero() {
		r.R_X667 = new(X667)
		r.R_X667.r_DITProfile = r.Profile()
		r.R_X667.r_root = r.r_root
	}

	return r.R_X667
}

/*
Profile returns the *[DITProfile] instance assigned to the receiver,
if set, else nil is returned.
*/
func (r *X667) profile() (prof *DITProfile) {
	if prof = r.r_DITProfile; !prof.Valid() {
		prof = &DITProfile{}
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *X667) IsZero() bool {
	return r == nil
}

func (r *X667) isEmpty() bool {
	return structEmpty(r)
}

func (r *X667) ldif() (l string) {
	if !r.IsZero() {
		l = toLDIF(r)
	}

	return
}

/*
Unmarshal returns an instance of map[string][]string bearing the contents
of the receiver.
*/
func (r *X667) unmarshal() map[string][]string {
	m := make(map[string][]string)
	return unmarshalStruct(r, m)
}

/*
Marshal returns an error following an attempt to execute the input meth
"func(any) error" method signature.

The meth value should be the (unexecuted) [go-ldap/v3 Entry.Unmarshal]
method instance for the [Entry] being transferred (marshaled) into the
receiver instance.

Alternatively, if the user has fashioned an alternative method of the
same signature, this may be supplied instead.

[go-ldap/v3 Entry.Unmarshal]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry.Unmarshal
[Entry]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry
*/
func (r *X667) marshal(meth func(any) error) (err error) {
	if !r.IsZero() {
		err = meth(r)
	}

	return
}

/*
RegisteredUUID returns the string UUID value assigned to the receiver
instance.
*/
func (r *X667) RegisteredUUID() string {
	return r.R_UUID
}

/*
SetRegisteredUUID assigns the string UUID value to the receiver instance.
*/
func (r *X667) SetRegisteredUUID(args ...any) error {
	return writeFieldByTag(`registeredUUID`, r.SetRegisteredUUID, r, args...)
}

/*
RegisteredUUIDGetFunc processes the underlying string UUIDfield value
through the provided [GetOrSetFunc] instance, returning an interface
value alongside an error.
*/
func (r *X667) RegisteredUUIDGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `registeredUUID`)
}
