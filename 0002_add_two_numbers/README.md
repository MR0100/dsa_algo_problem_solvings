# 0002 — Add Two Numbers

> LeetCode #2 · Difficulty: Medium
> **Categories:** Linked List, Math, Recursion, Simulation

---

## Problem Statement

You are given two **non-empty** linked lists representing two non-negative integers. The digits are stored in **reverse order**, and each of their nodes contains a single digit. Add the two numbers and return the sum as a linked list.

You may assume the two numbers do not contain any leading zero, except the number 0 itself.

**Example 1**
```
Input:  l1 = [2,4,3], l2 = [5,6,4]
Output: [7,0,8]
Explanation: 342 + 465 = 807.
```

**Example 2**
```
Input:  l1 = [0], l2 = [0]
Output: [0]
```

**Example 3**
```
Input:  l1 = [9,9,9,9,9,9,9], l2 = [9,9,9,9]
Output: [8,9,9,9,0,0,0,1]
Explanation: 9999999 + 9999 = 10009998.
```

**Constraints**
- The number of nodes in each list is in the range `[1, 100]`.
- `0 <= Node.val <= 9`
- It is guaranteed that the list represents a number that does not have leading zeros.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |
| Uber      | ★★☆☆☆ Low       | 2023          |
| LinkedIn  | ★★☆☆☆ Low       | 2022          |
| Netflix   | ★☆☆☆☆ Rare      | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community
> interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — digits are stored as nodes; we must traverse node-by-node and build a new list as output. → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Math / Carry Propagation** — elementary column addition: `sum = d1 + d2 + carry`, write `sum % 10`, propagate `carry = sum / 10`. → see [`/dsa/math.md`](/dsa/math.md)
- **Recursion** — the problem has a natural recursive structure: handle one digit column, recurse on the rest with the new carry.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Convert to Numbers (⚠️ overflows) | O(max(m,n)) | O(max(m,n)) | Never in production; shows the naive idea |
| 2 | Iterative Carry ✅ | O(max(m,n)) | O(1) aux | General case — clean and stack-safe |
| 3 | Recursive Carry | O(max(m,n)) | O(max(m,n)) | Elegant but risks stack overflow for very long lists |

---

## Approach 1 — Convert to Numbers, Add, Convert Back

### Intuition
The list `[2,4,3]` represents 342. Read it into an integer, do the same for the second list, add the integers, then convert the sum back into a reversed linked list.

### Algorithm
1. Walk `l1`, accumulating `num1 += digit × 10^position`.
2. Walk `l2` similarly into `num2`.
3. `sum = num1 + num2`.
4. Build result list: while `sum > 0`, append `sum % 10` as a node, then `sum /= 10`.

### Complexity
- **Time:** O(max(m,n)) — three linear passes.
- **Space:** O(max(m,n)) — the result list.

### ⚠️ Why this fails for large inputs
Go's `int` is 64-bit and holds up to ~9.2 × 10¹⁸. A list of 19 digits already overflows. LeetCode allows lists up to 100 nodes — `10^100` is astronomically beyond int64 range. **This approach produces wrong answers for long inputs.**

### Code
```go
func convertAndAdd(l1 *ListNode, l2 *ListNode) *ListNode {
    num1, mul := 0, 1
    for l1 != nil {
        num1 += l1.Val * mul
        mul *= 10
        l1 = l1.Next
    }
    num2, mul := 0, 1
    for l2 != nil {
        num2 += l2.Val * mul
        mul *= 10
        l2 = l2.Next
    }
    sum := num1 + num2
    if sum == 0 {
        return &ListNode{Val: 0}
    }
    dummy := &ListNode{}
    cur := dummy
    for sum > 0 {
        cur.Next = &ListNode{Val: sum % 10}
        cur = cur.Next
        sum /= 10
    }
    return dummy.Next
}
```

### Dry Run — Example 1: `l1=[2,4,3]`, `l2=[5,6,4]`
```
num1 = 2 + 4×10 + 3×100 = 342
num2 = 5 + 6×10 + 4×100 = 465
sum  = 807
Build: 807%10=7 → 80%10=0 → 8%10=8 → list=[7,0,8] ✓
```

---

## Approach 2 — Iterative Carry (Optimal)

### Intuition
Because the lists already store digits in reverse (LSB first), column addition can proceed left-to-right through the list — exactly the order of manual addition. A `carry` variable tracks what spills into the next column. A dummy head node avoids the special case of initialising the result head.

### Algorithm
1. `dummy = new node`, `cur = dummy`, `carry = 0`.
2. While `l1 != nil` OR `l2 != nil` OR `carry != 0`:
   - `sum = carry + (l1.Val if l1 else 0) + (l2.Val if l2 else 0)`
   - `carry = sum / 10`
   - Append `new node(sum % 10)` after `cur`, advance `cur`.
   - Advance `l1` and `l2` if non-nil.
