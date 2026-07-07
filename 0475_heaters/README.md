# 0475 — Heaters

> LeetCode #475 · Difficulty: Medium
> **Categories:** Array, Two Pointers, Binary Search, Sorting

---

## Problem Statement

Winter is coming! During the contest, your first job is to design a standard heater with a fixed warm radius to warm all the houses.

Every house can be warmed, as long as the house is within the heater's warm radius range.

Given the positions of `houses` and `heaters` on a horizontal line, return *the minimum radius standard of heaters* so that those heaters could cover all houses.

**Notice** that all the `heaters` follow your radius standard, and the warm radius will be the same.

**Example 1:**

```
Input: houses = [1,2,3], heaters = [2]
Output: 1
Explanation: The only heater was placed in the position 2, and if we use the radius 1 standard, then all the houses can be warmed.
```

**Example 2:**

```
Input: houses = [1,2,3,4], heaters = [1,4]
Output: 1
Explanation: The two heaters were placed at positions 1 and 4. We need to use a radius 1 standard, then all the houses can be warmed.
```

**Example 3:**

```
Input: houses = [1,5], heaters = [2]
Output: 3
```

**Constraints:**

- `1 <= houses.length, heaters.length <= 3 * 10^4`
- `1 <= houses[i], heaters[i] <= 10^9`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search on a sorted array** — with heaters sorted, the nearest heater to a house is one of the two straddling it; `lower_bound` finds them in `O(log K)` → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Two Pointers / merge sweep** — sorting *both* arrays lets a heater pointer advance monotonically as houses move right, giving each house its nearest heater in amortised O(1) → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Sorting** — every efficient approach begins by sorting so distances become locally monotone → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(H·K) | O(1) | Tiny inputs / sanity check; ~9·10⁸ ops at the limit → TLE |
| 2 | Sort + Binary Search | O((H+K) log K) | O(1) | Clear and robust; only heaters need sorting |
| 3 | Sort Both + Two Pointers (Optimal) | O(H log H + K log K) | O(1) | Fewest comparisons; the tidy linear-sweep answer |

---

## Approach 1 — Brute Force (Each House Scans Every Heater)

### Intuition

One global radius must reach the **worst-off** house. Each house's requirement is the distance to its *nearest* heater (the min over heaters). The whole line is covered exactly when the radius is at least the **maximum** of those per-house minima. So the answer is `max_house( min_heater |house − heater| )`. The brute force computes that double loop literally.

### Algorithm

1. `ans = 0`.
2. For each house `h`: compute `best = min over heaters of |h − heater|`.
3. `ans = max(ans, best)`.
4. Return `ans`.

### Complexity

- **Time:** O(H·K) — each of `H` houses probes all `K` heaters; at `3·10⁴` each that is ~`9·10⁸` operations.
- **Space:** O(1) — two running values.

### Code

```go
func bruteForce(houses []int, heaters []int) int {
	ans := 0 // the largest "distance to nearest heater" seen so far
	for _, h := range houses {
		best := -1 // distance from this house to its closest heater
		for _, ht := range heaters {
			d := abs(h - ht) // distance house→heater
			if best == -1 || d < best {
				best = d // this heater is closer
			}
		}
		if best > ans {
			ans = best // this house needs the radius bumped up
		}
	}
	return ans
}
```

### Dry Run

Example 3: `houses = [1,5]`, `heaters = [2]`.

| House | distances to heaters | min (needed radius) | ans after |
|-------|----------------------|---------------------|-----------|
| 1 | `|1−2| = 1` | 1 | max(0,1) = 1 |
| 5 | `|5−2| = 3` | 3 | max(1,3) = **3** |

Result: `3` ✔ — the far house at 5 forces radius 3.

---

## Approach 2 — Sort + Binary Search

### Intuition

Sort the heaters. Now the nearest heater to a house is one of exactly two candidates: the largest heater `≤` house, or the smallest heater `≥` house. A binary search (`lower_bound`) lands right between them in `O(log K)`; the house's required radius is the smaller of the two gaps (handling the ends where only one neighbour exists). Take the max over houses.

### Algorithm

1. Sort `heaters`.
2. For each house `h`: `pos = lower_bound(heaters, h)` (first heater `≥ h`).
   - Right neighbour `heaters[pos]` if `pos < K` → gap `heaters[pos] − h`.
   - Left neighbour `heaters[pos-1]` if `pos > 0` → gap `h − heaters[pos-1]`.
   - `dist` = min of the existing gaps.
3. `ans = max(ans, dist)`.

### Complexity

- **Time:** O((H + K) log K) — one sort of the heaters, then a `log K` search per house.
- **Space:** O(1) — sort in place; a few scalars.

### Code

