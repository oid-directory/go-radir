package radir

import (
	"fmt"
	"strings"
	"testing"
)

func ExampleRegistrant_Unmarshal() {
	reg := myDedicatedProfile.NewRegistrant()
	reg.FirstAuthority().SetCN("Some person")
	reg.FirstAuthority().SetCO("United States")
	reg.Sponsor().SetO("Benevolent Sponsors, Co.")
	um := reg.Unmarshal()
	fmt.Println(um[`sponsorOrg`][0])
	// Output: Benevolent Sponsors, Co.
}

func ExampleRegistrant_LDIF() {
	prof := myDedicatedProfile

	reg := prof.NewRegistrant()
	reg.SetDN("registrantID=16ddcdfddeb2e37," + prof.RegistrantBase())
	reg.FirstAuthority().SetCN("Some person")
	reg.FirstAuthority().SetCO("United States")
	reg.Sponsor().SetO("Benevolent Sponsors, Co.")
	fmt.Println(reg.LDIF())
	// Output: dn: registrantID=16ddcdfddeb2e37,ou=Registrants,o=rA
	// objectClass: top
	// objectClass: registrant
	// objectClass: firstAuthorityContext
	// objectClass: sponsorContext
	// firstAuthorityCountryName: United States
	// firstAuthorityCommonName: Some person
	// sponsorOrg: Benevolent Sponsors, Co.
}

/*
This example demonstrates use of a [GetOrSetFunc] instance to write the
input value to the receiver including a randomly-generated "[registrantID]".

This is contrary to use of [Registrant.DNGetFunc], which only alters the
presentation value, not the receiver field value.

[registrantID]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.34
*/
func ExampleRegistrant_SetDN_withSetOrGetFunc() {
	regi := myDedicatedProfile.NewRegistrant()

	if err := regi.SetDN(RegistrantDNGenerator); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Random registrant DN generated: %d bytes\n", len(regi.DN()))
	// Output: Random registrant DN generated: 43 bytes
}

/*
This example demonstrates use of a [GetOrSetFunc] instance to write the
input value to the receiver. For simplicity, this example merely introduces
a lowercase-normalizer as our [GetOrSetFunc].

This is contrary to use of [Registrant.SetDN], which will write the value
returned by the [GetOrSetFunc] to the receiver struct field value.
*/
func ExampleRegistrant_DNGetFunc() {
	regi := myDedicatedProfile.NewRegistrant()
	gosf := func(x ...any) (any, error) {
		ret := x[0].(string)
		ret = strings.ToLower(ret)
		return ret, nil
	}
	regi.SetDN("registrantID=EB3ED0aD39f,ou=Registrants,o=rA", gosf)
	out, err := regi.DNGetFunc(gosf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out.(string))
	// Output: registrantid=eb3ed0ad39f,ou=registrants,o=ra
}

func ExampleRegistrant_ID() {
	reg := myDedicatedProfile.NewRegistrant()
	reg.SetID("A35BEF04F")
	fmt.Println(reg.ID())
	// Output: A35BEF04F
}

/*
This example demonstrates the means of determining the descriptor of the
"[registrant]" STRUCTURAL class that is implemented through instances of
this type.

[registrant]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.16
*/
func ExampleRegistrant_StructuralObjectClass() {
	var reg *Registrant = myDedicatedProfile.NewRegistrant() // assume this is populated
	fmt.Println(reg.StructuralObjectClass())
	// Output: registrant
}

/*
This example demonstrates the means of applying a [GetOrSetFunc] instance
to a "[registrantID]" without changing the underlying value.

[registrantID]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.34
*/
func ExampleRegistrant_IDGetFunc() {
	reg := myDedicatedProfile.NewRegistrant()
	reg.SetID("A35BEF04F")
	gosf := func(x ...any) (any, error) {
		ret := x[0].(string)
		ret = strings.ToLower(ret)
		return ret, nil
	}
	out, err := reg.IDGetFunc(gosf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out.(string))
	// Output: a35bef04f
}

