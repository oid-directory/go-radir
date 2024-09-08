package radir

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var myCache *Cache // for go pkg examples

func ExampleCache_RegistrationLen() {
	fmt.Println(myCache.RegistrationLen())
	// Output: 2
}

func ExampleCache_Registration() {
	dn := "n=3,n=1,n=3,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA"
	threeD := myCache.Registration(dn)
	fmt.Println(threeD.Description())
	// Output: [The Three Dimensional Model]
}

func ExampleCache_TouchRegistration() {
	dn := "n=3,n=1,n=3,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA"
	myCache.TouchRegistration(dn, 2)
	fmt.Println(myCache.RegistrationExpired(dn))
	// Output: false
}

func ExampleCache_RegistrationCap() {
	fmt.Println(myCache.RegistrationCap())
	// Output: 0
}

func ExampleCache_RegistrantLen() {
	fmt.Println(myCache.RegistrantLen())
	// Output: 1
}

func ExampleCache_Registrant() {
	dn := "registrantID=jAHdfNm328,ou=Registrants,o=rA"
	alloc := myCache.Registrant(dn)
	fmt.Println(alloc.FirstAuthority().CN())
	// Output: Jesse Coretta
}

func ExampleCache_TouchRegistrant() {
	dn := "registrantID=jAHdfNm328,ou=Registrants,o=rA"
	myCache.TouchRegistrant(dn, 2)
	fmt.Println(myCache.RegistrantExpired(dn))
	// Output: false
}

func ExampleCache_RegistrantCap() {
	fmt.Println(myCache.RegistrantCap())
	// Output: 0
}

func TestCache_writes(t *testing.T) {
	// Create a temp file
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}
	defer os.RemoveAll(tempDir)

	var myCache3 *Cache
	myCache2 := NewCache(0, 0)

	fileName := `registrations.cache`
	path := filepath.Join(tempDir, fileName)
	myCache3.LoadRegistrations(``)
	myCache3.LoadRegistrations(path)
	myCache3.Freeze()
	myCache3.LoadRegistrations(path)
	myCache3.WriteRegistrations(path)
	myCache2.Freeze()
	myCache2.WriteRegistrations(``)
	myCache2.WriteRegistrations(path)
	myCache2.LoadRegistrations(`a`)
	myCache2.Thaw()
	myCache2.LoadRegistrations(path)
	myCache2.LoadRegistrations(`a`)

	fileName = `registrants.cache`
	path = filepath.Join(tempDir, fileName)
	myCache3.WriteRegistrants(``)
	myCache3.LoadRegistrants(path)
	myCache2.WriteRegistrants(path)
	myCache2.Freeze()
	myCache2.WriteRegistrants(path)
	myCache2.LoadRegistrants(`a`)
	myCache2.Thaw()
	myCache2.LoadRegistrants(path)
	myCache2.LoadRegistrants(`a`)
}

func TestCache_codecov(t *testing.T) {
	r := myDedicatedProfile.NewRegistration()
	r.X680().SetASN1Notation(`{joint-iso-itu-t(2) asn1(1)}`)
	r.X680().SetDotNotation(`2.1`)
	r.X680().SetIRI(`/ASN.1`)
	r.X660().SetLongArc(`/ASN.1`)	
	r.X660().LongArc()

	var c *Cache
	c.IsZero()
	c.RegistrationExpired("fargus")
	c.RegistrationExpired("")
	c.RegistrantExpired("fargus")
	c.RegistrantExpired("")
	c.Add(nil, -1)
	c.Add(&Registration{}, 1)
	c.Add(&Registration{R_DN: "fargus"}, 1)
	c.Add(&Registrant{}, 1)
	c.Add(&Registrant{R_DN: "fargus"}, 1)
	c.Registration("blarg")
	c.Registration("fargus")
	c.Registrant("blarg")
	c.Registrant("fargus")
	c.Flush()
	c.Tidy()
	c.TouchRegistration(`blarg`, 1)
	c.TouchRegistration(`who`, -1)
	c.RemoveRegistration(`blarg`)
	c.RemoveRegistration(``)
	c.TouchRegistrant(`blarg`, 1)
	c.TouchRegistrant(`who`, -1)
	c.RemoveRegistrant(`blarg`)
	c = NewCache(-5, -1)
	c = NewCache(1, 1)
	fakeR := cachedRegistration{
		Value: &Registration{R_DN: "fake"},
	}
	otherFakeR := cachedRegistration{
		Value: &Registration{R_DN: "otherFaker"},
	}
	otherFakeR2 := cachedRegistration{
		Value: &Registration{R_DN: "otherFaker2"},
	}
	c.registrations[`fake`] = fakeR
	c.Add(fakeR, 1)
	c.Registration(`faker`)
	c.Add(otherFakeR, 1)
	c.Add(otherFakeR2, 1)
	c.RemoveRegistrant(`bob`)
	c.RemoveRegistrant(``)
	c.Tidy()
	c.registrations[`fake`] = fakeR
	c.Registration(`fake`)
	c.registrations[`fake`] = fakeR
	c.Flush()
	c.RegistrationKeys()
	c.RegistrantKeys()
	fakeS := cachedRegistrant{
		Value: &Registrant{R_DN: "fake"},
	}
	c.registrants[`fake`] = fakeS
	c.Add(fakeS, 1)
	c.RemoveRegistrant(`bob`)
	c.Tidy()
	c.registrants[`fake`] = fakeS
	c.Registrant(`fake`)
	c.registrants[`fake`] = fakeS
	c.Flush()
}

