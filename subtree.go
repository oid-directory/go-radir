package radir

import (
	"errors"
)

/*
SubtreeSpecification implements the Subtree Specification construct.

From [Appendix A of RFC 3672]:

	SubtreeSpecification = "{" [ sp ss-base ]
	                           [ sep sp ss-specificExclusions ]
	                           [ sep sp ss-minimum ]
	                           [ sep sp ss-maximum ]
	                           [ sep sp ss-specificationFilter ]
	                                sp "}"

	ss-base                = id-base                msp LocalName
	ss-specificExclusions  = id-specificExclusions  msp SpecificExclusions
	ss-minimum             = id-minimum             msp BaseDistance
	ss-maximum             = id-maximum             msp BaseDistance
	ss-specificationFilter = id-specificationFilter msp Refinement

	BaseDistance = INTEGER-0-MAX

From [§ 6 of RFC 3642]:

	LocalName         = RDNSequence
	RDNSequence       = dquote *SafeUTF8Character dquote

	INTEGER-0-MAX   = "0" / positive-number
	positive-number = non-zero-digit *decimal-digit

	sp  =  *%x20  ; zero, one or more space characters
	msp = 1*%x20  ; one or more space characters
	sep = [ "," ]

	OBJECT-IDENTIFIER = numeric-oid / descr
	numeric-oid       = oid-component 1*( "." oid-component )
	oid-component     = "0" / positive-number

[Appendix A of RFC 3672]: https://datatracker.ietf.org/doc/html/rfc3672#appendix-A
[§ 6 of RFC 3642]: https://datatracker.ietf.org/doc/html/rfc3672#section-6
*/
type SubtreeSpecification struct {
	Base                LocalName
	SpecificExclusions  SpecificExclusions
	Min                 BaseDistance
	Max                 BaseDistance
	SpecificationFilter Refinement
}

/*
NewSubtreeSpecification returns an instance of [SubtreeSpecification] alongside
an error following an attempt to parse x.
*/
func NewSubtreeSpecification(raw string) (ss SubtreeSpecification, err error) {
	if len(raw) < 2 {
		return
	}

	if raw[0] != '{' || raw[len(raw)-1] != '}' {
		err = errors.New("SubtreeSpecification {} encapsulation error")
		return
	}

	if raw = trimS(raw[1 : len(raw)-1]); raw == `{}` {
		return
	}

	var ranges map[string][]int = make(map[string][]int, 0)

	var pos int
	if begin := idxs(raw, `base `); begin != -1 {
		var end int
		begin += 5
		if ss.Base, end, err = subtreeBase(raw[begin:]); err != nil {
			return
		}
		pos += begin
		end += pos + 1
		ranges[`base`] = []int{begin, end}
	}

	if begin := idxs(raw, `specificExclusions `); begin != -1 {
		var end int
		begin += 19
		if ss.SpecificExclusions, end, err = subtreeExclusions(raw, begin); err != nil {
			return
		}
		end = begin + end
		ranges[`specificExclusions`] = []int{begin, end}
	}

	if begin := idxs(raw, `minimum `); begin != -1 {
		var end int
		begin += 8
		if ss.Min, end, err = subtreeMinMax(raw, begin); err != nil {
			return
		}
		end = begin + end
		ranges[`minimum`] = []int{begin, end}
	}

	if begin := idxs(raw, `maximum `); begin != -1 {
		var end int
		begin += 8
		if ss.Max, end, err = subtreeMinMax(raw, begin); err != nil {
			return
		}
		end = begin + end
		ranges[`maximum`] = []int{begin, end}
	}

	if begin := idxs(raw, `specificationFilter `); begin != -1 {
		ss.SpecificationFilter = NewSpecificationFilter(raw[begin+20:])
		ranges[`specificationFilter`] = []int{begin, len(raw) - 1}
	}

	return
}

