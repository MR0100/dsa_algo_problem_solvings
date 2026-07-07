# 0423 — Reconstruct Original Digits from English

> LeetCode #423 · Difficulty: Medium
> **Categories:** Hash Table, Math, String, Counting

---

## Problem Statement

Given a string `s` containing an out-of-order English representation of digits `0-9`, return *the digits in ascending order*.

**Example 1:**

```
Input: s = "owoztneoer"
Output: "012"
```

**Example 2:**

```
Input: s = "fviefuro"
Output: "45"
```

**Constraints:**

- `1 <= s.length <= 10^5`
- `s[i]` is one of the characters `["e","g","f","i","h","o","n","s","r","u","t","w","v","x","z"]`.
- `s` is **guaranteed** to be valid.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Frequency counting (letter histogram)** — the input is a bag of letters; the whole solution is bookkeeping over a 26-entry count array → see [`/dsa/counting_sort.md`](/dsa/counting_sort.md)
- **Hash Table / character tally** — mapping each letter to its remaining count and reasoning about which digit-word "owns" it → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Math (elimination by uniqueness)** — spotting that `z,w,u,x,g` uniquely identify `0,2,4,6,8`, then peeling the rest in dependency order → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Backtracking** — the brute-force decomposition searches for an order to remove digit words that empties the letter multiset → see [`/dsa/backtracking.md`](/dsa/backtracking.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking | Exponential (worst) | O(#digits) | Correctness oracle / tiny inputs; unusable for `n = 10⁵` |
| 2 | Unique-Letter Counting (Optimal) | O(n) | O(1) | The intended answer; one pass and O(1) arithmetic |
| 3 | Counting with Self-Check | O(n) | O(1) | Same math, plus a reconstruction assertion to prove correctness |

---

## Approach 1 — Backtracking

### Intuition

Model `s` as a multiset of letters. Every spelled digit word (e.g. `"two"`) can be *peeled off* only when all its letters are still available; peeling subtracts them. A valid answer is any sequence of peels that empties the multiset. A depth-first search tries each removable digit, recurses, and undoes on failure. Because the answer must be in ascending digit order, we sort the collected digits at the end (order of removal doesn't matter for the final multiset).

### Algorithm

1. Build `count[26]` from `s`.
2. DFS:
   1. If all counts are zero → success; record the digits chosen along this path.
   2. Otherwise for each digit `0..9` whose word is fully available: subtract its letters, recurse; if it succeeds stop, else add the letters back and try the next.
3. On the first success, sort the recorded digits and render them as a string.

### Complexity

- **Time:** Exponential in the worst case — up to 10 branches per level of the search tree; acceptable only as a reference.
- **Space:** O(depth) recursion = O(number of digits), plus the O(1) counts.

### Code

```go
func backtracking(s string) string {
	var count [26]int
	for i := 0; i < len(s); i++ {
		count[s[i]-'a']++ // tally each letter of the shuffled input
	}

	var result []int    // digits of the first successful decomposition
	var chosen []int    // digits peeled on the current DFS path
	var dfs func() bool // returns true once the multiset is emptied

	dfs = func() bool {
		empty := true
		for _, c := range count {
			if c != 0 {
				empty = false // still letters left to consume
				break
			}
		}
		if empty {
			result = append(result, chosen...) // copy the winning path out
			return true
		}
		for d := 0; d <= 9; d++ {
			if canSubtract(&count, digitWords[d]) { // is digit d's word fully available?
				subtract(&count, digitWords[d], -1) // remove its letters
				chosen = append(chosen, d)
				if dfs() { // recurse on the smaller multiset
					return true // stop at the first complete decomposition
				}
				chosen = chosen[:len(chosen)-1] // undo the choice
				subtract(&count, digitWords[d], +1) // add the letters back
			}
		}
		return false // no digit word fits — dead end on this path
	}

	dfs()
	sort.Ints(result) // answer must be in ascending digit order
	var sb strings.Builder
	for _, d := range result {
		sb.WriteByte(byte('0' + d)) // render each digit as a character
	}
	return sb.String()
}
```

### Dry Run

Input `s = "owoztneoer"`. Letter counts: `e:2, n:1, o:3, r:1, t:1, w:1, z:1`.

| Depth | Try digit (word) | Subtractable? | count after | chosen |
|-------|------------------|---------------|-------------|--------|
| 0 | 0 `zero` (z,e,r,o) | yes         | e:1,n:1,o:2,t:1,w:1 | [0] |
| 1 | 0 `zero`          | no (`z` gone) | —          | [0] |
| 1 | 1 `one` (o,n,e)   | yes         | o:1,t:1,w:1 | [0,1] |
| 2 | 0/1 …             | no          | —           | [0,1] |
| 2 | 2 `two` (t,w,o)   | yes         | all zero    | [0,1,2] |
| 3 | multiset empty    | **success** | —           | record [0,1,2] |

Sort `[0,1,2]` → `"012"`.

---

## Approach 2 — Unique-Letter Counting (Optimal)

### Intuition

Some spelled digits contain a letter no other digit word has:

| letter | only in | digit |
|--------|---------|-------|
| `z` | zero | 0 |
| `w` | two | 2 |
| `u` | four | 4 |
| `x` | six | 6 |
| `g` | eight | 8 |

So the count of `z` **is** the number of 0s, and likewise for 2, 4, 6, 8. After those five are fixed, other letters become unique *relative to the digits still unknown*:

- `h` → three & eight → threes = `count(h) − eights`
- `f` → five & four → fives = `count(f) − fours`
- `s` → seven & six → sevens = `count(s) − sixes`
- `o` → one, zero, two, four → ones = `count(o) − zeros − twos − fours`
- `i` → nine, five, six, eight → nines = `count(i) − fives − sixes − eights`

Each digit is determined **exactly**, no search needed.

### Algorithm

1. Build the letter histogram of `s`.
2. `cnt[0]=z, cnt[2]=w, cnt[4]=u, cnt[6]=x, cnt[8]=g`.
3. `cnt[3]=h−cnt[8]`, `cnt[5]=f−cnt[4]`, `cnt[7]=s−cnt[6]`.
4. `cnt[1]=o−cnt[0]−cnt[2]−cnt[4]`, `cnt[9]=i−cnt[5]−cnt[6]−cnt[8]`.
5. Emit digit `d` exactly `cnt[d]` times, ascending.

### Complexity

- **Time:** O(n) — one pass to count, O(1) arithmetic, O(n) to write the output.
- **Space:** O(1) — fixed 26- and 10-entry arrays (output aside).

### Code

```go
func uniqueLetterCounting(s string) string {
	var c [26]int
	for i := 0; i < len(s); i++ {
		c[s[i]-'a']++ // frequency of each letter in the shuffled string
	}
	at := func(ch byte) int { return c[ch-'a'] }

	var cnt [10]int
	// Digits pinned by a letter unique to their word.
	cnt[0] = at('z') // zero: only word containing 'z'
	cnt[2] = at('w') // two:  only word containing 'w'
	cnt[4] = at('u') // four: only word containing 'u'
	cnt[6] = at('x') // six:  only word containing 'x'
	cnt[8] = at('g') // eight:only word containing 'g'

	// Digits whose defining letter is shared only with an already-known digit.
	cnt[3] = at('h') - cnt[8] // 'h' in three & eight
	cnt[5] = at('f') - cnt[4] // 'f' in five & four
	cnt[7] = at('s') - cnt[6] // 's' in seven & six

	// Digits determined after removing the contributions counted above.
	cnt[1] = at('o') - cnt[0] - cnt[2] - cnt[4] // 'o' in one, zero, two, four
	cnt[9] = at('i') - cnt[5] - cnt[6] - cnt[8] // 'i' in nine, five, six, eight

	var sb strings.Builder
	for d := 0; d <= 9; d++ {
		for k := 0; k < cnt[d]; k++ {
			sb.WriteByte(byte('0' + d)) // append digit d, cnt[d] times, in order
		}
	}
	return sb.String()
}
```

### Dry Run

Input `s = "owoztneoer"`. Histogram: `e:2, n:1, o:3, r:1, t:1, w:1, z:1` (all others 0).

| Step | Formula | Value |
|------|---------|-------|
| cnt[0] | `at('z')` | 1 |
| cnt[2] | `at('w')` | 1 |
| cnt[4] | `at('u')` | 0 |
| cnt[6] | `at('x')` | 0 |
| cnt[8] | `at('g')` | 0 |
| cnt[3] | `at('h') − cnt[8] = 0 − 0` | 0 |
| cnt[5] | `at('f') − cnt[4] = 0 − 0` | 0 |
| cnt[7] | `at('s') − cnt[6] = 0 − 0` | 0 |
| cnt[1] | `at('o') − cnt[0] − cnt[2] − cnt[4] = 3 − 1 − 1 − 0` | 1 |
| cnt[9] | `at('i') − cnt[5] − cnt[6] − cnt[8] = 0 − 0 − 0 − 0` | 0 |

Emit: one `0`, one `1`, one `2` → `"012"`.

---

## Approach 3 — Counting with Self-Check

### Intuition

The counting formulas are provably exact, but you can *demonstrate* correctness cheaply: re-spell every emitted digit, tally the letters, and assert the tally equals the original histogram. For valid input the assertion always holds; the value is identical to Approach 2. It shows how to guard against arithmetic slips in an interview.

### Algorithm

1. Compute `cnt[0..9]` exactly as in Approach 2.
2. Rebuild the expected letter histogram by summing each chosen digit's word `cnt[d]` times.
3. If the rebuilt histogram differs from the input's (never for valid input) return `""`; otherwise emit the digits ascending.

### Complexity

- **Time:** O(n) — counting plus a linear reconstruction pass.
- **Space:** O(1) — two fixed-size letter tables.

### Code

```go
func countingVerified(s string) string {
	var c [26]int
	for i := 0; i < len(s); i++ {
		c[s[i]-'a']++
	}
	at := func(ch byte) int { return c[ch-'a'] }

	var cnt [10]int
	cnt[0], cnt[2], cnt[4], cnt[6], cnt[8] = at('z'), at('w'), at('u'), at('x'), at('g')
	cnt[3] = at('h') - cnt[8]
	cnt[5] = at('f') - cnt[4]
	cnt[7] = at('s') - cnt[6]
	cnt[1] = at('o') - cnt[0] - cnt[2] - cnt[4]
	cnt[9] = at('i') - cnt[5] - cnt[6] - cnt[8]

	// Self-check: re-spell the digits and confirm the letters add back to s.
	var rebuilt [26]int
	for d := 0; d <= 9; d++ {
		for k := 0; k < cnt[d]; k++ {
			for i := 0; i < len(digitWords[d]); i++ {
				rebuilt[digitWords[d][i]-'a']++
			}
		}
	}
	if rebuilt != c {
		return "" // would indicate a logic error for valid inputs; never triggers
	}

	var sb strings.Builder
	for d := 0; d <= 9; d++ {
		for k := 0; k < cnt[d]; k++ {
			sb.WriteByte(byte('0' + d))
		}
	}
	return sb.String()
}
```

### Dry Run

Input `s = "fviefuro"`. Histogram: `f:2, v:1, i:1, e:1, u:1, r:1, o:1`.

Counting: `cnt[4]=at('u')=1` (four), `cnt[5]=at('f')−cnt[4]=2−1=1` (five); all other digits compute to 0. So digits = {4:1, 5:1}.

Self-check — re-spell `four` + `five`:

| word | letters added |
|------|----------------|
| four | f,o,u,r |
| five | f,i,v,e |

Rebuilt histogram: `f:2, o:1, u:1, r:1, i:1, v:1, e:1` — identical to the input. Assertion passes, emit `"45"`.

---

## Key Takeaways

- **Find an anchor with a unique feature, then subtract.** Whenever items overlap on features, look for a feature owned by exactly one item; it pins that item's count and unlocks a chain of subtractions.
- **Digit words hide five unique letters:** `z,w,u,x,g → 0,2,4,6,8`. Memorising this makes the problem O(n) and O(1) with zero searching.
- **Order matters for the elimination, not the output.** Compute counts in dependency order (evens first), but the final digits are just emitted ascending.
- The multiset/backtracking view is the honest brute force; the counting insight collapses that entire search into a handful of formulas.

---

## Related Problems

- LeetCode #451 — Sort Characters By Frequency (character histogram)
- LeetCode #387 — First Unique Character in a String (letter counts)
- LeetCode #383 — Ransom Note (multiset subtraction)
- LeetCode #438 — Find All Anagrams in a String (fixed-alphabet counting)
