package radir

/*
Map contains a 2D abstraction of a *[Registration] or *[Registrant].
Instances of this type are produced via the [Registration.Map] and
[Registrant.Map] methods.

For general values, no nesting occurs. All user content is "flattened",
and all values are string, []string and bool. The only exception to
this is the possible presence of a *[DITProfile] associated with the
"dITProfile" string key.

The key names for each entry are derived from the respective struct
field's "ldap" tag. For instance `ldap:"aSN1Notation"` equates to:

  m[`aSN1Notation`] = "{iso(1) identified-organization(3) ... }"

Field tag names bearing a ";tagName" suffix will be truncated
to just the proper attribute name. For instance, "c-rATTL;collective"
would manifest simply as "c-rATTL".

Case is significant with regards to key indexing.

No private fields are included during creation of instances of this type.

Use of this type may be of interest to those who wish to augment,
simplify or enrich a *[Registration] or *[Registrant] entry's
contents in some way -- without having to alter the original --
as well as those who prefer a simpler type with which to interact.
*/
type Map map[string]any

/*
IsZero returns a Boolean value indicative of a nil or zero-length
receiver instance.
*/
func (r Map) IsZero() bool { return r == nil || len(r) == 0 }

/*
DN returns the string representation of the distinguished name
assigned to the receiver, or an empty string.
*/
func (r Map) DN() string {
	var dn string
	if !r.IsZero() {
		dn, _ = r["dn"].(string)
	}
	return dn
}

/*
TTL returns the "[rATTL]" value assigned to the receiver, or
an empty string.

[rATTL]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.100
*/
func (r Map) TTL() string {
	var ttl string
	if !r.IsZero() {
		ttl, _ = r["rATTL"].(string)
	}
	return ttl
}

/*
CTTL returns the "[c-rATTL]" value assigned to the receiver, or
an empty string.

[c-rATTL]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.101
*/
func (r Map) CTTL() string {
	var ttl string
	if !r.IsZero() {
		ttl, _ = r["c-rATTL"].(string)
	}
	return ttl
}

/*
Profile returns an instance of *[DITProfile] assigned to the receiver instance
under the "dITProfile" index name.
*/
func (r Map) Profile() *DITProfile {
	var p *DITProfile = &DITProfile{}
	if !r.IsZero() {
		if _p, _ok := r[`dITProfile`]; _ok && _p != nil {
			p, _ = _p.(*DITProfile)
		}
	}
	return p
}

/*
Contains returns a Boolean value indicative of a successful
match of n as a map index.
*/
func (r Map) Contains(n string) bool {
	var c bool
	if !r.IsZero() {
		_, c = r[n]
	}
	return c
}

/*
Registration returns a Boolean value indicative of whether the receiver
instance reflects a *[radir.Registration] instance. This determination
is made based on the presence of the "[registration]" ABSTRACT object
class.

[registration]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.1
*/
func (r Map) Registration() bool {
	var reg bool
	if !r.IsZero() {
		if s, ok := r.StringsValue(`objectClass`); ok {
			for i := 0; i < len(s) && !reg; i++ {
				reg = s[i] == `registration` 
			}
		}
	}
	return reg
}

/*
Subentry returns a Boolean value indicative of whether the receiver
instance reflects a *[radir.Subentry] instance. This determination
is made based on the presence of the "[subentry]" STRUCTURAL object
class.

[subentry]: https://www.rfc-editor.org/info/rfc3672/#section-2.4
*/
func (r Map) Subentry() bool {
	var se bool
	if !r.IsZero() {
		if s, ok := r.StringsValue(`objectClass`); ok {
			for i := 0; i < len(s) && !se; i++ {
				se = s[i] == `subentry` 
			}
		}
	}
	return se
}

/*
Registrant returns a Boolean value indicative of whether the receiver
instance reflects a *[radir.Registant] instance. This determination
is made based on the presence of the "[registrant]" STRUCTURAL object
class.

[registrant]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.16
*/
func (r Map) Registrant() bool {
        var reg bool
	if !r.IsZero() {
        	if s, ok := r.StringsValue(`objectClass`); ok {
        	        for i := 0; i < len(s) && !reg; i++ {
        	                reg = s[i] == `registrant` 
        	        }
        	}
	}
        return reg
}

/*
StringValue returns a string alongside a Boolean value. This method
is used to call an index (n) and type assert as a non zero string.
A value of true being returned indicates both steps were successful. 
*/
func (r Map) StringValue(n string) (string, bool) {
	var s string
	var ok bool

	if !r.IsZero() {
		if _s, _ok := r[n]; _ok {
			if s, _ok = _s.(string); len(s) > 0 && _ok {
				ok = true
			}
		}
	}
	
	return s, ok
}

/*
StringsValue returns a []string instance alongside a Boolean value.
This method is used to call an index (n) and type assert as a non
zero and non nil []string. A value of true being returned indicates
both steps were successful. 
*/
func (r Map) StringsValue(n string) ([]string, bool) {
	var s []string
	var ok bool

	if !r.IsZero() {
		if _s, _ok := r[n]; _ok {
			if s, _ok = _s.([]string); len(s) > 0 && _ok {
				ok = true
			}
		}
	}
	
	return s, ok
}

/*
BoolValue returns two Boolean values. This method is used to call an
index (n) and type assert as a bool. The second value confirms both
steps were successful, while the first is the actual value instance
of interest.

Note that this method is only used for "[isLeafNode]" and "[isFrozen]"
*[Registration] values.

[isLeafNode]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.16
[isFrozen]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.17
*/
func (r Map) BoolValue(n string) (bool, bool) {
	var b, ok bool

	if !r.IsZero() {
		if _b, _ok := r[n]; _ok {
			if b, _ok = _b.(bool); _ok {
				ok = true
			}
		}
	}
	
	return b, ok
}

