package radir

/*
def.go contains a manifest of all RASCHEMA I-D defined attribute types,
object classes and name forms, which includes names and numeric OIDs.
*/

/*
Numeric OID constants represent official OID allocations which are defined
as follows:

  - The root OID prefix is defined in [Section 1.6 of the RADIR I-D]
  - The schema branches for Attribute Types, Object Classes and Name Forms are defined within [Sections 2.3], [2.5] and [2.7 of the RASCHEMA I-D] respectively
  - The two and three dimensional directory models are defined within [Sections 3.1.2] and [3.1.3 of the RADIT I-D] respectively

[Section 1.6 of the RADIR I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-roadmap#section-1.6
[Sections 3.1.2]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.1.2
[3.1.3 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.1.3
[Sections 2.3]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3
[2.5]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5
[2.7 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.7
*/
const (
	OIDPrefix               = `1.3.6.1.4.1.56521.101`
	IRIPrefix               = `/ISO/Identified-Organization/6/1/4/1/56521/101`
	AttributeTypesOIDPrefix = OIDPrefix + `.2.3`
	ObjectClassesOIDPrefix  = OIDPrefix + `.2.5`
	NameFormsOIDPrefix      = OIDPrefix + `.2.7`
	TwoDimensional          = OIDPrefix + `.3.1.2`
	ThreeDimensional        = OIDPrefix + `.3.1.3`
	ASN1Prefix              = `{iso identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 oid-directory(101)}`
)

/*
RegistrationAttributeTypes contains all Attribute Types defined within
[Section 2.3 of the RASCHEMA I-D] that are *[Registration] focused in
nature. This is by no means a complete manifest of all Attribute Types
which may be assigned to entries of this kind.

For reference, each leaf arc of the specified OIDs correlates to the given
subsection of its origin. For example, "aSN1Notation" is defined within
[Section 2.3.4 of the RASCHEMA I-D].

[Section 2.3 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3
[Section 2.3.4 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.4
*/
var RegistrationAttributeTypes map[string]string = map[string]string{
	"aSN1Notation":               AttributeTypesOIDPrefix + ".4",
	"additionalUnicodeValue":     AttributeTypesOIDPrefix + ".6",
	"c-discloseTo":               AttributeTypesOIDPrefix + ".33",
	"c-maxArc":                   AttributeTypesOIDPrefix + ".31",
	"c-minArc":                   AttributeTypesOIDPrefix + ".28",
	"c-supArc":                   AttributeTypesOIDPrefix + ".22",
	"c-topArc":                   AttributeTypesOIDPrefix + ".24",
	"discloseTo":                 AttributeTypesOIDPrefix + ".32",
	"dotEncoding":                AttributeTypesOIDPrefix + ".103",
	"dotNotation":                AttributeTypesOIDPrefix + ".2",
	"iRI":                        AttributeTypesOIDPrefix + ".3",
	"identifier":                 AttributeTypesOIDPrefix + ".7",
	"isFrozen":                   AttributeTypesOIDPrefix + ".17",
	"isLeafNode":                 AttributeTypesOIDPrefix + ".16",
	"leftArc":                    AttributeTypesOIDPrefix + ".26",
	"longArc":                    AttributeTypesOIDPrefix + ".20",
	"maxArc":                     AttributeTypesOIDPrefix + ".30",
	"minArc":                     AttributeTypesOIDPrefix + ".27",
	"n":                          AttributeTypesOIDPrefix + ".1",
	"nameAndNumberForm":          AttributeTypesOIDPrefix + ".19",
	"registeredUUID":             AttributeTypesOIDPrefix + ".102",
	"registrationClassification": AttributeTypesOIDPrefix + ".15",
	"registrationCreated":        AttributeTypesOIDPrefix + ".11",
	"registrationInformation":    AttributeTypesOIDPrefix + ".9",
	"registrationModified":       AttributeTypesOIDPrefix + ".12",
	"registrationRange":          AttributeTypesOIDPrefix + ".13",
	"registrationStatus":         AttributeTypesOIDPrefix + ".14",
	"registrationURI":            AttributeTypesOIDPrefix + ".10",
	"rightArc":                   AttributeTypesOIDPrefix + ".29",
	"secondaryIdentifier":        AttributeTypesOIDPrefix + ".8",
	"standardizedNameForm":       AttributeTypesOIDPrefix + ".18",
	"subArc":                     AttributeTypesOIDPrefix + ".25",
	"supArc":                     AttributeTypesOIDPrefix + ".21",
	"topArc":                     AttributeTypesOIDPrefix + ".23",
	"unicodeValue":               AttributeTypesOIDPrefix + ".5",
}

