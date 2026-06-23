package radir

/*
X660 implements [RASCHEMA § 2.5.4] and derives various concepts from
[ITU-T Rec. X.660].

	( 1.3.6.1.4.1.56521.101.2.5.4
	    NAME 'x660Context'
	    DESC 'X.660 contextual class'
	    SUP registration AUXILIARY
	    MAY ( additionalUnicodeValue $
	          currentAuthority $
	          firstAuthority $
	          secondaryIdentifier $
	          sponsor $
	          standardizedNameForm $
	          unicodeValue ) )

Instances of this type need not be initialized by the user directly.

[RASCHEMA § 2.5.4]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.4
[ITU-T Rec. X.660]: https://www.itu.int/rec/T-REC-X.660
*/
type X660 struct {
	R_UVal    string   `ldap:"unicodeValue"`           // RASCHEMA § 2.3.5
	R_AddlUV  []string `ldap:"additionalUnicodeValue"` // RASCHEMA § 2.3.6
	R_SecId   []string `ldap:"secondaryIdentifier"`    // RASCHEMA § 2.3.8
	R_StdNF   []string `ldap:"standardizedNameForm"`   // RASCHEMA § 2.3.18
	R_LongArc []string `ldap:"longArc"`                // RASCHEMA § 2.3.20

	// NON-COLLECTIVE DEDICATED entry DNs
	R_FAuthyDN []string `ldap:"firstAuthority"`   // RASCHEMA § 2.3.54
	R_CAuthyDN []string `ldap:"currentAuthority"` // RASCHEMA § 2.3.35
	R_SAuthyDN []string `ldap:"sponsor"`          // RASCHEMA § 2.3.74

	// COLLECTIVE DEDICATED registrant entry DN(s)
	RC_FAuthyDN []string `ldap:"c-firstAuthority;collective"`   // RASCHEMA § 2.3.55
	RC_CAuthyDN []string `ldap:"c-currentAuthority;collective"` // RASCHEMA § 2.3.36
	RC_SAuthyDN []string `ldap:"c-sponsor;collective"`          // RASCHEMA § 2.3.75

	// COMBINED registrant entry stubs.
	R_CFAuthy *FirstAuthority   // RASCHEMA § 2.3.37-53
	R_CCAuthy *CurrentAuthority // RASCHEMA § 2.3.56-73
	R_CSAuthy *Sponsor          // RASCHEMA § 2.3.76-93

	r_DITProfile *DITProfile
	r_root       *registeredRoot // linked from *Registration during init
	r_se         bool
}

/*
X660 returns (and if needed, initializes) the embedded instance of *[X660].
*/
func (r *Registration) X660() *X660 {
	if r.IsZero() {
		return &X660{}
	}

	if r.R_X660.IsZero() {
		r.R_X660 = new(X660)
		r.R_X660.r_DITProfile = r.Profile()
		r.R_X660.r_root = r.r_root
	}

	return r.R_X660
}

/*
X660 returns (and if needed, initializes) the embedded instance of *[X660].
*/
func (r *Subentry) X660() *X660 {
	if r.IsZero() {
		return &X660{}
	}

	if r.R_X660.IsZero() {
		r.R_X660 = new(X660)
		r.R_X660.r_DITProfile = r.Profile()
		r.R_X660.r_root = r.r_root
		r.R_X660.r_se = true
	}

	return r.R_X660
}

