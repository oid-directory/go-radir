package radir

/*
cfg.go handles all elements pertaining to RA DUA configuration.
*/

/*
SupportedModels is a string slice global variable meant to house all
directory model numeric OIDs which are supported for the implementation
in question.

At the time of this writing, only two (2) models exist:

  - [TwoDimensional]
  - [ThreeDimensional]

However, future amendments or extensions of the I-D series could include
other model OIDs, which may be added here.

Users may append to this variable instance, however they should not remove
slices during the runtime session, as it can diminish functionality.
*/
var SupportedModels []string

/*
DUAConfig implements [Section 2.2.2 of the RADUA I-D], [Sections 2.3.94]
through [2.3.100] and [2.5.17 of the RASCHEMA I-D], and [Sections 3.2.4.14]
and [3.2.4.15 of the RADIT I-D] for the purpose of facilitating an RA DUA
configuration stack, derived from the target DSA's Root DSE, and possibly
other sources.

See also [DITProfile] for the compartmentalized element used within
instances of this type.

[Section 2.2.2 of the RADUA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radua#section-2.2.2
[Sections 2.3.94]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.94
[2.3.100]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.100
[2.5.17 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.17
[Sections 3.2.4.14]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.2.4.14
[3.2.4.15 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.2.4.15
*/
type DUAConfig struct {
	R_DSE *DITProfile // for RA DSA's Root DSE (singular) profile

	// R_Profile[DN]s fields are only used for
	// non DSE-based profiling.
	R_ProfileDNs []string      `ldap:"rADITProfile"`
	R_Profiles   []*DITProfile // object forms of profiled entries
}

/*
DITProfile implements an abstraction of [Section 2.3.94 of the RASCHEMA I-D].

A sole instance of this type is derived from the RA DSA's Root DSE, assuming
the "[rADUAConfig]" AUXILIARY class is present within the DSE.

Additional instances of this type represent independent "[rADUAConfig]"
instances, likely derived from a special location within any backend of
the the RA DSA -- OTHER THAN the Root DSE.

Use of additional instances of this type is ONLY needed in particularly
special use-cases in which multiple separate RA DITs reside on the same
RA DSA and are governed independently of one another. Most implementations
(namely single DIT implementations) will be able to use a sole Root DSE
derived [DUAConfig].

This MAY apply in certain environments in which there are so-called
"staging" DITs alongside "active" DITs, where "unvetted" registrations
are written to "staging" for review (and alterations, if necessary) prior
to relocation to the "active" DIT following positive approval.

NOTE: Profile DNs can ONLY be referenced from within the Root DSE-derived
configuration instance. The act of "subprofiling" -- i.e.: referencing a
profile from another profile -- is not supported. This is precisely why
the "[rADITProfile]" type is NOT present as a struct field here.

[rADUAConfig]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.17
[rADITProfile]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.94
[Section 2.3.94 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.94
*/
type DITProfile struct {
	R_DN       string   `ldap:"dn"`
	R_GSR      string   `ldap:"governingStructureRule"`
	R_TTL      string   `ldap:"rATTL"`              // RASCHEMA 2.3.100
	R_Model    string   `ldap:"rADirectoryModel"`   // RASCHEMA 2.3.97
	R_RegBase  []string `ldap:"rARegistrationBase"` // RASCHEMA 2.3.95
	R_AthyBase []string `ldap:"rARegistrantBase"`   // RASCHEMA 2.3.96
	R_Mail     []string `ldap:"rAServiceMail"`      // RASCHEMA 2.3.98
	R_URI      []string `ldap:"rAServiceURIs"`      // RASCHEMA 2.3.99
	R_OC       []string `ldap:"objectClass"`

	// The collective form of the TTL (c-rATTL) is not
	// applicable to an instance of this type that is
	// an abstraction of the Root DSE.
	RC_TTL string `ldap:"c-rATTL;collective"` // RASCHEMA 2.3.101

	// No concept of arbitrary "settings" for an RA DUA is
	// officially defined in the RADUA I-D series, but is
	// nonetheless sensible to include here for reasons of
	// client optimization, if needed.
	R_Settings *ProfileSettings

	// Cache contains the profile's cache instance, if enabled.
	// Caching is covered in Section 2.2.3.4 of the RADUA I-D.
	r_Cache *Cache

	// Make a note of our dedicated type policy (draft or RFC).
	r_alt_types bool

	r_bsel [2]int
}

