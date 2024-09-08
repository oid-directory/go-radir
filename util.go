package radir

import (
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

var (
	sprintf  func(string, ...any) string                  = fmt.Sprintf
	fprintf  func(io.Writer, string, ...any) (int, error) = fmt.Fprintf
	fprintln func(io.Writer, string, ...any) (int, error) = fmt.Fprintf
	printf   func(string, ...any) (int, error)            = fmt.Printf

	atoi func(string) (int, error) = strconv.Atoi
	itoa func(int) string          = strconv.Itoa

	fields     func(string) []string               = strings.Fields
	hasPfx     func(string, string) bool           = strings.HasPrefix
	hasSfx     func(string, string) bool           = strings.HasSuffix
	idxRune    func(string, rune) int              = strings.IndexRune
	join       func([]string, string) string       = strings.Join
	lc         func(string) string                 = strings.ToLower
	uc         func(string) string                 = strings.ToUpper
	split      func(string, string) []string       = strings.Split
	eq         func(string, string) bool           = strings.EqualFold
	contains   func(string, string) bool           = strings.Contains
	splitAfter func(string, string) []string       = strings.SplitAfter
	splitN     func(string, string, int) []string  = strings.SplitN
	trimS      func(string) string                 = strings.TrimSpace
	trimL      func(string, string) string         = strings.TrimLeft
	trimR      func(string, string) string         = strings.TrimRight
	replaceAll func(string, string, string) string = strings.ReplaceAll

	isLetter func(rune) bool = unicode.IsLetter
	isDigit  func(rune) bool = unicode.IsDigit
	isLower  func(rune) bool = unicode.IsLower
	isUpper  func(rune) bool = unicode.IsUpper

	typeOf func(any) reflect.Type  = reflect.TypeOf
	valOf  func(any) reflect.Value = reflect.ValueOf
)

func newBuilder() strings.Builder {
	return strings.Builder{}
}

func isNumber(n string) bool {
	if len(n) == 0 {
		return false
	}

	for i := 0; i < len(n); i++ {
		if !isDigit(rune(n[i])) {
			return false
		}
	}
	return true
}

/*
isPtr returns a Boolean value indicative of whether kind
reflection revealed the presence of a pointer type.
*/
func isPtr(x any) bool {
	if x == nil {
		return false
	}

	return typeOf(x).Kind() == reflect.Ptr
}

func getReflectInstances(a any) (t reflect.Type, v reflect.Value, ok bool) {
	ot := typeOf(a) // reflect.Type
	ov := valOf(a)  // reflect.Value

	// If this is a pointer instance,
	// we'll need to dereference it.
	if isPtr(a) {
		ot = ot.Elem()
		ov = ov.Elem()
	}

	// Whether pointer or not, this
	// function only handles structs.
	if ov.Kind() != reflect.Struct {
		return
	}

	t = ot
	v = ov
	ok = true

	return
}

/*
structEmpty focuses on content as opposed to allocation. Thus, this function
is separate and distinct from the IsZero methods exported by most types in
this package.

Please note this is a self-executing function.
*/
func structEmpty(x any) (is bool) {
	is = true
	ot, ov, ok := getReflectInstances(x)
	if !ok {
		return
	}

	for i := 0; i < ot.NumField() && is; i++ {
		t := ot.Field(i)

		if !t.IsExported() {
			continue
		} else if t.Name == "R_DITProfile" {
			continue
		}

		v := ov.Field(i)

		// Dereference any ptr elements
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.String:
			is = len(v.String()) == 0
		case reflect.Slice:
			iface := v.Interface()
			switch tv := iface.(type) {
			case []string:
				is = len(tv) == 0
				//case []byte:
				//is = len(tv) == 0
			}
		case reflect.Struct:
			is = structEmpty(v.Interface())
		}
	}

	return
}

