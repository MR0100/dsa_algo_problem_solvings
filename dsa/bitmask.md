# Bitmask (Subset-as-Integer & Bitmask DP)

> **What it is:** representing a subset of a *small* universe (≤ ~20 elements)
> as a single integer — bit `i` set ⇔ element `i` is in the subset — and then
> **enumerating**, **transitioning between**, or **memoizing over** those
> subsets. The integer *is* the subset.
> **The signature move:** a DP whose state is a mask — `dp[mask]` — where the
> memo key and the game/search state are one and the same.

---

## What it is

A bitmask is an ordinary `int` reinterpreted as a subset of a fixed universe of
up to ~20 (comfortably), ~32 (with `int32`), or ~64 (with `int`/`uint64`)
elements. Bit `i` of the mask answers one yes/no question — "is element `i`
chosen / used / present?" — so the whole mask is a compact snapshot of *which*
elements are in play.

That compactness buys two things you cannot get from a slice or a `map[...]bool`:

1. **A hashable, array-indexable state.** A mask is just an `int`, so you can use
   it directly as `dp[mask]` (an array of length `2^n`) or as a `map[int]...`
   memo key. A slice cannot be a map key in Go and a fresh `map` per state is
   far too slow — the mask sidesteps both.
2. **O(1) whole-set operations.** Union `a|b`, intersection `a&b`, difference
   `a&^b`, "add element i" `mask|(1<<i)`, "is i in?" `mask>>i&1` — each is a
   single machine instruction regardless of how many elements the set holds.