/*
profile returns the *[DITProfile] instance assigned to the receiver,
if set, else nil is returned.
*/
func (r *X660) profile() (prof *DITProfile) {
	if prof = r.r_DITProfile; !prof.Valid() {
		prof = &DITProfile{}
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *X660) IsZero() bool {
	return r == nil
}

func (r *X660) isEmpty() bool {
	return structEmpty(r)
}

func (r *X660) ldif() (l string) {
	if !r.IsZero() {
		l = toLDIF(r)
	}

	return
}

/*
unmarshal returns an instance of map[string][]string bearing the contents
of the receiver.
*/
func (r *X660) unmarshal() map[string][]string {
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
func (r *X660) marshal(meth func(any) error) (err error) {
	if !r.IsZero() {
		if !r.r_DITProfile.Combined() {
			err = meth(r)
		} else {
			err = meth(r)
			for _, err = range []error{
				r.CombinedFirstAuthority().marshal(meth),
				r.CombinedCurrentAuthority().marshal(meth),
				r.CombinedSponsor().marshal(meth),
			} {
				if err != nil {
					break
				}
			}
		}
	}

	return
}

/*
LongArc returns the string long arc values assigned to the receiver instance,
each defining a "[longArc]" associated with the registration.

Not that only sub arcs below Joint-ISO-ITU-T (2) are permitted to possess
long arc values. Assignment of a "[longArc]" to an ITU-T (0) or ISO (1)
registration is considered illegal.

[longArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.20
*/
func (r *X660) LongArc() (larc []string) {
	if !r.IsZero() {
		if r.r_root.Depth > 1 && r.r_root.N == 2 {
			larc = r.R_LongArc
		}
	}

	return
}

/*
SetLongArc assigns one or more string long arc values to the receiver instance.
Note that if a slice is passed as X, the destination value will be clobbered.
*/
func (r *X660) SetLongArc(args ...any) error {
	return writeFieldByTag(`longArc`, r.SetLongArc, r, args...)
}

/*
LongArcGetFunc processes the underlying long arc field values through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *X660) LongArcGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `longArc`)
}

/*
UnicodeValue returns the string Unicode value assigned to the receiver
instance.
*/
func (r *X660) UnicodeValue() string {
	return r.R_UVal
}

/*
SetUnicodeValue assigns a string Unicode values to the receiver instance.
*/
func (r *X660) SetUnicodeValue(args ...any) error {
	return writeFieldByTag(`unicodeValue`, r.SetUnicodeValue, r, args...)
}

/*
UnicodeValueGetFunc processes the underlying string Unicode field value
through the provided [GetOrSetFunc] instance, returning an interface value
alongside an error.
*/
func (r *X660) UnicodeValueGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `unicodeValue`)
}

/*
AdditionalUnicodeValue returns the string Unicode values assigned to the
receiver instance, each defining a unicodeValue assigned to the registration.
*/
func (r *X660) AdditionalUnicodeValue() []string {
	return r.R_AddlUV
}

/*
SetAdditionalUnicodeValue appends one or more additional string Unicode
values to the receiver instance. Note that if a slice is passed as X, the
destination value will be clobbered.
*/
func (r *X660) SetAdditionalUnicodeValue(args ...any) error {
	return writeFieldByTag(`additionalUnicodeValue`, r.SetAdditionalUnicodeValue, r, args...)
}

/*
AdditionalUnicodeValueGetFunc processes the underlying additional string
Unicode field values through the provided [GetOrSetFunc] instance, returning
an interface value alongside an error.
*/
func (r *X660) AdditionalUnicodeValueGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `additionalUnicodeValue`)
}

/*
SecondaryIdentifier returns the secondary identifier string values assigned
to the receiver instance, each defining an alternative non-Unicode, non-integer
identifier by which the receiver is known.
*/
func (r *X660) SecondaryIdentifier() []string {
	return r.R_SecId
}

/*
SetSecondaryIdentifier appends one or more secondary identifier string
values to the receiver instance. Note that if a slice is passed as X, the
destination value will be clobbered.
*/
func (r *X660) SetSecondaryIdentifier(args ...any) error {
	return writeFieldByTag(`secondaryIdentifier`, r.SetSecondaryIdentifier, r, args...)
}

/*
SecondaryIdentifierGetFunc processes the underlying secondary identifier
field values through the provided [GetOrSetFunc] instance, returning an
interface value alongside an error.
*/
func (r *X660) SecondaryIdentifierGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `secondaryIdentifier`)
}

/*
StdNameForm returns the standardized name form string values assigned to
the receiver instance, each defining a standardized name form assigned to
the registration.
*/
func (r *X660) StdNameForm() []string {
	return r.R_StdNF
}

/*
SetStdNameForm assigns one or more standardized name form string values
to the receiver instance. Note that if a slice is passed as X, the
destination value will be clobbered.
*/
func (r *X660) SetStdNameForm(args ...any) error {
	return writeFieldByTag(`standardizedNameForm`, r.SetStdNameForm, r, args...)
}

/*
StdNameFormGetFunc processes the underlying standardized name form field
values through the provided [GetOrSetFunc] instance, returning an interface
value alongside an error.
*/
func (r *X660) StdNameFormGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `standardizedNameForm`)
}

