package radir

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/oid-directory/go-radir/oid"
)

/*
This example demonstrates using the *[Spatial] type and a
[GetOrSetFunc] to "chop" the RDN off a DN value and strip
away the attribute type, revealing only the lone RDN value,
which in this example would be a numberForm ("n").

For instance:

   n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA

... would become simply ...

   56521

This is useful for simplifying values associated with
"minArc", "leftArc", et al., rather than handling full
DN values unnecessarily.
*/
func ExampleSpatial_dNChopper() {
	// Define our "GetterOrSetter"
	getterOrSetter := func(x ...any) (any, error) {
		switch tv := x[0].(type) {
		case string:
			idxA := strings.IndexRune(tv, '=')
			idxC := strings.IndexRune(tv, ',')
			if idxA != -1 && idxC != -1 {
				// If you wanted to take this opportunity
				// to strconv the string to an actual number,
				// such as int or *big.Int, etc., you could
				// do it here :)
				return tv[idxA+1:idxC], nil
			}
		}
		return nil, fmt.Errorf("Some unexpected error occurred")
	}

	// Here is our test registration. This is very
	// watered down just for the sake of a simple
	// example.
	reg := &Registration{
		R_DN: "n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA",
		R_Spatial: &Spatial{
			R_LeftArc: "n=56520,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA",
			R_RightArc: "n=56522,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA",
		},
	}

	// Finally, call the desired method (LeftArcGetFunc) and
	// pass it our getterOrSetter definition.
	nf, err := reg.Spatial().LeftArcGetFunc(getterOrSetter)
	if err == nil {
		fmt.Printf("%s", nf)
	}
	// Output: 56520
}

func TestRegistrantDNGenerator(t *testing.T) {
	RegistrantDNGenerator()
	RegistrantDNGenerator(nil, nil)
	RegistrantDNGenerator(rune(10))
	RegistrantDNGenerator(``, ``)

	reg := myDedicatedProfile.NewRegistrant()
	reg.SetDN(`test`)
	RegistrantDNGenerator(reg)
}

func TestASN1NotationToMulti(t *testing.T) {
	for _, bogus := range []string{
		`{iso identified-organization(3)`,
		`{bad identified-organization(3)}`,
		`{iso(1) wrong(9838474)}`,
		`{iso(1) ""}`,
		`{ ,  , }`,
		`{  _  }`,
		``,
	} {
		cleanASN1(bogus)
		nanfToIdAndNF(bogus)
		if _, err := ASN1NotationToMulti(bogus); err == nil {
			t.Errorf("%s failed: expected err, got nothing", t.Name())
			return
		}
	}
}

