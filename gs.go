package radir

/*
gs.go contains stock GetOrSetFunc-qualified functions for use during the
composition or interrogation of Registration and/or Registrant instances,
as well as to offer practical examples regarding the creation of such
CUSTOM function instances by the user.
*/

import "time"

/*
GetOrSetFunc is a first class (closure) function signature that users
may adopt in order to write custom "setter or getter" functions, thereby
allowing complete control over the creation or interrogation of a value
assigned to *[Registration] or [Registrant]-qualified type instances.

All Set<*> and <*>GetFunc methods extended by [*Registration] or
[Registrant]-qualified type instances allow the GetOrSetFunc type to
be passed at runtime.

If no GetOrSetFunc is passed to a Set<*> method, the first (and only)
value is written as-is (type assertion permitting) with no special
processing.

However, if a GetOrSetFunc is passed as the second argument to a Set<*>
method, it will be executed and will process the first input value prior
to writing the value.

For each of the extended types within this package, struct fields capable
of being user-managed are either string or []string instances.

In the context of Set<*> executions, When a single string value is submitted
to a Set<*> method which interacts with a []string instance, the value will
be appended.  However, if a []string instance is submitted, it will clobber
the preexisting "[]string" instance. In other words, preexisting values will
be overwritten. You have been warned.

In the context of Set<*> return values, the first return value will be the
(post-processed) string value to be written. The second return value, an
error, will contain error information in the event of any encountered issues.

In the context of <*>GetFunc executions, the first (and only) input argument
will be the struct field value relevant to the executing method. This will
produce the value being "gotten" within functions/methods that conform to the
signature of this type. The second input argument will be the *[Registration]
or [Registrant]-qualified instance being interrogated.

In the context of <*>GetFunc return values, the first return value will
be the (processed) value being read. It will manifest as an instance of
'any', and will require manual type assertion. The second return value,
an error, will contain error information in the event of any encountered
issues.
*/
type GetOrSetFunc func(...any) (any, error)

const tfmt = `20060102150405`

var (
	since func(time.Time) time.Duration = time.Since
	until func(time.Time) time.Duration = time.Until
	now   func() time.Time              = time.Now
)

/*
ASN1NotationToMulti returns an instance of []string alongside an error.

The purpose of this function is to take an "[aSN1Notation]" value, such as ...

	{iso identified-organization(3) dod(6) internet(1) private(4) enterprise(1) 56521 oid-directory(101)}

... and create up to five (5) values

  - "[aSN1Notation]: (sanitized input)
  - "[dotNotation]": 1.3.6.1.4.1.56521.101
  - "[identifier]": oid-directory
  - "[nameAndNumberForm]": oid-directory(101)
  - "[n]": 101

In certain conditions, it may be possible to set yet another value -- the
DN of the *[Registration] instance. See the "SixValuesFromOne" example below.

If the leaf arc within the input value lacks an "[identifier]", the subsequent
"[nameAndNumberForm]" value will also be zero.

Note that the cleaned-up input value is returned within the []string return
instance as the first value unless deemed bogus, at which point a zero
instance is returned.

[n]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.1
[identifier]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.7
[dotNotation]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.2
[aSN1Notation]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.4
[nameAndNumberForm]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.19
*/
func ASN1NotationToMulti(input string) ([]string, error) {
	_, slices, err := cleanASN1(input)
	if err != nil {
		return nil, err
	}

	var root, rid, rnanf string
	if root, rid, rnanf, err = checkASN1Root(slices[0]); err != nil {
		return nil, err
	}

	slices[0] = rnanf
	if len(slices) == 1 {
		return []string{`{`+rnanf+`}`, ``, rid, rnanf, root}, nil
	}

	// append root number form to dot
	// but ONLY if there is more than
	// one (1) arc.
	var dot []string
	dot = append(dot, root)

	var nanf string
	var ident string
	var nf string

	for i := 1; i < len(slices); i++ {
		id, n := nanfToIdAndNF(slices[i])
		if len(n) == 0 {
			return nil, IllegalASN1NotationErr
		}

		if i == 1 && root != `2` {
			// Check to make sure the 2nd arc, if defined,
			// is no greater than 39 unless the root was 2
			if nf, err := atoi(n); err != nil || !(0 <= nf && nf <= 39) {
				return nil, IllegalASN1NotationErr
			}
		} else if i == len(slices)-1 {
			nf = n
			if ident = id; len(ident) > 0 {
				nanf = id + `(` + nf + `)`
			}
		}
		dot = append(dot, n)
	}

	return []string{`{`+join(slices,` `)+`}`, join(dot, `.`), ident, nanf, nf}, nil
}

