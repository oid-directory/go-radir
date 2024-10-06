package radir

/*
reg.go contains Registration methods.
*/

/*
Registration contains information either to be set upon, or derived from,
an LDAP entry that describes a registration.
*/
type Registration struct {
	R_DN      string   `ldap:"dn"`
	R_GSR     string   `ldap:"governingStructureRule"`
	R_TTL     string   `ldap:"rATTL"`
	RC_TTL    string   `ldap:"c-rATTL;collective"`
	R_SOC     string   `ldap:"structuralObjectClass"`
	R_CAS     []string `ldap:"collectiveAttributeSubentries"`
	R_OC      []string `ldap:"objectClass"`
	R_Desc    []string `ldap:"description"` // effective "title" of reg
	R_Also    []string `ldap:"seeAlso"`
	R_LongArc []string `ldap:"longArc"` // only permitted for subArcs of Joint-ISO-ITU-T (2).

	R_X660    *X660       // ITU-T Rec. X.660 types
	R_X667    *X667       // ITU-T Rec. X.667 types
	R_X680    *X680       // ITU-T Rec. X.680 types
	R_X690    *X690       // ITU-T Rec. X.690 types
	R_Extra   *Supplement // Non-standard: Supplemental types
	R_Spatial *Spatial    // Non-standard: Spatial types

	R_DITProfile *DITProfile
	r_Parent     *Registration
	r_Children   *Registrations
	r_root       *registeredRoot
}

/*
registeredRoot contains information about the nature and placement of
this registration. It is populated through subsequent X.680 input, and
will (likely) be queried at any point by other constructs, such as X.660.
*/
type registeredRoot struct {
	Depth      int    // 0 = root, >=1 sub arc (default: -1)
	N          int    // 0, 1 or 2 (default: -1)
	Id         string // identifier (name) of root
	NaNF       string // nameAndNumberForm ("Id(N)") of root
	Structural string // rootArc or arc
	Auxiliary  string // iTUTRegistration, iSORegistration or jointISOITUTRegistration
}

/*
Registrations contains slices of [Registration] instances.
*/
type Registrations []*Registration

/*
Get returns an instance of *[Registration] matching the input number form
string value, or a zero instance if not found.
*/
func (r Registrations) Get(n string) (reg *Registration) {
	for i := 0; i < r.Len(); i++ {
		if elem := r[i].R_X680; elem.R_N == n ||
			elem.R_NaNF == n ||
			elem.R_ASN1Not == n ||
			elem.R_DotNot == n {
			reg = r[i]
			break
		}
	}

	return
}

/*
Walk returns an instance of *[Registration] following an attempt to traverse
the receiver instance using the input dot notation string value. A zero
instance is returned if not found.
*/
func (r *Registration) Walk(id any) (reg *Registration) {
	switch tv := id.(type) {
	case string:
		if _, a, err := cleanASN1(tv); err != nil {
			if dot := trimL(tv, `.`); isNumericOID(dot) {
				reg = r.walkN(split(dot, `.`))
			}
		} else {
			nanfs := make([][]string, 0)
			for i := 0; i < len(a); i++ {
				nanfs = append(nanfs, nanfToSlice(a[i]))
			}
			reg = r.walkASN1(nanfs)
		}
		return
	}

	return
}

func (r *Registration) walkASN1(o [][]string) (reg *Registration) {
	if r.IsZero() {
		return
	}

	nanf := mknanf(o[0])
	nf := o[0][1]

	if nanf == r.X680().NameAndNumberForm() || nf == r.X680().N() {
		if len(o) == 1 {
			reg = r
		} else {
			reg = r.walkASN1(o[1:])
		}
	} else {
		kids := r.Children()
		for _, rg := range []*Registration{
			kids.Get(nanf),
			kids.Get(nf),
		} {
			if !rg.IsZero() {
				reg = rg.walkASN1(o)
				break
			}
		}
	}

	return
}

func (r *Registration) walkN(o []string) (reg *Registration) {
	if top := o[0]; top == r.X680().N() {
		if o = o[1:]; len(o) > 0 {
			reg = r.Children().Get(o[0]).walkN(o)
		} else {
			reg = r
		}
	}

	return
}

/*
Allocate will traverse the provided dot notation string value and allocate
each sub arc along the way, assigning an X.680 Number Form following child
initialization.
*/
func (r *Registration) Allocate(oid any, ident ...string) (reg *Registration) {
	var o []string
	var nanfs [][]string

	switch tv := oid.(type) {
	case string:
		// First just check to see if it is already allocated
		if _reg := r.Walk(tv); !_reg.IsZero() {
			reg = _reg
			return
		}

		if _, a, err := cleanASN1(tv); err != nil {
			o = split(trimL(tv, `.`), `.`)
		} else {
			nanfs = make([][]string, 0)
			for i := 0; i < len(a); i++ {
				nanfs = append(nanfs, nanfToSlice(a[i]))
			}
		}
	case []string:
		if len(tv) == 0 {
			return
		}
		o = tv
	case [][]string:
		if len(tv) == 0 {
			return
		}
		nanfs = tv
	}

	if len(nanfs) > 0 {
		reg = r.allocateASN1(nanfs)
	} else if len(o) > 0 {
		reg = r.allocateDotNot(o, ident...)
	}

	return
}