/*
Valid returns a boolean value indicative of whether the receiver
configuration instance is considered contextually valid and usable.

Usability is determined based on all of the following:

  - The receiver is not nil
  - The numeric OID of a supported directory model has been specified
  - The number of registration bases is non-zero
*/
func (r *DITProfile) Valid() (valid bool) {
	if !r.IsZero() {
		if strInSlice(r.R_Model, SupportedModels) {
			valid = len(r.R_RegBase) > 0
		}
	}

	return
}

/*
TTL returns the string time-to-live value associated with the entry of the
indicated profile.
*/
func (r *DITProfile) TTL() string {
	return selectTTL(r.R_TTL, r.RC_TTL)
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
func (r *DITProfile) GoverningStructureRule() (gsr string) {
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
func (r *DITProfile) GoverningStructureRuleGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `governingStructureRule`)
}

/*
TTLGetFunc processes the underlying TTL field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *DITProfile) TTLGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `rATTL`)
}

/*
NewDUAConfig returns an initialized instance of *[DUAConfig] containing
a single *[DITProfile] set as the default "DSE" struct field value.

No parameters are added, making this suitable for use in automatic DUA
configuration, as defined by the attribute value within the RA DSA's Root
DSE. This method is also suitable for those who configure their client
manually.

See also [NewFactoryDefaultDUAConfig].
*/
func NewDUAConfig() *DUAConfig {
	return &DUAConfig{
		R_DSE: &DITProfile{
			R_Settings: newProfileSettings(),
		},
		R_Profiles: make([]*DITProfile, 0),
	}
}

