package radir

/*
athy.go contains all top-level Registrant constructs.

See also ca.go, fa.go and sp.go for Current, First and Sponsor elements.
*/

/*
Registrant implements the "Dedicated Registrants" policy defined
in [Section 3.2.1.1.1 of the RADIT I-D].

In addition to "dn" and "registrantID" primitive fields, this type contains
three (3) embedded types:

  - *[FirstAuthority]
  - *[CurrentAuthority]
  - *[Sponsor]

[Section 3.2.1.1.1 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.2.1.1.1
*/
type Registrant struct {
	R_DN   string   `ldap:"dn"`
	R_GSR  string   `ldap:"governingStructureRule"`
	R_Id   string   `ldap:"registrantID"`       // RASCHEMA § 2.3.34
	R_TTL  string   `ldap:"rATTL"`              // RASCHEMA § 2.3.100
	RC_TTL string   `ldap:"c-rATTL;collective"` // RASCHEMA § 2.3.101
	R_OC   []string `ldap:"objectClass"`
	R_Desc []string `ldap:"description"`
	R_Also []string `ldap:"seeAlso"`

	R_DITProfile *DITProfile

	R_CA *CurrentAuthority
	R_FA *FirstAuthority
	R_SA *Sponsor
}

/*
DN returns the distinguished name value assigned to the receiver instance.
*/
func (r *Registrant) DN() (dn string) {
	if !r.IsZero() {
		dn = r.R_DN
	}

	return
}

/*
SetDN assigns the provided string value to the receiver instance.
*/
func (r *Registrant) SetDN(args ...any) error {
	if len(args) == 1 {
		if _, ok := args[0].(func(...any) (any, error)); ok {
			args = []any{``, args[0]}
			return writeFieldByTag(`dn`, r.SetDN, r, args...)
		}
	}

	return writeFieldByTag(`dn`, r.SetDN, r, args...)
}

/*
DNGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Registrant) DNGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `dn`)
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
func (r *Registrant) GoverningStructureRule() (gsr string) {
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
func (r *Registrant) GoverningStructureRuleGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `governingStructureRule`)
}

func (r *Registrant) refreshObjectClasses() {
	bools := []bool{
		r.FirstAuthority().isEmpty(),
		r.CurrentAuthority().isEmpty(),
		r.Sponsor().isEmpty(),
	}

	for tx, tag := range []string{
		`firstAuthorityContext`,
		`currentAuthorityContext`,
		`sponsorContext`,
	} {
		if bools[tx] {
			if strInSlice(tag, r.R_OC) {
				r.R_OC = removeStrInSlice(tag, r.R_OC)
			}
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
LDIF returns the string LDIF form of the receiver instance. Note that this
is a crude approximation of LDIF and should ideally be parsed through a
reliable LDIF parser such as [go-ldap/ldif] to verify integrity.

[go-ldap/ldif]: https://pkg.go.dev/github.com/go-ldap/ldif
*/
func (r *Registrant) LDIF() (l string) {
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
			bld.WriteString(r.FirstAuthority().ldif())
			bld.WriteString(r.CurrentAuthority().ldif())
			bld.WriteString(r.Sponsor().ldif())

			l = bld.String()
		}
	}

	return
}

/*
SeeAlso returns the string DN values assigned to the receiver instance.
*/
func (r *Registrant) SeeAlso() (also []string) {
	if !r.IsZero() {
		also = r.R_Also
	}

	return
}

/*
Dedicated returns a Boolean value indicative of the "Dedicated Registrants
Policy" being in-force.

If the underlying instance of *[DITProfile] within the receiver instance
is zero, false is returned.
*/
func (r *Registrant) Dedicated() bool {
	if !r.IsZero() {
		if dp := r.DITProfile(); !dp.IsZero() {
			return dp.Dedicated()
		}
	}

	return false
}

/*
SetSeeAlso appends one or more string DN values to the receiver instance.
Note that if a slice is passed as X, the destination value will be clobbered.
*/
func (r *Registrant) SetSeeAlso(args ...any) error {
	return writeFieldByTag(`seeAlso`, r.SetSeeAlso, r, args...)
}