func (r *Registration) allocateDotNot(o []string, ident ...string) (reg *Registration) {
	var identifier string
	if len(ident) > 0 {
		if isIdentifier(ident[0]) {
			identifier = ident[0]
		}
	}

	if o[0] == r.X680().N() {
		if len(o) == 1 {
			reg = r
			return
		}

		o = o[1:]
		if reg = r.Children().Get(o[0]); reg.IsZero() {
			if len(o) == 1 {
				reg = r.NewChild(o[0], identifier)
			} else {
				reg = r.NewChild(o[0], ``)
				return
			}
			reg = reg.allocateDotNot(o, ident...)
		} else {
			reg = reg.Allocate(o, ident...)
		}
	}

	return
}

func (r *Registration) allocateASN1(o [][]string) (reg *Registration) {
	nanf := mknanf(o[0])
	if o[0][1] == r.X680().N() || nanf == r.X680().NameAndNumberForm() {
		if len(o) == 1 {
			reg = r
			return
		}

		o = o[1:]
		id := o[0][0]
		nf := o[0][1]

		if reg = r.Children().Get(nf); reg.IsZero() {
			reg = r.NewChild(nf, id)
			if len(o) == 1 {
				return
			}
		}
		reg = reg.allocateASN1(o)
		//} else {
		//        reg = reg.Allocate(o)
		//}
	}

	return
}

/*
Contains returns a Boolean value indicative of whether an instance of
*[Registration] matching the input number form was found.
*/
func (r Registrations) Contains(n string) bool {
	return !r.Get(n).IsZero()
}

/*
SetYAxes wraps [Registrations.SetYAxes] for convenient invocation against
the underlying [Registration.Children] instance.
*/
func (r *Registration) SetYAxes(recursive ...bool) {
	if !r.IsZero() {
		r.Children().SetYAxes(recursive...)
	}
}

/*
SetXAxes wraps [Registrations.SetXAxes] for convenient invocation against
the underlying [Registration.Children] instance.
*/
func (r *Registration) SetXAxes(recursive ...bool) {
	if !r.IsZero() {
		r.Children().SetXAxes(recursive...)
	}
}

/*
SetXAxes will link ALL X-Axis (Horizontal) spatial references according
to the number form magnitude-ordered slice instances within the receiver.
This method is merely a convenient alternative to manual (and tedious)
X-Axis associations.

The recursive variadic Boolean value indicates whether the request should
span the entire progeny of X-Axis *[Registration] instances (downward),
or if this request is limited only to the receiver.

This process will result in attempting to set all of "[minArc]", "[maxArc]",
"[leftArc]" and "[rightArc]". It will not set collective variants of these types.

Note this operation has the potential to be convenient, but also quite
intensive if the receiver contains many slice instances.

See also the [Registrations.SetYAxes] method.

[leftArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.26
[rightArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.29
[minArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.27
[maxArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.30
*/
func (r *Registrations) SetXAxes(recursive ...bool) {
	var recurse bool
	if len(recursive) > 0 {
		recurse = recursive[0]
	}

	L := r.Len()
	switch {
	case L == 0:
		// Nothing to do.
		return
	case L == 1:
		// Its really not practical to set the
		// four X-axis spatial types upon less
		// than two slice members. However, if
		// recursion is requested, and if that
		// one element is a parent, then we'll
		// handle it manually.
		if single := r.Index(0); single.IsParent() && recurse {
			single.SetXAxes(recurse)
		}
		return
	}

	min := (*r)[0]
	max := (*r)[L-1]

	for i := 0; i < L; i++ {
		this := (*r)[i]
		if this.DN() == "" {
			// This absolutely cannot work.
			continue
		}

		if i > 0 {
			this.Spatial().SetLeftArc((*r)[i-1].DN())
		}

		if i < L-1 {
			this.Spatial().SetRightArc((*r)[i+1].DN())
		}

		this.setMinMax(i, L, min, max)

		if recurse && this.r_Children != nil {
			this.Children().SetXAxes(recurse)
		}
	}
}

func (r *Registration) setMinMax(i, l int, min, max *Registration) {
	if dn := min.DN(); dn != "" {
		r.Spatial().SetMinArc(dn)
		if r.Spatial().LeftArc() == "" && i == 0 {
			// In case we're at the FIRST element
			// and NO leftArc was set above ...
			r.Spatial().SetLeftArc(dn)
		}
	}

	if dn := max.DN(); dn != "" {
		r.Spatial().SetMaxArc(dn)
		if r.Spatial().RightArc() == "" && i == l-1 {
			// In case we're at the LAST element
			// and NO rightArc was set above ...
			r.Spatial().SetRightArc(dn)
		}
	}
}

