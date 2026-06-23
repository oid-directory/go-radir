package oid

import (
	"fmt"
	"testing"
)

func ExampleIsUnicodeValue() {
	fmt.Println(IsUnicodeValue("Identified-Organization"))
	// Output: true
}

func ExampleIRINotation_Ancestry() {
	want := `/ISO/Identified-Organization/6/1/4/1/56521`
	leaf, err := NewIRINotation(want)
	if err != nil {
		fmt.Println(err)
		return
	}

	ancestry := leaf.Ancestry()
	for i := 0; i < len(ancestry); i++ {
		fmt.Printf("%s\n", ancestry[i])
	}
	// Output:
	// /ISO/Identified-Organization/6/1/4/1/56521
	// /ISO/Identified-Organization/6/1/4/1
	// /ISO/Identified-Organization/6/1/4
	// /ISO/Identified-Organization/6/1
	// /ISO/Identified-Organization/6
	// /ISO/Identified-Organization
	// /ISO
}

func TestIRINotation_codecov(t *testing.T) {
	var i IRINotation
	_ = i.Type()
	_ = i.IsZero()
	_, _ = NewIRINotation("")
	_, _ = NewIRINotation("/")
	_, _ = NewIRINotation("/1")
	_, _ = NewIRINotation(nil)
	_, _ = NewIRINotation(IRIArc{ok: true, number: NumberForm{ok: true, native: 1}})
	_, _ = NewIRINotation(IRIArc{ok: false})
	_, _ = NewIRINotation(IRIArc{ok: false}, IRIArc{ok: false})
	_, _ = NewIRINotation(IRIArc{ok: true, name: "ISO"}, IRIArc{ok: true, name: "Identified-Organization"})
	_, _ = NewIRINotation([]any{}...)
	_, _ = NewIRINotation([]any{""}...)
	_, _ = NewIRINotation("//`")
	_, _ = NewIRINotation("/fjeskjkfdl/`")
	_, _ = NewIRINotation("/ASN.1()`")
	_, _ = NewIRINotation(struct{}{})
	_, _ = NewIRINotation(struct{}{}, struct{}{})
	_ = IsUnicodeValue("")
	_ = IsUnicodeValue("bogus")
	_ = IsUnicodeValue("6")

	num := NumberForm{ok: true, native: 3}
	iri := IRINotation{IRIArc{ok: true, number: num}}
	_, _ = iri.Index(-1)
	_, _ = iri.Index(-100)
	_, _ = iri.Index(0)
	_, _ = iri.Index(7)
	_ = iri.Root()
	_ = iri.Parent()
	_ = iri.Leaf()

	o := MustNewIRINotation(`/ISO/Identified-Organization/6/1/4/1/56521`)
	x := MustNewIRINotation(`/ISO/Identified-Organization/6/1/4/1/56521/101.2`)
	y := MustNewIRINotation(`/ISO/Identified-Organization/6/1/4/1/56521/999`)
	z := MustNewIRINotation(`/ISO/Identified-Organization/6/1/4/1/56521/101.3`)
	_ = o.ChildOf(y)
	_ = x.SiblingOf(z)
	_ = o.AncestorOf(x)

	_ = assertIRINotation(`/ISO/Identified-Organization`)

	// 0x80 is an invalid standalone UTF-8 byte
	bad := string([]byte{0x80})
	if _, err := NewIRIArc(bad); err == nil {
		t.Fatalf("NewIRIArc(%q) should have failed due to invalid UTF-8", bad)
	}
	three, _ := NewIRIArc(3)
	three.Eq(2)
	_ = isIRIChar('/')
}

