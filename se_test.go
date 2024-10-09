package radir

import (
	"fmt"
	"testing"
)

func TestSubentry_codecov(t *testing.T) {
	sents := make(Subentries, 0)

	var subentries *Subentries
	subentries.Len()
	subentries.Get(`hello`)
	subentries.Contains(`hello`)
	subentries.Index(21)
	subentries.IsZero()

	var subentry *Subentry
	subentry.Unmarshal()
	subentry.Marshal(func(any) error {
		return nil
	})

	subentry.Marshal(func(any) error {
		return fmt.Errorf("ERROR")
	})

	subentry.Supplement()
	subentry.X660()
	subentry.Spatial()
	subentry.IsZero()
	subentry.isEmpty()
	subentry = new(Subentry)
	subentry.refreshObjectClasses()
	subentry = myDedicatedProfile.NewSubentry()
	subentry.SetSubtreeSpecification(`{}`)
	subspec, _ := NewSubtreeSpecification(`{}`)
	subentry.SetSubtreeSpecification(subspec)
	subentry.SubtreeSpecification()
	subentry.SubtreeSpecificationGetFunc(nil)
	subentry.refreshObjectClasses()
	subentry.Supplement()
	subentry.X660().SetLongArc(`/Long`)
	subentry.X660().R_LongArc = []string{`fhdfsj`}
	subentry.R_OC = append(subentry.R_OC, `x690Context`)
	subentry.R_OC = append(subentry.R_OC, `x680Context`)
	subentry.R_OC = append(subentry.R_OC, `x667Context`)
	subentry.R_OC = append(subentry.R_OC, `x660Context`)
	subentry.R_OC = append(subentry.R_OC, `registrationSupplement`)
	subentry.R_OC = append(subentry.R_OC, `spatialContext`)
	subentry.R_OC = append(subentry.R_OC, `top`)
	subentry.R_OC = append(subentry.R_OC, "firstAuthorityContext")
	subentry.refreshObjectClasses()
	subentry.R_OC = append(subentry.R_OC, []string{"firstAuthorityContext", "top", "top", "subentry", "x660Context"}...)
	subentry.refreshObjectClasses()
	subentry.refreshObjectClasses()

	subentry.refreshObjectClasses()
	subentry.R_OC = removeStrInSlice(`x660Context`, subentry.R_OC)
	subentry.R_OC = removeStrInSlice(`spatialContext`, subentry.R_OC)
	subentry.R_OC = removeStrInSlice(`registrationSupplement`, subentry.R_OC)
	subentry.refreshObjectClasses()
	subentry.CN()
	subentry.SetCN(`spatialContext`)
	subentry.CNGetFunc(nil)
	subentry.X660()
	subentry.IsZero()
	subentry.Spatial()
	subentry.isEmpty()
	subentry.StructuralObjectClass()
	subentry.DN()
	subentry.ObjectClasses()
	subentry.ObjectClassesGetFunc(nil)
	subentry.Kind()
	subentry.r_DITProfile = &DITProfile{}
	subentry.SetDITProfile(&DITProfile{})
	subentry.Profile()
	subentry.SetDN(`n=1,ou=Registrations,o=rA`)
	subentry.SetDITProfile(myDedicatedProfile)
	subentry.TTL()
	subentry.TTLGetFunc(nil)
	subentry.SetTTL(5)
	subentry.CTTL()
	subentry.CTTLGetFunc(nil)
	subentry.SetCTTL(5)
	subentry.SetDN(`cn=spatialContext,n=6,n=3,n=1,ou=Registrations,o=rA`)
	subentry.R_STS = nil
	subentry.LDIF()
	subentry.refreshObjectClasses()
	subentry.SetObjectClasses(`x660Context`)
	subentry.refreshObjectClasses()
	subentry.DNGetFunc(nil)
	subentry.Unmarshal()
	subentry.Marshal(nil)
	subentry.R_STS = nil
	subentry.Marshal(func(any) error {
		return nil
	})
	subentry.Marshal(func(any) error {
		return fmt.Errorf("ERROR")
	})

	mtg := &Subentry{R_DN: `o=example-subentry,cn=Parent,ou=Container,dc=suffix`}
	mtg.LDIF()

	subentries = &sents
	subentries.Push(subentry)

	subentries.IsZero()
	subentries.Get(`hello`)
	subentries.Get(`spatialContext`)
	subentries.Contains(`hello`)
	subentries.Index(21)
	subentries.IsZero()

}