func TestAuthority_codecov(t *testing.T) {
	for idx, err := range []error{
		bogusRegistrant_codecov(),
		bogusFirstAuthority_codecov(),
		bogusCurrentAuthority_codecov(),
		bogusSponsor_codecov(),
	} {
		if err != nil {
			t.Errorf("%s[%d] failed: %v", t.Name(), idx, err)
			return
		}
	}

	var reg *Registrant // nil
	reg.isEmpty()

	ffunk := func(_ any) error {
		return nil
	}
	reg.Marshal(ffunk)

	reg.Unmarshal()
	reg.FirstAuthority()
	reg.CurrentAuthority()
	reg.Sponsor()
	reg.Dedicated()

	reg = &Registrant{}
	reg.Marshal(nil)
	reg.Marshal(ffunk)
	reg.isEmpty()
	reg.Dedicated()
	reg.RC_TTL = "5"
	reg.Unmarshal()
	reg.R_DITProfile = myDedicatedProfile
	reg.Marshal(ffunk)
	reg.Marshal(func(_ any) error {
		return fmt.Errorf("FAIL")
	})
	reg.R_OC = append(reg.R_OC, "firstAuthorityContext")
	reg.refreshObjectClasses()
	reg.R_OC = append(reg.R_OC, []string{"firstAuthorityContext", "top", "top", "registrant"}...)
	reg.refreshObjectClasses()
	reg.TTL()

	regs := Registrants{reg}
	regs.Len()
	regs.Push(&Registrant{R_Id: `meat`})
	regs.LDIF()
	regs.Index(0)
	regs.Index(-1)
	regs.Contains(`meat`)
	regs.Contains(`meats`)
	regs.Get(`meat`)
	regs.Get(`meats`)
	regs.Len()
	regs.Unmarshal()
	regs.Marshal(&DITProfile{}, ffunk)
	regs.Marshal(myCombinedProfile, ffunk)
	regs.Marshal(myDedicatedProfile, ffunk)
}

func bogusRegistrant_codecov() error {
	var nilReg *Registrant = myDedicatedProfile.NewRegistrant()

	if err := testBogusDedicatedSetters(nilReg); err != nil {
		return err
	}
	if err := testBogusDedicatedGetters(nilReg); err != nil {
		return err
	}

	nilReg.DN()
	nilReg.GoverningStructureRule()
	nilReg.CollectiveAttributeSubentries()
	nilReg.Description()
	nilReg.StructuralObjectClass()
	nilReg.ObjectClasses()
	nilReg.SeeAlso()
	nilReg.TTL()

	return nil
}

func testBogusDedicatedGetters(nilReg *Registrant) error {
	nilReg.R_GSR = "testing"
	nilReg.R_SOC = "testing"
	for _, funk := range []func(GetOrSetFunc) (any, error){
		nilReg.DNGetFunc,
		nilReg.DescriptionGetFunc,
		nilReg.GoverningStructureRuleGetFunc,
		nilReg.ObjectClassesGetFunc,
		nilReg.StructuralObjectClassGetFunc,
		nilReg.CollectiveAttributeSubentriesGetFunc,
		nilReg.SeeAlsoGetFunc,
		nilReg.TTLGetFunc,
	} {
		funk(nil)
		val, err := funk(func(v ...any) (any, error) {
			// fake, do whatever
			return v[0], nil
		})

		if err != nil {
			return err
		}

		txt := `testing`
		if assert, ok := val.(string); ok {
			if assert != txt {
				return errorf("Mismatched; want '%s', got '%s'",
					txt, assert)
			}
		} else if sassert, sok := val.([]string); sok {
			if len(sassert) == 1 {
				if sassert[0] != txt {
					return errorf("Mismatched; want '%s', got '%s'",
						txt, sassert[0])
				}
			}
		} else {
			return errorf("Unsupported type '%T'", val)
		}
	}

	return nil
}

func testBogusDedicatedSetters(nilReg *Registrant) error {
	var bogus any = []int{1, 2, 3, 4}

	for _, funk := range []func(...any) error{
		nilReg.SetDN,
		nilReg.SetDescription,
		nilReg.SetSeeAlso,
		nilReg.SetTTL,
		nilReg.SetObjectClasses,
	} {
		for _, err := range []error{
			funk(),
			funk(bogus),
		} {
			if err == nil {
				return mkerr("Expected error, got nothing")
			}
		}

		if err := funk(`testing`); err != nil {
			return err
		}
	}

	return nil
}