This file is specifically about **enumeration and DP over subsets encoded as
integers**. For general bit *tricks* — XOR pairing, parity, `n & (n-1)`, single-
number puzzles, power-of-two tests, swapping without a temp — see
[`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md). The two overlap (a mask
*is* bit manipulation), but the mental models differ:

| | `bit_manipulation.md` | `bitmask.md` (this file) |
|--|-----------------------|--------------------------|
| Core idea | operate on the bits of a number | treat the number as a **set** and iterate/DP over sets |
| Typical size | any width | small universe, ≤ ~20, because you touch `2^n` states |
| Headline technique | XOR / shifts / `n&(n-1)` | `for mask := 0; mask < 1<<n; mask++`, submask loops, `dp[mask]` |
| You reach for it when | "find the one non-repeated number", "count set bits" | "how many ways to partition into groups", "minimax over the set of remaining moves" |

---

## When to recognise it

Reach for a bitmask when **the state of the whole problem is 'which of a small
fixed set of things have I used / picked / visited'**, and order among them
either doesn't matter or is captured elsewhere.

| Signal in the problem | Why a bitmask fits |
|-----------------------|--------------------|
| A hard constraint like `n ≤ 15`, `n ≤ 20`, "at most 20 numbers" | The tiny bound is a nudge that `2^n` states are affordable — the classic "exponential is fine" tell |
| "In how many ways can you **choose a subset** such that …" | Enumerate all `2^n` subsets and test the predicate |
| "**Partition** the elements into groups (of equal sum / into the fewest teams / …)" | `dp[mask]` = best answer using exactly the elements in `mask`; build up from `0` |
| "Two players take turns removing items from a shared pool; who wins?" | The set of remaining items is the game state → memoize minimax on the mask |
| "Assign n tasks to n workers / seat n people" (assignment problems) | `dp[mask]` where `mask` = set of tasks already assigned; the number of set bits tells you which worker/row you're on |
| Order **does not** matter, only membership | A mask discards order automatically — the same subset reached via different sequences collapses to one state (huge dedup win) |
| Each element is used **at most once** and there are few of them | One bit per element cleanly encodes "used / unused" |
| A sub-search over "which positions to keep / cut" in a short string | Each of `2^m` keep/cut choices is a mask (abbreviations, gray-code style enumeration) |

**When *not* to use it:** the universe is large (n ≫ 25) — `2^n` explodes;
elements can be used multiple times or in unbounded quantity — one bit can't count;
or order genuinely matters and can't be recovered — a mask forgets it. In those
cases fall back to ordinary DP over indices/values, backtracking, or graph search.

---

## General templates (Go)

### Single-bit operations — the vocabulary

```go
// Universe of n elements, indices 0..n-1. mask is an int.
has   := mask>>i & 1 == 1     // is element i in the set?           (test)
mask |= 1 << i                // add element i                      (set)
mask &^= 1 << i               // remove element i (Go AND-NOT)      (clear)
mask ^= 1 << i                // toggle element i                   (flip)
low   := mask & -mask         // lowest set bit, isolated (a power of two)
i     := bits.TrailingZeros(uint(low)) // index of that lowest bit
mask &= mask - 1              // drop the lowest set bit             (pop it off)
cnt   := bits.OnesCount(uint(mask))    // popcount = size of the subset
full  := 1<<n - 1             // the "all elements chosen" mask
```

> `bits.TrailingZeros`, `bits.OnesCount`, and friends live in `math/bits` and
> compile to single CPU instructions. Prefer them over hand-rolled loops.

### Iterate every subset of an n-element universe

```go
for mask := 0; mask < 1<<n; mask++ { // all 2^n subsets, in increasing order
    for i := 0; i < n; i++ {
        if mask>>i&1 == 1 {
            // element i is present in this subset
        }
    }
    // ... evaluate this subset ...
}
```

### Iterate the set bits of one mask efficiently (skip the zeros)

```go
for m := mask; m > 0; m &= m - 1 { // visits only the set bits
    i := bits.TrailingZeros(uint(m)) // index of the current lowest set bit
    _ = i
}
```

### Iterate all **submasks** of a given mask (the `(sub-1)&mask` trick)

Enumerates every subset *of the elements already in `mask`* — the workhorse of
subset-partition DP. It runs in O(3^n) total across all masks (each element is
in one of three states: not in `mask`, in `mask` but not `sub`, in both).

```go
for sub := mask; sub > 0; sub = (sub - 1) & mask {
    comp := mask ^ sub // the complementary submask (mask split into sub | comp)
    _ = comp
    // ... use sub (a nonempty subset of mask) ...
}
// NOTE: this loop skips sub == 0. If you need the empty submask too,
// handle it once outside the loop, or use a do-while shape:
//   sub := mask
//   for {
//       // ... use sub, including 0 ...
//       if sub == 0 { break }
//       sub = (sub - 1) & mask
//   }
```

### Bitmask DP — `dp[mask]` where the memo key IS the state

The defining pattern. Solve every subset once, in an order where each `dp[mask]`
depends only on strictly smaller (fewer-bit) masks, and cache the result.

```go
// Assignment-style: dp[mask] = best cost to assign the tasks in `mask`.
// The number of tasks already assigned (popcount) tells us which worker we're on.
func assign(cost [][]int, n int) int {
    const INF = 1 << 30
    dp := make([]int, 1<<n)
    for m := 1; m < 1<<n; m++ {
        dp[m] = INF
    }
    // dp[0] = 0 already (empty assignment costs nothing).
    for mask := 0; mask < 1<<n; mask++ {
        if dp[mask] == INF {
            continue
        }
        worker := bits.OnesCount(uint(mask)) // this many tasks done ⇒ assign worker's task next
        if worker == n {
            continue
        }
        for task := 0; task < n; task++ {
            if mask>>task&1 == 0 { // task still free
                next := mask | 1<<task
                if c := dp[mask] + cost[worker][task]; c < dp[next] {
                    dp[next] = c
                }
            }
        }
    }
    return dp[1<<n-1] // all tasks assigned
}
```

### Bitmask memoization (top-down) — game / minimax on a shared pool

When the state is "which elements remain" and you branch over the current
player's moves, memoize on the mask so each reachable configuration is solved
once. `map[int]bool` (or a `[]int8` of {unknown, false, true}) is the memo.

```go
// canWin reports whether the player to move can force a win, given `used`
// = the set of already-taken elements. memo caches each mask's outcome.
func canWin(used int, /* ...target, n... */ memo map[int]int) bool {
    if v, ok := memo[used]; ok {
        return v == 1
    }
    win := false
    for i := 0; i < n; i++ {
        if used>>i&1 == 0 { // element i still available
            // Either taking i wins immediately, or it leaves the opponent
            // in a losing position (recurse with i marked used).
            if takingWinsNow(i) || !canWin(used|1<<i, memo) {
                win = true
                break
            }
        }
    }
    if win {
        memo[used] = 1
    } else {
        memo[used] = -1
    }
    return win
}
```

### Meet-in-the-middle — when n is ~40 (too big for one 2^n sweep)

`2^40` is out of reach, but `2^20` twice is not. Split the universe into two
halves, enumerate all `2^(n/2)` subset-sums of each half, sort one side, and for
every subset of the other half binary-search its complement. Turns `O(2^n)` into
`O(2^(n/2) · n)`.

```go
// Count subsets summing to target, n up to ~40.
func countSubsetSums(nums []int, target int) int {
    half := len(nums) / 2
    left, right := nums[:half], nums[half:]

    subsetSums := func(a []int) []int {
        sums := make([]int, 0, 1<<len(a))
        for mask := 0; mask < 1<<len(a); mask++ {
            s := 0
            for i := 0; i < len(a); i++ {
                if mask>>i&1 == 1 {
                    s += a[i]
                }
            }
            sums = append(sums, s)
        }
        return sums
    }

    ls, rs := subsetSums(left), subsetSums(right)
    sort.Ints(rs)
    count := 0
    for _, s := range ls { // need s + r == target ⇒ r == target - s
        want := target - s
        lo := sort.SearchInts(rs, want)
        hi := sort.SearchInts(rs, want+1)
        count += hi - lo // number of right-halves equal to `want`
    }
    return count
}
```

---

## Worked example — `dp[mask]` for Matchsticks to Square (LeetCode #473)

Problem: can the multiset of stick lengths be split into **4 groups of equal
sum**? Take `nums = [1, 1, 2, 2]`. Total = 6, so each side must sum to
`side = 6/4`… which is 1.5 — not an integer, so this particular set is trivially
`false`. Use `nums = [1, 1, 2, 2]` with target logic on a cleaner set instead:
`nums = [2, 2, 2, 2]`, total 8, `side = 2`, so every stick is its own side.

We DP on `dp[mask] = the leftover length on the *current* (partially filled)
side, given that the sticks in mask are already used` — where `dp[mask] = -1`
means "mask is unreachable / invalid". `n = 4`, `side = 2`, universe = sticks
`{0,1,2,3}` (all length 2). Masks range `0 .. 15`.

Rules per transition from a reachable `mask`:
- `dp[mask]` is the used-so-far total **modulo `side`** — i.e. how full the
  current side is. It must be `< side` for the mask to be extendable cleanly.
- To add stick `i` (with `mask>>i&1 == 0`), we require `dp[mask] + len[i] <= side`.
  The new leftover is `(dp[mask] + len[i]) % side`.

| Step | mask (bits 3210) | sticks used | dp[mask] (leftover on current side) |
|------|------------------|-------------|-------------------------------------|
| init | `0000` | {} | 0 (empty, current side 0/2 full) |
| add 0 | `0001` | {0} | (0+2)%2 = 0 (side 0 completed, next side empty) |
| add 1 | `0011` | {0,1} | 0 |
| add 2 | `0111` | {0,1,2} | 0 |
| add 3 | `1111` | {0,1,2,3} | 0 |

`dp[1111] = 0` and `1111` is the full mask, so **every side came out exactly
even → return true**. (Had any stick overshot `side`, that transition would be
rejected and the mask left at `-1`; if the full mask were unreachable, the
answer would be `false`.)

Why the mask beats plain backtracking here: the four groups are
**interchangeable**, so naive recursion revisits the same *set* of used sticks
along many different assignment orders. Keying the DP on `mask` collapses all of
those into a single cached entry — the exponential blowup is bounded at `2^n`
distinct states instead of `4^n` assignment orders.

---

## Complexity

Let `n` be the universe size (number of bits).

| Pattern | Time | Space | Reason |
|---------|------|-------|--------|
| Enumerate all subsets, inspect each bit | O(2ⁿ · n) | O(1) extra | `2^n` masks × `n` bits each |
| Enumerate all subsets, popcount/sum via `math/bits` | O(2ⁿ) amortized | O(1) | per-mask work is O(1) machine ops |
| Submask enumeration over **all** masks (`(sub-1)&mask`) | O(3ⁿ) | O(1) | element ∈ {out, in-mask-not-sub, in-both} ⇒ 3ⁿ (mask,sub) pairs |
| Bitmask DP `dp[mask]` with an inner O(n) transition | O(2ⁿ · n) | O(2ⁿ) | one state per mask, each with ≤ n edges |
| Bitmask DP with submask transition (partition DP) | O(3ⁿ) | O(2ⁿ) | sum over masks of `2^popcount(mask)` = 3ⁿ |
| Meet-in-the-middle | O(2^(n/2) · n) | O(2^(n/2)) | two halves of `2^(n/2)` subset-sums, sort + binary search |

Rule of thumb: `2^20 ≈ 10^6` (fine), `2^25 ≈ 3.3·10^7` (borderline),
`3^15 ≈ 1.4·10^7` (fine), `2^40 ≈ 10^12` (need meet-in-the-middle).

---

## Common pitfalls

1. **Operator precedence with shifts.** In Go, `<<` and `>>` bind *tighter* than
   `&`, `|`, `^`, so `mask & 1 << i` parses as `mask & (1<<i)` — usually what you
   want — but `mask>>i & 1` is `(mask>>i) & 1`, also fine. When in doubt,
   parenthesise: `(mask >> i) & 1`. A misread here silently tests the wrong bit.

2. **`1 << i` overflows for a large universe.** The untyped constant `1` is `int`;
   `1 << 40` is fine on 64-bit Go but `1 << 64` is not. For universes > 62–63
   elements you're past a single `int` anyway — that's the signal to rethink
   (bitmask is for *small* universes).

3. **Submask loop silently skips the empty set.** `for sub := mask; sub > 0; …`
   never yields `sub == 0`. If your partition DP needs the empty submask (e.g.
   "one group is empty"), handle `0` explicitly outside the loop.

4. **Wrong iteration order in bottom-up DP.** `dp[mask]` must be computed only
   after every mask it depends on. Iterating `mask` from `0` upward works when
   transitions **add** bits (dependencies have fewer bits ⇒ smaller integer).
   If your transitions *remove* bits, iterate downward — otherwise you read
   stale/zero entries.

5. **Confusing "index of set bit" with "the set-bit value".** `mask & -mask`
   gives the isolated lowest bit *as a power of two* (e.g. `0b1000` = 8), not its
   index (3). Convert with `bits.TrailingZeros`.

6. **Using popcount as the loop variable's meaning without checking it.** In
   assignment DPs the popcount of `mask` doubles as "which row/worker/step we're
   on". That only holds if you add **exactly one** element per DP layer; if a
   transition can add several bits, popcount no longer equals the step index.

7. **`map[int]` memo when a slice would do.** For `n ≤ ~22`, a `make([]int8,
   1<<n)` memo is far faster and lighter than `map[int]...`. Reach for the map
   only when the reachable masks are sparse or the universe is too big to
   allocate `2^n` slots.

8. **Forgetting three-valued memo for booleans.** A `map[int]bool` can't tell
   "computed and false" from "not yet computed" (both read as `false`). Use
   `map[int]int` with {absent, -1, +1}, a `map[int]bool` *paired* with a
   `seen` set, or a `[]int8` of `{0:unknown, 1:false, 2:true}`.

9. **Treating a bitmask as a multiset.** One bit is a boolean, not a counter.
   If an element can be used 2+ times, a single mask cannot represent it — you
   need a different encoding (base-k digits) or a different technique entirely.

---

## Problems in this repo that use Bitmask

- [0216 — Combination Sum III](/0216_combination_sum_iii/README.md) — the
  alternate approach enumerates all `2^9` subsets of the nine digits `1..9` as
  9-bit masks; keep a mask iff `popcount == k` and its digit-sum == `n`.
- [0320 — Generalized Abbreviation](/0320_generalized_abbreviation/README.md) —
  each of the `2^n` masks assigns every character to *keep* or *abbreviate*;
  iterating masks generates all abbreviations with no recursion.
- [0411 — Minimum Unique Word Abbreviation](/0411_minimum_unique_word_abbreviation/README.md)
  — encode each same-length dictionary word as a diff-mask vs. the target, then
  brute-force all `2^m` subsets of kept positions; a kept-mask is valid iff its
  AND with every diff-mask is nonzero.
- [0464 — Can I Win](/0464_can_i_win/README.md) — the set of still-available
  numbers (max ≤ 20) is a 20-bit mask; that mask is simultaneously the minimax
  branching state and the memo key (bitmask memoization).
- [0465 — Optimal Account Balancing](/0465_optimal_account_balancing/README.md)
  — subset-sum bitmask DP: partition the non-zero balances into the maximum
  number of disjoint zero-sum groups via submask enumeration; answer is
  `n − (max groups)`.
- [0473 — Matchsticks to Square](/0473_matchsticks_to_square/README.md) —
  `dp[mask]` over `2^n` states (`n ≤ 15`) tracks the leftover length on the
  current side; the interchangeable-group dedup is exactly what the mask buys.

### Related classics to know (not yet in repo)

- LeetCode #78 — Subsets · #90 — Subsets II (the plain `2^n` enumeration)
- LeetCode #698 — Partition to K Equal Sum Subsets (`dp[mask]` partition, sibling of #473)
- LeetCode #847 — Shortest Path Visiting All Nodes (BFS with a `(node, mask)` state)
- LeetCode #1349 — Maximum Students Taking Exam (row-by-row bitmask DP with adjacency masks)
- LeetCode #1799 — Maximize Score After N Operations (submask / pairing bitmask DP)
- Held–Karp Travelling Salesman — the canonical `dp[mask][last]` bitmask DP

See also [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md) for the
lower-level bit tricks (XOR, parity, `n&(n-1)`, power-of-two tests) that these
masks are built on, and [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
for DP fundamentals.