/*
unmarshalStruct is a general use struct-to-map unmarshaler for use by
this package to unmarshal (or "output") the contents of x into a new
instance of map[string][]string. The final return object is meant for
submission to go-ldap/v3.NewEntry.

Note that only EXPORTED fields are analyzed, and any field that houses
COLLECTIVE values will be skipped silently (collective attributes are
not added to any DIT in this manner).

Also note this is a self-executing function.
*/
func unmarshalStruct(x any, outer map[string][]string) map[string][]string {
	ot, ov, ok := getReflectInstances(x)
	if !ok {
		return outer
	}

	for i := 0; i < ot.NumField(); i++ {
		t := ot.Field(i)

		v := ov.Field(i)
		xt := t.Tag.Get(`ldap`)

		if unmarshalSkipField(lc(xt), t) {
			continue
		}

		// Dereference any ptr elements
		v = derefPtr(v)

		switch v.Kind() {
		case reflect.String:
			if len(xt) > 0 {
				if val := v.String(); len(val) > 0 {
					outer[xt] = []string{v.String()}
				}
			}
		case reflect.Slice:
			if len(xt) > 0 {
				if val, ok := v.Interface().([]string); ok && len(val) > 0 {
					outer[xt] = val
				}
			}
		case reflect.Struct:
			inner := make(map[string][]string)
			for k, v := range unmarshalStruct(v.Interface(), inner) {
				outer[k] = v
			}
		}
	}

	return outer
}

func derefPtr(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v
}

func unmarshalSkipField(tag string, t reflect.StructField) bool {
	return !t.IsExported() || hasPfx(tag, `c-`) || hasSfx(tag, `;collective`)
}

func selectTTL(lttl, cttl string) string {
	var ttl string
	if len(cttl) > 0 {
		ttl = cttl
	}

	if len(lttl) > 0 {
		ttl = lttl
	}

	if len(ttl) == 0 {
		return ``
	}

	return ttl
}

func assertTTL(ttl any) (t int) {
	switch tv := ttl.(type) {
	case string:
		t, _ = atoi(tv)
	case int:
		t = tv
	}

	return
}

func assertGetOrSetFunc(args ...any) (gosf GetOrSetFunc, err error) {
	// note: error only thrown if args > 1 and
	// 2nd arg is not a GetOrSetFunc type instance
	if len(args) > 1 {
		var ok bool
		if gosf, ok = args[1].(GetOrSetFunc); !ok {
			// try a fallback signature of the
			// same form, should interfacing
			// get in the way ...
			var _gosf func(...any) (any, error)
			if _gosf, ok = args[1].(func(...any) (any, error)); !ok {
				err = NilGetOrSetFuncErr
			} else {
				gosf = GetOrSetFunc(_gosf)
			}
		}
	}

	return
}

func chkSetterInput(args ...any) (X any, Y GetOrSetFunc, err error) {
	if len(args) == 0 {
		err = NilArgumentsErr
		return
	}

	if Y, err = assertGetOrSetFunc(args...); err == nil {
		X = args[0]
	}

	return
}

/*
strInSlice returns a Boolean value indicative of whether the specified
string (str) was present within the slice (sl).

Case is not significant in the matching process.
*/
func strInSlice(str string, sl []string) bool {
	for i := 0; i < len(sl); i++ {
		if eq(str, sl[i]) {
			return true
		}
	}

	return false
}

func removeStrInSlice(str string, sl []string) []string {
	if !strInSlice(str, sl) {
		return sl
	}

	var n []string
	for i := 0; i < len(sl); i++ {
		if eq(str, sl[i]) {
			continue
		}
		n = append(n, sl[i])
	}

	return n
}

func getFieldValueByNameTagAndGoSF(instance any, gosf GetOrSetFunc, tag string) (fv any, err error) {
	if gosf == nil {
		err = NilGetOrSetFuncErr
	} else if instance == nil {
		err = NilArgumentsErr
	} else {
		var _fv reflect.Value
		if _fv, err = getFieldByNameTag(instance, tag); err == nil {
			fv, err = gosf(_fv.Interface(), instance)
		}
	}

	return
}

