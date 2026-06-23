package radir

import (
	"errors"
)

/*
err.go contains predefined error instances that
describe certain known aberrant conditions.
*/

var (
	RegistrationValidityErr,
	UnsupportedInputTypeErr,
	EndTimeNotApplicableErr,
	IllegalASN1NotationErr,
	RegistrantValidityErr,
	DUAConfigValidityErr,
	IllegalNumberFormErr,
	InvalidDimensionErr,
	RegistrantPolicyErr,
	NilRegistrationErr,
	NilGetOrSetFuncErr,
	IllegalLongArcErr,
	MismatchedLeafErr,
	NilRegistrantErr,
	InvalidGTFracErr,
	NilArgumentsErr,
	FrozenCacheErr,
	NilInstanceErr,
	IllegalRootErr,
	InvalidOIDErr,
	InvalidGTErr,
	NilMethodErr,
	InvalidDNErr,
	NilCacheErr,
	LongArcErr error
)

var mkerr func(string) error = errors.New

func init() {
	RegistrationValidityErr = errors.New("Registration instance did not pass validity checks")
	UnsupportedInputTypeErr = errors.New("Unsupported value type provided without GetOrSetFunc instance")
	EndTimeNotApplicableErr = errors.New("EndTime is not applicable to a CurrentAuthority")
	IllegalASN1NotationErr = errors.New("ASN.1 Notation value is malformed or zero-length")
	RegistrantValidityErr = errors.New("Registrant instance did not pass validity checks")
	DUAConfigValidityErr = errors.New("DUAConfig instance did not pass validity checks, or is poorly formed")
	IllegalNumberFormErr = errors.New("N (Number Form) is malformed or zero length")
	InvalidDimensionErr = errors.New("Unknown dimension; must be TwoDimensional or ThreeDimensional")
	RegistrantPolicyErr = errors.New("Registrant Policy violation")
	NilRegistrationErr = errors.New("Registration instance is nil; initialization required")
	NilGetOrSetFuncErr = errors.New("GetOrSetFunc instance is nil")
	IllegalLongArcErr = errors.New("LongArc cannot be applied to this registration type or root")
	MismatchedLeafErr = errors.New("Mismatched NumberForm with leaf node of ASN.1 and/or DotNotation")
	NilRegistrantErr = errors.New("Registrant instance is nil")
	NilArgumentsErr = errors.New("Missing input arguments")
	FrozenCacheErr = errors.New("Cache is frozen")
	NilInstanceErr = errors.New("Instance is nil")
	IllegalRootErr = errors.New("Illegal root; must be 'name' or 'name(0|1|2)' or 0|1|2")
	InvalidOIDErr = errors.New("OID value is malformed or zero length")
	NilMethodErr = errors.New("Input method signature is nil")
	InvalidDNErr = errors.New("DN value is malformed, zero length or has an unknown suffix")
	InvalidGTErr = errors.New("Invalid generalized time value")
	NilCacheErr = errors.New("Cache subsystem not initialized")
	LongArcErr = errors.New("longArc values can only be assigned to sub arcs of Joint-ISO-ITU-T")
	InvalidGTFracErr = errors.New(InvalidGTErr.Error() +
		": Fraction exceeds Generalized Time fractional limit")
}

func errorf(msg any, x ...any) error {
	switch tv := msg.(type) {
	case string:
		if len(tv) > 0 {
			return errors.New(sprintf(tv, x...))
		}
	case error:
		if tv != nil {
			return errors.New(sprintf(tv.Error(), x...))
		}
	}

	return nil
}
