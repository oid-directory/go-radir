package radir

/*
reg.go contains Registration methods.
*/

/*
Registration contains information either to be set upon, or derived from,
an LDAP entry that describes a registration.
*/
type Registration struct {
	R_DN      string   `ldap:"dn"`
	R_GSR     string   `ldap:"governingStructureRule"`
	R_TTL     string   `ldap:"rATTL"`
	RC_TTL    string   `ldap:"c-rATTL;collective"`
	R_OC      []string `ldap:"objectClass"`
	R_Desc    []string `ldap:"description"` // effective "title" of reg
	R_Also    []string `ldap:"seeAlso"`
	R_LongArc []string `ldap:"longArc"` // only permitted for subArcs of Joint-ISO-ITU-T (2).

	R_X660    *X660       // ITU-T Rec. X.660 types
	R_X667    *X667       // ITU-T Rec. X.667 types
	R_X680    *X680       // ITU-T Rec. X.680 types
	R_X690    *X690       // ITU-T Rec. X.690 types
	R_Extra   *Supplement // Non-standard: Supplemental types
	R_Spatial *Spatial    // Non-standard: Spatial types

	R_DITProfile *DITProfile

	r_root *registeredRoot
}

/*
registeredRoot contains information about the nature and placement of
this registration. It is populated through subsequent X.680 input, and
will (likely) be queried at any point by other constructs, such as X.660.
*/
type registeredRoot struct {
	Depth      int    // 0 = root, >=1 sub arc (default: -1)
	N          int    // 0, 1 or 2 (default: -1)
	Id         string // identifier (name) of root
	NaNF       string // nameAndNumberForm ("Id(N)") of root
	Structural string // rootArc or arc
	Auxiliary  string // iTUTRegistration, iSORegistration or jointISOITUTRegistration
}

/*
Registrations contains slices of [Registration] instances.
*/
type Registrations []*Registration

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *Registration) IsZero() bool {
	return r == nil
}

func (r *Registration) isEmpty() bool {
	return structEmpty(r)
}

/*
DN returns the string-based LDAP Distinguished Name value, or a zero
string if unset.
*/
func (r *Registration) DN() string {
	return r.R_DN
}

/*
DNGetFunc executes the [GetOrSetFunc] instance and returns its own
return values. The 'any' value will require type assertion in order
to be accessed reliably. An error is returned if issues arise.
*/
func (r *Registration) DNGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `dn`)
}

/*
SetDN assigns the first provided value to the receiver instance as a string DN.

If additional arguments are present, the second value will be asserted as an
instance of [GetOrSetFunc], returning an error if this fails.
*/
func (r *Registration) SetDN(args ...any) error {
	return writeFieldByTag(`dn`, r.SetDN, r, args...)
}

/*
GoverningStructureRule returns the "[governingStructureRule]" assigned to
the receiver instance.

Note the "[governingStructureRule]" type is operational, and cannot be set
by the end user. It is also not collective.

Use of this method is only meaningful in environments which employ one or
more "[dITStructureRule]" definitions.

[governingStructureRule]: https://datatracker.ietf.org/doc/html/rfc4512#section-3.4.6
[dITStructureRule]: https://datatracker.ietf.org/doc/html/rfc4512#section-4.1.7.1
*/
func (r *Registration) GoverningStructureRule() (gsr string) {
	if !r.IsZero() {
		gsr = r.R_GSR
	}

	return
}

/*
GoverningStructureRuleGetFunc executes the [GetOrSetFunc] instance and
returns its own return values. The 'any' value will require type assertion
in order to be accessed reliably. An error is returned if issues arise.
*/
func (r *Registration) GoverningStructureRuleGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `governingStructureRule`)
}

/*
Description returns one or more string description values assigned to the
receiver instance.
*/
func (r *Registration) Description() (desc []string) {
	if !r.IsZero() {
		desc = r.R_Desc
	}

	return
}

/*
SetDescription appends one or more string description values to the receiver
instance. Note that if a slice is passed as X, the destination value will be
clobbered.
*/
func (r *Registration) SetDescription(args ...any) error {
	return writeFieldByTag(`description`, r.SetDescription, r, args...)
}

/*
DescriptionGetFunc processes the underlying field value(s) through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Registration) DescriptionGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `description`)
}

/*
SeeAlso returns the string DN values assigned to the receiver instance.
*/
func (r *Registration) SeeAlso() (also []string) {
	if !r.IsZero() {
		also = r.R_Also
	}

	return
}

/*
SetSeeAlso appends one or more string DN values to the receiver instance.
Note that if a slice is passed as X, the destination value will be clobbered.
*/
func (r *Registration) SetSeeAlso(args ...any) error {
	return writeFieldByTag(`seeAlso`, r.SetSeeAlso, r, args...)
}

