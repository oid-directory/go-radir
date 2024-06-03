package radir

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/JesseCoretta/go-objectid"
)

/*
GeneralizedTime aliases an instance of [time.Time] to implement [RFC4517 § 3.3.13].

        GeneralizedTime = century year month day hour
                             [ minute [ second / leap-second ] ]
                             [ fraction ]
                             g-time-zone

        century = 2(%x30-39) ; "00" to "99"
        year    = 2(%x30-39) ; "00" to "99"
        month   =   ( %x30 %x31-39 ) ; "01" (January) to "09"
                  / ( %x31 %x30-32 ) ; "10" to "12"
        day     =   ( %x30 %x31-39 )    ; "01" to "09"
                  / ( %x31-32 %x30-39 ) ; "10" to "29"
                  / ( %x33 %x30-31 )    ; "30" to "31"
        hour    = ( %x30-31 %x30-39 ) / ( %x32 %x30-33 ) ; "00" to "23"
        minute  = %x30-35 %x30-39                        ; "00" to "59"

        second      = ( %x30-35 %x30-39 ) ; "00" to "59"
        leap-second = ( %x36 %x30 )       ; "60"

        fraction        = ( DOT / COMMA ) 1*(%x30-39)
        g-time-zone     = %x5A  ; "Z"
                          / g-differential
        g-differential  = ( MINUS / PLUS ) hour [ minute ]
        MINUS           = %x2D  ; minus sign ("-")

See also [§ 2.3.37], [§ 2.3.56], [§ 2.3.57], [§ 2.3.76] and [§ 2.3.77] of
the [RASCHEMA ID].

[RFC4517 § 3.3.13]: https://datatracker.ietf.org/doc/html/rfc4517#section-3.3.13
[§ 2.3.37]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.37
[§ 2.3.56]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.56
[§ 2.3.57]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.57
[§ 2.3.76]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.76
[§ 2.3.77]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.77
[RASCHEMA ID]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema
*/
type GeneralizedTime time.Time

var (
	eq    func(string, string) bool = strings.EqualFold
	itoa  func(int) string          = strconv.Itoa
	mkerr func(string) error	= errors.New
)

func parseNumberForm(x any) (nf NumberForm, err error) {
	var _nf objectid.NumberForm

	if _nf, err = objectid.NewNumberForm(x); err == nil {
		nf = NumberForm(_nf)
	}

	return
}

/*
String returns the string representation of the receiver instance.
*/
func (r NumberForm) String() string {
	return objectid.NumberForm(r).String()
}

func parseDotNotation(x any) (dot DotNotation, err error) {
	var _dot *objectid.DotNotation

	if _dot, err = objectid.NewDotNotation(x); err == nil {
		dot = DotNotation(*_dot)
	}

	return
}

/*
String returns the string representation of the receiver instance.
*/
func (r DotNotation) String() string {
	return objectid.DotNotation(r).String()
}

func parseNameAndNumberForm(x any) (nanf NameAndNumberForm, err error) {
	var _nanf *objectid.NameAndNumberForm

	if _nanf, err = objectid.NewNameAndNumberForm(x); err == nil {
		nanf = NameAndNumberForm(*_nanf)
	}

	return
}

/*
String returns the string representation of the receiver instance.
*/
func (r NameAndNumberForm) String() string {
	return objectid.NameAndNumberForm(r).String()
}

func parseASN1Notation(x any) (anot ASN1Notation, err error) {
	var _anot *objectid.ASN1Notation

	if _anot, err = objectid.NewASN1Notation(x); err == nil {
		anot = ASN1Notation(*_anot)
	}

	return
}

/*
String returns the string representation of the receiver instance.
*/
func (r ASN1Notation) String() string {
	return objectid.ASN1Notation(r).String()
}

/*
String returns the string representation of the receiver instance.
*/
func (r GeneralizedTime) String() string {
	return time.Time(r).Format(`20060102150405`) + `Z`
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r GeneralizedTime) IsZero() bool {
	return time.Time(r).IsZero()
}

/*
parseGeneralizedTime returns an error following an analysis of x in the context
of a Generalized Time value.
*/
func parseGeneralizedTime(x any) (gt GeneralizedTime, err error) {
	var (
		format string = `20060102150405` // base format
		diff   string = `-0700`
		base   string
		raw    string
	)

        switch tv := x.(type) {
        case string:
                if l := len(tv); l < 15 {
                        err = mkerr("Invalid Generalized Time")
                        return
                }
                raw = tv
	case time.Time:
		if tv.IsZero() {
                        err = mkerr("Invalid Generalized Time (nil time.Time)")
                        return
		}
		gt = GeneralizedTime(tv)
		return
        default:
                err = mkerr("Invalid type for Generalized Time")
		return
        }

	raw = chopZulu(raw)

	// If we've got nothing left, must be zulu
	// without any fractional or differential
	// components
	if base = raw[14:]; len(base) == 0 {
		var _gt time.Time
		if _gt, err = time.Parse(format, raw); err == nil {
			gt = GeneralizedTime(_gt)
		}
		return
	}

	// Handle fractional component (up to six (6) digits)
	if format, err = genTimeFracDiffFormat(raw, base, diff, format); err != nil {
		return
	}

	var _gt time.Time
	if _gt, err = time.Parse(format, raw); err == nil {
		gt = GeneralizedTime(_gt)
	}

	return
}

// Handle generalizedTime fractional component (up to six (6) digits)
func genTimeFracDiffFormat(raw, base, diff, format string) (string, error) {
	var err error

        if base[0] == '.' || base[0] == ',' {
                format += string(".")
                for fidx, ch := range base[1:] {
                        if fidx > 6 {
                                err = mkerr(`Fraction exceeds Generalized Time fractional limit`)
                        } else if isDigit(ch) {
                                format += `0`
                                continue
                        }
                        break
                }
        }

        // Handle differential time, or bail out if not
        // already known to be zulu.
        if raw[len(raw)-5] == '+' || raw[len(raw)-5] == '-' {
                format += diff
        }

	return format, err
}

func chopZulu(raw string) string {
        if zulu := raw[len(raw)-1] == 'Z'; zulu {
                raw = raw[:len(raw)-1]
        }

	return raw
}
