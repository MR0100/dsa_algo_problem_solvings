# Stack

> **Pattern family:** Linear data structures · LIFO processing
> **Close relatives:** [Monotonic Stack](#3-monotonic-stack), Recursion / call stack, Queue (FIFO counterpart)

---

## 1. What a Stack Is

A **stack** is a linear collection with one rule: **Last In, First Out (LIFO)**.
You may only touch the *top* element. Three core operations, all **O(1)**:

| Operation | Meaning                              | Go (slice-backed)            |
|-----------|--------------------------------------|------------------------------|
| `Push(x)` | put `x` on top                       | `st = append(st, x)`         |
| `Pop()`   | remove & return the top element      | `x := st[len(st)-1]; st = st[:len(st)-1]` |
| `Peek()`  | read the top element without removal | `x := st[len(st)-1]`         |
| `Empty()` | is the stack empty?                  | `len(st) == 0`               |

Think of a stack of plates: the plate you placed last is the first one you can lift.

In Go there is no `Stack` type in the standard library — a **slice** is the
idiomatic stack. It gives amortised O(1) push/pop, cache-friendly memory, and
zero dependencies.

### Why LIFO matters

LIFO is exactly the order in which **nested or deferred things resolve**:

- The most recently opened bracket must be the first one closed.
- The most recent function call is the first one to return (the *call stack*).
- The most recent unresolved item ("a bar waiting for a taller bar") is the
  first one a new element can resolve.

Whenever a problem has this "most-recent-first" resolution structure, a stack
is the natural tool.

---

## 2. How to Recognise a Stack Problem

Signals in the problem statement:

1. **Matching / nesting / balancing** — "valid parentheses", "balanced
   brackets", nested tags, nested directories (`..` cancels the *last*
   directory ⇒ Simplify Path).
2. **Undo / cancel the most recent thing** — backspace in a string,
   `..` in a path, collision problems (Asteroid Collision) where the newest
   element may destroy earlier ones.
3. **"Nearest greater/smaller element to the left/right"** — for each element
   you need the closest previous/next element satisfying a comparison ⇒
   **monotonic stack** (Largest Rectangle in Histogram, Trapping Rain Water,
   Daily Temperatures, Next Greater Element).
4. **Parsing expressions** — evaluate/decode strings with operators,
   precedence, or nested structure (Basic Calculator, Decode String, RPN).
5. **Converting recursion to iteration** — any DFS / tree traversal can be
   made iterative with an explicit stack (Binary Tree Inorder Traversal,
   Flatten Binary Tree). If the recursive solution risks stack overflow or
   the interviewer asks "now do it without recursion", reach for a stack.
6. **Processing history you may need to revisit in reverse order** — browser
   back button, min-stack style questions ("design a stack that also …").

Rule of thumb: if while scanning left→right the *current* element interacts
only with the *most recent unresolved* earlier elements, it is a stack problem.

---

## 3. Templates (Go)

### 3.1 Basic slice-backed stack

```go
// A slice is the idiomatic Go stack. For type safety wrap it or use generics.
stack := []int{}                     // empty stack

stack = append(stack, 42)            // Push(42)

top := stack[len(stack)-1]           // Peek — MUST check len(stack) > 0 first!

stack = stack[:len(stack)-1]         // Pop (discard top)

x := stack[len(stack)-1]             // Pop-and-use idiom:
stack = stack[:len(stack)-1]         //   read top, then shrink
_ = x

if len(stack) == 0 { /* empty */ }   // Empty check
```

### 3.2 Matching / balancing template (Valid Parentheses shape)

```go
// isValid reports whether every opener is closed in the correct LIFO order.
//
// Pseudocode:
//   for each char c:
//     if c is an opener  -> push it
//     if c is a closer   -> stack must be non-empty AND top must be the
//                           matching opener; pop it, else fail
//   at the end the stack must be empty (no unclosed openers)
func isValid(s string) bool {
	pairs := map[byte]byte{')': '(', ']': '[', '}': '{'} // closer -> opener
	stack := []byte{}

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '(' || c == '[' || c == '{' {
			stack = append(stack, c) // opener: defer until its closer arrives
			continue
		}
		// closer: the MOST RECENT unclosed opener must match (LIFO)
		if len(stack) == 0 || stack[len(stack)-1] != pairs[c] {
			return false // nothing to close, or wrong nesting order
		}
		stack = stack[:len(stack)-1] // matched — resolve the opener
	}
	return len(stack) == 0 // leftover openers were never closed
}
```

### 3.3 Monotonic stack template (nearest smaller/greater element)