/*
SeeAlsoGetFunc processes the underlying string DN field value(s) through
the provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Registration) SeeAlsoGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `seeAlso`)
}

/*
TTL returns the effective time-to-live for the receiver instance, taking
into account *[DITProfile]-inherited values as well as any subtree-based
(COLLECTIVE) and entry literal values. The output can be used to instruct
instances of [Cache] when, and when not, to cache an instance.

See [Section 2.2.3.4 of the RADUA I-D] for details related to TTL precedence.

[Section 2.2.3.4 of the RADUA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radua#section-2.2.3.4
*/
func (r *Registration) TTL() string {
	ct := r.DITProfile().TTL()
	lt := selectTTL(r.R_TTL, r.RC_TTL)

	if lt == `` {
		// If no localized TTL or COLLECTIVE TTL, then
		// just return the DITProfile-based TTL.
		return ct
	}

	return lt
}

/*
SetTTL assigns the provided string TTL value to the receiver instance.
*/
func (r *Registration) SetTTL(args ...any) error {
	return writeFieldByTag(`rATTL`, r.SetTTL, r, args...)
}

/*
TTLGetFunc processes the underlying TTL field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Registration) TTLGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `rATTL`)
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
func (r *Registration) Marshal(meth func(any) error) (err error) {
	if meth == nil {
		err = NilMethodErr
		return
	} else if r == nil {
		err = NilRegistrationErr
		return
	}

	for _, err = range []error{
		meth(r),
		r.X660().marshal(meth),
		r.X667().marshal(meth),
		r.X680().marshal(meth),
		r.X690().marshal(meth),
		r.Supplement().marshal(meth),
		r.Spatial().marshal(meth),
	} {
		if err != nil {
			return
		}
	}

	// clean-up any duplicate objectClass
	// slices that may have appeared, and
	// remove any classes that aren't used.
	r.refreshObjectClasses()

	return
}

/*
Unmarshal transports values from the receiver into an instance of
map[string][]string, which can subsequently be fed to the [go-ldap/v3
NewEntry] function.

[go-ldap/v3 NewEntry]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#NewEntry
*/
func (r *Registration) Unmarshal() (outer map[string][]string) {
	if r.IsZero() {
		return
	}

	outer = make(map[string][]string)

	for _, inner := range []map[string][]string{
		unmarshalStruct(r, outer),
		r.X660().unmarshal(),
		r.X667().unmarshal(),
		r.X680().unmarshal(),
		r.X690().unmarshal(),
		r.Supplement().unmarshal(),
		r.Spatial().unmarshal(),
	} {
		if inner != nil {
			for k, v := range inner {
				outer[k] = v
			}
		}
	}

	return
}

/*
Unmarshal is a convenience method that returns slices of map[string][]string
instances, each representative of an individual *[Registration] instance
that was deemed valid for unmarshaling to map[string][]string.
*/
func (r *Registrations) Unmarshal() (maps []map[string][]string) {
	maps = make([]map[string][]string, 0)
	for i := 0; i < len(*r); i++ {
		if um := (*r)[i].Unmarshal(); len(um) > 0 {
			maps = append(maps, um)
		}
	}
	return
}

/*
Marshal returns an error following an attempt to execute all input method
[go-ldap/v3 Entry.Unmarshal] signatures. The result is a sequence of marshaled
*[Registration] instances being added to the receiver instance.

The input *[DITProfile] value will be used to initialize each *[Registration]
instance using the appropriate configuration.

[go-ldap/v3 Entry.Unmarshal]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry.Unmarshal
*/
func (r *Registrations) Marshal(profile *DITProfile, meths ...func(any) error) (err error) {
	if !profile.Valid() {
		err = DUAConfigValidityErr
		return
	}

	for i := 0; i < len(meths) && err == nil; i++ {
		reg := profile.NewRegistration()
		if err = reg.Marshal(meths[i]); err == nil {
			*r = append(*r, reg)
		}
	}

	return
}

/*
ObjectClasses returns string slices of "[objectClass]" in descriptor or
numeric OID forms currently held by the receiver instance.

[objectClass]: https://www.rfc-editor.org/rfc/rfc4512.html#section-3.3
*/
func (r *Registration) ObjectClasses() []string {
	return r.R_OC
}

/*
SetObjectClasses appends one or more string descriptor or numeric OID
"[objectClass]" values to the receiver instance. Note that if a slice is
passed as X, the destination value will be clobbered.

[objectClass]: https://www.rfc-editor.org/rfc/rfc4512.html#section-3.3
*/
func (r *Registration) SetObjectClasses(args ...any) error {
	return writeFieldByTag(`objectClass`, r.SetObjectClasses, r, args...)
}

/*
ObjectClassesGetFunc processes the underlying field value(s) through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Registration) ObjectClassesGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `objectClass`)
}

/*
Kind returns the static string value `STRUCTURAL` merely as a convenient
means to determine what [kind of "objectClass"] the receiver describes.

[kind of "objectClass"]: https://www.rfc-editor.org/rfc/rfc4512.html#section-4.1.1
*/
func (r *Registration) Kind() string {
	return `STRUCTURAL`
}

/*
Structural returns the STRUCTURAL objectClass of the receiver instance. The
expected return values are `arc` or `rootArc`.  A zero string is returned
if neither class is declared.
*/
func (r *Registration) Structural() (s string) {
	if !r.IsZero() {
		if strInSlice(`rootArc`, r.R_OC) {
			s = `rootArc`
		} else if strInSlice(`arc`, r.R_OC) {
			s = `arc`
		}
	}

	return
}

