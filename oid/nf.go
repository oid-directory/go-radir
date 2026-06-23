package oid

/*
nf.go contains methods and types for expressing X.680 number forms.
*/

import (
	"errors"
	"math"
	"math/big"
	"strconv"
)

var newBigInt func(int64) *big.Int = big.NewInt

/*
NumberForm implements the unbounded ASN.1 INTEGER type for use in the
context of X.680 number forms, which are present within [NameAndNumberForm],
[DotNotation] and [ASN1Notation] type instances.

Note that *[big.Int] is used internally ONLY if the number overflows uint64.

For safety reasons (with respect to ambiguity of default values), a zero
instance of this type is bogus.  Users MUST use the [NewNumberForm] or
[MustNewNumberForm] constructor to obtain valid instances of this type.
*/
type NumberForm struct {
	big, ok bool
	native  uint64   // Stores native unsigned integer values when possible
	bigInt  *big.Int // Stores big.Int values only when necessary
}

/*
NewNumberForm returns an instance of [NumberForm] alongside an error
following an attempt to marshal x as an X.680 number form.

Input types may be int, int32, int64, uint64, string, []byte or
*[big.Int]. In the case of []byte, the value is expected to
be the Big Endian representation of the desired number form.

Any unsigned magnitude is permitted. Number forms which overflow
uint64 are stored as *[big.Int].

When the input value is NOT a string and when NO constraints are
utilized, it is safe to shadow the return error.

See also [MustNewNumberForm].
*/
func NewNumberForm[T any](x T) (i NumberForm, err error) {
	i, err = assertNumberForm(x)
	return
}

/*
MustNewNumberForm returns an instance of [NumberForm] and panics if [NewNumberForm]
returned an error during processing of x.
*/
func MustNewNumberForm[T any](x T) NumberForm {
	i, err := NewNumberForm(x)
	if err != nil {
		panic(err)
	}
	return i
}

func assertNumberForm[T any](v T) (i NumberForm, err error) {
	if err = checkNegativeNF(any(v)); err != nil {
		return
	}

	switch value := any(v).(type) {
	case int:
		i = NumberForm{native: uint64(value)}
	case int64:
		i = NumberForm{native: uint64(value)}
	case uint64:
		i = uint64ToNumberForm(value)
	case []byte:
		if value == nil || len(value) == 0 {
			err = errorNFBadBE
		}
		i = bEToNumberForm(value)
	case *big.Int:
		i = bigToNumberForm(value)
	case int32:
		i = NumberForm{native: uint64(value)}
	case string:
		i, err = strToNumberForm(value)
	case NumberForm:
		if !value.ok {
			err = errorNFNil
		}
		i = value
	default:
		err = errorNFBadType
	}

	if err == nil {
		i.ok = true
	}

	return
}

