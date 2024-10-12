package radir

/*
extra.go contains all non-standard or supplemental elements related to
registry entries.
*/

/*
Supplement implements "[registrationSupplement]" AUXILIARY class, per
RASCHEMA § 2.5.12.

The purpose of instances of this type is to make certain miscellaneous
types available -- whether novel or derived from other standards -- for
optional assignment to *[Registration] instances.

	( 1.3.6.1.4.1.56521.101.2.5.12
	    NAME 'registrationSupplement'
	    DESC 'Supplemental registration class'
	    SUP registration AUXILIARY
	    MAY ( discloseTo $ isFrozen $ isLeafNode $
	          registrationClassification $
	          registrationCreated $
	          registrationInformation $
	          registrationModified $
	          registrationRange $
	          registrationStatus $
	          registrationURI ) )

[registrationSupplement]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.12
*/
type Supplement struct {
	R_Status   string   `ldap:"registrationStatus"`         // RASCHEMA § 2.3.14, RFC 2578 § 2
	R_Class    string   `ldap:"registrationClassification"` // RASCHEMA § 2.3.15
	R_Created  string   `ldap:"registrationCreated"`        // RASCHEMA § 2.3.11
	R_Frozen   string   `ldap:"isFrozen"`                   // RASCHEMA § 2.3.17
	R_LeafNode string   `ldap:"isLeafNode"`                 // RASCHEMA § 2.3.16
	R_Range    string   `ldap:"registrationRange"`          // RASCHEMA § 2.3.13
	R_Info     []string `ldap:"registrationInformation"`    // RASCHEMA § 2.3.9
	R_Modified []string `ldap:"registrationModified"`       // RASCHEMA § 2.3.12, RFC 2578 § 2
	R_URI      []string `ldap:"registrationURI"`            // RASCHEMA § 2.3.10

	R_DiscloseTo  []string `ldap:"discloseTo"`              // RASCHEMA § 2.3.32
	RC_DiscloseTo []string `ldap:"c-discloseTo;collective"` // RASCHEMA § 2.3.33

	r_DITProfile *DITProfile
	r_se         bool
}

/*
Supplement returns (and if needed, initializes) the embedded instance
of *[Supplement].
*/
func (r *Registration) Supplement() *Supplement {
	if r.IsZero() {
		return &Supplement{}
	}

	if r.R_Extra.IsZero() {
		r.R_Extra = new(Supplement)
		r.R_Extra.r_DITProfile = r.Profile()
	}

	return r.R_Extra
}

/*
Supplement returns (and if needed, initializes) the embedded instance
of *[Supplement].
*/
func (r *Subentry) Supplement() *Supplement {
	if r.IsZero() {
		return &Supplement{}
	}

	if r.R_Extra.IsZero() {
		r.R_Extra = new(Supplement)
		r.R_Extra.r_DITProfile = r.Profile()
		r.R_Extra.r_se = true
	}

	return r.R_Extra
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *Supplement) IsZero() bool {
	return r == nil
}

func (r *Supplement) isEmpty() bool {
	return structEmpty(r)
}

func (r *Supplement) ldif() (l string) {
	if !r.IsZero() {
		l = toLDIF(r)
	}

	return
}

/*
Unmarshal returns an instance of map[string][]string bearing the contents
of the receiver.
*/
func (r *Supplement) unmarshal() map[string][]string {
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
func (r *Supplement) marshal(meth func(any) error) (err error) {
	if !r.IsZero() {
		err = meth(r)
	}

	return
}

/*
CreateTime returns a string generalized time value assigned to the receiver
instance.
*/
func (r *Supplement) CreateTime() string {
	return r.R_Created
}

/*
CreateTimeGetFunc processes the underlying string generalized time field
value through the provided [GetOrSetFunc] instance, returning an interface
value alongside an error.
*/
func (r *Supplement) CreateTimeGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `registrationCreated`)
}

/*
SetCreateTime assigns a string generalized time value to the receiver instance.
*/
func (r *Supplement) SetCreateTime(args ...any) error {
	return writeFieldByTag(`registrationCreated`, r.SetCreateTime, r, args...)
}

/*
ModifyTime returns zero or more string generalized time values, each
reflecting a time and date at which the receiver was reported to have
been modified.
*/
func (r *Supplement) ModifyTime() (T []string) {
	return r.R_Modified
}

/*
ModifyTimeGetFunc processes the underlying string generalized time field
values through the provided [GetOrSetFunc] instance, returning an interface
value alongside an error.
*/
func (r *Supplement) ModifyTimeGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `registrationModified`)
}

/*
SetModifyTime appends one or more instances of string generalized time
values to the receiver instance. Note that if a slice is passed as X,
the destination value will be clobbered.
*/
func (r *Supplement) SetModifyTime(args ...any) error {
	return writeFieldByTag(`registrationModified`, r.SetModifyTime, r, args...)
}

/*
Range returns the string registration range terminus value assigned to
the receiver instance, or a zero string if unset.
*/
func (r *Supplement) Range() string {
	return r.R_Range
}

/*
SetRange assigns a string range terminus value, which can be any unsigned
number OR a negative -1.
*/
func (r *Supplement) SetRange(args ...any) error {
	return writeFieldByTag(`registrationRange`, r.SetRange, r, args...)
}

