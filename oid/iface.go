package oid

/*
OID is qualified through any instance of [DotNotation], [ASN1Notation]
and [IRINotation].
*/
type OID interface {
	// IsZero returns a Boolean value indicative of a nil
	// or zero length receiver.
	IsZero() bool

	// Len returns the integer length of the receiver.
	Len() int

	// ChildOf returns a Boolean value indicative of the receiver
	// being a child of the input value.
	ChildOf(any) bool

	// AncestorOf returns a Boolean value indicative of the receiver
	// being an ancestor of the input value.
	AncestorOf(any) bool

	// SiblingOf returns a Boolean value indicative of the receiver
	// being a sibling of the input value.
	SiblingOf(any) bool

	// Valid returns a Boolean value indicative of the receiver
	// representing a legal value.
	Valid() bool

	// String returns the string representation of the receiver.
	String() string

	// Type returns one (1) of the string literals "oiv", "dot" or "iri".
	Type() string
}

/*
Arc is qualified through any instance of [NumberForm], [NameAndNumberForm]
and [IRIArc].
*/
type Arc interface {
	IRIArc
	NameAndNumberForm
	NumberForm
}