func checkNegativeNF(v any) (err error) {
	switch value := v.(type) {
	case int:
		if value < 0 {
			err = errorNFNegative
		}
	case int32:
		if value < 0 {
			err = errorNFNegative
		}
	case int64:
		if value < 0 {
			err = errorNFNegative
		}
	case *big.Int:
		if value.Cmp(newBigInt(0)) == -1 {
			err = errorNFNegative
		}
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r NumberForm) IsZero() bool { return &r == nil }

/*
String returns the string representation of the receiver instance.
*/
func (r NumberForm) String() string {
	var s string
	if r.big {
		s = r.bigInt.String()
	} else {
		s = strconv.FormatUint(r.native, 10)
	}

	return s
}

/*
IsBig returns a Boolean value indicative of the underlying value
overflowing uint64.
*/
func (r NumberForm) IsBig() bool { return r.big }

/*
Native returns the underlying int64 value found within the receiver
instance. Note that this method should not be used unless a call of
[NumberForm.IsBig] beforehand returns false.
*/
func (r NumberForm) Native() uint64 { return r.native }

/*
Valid returns a Boolean value indicative of the receiver instance
being properly initialized via the [NewNumberForm] or [MustNewNumberForm]
constructor with an unambiguous (non-default) value.
*/
func (r NumberForm) Valid() bool { return r.ok }

/*
Big returns the *[big.Int] form of the receiver instance.

Note that use of this method constructs an entirely new instance of
*[big.Int] if the underlying value is an int64.  Thus, this method
should only usually be needed if a call to [NumberForm.IsBig] returns
true. In that case, the preexisting *[big.Int] value is returned, as
opposed to being generated on the fly.

When [NumberForm.IsBig] returns false, the return instance of *[big.Int]
is entirely independent of the receiver and does not replace the
underlying value. This can be useful, though potentially costly, in
cases where methods extended by *[big.Int] that are not wrapped in
this package directly need to be accessed for some reason.
*/
func (r NumberForm) Big() (i *big.Int) {
	if r.big {
		i = r.bigInt
	} else {
		i = newBigInt(0).SetUint64(r.native)
	}

	return
}

/*
Bytes returns the receiver instance expressed as Big Endian bytes.
*/
func (r NumberForm) Bytes() []byte {
	var buf []byte
	if r.big {
		buf = r.bigInt.Bytes()
	} else {
		buf = uint64ToBE(r.native)
	}

	return buf
}

/*
Eq returns a bool indicative of an equality match between the
receiver instance and x.
*/
func (r NumberForm) Eq(x any) bool { return r.cmpAny(x) == 0 }

/*
Ne returns a bool indicative of a negative equality match between
the receiver instance and x.
*/
func (r NumberForm) Ne(x any) bool { return r.cmpAny(x) != 0 }

/*
Gt returns a bool indicative of r being greater than x.
*/
func (r NumberForm) Gt(x any) bool { return r.cmpAny(x) > 0 }

/*
Ge returns a bool indicative of r being greater than or equal to x.
*/
func (r NumberForm) Ge(x any) bool { return r.cmpAny(x) >= 0 }

/*
Lt returns a bool indicative of r being less than x.
*/
func (r NumberForm) Lt(x any) bool { return r.cmpAny(x) < 0 }

/*
Le returns a bool indicative of r being less than or equal to x.
*/
func (r NumberForm) Le(x any) bool { return r.cmpAny(x) <= 0 }

func (r NumberForm) cmpAny(x any) (result int) {
	switch t := x.(type) {
	case NumberForm:
		result = cmpNumberForm(r, t)

	case int:
		result = r.cmpInt64(int64(t))

	case int32:
		result = r.cmpInt64(int64(t))

	case int64:
		result = r.cmpInt64(t)

	case uint64:
		result = r.cmpUint64(t)

	case []byte:
		if t == nil || len(t) == 0 {
			panic(errorNFBadBE)
		}
		result = cmpNumberForm(r, bEToNumberForm(t))

	case string:
		result = r.cmpNumberFormStr(t)

	case *big.Int:
		result = r.cmpBig(t)

	default:
		panic("NumberForm: unsupported type for comparison")
	}

	return
}

func (r NumberForm) cmpNumberFormStr(v string) int {
	nf, err := NewNumberForm(v)
	if err != nil {
		panic(err)
	}
	return cmpNumberForm(r, nf)
}

func cmpNumberForm(a, b NumberForm) int {
	if !a.big && !b.big {
		switch {
		case a.native < b.native:
			return -1
		case a.native > b.native:
			return +1
		default:
			return 0
		}
	}
	return a.Big().Cmp(b.Big())
}

func (r NumberForm) cmpInt64(v int64) int {
	if !r.big {
		switch {
		case r.native < uint64(v):
			return -1
		case r.native > uint64(v):
			return +1
		default:
			return 0
		}
	}
	return r.Big().Cmp(big.NewInt(v))
}

func (r NumberForm) cmpUint64(u uint64) int {
	if !r.big && u <= math.MaxInt64 {
		return r.cmpInt64(int64(u))
	}
	b := newBigInt(0).SetUint64(u)
	return r.Big().Cmp(b)
}

func (r NumberForm) cmpBig(b *big.Int) int {
	if !r.big {
		return newBigInt(0).SetUint64(r.native).Cmp(b)
	}
	return r.bigInt.Cmp(b)
}

/*
Encode returns a byte slice alongside an error following an attempt to
encode the receiver instance as Big Endian bytes.
*/
func (r NumberForm) Encode() ([]byte, error) {
	if !r.ok {
		return nil, errorNFNil
	}
	enc := vlqEncode(r)
	return enc, nil
}

/*
Decode returns an error following an attempt to decode the input buf bytes
into the receiver instance. Any data present in the receiver instance will
be destroyed.
*/
func (r *NumberForm) Decode(buf []byte) error {
	p := 0
	n, err := vlqDecode(buf, &p)
	if err == nil {
		*r = n
	}
	return err
}

func bEToUint64(b []byte) uint64 {
	n := len(b)
	if n > 8 {
		panic("bigEndianToUint64: buffer length must be ≤ 8")
	}

	var u uint64
	for i := 0; i < 8-n; i++ {
		u = (u << 8) | 0x00
	}
	for _, by := range b {
		u = (u << 8) | uint64(by)
	}
	return u
}

func uint64ToBE(n uint64) []byte {
	b := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		b[i] = byte(n & 0xff)
		n >>= 8
	}
	return b
}

func bEFitsUint64(b []byte) bool {
	n := len(b)
	if n > 8 {
		for i := 0; i < n-8; i++ {
			if b[i] != 0x00 {
				return false
			}
		}
	}
	return true
}

func bEToNumberForm(b []byte) (i NumberForm) {
	if i.big = !bEFitsUint64(b); i.big {
		i.bigInt = newBigInt(0).SetBytes(b)
	} else {
		i.native = bEToUint64(b)
	}

	return
}

/*
IsNumberForm returns a Boolean value indicative of the nf string
input value representing a valid [NumberForm], in that:

  - The number is one (1) or more valid digits, and ...
  - The number is base10 (e.g.: not octal), and ...
  - The number is not negative

Assuming the above requirements are satisfied, any unsigned magnitude
is considered valid.
*/
func IsNumberForm(nf string) bool {
	return numberFormCheck(nf) == nil
}

