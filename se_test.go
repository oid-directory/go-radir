package radir

import (
	"fmt"
	"testing"
)

func TestSubentry_codecov(t *testing.T) {
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
	subentry.X660()
	subentry.IsZero()
	subentry.Spatial()
	subentry.isEmpty()
	subentry.Structural()
	subentry.DN()
	subentry.ObjectClasses()
	subentry.ObjectClassesGetFunc(nil)
	subentry.Kind()
	subentry.r_DITProfile = &DITProfile{}
	subentry.SetDITProfile(&DITProfile{})
	subentry.DITProfile()
	subentry.SetDN(`n=1,ou=Registrations,o=rA`)
	subentry.SetDITProfile(myDedicatedProfile)
	subentry.TTL()
	subentry.TTLGetFunc(nil)
	subentry.SetTTL(5)
	subentry.CTTL()
	subentry.CTTLGetFunc(nil)
	subentry.SetCTTL(5)
	subentry.SetDN(`cn=spatialContext,n=6,n=3,n=1,ou=Registrations,o=rA`)
	subentry.LDIF()
	subentry.refreshObjectClasses()
	subentry.SetObjectClasses(`x660Context`)
	subentry.refreshObjectClasses()
	subentry.GoverningStructureRule()
	subentry.GoverningStructureRuleGetFunc(nil)
	subentry.DNGetFunc(nil)
	subentry.Unmarshal()
	subentry.Marshal(nil)
	subentry.Marshal(func(any) error {
		return nil
	})
	subentry.Marshal(func(any) error {
		return fmt.Errorf("ERROR")
	})

}