3. Return `dummy.Next`.

### Complexity
- **Time:** O(max(m,n)) — one pass; at most one extra iteration for the final carry.
- **Space:** O(1) auxiliary — output list is not counted; no extra storage beyond `carry` and pointers.

### Code
```go
func iterativeCarry(l1 *ListNode, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
    cur := dummy
    carry := 0
    for l1 != nil || l2 != nil || carry != 0 {
        sum := carry
        if l1 != nil { sum += l1.Val; l1 = l1.Next }
        if l2 != nil { sum += l2.Val; l2 = l2.Next }
        carry = sum / 10
        cur.Next = &ListNode{Val: sum % 10}
        cur = cur.Next
    }
    return dummy.Next
}
```

### Dry Run — Example 3: `l1=[9,9,9,9,9,9,9]`, `l2=[9,9,9,9]`

| Step | l1.Val | l2.Val | carry_in | sum | carry_out | node |
|------|--------|--------|----------|-----|-----------|------|
| 1 | 9 | 9 | 0 | 18 | 1 | 8 |
| 2 | 9 | 9 | 1 | 19 | 1 | 9 |
| 3 | 9 | 9 | 1 | 19 | 1 | 9 |
| 4 | 9 | 9 | 1 | 19 | 1 | 9 |
| 5 | 9 | nil | 1 | 10 | 1 | 0 |
| 6 | 9 | nil | 1 | 10 | 1 | 0 |
| 7 | 9 | nil | 1 | 10 | 1 | 0 |
| 8 | nil | nil | 1 | 1 | 0 | 1 |

Result: `[8,9,9,9,0,0,0,1]` ✓ (= 10009998 in forward order)

---

## Approach 3 — Recursive Carry

### Intuition
The iterative loop maps naturally to a recursion: process one digit column per call, pass the carry to the next call, return when both lists are exhausted and carry is zero.

### Algorithm
```
helper(l1, l2, carry):
  if l1==nil AND l2==nil AND carry==0 → return nil
  sum = carry + (l1.Val if l1) + (l2.Val if l2)
  node = new ListNode(sum % 10)
  node.Next = helper(l1.Next, l2.Next, sum/10)
  return node
```

### Complexity
- **Time:** O(max(m,n)) — one call per digit column.
- **Space:** O(max(m,n)) — call stack depth equals number of columns. For lists up to 100 nodes this is fine; for 10 000+ node lists it could stack-overflow.

### Code
```go
func recursiveCarry(l1 *ListNode, l2 *ListNode) *ListNode {
    return addHelper(l1, l2, 0)
}

func addHelper(l1, l2 *ListNode, carry int) *ListNode {
    if l1 == nil && l2 == nil && carry == 0 {
        return nil
    }
    sum := carry
    if l1 != nil { sum += l1.Val; l1 = l1.Next }
    if l2 != nil { sum += l2.Val; l2 = l2.Next }
    node := &ListNode{Val: sum % 10}
    node.Next = addHelper(l1, l2, sum/10)
    return node
}
```

### Dry Run — Example 1: `l1=[2,4,3]`, `l2=[5,6,4]`
```
call(l1=[2,4,3], l2=[5,6,4], carry=0): sum=7, node(7) → call([4,3],[6,4],0)
  call([4,3],[6,4],0): sum=10, node(0) → call([3],[4],1)
    call([3],[4],1): sum=8, node(8) → call(nil,nil,0)
      call(nil,nil,0): return nil
    node(8).Next = nil    → return node(8)
  node(0).Next = node(8) → return node(0)
node(7).Next = node(0)   → return node(7)

Result: 7 → 0 → 8  ✓
```

---

## Key Takeaways

- **Dummy head node** — when building a linked list from scratch, attaching a dummy sentinel at the start means the first real node is just `dummy.Next`, eliminating the `if head == nil` branch. This pattern appears in almost every linked-list construction problem.
- **Carry flush** — the loop condition `l1 != nil || l2 != nil || carry != 0` ensures a final carry (e.g. `9+9=18` at the last column) generates its own node. Forgetting `|| carry != 0` is the classic bug.
- **Reversed storage is a gift** — the problem could have stored digits MSB-first, requiring a full reversal before and after. Reversed storage means we can iterate directly from the head without any reversal.
- **Iterative vs recursive** — recursive code is more elegant but uses O(n) stack space. For interview purposes both are acceptable; in production, prefer iterative for long lists.
- **Never convert to int for big-number arithmetic** — any time a problem says "each node is a digit" with up to 100 nodes, integer conversion will overflow. Always simulate digit-by-digit.

---

## Related Problems

- LeetCode #445 — Add Two Numbers II (digits stored MSB-first; requires reversal or stack)
- LeetCode #43 — Multiply Strings (similar digit-by-digit simulation)
- LeetCode #415 — Add Strings (same carry logic on strings)
- LeetCode #67 — Add Binary (carry logic in base 2)
