# 0331 — Verify Preorder Serialization of a Binary Tree

> LeetCode #331 · Difficulty: Medium
> **Categories:** Stack, Tree, String, Binary Tree

---

## Problem Statement

One way to serialize a binary tree is to use preorder traversal. When we encounter a non-null node, we record the node's value. If it is a null node, we record using a sentinel value such as `'#'`.

For example, the below binary tree can be serialized to the string `"9,3,4,#,#,1,#,#,2,#,6,#,#"`, where `'#'` represents a null node.

```
        _9_
       /   \
      3     2
     / \   / \
    4   1  #  6
   / \ / \   / \
   # # # #   # #
```

Given a string of comma-separated values `preorder`, return `true` if it is a correct preorder traversal serialization of a binary tree.

It is **guaranteed** that each comma-separated value in the string is either an integer or a character `'#'` representing null pointer.

You may assume that the input format is always valid.

- For example, it could never contain two consecutive commas, such as `"1,,3"`.

**Note:** You are not allowed to reconstruct the tree.

**Example 1:**

```
Input: preorder = "9,3,4,#,#,1,#,#,2,#,6,#,#"
Output: true
```

**Example 2:**

```
Input: preorder = "1,#"
Output: false
```

**Example 3:**

```
Input: preorder = "9,#,#,1"
Output: false
```

**Constraints:**

- `1 <= preorder.length <= 10^4`
- `preorder` consists of integers in the range `[0, 100]` and `'#'` separated by commas `','`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Stack** — the fold approach repeatedly collapses a completed `x,#,#` subtree pattern by popping three tokens and pushing one, exactly the reduce-on-a-stack idiom → see [`/dsa/stack.md`](/dsa/stack.md)
- **Tree Traversal** — the string *is* a preorder DFS; understanding how preorder lays out node/null slots is the whole insight → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **String Algorithms** — parsing a comma-separated token stream in a single pass → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Stack (Fold Complete Subtrees) | O(n) | O(n) | Very intuitive; mirrors how the tree is built bottom-up |
| 2 | Slot Counting (Optimal) | O(n) | O(1) extra | The slick answer; constant space, one arithmetic invariant |

---

## Approach 1 — Stack (Fold Complete Subtrees)

### Intuition

A valid preorder string is assembled from the recurring pattern **node, leftSubtree, rightSubtree**. The smallest complete unit is `number,#,#`: a node whose two children both resolved to null. From its parent's perspective, that finished subtree behaves *exactly* like a single null — it fills one child slot. So repeatedly find any `number,#,#` and replace it with `#`. Keep folding. A well-formed tree collapses all the way down to a single `#` — the one null slot the root occupies relative to its (imaginary) parent.

### Algorithm

1. Split `preorder` into tokens and process left→right, pushing each onto a stack.
2. After every push, while the top three tokens (in reading order) are `number,#,#`, pop those three and push a single `#`.
3. After all tokens are consumed, the string is valid iff the stack is exactly `["#"]`.

### Complexity

- **Time:** O(n) — each token is pushed once; each fold removes three and adds one, so total pushes/pops are linear.
- **Space:** O(n) — the stack can hold up to O(n) tokens before folds trigger.

### Code

```go
func stackFold(preorder string) bool {
	tokens := strings.Split(preorder, ",") // split the CSV serialization
	stack := []string{}                    // holds unresolved tokens
	for _, t := range tokens {
		stack = append(stack, t) // read one token
		// Collapse "number,#,#" (top-of-stack order: #, #, number) into "#".
		for len(stack) >= 3 &&
			stack[len(stack)-1] == "#" &&
			stack[len(stack)-2] == "#" &&
			stack[len(stack)-3] != "#" {
			stack = stack[:len(stack)-3] // pop the two nulls and the node
			stack = append(stack, "#")   // the completed subtree acts as a null
		}
	}
	// A valid full tree reduces to a single null placeholder.
	return len(stack) == 1 && stack[0] == "#"
}
```

### Dry Run

Example 1: `preorder = "9,3,4,#,#,1,#,#,2,#,6,#,#"`.

