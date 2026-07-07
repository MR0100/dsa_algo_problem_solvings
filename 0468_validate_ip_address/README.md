# 0468 — Validate IP Address

> LeetCode #468 · Difficulty: Medium
> **Categories:** String, Simulation / Parsing

---

## Problem Statement

Given a string `queryIP`, return `"IPv4"` if IP is a valid IPv4 address, `"IPv6"` if IP is a valid IPv6 address or `"Neither"` if IP is not a correct IP of any type.

**A valid IPv4** address is an IP in the form `"x1.x2.x3.x4"` where `0 <= xi <= 255` and `xi` **cannot contain** leading zeros. For example, `"192.168.1.1"` and `"192.168.1.0"` are valid IPv4 addresses while `"192.168.01.1"`, `"192.168.1.00"`, and `"192.168@1.1"` are invalid IPv4 addresses.

**A valid IPv6** address is an IP in the form `"x1:x2:x3:x4:x5:x6:x7:x8"` where:

- `1 <= xi.length <= 4`
- `xi` is a **hexadecimal string** which may contain digits, lowercase English letter (`'a'` to `'f'`) and upper-case English letters (`'A'` to `'F'`).
- Leading zeros are allowed in `xi`.

For example, `"2001:0db8:85a3:0:0:8A2E:0370:7334"` and `"2001:db8:85a3:0:0:8A2E:0370:7334"` are valid IPv6 addresses, while `"2001:0db8:85a3::8A2E:037j:7334"` and `"02001:0db8:85a3:0000:0000:8a2e:0370:7334"` are invalid IPv6 addresses.

**Example 1:**

```
Input: queryIP = "172.16.254.1"
Output: "IPv4"
Explanation: This is a valid IPv4 address, return "IPv4".
```

**Example 2:**

```
Input: queryIP = "2001:0db8:85a3:0:0:8A2E:0370:7334"
Output: "IPv6"
Explanation: This is a valid IPv6 address, return "IPv6".
```

**Example 3:**

```
Input: queryIP = "256.256.256.256"
Output: "Neither"
Explanation: This is neither a IPv4 address nor a IPv6 address.
```

**Constraints:**

- `queryIP` consists only of English letters, digits and the characters `'.'` and `':'`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Cisco      | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String parsing / tokenisation** — the task is a rules engine over tokens: choose a family by separator, split into groups, and validate each group's length, alphabet, and value; classic parsing with careful edge cases (empty groups, leading zeros, group counts) → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Split and Validate (strings package) | O(L) | O(L) | Cleanest to read; leans on `strings.Split` + per-group predicates |
| 2 | Manual Single-Pass Parser (No Split) | O(L) | O(1) | No allocations; the whiteboard state-machine version |

`L = len(queryIP)`.

---

## Approach 1 — Split and Validate (strings package)

### Intuition

The two formats are disjoint: IPv4 uses dots, IPv6 uses colons. Detect the intended family from which separator is present, split on it, require the exact group count (4 for IPv4, 8 for IPv6), and validate each group against that family's rules. Crucially, `strings.Split` produces an **empty** element for a leading/trailing/double separator (`"1..1.1"` → `["1","","1","1"]`), and the group validators reject empty strings — so malformed separators are handled for free.

### Algorithm

1. If `queryIP` contains `'.'`: split on `'.'`; accept as `"IPv4"` iff there are exactly 4 groups and each passes `isIPv4Group`.
2. Else if it contains `':'`: split on `':'`; accept as `"IPv6"` iff there are exactly 8 groups and each passes `isIPv6Group`.
3. `isIPv4Group`: 1-3 digits, digits only, no leading zero unless the group is `"0"`, value `≤ 255`.
4. `isIPv6Group`: 1-4 characters, all hexadecimal (`0-9a-fA-F`); leading zeros allowed.
5. Otherwise return `"Neither"`.

### Complexity

- **Time:** O(L) — split and per-group validation each touch every character a constant number of times.
- **Space:** O(L) — the slice of group substrings.

### Code

```go
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
```

### Dry Run

Example 2: `queryIP = "2001:0db8:85a3:0:0:8A2E:0370:7334"`.

1. Contains `'.'`? No. Contains `':'`? Yes → try IPv6.
2. Split on `':'` → `["2001","0db8","85a3","0","0","8A2E","0370","7334"]` → 8 groups ✓.
3. Validate each group with `isIPv6Group`:

| group | length | all hex? | valid? |
|-------|--------|----------|--------|
| 2001 | 4 | yes | ✓ |
| 0db8 | 4 | yes (0,d,b,8) | ✓ |
| 85a3 | 4 | yes | ✓ |
| 0 | 1 | yes | ✓ |
| 0 | 1 | yes | ✓ |
| 8A2E | 4 | yes (uppercase ok) | ✓ |
| 0370 | 4 | yes | ✓ |
| 7334 | 4 | yes | ✓ |

