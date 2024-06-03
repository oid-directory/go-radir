package radir

import (
	"testing"
	//oID "github.com/JesseCoretta/go-objectid"
)

func TestConfig(t *testing.T) {
	var d DUAConfig = DUAConfig{
		DITProfile{
			R_DN: "",
			R_RegB: []string{"ou=Registration,o=rA"},
			R_AthB: []string{"ou=Registrant,o=rA"},
			R_Model: "1.3.6.1.4.1.56521.101.3.1.3",
			R_Desc: []string{"Registration Authority"},
		},
	}

	t.Logf("%v\n", d)
}

func TestRegistration(t *testing.T) {
	var arc JointISOITUT = JointISOITUT{
		Common: RegistrationCommon{
			R_N:    "14",
			R_Desc: "This is a registration",
		},
		X680: X680{
			R_DNot: "1.3.6.1.4.1.56521.14",
			R_ANot: "{iso identified-organization(3) 6 1 4 1 56521 14}",
			R_IRI: []string{
				`/ISO/Identified-Organization/6/1/4/1/56521/14`,
			},
		},
		R_LArc: []string{
			`/ASN.1`,
		},
		Extra: Supplemental{
			R_CTime: `19910212170805Z`,
			R_MTime: []string{
				`19910422134505Z`,
				`20080701040256.019485Z`,
			},
		},
	}

	arc.SetIdentifier(`testReg`)

	t.Logf("(%s) - %s, %s, %s, %s\n", arc.Type(), arc.N(), arc.DotNotation(), arc.Identifier(), arc.ASN1Notation())

	var arc2 RootArc = ISORoot
	t.Logf("(%s) - %s, %s, %s\n", arc2.Type(), arc2.N(), arc2.Identifier(), arc2.ASN1Notation())

	t.Logf("CT: %s\n", arc.CreateTime())
	t.Logf("MT: %v\n", arc.ModifyTime())
}

func TestAuthority(t *testing.T) {
	fa := FirstAuthority{
		R_StartTime: `20230601114005Z`,
		R_EndTime: `20240601114005Z`,
	}

	stime := fa.StartTime()
	etime := fa.EndTime()
	t.Logf("%s (%T)\n", stime, stime)
	t.Logf("%s (%T)\n", etime, etime)
}