/*
SeeAlsoGetFunc processes the underlying string DN field value(s) through
the provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Registrant) SeeAlsoGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `seeAlso`)
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *Registrant) IsZero() bool {
	return r == nil
}

func (r *Registrant) isEmpty() bool {
	return structEmpty(r)
}

/*
FirstAuthority returns (and if needed, initializes) the embedded instance of
*[FirstAuthority].

This method is intended solely for use under the terms of the "Dedicated
Registrants policy".
*/
func (r *Registrant) FirstAuthority() (dr *FirstAuthority) {
	if r.Dedicated() {
		if r.R_FA.IsZero() {
			r.R_FA = new(FirstAuthority)
			r.R_FA.r_alt_types = r.R_DITProfile.r_alt_types
		}

		return r.R_FA
	}

	return &FirstAuthority{}
}

/*
CurrentAuthority returns (and if needed, initializes) the embedded instance of
*[CurrentAuthority].

This method is intended solely for use under the terms of the "Dedicated
Registrants policy".
*/
func (r *Registrant) CurrentAuthority() *CurrentAuthority {
	if r.Dedicated() {
		if r.R_CA.IsZero() {
			r.R_CA = new(CurrentAuthority)
			r.R_CA.r_alt_types = r.R_DITProfile.r_alt_types
		}

		return r.R_CA
	}

	return &CurrentAuthority{}
}

/*
Sponsor returns (and if needed, initializes) the embedded instance of
*[Sponsor].

This method is intended solely for use under the terms of the "Dedicated
Registrants policy".
*/
func (r *Registrant) Sponsor() *Sponsor {
	if r.Dedicated() {
		if r.R_SA.IsZero() {
			r.R_SA = new(Sponsor)
			r.R_SA.r_alt_types = r.R_DITProfile.r_alt_types
		}

		return r.R_SA
	}

	return &Sponsor{}
}

/*
TTL returns the effective time-to-live for the receiver instance, taking
into account *[DITProfile]-inherited values as well as subtree-based
(COLLECTIVE) and entry literal values. The output can be used to instruct
instances of [Cache] when, and when not, to cache an instance.

See [Section 2.2.3.4 of the RADUA I-D] for details related to TTL precedence
when handling multiple TTL directives.

[Section 2.2.3.4 of the RADUA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radua#section-2.2.3.4
*/
func (r *Registrant) TTL() string {
	ct := r.DITProfile().TTL()
	lt := selectTTL(r.R_TTL, r.R_TTL)

	if lt == `` {
		return ct
	}

	return lt
}

/*
SetTTL assigns the provided string value to the receiver instance.
*/
func (r *Registrant) SetTTL(args ...any) error {
	return writeFieldByTag(`rATTL`, r.SetTTL, r, args...)
}

