package oid

/*
oid.go contains all types and methods pertaining to the ASN.1
OBJECT IDENTIFIER type.
*/

import (
	"errors"
	"math/big"
	"strconv"
	"strings"
)

/*
DotNotation implements an unbounded ASN.1 OBJECT IDENTIFIER (tag 6),
which is convertible to both the [encoding/asn1.ObjectIdentifier] and
[crypto/x509.OID] types.

See the [DotNotation.IntSlice] and [DotNotation.Uint64Slice]
methods for details.
*/
type DotNotation []NumberForm

/*
ASN1Notation returns an instance of [ASN1Notation]
alongside an error following an attempt to craft a series of individual
[NameAndNumberForm] instances using the number forms derived from the
receiver instance and the input string nameForms value.

The input nameForms value MUST be of an identical length to the receiver
instance. A zero slice indicates the associated number form has no name
form. A non-zero instance is processed for X.680 name form validity.
*/
func (r DotNotation) ASN1Notation(nameForms ...string) (
	ASN1Notation,
	error,
) {
	var (
		L1  int = len(r)
		L2  int = len(nameForms)
		err error
		_oiv,
		oiv ASN1Notation
	)

	if L2 == 0 {
		err = errorOIDOIVBadNames
		return oiv, err
	} else if L1 != L2 {
		err = errorOIDOIVBadNamesLen
		return oiv, err
	}

	for i := 0; i < len(nameForms) && err == nil; i++ {
		nf := nameForms[i]
		if nf == "" {
			_oiv = append(_oiv, NameAndNumberForm{
				numberForm: r[i],
				ok:         true,
			})
		} else if !IsNameForm(nf) {
			err = errorOIDOIVBadNames
		} else {
			_oiv = append(_oiv, NameAndNumberForm{
				nameForm:   nf,
				numberForm: r[i],
				ok:         true,
			})
		}
	}

	if err == nil {
		oiv = _oiv
	}

	return oiv, err
}

/*
Type returns the string literal "dot".
*/
func (r DotNotation) Type() string { return `dot` }

/*
String returns the string representation of the receiver instance.
*/
func (r DotNotation) String() (s string) {
	if r.Valid() {
		var x []string = make([]string, len(r))
		for i := 0; i < len(r); i++ {
			x[i] = r[i].String()
		}

		s = strings.Join(x, `.`)
	}
	return
}

/*
Eq returns a Boolean value indicative of an equality match between
the receiver and input [DotNotation] instances.
*/
func (r DotNotation) Eq(o DotNotation) bool {
	var ok bool
	if ok = r.Len() == o.Len(); ok {
		// compare each NumberForm slice
		for i := 0; i < r.Len() && ok; i++ {
			ok = r[i].Eq(o[i])
		}
	}

	return ok
}

/*
Ancestry returns slices of [DotNotation] instances ordered as leaf
node first.

Empty slices of [DotNotation] are returned if the value within the
receiver is less than two (2) [NumberForm] values in length.

Note that, unlike [ASN1Notation], the root [DotNotation] is
never added to the return instance by itself, as any valid [DotNotation]
MUST have a root AND at least one child. This means that for a receiver
value of "1.3.6.1.4.1", the final ancestor present within the return value
will be "1.3".
*/
func (r DotNotation) Ancestry() (anc []DotNotation) {
	if r.Len() >= 2 {
		for i := r.Len(); i > 1; i-- {
			anc = append(anc, r[:i])
		}
	}

	return
}

/*
Tag returns the integer 6. Note that this method is present merely for
convenience, and is not added as a header to any encoding of the receiver
instance.
*/
func (_ DotNotation) Tag() int { return 6 }

/*
Len returns the integer length of the receiver instance.
*/
func (r DotNotation) Len() int { return len(r) }

/*
IsZero returns a Boolean value indicative of a nil or zero length receiver state.
*/
func (r DotNotation) IsZero() (is bool) {
	if is = &r == nil; !is {
		is = r.Len() == 0
	}
	return
}

/*
IsDotNotation returns a Boolean value indicative of a call to
[NewDotNotation] returning an nil error following processing
of input value x.
*/
func IsDotNotation(x string) bool {
	_, err := newDotNotationStr(x)
	return err == nil
}

func assertDotNotation(id any) (A DotNotation) {
	switch tv := id.(type) {
	case string:
		A, _ = NewDotNotation(tv)
	case DotNotation:
		if tv.Len() >= 0 {
			A = tv
		}
	}

	return
}