| Step | Token read | Stack after push | Fold applied? | Stack after fold |
|------|-----------|------------------|---------------|------------------|
| 1 | 9 | `[9]` | no | `[9]` |
| 2 | 3 | `[9,3]` | no | `[9,3]` |
| 3 | 4 | `[9,3,4]` | no | `[9,3,4]` |
| 4 | # | `[9,3,4,#]` | no | `[9,3,4,#]` |
| 5 | # | `[9,3,4,#,#]` | yes → `4,#,#`→`#` | `[9,3,#]` |
| 6 | 1 | `[9,3,#,1]` | no | `[9,3,#,1]` |
| 7 | # | `[9,3,#,1,#]` | no | `[9,3,#,1,#]` |
| 8 | # | `[9,3,#,1,#,#]` | yes → `1,#,#`→`#`, then `3,#,#`→`#` | `[9,#]` |
| 9 | 2 | `[9,#,2]` | no | `[9,#,2]` |
| 10 | # | `[9,#,2,#]` | no | `[9,#,2,#]` |
| 11 | 6 | `[9,#,2,#,6]` | no | `[9,#,2,#,6]` |
| 12 | # | `[9,#,2,#,6,#]` | no | `[9,#,2,#,6,#]` |
| 13 | # | `[9,#,2,#,6,#,#]` | `6,#,#`→`#`, then `2,#,#`→`#`, then `9,#,#`→`#` | `[#]` |

Final stack `["#"]`. Result: `true` ✔

---

## Approach 2 — Slot Counting (Optimal)

### Intuition

Model the tree as a supply of **open child slots**. Start with `slots = 1` (the root needs somewhere to hang — its incoming edge). Every token you read *consumes one slot*. A `#` (null) opens no new slots; a real node opens **two** (its left and right children). Two invariants nail correctness: you must never read a token when `slots == 0` (nowhere to attach it → invalid), and after the last token `slots` must be exactly `0` (no dangling unfilled edges). This never builds the tree — it just balances the slot budget.

### Algorithm

1. Set `slots = 1`.
2. For each token: if `slots == 0` before consuming, return `false`. Otherwise `slots--`; if the token is a number, `slots += 2`.
3. Return `slots == 0`.

### Complexity

- **Time:** O(n) — a single scan of the tokens.
- **Space:** O(1) extra — only the integer `slots` (O(n) if you materialize the split; the scan itself is constant space).

### Code

```go
func slotCounting(preorder string) bool {
	tokens := strings.Split(preorder, ",") // the sequence of nodes
	slots := 1                             // the root occupies one incoming slot
	for _, t := range tokens {
		if slots == 0 {
			return false // a token arrived with nowhere to attach
		}
		slots-- // this token fills one open slot
		if t != "#" {
			slots += 2 // a real node opens two child slots
		}
	}
	return slots == 0 // every slot must be exactly filled
}
```

### Dry Run

Example 1: `preorder = "9,3,4,#,#,1,#,#,2,#,6,#,#"`, start `slots = 1`.

| Step | Token | slots == 0? | slots-- | number? +2 | slots after |
|------|-------|-------------|---------|------------|-------------|
| 1 | 9 | no | 0 | +2 | 2 |
| 2 | 3 | no | 1 | +2 | 3 |
| 3 | 4 | no | 2 | +2 | 4 |
| 4 | # | no | 3 | — | 3 |
| 5 | # | no | 2 | — | 2 |
| 6 | 1 | no | 1 | +2 | 3 |
| 7 | # | no | 2 | — | 2 |
| 8 | # | no | 1 | — | 1 |
| 9 | 2 | no | 0 | +2 | 2 |
| 10 | # | no | 1 | — | 1 |
| 11 | 6 | no | 0 | +2 | 2 |
| 12 | # | no | 1 | — | 1 |
| 13 | # | no | 0 | — | 0 |

Ends with `slots == 0`. Result: `true` ✔ — (contrast Example 3 `"9,#,#,1"`: after `9,#,#` slots hits 0, but token `1` still remains → the `slots == 0` guard fires → `false`).

---

## Key Takeaways

- **Serialized trees have a conserved quantity.** In-degree vs out-degree of the implied graph must balance: each node contributes +2 out-edges and −1 in-edge; the whole string must net to a fixed value. That is exactly the `slots` invariant (`slots = 1` start, `slots = 0` end).
- **"Don't reconstruct the tree" is a hint to find the invariant.** When a problem forbids building the obvious structure, look for a counting/parity property that the structure would have satisfied.
- **Stack-fold = bottom-up subtree completion.** The `x,#,#` → `#` reduction is a reusable pattern for validating or evaluating any preorder/prefix expression (see also evaluating Polish notation on a stack).
- The slot-counting mid-stream check (`slots == 0` before consuming) is essential — without it `"9,#,#,1"` would wrongly balance to 0 only at the very end.

---

## Related Problems

- LeetCode #297 — Serialize and Deserialize Binary Tree (the reconstruction this problem forbids)
- LeetCode #106 — Construct Binary Tree from Inorder and Postorder Traversal (traversal→tree)
- LeetCode #536 — Construct Binary Tree from String (parsing a tree encoding)
- LeetCode #150 — Evaluate Reverse Polish Notation (stack-fold of a token stream)
- LeetCode #1008 — Construct BST from Preorder Traversal (preorder structure)