/*
SetYAxes will link all Y-Axis (Vertical) spatial references according
to vertical (root, parent, child) association. This method is merely a
convenient alternative to manual (and tedious) Y-Axis associations.

The recursive variadic Boolean value indicates whether the request should
span the entire progeny of Y-Axis *[Registration] instances (downward),
or if this request is limited only to the receiver.

This process will result in attempting to set all of "[topArc]", "[supArc]",
and "[subArc]". It will not set collective variants of these types.

Note this operation has the potential to be convenient, but also quite
intensive if the receiver contains many slice instances.

See also the [Registrations.SetXAxes] method.

[supArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.21
[topArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.23
[subArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.25
*/
func (r *Registrations) SetYAxes(recursive ...bool) {
	var recurse bool
	if len(recursive) > 0 {
		recurse = recursive[0]
	}

	for i := 0; i < r.Len(); i++ {
		if reg := (*r)[i]; !reg.IsZero() {
			sdn := reg.DN()
			tdn := reg.Spatial().TopArc()

			children := reg.Children()
			LK := children.Len()
			for i := 0; i < LK; i++ {
				if child := (*children)[i]; !child.IsZero() {
					if cdn := child.DN(); cdn != "" {
						reg.Spatial().SetSubArc(cdn)
					}
					if sdn != "" {
						child.Spatial().SetSupArc(sdn)
					}
					if tdn != "" {
						child.Spatial().SetTopArc(tdn)
					}
					if recurse && child.r_Children != nil {
						child.Children().SetYAxes(recurse)
					}
				}
			}
		}
	}
}

/*
Less returns a Boolean value indicative of whether the slice at index
i is deemed to be less than the slice at j in the context of ordering.

Ordering is implemented according to individual number form magnitudes.
*/
func (r *Registrations) Less(i, j int) (less bool) {
	L := r.Len()
	if L <= i || L <= j {
		return
	}

	S1 := (*r)[i]
	S2 := (*r)[j]

	if S1.R_X680 == nil || S2.R_X680 == nil {
		return
	}

	// avoid needless initialization
	N1 := S1.R_X680.R_N
	N2 := S2.R_X680.R_N

	if N1 == "" || N2 == "" {
		return
	}

	bint1, ok1 := atobig(N1)
	bint2, ok2 := atobig(N2)

	if ok1 && ok2 {
		less = bint1.Cmp(bint2) == -1
	}

	return
}

/*
Push appends the non-zero input *[Registration] instance to the receiver
slice instance.
*/
func (r *Registrations) Push(reg *Registration) {
	if !r.IsZero() {
		if !structEmpty(reg) {
			*r = append(*r, reg)
		}
	}
}

/*
Len returns the integer length of the receiver instance.
*/
func (r *Registrations) Len() (l int) {
	if !r.IsZero() {
		l = len(*r)
	}

	return
}

/*
Index returns the Nth *[Registration] instance within the receiver instance,
or a zero instance if not found.
*/
func (r *Registrations) Index(idx int) (reg *Registration) {
	if !r.IsZero() {
		if 0 <= idx && idx < r.Len() {
			reg = (*r)[idx]
		}
	}

	return
}

/*
IsParent returns a Boolean value indicative of the receiver containing
one or more child *[Registration] instances.

See also [Registrations.HasParents].
*/
func (r *Registration) IsParent() (is bool) {
	if !r.IsZero() {
		is = r.Children().Len() > 0
	}

	return
}

/*
HasParents returns a Boolean value indicative of the receiver instance
containing one or more *[Registration] instances who are parents themselves.

See also [Registration.IsParent].
*/
func (r *Registrations) HasParents() bool {
	parents := 0
	if !r.IsZero() {
		for i := 0; i < r.Len(); i++ {
			if r.Index(i).IsParent() {
				parents++
			}
		}
	}

	return parents > 0
}

/*
Swap implements the func(int,int) signature required by the [sort.Interface]
signature. The result is the swapping of the specified receiver slice indices.
*/
func (r *Registrations) Swap(i, j int) {
	L := r.Len()
	if (-1 < i && i < L) && (-1 < j && j < L) && i != j {
		(*r)[i], (*r)[j] = (*r)[j], (*r)[i]
	}
}

/*
SortByNumberForm wraps *[Registrations.SortByNumberForm] for convenient
invocation of [sort.Stable] sorting of any underlying *[Registration]
(child) instances.

The recursive variadic Boolean value indicates whether the request should
span the entire progeny of *[Registration] instances (downward), or if this
request is limited only to the receiver.

Note that this particular wrapper serves no purpose when not executed with
positive recursion.
*/
func (r *Registration) SortByNumberForm(recursive ...bool) {
	if !r.IsZero() {
		if L := r.Children().Len(); L != 0 {
			r.Children().SortByNumberForm(recursive...)
		}
	}
}

