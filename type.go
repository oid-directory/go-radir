package radir

import (
	"github.com/JesseCoretta/go-objectid"
)

type (
	NumberForm        objectid.NumberForm        // RASCHEMA § 2.3.1
	DotNotation       objectid.DotNotation       // RASCHEMA § 2.3.2
	ASN1Notation      objectid.ASN1Notation      // RASCHEMA § 2.3.4
	NameAndNumberForm objectid.NameAndNumberForm // RASCHEMA § 2.3.19
	Dimension         uint8			     // RADIT § 1.4, 3.1.2, 3.1.3
)

const (
	_ Dimension = iota + 1
	Two	// RADIT § 1.4, 3.1.2
	Three   // RADIT § 1.4, 3.1.3
)

/*
RootArc implements prefabrications of [RASCHEMA § 2.5.2].

[RASCHEMA § 2.5.2]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.2
*/
type RootArc uint8

const (
	ITUTRoot  RootArc = iota // ITU-T Root Arc (0)
	ISORoot                  // ISO Root Arc (1)
	JointRoot                // Joint-ISO-ITU-T Root Arc (2)
)

/*
PrefixOID is a constant definition that mirrors the prefix of all numeric OIDs
defined throughout the OID Directory ID series.  The numeric OID is defined in
[RADIR § 1.6].

The ASN.1 Notation for this definition is as follows:

  {iso identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 oid-directory(101)}

[RADIR § 1.6]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-roadmap#section-1.6
*/
const PrefixOID string = `1.3.6.1.4.1.56521.101`

/*                                                                      
GetOrSetFunc is a first class (closure) "getter or setter" function. This
allows complete control over the creation or interrogation of a value 
assigned to [Registration] or [Registrant] type instances.                  
                                                                        
All Set<*> and <*>GetFunc methods extended by [Registration] or [Registrant]
type instances allow the [GetOrSetFunc] type to be passed at runtime.     
                                                                        
If no [GetOrSetFunc] is passed to a Set<*> method, the value is written   
as-is (type assertion permitting) with no special processing.           
                                                                        
In the context of Set<*> executions, the first input argument will be   
the value to be written to the appropriate struct field. The second     
input argument will be the (non-nil) POINTER receiver instance that     
contains the target field(s) (i.e.: the object to which something is    
being written).                                                         
                                                                        
A Set<*> function that interacts with a struct field of type []string   
can allow append operations (with an individual input value), as well   
as so-called "clobber" operations (with a slice ([]) input value) in    
which any values already present are overwritten (clobbered).  If an    
append operation is needed for multiple values, and clobbering is NOT   
desired, the submission must be done in an iterative (looped) manner.   
You have been warned.                                                   
                                                                        
In the context of Set<*> return values, the first return value will be  
the (processed) value to be written. The second return value, an error, 
will contain error information in the event of any encountered issues.  
                                                                        
In the context of <*>GetFunc executions, the first input argument will  
be the struct field value relevant to the executing method. This will   
produce the value being "gotten" within functions/methods that conform  
to the signature of this type. The second input argument will be the    
non-nil [Registration] or [Registrant] instance being interrogated, which   
may or may not be a POINTER instance.                                   
                                                                        
In the context of <*>GetFunc return values, the first return value will 
be the (processed) value being read. It will manifest as an instance of 
'any', and will require manual type assertion. The second return value, 
an error, will contain error information in the event of any encountered
issues.                                                                 
*/
type GetOrSetFunc func(any, any) (any, error)                           

/*
ITUT implements [RASCHEMA § 2.5.8].

	( 1.3.6.1.4.1.56521.101.2.5.8
	    NAME 'iTUTRegistration'
	    DESC 'X.660, cl. A.2: ITU-T'
	    SUP ( x660Context $
	          x680Context $
	          x690Context )
	    AUXILIARY )

[RASCHEMA § 2.5.8]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.8
*/
type ITUT struct {
	Common RegistrationCommon
	Extra  Supplemental		// RASCHEMA § 2.5.12

	X660	// RASCHEMA § 2.5.4
	X680	// RASCHEMA § 2.5.6
	X690	// RASCHEMA § 2.5.7
	Spatial	// RASCHEMA § 2.5.11
}

/*
ISO implements [RASCHEMA § 2.5.9].

	( 1.3.6.1.4.1.56521.101.2.5.9
	    NAME 'iSORegistration'
	    DESC 'X.660, cl. A.2: ISO'
	    SUP ( x660Context $
	          x680Context $
	          x690Context )
	    AUXILIARY )

[RASCHEMA § 2.5.9]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.9
*/
type ISO struct {
	Common RegistrationCommon
	Extra  Supplemental             // RASCHEMA § 2.5.12

	X660	// RASCHEMA § 2.5.4
	X680	// RASCHEMA § 2.5.6
	X690	// RASCHEMA § 2.5.7
	Spatial	// RASCHEMA § 2.5.11
}

