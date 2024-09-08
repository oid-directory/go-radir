package radir

import (
	"fmt"
	"time"
)

func ExampleNewFactoryDefaultDUAConfig() {
	cfg := NewFactoryDefaultDUAConfig()
	fmt.Println(cfg.Profile().LDIF())
	// Output: dn:
	// rADirectoryModel: 1.3.6.1.4.1.56521.101.3.1.3
	// rARegistrationBase: ou=Registrations,o=rA
	// rARegistrantBase: ou=Registrants,o=rA
}

func ExampleSupplement_ModifyTime() {
	var X *Registration = myDedicatedProfile.NewRegistration()
	X.SetDN(`n=11,n=3,n=2,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	X.X680().SetN(`11`)
	X.X680().SetASN1Notation(`{iso(1) org(3) dod(6) internet(1) private(4) enterprise(1) 56521 101 2 3 11}`)
	X.X680().SetDotNotation(`1.3.6.1.4.1.56521.101.2.3.11`)
	X.Supplement().SetModifyTime(`19711103040901Z`)
	X.Supplement().SetModifyTime(`19860110120110Z`)
	X.Supplement().SetModifyTime(`20080414090555Z`)

	fmt.Printf("%s", X.Supplement().ModifyTime()[0])
	// Output: 19711103040901Z
}

func ExampleSupplement_CreateTime() {
	var X *Registration = myDedicatedProfile.NewRegistration()
	X.SetDN(`n=11,n=3,n=2,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	X.X680().SetN(`11`)
	X.X680().SetASN1Notation(`{iso(1) org(3) dod(6) internet(1) private(4) enterprise(1) 56521 101 2 3 11}`)
	X.X680().SetDotNotation(`1.3.6.1.4.1.56521.101.2.3.11`)
	X.Supplement().SetCreateTime(`19860110120110Z`)

	fmt.Printf("%s", X.Supplement().CreateTime())
	// Output: 19860110120110Z
}

func ExampleRegistration_Kind() {
	var X *Registration = myDedicatedProfile.NewRegistration()
	X.SetDN(`n=11,n=3,n=2,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	X.X680().SetN(`11`)
	X.X680().SetASN1Notation(`{iso(1) org(3) dod(6) internet(1) private(4) enterprise(1) 56521 101 2 3 11}`)
	X.X680().SetDotNotation(`1.3.6.1.4.1.56521.101.2.3.11`)
	X.Supplement().SetCreateTime(`19860110120110Z`)

	fmt.Printf("%s", X.Kind())
	// Output: STRUCTURAL
}

func ExampleX680_SetDotNotation_withGetOrSetFunc3D() {
	var X *Registration = myDedicatedProfile.NewRegistration()
	X.SetDN(`n=11,n=3,n=2,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	if err := X.X680().SetDotNotation(X.DN(), DNToDotNot3D); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", X.X680().DotNotation())
	// Output: 1.3.6.1.4.1.56521.101.2.3.11
}

func ExampleX680_SetDotNotation_withGetOrSetFunc2D() {
	duac := &DITProfile{
		R_RegBase: []string{`ou=Registrations,o=rA`},
		R_Model:   TwoDimensional,
	}

	X := duac.NewRegistration()
	X.SetDN(`dotNotation=1.3.6.1.4.1.56521.101.2.3.11,ou=Registrations,o=rA`)

	if err := X.X680().SetDotNotation(X.DN(), DNToDotNot2D); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", X.X680().DotNotation())
	// Output: 1.3.6.1.4.1.56521.101.2.3.11
}

func ExampleRegistration_SetDN_withGetOrSetFunc3D() {
	var X *Registration = myDedicatedProfile.NewRegistration()
	X.X680().SetDotNotation(`1.3.6.1.4.1.56521.101.2.3.11`)

	if err := X.SetDN(X.X680().DotNotation(), DotNotToDN3D); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", X.DN())
	// Output: n=11,n=3,n=2,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA
}

func ExampleRegistration_SetDN_withGetOrSetFunc2D() {
	duac := &DITProfile{
		R_RegBase: []string{`ou=Registrations,o=rA`},
		R_Model:   TwoDimensional,
	}

	X := duac.NewRegistration()
	X.X680().SetDotNotation(`1.3.6.1.4.1.56521.101.2.3.11`)

	if err := X.SetDN(X.X680().DotNotation(), DotNotToDN2D); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", X.DN())
	// Output: dotNotation=1.3.6.1.4.1.56521.101.2.3.11,ou=Registrations,o=rA
}

/*
This example demonstrates use of the [GetOrSetFunc] type to alter the default
return value of a "getter".

In this case, instead of requesting the "registrationModified" value (which is
string or []string by default), we feed a compliant [GetOrSetFunc] instance --
the [GeneralizedTimeToTime] function -- to the [Registration.ModifyTimeGetFunc]
method in order to marshal the string value into one or more [time.Time] values.
*/
func ExampleSupplement_ModifyTimeGetFunc_withGetOrSetFuncGeneralizedTimeToTime() {
	var X *Registration = new(Registration)

	// give our instance a string val
	X.Supplement().SetModifyTime(`20080114214613Z`)

	// Hand the ModifyTimeGetFunc method the
	// desired GetOrSetFunc function instance
	// to be executed.
	t, err := X.Supplement().ModifyTimeGetFunc(GeneralizedTimeToTime)
	if err != nil {
		fmt.Println(err)
		return
	}

	if T, ok := t.(time.Time); ok && !T.IsZero() {
		// There is only one timestamp
		fmt.Printf("%s", T)
	} else if Ts, oks := t.([]time.Time); oks && len(Ts) > 0 {
		// There are multiple timestamps,
		// let's just take the first one.
		fmt.Printf("%s", Ts[0])
	}
	// Output: 2008-01-14 21:46:13 +0000 UTC
}

func ExampleRegistration_dedicated() {
	X := myDedicatedProfile.NewRegistration()

	X.SetDN(`n=11,n=3,n=2,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	X.X680().SetN(`11`)
	X.X680().SetASN1Notation(`{iso(1) org(3) dod(6) internet(1) private(4) enterprise(1) 56521 101 2 3 11}`)
	X.X680().SetDotNotation(`1.3.6.1.4.1.56521.101.2.3.11`)
	X.Supplement().SetCreateTime(`20210110120110Z`)
	X.X660().SetFirstAuthorities(`registrantID=Jesse Coretta,ou=Registrants,o=rA`)
	X.X660().SetCurrentAuthorities(`registrantID=ABC Authority,ou=Registrants,o=rA`)

	fmt.Printf("%s (first), %s (current)",
		X.X660().FirstAuthorities()[0], X.X660().CurrentAuthorities()[0])
	// Output: registrantID=Jesse Coretta,ou=Registrants,o=rA (first), registrantID=ABC Authority,ou=Registrants,o=rA (current)
}

func ExampleRegistration_combined() {
	X := myCombinedProfile.NewRegistration()

	X.SetDN(`n=11,n=3,n=2,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	X.X680().SetN(`11`)
	X.X680().SetASN1Notation(`{iso(1) org(3) dod(6) internet(1) private(4) enterprise(1) 56521 101 2 3 11}`)
	X.X680().SetDotNotation(`1.3.6.1.4.1.56521.101.2.3.11`)
	X.Supplement().SetCreateTime(`20210110120110Z`)
	X.X660().CombinedFirstAuthority().SetCN(`Mister Authority`)
	X.X660().CombinedFirstAuthority().SetO(`Authority, Co.`)

	fmt.Printf("%s at %s", X.X660().CombinedFirstAuthority().CN(), X.X660().CombinedFirstAuthority().O())
	// Output: Mister Authority at Authority, Co.
}

func ExampleRegistration_Unmarshal() {
	var X *Registration = myDedicatedProfile.NewRegistration()
	X.SetDN(`n=11,n=3,n=2,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	X.X680().SetN(`11`)
	X.X680().SetASN1Notation(`{iso(1) org(3) dod(6) internet(1) private(4) enterprise(1) 56521 101 2 3 11}`)
	m := X.Unmarshal()
	fmt.Printf("%s", m[`aSN1Notation`][0])
	// Output: {iso(1) org(3) dod(6) internet(1) private(4) enterprise(1) 56521 101 2 3 11}
}

/*
func Registrants_Unmarshal() {
	sp := new(Sponsor)
	sp.SetDN(`registrantID=430a8727-c8b0-4734-83d7-0a104ab00a69,ou=Registrants,o=rA`)
	sp.SetCN(`Mister Sponsor`)
	sp.SetO(`Benevolent Sponsors, Inc.`)
	sp.SetTel(`+1 555 555 0134`)
	sp.SetEmail(`sponsors@example.com`)

	fa := new(FirstAuthority)
	fa.SetDN(`registrantID=90b29aac-7771-1249-01a6-143f5d4d2500,ou=Registrants,o=rA`)
	fa.SetCN(`Jesse Coretta`)
	fa.SetTel(`+1 555 555 9812`)

	var regs Registrants = Registrants{sp, fa}

	maps := regs.Unmarshal()
	fmt.Printf("len:%d", len(maps))
	// Output: 2
}
*/

func ExampleGeneralizedTimeToTime() {
	t := `20080114154613-0600`
	if T, err := GeneralizedTimeToTime(t); err == nil {
		// Reset timezone to UTC (optional)
		T = T.(time.Time).UTC()
		fmt.Printf("%s", T)
	}

	// Output: 2008-01-14 21:46:13 +0000 UTC
}

func ExampleTimeToGeneralizedTime() {
	// Use the GenTimeToTime example for brevity.
	t := `20080114154613-0600`
	T, _ := GeneralizedTimeToTime(t)
	T = T.(time.Time).UTC() // normalize to UTC

	// Use the same time, but with utc (skip
	// bool checking for brevity).
	tu := `20080114214613Z`
	Tu, _ := GeneralizedTimeToTime(tu) // parse into time.Time

	// Parse the respective times back into
	// string-based generalizedTime values.
	nt, _ := TimeToGeneralizedTime(Tu) // parse into time.Time
	ntu, _ := TimeToGeneralizedTime(T) // parse into time.Time

	// They should be valid and identical
	fmt.Printf("Times are valid and equal: %t (%s)", nt == ntu, nt)
	// Output: Times are valid and equal: true (20080114214613Z)
}
