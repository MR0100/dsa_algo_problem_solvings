# 0497 — Random Point in Non-overlapping Rectangles

> LeetCode #497 · Difficulty: Medium
> **Categories:** Reservoir Sampling, Math, Binary Search, Prefix Sum, Ordered Set, Randomized

---

## Problem Statement

You are given an array of non-overlapping axis-aligned rectangles `rects` where `rects[i] = [aᵢ, bᵢ, xᵢ, yᵢ]` indicates that `(aᵢ, bᵢ)` is the bottom-left corner point of the `iᵗʰ` rectangle and `(xᵢ, yᵢ)` is the top-right corner point of the `iᵗʰ` rectangle. Design an algorithm to pick a random integer point inside the space covered by one of the given rectangles. A point on the perimeter of a rectangle is included in the space covered by the rectangle.

Any integer point inside the space covered by one of the given rectangles should be equally likely to be returned.

**Note** that an integer point is a point that has integer coordinates.

Implement the `Solution` class:

- `Solution(int[][] rects)` Initializes the object with the given rectangles `rects`.
- `int[] pick()` Returns a random integer point `[u, v]` inside the space covered by one of the given rectangles.

**Example 1:**

```
Input
["Solution", "pick", "pick", "pick", "pick", "pick"]
[[[[-2, -2, 1, 1], [2, 2, 4, 6]]], [], [], [], [], []]
Output
[null, [1, -2], [1, -1], [-1, -2], [-2, -2], [0, 0]]

Explanation
Solution solution = new Solution([[-2, -2, 1, 1], [2, 2, 4, 6]]);
solution.pick(); // return [1, -2]
solution.pick(); // return [1, -1]
solution.pick(); // return [-1, -2]
solution.pick(); // return [-2, -2]
solution.pick(); // return [0, 0]
```

**Constraints:**

- `1 <= rects.length <= 100`
- `rects[i].length == 4`
- `-10^9 <= aᵢ < xᵢ <= 10^9`
- `-10^9 <= bᵢ < yᵢ <= 10^9`
- `xᵢ - aᵢ <= 2000`
- `yᵢ - bᵢ <= 2000`
- All the rectangles do not overlap.
- At most `10^4` calls will be made to `pick`.

> Because `pick()` is random, the concrete coordinates above are just *one*
> valid outcome. Any point inside the union, chosen uniformly, is correct.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Weighted Random / Reservoir Sampling** — the union is picked *by weight* (integer-point count per rectangle); weighted reservoir sampling (A-Res with k=1) selects a rectangle proportionally in one streaming pass with O(1) memory → see [`/dsa/reservoir_sampling.md`](/dsa/reservoir_sampling.md)
- **Prefix Sum** — cumulative sums of the per-rectangle point counts lay the rectangles end-to-end on a number line so a single random id maps to an owning rectangle → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Binary Search** — given the prefix table, the owning rectangle for a random id is found by binary-searching the first cumulative count that exceeds it, in O(log k) → see [`/dsa/binary_search.md`](/dsa/binary_search.md)

---

## Approaches Overview

| # | Approach | Construct | Pick | Space | When to use |
|---|----------|-----------|------|-------|-------------|
| 1 | Prefix-Sum of Weights + Binary Search (Optimal) | O(k) | O(log k) | O(k) | Many `pick()` calls; the standard answer |
| 2 | Weighted Reservoir Sampling | O(1) | O(k) | O(1) | Streaming / when rectangles arrive online or memory is tight |

*(k = number of rectangles)*

---

## Approach 1 — Prefix-Sum of Weights + Binary Search (Optimal)

### Intuition

Uniform over **points**, not over rectangles: a big rectangle must be chosen more often. Count the integer points each rectangle holds — for an inclusive rectangle `[x1,y1,x2,y2]` that is `(x2−x1+1)·(y2−y1+1)` — and imagine giving every point a global id `0 … total−1`, with rectangle 0 owning the first block of ids, rectangle 1 the next block, and so on. Draw a uniform id in `[0, total)` and find which block (rectangle) it lands in; that picks each rectangle with probability `weight/total`, i.e. uniformly over points. Prefix sums of the weights make "which block" an O(log k) binary search. Finally pick a uniform lattice point inside the chosen rectangle.