/*
SortByNumberForm executes [sort.Stable] to sort the contents of the receiver
slice instance according to NumberForm magnitude, ordered lowest to highest.

The recursive variadic Boolean value indicates whether the request should
span the entire progeny of *[Registration] instances (downward), or if this
request is limited only to the receiver.
*/
func (r *Registrations) SortByNumberForm(recursive ...bool) {
	if L := r.Len(); L > 0 {
		stabSort(r)

		if len(recursive) > 0 {
			if recurse := recursive[0]; recurse {
				// recurse through each slice reg and sort their
				// children, descending indefinitely until done.
				for i := 0; i < L; i++ {
					r.Index(i).SortByNumberForm(recursive...)
				}
			}
		}
	}
}

/*
Parent returns the underlying instance of *[Registration] present within
the receiver, or a zero instance if unset.

The parent instance is set automatically through use of the [Registration.NewChild]
method.
*/
func (r *Registration) Parent() (reg *Registration) {
	if !r.IsZero() {
		reg = r.r_Parent
	}

	return
}

/*
Size traverses the extent of the underlying OID tree and returns the
integer total of all non-nil *[Registration] instances found en route.

Note that this process does not traverse "upwards" -- only "downwards".
This means that for a COMPLETE total of an entire root's OIDs, this
method should be executed upon said root *[Registration] instance.

Otherwise, it will only count the number of instances from that point
downward.
*/
func (r *Registration) Size() (size int) {
	if !r.IsZero() {
		size++
		K := r.Children()
		LK := K.Len()
		for i := 0; i < LK; i++ {
			if sub := (*K)[i]; !sub.IsZero() {
				size += sub.Size()
			}
		}
	}

	return
}

/*
Children returns the underlying instance of *[Registrations] present within
the receiver's child slice instance, or a zero instance if unset.

The child slice instances are set automatically through use of the
[Registration.NewChild] method.
*/
func (r *Registration) Children() (regs *Registrations) {
	if !r.IsZero() {
		if r.r_Children == nil {
			k := make(Registrations, 0)
			r.r_Children = &k
		}

		regs = r.r_Children
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *Registration) IsZero() bool {
	return r == nil
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *Registrations) IsZero() bool {
	return r == nil
}

func (r *Registration) isEmpty() bool {
	return structEmpty(r)
}

/*
StructuralObjectClass returns the string-based STRUCTURAL object class,
or a zero string if unset.

Note this value is not specified manually by users.
*/
func (r *Registration) StructuralObjectClass() (soc string) {
	if !r.IsZero() {
		soc = r.R_SOC
	}

	return
}

/*
StructuralObjectClassGetFunc executes the [GetOrSetFunc] instance and
returns its own return values. The 'any' value will require type assertion
in order to be accessed reliably. An error is returned if issues arise.
*/
func (r *Registration) StructuralObjectClassGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `structuralObjectClass`)
}

/*
CollectiveAttributeSubentries returns one or more LDAP distinguished names
which identify all "[collectiveAttributeSubentries]" references that serve
to populate the *[Registration] entry.

Note this value is not specified manually by users.

[collectiveAttributeSubentries]: https://www.rfc-editor.org/rfc/rfc3671.html#section-2.2
*/
func (r *Registration) CollectiveAttributeSubentries() (cas []string) {
	if !r.IsZero() {
		cas = r.R_CAS
	}

	return
}

/*
CollectiveAttributeSubentriesGetFunc executes the [GetOrSetFunc] instance
and returns its own return values. The 'any' value will require type assertion
in order to be accessed reliably. An error is returned if issues arise.
*/
func (r *Registration) CollectiveAttributeSubentriesGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `collectiveAttributeSubentries`)
}

/*
DN returns the string-based LDAP Distinguished Name value, or a zero
string if unset.
*/
func (r *Registration) DN() (dn string) {
	if !r.IsZero() {
		dn = r.R_DN
	}

	return
}

/*
DNGetFunc executes the [GetOrSetFunc] instance and returns its own
return values. The 'any' value will require type assertion in order
to be accessed reliably. An error is returned if issues arise.
*/
func (r *Registration) DNGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `dn`)
}

/*
SetDN assigns the first provided value to the receiver instance as a string DN.

If additional arguments are present, the second value will be asserted as an
instance of [GetOrSetFunc], returning an error if this fails.
*/
func (r *Registration) SetDN(args ...any) error {
	return writeFieldByTag(`dn`, r.SetDN, r, args...)
}

/*
GoverningStructureRule returns the "[governingStructureRule]" assigned to
the receiver instance.

Note the "[governingStructureRule]" type is operational, and cannot be set
by the end user. It is also not collective.

Use of this method is only meaningful in environments which employ one or
more "[dITStructureRule]" definitions.

[governingStructureRule]: https://datatracker.ietf.org/doc/html/rfc4512#section-3.4.6
[dITStructureRule]: https://datatracker.ietf.org/doc/html/rfc4512#section-4.1.7.1
*/
func (r *Registration) GoverningStructureRule() (gsr string) {
	if !r.IsZero() {
		gsr = r.R_GSR
	}

	return
}

/*
GoverningStructureRuleGetFunc executes the [GetOrSetFunc] instance and
returns its own return values. The 'any' value will require type assertion
in order to be accessed reliably. An error is returned if issues arise.
*/
func (r *Registration) GoverningStructureRuleGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `governingStructureRule`)
}

