package radir

import (
	"fmt"
	"testing"
)

/*
This example demonstrates use of the [ProfileSettings.Keys] method to
access slices of string values, each representing a key name found
within the receiver instance.
*/
func ExampleProfileSettings_Keys() {
	prof := myDedicatedProfile
	prof.Settings().Set(`KeyX`, true) // set it
	fmt.Println(prof.Settings().Keys())
	// Output: [KeyX]
}

/*
This example demonstrates use of the [ProfileSettings.Value] method to
access an unasserted value of any kind.
*/
func ExampleProfileSettings_Value() {
	prof := myDedicatedProfile
	prof.Settings().Set(`KeyX`, []rune{'H', 'E', 'L', 'L', 'O', '.'}) // set it
	value, found := prof.Settings().Value(`KeyX`)
	if !found {
		fmt.Println("Failed to locate 'KeyX'")
		return
	}
	fmt.Printf("%#v\n", value)
	// Output: []int32{72, 69, 76, 76, 79, 46}
}

/*
This example demonstrates use of the [ProfileSettings.IsZero] method to
ascertain whether the instance is uninitialized.
*/
func ExampleProfileSettings_IsZero() {
	prof := myDedicatedProfile
	fmt.Println(prof.Settings().IsZero())
	// Output: false
}

/*
This example demonstrates the means of accessing a known Boolean key/value
pair within an instance of [ProfileSettings].

Case IS significant in the matching process.

In the case of the return values, the first Boolean value is the actual
value set earlier. The second Boolean value is the presence indicator,
which removes the ambiguity that would normally arise had there only been
a single (value) returned.
*/
func ExampleProfileSettings_BoolValue() {
	prof := myDedicatedProfile
	s := prof.Settings()
	s.Set(`KeyX`, true)                 // set it
	value, found := s.BoolValue(`KeyX`) // call it
	fmt.Printf("KeyX value: %t, found: %t", value, found)
	// Output: KeyX value: true, found: true
}

/*
This example demonstrates the means of accessing a known string key/value
pair within an instance of [ProfileSettings].

Case IS significant in the matching process.
*/
func ExampleProfileSettings_StringValue() {
	prof := myDedicatedProfile
	s := prof.Settings()
	s.Set(`KeyX`, `abcxyz`)               // set it
	value, found := s.StringValue(`KeyX`) // call it
	fmt.Printf("KeyX value: %s, found: %t", value, found)
	// Output: KeyX value: abcxyz, found: true
}

/*
This example demonstrates the means of accessing a known string slices
key/value pair within an instance of [ProfileSettings].

Case IS significant in the matching process.
*/
func ExampleProfileSettings_StringSliceValue() {
	prof := myDedicatedProfile
	s := prof.Settings()
	s.Set(`KeyX`, []string{`abc`, `xyz`})      // set it
	value, found := s.StringSliceValue(`KeyX`) // call it
	fmt.Printf("KeyX value: %v, found: %t", value, found)
	// Output: KeyX value: [abc xyz], found: true
}

/*
This example demonstrates the means for removing a key and its value.

Case IS significant in the matching process.
*/
func ExampleProfileSettings_Delete() {
	prof := myDedicatedProfile
	s := prof.Settings()
	s.Set(`KeyX`, true) // Assign Boolean true to map key "KeyX"
	//s.Delete(`keyX`)		 // This deletion would fail (case mismatch)
	s.Delete(`KeyX`)                // This deletion succeeds (case match)
	_, found := s.BoolValue(`KeyX`) // Call it, and focus only on presence
	fmt.Printf("KeyX found: %t", found)
	// Output: KeyX found: false
}

/*
This example demonstrates the means for accessing the [ProfileSettings]
instance found within the receiver instance of *[DITProfile].
*/
func ExampleDITProfile_Settings() {
	prof := myDedicatedProfile
	s := prof.Settings()
	s.Set(`Key1`, `some value`)
	s.Set(`Key2`, true)
	s.Set(`Key3`, []string{`many`, `values`})
	s.Set(`Key4`, struct {
		Field string
	}{
		Field: `another value`,
	})

	fmt.Printf("%d fields found", prof.Settings().Len())
	// Output: 4 fields found
}

func TestDITProfile_codecov(t *testing.T) {
	myDedicatedProfile.Marshal(func(_ any) error {
		return nil
	})
	ex := *myDedicatedProfile
	ex.R_DN = `adn`
	ex.R_OC = []string{`top`, `other`}
	X := &ex
	X.LDIF()
	X.AllowsRegistrants()
	X.RegistrantSuffixEqual("registrantID=12345,ou=Registrants,o=rA")
	X.RegistrationSuffixEqual("n=1,n=2,ou=Registrations,o=rA")

	var fx *DITProfile
	fx.NewRegistrant()
	fx.AllowsRegistrants()
	fx = new(DITProfile)
	fx.RegistrationBase(13)
	fx.NumRegistrationBase()
	fx.RegistrantBase(13)
	fx.NumRegistrantBase()
	fx.Model()
	fx.SetModel(ThreeDimensional)
	fx.Model()
	fx.SetRegistrantBase(`12345`)
	fx.AllowsRegistrants()
	fx.Mail()
	fx.URI()
	fx.SetMail("noreply@example.com")
	fx.SetURI("http://example.com")

	fx.NumMail()
	fx.Mail(15)
	fx.URI(15)
	fx.Mail(0)
	fx.URI(0)

	X.RegistrationBase(-1)
	X.RegistrantBase(-1)
	X.R_RegBase = []string{`blarg`}
	X.R_AthyBase = []string{`blarg`}
	X.RegistrationBase(-1)
	X.RegistrantBase(-1)

	D := &DUAConfig{R_DSE: X}
	D.NumProfile()
	D.Profile(15)
	D.Profile()
	D.Profile(0)
	D = &DUAConfig{R_DSE: X, R_Profiles: []*DITProfile{X}}
	E := &DUAConfig{R_Profiles: []*DITProfile{X}}
	D.Profile(-1)
	D.Profile(0).UseAltAuthorityTypes(true)
	D.Profile(0)
	D.Profile(15)
	D.Profile(0).GoverningStructureRule()
	D.Profile(0).GoverningStructureRuleGetFunc(nil)
	D.Profile(0).TTLGetFunc(nil)
	E.Profile(0)
	E.NumProfile()

	E.Profile(0).RegistrationBaseGetFunc(nil)
	E.Profile(0).RegistrantBaseGetFunc(nil)
	E.Profile(0).MailGetFunc(nil)
	E.Profile(0).URIGetFunc(nil)
	E.Profile(0).ModelGetFunc(nil)

	X.NewRegistration(true)
	X.MakeCache(100, 100)
	X.DropCache()
	X.Settings().Len()
	X.Settings().Set(`key0`, `value`)
	X.Settings().Set(`key1`, true)
	X.Settings().Set(`key2`, []string{`this`, `that`})
	for _, key := range X.Settings().Keys() {
		X.Settings().Delete(key)
	}

	X.Settings().IsZero()
	X.Settings().Keys()
	X.Settings().Value(`key0`)
	X.Settings().BoolValue(`key1`)
	X.Settings().StringValue(`key0`)
	X.Settings().StringSliceValue(`key2`)
}
