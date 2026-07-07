# Monotonic Stack

> A stack whose elements are kept in sorted (monotonic) order — either
> non-decreasing or non-increasing — from bottom to top. Elements that would
> break the invariant are popped **before** the new element is pushed, and the
> pop moment is exactly when we learn something useful about the popped element.

---

## What it is

A monotonic stack is an ordinary stack (usually of **indices**, not values) with
one extra rule enforced on every push:

- **Monotonic increasing stack** — values from bottom to top are increasing
  (or non-decreasing). Before pushing `x`, pop every element `≥ x` (or `> x`).
- **Monotonic decreasing stack** — values from bottom to top are decreasing
  (or non-increasing). Before pushing `x`, pop every element `≤ x` (or `< x`).

The magic is in the **pop event**. When element `s = stack.top` is popped
because the incoming element `x` at index `i` violates the invariant:

- `x` is the **first element to the right** of `s` that beats it
  (first smaller / first greater, depending on stack direction), and
- the element now underneath `s` on the stack is the **nearest element to the
  left** of `s` that beats it.

So in a single left-to-right pass, every element learns both its *next* and
*previous* smaller/greater neighbour — in **O(n)** total, because each index is
pushed exactly once and popped at most once.

| Stack keeps values… | Pop condition (strict) | A pop answers…                          |
|---------------------|------------------------|-----------------------------------------|
| Increasing          | `x < top`              | **next smaller** to the right of `top`  |
| Increasing          | `x <= top`             | next smaller-or-equal to the right      |
| Decreasing          | `x > top`              | **next greater** to the right of `top`  |
| Decreasing          | `x >= top`             | next greater-or-equal to the right      |

Mnemonic: the stack is monotonic in the *opposite* direction of what you are
looking for. Hunting "next **greater**" → keep a **decreasing** stack (a
greater element is exactly what breaks it). Hunting "next **smaller**" → keep
an **increasing** stack.

---

## How to recognise it — signals in the problem statement

Reach for a monotonic stack when you see:

