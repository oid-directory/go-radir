/*
Package oid contains convenient helpers for [DotNotation], [IRINotation], [ASN1Notation], [NumberForm], [IRIArc] and [NameAndNumberForm] types, allowing for simple syntax checks and reliable comparison operations. This package subdirectory is mainly intended for use within the context of the OID Directory I-D series.

# NumberForms

The [DotNotation], [IRIArc] and [NameAndNumberForm] types contain instances of [NumberForm], which is an unbounded, unsigned ASN.1 INTEGER implementation. In this context, this means there is NO limit on upper magnitude.

By default, [NumberForm] uses a native uint64 for storage of the decimal value. If that decimal value overflows uint64, an instance of *[math/big.Int] is used for its support of arbitrary precision. Performance penalties may apply, given sufficiently large values.

In the world-wide OID tree, most number forms can fit within uint64, but not all of them will (hint: see the leaf nodes beneath "{itu-t(0) uuid(25)}" for real-world examples of this exception).

# Encoding

The [NumberForm] and [DotNotation] types support seamless payload encoding and decoding, however no ASN.1 encoding rules are implemented (e.g.: BER, DER, et al) in this package. This means no so-called TLV (Type-Length-Value) instances are produced by this package.

In other words, when encoding the [DotNotation] "1.3.6.1.4.1", the encoding will manifest as ...

	[]byte{0x2b, 0x06, 0x01, 0x04, 0x01}
	       ----------------------------
	                    V

... and not ...

	[]byte{0x06, 0x05, 0x2b, 0x06, 0x01, 0x04, 0x01}
	       ----- ----- ----------------------------
	         T     L                V

If users require specific encoding rule capabilities, they are expected to implement the needed specification themselves.

No other types in this package are eligible for encoding in this manner.

# Type Constructors and Criticality

Each constructor function or method -- for example, [NewDotNotation] -- has a critical counterpart of the same name, but with the term "Must" added at the beginning (e.g.: [MustNewDotNotation]).

The difference is that the critical constructor only returns the indicated value, and panics on any error versus returning that error as an instance.

This is useful for direct assignment of freshly constructed values to struct fields, et al, without needing to follow Go's standard idiomatic practice with regards to adhoc error handling.

For instance:

	// Non-critical behavior
	id, err := NewDotNotation(...)
	<check your error>

	// Critical behavior
	id := MustNewDotNotation(...)
	<panics if error was encountered>
*/
package oid
