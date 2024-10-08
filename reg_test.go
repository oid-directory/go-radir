package radir

import (
	"fmt"
	"testing"
)

/*
This example demonstrates an ill-fated attempt to write a "[longArc]"
value via *[X660.SetLongArc] upon a *[Registration] instance that extends
from something other than the requisite "Joint-ISO-ITU-T" root (in this
case, "ISO").

[longArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.20
*/
func ExampleX660_LongArc_violation() {
	reg := myDedicatedProfile.NewRegistration()

	// Here we use the official OID/ASN1 prefix
	// constants for input simply for brevity.
	// Any value would do ...
	reg.X680().SetASN1Notation(ASN1Prefix)
	reg.X680().SetDotNotation(OIDPrefix) // 1.3.6.1.4.1.56521.101 (an ISO allocation)

	// Note also that this would fail even if we did
	// it before the previous Set<X> execs.
	if err := reg.X660().SetLongArc("/OID-Directory"); err != nil {
		fmt.Println(err)
	}
	// Output: longArc values can only be assigned to sub arcs of Joint-ISO-ITU-T

}

/*
This example demonstrates accessing information about the root through
an arc *[Registration] instance. Use of this method is only meaningful
on non-root arcs.
*/
func ExampleRegistration_Root_arc() {
	reg := myDedicatedProfile.NewRegistration()
	// Set either the aSN1Notation
	// OR the dotNotation form to
	// activate root awareness.
	// Once set, it cannot be unset.
	reg.X680().SetASN1Notation(ASN1Prefix) // use I-D prefix for simplicity
	root, class := reg.Root()

	fmt.Printf("%d = %s", root, class)
	// Output: 1 = iSORegistration
}

func ExampleRegistration_NewChild() {
	dad := myDedicatedProfile.NewRegistration()
	dad.X680().SetN(`5`)
	dad.X680().SetASN1Notation(`{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 test(5)}`)
	dad.X680().SetIdentifier(`test`)
	dad.X680().SetDotNotation(`1.3.6.1.4.1.56521.5`)
	dad.X680().SetNameAndNumberForm(`test(5)`)
	dad.SetDN(`n=5,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)

	son := dad.NewChild(`10`, `child`)
	fmt.Printf("%s :: %s\n", son.DN(), son.X680().NameAndNumberForm())
	// Output: n=10,n=5,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA :: child(10)
}

func ExampleRegistration_NewSibling() {
	bro := myDedicatedProfile.NewRegistration()
	bro.X680().SetN(`5`)
	bro.X680().SetASN1Notation(`{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 test(5)}`)
	bro.X680().SetIdentifier(`test`)
	bro.X680().SetNameAndNumberForm(`test(5)`)
	bro.SetDN(`n=5,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)

	sis := bro.NewSibling(`999`, `example`)
	fmt.Printf("%s\n", sis.DN())
	// Output: n=999,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA
}

/*
This example demonstrates creating a new child *[Registration] using a root
*[Registration] as the parent (source).
*/
func ExampleRegistration_NewChild_fromRoot() {
	// iso root
	root1 := myDedicatedProfile.NewRegistration(true)
	root1.SetDN(`n=1,ou=Registrations,o=rA`)
	root1.X680().SetASN1Notation(`{iso(1)}`)
	root1.X680().SetIdentifier(`iso`)
	root1.X680().SetN(`1`)

	// Here we initialize the Spatial type, but won't need
	// to populate it. This trick will result in the source
	// instance's DN becoming the topArc and/or supArc when
	// appropriate, but only if already initialized.
	root1.Spatial()

	root1dot3 := root1.NewChild(`3`, `identified-organization`)

	fmt.Printf("DN:%s\nNumberForm:%s\nASN1Notation:%s\nNameAndNumberForm:%s",
		root1dot3.DN(),
		root1dot3.X680().N(),
		root1dot3.X680().ASN1Notation(),
		root1dot3.X680().NameAndNumberForm())
	// Output: DN:n=3,n=1,ou=Registrations,o=rA
	// NumberForm:3
	// ASN1Notation:{iso(1) identified-organization(3)}
	// NameAndNumberForm:identified-organization(3)
}

