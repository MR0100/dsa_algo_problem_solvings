# 0299 — Bulls and Cows

> LeetCode #299 · Difficulty: Medium
> **Categories:** Hash Table, String, Counting

---

## Problem Statement

You are playing the **Bulls and Cows** game with your friend.

You write down a secret number and ask your friend to guess what the number is. When your friend makes a guess, you provide a hint with the following info:

- The number of "bulls", which are digits in the guess that are in the correct position.
- The number of "cows", which are digits in the guess that are in your secret number but are located in the wrong position. Specifically, the non-bull digits in the guess that could be rearranged such that they become bulls.

Given the secret number `secret` and your friend's guess `guess`, return *the hint for your friend's guess*.

The hint should be formatted as `"xAyB"`, where `x` is the number of bulls and `y` is the number of cows. Note that both `secret` and `guess` may contain duplicate digits.

**Example 1:**

```
Input: secret = "1807", guess = "7810"
Output: "1A3B"
Explanation: Bulls are connected with a '|' and cows are underlined:
"1807"
  |
"7810"
```

**Example 2:**

```
Input: secret = "1123", guess = "0111"
Output: "1A1B"
Explanation: Bulls are connected with a '|' and cows are underlined:
"1123"        "1123"
  |      or     |
"0111"        "0111"
Note that only one of the two unmatched 1s from the guess is counted as a cow
since the non-bull digits can only be rearranged to allow one 1 to be a bull.
```

**Constraints:**

- `1 <= secret.length, guess.length <= 1000`
- `secret.length == guess.length`
- `secret` and `guess` consist of digits only.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Digit frequency counting** — fixed size-10 arrays to tally leftover digits → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String iteration** — single index walk over equal-length strings → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two-Pass with Digit Counts | O(n) | O(1) | Easiest to explain; separate bull/cow phases |
| 2 | Single-Pass Signed Balance | O(n) | O(1) | Slickest; one loop, one array |

---

## Approach 1 — Two-Pass with Digit Counts

### Intuition
Bulls are trivial: same digit at the same index. For cows, ignore all bull positions, then tally how many of each digit remain in the secret and in the guess **separately**. For each digit `0..9`, the cows it can form is the **minimum** of its two leftover counts.

### Algorithm
1. Loop indices: if `secret[i] == guess[i]`, `bulls++`; else increment `secretCount[secret[i]]` and `guessCount[guess[i]]`.
2. `cows = Σ min(secretCount[d], guessCount[d])` over `d = 0..9`.
3. Return `"{bulls}A{cows}B"`.

### Complexity
- **Time:** O(n) — one pass over the strings plus a fixed 10-way merge.
- **Space:** O(1) — two size-10 arrays.

### Code
```go
func twoPassCount(secret string, guess string) string {
	var secretCount, guessCount [10]int
	bulls := 0

	for i := 0; i < len(secret); i++ {
		if secret[i] == guess[i] {
			bulls++
		} else {
			secretCount[secret[i]-'0']++
			guessCount[guess[i]-'0']++
		}
	}

	cows := 0
	for d := 0; d < 10; d++ {
		cows += min(secretCount[d], guessCount[d])
	}

	return strconv.Itoa(bulls) + "A" + strconv.Itoa(cows) + "B"
}
```

### Dry Run
Example 1: `secret = "1807"`, `guess = "7810"`.

| i | secret[i] | guess[i] | bull? | secretCount update | guessCount update |
|---|-----------|----------|-------|--------------------|-------------------|
| 0 | 1         | 7        | no    | s[1]=1             | g[7]=1            |
| 1 | 8         | 8        | **yes** (bulls=1) | —      | —                 |
| 2 | 0         | 1        | no    | s[0]=1             | g[1]=1            |
| 3 | 7         | 0        | no    | s[7]=1             | g[0]=1            |

Cows: digit 0 → min(1,1)=1; digit 1 → min(1,1)=1; digit 7 → min(1,1)=1. Total cows = 3.
Result = `"1A3B"`.

---

## Approach 2 — Single-Pass Signed Balance (Optimal)

### Intuition
Use **one** array `count`. Its sign for a digit tells which side currently has a surplus. When a non-bull secret digit appears, if the guess had already "requested" it (`count < 0`), a cow forms. Symmetrically for a non-bull guess digit when the secret has a surplus (`count > 0`). Each pending surplus is matched at most once.

### Algorithm
1. For each `i`:
   - if `secret[i] == guess[i]`: `bulls++`, continue.
   - `s = secret[i]`: if `count[s] < 0` → `cows++`; then `count[s]++`.
   - `g = guess[i]`: if `count[g] > 0` → `cows++`; then `count[g]--`.
2. Return `"{bulls}A{cows}B"`.

### Complexity
- **Time:** O(n) — a single pass.
- **Space:** O(1) — one size-10 array.

### Code
```go
func singlePassBalance(secret string, guess string) string {
	var count [10]int
	bulls, cows := 0, 0

	for i := 0; i < len(secret); i++ {
		if secret[i] == guess[i] {
			bulls++
			continue
		}
		s := secret[i] - '0'
		if count[s] < 0 {
			cows++
		}
		count[s]++

		g := guess[i] - '0'
		if count[g] > 0 {
			cows++
		}
		count[g]--
	}

	return strconv.Itoa(bulls) + "A" + strconv.Itoa(cows) + "B"
}
```

### Dry Run
Example 1: `secret = "1807"`, `guess = "7810"`. `count` starts all zeros.

| i | s,g | bull? | count[s]<0? cow | count after s++ | count[g]>0? cow | count after g-- | cows |
|---|-----|-------|-----------------|-----------------|-----------------|-----------------|------|
| 0 | 1,7 | no    | count[1]=0 no   | count[1]=1      | count[7]=0 no   | count[7]=-1     | 0    |
| 1 | 8,8 | yes   | —               | —               | —               | —               | 0    |
| 2 | 0,1 | no    | count[0]=0 no   | count[0]=1      | count[1]=1>0 ✔  | count[1]=0      | 1    |
| 3 | 7,0 | no    | count[7]=-1<0 ✔ | count[7]=0      | count[0]=1>0 ✔  | count[0]=0      | 3    |

bulls=1, cows=3 → `"1A3B"`.

---

## Key Takeaways
- **Bulls first, cows from leftovers**: a digit counted as a bull must never also count as a cow, so exclude bull positions before tallying.
- Cows honour multiplicity via `min(secretLeftover, guessLeftover)` per digit — the classic "how many of a digit can pair up" pattern.
- The single-pass trick collapses both count tables into **one signed array**: positive means the secret is ahead, negative means the guess is ahead, and each crossing of zero-from-the-other-side yields a cow.
- Only 10 possible digits ⇒ counting space is O(1) regardless of input length.

---

## Related Problems
- LeetCode #242 — Valid Anagram (digit/char frequency matching)
- LeetCode #383 — Ransom Note (min-count style matching)
- LeetCode #387 — First Unique Character in a String (frequency table)