/*
DITProfile returns the *[DITProfile] instance assigned to the receiver,
if set, else nil is returned.
*/
func (r *Registration) DITProfile() (prof *DITProfile) {
	if !r.IsZero() {
		if prof = r.R_DITProfile; !prof.Valid() {
			prof = &DITProfile{}
		}
	}
	return
}

/*
SetDITProfile assigns *[DITProfile] d to the receiver instance if, and
only if, both of the following evaluate as true:

  - [DITProfile.Valid] returns true for d, AND ...
  - the DN of the receiver instance matches a *[Registration] search base within d

Case is not significant in the matching process.
*/
func (r *Registration) SetDITProfile(d *DITProfile) {
	if d.Valid() {
		if idx := d.RegistrationSuffixEqual(r.DN()); idx != -1 {
			r.R_DITProfile = d
		}
	}
}

/*
Combined returns a Boolean value indicative of the "Combined Registrants
Policy" being in-force.

If the underlying instance of *[DITProfile] within the receiver instance
is zero, false is returned.
*/
func (r *Registration) Combined() bool {
	if !r.IsZero() {
		if dp := r.DITProfile(); !dp.IsZero() {
			return dp.Combined()
		}
	}

	return false
}

/*
Dedicated returns a Boolean value indicative of the "Dedicated Registrants
Policy" being in-force.

If the underlying instance of *[DITProfile] within the receiver instance
is zero, false is returned.
*/
func (r *Registration) Dedicated() bool {
	if !r.IsZero() {
		if dp := r.DITProfile(); !dp.IsZero() {
			return dp.Dedicated()
		}
	}

	return false
}

func (r *Registration) refreshObjectClasses() {
	bools := []bool{
		r.X660().isEmpty(),
		r.X667().isEmpty(),
		r.X680().isEmpty(),
		r.X690().isEmpty(),
		r.Supplement().isEmpty(),
		r.Spatial().isEmpty(),
	}

	for tx, tag := range []string{
		`x660Context`,
		`x667Context`,
		`x680Context`,
		`x690Context`,
		`registrationSupplement`,
		`spatialContext`,
	} {
		if bools[tx] {
			r.R_OC = removeStrInSlice(tag, r.R_OC)
		} else {
			if !strInSlice(tag, r.R_OC) {
				r.R_OC = append(r.R_OC, tag)
			}
		}
	}

	// go-ldap/v3.Entry.Unmarshal is sloppy about adding
	// duplicate objectClasses, so let's clean up any
	// that may have appeared.
	var tmp []string
	for _, oc := range r.R_OC {
		if !strInSlice(oc, tmp) {
			tmp = append(tmp, oc)
		}
	}
	if len(tmp) != len(r.R_OC) {
		r.R_OC = tmp
	}
}

/*
Root returns the official integer value for any of the three (3) possible
roots alongside an AUXILIARY class descriptor. The possible integer roots
are as follows:

  - itu-t (0)
  - iso (1)
  - joint-iso-itut (2)

If the receiver instance is determined to be associated with a root
*[Registration] instance, no descriptor class is returned.

If the receiver is not yet populated in this regard, a -1 integer is
returned with no descriptor class.

Otherwise, for any non-root *[Registration], one (1) of the following
AUXILIARY class descriptors will be returned:

  - "[iTUTRegistration]"
  - "[iSORegistration]"
  - "[jointISOITUTRegistration]"

[iTUTRegistration]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.8
[iSORegistration]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.9
[jointISOITUTRegistration]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.10
*/
func (r *Registration) Root() (n int, class string) {
	n = -1
	if !r.IsZero() {
		if !structEmpty(r.r_root) {
			n = r.r_root.N
			if r.r_root.Depth > 1 {
				class = r.r_root.Auxiliary
			}
		}
	}

	return
}

/*
LDIF returns the string LDIF form of the receiver instance. Note that this
is a crude approximation of LDIF and should ideally be parsed through a
reliable LDIF parser such as [go-ldap/ldif] to verify integrity.

[go-ldap/ldif]: https://pkg.go.dev/github.com/go-ldap/ldif
*/
func (r *Registration) LDIF() (l string) {
	if !r.IsZero() {
		dn := readFieldByTag(`dn`, r)
		if len(dn) > 0 {
			r.refreshObjectClasses()

			oc := readFieldByTag(`objectClass`, r)
			bld := newBuilder()

			bld.WriteString(`dn: ` + dn[0])
			bld.WriteRune(10)

			for i := 0; i < len(oc); i++ {
				bld.WriteString(`objectClass: ` + oc[i])
				bld.WriteRune(10)
			}

			bld.WriteString(toLDIF(r))
			bld.WriteString(r.X660().ldif())
			bld.WriteString(r.X667().ldif())
			bld.WriteString(r.X680().ldif())
			bld.WriteString(r.X690().ldif())
			bld.WriteString(r.Spatial().ldif())

			l = bld.String()
		}
	}

	return
}
