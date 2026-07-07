# Rejection Sampling

> **What it does:** draws a **uniform** sample over a target set `T` that is
> awkward to sample directly, by sampling from a **larger, easy** set `S ⊇ T`
> and **throwing away** every draw that lands outside `T` — retrying until one
> lands inside.
> **The one guarantee:** conditioning a uniform draw on an event keeps it
> uniform. If `X` is uniform on `S`, then `X | (X ∈ T)` is uniform on `T`.
> **The one cost:** wasted draws. Expected retries = `1 / P(accept)` =
> `|S| / |T|`.

---

## What it is

Rejection sampling is a **randomized** technique for producing a uniform sample
from a distribution/region you can't (or don't want to) invert directly. You
already know how to sample *something bigger* — a square that contains your
circle, a range `[0, 49)` when you only have `rand7` — so you:

1. **Propose** a draw from the big, easy set `S`.
2. **Test** whether it lies in the target `T`.
3. **Accept** it if so (return it); **reject** and **go back to step 1** if not.

Because every point of `S` was equally likely and you keep points *only* inside
`T`, the survivors are equally likely across `T` — i.e. uniform on `T`. The
rejected tail carries away exactly the excess probability that would otherwise
bias the result.

Two flavours dominate interview problems:

- **Geometric** — sample a point in an easy-to-sample superset of a region and
  reject points outside the region. Canonical case: uniform point in a **disk**
  by sampling the bounding **square** and rejecting the corners.
- **Integer (rand-from-rand)** — build a bigger uniform *integer* range from a
  smaller RNG, keep the **largest multiple** of your target size, reject the
  leftover tail, and **reduce mod** the target. Canonical case: `rand10()` from
  `rand7()`.

### The integer pattern in one line

> **Build a bigger uniform range → keep the largest multiple of the target →
> reject the tail → reduce mod target.**