/*
RangeGetFunc processes the underlying range terminus field value through
the provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Supplement) RangeGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `registrationRange`)
}

/*
LeafNode returns a Boolean value indicative of whether the receiver instance
has been marked as a leaf-node, or false if unset.
*/
func (r *Supplement) LeafNode() bool {
	return eq(r.R_LeafNode, `TRUE`)
}

/*
SetLeafNode assigns the provided leaf node Boolean value to the receiver
instance.
*/
func (r *Supplement) SetLeafNode(args ...any) error {
	return writeFieldByTag(`isLeafNode`, r.SetLeafNode, r, args...)
}

/*
LeafNodeGetFunc processes the underlying leaf node field value through the
provided [GetOrSetFunc] instance, returning an interface value alongside an
error.
*/
func (r *Supplement) LeafNodeGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `isLeafNode`)
}

/*
Frozen returns a Boolean value indicative of whether the receiver has been
marked as frozen, or false if unset.
*/
func (r *Supplement) Frozen() bool {
	return eq(r.R_Frozen, `TRUE`)
}

/*
SetFrozen assigns the provided frozen Boolean value to the receiver instance.
*/
func (r *Supplement) SetFrozen(args ...any) error {
	return writeFieldByTag(`isFrozen`, r.SetFrozen, r, args...)
}

/*
FrozenGetFunc processes the underlying frozen field value through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Supplement) FrozenGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `isFrozen`)
}

/*
Classification returns a string status value that defines the classification
the receiver instance.
*/
func (r *Supplement) Classification() string {
	return r.R_Class
}

/*
SetClassification assigns a string status value to the receiver instance.
*/
func (r *Supplement) SetClassification(args ...any) error {
	return writeFieldByTag(`registrationClassification`, r.SetClassification, r, args...)
}

/*
ClassificationGetFunc processes the underlying string status field value through the
provided [GetOrSetFunc] instance, returning an interface value alongside an
error.
*/
func (r *Supplement) ClassificationGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `registrationClassification`)
}

/*
Status returns a string status value that defines the current status of
the receiver instance.
*/
func (r *Supplement) Status() string {
	return r.R_Status
}

/*
SetStatus assigns a string status value to the receiver instance.
*/
func (r *Supplement) SetStatus(args ...any) error {
	return writeFieldByTag(`registrationStatus`, r.SetStatus, r, args...)
}

/*
StatusGetFunc processes the underlying string status field value through the
provided [GetOrSetFunc] instance, returning an interface value alongside an
error.
*/
func (r *Supplement) StatusGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `registrationStatus`)
}

/*
URI returns the string URI values assigned to the receiver instance, each
defining a uniform resource identifier assigned to the registration.
*/
func (r *Supplement) URI() []string {
	return r.R_URI
}

/*
SetURI assigns one or more string URI values to the receiver instance. Note
that if a slice is passed as X, the destination value will be clobbered.
*/
func (r *Supplement) SetURI(args ...any) error {
	return writeFieldByTag(`registrationURI`, r.SetURI, r, args...)
}

/*
URIGetFunc processes the underlying URI field values through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Supplement) URIGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `registrationURI`)
}

/*
CDiscloseTo returns string DN slice values assigned to the receiver instance,
each defining an identity, group or some other DN-based reference related
to the "[c-discloseTo]" attribute type mechanics.

[c-discloseTo]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.33
*/
func (r *Supplement) CDiscloseTo() []string {
	return r.RC_DiscloseTo
}

/*
CDiscloseToGetFunc processes the underlying string DN field values through
the provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Supplement) CDiscloseToGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `c-discloseTo`)
}

/*
Info returns the string data values assigned to the receiver, each defining
an arbitrary block of textual information assigned to the registration.
*/
func (r *Supplement) Info() []string {
	return r.R_Info
}

/*
SetInfo assigns one or more string data values to the receiver instance.
Note that if a slice is passed as X, the destination value will be clobbered.
*/
func (r *Supplement) SetInfo(args ...any) error {
	return writeFieldByTag(`registrationInformation`, r.SetInfo, r, args...)
}

/*
InfoGetFunc processes the underlying data field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Supplement) InfoGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `registrationInformation`)
}

/*
DiscloseTo returns the string DN slice values assigned to the receiver
instance, each defining an identity, group or some other DN-based reference
related to the "[discloseTo]" attribute type mechanics.

[discloseTo]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.32
*/
func (r *Supplement) DiscloseTo() []string {
	return r.R_DiscloseTo
}

/*
SetDiscloseTo appends one or more string DN values to the receiver instance.
Note that if a slice is passed as X, the destination value will be clobbered.
*/
func (r *Supplement) SetDiscloseTo(args ...any) error {
	return writeFieldByTag(`discloseTo`, r.SetDiscloseTo, r, args...)
}

/*
SetCDiscloseTo appends one or more string DN values to the receiver instance.
Note that if a slice is passed as X, the destination value will be clobbered.
Also note that this method will only have an effect if executed upon an instance
which was initialized within an instance of *[Subentry].
*/
func (r *Supplement) SetCDiscloseTo(args ...any) error {
	return writeFieldByTag(`c-discloseTo;collective`, r.SetDiscloseTo, r, args...)
}

/*
DiscloseToGetFunc processes the underlying string DN field values through
the provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Supplement) DiscloseToGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `discloseTo`)
}
