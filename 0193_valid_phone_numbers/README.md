# 0193 — Valid Phone Numbers

> LeetCode #193 · Difficulty: Easy
> **Categories:** Shell, String, Regular Expressions

---

## Problem Statement

Given a text file `file.txt` that contains a list of phone numbers (one per line), write a one-line bash script to print all valid phone numbers.

You may assume that a valid phone number must appear in one of the following two formats: `(xxx) xxx-xxxx` or `xxx-xxx-xxxx`. (`x` means a digit)

You may also assume each line in the text file must not contain leading or trailing white spaces.

**Example:**

Assume that `file.txt` has the following content:

```
987-123-4567
123 456 7890
(123) 456-7890
```

Your script should output the following valid phone numbers:

```
987-123-4567
(123) 456-7890
```

> **Repo note:** #193 is one of LeetCode's four Shell problems. Following this
> repo's Go-only convention, every approach below re-implements the identical
> line filter in Go, treating the file's content as an in-memory string
> (`file.txt` → the `fileTxt` constant). The canonical one-line bash answers
> appear in Key Takeaways for interview completeness.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Cisco      | ★★☆☆☆ Low        | 2023          |
| Oracle     | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String Algorithms (validation / tokenisation)** — the problem is pure format validation: fixed-width position checks or separator-based token checks over each line → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (position-by-position template check) | O(N) | O(1) | No regex allowed; shows you understand exactly what the pattern demands |
| 2 | Split and Validate Tokens | O(N) | O(L) per line | Cleaner rule-per-group decomposition; easy to extend to new formats |
| 3 | Regular Expression (Optimal) | O(N) | O(1) per line | The intended answer; one anchored pattern covers both formats |

*N = total characters in the file, L = length of one line.*

---

## Approach 1 — Brute Force (Position-by-Position Template Check)

### Intuition
Both valid formats are completely rigid: every character position is either "must be a digit" or "must be this exact separator". `xxx-xxx-xxxx` is exactly 12 characters with dashes at indices 3 and 7; `(xxx) xxx-xxxx` is exactly 14 characters with `(` `)` ` ` `-` at indices 0, 4, 5, 9. So a line is valid iff its length matches one template and every index passes that template's check — this is what a regex engine would do, written out by hand.

### Algorithm
1. Split the file into lines.
2. For each line, run `matchesDashed`: length must be 12; indices 3 and 7 must be `'-'`; all other indices must be digits.
3. Failing that, run `matchesParen`: length must be 14; indices 0, 4, 5, 9 must be `'('`, `')'`, `' '`, `'-'`; the ten remaining indices must be digits.
4. Keep the line if either template accepts it.

### Complexity
- **Time:** O(N) — every character of the file is examined a constant number of times (each line at most twice).
- **Space:** O(1) beyond the output slice — only index variables, no intermediate structures.

### Code
```go
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func matchesDashed(s string) bool {
	if len(s) != 12 { // template is exactly 12 characters
		return false
	}
	for i := 0; i < 12; i++ {
		if i == 3 || i == 7 { // separator slots
			if s[i] != '-' {
				return false
			}
		} else if !isDigit(s[i]) { // every other slot must be a digit
			return false
		}
	}
	return true
}

func matchesParen(s string) bool {
	if len(s) != 14 { // template is exactly 14 characters
		return false
	}
	if s[0] != '(' || s[4] != ')' || s[5] != ' ' || s[9] != '-' {
		return false // a literal separator is out of place
	}
	for _, i := range []int{1, 2, 3, 6, 7, 8, 10, 11, 12, 13} { // digit slots
		if !isDigit(s[i]) {
			return false
		}
	}
	return true
}

func bruteForce(fileContent string) []string {
	valid := []string{}
	for _, line := range strings.Split(fileContent, "\n") { // one candidate per line
		if matchesDashed(line) || matchesParen(line) { // accept if either template fits
			valid = append(valid, line)
		}
	}
	return valid
}
```

### Dry Run
Example, line by line:

