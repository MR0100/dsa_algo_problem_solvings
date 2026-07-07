# Fisher–Yates (Knuth) Shuffle

> **What it does:** permutes an array **uniformly at random** — every one of the `n!` orderings is equally likely — **in place**, in **O(n)** time with **O(1)** extra space, using exactly `n-1` random draws.
> **The one rule to remember:** at step `i`, swap element `i` with a random element in `[i, n)` — the range must **include `i` itself**, and must **never reach back into the already-fixed prefix**.

---

## What it is

The Fisher–Yates shuffle (modernised by Durstenfeld and popularised by Knuth in *TAOCP*, hence "Knuth shuffle") produces a uniformly random permutation of a finite sequence. It walks the array once; at each position it swaps the current element with a uniformly chosen element from the *unshuffled remainder*.

Two equivalent directions:

- **Forward:** for `i = 0 .. n-2`, pick `j` uniformly in `[i, n)` and swap `a[i], a[j]`. After step `i`, position `i` is finalised.
- **Backward (Durstenfeld's, the most common):** for `i = n-1 .. 1`, pick `j` uniformly in `[0, i]` and swap `a[i], a[j]`. After step `i`, position `i` is finalised.

It sits in the family of **randomized / sampling algorithms** alongside reservoir sampling and random-pick problems — all built on the same "each candidate gets exactly the right probability, maintained incrementally" reasoning.

### Why every permutation is equally likely (`n!` outcomes, each with probability `1/n!`)

Think of the algorithm as *filling positions one at a time from the pool of remaining elements*:

- Position 0 is chosen uniformly from all `n` elements → probability `1/n` for any specific element.
- Position 1 is chosen uniformly from the remaining `n-1` → probability `1/(n-1)` for any specific remaining element.
- … position `k` from the remaining `n-k` → `1/(n-k)`.

For a **specific** target permutation, the probability that the algorithm produces it is the product of those independent choices:

```
1/n · 1/(n-1) · 1/(n-2) · … · 1/1  =  1/n!
```

There are exactly `n!` permutations, each hit with probability `1/n!`, so the distribution is uniform. Equivalently, a clean **induction**: assume the last `k` positions already hold a uniformly random `k`-permutation of the whole set; the next backward step picks position `n-k-1`'s occupant uniformly from all not-yet-placed elements, extending the invariant to `k+1`. Base case `k = 1` is trivially uniform. Hence the final array is a uniform `n`-permutation.

### The naive version and its off-by-one bias

The tempting-but-wrong variant swaps with a random index over the **entire** array every time:

```go
// BIASED — do NOT use. j ranges over [0, n) instead of [i, n).
func badShuffle(a []int) {
    n := len(a)
    for i := 0; i < n; i++ {
        j := rand.Intn(n) // WRONG: should be i + rand.Intn(n-i)
        a[i], a[j] = a[j], a[i]
    }
}
```

Why it's biased: it makes `n` independent choices out of `n` options, so it can produce `n^n` equally-likely *execution paths*. But there are only `n!` permutations, and `n^n` is **not divisible by `n!`** for `n ≥ 3`. By pigeonhole, some permutations are reachable via more execution paths than others, so they come out **more often** — the distribution is provably non-uniform. Concretely for `n = 3`: `3^3 = 27` execution paths spread over `3! = 6` permutations — `27/6 = 4.5`, not an integer, so the outcomes cannot be equal. (A famous real-world instance of this bug skewed a browser vendor's "random" ballot ordering.)

The correct version draws from a **shrinking** range (`n·(n-1)·…·1 = n!` execution paths), which matches the `n!` permutations exactly — one path per outcome — giving perfect uniformity.

---

## When to recognise it