/*
This example demonstrates a means of quickly initializing a new (root)
instance of *[Registration] using another root instance as an initializer.

One need only feed the desired number form and identifier string values
to the [Registration.NewSibling] method to receive the new root instance.

Note that all *[Registration] instances which represent root allocations
MUST bear a number form of 0, 1 or 2, else it is considered invalid and
a nil sibling instance would be returned. Also note that the receiver
instance (the source) must have a DN prior to creating a sibling.
*/
func ExampleRegistration_NewSibling_fromRoot() {
	// iso root
	root1 := myDedicatedProfile.NewRegistration(true)
	root1.SetDN(`n=1,ou=Registrations,o=rA`)
	root1.X680().SetIdentifier(`iso`)
	root1.X680().SetN(`1`)

	// itu-t root
	root0 := root1.NewSibling(`0`, `itu-t`)
	fmt.Printf("DN:%s\nNumberForm:%s\nASN1Notation:%s\nNameAndNumberForm:%s",
		root0.DN(),
		root0.X680().N(),
		root0.X680().ASN1Notation(),
		root0.X680().NameAndNumberForm())
	// Output: DN:n=0,ou=Registrations,o=rA
	// NumberForm:0
	// ASN1Notation:{itu-t(0)}
	// NameAndNumberForm:itu-t(0)
}

func ExampleX680_Depth() {
	reg := myDedicatedProfile.NewRegistration()
	reg.X680().SetASN1Notation(ASN1Prefix) // use I-D prefix for simplicity
	fmt.Println(reg.X680().Depth())
	// Output: 8
}

func ExampleASN1NotationToMulti() {
	slice, err := ASN1NotationToMulti(ASN1Prefix) // use I-D prefix for simplicity
	if err != nil {
		fmt.Println(err)
		return
	}

	// We skip slice 0 since that is the original input value (a).
	// Slice 1 is dotNotation
	// Slice 2 is identifier (if defined)
	// Slice 3 is nameAndNumberForm (if identifier was defined)
	// Slice 4 is the number form
	fmt.Println(slice[1:])
	// Output: [1.3.6.1.4.1.56521.101 oid-directory oid-directory(101) 101]
}

func ExampleASN1NotationToMulti_root() {
	a := `{iso}`
	slice, err := ASN1NotationToMulti(a)
	if err != nil {
		fmt.Println(err)
		return
	}

	// We skip slice 0 since that is the original input value (a).
	// Slice 1 is <nil>; roots do not use dotNotation
	// Slice 2 is identifier (if defined)
	// Slice 3 is nameAndNumberForm (if identifier was defined)
	// Slice 4 is the number form
	fmt.Printf("%#v\n", slice)
	// Output: []string{"{iso(1)}", "", "iso", "iso(1)", "1"}
}

