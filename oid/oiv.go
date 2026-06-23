package oid

/*
oiv.go contains all types and methods pertaining to the ASN.1
ASN1Notation notation.
*/

import (
	"errors"
	"math/big"
	"strconv"
	"strings"
)

/*
ASN1Notation implements a type containing any number of
[NameAndNumberForm] instances, e.g.:

	{iso(1) identified-organization(3)}
	{0 0 1}
	{iso(1) identified-organization(3) 6 1 4 1 56521}

Note that instances of this type are NOT intended to be encoded.
*/
type ASN1Notation []NameAndNumberForm

/*
DotNotation returns an instance of [DotNotation] alongside
an error following an attempt to extract all number form components
from the receiver instance to form an equivalent numeric OID. Note
that only minimal validity checking is performed.
*/
func (r ASN1Notation) DotNotation() (DotNotation, error) {
	var (
		oid,
		_oid DotNotation
		err error
	)

	L := len(r)

	if L < 2 {
		// OIV must be 2+ nanfs
		err = errorOIVBadCvt
		return oid, err
	}

	for i := 0; i < L && err == nil; i++ {
		if nanf := r[i]; !nanf.ok {
			err = errorOIVBadNaNF
		} else {
			_oid = append(_oid, nanf.numberForm)
		}
	}

	if err == nil {
		oid = _oid
	}

	return oid, err
}

/*
Type returns the string literal "oiv".
*/
func (r ASN1Notation) Type() string { return `oiv` }

/*
IsASN1Notation returns a Boolean value indicative of
whether the oiv input string value represents a valid instance
of [ASN1Notation], in that:

  - The whole value is encapsulated within curly braces, and ...
  - The number of arcs is non zero, and ...
  - Each arc is a [NumberForm] or [NameAndNumberForm] instance
*/
func IsASN1Notation(oiv string) bool {
	idxF := strings.IndexRune(oiv, '{')
	idxL := strings.IndexRune(oiv, '}')

	if idxF == -1 || idxL == -1 {
		// bogus brace encaps
		return false
	}

	o := condenseWHSP(oiv[idxF+1 : idxL]) // extra space/NL delims OK
	arcs := strings.Split(o, ` `)

	if len(arcs) == 1 && arcs[0] == "" {
		return false
	}

	for _, arc := range arcs {
		if !IsNumberForm(arc) && !IsNameAndNumberForm(arc) {
			return false
		}
	}

	return true
}

/*
Len returns the integer length of the receiver instance.
*/
func (r ASN1Notation) Len() int { return len(r) }

/*
Index returns the Nth [NameAndNumberForm] index from the receiver alongside
a Boolean value indicative of success. This method supports the use of negative
indices.
*/
func (r ASN1Notation) Index(idx int) (NameAndNumberForm, bool) {
	var (
		ok   bool
		nanf NameAndNumberForm
	)
	if L := len(r); L > 0 {
		if idx < 0 {
			nanf = r[0]
			if x := L + idx; x >= 0 {
				nanf = r[x]
			}
		} else if idx > L {
			nanf = r[L-1]
		} else if idx < L {
			nanf = r[idx]
		}
		ok = nanf.ok
	}

	return nanf, ok
}

/*
String returns the string representation of the receiver instance.
*/
func (r ASN1Notation) String() string {
	var x []string
	for i := 0; i < len(r); i++ {
		x = append(x, r[i].String())
	}
	return `{` + strings.Join(x, ` `) + `}`
}

/*
NewASN1Notation returns an instance of [ASN1Notation]
alongside an error following an attempt to marshal x. x may be a string
or []string instance.

See also [MustNewASN1Notation].
*/
func NewASN1Notation(x any) (ASN1Notation, error) {
	var (
		oiv ASN1Notation
		nfs []string
		err error
	)

	switch tv := x.(type) {
	case string:
		nfs = strings.Fields(condenseWHSP(strings.TrimRight(strings.TrimLeft(tv, `{`), `}`)))
	case []string:
		nfs = tv
	default:
		err = errorOIVBadType
		return oiv, err
	}

	// Ensure at least one valid arc resides within the slice type.
	if len(nfs) == 0 || len(nfs) == 1 && nfs[0] == "" {
		err = errorOIVNil
		return oiv, err
	}

	// If the input has a nameForm as the first element (with
	// NO number form), ensure it matches the known root names.
	// If so, replace this "shorthand" with the fully qualified
	// NameAndNumberForm. If unknown, we die here.
	if !strings.Contains(nfs[0], `(`) {
		nf, found := oivRoots[nfs[0]]
		if !found {
			err = errorOIVBadRoot
			return oiv, err
		}
		nfs[0] += `(` + strconv.Itoa(nf) + `)`
	}

	var _oiv ASN1Notation
	for i := 0; i < len(nfs) && err == nil; i++ {
		var nanf NameAndNumberForm
		if nanf, err = NewNameAndNumberForm(nfs[i]); err == nil {
			_oiv = append(_oiv, nanf)
		}
	}

	if err == nil {
		oiv = _oiv
	}
	return oiv, err
}

