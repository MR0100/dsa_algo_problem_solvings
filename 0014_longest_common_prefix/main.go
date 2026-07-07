package main

import (
	"fmt"
	"sort"
	"strings"
)

// ── Approach 1: Horizontal Scanning ──────────────────────────────────────────
//
// horizontalScan starts with strs[0] as the prefix and progressively trims it
// until it is a prefix of every subsequent string.
//
// Intuition:
//   If prefix is "flower" and the next string is "flow", trim prefix to "flow".
//   Keep trimming until strings.HasPrefix returns true, or prefix is empty.
//
// Time:  O(S) where S = sum of all character lengths in strs.
// Space: O(1) extra beyond the output.
func horizontalScan(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, s := range strs[1:] {
		// Trim prefix from the right until it matches s's beginning.
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix
}

// ── Approach 2: Vertical Scanning ────────────────────────────────────────────
//
// verticalScan compares all strings column by column (character by character).
// The first column where any string differs (or ends) stops the scan.
//
// Intuition:
//   For each character position i, check strs[j][i] for all j.
//   If any string is shorter than i, or strs[j][i] != strs[0][i], return s[0:i].
//
// Time:  O(S) worst case, but short-circuits at the first mismatch column —
//        faster in practice when the common prefix is short.
// Space: O(1).
func verticalScan(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	for i := 0; i < len(strs[0]); i++ {
		ch := strs[0][i]
		for _, s := range strs[1:] {
			// Mismatch: either s is shorter or characters differ.
			if i >= len(s) || s[i] != ch {
				return strs[0][:i]
			}
		}
	}
	return strs[0]
}

// ── Approach 3: Sort + Compare First and Last ─────────────────────────────────
//
// sortEndpoints sorts the slice lexicographically and compares only the
// first and last strings — the LCP of all strings equals the LCP of these two.
//
// Intuition:
//   After sorting, the first and last strings are maximally different (one is
//   the lexicographic minimum, the other the maximum). The LCP of the entire
//   set is bounded by LCP(first, last): any character where they disagree will
//   also disagree in at least one other pair. So we only need one comparison.
//
// Time:  O(n log n) for the sort + O(m) for the comparison, where m = min length.
// Space: O(1) extra (sort may be in-place depending on implementation).
func sortEndpoints(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	sort.Strings(strs)
	first, last := strs[0], strs[len(strs)-1]

	i := 0
	for i < len(first) && i < len(last) && first[i] == last[i] {
		i++
	}
	return first[:i]
}

// ── Approach 4: Binary Search on Prefix Length ────────────────────────────────
//
// binarySearchLen binary-searches the length of the LCP in range [0, minLen].
// For a given length L, check if strs[0][:L] is a prefix of all strings.
//
// Intuition:
//   The LCP length has a monotonic property: if a prefix of length L works,
//   all shorter lengths also work. If length L fails, all longer lengths also
//   fail. Binary search exploits this monotonicity.
//
// Time:  O(S log m) — O(log m) iterations × O(S/n) per check, where m = min length.
// Space: O(1).
func binarySearchLen(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	// Find the minimum string length to bound the search.
	minLen := len(strs[0])
	for _, s := range strs[1:] {
		if len(s) < minLen {
			minLen = len(s)
		}
	}

	lo, hi := 0, minLen
	for lo < hi {
		mid := (lo + hi + 1) / 2 // bias up so we converge on the maximum valid length
		if isCommonPrefix(strs, mid) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return strs[0][:lo]
}

// isCommonPrefix checks whether strs[0][:length] is a prefix of every string.
func isCommonPrefix(strs []string, length int) bool {
	prefix := strs[0][:length]
	for _, s := range strs[1:] {
		if !strings.HasPrefix(s, prefix) {
			return false
		}
	}
	return true
}

func main() {
	examples := []struct {
		strs   []string
		expect string
	}{
		{[]string{"flower", "flow", "flight"}, "fl"},
		{[]string{"dog", "racecar", "car"}, ""},
		{[]string{"interview", "interact", "interior"}, "inter"},
		{[]string{"a"}, "a"},
		{[]string{"ab", "ab"}, "ab"},
	}

	approaches := []struct {
		name string
		fn   func([]string) string
	}{
		{"Approach 1: Horizontal Scan        O(S)       T | O(1) S", horizontalScan},
		{"Approach 2: Vertical Scan        ✅ O(S)       T | O(1) S", verticalScan},
		{"Approach 3: Sort + Compare Ends    O(n log n) T | O(1) S", sortEndpoints},
		{"Approach 4: Binary Search on Len   O(S log m) T | O(1) S", binarySearchLen},
	}

	for _, ex := range examples {
		fmt.Printf("strs=%v  expect=%q\n", ex.strs, ex.expect)
		for _, ap := range approaches {
			// Work on a copy so sort doesn't mutate the shared slice.
			cp := make([]string, len(ex.strs))
			copy(cp, ex.strs)
			fmt.Printf("  %-60s → %q\n", ap.name, ap.fn(cp))
		}
		fmt.Println()
	}
}
