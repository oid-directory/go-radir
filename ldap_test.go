package radir

import (
	"fmt"
	"testing"
)

func ExampleAttributeSelector_All() {
	attr := AttributeSelector{}
	fmt.Println(attr.All())
	// Output: [* +]
}

func ExampleAttributeSelector_AllUser() {
	attr := AttributeSelector{}
	fmt.Println(attr.AllUser())
	// Output: [*]
}

func ExampleAttributeSelector_AllOper() {
	attr := AttributeSelector{}
	fmt.Println(attr.AllOper())
	// Output: [+]
}

func TestResolveAltType_codecov(t *testing.T) {

	instances := []authority{
		&FirstAuthority{r_alt_types: true},
		&FirstAuthority{r_alt_types: false},
		&CurrentAuthority{r_alt_types: true},
		&CurrentAuthority{r_alt_types: false},
		&Sponsor{r_alt_types: true},
		&Sponsor{r_alt_types: false},
	}

	for _, instance := range instances {
		for typ, choice := range []bool{true, false} {
			instance.CN()
			instance.CO()
			instance.L()
			instance.O()
			instance.C()
			instance.Fax()
			instance.Tel()
			instance.ST()
			instance.Mobile()
			instance.Street()
			instance.Email()
			instance.Title()
			instance.URI()
			instance.PostalAddress()
			instance.PostalCode()
			instance.POBox()

			for _, tag := range []string{
				`sponsorCommonName`,
				`currentAuthorityOrg`,
				`sponsorLocality`,
				`firstAuthorityFax`,
				`currentAuthorityState`,
				`firstAuthorityStartTime`,
				`currentAuthorityEndTime`, // fake!
			} {
				resolveAltType(tag, typ, choice)
			}
		}
	}
}

/*
This example demonstrates the means of composing a pre-allocation range
check search filter.
*/
func ExampleRangeCheckSearchFilter() {
	fmt.Println(RangeCheckSearchFilter("56521"))
	// Output: (|(&(n<=56521)(|(registrationRange>=56521)(registrationRange=-1)))(n=56521))
}

/*
This example demonstrates a convenient means of obtaining the needed
parameters for a search request which checks to ensure we're not
encroaching upon a ranged allocation.

In this context, zero (0) results would mean we are not encroaching upon
any ranges, thus is likely the desired outcome.
*/
func ExampleRangeCheckSearchRequest() {
	X := "56521"                                         // the intended number form
	P := "n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA" // parent of siblings to be checked

	sup, scope, typesOnly, filter, at := RangeCheckSearchRequest(X, P)

	fmt.Printf("Target DN: %s\n", sup)
	fmt.Printf("Search Scope: %d\n", scope)
	fmt.Printf("Types Only: %t\n", typesOnly)
	fmt.Printf("Filter: %s\n", filter)
	fmt.Printf("Selected Attribute(s): %v\n", at)
	// Output: Target DN: n=1,n=4,n=1,n=6,n=3,n=1,ou=Registrations,o=rA
	// Search Scope: 1
	// Types Only: true
	// Filter: (|(&(n<=56521)(|(registrationRange>=56521)(registrationRange=-1)))(n=56521))
	// Selected Attribute(s): [registrationRange]
}

func TestLDAP_codecov(t *testing.T) {
	toLDIF(nil)
	toLDIF(struct{}{})
	tokenizeDN("")
	tokenizeDN("o")
	tokenizeDN("o=example")
	tokenizeDN("gidNumber=5042+o=example")
	tokenizeDN("uid=jesse+gidNumber=5042,ou=People,o=example")
	tokenizeDN(`cn=acme\, co,ou=Organizations,dc=example,dc=com`)
}