/*
JointISOITUT implements [RASCHEMA § 2.5.10].

	( 1.3.6.1.4.1.56521.101.2.5.10
	    NAME 'jointISOITUTRegistration'
	    DESC 'X.660, cl. A.2: Joint ISO/ITU-T Administration'
	    SUP ( x660Context $
	          x680Context $
	          x690Context )
	    AUXILIARY
	    MAY longArc )

[RASCHEMA § 2.5.10]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.10
*/
type JointISOITUT struct {
	Common RegistrationCommon
	Extra  Supplemental		 // RASCHEMA § 2.5.12
	R_LArc []string `ldap:"longArc"` // RASCHEMA § 2.3.20

	X660	// RASCHEMA § 2.5.4
	X667	// RASCHEMA § 2.5.5
	X680	// RASCHEMA § 2.5.6
	X690	// RASCHEMA § 2.5.7
	Spatial	// RASCHEMA § 2.5.11
}

/*
Registration is an interface type qualified through instances of any of the following types:

  - [ITUT]
  - [ISO]
  - [JointISOITUT]
  - [RootArc]
*/
type Registration interface {
	Identifier() string			  // RASCHEMA § 2.3.7
	DotNotation() DotNotation		  // RASCHEMA § 2.3.2
	ASN1Notation() ASN1Notation		  // RASCHEMA § 2.3.4
	CreateTime() GeneralizedTime		  // RASCHEMA § 2.3.11
	ModifyTime() []GeneralizedTime		  // RASCHEMA § 2.3.12
	N() NumberForm		  		  // RASCHEMA § 2.3.1
	NameAndNumberForm() NameAndNumberForm	  // RASCHEMA § 2.3.19
	SetIdentifier(any, ...GetOrSetFunc) error
	Type() string
}

/*
RegistrationCommon implements struct fields common to all registrations through a
terse implementation of [RASCHEMA § 2.5.1].

	( 1.3.6.1.4.1.56521.101.2.5.1
	    NAME 'registration'
	    DESC 'Abstract OID arc class'
	    SUP top ABSTRACT
	    MUST n
	    MAY ( description $ seeAlso $ rATTL ) )

[RASCHEMA § 2.5.1]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.1
*/
type RegistrationCommon struct {
	R_DN      string   `ldap:"dn"`		// RFC 4511 § 4.1.3
	R_N       string   `ldap:"n"`		// RASCHEMA § 2.3.1
	R_Desc    string   `ldap:"description"` // RFC 4519 § 2.5
	R_RATTL   string   `ldap:"rATTL"`	// RASCHEMA § 2.3.100
	R_OC      []string `ldap:"objectClass"`	// RFC 4512 § 3.3
	R_SeeAlso []string `ldap:"seeAlso"`	// RFC 4519 § 2.30

	RC_RATTL  string   `ldap:"c-rATTL"`	// RASCHEMA § 2.3.101
}

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

[RASCHEMA § 2.5.4]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.4
[ITU-T Rec. X.660]: https://www.itu.int/rec/T-REC-X.660                            
*/
type X660 struct {
	R_UVal     []byte   `ldap:"unicodeValue"`		// RASCHEMA § 2.3.5
	R_AddlUVal []string `ldap:"additionalUnicodeValue"`	// RASCHEMA § 2.3.6
	R_SecId    []string `ldap:"secondaryIdentifier"`	// RASCHEMA § 2.3.8
	R_StdNF    []string `ldap:"standardizedNameForm"`	// RASCHEMA § 2.3.18

	R_FirstRA   []string `ldap:"firstAuthority"`		// RASCHEMA § 2.3.54 -- Dedicated Registrant use only [RECOMMENDED]!
	R_CurrentRA []string `ldap:"currentAuthority"`		// RASCHEMA § 2.3.35 -- Dedicated Registrant use only [RECOMMENDED]!
	R_SponsorRA []string `ldap:"sponsor"`			// RASCHEMA § 2.3.74 -- Dedicated Registrant use only [RECOMMENDED]!

	GenericAuthority // Combined Registrant use only [NOT RECOMMENDED]!
	CurrentAuthority // RASCHEMA § 2.3.37-53 -- Combined Registrant use only [NOT RECOMMENDED]!
	FirstAuthority   // RASCHEMA § 2.3.56-73 -- Combined Registrant use only [NOT RECOMMENDED]!
	Sponsor          // RASCHEMA § 2.3.76-93 -- Combined Registrant use only [NOT RECOMMENDED]!
}

/*
X667 implements [RASCHEMA § 2.5.5] and derives various concepts from
[ITU-T Rec. X.667].

	( 1.3.6.1.4.1.56521.101.2.5.5
	    NAME 'x667Context'
	    DESC 'X.667 contextual class'
	    SUP registration AUXILIARY
	    MUST registeredUUID )

[RASCHEMA § 2.5.5]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.5
[ITU-T Rec. X.667]: https://www.itu.int/rec/T-REC-X.667                            
*/
type X667 struct {
	R_UUID []byte `ldap:"registeredUUID"`	 // RASCHEMA § 2.3.102
}