A **monotonic stack** keeps its elements sorted (by value) from bottom to top.
Pushing a new element first pops every element that breaks the order — and the
moment an element is popped, the new element is its answer ("first
greater/smaller to the right").

```go
// nextGreater[i] = index of the first element to the RIGHT of i that is
// strictly greater than nums[i], or -1 if none exists.
//
// Invariant: the stack holds INDICES whose values are strictly DECREASING
// from bottom to top — i.e. indices still waiting for a greater element.
//
// Every index is pushed once and popped at most once -> O(n) total.
func nextGreaterIndices(nums []int) []int {
	n := len(nums)
	res := make([]int, n)
	for i := range res {
		res[i] = -1 // default: no greater element to the right
	}
	stack := []int{} // indices with decreasing values, all "unresolved"

	for i := 0; i < n; i++ {
		// nums[i] resolves every waiting index with a smaller value
		for len(stack) > 0 && nums[stack[len(stack)-1]] < nums[i] {
			j := stack[len(stack)-1]     // most recent unresolved index
			stack = stack[:len(stack)-1] // pop it — it's resolved now
			res[j] = i                   // nums[i] is its next greater
		}
		stack = append(stack, i) // i now waits for ITS next greater
	}
	return res
}
```

Variants (memorise the mapping):

| Question                          | Stack order (bottom→top) | Pop while                |
|-----------------------------------|--------------------------|--------------------------|
| next **greater** to the right     | decreasing               | `nums[top] <  nums[i]`   |
| next **greater or equal**         | strictly decreasing pops | `nums[top] <= nums[i]`   |
| next **smaller** to the right     | increasing               | `nums[top] >  nums[i]`   |
| previous greater / smaller        | same, but the answer is what remains on the stack *below* after popping |

### 3.4 Iterative DFS / tree traversal template

```go
// inorder traversal without recursion: the explicit stack replaces the
// call stack — push the path of left children, then visit, then go right.
func inorder(root *TreeNode) []int {
	res := []int{}
	stack := []*TreeNode{}
	curr := root

	for curr != nil || len(stack) > 0 {
		for curr != nil { // dive left, deferring each node
			stack = append(stack, curr)
			curr = curr.Left
		}
		curr = stack[len(stack)-1] // deepest deferred node
		stack = stack[:len(stack)-1]
		res = append(res, curr.Val) // visit (inorder position)
		curr = curr.Right           // then explore its right subtree
	}
	return res
}
```

---

## 4. Worked Example — Valid Parentheses (LeetCode #20)

Input: `s = "([{}])"` → expected `true`.

Using the matching template (section 3.2), trace character by character:

| Step | Char | Action                                        | Stack (bottom→top) |
|------|------|-----------------------------------------------|--------------------|
| 1    | `(`  | opener → push                                 | `(`                |
| 2    | `[`  | opener → push                                 | `( [`              |
| 3    | `{`  | opener → push                                 | `( [ {`            |
| 4    | `}`  | closer; top `{` == pairs[`}`] → pop           | `( [`              |
| 5    | `]`  | closer; top `[` == pairs[`]`] → pop           | `(`                |
| 6    | `)`  | closer; top `(` == pairs[`)`] → pop           | *(empty)*          |
| end  | —    | stack empty → **return `true`**               |                    |

Counter-example `s = "(]"`:

| Step | Char | Action                                              | Stack |
|------|------|-----------------------------------------------------|-------|
| 1    | `(`  | opener → push                                       | `(`   |
| 2    | `]`  | closer; top `(` ≠ pairs[`]`] = `[` → **return `false`** | `(` |

The LIFO order is doing all the work: step 4 *must* close the most recently
opened bracket, and the stack top is exactly that bracket.

**Complexity:** Time O(n) — each char pushed/popped at most once.
Space O(n) — worst case all openers, e.g. `"((((("`.

---

## 5. Common Pitfalls (and fixes)

1. **Popping / peeking an empty stack.**
   `stack[len(stack)-1]` on an empty slice panics with index out of range.
   *Fix:* always guard with `len(stack) > 0` — in matching problems an early
   closer means the input is invalid, not that you should skip the check.

2. **Forgetting the final emptiness check.**
   `"((("` passes every per-character check but is invalid.
   *Fix:* the last line of a matching solution is `return len(stack) == 0`.

3. **Wrong pop condition in monotonic stacks (`<` vs `<=`).**
   Strict vs non-strict comparison decides how duplicates are handled and can
   silently double-count (classic bug in Largest Rectangle in Histogram).
   *Fix:* decide explicitly whether equal elements resolve each other; test
   with an input containing duplicates, e.g. `[2,2,2]`.

4. **Storing values when you need indices.**
   Most monotonic-stack answers need distances or positions (`i - j`), so
   push **indices** and look values up via `nums[idx]`.

5. **Forgetting to drain the stack after the scan.**
   In histogram-style problems, bars left on the stack still need processing.
   *Fix:* either loop `for i <= n` with a virtual 0-height sentinel bar at
   `i == n`, or add an explicit drain loop after the main pass.

6. **Believing the nested loop makes it O(n²).**
   The inner `for` of a monotonic stack pops each element at most once over
   the *whole* run ⇒ amortised O(n). Say this in interviews.

7. **Aliasing bugs when copying slice-backed stacks.**
   `st2 := st` shares the backing array; a later `append` on one may or may
   not affect the other. *Fix:* copy explicitly
   (`st2 := append([]int(nil), st...)`) if you need an independent snapshot.

8. **Using `container/list` for a stack in Go.**
   It works but is slower (pointer chasing, allocations) and clunkier than a
   slice. Interviewers expect the slice idiom.

---

## 6. Problems in This Repo

| Problem | Stack usage |
|---------|-------------|
| [0020 — Valid Parentheses](../0020_valid_parentheses/README.md) | The canonical matching/balancing stack |
| [0032 — Longest Valid Parentheses](../0032_longest_valid_parentheses/README.md) | Stack of indices to measure valid spans |
| [0042 — Trapping Rain Water](../0042_trapping_rain_water/README.md) | Monotonic decreasing stack fills water layer by layer |
| [0071 — Simplify Path](../0071_simplify_path/README.md) | Stack of directory names; `..` pops the most recent |
| [0084 — Largest Rectangle in Histogram](../0084_largest_rectangle_in_histogram/README.md) | Monotonic increasing stack finds nearest smaller bars |
| [0085 — Maximal Rectangle](../0085_maximal_rectangle/README.md) | Reduces each row to Largest Rectangle in Histogram |
| [0094 — Binary Tree Inorder Traversal](../0094_binary_tree_inorder_traversal/README.md) | Explicit stack replaces recursion for iterative DFS |
| [0114 — Flatten Binary Tree to Linked List](../0114_flatten_binary_tree_to_linked_list/README.md) | Iterative preorder flattening with an explicit stack |

<!-- Problems 0131+ to be added in a later pass. -->