All 8 groups valid → return `"IPv6"` ✔

---

## Approach 2 — Manual Single-Pass Parser (No Split)

### Intuition

Same rules, but validate without allocating substrings. Decide the family from the first separator seen, then walk the string once, keeping per-group accumulators — for IPv4: digit count, running value, and a leading-zero flag; for IPv6: hex-digit count. On each separator, "close" the current group by applying the family rules and reset the accumulators; count separators to enforce the group count. Validate the trailing group after the loop (there is no trailing separator).

### Algorithm

1. Detect `'.'` vs `':'`; reject if both or neither appear.
2. **IPv4 pass:** for each char — on `'.'`, check the octet (`length 1..3`, no leading zero unless `"0"`, `value ≤ 255`) and reset; on a digit, accumulate `value` and `length` with early rejection of `length > 3` or `value > 255`; on any non-digit, reject. Require exactly 4 octets.
3. **IPv6 pass:** for each char — on `':'`, check `length 1..4` and reset; on a hex digit, increment `length` (reject `> 4`); on any non-hex, reject. Require exactly 8 groups.

### Complexity

- **Time:** O(L) — one linear scan (plus a constant scan to pick the family).
- **Space:** O(1) — only scalar accumulators; nothing allocated.

### Code

```go
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

func parseIPv4(q string) string {
	groups := 0    // completed octets seen
	length := 0    // digits in the current octet
	value := 0     // numeric value of the current octet
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
	if !ipv4GroupOK(length, value, leadZero) { // final octet, no trailing '.'
		return "Neither"
	}
	groups++
	if groups == 4 { // exactly four octets required
		return "IPv4"
	}
	return "Neither"
}

func ipv4GroupOK(length, value int, leadZero bool) bool {
	if length == 0 || length > 3 { // empty or too many digits
		return false
	}
	if leadZero && length > 1 { // "01", "007" invalid; "0" is fine
		return false
	}
	return value <= 255
}

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
	if length == 0 || length > 4 { // final group
		return "Neither"
	}
	groups++
	if groups == 8 { // exactly eight groups required
		return "IPv6"
	}
	return "Neither"
}
```

### Dry Run

Example 1: `queryIP = "172.16.254.1"`.

Family: contains `'.'`, no `':'` → `parseIPv4`.

| i | char | action | length | value | leadZero | groups |
|---|------|--------|--------|-------|----------|--------|
| 0 | 1 | digit | 1 | 1 | false | 0 |
| 1 | 7 | digit | 2 | 17 | false | 0 |
| 2 | 2 | digit | 3 | 172 | false | 0 |
| 3 | . | close octet "172" (1≤len≤3, ≤255) ✓, reset | 0 | 0 | false | 1 |
| 4 | 1 | digit | 1 | 1 | false | 1 |
| 5 | 6 | digit | 2 | 16 | false | 1 |
| 6 | . | close "16" ✓, reset | 0 | 0 | false | 2 |
| 7 | 2 | digit | 1 | 2 | false | 2 |
| 8 | 5 | digit | 2 | 25 | false | 2 |
| 9 | 4 | digit | 3 | 254 | false | 2 |
| 10 | . | close "254" ✓, reset | 0 | 0 | false | 3 |
| 11 | 1 | digit | 1 | 1 | false | 3 |

After loop: close final octet "1" ✓ → `groups = 4`. Exactly four octets → return `"IPv4"` ✔

---

## Key Takeaways

- **Choose the family by separator first**, then validate — the two formats never mix (`.` ⟹ IPv4, `:` ⟹ IPv6), so a single classifier gate simplifies everything downstream.
- **`strings.Split` exposes empty groups** for leading/trailing/consecutive separators; letting the group validator reject empty strings removes a whole class of edge cases for free.
- **The three IPv4 gotchas**: leading zero (`"01"` invalid but `"0"` valid), value `> 255`, and non-digit characters. The three IPv6 gotchas: group length outside 1-4, non-hex characters, and wrong group count. Enumerate these explicitly.
- **Parse-and-validate vs. split-and-validate** is a general trade-off: the single-pass parser is O(1) space and avoids allocations; the split version is shorter and more obviously correct. Know both.

---

## Related Problems

- LeetCode #93 — Restore IP Addresses (generate all valid IPv4 splits — backtracking cousin)
- LeetCode #65 — Valid Number (another rules-heavy string validator / state machine)
- LeetCode #8 — String to Integer (atoi) (careful character-by-character parsing)
- LeetCode #468 companion pattern — any "validate this format" task (dates, phone numbers, JSON-ish tokens)