### Algorithm

1. **Construct:** build `prefix` with `prefix[0]=0` and `prefix[i+1]=prefix[i]+weightᵢ`; `total = prefix[k]`.
2. **Pick:**
   1. Draw `target = rand.Intn(total)`.
   2. Binary-search the least `i` with `prefix[i+1] > target` → rectangle `i` owns `target`.
   3. Return `[x1 + rand.Intn(x2−x1+1), y1 + rand.Intn(y2−y1+1)]` for rectangle `i`.

### Complexity

- **Time:** O(k) once to build the prefix table; O(log k) per `pick()` for the binary search (plus O(1) to draw the point).
- **Space:** O(k) for the prefix array.

### Code

```go
type SolutionBinary struct {
	rects  [][]int // original rectangles
	prefix []int   // prefix[i] = total integer points in rects[0..i-1]
	total  int     // total integer points across all rectangles
}

func ConstructorBinary(rects [][]int) SolutionBinary {
	prefix := make([]int, len(rects)+1) // prefix[0]=0; prefix[k]=total
	for i, r := range rects {
		// integer points in an inclusive rectangle = (width+1)*(height+1)
		w := (r[2] - r[0] + 1) * (r[3] - r[1] + 1)
		prefix[i+1] = prefix[i] + w // running cumulative count
	}
	return SolutionBinary{rects: rects, prefix: prefix, total: prefix[len(rects)]}
}

func (s *SolutionBinary) Pick() []int {
	target := rand.Intn(s.total) // a uniform "point id" in [0, total)
	// Find the first rectangle whose cumulative count strictly exceeds target;
	// sort.Search returns the least i in [1,k] with prefix[i] > target.
	i := sort.Search(len(s.rects), func(i int) bool {
		return s.prefix[i+1] > target
	})
	r := s.rects[i] // the chosen rectangle
	// draw x uniformly in [x1, x2] and y uniformly in [y1, y2] (inclusive)
	x := r[0] + rand.Intn(r[2]-r[0]+1)
	y := r[1] + rand.Intn(r[3]-r[1]+1)
	return []int{x, y}
}
```

### Dry Run

Rectangles `[[1,1,5,5], [-2,-2,-1,-1]]`. Weights: rect0 = `(5−1+1)·(5−1+1)=25`, rect1 = `(−1−(−2)+1)·(−1−(−2)+1)=2·2=4`.

Construct `prefix`:

| i | rectangle | weight | prefix[i+1] |
|---|-----------|--------|-------------|
| 0 | `[1,1,5,5]` | 25 | 25 |
| 1 | `[-2,-2,-1,-1]` | 4 | 29 |

`total = 29`. Suppose `pick()` draws `target = 26`:

| condition tested | prefix[i+1] > 26? | verdict |
|------------------|-------------------|---------|
| i=0 → prefix[1]=25 | 25 > 26? no | keep searching right |
| i=1 → prefix[2]=29 | 29 > 26? yes | owner = rectangle 1 |

Rectangle 1 = `[-2,-2,-1,-1]`; draw `x ∈ {-2,-1}`, `y ∈ {-2,-1}`, e.g. `[-1,-2]`. A point inside the union ✔. Over many draws, rect0 is chosen `25/29 ≈ 0.862` of the time — matching the program's measured share.

---

## Approach 2 — Weighted Reservoir Sampling

### Intuition

Skip the prefix array entirely and select a rectangle in a single pass. Track the running total of points seen so far; when rectangle `i` (weight `wᵢ`) arrives, adopt it as the current choice with probability `wᵢ / running`. By the standard reservoir induction, after the last rectangle the held one is rectangle `i` with probability `wᵢ / total` — the identical uniform-over-points distribution, achieved with O(1) memory and no precomputation. Then draw a uniform point inside the survivor.

