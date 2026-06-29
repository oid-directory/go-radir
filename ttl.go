package radir

/*
TTLPrecedence returns the highest precedence of all the provided
"time-to-live" int values:

 (lowest) -----> (highest)
   dTTL -> cTTL -> eTTL -> fTTL
                        (fallback)  

dTTL refers to the "[rATTL]" attribute value assigned explicitly
to the "[rADUAConfig]" entry, possibly derived from the RootDSE
of the rA DIT, or some other location. This value has the lowest
precedence.

cTTL refers to the "[c-rATTL]" COLLECTIVE attribute value assigned
virtually to a "[registration]" or "[registrant]" entry.

eTTL refers to the "[rATTL]" attribute value assigned explicitly
to a "[registration]" or "[registrant]" entry. This value has the
highest precedence.

fTTL refers to an optional fallback TTL value, should ALL of the
previous three values evaluate to zero (0).

Input values may be string or int. Any string value is assumed be
the string representation of a base 10 number.

If any of eTTL, cTTL or dTTL are unused -- for example if "[c-rATTL]"
is not used in the directory environment -- simply use "" or 0 in its
place.

[rATTL]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.100
[c-rATTL]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.101
[rADUAConfig]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.17
[registration]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.1
[registrant]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.16
*/
func TTLPrecedence(dTTL, cTTL, eTTL any, fTTL ...any) int {
	var p int    // precedence
	u := [4]any{dTTL,cTTL,eTTL,0}
	if len(fTTL) > 0 {
		u[3] = fTTL[0]
	}

	for i := 0; i < 4; i++ {
		if i < len(u) {
			var _p int
			switch tv := u[i].(type) {
			case int:
				_p = tv
			case string:
				_p, _ = atoi(tv)
			}

			if _p > 0 {
				if i < 3 {
					p = _p
				} else if i == 3 && p == 0 {
					// dont use fallback unless
					// precedence was always 0.
					p = _p
				}
			}
		}
	}

	return p
}

/*
TTLPrecedenceFromEntry is a convenience wrapper which reads TTL information
from the input instance of [Entry] and calls [TTLPrecedence] automatically.

The (optional) variadic fTTL argument semantics are identical to those
documented in the [TTLPrecedence] function.
*/
func TTLPrecedenceFromEntry(entry Entry, fTTL ...any) int {
	return TTLPrecedence(
		entry.Profile().TTL(),
		entry.CTTL(),
		entry.TTL(),
		fTTL...)
}