/*
Description returns one or more string description values assigned to the
receiver instance.
*/
func (r *Registration) Description() (desc []string) {
	if !r.IsZero() {
		desc = r.R_Desc
	}

	return
}

/*
SetDescription appends one or more string description values to the receiver
instance. Note that if a slice is passed as X, the destination value will be
clobbered.
*/
func (r *Registration) SetDescription(args ...any) error {
	return writeFieldByTag(`description`, r.SetDescription, r, args...)
}

/*
DescriptionGetFunc processes the underlying field value(s) through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Registration) DescriptionGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `description`)
}

/*
SeeAlso returns the string DN values assigned to the receiver instance.
*/
func (r *Registration) SeeAlso() (also []string) {
	if !r.IsZero() {
		also = r.R_Also
	}

	return
}

/*
SetSeeAlso appends one or more string DN values to the receiver instance.
Note that if a slice is passed as X, the destination value will be clobbered.
*/
func (r *Registration) SetSeeAlso(args ...any) error {
	return writeFieldByTag(`seeAlso`, r.SetSeeAlso, r, args...)
}

/*
SeeAlsoGetFunc processes the underlying string DN field value(s) through
the provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Registration) SeeAlsoGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `seeAlso`)
}

/*
TTL returns the effective time-to-live for the receiver instance, taking
into account *[DITProfile]-inherited values as well as any subtree-based
(COLLECTIVE) and entry literal values. The output can be used to instruct
instances of [Cache] when, and when not, to cache an instance.

See [Section 2.2.3.4 of the RADUA I-D] for details related to TTL precedence.

[Section 2.2.3.4 of the RADUA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radua#section-2.2.3.4
*/
func (r *Registration) TTL() string {
	ct := r.DITProfile().TTL()
	lt := selectTTL(r.R_TTL, r.RC_TTL)

	if lt == `` {
		// If no localized TTL or COLLECTIVE TTL, then
		// just return the DITProfile-based TTL.
		return ct
	}

	return lt
}

/*
SetTTL assigns the provided string TTL value to the receiver instance.
*/
func (r *Registration) SetTTL(args ...any) error {
	return writeFieldByTag(`rATTL`, r.SetTTL, r, args...)
}

/*
TTLGetFunc processes the underlying TTL field value through the provided
[GetOrSetFunc] instance, returning an interface value alongside an error.
*/
func (r *Registration) TTLGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `rATTL`)
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
func (r *Registration) Marshal(meth func(any) error) (err error) {
	if meth == nil {
		err = NilMethodErr
		return
	} else if r == nil {
		err = NilRegistrationErr
		return
	}

	for _, err = range []error{
		meth(r),
		r.X660().marshal(meth),
		r.X667().marshal(meth),
		r.X680().marshal(meth),
		r.X690().marshal(meth),
		r.Supplement().marshal(meth),
		r.Spatial().marshal(meth),
	} {
		if err != nil {
			return
		}
	}

	// clean-up any duplicate objectClass
	// slices that may have appeared, and
	// remove any classes that aren't used.
	r.refreshObjectClasses()

	return
}

/*
Unmarshal transports values from the receiver into an instance of
map[string][]string, which can subsequently be fed to the [go-ldap/v3
NewEntry] function.

[go-ldap/v3 NewEntry]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#NewEntry
*/
func (r *Registration) Unmarshal() (outer map[string][]string) {
	if r.IsZero() {
		return
	}

	outer = make(map[string][]string)

	for _, inner := range []map[string][]string{
		unmarshalStruct(r, outer),
		r.X660().unmarshal(),
		r.X667().unmarshal(),
		r.X680().unmarshal(),
		r.X690().unmarshal(),
		r.Supplement().unmarshal(),
		r.Spatial().unmarshal(),
	} {
		if inner != nil {
			for k, v := range inner {
				outer[k] = v
			}
		}
	}

	return
}

/*
Unmarshal is a convenience method that returns slices of map[string][]string
instances, each representative of an individual *[Registration] instance
that was deemed valid for unmarshaling to map[string][]string.
*/
func (r *Registrations) Unmarshal() (maps []map[string][]string) {
	maps = make([]map[string][]string, 0)
	for i := 0; i < len(*r); i++ {
		if um := (*r)[i].Unmarshal(); len(um) > 0 {
			maps = append(maps, um)
		}
	}
	return
}

/*
Marshal returns an error following an attempt to execute all input method
[go-ldap/v3 Entry.Unmarshal] signatures. The result is a sequence of marshaled
*[Registration] instances being added to the receiver instance.

The input *[DITProfile] value will be used to initialize each *[Registration]
instance using the appropriate configuration.

[go-ldap/v3 Entry.Unmarshal]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry.Unmarshal
*/
func (r *Registrations) Marshal(profile *DITProfile, meths ...func(any) error) (err error) {
	if !profile.Valid() {
		err = DUAConfigValidityErr
		return
	}

	for i := 0; i < len(meths) && err == nil; i++ {
		reg := profile.NewRegistration()
		if err = reg.Marshal(meths[i]); err == nil {
			*r = append(*r, reg)
		}
	}

	return
}