/*
NewFactoryDefaultDUAConfig is a convenience function that simply inputs
the default DUA configuration parameters suggested throughout the examples
defined within the OID Directory I-D series.

This method wraps [NewDUAConfig] and adjusts the parameters as follows:

  - Assume Root DSE is used (single profile)
  - Directory Model: [ThreeDimensional]
  - Registration Base (1): ou=Registrations,o=rA
  - Registrant Base (1): ou=Registrants,o=rA

The result is a 3D configuration using the "Dedicated Registrants Policy"
that is suitable for use in Quick-Start scenarios.
*/
func NewFactoryDefaultDUAConfig() (d *DUAConfig) {
	d = NewDUAConfig()
	d.R_DSE.SetModel(ThreeDimensional)
	d.R_DSE.SetRegistrationBase("ou=Registrations,o=rA")
	d.R_DSE.SetRegistrantBase("ou=Registrants,o=rA")

	return d
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
func (r *DITProfile) Marshal(meth func(any) error) (err error) {
	if !r.IsZero() {
		err = meth(r)
	}

	return
}

/*
LDIF returns the string LDIF form of the receiver instance. Note that this
is a crude approximation of LDIF and should ideally be parsed through a
reliable LDIF parser such as [go-ldap/ldif] to verify integrity.

[go-ldap/ldif]: https://pkg.go.dev/github.com/go-ldap/ldif
*/
func (r *DITProfile) LDIF() (l string) {
	if !r.IsZero() {
		dn := readFieldByTag(`dn`, r)
		oc := readFieldByTag(`objectClass`, r)

		bld := newBuilder()

		if len(dn) > 0 {
			bld.WriteString(`dn: ` + dn[0])
			bld.WriteRune(10)
		} else {
			// Root DSE has no DN.
			bld.WriteString(`dn:`)
			bld.WriteRune(10)
		}

		for i := 0; i < len(oc); i++ {
			bld.WriteString(`objectClass: ` + oc[i])
			bld.WriteRune(10)
		}

		bld.WriteString(toLDIF(r))
		bld.WriteRune(10)

		l = bld.String()
	}

	return
}

/*
Profile returns the Nth instance of *[DITProfile] found within the receiver
instance.

If an index of 0 is provided and the receiver only contains a single profile
derived from the Root DSE, that instance is returned. This is equivalent to
simply accessing the "DSE" struct field, and represents the most common use case.

If, however, the number of separate profiles is non-zero, the input index
integer i will result in the return of the Nth underlying *[DITProfile].
In this case, the "DSE" struct field will be nil, as any reference to
external *[DITProfile] instances supersedes use of the "DSE"-based profile.

If no input is provided whatsoever, this is equivalent to an input of 0,
regardless of the above conditions.

A nil instance is returned in all other scenarios.
*/
func (r *DUAConfig) Profile(idx ...int) (dp *DITProfile) {
	var i int
	if len(idx) > 0 {
		i = idx[0]
	}

	if dp = r.R_DSE; dp.IsZero() {
		if lp := r.NumProfile(); lp > 0 {
			if 0 <= i && i <= lp {
				dp = r.R_Profiles[i]
			}
		}
	}

	return
}

/*
NumProfile returns the integer number of underlying (non-nil) *[DITProfile]
instances.
*/
func (r *DUAConfig) NumProfile() (num int) {
	if num = len(r.R_Profiles); num == 0 {
		if !r.R_DSE.IsZero() {
			num = 1
		}
	}

	return
}

/*
Settings returns the underlying instance of [ProfileSettings].
*/
func (r *DITProfile) Settings() *ProfileSettings {
	return r.R_Settings
}

/*
UseAltAuthorityTypes assigns the input alt Boolean value as the Authority
Attribute Type Policy "indicator". If and when this setting is used for a
specific *[DITProfile] instance, it should NOT be changed after the fact,
as this will result not only in quirky RA DUA behavior, but also an RA DIT
that is less than orderly. It is best to decide this policy prior to its
engagement.

An input value of true will result in the draft-provided types, such as
"[sponsorCommonName]", remaining unused in favor of standard definitions
like the RFC 4519 "[cn]" type.

Note that novel types -- namely the draft-provided "[Generalized Time]"
types -- have no official equivalents. As such, their continued use is
almost certainly required.

Also note that continued use of the appropriate AUXILIARY authority class(es)
is also likely required, specifically ONE (1) of:

  - "[firstAuthorityContext]"
  - "[currentAuthorityContext]"
  - "[sponsorContext]"

By utilizing the alternative type policy, it is no longer possible to have
more than one (1) of the above contexts assigned to a single *[Registrant]
or *[Registration] instance due to ambiguity issues. For example, it would
not be practical to use "[o]" for both "[firstAuthorityContext]" in addition
to "[sponsorContext]", as there is no reliable means for determining which
context "owns" which value of "[o]".  This is mainly only an issue in cases
where one is somehow unmarshaling information to an instance of [go-ldap/v3
Entry], but should ideally be avoided in all cases.

Lack of any of the necessary contextual classes (when expected) may result
in both diminished RA DUA functionality and related search request effectiveness.

This setting has no effect when the *[DITProfile] instance involved does
not support any *[Registrant] operations (i.e.: in *[Registration]-only
scenarios).

[Generalized Time]: https://datatracker.ietf.org/doc/html/rfc4517#section-3.3.13
[sponsorCommonName]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.78
[cn]: https://datatracker.ietf.org/doc/html/rfc4519#section-2.3
[o]: https://datatracker.ietf.org/doc/html/rfc4519#section-2.19
[firstAuthorityContext]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.13
[currentAuthorityContext]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.14
[sponsorContext]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.15
[go-ldap/v3 Entry]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry
*/
func (r *DITProfile) UseAltAuthorityTypes(alt bool) {
	if !r.IsZero() {
		r.r_alt_types = alt
	}
}

/*
RegistrationTarget sets the preferred Registration Base within the receiver.

Subsequent calls of [DITProfile.RegistrationBase] will return the string
representation of the selected Registration Base DN. Note that this method
has no effect when only a single such base exists.
*/
func (r *DITProfile) RegistrationTarget(base int) {
	if !r.IsZero() {
		r.r_bsel[0] = base
	}
}

/*
RegistrantBase returns the string representation of the Registration Base
DN currently selected within the receiver instance.
*/
func (r DITProfile) RegistrationBase() (base string) {
	if !r.IsZero() {
		if r.NumRegistrationBase() == 1 {
			base = r.registrationBase()
		} else {
			base = r.registrationBase(r.r_bsel[0])
		}
	}

	return
}

/*
RegistrantTarget sets the preferred Registrant Base within the receiver.

Subsequent calls of [DITProfile.RegistrantBase] will return the string
representation of the selected Registrant Base DN. Note that this method
has no effect when only a single such base exists.
*/
func (r *DITProfile) RegistrantTarget(base int) {
	if !r.IsZero() {
		r.r_bsel[1] = base
	}
}

/*
RegistrantBase returns the string representation of the Registrant Base
DN currently selected within the receiver instance.
*/
func (r DITProfile) RegistrantBase() (base string) {
	if !r.IsZero() {
		if r.NumRegistrantBase() == 1 {
			base = r.registrantBase()
		} else {
			base = r.registrantBase(r.r_bsel[1])
		}
	}

	return
}

/*
registrationBase returns the Nth "[rARegistrationBase]" value. If no index
value is provided, the 0th slice will be returned.

Any negative value results in a zero return value.

If an index value greater than the length of the "[rARegistrationBase]",
a zero string is returned.

If no "[rARegistrationBase]" slices are present, a zero string is returned.

[rARegistrationBase]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.95
*/
func (r *DITProfile) registrationBase(idx ...int) (base string) {
	L := r.NumRegistrationBase()
	if L == 0 {
		return
	}

	var i int
	if len(idx) > 0 {
		i = idx[0]
	}

	if L >= i && i >= 0 {
		base = r.R_RegBase[i]
	}

	return
}

/*
RegistrationBaseGetFunc executes the [GetOrSetFunc] instance and returns
its own return values. The 'any' value will require type assertion in
order to be accessed reliably. An error is returned if issues arise.
*/
func (r *DITProfile) RegistrationBaseGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `rARegistrationBase`)
}

