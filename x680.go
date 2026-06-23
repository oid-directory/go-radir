package radir

/*
X680 implements [RASCHEMA § 2.5.6] and derives various concepts from
[ITU-T Rec. X.680].

	( 1.3.6.1.4.1.56521.101.2.5.6
	    NAME 'x680Context'
	    DESC 'X.680 contextual class'
	    SUP registration AUXILIARY
	    MAY ( aSN1Notation $
	          dotNotation $
	          identifier $
	          iRI $
	          nameAndNumberForm ) )

Instances of this type need not be initialized by the user directly.

[RASCHEMA § 2.5.6]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.6
[ITU-T Rec. X.680]: https://www.itu.int/rec/T-REC-X.680
*/
type X680 struct {
	R_N          string   `ldap:"n"`                 // RASCHEMA § 2.3.1
	R_ASN1Not    string   `ldap:"aSN1Notation"`      // RASCHEMA § 2.3.4
	R_DotNot     string   `ldap:"dotNotation"`       // RASCHEMA § 2.3.2
	R_Id         string   `ldap:"identifier"`        // RASCHEMA § 2.3.7
	R_NaNF       string   `ldap:"nameAndNumberForm"` // RASCHEMA § 2.3.19
	R_IRI        []string `ldap:"iRI"`               // RASCHEMA § 2.3.3
	r_DITProfile *DITProfile
	r_root       *registeredRoot
}

/*
X680 returns (and if needed, initializes) the embedded instance of *[X680].
*/
func (r *Registration) X680() *X680 {
	if r.IsZero() {
		return &X680{}
	}

	if r.R_X680.IsZero() {
		r.R_X680 = new(X680)
		r.R_X680.r_DITProfile = r.Profile()
		r.R_X680.r_root = r.r_root
	}

	return r.R_X680
}

/*
profile returns the *[DITProfile] instance assigned to the receiver,
if set, else nil is returned.
*/
func (r *X680) profile() (prof *DITProfile) {
	if prof = r.r_DITProfile; !prof.Valid() {
		prof = &DITProfile{}
	}

	return
}

func (r *X680) ldif() (l string) {
	if !r.IsZero() {
		l = toLDIF(r)
	}

	return
}

/*
Unmarshal returns an instance of map[string][]string bearing the contents
of the receiver.
*/
func (r *X680) unmarshal() map[string][]string {
	m := make(map[string][]string)
	return unmarshalStruct(r, m)
}

/*
N returns the number form, or primary identifier, of the receiver. If unset,
a zero string is returned.
*/
func (r *X680) N() string {
	return r.R_N
}

/*
NGetFunc executes the [GetOrSetFunc] instance and returns its own return values.
*/
func (r *X680) NGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `n`)
}

/*
SetN assigns the string number form value to the receiver instance.
*/
func (r *X680) SetN(args ...any) error {
	return writeFieldByTag(`n`, r.SetN, r, args...)
}

/*
IRI returns the string "[iRI]" values assigned to the receiver instance,
each defining an internationalized resource identifier.

[iRI]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.3
*/
func (r *X680) IRI() []string {
	return r.R_IRI
}

/*
SetIRI appends one or more string "[iRI]" values to the receiver instance.
Note that if a slice is passed as X, the destination value will be clobbered.

[iRI]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.3
*/
func (r *X680) SetIRI(args ...any) error {
	return writeFieldByTag(`iRI`, r.SetIRI, r, args...)
}

/*
IRIGetFunc processes the underlying string "[iRI]" field values through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.

[iRI]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.3
*/
func (r *X680) IRIGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `iRI`)
}

/*
NameAndNumberForm returns the string "[nameAndNumberForm]" value assigned
to the receiver, or a zero string if unset.

[nameAndNumberForm]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.19
*/
func (r *X680) NameAndNumberForm() string {
	return r.R_NaNF
}

/*
SetNameAndNumberForm assigns the string "[nameAndNumberForm]" value to the
receiver instance.

[nameAndNumberForm]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.19
*/
func (r *X680) SetNameAndNumberForm(args ...any) error {
	return writeFieldByTag(`nameAndNumberForm`, r.SetNameAndNumberForm, r, args...)
}

/*
NameAndNumberFormGetFunc processes the underlying "[nameAndNumberForm]"
value through the provided [GetOrSetFunc] instance, returning an interface
value alongside an error.

[nameAndNumberForm]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.19
*/
func (r *X680) NameAndNumberFormGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `nameAndNumberForm`)
}

/*
Identifier returns the string "[identifier]", or "name form", of the receiver.
If unset, a zero string is returned.

[identifier]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.7
*/
func (r *X680) Identifier() string {
	return r.R_Id
}