/*
CurrentAuthorities returns the string DN values assigned to the receiver
instance, each referencing a specific *[CurrentAuthority] entry by DN.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) CurrentAuthorities() (val []string) {
	if r.profile().Dedicated() {
		val = r.R_CAuthyDN
	}

	return
}

/*
SetCurrentAuthorities appends one or more string DN values to the receiver
instance. Note that if a slice is passed as X, the destination value will
be clobbered.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) SetCurrentAuthorities(args ...any) error {
	if !r.profile().Dedicated() {
		return RegistrantPolicyErr
	}
	return writeFieldByTag(`currentAuthority`, r.SetCurrentAuthorities, r, args...)
}

/*
SetCCurrentAuthorities appends one or more string DN values to the receiver
instance. Note that if a slice is passed as X, the destination value will
be clobbered.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy" and only when the receiver
instance was initialized within an instance of *[Subentry].
*/
func (r *X660) SetCCurrentAuthorities(args ...any) error {
	if !r.profile().Dedicated() {
		return RegistrantPolicyErr
	}
	return writeFieldByTag(`c-currentAuthority;collective`, r.SetCurrentAuthorities, r, args...)
}

/*
CurrentAuthoritiesGetFunc processes the underlying string DN field values
through the provided [GetOrSetFunc] instance, returning an interface value
alongside an error.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) CurrentAuthoritiesGetFunc(getfunc GetOrSetFunc) (any, error) {
	if !r.profile().Dedicated() {
		return nil, RegistrantPolicyErr
	}
	return getFieldValueByNameTagAndGoSF(r, getfunc, `currentAuthority`)
}

/*
FirstAuthorities returns the string DN values assigned to the receiver
instance, each referencing a specific *[FirstAuthority] entry by DN.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) FirstAuthorities() (val []string) {
	if r.profile().Dedicated() {
		val = r.R_FAuthyDN
	}

	return
}

/*
SetFirstAuthorities appends one or more string DN values to the receiver
instance. Note that if a slice is passed as X, the destination value will
be clobbered.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) SetFirstAuthorities(args ...any) error {
	if !r.profile().Dedicated() {
		return RegistrantPolicyErr
	}

	return writeFieldByTag(`firstAuthority`, r.SetFirstAuthorities, r, args...)
}

/*
SetCFirstAuthorities appends one or more string DN values to the receiver
instance. Note that if a slice is passed as X, the destination value will
be clobbered.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy", and only when the receiver
instance was initialized within an instance of *[Subentry].
*/
func (r *X660) SetCFirstAuthorities(args ...any) error {
	if !r.profile().Dedicated() {
		return RegistrantPolicyErr
	}

	return writeFieldByTag(`c-firstAuthority;collective`, r.SetFirstAuthorities, r, args...)
}

/*
FirstAuthoritiesGetFunc processes the underlying string DN field values
through the provided [GetOrSetFunc] instance, returning an interface
value alongside an error.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) FirstAuthoritiesGetFunc(getfunc GetOrSetFunc) (any, error) {
	if !r.profile().Dedicated() {
		return nil, RegistrantPolicyErr
	}
	return getFieldValueByNameTagAndGoSF(r, getfunc, `firstAuthority`)
}

/*
CFirstAuthorities returns the string DN values assigned to the receiver instance,
each referencing a specific *[FirstAuthority] entry by DN.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) CFirstAuthorities() (val []string) {
	if r.profile().Dedicated() {
		val = r.RC_FAuthyDN
	}

	return
}

/*
CFirstAuthoritiesGetFunc processes the underlying string DN field values
through the provided [GetOrSetFunc] instance, returning an interface
value alongside an error.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) CFirstAuthoritiesGetFunc(getfunc GetOrSetFunc) (any, error) {
	if !r.profile().Dedicated() {
		return nil, RegistrantPolicyErr
	}

	return getFieldValueByNameTagAndGoSF(r, getfunc, `c-firstAuthority`)
}

/*
CCurrentAuthorities returns the string DN values assigned to the receiver instance,
each referencing a specific *[CurrentAuthority] entry by DN.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) CCurrentAuthorities() (val []string) {
	if r.profile().Dedicated() {
		val = r.RC_CAuthyDN
	}

	return
}

/*
CCurrentAuthoritiesGetFunc processes the underlying string DN field values
through the provided [GetOrSetFunc] instance, returning an interface
value alongside an error.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) CCurrentAuthoritiesGetFunc(getfunc GetOrSetFunc) (any, error) {
	if !r.profile().Dedicated() {
		return nil, RegistrantPolicyErr
	}

	return getFieldValueByNameTagAndGoSF(r, getfunc, `c-currentAuthority`)
}

/*
Sponsors returns the string DN values assigned to the receiver instance,
each referencing a specific *[Sponsor] entry by DN.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) Sponsors() (val []string) {
	if r.profile().Dedicated() {
		val = r.R_SAuthyDN
	}

	return
}

/*
SetSponsors appends one or more string DN values to the receiver instance.
Note that if a slice is passed as X, the destination value will be clobbered.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) SetSponsors(args ...any) error {
	if !r.profile().Dedicated() {
		return RegistrantPolicyErr
	}

	return writeFieldByTag(`sponsor`, r.SetSponsors, r, args...)
}

/*
SetCSponsors appends one or more string DN values to the receiver instance.
Note that if a slice is passed as X, the destination value will be clobbered.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy" and only when the receiver
instance was initialized within an instance of *[Subentry].
*/
func (r *X660) SetCSponsors(args ...any) error {
	if !r.profile().Dedicated() {
		return RegistrantPolicyErr
	}

	return writeFieldByTag(`c-sponsor;collective`, r.SetSponsors, r, args...)
}