/*
MustNewASN1Notation returns an instance of [ASN1Notation]
and panics if [NewASN1Notation] returned an error during processing
of x.
*/
func MustNewASN1Notation(x any) ASN1Notation {
	b, err := NewASN1Notation(x)
	if err != nil {
		panic(err)
	}
	return b
}

/*
Ancestry returns slices of [ASN1Notation] instances ordered from leaf
node (first) to root node (last).

Empty slices of [ASN1Notation] are returned if the value within the
receiver is less than two (2) [NameAndNumberForm] or [NumberForm] values in length.
*/
func (r ASN1Notation) Ancestry() (anc []ASN1Notation) {
	if r.Len() >= 2 {
		for i := r.Len(); i > 0; i-- {
			anc = append(anc, r[:i])
		}
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r ASN1Notation) IsZero() bool {
	return &r == nil || r.Len() == 0
}

/*
Valid returns a Boolean instance following a validity check
upon the receiver instance. The criterion for validity are
as follows:

  - Must be one (1) or more [NameAndNumberForm] instances in length
  - Each [NameAndNumberForm] instance must be valid
*/
func (r ASN1Notation) Valid() bool {
	L := r.Len()
	if L == 0 {
		return false
	}

	for i := 0; i < L; i++ {
		if !r[i].ok {
			return false
		}
	}

	return true
}

/*
MustNewSibling returns an instance of [ASN1Notation] and panics if
[ASN1Notation.NewSibling] returned an error during processing of
nanf.
*/
func (r ASN1Notation) MustNewSibling(nanf any) ASN1Notation {
	A, err := r.NewSibling(nanf)
	if err != nil {
		panic(err)
	}
	return A
}

/*
NewSibling returns a new instance of [ASN1Notation] based upon the
contents of the receiver as well as the input [NameAndNumberForm] sibling
value. The result is the last being arc replaced with the input [NameAndNumberForm]
in the return instance.

The input nanf syntax is controlled through the underlying call to the
[NewNameAndNumberForm] or [MustNewNameAndNumberForm] constructor.

See also [ASN1Notation.MustNewSibling].
*/
func (r ASN1Notation) NewSibling(nanf ...any) (A ASN1Notation, err error) {
	if !r.Valid() {
		err = errorOIDMinLen
		return
	}

	var n NameAndNumberForm
	// Prepare the new leaf nameAndNumberForm, or die trying.
	if n, err = NewNameAndNumberForm(nanf...); err == nil {
		L := r.Len()
		A = make(ASN1Notation, L, L)
		for i := 0; i < L-1; i++ {
			A[i] = r[i]
		}
		A[L-1] = n
	}

	return
}

/*
MustNewSubordinate returns an instance of [ASN1Notation] and panics
if [NewASN1Notation] returned an error during processing of x.
*/
func (r ASN1Notation) MustNewSubordinate(nanf any) ASN1Notation {
	A, err := r.NewSubordinate(nanf)
	if err != nil {
		panic(err)
	}
	return A
}

/*
NewSubordinate returns a new instance of [ASN1Notation] alongside an error
following an attempt to create a new "child" instance of the "parent" receiver; that
is, an effective copy of the receiver, plus a new leaf instance of [NameAndNumberForm].
*/
func (r ASN1Notation) NewSubordinate(nanf any) (A ASN1Notation, err error) {
	if r.Len() > 0 {
		alloc := r.Len() + 1
		_A := make(ASN1Notation, alloc, alloc)

		var leaf NameAndNumberForm // the new child nanf or nf
		switch tv := nanf.(type) {
		case NumberForm:
			if !tv.Valid() {
				err = errorNFNil
				break
			}
			leaf = NameAndNumberForm{numberForm: tv, ok: true}
		case NameAndNumberForm:
			if !tv.ok {
				err = errorNaNFNil
				break
			}
			leaf = tv
		case string:
			if len(tv) > 0 && ('0' <= tv[0] && tv[0] <= '9') {
				var n NumberForm
				n, err = NewNumberForm(tv)
				leaf = NameAndNumberForm{numberForm: n, ok: err == nil}
			} else {
				var n NameAndNumberForm
				n, err = NewNameAndNumberForm(tv)
				leaf = n
			}
		case int, int32, int64, uint64, *big.Int:
			var n NumberForm
			n, err = NewNumberForm(tv)
			leaf = NameAndNumberForm{numberForm: n}
		default:
			err = errorOIVBadType
		}

		if err == nil {
			for i := 0; i < r.Len(); i++ {
				_A[i] = r[i]
			}
			_A[alloc-1] = leaf
			A = _A
		}
	}

	return A, err
}

/*
AncestorOf returns a Boolean value indicative of whether the receiver
is an ancestor of the input value, which can be string or [ASN1Notation].
*/
func (r ASN1Notation) AncestorOf(oiv any) (anc bool) {
	if !r.IsZero() {
		if A := assertASN1Notation(oiv); !A.IsZero() {
			if A.Len() > r.Len() {
				anc = r.matchOIV(A, 0)
			}
		}
	}

	return
}

/*
ChildOf returns a Boolean value indicative of whether the receiver is
a direct superior (parent) of the input value, which can be string or
[ASN1Notation].
*/
func (r ASN1Notation) ChildOf(oiv any) (cof bool) {
	if !r.IsZero() {
		if A := assertASN1Notation(oiv); !A.IsZero() {
			if A.Len()-1 == r.Len() {
				cof = r.matchOIV(A, 0)
			}
		}
	}

	return
}

/*
SiblingOf returns a Boolean value indicative of whether the receiver is
a sibling of the input value, which can be string or [ASN1Notation].
*/
func (r ASN1Notation) SiblingOf(oiv any) (sof bool) {
	if !r.IsZero() {
		if A := assertASN1Notation(oiv); !A.IsZero() {
			if A.Len() == r.Len() && !A.Leaf().Eq(r.Leaf()) {
				sof = r.matchOIV(A, -1)
			}
		}
	}

	return
}

/*
Root returns the root node (0) string value from the receiver.
*/
func (r ASN1Notation) Root() NameAndNumberForm {
	x, _ := r.Index(0)
	return x
}

/*
Leaf returns the leaf node (-1) string value from the receiver.
*/
func (r ASN1Notation) Leaf() NameAndNumberForm {
	x, _ := r.Index(-1)
	return x
}

/*
Parent returns the leaf node's parent (-2) string value from the receiver.
*/
func (r ASN1Notation) Parent() NameAndNumberForm {
	x, _ := r.Index(-2)
	return x
}

func assertASN1Notation(oiv any) (A ASN1Notation) {
	switch tv := oiv.(type) {
	case string:
		A, _ = NewASN1Notation(tv)
	case ASN1Notation:
		if tv.Len() >= 0 {
			A = tv
		}
	}

	return
}

func (r ASN1Notation) matchOIV(oiv ASN1Notation, off int) (matched bool) {
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

var oivRoots = map[string]int{
	`itu-t`:           0,
	`itu-r`:           0,
	`ccitt`:           0,
	`iso`:             1,
	`joint-iso-itu-t`: 2,
	`joint-iso-ccitt`: 2,
}

var (
	errorOIVBadNaNF = errors.New("OBJECT IDENTIFIER VALUE: invalid NameAndNumberForm instance")
	errorOIVBadCvt  = errors.New("OBJECT IDENTIFIER VALUE: conversion requires two (2) or more arcs")
	errorOIVNil     = errors.New("OBJECT IDENTIFIER VALUE: nil or bogus instance")
	errorOIVBadType = errors.New("OBJECT IDENTIFIER VALUE: unsupported input type")
	errorOIVBadRoot = errors.New("OBJECT IDENTIFIER VALUE: bad root")
)