/*
TTLGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Registrant) TTLGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `rATTL`)
}

/*
Description returns one or more string description values assigned to the
receiver instance.
*/
func (r *Registrant) Description() (desc []string) {
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
func (r *Registrant) SetDescription(args ...any) error {
	return writeFieldByTag(`description`, r.SetDescription, r, args...)
}

/*
DescriptionGetFunc processes the underlying field value(s) through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Registrant) DescriptionGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `description`)
}

/*
ID returns the "[registrantID]" value assigned to the receiver instance.

[registrantID]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.34
*/
func (r *Registrant) ID() (id string) {
	if !r.IsZero() {
		id = r.R_Id
	}

	return
}

/*
SetID assigns the provided string value to the receiver instance.
*/
func (r *Registrant) SetID(args ...any) error {
	return writeFieldByTag(`registrantID`, r.SetID, r, args...)
}

/*
IDGetFunc processes the underlying field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Registrant) IDGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `registrantID`)
}

/*
ObjectClasses returns the string objectClass descriptor or numeric OID
values held by the receiver instance.
*/
func (r *Registrant) ObjectClasses() (oc []string) {
	if !r.IsZero() {
		oc = r.R_OC
	}

	return
}

/*
SetObjectClasses appends one or more string descriptor or numeric OID values
to the receiver instance. Note that if a slice is passed as X, the destination
value will be clobbered.
*/
func (r *Registrant) SetObjectClasses(args ...any) error {
	return writeFieldByTag(`objectClass`, r.SetObjectClasses, r, args...)
}

/*
ObjectClassesGetFunc processes the underlying field value(s) through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Registrant) ObjectClassesGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `objectClass`)
}

/*
Structural returns the string literal `registrant`, offering a convenient
means of identifying the STRUCTURAL objectClass of the receiver.
*/
func (r *Registrant) Structural() (oc string) {
	if !r.IsZero() {
		oc = `registrant`
	}

	return
}

/*
DITProfile returns the *[DITProfile] instance assigned to the receiver,
if set, else a freshly initialized instance is returned.
*/
func (r *Registrant) DITProfile() (prof *DITProfile) {
	if prof = r.R_DITProfile; prof == nil {
		prof = &DITProfile{}
	}

	return
}

/*
Unmarshal transports values from the receiver into an instance of
map[string][]string, which can subsequently be fed to go-ldap's
NewEntry function.
*/
func (r *Registrant) Unmarshal() (outer map[string][]string) {
	if r.IsZero() {
		return
	}

	if dc := r.DITProfile(); !dc.Dedicated() {
		// You're trying to use something that is
		// antithetical to "combined" reg pols.
		return
	}

	outer = make(map[string][]string)

	for _, inner := range []map[string][]string{
		unmarshalStruct(r, outer),
		r.FirstAuthority().unmarshal(),
		r.CurrentAuthority().unmarshal(),
		r.Sponsor().unmarshal(),
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
Marshal returns an error following an attempt to execute the input meth
signature upon the receiver instance.

See the [Registrant] documentation for details on how and why
this method was implemented, as opposed to encouraging use of the
[go-ldap/v3 Entry.Unmarshal] method alone.

[go-ldap/v3 Entry.Unmarshal]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry.Unmarshal
*/
func (r *Registrant) Marshal(meth func(any) error) (err error) {
	if meth == nil {
		err = NilMethodErr
		return
	} else if r == nil {
		err = NilRegistrationErr
		return
	}

	if dc := r.DITProfile(); !dc.Dedicated() {
		// You're trying to use something that is
		// antithetical to "combined" reg pols.
		err = RegistrantPolicyErr
		return
	}

	for _, err = range []error{
		meth(r),
		r.CurrentAuthority().marshal(meth),
		r.FirstAuthority().marshal(meth),
		r.Sponsor().marshal(meth),
	} {
		if err != nil {
			break
		}
	}

	return
}

/*
Registrants contains slices of *[Registrant] instances.
*/
type Registrants []*Registrant

/*
Unmarshal is a convenience method that returns slices of map[string][]string
instances, each representative of an individual unmarshaled *[Registrant]
instance.
*/
func (r *Registrants) Unmarshal() (maps []map[string][]string) {
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

Note that use of this method only makes sense when operating under the
terms of the "Dedicated Registrant Policy".

[go-ldap/v3 Entry.Unmarshal]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry.Unmarshal
*/
func (r *Registrants) Marshal(profile *DITProfile, meths ...func(any) error) (err error) {
	if !profile.Valid() {
		err = DUAConfigValidityErr
		return
	} else if !profile.Dedicated() {
		err = RegistrantPolicyErr
		return
	}

	for i := 0; i < len(meths) && err == nil; i++ {
		reg := profile.NewRegistrant()
		if err = reg.Marshal(meths[i]); err == nil {
			*r = append(*r, reg)
		}
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *Registrants) IsZero() bool {
	return r == nil
}

/*
Push appends a non-zero instance of *[Registrant] to the receiver instance.
*/
func (r *Registrants) Push(ath *Registrant) {
	if !r.IsZero() && !ath.IsZero() {
		*r = append(*r, ath)
	}
}

/*
Len returns the integer length of the receiver instance.
*/
func (r Registrants) Len() int {
	return len(r)
}

/*
Get returns an instance of *[Registrant] following a search for a matching
"[registrantID]". Case is significant in the matching process. A zero
instance is returned if not found.

[registrantID]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.34
*/
func (r Registrants) Get(id string) (ath *Registrant) {
	for i := 0; i < r.Len(); i++ {
		if id == r[i].ID() {
			ath = r[i]
			break
		}
	}

	return
}

/*
Contains returns a Boolean value indicative of a positive match between
the input "[registrantID]" value and one of the *[Registrant] instances
in the receiver instance. Case is significant in the matching process. A
zero instance is returned if not found.

[registrantID]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.34
*/
func (r Registrants) Contains(id string) bool {
	return !r.Get(id).IsZero()
}

/*
for test streamlining only; no practical use otherwise.
*/
type authority interface {
	CN() string
	L() string
	O() string
	C() string
	CO() string
	ST() string
	Tel() string
	Fax() string
	Title() string
	Email() string
	POBox() string
	PostalAddress() string
	PostalCode() string
	Mobile() string
	Street() string
	URI() []string
}