/*
SetRegistrationBase assigns input value b as an "[rARegistrationBase]" instance.

[rARegistrationBase]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.95
*/
func (r *DITProfile) SetRegistrationBase(b string) {
	if !r.IsZero() {
		if !strInSlice(b, r.R_RegBase) && len(b) > 0 {
			r.R_RegBase = append(r.R_RegBase, b)
		}
	}
}

/*
NumRegistrationBase returns the integer number of "[rARegistrationBase]"
instances present within within the receiver instance.

[rARegistrationBase]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.95
*/
func (r *DITProfile) NumRegistrationBase() (l int) {
	if !r.IsZero() {
		l = len(r.R_RegBase)
	}

	return
}

/*
registrantBase returns the Nth "[rARegistrantBase]" value. If no index value
is provided, the 0th slice will be returned.

Any negative value results in a zero return value.

If an index value greater than the length of the "[rARegistrantBase]", a zero
string is returned.

If no "[rARegistrantBase]" slices are present, a zero string is returned.

[rARegistrantBase]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.96
*/
func (r *DITProfile) registrantBase(idx ...int) (base string) {
	L := r.NumRegistrantBase()
	if L == 0 {
		return
	}

	var i int
	if len(idx) > 0 {
		i = idx[0]
	}

	if L >= i && i >= 0 {
		base = r.R_AthyBase[i]
	}

	return
}

/*
RegistrantBaseGetFunc executes the [GetOrSetFunc] instance and returns
its own return values. The 'any' value will require type assertion in
order to be accessed reliably. An error is returned if issues arise.
*/
func (r *DITProfile) RegistrantBaseGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `rARegistrantBase`)
}

