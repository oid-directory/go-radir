package radir

/*
se.go implements a generic subentry type with basic methods.
*/

/*
Subentry contains contents derived from, or for use in creating, LDAP
subentries within the context of the I-D series.

In cases where Collective Attributes, or other virtualization services,
are not available, and where manual population of values meant for broad
assignment to many entries is simply not practical, this type provides
for a crude alternative.

This type incorporates certain struct types, such as *[Spatial],
*[FirstAuthority], *[CurrentAuthority] and *[Sponsor], which are used
in cases where many entries share common values.

The drawback is that the RA DUA must perform a separate LDAP search
request for subentries, and must understand how to utilize the values
found as "fallbacks" for missing (explicit) values normally assigned
to entries directly.

In cases where Collective Attributes are supported, this type can provide
for convenient creation of subentries for submission through an LDAP Add
request or other means.

Instances of this type are created using the *[DITProfile.NewSubentry]
method.
*/
type Subentry struct {
	R_DN   string   `ldap:"dn"`
	R_TTL  string   `ldap:"rATTL"`
	RC_TTL string   `ldap:"c-rATTL;collective"`
	R_GSR  string   `ldap:"governingStructureRule"`
	R_OC   []string `ldap:"objectClass"`
	R_STS  []string `ldap:"subtreeSpecification"`

	R_Spatial    *Spatial
	R_X660       *X660
	R_Extra      *Supplement
	r_DITProfile *DITProfile
	r_root       *registeredRoot
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *Subentry) IsZero() bool {
	return r == nil
}

func (r *Subentry) isEmpty() bool {
	return structEmpty(r)
}

/*
Structural returns the STRUCTURAL objectClass `subentry`.
*/
func (r *Subentry) Structural() (s string) {
	return `subentry`
}

/*
SubtreeSpecification returns the underlying Subtree Specification slice
values, if set, else zero slices are returned.
*/
func (r *Subentry) SubtreeSpecification() (sts []string) {
	if !r.IsZero() {
		sts = r.R_STS
	}

	return
}

/*
SetSubtreeSpecification appends the provided string value to the receiver
instance as a Subtree Specification. If an instance of []string is provided,
the receiver value is clobbered (overwritten).
*/
func (r *Subentry) SetSubtreeSpecification(args ...any) (err error) {
	return writeFieldByTag(`subtreeSpecification`, r.SetSubtreeSpecification, r, args...)
}

/*
SubtreeSpecificationGetFunc executes the [GetOrSetFunc] instance and
returns its own return values. The 'any' value will require type assertion
in order to be accessed reliably. An error is returned if issues arise.
*/
func (r *Subentry) SubtreeSpecificationGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `subtreeSpecification`)
}

/*
LDIF returns the string LDIF form of the receiver instance. Note that
this is a crude approximation of LDIF and should ideally be parsed
through a reliable LDIF parser such as [go-ldap/ldif] to verify integrity.

[go-ldap/ldif]: https://pkg.go.dev/github.com/go-ldap/ldif
*/
func (r *Subentry) LDIF() (l string) {
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
			bld.WriteString(r.Spatial().ldif())

			l = bld.String()
		}
	}

	return
}

func (r *Subentry) refreshObjectClasses() {
	bools := []bool{
		r.X660().isEmpty(),
		r.Spatial().isEmpty(),
		r.Supplement().isEmpty(),
	}

	for tx, tag := range []string{
		`x660Context`,
		`spatialContext`,
		`registrationSupplement`,
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
DN returns the string-based LDAP Distinguished Name value, or a zero
string if unset.
*/
func (r *Subentry) DN() (dn string) {
	if !r.IsZero() {
		dn = r.R_DN
	}

	return
}

/*
DNGetFunc executes the [GetOrSetFunc] instance and returns its own
return values. The 'any' value will require type assertion in order
to be accessed reliably. An error is returned if issues arise.
*/
func (r *Subentry) DNGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `dn`)
}

/*
SetDN assigns the first provided value to the receiver instance as a string DN.

If additional arguments are present, the second value will be asserted as an
instance of [GetOrSetFunc], returning an error if this fails.
*/
func (r *Subentry) SetDN(args ...any) error {
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
func (r *Subentry) GoverningStructureRule() (gsr string) {
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
func (r *Subentry) GoverningStructureRuleGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `governingStructureRule`)
}

/*
TTL returns the raw time-to-live within the receiver instance.
*/
func (r *Subentry) TTL() (ttl string) {
	if !r.IsZero() {
		ttl = r.R_TTL
	}

	return
}

/*
TTL returns the raw (collective) time-to-live within the receiver instance.
*/
func (r *Subentry) CTTL() (cttl string) {
	if !r.IsZero() {
		cttl = r.RC_TTL
	}

	return
}

/*
SetCTTL assigns the provided string TTL (collective) value to the receiver instance.
*/
func (r *Subentry) SetCTTL(args ...any) error {
	return writeFieldByTag(`c-rATTL;collective`, r.SetCTTL, r, args...)
}

/*
SetTTL assigns the provided string TTL value to the receiver instance.
*/
func (r *Subentry) SetTTL(args ...any) error {
	return writeFieldByTag(`rATTL`, r.SetTTL, r, args...)
}

/*
TTLGetFunc processes the underlying TTL field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Subentry) TTLGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `rATTL`)
}

/*
CTTLGetFunc processes the underlying (collective) TTL field value through
the provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Subentry) CTTLGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `c-rATTL`)
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
func (r *Subentry) Marshal(meth func(any) error) (err error) {
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
func (r *Subentry) Unmarshal() (outer map[string][]string) {
	if r.IsZero() {
		return
	}

	outer = make(map[string][]string)

	for _, inner := range []map[string][]string{
		unmarshalStruct(r, outer),
		r.X660().unmarshal(),
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
ObjectClasses returns string slices of "[objectClass]" in descriptor or
numeric OID forms currently held by the receiver instance.

[objectClass]: https://www.rfc-editor.org/rfc/rfc4512.html#section-3.3
*/
func (r *Subentry) ObjectClasses() []string {
	return r.R_OC
}

/*
SetObjectClasses appends one or more string descriptor or numeric OID
"[objectClass]" values to the receiver instance. Note that if a slice is
passed as X, the destination value will be clobbered.

[objectClass]: https://www.rfc-editor.org/rfc/rfc4512.html#section-3.3
*/
func (r *Subentry) SetObjectClasses(args ...any) error {
	return writeFieldByTag(`objectClass`, r.SetObjectClasses, r, args...)
}

/*
ObjectClassesGetFunc processes the underlying field value(s) through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Subentry) ObjectClassesGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `objectClass`)
}

/*
Kind returns the static string value `STRUCTURAL` merely as a convenient
means to determine what [kind of "objectClass"] the receiver describes.

[kind of "objectClass"]: https://www.rfc-editor.org/rfc/rfc4512.html#section-4.1.1
*/
func (r *Subentry) Kind() string {
	return `STRUCTURAL`
}

/*
DITProfile returns the *[DITProfile] instance assigned to the receiver,
if set, else nil is returned.
*/
func (r *Subentry) DITProfile() (prof *DITProfile) {
	if !r.IsZero() {
		if prof = r.r_DITProfile; !prof.Valid() {
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
func (r *Subentry) SetDITProfile(d *DITProfile) {
	if d.Valid() {
		if idx := d.RegistrationSuffixEqual(r.DN()); idx != -1 {
			r.r_DITProfile = d
		}
	}
}
