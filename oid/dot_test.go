package oid

import (
	"fmt"
	"testing"
)

func ExampleIsDotNotation() {
	fmt.Println(IsDotNotation(`1.3.6`))
	// Output: true
}

func ExampleDotNotation_IntSlice() {
	id := `1.3.6.1.4.1.56521.999`
	o, err := NewDotNotation(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	var slice []int
	if slice, err = o.IntSlice(); err == nil {
		fmt.Println(slice)
	}
	// Output: [1 3 6 1 4 1 56521 999]
}

func ExampleDotNotation_Uint64Slice() {
	id := `1.3.6.1.4.1.56521.999`
	o, err := NewDotNotation(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	var slice []uint64
	if slice, err = o.Uint64Slice(); err == nil {
		fmt.Println(slice)
	}
	// Output: [1 3 6 1 4 1 56521 999]
}

func ExampleDotNotation_MustNewSibling() {
	id := `1.3.6.1.4.1.56521.999`
	o, err := NewDotNotation(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	sib := o.MustNewSibling(101)
	fmt.Println(sib)
	// Output: 1.3.6.1.4.1.56521.101
}

func ExampleDotNotation_MustNewSubordinate() {
	id := `1.3.6.1.4.1.56521.101`
	o, err := NewDotNotation(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	sub := o.MustNewSubordinate(2)
	fmt.Println(sub)
	// Output: 1.3.6.1.4.1.56521.101.2
}

func ExampleDotNotation_Eq() {
	one, err := NewDotNotation(`1.3.6.1.5.1`)
	if err != nil {
		fmt.Println(err)
		return
	}
	two, err := NewDotNotation(`1.3.6.1.5.2`)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("OIDs match: %t", one.Eq(two))
	// Output: OIDs match: false
}

func ExampleDotNotation_ASN1Notation() {
	o, err := NewDotNotation(`1.3.6.1.4.1`)
	if err != nil {
		fmt.Println(err)
		return
	}

	var oiv ASN1Notation
	oiv, err = o.ASN1Notation(`iso`,
		`identified-organization`, `dod`,
		`internet`, `private`, `enterprise`)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(oiv)
	// Output: {iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1)}
}

func TestDotNotation_Ancestry(t *testing.T) {
	want := `1.3.6.1.4.1.56521`
	leaf, err := NewDotNotation(want)
	if err != nil {
		t.Errorf("%s parse failed: %v", t.Name(), err)
	}

	ancestry := leaf.Ancestry()
	if L := len(ancestry); L != 6 {
		t.Errorf("%s count failed: want: 6, got: %d", t.Name(), L)
	} else if got := ancestry[0].String(); got != want {
		t.Errorf("%s compare failed:\n\tWant: %s\n\tGot:  %s", t.Name(), want, got)
	}

	if !leaf[0].Eq(ancestry[len(ancestry)-1][0]) {
		t.Errorf("%s root compare failed:\n\tWant: %s\n\tGot:  %s",
			t.Name(), leaf[0], ancestry[len(ancestry)-1][0])
	}
}

func ExampleDotNotation_NewSibling() {
	want := `1.3.6.1.4.1.56521`
	leaf, err := NewDotNotation(want)
	if err != nil {
		fmt.Println(err)
		return
	}

	var sibling DotNotation
	if sibling, err = leaf.NewSibling(994734); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(sibling)
	// Output: 1.3.6.1.4.1.994734
}

func ExampleDotNotation_Ancestry() {
	want := `1.3.6.1.4.1.56521`
	leaf, err := NewDotNotation(want)
	if err != nil {
		fmt.Println(err)
		return
	}

	ancestry := leaf.Ancestry()
	for i := 0; i < len(ancestry); i++ {
		fmt.Printf("%s\n", ancestry[i])
	}
	// Output:
	// 1.3.6.1.4.1.56521
	// 1.3.6.1.4.1
	// 1.3.6.1.4
	// 1.3.6.1
	// 1.3.6
	// 1.3
}

func ExampleDotNotation_Encode() {
	var orig string = `1.3.6.1.4.1`
	o, err := NewDotNotation(orig)
	if err != nil {
		fmt.Println(err)
		return
	}

	var enc []byte
	if enc, err = o.Encode(); err == nil {
		fmt.Printf("%#v\n", enc)
	}
	// Output: []byte{0x2b, 0x6, 0x1, 0x4, 0x1}
}

func ExampleDotNotation_Decode() {
	// pre-encoded OID bytes for 1.3.6.1.4.1
	enc := []byte{0x2b, 0x6, 0x1, 0x4, 0x1}

	var dest DotNotation
	if err := dest.Decode(enc); err == nil {
		fmt.Println(dest.String())
	}
	// Output: 1.3.6.1.4.1
}

func ExampleDotNotation_Decode_bogus() {
	// pre-encoded OID bytes for (illegal) 3.3.6.1.4.1
	enc := []byte{0x7B, 0x6, 0x1, 0x4, 0x1}

	var dest DotNotation
	fmt.Println(dest.Decode(enc))
	// Output: OBJECT IDENTIFIER: illegal first and/or second level arcs
}

func TestDotNotation(t *testing.T) {
	for _, oid := range []string{
		`1.3.6.1.4.1`,
		`2.5`,
		`0.0.4`,
	} {
		if _, err := NewDotNotation(oid); err != nil {
			t.Errorf("%s failed: genuine OID flagged as bogus (%s)", t.Name(), oid)
			return
		}
	}

	// test pre-encoded byte slice decoding
	// via constructor shortcut.
	dec, err := NewDotNotation([]byte{0x2b, 0x06, 0x01, 0x04, 0x01})
	if err != nil {
		t.Errorf("%s failed: decoding of pre-encoded bytes encountered an error: %v", t.Name(), err)
	} else if raw := dec.String(); raw != `1.3.6.1.4.1` {
		t.Errorf("%s failed: want 1.3.6.1.4.1, got %s", t.Name(), raw)
	}

	for _, bad := range []string{
		`2`,
		``,
		`$.3`,
		`1.2.3..4.5`,
		`1.2.3.4.5.`,
		`4.2.3.t.5`,
		`3.1`,
		`_`,
		`1.50`,
		`1.S0`,
		`.1.3`,
		`1.3.`,
	} {
		if _, err := NewDotNotation(bad); err == nil {
			t.Errorf("%s failed: bogus OID flagged as genuine (%s)", t.Name(), bad)
			return
		}
	}
}

func TestDotNotation_basicAny(t *testing.T) {
	bigPen := newBigInt(int64(56521))
	nfPen := NumberForm{
		ok:     true,
		native: uint64(56521),
	}

	for idx, this := range [][]any{
		{0, 0},                              // ITU-T Recommendation
		{1, 3, int64(6)},                    // DoD
		{1, 3, 6, 1, 4, 1, int32(56521)},    // My PEN OID
		{1, 3, 6, 1, 4, 1, bigPen},          // "" (big.Int)
		{1, 3, 6, 1, 4, 1, nfPen},           // "" (NumberForm)
		{1, 3, 6, 1, 4, 1, `56521`},         // "" (string)
		{1, 3, 6, 1, 5, uint64(5), 7, 3, 1}, // TLS ServerAuth
	} {
		oid, err := NewDotNotation(this...)
		if err != nil {
			t.Errorf("%s [%d] failed: %v", t.Name(), idx, err)
		}

		var enc []byte
		if enc, err = oid.Encode(); err != nil {
			t.Errorf("%s [%d] encode failed: %v", t.Name(), idx, err)
		}

		var dest DotNotation
		if err = dest.Decode(enc); err != nil {
			t.Errorf("%s [%d] decode failed: %v", t.Name(), idx, err)
		}

		// Use strings only for simple comparison
		want := oid.String()
		if got := dest.String(); got != want {
			t.Errorf("%s [%d] compare failed: want %s, got %s", t.Name(), idx, want, got)
		}
	}
}

func TestDotNotation_basicString(t *testing.T) {
	for _, want := range []string{
		`0.0`,               // ITU-T recommendation
		`1.3.6`,             // DoD
		`1.3.6.1.4.1.56521`, // My PEN OID
		`1.3.6.1.5.5.7.3.1`, // TLS ServerAuth
	} {
		oid, err := NewDotNotation(want)
		if err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
		}

		var enc []byte
		if enc, err = oid.Encode(); err != nil {
			t.Errorf("%s encode failed: %v", t.Name(), err)
		}

		var dest DotNotation
		if err = dest.Decode(enc); err != nil {
			t.Errorf("%s decode failed: %v", t.Name(), err)
		}

		if got := dest.String(); got != want {
			t.Errorf("%s compare failed: want %s, got %s", t.Name(), want, got)
		}
	}
}

func TestDotNotation_bigString(t *testing.T) {
	for _, want := range []string{
		`2.25.987895962269883002155146617097157934`,
		`2.25.923487482374589327592759723372975730`,
		`2.999999999999999999999999999999999999`, // not a real OID, but valid
	} {
		oid, err := NewDotNotation(want)
		if err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
		}

		var enc []byte
		if enc, err = oid.Encode(); err != nil {
			t.Errorf("%s encode failed: %v", t.Name(), err)
		}

		var dest DotNotation
		if err = dest.Decode(enc); err != nil {
			t.Errorf("%s decode failed: %v", t.Name(), err)
		}

		if got := dest.String(); got != want {
			t.Errorf("%s compare failed: want %s, got %s", t.Name(), want, got)
		}
	}
}

func TestDotNotation_NewSubordinate(t *testing.T) {
	id := `1.3.6.1.4.1`
	parent, err := NewDotNotation(id)
	if err != nil {
		t.Errorf("%s parse failed: %v", t.Name(), err)
	}
	if !parent.Valid() {
		t.Errorf("%s parent init failed: syntax/parser error", t.Name())
	}

	for idx, leaf := range []any{
		`56521`,
		56521,
		int32(56521),
		int64(56521),
		uint64(56521),
		newBigInt(56521),
		NumberForm{ok: true, native: uint64(56521)},
	} {
		want := id + ".56521"
		var child DotNotation
		if child, err = parent.NewSubordinate(leaf); err != nil {
			t.Errorf("%s [%d](%T)new subordinate failed: %v",
				t.Name(), idx, leaf, err)
		} else if !child.Valid() {
			t.Logf("%#v\n", child)
			t.Errorf("%s [%d](%T) child init failed: syntax/parser error",
				t.Name(), idx, leaf)
		} else if got := child.String(); got != want {
			t.Errorf("%s [%d](%T) compare failed:\n\tWant: %s\n\tGot:  %s",
				t.Name(), idx, leaf, want, got)
		}
	}
}

func TestDotNotation_codecov(t *testing.T) {
	var o DotNotation
	_ = o.Type()
	_, _ = o.Uint64Slice()
	_, _ = o.IntSlice()
	_ = o.IsZero()
	_, _ = o.ASN1Notation()
	_, _ = o.Encode()
	_ = o.Decode([]byte{})
	_ = o.Decode([]byte{0xff, 0xff})
	_ = o.Decode([]byte{0x81, 0xbe, 0xa1, 0xc2, 0x81, 0xc8, 0xc8, 0xc2, 0x8b, 0x8e, 0xd0, 0x80, 0x80, 0xaa, 0xae, 0xd7, 0xfa, 0x2e})
	o = MustNewDotNotation(`1.3.6.1.4.1.56521.101`)
	_ = o.Root()
	_ = o.Parent()
	_ = o.Leaf()
	_, _ = o.ASN1Notation(`iso`, `identified-organization`, `dod`, ``, ``, ``, ``, ``)
	_, _ = o.ASN1Notation(`iso`, `identified-organization`, `DoD`, ``, ``, ``, ``, ``)
	_, _ = o.ASN1Notation(`iso`, `identified-organization`)
	_, _ = o.ASN1Notation()
	_ = o.IsZero()
	_ = o.Tag()
	_ = o.Len()
	assertDotNotation(struct{}{})
	assertDotNotation(``)
	assertDotNotation(o)
	_, _ = NewDotNotation()
	_, _ = NewDotNotation(struct{}{})
	_, _ = NewDotNotation(struct{}{}, struct{}{})
	_, _ = NewDotNotation(o)
	_, _ = NewDotNotation(DotNotation{NumberForm{}})
	newDotNotationStr(`1.3.6._.4.1`)
	x := MustNewDotNotation(`1.3.6.1.4.1.56521.101.2`)
	y := MustNewDotNotation(`1.3.6.1.4.1.56521.999`)
	z := MustNewDotNotation(`1.3.6.1.4.1.56521.101.3`)
	_ = o.ChildOf(x)
	_ = o.SiblingOf(y)
	_ = o.AncestorOf(z)
	z.Index(10)
	z.Index(0)
	z.Index(-1)
	z.Index(-2)
	b := DotNotation{NumberForm{}}
	_, _ = b.Uint64Slice()
	_, _ = b.IntSlice()
}

func TestDotNotation_MustNewSibling(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	var x DotNotation
	x.MustNewSibling(2)
	id := `1.3.6.1.4.1.56521.999`
	o, _ := NewDotNotation(id)
	_ = o.MustNewSibling(struct{}{})
}

func TestDotNotation_MustNewSubordinate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	var x DotNotation
	x.MustNewSubordinate(2)
	id := `1.3.6.1.4.1.56521.999`
	o, _ := NewDotNotation(id)
	_ = o.MustNewSubordinate(struct{}{})
}

func TestMustNewDotNotation(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	_ = MustNewDotNotation(struct{}{})
}
