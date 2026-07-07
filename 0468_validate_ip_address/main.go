package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ── Approach 1: Split and Validate with strings package ──────────────────────
//
// splitValidate solves Validate IP Address by first deciding which family the
// string is claiming to be (a '.' means IPv4, a ':' means IPv6) and then running
// the matching, rule-by-rule validator on the pieces produced by strings.Split.
//
// Intuition:
//
//	The two formats are disjoint: IPv4 uses dots, IPv6 uses colons. So look at
//	which separator is present, split on it, require exactly the right number of
//	groups, and validate each group against that family's rules. IPv4 groups are
//	1-3 decimal digits forming 0..255 with NO leading zero (except "0"); IPv6
//	groups are 1-4 hex digits (0-9a-fA-F), leading zeros allowed. Anything else
//	is "Neither".
//
// Algorithm:
//  1. If queryIP contains '.', try IPv4: split on '.', need 4 groups, each a
//     valid decimal byte.
//  2. Else if it contains ':', try IPv6: split on ':', need 8 groups, each a
//     valid hex group.
//  3. strings.Split gives an empty string for leading/trailing/double separators
//     (e.g. "1..1.1"), which the group validators reject.
//  4. Return "IPv4", "IPv6", or "Neither".
//
// Time:  O(L) — L = len(queryIP); split and per-group checks each touch every
//
//	character a constant number of times.
//
// Space: O(L) — the slice of group substrings returned by Split.
func splitValidate(queryIP string) string {
	if strings.Contains(queryIP, ".") { // dotted → candidate IPv4
		groups := strings.Split(queryIP, ".")
		if len(groups) == 4 && allTrue(groups, isIPv4Group) {
			return "IPv4"
		}
		return "Neither"
	}
	if strings.Contains(queryIP, ":") { // coloned → candidate IPv6
		groups := strings.Split(queryIP, ":")
		if len(groups) == 8 && allTrue(groups, isIPv6Group) {
			return "IPv6"
		}
		return "Neither"
	}
	return "Neither" // neither separator present
}

// allTrue reports whether pred holds for every group (helper for readability).
func allTrue(groups []string, pred func(string) bool) bool {
	for _, g := range groups {
		if !pred(g) { // one bad group disqualifies the whole address
			return false
		}
	}
	return true
}

// isIPv4Group validates a single IPv4 octet: 1-3 digits, value 0..255, and no
// leading zero unless the group is exactly "0".
func isIPv4Group(g string) bool {
	if len(g) == 0 || len(g) > 3 { // empty ("1..1") or too long
		return false
	}
	for i := 0; i < len(g); i++ {
		if g[i] < '0' || g[i] > '9' { // only decimal digits allowed
			return false
		}
	}
	if g[0] == '0' && len(g) > 1 { // leading zero like "01" or "00" is invalid
		return false
	}
	v, _ := strconv.Atoi(g) // safe: we proved g is all digits, length <= 3
	return v <= 255         // must fit in a byte
}

// isIPv6Group validates a single IPv6 group: 1-4 hexadecimal digits.
func isIPv6Group(g string) bool {
	if len(g) == 0 || len(g) > 4 { // empty or more than 4 hex digits
		return false
	}
	for i := 0; i < len(g); i++ {
		c := g[i]
		isHex := (c >= '0' && c <= '9') ||
			(c >= 'a' && c <= 'f') ||
			(c >= 'A' && c <= 'F')
		if !isHex { // any non-hex character disqualifies the group
			return false
		}
	}
	return true // leading zeros are explicitly allowed in IPv6
}

// ── Approach 2: Manual Single-Pass Parser (No Split) ─────────────────────────
//
// manualParse solves Validate IP Address without allocating substrings: it scans
// queryIP once, tracking the current group's characteristics, and validates each
// group as it closes on a separator. This is the "interview whiteboard" version
// that avoids library helpers and shows the state machine explicitly.
//
// Intuition:
//
//	Decide the family from the first separator seen. Then walk character by
//	character maintaining, for the current group: its length, whether it is all
//	decimal / all hex, its numeric value (IPv4), and its first character (for the
//	leading-zero rule). On each separator, "close" the group by checking the
//	family rules and reset the accumulators; count separators to enforce the group
//	count. Validate the final trailing group after the loop.
//
// Algorithm:
//  1. Scan once to detect '.' vs ':'; reject if both or neither appear.
//  2. Walk the string; for IPv4 accumulate a value and digit-count per group,
//     enforcing digits-only, <=255, no-leading-zero, on each '.' boundary.
//     For IPv6 accumulate a hex-length per group, enforcing hex-only and length
//     1..4 on each ':' boundary.
//  3. Require exactly 3 dots (4 groups) / 7 colons (8 groups) and validate the
//     last group after the loop.
//
// Time:  O(L) — a single linear scan (plus one initial scan to pick the family).
// Space: O(1) — only scalar accumulators; no substrings allocated.
func manualParse(queryIP string) string {
	hasDot := strings.IndexByte(queryIP, '.') >= 0
	hasColon := strings.IndexByte(queryIP, ':') >= 0
	if hasDot && !hasColon {
		return parseIPv4(queryIP)
	}
	if hasColon && !hasDot {
		return parseIPv6(queryIP)
	}
	return "Neither" // both separators, or neither
}