/*
ConfigurationAttributeTypes contains all Attribute Types defined within
[Section 2.3 of the RASCHEMA I-D] that are *[DUAConfig] and *[DITProfile]
focused in nature. Generally, values of instances of these types will be
found within the Root DSE, or a separate profile entry to which users are
referred from the Root DSE.

In the case of 'rATTL' and 'c-rATTL', instances of these types may appear
virtually anywhere and would be influenced by the overall configuration
and design of the directory service(s) in question.

For reference, each leaf arc of the specified OIDs correlates to the given
subsection of its origin. For example, "c-rATTL" is defined within [Section
2.3.101 of the RASCHEMA I-D].

[Section 2.3 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3
[Section 2.3.101 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.101
*/
var ConfigurationAttributeTypes map[string]string = map[string]string{
	"c-rATTL":            AttributeTypesOIDPrefix + ".101",
	"rADITProfile":       AttributeTypesOIDPrefix + ".94",
	"rADirectoryModel":   AttributeTypesOIDPrefix + ".97",
	"rARegistrantBase":   AttributeTypesOIDPrefix + ".96",
	"rARegistrationBase": AttributeTypesOIDPrefix + ".95",
	"rAServiceMail":      AttributeTypesOIDPrefix + ".98",
	"rAServiceURI":       AttributeTypesOIDPrefix + ".99",
	"rATTL":              AttributeTypesOIDPrefix + ".100",
}

/*
RegistrantAttributeTypes contains all top-level Attribute Types that
are used purely in cases where the RA DSA's registrant policy is "Dedicated".

For reference, each leaf arc of the specified OIDs correlates to the given
subsection of its origin. For example, "c-currentAuthority" is defined within
[Section 2.3.36 of the RASCHEMA I-D].

[Section 2.3.36 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.36
*/
var RegistrantAttributeTypes map[string]string = map[string]string{
	"c-currentAuthority": AttributeTypesOIDPrefix + ".36",
	"c-firstAuthority":   AttributeTypesOIDPrefix + ".55",
	"c-sponsor":          AttributeTypesOIDPrefix + ".75",
	"currentAuthority":   AttributeTypesOIDPrefix + ".35",
	"firstAuthority":     AttributeTypesOIDPrefix + ".54",
	"registrantID":       AttributeTypesOIDPrefix + ".34",
	"sponsor":            AttributeTypesOIDPrefix + ".74",
}

/*
FirstAuthorityAttributeTypes contains all contact-related Attribute Types which
may be assigned to *[FirstAuthority] entries, whether "Dedicated" or "Combined".

For reference, each leaf arc of the specified OIDs correlates to the given
subsection of its origin. For example, "firstAuthorityCommonName" is defined
in [Section 2.3.58 of the RASCHEMA I-D].

[Section 2.3.58 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.58
*/
var FirstAuthorityAttributeTypes map[string]string = map[string]string{
	"firstAuthorityCommonName":     AttributeTypesOIDPrefix + ".58",
	"firstAuthorityCountryCode":    AttributeTypesOIDPrefix + ".59",
	"firstAuthorityCountryName":    AttributeTypesOIDPrefix + ".60",
	"firstAuthorityEmail":          AttributeTypesOIDPrefix + ".61",
	"firstAuthorityEndTimestamp":   AttributeTypesOIDPrefix + ".57",
	"firstAuthorityFax":            AttributeTypesOIDPrefix + ".62",
	"firstAuthorityLocality":       AttributeTypesOIDPrefix + ".63",
	"firstAuthorityMobile":         AttributeTypesOIDPrefix + ".64",
	"firstAuthorityOrg":            AttributeTypesOIDPrefix + ".65",
	"firstAuthorityPOBox":          AttributeTypesOIDPrefix + ".66",
	"firstAuthorityPostalAddress":  AttributeTypesOIDPrefix + ".67",
	"firstAuthorityPostalCode":     AttributeTypesOIDPrefix + ".68",
	"firstAuthorityStartTimestamp": AttributeTypesOIDPrefix + ".56",
	"firstAuthorityState":          AttributeTypesOIDPrefix + ".69",
	"firstAuthorityStreet":         AttributeTypesOIDPrefix + ".70",
	"firstAuthorityTelephone":      AttributeTypesOIDPrefix + ".71",
	"firstAuthorityTitle":          AttributeTypesOIDPrefix + ".72",
	"firstAuthorityURI":            AttributeTypesOIDPrefix + ".73",
}

