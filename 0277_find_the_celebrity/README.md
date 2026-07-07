# 0277 — Find the Celebrity

> LeetCode #277 · Difficulty: Medium
> **Categories:** Graph, Two Pointers, Interactive

---

## Problem Statement

Suppose you are at a party with `n` people labeled from `0` to `n - 1` and among
them there may exist one celebrity. The definition of a celebrity is that all
the other `n - 1` people know the celebrity, but the celebrity does not know any
of them.

Now you want to find out who the celebrity is or verify that there is not one.
You are only allowed to ask questions like: "Hi, A. Do you know B?" to get
information about whether A knows B. You need to find out the celebrity (or
verify there is not one) by asking as few questions as possible (in the
asymptotic sense).

You are given a helper function `bool knows(a, b)` which tells you whether A
knows B. Implement a function `int findCelebrity(n)`. There will be exactly one
celebrity if they are at the party.

Return _the celebrity's label if there is a celebrity at the party_. If there is
no celebrity, return `-1`.

**Example 1:**

```
Input: graph = [[1,1,0],[0,1,0],[1,1,1]]
Output: 1
Explanation: There are three persons labeled with 0, 1 and 2.
graph[i][j] = 1 means person i knows person j, otherwise graph[i][j] = 0 means
person i does not know person j. The celebrity is the person labeled as 1
because both 0 and 2 know him but 1 does not know anybody.
```

**Example 2:**

```
Input: graph = [[1,0,1],[1,1,0],[0,1,1]]
Output: -1
Explanation: There is no celebrity.
```

**Constraints:**

- `n == graph.length == graph[i].length`
- `2 <= n <= 100`
- `graph[i][j]` is `0` or `1`.
- `graph[i][i] == 1`

**Follow up:** If the maximum number of allowed calls to the API `knows` is
`3 * n`, could you find a solution without exceeding the maximum number of calls?

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Facebook  | ★★★★☆ High       | 2023          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Google    | ★★★☆☆ Medium     | 2022          |
| LinkedIn  | ★★★☆☆ Medium     | 2022          |
| Microsoft | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Graph modeling** — the `knows` relation is a directed graph; the celebrity
  is the unique node with in-degree `n-1` and out-degree `0` → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Two Pointers / Elimination** — one candidate pointer sweeps the array,
  eliminating one person per comparison → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Check the definition directly for each candidate |
| 2 | Two-Pass Elimination (Optimal) | O(n) | O(1) | Narrow to one candidate, then verify — meets the 3n bound |

---

## Approach 1 — Brute Force

### Intuition
Person `c` is the celebrity iff, for every other person `i`, `knows(i, c)` is
true (everyone knows `c`) and `knows(c, i)` is false (`c` knows no one). Just
test the definition for every candidate.

### Algorithm
1. For each candidate `c` in `0..n-1`:
2. For each other `i`: if `knows(c, i)` OR `!knows(i, c)`, disqualify `c`.
3. If `c` survived all checks, return `c`.
4. If nobody qualifies, return `-1`.

### Complexity
- **Time:** O(n²) — `n` candidates × `n` checks (each `knows` is O(1)).
- **Space:** O(1) — no extra structures.

### Code
```go
func bruteForce(n int) int {
	for c := 0; c < n; c++ { // try every person as the celebrity
		isCelebrity := true
		for i := 0; i < n; i++ {
			if i == c {
				continue // skip self
			}
			// c must NOT know i, and i MUST know c
			if knows(c, i) || !knows(i, c) {
				isCelebrity = false
				break // this candidate is disqualified
			}
		}
		if isCelebrity {
			return c
		}
	}
	return -1 // no celebrity present
}
```

### Dry Run
`graph = [[1,1,0],[0,1,0],[1,1,1]]`, `n = 3`:

| c | i | knows(c,i) | knows(i,c) | Verdict |
|---|---|-----------|-----------|---------|
| 0 | 1 | true | — | `knows(0,1)` true → 0 disqualified |
| 1 | 0 | false | true | ok |
| 1 | 2 | false | true | ok → **1 is celebrity** |

Return **1**.

---

## Approach 2 — Two-Pass Elimination (Optimal)

### Intuition
Every single `knows(a, b)` call eliminates exactly one person:
- if `knows(a, b)` is **true**, `a` knows someone → `a` can't be the celebrity.
- if `knows(a, b)` is **false**, nobody-knows-`b` → `b` can't be the celebrity.

Sweep one `candidate` pointer across everyone; after `n-1` calls only one person
can still possibly be the celebrity. But "possible" ≠ "confirmed" (there may be
no celebrity at all), so verify the survivor against everyone.

### Algorithm
1. `candidate = 0`.
2. For `i = 1..n-1`: if `knows(candidate, i)` then `candidate = i` (old one out).
3. Verify `candidate`: it must know no one and be known by all others.
4. Return `candidate` if verified, else `-1`.

### Complexity
- **Time:** O(n) — `n-1` calls to narrow + up to `2n` to verify (≤ `3n` total).
- **Space:** O(1) — a single pointer.

### Code
```go
func twoPointers(n int) int {
	candidate := 0 // start by assuming person 0 might be the celebrity
	// Phase 1: narrow down to one candidate.
	for i := 1; i < n; i++ {
		if knows(candidate, i) {
			// candidate knows someone → candidate is disqualified; i survives
			candidate = i
		}
		// else i is disqualified (nobody-knows-i via candidate); keep candidate
	}
	// Phase 2: verify the survivor is really the celebrity.
	for i := 0; i < n; i++ {
		if i == candidate {
			continue
		}
		// celebrity knows no one AND is known by everyone
		if knows(candidate, i) || !knows(i, candidate) {
			return -1 // survivor failed verification → no celebrity
		}
	}
	return candidate
}
```

### Dry Run
`graph = [[1,1,0],[0,1,0],[1,1,1]]`, `n = 3`:

**Phase 1 (narrow):**

| i | candidate before | knows(candidate,i) | candidate after |
|---|------------------|--------------------|-----------------|
| 1 | 0 | knows(0,1)=true | 1 |
| 2 | 1 | knows(1,2)=false | 1 |

Survivor candidate = 1.

**Phase 2 (verify candidate 1):**

| i | knows(1,i) | knows(i,1) | ok? |
|---|-----------|-----------|-----|
| 0 | false | true | yes |
| 2 | false | true | yes |

Verified → return **1**.

---

## Key Takeaways

- The elimination trick — "one comparison rules out exactly one candidate" — is
  the reusable core; it appears in majority-element voting and tournament
  problems too.
- Narrowing to a single candidate is O(n) but you still must **verify**, because
  the narrowing only guarantees "no one else can be it", not "this one is it".
- The `3n` call budget = `n-1` (narrow) + `≤ 2n` (verify), matching the follow-up.

---

## Related Problems

- LeetCode #169 — Majority Element (Boyer–Moore, same one-cancels-one idea)
- LeetCode #997 — Find the Town Judge (in-degree n-1, out-degree 0 restated)
- LeetCode #1971 — Find if Path Exists in Graph (directed-graph reasoning)