/*
SetRegistrantBase assigns input value b as an "[rARegistrantBase]" instance.

[rARegistrantBase]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.96
*/
func (r *DITProfile) SetRegistrantBase(b string) {
	if !r.IsZero() {
		if !strInSlice(b, r.R_AthyBase) && len(b) > 0 {
			r.R_AthyBase = append(r.R_AthyBase, b)
		}
	}
}

/*
NumRegistrantBase returns the integer number of "[rARegistrantBase]" instances
present within within the receiver instance.

[rARegistrantBase]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.96
*/
func (r *DITProfile) NumRegistrantBase() (l int) {
	if !r.IsZero() {
		l = len(r.R_AthyBase)
	}

	return
}

/*
Mail returns the Nth "[rAServiceMail]" value. If no index value is provided,
the 0th slice will be returned.

Any negative value results in a zero return value.

If an index value greater than the length of the "[rAServiceMail]", a zero
string is returned.

If no "[rAServiceMail]" slices are present, a zero string is returned.

[rAServiceMail]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.98
*/
func (r *DITProfile) Mail(idx ...int) (m string) {
	L := r.NumMail()
	if L == 0 {
		return
	}

	var i int
	if len(idx) > 0 {
		i = idx[0]
	}

	if L >= i && i >= 0 {
		m = r.R_Mail[i]
	}

	return
}

/*
MailGetFunc executes the [GetOrSetFunc] instance and returns its own
return values. The 'any' value will require type assertion in order
to be accessed reliably. An error is returned if issues arise.
*/
func (r *DITProfile) MailGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `rAServiceMail`)
}

/*
SetMail assigns input value m as an "[rAServiceMail]" instance.

[rAServiceMail]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.98
*/
func (r *DITProfile) SetMail(m string) {
	if !strInSlice(m, r.R_Mail) && len(m) > 0 {
		r.R_Mail = append(r.R_Mail, m)
	}
}

/*
NumURI returns the integer number of "[rAServiceMail]" instances present
within within the receiver instance.

[rAServiceMail]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.98
*/
func (r *DITProfile) NumMail() int {
	return len(r.R_Mail)
}

/*
URI returns the Nth "[rAServiceURI]" value. If no index value is provided,
the 0th slice will be returned.

Any negative value results in a zero return value.

If an index value greater than the length of the "[rAServiceURI]", a zero
string is returned.

If no "[rAServiceURI]" slices are present, a zero string is returned.

[rAServiceURI]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.99
*/
func (r *DITProfile) URI(idx ...int) (u string) {
	L := r.NumURI()
	if L == 0 {
		return
	}

	var i int
	if len(idx) > 0 {
		i = idx[0]
	}

	if L >= i && i >= 0 {
		u = r.R_URI[i]
	}

	return
}

/*
URIGetFunc executes the [GetOrSetFunc] instance and returns its own
return values. The 'any' value will require type assertion in order
to be accessed reliably. An error is returned if issues arise.
*/
func (r *DITProfile) URIGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `rAServiceURI`)
}

/*
SetURI assigns input value u as an "[rAServiceURI]" instance.

[rAServiceURI]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.99
*/
func (r *DITProfile) SetURI(u string) {
	if !strInSlice(u, r.R_URI) && len(u) > 0 {
		r.R_URI = append(r.R_URI, u)
	}
}

/*
NumURI returns the integer number of "[rAServiceURI]" instances present
within within the receiver instance.

[rAServiceURI]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.99
*/
func (r *DITProfile) NumURI() int {
	return len(r.R_URI)
}

/*
Model returns the numeric OID of the "[rADirectoryModel]" employed to
advertise the structural layout of the RA DIT in question.

See [Sections 3.1.2] and [3.1.3 of the RADIT I-D] for (currently) valid
values for instances of this type.

Additional models extended within future amendments or extensions of the
series may be supported by adding the associated numeric OID to the global
[SupportedModels] slice variable.

If the model assigned to the receiver instance is unknown, a zero string
is returned.

[rADirectoryModel]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.97
[Sections 3.1.2]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.1.2
[3.1.3 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.1.3
*/
func (r *DITProfile) Model() (model string) {
	if strInSlice(r.R_Model, SupportedModels) {
		model = r.R_Model
	}

	return
}