func init() {
	// add some items to the cache just
	// for the sake of examples.
	threeDeeDN := `n=3,n=1,n=3,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`
	twoDeeDN := `n=2,n=1,n=3,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`

	myCache = NewCache(0, 0)
	myCache.expired(``, 1)
	myCache.expired(``, 0)
	myCache.expired(`&*((*&(&^R`, 0)
	threeDee := myDedicatedProfile.NewRegistration()
	threeDee.SetDN(threeDeeDN)
	threeDee.X680().SetN("3")
	threeDee.X680().SetASN1Notation("{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 oid-directory(101) rA-DIT(3) models(1) threeDimensional(3)}")
	threeDee.X680().SetDotNotation("1.3.6.1.4.1.56521.101.3.1.3")
	threeDee.X680().SetIdentifier("threeDimensional")
	threeDee.Supplement().SetLeafNode(true)
	threeDee.X680().SetNameAndNumberForm("threeDimensional(3)")
	threeDee.Spatial().SetSupArc("1")
	threeDee.Spatial().SetTopArc("1")
	threeDee.X660().SetCurrentAuthorities("registrantID=jAHdfNm328,ou=Registrants,o=rA")
	threeDee.SetObjectClasses([]string{
		"x680Context",
		"x660Context",
		"registrationSupplement",
		"spatialContext",
		"iSORegistration",
	})
	threeDee.SetDescription("The Three Dimensional Model")
	threeDee.X680().SetIRI("/ISO/Identified-Organization/6/1/4/1/56521/101/3/1/3")

	twoDee := myDedicatedProfile.NewRegistration()
	twoDee.SetDN(twoDeeDN)
	threeDee.X680().SetN("2")
	threeDee.X680().SetASN1Notation("{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 oid-directory(101) rA-DIT(3) models(1) twoDimensional(2)}")
	threeDee.X680().SetDotNotation("1.3.6.1.4.1.56521.101.3.1.2")
	threeDee.X680().SetIdentifier("twoDimensional")
	threeDee.Supplement().SetLeafNode(true)
	threeDee.X680().SetNameAndNumberForm("twoDimensional(2)")
	threeDee.Spatial().SetSupArc("1")
	threeDee.Spatial().SetTopArc("1")
	threeDee.X660().SetCurrentAuthorities("registrantID=jAHdfNm328,ou=Registrants,o=rA")
	threeDee.SetObjectClasses([]string{
		"x680Context",
		"x660Context",
		"registrationSupplement",
		"spatialContext",
		"iSORegistration",
	})
	twoDee.SetDescription("The Two Dimensional Model")
	twoDee.X680().SetIRI("/ISO/Identified-Organization/6/1/4/1/56521/101/3/1/2")

	jesse := myDedicatedProfile.NewRegistrant()
	jesse.SetDN("registrantID=jAHdfNm328,ou=Registrants,o=rA")
	jesse.FirstAuthority().SetCN("Jesse Coretta")
	jesse.FirstAuthority().SetO("Individual")
	jesse.FirstAuthority().SetEmail("jesse.coretta@icloud.com")

	courtney := myDedicatedProfile.NewRegistrant()
	courtney.SetDN("registrantID=N0389j4DFs,ou=Registrants,o=rA")
	courtney.FirstAuthority().SetCN("Courtney Tolana")
	courtney.FirstAuthority().SetO("Individual")


	limited := NewCache(1,1)	
	limited.Add(threeDee, 0)
	limited.Add(threeDee, 1)
	limited.Add(twoDee, 0)
	limited.Add(twoDee, 1)
	limited.Add(jesse,1)
	limited.Add(courtney,0)
	limited.Add(courtney,1)

	myCache.Add(nil, -1)
	myCache.Add(threeDee, 0)
	myCache.Add(threeDee, 1)
	myCache.Add(twoDee, 0)
	myCache.Add(twoDee, 1)
	myCache.RegistrationKeys()
	myCache.expired(threeDeeDN, 1)
	myCache.touch(threeDeeDN, 1, 0)
	myCache.touch(``, 1, 0)
	myCache.remove(0)
	myCache.remove(1)
	myCache.remove(0, []string{`this`, `that`}...)
	myCache.remove(1, []string{`this`, `that`}...)
	myCache.Freeze()
	myCache.remove(0, []string{`this`, `that`}...)
	myCache.remove(1, []string{`this`, `that`}...)
	myCache.Thaw()

	myCache.Add(jesse, 1)
	myCache.Add(jesse, 0)
	myCache.RegistrantKeys()
	myCache.expired(`registrantID=jAHdfNm328,ou=Registrants,o=rA`, 1)
	myCache.touch(`n=3,n=1,n=3,n=101,n=56521,n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA`, 1, 1)
	myCache.touch(``, 1, 1)
}