/*
NewDotNotation returns an instance of [DotNotation] alongside
an error following an attempt to marshal the variadic x inputs as an ASN.1
OBJECT IDENTIFIER.

Variadic input allows for slice mixtures of all of the following types,
with each treated as an individual [NumberForm] instance:

  - *[big.Int]
  - [NumberForm]
  - string
  - []byte
  - uint64
  - int64
  - int32
  - int

If a string primitive is the only input option, it will be treated as a
complete [DotNotation] (e.g.: "1.3.6.1"). A single input value that
is NOT a string returns an error, as [DotNotation] instances MUST
have two (2) or more number form arcs at any given time.

If an [DotNotation] is the only input option, it is checked for
validity and returned without further processing.

If a byte slice is the only input option, it will be treated as encoded
bytes and an attempt to decode it into the return instance shall be made.
This is essentially the same as doing the following:

	var other DotNotation
	err = other.Decode(<bytes>) // write decoded contents to new var
*/
func NewDotNotation(x ...any) (r DotNotation, err error) {
	var _d DotNotation = make(DotNotation, 0)

	if len(x) == 1 {
		if slice, ok := x[0].(string); ok {
			// single string input
			r, err = newDotNotationStr(slice)
			return
		} else if slice2, ok := x[0].(DotNotation); ok {
			// check OID as valid
			if !slice2.Valid() {
				err = errorOIDNil
			} else {
				r = slice2
			}
			return
		} else if encb, ok := x[0].([]byte); ok {
			// pre-encoded bytes shortcut
			if err = _d.Decode(encb); err == nil {
				r = _d
				return
			}
		} else {
			err = errorOIDMinLen
			return
		}
	}

	for i := 0; i < len(x) && err == nil; i++ {
		var nf NumberForm
		switch tv := x[i].(type) {
		case *big.Int, NumberForm, string, int32, int64, uint64, int:
			nf, err = NewNumberForm(tv)
		default:
			err = errorOIDBadType
		}

		_d = append(_d, nf)
		if _d.Len() == 2 {
			// run validity check for first two
			// arcs before proceeding any further.
			if !_d.Valid() {
				err = errorOIDBadFirstArcs
			}
		}
	}

	if err == nil {
		r = _d
	}

	return
}

/*
MustNewDotNotation returns an instance of [DotNotation] and
panics if [NewDotNotation] returned an error during processing
of x.
*/
func MustNewDotNotation(x ...any) DotNotation {
	b, err := NewDotNotation(x...)
	if err != nil {
		panic(err)
	}
	return b
}

func newDotNotationStr(s string) (DotNotation, error) {
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return nil, errorOIDMinLen
	}

	args := make([]any, len(parts))
	for i, p := range parts {
		args[i] = p
	}

	o, err := NewDotNotation(args...)
	if !o.Valid() {
		err = errorOIDNil
	}
	return o, err
}

/*
Encode returns byte slices alongside an error following an attempt to
encode the receiver instance.
*/
func (r DotNotation) Encode() ([]byte, error) {
	if !r.Valid() {
		return nil, errorOIDNil
	}

	first := r[0]
	second := r[1]

	// First-two-arc compression: (first * 40) + second
	var combined NumberForm

	if !second.big {
		// Fast path: first is ALWAYS native, second is native here
		combined = NumberForm{
			ok:     true,
			native: first.native*40 + second.native,
		}
	} else {
		// Slow path: second is big, so we must use big.Int math
		tmp := newBigInt(0).Mul(newBigInt(0).SetUint64(first.native), newBigInt(40))
		tmp.Add(tmp, second.bigInt)
		combined = NumberForm{
			ok:     true,
			big:    true,
			bigInt: tmp,
		}
	}

	// Encode first two arcs
	wire := vlqEncode(combined)

	// Encode remaining arcs
	for i := 2; i < r.Len(); i++ {
		wire = append(wire, vlqEncode(r[i])...)
	}

	return wire, nil
}