/*
ModelGetFunc executes the [GetOrSetFunc] instance and returns its own
return values. The 'any' value will require type assertion in order
to be accessed reliably. An error is returned if issues arise.
*/
func (r *DITProfile) ModelGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `rADirectoryModel`)
}

/*
SetModel assigns the input model to the receiver instance, but only if
all of the following conditions evaluate as true:

  - The input model is a member of [SupportedModels]
  - The input model is non-zero in length
  - The receiver does not already have a model specified
*/
func (r *DITProfile) SetModel(model string) {
	if strInSlice(model, SupportedModels) &&
		len(r.R_Model) == 0 && len(model) > 0 {
		r.R_Model = model
	}
}

/*
RegistrationSuffixEqual returns an integer value indicative of which
configured *[Registration] search base within the receiver instance is
the suffix of the input distinguished name.

A -1 is returned if the input distinguished name is zero length, or
does not match any *[Registration] search base value.

Case is not significant in the matching process.
*/
func (r *DITProfile) RegistrationSuffixEqual(dn string) (index int) {
	return r.suffixBaseEqual(dn, 0)
}

/*
RegistrantSuffixEqual returns an integer value indicative of which
configured *[Registrant] search base within the receiver
instance is the suffix of the input distinguished name.

A -1 is returned if the input distinguished name is zero length, or
does not match any *[Registrant] search base value.

Case is not significant in the matching process.
*/
func (r *DITProfile) RegistrantSuffixEqual(dn string) (index int) {
	return r.suffixBaseEqual(dn, 1)
}

func (r *DITProfile) suffixBaseEqual(dn string, t int) (index int) {
	index = -1

	if dn = lc(dn); len(dn) > 0 {
		var RL int
		var funk func(...int) string
		if t == 0 {
			RL = r.NumRegistrationBase()
			funk = r.registrationBase
		} else {
			RL = r.NumRegistrantBase()
			funk = r.registrantBase
		}

		for i := 0; i < RL; i++ {
			if base := funk(i); hasSfx(dn, lc(base)) {
				index = i
				break
			}
		}
	}

	return
}

/*
Dedicated returns a Boolean value indicative of whether the receiver
instance operates on a "Dedicated Registrant" basis, whereas any and all
authority-related content uses dedicated directory entries which bear the
"[registrant]" STRUCTURAL class.

These authorities are referenced by any number of *[Registration] instances
by way of the one or more appropriate DN attribute values assigned to any of:

  - "[firstAuthority]" and/or "[c-firstAuthority]"
  - "[currentAuthority]" and/or "[c-currentAuthority]"
  - "[sponsor]" and/or "[c-sponsor]"

A value of false is returned if the "Combined Registrant Policy" is in
force, or if no [Registrant] content is permitted of any kind, per the
lack of any [Registrant] search bases present within the receiver instance.

This method is the mutex inverse of the [DITProfile.Combined] method.

[registrant]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.16
[firstAuthority]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.54
[c-firstAuthority]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.55
[currentAuthority]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.35
[c-currentAuthority]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.36
[sponsor]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.74
[c-sponsor]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.75
*/
func (r *DITProfile) Dedicated() bool {
	c, d := r.determineRegistrantPolicy()
	return c == 0 && d > 0
}

/*
Combined returns a Boolean value indicative of whether the receiver
instance operates under the terms of the "Combined Registrant Policy",
whereas any and all authority content is assigned "in-line" to the
*[Registration] entry with which it is associated directly. As such,
instances of the [Registrant] type are unused.

A value of false is returned the "Dedicated Registrant Policy" is in
force, or if no [Registrant] content is permitted of any kind, per the
lack of any [Registrant] search bases present within the receiver instance.

This method is the mutex inverse of the [DITProfile.Dedicated] method.
*/
func (r *DITProfile) Combined() bool {
	c, d := r.determineRegistrantPolicy()
	return c > 0 && d == 0
}