/*
X680 implements [RASCHEMA § 2.5.6] and derives various concepts from
[ITU-T Rec. X.680].

	( 1.3.6.1.4.1.56521.101.2.5.6
	    NAME 'x680Context'
	    DESC 'X.680 contextual class'
	    SUP registration AUXILIARY
	    MAY ( aSN1Notation $
	          dotNotation $
	          identifier $
	          iRI $
	          nameAndNumberForm ) )

[RASCHEMA § 2.5.6]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.6
[ITU-T Rec. X.680]: https://www.itu.int/rec/T-REC-X.680                            
*/
type X680 struct {
	R_ANot string   `ldap:"aSN1Notation"`	   // RASCHEMA § 2.3.4
	R_DNot string   `ldap:"dotNotation"`	   // RASCHEMA § 2.3.2
	R_Id   string   `ldap:"identifier"`	   // RASCHEMA § 2.3.7
	R_NaNF string   `ldap:"nameAndNumberForm"` // RASCHEMA § 2.3.19
	R_IRI  []string `ldap:"iRI"`		   // RASCHEMA § 2.3.3
}

/*
X690 implements [RASCHEMA § 2.5.7] and derives various concepts from
[ITU-T Rec. X.690].

	( 1.3.6.1.4.1.56521.101.2.5.7
	    NAME 'x690Context'
	    DESC 'X.690 contextual class'
	    SUP registration AUXILIARY
	    MAY dotEncoding )

[RASCHEMA § 2.5.7]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.7
[ITU-T Rec. X.690]: https://www.itu.int/rec/T-REC-X.690 
*/
type X690 struct {
	R_DEnc []byte `ldap:"dotEncoding"`	// RASCHEMA § 2.3.103
}

/*
Spatial implements [RASCHEMA § 2.5.11].

	( 1.3.6.1.4.1.56521.101.2.5.11
	    NAME 'spatialContext'
	    DESC 'Logical spatial orientation and association class'
	    SUP registration AUXILIARY
	    MAY ( topArc $
	          supArc $
	          subArc $
	          minArc $
	          maxArc $
	          leftArc $
	          rightArc ) )

[RASCHEMA § 2.5.11]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.11
*/
type Spatial struct {
	R_TopArc   string   `ldap:"topArc"`	// RASCHEMA § 2.3.23
	R_SupArc   string   `ldap:"supArc"`	// RASCHEMA § 2.3.21
	R_MinArc   string   `ldap:"minArc"`	// RASCHEMA § 2.3.27
	R_MaxArc   string   `ldap:"maxArc"`	// RASCHEMA § 2.3.30
	R_LeftArc  string   `ldap:"leftArc"`	// RASCHEMA § 2.3.26
	R_RightArc string   `ldap:"rightArc"`	// RASCHEMA § 2.3.29
	R_SubArc   []string `ldap:"subArc"`	// RASCHEMA § 2.3.25

	RC_TopArc string `ldap:"c-topArc"`	// RASCHEMA § 2.3.24
	RC_SupArc string `ldap:"c-supArc"`	// RASCHEMA § 2.3.22
	RC_MinArc string `ldap:"c-minArc"`	// RASCHEMA § 2.3.28
	RC_MaxArc string `ldap:"c-maxArc"`	// RASCHEMA § 2.3.31
}

/*
Supplemental implements [RASCHEMA § 2.5.12] to make miscellaneous
types available -- whether novel or derived from other standards --
for assignment to registration instances.

	( 1.3.6.1.4.1.56521.101.2.5.12
	    NAME 'registrationSupplement'
	    DESC 'Supplemental registration class'
	    SUP registration AUXILIARY
	    MAY ( discloseTo $ isFrozen $ isLeafNode $
	          registrationClassification $
	          registrationCreated $
	          registrationInformation $
	          registrationModified $
	          registrationRange $
	          registrationStatus $
	          registrationURI ) )

[RASCHEMA § 2.5.12]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.12
*/
type Supplemental struct {
	R_Status string   `ldap:"registrationStatus"`		// RASCHEMA § 2.3.14, RFC 2578 § 2
	R_Class  string   `ldap:"registrationClassification"`	// RASCHEMA § 2.3.15
	R_CTime  string   `ldap:"registrationCreated"`		// RASCHEMA § 2.3.11
	R_Frozen string   `ldap:"isFrozen"`			// RASCHEMA § 2.3.17
	R_Leaf   string   `ldap:"isLeafNode"`			// RASCHEMA § 2.3.16
	R_Range  string   `ldap:"registrationRange"`		// RASCHEMA § 2.3.13
	R_Info   []string `ldap:"registrationInformation"`	// RASCHEMA § 2.3.9
	R_MTime  []string `ldap:"registrationModified"`		// RASCHEMA § 2.3.12, RFC 2578 § 2
	R_URI    []string `ldap:"registrationURI"`		// RASCHEMA § 2.3.10
	R_DsclTo []string `ldap:"discloseTo"`			// RASCHEMA § 2.3.32

	RC_DsclTo []string `ldap:"c-discloseTo"`		// RASCHEMA § 2.3.33
}