```go
func binarySearch(houses []int, heaters []int) int {
	sort.Ints(heaters) // sorted heaters enable straddle lookup
	k := len(heaters)
	ans := 0
	for _, h := range houses {
		// pos = first index whose heater value is >= h (lower bound).
		pos := sort.SearchInts(heaters, h)
		dist := 1 << 62 // +inf placeholder
		if pos < k {
			// Heater at pos is >= h: gap to the right neighbour.
			if d := heaters[pos] - h; d < dist {
				dist = d
			}
		}
		if pos > 0 {
			// Heater at pos-1 is < h: gap to the left neighbour.
			if d := h - heaters[pos-1]; d < dist {
				dist = d
			}
		}
		if dist > ans {
			ans = dist // worst house so far dictates the radius
		}
	}
	return ans
}
```

### Dry Run

Example 2: `houses = [1,2,3,4]`, `heaters = [1,4]` (already sorted).

| House h | pos = lower_bound | left nbr gap | right nbr gap | dist | ans after |
|---------|-------------------|--------------|---------------|------|-----------|
| 1 | 0 (heaters[0]=1 ≥ 1) | none (pos=0) | `1−1 = 0` | 0 | 0 |
| 2 | 1 (heaters[1]=4 ≥ 2) | `2−1 = 1` | `4−2 = 2` | 1 | 1 |
| 3 | 1 | `3−1 = 2` | `4−3 = 1` | 1 | 1 |
| 4 | 1 (heaters[1]=4 ≥ 4) | `4−1 = 3` | `4−4 = 0` | 0 | 1 |

Result: `1` ✔

---

## Approach 3 — Sort Both + Two-Pointer Sweep (Optimal)

### Intuition

Sort **both** arrays and merge-walk them. Process houses left → right with a heater index `j`. For the current house, keep advancing `j` while the *next* heater is no farther than the current one — i.e. while moving right still helps (or ties). The moment the next heater is strictly farther, `heaters[j]` is this house's nearest heater. Because houses only increase and the best heater for a larger house is never to the left of the best heater for a smaller one, `j` never moves backward — the whole sweep is linear after sorting.

### Algorithm

1. Sort `houses` and `heaters`.
2. `j = 0`, `ans = 0`.
3. For each house `h` in sorted order:
   - While `j+1 < K` and `|heaters[j+1] − h| <= |heaters[j] − h|`: `j++`.
   - `ans = max(ans, |heaters[j] − h|)`.
4. Return `ans`.

### Complexity

- **Time:** O(H log H + K log K) — the two sorts dominate; the pointer sweep is O(H + K) since `j` advances monotonically.
- **Space:** O(1) — in-place sorts and two indices.

### Code

```go
func twoPointers(houses []int, heaters []int) int {
	sort.Ints(houses)  // sweep houses in increasing order
	sort.Ints(heaters) // heaters aligned so the pointer only moves forward
	j := 0             // index of the current candidate heater
	ans := 0
	k := len(heaters)
	for _, h := range houses {
		// Advance while the next heater is no farther than the current one.
		// Once the next heater is strictly farther, heaters[j] is the closest.
		for j+1 < k && abs(heaters[j+1]-h) <= abs(heaters[j]-h) {
			j++
		}
		if d := abs(heaters[j] - h); d > ans {
			ans = d // this house needs at least distance d
		}
	}
	return ans
}
```

### Dry Run

Example 2: `houses = [1,2,3,4]`, `heaters = [1,4]` (both sorted). Start `j = 0`, `ans = 0`.

| House h | advance check: `|heaters[j+1]−h| ≤ |heaters[j]−h|`? | j after | dist `|heaters[j]−h|` | ans after |
|---------|------------------------------------------------------|---------|------------------------|-----------|
| 1 | `|4−1|=3 ≤ |1−1|=0`? no | 0 | `|1−1| = 0` | 0 |
| 2 | `|4−2|=2 ≤ |1−2|=1`? no | 0 | `|1−2| = 1` | 1 |
| 3 | `|4−3|=1 ≤ |1−3|=2`? yes → j=1; then j+1=2 not < K, stop | 1 | `|4−3| = 1` | 1 |
| 4 | j+1=2 not < K, stop | 1 | `|4−4| = 0` | 1 |

Result: `1` ✔ — the pointer advances from heater 1 to heater 4 exactly once, at house 3.

---

## Key Takeaways

- **"Minimum radius to cover all" = minimise a maximum.** The radius must satisfy the worst house, so compute each house's nearest-heater distance and take the max. This *max-of-mins* framing recurs across covering problems.
- **Sort to make distance locally monotone.** Whether you binary-search the straddling pair or two-pointer-sweep, sorting turns "nearest of K" into an O(log K) or amortised O(1) lookup.
- **Two sorted sequences ⇒ merge with a monotone pointer.** The heater index only moves forward, which is why the sweep beats per-house binary search on constant factors.
- Watch the boundaries: a house before the first heater or after the last has only one neighbour — the `<inf>` seed and the `pos > 0` / `pos < K` guards handle it.

---

## Related Problems

- LeetCode #35 — Search Insert Position (`lower_bound` on a sorted array)
- LeetCode #167 — Two Sum II (two-pointer sweep on sorted input)
- LeetCode #1102 — Path With Maximum Minimum Value (max-of-min objective)
- LeetCode #1011 — Capacity To Ship Packages Within D Days (binary search on the answer / covering)
