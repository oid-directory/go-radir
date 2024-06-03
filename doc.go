/*
Package radir implements a basic entry framework relating to the OID
Directory -- an experimental ID series of which I am the author.

This framework is suitable for use in the development of applications or
operations related to the ID series, whether client or server focused in
nature.

# No LDAP Functionality

This package offers instances which are compatible with go-ldap's
[LDAPv3 Entry.Unmarshal] method by design. However, the go-ldap
package is not imported directly, rather this is left to the user
if and when needed.

# Basic Usage

	   import (
		"github.com/oid-directory/go-radir"
	   )

	   func main(){
		var reg Arc = Arc[ISO]{
		}
	   }

# Experimental Advisory

The following text is prominently displayed on the relevant IETF DataTracker
pages for this ID series:

	This document is an Internet-Draft (I-D). Anyone may submit an I-D to the
	IETF. This I-D is not endorsed by the IETF and has no formal standing in
	the IETF standards process.

It is posted here to reiterate that NO part of the ID series, however well
received by the community, is officially sanctioned or recommended for use
outside of proof-of-concept or testing purposes.  Use at your own risk.

# Abstraction Notice

The OID Directory ID series is written to allow a reasonable degree of
flexibility in terms of the manner in which an implementation occurs.
In short, there is more than one way to accomplish a variety of tasks.

As such, the reader is advised to consider that this framework, even if
written by the respective author, is but one interpretation of this
specification.  Other interpretations that may develop in the wild are
not necessarily invalid nor contradictory to the series.

# Resources

See the following resources for the full text of this ID series:

  - [The OID Directory: A Technical Roadmap] (RADIR)
  - [The OID Directory: The RA DIT] (RADIT)
  - [The OID Directory: The RA DSA] (RADSA)
  - [The OID Directory: The RA DUA] (RADUA)
  - [The OID Directory: The Schema] (RASCHEMA)

See also [The OID Directory Page].

[The OID Directory: A Technical Roadmap]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-roadmap
[The OID Directory: The RA DIT]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit
[The OID Directory: The RA DSA]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radsa
[The OID Directory: The RA DUA]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radua
[The OID Directory: The Schema]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema]
[RASCHEMA]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema]
[The OID Directory Page]: https://oid.directory
[LDAPv3 Entry.Unmarshal]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry.Unmarshal
*/
package radir