/*
Registrant implements the dedicated registrant per [RASCHEMA § 2.5.16].

	( 1.3.6.1.4.1.56521.101.2.5.16
	    NAME 'registrant'
	    DESC 'Generalized auxiliary class for registrant data'
	    SUP top STRUCTURAL
	    MUST registrantID
	    MAY ( description $
	              seeAlso $
	              rATTL ) )

[RASCHEMA § 2.5.16]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.16
*/
type Registrant struct {
	R_DN      string   `ldap:"dn"`		 // RFC 4511 § 4.1.3
	R_Id      string   `ldap:"registrantID"` // RASCHEMA § 2.3.34
	R_Desc    string   `ldap:"description"`  // RFC 4519 § 2.5
	R_RATTL   string   `ldap:"rATTL"`	 // RASCHEMA § 2.3.100
	R_OC      []string `ldap:"objectClass"`  // RFC 4512 § 3.3
	R_SeeAlso []string `ldap:"seeAlso"`	 // RFC 4519 § 2.30

	GenericAuthority
	CurrentAuthority
	FirstAuthority
	Sponsor
}

/*
GenericAuthority implements [RADIT § 3.2.1.1.1].  With respect to use of
the registrant types, such as "sponsorCommonName", in the context of
Dedicated Registrants:

	Adopters of the RA DIT MAY choose to forego use of some or all of
	these [registrant] types in favor of those defined within official
	standards such as RFC4519 and RFC4524.

	For example, the appropriate replacement type for 'sponsorCommonName'
	is 'cn', because 'cn' is its super type.

This struct type implements a basic means of marshaling values assigned
to these well-known types instead of those defined within the [RASCHEMA ID].

[RASCHEMA ID]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema
[RADIT § 3.2.1.1.1]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.2.1.1.1
*/
type GenericAuthority struct {
	R_L      string   `ldap:"l"`				// RFC 4519 § 2.16
	R_O      string   `ldap:"o"`				// RFC 4519 § 2.19
	R_C      string   `ldap:"c"`				// RFC 4519 § 2.2
	R_CO     string   `ldap:"co"`				// RFC 4524 § 2.4
	R_ST     string   `ldap:"st"`				// RFC 4519 § 2.33
	R_CN     string   `ldap:"cn"`				// RFC 4519 § 2.3
	R_Tel    string   `ldap:"telephoneNumber"`		// RFC 4519 § 2.35
	R_Fax    string   `ldap:"facsimileTelephoneNumber"`	// RFC 4519 § 2.10
	R_Title  string   `ldap:"title"`			// RFC 4519 § 2.38
	R_Email  string   `ldap:"mail"`				// RFC 4524 § 2.16
	R_POBox  string   `ldap:"postOfficeBox"`		// RFC 4519 § 2.25
	R_PCode  string   `ldap:"postalCode"`			// RFC 4519 § 2.24
	R_PAddr  string   `ldap:"postalAddress"`		// RFC 4519 § 2.23
	R_Street string   `ldap:"street"`			// RFC 4519 § 2.34
	R_Mobile string   `ldap:"mobile"`			// RFC 4524 § 2.18
	R_URI    []string `ldap:"labeledURI"`			// RFC 2079 Pg. 2
}

/*
Sponsor implements [§ 2.5.15] and [§ 2.3.76] through [§ 2.3.93] of the
[RASCHEMA ID].

	( 1.3.6.1.4.1.56521.101.2.5.15
	    NAME 'sponsorContext'
	    DESC 'Registration sponsoring authority class'
	    SUP top AUXILIARY
	    MAY ( sponsorCommonName $
	          sponsorCountryCode $
	          sponsorCountryName $
	          sponsorEmail $
	          sponsorEndTimestamp $
	          sponsorFax $
	          sponsorLocality $
	          sponsorMobile $
	          sponsorOrg $
	          sponsorPOBox $
	          sponsorPostalAddress $
	          sponsorPostalCode $
	          sponsorStartTimestamp $
	          sponsorState $
	          sponsorStreet $
	          sponsorTelephone $
	          sponsorTitle $
	          sponsorURI ) )

[RASCHEMA ID]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema
[§ 2.5.15]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.15
[§ 2.3.76]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.76
[§ 2.3.93]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.93
*/
type Sponsor struct {
	R_L         string   `ldap:"sponsorLocality"`	    // RASCHEMA § 2.3.83
	R_O         string   `ldap:"sponsorOrg"`	    // RASCHEMA § 2.3.85
	R_C         string   `ldap:"sponsorCountryCode"`    // RASCHEMA § 2.3.79
	R_CO        string   `ldap:"sponsorCountryName"`    // RASCHEMA § 2.3.80
	R_ST        string   `ldap:"sponsorState"`	    // RASCHEMA § 2.3.89
	R_CN        string   `ldap:"sponsorCommonName"`	    // RASCHEMA § 2.3.78
	R_Tel       string   `ldap:"sponsorTelephone"`	    // RASCHEMA § 2.3.91
	R_Fax       string   `ldap:"sponsorFax"`	    // RASCHEMA § 2.3.82
	R_Title     string   `ldap:"sponsorTitle"`	    // RASCHEMA § 2.3.92
	R_Email     string   `ldap:"sponsorEmail"`	    // RASCHEMA § 2.3.81
	R_POBox     string   `ldap:"sponsorPOBox"`	    // RASCHEMA § 2.3.86
	R_PCode     string   `ldap:"sponsorPostalCode"`     // RASCHEMA § 2.3.88
	R_PAddr     string   `ldap:"sponsorPostalAddress"`  // RASCHEMA § 2.3.87
	R_Street    string   `ldap:"sponsorStreet"`	    // RASCHEMA § 2.3.90
	R_Mobile    string   `ldap:"sponsorMobile"`	    // RASCHEMA § 2.3.84
	R_StartTime string   `ldap:"sponsorStartTimestamp"` // RASCHEMA § 2.3.76
	R_EndTime   string   `ldap:"sponsorEndTimestamp"`   // RASCHEMA § 2.3.77
	R_URI       []string `ldap:"sponsorURI"`	    // RASCHEMA § 2.3.93
}