/*
FirstAuthorityAltAttributeTypes defines the standard RFC-based first
authority equivalent Attribute Types which may be used as opposed to
those defined throughout [Section 2.3 of the RASCHEMA I-D].

Note that novel types, namely the start and end timestamp values, have no
standard (user-managed) equivalents, and thus the novel types are required.

This variable is exposed merely for informational reasons. Users are
generally not required to interact with this instance.

See [Section 3.2.1.1.1 of the RADIT I-D] for details relating to this
alternative Attribute Type policy.

Users are expected to manage any needed Object Classes with respect to
these alternatives themselves.

[Section 2.3 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3
[Section 3.2.1.1.1 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.2.1.1.1
*/
var FirstAuthorityAltAttributeTypes map[string]string = map[string]string{
	"firstAuthorityCommonName":     "cn",
	"firstAuthorityCountryCode":    "c",
	"firstAuthorityCountryName":    "co",
	"firstAuthorityEmail":          "mail",
	"firstAuthorityEndTimestamp":   "firstAuthorityEndTimestamp", // novel type
	"firstAuthorityFax":            "facsimileTelephoneNumber",
	"firstAuthorityLocality":       "l",
	"firstAuthorityMobile":         "mobile",
	"firstAuthorityOrg":            "o",
	"firstAuthorityPOBox":          "postOfficeBox",
	"firstAuthorityPostalAddress":  "postalAddress",
	"firstAuthorityPostalCode":     "postalCode",
	"firstAuthorityStartTimestamp": "firstAuthorityStartTimestamp", // novel type
	"firstAuthorityState":          "st",
	"firstAuthorityStreet":         "street",
	"firstAuthorityTelephone":      "telephoneNumber",
	"firstAuthorityTitle":          "title",
	"firstAuthorityURI":            "labeledURI",
}

/*
CurrentAuthorityAttributeTypes contains all contact-related Attribute Types which
may be assigned to *[CurrentAuthority] entries, whether "Dedicated" or "Combined".

For reference, each leaf arc of the specified OIDs correlates to the given
subsection of its origin. For example, "currentAuthorityCommonName" is defined
in [Section 2.3.38 of the RASCHEMA I-D].

[Section 2.3.38 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.38
*/
var CurrentAuthorityAttributeTypes map[string]string = map[string]string{
	"currentAuthorityCommonName":     AttributeTypesOIDPrefix + ".38",
	"currentAuthorityCountryCode":    AttributeTypesOIDPrefix + ".39",
	"currentAuthorityCountryName":    AttributeTypesOIDPrefix + ".40",
	"currentAuthorityEmail":          AttributeTypesOIDPrefix + ".41",
	"currentAuthorityFax":            AttributeTypesOIDPrefix + ".42",
	"currentAuthorityLocality":       AttributeTypesOIDPrefix + ".43",
	"currentAuthorityMobile":         AttributeTypesOIDPrefix + ".44",
	"currentAuthorityOrg":            AttributeTypesOIDPrefix + ".45",
	"currentAuthorityPOBox":          AttributeTypesOIDPrefix + ".46",
	"currentAuthorityPostalAddress":  AttributeTypesOIDPrefix + ".47",
	"currentAuthorityPostalCode":     AttributeTypesOIDPrefix + ".48",
	"currentAuthorityStartTimestamp": AttributeTypesOIDPrefix + ".37",
	"currentAuthorityState":          AttributeTypesOIDPrefix + ".49",
	"currentAuthorityStreet":         AttributeTypesOIDPrefix + ".50",
	"currentAuthorityTelephone":      AttributeTypesOIDPrefix + ".51",
	"currentAuthorityTitle":          AttributeTypesOIDPrefix + ".52",
	"currentAuthorityURI":            AttributeTypesOIDPrefix + ".53",
}