/*
Decode returns an error following an attempt to decode the input buf bytes
into the receiver instance. Any data present in the receiver instance will
be destroyed.
*/
func (r *DotNotation) Decode(buf []byte) error {
	if len(buf) < 1 {
		return errorOIDBadEnc
	}

	p := 0

	// Decode combined first+second arc as NumberForm
	combined, err := vlqDecode(buf, &p)
	if err != nil {
		return err
	}

	// Reject combined values 120..159 (these correspond to illegal root 3.x)
	if !combined.big {
		v := combined.native
		if v >= 120 && v < 160 {
			return errorOIDBadFirstArcs
		}
	}

	var firstNF, secondNF NumberForm

	if !combined.big {
		// Native path
		v := combined.native

		switch {
		case v < 40:
			// 0.x
			firstNF = NumberForm{ok: true, native: 0}
			secondNF = NumberForm{ok: true, native: v}
		case v < 80:
			// 1.(v-40)
			firstNF = NumberForm{ok: true, native: 1}
			secondNF = NumberForm{ok: true, native: v - 40}
		default:
			// 2.(v-80)
			firstNF = NumberForm{ok: true, native: 2}
			secondNF = NumberForm{ok: true, native: v - 80}
		}
	} else {
		// Big path: combined is big.Int
		// For big combined, it must be >= 80: first = 2, second = combined - 80
		tmp := newBigInt(0).Set(combined.bigInt)
		tmp.Sub(tmp, newBigInt(80))

		firstNF = NumberForm{ok: true, native: 2}
		secondNF = NumberForm{ok: true, big: true, bigInt: tmp}
	}

	arcs := make(DotNotation, 0, 4)
	arcs = append(arcs, firstNF, secondNF)

	// Remaining arcs
	for p < len(buf) && err == nil {
		var nf NumberForm
		if nf, err = vlqDecode(buf, &p); err == nil {
			arcs = append(arcs, nf)
		}
	}

	if err == nil {
		*r = arcs
	}

	return err
}

/*
IntSlice returns slices of integer values and an error. The integer values are based
upon the contents of the receiver. Note that if any single arc number overflows int,
a zero slice is returned.

Successful output can be cast as an instance of [encoding/asn1.ObjectIdentifier], if desired.
*/
func (r DotNotation) IntSlice() (slice []int, err error) {
	if r.IsZero() {
		err = errorOIDNil
		return
	} else if r.Len() < 2 {
		err = errorOIDMinLen
		return
	}

	var t []int
	for i := 0; i < len(r) && err == nil; i++ {
		var n int
		if n, err = strconv.Atoi(r[i].String()); err == nil {
			t = append(t, n)
		}
	}

	if len(t) > 0 && err == nil {
		slice = t[:]
	}

	return
}

/*
Uint64Slice returns slices of uint64 values and an error. The uint64
values are based upon the contents of the receiver.

Note that if any single arc number overflows uint64, a zero slice is
returned alongside an error.

Successful output can be cast as an instance of [crypto/x509.OID], if
desired.
*/
func (r DotNotation) Uint64Slice() (slice []uint64, err error) {
	if r.IsZero() {
		err = errorOIDNil
		return
	} else if r.Len() < 2 {
		err = errorOIDMinLen
		return
	}

	var t []uint64
	for i := 0; i < len(r) && err == nil; i++ {
		var n uint64
		if n, err = strconv.ParseUint(r[i].String(), 10, 64); err == nil {
			t = append(t, n)
		}
	}

	if len(t) > 0 && err == nil {
		slice = t[:]
	}

	return
}

/*
Index returns the Nth index from the receiver, alongside a Boolean
value indicative of success.

This method supports the use of negative indices. For example, an
input index of -1 returns the final (far right) [NumberForm].
*/
func (r DotNotation) Index(idx int) (a NumberForm, ok bool) {
	if L := len(r); L > 0 {
		if idx < 0 {
			a = r[0]
			if x := L + idx; x >= 0 {
				a = r[x]
			}
		} else if idx > L {
			a = r[L-1]
		} else if idx < L {
			a = r[idx]
		}
		ok = !a.IsZero()
	}

	return
}

/*
Valid returns a Boolean value indicative of the following:

  - Receiver's length is greater than or equal to two (2) slice members, and ...
  - The first slice in the receiver contains an unsigned decimal value that is less than three (3), and ...
  - If root arc is 0 or 1, the second arc must be no greater than thirty nine (39)
*/
func (r DotNotation) Valid() (is bool) {
	if L := r.Len(); L > 0 {
		if is = r[0].Lt(3) && L >= 2; is {
			for i := 1; i < L && is; i++ {
				if i == 1 && r[0].Lt(2) {
					if is = r[1].Lt(40); !is {
						break
					}
				}
				is = r[i].ok
			}
		}
	}

	return
}

/*
MustNewSibling returns an instance of [DotNotation] and panics if
[DotNotation.NewSibling] returned an error during processing of x.
*/
func (r DotNotation) MustNewSibling(x any) DotNotation {
	O, err := r.NewSibling(x)
	if err != nil {
		panic(err)
	}
	return O
}