/*
SetIdentifier assigns the string "[identifier]" value to the receiver instance.

[identifier]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.7
*/
func (r *X680) SetIdentifier(args ...any) error {
	return writeFieldByTag(`identifier`, r.SetIdentifier, r, args...)
}

/*
IdentifierGetFunc processes the underlying "[identifier]" field value through
the provided [GetOrSetFunc] instance, returning an interface value alongside
an error.

[identifier]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.7
*/
func (r *X680) IdentifierGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `identifier`)
}

/*
DotNotation returns the string "[dotNotation]" value assigned to the receiver,
or a zero string if unset.

[dotNotation]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.2
*/
func (r *X680) DotNotation() (dot string) {
	if !r.IsZero() {
		dot = r.R_DotNot
	}

	return
}

/*
DotNotationGetFunc processes the underlying "[dotNotation]" field value through
the provided [GetOrSetFunc] instance, returning an interface value alongside
an error.

[dotNotation]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.2
*/
func (r *X680) DotNotationGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `dotNotation`)
}

/*
SetDotNotation assigns the string "[dotNotation]" value to the receiver instance.

[dotNotation]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.2
*/
func (r *X680) SetDotNotation(args ...any) error {
	return writeFieldByTag(`dotNotation`, r.SetDotNotation, r, args...)
}

/*
ASN1Notation returns the string "[aSN1Notation]" value assigned to the receiver,
or a zero string if unset.

[aSN1Notation]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.4
*/
func (r *X680) ASN1Notation() (asn string) {
	if !r.IsZero() {
		asn = r.R_ASN1Not
	}

	return
}

/*
SetASN1Notation assigns the string "[aSN1Notation]" notation value to the receiver
instance.

[aSN1Notation]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.4
*/
func (r *X680) SetASN1Notation(args ...any) error {
	return writeFieldByTag(`aSN1Notation`, r.SetASN1Notation, r, args...)
}

/*
ASN1NotationGetFunc processes the underlying "[aSN1Notation]" field value
through the provided [GetOrSetFunc] instance, returning an interface value
alongside an error.

[aSN1Notation]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.4
*/
func (r *X680) ASN1NotationGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `aSN1Notation`)
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *X680) IsZero() bool {
	return r == nil
}

func (r *X680) isEmpty() bool {
	return structEmpty(r)
}

/*
marshal returns an error following an attempt to execute the input meth
"func(any) error" method signature.

The meth value should be the (unexecuted) [go-ldap/v3 Entry.Unmarshal]
method instance for the [Entry] being transferred (marshaled) into the
receiver instance.

Alternatively, if the user has fashioned an alternative method of the
same signature, this may be supplied instead.

[go-ldap/v3 Entry.Unmarshal]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry.Unmarshal
[Entry]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry
*/
func (r *X680) marshal(meth func(any) error) (err error) {
	if !r.IsZero() {
		err = meth(r)
	}

	return
}

func (r *X680) specialHandling(tag string, value any) {
	if value != nil && !r.IsZero() {
		switch lc(tag) {
		case `dotnotation`:
			r.dotNotationHandler(value.(string))
		case `asn1notation`:
			r.asn1NotationHandler(value.(string))
		}
	}
}

func (r *X680) asn1NotationHandler(ant string) {
	if !structEmpty(r.r_root) || len(ant) == 0 {
		return
	}

	if _, slices, err := cleanASN1(ant); err == nil {
		var root, id, nanf string
		if root, id, nanf, err = checkASN1Root(slices[0]); err == nil {
			r.r_root.N = getRoot(rune(root[0]))
			r.r_root.Depth = len(slices)
			r.r_root.NaNF = nanf
			r.r_root.Id = id
			r.r_root.Auxiliary = rootClass(r.r_root.N)
			if r.r_root.Depth == 1 {
				r.r_root.Structural = `rootArc`
			} else if r.r_root.Depth > 1 {
				r.r_root.Structural = `arc`
			}
		}
	}
}

func (r *X680) dotNotationHandler(dot string) {
	spl := dotSplit(dot)
	if len(spl) == 0 {
		return
	} else if len(spl[0]) == 0 {
		return
	}

	if structEmpty(r.r_root) {
		if n := getRoot(rune(spl[0][0])); 0 <= n && n <= 2 {
			r.r_root.N = n
			r.r_root.Depth = len(spl)
		}
	}
}

/*
Depth returns the integer number of dotted decimal values present within
the underlying OID. For example, the OID "1.3.6.1.4.1.56521.101" (the
[OIDPrefix] constant) has a depth of eight (8), while any root arc --
such as joint-iso-itu-t(2) -- will always have a depth of one (1).

If no OID is present, zero (0) is returned.
*/
func (r *X680) Depth() (d int) {
	d = -1
	if !r.IsZero() {
		d = r.r_root.Depth
	}

	return
}
