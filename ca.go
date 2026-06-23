package radir

/*
ca.go contains all CurrentAuthority methods and types.
*/

/*
CurrentAuthority describes an initial or previous registration authority.

Instances of this type should not be initialized by the user directly.
Instead, see:

  - *[Registrant.CurrentAuthority] (Dedicated Registrants Policy)
  - *[Registration.X660] to access [X660.CombinedCurrentAuthority] (Combined Registrants Policy)
*/
type CurrentAuthority struct {
	// Primary draft-based attribute types for authorities. These
	// represent the default types/fields that will be used for an
	// authority of this form.
	R_L         string   `ldap:"currentAuthorityLocality"`
	R_O         string   `ldap:"currentAuthorityOrg"`
	R_C         string   `ldap:"currentAuthorityCountryCode"`
	R_CO        string   `ldap:"currentAuthorityCountryName"`
	R_ST        string   `ldap:"currentAuthorityState"`
	R_CN        string   `ldap:"currentAuthorityCommonName"`
	R_Tel       string   `ldap:"currentAuthorityTelephone"`
	R_Fax       string   `ldap:"currentAuthorityFax"`
	R_Title     string   `ldap:"currentAuthorityTitle"`
	R_Email     string   `ldap:"currentAuthorityEmail"`
	R_POBox     string   `ldap:"currentAuthorityPOBox"`
	R_PCode     string   `ldap:"currentAuthorityPostalCode"`
	R_PAddr     string   `ldap:"currentAuthorityPostalAddress"`
	R_Street    string   `ldap:"currentAuthorityStreet"`
	R_Mobile    string   `ldap:"currentAuthorityMobile"`
	R_StartTime string   `ldap:"currentAuthorityStartTimestamp"`
	R_URI       []string `ldap:"currentAuthorityURI"`

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
	// Also note that the 'currentAuthorityContext' AUXILIARY class will
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
func (r *CurrentAuthority) marshal(meth func(any) error) error {
	if r.IsZero() {
		r = new(CurrentAuthority)
	} else if meth == nil {
		return NilMethodErr
	}

	return meth(r)
}

/*
unmarshal returns an instance of map[string][]string bearing the contents
of the receiver.
*/
func (r *CurrentAuthority) unmarshal() map[string][]string {
	m := make(map[string][]string)
	return unmarshalStruct(r, m)
}

func (r *CurrentAuthority) ldif() (l string) {
	if !r.IsZero() {
		l = toLDIF(r)
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *CurrentAuthority) IsZero() bool {
	return r == nil
}

func (r *CurrentAuthority) isEmpty() bool {
	return structEmpty(r)
}

/*
CN returns the common name value assigned to the receiver instance.
*/
func (r *CurrentAuthority) CN() (val string) {
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
func (r *CurrentAuthority) SetCN(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityCommonName`,
		1, r.r_alt_types), r.SetCN, r, args...)
}

/*
CNGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) CNGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityCommonName`, 1, r.r_alt_types))
}

/*
L returns the locality name value assigned to the receiver instance.
*/
func (r *CurrentAuthority) L() (val string) {
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
func (r *CurrentAuthority) SetL(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityLocality`,
		1, r.r_alt_types), r.SetL, r, args...)
}

/*
LGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) LGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityLocality`, 1, r.r_alt_types))
}

/*
O returns the organization name value assigned to the receiver instance.
*/
func (r *CurrentAuthority) O() (val string) {
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
func (r *CurrentAuthority) SetO(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityOrg`,
		1, r.r_alt_types), r.SetO, r, args...)
}

/*
OGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) OGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityOrg`, 1, r.r_alt_types))
}

/*
C returns the 2-letter country code value assigned to the receiver instance.
*/
func (r *CurrentAuthority) C() (val string) {
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
func (r *CurrentAuthority) SetC(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityCountryCode`,
		1, r.r_alt_types), r.SetC, r, args...)
}

/*
CGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) CGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityCountryCode`, 1, r.r_alt_types))
}

/*
CO returns the so-called "friendly country name" value assigned to the receiver instance.
*/
func (r *CurrentAuthority) CO() (val string) {
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
func (r *CurrentAuthority) SetCO(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityCountryName`,
		1, r.r_alt_types), r.SetCO, r, args...)
}

/*
COGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) COGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityCountryName`, 1, r.r_alt_types))
}

/*
ST returns the state or province name value assigned to the receiver instance.
*/
func (r *CurrentAuthority) ST() (val string) {
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
func (r *CurrentAuthority) SetST(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityState`,
		1, r.r_alt_types), r.SetST, r, args...)
}

/*
STGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) STGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityState`, 1, r.r_alt_types))
}

