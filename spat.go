package radir

/*
Spatial implements an abstraction of the "[spatialContext]" class, as well
as the mechanics described within [Section 3.2.4.19 of the RADIT I-D].

[spatialContext]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.11
[Section 3.2.4.19 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.2.4.19
*/
type Spatial struct {
	R_SupArc   string   `ldap:"supArc"`   // RASCHEMA § 2.3.21
	R_TopArc   string   `ldap:"topArc"`   // RASCHEMA § 2.3.23
	R_MinArc   string   `ldap:"minArc"`   // RASCHEMA § 2.3.27
	R_MaxArc   string   `ldap:"maxArc"`   // RASCHEMA § 2.3.30
	R_LeftArc  string   `ldap:"leftArc"`  // RASCHEMA § 2.3.26
	R_RightArc string   `ldap:"rightArc"` // RASCHEMA § 2.3.29
	R_SubArc   []string `ldap:"subArc"`   // RASCHEMA § 2.3.25

	// COLLECTIVE spatial types
	RC_SupArc string `ldap:"c-supArc;collective"` // RASCHEMA § 2.3.22
	RC_TopArc string `ldap:"c-topArc;collective"` // RASCHEMA § 2.3.24
	RC_MinArc string `ldap:"c-minArc;collective"` // RASCHEMA § 2.3.28
	RC_MaxArc string `ldap:"c-maxArc;collective"` // RASCHEMA § 2.3.31
}

/*
Spatial returns (and if needed, initializes) the embedded instance of
*[Spatial].
*/
func (r *Registration) Spatial() *Spatial {
	if r.IsZero() {
		return &Spatial{}
	}

	if r.R_Spatial.IsZero() {
		r.R_Spatial = new(Spatial)
	}

	return r.R_Spatial
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *Spatial) IsZero() bool {
	return r == nil
}

func (r *Spatial) isEmpty() bool {
	return structEmpty(r)
}

func (r *Spatial) ldif() (l string) {
	if !r.IsZero() {
		l = toLDIF(r)
	}

	return
}

/*
Unmarshal returns an instance of map[string][]string bearing the contents
of the receiver.
*/
func (r *Spatial) unmarshal() map[string][]string {
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
func (r *Spatial) marshal(meth func(any) error) (err error) {
	if !r.IsZero() {
		err = meth(r)
	}

	return
}

/*
LeftArc returns the string leftArc DN value assigned to the receiver, or a
zero string if unset.
*/
func (r *Spatial) LeftArc() string {
	return r.R_LeftArc
}

/*
SetLeftArc assigns the string DN value to the receiver instance.
*/
func (r *Spatial) SetLeftArc(args ...any) error {
	return writeFieldByTag(`leftArc`, r.SetLeftArc, r, args...)
}

/*
LeftArcGetFunc processes the underlying DN field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Spatial) LeftArcGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `leftArc`)
}

/*
MinArc returns the string DN value assigned to the receiver, or a zero
string if unset.
*/
func (r *Spatial) MinArc() string {
	return r.R_MinArc
}

/*
CMinArc returns the string DN value assigned to the receiver, or a zero
string if unset.
*/
func (r *Spatial) CMinArc() string {
	return r.RC_MinArc
}

/*
CMinArcGetFunc processes the underlying DN field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Spatial) CMinArcGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `c-minArc;collective`)
}

/*
SetMinArc assigns the string DN value to the receiver instance.
*/
func (r *Spatial) SetMinArc(args ...any) error {
	return writeFieldByTag(`minArc`, r.SetMinArc, r, args...)
}

/*
MinArcGetFunc processes the underlying DN field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Spatial) MinArcGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `minArc`)
}

/*
MaxArc returns the string DN value assigned to the receiver instance, or
a zero string if unset.
*/
func (r *Spatial) MaxArc() string {
	return r.R_MaxArc
}

/*
CMaxArc returns the string DN value assigned to the receiver, or a zero
string if unset.
*/
func (r *Spatial) CMaxArc() string {
	return r.RC_MaxArc
}

/*
CMaxArcGetFunc processes the underlying DN field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Spatial) CMaxArcGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `c-maxArc;collective`)
}

/*
SetMaxArc assigns the string DN value to the receiver instance.
*/
func (r *Spatial) SetMaxArc(args ...any) error {
	return writeFieldByTag(`maxArc`, r.SetMaxArc, r, args...)
}

/*
MaxArcGetFunc processes the underlying DN field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Spatial) MaxArcGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `maxArc`)
}

/*
RightArc returns the string DN value assigned to the receiver instance,
or a zero string if unset.
*/
func (r *Spatial) RightArc() string {
	return r.R_RightArc
}

/*
SetRightArc assigns the string DN value to the receiver instance.
*/
func (r *Spatial) SetRightArc(args ...any) error {
	return writeFieldByTag(`rightArc`, r.SetRightArc, r, args...)
}

/*
RightArcGetFunc processes the underlying DN field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Spatial) RightArcGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `rightArc`)
}

/*
TopArc returns the string DN value assigned to the receiver instance, or
a zero string if unset.
*/
func (r *Spatial) TopArc() string {
	return r.R_TopArc
}

/*
CTopArc returns the string DN value assigned to the receiver instance,
or a zero string if unset.
*/
func (r *Spatial) CTopArc() string {
	return r.RC_TopArc
}

/*
CTopArcGetFunc processes the underlying DN field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Spatial) CTopArcGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `c-topArc`)
}

/*
TopArcGetFunc processes the underlying DN field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Spatial) TopArcGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `topArc`)
}

/*
SetTopArc assigns the string DN value to the receiver instance.
*/
func (r *Spatial) SetTopArc(args ...any) error {
	return writeFieldByTag(`topArc`, r.SetTopArc, r, args...)
}

/*
SupArc returns the string supArc DN value assigned to the receiver, or
a zero string if unset.
*/
func (r *Spatial) SupArc() string {
	return r.R_SupArc
}

/*
CSupArc returns the string DN value assigned to the receiver instance,
or a zero string if unset.
*/
func (r *Spatial) CSupArc() string {
	return r.RC_SupArc
}

/*
CSupArcGetFunc processes the underlying DN field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Spatial) CSupArcGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `c-supArc`)
}

/*
SetSupArc assigns the string DN value to the receiver instance.
*/
func (r *Spatial) SetSupArc(args ...any) error {
	return writeFieldByTag(`supArc`, r.SetSupArc, r, args...)
}

/*
SupArcGetFunc processes the underlying DN field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Spatial) SupArcGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `supArc`)
}

/*
SubArc returns zero or more string DN slice values, each describing a
"[subArc]" value assigned to the receiver.

When populated using the search result of an LDAP Search, the values will
be DNs, however users may opt to replace the contents with raw number form
leaf values, thereby reducing client memory usage significantly.

[subArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.25
*/
func (r *Spatial) SubArc() []string {
	return r.R_SubArc
}

/*
SetSubArc appends one or more string DN values to the receiver instance.
Note that if a slice is passed as X, the destination value will be
clobbered.
*/
func (r *Spatial) SetSubArc(args ...any) error {
	return writeFieldByTag(`subArc`, r.SetSubArc, r, args...)
}

/*
SubArcGetFunc processes the underlying string DN field values through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Spatial) SubArcGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `subArc`)
}