/*
NewSpecificationFilter returns an instance of [Refinement] based upon the
string input value.

Instances of [Refinement] are intended for assignment to the [SubtreeSpecification]
"SpecificationFilter" field, though this is done automatically when using the
[NewSubtreeSpecification] function.
*/
func NewSpecificationFilter(input string) (r Refinement) {
	if hasPfx(input, "and:") {
		r.And = parseRefinementList(trimPfx(input, "and:"))
	} else if hasPfx(input, "or:") {
		r.Or = parseRefinementList(trimPfx(input, "or:"))
	} else if hasPfx(input, "not:") {
		notRefinement := NewSpecificationFilter(trimPfx(input, "not:"))
		r.Not = &notRefinement
	} else if hasPfx(input, "item:") {
		r.Item = trimPfx(input, "item:")
	}

	return
}

func parseRefinementList(input string) (refs []Refinement) {

	t := trim(input, "{} ")

	var partStart, braceLevel int
	for i, r := range t {
		if r == ',' && braceLevel == 0 {
			refs = append(refs, NewSpecificationFilter(t[partStart:i]))
			partStart = i + 1
		} else if r == '{' {
			braceLevel++
		} else if r == '}' {
			braceLevel--
		}
	}

	if partStart < len(t) {
		refs = append(refs, NewSpecificationFilter(t[partStart:]))
	}

	return
}

/*
SpecificExclusions implements the Subtree Specification exclusions construct.

From Appendix A of RFC 3672:

	SpecificExclusions = "{" [ sp SpecificExclusion *( "," sp SpecificExclusion ) ] sp "}"
*/
type SpecificExclusions []SpecificExclusion

/*
SpecificExclusion implements the Subtree Specification exclusion construct.

From Appendix A of RFC 3672:

	SpecificExclusion  = chopBefore / chopAfter
	chopBefore         = id-chopBefore ":" LocalName
	chopAfter          = id-chopAfter  ":" LocalName
	id-chopBefore      = %x63.68.6F.70.42.65.66.6F.72.65 ; "chopBefore"
	id-chopAfter       = %x63.68.6F.70.41.66.74.65.72    ; "chopAfter"
*/
type SpecificExclusion struct {
	Name  LocalName
	After bool // false = Before
}

/*
BaseDistance implements the Minimum and Maximum base distance specifiers
for a [SubtreeSpecification].
*/
type BaseDistance int

/*
LocalName implements an RDNSequence.
*/
type LocalName string

/*
String returns the string representation of the receiver instance.
*/
func (r SpecificExclusions) String() string {
	if len(r) == 0 {
		return `{ }`
	}

	var _s []string
	for i := 0; i < len(r); i++ {
		_s = append(_s, r[i].String())
	}

	return `{ ` + join(_s, `, `) + ` }`
}

/*
String returns the string representation of the receiver instance.
*/
func (r SpecificExclusion) String() (s string) {
	if len(r.Name) > 0 {
		if r.After {
			s = `chopAfter ` + `"` + string(r.Name) + `"`
		} else {
			s = `chopBefore ` + `"` + string(r.Name) + `"`
		}
	}

	return
}

func subtreeExclusions(raw string, begin int) (excl SpecificExclusions, end int, err error) {
	end = -1

	if raw[begin] != '{' {
		err = errors.New("Bad exclusion encapsulation")
		return
	}

	var pos int
	if pos, end, err = deconstructExclusions(raw, begin); err != nil {
		return
	}

	values := fields(raw[pos:end])
	excl = make(SpecificExclusions, 0)

	for i := 0; i < len(values); i += 2 {
		var ex SpecificExclusion
		if !strInSlice(values[i], []string{`chopBefore`, `chopAfter`}) {
			err = errors.New("Unexpected key '" + values[i] + "'")
			break

		}
		ex.After = values[i] == `chopAfter`

		localName := trim(trimR(values[i+1], `,`), `"`)
		//if err = isSafeUTF8(localName); err == nil {
		ex.Name = LocalName(localName)
		excl = append(excl, ex)
		//}
	}

	return
}

func deconstructExclusions(raw string, begin int) (pos, end int, err error) {
	pos = -1
	if idx := idxs(raw[begin:], `chop`); idx != -1 {
		var (
			before int = -1
			after  int = -1
		)

		if hasPfx(raw[begin+idx+4:], `Before`) {
			before = begin + idx
		}

		if hasPfx(raw[begin+idx+4:], `After`) {
			after = begin + idx
		}

		if after == -1 && before > after {
			pos = before
		} else if before == -1 && before < after {
			pos = after
		}
	}

	if pos == -1 {
		err = errors.New("No chop directive found in value")
		return
	}

	for i, char := range raw[pos:] {
		switch char {
		case '}':
			end = pos + i
			break
		}
	}

	return
}