| line | len | matchesDashed | matchesParen | verdict | valid slice after |
|------|-----|---------------|--------------|---------|-------------------|
| `987-123-4567` | 12 | idx 3=`-` ✓, idx 7=`-` ✓, rest digits ✓ → **true** | (not evaluated) | keep | `[987-123-4567]` |
| `123 456 7890` | 12 | idx 3 is `' '` ≠ `'-'` → false | len 12 ≠ 14 → false | drop | `[987-123-4567]` |
| `(123) 456-7890` | 14 | len 14 ≠ 12 → false | idx 0=`(` ✓, 4=`)` ✓, 5=`' '` ✓, 9=`-` ✓, 10 digit slots ✓ → **true** | keep | `[987-123-4567, (123) 456-7890]` |

Output: `987-123-4567`, `(123) 456-7890`. ✓

---

## Approach 2 — Split and Validate Tokens

### Intuition
Instead of walking indices, think in tokens: a dashed number is exactly three `'-'`-separated groups of sizes 3/3/4; a parenthesised number is `"(ddd)"` + one space + `"ddd-dddd"`. Splitting expresses the format rules at a higher level, localises each rule to one small check, and would extend gracefully if new formats were added (e.g. country codes).

### Algorithm
1. For each line, branch on whether it starts with `'('`.
2. Parenthesised branch: `SplitN` on the first space into head + rest; head must be exactly `'('` + 3 digits + `')'` (length 5); rest must split on `'-'` into exactly two all-digit groups of lengths 3 and 4.
3. Dashed branch: split on `'-'`; accept exactly three all-digit groups with lengths 3, 3, 4.
4. Collect the lines that pass either branch.

### Complexity
- **Time:** O(N) — each line is split and scanned a constant number of times.
- **Space:** O(L) per line — the transient token slices produced by the splits.

### Code
```go
func allDigits(s string) bool {
	if len(s) == 0 {
		return false // an empty group can never be a digit block
	}
	for i := 0; i < len(s); i++ {
		if !isDigit(s[i]) {
			return false
		}
	}
	return true
}

func validBySplit(line string) bool {
	if strings.HasPrefix(line, "(") { // candidate for the (xxx) xxx-xxxx form
		parts := strings.SplitN(line, " ", 2) // split into "(xxx)" and "xxx-xxxx"
		if len(parts) != 2 {
			return false // no space → cannot match the parenthesised form
		}
		head, rest := parts[0], parts[1]
		// head must be exactly "(" + 3 digits + ")".
		if len(head) != 5 || head[0] != '(' || head[4] != ')' || !allDigits(head[1:4]) {
			return false
		}
		// rest must be exactly 3 digits, '-', 4 digits.
		groups := strings.Split(rest, "-")
		return len(groups) == 2 &&
			len(groups[0]) == 3 && allDigits(groups[0]) &&
			len(groups[1]) == 4 && allDigits(groups[1])
	}

	// Candidate for the xxx-xxx-xxxx form: exactly three digit groups 3/3/4.
	groups := strings.Split(line, "-")
	return len(groups) == 3 &&
		len(groups[0]) == 3 && allDigits(groups[0]) &&
		len(groups[1]) == 3 && allDigits(groups[1]) &&
		len(groups[2]) == 4 && allDigits(groups[2])
}

func splitAndValidate(fileContent string) []string {
	valid := []string{}
	for _, line := range strings.Split(fileContent, "\n") {
		if validBySplit(line) {
			valid = append(valid, line)
		}
	}
	return valid
}
```

### Dry Run
Example, line by line:

| line | branch | tokens | group checks | verdict |
|------|--------|--------|--------------|---------|
| `987-123-4567` | dashed (no `(` prefix) | `["987","123","4567"]` | 3 groups ✓, lens 3/3/4 ✓, all digits ✓ | keep |
| `123 456 7890` | dashed | `["123 456 7890"]` (no `-` to split on) | 1 group ≠ 3 | drop |
| `(123) 456-7890` | paren | head=`"(123)"`, rest=`"456-7890"` → `["456","7890"]` | head len 5 ✓, `(`/`)` ✓, `123` digits ✓; groups lens 3/4 ✓, digits ✓ | keep |