func getAttributeTypeFieldTags(instance any) (tags []string) {
	var t reflect.Type = typeOf(instance)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tags = append(tags, field.Tag.Get("ldap"))
	}

	return
}

func readFieldByTag(tag string, instance any) (values []string) {
	fieldValue, err := getFieldByNameTag(instance, tag)
	if err != nil {
		return
	}

	if fieldValue.Kind() == reflect.Slice {
		x := fieldValue.Interface()
		if _values, _ := x.([]string); len(_values) > 0 {
			values = append(values, _values...)
		}
	} else if fieldValue.Kind() == reflect.String {
		if str := fieldValue.String(); len(str) > 0 {
			values = append(values, str)
		}
	}

	return
}

func getFieldByNameTag(instance any, tag string) (reflect.Value, error) {
	v := valOf(instance)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, errorf("not a struct pointer")
	}

	t := typeOf(instance).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		itag := field.Tag.Get("ldap")
		if eq(itag, tag) || (hasPfx(itag, tag) && hasSfx(lc(itag), `;collective`)) {
			return v.Elem().Field(i), nil
		}
	}

	return reflect.Value{}, errorf("field with tag %s not found", tag)
}

func writeFieldByTag(tag string, funk func(...any) error, instance any, args ...any) error {
	X, Y, err := chkSetterInput(args...)
	if err != nil {
		return err
	}

	if Y == nil {
		return writeValue(instance, X, tag)
	}

	var v any
	if v, err = Y(X, instance); err == nil {
		err = funk(v)
	}

	return err
}

func writeValue(instance, value any, tag string) error {
	fieldValue, err := getFieldByNameTag(instance, tag)
	if err != nil {
		return err
	}

	if !fieldValue.CanSet() {
		return errorf("field with tag %s is not settable", tag)
	}

	switch v := value.(type) {
	case int:
		if err = writeEligible(tag, v, instance); err == nil {
			err = writeInt(tag, v, fieldValue)
		}
	case *big.Int:
		if err = writeEligible(tag, v, instance); err == nil {
			err = writeBigInt(tag, v, fieldValue)
		}
	case string:
		if err = writeEligible(tag, v, instance); err == nil {
			if err = writeString(tag, v, fieldValue); err == nil {
				specialHandling(tag, v, instance)
			}
		}
	case []string:
		if err = writeEligible(tag, v, instance); err == nil {
			if err = writeStrings(tag, v, fieldValue); err == nil {
				specialHandling(tag, v, instance)
			}
		}
	default:
		err = UnsupportedInputTypeErr
	}

	return err
}

func writeInt(tag string, v int, fieldValue reflect.Value) (err error) {
	if fieldValue.Kind() != reflect.String {
		err = errorf("field with tag %s is not a string", tag)
	} else if !strInSlice(tag, []string{`n`, `registrationRange`, `rATTL`}) {
		err = errorf("field with tag %s is not a integer compatible", tag)
	} else {
		fieldValue.SetString(itoa(v))
	}

	return
}

func writeString(tag string, v string, fieldValue reflect.Value) (err error) {
	if fieldValue.Kind() == reflect.Slice {
		if !fieldValue.IsZero() {
			if fvi, ok := fieldValue.Interface().([]string); ok {
				if !strInSlice(v, fvi) {
					fieldValue.Set(reflect.Append(fieldValue, valOf(v)))
				}
			}
		} else {
			fieldValue.Set(reflect.Append(fieldValue, valOf(v)))
		}
	} else if fieldValue.Kind() == reflect.String {
		fieldValue.SetString(v)
	} else {
		err = errorf("field with tag %s is not a string", tag)
	}

	return
}

func writeStrings(tag string, v []string, fieldValue reflect.Value) error {
	if fieldValue.Kind() != reflect.Slice || fieldValue.Type().Elem().Kind() != reflect.String {
		return errorf("field with tag %s is not a []string", tag)
	}

	fieldValue.Set(valOf(v))

	return nil
}

