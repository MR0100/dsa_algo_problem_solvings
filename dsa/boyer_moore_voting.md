# Boyer–Moore Majority Vote

> **Core idea:** find an element that appears **more than n/k times** in `O(n)` time and `O(k)` extra space — without sorting and without a hash map — by keeping `k-1` running (candidate, counter) pairs that cancel each other out.
> **Two phases:** (1) a linear *voting* scan produces the candidate(s); (2) a second *verification* scan confirms each candidate actually clears the threshold.

---

## What it is

The Boyer–Moore Majority Vote algorithm answers: **"which element appears more than ⌊n/2⌋ times?"** in a single pass, using only two variables — a `candidate` and a `count`.

The intuition is *pairwise cancellation*. Imagine every element is a voter. Whenever two voters disagree, they cancel each other and both leave the room. A strict majority element (> n/2 copies) can never be fully cancelled: even if every one of its copies is paired against a different element, there are simply not enough other elements to eliminate them all. So whoever is left standing at the end must be the majority — *if one exists*.

That last clause is the crux: the voting phase always **produces** a candidate, but only guarantees it is correct **when a majority actually exists**. On input with no majority it still returns *some* element (garbage). That is why the honest, general version of the algorithm is **two-phase**: vote to get the candidate, then re-scan to verify.

### Generalisation to > n/k (the "k-1 counters" theorem)

The n/2 case uses **1** counter. The pattern generalises beautifully:

> There can be **at most `k-1`** elements that each appear **more than `n/k`** times (because `k` elements each appearing `> n/k` times would need `> n` slots total).

So to find all elements appearing more than `n/k` times, maintain **`k-1` (candidate, counter) pairs**. The cancellation rule becomes: an incoming element either matches an existing candidate (increment its counter), fills an empty slot (become a new candidate), or — if it matches none and all slots are occupied — **decrements every counter by one** (a `k`-way mutual cancellation; when a counter hits 0 its slot is freed).