/*
CurrentAuthority implements [§ 2.5.14] and [§ 2.3.37] through [§ 2.3.53] of the
[RASCHEMA ID].

	( 1.3.6.1.4.1.56521.101.2.5.14
	    NAME 'currentAuthorityContext'
	    DESC 'Current registration authority class'
	    SUP top AUXILIARY
	    MAY ( currentAuthorityCommonName $
	          currentAuthorityCountryCode $
	          currentAuthorityCountryName $
	          currentAuthorityEmail $
	          currentAuthorityFax $
	          currentAuthorityLocality $
	          currentAuthorityMobile $
	          currentAuthorityOrg $
	          currentAuthorityPOBox $
	          currentAuthorityPostalAddress $
	          currentAuthorityPostalCode $
	          currentAuthorityStartTimestamp $
	          currentAuthorityState $
	          currentAuthorityStreet $
	          currentAuthorityTelephone $
	          currentAuthorityTitle $
	          currentAuthorityURI ) )

[RASCHEMA ID]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema
[§ 2.5.14]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.14
[§ 2.3.37]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.37
[§ 2.3.53]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.53
*/
type CurrentAuthority struct {
	R_L         string   `ldap:"currentAuthorityLocality"`	     // RASCHEMA § 2.3.43
	R_O         string   `ldap:"currentAuthorityOrg"`	     // RASCHEMA § 2.3.45
	R_C         string   `ldap:"currentAuthorityCountryCode"`    // RASCHEMA § 2.3.39
	R_CO        string   `ldap:"currentAuthorityCountryName"`    // RASCHEMA § 2.3.40
	R_ST        string   `ldap:"currentAuthorityState"`	     // RASCHEMA § 2.3.49
	R_CN        string   `ldap:"currentAuthorityCommonName"`     // RASCHEMA § 2.3.38
	R_Tel       string   `ldap:"currentAuthorityTelephone"`	     // RASCHEMA § 2.3.51
	R_Fax       string   `ldap:"currentAuthorityFax"`	     // RASCHEMA § 2.3.42
	R_Title     string   `ldap:"currentAuthorityTitle"`	     // RASCHEMA § 2.3.52
	R_Email     string   `ldap:"currentAuthorityEmail"`	     // RASCHEMA § 2.3.41
	R_POBox     string   `ldap:"currentAuthorityPOBox"`	     // RASCHEMA § 2.3.46
	R_PCode     string   `ldap:"currentAuthorityPostalCode"`     // RASCHEMA § 2.3.48
	R_PAddr     string   `ldap:"currentAuthorityPostalAddress"`  // RASCHEMA § 2.3.47
	R_Street    string   `ldap:"currentAuthorityStreet"`	     // RASCHEMA § 2.3.50
	R_Mobile    string   `ldap:"currentAuthorityMobile"`	     // RASCHEMA § 2.3.44
	R_StartTime string   `ldap:"currentAuthorityStartTimestamp"` // RASCHEMA § 2.3.37
	R_URI       []string `ldap:"currentAuthorityURI"`	     // RASCHEMA § 2.3.53
}