/*
ObjectClasses returns string slices of "[objectClass]" in descriptor or
numeric OID forms currently held by the receiver instance.

[objectClass]: https://www.rfc-editor.org/rfc/rfc4512.html#section-3.3
*/
func (r *Registration) ObjectClasses() []string {
	return r.R_OC
}

/*
SetObjectClasses appends one or more string descriptor or numeric OID
"[objectClass]" values to the receiver instance. Note that if a slice is
passed as X, the destination value will be clobbered.

[objectClass]: https://www.rfc-editor.org/rfc/rfc4512.html#section-3.3
*/
func (r *Registration) SetObjectClasses(args ...any) error {
	return writeFieldByTag(`objectClass`, r.SetObjectClasses, r, args...)
}

/*
ObjectClassesGetFunc processes the underlying field value(s) through the
provided [GetOrSetFunc] instance, returning an interface value alongside
an error.
*/
func (r *Registration) ObjectClassesGetFunc(getfunc GetOrSetFunc) (any, error) {
	return getFieldValueByNameTagAndGoSF(r, getfunc, `objectClass`)
}

/*
Kind returns the static string value `STRUCTURAL` merely as a convenient
means to determine what [kind of "objectClass"] the receiver describes.

[kind of "objectClass"]: https://www.rfc-editor.org/rfc/rfc4512.html#section-4.1.1
*/
func (r *Registration) Kind() string {
	return `STRUCTURAL`
}

/*
IsRoot returns a Boolean value indicative of whether the receiver instance
represents an official root registration, as indicated by presence of the
"[rootArc]" STRUCTURAL class (in descriptor or numeric OID).

See also the [Registration.IsNonRoot] method.

[rootArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.2
*/
func (r *Registration) IsRoot() (is bool) {
	if !r.IsZero() {
		rn := []string{`0`, `1`, `2`}
		roid := `1.3.6.1.4.1.56521.101.2.5.2`
		if is = eq(r.R_SOC, `rootarc`) || eq(r.R_SOC, roid); is {
			is = strInSlice(r.X680().N(), rn)
		}
	}

	return
}

/*
IsNonRoot returns a Boolean value indicative of the presence of the "[arc]"
STRUCTURAL object class (in descriptor or numeric OID form) within the
receiver instance.

See also the [Registration.IsRoot] method.

[arc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.3
*/
func (r *Registration) IsNonRoot() (is bool) {
	if !r.IsZero() {
		roid := `1.3.6.1.4.1.56521.101.2.5.3`
		is = eq(r.R_SOC, `arc`) || eq(r.R_SOC, roid)
	}

	return
}

/*
DITProfile returns the *[DITProfile] instance assigned to the receiver,
if set, else nil is returned.
*/
func (r *Registration) DITProfile() (prof *DITProfile) {
	if !r.IsZero() {
		if prof = r.R_DITProfile; !prof.Valid() {
			prof = &DITProfile{}
		}
	}
	return
}

/*
SetDITProfile assigns *[DITProfile] d to the receiver instance if, and
only if, both of the following evaluate as true:

  - [DITProfile.Valid] returns true for d, AND ...
  - the DN of the receiver instance matches a *[Registration] search base within d

Case is not significant in the matching process.
*/
func (r *Registration) SetDITProfile(d *DITProfile) {
	if d.Valid() {
		if idx := d.RegistrationSuffixEqual(r.DN()); idx != -1 {
			r.R_DITProfile = d
		}
	}
}

/*
Combined returns a Boolean value indicative of the "Combined Registrants
Policy" being in-force.

If the underlying instance of *[DITProfile] within the receiver instance
is zero, false is returned.
*/
func (r *Registration) Combined() bool {
	if !r.IsZero() {
		if dp := r.DITProfile(); !dp.IsZero() {
			return dp.Combined()
		}
	}

	return false
}

/*
Dedicated returns a Boolean value indicative of the "Dedicated Registrants
Policy" being in-force.

If the underlying instance of *[DITProfile] within the receiver instance
is zero, false is returned.
*/
func (r *Registration) Dedicated() bool {
	if !r.IsZero() {
		if dp := r.DITProfile(); !dp.IsZero() {
			return dp.Dedicated()
		}
	}

	return false
}

func (r *Registration) refreshObjectClasses() {
	bools := []bool{
		r.X660().isEmpty(),
		r.X667().isEmpty(),
		r.X680().isEmpty(),
		r.X690().isEmpty(),
		r.Supplement().isEmpty(),
		r.Spatial().isEmpty(),
	}

	for tx, tag := range []string{
		`x660Context`,
		`x667Context`,
		`x680Context`,
		`x690Context`,
		`registrationSupplement`,
		`spatialContext`,
	} {
		if bools[tx] {
			r.R_OC = removeStrInSlice(tag, r.R_OC)
		} else {
			if !strInSlice(tag, r.R_OC) {
				r.R_OC = append(r.R_OC, tag)
			}
		}
	}

	// go-ldap/v3.Entry.Unmarshal is sloppy about adding
	// duplicate objectClasses, so let's clean up any
	// that may have appeared.
	var tmp []string
	for _, oc := range r.R_OC {
		if !strInSlice(oc, tmp) {
			tmp = append(tmp, oc)
		}
	}
	if len(tmp) != len(r.R_OC) {
		r.R_OC = tmp
	}
}