Output: `987-123-4567`, `(123) 456-7890`. ✓

---

## Approach 3 — Regular Expression (Optimal)

### Intuition
Both formats share the tail `xxx-xxxx` and differ only in the prefix: `(xxx) ` versus `xxx-`. One alternation captures exactly that: `^(\(\d{3}\) |\d{3}-)\d{3}-\d{4}$`. The `^...$` anchors are the crux — an unanchored pattern would happily accept garbage like `0(001) 345-0000` or `123-456-78901` because a valid substring hides inside. This is the direct translation of the intended one-line `grep -P` answer.

### Algorithm
1. Compile `^(\(\d{3}\) |\d{3}-)\d{3}-\d{4}$` once at package level (compiling per line would be wasteful).
2. For each line, keep it iff `MatchString` reports a full-line match.

### Complexity
- **Time:** O(N) — Go's RE2 engine guarantees linear-time matching with no backtracking, so the whole file is one linear scan.
- **Space:** O(1) per line — the compiled automaton is fixed-size and shared across all lines.

### Code
```go
// phoneRe encodes both formats in one anchored pattern:
//   ^        start of line (nothing before the number)
//   \(\d{3}\) ␣  →  "(xxx) "  … or …  \d{3}-  →  "xxx-"
//   \d{3}-\d{4}  →  the shared "xxx-xxxx" tail
//   $        end of line (nothing after the number)
var phoneRe = regexp.MustCompile(`^(\(\d{3}\) |\d{3}-)\d{3}-\d{4}$`)

func regexMatch(fileContent string) []string {
	valid := []string{}
	for _, line := range strings.Split(fileContent, "\n") {
		if phoneRe.MatchString(line) { // anchored full-line match
			valid = append(valid, line)
		}
	}
	return valid
}
```

### Dry Run
Example, line by line against `^(\(\d{3}\) |\d{3}-)\d{3}-\d{4}$`:

| line | prefix alternation | tail `\d{3}-\d{4}` | anchors | match? |
|------|--------------------|--------------------|---------|--------|
| `987-123-4567` | `987-` matches `\d{3}-` | `123-4567` ✓ | consumes whole line ✓ | yes → keep |
| `123 456 7890` | `123 ` — space where `-` or `(` required | — | — | no → drop |
| `(123) 456-7890` | `(123) ` matches `\(\d{3}\) ` | `456-7890` ✓ | consumes whole line ✓ | yes → keep |

Output: `987-123-4567`, `(123) 456-7890`. ✓

---

## Key Takeaways

- **Anchor your validation** — `^` and `$` (or an explicit length check) are what separate "contains a phone number" from "is a phone number"; forgetting them is the classic bug here.
- **Fixed-format validation has three equivalent shapes**: index templates (fastest, most verbose), token splits (most readable/extensible), regex (most concise). Know how to move between them when an interviewer bans one.
- **Compile regexes once** (package-level `MustCompile`), never inside the per-line loop.
- **Go's RE2 is linear-time** — no catastrophic backtracking, unlike PCRE; a useful production talking point.
- The intended one-line bash answers, worth knowing verbatim:
  ```bash
  grep -P '^(\(\d{3}\) |\d{3}-)\d{3}-\d{4}$' file.txt
  # or POSIX-portable:
  grep -E '^([0-9]{3}-|\([0-9]{3}\) )[0-9]{3}-[0-9]{4}$' file.txt
  # or: sed -n -E '/^([0-9]{3}-|\([0-9]{3}\) )[0-9]{3}-[0-9]{4}$/p' file.txt
  ```

---

## Related Problems

- LeetCode #65 — Valid Number (harder cousin: multi-branch string format validation)
- LeetCode #468 — Validate IP Address (same anchored-template validation pattern)
- LeetCode #192 — Word Frequency (Shell problem set)
- LeetCode #194 — Transpose File (Shell problem set)
- LeetCode #195 — Tenth Line (Shell problem set)
