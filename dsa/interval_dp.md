# Interval DP (Range DP)

> **Also known as:** range DP, "DP on subarrays / substrings".
> **State shape:** `dp[i][j]` = the best answer for the sub-interval `[i, j]`.
> **Fill order:** by **increasing interval length**, so every smaller interval a transition
> needs is already computed.
> **Transition:** pick a **split point** `k` (or a **last / first operated element**) inside
> `[i, j]` and combine the two sides — an O(n) choice, giving the classic **O(n³)** family.

---

## What it is

Interval DP solves optimisation problems over a **contiguous range** where the answer for a
range is built from the answers of its **sub-ranges**. The signature state is two-dimensional:

```
dp[i][j] = optimal value achievable on the interval [i, j]
```

The recurrence chooses some **pivot** `k` with `i ≤ k ≤ j` that partitions the work, then
combines the two resulting sub-intervals:

```
dp[i][j] = best over k of ( dp[i][k-1]  ⊕  cost(k)  ⊕  dp[k+1][j] )
```

where `⊕` and `cost(k)` are problem-specific. The pivot is interpreted as one of:

- a **split point** — cut `[i, j]` into `[i, k]` and `[k+1, j]` (matrix-chain multiplication,
  palindrome partitioning, optimal BST); or
- the **last element to be operated on** — everything else in `[i, j]` is resolved *before* `k`,
  so `k`'s two neighbours are exactly `i-1` and `j+1` (Burst Balloons is the canonical example —
  thinking "which balloon bursts **last**" is what makes it tractable); or
- a **first split into two halves** that are then solved recursively and possibly swapped
  (Scramble String).

Two things make the fill order non-obvious and are the crux of getting interval DP right:

1. **You cannot fill row by row or column by column naively.** `dp[i][j]` depends on strictly
   *shorter* intervals (`[i][k-1]`, `[k+1][j]`), so you must iterate by **increasing length**:
   length 1 (or the base case) first, then 2, then 3, … up to `n`. Only then are the pieces
   ready.
2. **The base case sits on the diagonal** (`dp[i][i]`, single elements) or just off it
   (`dp[i][i+1]`, adjacent pairs), and the final answer is the **top-right corner** `dp[0][n-1]`
   (the whole interval).

### Mental model

Picture the DP table as an **upper triangle** (only `i ≤ j` matters). You fill it in **bands
parallel to the main diagonal**: the diagonal itself (length-1 intervals), then the next band
up (length-2), sweeping toward the single top-right cell that answers the whole problem. Each
cell asks: "given that both shorter pieces on either side of some pivot are already solved,
which pivot is best?"

---

## When to recognise it