func numberFormCheck(num string) (err error) {
	if len(num) == 0 {
		err = errorNFNoInput
		return
	}

	if num[0] == '-' {
		err = errorNFNegative
		return
	} else if len(num) > 1 && num[0] == '0' {
		err = errorNFOctal
		return
	}

	for i := 0; i < len(num); i++ {
		if ch := num[i]; !('0' <= ch && ch <= '9') {
			err = errorNFNaN
			break
		}
	}

	return
}

func strToNumberForm(num string) (i NumberForm, err error) {
	if err = numberFormCheck(num); err != nil {
		return
	}

	_i, _ := newBigInt(0).SetString(num, 10)
	if _i.IsUint64() {
		i = NumberForm{native: _i.Uint64()}
	} else {
		i = NumberForm{big: true, bigInt: _i}
	}

	return
}

func bigToNumberForm(num *big.Int) (i NumberForm) {
	if i.big = !num.IsUint64(); i.big {
		i.bigInt = num
	} else {
		i.native = num.Uint64()
	}

	return
}

func uint64ToNumberForm(num uint64) (i NumberForm) {
	if i.big = num > uint64(math.MaxInt64); i.big {
		i.bigInt = newBigInt(0).SetUint64(num)
	} else {
		i.native = num
	}

	return
}

func vlqEncode(nf NumberForm) []byte {
	if !nf.ok {
		return nil
	}

	if nf.big {
		return vlqEncodeBig(nf.bigInt)
	}
	return vlqEncodeUint64(nf.native)
}

func vlqEncodeUint64(n uint64) []byte {
	if n == 0 {
		return []byte{0}
	}

	var buf [16]byte
	i := len(buf)

	for n > 0 {
		i--
		b := byte(n & 0x7F) // take 7 bits
		n >>= 7

		if len(buf)-i > 1 { // set continuation bit except on last octet
			b |= 0x80
		}
		buf[i] = b
	}

	return buf[i:]
}

func vlqEncodeBig(n *big.Int) []byte {
	if n.Sign() == 0 {
		return []byte{0}
	}

	tmp := newBigInt(0).Set(n)
	rem := newBigInt(0)

	out := make([]byte, 0, 16)

	for tmp.Sign() != 0 {
		tmp.DivMod(tmp, newBigInt(128), rem)
		b := byte(rem.Uint64())
		out = append(out, b)
	}

	// reverse and set continuation bits
	for i := 0; i < len(out)/2; i++ {
		out[i], out[len(out)-1-i] = out[len(out)-1-i], out[i]
	}

	for i := 0; i < len(out)-1; i++ {
		out[i] |= 0x80
	}

	return out
}

// Top-level dispatcher: choose native or big based on magnitude.
func vlqDecode(buf []byte, p *int) (NumberForm, error) {
	u, done, err := vlqDecodeUint64(buf, p)
	if err != nil {
		return NumberForm{}, err
	}
	if done {
		return NumberForm{ok: true, native: u}, nil
	}
	return vlqDecodeBig(buf, p, u)
}

func vlqDecodeUint64(buf []byte, p *int) (uint64, bool, error) {
	var u uint64

	for {
		if *p >= len(buf) {
			return 0, false, errorNFBadVLQ
		}

		b := buf[*p]
		seven := uint64(b & 0x7F)

		// Would shifting overflow?
		if u > (^(uint64(0)) >> 7) {
			// Do NOT consume this byte; let big path handle it.
			return u, false, nil
		}

		// Safe to consume
		*p++
		u = (u << 7) | seven

		if b&0x80 == 0 {
			return u, true, nil
		}
	}
}

func vlqDecodeBig(buf []byte, p *int, prefix uint64) (NumberForm, error) {
	n := newBigInt(0).SetUint64(prefix)

	for {
		if *p >= len(buf) {
			return NumberForm{}, errorNFBadVLQ
		}

		b := buf[*p]
		*p++

		seven := int64(b & 0x7F)

		n.Lsh(n, 7)
		n.Or(n, newBigInt(seven))

		if b&0x80 == 0 {
			return NumberForm{ok: true, big: true, bigInt: n}, nil
		}
	}
}

var (
	errorNFNegative = errors.New("NUMBER FORM: negative numbers prohibited")
	errorNFNil      = errors.New("NUMBER FORM: nil or bogus instance")
	errorNFBadType  = errors.New("NUMBER FORM: unsupported input type")
	errorNFNoInput  = errors.New("NUMBER FORM: nil or zero input")
	errorNFOctal    = errors.New("NUMBER FORM: leading zeroes (octal numbers) prohibited")
	errorNFBadVLQ   = errors.New("NUMBER FORM: truncated VLQ")
	errorNFBadBE    = errors.New("NUMBER FORM: invalid BE input bytes")
	errorNFNaN      = errors.New("NUMBER FORM: non numeric character found")
)