/*
This example demonstrates a means of setting up to six (6) values using
the output from processing a single "[aSN1Notation]" value. YMMV.

[aSN1Notation]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.4
*/
func ExampleASN1NotationToMulti_sixValuesFromOne() {
	reg := myDedicatedProfile.NewRegistration()

	// This is the same as our ASN1Prefix constant,
	// except we've added tabs and linebreaks for
	// readability here.
	a := `{
		iso
		identified-organization(3)
		dod(6)
		internet(1)
		private(4)
		enterprise(1)
		56521
		oid-directory(101)
	}`

	// Here we parse our ASN.1 notation value (a).
	slice, err := ASN1NotationToMulti(a)
	if err != nil {
		fmt.Println(err)
		return
	}

	// The order of the return value is known
	// in advance ...
	//
	//   - original (cleaned-up) ASN.1 value
	//   - dotNotation value
	//   - identifier value
	//   - nameAndNumberForm value
	//   - numberForm value
	//
	// ... so we'll mirror this order when
	// creating a list of "Set<X>" methods
	// (of equal length) for execution ...
	for idx, funk := range []func(...any) error{
		reg.X680().SetASN1Notation,
		reg.X680().SetDotNotation,
		reg.X680().SetIdentifier,
		reg.X680().SetNameAndNumberForm,
		reg.X680().SetN,
	} {
		// Don't process zero values
		if slice[idx] == "" {
			continue
		}
		// execute funk #idx upon slice #idx
		if err := funk(slice[idx]); err != nil {
			fmt.Println(err)
			return
		}
	}

	// Since we now have a dotNotation value set within
	// our Registration instance, why not set the DN?
	//
	// Simply feed the SetDN method our freshly-set dot
	// notation. Additionally, since our configuration
	// profile is 3D, we'll include the DotNotToDN3D
	// method for specialized handling.
	if err := reg.SetDN(reg.X680().DotNotation(), DotNotToDN3D); err != nil {
		fmt.Println(err)
		return
	}

	subent := myDedicatedProfile.NewSubentry()
	subent.SetDN(`cn=test-subentry,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	subent.SetCN(`test-subentry`)
	//subent.SetSubtreeSpecification(`{minimum 1, maximum 1, specificationFilter {}}`)
	reg.Subentries().Push(subent)

	// Take a look at what we got so far. There are
	// other ways of examining the values, but LDIF
	// is quick and easy (not to mention relevant).
	fmt.Println(reg.LDIF(0))
	// Output: dn: n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA
	// objectClass: top
	// objectClass: registration
	// objectClass: arc
	// objectClass: x680Context
	// n: 101
	// aSN1Notation: {iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 oid-directory(101)}
	// dotNotation: 1.3.6.1.4.1.56521.101
	// identifier: oid-directory
	// nameAndNumberForm: oid-directory(101)
}

func TestRegistration_Allocate(t *testing.T) {
	prof := myDedicatedProfile

	// Create the ISO arc
	iso := prof.NewRegistration(true)
	iso.SetDN(`n=1,ou=Registrations,o=rA`)
	iso.X680().SetN(`1`)
	iso.X680().SetASN1Notation(`{iso(1)}`)
	iso.X680().SetIRI(`/ISO`)
	iso.X660().SetUnicodeValue(`ISO`)

	for _, str := range [][]string{
		{`ASN.1`, `{iso(1)
			member-body(2) 56521}`},
		{`ASN.1`,
			`{iso(1)
			identified-organization(3)	
			dod(6)
			internet(1)
			private(4)
			enterprise(1)
			56521}`},
		{`Dot`, `1.3.6.1.4.1.56521`},
	} {
		//alloc := iso.Allocate(str[1])
		if walked := iso.Allocate(str[1]); walked.X680().N() != `56521` {
			t.Errorf("%s [%s alloc/walk] failed: want 56521, got '%s'",
				t.Name(), str[0], walked.X680().N())
			return
		}
	}

	iso.allocateDotNot([]string{`1`, `3`, `6`, `1`, `5`, `5`, `7`, `2`, `0`}, `reserved`)
	iso.allocateDotNot([]string{`1`}, `reserved`)
	iso.allocateDotNot([]string{`1`, `2`}, `reserved`)
	iso.allocateDotNot([]string{`1`, `2`, `3`}, `reserved`)
	iso.allocateASN1([][]string{{`iso`, `1`}})
	nanfToSlice(`iso`)
	nanfToSlice(`itu-t`)
	nanfToSlice(`joint-iso-itu-t`)
	iso.Allocate([]string{})
	iso.Allocate([]string{"3"})
	iso.Allocate([][]string{})
	iso.Allocate(`0.0.4`)
	iso.Allocate([][]string{
		{"identified-organization", "3"},
		{"dod", "6"}})
	iso.Allocate(nil)

	// Create the ITU-T arc
	itu := prof.NewRegistration(true)
	itu.SetDN(`n=0,ou=Registrations,o=rA`)
	itu.X680().SetN(`0`)
	itu.X680().SetASN1Notation(`{itu-t(0)}`)
	itu.X680().SetIRI(`/ITU-T`)
	itu.X660().SetUnicodeValue(`ITU-T`)

	//memberID := `{itu-t}`
	//memberN := `0`

	for _, str := range [][]string{
		{`ASN.1`, `{itu-t(0) member-body(2)}`},
		{`ASN.1`, `{itu-t(0) 2}`},
		{`Dot`, `0.2`},
	} {
		itu.Allocate(str[1])
		if walked := itu.Walk(str[1]); walked.X680().N() != `2` {
			t.Errorf("%s [%s walk] failed: want 2, got '%s'",
				t.Name(), str[0], walked.X680().N())
		}
	}
}

func TestRegistration_Walk(t *testing.T) {
	prof := myDedicatedProfile
	iso := prof.NewRegistration(true)
	iso.SetDN(`n=1,ou=Registrations,o=rA`)
	iso.X680().SetN(`1`)
	iso.X680().SetASN1Notation(`{iso(1)}`)
	iso.X680().SetIRI(`/ISO`)
	iso.X660().SetUnicodeValue(`ISO`)

	org := iso.NewChild(`3`, `identified-organization`)
	org.X660().SetUnicodeValue(`Identified-Organization`)
	org.X680().SetIRI(`/ISO/Identified-Organization`)

	dod := org.NewChild(`6`, `dod`)
	dod.X680().SetIRI(`/ISO/Identified-Organization/6`)

	dodOID := `1.3.6`

	var none *Registration
	none.walkASN1(nil)

	defense := iso.Walk(dodOID)
	if defense.X680().DotNotation() != dodOID {
		t.Errorf("%s failed: could not find dod", t.Name())
	}
}

func TestITUXSeries_unmarshal(t *testing.T) {
	w := &X690{r_root: new(registeredRoot)}
	w.SetDotEncoding(`BgEr`)
	wval := w.unmarshal()[`dotEncoding`][0]
	want := `BgEr`
	if wval != want {
		t.Errorf("%s X.690 failed: want '%s', got '%s'", t.Name(), want, wval)
		return
	}

	x := &X680{r_root: new(registeredRoot)}
	x.SetASN1Notation(`{iso identified-organization(3)}`)
	x.SetDotNotation(`1.3`)
	x.SetIdentifier(`identified-organization`)
	x.SetNameAndNumberForm(`{org}`)
	x.SetIRI(`/ISO/Identified-Organization`)
	xval := x.unmarshal()[`aSN1Notation`][0]
	want = `{iso identified-organization(3)}`
	if xval != want {
		t.Errorf("%s X.680 failed: want '%s', got '%s'", t.Name(), want, xval)
		return
	}

	y := &X667{r_root: new(registeredRoot)}
	y.SetRegisteredUUID(`4a3cf69e-3b6d-40c5-9b72-25b28f431e16`)
	yval := y.unmarshal()[`registeredUUID`][0]
	want = `4a3cf69e-3b6d-40c5-9b72-25b28f431e16`
	if yval != want {
		t.Errorf("%s X.667 failed: want '%s', got '%s'", t.Name(), want, yval)
		return
	}

	var yourReg *Registration = myDedicatedProfile.NewRegistration()
	z := yourReg.X660()
	z.SetUnicodeValue(`Identified-Organization`)
	z.SetAdditionalUnicodeValue(`Org`) // fake
	z.SetSecondaryIdentifier(`org`)
	z.SetStdNameForm(`{org}`)
	z.SetFirstAuthorities(`cn=Predecessor,...`)
	z.SetCurrentAuthorities(`cn=You,...`)
	z.SetSponsors(`cn=Your Sponsor,...`)
	z.CombinedSponsor().SetCN("cn=Your Combined Sponsor,...")
	um := z.unmarshal()
	zval := um[`unicodeValue`][0]
	want = `Identified-Organization`
	if zval != want {
		t.Errorf("%s X.660 failed: want '%s', got '%s'", t.Name(), want, zval)
		return
	}

	yourReg.SetDN(`n=999,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	yourReg.X680().SetDotNotation(`1.3.6.1.4.1.56521.999`)
	_, _, _, _, _ = yourReg.sibOrSub(`-1`, ``, true)
	_, _, _, _, _ = yourReg.sibOrSub(`1`, `$TUPID`, false)
	_, _, _, _, _ = yourReg.sibOrSub(`15`, `fake`, true)
	yourReg.Spatial().SetTopArc(`n=1,ou=Registration,o=rA`)
	_ = yourReg.NewSibling(`15`, `fake`)
	_ = yourReg.NewSibling(`-1`, ``)
	_ = yourReg.NewChild(`15`, `fake`)
	_ = yourReg.NewChild(`-15`, ``)
	yourReg.R_OC = append(yourReg.R_OC, `rootArc`)
	yourReg.X680().R_N = `0`
	_ = yourReg.NewSibling(`15`, `fake`)
	yourReg.R_OC = []string{}
	_ = yourReg.NewChild(`15`, `fake`)
	_ = yourReg.NewSibling(`15`, `fake`)

}

func TestSpatial(t *testing.T) {
	var nilSpat *Spatial = new(Spatial)

	var val string = `testing`
	for _, funk := range []func(...any) error{
		nilSpat.SetLeftArc,
		nilSpat.SetRightArc,
		nilSpat.SetTopArc,
		nilSpat.SetSupArc,
		nilSpat.SetMinArc,
		nilSpat.SetMaxArc,
		nilSpat.SetSubArc,
	} {
		if err := funk(val); err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
			return
		}
	}

	// Test unmarshal
	if um := nilSpat.unmarshal(); um == nil {
		t.Errorf("%s failed: nil spatial unmarshal", t.Name())
		return
	}

	// We have to artificially add values to
	// collective spatial types because they
	// would only come from marshaling entry
	// content obtained from LDAP search.
	nilSpat.RC_MaxArc = val
	nilSpat.RC_MinArc = val
	nilSpat.RC_SupArc = val
	nilSpat.RC_TopArc = val

	for _, funk := range []func() string{
		nilSpat.CMaxArc,
		nilSpat.CMinArc,
		nilSpat.CSupArc,
		nilSpat.CTopArc,
		nilSpat.LeftArc,
		nilSpat.MaxArc,
		nilSpat.MinArc,
		nilSpat.RightArc,
		nilSpat.SupArc,
		nilSpat.TopArc,
	} {
		if value := funk(); val != value {
			t.Errorf("%s failed: want '%s', got '%s'", t.Name(), val, value)
			return
		}
	}

	if value := nilSpat.SubArc(); value[0] != val {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), val, value[0])
		return
	}

	for idx, funk := range []func(GetOrSetFunc) (any, error){
		nilSpat.MaxArcGetFunc,
		nilSpat.LeftArcGetFunc,
		nilSpat.MinArcGetFunc,
		nilSpat.SubArcGetFunc,
		nilSpat.SupArcGetFunc,
		nilSpat.RightArcGetFunc,
		nilSpat.TopArcGetFunc,
		nilSpat.CMaxArcGetFunc,
		nilSpat.CMinArcGetFunc,
		nilSpat.CSupArcGetFunc,
		nilSpat.CTopArcGetFunc,
	} {
		funk(nil)
		out, err := funk(func(v ...any) (any, error) {
			// fake, do whatever
			return v[0], nil
		})

		if err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
			return
		}

		if assert, ok := out.(string); ok {
			if assert != val {
				t.Errorf("%s: mismatched slice %d; want '%s', got '%s'",
					t.Name(), idx, val, assert)
				break
			}
		} else if sassert, sok := out.([]string); sok {
			if len(sassert) == 1 {
				if sassert[0] != val {
					t.Errorf("%s: Mismatched; want '%s', got '%s'",
						t.Name(), val, sassert[0])
					break
				}
			}
		} else {
			t.Errorf("%s: Unsupported type '%T'", t.Name(), out)
			break
		}
	}
}

