package oid

/*
nanf.go contains all types and methods pertaining to the X.680
NameAndNumberForm type.
*/

import (
	"errors"
	"math/big"
	"strings"
)

/*
NameAndNumberForm implements the ITU-T rec. X.680 "name and number form".
Instances of this type are created using the [NewNameAndNumberForm] or
[MustNewNameAndNumberForm] constructor.
*/
type NameAndNumberForm struct {
	nameForm   string     // e.g.: "iso"; optional, must conform to X.680 name form notation
	numberForm NumberForm // number form type; always required, never negative
	ok         bool       // avoid ambiguity with zero instances
}

/*
NumberForm returns the underlying name form string instance, if specified.
A zero string is returned otherwise.
*/
func (r NameAndNumberForm) NameForm() string {
	var nf string
	if r.ok {
		nf = r.nameForm
	}
	return nf
}

/*
NumberForm returns the underlying number form [NumberForm] instance.

If the value indicates the receiver is in an aberrant state, or was
not created properly using the [NewNameAndNumberForm] or
[MustNewNameAndNumberForm] constructor, this method will panic.
*/
func (r NameAndNumberForm) NumberForm() NumberForm {
	var nf NumberForm
	if !r.ok {
		panic("NumberForm: invalid instance")
	}
	nf = r.numberForm
	return nf
}

/*
String returns the string representation of the receiver instance.

If the receiver is in an aberrant state, or was not created properly
using the [NewNameAndNumberForm] or [MustNewNameAndNumberForm]
constructor, the return string value will be zero.
*/
func (r NameAndNumberForm) String() string {
	var n string
	if r.ok {
		n = r.numberForm.String()
		if len(r.nameForm) > 0 {
			n = r.nameForm + `(` + n + `)`
		}
	}
	return n
}

/*
Eq returns a Boolean value indicative of receiver r matching input value
other. This method evaluates the state of both the underlying [NumberForm]
and string nameForm values.
*/
func (r NameAndNumberForm) Eq(other NameAndNumberForm) bool {
	return r.nameForm == other.nameForm &&
		r.numberForm.Eq(other.numberForm)
}

/*
MustNewNameAndNumberForm returns an instance of [NameAndNumberForm]
by calling [NewNameAndNumberForm], and will panic if an error is
encountered.
*/
func MustNewNameAndNumberForm(nanf ...any) NameAndNumberForm {
	nf, err := NewNameAndNumberForm(nanf...)
	if err != nil {
		panic(err)
	}
	return nf
}

/*
NewNameAndNumberForm returns an instance of [NameAndNumberForm]
alongside an error following an attempt to marshal nanf.

This function is not normally needed by the end user, and only
exists to be used in iterative fashion by [NewASN1Notation]
and [MustNewASN1Notation].

If only one argument is provided, the following requirements are
imposed:

  - If it is a string, must be purely numerical OR must represent a complete [NameAndNumberForm] value
  - If it is a supported number (e.g.: int, etc), it is interpreted as a solitary [NumberForm]

If two arguments are provided, the following requirements are
imposed:

  - First argument must be a [NumberForm]-compliant value, whether a supported number type or string number
  - Second argument must be a valid name form string (optional).

Supplying any other number of arguments, or unsupported types, is an error.

See [NewNumberForm] for supported number types.
*/
func NewNameAndNumberForm(nanf ...any) (NameAndNumberForm, error) {
	var (
		r   NameAndNumberForm
		err error
	)

	switch len(nanf) {
	case 1:
		switch tv := nanf[0].(type) {
		case string:
			var nf NumberForm
			if nf, err = NewNumberForm(tv); err == nil {
				// seems to be just digits (numberForm only)
				r.numberForm = nf
				r.ok = true
			} else {
				// seems to be a complete NaNF string
				r, err = newNameAndNumberFormStr(tv)
			}
		case int, int32, int64, uint64, *big.Int, NumberForm:
			var nf NumberForm
			if nf, err = NewNumberForm(tv); err == nil {
				r.numberForm = nf
				r.ok = true
			}
		default:
			err = errorNaNFBadType
		}
	case 2:
		var nf NumberForm
		if nf, err = NewNumberForm(nanf[0]); err == nil {
			switch tv := nanf[1].(type) {
			case string:
				if !IsNameForm(tv) {
					err = errorNameFNil
					break
				}
				r.numberForm = nf
				r.nameForm = tv
				r.ok = true
			default:
				err = errorNaNFBadType
			}
		}
	default:
		err = errorNaNFBadArgCt
	}

	return r, err
}