func TestGS_codecov(t *testing.T) {
	ditp2 := &DITProfile{
		R_Model:   TwoDimensional,
		R_RegBase: []string{`ou=Registrations,o=rA`},
	}
	ditp3 := &DITProfile{
		R_Model:   ThreeDimensional,
		R_RegBase: []string{`ou=Registrations,o=rA`},
	}
	twoDc := &Registration{R_DITProfile: ditp2}

	_, _ = DNToDotNot3D()
	_, _ = DNToDotNot3D(nil, &X667{})
	_, _ = DNToDotNot3D(nil, &X680{r_DITProfile: ditp3})
	_, _ = DNToDotNot3D([]string{}, &X680{r_DITProfile: ditp3})
	_, _ = DNToDotNot3D(``, &X680{r_DITProfile: ditp2})
	_, _ = DNToDotNot3D(`ou=Fake,dc=example,dc=com`, &X680{r_DITProfile: ditp3})
	_, _ = DNToDotNot3D(`dotNotation=Xdc=exampledc=com`, &X680{r_DITProfile: ditp3})
	_, _ = DNToDotNot3D(`n=56521,1,n=4,n=1,6,3,n=1,ou=Registrationso=rA`, &X680{r_DITProfile: ditp3})
	_, _ = DNToDotNot3D(`n=56521,1,n=4,n=1,6,3,n=1,ou=Registrations,o=rA`, &X680{r_DITProfile: ditp3})
	_, _ = DNToDotNot3D(`n=,ou=Registrations,o=rA`, &X680{r_DITProfile: ditp3})
	_, _ = DNToDotNot3D(`,ou=Registrations,o=rA`, &X680{r_DITProfile: ditp3})

	_, _ = DNToDotNot2D()
	_, _ = DNToDotNot2D(nil, &X667{})
	_, _ = DNToDotNot2D(nil, &X680{r_DITProfile: ditp2})
	_, _ = DNToDotNot2D([]string{}, &X680{r_DITProfile: ditp2})
	_, _ = DNToDotNot2D(``, &X680{r_DITProfile: ditp3})
	_, _ = DNToDotNot2D(`ou=Fake,dc=example,dc=com`, &X680{r_DITProfile: ditp2})
	_, _ = DNToDotNot2D(`dotNotation=Xdc=exampledc=com`, &X680{r_DITProfile: ditp2})
	_, _ = DNToDotNot2D(`dotRotation=X,ou=Registrations,o=rA`, &X680{r_DITProfile: ditp2})
	_, _ = DNToDotNot2D(`dotnotation=X,ou=Registrationso=rA`, &X680{r_DITProfile: ditp2})
	_, _ = DNToDotNot2D(`,ou=Registrations,o=rA`, &X680{r_DITProfile: ditp2})
	_, _ = DNToDotNot2D(`dotRotation=X,ou=Registrations,o=rA`, &X680{r_DITProfile: ditp2})
	_, _ = DotNotToDN2D(nil)
	_, _ = DotNotToDN2D(nil, &X680{})
	_, _ = DotNotToDN2D(nil, &Registration{})
	_, _ = DotNotToDN2D(rune(10), twoDc)
	_, _ = DotNotToDN2D(`1.2.X`, twoDc)
	_, _ = DotNotToDN3D()
	_, _ = DotNotToDN3D(nil, &X680{})
	_, _ = DotNotToDN3D(nil, &Registration{})
	_, _ = DotNotToDN3D(nil, &Registration{R_DITProfile: ditp2})
	_, _ = DotNotToDN3D([]string{`g`}, &Registration{R_DITProfile: ditp3})
	_, _ = DotNotToDN3D(`1.3.6.1`, &Registration{R_DITProfile: ditp3})
	_, _ = DotNotToDN3D(`1.3.F.1`, &Registration{R_DITProfile: ditp3})

	for _, ident := range []string{
		`_hello`,
		`JERRY.HELLO`,
		`hh世界hgd`,
		`JERRY?HELLO`,
		`8fj`,
		`8&j`,
		`hell-o-`,
		``,
		`_`,
		`hell--o`,
	} {
		if oid.IsNameForm(ident) {
			t.Errorf("%s: expected error, got nothing", t.Name())
			return
		}
	}

	reg := myDedicatedProfile.NewRegistration()
	bint := big.NewInt(0)
	bint.SetString("4389248392084903280582390589038520985902", 10)
	reg.Supplement().SetInfo(bint)
	reg.X680().SetN(bint)
	reg.X680().SetN(1)
	bint.SetString("-4385902", 10)
	reg.X680().SetN(bint)
	reg.X680().SetASN1Notation([]string{`trash`})
	_, _, _ = cleanASN1(`{{}}`)
	_, _, _ = cleanASN1(`{iso someThing(X)}`)
	isNumber(`x`)
	isNumber(``)
	rootClass(-1)
	getRoot('E')
	getRoot('0')
	getRoot('1')
	getRoot('2')
	rootClass(0)
	rootClass(1)
	rootClass(2)
	rootClass(3)
	assertGetOrSetFunc(nil, nil)
	getFieldByNameTag(nil, ``)
	getFieldByNameTag(struct {
		Fake string `ldap:"fake"`
	}{
		Fake: `value`,
	}, `fayke`)
	getFieldValueByNameTagAndGoSF(nil, nil, ``)
	getFieldValueByNameTagAndGoSF(nil, func(_ ...any) (any, error) {
		return nil, nil
	}, ``)
	getFieldByNameTag(&struct{}{}, `l`)
	writeInt(`c`, 5, valOf(&struct{}{}))
	writeInt(`c`, 5, valOf(`hello`))
	writeString(`cn`, `hi`, valOf(rune(13)))
	writeValue(nil, nil, `n`)
	var nothing any = `4`
	writeValue(&struct {
		n *any `ldap:"n"`
	}{
		n: &nothing,
	}, &nothing, `n`)

	readFieldByTag(``, nil)
	unmarshalStruct(nil, nil)
	nanfToIdAndNF(`someThing(X)`)
	structEmpty(
		struct {
			Field []any
		}{
			Field: []any{`a`, `b`},
		})
}

func TestTime_codecov(t *testing.T) {
	var ts []any = []any{
		`20010718155634-0600.019283Z`,
		`20010718155634.019283432Z`,
		nil,
		13,
		`200`,
		``,
		`0600.0192834378297`,
	}

	for _, thyme := range ts {
		if assert, ok := thyme.(string); ok {
			if _, err := gt2t(assert); err == nil {
				t.Errorf("%s failed: expected error, got nothing", t.Name())
				return
			}
		}

		if _, err := GeneralizedTimeToTime(thyme); err == nil {
			t.Errorf("%s failed: expected error, got nothing", t.Name())
			return
		}
	}

	if _, err := GeneralizedTimeToTime(); err == nil {
		t.Errorf("%s failed: expected error, got nothing", t.Name())
		return
	}

	if _, err := GeneralizedTimeToTime([]string{`1`, `2`, `3`}); err == nil {
		t.Errorf("%s failed: expected error, got nothing", t.Name())
		return
	}

	if _, err := TimeToGeneralizedTime(); err == nil {
		t.Errorf("%s failed: expected error, got nothing", t.Name())
		return
	}

	if _, err := TimeToGeneralizedTime(ts); err == nil {
		t.Errorf("%s failed: expected error, got nothing", t.Name())
		return
	}

	if _, err := TimeToGeneralizedTime([]time.Time{now()}); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}
}