/*
NewSibling returns a new instance of [DotNotation] based upon the
contents of the receiver as well as the input [NumberForm] sibling value.
The result is the last arc of the new OID being replaced with the input
[NumberForm] in the return instance.

See also [DotNotation.MustNewSibling].
*/
func (r DotNotation) NewSibling(nf any) (O DotNotation, err error) {
	if !r.Valid() {
		err = errorOIDMinLen
		return
	}

	var n NumberForm
	// Prepare the new leaf numberForm, or die trying.
	if n, err = NewNumberForm(nf); err == nil {
		L := r.Len()
		O = make(DotNotation, L, L)
		for i := 0; i < L-1; i++ {
			O[i] = r[i]
		}
		O[L-1] = n
	}

	return
}

/*
MustNewSubordinate returns an instance of [DotNotation] and panics if
[DotNotation.NewSubordinate] returned an error during processing of x.
*/
func (r DotNotation) MustNewSubordinate(x any) DotNotation {
	O, err := r.NewSubordinate(x)
	if err != nil {
		panic(err)
	}
	return O
}

/*
NewSubordinate returns a new instance of [DotNotation] based upon the
contents of the receiver as well as the input [NumberForm] subordinate value.
This creates a fully-qualified child [DotNotation] value of the receiver.

See also [DotNotation.MustNewSubordinate].
*/
func (r DotNotation) NewSubordinate(nf any) (O DotNotation, err error) {
	if !r.Valid() {
		err = errorOIDMinLen
		return
	}

	var n NumberForm
	// Prepare the new leaf numberForm, or die trying.
	if n, err = NewNumberForm(nf); err == nil {
		L := r.Len()
		alloc := L + 1
		O = make(DotNotation, alloc, alloc)
		for i := 0; i < L; i++ {
			O[i] = r[i]
		}
		O[alloc-1] = n
	}

	return
}

/*
AncestorOf returns a Boolean value indicative of whether the receiver is
an ancestor of the input value, which can be string or [DotNotation].
*/
func (r DotNotation) AncestorOf(id any) (anc bool) {
	if !r.IsZero() {
		if A := assertDotNotation(id); !A.IsZero() {
			if A.Len() > r.Len() {
				anc = r.matchOID(A, 0)
			}
		}
	}

	return
}

/*
ChildOf returns a Boolean value indicative of whether the receiver is
a direct superior (parent) of the input value, which can be string or
[DotNotation].
*/
func (r DotNotation) ChildOf(id any) (cof bool) {
	if !r.IsZero() {
		if A := assertDotNotation(id); !A.IsZero() {
			if A.Len()-1 == r.Len() {
				cof = r.matchOID(A, 0)
			}
		}
	}

	return
}

/*
SiblingOf returns a Boolean value indicative of whether the receiver is
a sibling of the input value, which can be string or [DotNotation].
*/
func (r DotNotation) SiblingOf(id any) (sof bool) {
	if !r.IsZero() {
		if A := assertDotNotation(id); !A.IsZero() {
			if A.Len() == r.Len() && !A.Leaf().Eq(r.Leaf()) {
				sof = r.matchOID(A, -1)
			}
		}
	}

	return
}

/*
Root returns the root node (0) [NumberForm] value from the receiver.
*/
func (r DotNotation) Root() NumberForm {
	x, _ := r.Index(0)
	return x
}

/*
Leaf returns the leaf node (-1) [NumberForm] value from the receiver.
*/
func (r DotNotation) Leaf() NumberForm {
	x, _ := r.Index(-1)
	return x
}

/*
Parent returns the leaf node's parent (-2) [NumberForm] value from the receiver.
*/
func (r DotNotation) Parent() NumberForm {
	x, _ := r.Index(-2)
	return x
}

func (r DotNotation) matchOID(oiv DotNotation, off int) (matched bool) {
	L := r.Len()
	ct := 0
	for i := 0; i < L; i++ {
		x, _ := r.Index(i)
		if x.Eq(oiv[i]) {
			ct++
		} else if off == -1 && L-1 == i {
			// sibling check should end in
			// a FAILED match for the final
			// arcs.
			ct++
		}
	}

	return ct == L
}

var (
	errorOIDMinLen         = errors.New("OBJECT IDENTIFIER: two (2) or more arcs required")
	errorOIDNil            = errors.New("OBJECT IDENTIFIER: nil or bogus instance")
	errorOIDBadType        = errors.New("OBJECT IDENTIFIER: unsupported input type")
	errorOIDBadFirstArcs   = errors.New("OBJECT IDENTIFIER: illegal first and/or second level arcs")
	errorOIDBadEnc         = errors.New("OBJECT IDENTIFIER: bad encoding")
	errorOIDOIVBadNames    = errors.New("OBJECT IDENTIFIER: no nameForms at input for OIV init")
	errorOIDOIVBadNamesLen = errors.New("OBJECT IDENTIFIER: nameForm count MUST be equal length for OIV init")
)