| Signal in the problem | Why interval DP fits |
|-----------------------|----------------------|
| The answer for a range is composed of answers for **sub-ranges** of a **contiguous** array/string | Direct `dp[i][j]` decomposition. |
| "Choose an order of operations on a sequence" — **matrix-chain multiplication**, merging stones, adding parentheses | The last (or split) operation partitions the range into two independent halves. |
| "**Burst / remove** elements one at a time; each removal's cost depends on its **current neighbours**" (#312) | Reframe as *which element is handled **last*** — then its neighbours are the fixed interval ends, decoupling the two sides. |
| "Partition a string into pieces each satisfying a property; minimise cuts / count ways" (#132) | `dp` over prefixes with a palindrome-interval helper `isPal[i][j]` filled by increasing length. |
| "Optimal binary search tree / minimise worst-case guessing cost over `[i, j]`" (#375) | Try each value as the **root/guess**; cost = value + worst of the two sub-intervals. Interval DP + minimax. |
| Palindrome questions on substrings — **longest palindromic subsequence**, count palindromic substrings, min insertions to make a palindrome | `dp[i][j]` built from `dp[i+1][j-1]` (shrink from both ends). |
| "Can string A be a **scrambled** version of string B?" (#87) | Split at every `k`; recurse on matched halves, and on swapped halves. Interval-style (though keyed by `(i, j, length)`). |

**When *not* to use it:** the sub-problems overlap on a **non-contiguous** index set (that is
usually a subset / bitmask DP, not interval); or a single **left-to-right 1-D DP** already
captures the recurrence (e.g. longest-increasing-subsequence, coin change) — no need for a 2-D
interval table. Reach for interval DP specifically when *a range's answer needs a split inside
that range*.

---

## General templates (Go)

### Template A — canonical interval DP (increasing-length sweep)

The reusable skeleton. Fill by length, choose the best split `k`.

```go
// intervalDP fills dp[i][j] = best answer on [i, j] over a sequence of length n.
// Replace `combine` and `cost` with the problem's rule.
func intervalDP(n int, cost func(i, k, j int) int) [][]int {
    dp := make([][]int, n)
    for i := range dp {
        dp[i] = make([]int, n)
        // base case: dp[i][i] for single elements (0 here; set per problem)
    }

    // length = size of the interval being solved (2, 3, ..., n).
    for length := 2; length <= n; length++ {
        for i := 0; i+length-1 < n; i++ {
            j := i + length - 1
            best := 1 << 30 // +inf for a minimisation; use -inf / 0 for maximisation
            // Try every split / pivot k inside [i, j].
            for k := i; k <= j; k++ {
                // Left = dp[i][k-1], Right = dp[k+1][j]; guard the ends.
                left, right := 0, 0
                if k-1 >= i {
                    left = dp[i][k-1]
                }
                if k+1 <= j {
                    right = dp[k+1][j]
                }
                val := left + cost(i, k, j) + right
                if val < best {
                    best = val
                }
            }
            dp[i][j] = best
        }
    }
    return dp
}
```

### Template B — Burst Balloons (#312): pivot = the element burst **last**

The insight that unlocks the problem: if balloon `k` is the **last** one popped in the open
interval `(i, j)`, then at the moment it pops its neighbours are exactly `i` and `j` (all
others in between are already gone). That makes the two sides independent.

```go
// maxCoins returns the maximum coins from bursting all balloons.
// Trick: pad with 1s at both ends, and let k be the LAST balloon burst in (i, j).
func maxCoins(nums []int) int {
    n := len(nums)
    // vals = [1, nums..., 1]; the padding 1s are the "virtual" boundary balloons.
    vals := make([]int, n+2)
    vals[0], vals[n+1] = 1, 1
    copy(vals[1:], nums)

    // dp[i][j] = max coins from bursting all balloons strictly between i and j.
    dp := make([][]int, n+2)
    for i := range dp {
        dp[i] = make([]int, n+2)
    }

    // Iterate by increasing gap between the open boundaries i and j.
    for length := 2; length <= n+1; length++ {
        for i := 0; i+length <= n+1; i++ {
            j := i + length
            best := 0
            // k = the last balloon to burst inside (i, j).
            for k := i + 1; k < j; k++ {
                // When k bursts last, its neighbours are the boundaries i and j.
                coins := vals[i]*vals[k]*vals[j] + dp[i][k] + dp[k][j]
                if coins > best {
                    best = coins
                }
            }
            dp[i][j] = best
        }
    }
    return dp[0][n+1]
}
```

### Template C — Palindrome Partitioning II (#132): interval helper + prefix DP

Two DPs cooperate: an interval table `isPal[i][j]` (built by increasing length) feeds a 1-D
"minimum cuts for the prefix ending at `i`" DP.

```go
// minCut returns the minimum number of cuts so every piece of s is a palindrome.
func minCut(s string) int {
    n := len(s)
    // isPal[i][j] = true iff s[i..j] is a palindrome.
    isPal := make([][]bool, n)
    for i := range isPal {
        isPal[i] = make([]bool, n)
    }
    // Fill by increasing length so isPal[i+1][j-1] is ready before isPal[i][j].
    for length := 1; length <= n; length++ {
        for i := 0; i+length-1 < n; i++ {
            j := i + length - 1
            if s[i] == s[j] && (length <= 2 || isPal[i+1][j-1]) {
                isPal[i][j] = true
            }
        }
    }

    // cuts[i] = minimum cuts needed for the prefix s[0..i].
    cuts := make([]int, n)
    for i := 0; i < n; i++ {
        if isPal[0][i] {
            cuts[i] = 0 // whole prefix is a palindrome: no cut
            continue
        }
        cuts[i] = i // worst case: cut before every character
        for j := 1; j <= i; j++ {
            // If s[j..i] is a palindrome, one cut after cuts[j-1] suffices.
            if isPal[j][i] && cuts[j-1]+1 < cuts[i] {
                cuts[i] = cuts[j-1] + 1
            }
        }
    }
    return cuts[n-1]
}
```

### Template D — Scramble String (#87): split, recurse, and try the swap

Interval-flavoured but keyed by `(i, j, length)` — two start indices and a shared length. At
each split you check both the **no-swap** and **swapped** pairings of the halves.

```go
// isScramble reports whether s2 is a scrambled version of s1 (equal length).
func isScramble(s1, s2 string) bool {
    memo := map[string]bool{}
    var solve func(a, b string) bool
    solve = func(a, b string) bool {
        if a == b {
            return true
        }
        if len(a) != len(b) || !sameLetters(a, b) {
            return false // quick prune: different multiset of characters
        }
        key := a + "#" + b
        if v, ok := memo[key]; ok {
            return v
        }
        n := len(a)
        res := false
        for k := 1; k < n && !res; k++ { // k = size of the left part
            // Case 1: no swap — left with left, right with right.
            if solve(a[:k], b[:k]) && solve(a[k:], b[k:]) {
                res = true
            }
            // Case 2: swap — left of a matches the RIGHT (last k) of b, etc.
            if solve(a[:k], b[n-k:]) && solve(a[k:], b[:n-k]) {
                res = true
            }
        }
        memo[key] = res
        return res
    }
    return solve(s1, s2)
}
```

---

## Worked example — Burst Balloons on `nums = [3, 1, 5]` (#312)

Pad to `vals = [1, 3, 1, 5, 1]` (indices 0..4). `dp[i][j]` = max coins from bursting everything
strictly between boundaries `i` and `j`.

Base: intervals with no balloon between them (`j = i+1`) give `dp[i][i+1] = 0`.

**Length 2 (one balloon between the boundaries):**

- `dp[0][2]` — only balloon `k=1` (value 3) between boundaries 0 and 2:
  `vals[0]*vals[1]*vals[2] = 1*3*1 = 3`. → `dp[0][2] = 3`.
- `dp[1][3]` — only balloon `k=2` (value 1) between 1 and 3:
  `vals[1]*vals[2]*vals[3] = 3*1*5 = 15`. → `dp[1][3] = 15`.
- `dp[2][4]` — only balloon `k=3` (value 5) between 2 and 4:
  `vals[2]*vals[3]*vals[4] = 1*5*1 = 5`. → `dp[2][4] = 5`.

**Length 3 (two balloons between the boundaries):**

- `dp[0][3]` — boundaries 0 and 3, candidates `k ∈ {1, 2}` (burst **last**):
  - `k=1` last: `vals[0]*vals[1]*vals[3] + dp[0][1] + dp[1][3] = 1*3*5 + 0 + 15 = 30`.
  - `k=2` last: `vals[0]*vals[2]*vals[3] + dp[0][2] + dp[2][3] = 1*1*5 + 3 + 0 = 8`.
  - → `dp[0][3] = max(30, 8) = 30`.
- `dp[1][4]` — boundaries 1 and 4, candidates `k ∈ {2, 3}`:
  - `k=2` last: `vals[1]*vals[2]*vals[4] + dp[1][2] + dp[2][4] = 3*1*1 + 0 + 5 = 8`.
  - `k=3` last: `vals[1]*vals[3]*vals[4] + dp[1][3] + dp[3][4] = 3*5*1 + 15 + 0 = 30`.
  - → `dp[1][4] = max(8, 30) = 30`.

**Length 4 (the whole thing), `dp[0][4]`** — boundaries 0 and 4, candidates `k ∈ {1, 2, 3}`:

| `k` (burst last) | `vals[0]*vals[k]*vals[4]` | `+ dp[0][k]` | `+ dp[k][4]` | total |
|------------------|---------------------------|--------------|--------------|-------|
| 1 (val 3) | `1*3*1 = 3` | `dp[0][1]=0` | `dp[1][4]=30` | **33** |
| 2 (val 1) | `1*1*1 = 1` | `dp[0][2]=3` | `dp[2][4]=5`  | 9 |
| 3 (val 5) | `1*5*1 = 5` | `dp[0][3]=30`| `dp[3][4]=0`  | 35 |

`dp[0][4] = max(33, 9, 35) = 35`. ✓

Interpretation of the winner (`k=3` last): burst 3 and 1 first (in the best internal order),
*then* the 5 last with both padding-1 boundaries as neighbours. Coins:
`3·1·5 (=15, bursting the 1) + 3·1 (…)` — the DP has already accounted for the optimal internal
order; the total maximum is **35**.

---

## Complexity

Let `n` be the sequence length.

| Quantity | Value | Reason |
|----------|-------|--------|
| Number of states | **O(n²)** | one `dp[i][j]` per interval, upper triangle. |
| Work per state | **O(n)** | scan every split / pivot `k` in `[i, j]`. |
| **Time (typical)** | **O(n³)** | states × per-state pivot scan. |
| **Space** | **O(n²)** | the `dp` table (plus any `isPal` helper). |

Specifics:

- **Burst Balloons (#312), min-cut (#132) inner scan, Guess Number II (#375):** O(n³) time,
  O(n²) space.
- **Palindrome-substring interval DPs** (longest palindromic subsequence, `isPal` table):
  O(n²) time and space — the transition is O(1) (`dp[i+1][j-1]`), no pivot scan.
- **Scramble String (#87):** O(n⁴) in the worst case — O(n³) `(i, j, len)` states × O(n) split
  choices — with memoisation keeping it from re-exploring; strong character-multiset pruning
  makes it fast in practice.

O(n³) is fine for `n` up to a few hundred (LeetCode's usual interval-DP bound). For larger
`n`, look for a **Knuth–Yao** optimisation (monotone split points) that can shave the middle
loop to bring some O(n³) DPs down to O(n²).

---

## Common pitfalls

1. **Wrong fill order.** Iterating `i` then `j` in plain increasing order reads `dp[k+1][j]`
   before it exists. You **must** iterate by increasing interval **length** (or memoise
   recursively so dependencies resolve on demand). This is the #1 interval-DP bug.

2. **Split-point vs last-element framing.** For Burst Balloons, choosing "which balloon to
   burst **first**" couples the two sides (a first burst changes both neighbours). Choosing
   "which balloon bursts **last**" decouples them — its neighbours are the fixed boundaries.
   Picking the wrong framing makes the recurrence non-decomposable.

3. **Boundary padding.** #312 relies on virtual `1`s at both ends so edge balloons have
   well-defined neighbours. Forgetting the padding forces messy special cases for `i=0` and
   `j=n-1` and usually produces wrong sums.

4. **Off-by-one on the pivot range and sub-interval indices.** Depending on whether the pivot
   is *inside* `[i, j]` (split) or an *open-interval* index (Burst Balloons uses `k` in
   `(i, j)`), the loop bounds and the `dp[i][k]/dp[k][j]` vs `dp[i][k-1]/dp[k+1][j]` shapes
   differ. Write out a length-2 and length-3 case to pin them down.

5. **Base-case placement.** Length-1 intervals (`dp[i][i]`) or empty ranges must be seeded
   correctly — 0 for "no work", `+∞` for an impossible minimisation, or a problem-specific
   value. A wrong diagonal poisons everything above it.

6. **Palindrome table dependency direction.** `isPal[i][j]` needs `isPal[i+1][j-1]` (the inner
   substring), so it too must be filled by **increasing length**, and the `length <= 2` special
   case (single char / adjacent pair) must short-circuit before indexing `i+1 > j-1`.

7. **Confusing interval DP with subset DP.** If the sub-problems are over arbitrary subsets
   rather than contiguous ranges, `dp[i][j]` is the wrong model — you likely need bitmask DP.
   Interval DP requires the *contiguity* of `[i, j]`.

8. **Memory when only the answer is needed.** The full O(n²) table is required for most
   interval DPs (you read arbitrary sub-intervals), so the usual "keep two rows" 1-D trick
   does **not** apply. Budget O(n²) space.

---

## Problems in this repo that use it

- [0312 — Burst Balloons](/0312_burst_balloons/README.md) — the canonical "pivot = element burst **last**" interval DP with boundary padding; O(n³).
- [0132 — Palindrome Partitioning II](/0132_palindrome_partitioning_ii/README.md) — `isPal[i][j]` interval table (filled by length) feeding a 1-D min-cut prefix DP.
- [0375 — Guess Number Higher or Lower II](/0375_guess_number_higher_or_lower_ii/README.md) — interval DP + minimax: `dp[i][j] = min over k of (k + max(dp[i][k-1], dp[k+1][j]))`; see also `dsa/game_theory.md`.
- [0087 — Scramble String](/0087_scramble_string/README.md) — split at every `k`, recurse on matched **and** swapped halves; keyed by `(i, j, length)` with multiset pruning + memo.

### Related in this repo (sibling / non-interval variant)

- [0131 — Palindrome Partitioning](/0131_palindrome_partitioning/README.md) — the *enumerate all* partitions sibling of #132 (backtracking, not interval DP), but it uses the same `isPal[i][j]` interval helper to prune.

### Related classics to know (not yet in repo)

- LeetCode #516 — Longest Palindromic Subsequence (`dp[i][j]` from `dp[i+1][j-1]`, O(n²))
- LeetCode #1039 — Minimum Score Triangulation of Polygon (textbook split-point interval DP)
- LeetCode #1000 — Minimum Cost to Merge Stones (interval DP with a `k`-way merge constraint)
- LeetCode #1547 — Minimum Cost to Cut a Stick (add sentinel ends, then interval DP over cuts)
- LeetCode #96 / #312-style optimal-BST — Knuth–Yao O(n²) speed-up of an O(n³) interval DP