1. **"Next / previous greater / smaller element"** — stated literally
   (LC #496, #503) or in disguise: "how many days until a warmer temperature"
   (LC #739), "span of stock prices" (LC #901).
2. **Bars / heights / histograms / skylines** — "largest rectangle in
   histogram" (LC #84), "maximal rectangle" (LC #85), "trapping rain water"
   (LC #42), "buildings with an ocean view". Anything drawn as vertical bars
   where a bar's reach is limited by the first shorter/taller bar on each side.
3. **"For each element, find the range where it is the min/max"** — sum of
   subarray minimums (LC #907), subarray ranges. The monotonic stack finds, for
   each element, the maximal window in which it dominates.
4. **Greedy removal to keep the best sequence** — "remove k digits to make the
   smallest number" (LC #402), "remove duplicate letters, smallest
   lexicographic result" (LC #316), "create maximum number" (LC #321). You keep
   a stack and pop earlier elements when a better one arrives and you can still
   afford to drop them.
5. **A brute force that repeatedly rescans left/right** for the nearest
   dominating element — the O(n²) inner loop of "walk left until you find
   something taller" is exactly what a monotonic stack amortises to O(n).
6. **"Visible" elements** — who can see whom over intermediate obstacles
   (visible people in a queue, LC #1944).

Counter-signal: if the property is about *sums* in a window, you likely want a
sliding window or prefix sums; if it is a *sliding-window min/max*, you want the
cousin structure, a **monotonic deque** (LC #239).

---

## General templates (Go)

### Template 1 — Next smaller element to the right (increasing stack)

```go
// nextSmaller[i] = index of the first element to the RIGHT of i with a
// strictly smaller value; n if none exists.
func nextSmallerRight(nums []int) []int {
    n := len(nums)
    res := make([]int, n)
    for i := range res {
        res[i] = n // default: no smaller element to the right
    }
    stack := []int{} // holds indices; nums[stack] is increasing bottom→top

    for i := 0; i < n; i++ {
        // Incoming nums[i] breaks the invariant for every top that is bigger.
        // The pop moment IS the answer moment: nums[i] is the first value to
        // the right of stack-top that is smaller than it.
        for len(stack) > 0 && nums[i] < nums[stack[len(stack)-1]] {
            top := stack[len(stack)-1]
            stack = stack[:len(stack)-1] // pop
            res[top] = i                 // i is top's "next smaller"
        }
        stack = append(stack, i) // push; invariant restored
    }
    return res
}
```

Variants, all the same skeleton:

- **Next greater**: flip the comparison to `nums[i] > nums[stack.top]`
  (stack becomes decreasing).
- **Previous smaller/greater**: after the pop loop, the element below the new
  top (i.e. the current `stack.top` just before pushing `i`) is `i`'s nearest
  smaller (resp. greater) to the **left** — record it at push time.
- **Circular array** (LC #503): loop `i` from `0` to `2n-1`, index with `i % n`,
  only push while `i < n`.

### Template 2 — Answer-at-pop with both boundaries (histogram pattern)

```go
// For each bar, when it is popped we simultaneously know:
//   right boundary = current index i        (first bar shorter on the right)
//   left  boundary = new stack top + 1      (first bar shorter on the left)
// A sentinel pass with height 0 at i == n flushes the stack.
func largestRectangleArea(heights []int) int {
    best := 0
    stack := []int{} // indices; heights[stack] increasing bottom→top

    for i := 0; i <= len(heights); i++ {
        h := 0 // sentinel: a bar of height 0 pops everything at the end
        if i < len(heights) {
            h = heights[i]
        }
        for len(stack) > 0 && h < heights[stack[len(stack)-1]] {
            top := stack[len(stack)-1]
            stack = stack[:len(stack)-1]

            height := heights[top]
            left := -1 // if stack empty, top was the minimum of prefix [0..i)
            if len(stack) > 0 {
                left = stack[len(stack)-1]
            }
            width := i - left - 1 // bars in (left, i) are all >= height
            if area := height * width; area > best {
                best = area
            }
        }
        stack = append(stack, i)
    }
    return best
}
```

### Template 3 — Greedy "keep the best prefix" (build-a-sequence pattern)

```go
// removeKdigits-style: keep result monotonic while a budget allows popping.
// Pseudocode:
//   for each ch in input:
//       while stack not empty AND budget > 0 AND stack.top > ch:
//           pop; budget--          // dropping a bigger earlier digit helps
//       push ch (subject to length / "must keep" constraints)
//   trim / flush remaining budget from the tail
```

Here the stack **is the answer** being constructed, and the monotonic invariant
encodes optimality ("no larger digit should precede a smaller one while we can
still delete").

---

## Worked example — full trace

Problem: **next greater element to the right** for
`nums = [2, 1, 2, 4, 3]` (want `[4, 2, 4, -1, -1]`).
We keep a **decreasing** stack of indices; pop while `nums[i] > nums[top]`.

| Step | i | nums[i] | Pops (index:value → answer)         | Stack after push (idx:val) | res so far              |
|------|---|---------|-------------------------------------|-----------------------------|-------------------------|
| 1    | 0 | 2       | — (stack empty)                     | `[0:2]`                     | `[-1,-1,-1,-1,-1]`      |
| 2    | 1 | 1       | — (1 < 2, invariant holds)          | `[0:2, 1:1]`                | `[-1,-1,-1,-1,-1]`      |
| 3    | 2 | 2       | pop `1:1` → `res[1] = 2`            | `[0:2, 2:2]`                | `[-1, 2,-1,-1,-1]`      |
|      |   |         | (2 is NOT > 2 → stop; strict pops give "next strictly greater") | | |
| 4    | 3 | 4       | pop `2:2` → `res[2] = 4`; pop `0:2` → `res[0] = 4` | `[3:4]`      | `[ 4, 2, 4,-1,-1]`      |
| 5    | 4 | 3       | — (3 < 4)                           | `[3:4, 4:3]`                | `[ 4, 2, 4,-1,-1]`      |
| end  |   |         | indices left on stack (3, 4) have no greater element → stay -1 | | `[4, 2, 4, -1, -1]` ✓ |

Observe the two invariant payoffs at step 4: the *pop moment* pairs each popped
index with its next greater value (4), and each element was pushed/popped at
most once — 5 pushes, 3 pops, O(n) total.

For the *boundary-pair* flavour of the same mechanics (left and right boundary
learned at pop time), see the dry-run tables in
[0084 — Largest Rectangle in Histogram](../0084_largest_rectangle_in_histogram/README.md)
and [0042 — Trapping Rain Water](../0042_trapping_rain_water/README.md)
(Approach 4 traces water being added layer-by-layer at each pop).

---

## Common pitfalls and how to avoid them

1. **Storing values instead of indices.** Values are rarely enough — widths,
   distances ("how many days"), and boundary computations all need positions.
   Default to pushing indices; the value is always `nums[stack[k]]`.
2. **Strict vs non-strict pop condition.** `<` vs `<=` decides how duplicates
   are handled. For "next strictly greater" pop on `>`. For counting problems
   like sum of subarray minimums, one side must be strict and the other
   non-strict or duplicated subarrays get counted twice (or missed). Decide
   deliberately; test with an input containing equal adjacent values.
3. **Forgetting to flush the stack at the end.** Elements still on the stack
   never met their "breaker". Either post-process them (answer = -1 / n), or
   use a **sentinel**: iterate one index past the end with a virtual value
   (`0` for histogram, `-∞`/`+∞` generally) that pops everything — as in
   Template 2. Alternatively append the sentinel to the input, but remember it
   changes `len`.
4. **Wrong stack direction.** If your pops never fire (stack keeps growing) or
   fire on every element, you almost certainly inverted the comparison.
   Re-derive: *the thing I'm searching for is the thing that pops.*
5. **Empty-stack width bugs.** In boundary computations, when the stack is
   empty after a pop, the popped element was the minimum (or maximum) of the
   whole prefix — the left boundary is `-1`, giving `width = i - (-1) - 1 = i`.
   Off-by-one here is the #1 histogram bug; verify with a single-bar input.
6. **Circular arrays.** For "next greater in a circular array", don't rotate or
   copy — iterate `2n` times with `i % n` and only push during the first pass.
7. **Confusing it with a monotonic deque.** Sliding-window max/min needs
   removal from **both** ends (expire indices leaving the window from the
   front) — that's a deque, not a stack. Same invariant idea, different
   structure.
8. **Believing the nested loop makes it O(n²).** The inner `for` pops each
   index at most once over the entire run, so the amortised cost is O(n).
   State this in interviews — it's a standard follow-up question.

---

## When you spot it, say this in the interview

> "Each element needs the nearest element to its left/right that is
> smaller/greater. A brute force rescans and is O(n²). A monotonic stack gives
> every element that answer at the moment it is popped, each index is pushed
> and popped at most once, so the whole pass is O(n) time, O(n) space."

---

## Problems in this repo

| # | Problem | Role of the monotonic stack |
|---|---------|-----------------------------|
| 0042 | [Trapping Rain Water](../0042_trapping_rain_water/README.md) | Decreasing stack; each pop fills one horizontal water layer between the popped valley floor and its two taller walls |
| 0084 | [Largest Rectangle in Histogram](../0084_largest_rectangle_in_histogram/README.md) | Increasing stack of indices; a pop yields the bar's full extent (left & right first-shorter boundaries) → area in O(n) |
| 0085 | [Maximal Rectangle](../0085_maximal_rectangle/README.md) | Reduces each matrix row to a histogram and reuses the #84 monotonic-stack routine row by row |

> Problems #0131–#0400 are being added concurrently; a later pass will extend
> this list (expect: #150 note — plain stack; #316 Remove Duplicate Letters,
> #321 Create Maximum Number, #402 Remove K Digits when they land, plus later
> classics #496/#503 Next Greater Element, #739 Daily Temperatures, #901 Online
> Stock Span, #907 Sum of Subarray Minimums).

Related reference: for the plain LIFO structure and non-monotonic uses, see
[`/dsa/stack.md`](/dsa/stack.md) *(referenced by problem READMEs; created when
that concept file is written)*.
