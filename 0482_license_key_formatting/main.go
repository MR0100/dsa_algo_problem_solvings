package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Clean, Then Slice From the Left (Brute Force) ─────────────────
//
// bruteForceLeftSlice removes dashes, upper-cases, computes the (possibly
// short) first-group size, then walks left-to-right cutting fixed k-length
// groups.
//
// Intuition:
//
//	Strip everything down to the raw alphanumeric characters (uppercased).
//	Grouping is defined from the RIGHT — every group is exactly k except the
//	first, which is the remainder. If there are L clean characters, the first
//	group has size firstLen = L % k (or k when L % k == 0 and L > 0). Once we
//	know firstLen, we can emit that leading chunk and then every subsequent
//	k-sized chunk in normal left-to-right order, joining with dashes.
//
// Algorithm:
//  1. Build `clean` = all non-dash chars of s, upper-cased.
//  2. If clean is empty, return "".
//  3. firstLen = len(clean) % k; if firstLen == 0, firstLen = k.
//  4. Emit clean[0:firstLen], then clean[firstLen:firstLen+k], … each as a group.
//  5. Join the groups with "-".
//
// Time:  O(n) — one pass to clean, one pass to slice (n = len(s)).
// Space: O(n) — the cleaned string plus the output.
func bruteForceLeftSlice(s string, k int) string {
	// Step 1: strip dashes and upper-case in a single pass.
	var clean strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' {
			continue // dashes carry no information; drop them
		}
		if c >= 'a' && c <= 'z' {
			c -= 32 // ASCII lower→upper ('a'-'A' == 32)
		}
		clean.WriteByte(c)
	}
	cs := clean.String()
	if len(cs) == 0 {
		return "" // nothing but dashes → empty result
	}

	// Step 3: size of the (possibly short) first group.
	firstLen := len(cs) % k
	if firstLen == 0 {
		firstLen = k // exact multiple → first group is a full k, not empty
	}

	// Step 4-5: emit the first group, then successive full k-groups.
	groups := []string{cs[:firstLen]}        // leading, possibly-short chunk
	for i := firstLen; i < len(cs); i += k { // every remaining chunk is exactly k
		groups = append(groups, cs[i:i+k])
	}
	return strings.Join(groups, "-")
}

// ── Approach 2: Build From the Right, One Char at a Time (Optimal) ────────────
//
// rightToLeftBuild scans the original string from the end, appending each
// alphanumeric character (upper-cased) and dropping in a dash every k
// characters — then reverses once.
//
// Intuition:
//
//	Because groups are anchored at the RIGHT, the cleanest construction is to
//	walk s backwards. Maintain a counter of how many real characters we have
//	placed since the last dash; each time it reaches k, insert a dash before
//	the next character. This automatically leaves the leftmost group short
//	(whatever is left over) and needs no length arithmetic or slicing. Build
//	the answer reversed, then flip it once at the end.
//
// Algorithm:
//  1. Walk i from len(s)-1 down to 0.
//  2. Skip dashes. For each real char: if count > 0 and count % k == 0, append
//     '-'; then append the upper-cased char and increment count.
//  3. Reverse the accumulated bytes to restore left-to-right order.
//
// Time:  O(n) — one backward pass plus one reversal.
// Space: O(n) — the output buffer.
func rightToLeftBuild(s string, k int) string {
	buf := make([]byte, 0, len(s)) // reversed output accumulator
	count := 0                     // real chars placed since the last dash

	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c == '-' {
			continue // ignore existing separators
		}
		if count > 0 && count%k == 0 {
			buf = append(buf, '-') // completed a group of k → separator goes here
		}
		if c >= 'a' && c <= 'z' {
			c -= 32 // normalise to uppercase
		}
		buf = append(buf, c)
		count++ // one more real character in the current group
	}

	// buf currently holds the answer reversed; flip it in place.
	for l, r := 0, len(buf)-1; l < r; l, r = l+1, r-1 {
		buf[l], buf[r] = buf[r], buf[l]
	}
	return string(buf)
}

func main() {
	fmt.Println("=== Approach 1: Clean, Then Slice From the Left (Brute Force) ===")
	fmt.Printf("s=%q k=4  got=%q  expected \"5F3Z-2E9W\"\n", "5F3Z-2e-9-w", bruteForceLeftSlice("5F3Z-2e-9-w", 4)) // "5F3Z-2E9W"
	fmt.Printf("s=%q k=2  got=%q  expected \"2-5G-3J\"\n", "2-5g-3-J", bruteForceLeftSlice("2-5g-3-J", 2))         // "2-5G-3J"
	fmt.Printf("s=%q k=3  got=%q  expected \"ABC-DEF\"\n", "---abc-def", bruteForceLeftSlice("---abc-def", 3))     // "ABC-DEF"
	fmt.Printf("s=%q k=1  got=%q  expected \"\"\n", "----", bruteForceLeftSlice("----", 1))                        // "" (only dashes)

	fmt.Println("=== Approach 2: Build From the Right, One Char at a Time (Optimal) ===")
	fmt.Printf("s=%q k=4  got=%q  expected \"5F3Z-2E9W\"\n", "5F3Z-2e-9-w", rightToLeftBuild("5F3Z-2e-9-w", 4)) // "5F3Z-2E9W"
	fmt.Printf("s=%q k=2  got=%q  expected \"2-5G-3J\"\n", "2-5g-3-J", rightToLeftBuild("2-5g-3-J", 2))         // "2-5G-3J"
	fmt.Printf("s=%q k=3  got=%q  expected \"ABC-DEF\"\n", "---abc-def", rightToLeftBuild("---abc-def", 3))     // "ABC-DEF"
	fmt.Printf("s=%q k=1  got=%q  expected \"\"\n", "----", rightToLeftBuild("----", 1))                        // "" (only dashes)
}