func ExampleIRINotation_MustNewSibling() {
	id := `/ISO/Member-Body`
	o, err := NewIRINotation(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	sib := o.MustNewSibling(`Identified-Organization`)
	fmt.Println(sib)
	// Output: /ISO/Identified-Organization
}

func ExampleIRINotation_MustNewSubordinate() {
	id := `/ISO/Identified-Organization`
	o, err := NewIRINotation(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	sub := o.MustNewSubordinate(6)
	fmt.Println(sub)
	// Output: /ISO/Identified-Organization/6
}

func ExampleIRINotation_NewSibling() {
	want := `/ISO/Identified-Organization`
	leaf, err := NewIRINotation(want)
	if err != nil {
		fmt.Println(err)
		return
	}

	var sibling IRINotation
	if sibling, err = leaf.NewSibling(`Member-Body`); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(sibling)
	// Output: /ISO/Member-Body
}

func TestIRINotation_MustNewSibling(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	var x IRINotation
	x.MustNewSibling(2)
	id := `1.3.6.1.4.1.56521.999`
	o, _ := NewIRINotation(id)
	_ = o.MustNewSibling(struct{}{})
}

func TestIRINotation_MustNewSubordinate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	var x IRINotation
	x.MustNewSubordinate(2)
	id := `1.3.6.1.4.1.56521.999`
	o, _ := NewIRINotation(id)
	_ = o.MustNewSubordinate(struct{}{})
}

func ExampleIRINotation_Eq() {
	one, err := NewIRINotation(`/ISO/Identified-Organization/6/1/5/4`)
	if err != nil {
		fmt.Println(err)
		return
	}
	two, err := NewIRINotation(`/ISO/Identified-Organization/6/1/5/5`)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("IRIs match: %t", one.Eq(two))
	// Output: IRIs match: false
}

func TestIRIArc_codecov(t *testing.T) {
	_, _ = NewIRIArc("")
	_, _ = NewIRIArc(struct{}{})
}

func TestNewIRIArc_ValidNames(t *testing.T) {
	tests := []string{
		"ISO",
		"Identified-Organization",
		"猫科",
		"Συνάρτηση",
		"пример",
	}

	for _, s := range tests {
		arc, err := NewIRIArc(s)
		if err != nil {
			t.Fatalf("NewIRIArc(%q) returned error: %v", s, err)
		}
		if !arc.Valid() {
			t.Fatalf("NewIRIArc(%q) produced invalid arc", s)
		}
		if arc.String() != s {
			t.Fatalf("NewIRIArc(%q).String() = %q, want %q", s, arc.String(), s)
		}
	}
}

func TestNewIRIArc_ValidNumbers(t *testing.T) {
	tests := []any{
		"0", int32(0), int64(0), uint64(0),
		"1",
		"999",
		"123456789012345678901234567890",
	}

	for _, s := range tests {
		arc, err := NewIRIArc(s)
		if err != nil {
			t.Fatalf("NewIRIArc(%q) returned error: %v", s, err)
		}
		if !arc.Valid() {
			t.Fatalf("NewIRIArc(%q) produced invalid arc", s)
		}
		if !arc.number.Eq(s) {
			t.Fatalf("NewIRIArc(%q).String() = %q, want %q", s, arc.String(), s)
		}
	}
}

func TestNewIRIArc_Invalid(t *testing.T) {
	tests := []string{
		"",         // empty
		"foo(3)",   // illegal in IRI
		"ISO(1)",   // illegal in IRI
		"bad char", // space not allowed
		"foo/bar",  // slash not allowed
		"\x00",     // control char
	}

	for _, s := range tests {
		if arc, err := NewIRIArc(s); err == nil {
			t.Fatalf("NewIRIArc(%q) should have failed, got arc=%v", s, arc)
		}
	}
}

func TestNewIRINotation_ValidString(t *testing.T) {
	s := "/ISO/Identified-Organization/6/1/4/1"

	iri, err := NewIRINotation(s)
	if err != nil {
		t.Fatalf("NewIRINotation(%q) returned error: %v", s, err)
	}

	if !iri.Valid() {
		t.Fatalf("NewIRINotation(%q) produced invalid IRI", s)
	}

	if iri.String() != s {
		t.Fatalf("IRI.String() = %q, want %q", iri.String(), s)
	}
}

func TestNewIRINotation_ValidSlice(t *testing.T) {
	arcs := []any{"ISO", "Identified-Organization", 6, int32(1), uint64(4), int64(1)}
	iri, err := NewIRINotation(arcs...)
	if err != nil {
		t.Fatalf("%s failed: returned error: %v", t.Name(), err)
	}

	if !iri.Valid() {
		t.Fatalf("%s failed: produced invalid IRI", t.Name())
	}

	if iri.String() != "/ISO/Identified-Organization/6/1/4/1" {
		t.Fatalf("%s failed: IRI.String() mismatch: %q", t.Name(), iri.String())
	}
}

func TestNewIRINotation_Invalid(t *testing.T) {
	tests := []string{
		"",            // empty
		"ISO/1/2",     // missing leading slash
		"/",           // no arcs
		"/foo//bar",   // empty arc
		"/foo/bar(3)", // illegal arc
	}

	for _, s := range tests {
		if iri, err := NewIRINotation(s); err == nil {
			t.Fatalf("NewIRINotation(%q) should have failed, got %v", s, iri)
		}
	}
}

func TestIRINotation_Index(t *testing.T) {
	iri, _ := NewIRINotation("/ISO/Identified-Organization/6")

	arc, ok := iri.Index(0)
	if !ok || arc.String() != "ISO" {
		t.Fatalf("Index(0) = %v, ok=%v", arc, ok)
	}

	arc, ok = iri.Index(-1)
	if !ok || arc.String() != "6" {
		t.Fatalf("Index(-1) = %v, ok=%v", arc, ok)
	}
}

func TestIRINotation_Valid(t *testing.T) {
	iri, _ := NewIRINotation("/ISO/Identified-Organization/6")
	if !iri.Valid() {
		t.Fatalf("Valid() returned false for valid IRI")
	}

	// Inject an invalid arc
	iri[1] = IRIArc{} // ok=false
	if iri.Valid() {
		t.Fatalf("Valid() returned true for invalid IRI")
	}
}

func ExampleMustNewIRINotation() {
	iri := MustNewIRINotation("/ITU-T/Recommendation")
	fmt.Println(iri)
	// Output: /ITU-T/Recommendation
}

func ExampleMustNewIRIArc() {
	iri := MustNewIRIArc("Recommendation")
	fmt.Println(iri)
	// Output: Recommendation
}

func TestMustNewIRINotation(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	_ = MustNewIRINotation(struct{}{})
}

func TestMustNewIRIArc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	_ = MustNewIRIArc(struct{}{})
}