func checkASN1Root(slice string) (root, rid, rnanf string, err error) {
	// list the root forms we will support
	roots := [][]string{
		{`itu-t`, `itu-t(0)`, `0`},
		{`iso`, `iso(1)`, `1`},
		{`joint-iso-itu-t`, `joint-iso-itu-t(2)`, `2`},
	}

	// make sure we got a valid root
	_root := -1
	for idx, r := range roots {
		if strInSlice(slice, r) {
			_root = idx
			root = itoa(idx)
			rid = roots[idx][0]
			rnanf = rid + `(` + root + `)`
			break
		}
	}

	if _root == -1 {
		err = IllegalRootErr
	}

	return
}

func cleanASN1(anot string) (a string, sl []string, err error) {
	if len(anot) == 0 {
		err = NilArgumentsErr
		return
	}

	if anot[0] != '{' || anot[len(anot)-1] != '}' {
		err = IllegalASN1NotationErr
		return
	}

	// clean up input value
	a = replaceAll(anot, string(rune(10)), ``)
	a = condenseWHSP(trimS(trimR(trimL(trimS(a), `{`), `}`)))

	// split input into slices
	sl = split(a, ` `)
	return
}

func nanfToIdAndNF(in string) (id, n string) {
	if len(in) == 0 {
		return
	}

	if idx := idxRune(in, '('); idx != -1 {
		final := len(in) - 1
		if in[final] != ')' {
			return
		}
		var _id, _n string
		if _id = in[:idx]; !isIdentifier(_id) {
			return
		}
		if _n = in[idx+1 : final]; !isNumber(_n) {
			return
		}
		id = _id
		n = _n
	} else {
		if !isNumber(in) {
			return
		}
		n = in
	}

	return
}

/*
gt2t marshals the input generalized time string value (in) into an
instance of [time.Time] (out).
*/
func gt2t(in string) (out time.Time, err error) {
	var (
		format string = tfmt // base format
		diff   string = `-0700`
		base   string
	)

	if len(in) < 15 {
		err = InvalidGTErr
		return
	}

	if zulu := in[len(in)-1] == 'Z'; zulu {
		in = in[:len(in)-1]
	}

	// If we've got nothing left, must be zulu
	// without any fractional or differential
	// components
	if base = in[14:]; len(base) == 0 {
		out, err = time.Parse(format, in)
		return
	}

	// Handle fractional component (up to six (6) digits)
	if base[0] == '.' || base[0] == ',' {
		format += string(".")
		for fidx, ch := range base[1:] {
			if fidx > 6 {
				err = InvalidGTFracErr
				return
			} else if isDigit(ch) {
				format += `0`
			}
		}
	}

	// Handle differential time, or bail out if not
	// already known to be zulu.
	if in[len(in)-5] == '+' || in[len(in)-5] == '-' {
		format += diff
	}

	out, err = time.Parse(format, in)

	return
}

/*
t2gt marshals the input [time.Time] instance (in) into a generalized
time string instance (out).
*/
func t2gt(in time.Time) (out string) {
	out = in.Format(tfmt) + `Z`
	return
}

/*
GeneralizedTimeToTime converts one or more string-based generalized time
values into a UTC-aligned [time.Time] instances, which are then added as
slices to an instance of [][time.Time] (as an interface type), and then
returned alongside an error. This function qualifies for the [GetOrSetFunc]
type signature.

For more information about functions such as this one, as well as details
on writing your own speciality functions/methods, see the documentation
for the [GetOrSetFunc] closure type.
*/
func GeneralizedTimeToTime(args ...any) (T any, err error) {
	if len(args) == 0 {
		err = NilArgumentsErr
		return
	}
	X := args[0]

	switch tv := X.(type) {
	case string:
		var t time.Time
		if t, err = gt2t(tv); err != nil {
			return
		}

		if !t.IsZero() {
			T = t
			return
		}
	case []string:
		var Tx []time.Time
		for i := 0; i < len(tv); i++ {
			var t time.Time
			if t, err = gt2t(tv[i]); err != nil {
				return
			}
			Tx = append(Tx, t)
		}

		if len(Tx) > 0 {
			T = Tx
			return
		}
	default:
		err = InvalidGTErr
	}

	return
}

/*
TimeToGeneralizedTime converts an instance of [time.Time] into a string
value instance of generalized time, which is returned alongside an error.
This function qualifies for the [GetOrSetFunc] type signature.

For more information about functions such as this one, as well as details
on writing your own specialty functions/methods, see the documentation
for the [GetOrSetFunc] closure type.
*/
func TimeToGeneralizedTime(args ...any) (G any, err error) {
	if len(args) == 0 {
		err = NilArgumentsErr
		return
	}
	X := args[0]

	switch tv := X.(type) {
	case time.Time:
		if g := t2gt(tv); len(g) > 0 && !tv.IsZero() {
			G = g
			return
		}
	case []time.Time:
		var Gx []string
		for i := 0; i < len(tv); i++ {
			if g := t2gt(tv[i]); len(g) > 0 && !tv[i].IsZero() {
				Gx = append(Gx, g)
			}
		}

		if len(Gx) > 0 {
			G = Gx
			return
		}
	}

	err = InvalidGTErr

	return
}

