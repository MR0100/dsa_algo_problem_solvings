# Reservoir Sampling

> **Category:** Randomized Algorithms / Streaming
> **Difficulty to master:** Medium — the algorithm is short, the proof is the interview.

---

## What it is

**Reservoir sampling** is a family of randomized algorithms for choosing a
uniform random sample of `k` items from a stream of `n` items, where:

- `n` is **unknown in advance** (or too large to store), and
- you may only make **one pass** over the data, and
- you may only use **O(k) extra memory** (the "reservoir").

Every item in the stream must end up in the sample with probability exactly
`k/n`, even though you never know `n` until the stream ends.

The classic single-item case (`k = 1`) is **Algorithm R** (Vitter): keep one
candidate; when you see the *i*-th item (1-indexed), replace the candidate
with it with probability `1/i`.

### Why it works (the one-line proof you must be able to say)

For item `i` to be the final answer it must be **picked at step i**
(probability `1/i`) and then **survive** every later step `j > i`
(probability `(j-1)/j` each). The product telescopes:

```
P(item i survives) = (1/i) · (i/(i+1)) · ((i+1)/(i+2)) · … · ((n-1)/n) = 1/n
```

Every item ends with probability `1/n` — uniform, regardless of `n`.

For general `k`: fill the reservoir with the first `k` items; for item
`i > k`, generate `j = rand(1..i)`; if `j ≤ k`, replace `reservoir[j]`.
Each item survives with probability `k/n` (same telescoping argument).

---

## How to recognise it — signals in the problem statement

Reach for reservoir sampling when you see any of these:

| Signal | Example phrasing |
|--------|------------------|
| **Pick a random element uniformly** | "Return a random node's value; each node must have the same probability of being chosen" |
| **Unknown / unbounded length** | "The linked list is extremely large and its length is unknown to you" |
| **Single pass / streaming constraint** | "Could you solve this efficiently without using extra space?" or "the input is a stream" |
| **Follow-up forbids O(n) memory** | "What if the array is too big to fit in memory?" |
| **Pick a random index among duplicates** | "Given a target, return a random index i such that nums[i] == target" (the matches form an implicit stream of unknown count) |
| **Sample k of n** | "Randomly select k distinct items/lines/records from a huge file" |

Classic LeetCode problems: **#382 Linked List Random Node**, **#398 Random
Pick Index**, **#519 Random Flip Matrix** (related), and the system-design
staple "pick k random lines from a huge log file."

