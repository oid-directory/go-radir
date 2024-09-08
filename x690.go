package radir

/*
X690 implements [RASCHEMA § 2.5.7] and derives various concepts from
[ITU-T Rec. X.690].

	( 1.3.6.1.4.1.56521.101.2.5.7
	    NAME 'x690Context'
	    DESC 'X.690 contextual class'
	    SUP registration AUXILIARY
	    MAY dotEncoding )

Instances of this type need not be initialized by the user directly.

[RASCHEMA § 2.5.7]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.7
[ITU-T Rec. X.690]: https://www.itu.int/rec/T-REC-X.690
*/
type X690 struct {
	R_DotEnc     string `ldap:"dotEncoding"` // RASCHEMA § 2.3.103
	r_DITProfile *DITProfile
	r_root	     *registeredRoot
}

/*
X690 returns (and if needed, initializes) the embedded instance of *[X690].
*/
func (r *Registration) X690() *X690 {
	if r.IsZero() {
		return &X690{}
	}

	if r.R_X690.IsZero() {
		r.R_X690 = new(X690)
		r.R_X690.r_DITProfile = r.DITProfile()
		r.R_X690.r_root = r.r_root
	}

	return r.R_X690
}

/*
DITProfile returns the *[DITProfile] instance assigned to the receiver,
if set, else nil is returned.
*/
func (r *X690) DITProfile() (prof *DITProfile) {
	if prof = r.r_DITProfile; !prof.Valid() {
		prof = &DITProfile{}
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *X690) IsZero() bool {
	return r == nil
}

func (r *X690) isEmpty() bool {
	return structEmpty(r)
}

func (r *X690) ldif() (l string) {
	if !r.IsZero() {
		l = toLDIF(r)
	}

	return
}

/*
Unmarshal returns an instance of map[string][]string bearing the contents
of the receiver.
*/
func (r *X690) unmarshal() map[string][]string {
	m := make(map[string][]string)
	return unmarshalStruct(r, m)
}

/*
Marshal returns an error following an attempt to execute the input meth
"func(any) error" method signature.

The meth value should be the (unexecuted) [go-ldap/v3 Entry.Unmarshal]
method instance for the [Entry] being transferred (marshaled) into the
receiver instance.

Alternatively, if the user has fashioned an alternative method of the
same signature, this may be supplied instead.

[go-ldap/v3 Entry.Unmarshal]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry.Unmarshal
[Entry]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry
*/
func (r *X690) marshal(meth func(any) error) (err error) {
	if !r.IsZero() {
		err = meth(r)
	}

	return
}

/*
DotEncoding returns the string dotEncoding value assigned to the receiver,
or a zero string if unset.
*/
func (r *X690) DotEncoding() string {
	return r.R_DotEnc
}

/*
DotEncodingGetFunc processes underlying dotEncoding field value through
the provided [GetOrSetFunc] instance, returning an interface value
alongside an error.
*/
func (r *X690) DotEncodingGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `dotEncoding`)
}

/*
SetDotEncoding assigns the string dotEncoding value to the receiver instance.
*/
func (r *X690) SetDotEncoding(args ...any) error {
	return writeFieldByTag(`dotEncoding`, r.SetDotEncoding, r, args...)
}
