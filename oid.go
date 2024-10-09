package radir

import (
	"github.com/oid-directory/go-objectid"
)

var (
	NewOID               func(any) (*objectid.OID, error)               = objectid.NewOID
	NewASN1Notation      func(any) (*objectid.ASN1Notation, error)      = objectid.NewASN1Notation
	NewNameAndNumberForm func(any) (*objectid.NameAndNumberForm, error) = objectid.NewNameAndNumberForm
	NewNumberForm        func(any) (objectid.NumberForm, error)         = objectid.NewNumberForm
	NewDotNotation       func(...any) (*objectid.DotNotation, error)    = objectid.NewDotNotation
	IsIdentifier         func(string) bool                              = objectid.IsIdentifier
)

func IsNumericOID(id string) bool {
	_, err := NewDotNotation(id)
	return err == nil
}

func IsASN1Notation(asn string) bool {
	_, err := NewASN1Notation(asn)
	return err == nil
}

func IsNumberForm(nf string) bool {
	_, err := NewNumberForm(nf)
	return err == nil
}

func IsNameAndNumberForm(nanf string) bool {
	_, err := NewNameAndNumberForm(nanf)
	return err == nil
}
