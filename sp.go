package radir

/*
sp.go contains all Sponsor methods and types.
*/

/*
Sponsor describes an initial or previous registration authority.

Instances of this type should not be initialized by the user directly.
Instead, see:

  - *[Registrant.Sponsor] (Dedicated Registrants Policy)
  - *[Registration.X660] to access [X660.CombinedSponsor] (Combined Registrants Policy)
*/
type Sponsor struct {
	// Primary draft-based attribute types for authorities. These
	// represent the default types/fields that will be used for an
	// authority of this form.
	R_L         string   `ldap:"sponsorLocality"`
	R_O         string   `ldap:"sponsorOrg"`
	R_C         string   `ldap:"sponsorCountryCode"`
	R_CO        string   `ldap:"sponsorCountryName"`
	R_ST        string   `ldap:"sponsorState"`
	R_CN        string   `ldap:"sponsorCommonName"`
	R_Tel       string   `ldap:"sponsorTelephone"`
	R_Fax       string   `ldap:"sponsorFax"`
	R_Title     string   `ldap:"sponsorTitle"`
	R_Email     string   `ldap:"sponsorEmail"`
	R_POBox     string   `ldap:"sponsorPOBox"`
	R_PCode     string   `ldap:"sponsorPostalCode"`
	R_PAddr     string   `ldap:"sponsorPostalAddress"`
	R_Street    string   `ldap:"sponsorStreet"`
	R_Mobile    string   `ldap:"sponsorMobile"`
	R_StartTime string   `ldap:"sponsorStartTimestamp"`
	R_EndTime   string   `ldap:"sponsorEndTimestamp"`
	R_URI       []string `ldap:"sponsorURI"`

	// Alternative RFC-based attribute types for authorities. See Section
	// 3.2.1.1.1 of the RADIT I-D for strategy details and caveats.
	//
	// By utilizing this strategy, the users in question are expected to
	// manage any custom object class chain elements, such as the 'person',
	// 'organizationalRole', etc. This package will not assist in this task
	// but will not stand in your way, either.
	//
	// Note that is is possible for these standard types to replace all of
	// the above *EXCEPT* for start and end time, as there is no standard
	// user-managed timestamp type of this nature.
	//
	// Also note that the 'sponsorContext' AUXILIARY class will
	// still be used for entries of this kind, regardless of attribute
	// content strategy.
	R_L_alt      string   `ldap:"l"`                        // RFC 4519 § 2.16
	R_O_alt      string   `ldap:"o"`                        // RFC 4519 § 2.19
	R_C_alt      string   `ldap:"c"`                        // RFC 4519 § 2.2
	R_CO_alt     string   `ldap:"co"`                       // RFC 4524 § 2.4
	R_ST_alt     string   `ldap:"st"`                       // RFC 4519 § 2.33
	R_CN_alt     string   `ldap:"cn"`                       // RFC 4519 § 2.3
	R_Tel_alt    string   `ldap:"telephoneNumber"`          // RFC 4519 § 2.35
	R_Fax_alt    string   `ldap:"facsimileTelephoneNumber"` // RFC 4519 § 2.10
	R_Title_alt  string   `ldap:"title"`                    // RFC 4519 § 2.38
	R_Email_alt  string   `ldap:"mail"`                     // RFC 4524 § 2.16
	R_POBox_alt  string   `ldap:"postOfficeBox"`            // RFC 4519 § 2.25
	R_PCode_alt  string   `ldap:"postalCode"`               // RFC 4519 § 2.24
	R_PAddr_alt  string   `ldap:"postalAddress"`            // RFC 4519 § 2.23
	R_Street_alt string   `ldap:"street"`                   // RFC 4519 § 2.34
	R_Mobile_alt string   `ldap:"mobile"`                   // RFC 4524 § 2.18
	R_URI_alt    []string `ldap:"labeledURI"`               // RFC 2079 § 2

	r_alt_types bool
}

/*
marshal returns an error following an attempt to execute the input meth
signature upon the receiver instance. Generally, this method need not
be called directly by the end-user for instances of this type.
*/
func (r *Sponsor) marshal(meth func(any) error) error {
	if r.IsZero() {
		r = new(Sponsor)
	} else if meth == nil {
		return NilMethodErr
	}

	return meth(r)
}

/*
unmarshal returns an instance of map[string][]string bearing the contents
of the receiver.
*/
func (r *Sponsor) unmarshal() map[string][]string {
	m := make(map[string][]string)
	return unmarshalStruct(r, m)
}

