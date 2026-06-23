package radir

import "testing"

func TestSubtreeSpecification(t *testing.T) {
	// Verify parsing of valid string-based SubSpec values
	for idx, raw := range []string{
		`{minimum 1, maximum 1}`,
		`{base "n=1,n=4,n=1,n=6,n=3,n=1", minimum 1}`,
		`{specificExclusions { chopBefore "n=14", chopAfter "n=555", chopAfter "n=74,n=6" }, minimum 1, maximum 1}`,
		`{specificationFilter and:{item:1.3.6.1.4.1,or:{item:cn,item:2.5.4.7}}}`,
		`{base "n=1,n=4,n=1,n=6,n=3,n=1", minimum 1, maximum 1, specificationFilter and:{item:1.3.6.1.4.1,or:{item:cn,item:2.5.4.7}}}`,
		`{base "n=1,n=4,n=1,n=6,n=3,n=1", minimum 1, maximum 1, specificationFilter and:{item:1.3.6.1.4.1,not:item:1.3.6.1.5.5,or:{item:cn,item:2.5.4.7}}}`,
		`{base "n=1,n=4,n=1,n=6,n=3,n=1", minimum 1, maximum 1, specificationFilter item:1.3.6.1.4.1.56521}`,
		`{base "n=1,n=4,n=1,n=6,n=3,n=1", minimum 1, maximum 1, specificationFilter not:item:1.3.6.1.4.1.56521}`,
		`{}`,
	} {
		if v, err := NewSubtreeSpecification(raw); err != nil {
			t.Errorf("%s[%d] failed: %v", t.Name(), idx, err)
		} else if str := v.String(); str != raw {
			t.Errorf("%s[%d] failed:\nwant: %s\ngot%s", t.Name(), idx, raw, str)
		}
	}
}

func TestSubtreeSpecification_codecov(t *testing.T) {
	// largely focused on avoiding panics.

	subtreeBase(``)            // zero
	subtreeBase(`n=1,n=3,n=6`) // missing quotes
	subtreeBase(`"n=1,n=3,n=6"`)

	spexcl := SpecificExclusions{}
	_ = spexcl.String()

	parseRefinementList(`{and:{item:1.2.3.4,or:}{item:1.2.3.4}}}`)

	subtreeExclusions(`{}`, 0)
	subtreeExclusions(`chopBefore "n=1"`, 0)
	subtreeExclusions(`{chopBefore "n=1}"`, 0)
	subtreeExclusions(`{chopBeforee "n=1}"`, 0)

	deconstructExclusions(``, 0)
	deconstructExclusions(`chopBefore "n=1"`, 0)
	deconstructExclusions(`chopAfter "n=1"`, 0)

	_, _ = NewSubtreeSpecification(`{and:{item:1.2.3.4,or:}{item:1.2.3.4}}}`)
	_, _ = NewSubtreeSpecification(`{`)
	_, _ = NewSubtreeSpecification(`?{}`)
	_, _ = NewSubtreeSpecification(`{}`)
	_, _ = NewSubtreeSpecification(`{{}}`)
	_, _ = NewSubtreeSpecification(`{and:{item:1.2.3.4,or:}{item:1.2.3.4}}}`)
	_, _ = NewSubtreeSpecification(`{base __, specificationFilter and:{item:1.2.3.4,or:}{item:1.2.3.4}}}`)
	_, _ = NewSubtreeSpecification(`{specificExclusions {hello "test"}, specificationFilter and:{item:1.2.3.4,or:}{item:1.2.3.4}}}`)
	_, _ = NewSubtreeSpecification(`{minimum A}`)
	_, _ = NewSubtreeSpecification(`{maximum J}`)

}