**Anti-signal:** if the length is known and the data fits in memory, plain
`rand.Intn(n)` indexing (or Fisher–Yates shuffle for k samples) is simpler
and faster — reservoir sampling buys you nothing there. Weighted random pick
(#528) is **prefix sum + binary search**, not reservoir sampling.

---

## General templates (Go)

### Template 1 — pick 1 item from a stream (Algorithm R, k = 1)

```go
// reservoirSampleOne returns one uniformly random element from a stream.
//
// Time:  O(n) — one pass, O(1) work per item.
// Space: O(1) — a single candidate slot.
func reservoirSampleOne(stream []int) int {
	// result holds the current candidate ("the reservoir of size 1").
	var result int

	// i counts how many items we have seen so far (1-indexed).
	i := 0

	for _, x := range stream { // pseudocode: for each item x arriving...
		i++ // ...we have now seen i items

		// Replace the candidate with probability 1/i.
		// rand.Intn(i) is uniform over [0, i), so it equals 0 with prob 1/i.
		if rand.Intn(i) == 0 {
			result = x // item i is picked at step i
		}
		// Otherwise the old candidate survives step i with prob (i-1)/i.
	}
	return result
}
```

### Template 2 — pick k items from a stream (general Algorithm R)

```go
// reservoirSampleK returns k uniformly random elements (without replacement)
// from a stream of unknown length.
//
// Time:  O(n)
// Space: O(k)
func reservoirSampleK(stream []int, k int) []int {
	reservoir := make([]int, 0, k)

	for i, x := range stream { // i is 0-indexed position
		if i < k {
			// Phase 1: the first k items fill the reservoir unconditionally.
			reservoir = append(reservoir, x)
		} else {
			// Phase 2: item i+1 (1-indexed) should be kept with prob k/(i+1).
			// rand.Intn(i+1) is uniform over [0, i]; it lands in [0, k) with
			// probability exactly k/(i+1).
			j := rand.Intn(i + 1)
			if j < k {
				// Evict a uniformly random current occupant — reusing j as the
				// eviction slot keeps the distribution uniform.
				reservoir[j] = x
			}
		}
	}
	return reservoir
}
```

### Template 3 — random node from a linked list (LeetCode #382 shape)

```go
// getRandom returns a uniformly random node value from a singly linked list
// whose length is unknown, in one pass and O(1) space.
func (s *Solution) getRandom() int {
	result, i := 0, 0
	for node := s.head; node != nil; node = node.Next {
		i++                      // node index, 1-based
		if rand.Intn(i) == 0 {   // keep this node with probability 1/i
			result = node.Val
		}
	}
	return result
}
```

### Template 4 — random index of a target among duplicates (LeetCode #398 shape)

The positions where `nums[i] == target` form an implicit stream whose count
you don't want to precompute — sample from it on the fly:

```go
// pick returns a uniformly random index i with nums[i] == target,
// using O(1) extra space per query.
func (s *Solution) pick(target int) int {
	result, count := -1, 0
	for i, v := range s.nums {
		if v != target {
			continue // not part of the stream of matches
		}
		count++                      // this is the count-th match seen
		if rand.Intn(count) == 0 {   // keep it with probability 1/count
			result = i
		}
	}
	return result
}
```

---

## Worked example — full trace

Stream: `[10, 20, 30, 40]`, `k = 1`. We trace Template 1 and show the
survival probability of each candidate.

| Step i | Item | rand.Intn(i) range | P(replace) = 1/i | If replaced, result = | P(this item is FINAL answer) |
|--------|------|--------------------|------------------|----------------------|------------------------------|
| 1 | 10 | {0} | 1/1 = always | 10 | 1 · (1/2) · (2/3) · (3/4) = **1/4** |
| 2 | 20 | {0,1} | 1/2 | 20 | (1/2) · (2/3) · (3/4) = **1/4** |
| 3 | 30 | {0,1,2} | 1/3 | 30 | (1/3) · (3/4) = **1/4** |
| 4 | 40 | {0,1,2,3} | 1/4 | 40 | 1/4 = **1/4** |

Concrete run with dice rolls `rand.Intn` = `[0, 1, 0, 2]`:

1. `i=1`, item 10: roll 0 → 0 == 0, **replace**. `result = 10`.
2. `i=2`, item 20: roll 1 → 1 != 0, keep. `result = 10` (10 survived its 1/2 coin flip... precisely: survived with prob 1/2).
3. `i=3`, item 30: roll 0 → **replace**. `result = 30`.
4. `i=4`, item 40: roll 2 → keep. `result = 30`.

Final answer this run: `30`. Over many runs each of 10/20/30/40 appears
exactly 25% of the time — verify empirically by running the sampler ~100k
times and counting frequencies (a great `main()` demo for this repo's style).

---

## Common pitfalls and how to avoid them

1. **Off-by-one in the probability.** The *i*-th item (1-indexed) must be
   kept with probability `1/i`, i.e. `rand.Intn(i) == 0` — not
   `rand.Intn(i+1)` or `rand.Intn(i-1)`. If your counter is 0-indexed as
   `idx`, the test is `rand.Intn(idx+1) == 0`. Mixing conventions silently
   skews the distribution; always dry-run n = 2 by hand.

2. **Forgetting the first item is always taken.** At `i = 1`,
   `rand.Intn(1)` is always 0 — the first item must unconditionally seed
   the reservoir. If your loop structure can skip it, `result` may be
   garbage when the stream has one element.

3. **k-case: evicting the wrong slot.** In Template 2 you must evict a
   *uniformly random* occupant. Reusing `j = rand.Intn(i+1)` (when `j < k`)
   as the eviction index is the standard, correct trick — evicting always
   slot 0, or rolling a *second* random number conditioned on the first,
   changes the distribution unless done very carefully.

4. **Testing exact outputs.** Randomized output can't be asserted directly.
   Test the **distribution**: run 100k trials, check each item's frequency
   is within a few percent of `k/n` (chi-squared if you're fancy). In this
   repo's `main()`, print observed frequencies with the expected value as
   the inline comment.

5. **Reaching for it when n is known.** If you can index the data
   (`nums[rand.Intn(len(nums))]`), reservoir sampling is over-engineering.
   In interviews, state both: "with O(n) preprocessing I'd use a hash map of
   indices; with the follow-up's memory constraint I'd reservoir-sample."
   #398 is exactly this trade-off (hash map of target→indices vs O(1)-space
   sampling per query).

6. **Seeding / RNG misuse in Go.** Since Go 1.20, the global `math/rand` is
   auto-seeded — no `rand.Seed(time.Now().UnixNano())` needed (it's
   deprecated). For reproducible dry runs use `rand.New(rand.NewSource(42))`.
   Never use `math/rand` where crypto-grade randomness matters (not an issue
   for LeetCode).

7. **Confusing it with weighted sampling.** "Pick index with probability
   proportional to weight" (#528) is prefix sums + binary search. A weighted
   *reservoir* variant exists (Efraimidis–Spirakis: key = `u^(1/w)`, keep max
   keys), but plain Algorithm R is uniform-only.

---

## Interview follow-ups worth knowing

- **Distributed streams:** run reservoir sampling per shard, then merge —
  weight each shard's sample by its item count.
- **Algorithm L:** skips ahead geometrically instead of flipping a coin per
  item; expected `O(k(1 + log(n/k)))` random numbers instead of `O(n)`.
  Mention it if asked "can you call the RNG fewer times?"
- **Random point in a stream of unknown length appears in Google/Meta phone
  screens** as "pick a random line from a huge file" — same one-liner proof.

---

## Problems in this repo

_No problems in this repo (0001–0130 range) use reservoir sampling yet._
The canonical problems — LeetCode **#382 Linked List Random Node** and
**#398 Random Pick Index** — fall in the 0131–0400 batch currently being
written; link them here once their folders exist:

- `../0382_linked_list_random_node/README.md` *(pending)*
- `../0398_random_pick_index/README.md` *(pending)*