/*
Root returns the official integer value for any of the three (3) possible
roots alongside an AUXILIARY class descriptor. The possible integer roots
are as follows:

  - itu-t (0)
  - iso (1)
  - joint-iso-itut (2)

If the receiver instance is determined to be associated with a root
*[Registration] instance, no descriptor class is returned.

If the receiver is not yet populated in this regard, a -1 integer is
returned with no descriptor class.

Otherwise, for any non-root *[Registration], one (1) of the following
AUXILIARY class descriptors will be returned:

  - "[iTUTRegistration]"
  - "[iSORegistration]"
  - "[jointISOITUTRegistration]"

[iTUTRegistration]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.8
[iSORegistration]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.9
[jointISOITUTRegistration]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.10
*/
func (r *Registration) Root() (n int, class string) {
	n = -1
	if !r.IsZero() {
		if !structEmpty(r.r_root) {
			n = r.r_root.N
			if r.r_root.Depth > 1 {
				class = r.r_root.Auxiliary
			}
		}
	}

	return
}

/*
LDIF returns the string LDIF form of the receiver instance. Note that
this is a crude approximation of LDIF and should ideally be parsed
through a reliable LDIF parser such as [go-ldap/ldif] to verify integrity.

[go-ldap/ldif]: https://pkg.go.dev/github.com/go-ldap/ldif
*/
func (r *Registration) LDIF() (l string) {
	if !r.IsZero() {
		dn := readFieldByTag(`dn`, r)
		if len(dn) > 0 {
			r.refreshObjectClasses()

			oc := readFieldByTag(`objectClass`, r)
			bld := newBuilder()

			bld.WriteString(`dn: ` + dn[0])
			bld.WriteRune(10)

			for i := 0; i < len(oc); i++ {
				bld.WriteString(`objectClass: ` + oc[i])
				bld.WriteRune(10)
			}

			bld.WriteString(toLDIF(r))
			bld.WriteString(r.X660().ldif())
			bld.WriteString(r.X667().ldif())
			bld.WriteString(r.X680().ldif())
			bld.WriteString(r.X690().ldif())
			bld.WriteString(r.Spatial().ldif())
			bld.WriteString(r.Supplement().ldif())

			l = bld.String()
		}

		if len(l) > 0 {
			// if the very last character is
			// a newline, remove it.
			if rune(l[len(l)-1]) == rune(10) {
				l = l[:len(l)-1]
			}
		}
	}

	return
}

/*
LDIFs traverses the extent of the underlying OID tree and returns the
content as an LDIF payload containing all non-nil *[Registration]
instances found en route.

Note that this process does not traverse "upwards" -- only "downwards".
This means that for a COMPLETE set of an entire root's LDIFs, this
method should be executed upon said root *[Registration] instance.

Otherwise, it will only write the instances from that point downward.
*/
func (r *Registration) LDIFs() (l string) {
	if !r.IsZero() {
		_l := newBuilder()
		_l.WriteString(r.LDIF())
		_l.WriteRune(10)
		K := r.Children()
		LK := K.Len()
		for i := 0; i < LK; i++ {
			if sub := (*K)[i]; !sub.IsZero() {
				_l.WriteString(sub.LDIFs())
				_l.WriteRune(10)
			}
		}
		l = _l.String()
	}

	return
}

/*
NewChild initializes a new instance of *[Registration] bearing superior
values to that of the receiver. If successful, the return value will be
automatically added to the underlying children *[Registrations] instance.
*/
func (r *Registration) NewChild(nf, id string) (s *Registration) {
	if !r.IsZero() {
		ident, nanf, dotp, _, ok := r.sibOrSub(nf, id, false)
		if !ok {
			return
		}

		switch {
		case r.IsRoot():
			// Start a new dotNotation, since the
			// source (root) wouldn't have one to
			// begin with ...
			dotp = r.X680().N()
		case r.IsNonRoot():
		default:
			return
		}

		dotp += `.` + nf // complete the new dotNotation

		var oiv string
		_s := r.DITProfile().NewRegistration()
		if a1 := r.X680().ASN1Notation(); len(nanf) > 0 && len(a1) > 0 {
			oiv = trimR(a1, `}`) + ` ` + nanf + `}`
		}

		var this string = r.DN()
		var cdn string
		if r.DITProfile().Model() == TwoDimensional {
			if idx := idxr(this, ','); idx != -1 {
				cdn = `dotNotation=` + dotp + `,` + this[idx+1:]
			}
		} else {
			cdn = `n=` + nf + `,` + r.DN()
		}

		// Use the source's objectClass slices as a template,
		// but swap rootArc for arc.
		_s.SetObjectClasses(removeStrInSlice(`rootArc`, r.ObjectClasses()))
		_s.SetObjectClasses(`arc`)
		_s.R_SOC = `arc`

		for val, funk := range map[string]func(...any) error{
			ident: _s.X680().SetIdentifier,
			nanf:  _s.X680().SetNameAndNumberForm,
			nf:    _s.X680().SetN,
			dotp:  _s.X680().SetDotNotation,
			oiv:   _s.X680().SetASN1Notation,
			cdn:   _s.SetDN,
		} {
			if len(val) > 0 {
				funk(val)
			}
		}

		_s.r_Parent = r
		s = _s

		r.Children().Push(s)
	}

	return
}

