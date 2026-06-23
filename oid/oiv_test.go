package oid

import (
	"fmt"
	"testing"
)

func ExampleASN1Notation_MustNewSibling() {
	oiv := `{iso(1) identified-organization(3) 6 1 4 1 56521 example(999)}`
	o, err := NewASN1Notation(oiv)
	if err != nil {
		fmt.Println(err)
		return
	}

	sib := o.MustNewSibling("oid-directory(101)")
	fmt.Println(sib)
	// Output: {iso(1) identified-organization(3) 6 1 4 1 56521 oid-directory(101)}
}

func ExampleASN1Notation_MustNewSubordinate() {
	oiv := `{iso(1) identified-organization(3) 6 1 4 1 56521 oid-directory(101)}`
	o, err := NewASN1Notation(oiv)
	if err != nil {
		fmt.Println(err)
		return
	}

	sub := o.MustNewSubordinate("schema(2)")
	fmt.Println(sub)
	// Output: {iso(1) identified-organization(3) 6 1 4 1 56521 oid-directory(101) schema(2)}
}

func ExampleMustNewASN1Notation() {
	oiv := `{iso(1) identified-organization(3) 6 1 4 1 56521 oid-directory(101)}`
	o := MustNewASN1Notation(oiv)
	fmt.Println(o)
	// Output: {iso(1) identified-organization(3) 6 1 4 1 56521 oid-directory(101)}
}

func TestASN1Notation_MustNewSibling(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	oiv := `{iso(1) identified-organization(3) 6 1 4 1 56521 example(999)}`
	o, _ := NewASN1Notation(oiv)
	_ = o.MustNewSibling(struct{}{})
}

func TestASN1Notation_MustNewSubordinate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	oiv := `{iso(1) identified-organization(3) 6 1 4 1 56521 example(999)}`
	o, _ := NewASN1Notation(oiv)
	_ = o.MustNewSubordinate(struct{}{})
}

func TestMustNewASN1Notation(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	_ = MustNewASN1Notation(struct{}{})
}

func TestASN1Notation_codecov(t *testing.T) {
	var o ASN1Notation
	o.Type()
	o.IsZero()
	o.Len()
	o.Valid()
	o.NewSibling()
	o.NewSubordinate("barf(0)")

	_, _ = o.DotNotation()

	o = append(o, NameAndNumberForm{ok: true, numberForm: NumberForm{ok: true, native: 1}})
	_, _ = o.DotNotation()

	o[0] = NameAndNumberForm{}
	o = append(o, NameAndNumberForm{})
	o.Valid()
	o.NewSibling()
	o.NewSubordinate("barf(0)")
	_, _ = o.DotNotation()

	IsASN1Notation("")
	IsASN1Notation("{}")
	IsASN1Notation("{iso()}")

	o = ASN1Notation{NameAndNumberForm{numberForm: NumberForm{ok: true, native: 1}, ok: true}}
	o.Index(-1)
	o.Index(0)
	o.Index(1)
	o.Index(2)
	o.NewSubordinate("barf")
	o.NewSubordinate(0)
	o.NewSubordinate(NumberForm{})
	o.NewSubordinate(NameAndNumberForm{})
	o.NewSubordinate(int64(0))

	x, _ := NewASN1Notation([]string{"iso(1)", "identified-organization(3)"})
	y, _ := NewASN1Notation([]string{"iso(1)", "member-body(2)"})
	_, _ = NewASN1Notation(struct{}{})
	_, _ = NewASN1Notation(nil)
	_, _ = NewASN1Notation([]string{})
	_, _ = NewASN1Notation([]string{""})
	_, _ = NewASN1Notation("{barf}")
	_, _ = NewASN1Notation("{iso}")

	assertASN1Notation(``)
	o.Root()
	o.Parent()
	o.Leaf()
	o.AncestorOf(x)
	o.ChildOf(x)
	o.SiblingOf(x)
	o.SiblingOf(x[0])
	o.NewSubordinate(x[0])
	o.NewSubordinate(x[0].NumberForm())
	o.NewSubordinate(struct{}{})
	o.matchOIV(o, 0)
	o.matchOIV(x, -1)
	y.SiblingOf(x)
}

