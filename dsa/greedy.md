# Greedy

## What it is

A **greedy algorithm** builds a solution one step at a time, and at every step
it makes the choice that looks best *right now* (the **locally optimal
choice**) — without looking ahead, without backtracking, and without ever
revisiting a decision. If the problem has the right structure, this sequence of
locally optimal choices provably produces a **globally optimal** answer.

Greedy is not one algorithm; it is a *design strategy*. What changes from
problem to problem is:

1. **The greedy choice** — what "best right now" means (largest value, earliest
   finish time, farthest reach, smallest cost, ...).
2. **The proof** — why that local choice can never lock you out of the global
   optimum.

The two classical properties a problem must have for greedy to be correct:

- **Greedy-choice property** — there exists an optimal solution that *starts*
  with the greedy choice. In other words, you can always "exchange" the first
  decision of any optimal solution for the greedy decision without making it
  worse (this is the basis of the standard **exchange argument** proof).
- **Optimal substructure** — after committing the greedy choice, the remaining
  problem is a smaller instance of the same problem, and an optimal solution to
  the remainder plus the greedy choice is optimal overall.

If either property fails, greedy silently produces a wrong (suboptimal) answer
— which is why greedy is the pattern where *correctness arguments matter most*
in interviews.

### Greedy vs Dynamic Programming

| | Greedy | DP |
|---|---|---|
| Choices considered per step | exactly one (the best-looking one) | all of them |
| Revisits decisions? | never | effectively yes (memoized subproblems) |
| Requirement | greedy-choice property + optimal substructure | optimal substructure only |
| Typical complexity | O(n) or O(n log n) (sort first) | O(n²), O(n·k), ... |