/*
AllowsRegistrants returns a Boolean value indicative of whether the receiver
instance supports [Registrant] handling of any kind. It is simply a succinct
alternative to determining whether any *[Registrant] policy is used.

A value of true indicates [Registrant] handling is allowed, and that one
(1) of [DITProfile.Dedicated] OR [DITProfile.Combined] should return true.

A value of false indicates NO [Registrant] handling of any kind is supported,
or that the policy is ambiguous. For example, the receiver instance can only
support "Combined Registrants" OR "Dedicated Registrants" OR neither -- but
NEVER both.

If there are variations in the individual search bases that would require
both, use of a separate *[DITProfile] may be indicated.
*/
func (r *DITProfile) AllowsRegistrants() bool {
	c, d := r.determineRegistrantPolicy()
	return c == 0 && d > 0 || c > 0 && d == 0
}

func (r *DITProfile) determineRegistrantPolicy() (int, int) {
	if r.IsZero() {
		return 0, 0
	} else if len(r.R_AthyBase) == 0 {
		return 0, 0
	} else if rb := len(r.R_RegBase); rb == 0 {
		// If there are NO registration
		// bases, but there is at least
		// one authority base, this HAS
		// to mean dedicated registrant
		// operation is in use. This
		// does NOT indicate the config
		// is otherwise valid!
		return 0, rb
	}

	var c, d int
	for _, athy := range r.R_AthyBase {
		if strInSlice(athy, r.R_RegBase) {
			// a registrant base matching a slice
			// within the registration base means
			// COMBINED authorities are in use.
			// Case-folding is not significant.
			c++
		} else if len(athy) > 0 {
			// Contrary to the above, a registrant
			// base that is NOT present within the
			// registration base list is considered
			// to be a DEDICATED authority base.
			d++
		}
	}

	return c, d
}

/*
NewSubentry initializes and returns a new instance of *[Subentry].
*/
func (r *DITProfile) NewSubentry() *Subentry {
	var oc []string = make([]string, 0)
	oc = append(oc, []string{
		`top`,
		`subentry`,
	}...)

	return &Subentry{
		R_OC:         oc,
		r_DITProfile: r,
	}
}

/*
NewRegistration returns a freshly initialized instance of *[Registration].

The variadic Boolean input value influences which STRUCTURAL "[objectClass]"
is used for the new instance:

  - "[rootArc]" (true)
  - "[arc]" (false)

As there are only three (3) possible root entries (0, 1 or 2), the most
common use-case involves no variadic input -- which is equivalent to an
explicit false -- to select the "[arc]" class.

A reference to the receiver instance is stored within the return.

[rootArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.2
[arc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.3
[objectClass]: https://www.rfc-editor.org/rfc/rfc4512.html#section-3.3
*/
func (r *DITProfile) NewRegistration(root ...bool) *Registration {
	var oc []string = make([]string, 0)
	oc = append(oc, []string{
		`top`,
		`registration`,
	}...)

	var soc string = `arc` // default
	if len(root) > 0 {
		if root[0] {
			soc = `rootArc`
		}
	}

	oc = append(oc, soc)

	return &Registration{
		R_OC:         oc,
		R_SOC:        soc,
		R_DITProfile: r,
		r_root:       new(registeredRoot),
	}
}

/*
NewRegistrant returns a freshly initialized instance of *[Registrant].

A reference to the receiver instance is stored within the return.

Note that a bogus instance is returned if the profile operates under the
terms of the "Combined Registrants Policy".
*/
func (r *DITProfile) NewRegistrant() *Registrant {
	if !r.Dedicated() {
		return &Registrant{}
	}

	var oc []string = make([]string, 0)
	var soc string = `registrant`
	oc = append(oc, []string{
		`top`,
		soc,
	}...)

	return &Registrant{
		R_OC:         oc,
		R_SOC:        soc,
		R_DITProfile: r,
	}
}