// parseIPv4 validates a dotted string as IPv4 in one pass.
func parseIPv4(q string) string {
	groups := 0 // completed octets seen
	length := 0 // digits in the current octet
	value := 0  // numeric value of the current octet
	leadZero := false
	for i := 0; i < len(q); i++ {
		c := q[i]
		if c == '.' { // close the current octet
			if !ipv4GroupOK(length, value, leadZero) {
				return "Neither"
			}
			groups++
			length, value, leadZero = 0, 0, false // reset for next octet
			continue
		}
		if c < '0' || c > '9' { // IPv4 octets are decimal only
			return "Neither"
		}
		if length == 0 && c == '0' {
			leadZero = true // remember a leading zero to reject "01"
		}
		value = value*10 + int(c-'0') // build the octet value
		length++
		if length > 3 || value > 255 { // early reject overlong / oversized
			return "Neither"
		}
	}
	// close and validate the final octet (there is no trailing '.')
	if !ipv4GroupOK(length, value, leadZero) {
		return "Neither"
	}
	groups++
	if groups == 4 { // exactly four octets required
		return "IPv4"
	}
	return "Neither"
}

// ipv4GroupOK applies the octet rules given its running length/value/leadZero.
func ipv4GroupOK(length, value int, leadZero bool) bool {
	if length == 0 || length > 3 { // empty or too many digits
		return false
	}
	if leadZero && length > 1 { // "01", "007" invalid; "0" is fine
		return false
	}
	return value <= 255
}

// parseIPv6 validates a coloned string as IPv6 in one pass.
func parseIPv6(q string) string {
	groups := 0 // completed hex groups seen
	length := 0 // hex digits in the current group
	for i := 0; i < len(q); i++ {
		c := q[i]
		if c == ':' { // close the current group
			if length == 0 || length > 4 { // 1..4 hex digits required
				return "Neither"
			}
			groups++
			length = 0 // reset for next group
			continue
		}
		isHex := (c >= '0' && c <= '9') ||
			(c >= 'a' && c <= 'f') ||
			(c >= 'A' && c <= 'F')
		if !isHex { // non-hex char anywhere is invalid
			return "Neither"
		}
		length++
		if length > 4 { // early reject overlong group
			return "Neither"
		}
	}
	// validate the final group (no trailing ':')
	if length == 0 || length > 4 {
		return "Neither"
	}
	groups++
	if groups == 8 { // exactly eight groups required
		return "IPv6"
	}
	return "Neither"
}

func main() {
	fmt.Println("=== Approach 1: Split and Validate (strings package) ===")
	fmt.Printf("172.16.254.1                        got=%q  expected \"IPv4\"\n", splitValidate("172.16.254.1"))
	fmt.Printf("2001:0db8:85a3:0:0:8A2E:0370:7334   got=%q  expected \"IPv6\"\n", splitValidate("2001:0db8:85a3:0:0:8A2E:0370:7334"))
	fmt.Printf("256.256.256.256                     got=%q  expected \"Neither\"\n", splitValidate("256.256.256.256"))
	fmt.Printf("192.168.01.1                        got=%q  expected \"Neither\"\n", splitValidate("192.168.01.1"))                       // leading zero
	fmt.Printf("02001:0db8:85a3:0:0:8A2E:0370:7334  got=%q  expected \"Neither\"\n", splitValidate("02001:0db8:85a3:0:0:8A2E:0370:7334")) // 5 hex digits
	fmt.Printf("1e1.4.5.6                           got=%q  expected \"Neither\"\n", splitValidate("1e1.4.5.6"))                          // non-digit in IPv4

	fmt.Println("=== Approach 2: Manual Single-Pass Parser (No Split) ===")
	fmt.Printf("172.16.254.1                        got=%q  expected \"IPv4\"\n", manualParse("172.16.254.1"))
	fmt.Printf("2001:0db8:85a3:0:0:8A2E:0370:7334   got=%q  expected \"IPv6\"\n", manualParse("2001:0db8:85a3:0:0:8A2E:0370:7334"))
	fmt.Printf("256.256.256.256                     got=%q  expected \"Neither\"\n", manualParse("256.256.256.256"))
	fmt.Printf("192.168.01.1                        got=%q  expected \"Neither\"\n", manualParse("192.168.01.1"))
	fmt.Printf("02001:0db8:85a3:0:0:8A2E:0370:7334  got=%q  expected \"Neither\"\n", manualParse("02001:0db8:85a3:0:0:8A2E:0370:7334"))
	fmt.Printf("1.0.1.                              got=%q  expected \"Neither\"\n", manualParse("1.0.1."))                            // trailing dot → empty group
	fmt.Printf("20EE:FGb8:85a3:0:0:8A2E:0370:7334   got=%q  expected \"Neither\"\n", manualParse("20EE:FGb8:85a3:0:0:8A2E:0370:7334")) // 'G' not hex
}