Rule of thumb: **every correct greedy problem is also solvable by DP** (greedy
is a DP where you can prove only one branch matters). So a safe interview path
is: sketch the DP, notice that one choice always dominates, collapse it to
greedy, and state the exchange argument. Problems like Jump Game (#55) and
Best Time to Buy and Sell Stock II (#122) are literally tagged both DP and
Greedy for this reason.

---

## How to recognise it — signals in the problem statement

Strong signals that a problem wants greedy:

- **"Minimum number of ..." / "maximum number of ..."** — min jumps, min
  intervals removed, max events attended, max profit — *especially* when the
  answer is a count or a single scalar rather than the actual combination.
- **Intervals / scheduling** — meetings, events, balloons, merging ranges.
  Almost always: *sort by one endpoint, sweep, keep a running boundary*.
- **"Reach" / "cover" / "jump"** — can you get to the end, minimum steps to
  cover a range → track the *farthest reachable point* while sweeping.
- **Canonical decomposition** — represent a number/string using the fewest
  pieces from a fixed system that is *designed* to be greedy-friendly (Roman
  numerals #12, coin change with canonical denominations, factorial number
  system #60).
- **"Pack as many as possible"** stated verbatim in the problem (Text
  Justification #68 literally says "greedy approach").
- **Order can be fixed up-front** — if sorting the input by some key makes
  overlapping/conflicting items adjacent, a single greedy sweep usually
  finishes the job (Merge Intervals #56).
- **Local decision is irreversible but provably safe** — e.g. Container With
  Most Water #11: moving the shorter wall inward discards only pairs that can
  never be better.
- **Every profitable micro-step can be taken independently** — Stock II #122:
  total profit decomposes into the sum of all positive day-to-day deltas.

Signals that greedy is probably **wrong** and you need DP/backtracking:

- Choices interact non-locally (taking item A now changes the *value* of item
  B later) — e.g. 0/1 Knapsack, Coin Change with arbitrary denominations
  (`[1, 3, 4]`, target 6: greedy gives 4+1+1 = 3 coins, optimal is 3+3 = 2).
- You must output *all* solutions or count them.
- Small constraints (n ≤ 20) hinting at exponential search.
- You try three natural greedy keys and can find a counterexample for each.

**Litmus test:** before coding, actively hunt for a counterexample to your
greedy rule. If you can't break it in 2–3 adversarial attempts *and* you can
sketch an exchange argument, proceed.

---

## General templates (Go)

### Template 1 — Sort, then single greedy sweep (intervals / scheduling)

The most common greedy shape on LeetCode. Sorting makes conflicts adjacent, so
one linear pass with a running "boundary" resolves everything.

```go
// greedySweep: sort by the key that makes the greedy choice safe,
// then scan once, maintaining a boundary/accumulator.
func greedySweep(items [][]int) [][]int {
    // 1. Sort by the decisive key.
    //    - merge intervals:        sort by START
    //    - max non-overlapping /   sort by END  (earliest finish time —
    //      min removals / arrows:   leaves the most room for the rest)
    sort.Slice(items, func(i, j int) bool { return items[i][0] < items[j][0] })

    res := [][]int{items[0]} // commit the first item — greedy never undoes it
    for _, cur := range items[1:] {
        last := res[len(res)-1]
        if cur[0] <= last[1] {
            // conflict with the committed boundary → resolve greedily
            // (merge: extend end; scheduling: skip/count a removal)
            last[1] = max(last[1], cur[1])
        } else {
            // no conflict → commit as-is and advance the boundary
            res = append(res, cur)
        }
    }
    return res
}
```

### Template 2 — Running-extremum sweep (reach / profit / Kadane-style)

No sorting; the input order *is* the timeline. Maintain the best value seen so
far and make one irreversible decision per element.

```go
// runningExtremum: one pass, O(1) space. Each element either improves the
// tracked extremum or is answered against it — never both revisited.
func runningExtremum(nums []int) int {
    best := 0            // global answer accumulator
    reach := 0           // running extremum (farthest reach / min price / max-ending-here)
    for i, v := range nums {
        if i > reach {   // greedy feasibility check: current position unreachable
            return -1    // (Jump Game: stuck → false)
        }
        reach = max(reach, i+v) // greedy update: extend the frontier as far as possible
        best = max(best, reach) // fold into the global answer
    }
    return best
}
```

Variant — **BFS-style level expansion** (Jump Game II #45): sweep with *two*
frontiers, `curEnd` (end of the current jump's range) and `farthest`; when `i`
hits `curEnd`, you are forced to take one more jump and `curEnd = farthest`.

### Template 3 — Greedy on a canonical value system (largest-piece-first)

```go
// canonicalDecompose: repeatedly take the largest denomination that fits.
// Correct ONLY when the value system is canonical (Roman numerals, factorial
// base, US coins) — prove or know this before using it!
func canonicalDecompose(num int) string {
    values := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
    symbols := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
    var sb strings.Builder
    for i, v := range values {      // scan denominations largest → smallest
        for num >= v {              // take this piece as many times as it fits
            sb.WriteString(symbols[i])
            num -= v                // shrink the remaining problem (optimal substructure)
        }
    }
    return sb.String()
}
```

### Template 4 — Greedy with a heap (when "best right now" changes dynamically)

When the greedy choice at each step depends on a *dynamic* candidate set (e.g.
"always process the smallest/largest available item"), keep candidates in a
`container/heap` priority queue: pop the best, process, push new candidates.
O(n log n) overall. (Appears in problems like Task Scheduler #621 and Merge
k Sorted Lists #23 — see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md).)

**Pseudocode skeleton common to all four templates:**

```
1. (optional) sort input by the decisive key, or build a heap
2. initialise the committed state (boundary / extremum / result)
3. for each item in order:
       a. decide using ONLY the committed state and the current item
       b. commit — update state; never look back
4. return the accumulated answer
```

---

## Worked example — Jump Game (LeetCode #55), traced step by step

**Problem:** `nums[i]` is the maximum jump length from index `i`. Starting at
index 0, can you reach the last index?

**Greedy insight:** you don't care *which* jumps you take — only *how far* it
is possible to get. Sweep left to right maintaining `maxReach`, the farthest
index reachable using any combination of jumps so far. If the sweep ever
arrives at an index beyond `maxReach`, that index is unreachable and the answer
is `false`.

**Why it's correct (exchange argument):** any strategy that reaches index `i`
can be exchanged, step for step, with the choice that maximises reach without
ever decreasing the set of reachable indices — so tracking the single maximum
loses nothing.

```go
// canJump reports whether the last index is reachable.
// Time: O(n) — one pass. Space: O(1).
func canJump(nums []int) bool {
    maxReach := 0                       // farthest index reachable so far
    for i, v := range nums {
        if i > maxReach {               // gap: we can never stand on i
            return false
        }
        maxReach = max(maxReach, i+v)   // greedy: extend the frontier
        if maxReach >= len(nums)-1 {    // early exit: last index already covered
            return true
        }
    }
    return true
}
```

**Trace on `nums = [2,3,1,1,4]`** (expected `true`):

| step | i | nums[i] | i > maxReach? | i + nums[i] | maxReach after | note |
|------|---|---------|---------------|-------------|----------------|------|
| 1 | 0 | 2 | 0 > 0? no | 2 | **2** | can reach up to index 2 |
| 2 | 1 | 3 | 1 > 2? no | 4 | **4** | 4 ≥ 4 → last index covered → return `true` |

**Trace on `nums = [3,2,1,0,4]`** (expected `false`):

| step | i | nums[i] | i > maxReach? | i + nums[i] | maxReach after | note |
|------|---|---------|---------------|-------------|----------------|------|
| 1 | 0 | 3 | no | 3 | **3** | |
| 2 | 1 | 2 | no | 3 | **3** | no improvement |
| 3 | 2 | 1 | no | 3 | **3** | no improvement |
| 4 | 3 | 0 | no | 3 | **3** | stuck at the 0 |
| 5 | 4 | 4 | **4 > 3 → return `false`** | — | — | index 4 unreachable |

Note how the greedy never enumerates jump combinations (the brute force is
O(2ⁿ), the DP is O(n²)) — one integer of state suffices.

---

## Complexity

| Shape | Time | Space |
|---|---|---|
| Pure sweep (Templates 2, 3) | O(n) | O(1) |
| Sort + sweep (Template 1) | O(n log n) | O(1)–O(n) (sort / output) |
| Heap-driven (Template 4) | O(n log n) | O(n) |

The near-universal O(n)/O(n log n) time with O(1) extra space is exactly *why*
interviewers love asking "can you do better than DP?" — greedy is usually the
intended final answer.

---

## Common pitfalls and how to avoid them

1. **Assuming greedy works without proof.** The #1 pitfall. Coin change with
   denominations `[1, 3, 4]` and target 6 breaks largest-first greedy
   (4+1+1 vs 3+3). *Avoid:* hunt for a counterexample first; if none found,
   sketch an exchange argument ("swap the first non-greedy decision of an
   optimal solution for the greedy one — the result is no worse").

2. **Sorting by the wrong key.** Interval problems flip between sort-by-start
   (merging, #56) and sort-by-end (max non-overlapping selection, min
   removals). Sorting activity-selection by *start* time is wrong. *Avoid:*
   ask "which key makes the safe choice the *first* element?" — for "fit the
   most items", earliest **finish** leaves maximum room for the rest.

3. **Confusing local metric with global objective.** In Container With Most
   Water (#11), the naive greedy "move the pointer at the *wider* side" is
   wrong; moving the **shorter** wall is the provably safe move because the
   shorter wall caps every remaining pair it could form. *Avoid:* justify each
   discard — "everything I skip is dominated by something I keep."

4. **Off-by-one on boundaries in sweep greedies.** In Jump Game II (#45) the
   loop must stop at `len(nums)-2` (jumping *from* the last index adds a bogus
   extra jump); in merge intervals, touching endpoints (`cur.start == last.end`)
   usually count as overlapping. *Avoid:* dry-run the 1-element and 2-element
   inputs by hand.

5. **Forgetting the tie-break / secondary ordering.** Many sort-then-sweep
   greedies need a secondary sort key (e.g. same start → longer end first) or
   the sweep makes an arbitrary and wrong first commit. *Avoid:* ask "what if
   two items share the primary key?" while writing the comparator.

6. **Mutating the answer retroactively.** Greedy's contract is *never look
   back*. If you find yourself patching earlier commits, the problem likely
   needs DP or a different greedy invariant. (Exception: the *bookmark* trick —
   Wildcard Matching #44 records the last `*` position as a deliberate,
   O(1)-state fallback; that is still greedy because only the single most
   recent bookmark ever matters.)

7. **Greedy on non-canonical systems.** Largest-symbol-first works for Roman
   numerals (#12) and the factorial number system (#60) only because those
   systems are constructed to be canonical. Never port that greedy to arbitrary
   denominations without re-proving it.

8. **Missing the greedy inside a "simulation" problem.** Text Justification
   (#68) looks like pure string simulation, but the packing rule ("as many
   words as fit") is a greedy commitment — mis-stating it (e.g. trying to
   balance line lengths) changes the answer.

---

## Problems in this repo that use Greedy

*(Problems 0001–0130 at the time of writing; later problems will be appended
in a follow-up pass.)*

- [0011 — Container With Most Water](/0011_container_with_most_water/README.md)
  — converging two pointers; greedily discard the shorter wall (exchange
  argument: every pair skipped is dominated).
- [0012 — Integer to Roman](/0012_integer_to_roman/README.md) — canonical
  decomposition: always take the largest value-symbol that fits (Template 3).
- [0031 — Next Permutation](/0031_next_permutation/README.md) — greedily make
  the *smallest* possible change: rightmost pivot, smallest larger successor,
  then sort the suffix ascending.
- [0044 — Wildcard Matching](/0044_wildcard_matching/README.md) — greedy `*`
  bookmarking: assume `*` matches zero characters, extend only when forced;
  the last `*` is always the best fallback.
- [0045 — Jump Game II](/0045_jump_game_ii/README.md) — greedy / BFS level
  expansion: minimum jumps = number of frontier expansions (Template 2
  variant).
- [0053 — Maximum Subarray](/0053_maximum_subarray/README.md) — Kadane's
  algorithm: greedily drop any negative running prefix.
- [0055 — Jump Game](/0055_jump_game/README.md) — farthest-reach sweep; the
  worked example above.
- [0056 — Merge Intervals](/0056_merge_intervals/README.md) — sort by start,
  then greedily extend-or-append (Template 1).
- [0060 — Permutation Sequence](/0060_permutation_sequence/README.md) — greedy
  digit picking via factorial number system: each digit chosen by
  `k / (n-1)!`.
- [0068 — Text Justification](/0068_text_justification/README.md) — greedy
  line packing: fit as many words as possible per line, then distribute
  spaces.
- [0121 — Best Time to Buy and Sell Stock](/0121_best_time_to_buy_and_sell_stock/README.md)
  — one-pass running minimum: greedily buy at the cheapest price seen so far.
- [0122 — Best Time to Buy and Sell Stock II](/0122_best_time_to_buy_and_sell_stock_ii/README.md)
  — sum every positive day-to-day delta; profit decomposes into independent
  micro-transactions.

Related references: [`/dsa/sorting.md`](/dsa/sorting.md) (most greedies start
with a sort), [`/dsa/two_pointers.md`](/dsa/two_pointers.md) (converging-pointer
greedies), [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
(dynamic greedy choice sets), [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
(Kadane-adjacent running accumulators).