func (r *Sponsor) ldif() (l string) {
	if !r.IsZero() {
		l = toLDIF(r)
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *Sponsor) IsZero() bool {
	return r == nil
}

func (r *Sponsor) isEmpty() bool {
	return structEmpty(r)
}

/*
CN returns the common name value assigned to the receiver instance.
*/
func (r *Sponsor) CN() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_CN
		} else {
			val = r.R_CN_alt
		}
	}

	return
}

/*
SetCN assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetCN(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorCommonName`,
		2, r.r_alt_types), r.SetCN, r, args...)
}

/*
CNGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) CNGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorCommonName`, 2, r.r_alt_types))
}

/*
L returns the locality name value assigned to the receiver instance.
*/
func (r *Sponsor) L() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_L
		} else {
			val = r.R_L_alt
		}
	}

	return
}

/*
SetL assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetL(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorLocality`,
		2, r.r_alt_types), r.SetL, r, args...)
}

/*
LGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) LGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorLocality`, 2, r.r_alt_types))
}

/*
O returns the organization name value assigned to the receiver instance.
*/
func (r *Sponsor) O() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_O
		} else {
			val = r.R_O_alt
		}
	}

	return
}

/*
SetO assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetO(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorOrg`,
		2, r.r_alt_types), r.SetO, r, args...)
}

/*
OGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) OGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorOrg`, 2, r.r_alt_types))
}

/*
C returns the 2-letter country code value assigned to the receiver instance.
*/
func (r *Sponsor) C() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_C
		} else {
			val = r.R_C_alt
		}
	}

	return
}

/*
SetC assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetC(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorCountryCode`,
		2, r.r_alt_types), r.SetC, r, args...)
}

/*
CGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) CGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorCountryCode`, 2, r.r_alt_types))
}

/*
CO returns the so-called "friendly country name" value assigned to the receiver instance.
*/
func (r *Sponsor) CO() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_CO
		} else {
			val = r.R_CO_alt
		}
	}

	return
}

/*
SetCO assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetCO(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorCountryName`,
		2, r.r_alt_types), r.SetCO, r, args...)
}

/*
COGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) COGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorCountryName`, 2, r.r_alt_types))
}

/*
ST returns the state or province name value assigned to the receiver instance.
*/
func (r *Sponsor) ST() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_ST
		} else {
			val = r.R_ST_alt
		}
	}

	return
}

/*
SetST assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetST(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorState`,
		2, r.r_alt_types), r.SetST, r, args...)
}

/*
STGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) STGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorState`, 2, r.r_alt_types))
}

/*
Tel returns the telephone number value assigned to the receiver instance.
*/
func (r *Sponsor) Tel() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_Tel
		} else {
			val = r.R_Tel_alt
		}
	}

	return
}

/*
SetTel assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetTel(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorTelephone`,
		2, r.r_alt_types), r.SetTel, r, args...)
}

/*
TelGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) TelGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorTelephone`, 2, r.r_alt_types))
}

/*
Fax returns the facsimile telephone number value assigned to the receiver instance.
*/
func (r *Sponsor) Fax() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_Fax
		} else {
			val = r.R_Fax_alt
		}
	}

	return
}

/*
SetFax assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetFax(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorFax`,
		2, r.r_alt_types), r.SetFax, r, args...)
}

/*
FaxGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) FaxGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorFax`, 2, r.r_alt_types))
}

/*
Title returns the title value assigned to the receiver instance.
*/
func (r *Sponsor) Title() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_Title
		} else {
			val = r.R_Title_alt
		}
	}

	return
}

/*
SetTitle assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetTitle(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorTitle`,
		2, r.r_alt_types), r.SetTitle, r, args...)
}

/*
TitleGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) TitleGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorTitle`, 2, r.r_alt_types))
}

/*
Email returns the email address value assigned to the receiver instance.
*/
func (r *Sponsor) Email() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_Email
		} else {
			val = r.R_Email_alt
		}
	}

	return
}

/*
SetEmail assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetEmail(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorEmail`,
		2, r.r_alt_types), r.SetEmail, r, args...)
}

/*
EmailGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) EmailGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorEmail`, 2, r.r_alt_types))
}

/*
POBox returns the postal office box value assigned to the receiver instance.
*/
func (r *Sponsor) POBox() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_POBox
		} else {
			val = r.R_POBox_alt
		}
	}

	return
}

