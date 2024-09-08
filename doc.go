/*
Package radir implements a large subset of "The OID Directory" -- an EXPERIMENTAL Internet-Draft (I-D) series by Jesse Coretta.

# Draft Information

The Internet-Drafts (henceforth referred to as "the I-D series") is made up of the following individual drafts:

  - draft-coretta-oiddir-roadmap
  - draft-coretta-oiddir-radit
  - draft-coretta-oiddir-radsa
  - draft-coretta-oiddir-radua
  - draft-coretta-oiddir-schema

These drafts can be viewed on the [IETF Data Tracker site], or via the official [OID Directory site] and [GitHub repositories]. At present, the current revisions are set to expire on February 23, 2025.

# Experimental Status

The I-D series, and by necessity this package, is thoroughly EXPERIMENTAL. It is not yet approved by the IETF, and thus should NEVER be used in any capacity beyond proof-of-concept or work-in-progress efforts.

# What this package is

This package is an abstract, general-use framework supplement. It will aid in the marshaling and unmarshaling of OID registration and registrant (contact) constructs, whether using a proper [go-ldap/v3 Entry] instance, or through manual assembly.

The package can aid in the bidirectional conversion of certain values, such as "dotNotation" and "dn" values, and offers many other useful features in service to the I-D series mentioned above.

Implementations which use this package may be of a server-side or client-side nature, or neither. There is no singular use-case for this package.

TLDR; its a nifty toolbox; what you build is what the package serves.

# What this package is not

As the terms are defined throughout the OID Directory I-D series, this package is absolutely not a complete DUA, DIT or DSA. While it can serve as a valuable component in such constructs, its current state does not allow drop-in functionality of that nature, nor was this intended.

For instance, those designing a compliant RA DUA, per the [RADUA I-D], are expected to install and utilize the [go-ldap/v3] package on their own terms and in service to their particular environment or infrastructure.

This is done to maximize compatibility across the many potential use-cases and directory products, as well as to limit potential security vulnerabilities relating to this package itself. This approach also has the secondary effect of making potential integration efforts much simpler and far less disruptive.

TLDR; this package works with [go-ldap/v3], but it does NOT import it directly. Do it yourself.

# GetOrSetFunc Extensibility

Thanks to the [GetOrSetFunc] closure type, this package is supremely extensible.

Virtually all [Registration] and [Registrant] methods -- such as [Registration.SetDN] or [FirstAuthority.POBoxGetFunc] -- allow for closure-based behavioral overrides. This allows limitless control over how values manifest during presentation, as well as how they are written to instances of the aforementioned types.

For additional information, see the [GetOrSetFunc] type documentation, as well as the package examples for all methods which allow input of instances of this type. See also the next section regarding storage space considerations with regards to especially -- and unnecessarily -- large values.

TLDR; Control value I/O using a closure signature of "func(...any) (any, error)" ([GetOrSetFunc]) for any "Set<X>" or "<X>GetFunc" methods.

# Instance caching

Per [Section 2.2.3.4 of the RADUA I-D], this package provides a thread-safe, memory-based [Cache] facility for use by a client.

The primary purpose of this facility is to cache, or store temporarily, all *[Registration] and *[Registrant] instances that have either been crafted manually, or marshaled by way of a [go-ldap/v3] entry instance. While crude, it can help provide considerable I/O savings in terms of LDAP search requests, which may or may not be transmitted over-the-wire.

Lifespans of cached entries is directed by manual specification (e.g.: by the end user), or by way of a literal or collectively-inherited TTL obtained within the RA DIT or via the appropriate *[DITProfile] instance as a global fallback. See the aforementioned RA DUA section for details regarding TTL precedence and other mechanics.

Use of this facility is not required to comply with the aforementioned specification. Adopters may freely supplant the package-provided [Cache] with a caching system of their own choosing or design.

TLDR; Caching eligible instances reduces network (LDAP) I/O at the expense of memory. You can use the [Cache] type, or a third-party one, or abstain from caching entirely.

# Directory model

The I-D series offers two (2) directory models in terms of [Registration] structure and layout, each of which are implemented in this package.

  - Two dimensional model
  - Three dimensional model

The two dimensional model is discussed in [Section 3.1.2 of the RADIT I-D]. The three dimensional model is discussed in [Section 3.1.3 of the RADIT I-D]. In most scenarios, use of the three dimensional model is the preferred strategy.

TLDR; Use the [ThreeDimensional] directory model.

# Registrant entry policy

The I-D series offers two (2) registrant entry policies, each of which are implemented in this package.

  - Dedicated: Authorities have their own entries, and [Registration] instances link to these authorities using DN-based attributes
  - Combined: [Registration] entries incorporate authority information in-line within a single entry

Dedicated registrants are covered in [Section 3.2.1.1.1 of the RADIT I-D]. Combined registrants are briefly covered in [Section 3.2.1.1.2 of the RADIT I-D]. In most scenarios, use of dedicated registrants is the preferred strategy.

TLDR; Use *[Registrant] instances instead of "combining" registrant content with *[Registration] instances (in-line).

# Registrant attribute type policy

As stated in [Section 3.2.1.1.1 of the RADIT I-D], it is possible to forego use of the draft-based authority types, such as "[sponsorCommonName]", in favor of the traditional "[cn]" type. This logic applies may extend to either "Combined" or "Dedicated" Registrant Policies.

See the [DITProfile.UseAltAuthorityTypes] method for a means of enabling this behavior. Note there are caveats with either standpoint, and thus the reader is advised to review the aforementioned section of the draft to ensure they understand the ramifications of their decision.

Please also note it is inadvisable to change this value without a good reason, and inappropriate alteration will result in degraded client behavior and likely a deviation from the established content policies enforced within the directory. You have been warned.

See the [FirstAuthorityAltAttributeTypes], [CurrentAuthorityAltAttributeTypes] and [SponsorAltAttributeTypes] map variables for a complete list of the types that are -- and are not -- subject to the influence of the aforementioned method.

TLDR; You may use, for example, "[cn]" instead of "[sponsorCommonName]" ... but there are caveats.

# (Un)marshaling support

This package makes conversion (in either direction) between [go-ldap/v3 Entry] and *[Registration] or *[Registrant] instances a breeze!

When unmarshaling FROM an instance of [go-ldap/v3 Entry] TO an instance of *[Registration], rather than using the [go-ldap/v3 Entry.Unmarshal] method directly, simply feed the method to *[Registration.Marshal] to achieve the same effect:

	var entry *ldap.Entry // assume this was populated by an LDAP Search already
	var reg *Registration // assume this was initialized already

	// Note this is a closure scenario (no "()")
	if err := reg.Marshal(entry.Unmarshal); err != nil {
	     fmt.Println(err)
	     return
	}

This is necessary because the [go-ldap/v3 Entry.Unmarshal] method only supports a limited number of struct field value types. To get around this issue, radir simply performs independent marshaling upon any individual struct components within the destination instance (*[Registration]). In other words, if there are four fields that contain struct values, each of these fields is marshaled independently. This ensures that all of the needed attribute values are collected from the source [go-ldap/v3 Entry] instance.

When unmarshaling FROM an instance of *[Registration] (or *[Registrant]) TO an instance of [go-ldap/v3 Entry], simply use the [Registration.Unmarshal] (or [Registrant.Unmarshal]) method. Feed the output to the [go-ldap/v3 NewEntry] function:

	var reg *Registration // assume this was already populated somehow

	// Note, contrary to the above, this is NOT closure.
	// Unmarshal simply outputs map[string][]string.
	entry := ldap.NewEntry(reg.Unmarshal())

TLDR; Excellent marshal and unmarshal features. And while [go-ldap/v3 Entry.Unmarshal] is very limited, we have a most suitable workaround: don't "use" it, let us handle it for you.

# Spatial support

OIDs are virtually infinite in size. Large pools of sibling registrations can be exceedingly difficult to navigate manually; the sequence of number forms may not be contiguous, and there is no guarantee the entries which bear these values will be ordered correctly in directory search results.

To that end, the "[spatialContext]" AUXILIARY class defined within the I-D series is implemented within this package as the *[Spatial] struct type.

Use of this type can help mitigate some of this tedium by allowing any given registration entry to bear direct DN-based references to other spatially-relevant registrations.

Specifically, this produces an abstraction of directional movement in the following contexts:

  - "[supArc]" / "[c-supArc]" - Up (immediate ancestor (parent) registration)
  - "[topArc]" / "[c-topArc]" - Top (root ancestor)
  - "[minArc]" / "[c-minArc]" - Far Left (absolute lowest numbered sibling)
  - "[maxArc]" / "[c-maxArc]" - Far Right (absolute highest numbered sibling)
  - "[leftArc]" - Left (nearest sibling numbered less than the target)
  - "[rightArc]" - Right (nearest sibling numbered greater than the target)
  - "[subArc]" - Down (child registration(s))

Non-collective *[Spatial] attribute types may be set manually, or they may be present within entries marshaled into [Registration] instances as literal or collective values. Collective values are not meant for manual assignment, thus no related "set" methods exist in that regard.

Like virtually all other methods in this package, the relevant *[Spatial] methods allow for [GetOrSetFunc] closure use, thereby letting the user enhance the behavior of instances of this type in a variety of ways:

  - Present a bulky spatial DN as a single (leaf) number form, such as "5"
  - Convert a (string) number form to a *[big.Int] for [ITU-T Rec. X.667] activities
  - Re-order a collection of horizontal (sibling) *[Registration] instances artificially according to number form magnitude using [sort] or some other sorting package
  - Extrapolate ancestral spatial DNs in cases where a directory lacks such literal values

TLDR; RA DIT navigation with a "🕹️" duct-taped on to it.

[OID Directory site]: https://oid.directory
[ITU-T Rec. X.667]: https://www.itu.int/rec/T-REC-X.667
[IETF Data Tracker site]: https://datatracker.ietf.org
[GitHub repositories]: https://github.com/oid-directory
[RADUA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radua
[go-ldap/v3]: https://pkg.go.dev/github.com/go-ldap/ldap/v3
[go-ldap/v3 Entry]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry
[go-ldap/v3 NewEntry]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#NewEntry
[go-ldap/v3 Entry.Unmarshal]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry.Unmarshal
[cn]: https://datatracker.ietf.org/doc/html/rfc4519#section-2.3
[sponsorCommonName]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.78
[spatialContext]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.5.11
[supArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.21
[subArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.25
[c-supArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.22
[topArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.23
[c-topArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.24
[leftArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.26
[rightArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.29
[minArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.27
[c-minArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.28
[maxArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.30
[c-maxArc]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2.3.31
[Section 3.1.2 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.1.2
[Section 3.1.3 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.1.3
[Section 3.2.1.1.2 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.2.1.1.2
[Section 3.2.1.1.1 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.2.1.1.1
[Section 2.2.3.4 of the RADUA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radua#section-2.2.3.4

[RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit
[go-ldap/v3 Entry]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry
[Entry.Unmarshal]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#Entry.Unmarshal
[Section 3.2.4.19 of the RADIT I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit#section-3.2.4.19
[Section 2 of the RASCHEMA I-D]: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema#section-2
*/
package radir