/*
SponsorGetFunc processes the underlying string DN field values through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) SponsorsGetFunc(getfunc GetOrSetFunc) (any, error) {
	if !r.profile().Dedicated() {
		return nil, RegistrantPolicyErr
	}
	return getFieldValueByNameTagAndGoSF(r, getfunc, `sponsor`)
}

/*
CSponsors returns the string DN values assigned to the receiver instance,
each referencing a specific *[Sponsor] entry by DN.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) CSponsors() (val []string) {
	if r.profile().Dedicated() {
		val = r.RC_SAuthyDN
	}

	return
}

/*
CSponsorGetFunc processes the underlying string DN field values through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.

Note that use of this method is only useful when operating under the
terms of the "Dedicated Registrants Policy".
*/
func (r *X660) CSponsorsGetFunc(getfunc GetOrSetFunc) (any, error) {
	if !r.profile().Dedicated() {
		return nil, RegistrantPolicyErr
	}

	return getFieldValueByNameTagAndGoSF(r, getfunc, `c-sponsor`)
}

/*
CombinedFirstAuthority returns (and if needed, initializes) the embedded
instance of *[FirstAuthority].

This method is intended solely for use under the terms of the "Combined
Registrants Policy".
*/
func (r *X660) CombinedFirstAuthority() (dr *FirstAuthority) {
	if !r.r_DITProfile.Combined() {
		return &FirstAuthority{}
	}

	if r.R_CFAuthy.IsZero() {
		r.R_CFAuthy = new(FirstAuthority)
		r.R_CFAuthy.r_alt_types = r.profile().r_alt_types
	}

	return r.R_CFAuthy
}

/*
CombinedCurrentAuthority returns (and if needed, initializes) the embedded
instance of *[CurrentAuthority].

This method is intended solely for use under the terms of the "Combined
Registrants Policy".
*/
func (r *X660) CombinedCurrentAuthority() *CurrentAuthority {
	if !r.r_DITProfile.Combined() {
		return &CurrentAuthority{}
	}

	if r.R_CCAuthy.IsZero() {
		r.R_CCAuthy = new(CurrentAuthority)
		r.R_CCAuthy.r_alt_types = r.profile().r_alt_types
	}

	return r.R_CCAuthy
}

/*
CombinedSponsor returns (and if needed, initializes) the embedded instance
of *[Sponsor].

This method is intended solely for use under the terms of the "Combined
Registrants Policy".
*/
func (r *X660) CombinedSponsor() *Sponsor {
	if !r.r_DITProfile.Combined() {
		return &Sponsor{}
	}

	if r.R_CSAuthy.IsZero() {
		r.R_CSAuthy = new(Sponsor)
		r.R_CSAuthy.r_alt_types = r.profile().r_alt_types
	}

	return r.R_CSAuthy
}

// Ensure certain values are permitted to be set
func (r *X660) writeEligible(tag string, value any) (err error) {
	switch lc(tag) {
	case `longarc`:
		if r.r_root == nil {
			err = NilInstanceErr
		} else if !(r.r_root.Depth > 1 && r.r_root.N == 2) {
			err = LongArcErr
		}
	}

	return
}