/*
This example demonstrates the means for creating a new [DotNotation] instance
using a preexisting [ASN1Notation] instance.
*/
func ExampleASN1Notation_DotNotation() {
	oiv := `{iso(1) identified-organization(3) 6 1 4 1 56521 example(999)}`
	o, err := NewASN1Notation(oiv)
	if err != nil {
		fmt.Println(err)
		return
	}

	var id DotNotation
	if id, err = o.DotNotation(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(id)
	// Output: 1.3.6.1.4.1.56521.999

}

func ExampleIsASN1Notation() {
	oiv := `{iso(1) identified-organization(3) 6 1 4 1 56521 example(999)}`
	fmt.Printf("ASN1Notation (%q) is valid: %t", oiv, IsASN1Notation(oiv))
	// Output: ASN1Notation ("{iso(1) identified-organization(3) 6 1 4 1 56521 example(999)}") is valid: true
}

func TestASN1Notation_Ancestry(t *testing.T) {
	want := `{iso(1) identified-organization(3) 6 1 4 1 56521}`
	leaf, err := NewASN1Notation(want)
	if err != nil {
		t.Errorf("%s parse failed: %v", t.Name(), err)
	}

	ancestry := leaf.Ancestry()
	if L := len(ancestry); L != 7 {
		t.Errorf("%s count failed: want: 7, got: %d", t.Name(), L)
	} else if got := ancestry[0].String(); got != want {
		t.Errorf("%s compare failed:\n\tWant: %s\n\tGot:  %s", t.Name(), want, got)
	}

	if !leaf[0].Eq(ancestry[len(ancestry)-1][0]) {
		t.Errorf("%s root compare failed:\n\tWant: %s\n\tGot:  %s",
			t.Name(), leaf[0], ancestry[len(ancestry)-1][0])
	}
}

func ExampleASN1Notation_Ancestry() {
	want := `{iso(1) identified-organization(3) 6 1 4 1 56521}`
	leaf, err := NewASN1Notation(want)
	if err != nil {
		fmt.Println(err)
		return
	}

	ancestry := leaf.Ancestry()
	for i := 0; i < len(ancestry); i++ {
		fmt.Printf("%s\n", ancestry[i])
	}
	// Output:
	// {iso(1) identified-organization(3) 6 1 4 1 56521}
	// {iso(1) identified-organization(3) 6 1 4 1}
	// {iso(1) identified-organization(3) 6 1 4}
	// {iso(1) identified-organization(3) 6 1}
	// {iso(1) identified-organization(3) 6}
	// {iso(1) identified-organization(3)}
	// {iso(1)}
}

func ExampleASN1Notation_NewSibling() {
	want := `{iso(1) identified-organization(3) 6 1 4 1 56521 oid-directory(101)}`
	leaf, err := NewASN1Notation(want)
	if err != nil {
		fmt.Println(err)
		return
	}

	var sibling ASN1Notation
	if sibling, err = leaf.NewSibling("example(999)"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(sibling)
	// Output: {iso(1) identified-organization(3) 6 1 4 1 56521 example(999)}
}

func TestASN1Notation_NewSubordinate(t *testing.T) {
	oiv := `iso(1) identified-organization(3) 6 1 4 1 56521`
	parent, err := NewASN1Notation("{" + oiv + "}")
	if err != nil {
		t.Errorf("%s parse failed: %v", t.Name(), err)
	}
	if !parent.Valid() {
		t.Errorf("%s parent init failed: syntax/parser error", t.Name())
	}

	for _, leaf := range []string{
		`example(999)`,
		`999`,
	} {
		want := "{" + oiv + " " + leaf + "}"
		var child ASN1Notation
		if child, err = parent.NewSubordinate(leaf); err != nil {
			t.Errorf("%s new subordinate failed: %v", t.Name(), err)
		} else if !child.Valid() {
			t.Errorf("%s child init failed: syntax/parser error", t.Name())
		} else if got := child.String(); got != want {
			t.Errorf("%s compare failed:\n\tWant: %s\n\tGot:  %s", t.Name(), want, got)
		}
	}
}