/*
SetPOBox assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetPOBox(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorPOBox`,
		2, r.r_alt_types), r.SetPOBox, r, args...)
}

/*
POBoxGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) POBoxGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorPOBox`, 2, r.r_alt_types))
}

/*
PostalAddress returns the postal address value assigned to the receiver instance.
*/
func (r *Sponsor) PostalAddress() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_PAddr
		} else {
			val = r.R_PAddr_alt
		}
	}

	return
}

/*
SetPostalAddress assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetPostalAddress(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorPostalAddress`,
		2, r.r_alt_types), r.SetPostalAddress, r, args...)
}

/*
PostalAddressGetFunc processes the underlying field value(s) through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Sponsor) PostalAddressGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorPostalAddress`, 2, r.r_alt_types))
}

/*
PostalCode returns the postal code value assigned to the receiver instance.
*/
func (r *Sponsor) PostalCode() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_PCode
		} else {
			val = r.R_PCode_alt
		}
	}

	return
}

/*
SetPostalCode assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetPostalCode(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorPostalCode`,
		2, r.r_alt_types), r.SetPostalCode, r, args...)
}

/*
PostalCodeGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) PostalCodeGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorPostalCode`, 2, r.r_alt_types))
}

/*
Mobile returns the mobile telephone number value assigned to the receiver instance.
*/
func (r *Sponsor) Mobile() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_Mobile
		} else {
			val = r.R_Mobile_alt
		}
	}

	return
}

/*
SetMobile assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetMobile(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorMobile`,
		2, r.r_alt_types), r.SetMobile, r, args...)
}

/*
MobileGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) MobileGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorMobile`, 2, r.r_alt_types))
}

/*
Street returns the street value assigned to the receiver instance.
*/
func (r *Sponsor) Street() (val string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_Street
		} else {
			val = r.R_Street_alt
		}
	}

	return
}

/*
SetStreet assigns the provided string value to the receiver instance.
*/
func (r *Sponsor) SetStreet(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorStreet`,
		2, r.r_alt_types), r.SetStreet, r, args...)
}

/*
StreetGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) StreetGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorStreet`, 2, r.r_alt_types))
}

/*
URI returns slices of string URIs assigned to the receiver instance.
*/
func (r *Sponsor) URI() (val []string) {
	if !r.IsZero() {
		if !r.r_alt_types {
			val = r.R_URI
		} else {
			val = r.R_URI_alt
		}
	}

	return
}

/*
SetURI appends one or more string slice values to the receiver instance.
Note that if a slice is passed as X, the destination value will be clobbered.
*/
func (r *Sponsor) SetURI(args ...any) error {
	return writeFieldByTag(resolveAltType(`sponsorURI`,
		2, r.r_alt_types), r.SetURI, r, args...)
}

/*
URIGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) URIGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`sponsorURI`, 2, r.r_alt_types))
}

/*
StartTime returns the string-based generalized time value that reflects
the time at which the receiver was (or will be) officiated.
*/
func (r *Sponsor) StartTime() (when string) {
	if !r.IsZero() {
		when = r.R_StartTime
	}

	return
}

/*
SetStartTime assigns the string input value to the receiver instance.
*/
func (r *Sponsor) SetStartTime(args ...any) error {
	return writeFieldByTag(`sponsorStartTimestamp`, r.SetStartTime, r, args...)
}

/*
StartTimeGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) StartTimeGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `sponsorStartTimestamp`)
}

/*
EndTime returns the string-based generalized time value that reflects the
time at which the receiver was (or will be) terminated.
*/
func (r *Sponsor) EndTime() (when string) {
	if !r.IsZero() {
		when = r.R_EndTime
	}

	return
}

/*
SetEndTime appends one or more string slice value to the receiver instance.
*/
func (r *Sponsor) SetEndTime(args ...any) error {
	return writeFieldByTag(`sponsorEndTimestamp`, r.SetEndTime, r, args...)
}

/*
EndTimeGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Sponsor) EndTimeGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `sponsorEndTimestamp`)
}

/*
Auxiliary returns the static string value "[sponsorContext]" as a
convenient means of determining the AUXILIARY class associated with an
instance of this type.

[sponsorContext]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.13
*/
func (r *Sponsor) Auxiliary() string {
	return `sponsorContext`
}