/*
DotNotToDN2D returns a string-based LDAP distinguished name value (dn)
based upon the contents of the input ASN.1 dotNotation value (X) alongside
an error. This function qualifies for the [GetOrSetFunc] type signature.

This conforms to the two dimensional DN syntax, as described in [Section
3.1.2 of the RADIT I-D]. This function will output a distinguished name
value that uses the dotNotation for the RDN type. Individual number forms
present within the dotNotation are verified as non-negative numbers, but
are not modified.

For more information about functions such as this one, as well as details
on writing your own speciality functions/methods, see the documentation
for the [GetOrSetFunc] closure type.

[Section 3.1.2 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.1.2
*/
func DotNotToDN2D(args ...any) (dn any, err error) {
	if len(args) < 2 {
		err = NilArgumentsErr
		return
	}

	r, ok := args[1].(*Registration)
	if !ok {
		err = NilRegistrationErr
		return
	}

        // Grab our DITProfile instance if defined, else throw          
        // an error.  This is required because we need to know          
        // the reg. base string value as well as the directory          
        // model in use.                                                
        var duaConf *DITProfile = r.DITProfile()                        
        // We want at least one Registration Base                       
        // and our model MUST be 3D. Return error                       
        // value otherwise.                                             
        if !duaConf.Valid() || duaConf.Model() != TwoDimensional {      
                err = DUAConfigValidityErr                              
                return                                                  
        }                                                               
                                                                        
	// store the original OID here, as we'll
	// just slap it on our composite DN later.
	var O string

	// Make sure input args[0] type is supported, else bail.
	switch tv := args[0].(type) {
	case string:
		O = tv
	default:
		// This stock function only allows string-based OIDs as
		// input. If you need something more specialized, such
		// as asn1.ObjectIdentifier, write your own GetOrSetFunc
		err = InvalidOIDErr
		return
	}

	// Just scan the slices and verify as
	// a number; no alterations needed.
	var D []string = split(O, `.`) // temporary storage for verification
	for i := 0; i < len(D); i++ {
		if !isNumber(D[i]) {
			err = InvalidOIDErr
			return
		}
	}

	// Prepare return value.
	dn = `dotNotation=` + O + `,` + duaConf.RegistrationBase()
	return
}

/*
DNtoDotNot2D returns a dotNotation-based ASN.1 Object Identifier (id)
based upon the contents of the input string distinguished name value (X)
alongside an error. This function qualifies for the [GetOrSetFunc] type
signature.

This conforms to the two dimensional DN syntax, as described in [Section
3.1.2 of the RADIT I-D]. This function expects the use of dotNotation in
the RDN.

For more information about functions such as this one, as well as details
on writing your own speciality functions/methods, see the documentation
for the [GetOrSetFunc] closure type.

[Section 3.1.2 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.1.2
*/
func DNToDotNot2D(args ...any) (id any, err error) {
	if len(args) < 2 {
		err = NilArgumentsErr
		return
	}

        r, ok := args[1].(*X680)
        if !ok {
                err = NilRegistrationErr
                return
        }

	var D string
	switch tv := args[0].(type) {
	case string:
		D = tv
	default:
		// This stock function only allows string-based DNs as
		// input. If you need something more specialized, such
		// as *ldap.DN, write your own GetOrSetFunc
		err = InvalidDNErr
		return
	}

        // Grab our DITProfile instance if defined, else throw
        // an error.  This is required because we need to know
        // the reg. base string value as well as the directory
        // model in use.
        var duaConf *DITProfile = r.DITProfile()
        // We want at least one Registration Base
        // and our model MUST be 3D. Return error
        // value otherwise.
        if !duaConf.Valid() || duaConf.Model() != TwoDimensional {
                err = DUAConfigValidityErr
                return
        }

        bidx := duaConf.RegistrationSuffixEqual(D)
        if bidx == -1 {
                err = InvalidDNErr
                return
        }

	var N string
	if idx := idxRune(D, ','); idx != -1 {
		N = replaceAll(lc(D[:idx]), `dotnotation=`, ``)
	}

	// In case the input DN was whacky,
	// or flat-out zero, let's put a
	// stop to this madness now.
	if len(N) == 0 {
		err = InvalidDNErr
		return
	}

	// Just scan the slices and verify as
	// a number; no alterations needed.
	var S []string = split(N, `.`)
	for i := 0; i < len(S); i++ {
		if !isNumber(S[i]) {
			err = InvalidDNErr
			return
		}
	}

	id = N
	return
}

