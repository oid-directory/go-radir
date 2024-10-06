package radir

/*
ldap.go contains various values meant to aid in composing LDAP search
requests, et al. This file does not import go-ldap/v3.
*/

/*
RangeCheckFilter returns a string LDAP filter value that implements
[Section 2.2.4.1.3 of the RADUA I-D] with regards to pre-allocation
range checks, generally involving three dimensional implementations
of the I-D series.

This filter is intended to be used in conjunction with a singleLevel
(onelevel) LDAP Search Scope upon the intended parent context of the
new entry. A return of zero (0) entries following use of this filter
within a submitted LDAP Search Request would seem to indicate that no
range violation was detected, however this cannot be guaranteed if
suitable read privileges are not granted to the user for the directory
context(s) in question.

See also [RangeCheckSearchRequest].

[Section 2.2.4.1.3 of the RADUA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radua#section-2.2.4.1.3
*/
func RangeCheckSearchFilter(X string) string {
	return "(|(&(n<=" + X + ")(|(registrationRange>=" + X +
		")(registrationRange=-1)))(n=" + X + "))"
}

/*
RangeCheckSearchRequest provides default LDAP Search Request values that
are compatible as input to the [ldap/v3.NewSearchRequest] function. Note
that, depending on the nature of the request as well as the environment
in which it is made, one or more of the return values may require some
adjustment prior to use.

The input X parameter indicates the string-represented number form of the
new *[Registration] attempting to be created subordinate to P, the string
form of the intended parent DN.

See [Section 2.2.4.1.3 of the RADUA I-D] for details.

[Section 2.2.4.1.3 of the RADUA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radua#section-2.2.4.1.3
[ldap/v3.NewSearchRequest]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#NewSearchRequest
*/
func RangeCheckSearchRequest(X, P string) (dn string, scope int, typesOnly bool, filter string, at []string) {
	return P, 1, true, RangeCheckSearchFilter(X), []string{`registrationRange`}
}

const (
	DefaultSearchScope            = 0
	DefaultRegistrantSearchItem   = `(objectClass=registrant)` // DEDICATED use only!
	DefaultRegistrationSearchItem = `(objectClass=registration)`
	DefaultRootArcSearchItem      = `(objectClass=rootArc)`
	DefaultArcSearchItem          = `(objectClass=arc)`
)

/*
AttributeSelector is a convenience type which extends methods meant to
streamline the attribute selection process during LDAP Search Request
composure.

Use of instances of this type is entirely optional.
*/
type AttributeSelector struct{}

/*
All returns string slices `*` and `+`, which are used to request all
user attributes and operational attributes respectively.
*/
func (r AttributeSelector) All() []string     { return []string{`*`, `+`} }
func (r AttributeSelector) AllUser() []string { return []string{`*`} }
func (r AttributeSelector) AllOper() []string { return []string{`+`} }

/*
toLDIF implements a crude approximation of LDIF content.
*/
func toLDIF(in any) (out string) {
	if in == nil {
		return
	}

	bld := newBuilder()

	tags := getAttributeTypeFieldTags(in)
	if len(tags) == 0 {
		return
	}

	for _, tag := range tags {
		if eq(tag, `objectclass`) || eq(tag, `dn`) || eq(tag, `structuralObjectClass`) {
			// handled at higher level
			continue
		}

		if values := readFieldByTag(tag, in); len(values) > 0 {
			for i := 0; i < len(values); i++ {
				bld.WriteString(tag + `: ` + values[i])
				bld.WriteRune(10)
			}
		}
	}

	out = bld.String()

	return
}