/*
FirstAuthority implements [§ 2.5.13] and [§ 2.3.56] through [§ 2.3.73] of the
[RASCHEMA ID].

	( 1.3.6.1.4.1.56521.101.2.5.13
	    NAME 'firstAuthorityContext'
	    DESC 'Initial registration authority class'
	    SUP top AUXILIARY
	    MAY ( firstAuthorityCommonName $
	          firstAuthorityCountryCode $
	          firstAuthorityCountryName $
	          firstAuthorityEmail $
	          firstAuthorityEndTimestamp $
	          firstAuthorityFax $
	          firstAuthorityLocality $
	          firstAuthorityMobile $
	          firstAuthorityOrg $
	          firstAuthorityPOBox $
	          firstAuthorityPostalAddress $
	          firstAuthorityPostalCode $
	          firstAuthorityStartTimestamp $
	          firstAuthorityState $
	          firstAuthorityStreet $
	          firstAuthorityTelephone $
	          firstAuthorityTitle $
	          firstAuthorityURI ) )

[RASCHEMA ID]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema
[§ 2.5.13]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.13
[§ 2.3.56]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.56
[§ 2.3.73]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.73
*/
type FirstAuthority struct {
	R_L         string   `ldap:"firstAuthorityLocality"`	   // RASCHEMA § 2.3.63
	R_O         string   `ldap:"firstAuthorityOrg"`	   	   // RASCHEMA § 2.3.65
	R_C         string   `ldap:"firstAuthorityCountryCode"`	   // RASCHEMA § 2.3.59
	R_CO        string   `ldap:"firstAuthorityCountryName"`	   // RASCHEMA § 2.3.60
	R_ST        string   `ldap:"firstAuthorityState"`	   // RASCHEMA § 2.3.69
	R_CN        string   `ldap:"firstAuthorityCommonName"`	   // RASCHEMA § 2.3.58
	R_Tel       string   `ldap:"firstAuthorityTelephone"`	   // RASCHEMA § 2.3.71
	R_Fax       string   `ldap:"firstAuthorityFax"`	   	   // RASCHEMA § 2.3.62
	R_Title     string   `ldap:"firstAuthorityTitle"`	   // RASCHEMA § 2.3.72
	R_Email     string   `ldap:"firstAuthorityEmail"`	   // RASCHEMA § 2.3.61
	R_POBox     string   `ldap:"firstAuthorityPOBox"`	   // RASCHEMA § 2.3.66
	R_PCode     string   `ldap:"firstAuthorityPostalCode"`	   // RASCHEMA § 2.3.68
	R_PAddr     string   `ldap:"firstAuthorityPostalAddress"`  // RASCHEMA § 2.3.67
	R_Street    string   `ldap:"firstAuthorityStreet"`	   // RASCHEMA § 2.3.70
	R_Mobile    string   `ldap:"firstAuthorityMobile"`	   // RASCHEMA § 2.3.64
	R_StartTime string   `ldap:"firstAuthorityStartTimestamp"` // RASCHEMA § 2.3.56
	R_EndTime   string   `ldap:"firstAuthorityEndTimestamp"`   // RASCHEMA § 2.3.57
	R_URI       []string `ldap:"firstAuthorityURI"`	   	   // RASCHEMA § 2.3.73
}

/*
DUAConfig implements [RASCHEMA § 2.5.17].

      ( 1.3.6.1.4.1.56521.101.2.5.17
          NAME 'rADUAConfig'
          DESC 'RA DUA configuration advertisement class'
          SUP top AUXILIARY
          MAY ( description $
                rADITProfile $
                rADirectoryModel $
                rARegistrationBase $
                rARegistrantBase $
                rAServiceMail $
                rAServiceURI $
                rATTL ) )

Instances of this type generally manifest based on the contents of the
Root DSE, thus does not possess a non-zero DN.

If the R_Prof field is of a non-zero length, all other tagged values
are ignored if populated.

Each R_Prof DN slice refers to the [DITProfile] slice of a like index.
Each [DITProfile] slice represents a distinct DIT profile configuration
to be used by capable clients attempting to interact with the associated
RA DIT.

[RASCHEMA § 2.5.17]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.17
*/

type DITProfile struct {
	R_DN    string   `ldap:"dn"`			// RFC 4511 § 4.1.3 
	R_Model string   `ldap:"rADirectoryModel"`	// RASCHEMA § 2.3.97
	R_TTL   string   `ldap:"rATTL"`			// RASCHEMA § 2.3.100
	R_URI   []string `ldap:"rAServiceURI"`		// RASCHEMA § 2.3.99
	R_RegB  []string `ldap:"rARegistrationBase"`	// RASCHEMA § 2.3.95
	R_AthB  []string `ldap:"rARegistrantBase"`	// RASCHEMA § 2.3.96
	R_Mail  []string `ldap:"rAServiceMail"`		// RASCHEMA § 2.3.98
	R_Desc  []string `ldap:"description"`		// RFC 4511 § 4.1.3 
	R_Prof  []string `ldap:"rADITProfile"`		// RASCHEMA § 2.3.94

	RC_TTL string `ldap:"c-rATTL"`	// RASCHEMA § 2.3.101
}