func writeBigInt(tag string, v *big.Int, fieldValue reflect.Value) (err error) {
	if fieldValue.Kind() != reflect.String {
		err = errorf("field with tag %s is not a string", tag)
	} else if eq(tag, `n`) {
		if v.Cmp(big.NewInt(0)) > -1 {
			fieldValue.SetString(v.String())
			return
		}
		err = errorf("number assignment with tag %s cannot be negative", tag)
	}

	return
}

/*
Perform special pre-write checks to ensure sanity and compliance.
*/
func writeEligible(tag string, value, instance any) (err error) {
        switch tv := instance.(type) {
        case *X660:
                err = tv.writeEligible(tag, value)
        }

	return
}

func condenseWHSP(b string) string {
	b = trimS(b)

	var last bool // previous char was WHSP or HTAB.
	var bld strings.Builder

	for i := 0; i < len(b); i++ {
		c := rune(b[i])
		switch c {
		case rune(9), rune(32): // match either WHSP or horizontal tab
			if !last {
				last = true
				bld.WriteRune(rune(32)) // Add WHSP
			}
		default: // match all other characters
			if last {
				last = false
			}
			bld.WriteRune(c)
		}
	}

	return bld.String()
}

/*
isIdentifier scans the input string val and judges whether it
qualifies as an X.680 identifier, or name form. All of the
following MUST evaluate as true:

  - Non-zero in length
  - Begins with a lowercase alphabetical character
  - Ends in an alphanumeric character
  - Contains only alphanumeric characters or hyphens
  - No contiguous hyphens

This function is an alternative to engaging the [antlr4512]
parsing subsystem.
*/
func isIdentifier(val string) bool {
	if len(val) == 0 {
		return false
	}

	// must begin with a lower alpha.
	if !isLower(rune(val[0])) {
		return false
	}

	// can only end in alnum.
	if !isAlnum(rune(val[len(val)-1])) {
		return false
	}

	// watch hyphens to avoid contiguous use
	var lastHyphen bool

	// iterate all characters in val, checking
	// each one for validity.
	for i := 0; i < len(val); i++ {
		ch := rune(val[i])
		switch {
		case isAlnum(ch):
			lastHyphen = false
		case ch == '-':
			if lastHyphen {
				// cannot use consecutive hyphens
				return false
			}
			lastHyphen = true
		default:
			// invalid character (none of [a-zA-Z0-9\-])
			return false
		}
	}

	return true
}

/*
isAlnum returns a Boolean value indicative of whether rune r represents
an alphanumeric character. Specifically, one (1) of the following ranges
must evaluate as true:

  - 0-9 (ASCII characters 48 through 57)
  - A-Z (ASCII characters 65 through 90)
  - a-z (ASCII characters 97 through 122)
*/
func isAlnum(r rune) bool {
	return isLower(r) || isUpper(r) || isDigit(r)
}

/*
take additional steps for attribute types of note, such as measuring the
length (depth) of a freshly-set 'dotNotation' value upon an instance of
*X680.
*/
func specialHandling(tag string, value, instance any) {
	switch tv := instance.(type) {
	case *X680:
		tv.specialHandling(tag, value)
	}
}

func dotJoin(sp []string) (dot string) {
	if len(sp) > 0 {
		dot = join(sp,`.`)
	}

	return
}

func dotSplit(dot string) (sp []string) {
	if len(dot) > 0 {
		sp = split(dot, `.`)
	}

	return
}

func getRoot(in rune) (out int) {
	out = -1
	// inc roots by 1 to avoid
	// false positives.
	switch in {
	case '0':
		out = 0
	case '1':
		out = 1
	case '2':
		out = 2
	}

	return
}

func rootClass(n int) (class string) {
	// note: case values adjusted by +1.
	switch n {
	case 0:
		class = `iTUTRegistration` // RASCHEMA § 2.5.8
	case 1:
		class = `iSORegistration` // RASCHEMA § 2.5.9
	case 2:
		class = `jointISOITUTRegistration` // RASCHEMA § 2.5.10
	}

	return
}