func subtreeMinMax(raw string, begin int) (minmax BaseDistance, end int, err error) {
	end = -1

	var (
		max string
		m   int
	)

	for i := 0; i < len(raw[begin:]); i++ {
		if isDigit(rune(raw[begin+i])) {
			max += string(raw[begin+i])
			continue
		}
		break
	}

	if m, err = atoi(max); err == nil {
		minmax = BaseDistance(m)
		end = len(max)
	}

	return
}

/*
String returns the string representation of the receiver instance.
*/
func (r SubtreeSpecification) String() (s string) {
	var _s []string
	if len(r.Base) > 0 {
		_s = append(_s, `base `+`"`+string(r.Base)+`"`)
	}

	if x := r.SpecificExclusions; len(x) > 0 {
		_s = append(_s, `specificExclusions `+x.String())
	}

	if r.Min > 0 {
		_s = append(_s, `minimum `+itoa(int(r.Min)))

	}

	if r.Max > 0 {
		_s = append(_s, `maximum `+itoa(int(r.Max)))

	}

	if x := r.SpecificationFilter.String(); len(x) > 0 {
		_s = append(_s, `specificationFilter `+x)
	}

	s = `{` + join(_s, `, `) + `}`

	return
}

/*
String returns the string representation of the receiver instance.
*/
func (r Refinement) String() string {
	sb := newBuilder()

	if r.Item != "" {
		sb.WriteString("item:" + r.Item)
	} else if len(r.And) > 0 {
		sb.WriteString("and:{")
		for i, andRef := range r.And {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(andRef.String())
		}
		sb.WriteString("}")
	} else if len(r.Or) > 0 {
		sb.WriteString("or:{")
		for i, orRef := range r.Or {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(orRef.String())
		}
		sb.WriteString("}")
	} else if r.Not != nil {
		sb.WriteString("not:" + r.Not.String())
	}

	return sb.String()
}

/*
Refinement implements the Subtree Specification Refinement construct.

From [Appendix A of RFC 3672]:

	Refinement  = item / and / or / not
	item        = id-item ":" OBJECT-IDENTIFIER
	and         = id-and  ":" Refinements
	or          = id-or   ":" Refinements
	not         = id-not  ":" Refinement

	Refinements = "{" [ sp Refinement *( "," sp Refinement ) ] sp "}"
	id-item     = %x69.74.65.6D ; "item"
	id-and      = %x61.6E.64    ; "and"
	id-or       = %x6F.72       ; "or"
	id-not      = %x6E.6F.74    ; "not"

From [clause 12.3.5 of ITU-T Rec. X.501]:

	Refinement ::= CHOICE {
		item [0] OBJECT-CLASS.&id,
		and  [1] SET SIZE (1..MAX) OF Refinement,
		or   [2] SET SIZE (1..MAX) OF Refinement,
		not  [3] Refinement,
		... }

[Appendix A of RFC 3672]: https://datatracker.ietf.org/doc/html/rfc3672#appendix-A
[clause 12.3.5 of ITU-T Rec. X.501]: https://www.itu.int/rec/T-REC-X.501
*/
type Refinement struct {
	Item string
	And  []Refinement
	Or   []Refinement
	Not  *Refinement
}

func subtreeBase(raw string) (base LocalName, end int, err error) {
	end = -1
	if len(raw) == 0 {
		err = errors.New("Empty LocalName")
		return
	}

	if raw[0] != '"' {
		err = errors.New("Missing encapsulation (\") for LocalName")
		return
	}

	for i := 1; i < len(raw) && end == -1; i++ {
		switch char := rune(raw[i]); char {
		case '"':
			end = i
			break
		}
	}

	//if err = isSafeUTF8(raw[1:end]); err == nil {
	base = LocalName(raw[1:end])
	//}

	return
}