/*
NewSibling initializes a new instance of *[Registration] bearing similar
or "parallel" values to those held by the receiver.

The input value n reflects the desired number form to be held by the new
*[Registration] instance. If the value is not a number, or is identical
to that held by the receiver instance, a nil instance is returned. If
the receiver lacks a DN, a nil instance is returned.

In the case of non-root instances, this value will also serve as the leaf
"[iRI]" component, if defined within the (source) receiver instance.

The input value id reflects the desired identifier, or name form, to be
held by the new *[Registration] instance. In turn this also reflects the
"[nameAndNumberForm]" identifier component to be set.

In the case of non-root instances, the identifier will also serve as the
"[aSN1Notation]" leaf identifier component. In the case of root instances,
the identifier serves as the sole "[aSN1Notation]" component, identical
to the "[nameAndNumberForm]" value.

If the id input value is zero length, all of the above identifier handling
procedures will be skipped. A valid instance will still be returned. However
if a non-compliant identifier value is passed, a nil instance is returned.

If the receiver instance possesses any of "[supArc]", "[c-supArc]", "[topArc]",
"[c-topArc]", "[minArc]", "[c-minArc]", "[maxArc]" or "[c-maxArc]" attribute
values, they will be transferred automatically, as these are common to all
relative siblings.

Use of this method is merely a convenient alternative to manual composition
of a new instance, but will still require additional configuration for cases
in which the appropriate values cannot be "extrapolated" using receiver input,
such as "[leftArc]", "[rightArc]", "[subArc]" and others.
*/
func (r *Registration) NewSibling(nf, id string) (s *Registration) {
	if !r.IsZero() {
		ident, nanf, dotp, dnp, ok := r.sibOrSub(nf, id, true)
		if !ok {
			return
		}

		var oiv string
		var _s *Registration
		switch {
		case r.IsRoot():
			dotp = ``
			_s = r.DITProfile().NewRegistration(true)
			if len(nanf) > 0 {
				oiv = `{` + nanf + `}`
			}
		case r.IsNonRoot():
			dotp += `.` + nf
			_s = r.DITProfile().NewRegistration()
			a1 := r.X680().ASN1Notation()
			if len(nanf) > 0 && len(a1) > 0 {
				poiv := split(a1[1:len(a1)-1], ` `)
				oiv = `{` + join(poiv[:len(poiv)-1], ` `) + ` ` + nanf + `}`
			}
		default:
			return
		}

		sdn := `n=` + nf + `,` + dnp

		_s.SetObjectClasses(r.ObjectClasses())
		_s.R_SOC = r.R_SOC // structuralObjectClass

		for val, funk := range map[string]func(...any) error{
			ident: _s.X680().SetIdentifier,
			nanf:  _s.X680().SetNameAndNumberForm,
			nf:    _s.X680().SetN,
			dotp:  _s.X680().SetDotNotation,
			oiv:   _s.X680().SetASN1Notation,
			sdn:   _s.SetDN,
		} {
			if len(val) > 0 {
				funk(val)
			}
		}

		// Take any common spatial types
		if !r.R_Spatial.IsZero() {
			_s.Spatial().SetSupArc(r.Spatial().SupArc())
			_s.Spatial().SetTopArc(r.Spatial().TopArc())
			_s.Spatial().SetMinArc(r.Spatial().MinArc())
			_s.Spatial().SetMaxArc(r.Spatial().MaxArc())
		}
		s = _s
	}

	return
}

func (r *Registration) sibOrSub(nf, id string, sib bool) (ident, nanf, dotp, dnp string, ok bool) {
	if len(r.DN()) == 0 {
		// no DN, no service
		return
	}

	if !isNumber(nf) || (nf == r.X680().N() && sib) {
		// n is not a number, OR it is
		// identical to receiver's n.
		return
	}

	if len(id) > 0 {
		if !isIdentifier(id) {
			// bogus non-zero identifier
			return
		}
		ident = id
		nanf = id + `(` + nf + `)`
	}

	if dot := r.X680().DotNotation(); len(dot) > 2 {
		sp := split(dot, `.`)
		if sib {
			dotp = join(sp[:len(sp)-1], `.`)
		} else {
			dotp = dot
		}
	}

	// Set the DN last
	var bdn string = r.DN()
	if x := idxr(bdn, ','); x != -1 {
		dnp = bdn[x+1:]
	}

	ok = len(dnp) > 0

	return
}