/*
CurrentAuthorityAltAttributeTypes defines the standard RFC-based current
authority equivalent Attribute Types which may be used as opposed to those
defined throughout [Section 2.3 of the RASCHEMA I-D].

Note that novel types, namely the start timestamp value, have no standard
(user-managed) equivalents, and thus the novel types are required.

This variable is exposed merely for informational reasons. Users are
generally not required to interact with this instance.

See [Section 3.2.1.1.1 of the RADIT I-D] for details relating to this
alternative Attribute Type policy.

Users are expected to manage any needed Object Classes with respect to
these alternatives themselves.

[Section 2.3 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3
[Section 3.2.1.1.1 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.2.1.1.1
*/
var CurrentAuthorityAltAttributeTypes map[string]string = map[string]string{
	"currentAuthorityCommonName":     "cn",
	"currentAuthorityCountryCode":    "c",
	"currentAuthorityCountryName":    "co",
	"currentAuthorityEmail":          "mail",
	"currentAuthorityFax":            "facsimileTelephoneNumber",
	"currentAuthorityLocality":       "l",
	"currentAuthorityMobile":         "mobile",
	"currentAuthorityOrg":            "o",
	"currentAuthorityPOBox":          "postOfficeBox",
	"currentAuthorityPostalAddress":  "postalAddress",
	"currentAuthorityPostalCode":     "postalCode",
	"currentAuthorityStartTimestamp": "currentAuthorityStartTimestamp", // novel type
	"currentAuthorityState":          "st",
	"currentAuthorityStreet":         "street",
	"currentAuthorityTelephone":      "telephoneNumber",
	"currentAuthorityTitle":          "title",
}

/*
SponsorAttributeTypes contains all contact-related Attribute Types which
may be assigned to *[Sponsor] entries, whether "Dedicated" or "Combined".

For reference, each leaf arc of the specified OIDs correlates to the given
subsection of its origin. For example, "sponsorCommonName" is defined in
[Section 2.3.78 of the RASCHEMA I-D].

[Section 2.3.78 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.78
*/
var SponsorAttributeTypes map[string]string = map[string]string{
	"sponsorCommonName":     AttributeTypesOIDPrefix + ".78",
	"sponsorCountryCode":    AttributeTypesOIDPrefix + ".79",
	"sponsorCountryName":    AttributeTypesOIDPrefix + ".80",
	"sponsorEmail":          AttributeTypesOIDPrefix + ".81",
	"sponsorEndTimestamp":   AttributeTypesOIDPrefix + ".77",
	"sponsorFax":            AttributeTypesOIDPrefix + ".82",
	"sponsorLocality":       AttributeTypesOIDPrefix + ".83",
	"sponsorMobile":         AttributeTypesOIDPrefix + ".84",
	"sponsorOrg":            AttributeTypesOIDPrefix + ".85",
	"sponsorPOBox":          AttributeTypesOIDPrefix + ".86",
	"sponsorPostalAddress":  AttributeTypesOIDPrefix + ".87",
	"sponsorPostalCode":     AttributeTypesOIDPrefix + ".88",
	"sponsorStartTimestamp": AttributeTypesOIDPrefix + ".76",
	"sponsorState":          AttributeTypesOIDPrefix + ".89",
	"sponsorStreet":         AttributeTypesOIDPrefix + ".90",
	"sponsorTelephone":      AttributeTypesOIDPrefix + ".91",
	"sponsorTitle":          AttributeTypesOIDPrefix + ".92",
	"sponsorURI":            AttributeTypesOIDPrefix + ".93",
}

