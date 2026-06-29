package radir

import (
	"fmt"
	"testing"
)

func ExampleTTLPrecedence() {
	ttl := TTLPrecedence(
		0,     // from config (unset), lowest precedence
		"120", // from collective attribute assignment
		720,   // from explicit entry assignment, highest precedence
		"60",  // OPTIONAL: fallback if all of the above are 0
	)
	fmt.Printf("TTL precedence: %d", ttl)
	// Output: TTL precedence: 720
}

func TestTTLPrecedence(t *testing.T) {
	for try, sl := range []struct{
		args []any
		want int
	}{
		// ints
		{args: []any{60, 0, 0, 0}, want: 60},
		{args: []any{0, 120, 0, 5}, want: 120},
		{args: []any{0, 0, 0, 300}, want: 300},
		{args: []any{60, 0, 5, 0}, want: 5},
		{args: []any{0, 0, 0, 0}, want: 0},

		// strings
		{args: []any{"60", "0", "0", "0"}, want: 60},
		{args: []any{"0", "120", "0", "5"}, want: 120},
		{args: []any{"0", "0", "0", "300"}, want: 300},
		{args: []any{"60", "0", "5", "0"}, want: 5},
		{args: []any{"0", "0", "0", "0"}, want: 0},
	}{
		if got := TTLPrecedence(sl.args[0],sl.args[1],sl.args[2],sl.args[3]); got != sl.want {
			t.Errorf("%s [index %d] failed:\n\twant: %d min.\n\tgot:  %d min.",
				t.Name(), try, sl.want, got)
		}
	}
}