/*
DUAConfig implements an abstraction of [RADUA § 2.2.2.1.1] and [RASCHEMA § 2.5.17].

The most common use case is the advertisement of a single RA DIT configuration
by way of the Root DSE. The following is a basic example configuration:

	  DUAConfig = []DUAProfile{
		DUAProfile{
			R_DN:	 "",	// zero length DN of Root DSE
			R_Model: "1.3.6.1.4.1.56521.101.3.1.3,
			R_RegB: "ou=Registrations,o=rA",
			R_AthB: "ou=Registrants,o=rA",
		},
	  }

For reference, the following is the equivalent (generalized) LDIF entry for the above:

	dn:
	objectClass: rADUAConfig
	rADirectoryModel: 1.3.6.1.4.1.56521.101.3.1.3
	rARegistrationBase: ou=Registrations,o=rA
	rARegistrantBase: ou=Registrants,o=rA

Given the above configuration data, clients capable of auto-configuration in
the terms described within the OID Directory ID series will have the means of
efficient location of, and interaction with, registration content.

If two (2) or more slices reside within an instance of this type, the first
slice serves only to refer to subsequence slices by way of the R_Prof []slice
field:

	  DUAConfig = []DUAProfile{
		DUAProfile{
			// Root DSE "rADITProfile" references
			R_Prof: []string{
				"dc=example,dc=net", // rADUAConfig entry
				"o=rA",              // rADUAConfig entry
			},
		},
		DUAProfile{
			R_DN:    "dc=example,dc=net"
			R_Desc:  "Legacy Registration Authority",
			R_Model: "1.3.6.1.4.1.56521.101.3.1.2,
			R_RegB:  "ou=OIDs,dc=example,dc=net",
			R_AthB:  "ou=Authority,dc=example,dc=net",
		},
		DUAProfile{
			R_DN:	 "o=rA",
			R_Desc:  "Standard Registration Authority",
			R_Model: "1.3.6.1.4.1.56521.101.3.1.3",
			R_RegB:  "ou=Registrations,o=rA",
			R_AthB:  "ou=Registrants,o=rA",
		},
	  }

For reference, the following are the equivalent (generalized) LDIF entries for the above:

	dn:
	objectClass: rADUAConfig
	rADITProfile: dc=example,dc=net
	rADITProfile: o=rA

	dn: dc=example,dc=net
	objectClass: rADUAConfig
	description: Legacy Registration Authority
	rADirectoryModel: 1.3.6.1.4.1.56521.101.3.1.2
	rARegistrationBase: ou=OIDs,dc=example,dc=net
	rARegistrantBase: ou=Authority,dc=example,dc=net

	dn: o=rA
	objectClass: rADUAConfig
	description: Standard Registration Authority
	rADirectoryModel: 1.3.6.1.4.1.56521.101.3.1.3
	rARegistrationBase: ou=Registrations,o=rA
	rARegistrantBase: ou=Registrants,o=rA

Use of multiple profiles as shown is only necessary for cornercases in which different segments
of the OID tree are served by the same directory server, but operate in different models and/or
governing structures.  In nearly all cases, a single 'rADUAConfig' implementation served via the
Root DSE is sufficient and recommended.

[RADUA § 2.2.2.1.1]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radua#section-2.2.2.1.1
[RASCHEMA § 2.5.17]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.17

*/
type DUAConfig []DITProfile

var schema map[string]map[string][]string
var attributeTypes map[string][]string
var objectClasses map[string][]string
var nameForms map[string][]string