| Signal in the problem | Why Fisher–Yates |
|-----------------------|------------------|
| "Return a **random permutation** / shuffle the array so all orderings are equally likely" | this is the textbook use (#384) |
| "Randomly reorder in place, O(1) extra space" | Fisher–Yates is the only O(n)/O(1) uniform shuffle |
| "Deal / draw without replacement", "random sample of k *distinct* elements" | run a **partial** shuffle for `k` steps — the first `k` slots are a uniform sample (Fisher–Yates truncated) |
| "Reset to original, then reshuffle repeatedly" | keep an immutable original copy; shuffle a working copy (#384's `reset`/`shuffle` API) |
| "Random pick / equal probability among a stream of unknown length" | that's **reservoir sampling**, a close cousin — see the contrast below |

**When it's *not* this:** if you need a *sample from a stream whose length you don't know in advance*, use **reservoir sampling** (you can't index a random position if you don't yet know `n`). If you need a random element *with replacement*, a single `rand.Intn(n)` suffices — no shuffle needed.

---

## General template / pseudocode

### Canonical in-place shuffle (backward / Durstenfeld)

```go
import "math/rand"

// shuffle permutes a uniformly at random, in place, in O(n) time / O(1) space.
//
// Backward pass: i counts down from the last index. At each i we pick j
// uniformly from [0, i] — the still-unfixed prefix PLUS i itself — and swap.
// After the swap, a[i] is finalised and never touched again.
func shuffle(a []int) {
    for i := len(a) - 1; i > 0; i-- {
        // rand.Intn(i+1) returns a uniform int in [0, i]; note the +1 so j can equal i.
        j := rand.Intn(i + 1)
        a[i], a[j] = a[j], a[i] // an element may swap with itself (j == i) — that's fine
    }
}
```

### Forward pass (equivalent)

```go
// Forward: i counts up; pick j uniformly from [i, n) — i itself plus the untouched suffix.
func shuffleForward(a []int) {
    n := len(a)
    for i := 0; i < n-1; i++ {
        j := i + rand.Intn(n-i) // uniform in [i, n)
        a[i], a[j] = a[j], a[i]
    }
}
```

### The #384 "Shuffle an Array" design (reset + repeatable shuffle)

Keep the pristine original; always shuffle a copy so `reset()` is O(n) and `shuffle()` is independent each call.

```go
import "math/rand"

type Solution struct {
    original []int // never mutated after construction
    current  []int // the working copy we shuffle
}

func Constructor(nums []int) Solution {
    orig := make([]int, len(nums))
    copy(orig, nums)
    cur := make([]int, len(nums))
    copy(cur, nums)
    return Solution{original: orig, current: cur}
}

// Reset restores and returns the original configuration.
func (s *Solution) Reset() []int {
    copy(s.current, s.original) // O(n) restore from the untouched original
    return s.current
}

// Shuffle returns a uniformly random permutation of the array.
func (s *Solution) Shuffle() []int {
    for i := len(s.current) - 1; i > 0; i-- {
        j := rand.Intn(i + 1)
        s.current[i], s.current[j] = s.current[j], s.current[i]
    }
    return s.current
}
```

### Partial shuffle — a uniform sample of `k` distinct elements in O(k)

Stop after `k` steps; the first `k` positions are a uniform `k`-subset in random order.

```go
// sampleK returns k uniformly-chosen distinct elements (as a fresh slice),
// touching only k elements of a copy — O(k) work beyond the O(n) copy.
func sampleK(a []int, k int) []int {
    b := make([]int, len(a))
    copy(b, a)
    n := len(b)
    for i := 0; i < k; i++ {
        j := i + rand.Intn(n-i) // uniform in [i, n)
        b[i], b[j] = b[j], b[i]
    }
    return b[:k]
}
```

---

## Worked example

Shuffle `a = [A, B, C, D]` (n = 4) with the backward pass. Suppose the RNG yields the draws below (each `j = rand.Intn(i+1)`):

| Step `i` | Range `[0, i]` | Drawn `j` | Swap | Array after | Finalised |
|----------|----------------|-----------|------|-------------|-----------|
| start | — | — | — | `[A, B, C, D]` | — |
| i = 3 | [0,3] | 1 | swap a[3]↔a[1] | `[A, D, C, B]` | position 3 = **B** |
| i = 2 | [0,2] | 2 | swap a[2]↔a[2] (self) | `[A, D, C, B]` | position 2 = **C** |
| i = 1 | [0,1] | 0 | swap a[1]↔a[0] | `[D, A, C, B]` | position 1 = **A** |
| (i = 0 loop ends — position 0 is whatever remains) | | | | `[D, A, C, B]` | position 0 = **D** |

Result: `[D, A, C, B]`. Observations that show *why* it's uniform:

- At `i = 3`, each of the 4 elements had probability `1/4` of landing in slot 3 (any `j ∈ {0,1,2,3}`).
- Given slot 3 is fixed, at `i = 2` each of the remaining 3 elements had probability `1/3` for slot 2 — and so on.
- The product `1/4 · 1/3 · 1/2 · 1 = 1/24 = 1/4!` is the probability of this exact ordering, identical for every ordering.
- The **self-swap** at `i = 2` (`j == i`) is not wasted — it's a legitimate outcome (element C "chose" to stay), and excluding it (drawing from `[0, i-1]`) is exactly the bias to avoid.

---

## Complexity

| Aspect | Cost | Reason |
|--------|------|--------|
| Time (shuffle) | **O(n)** | one pass, one RNG draw + one swap per element |
| Extra space | **O(1)** | swaps happen in place (the #384 copy is O(n) storage for the API, not for the algorithm) |
| RNG draws | exactly **n − 1** | positions `n-1` down to `1`; position `0` needs no draw |
| Partial `k`-sample | **O(k)** time | stop after `k` steps |
| `reset()` in #384 | **O(n)** | copy the pristine original back |

---

## Common pitfalls

1. **Wrong random range → the classic off-by-one bias.** Drawing `j` from `[0, n)` (or `[0, i)` instead of `[0, i]`) every iteration is *not* uniform. Forward pass must use `[i, n)`; backward pass must use `[0, i]` (i.e. `rand.Intn(i+1)` — the `+1` is essential). The self-swap is a feature, not a bug.

2. **Iterating the wrong direction with the wrong bound.** Backward loop `for i := n-1; i > 0; i--` pairs with `rand.Intn(i+1)`. If you instead count *up* you must switch to `i + rand.Intn(n-i)`. Mixing an up-loop with `rand.Intn(i+1)` re-introduces bias.

3. **Off-by-one at the boundary.** The backward loop stops at `i > 0` (position 0 is settled by elimination); an inclusive `i >= 0` does one pointless `rand.Intn(1) == 0` self-swap — harmless but a sign of a copy-paste slip. Forward loop stops at `i < n-1` for the same reason.

4. **Sorting with a random comparator instead of shuffling.** `sort.Slice(a, func(i,j) bool { return rand.Intn(2) == 0 })` is **not** a uniform shuffle — comparison sorts assume a consistent order, and a random comparator yields a skewed, implementation-dependent distribution (and can even violate sort invariants). Use Fisher–Yates.

5. **Forgetting to preserve the original for `reset` (#384).** If `shuffle` mutates the only copy you hold, `reset` can't restore it. Keep an immutable `original` and shuffle a separate `current`.

6. **Reusing / re-seeding the RNG wrongly.** Re-seeding `rand` with a fixed value before each shuffle makes every "shuffle" identical. In Go 1.20+ the top-level `rand` funcs are auto-seeded; for reproducible tests seed a *local* `*rand.Rand` once, not per call. Also prefer one shared source over constructing a new `rand.New` inside a hot loop.

7. **Using a modulo-biased RNG.** Rolling your own `rawRandom() % (i+1)` can bias toward small remainders when the raw range isn't a multiple of `i+1`. `rand.Intn` already handles this correctly — use it rather than manual modulo.

8. **Assuming shuffle == sample-with-replacement.** A shuffle draws **without** replacement (a permutation, all distinct positions). If the problem wants repeated independent picks, that's `rand.Intn(n)` per pick, a different tool.

---

## Relationship to reservoir sampling

Fisher–Yates and **reservoir sampling** are cousins in the randomized-sampling family, distinguished by whether you know `n` up front and whether you can random-access the data:

| | Fisher–Yates shuffle | Reservoir sampling |
|---|----------------------|--------------------|
| Goal | uniform *permutation* of all `n` (or a `k`-sample from a known array) | uniform `k`-*sample* from a stream of **unknown / unbounded** length |
| Needs `n` known? | **yes** — it indexes random positions in `[i, n)` | **no** — that's its whole reason to exist |
| Access pattern | random access, in place | single sequential pass, keep a size-`k` reservoir |
| Core step | swap `a[i]` with a random earlier/later element | the `t`-th item replaces a reservoir slot with prob `k/t` |
| Space | O(1) in place | O(k) for the reservoir |

If you *can* materialise the whole array and know `n`, a **partial Fisher–Yates** (`k` steps) samples `k` distinct elements in O(k) and is simpler. Reservoir sampling wins precisely when you *can't* — streaming data, unknown length, or memory too small to hold everything. See [`/dsa/reservoir_sampling.md`](/dsa/reservoir_sampling.md) for the streaming side of the family.

---

## Problems in this repo that use it

- [0384 — Shuffle an Array](/0384_shuffle_an_array/README.md) — the canonical Fisher–Yates problem: `reset()` restores the pristine array, `shuffle()` returns a uniform permutation in O(n)/O(1).

### Related randomized-sampling problems in this repo

- [0382 — Linked List Random Node](/0382_linked_list_random_node/README.md) — size-1 reservoir sampling (the streaming cousin; can't index a list of unknown length).
- [0398 — Random Pick Index](/0398_random_pick_index/README.md) — reservoir sampling to pick a uniform index among matches in one pass.
- [0381 — Insert Delete GetRandom O(1) — Duplicates Allowed](/0381_insert_delete_getrandom_o1_duplicates_allowed/README.md) — uniform pick from a backing slice via `rand.Intn(n)`.

### Related classics to know (may not be in repo)

- LeetCode #470 — Implement Rand10() Using Rand7() (rejection sampling — building a uniform draw from a smaller one, the bias-avoidance mindset).
- LeetCode #528 — Random Pick with Weight (prefix sums + binary search — non-uniform sampling).