- `k = 2` → 1 pair → classic Majority Element (#169).
- `k = 3` → 2 pairs → Majority Element II (> n/3, at most 2 answers, #229).

Verification is even more essential here: the voting phase yields `k-1` *suspects*, but the problem may have fewer (or zero) real answers, so the second scan filters them.

---

## When to recognise it

| Signal in the problem | Why Boyer–Moore fits |
|-----------------------|----------------------|
| "element that appears **more than ⌊n/2⌋** times" | The definition of a strict majority → 1-counter vote |
| "all elements appearing **more than ⌊n/3⌋** times" | > n/k with k=3 → 2-counter generalisation, ≤ 2 answers |
| **O(1) / constant extra space required** | Hash-map counting is O(n) space; sorting mutates or costs O(n log n). Voting is O(k) space, O(n) time |
| "the majority is **guaranteed to exist**" | You can even skip verification (phase 1 alone is correct) |
| **Streaming / single-pass** constraints | Only a handful of counters are retained; the input need not be stored |
| "find the element left after cancelling unequal pairs" | This *is* the cancellation model, sometimes disguised |

**When *not* to reach for it:** you need the *count* of every element, the top-K by frequency for general K (use a hash map + heap / bucket sort), or the majority is defined as ≥ n/2 (not strictly >). For "> n/k" with large k, `k-1` counters still work but a hash map is usually simpler once k stops being tiny.

---

## General template / pseudocode

### Phase 1 + 2 — classic majority (> n/2)

```go
// majorityElement returns the element appearing more than len(nums)/2 times.
// Assumes such an element exists; if not, verify with a second pass (below).
//
// Phase 1 (vote): one candidate, one counter, cancelling unequal elements.
// Time: O(n)   Space: O(1)
func majorityElement(nums []int) int {
    candidate, count := 0, 0
    for _, x := range nums {
        if count == 0 {
            candidate = x // no one standing → x takes the floor
        }
        if x == candidate {
            count++ // a supporter arrives
        } else {
            count-- // an opponent cancels one supporter
        }
    }
    return candidate
}
```

If a majority is **not** guaranteed, add the verification pass:

```go
// Phase 2 (verify): confirm the candidate truly exceeds n/2.
func majorityElementChecked(nums []int) (int, bool) {
    cand := majorityElement(nums)
    cnt := 0
    for _, x := range nums {
        if x == cand {
            cnt++
        }
    }
    if cnt > len(nums)/2 {
        return cand, true
    }
    return 0, false // no strict majority exists
}
```

### Generalised — all elements > n/k (here k = 3, so ≤ 2 answers)

```go
// majorityElementII returns every element appearing more than len(nums)/3 times.
// Maintains k-1 = 2 (candidate, counter) pairs.
// Time: O(n)   Space: O(1)
func majorityElementII(nums []int) []int {
    // Two candidate slots. cand1 != cand2 is kept as an invariant.
    var cand1, cand2 int
    var cnt1, cnt2 int

    // ── Phase 1: voting ──────────────────────────────────────────────
    for _, x := range nums {
        switch {
        case cnt1 > 0 && x == cand1:
            cnt1++ // matches candidate 1
        case cnt2 > 0 && x == cand2:
            cnt2++ // matches candidate 2
        case cnt1 == 0:
            cand1, cnt1 = x, 1 // slot 1 is free → claim it
        case cnt2 == 0:
            cand2, cnt2 = x, 1 // slot 2 is free → claim it
        default:
            // x matches neither and both slots are taken:
            // 3-way cancellation — drop one vote from each.
            cnt1--
            cnt2--
        }
    }

    // ── Phase 2: verification (both suspects may be false) ───────────
    cnt1, cnt2 = 0, 0
    for _, x := range nums {
        if x == cand1 {
            cnt1++
        } else if x == cand2 { // else-if: don't double count if cand2 was never set
            cnt2++
        }
    }

    res := []int{}
    if cnt1 > len(nums)/3 {
        res = append(res, cand1)
    }
    if cnt2 > len(nums)/3 {
        res = append(res, cand2)
    }
    return res
}
```

Two subtleties encoded above, both classic bug sources:
- The `cnt > 0 && x == cand` order matters — checking `x == cand1` *before* confirming the slot is occupied can wrongly match a stale candidate whose counter is 0.
- Slot-claiming must come **after** the match checks, otherwise an element equal to an existing candidate could open a second slot for the same value, breaking the `cand1 != cand2` invariant.

---

## Worked example — step-by-step trace

### Classic (> n/2): `nums = [2, 2, 1, 1, 1, 2, 2]`, n = 7, majority needs > 3 copies.

| step | x | before (cand,count) | rule applied | after (cand,count) |
|------|---|---------------------|--------------|--------------------|
| 1 | 2 | (_, 0) | count==0 → adopt 2; match → +1 | (2, 1) |
| 2 | 2 | (2, 1) | match → +1 | (2, 2) |
| 3 | 1 | (2, 2) | mismatch → −1 | (2, 1) |
| 4 | 1 | (2, 1) | mismatch → −1 | (2, 0) |
| 5 | 1 | (2, 0) | count==0 → adopt 1; match → +1 | (1, 1) |
| 6 | 2 | (1, 1) | mismatch → −1 | (1, 0) |
| 7 | 2 | (1, 0) | count==0 → adopt 2; match → +1 | (2, 1) |

Candidate = **2**. Verify: 2 appears 4 times, 4 > 3 ✓. Answer `2`.

Notice how the candidate flipped twice (2 → 1 → 2) yet still landed on the true majority — the cancellations that dethroned it were always paid for by minority elements, and 2 had one copy to spare.

### Generalised (> n/3): `nums = [1, 1, 1, 3, 3, 2, 2, 2]`, n = 8, need > 2 copies.

| step | x | (cand1,cnt1) | (cand2,cnt2) | action |
|------|---|--------------|--------------|--------|
| 1 | 1 | (1,1) | (_,0) | slot1 free → claim 1 |
| 2 | 1 | (1,2) | (_,0) | match cand1 |
| 3 | 1 | (1,3) | (_,0) | match cand1 |
| 4 | 3 | (1,3) | (3,1) | slot2 free → claim 3 |
| 5 | 3 | (1,3) | (3,2) | match cand2 |
| 6 | 2 | (1,2) | (3,1) | matches neither, both full → −1, −1 |
| 7 | 2 | (1,1) | (3,0) | matches neither, both full → −1, −1 |
| 8 | 2 | (1,1) | (2,1) | slot2 free (cnt2==0) → claim 2 |

Suspects: **1** and **2**. Verify over the array: 1 appears 3× (3 > 2 ✓), 2 appears 3× (3 > 2 ✓), 3 appears 2× (would-be false positive, correctly *not* a suspect here). Answer `[1, 2]`.

This trace shows exactly why verification is non-negotiable: at step 8, `3` had been fully cancelled and evicted even though it was a legitimate candidate earlier — had the counts differed slightly, a non-answer could have survived in a slot and the second pass is what rejects it.

---

## Complexity

| Variant | Time | Extra space | Passes |
|---------|------|-------------|--------|
| Majority > n/2 (vote only, majority guaranteed) | O(n) | O(1) | 1 |
| Majority > n/2 (vote + verify) | O(n) | O(1) | 2 |
| Majority > n/k (`k-1` counters, vote + verify) | O(n·k) | O(k) | 2 |

- **Time** — each phase is a single linear scan. For the generalised version, every element is compared against up to `k-1` candidates, giving `O(n·k)`; with `k` a small constant (2 or 3) this is `O(n)`.
- **Space** — only the counters are stored: `O(1)` for the classic case, `O(k)` in general. This is the whole selling point over the hash-map approach's `O(n)`.

---

## Common pitfalls

1. **Skipping verification when no majority is guaranteed.** Phase 1 *always* returns a candidate. On `[1, 2, 3]` the classic algorithm happily returns `3` — which is not a majority. If the problem does not promise a majority, the second pass is mandatory.
2. **Checking `x == candidate` before `count == 0`.** The adopt-a-new-candidate step must run when the counter is empty; interleaving the comparisons wrong can strand a stale candidate. Follow the template's order.
3. **(Generalised) opening a second slot for a value already held.** Match against *both* existing candidates before claiming an empty slot, or you can end up with `cand1 == cand2`, silently halving your capacity.
4. **(Generalised) forgetting the `else if` in verification.** If you count `x == cand1` and `x == cand2` with two independent `if`s and the two candidates happen to be equal (from a bug, or an unset slot defaulting to 0 that matches a real 0 in the data), you double-count. Use `else if`, and be wary of zero-value default candidates matching a literal `0` in the array.
5. **Using it for ≥ n/2 or for general top-K.** The proof relies on a **strict** majority (> n/k). "At least half" or "the two most frequent regardless of threshold" are different problems — reach for hashing/heaps.
6. **Assuming the candidate's counter equals its frequency.** `count` at the end is a *residual* after cancellations, not the element's actual count. Never report `count` as the frequency; recompute it in phase 2 if you need it.

---

## Problems in this repo that use it

- [0169 — Majority Element](/0169_majority_element/README.md) — the canonical single-counter vote (> n/2); majority is guaranteed, so phase 1 alone suffices, with verification shown as the robust variant.
- [0229 — Majority Element II](/0229_majority_element_ii/README.md) — the `k = 3` generalisation with two (candidate, counter) pairs and a mandatory verification pass; at most two elements can exceed n/3.

### Related classics to know

- LeetCode #169 / #229 — the two above are the standard interview pair for this technique.
- The generalisation ("find all elements > n/k") is a common follow-up: the answer is "keep `k-1` counters, then verify."