func init() {
	attributeTypes = map[string][]string{
		// §	   NAME
		`2.3.1`:   {`n`, `numberForm`},
		`2.3.2`:   {`dotNotation`},
		`2.3.3`:   {`iRI`},
		`2.3.4`:   {`aSN1Notation`},
		`2.3.5`:   {`unicodeValue`},
		`2.3.6`:   {`additionalUnicodeValue`},
		`2.3.7`:   {`identifier`, `nameForm`},
		`2.3.8`:   {`secondaryIdentifier`},
		`2.3.9`:   {`registrationInformation`},
		`2.3.10`:  {`registrationURI`},
		`2.3.11`:  {`registrationCreated`},
		`2.3.12`:  {`registrationModified`},
		`2.3.13`:  {`registrationRange`},
		`2.3.14`:  {`registrationStatus`},
		`2.3.15`:  {`registrationClassification`},
		`2.3.16`:  {`isLeafNode`},
		`2.3.17`:  {`isFrozen`},
		`2.3.18`:  {`standardizedNameForm`},
		`2.3.19`:  {`nameAndNumberForm`},
		`2.3.20`:  {`longArc`},
		`2.3.21`:  {`supArc`},
		`2.3.22`:  {`c-supArc`},
		`2.3.23`:  {`topArc`},
		`2.3.24`:  {`c-topArc`},
		`2.3.25`:  {`subArc`},
		`2.3.26`:  {`leftArc`},
		`2.3.27`:  {`minArc`},
		`2.3.28`:  {`c-minArc`},
		`2.3.29`:  {`rightArc`},
		`2.3.30`:  {`maxArc`},
		`2.3.31`:  {`c-maxArc`},
		`2.3.32`:  {`discloseTo`},
		`2.3.33`:  {`c-discloseTo`},
		`2.3.34`:  {`registrantID`},
		`2.3.35`:  {`currentAuthority`},
		`2.3.36`:  {`c-currentAuthority`},
		`2.3.37`:  {`currentAuthorityStartTimestamp`},
		`2.3.38`:  {`currentAuthorityCommonName`},
		`2.3.39`:  {`currentAuthorityCountryCode`},
		`2.3.40`:  {`currentAuthorityCountryName`},
		`2.3.41`:  {`currentAuthorityEmail`},
		`2.3.42`:  {`currentAuthorityFax`},
		`2.3.43`:  {`currentAuthorityLocality`},
		`2.3.44`:  {`currentAuthorityMobile`},
		`2.3.45`:  {`currentAuthorityOrg`},
		`2.3.46`:  {`currentAuthorityPOBox`},
		`2.3.47`:  {`currentAuthorityPostalAddress`},
		`2.3.48`:  {`currentAuthorityPostalCode`},
		`2.3.49`:  {`currentAuthorityState`},
		`2.3.50`:  {`currentAuthorityStreet`},
		`2.3.51`:  {`currentAuthorityTelephone`},
		`2.3.52`:  {`currentAuthorityTitle`},
		`2.3.53`:  {`currentAuthorityURI`},
		`2.3.54`:  {`firstAuthority`},
		`2.3.55`:  {`c-firstAuthority`},
		`2.3.56`:  {`firstAuthorityStartTimestamp`},
		`2.3.57`:  {`firstAuthorityEndTimestamp`},
		`2.3.58`:  {`firstAuthorityCommonName`},
		`2.3.59`:  {`firstAuthorityCountryCode`},
		`2.3.60`:  {`firstAuthorityCountryName`},
		`2.3.61`:  {`firstAuthorityEmail`},
		`2.3.62`:  {`firstAuthorityFax`},
		`2.3.63`:  {`firstAuthorityLocality`},
		`2.3.64`:  {`firstAuthorityMobile`},
		`2.3.65`:  {`firstAuthorityOrg`},
		`2.3.66`:  {`firstAuthorityPOBox`},
		`2.3.67`:  {`firstAuthorityPostalAddress`},
		`2.3.68`:  {`firstAuthorityPostalCode`},
		`2.3.69`:  {`firstAuthorityState`},
		`2.3.70`:  {`firstAuthorityStreet`},
		`2.3.71`:  {`firstAuthorityTelephone`},
		`2.3.72`:  {`firstAuthorityTitle`},
		`2.3.73`:  {`firstAuthorityURI`},
		`2.3.74`:  {`sponsor`},
		`2.3.75`:  {`c-sponsor`},
		`2.3.76`:  {`sponsorStartTimestamp`},
		`2.3.77`:  {`sponsorEndTimestamp`},
		`2.3.78`:  {`sponsorCommonName`},
		`2.3.79`:  {`sponsorCountryCode`},
		`2.3.80`:  {`sponsorCountryName`},
		`2.3.81`:  {`sponsorEmail`},
		`2.3.82`:  {`sponsorFax`},
		`2.3.83`:  {`sponsorLocality`},
		`2.3.84`:  {`sponsorMobile`},
		`2.3.85`:  {`sponsorOrg`},
		`2.3.86`:  {`sponsorPOBox`},
		`2.3.87`:  {`sponsorPostalAddress`},
		`2.3.88`:  {`sponsorPostalCode`},
		`2.3.89`:  {`sponsorState`},
		`2.3.90`:  {`sponsorStreet`},
		`2.3.91`:  {`sponsorTelephone`},
		`2.3.92`:  {`sponsorTitle`},
		`2.3.93`:  {`sponsorURI`},
		`2.3.94`:  {`rADITProfile`},
		`2.3.95`:  {`rARegistrationBase`},
		`2.3.96`:  {`rARegistrantBase`},
		`2.3.97`:  {`rADirectoryModel`},
		`2.3.98`:  {`rAServiceMail`},
		`2.3.99`:  {`rAServiceURI`},
		`2.3.100`: {`rATTL`},
		`2.3.101`: {`c-rATTL`},
		`2.3.102`: {`registeredUUID`},
		`2.3.103`: {`dotEncoding`},
	}

	objectClasses = map[string][]string{
		// §	   NAME
		`2.5.1`:  {`registration`},
		`2.5.2`:  {`rootArc`},
		`2.5.3`:  {`arc`},
		`2.5.4`:  {`x660Context`},
		`2.5.5`:  {`x667Context`},
		`2.5.6`:  {`x680Context`},
		`2.5.7`:  {`x690Context`},
		`2.5.8`:  {`iTUTRegistration`},
		`2.5.9`:  {`iSORegistration`},
		`2.5.10`: {`jointISOITUTRegistration`},
		`2.5.11`: {`spatialContext`},
		`2.5.12`: {`registrationSupplement`},
		`2.5.13`: {`firstAuthorityContext`},
		`2.5.14`: {`currentAuthorityContext`},
		`2.5.15`: {`sponsorContext`},
		`2.5.16`: {`registrant`},
		`2.5.17`: {`rADUAConfig`},
	}

	nameForms = map[string][]string{
		// §	   NAME
		`2.7.1`: {`nRootArcForm`},
		`2.7.2`: {`nArcForm`},
		`2.7.3`: {`dotNotationForm`},
	}

	schema = map[string]map[string][]string{
		`at`: attributeTypes,
		`oc`: objectClasses,
		`nf`: nameForms,
	}
}
