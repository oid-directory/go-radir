package radir

/*
Entry is qualified through any instance of [Registration],
[Registant], [Subentry] or [Map].
*/
type Entry interface {
	DN() string
	TTL() string
	CTTL() string
	Profile() *DITProfile
}