/*
SponsorAltAttributeTypes defines the standard RFC-based sponsor authority
equivalent Attribute Types which may be used as opposed to those defined
throughout [Section 2.3 of the RASCHEMA I-D].

Note that novel types, namely the start and end timestamp values, have no
standard (user-managed) equivalents, and thus the novel types are required.

This variable is exposed merely for informational reasons. Users are
generally not required to interact with this instance.

See [Section 3.2.1.1.1 of the RADIT I-D] for details relating to this
alternative Attribute Type policy.

Users are expected to manage any needed Object Classes with respect to
these alternatives themselves.

[Section 2.3 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3
[Section 3.2.1.1.1 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.2.1.1.1
*/
var SponsorAltAttributeTypes map[string]string = map[string]string{
	"sponsorCommonName":     "cn",
	"sponsorCountryCode":    "c",
	"sponsorCountryName":    "co",
	"sponsorEmail":          "mail",
	"sponsorEndTimestamp":   "sponsorEndTimestamp", // novel type
	"sponsorFax":            "facsimileTelephoneNumber",
	"sponsorLocality":       "l",
	"sponsorMobile":         "mobile",
	"sponsorOrg":            "o",
	"sponsorPOBox":          "postOfficeBox",
	"sponsorPostalAddress":  "postalAddress",
	"sponsorPostalCode":     "postalCode",
	"sponsorStartTimestamp": "sponsorStartTimestamp", // novel type
	"sponsorState":          "st",
	"sponsorStreet":         "street",
	"sponsorTelephone":      "telephoneNumber",
	"sponsorTitle":          "title",
	"sponsorURI":            "labeledURI",
}

/*
ObjectClasses contains all Object Classes defined with [Section 2.5 of the
RASCHEMA I-D].

For reference, each leaf arc of the specified OIDs correlates to the given
subsection of its origin. For example, "arc" is defined in [Section 2.5.3
of the RASCHEMA I-D].

[Section 2.5 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5
[Section 2.5.3 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.3
*/
var ObjectClasses map[string]string = map[string]string{
	"arc":                      ObjectClassesOIDPrefix + ".3",
	"currentAuthorityContext":  ObjectClassesOIDPrefix + ".14",
	"firstAuthorityContext":    ObjectClassesOIDPrefix + ".13",
	"iSORegistration":          ObjectClassesOIDPrefix + ".9",
	"iTUTRegistration":         ObjectClassesOIDPrefix + ".8",
	"jointISOITUTRegistration": ObjectClassesOIDPrefix + ".10",
	"rADUAConfig":              ObjectClassesOIDPrefix + ".17",
	"registrant":               ObjectClassesOIDPrefix + ".16",
	"registration":             ObjectClassesOIDPrefix + ".1",
	"registrationSupplement":   ObjectClassesOIDPrefix + ".12",
	"rootArc":                  ObjectClassesOIDPrefix + ".2",
	"spatialContext":           ObjectClassesOIDPrefix + ".11",
	"sponsorContext":           ObjectClassesOIDPrefix + ".15",
	"x660Context":              ObjectClassesOIDPrefix + ".4",
	"x667Context":              ObjectClassesOIDPrefix + ".5",
	"x680Context":              ObjectClassesOIDPrefix + ".6",
	"x690Context":              ObjectClassesOIDPrefix + ".7",
}

/*
NameForms contains all Name Forms defined with [Section 2.7 of the RASCHEMA I-D].

For reference, each leaf arc of the specified OIDs correlates to the given
subsection of its origin. For example, "dotNotationForm" is defined in [Section
2.7.3 of the RASCHEMA I-D].

[Section 2.7 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.7
[Section 2.7.3 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.7.3
*/
var NameForms map[string]string = map[string]string{
	"dotNotationForm": NameFormsOIDPrefix + ".3",
	"nArcForm":        NameFormsOIDPrefix + ".2",
	"nRootArcForm":    NameFormsOIDPrefix + ".1",
}

func resolveAltType(tag string, typ int, alt bool) string {
	if !alt || !(0 <= typ && typ <= 2) {
		// If alt attr policy is not used
		// or, if the type index is bogus,
		// return the same tag.
		return tag
	}

	Maps := []map[string]string{
		FirstAuthorityAltAttributeTypes,   // 0
		CurrentAuthorityAltAttributeTypes, // 1
		SponsorAltAttributeTypes,          // 2
	}

	for k, v := range Maps[typ] {
		if eq(tag, k) {
			// Alternative found. Return
			// it instead of original.
			return v
		}
	}

	// Just return whatever was originally
	// requested, as there seem to be no
	// alternatives for it.
	return tag
}