func newNameAndNumberFormStr(x string) (NameAndNumberForm, error) {
	var (
		err error
		r   NameAndNumberForm
	)
	// Don't waste time on bogus values.
	if len(x) == 0 {
		err = errorNaNFNil
		return r, err
	} else if x[len(x)-1] != ')' {
		// if there was no closing paren, this COULD
		// mean there is only a number form (with no
		// name). If this is the case, parse as an
		// integer by itself, else throw an error.
		var n NumberForm
		if n, err = NewNumberForm(x); err == nil {
			r = NameAndNumberForm{
				numberForm: n,
				ok:         true,
			}
		}
		return r, err
	}

	// index the rune for '(', indicating the
	// identifier (nameForm) has ended, and the
	// numberForm is beginning.
	idx := strings.IndexRune(x, '(')
	if idx == -1 {
		err = errorNaNFNil
		return r, err
	}

	var n NumberForm
	if n, err = NewNumberForm(x[idx+1 : len(x)-1]); err == nil {
		// Parse/verify what appears to be the
		// identifier string value.
		if id := x[:idx]; !IsNameForm(id) {
			err = errorNameFNil
		} else {
			r = NameAndNumberForm{
				nameForm:   id,
				numberForm: n,
				ok:         true,
			}
		}
	}

	return r, err
}

/*
IsNameForm returns a Boolean value indicative of the val string
input value being a legal X.680 name form, in that:

  - The name is at least one character long, and ...
  - The first character is a lowercase letter, and ...
  - The final character is alphanumeric, and ...
  - Any hyphens present are NOT contiguous (e.g.: "--")
*/
func IsNameForm(val string) bool {
	if len(val) == 0 {
		return false
	}

	isLower := func(r rune) bool {
		return ('a' <= r && r <= 'z')
	}

	if !isLower(rune(val[0])) {
		// must begin with a lower alpha.
		return false
	}

	isAlnum := func(r rune) bool {
		return isLower(r) || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9')
	}

	L := len(val) - 1
	if rune(val[L]) == '-' {
		// can't end in a hyphen
		return false
	}

	// watch hyphens to avoid contiguous use.
	// A value of true at any point means the
	// PREVIOUS char was a hyphen.
	var lastHyphen bool

	// iterate all characters in val (except for the
	// first and final chars already checked above),
	// checking each one for "descr" validity.
	for i := 1; i < L; i++ {
		ch := rune(val[i])
		switch {
		case isAlnum(ch):
			lastHyphen = false
		case ch == '-':
			if lastHyphen {
				// cannot use consecutive hyphens
				return false
			}
			lastHyphen = true
		default:
			// invalid character (none of [a-zA-Z0-9\-])
			return false
		}
	}

	return true
}

/*
IsNameAndNumberForm returns a Boolean value indicative of the nanf
string input value representing a valid [NameAndNumberForm] in syntax.

Use of the variadic relax bool, when false, will require the presence
of an underlying name form in addition to the requisite number form.
When true, a lone number form will suffice. The default is false.

While the standard expects each [NameAndNumberForm] value to be single-space
delimited, this handler accepts multiple spaces, newlines and horizontal
tabs as such delimiters without issue. This can be useful if the input
value has been whitespace-formatted for readability.

See also [IsNumberForm] and [IsNameForm].
*/
func IsNameAndNumberForm(nanf string, relax ...bool) bool {

	idxF := strings.IndexRune(nanf, '(')
	idxL := strings.IndexRune(nanf, ')')

	if idxF == -1 || idxL == -1 {
		if len(relax) > 0 && relax[0] {
			// possible number form only
			return IsNumberForm(nanf)
		}
		// unbalanced paren encaps
		return false
	} else if idxL != len(nanf)-1 {
		// closing paren MUST be final char
		return false
	}

	if nameForm := nanf[:idxF]; !IsNameForm(nameForm) {
		// nameForm component is bogus
		return false
	}

	// finally, check if numberForm component is bogus
	return IsNumberForm(nanf[idxF+1 : idxL])
}

var (
	errorNaNFNil      = errors.New("NAME AND NUMBER FORM: nil or bogus instance")
	errorNameFNil     = errors.New("NAME FORM: nil or bogus instance")
	errorNaNFBadType  = errors.New("NAME AND NUMBER FORM: unsupported input type")
	errorNaNFBadArgCt = errors.New("NAME AND NUMBER FORM: unexpected number of input arguments")
)