Why the induction holds (k=1, weighted): after processing prefix of total weight `Wₘ`, each seen rectangle is held with prob `wᵢ/Wₘ`. Adding rectangle `m+1` (weight `w`): it is adopted with prob `w/Wₘ₊₁`; an earlier `i` survives with prob `(wᵢ/Wₘ)·(1 − w/Wₘ₊₁) = (wᵢ/Wₘ)·(Wₘ/Wₘ₊₁) = wᵢ/Wₘ₊₁`. Invariant preserved.

### Algorithm

1. `running = 0`, `chosen = -1`.
2. For each rectangle `i` with weight `wᵢ`:
   1. `running += wᵢ`.
   2. With probability `wᵢ/running` (i.e. `rand.Intn(running) < wᵢ`) set `chosen = i`.
3. Return a uniform integer point inside rectangle `chosen`.

### Complexity

- **Time:** O(k) per `pick()` — every rectangle is streamed on each call.
- **Space:** O(1) beyond the stored rectangles (no prefix table).

### Code

```go
type SolutionReservoir struct {
	rects [][]int // original rectangles
}

func ConstructorReservoir(rects [][]int) SolutionReservoir {
	return SolutionReservoir{rects: rects}
}

func (s *SolutionReservoir) Pick() []int {
	running := 0  // total integer points seen so far in the stream
	chosen := -1  // index of the currently-held rectangle
	for i, r := range s.rects {
		w := (r[2] - r[0] + 1) * (r[3] - r[1] + 1) // this rectangle's point count
		running += w                               // extend the seen-so-far total
		// keep rectangle i with probability w/running (rand.Intn(running) < w)
		if rand.Intn(running) < w {
			chosen = i
		}
	}
	r := s.rects[chosen] // the survivor of the reservoir
	x := r[0] + rand.Intn(r[2]-r[0]+1) // uniform x in [x1, x2]
	y := r[1] + rand.Intn(r[3]-r[1]+1) // uniform y in [y1, y2]
	return []int{x, y}
}
```

### Dry Run

Rectangles `[[1,1,5,5], [-2,-2,-1,-1]]`, weights 25 and 4. One `pick()`:

| i | wᵢ | running after += | adopt prob wᵢ/running | example draw < wᵢ? | chosen |
|---|----|------------------|-----------------------|--------------------|--------|
| 0 | 25 | 25 | 25/25 = 1.00 | always (Intn(25)<25) | 0 |
| 1 | 4 | 29 | 4/29 ≈ 0.138 | e.g. Intn(29)=17 → 17<4? no | 0 |

Survivor = rectangle 0 = `[1,1,5,5]`; draw `x ∈ [1,5]`, `y ∈ [1,5]`, e.g. `[3,3]` — inside the union ✔. Rectangle 0 ends up chosen `25/29 ≈ 0.862` of the time across calls, exactly as measured.

---

## Key Takeaways

- **"Uniform over the union" = weighted-by-size selection.** Convert each region to a weight (here, integer-point count `(dx+1)(dy+1)` for inclusive corners), pick a region proportional to weight, then pick uniformly inside it.
- **Two interchangeable weighting engines:** prefix-sum + binary search (O(log k) pick, O(k) memory) vs. weighted reservoir sampling (O(k) pick, O(1) memory). Choose by call-count vs. memory/streaming constraints.
- **`rand.Intn(running) < w`** is the clean integer test for "adopt with probability `w/running`", avoiding floating point.
- **Inclusive-corner counting** is the classic off-by-one trap: a rectangle from `x1` to `x2` spans `x2−x1+1` integer x-values, not `x2−x1`.
- Randomised outputs are verified by **invariants** (every point lands in the union; empirical frequencies match the weights), never by fixed expected coordinates.

---

## Related Problems

- LeetCode #528 — Random Pick with Weight (prefix-sum + binary search, the 1-D core)
- LeetCode #398 — Random Pick Index (unweighted reservoir sampling, k=1)
- LeetCode #382 — Linked List Random Node (reservoir sampling over a stream)
- LeetCode #470 — Implement Rand10() Using Rand7() (uniform sampling construction)
- LeetCode #710 — Random Pick with Blacklist (uniform pick over a remapped domain)
