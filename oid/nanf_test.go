package oid

import (
	"fmt"
	"testing"
)

func ExampleMustNewNameAndNumberForm() {
	nanf := MustNewNameAndNumberForm("enterprise(1)")
	fmt.Println(nanf)
	// Output: enterprise(1)
}

func ExampleNameAndNumberForm_NumberForm() {
	n := "dod(6)"
	nanf, err := NewNameAndNumberForm(n)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(nanf.NumberForm())
	// Output: 6
}

func ExampleNameAndNumberForm_NameForm() {
	n := "dod(6)"
	nanf, err := NewNameAndNumberForm(n)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(nanf.NameForm())
	// Output: dod
}

func ExampleIsNameAndNumberForm() {
	nanf := "dod(6)"
	fmt.Printf("NameAndNumberForm (%q) is valid: %t", nanf, IsNameAndNumberForm(nanf))
	// Output: NameAndNumberForm ("dod(6)") is valid: true
}

func ExampleIsNameForm_bogus() {
	name := "-bogus" // cannot begin with '-'
	fmt.Printf("NameForm (%q) is valid: %t", name, IsNameForm(name))
	// Output: NameForm ("-bogus") is valid: false
}

func TestNameForm_bogus(t *testing.T) {
	for idx, this := range []string{
		``,
		`A`,
		`#`,
		`a--bc`,
		`aBC-`,
		`-aBC`,
		`3com`,
	} {
		if is := IsNameForm(this); is {
			t.Errorf("%s [%d] failed handling value (%q): expected %t, got %t",
				t.Name(), idx, this, !is, is)
		}
	}
}

func ExampleNewNameAndNumberForm_singleString() {
	// Input single string containing a valid
	// nameAndNumberForm value.
	nanf, err := NewNameAndNumberForm(`example(999)`)
	if err == nil {
		fmt.Println(nanf)
	}
	// Output: example(999)
}

func ExampleNewNameAndNumberForm_alternative() {
	// Any of int, int32, int64, uint64, *big.Int, string or
	// NumberForm are permitted for the numberForm value (arg 0).
	// Supplying a valid numberForm is always required, regardless
	// of input type.
	//
	// The second argument (arg 1, if present) must be a string
	// that fully conforms to X.680 nameForm syntax requirements.
	// Supplying a nameForm is optional.
	nanf, err := NewNameAndNumberForm(999, `example`)
	if err == nil {
		fmt.Println(nanf)
	}
	// Output: example(999)
}

func TestNameAndNumberForm_codecov(t *testing.T) {
	bigOne := newBigInt(1)
	nfOne := NumberForm{native: 999, ok: true}
	for _, Any := range [][]any{
		{1}, {`1`}, {int32(1)}, {int64(1)}, {uint64(1)},
		{bigOne}, {nfOne}, {999, "example"},
		{int32(999), "example"}, {int64(999), "example"},
		{uint64(999), "example"}, {bigOne, "example"},
		{nfOne, "example"}, {`999`, `example`},
	} {
		_, _ = NewNameAndNumberForm(Any...)
	}
	_, _ = newNameAndNumberFormStr(`909`)

	// do some bogus calls
	bigMinusOne := newBigInt(1)
	for _, Any := range [][]any{
		{-1}, {int32(-1)}, {int64(-1)},
		{bigMinusOne},
		{-999, "Example"}, {int32(-999), "-example"},
		{int64(-999), "example-"},
		{int64(-999), []byte("example-")},
		{bigMinusOne, "ex--ample"},
		{bigMinusOne, []rune("ex--ample")},
		{},
	} {
		_, _ = NewNameAndNumberForm(Any...)
	}
	_, _ = newNameAndNumberFormStr(``)
	_, _ = newNameAndNumberFormStr(`-string`)
	_, _ = newNameAndNumberFormStr(`string-`)
	_, _ = newNameAndNumberFormStr(`-string-`)
	_, _ = newNameAndNumberFormStr(`str--ing`)
	_, _ = newNameAndNumberFormStr(`9_09`)
	_, _ = newNameAndNumberFormStr(`-example(999)`)
	_, _ = newNameAndNumberFormStr(`example(999`)
	_, _ = newNameAndNumberFormStr(`example999)`)
	_ = IsNameForm("exa#mple")
	_ = IsNameAndNumberForm("exam#ple(999)")
	_ = IsNameAndNumberForm("example(999))")
	_ = IsNameAndNumberForm("example(999")
	_ = IsNameAndNumberForm("example999)")
	_ = IsNameAndNumberForm("_909")
	_ = IsNameAndNumberForm("_909", true)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	bogus := NameAndNumberForm{}
	_ = bogus.NumberForm()
}

func TestMustNewAndNumberForm(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	_ = MustNewNameAndNumberForm(struct{}{})
}
