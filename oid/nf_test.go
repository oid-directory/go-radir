package oid

import (
	"fmt"
	"math/big"
	"testing"
)

var (
	testBigPlusOne  *big.Int   = newBigInt(1)
	testBigPlusHuge *big.Int   = newBigInt(0)
	testBigNegOne   *big.Int   = newBigInt(-1)
	testBigNegHuge  *big.Int   = newBigInt(0)
	testNFPlusOne   NumberForm = MustNewNumberForm(1)
	testNFPlusHuge  NumberForm = MustNewNumberForm(`58742589373894573475934758934758934759347534789`)
)

func ExampleNumberForm_Encode() {
	nf, err := NewNumberForm(56521)
	if err != nil {
		fmt.Println(err)
		return
	}

	var enc []byte
	if enc, err = nf.Encode(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%#v\n", enc)
	// Output: []byte{0x83, 0xb9, 0x49}
}

func ExampleNumberForm_Decode() {
	// pre-encoded number form value 56521
	enc := []byte{0x83, 0xb9, 0x49}

	var nf NumberForm
	if err := nf.Decode(enc); err == nil {
		fmt.Println(nf)
	}
	// Output: 56521
}

func ExampleIsNumberForm_bogus() {
	nf := "06" // bogus (octal)
	fmt.Printf("NumberForm (%q) is valid: %t", nf, IsNumberForm(nf))
	// Output: NumberForm ("06") is valid: false
}

func TestNumberForm_valid(t *testing.T) {
	for _, nf := range []any{
		0, 1, 2,
		int32(0), int32(1), int32(2),
		int64(0), int64(1), int64(2),
		uint64(0), uint64(1), uint64(2),
		testBigPlusOne, testBigPlusHuge,
		testNFPlusOne, testNFPlusHuge,
		`58742589373894573475934758934758934759347534789`,
	} {
		nf, err := NewNumberForm(nf)
		if err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
		}
		_ = nf.IsBig()
		_ = nf.Native()
		_ = nf.Big()
		_ = nf.Bytes()
		_ = nf.Eq(1)
		_ = nf.Ne(1)
		_ = nf.Le(1)
		_ = nf.Ge(1)
		_ = nf.Lt(1)
		_ = nf.Gt(1)
		_ = nf.cmpAny(`1`)
		_ = nf.cmpAny(1)
		_ = nf.cmpAny([]byte{0x02})
		_ = nf.cmpAny(int32(1))
		_ = nf.cmpAny(int64(1))
		_ = nf.cmpAny(uint64(1))
		_ = nf.cmpAny(testBigPlusHuge)
		_ = nf.cmpAny(testBigNegHuge)
		_ = nf.cmpAny(NumberForm{ok: true, native: 1})
	}
}

func TestNumberForm_invalid(t *testing.T) {
	for idx, nf := range []any{
		-1, -2,
		int32(-1), int64(-1), int64(-2),
		testBigNegOne, testBigNegHuge,
		`-58742589373894573475934758934758934759347534789`,
		``, NumberForm{}, nil, []byte{},
	} {
		if _, err := NewNumberForm(nf); err == nil {
			t.Errorf("%s [%d] failed: expected error, got nil", t.Name(), idx)
		}
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	var nf NumberForm
	_ = nf.cmpAny([]byte{})
}

func TestNumberForm_panicCmp(t *testing.T) {
	var nf NumberForm
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	_ = nf.cmpAny(struct{}{})
}

func TestNumberForm_panicStr(t *testing.T) {
	var nf NumberForm
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	_ = nf.cmpNumberFormStr(`-`)
}

func TestMustNewNumberForm(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	_ = MustNewNumberForm(struct{}{})
}

func TestNumberForm_codecCoverage(t *testing.T) {
	p := 7
	_, _ = vlqDecodeBig([]byte{0x01}, &p, 0)
	_ = vlqEncodeBig(&big.Int{})
	_ = vlqEncode(NumberForm{})
	_ = uint64ToNumberForm(uint64(9223372036854775808))
	_ = bEToNumberForm([]byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01})
	_ = bEFitsUint64([]byte{0x01, 0x01, 0x01, 0x01})

	var nf NumberForm
	_, _ = nf.Encode()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s failed: %v", t.Name(), "expected panic")
		}
	}()
	_ = bEToUint64([]byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01})
}

func init() {
	_, _ = testBigPlusHuge.SetString(`58742589373894573475934758934758934759347534789`, 10)
	_, _ = testBigNegHuge.SetString(`-58742589373894573475934758934758934759347534789`, 10)
}