func TestRegistration_codecov(t *testing.T) {
	if err := bogusRegistration_codecov(); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}
}

func TestRegistrations(t *testing.T) {
	var regs Registrations

	regs.Less(1, 2)
	regs.Swap(1, 2)

	nreg1 := myDedicatedProfile.NewRegistration()
	nreg2 := myDedicatedProfile.NewRegistration()

	o1 := `1.3.6.1.4.1.56521.999.18.1`
	o2 := `1.3.6.1.4.1.56521.999.18.2`

	nreg1.SetDN(`n=1,n=18,n=999,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	nreg1.X680().SetDotNotation(o1)
	nreg1.X680().SetN(`1`)

	nreg2.SetDN(`n=2,n=18,n=999,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	nreg2.X680().SetDotNotation(o2)
	nreg2.X680().SetN(`2`)
	nreg2.Spatial().SetTopArc("n=1,ou=Registrations,o=rA")
	nreg2.Spatial().SetSupArc(`n=18,n=999,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	_ = nreg2.NewChild(`33`, `thisName`) // no need for var
	nreg2.Children().SetYAxes(true)
	nreg2.Walk(nil)
	nreg2.Size()
	nreg2.LDIF(-1)
	nreg2.LDIF(0)
	nreg2.LDIF(1)
	nreg2.LDIF(2)
	sube := myDedicatedProfile.NewSubentry()
	sube.SetDN(`cn=test-subentry,n=1,n=18,n=999,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`)
	sube.SetCN(`test-subentry`)

	_ = IsASN1Notation(`{iso(1)}`)
	_ = IsNumericOID(`1.2.3`)
	_ = IsIdentifier(`l`)
	_ = IsNameAndNumberForm(`l(1)`)
	_ = IsNumberForm(`1`)

	twoDPro := NewFactoryDefaultDUAConfig()
	twoDPro.R_DSE.R_Model = TwoDimensional
	nreg3 := twoDPro.Profile().NewRegistration()
	nreg3.X680().SetDotNotation(`1.3.6.1.4.1.56521.999.8`)
	nreg3.SetDN(nreg3.X680().DotNotation(), DotNotToDN2D)
	nreg3.CollectiveAttributeSubentries()

	nreg2.Children().Push(nreg1)
	nreg2.Children().Push(nreg3)
	nreg2.Subentries().Push(sube)
	nreg3.Subentries().Push(sube)

	nreg2.LDIF(0, true)
	nreg2.LDIF(0, false)

	nreg2.LDIF(2, true)
	nreg2.LDIF(2, false)

	nreg2.subtreeLDIF(true)
	nreg3.subtreeLDIF(true)
	nreg2.subtreeLDIF(false)
	nreg3.subtreeLDIF(false)

	nreg3.SetXAxes(true)
	nreg3.SetXAxes()
	nreg3.SetYAxes(true)
	nreg3.SetYAxes()
	nreg3.X680().SetIdentifier(`dad`)
	nreg3.X680().SetASN1Notation(`{iso identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 example(999) dad(8)}`)
	nreg3.X680().SetN(`8`)
	_ = nreg3.NewChild(`14`, `son`)

	regs.SetXAxes()
	regs.Push(nreg2) // ordered incorrectly
	regs.SetXAxes(true)
	regs.SetYAxes(true)
	regs.Push(nreg1)

	if L := regs.Len(); L != 2 {
		t.Errorf("%s failed; want '%d', got '%d'", t.Name(), 2, L)
		return
	}

	regs.SortByNumberForm() // reorder
	if N := regs.Get(`1`).X680().N(); N != `1` {
		t.Errorf("%s failed; want '%s', got '%s'", t.Name(), `1`, N)
		return
	}

	regs.SetXAxes()
	regs.SetYAxes(true)

	if !regs.Contains(`2`) {
		t.Errorf("%s failed; '%s' not found", t.Name(), o2)
		return
	}

	dad := regs.Get(o2)
	_, as, err := cleanASN1(`{1 3 6 1 4 1 56521 999 18 2}`)
	if err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	nanfs := make([][]string, 0)
	for i := 0; i < len(as); i++ {
		nanfs = append(nanfs, nanfToSlice(as[i]))
	}

	dad.allocateASN1(nanfs)
	dad.Size()
	dad.LDIF(0)
	dad.LDIF(1)
	dad.LDIF(2)

	// codecov
	regs.Push(&Registration{R_X680: &X680{R_N: "2"}})
	regs.Push(&Registration{R_X680: &X680{R_N: "1"}})

	altRegs := &Registrations{
		&Registration{},
		&Registration{},
	}

	altRegs.Less(0, 1)
	(*altRegs)[0].R_X680 = &X680{}
	(*altRegs)[1].R_X680 = &X680{}
	altRegs.Less(0, 1)

	regs.SetXAxes(true)
	regs.SetYAxes(true)

	if regs.Less(2, 3) {
		t.Errorf("%s failed: want false, got true", t.Name())
		return
	}
	regs.SortByNumberForm(true)

	var upreg *Registration
	upreg.Subentries()

	regs.Push(&Registration{})
	regs.Push(&Registration{})
	regs.Less(8, 9)
	regs.HasParents()

	regs.Push(&Registration{R_X680: &X680{}})
	regs.Push(&Registration{R_X680: &X680{}})

	if regs.Less(10, 11) {
		t.Errorf("%s failed: want false, got true", t.Name())
	}

	regs = append(regs, &Registration{})
	regs.SetXAxes()
}

func TestNumericOID(t *testing.T) {
	for _, oid := range []string{
		`1.3.6.1.4.1`,
		`2.5`,
		`0.0.4`,
	} {
		if !IsNumericOID(oid) {
			t.Errorf("%s failed: genuine OID flagged as bogus (%s)", t.Name(), oid)
			return
		}
	}

	for _, bad := range []string{
		`2`,
		``,
		`3.1`,
		`$.3`,
		`1.2.3..4.5`,
		`1.2.3.4.5.`,
		`4.2.3.t.5`,
		`_`,
		`1.50`,
		`1.S0`,
	} {
		if IsNumericOID(bad) {
			t.Errorf("%s failed: bogus OID flagged as genuine (%s)", t.Name(), bad)
			return
		}
	}
}

func bogusRegistration_codecov() error {
	_ = errorf(MismatchedLeafErr)
	_ = errorf("MismatchedLeafErr")
	_ = errorf("")

	efunk := func(_ any) error {
		return fmt.Errorf("FAIL")
	}

	ffunk := func(_ any) error {
		return nil
	}

	var empty *Registration
	_, _, _, _, _ = empty.sibOrSub(`-1`, ``, false)
	empty.Unmarshal()
	empty.Dedicated()
	empty.Combined()
	empty.StructuralObjectClass()
	empty.CollectiveAttributeSubentries()
	regs := Registrations{empty}
	regs.Unmarshal()
	regs.Marshal(&DITProfile{}, efunk)
	regs.Marshal(myDedicatedProfile, efunk)
	regs.Marshal(myCombinedProfile, ffunk)

	for _, nilReg := range []*Registration{
		myDedicatedProfile.NewRegistration(true),
		myDedicatedProfile.NewRegistration(),
		myCombinedProfile.NewRegistration(true),
		myCombinedProfile.NewRegistration(),
	} {
		if err := testBogusRegistrationSetters(nilReg); err != nil {
			return err
		}
		if err := testBogusRegistrationGetters(nilReg); err != nil {
			return err
		}

		regs = append(regs, nilReg)
		regs.Unmarshal()
		regs.Index(0)
		regs.Index(0).Size()

		nilReg.X680().dotNotationHandler(``)
		nilReg.X680().dotNotationHandler(`.`)
		nilReg.X680().asn1NotationHandler(`{iso}`)

		nilReg.CollectiveAttributeSubentries()
		regs.Push(nilReg)

		nilReg.R_OC = append(nilReg.R_OC, `x690Context`)
		nilReg.R_OC = append(nilReg.R_OC, `x680Context`)
		nilReg.R_OC = append(nilReg.R_OC, `x667Context`)
		nilReg.R_OC = append(nilReg.R_OC, `x660Context`)
		nilReg.R_OC = append(nilReg.R_OC, `top`)
		nilReg.refreshObjectClasses()
		nilReg.X680().SetDotNotation(`1.3.6.1.4.1.56521.999.0.8`)
		nilReg.X680().SetN(`8`)
		nilReg.X680().SetASN1Notation(`{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 reserved(0) eight(8)}`)
		nilReg.SetDN(nilReg.X680().DotNotation(), DotNotToDN3D)
		nilReg.R_OC = []string{}
		nilReg.R_SOC = ``
		nilReg.NewChild(`1`, `this`)
		nilReg.NewSibling(`2`, `that`)
		nilReg.NewSubentry(`subentry`)

		var em *Registration
		em.X660()
		em.X660().profile()
		em.X667()
		em.X667().profile()
		em.X680()
		em.X680().profile()
		em.X690()
		em.X690().profile()
		em.Profile()
		em.Root()
		em.Spatial()
		em.Supplement()
		em.GoverningStructureRule()
		nilReg.CollectiveAttributeSubentries()
		nilReg.GoverningStructureRule()
		nilReg.DN()
		nilReg.Root()
		nilReg.Dedicated()
		nilReg.Combined()
		//nilReg.R_TTL = "5"
		nilReg.Marshal(nil)
		nilReg.Marshal(ffunk)
		nilReg.Marshal(efunk)
		nilReg.StructuralObjectClass()
		em.Marshal(ffunk)
		em.Marshal(efunk)
		//nilReg.R_DITProfile = &DITProfile{R_TTL: "5"}
		nilReg.isEmpty()
		nilReg.Description()
		nilReg.ObjectClasses()
		nilReg.SeeAlso()
		nilReg.LDIF(0)
		nilReg.LDIF(0, true)
		nilReg.LDIF(0, false)
		nilReg.LDIF(2)
		nilReg.LDIF(2, true)
		nilReg.LDIF(2, false)
		nilReg.TTL()
		nilReg.X690().DotEncoding()
		nilReg.X680().ASN1Notation()
		nilReg.X680().DotNotation()
		nilReg.X680().Identifier()
		nilReg.X680().NameAndNumberForm()
		nilReg.X680().IRI()
		nilReg.X680().N()
		nilReg.X667().profile()
		nilReg.X667().RegisteredUUID()
		nilReg.X660().LongArc()
		nilReg.X660().marshal(ffunk)
		nilReg.X660().marshal(efunk)
		nilReg.X660().StdNameForm()
		nilReg.X660().AdditionalUnicodeValue()
		nilReg.X660().CCurrentAuthorities()
		nilReg.X660().CurrentAuthorities()
		nilReg.X660().CombinedCurrentAuthority()
		nilReg.X660().CFirstAuthorities()
		nilReg.X660().FirstAuthorities()
		nilReg.X660().CombinedFirstAuthority()
		nilReg.X660().SecondaryIdentifier()
		nilReg.X660().Sponsors()
		nilReg.X660().CombinedSponsor()
		nilReg.X660().CSponsors()
		nilReg.X660().UnicodeValue()
		nilReg.Supplement().Status()
		nilReg.Supplement().Info()
		nilReg.Supplement().Range()
		nilReg.Supplement().ModifyTime()
		nilReg.Supplement().URI()
		nilReg.Supplement().CreateTime()
		nilReg.Supplement().Classification()
		nilReg.Supplement().CDiscloseTo()
		nilReg.Supplement().DiscloseTo()
		nilReg.Supplement().Frozen()
		nilReg.Supplement().LeafNode()
	}

	return nil
}

func testBogusRegistrationGetters(nilReg *Registration) error {
	nilReg.Profile()
	nilReg.X660().profile()
	nilReg.X667().profile()
	nilReg.X680().profile()
	nilReg.X690().profile()
	nilReg.SetDITProfile(myCombinedProfile)
	nilReg.Supplement().RC_DiscloseTo = []string{"testing"}
	nilReg.Supplement().R_Class = "testing"
	nilReg.R_GSR = "testing"
	nilReg.R_SOC = "testing"
	nilReg.R_TTL = ""
	nilReg.RC_TTL = ""
	nilReg.TTL()
	nilReg.Parent()
	nilReg.R_DN = "n=1,n=3,n=1,ou=Registrations,o=rA"
	nilReg.SetDITProfile(myCombinedProfile)
	nilReg.Profile().R_TTL = "5"
	nilReg.R_DN = "testing"
	nilReg.TTL()
	nilReg.R_TTL = "testing"
	nilReg.RC_TTL = "testing"
	nilReg.X660()
	nilReg.X660().r_root.Depth = 3
	nilReg.X660().r_root.N = 2
	nilReg.X660().R_LongArc = []string{"testing"}
	nilReg.X660().LongArc()

	for idx, funk := range []func(GetOrSetFunc) (any, error){
		nilReg.DNGetFunc,
		nilReg.DescriptionGetFunc,
		nilReg.ObjectClassesGetFunc,
		nilReg.StructuralObjectClassGetFunc,
		nilReg.GoverningStructureRuleGetFunc,
		nilReg.CollectiveAttributeSubentriesGetFunc,
		nilReg.SeeAlsoGetFunc,
		nilReg.TTLGetFunc,
		nilReg.X690().DotEncodingGetFunc,
		nilReg.X680().NGetFunc,
		nilReg.X680().DotNotationGetFunc,
		nilReg.X680().ASN1NotationGetFunc,
		nilReg.X680().IdentifierGetFunc,
		nilReg.X680().IRIGetFunc,
		nilReg.X680().NameAndNumberFormGetFunc,
		nilReg.X667().RegisteredUUIDGetFunc,
		nilReg.X660().LongArcGetFunc,
		nilReg.X660().SponsorsGetFunc,
		nilReg.X660().SecondaryIdentifierGetFunc,
		nilReg.X660().StdNameFormGetFunc,
		nilReg.X660().UnicodeValueGetFunc,
		nilReg.X660().AdditionalUnicodeValueGetFunc,
		nilReg.X660().CurrentAuthoritiesGetFunc,
		nilReg.X660().FirstAuthoritiesGetFunc,
		nilReg.X660().CFirstAuthoritiesGetFunc,
		nilReg.X660().CCurrentAuthoritiesGetFunc,
		nilReg.X660().CSponsorsGetFunc,
		nilReg.Supplement().CreateTimeGetFunc,
		nilReg.Supplement().InfoGetFunc,
		nilReg.Supplement().ClassificationGetFunc,
		nilReg.Supplement().CDiscloseToGetFunc,
		nilReg.Supplement().DiscloseToGetFunc,
		nilReg.Supplement().FrozenGetFunc,
		nilReg.Supplement().LeafNodeGetFunc,
		nilReg.Supplement().ModifyTimeGetFunc,
		nilReg.Supplement().RangeGetFunc,
		nilReg.Supplement().StatusGetFunc,
		nilReg.Supplement().URIGetFunc,
	} {

		funk(nil)
		val, err := funk(func(v ...any) (any, error) {
			// fake, do whatever
			return v[0], nil
		})

		if err != nil && err != RegistrantPolicyErr {
			return err
		}

		txt := `testing`
		if assert, ok := val.(string); ok {
			if assert != txt {
				return errorf("Mismatched string slice %d; want '%s', got '%s'",
					idx, txt, assert)
			}
		} else if sassert, sok := val.([]string); sok {
			if len(sassert) == 1 {
				if sassert[0] != txt {
					return errorf("Mismatched []string slice; want '%s', got '%s'",
						txt, sassert[0])
				}
			}
		} else {
			if idx == 21 && !nilReg.Combined() {
				return errorf("Unsupported type '%T' at idx %d", val, idx)
			}
		}
	}

	return nil
}

func testBogusRegistrationSetters(nilReg *Registration) error {
	var bogus any = []int{1, 2, 3, 4}

	for idx, funk := range []func(...any) error{
		nilReg.SetDN,
		nilReg.SetSeeAlso,
		nilReg.SetDescription,
		nilReg.SetObjectClasses,
		nilReg.SetTTL,
		nilReg.X690().SetDotEncoding,
		nilReg.X680().SetIRI,
		nilReg.X680().SetN,
		nilReg.X680().SetIdentifier,
		nilReg.X680().SetASN1Notation,
		nilReg.X680().SetNameAndNumberForm,
		nilReg.X680().SetDotNotation,
		nilReg.X667().SetRegisteredUUID,
		nilReg.X660().SetSecondaryIdentifier,
		nilReg.X660().SetUnicodeValue,
		nilReg.X660().SetCurrentAuthorities,
		nilReg.X660().SetFirstAuthorities,
		nilReg.X660().SetSponsors,
		nilReg.X660().SetAdditionalUnicodeValue,
		nilReg.X660().SetStdNameForm,
		nilReg.Supplement().SetDiscloseTo,
		nilReg.Supplement().SetInfo,
		nilReg.Supplement().SetCreateTime,
		nilReg.Supplement().SetModifyTime,
		nilReg.Supplement().SetLeafNode,
		nilReg.Supplement().SetFrozen,
		nilReg.Supplement().SetStatus,
		nilReg.Supplement().SetClassification,
		nilReg.Supplement().SetRange,
		nilReg.Supplement().SetURI,
	} {
		for _, err := range []error{
			funk(),
			funk(bogus),
		} {
			if err == nil {
				return mkerr("Expected error, got nothing, at slice " + itoa(idx))
			}
		}

		if err := funk(`testing`); err != nil && err != RegistrantPolicyErr {
			return err
		}
	}

	return nil
}