/*
DotNotToDN3D returns a string-based LDAP distinguished name value (dn)
based upon the contents of the input ASN.1 dotNotation value (X) alongside
an error. This function qualifies for the [GetOrSetFunc] type signature.

This conforms to the three dimensional DN syntax, as described in [Section
3.1.3 of the RADIT I-D]. This function will output relative distinguished
name values, each of whom describe specific number form values, using the
preferred RDN type descriptor 'n'.

For more information about functions such as this one, as well as details
on writing your own speciality functions/methods, see the documentation
for the [GetOrSetFunc] closure type.

[Section 3.1.3 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.1.3
*/
func DotNotToDN3D(args ...any) (dn any, err error) {
	if len(args) < 2 {
		err = NilArgumentsErr
		return
	}

	r, ok := args[1].(*Registration)
	if !ok {
		err = NilRegistrationErr
		return
	}

        // Grab our DITProfile instance if defined, else throw
        // an error.  This is required because we need to know
        // the reg. base string value as well as the directory
        // model in use.
        var duaConf *DITProfile = r.DITProfile()
        // We want at least one Registration Base
        // and our model MUST be 3D. Return error
        // value otherwise.
        if !duaConf.Valid() || duaConf.Model() != ThreeDimensional {
                err = DUAConfigValidityErr
                return
        }

	var D []string // store a portion of the original dotNotation

	// Make sure input args[0] type is supported, else bail.
	switch tv := args[0].(type) {
	case string:
		D = split(tv, `.`)
	default:
		// This stock function only allows string-based DNs as
		// input. If you need something more specialized, such
		// as *ldap.DN, write your own GetOrSetFunc
		err = InvalidOIDErr
		return
	}

	var nfs []string

	// reverse iterate over our D slices,
	// verifying each slice as a number
	// and appending the composite RDN value
	// (n=N) to the nfs string slice type.
	for i := len(D); i > 0; i-- {
		if !isNumber(D[i-1]) {
			err = InvalidOIDErr
			return
		}
		nfs = append(nfs, `n=`+D[i-1])
	}

	// prepare return value, and just grab
	// the first reg. base ...
	dn = join(nfs, `,`) + `,` + duaConf.RegistrationBase()
	return
}

/*
DNToDotNot3D returns a dotNotation-based ASN.1 Object Identifier (id)
based upon the contents of the input string distinguished name value (X)
alongside an error. This function qualifies for the [GetOrSetFunc] type
signature.

This conforms to the three dimensional DN syntax, as described in [Section
3.1.3 of the RADIT I-D]. This function offers positive support for the RDN
type descriptor 'n'.

For more information about functions such as this one, as well as details
on writing your own speciality functions/methods, see the documentation
for the [GetOrSetFunc] closure type.

See also [this gist].

[this gist]: https://gist.github.com/oid-directory/504af4aa5b90c6e235c589b5d647cf2e
[Section 3.1.3 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.1.3
*/
func DNToDotNot3D(args ...any) (id any, err error) {
        if len(args) < 2 {
                err = NilArgumentsErr
                return
        }

        r, ok := args[1].(*X680)
        if !ok {
                err = NilRegistrationErr
                return
        }

        // Grab our DITProfile instance if defined, else throw
	// an error.  This is required because we need to know
	// the reg. base string value as well as the directory
        // model in use.
        var duaConf *DITProfile = r.DITProfile()
        // We want at least one Registration Base
        // and our model MUST be 3D. Return error
        // value otherwise.
        if !duaConf.Valid() || duaConf.Model() != ThreeDimensional {
                err = DUAConfigValidityErr
                return
        }

        // Make sure input args[0] type is supported, else bail.
	var D string
        switch tv := args[0].(type) {
        case string:
		D = lc(tv)
	default:
		err = UnsupportedInputTypeErr
		return
	}

	bidx := duaConf.RegistrationSuffixEqual(D)
	if bidx == -1 {
		err = InvalidDNErr
		return
	}
	base := lc(duaConf.RegistrationBase(bidx))
	component := split(trimR(D, `,`+base), `,`)

	var _dot []string
	// Iterate DN in reverse order to produce the
	// correctly-ordered dotNotation.
	for i := len(component); i > 0; i-- {
		n := split(component[i-1], `=`)
		if len(n) != 2 {
			// The slices should never be anything other
			// than n=<digit+> (exactly 2 values)
			err = InvalidDNErr
			return
		}
		_dot = append(_dot, n[1])
	}

	id = dotJoin(_dot)

	return
}

