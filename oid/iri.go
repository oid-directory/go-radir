package oid

import (
	"errors"
	"math/big"
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
IsUnicodeValue returns a Boolean value indicative of the input
string name being a valid UCS string. Pure numbers, nameForms
(identifiers) and zero strings are not valid unicode values.
*/
func IsUnicodeValue(name string) bool {
	if len(name) == 0 || IsNameForm(name) || IsNumberForm(name) {
		return false
	}

	// ASCII-only lowercase alpha check
	first := rune(name[0])
	if 'a' <= first && first <= 'z' {
		return false
	}

	var (
		is = true
		i  = 0
	)

	for i < len(name) && is {
		r, size := utf8.DecodeRuneInString(name[i:])
		if r == utf8.RuneError && size == 1 {
			return false
		}
		is = isIRIChar(r)
		i += size
	}

	return is
}

/*
IRIArc implements the slice type of an [IRINotation] instance.
*/
type IRIArc struct {
	name   string
	number NumberForm
	ok     bool
}

/*
Valid returns a Boolean value indicative of the receiver instance
being a valid [IRIArc].
*/
func (r IRIArc) Valid() bool { return r.ok }

/*
String returns the string representation of the receiver instance.
*/
func (r IRIArc) String() string {
	var s string
	if r.Valid() {
		if len(r.name) > 0 {
			s = r.name
		} else {
			s = r.number.String()
		}
	}

	return s
}

/*
MustNewIRIArc returns an instance of [IRIArc] and panics if the underlying
call to [NewIRIArc] returns an error.
*/
func MustNewIRIArc(x any) IRIArc {
	i, err := NewIRIArc(x)
	if err != nil {
		panic(err)
	}
	return i
}

/*
NewIRIArc returns an instance of [IRIArc] alongside an error following an
attempt to parse input value x as the slice type of an [IRINotation].
*/
func NewIRIArc(x any) (IRIArc, error) {
	var (
		arc IRIArc
		err error
	)

	switch tv := x.(type) {
	case string:
		if tv == "" {
			return arc, errorIRINilArc
		}

		if strings.ContainsAny(tv, "()") {
			return arc, errorIRINilArc
		}

		isNumber := IsNumberForm(tv)

		// Validate Unicode characters
		if !IsUnicodeValue(tv) && !isNumber {
			return arc, errorIRIArcBadUCS
		}

		// Purely numeric?
		if isNumber {
			var nf NumberForm
			if nf, err = NewNumberForm(tv); err == nil {
				arc.number = nf
				arc.ok = true
			}
			break
		}

		// Otherwise it's a Unicode name
		arc.name = tv
		arc.ok = true

	case int, int32, int64, uint64, *big.Int, NumberForm:
		var nf NumberForm
		if nf, err = NewNumberForm(tv); err == nil {
			arc.number = nf
			arc.ok = true
		}

	default:
		err = errorIRIBadType
	}

	return arc, err
}

/*
IRINotation represents an OID-IRI, e.g. /ISO/Identified-Organization/6/1/4/1
*/
type IRINotation []IRIArc

/*
Len returns the integer length of the receiver.
*/
func (r IRINotation) Len() int { return len(r) }

/*
Index returns the Nth [IRIArc] value alongside a success-indicative Boolean
value following an attempt to read the indicated slice index from the receiver.

This method supports the use of negative indices.

See also the [IRINotation.Root], [IRINotation.Parent] and [IRINotation.Leaf]
convenience methods.
*/
func (r IRINotation) Index(idx int) (IRIArc, bool) {
	var (
		ok  bool
		arc IRIArc
	)

	if L := len(r); L > 0 {
		if idx < 0 {
			arc = r[0]
			if x := L + idx; x >= 0 {
				arc = r[x]
			}
		} else if idx > L {
			arc = r[L-1]
		} else if idx < L {
			arc = r[idx]
		}
		ok = arc.ok // found
	}

	return arc, ok
}

/*
Valid returns a Boolean value indicative of the receiver
instance being a valid [IRINotation].
*/
func (r IRINotation) Valid() bool {
	L := r.Len()
	var is bool
	if is = L > 0; is {
		for i := 0; i < L && is; i++ {
			is = r[i].Valid()
		}
	}

	return is
}

/*
String returns the string representation of the receiver instance.
*/
func (r IRINotation) String() string {
	var s []string
	if r.Valid() {
		for i := 0; i < r.Len(); i++ {
			s = append(s, r[i].String())
		}
	}
	return "/" + strings.Join(s, "/")
}

/*
IsIRINotation returns a Boolean value indicative of the input
string value s being a legal OID-IRI.
*/
func IsIRINotation(s string) bool {
	if len(s) == 0 || s[0] != '/' {
		return false
	}
	parts := strings.Split(s[1:], "/")
	if len(parts) == 1 && parts[0] == "" {
		return false
	}

	var is bool = true // assume OK by default
	for i := 0; i < len(parts) && is; i++ {
		is = IsUnicodeValue(parts[i]) || IsNumberForm(parts[i])
	}
	return is
}

/*
MustNewIRINotation returns an instance of [IRINotation] and panics if the underlying
call to [NewIRINotation] returns an error.
*/
func MustNewIRINotation(x any) IRINotation {
	i, err := NewIRINotation(x)
	if err != nil {
		panic(err)
	}
	return i
}

/*
NewIRINotation returns an instance of [IRINotation] alongside an
error following an attempt to parse input value x.
*/
func NewIRINotation(x ...any) (IRINotation, error) {
	var (
		iri IRINotation
		err error
	)

	switch len(x) {
	case 0:
		err = errorIRINoInput
	case 1:
		iri, err = assertIRINotationSingleArg(x[0])
	default:
		iri, err = assertIRINotationVariadic(x...)
	}

	return iri, err
}

func assertIRINotationSingleArg(x any) (iri IRINotation, err error) {
	// single string containing full IRI
	switch tv := x.(type) {
	case string:
		if len(tv) > 0 && IsNumberForm(tv[1:]) {
			// pure iri number (e.g.: /1)
			var nf NumberForm
			if nf, err = NewNumberForm(tv[1:]); err == nil {
				iri = append(iri, IRIArc{ok: true, number: nf})
				break
			}
		}

		if !IsIRINotation(tv) {
			return iri, errorIRINil
		}

		parts := strings.Split(tv[1:], "/")
		for i := 0; i < len(parts) && err == nil; i++ {
			var arc IRIArc
			if arc, err = NewIRIArc(parts[i]); err == nil {
				iri = append(iri, arc)
			}
		}
	case IRIArc:
		if !tv.ok {
			err = errorIRINil
			break
		}
		iri = append(iri, tv)
	default:
		err = errorIRIBadType
	}

	return
}

func assertIRINotationVariadic(x ...any) (iri IRINotation, err error) {
	// variadic arcs to comprise a full IRI
	for i := 0; i < len(x) && err == nil; i++ {
		switch tv := x[i].(type) {
		case string, int, int32, int64, uint64, *big.Int, NumberForm:
			var arc IRIArc
			if arc, err = NewIRIArc(tv); err == nil {
				iri = append(iri, arc)
			}

		case IRIArc:
			if !tv.ok {
				err = errorIRINil
				break
			}
			iri = append(iri, tv)

		default:
			err = errorIRIBadType
		}
	}

	return
}

/*
MustNewSibling returns an instance of [DotNotation] and panics if
[IRINotation.NewSibling] returned an error during processing of x.
*/
func (r IRINotation) MustNewSibling(x any) IRINotation {
	O, err := r.NewSibling(x)
	if err != nil {
		panic(err)
	}
	return O
}

/*
NewSibling returns a new instance of [IRINotation] based upon the contents
of the receiver as well as the input [IRIArc] sibling value. The result is
the last arc of the new OID being replaced with the input [IRIArc] in the
return instance.

See also [IRINotation.MustNewSibling].
*/
func (r IRINotation) NewSibling(arc any) (
	O IRINotation,
	err error,
) {
	if !r.Valid() {
		err = errorIRINil
		return
	}

	var n IRIArc
	// Prepare the new leaf IRIArc, or die trying.
	if n, err = NewIRIArc(arc); err == nil {
		L := r.Len()
		O = make(IRINotation, L, L)
		for i := 0; i < L-1; i++ {
			O[i] = r[i]
		}
		O[L-1] = n
	}

	return
}

/*
MustNewSubordinate returns an instance of [IRINotation]
and panics if [IRINotation.NewSubordinate] returned an
error during processing of x.
*/
func (r IRINotation) MustNewSubordinate(x any) IRINotation {
	O, err := r.NewSubordinate(x)
	if err != nil {
		panic(err)
	}
	return O
}

/*
NewSubordinate returns a new instance of [IRINotation]
based upon the contents of the receiver as well as the input [IRIArc] subordinate
value. This creates a fully-qualified child [IRINotation]
value of the receiver.

See also [IRINotation.MustNewSubordinate].
*/
func (r IRINotation) NewSubordinate(nf any) (
	O IRINotation,
	err error,
) {
	if !r.Valid() {
		err = errorIRINil
		return
	}

	var n IRIArc
	// Prepare the new leaf iRIArc, or die trying.
	if n, err = NewIRIArc(nf); err == nil {
		L := r.Len()
		alloc := L + 1
		O = make(IRINotation, alloc, alloc)
		for i := 0; i < L; i++ {
			O[i] = r[i]
		}
		O[alloc-1] = n
	}

	return
}

/*
Ancestry returns slices of [IRINotation] instances
ordered as leaf node first.

Empty slices of [IRINotation] are returned if the
value within the receiver is less than two (2) [IRIArc] values in length.

Note that, unlike [ASN1Notation], the root [IRINotation] is
never added to the return instance by itself, as any valid [IRINotation]
MUST have a root AND at least one child. This means that for a receiver value
of "1.3.6.1.4.1", the final ancestor present within the return value will be "1.3".
*/
func (r IRINotation) Ancestry() (
	anc []IRINotation,
) {
	if r.Len() >= 2 {
		for i := r.Len(); i > 0; i-- {
			anc = append(anc, r[:i])
		}
	}

	return
}

/*
Eq returns a Boolean value indicative of an equality match between the
receiver and input [IRINotation] instances.
*/
func (r IRINotation) Eq(o IRINotation) bool {
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
Type returns the string literal "iri".
*/
func (r IRINotation) Type() string { return `iri` }

/*
Eq returns a Boolean value indicative of receiver r matching input value
other. This method evaluates the state of both the underlying [IRIArc]
and string name values.
*/
func (r IRIArc) Eq(other any) bool {
	var same bool
	switch tv := other.(type) {
	case int, int32, int64, uint64, *big.Int, NumberForm, string:
		if arc, err := NewIRIArc(tv); err == nil {
			same = r.Eq(arc)
		}
	case IRIArc:
		same = r.name == tv.name && r.number.Eq(tv.number)
	}

	return same
}

/*
IsZero returns a Boolean value indicative of a nil or zero length receiver state.
*/
func (r IRINotation) IsZero() (is bool) {
	if is = &r == nil; !is {
		is = r.Len() == 0
	}
	return
}

/*
Root returns the root node (0) [IRIArc] value from the receiver.
*/
func (r IRINotation) Root() IRIArc {
	x, _ := r.Index(0)
	return x
}

/*
Leaf returns the leaf node (-1) [IRIArc] value from the receiver.
*/
func (r IRINotation) Leaf() IRIArc {
	x, _ := r.Index(-1)
	return x
}

/*
Parent returns the leaf node's parent (-2) [IRIArc] value from the receiver.
*/
func (r IRINotation) Parent() IRIArc {
	x, _ := r.Index(-2)
	return x
}

/*
AncestorOf returns a Boolean value indicative of whether the receiver is an
ancestor of the input value, which can be string or [IRINotation].
*/
func (r IRINotation) AncestorOf(id any) (anc bool) {
	if !r.IsZero() {
		if A := assertIRINotation(id); !A.IsZero() {
			if A.Len() > r.Len() {
				anc = r.matchIRI(A, 0)
			}
		}
	}

	return
}

/*
ChildOf returns a Boolean value indicative of whether the receiver is
a direct superior (parent) of the input value, which can be string or
[IRINotation].
*/
func (r IRINotation) ChildOf(id any) (cof bool) {
	if !r.IsZero() {
		if A := assertIRINotation(id); !A.IsZero() {
			if A.Len()-1 == r.Len() {
				cof = r.matchIRI(A, 0)
			}
		}
	}

	return
}

/*
SiblingOf returns a Boolean value indicative of whether the receiver is
a sibling of the input value, which can be string or [IRINotation].
*/
func (r IRINotation) SiblingOf(id any) (sof bool) {
	if !r.IsZero() {
		if A := assertIRINotation(id); !A.IsZero() {
			if A.Len() == r.Len() && !A.Leaf().Eq(r.Leaf()) {
				sof = r.matchIRI(A, -1)
			}
		}
	}

	return
}

/*
isIRIChar returns a Boolean value indicative of input value r being
a legal IRI character.

See also the iRISpec var.
*/
func isIRIChar(r rune) bool {
	if r == '/' {
		return false
	}
	return unicode.Is(&iRISpec, r)
}

func (r IRINotation) matchIRI(iri IRINotation, off int) (
	matched bool,
) {
	L := r.Len()
	ct := 0
	for i := 0; i < L; i++ {
		x, _ := r.Index(i)
		if x.Eq(iri[i]) {
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

func assertIRINotation(id any) (A IRINotation) {
	switch tv := id.(type) {
	case string:
		A, _ = NewIRINotation(tv)
	case IRINotation:
		if tv.Len() >= 0 {
			A = tv
		}
	}

	return
}

/*
iRISpec circumscribes unicode range tables for the purpose of
identifying bogus characters during IRI parsing.
*/
var iRISpec = unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0021, 0x007E, 1},
		{0x00A0, 0xD7FF, 1},
		{0xF900, 0xFDCF, 1},
		{0xFDF0, 0xFFEF, 1},
	},
	R32: []unicode.Range32{
		{0x10000, 0x1FFFD, 1},
		{0x20000, 0x2FFFD, 1},
		{0x30000, 0x3FFFD, 1},
		{0x40000, 0x4FFFD, 1},
		{0x50000, 0x5FFFD, 1},
		{0x60000, 0x6FFFD, 1},
		{0x70000, 0x7FFFD, 1},
		{0x80000, 0x8FFFD, 1},
		{0x90000, 0x9FFFD, 1},
		{0xA0000, 0xAFFFD, 1},
		{0xB0000, 0xBFFFD, 1},
		{0xC0000, 0xCFFFD, 1},
		{0xD0000, 0xDFFFD, 1},
		{0xE0000, 0xEFFFD, 1},
	},
}

var (
	errorIRINil       = errors.New("INTERNATIONALIZED RESOURCE IDENTIFIER: nil or bogus instance")
	errorIRINilArc    = errors.New("INTERNATIONALIZED RESOURCE IDENTIFIER: nil or bogus arc instance")
	errorIRIBadType   = errors.New("INTERNATIONALIZED RESOURCE IDENTIFIER: unsupported type")
	errorIRINoInput   = errors.New("INTERNATIONALIZED RESOURCE IDENTIFIER: no input")
	errorIRIArcBadUCS = errors.New("INTERNATIONALIZED RESOURCE IDENTIFIER: invalid UCS character(s) found")
)