/*
MakeCache assigns a freshly initialized instance of *[Cache] to the receiver
instance. The *[Cache] instance may be accessed via the [DITProfile.Cache]
method.
*/
func (r *DITProfile) MakeCache(limits ...int) {
	var ai, ab int
	if len(limits) > 0 {
		ai = limits[0]
		if len(limits) > 1 {
			ab = limits[1]
		}
	}

	r.r_Cache = newCache(ai, ab)
}

/*
DropCache purges and frees the underlying instance of *[Cache] within the
receiver instance.
*/
func (r *DITProfile) DropCache() {
	if c := r.Cache(); !c.IsZero() {
		c.Flush()
		c.Free()
	}
}

/*
Cache returns the underlying instance of *[Cache] within the receiver
instance.
*/
func (r *DITProfile) Cache() *Cache {
	return r.r_Cache
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *DITProfile) IsZero() bool {
	return r == nil
}

/*
ProfileSettings is an optional, general-use map[string]any derivative
type. Instances of this type are found within instances of *[DITProfile].

This type merely exists for the sake of client optimization and convenience,
and does not extend from any logic in the I-D series.

Instances of this type are not thread-safe.
*/
type ProfileSettings map[string]any

func newProfileSettings() *ProfileSettings {
	s := make(ProfileSettings, 0)
	return &s
}

/*
Value returns the unasserted value associated with key alongside a
Boolean value indicative of a successful key match.

Case is significant in the matching process.
*/
func (r *ProfileSettings) Value(key string) (value any, ok bool) {
	value, ok = (*r)[key]
	return
}

/*
BoolValue returns an asserted Boolean value associated with key alongside
a Boolean value indicative of a successful key match and value assertion.

Case is significant in the matching process.
*/
func (r *ProfileSettings) BoolValue(key string) (value, ok bool) {
	if avalue, aok := (*r)[key]; aok {
		value, ok = avalue.(bool)
	}

	return
}

/*
StringValue returns an asserted string value associated with key alongside
a Boolean value indicative of a successful key match and value assertion.

Case is significant in the matching process.
*/
func (r *ProfileSettings) StringValue(key string) (value string, ok bool) {
	if avalue, aok := (*r)[key]; aok {
		value, ok = avalue.(string)
	}

	return
}

/*
StringSliceValue returns an asserted Boolean value associated with key
alongside a Boolean value indicative of a successful key match and value
assertion.

Case is significant in the matching process.
*/
func (r *ProfileSettings) StringSliceValue(key string) (value []string, ok bool) {
	if avalue, aok := (*r)[key]; aok {
		value, ok = avalue.([]string)
	}

	return
}

/*
Keys returns slices of string keys present within the receiver instance.
*/
func (r *ProfileSettings) Keys() (keys []string) {
	for k := range *r {
		keys = append(keys, k)
	}

	return
}

/*
Len returns the integer length of the receiver instance.
*/
func (r ProfileSettings) Len() int {
	return len(r)
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *ProfileSettings) IsZero() bool {
	return r == nil
}

/*
Set associates key with value, and assigns value to the receiver instance.
*/
func (r *ProfileSettings) Set(key string, value any) {
	if !r.IsZero() {
		(*r)[key] = value
	}
}

func (r *ProfileSettings) Delete(key string) {
	delete(*r, key)
}

var myDedicatedProfile, myCombinedProfile *DITProfile

func init() {
	SupportedModels = append(SupportedModels, TwoDimensional)
	SupportedModels = append(SupportedModels, ThreeDimensional)

	myDedicatedProfile = &DITProfile{
		R_Settings: newProfileSettings(),
	}

	myDedicatedProfile.SetModel(ThreeDimensional)
	myDedicatedProfile.SetRegistrationBase("ou=Registrations,o=rA")
	myDedicatedProfile.SetRegistrantBase("ou=Registrants,o=rA")

	myCombinedProfile = &DITProfile{
		R_Settings: newProfileSettings(),
	}

	myCombinedProfile.SetModel(ThreeDimensional)
	myCombinedProfile.SetRegistrationBase("ou=Registrations,o=rA")
	myCombinedProfile.SetRegistrantBase("ou=Registrations,o=rA")
}