/*
Tel returns the telephone number value assigned to the receiver instance.
*/
func (r *CurrentAuthority) Tel() (val string) {
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
func (r *CurrentAuthority) SetTel(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityTelephone`,
		1, r.r_alt_types), r.SetTel, r, args...)
}

/*
TelGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) TelGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityTelephone`, 1, r.r_alt_types))
}

/*
Fax returns the facsimile telephone number value assigned to the receiver instance.
*/
func (r *CurrentAuthority) Fax() (val string) {
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
func (r *CurrentAuthority) SetFax(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityFax`,
		1, r.r_alt_types), r.SetFax, r, args...)
}

/*
FaxGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) FaxGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityFax`, 1, r.r_alt_types))
}

/*
Title returns the title value assigned to the receiver instance.
*/
func (r *CurrentAuthority) Title() (val string) {
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
func (r *CurrentAuthority) SetTitle(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityTitle`,
		1, r.r_alt_types), r.SetTitle, r, args...)
}

/*
TitleGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) TitleGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityTitle`, 1, r.r_alt_types))
}

/*
Email returns the email address value assigned to the receiver instance.
*/
func (r *CurrentAuthority) Email() (val string) {
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
func (r *CurrentAuthority) SetEmail(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityEmail`,
		1, r.r_alt_types), r.SetEmail, r, args...)
}

/*
EmailGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) EmailGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityEmail`, 1, r.r_alt_types))
}

/*
POBox returns the postal office box value assigned to the receiver instance.
*/
func (r *CurrentAuthority) POBox() (val string) {
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
func (r *CurrentAuthority) SetPOBox(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityPOBox`,
		1, r.r_alt_types), r.SetPOBox, r, args...)
}

/*
POBoxGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) POBoxGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityPOBox`, 1, r.r_alt_types))
}

/*
PostalAddress returns the postal address value assigned to the receiver instance.
*/
func (r *CurrentAuthority) PostalAddress() (val string) {
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
func (r *CurrentAuthority) SetPostalAddress(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityPostalAddress`,
		1, r.r_alt_types), r.SetPostalAddress, r, args...)
}

/*
PostalAddressGetFunc processes the underlying field value(s) through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *CurrentAuthority) PostalAddressGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityPostalAddress`, 1, r.r_alt_types))
}

/*
PostalCode returns the postal code value assigned to the receiver instance.
*/
func (r *CurrentAuthority) PostalCode() (val string) {
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
func (r *CurrentAuthority) SetPostalCode(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityPostalCode`,
		1, r.r_alt_types), r.SetPostalCode, r, args...)
}

/*
PostalCodeGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) PostalCodeGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityPostalCode`, 1, r.r_alt_types))
}

/*
Mobile returns the mobile telephone number value assigned to the receiver instance.
*/
func (r *CurrentAuthority) Mobile() (val string) {
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
func (r *CurrentAuthority) SetMobile(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityMobile`,
		1, r.r_alt_types), r.SetMobile, r, args...)
}

/*
MobileGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) MobileGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityMobile`, 1, r.r_alt_types))
}

/*
Street returns the street value assigned to the receiver instance.
*/
func (r *CurrentAuthority) Street() (val string) {
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
func (r *CurrentAuthority) SetStreet(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityStreet`,
		1, r.r_alt_types), r.SetStreet, r, args...)
}

/*
StreetGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) StreetGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityStreet`, 1, r.r_alt_types))
}

/*
URI returns slices of string URIs assigned to the receiver instance.
*/
func (r *CurrentAuthority) URI() (val []string) {
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
func (r *CurrentAuthority) SetURI(args ...any) error {
	return writeFieldByTag(resolveAltType(`currentAuthorityURI`,
		1, r.r_alt_types), r.SetURI, r, args...)
}

/*
URIGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) URIGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc,
		resolveAltType(`currentAuthorityURI`, 1, r.r_alt_types))
}

/*
StartTime returns the string-based generalized time value that reflects
the time at which the receiver was (or will be) officiated.
*/
func (r *CurrentAuthority) StartTime() (when string) {
	if !r.IsZero() {
		when = r.R_StartTime
	}

	return
}

/*
SetStartTime assigns the string input value to the receiver instance.
*/
func (r *CurrentAuthority) SetStartTime(args ...any) error {
	return writeFieldByTag(`currentAuthorityStartTimestamp`, r.SetStartTime, r, args...)
}

/*
StartTimeGetFunc processes the underlying field value(s) through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *CurrentAuthority) StartTimeGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `currentAuthorityStartTimestamp`)
}

/*
Auxiliary returns the static string value "[currentAuthorityContext]" as a
convenient means of determining the AUXILIARY class associated with an
instance of this type.

[currentAuthorityContext]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.13
*/
func (r *CurrentAuthority) Auxiliary() string {
	return `currentAuthorityContext`
}