You have `randK()` uniform on `{1..K}` and want `randN()` uniform on `{1..N}`
with `N` not a nice factor of `K`. Combine several `randK` calls to get a uniform
value on `{0 .. K^t − 1}` (a range of size `M = K^t ≥ N`). Let
`usable = M − (M mod N)` — the largest multiple of `N` that fits. If the draw is
`≥ usable`, **reject** (that tail can't map to all `N` outcomes evenly); else
return `draw mod N` (+1 to shift into `{1..N}`). Every one of the `N` outcomes
gets exactly `usable / N` of the `M` equally likely values → perfectly uniform.

### Why you can't just take `mod` without rejecting

If `M` is not a multiple of `N`, then `draw mod N` hands the low residues one
extra representative each. With `rand7`→`rand10` on a raw 49-value range,
`49 mod 10 = 9` values are "extra": residues would be lopsided and the result
biased. Chopping off the top `9` values (keeping `40 = 4·10`) restores exact
uniformity. **The rejection is not optional — it is what makes the sample
uniform.**

---

## When to recognise it

| Signal in the problem | Why rejection sampling fits |
|-----------------------|-----------------------------|
| "Implement `randN()` using only `randK()`" where `N ∤ K` | Build a range of size `K^t ≥ N`, keep the largest multiple of `N`, reject the rest, take mod |
| "Generate a uniform random **point inside** a circle / disk / triangle / odd region" | Sample the bounding box (trivially uniform), reject points outside the region |
| "Uniform over a set that's a **subset** of an easy-to-sample set" | Propose from the easy superset, accept iff inside the target |
| "You may only call `<given RNG>`; no built-in random" | Forces you to *construct* the distribution — rejection is the construction tool |
| Direct inversion (CDF / closed-form transform) is ugly or unknown | Rejection needs only a membership test, not an inverse |
| The acceptance region is a clean fraction of the proposal region | Few retries → cheap; e.g. disk fills π/4 ≈ 79% of its square |

**When *not* to use it:** when a direct transform is easy and exact — e.g. a
uniform point in a disk via polar coordinates `r = R·√u, θ = 2πv` (no rejection,
constant time, note the `√` that fixes the area weighting; see
[`/dsa/geometry.md`](/dsa/geometry.md)). Also avoid it when the acceptance rate
is tiny (a thin sliver inside a huge box) — expected retries `= 1/P(accept)`
blows up; reach for a tighter proposal region or a direct method instead.

---

## Expected number of retries

Each independent proposal is accepted with probability
`p = |T| / |S|` (measure of target ÷ measure of proposal). The number of draws
until the first acceptance is a **geometric** random variable, so:

$$\mathbb{E}[\text{draws}] = \frac{1}{p} = \frac{|S|}{|T|}$$

- Disk in its bounding square: `p = πR² / (2R)² = π/4 ≈ 0.785` ⇒ ≈ **1.27** draws
  of the pair `(x, y)` on average.
- `rand10` from `rand7` via a 49-value grid keeping 40: `p = 40/49 ≈ 0.816` ⇒
  ≈ **1.225** grid draws, i.e. ≈ **2.45** `rand7()` calls (2 per grid draw).

The tail is exponentially thin: `P(more than k draws) = (1−p)^k`. With
`p ≈ 0.8`, ten straight rejections have probability `0.2^10 ≈ 10^-7` — so a
`for {}` loop terminates with probability 1 and almost always in the first few
iterations. It is **expected**-time O(1), not worst-case bounded (there is no
hard cap on iterations), which is fine for these problems.

---

## General templates (Go)

### Integer: `rand10()` from `rand7()` — build 49, keep 40, reject, mod

```go
// rand7 returns a uniform integer in [1, 7]. (Given.)
func rand7() int

// rand10 returns a uniform integer in [1, 10] using only rand7.
func rand10() int {
    for {
        // Two rand7 calls form a uniform value on the 7x7 grid: [0, 48].
        row := rand7() - 1        // 0..6
        col := rand7() - 1        // 0..6
        idx := row*7 + col        // 0..48, each of the 49 cells equally likely

        // Keep the largest multiple of 10 that fits in [0, 48], i.e. [0, 39].
        if idx < 40 {             // 40 = 4 * 10 usable values
            return idx%10 + 1     // fold to 0..9, shift to 1..10 — exactly uniform
        }
        // idx in 40..48 (9 values): the tail. Reject and re-draw.
    }
}
```

Generalised helper — any `randN` from any `randK`:

```go
// randFromRand builds randN (uniform [1,N]) from randK (uniform [1,K]).
// It draws enough randK values to cover N, then rejects the non-uniform tail.
func randFromRand(randK func() int, k, n int) int {
    for {
        m := 1     // size of the range built so far (M)
        x := 0     // uniform value in [0, m-1]
        for m < n {
            x = x*k + (randK() - 1) // extend: uniform on [0, m*k - 1]
            m *= k
        }
        usable := m - m%n          // largest multiple of n that fits in [0, m-1]
        if x < usable {
            return x%n + 1         // uniform in [1, n]
        }
        // x in [usable, m-1] — reject and rebuild.
    }
}
```

### Geometric: uniform point in a circle — sample the square, reject the corners

```go
type Solution struct {
    radius, xCenter, yCenter float64
}

// randPoint returns a uniform random point [x, y] inside (or on) the circle.
func (s *Solution) randPoint() []float64 {
    for {
        // Propose uniformly in the bounding square [-r, r] x [-r, r].
        x := (rand.Float64()*2 - 1) * s.radius // uniform in [-r, r]
        y := (rand.Float64()*2 - 1) * s.radius // uniform in [-r, r]

        // Accept iff inside the disk. Compare squared distances — no sqrt needed.
        if x*x+y*y <= s.radius*s.radius {
            return []float64{s.xCenter + x, s.yCenter + y}
        }
        // Outside the circle (one of the four corners) — reject and re-draw.
    }
}
```

> Note the two small but important details: compare `x*x + y*y <= r*r` (avoid a
> `sqrt`), and `<=` (not `<`) so the boundary is included and the disk is
> sampled uniformly by *area*.

---

## Worked example — `rand10()` from `rand7()`, traced

Universe: two `rand7()` calls → `idx = (row)*7 + col`, `row, col ∈ {0..6}`, so
`idx ∈ {0..48}`, each of the 49 cells equally likely (probability `1/49`).

| Draw | `row` | `col` | `idx = row*7+col` | In `[0,39]`? | Result |
|------|-------|-------|-------------------|--------------|--------|
| 1 | 6 | 5 | 47 | no (≥ 40) | **reject**, redraw |
| 2 | 0 | 3 | 3 | yes | `3 % 10 + 1 = 4` |

Why `4` is uniform: of the 49 equally likely cells, we keep the 40 with
`idx ∈ {0..39}`. Exactly **4** of those (`3, 13, 23, 33`) satisfy
`idx % 10 + 1 == 4`, and the same count — 4 — maps to every one of the ten
outputs `1..10`. So each output has probability `4/40 = 1/10`. The 9 rejected
cells (`40..48`) carry away exactly the `9/49` of probability that couldn't be
split evenly into tenths; without discarding them, outputs `1..9` would each
have gotten one extra cell and been over-represented.

Expected `rand7()` calls: each loop iteration makes 2 calls and is accepted with
probability `40/49`. Expected iterations `= 49/40 = 1.225`, so expected calls
`= 2 × 1.225 = 2.45`.

---

## Complexity

| Quantity | Value | Reason |
|----------|-------|--------|
| Time (expected) | O(1) | geometric # of iterations with mean `1/p`, each O(1) |
| Time (worst case) | unbounded | no hard cap on retries; `P(>k)=(1−p)^k → 0` |
| Space | O(1) | a couple of scalars; no allocation in the loop |
| `rand10`/`rand7` draws | `2 · 49/40 ≈ 2.45` | 2 calls per grid draw × `1/p` iterations |
| disk point / `rand.Float64` draws | `2 · 4/π ≈ 2.55` | 2 coords per proposal × `1/p = 4/π` iterations |

The technique trades a small, bounded-in-expectation amount of extra work for
**exact** uniformity and **trivial** code (only a membership test, no CDF
inversion).

---

## Common pitfalls

1. **Taking `mod` without rejecting the tail.** `x % N` on a range whose size
   isn't a multiple of `N` biases the low residues. You **must** discard the
   `M mod N` leftover values first. This is the single most common bug.

2. **Wrong "usable" cutoff.** Keep `M − (M mod N)` values, i.e. indices
   `[0, usable)`. Off-by-one here (e.g. keeping 41 instead of 40 out of 49)
   silently reintroduces bias that no test with small samples will catch.

3. **Comparing with `sqrt` or the wrong inequality in the disk test.** Use
   `x*x + y*y <= r*r`; calling `math.Sqrt` is slower and can misclassify
   boundary points due to floating error. Use `<=` so the rim is included.

4. **Half-open vs. closed proposal interval.** `rand.Float64()` yields `[0, 1)`;
   `(rand.Float64()*2 - 1) * r` gives `[-r, r)`. That's fine for a continuous
   region (a single edge has measure zero), but don't reuse the same code for a
   *discrete* grid without re-checking which endpoints are included.

5. **Assuming a fixed number of iterations / adding a false cap.** Rejection
   sampling is expected-O(1), not bounded. Don't cap the loop and return a
   "fallback" value on the last iteration — that breaks uniformity. A plain
   `for {}` is correct; it terminates with probability 1.

6. **Re-seeding the RNG inside the loop (or per call).** Seed once at start-up.
   Re-seeding per draw (especially from a low-resolution clock) destroys
   independence and can freeze the output.

7. **Sampling radius uniformly when doing the *direct* alternative.** Not a
   rejection bug per se, but the classic trap in the same problem: `r = R·u` is
   **wrong** (over-samples the center); the area-correct transform is
   `r = R·√u`. Rejection sampling sidesteps this entirely — which is exactly why
   it's the easier-to-prove approach. (See [`/dsa/geometry.md`](/dsa/geometry.md).)

8. **Proposal region far larger than the target.** If `T` is a thin sliver of
   `S`, `1/p` is huge and you'll loop forever in practice. Tighten the proposal
   (a snugger bounding shape) or switch to a direct method.

---

## Rejection sampling vs. its cousins

All three are randomized-sampling staples; know which problem each solves.

| Technique | Problem it solves | Key idea | This repo |
|-----------|-------------------|----------|-----------|
| **Rejection sampling** (this file) | uniform sample over a hard-to-sample **set/region** | over-sample an easy superset, throw away out-of-target draws | #470, #478 |
| **Reservoir sampling** | pick `k` uniform items from a **stream of unknown length** in one pass, O(k) memory | keep a size-`k` reservoir; item `i` replaces a random slot with prob `k/i` | see [`/dsa/reservoir_sampling.md`](/dsa/reservoir_sampling.md) |
| **Fisher–Yates shuffle** | a uniform random **permutation** of a *known* array, in place | at step `i`, swap `i` with a random index in `[i, n)` | see [`/dsa/shuffle.md`](/dsa/shuffle.md) |

The distinction: **rejection** answers "give me a uniform *value* in this
region"; **reservoir** answers "give me uniform *samples* from a stream I can't
store"; **shuffle** answers "give me a uniform *ordering* of what I already
hold." They're not interchangeable, but they share the DNA of "make every
outcome equally likely, provably."

---

## Problems in this repo that use Rejection Sampling

- [0470 — Implement Rand10() Using Rand7()](/0470_implement_rand10_using_rand7/README.md)
  — the textbook integer case: two `rand7()` calls build a uniform `[0, 48]`,
  keep the largest multiple of 10 (`[0, 39]`), reject the 9-value tail, and
  reduce mod 10.
- [0478 — Generate Random Point in a Circle](/0478_generate_random_point_in_a_circle/README.md)
  — the textbook geometric case: sample the bounding square, accept iff
  `x² + y² ≤ r²`, reject the corners (≈ 1.27 proposals on average).

### Related classics to know (not yet in repo)

- LeetCode #528 — Random Pick with Weight (prefix sums + binary search; a *direct*
  method, contrast with rejection)
- LeetCode #710 — Random Pick with Blacklist (remap the blacklist, or reject
  blacklisted draws from `[0, N)`)
- LeetCode #398 — Random Pick Index (reservoir sampling — the streaming cousin)
- LeetCode #384 — Shuffle an Array (Fisher–Yates — the permutation cousin)
- Box–Muller / general Monte-Carlo integration (rejection sampling's continuous,
  non-uniform generalisation)
