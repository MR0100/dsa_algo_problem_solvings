# 0043 — Multiply Strings

> LeetCode #43 · Difficulty: Medium
> **Categories:** Math, String, Simulation

---

## Problem Statement

Given two non-negative integers `num1` and `num2` represented as strings, return the product of `num1` and `num2`, also represented as a string.

**Note:** You must not use any built-in BigInteger library or convert the inputs to integer directly.

**Example 1**
```
Input:  num1 = "2", num2 = "3"
Output: "6"
```

**Example 2**
```
Input:  num1 = "123", num2 = "456"
Output: "56088"
```

**Constraints**
- `1 <= num1.length, num2.length <= 200`
- `num1` and `num2` consist of digits only.
- Both `num1` and `num2` do not contain any leading zeros, except the number `0` itself.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Grade-School Multiplication** — multiply each digit pair and accumulate into a position array.
- **Carry Propagation** — at position `p2`, accumulate products and propagate carries to `p1`.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Grade-School Digit × Digit ✅ | O(n × m) | O(n+m) | The standard and only needed approach |

n = len(num1), m = len(num2).

---

## Approach 1 — Grade-School Multiplication (Recommended ✅)

### Intuition
When multiplying two numbers manually:
```
    123
  × 456
  -----
    738  (123 × 6)
   615   (123 × 5, shifted)
  492    (123 × 4, shifted)
```

For each digit pair `(num1[i], num2[j])`, their product contributes to a specific position in the result. An n-digit × m-digit multiplication produces at most n+m digits.

**Position mapping:**
- `num1[i]` is at position `n-1-i` from the right (0-indexed from right).
- `num2[j]` is at position `m-1-j` from the right.
- Their product lands at position `(n-1-i) + (m-1-j) = n+m-2-i-j` from the right.
- In a 0-indexed array of size n+m: `p2 = i+j+1` (units digit), `p1 = i+j` (carry).

### Algorithm
```
pos = array of n+m zeros
for i = n-1 downto 0:
  for j = m-1 downto 0:
    mul = (num1[i]-'0') * (num2[j]-'0')
    p1, p2 = i+j, i+j+1
    sum = mul + pos[p2]
    pos[p2] = sum % 10
    pos[p1] += sum / 10   // carry
strip leading zeros from pos
return as string
```

### Complexity
- **Time:** O(n × m) — nested loops over digit pairs.
- **Space:** O(n+m) — the result array.

### Code
```go
func multiply(num1, num2 string) string {
    if num1 == "0" || num2 == "0" { return "0" }
    n, m := len(num1), len(num2)
    pos := make([]int, n+m)
    for i := n-1; i >= 0; i-- {
        for j := m-1; j >= 0; j-- {
            mul := int(num1[i]-'0') * int(num2[j]-'0')
            p1, p2 := i+j, i+j+1
            sum := mul + pos[p2]
            pos[p2] = sum % 10
            pos[p1] += sum / 10
        }
    }
    var sb strings.Builder
    for _, d := range pos {
        if sb.Len() == 0 && d == 0 { continue }
        sb.WriteByte(byte('0' + d))
    }
    return sb.String()
}
```

### Dry Run — `num1 = "123"`, `num2 = "456"`
```
pos = [0,0,0,0,0,0] (size 6)

i=2 (3), j=2 (6): mul=18. p1=4,p2=5. sum=18. pos[5]=8, pos[4]+=1 → [0,0,0,0,1,8]
i=2 (3), j=1 (5): mul=15. p1=3,p2=4. sum=15+1=16. pos[4]=6, pos[3]+=1 → [0,0,0,1,6,8]
i=2 (3), j=0 (4): mul=12. p1=2,p2=3. sum=12+1=13. pos[3]=3, pos[2]+=1 → [0,0,1,3,6,8]
i=1 (2), j=2 (6): mul=12. p1=3,p2=4. sum=12+6=18. pos[4]=8, pos[3]+=1 → [0,0,1,4,8,8]
i=1 (2), j=1 (5): mul=10. p1=2,p2=3. sum=10+4=14. pos[3]=4, pos[2]+=1+1=2 → [0,0,2,4,8,8]
i=1 (2), j=0 (4): mul=8. p1=1,p2=2. sum=8+2=10. pos[2]=0, pos[1]+=1 → [0,1,0,4,8,8]
i=0 (1), j=2 (6): mul=6. p1=2,p2=3. sum=6+4=10. pos[3]=0, pos[2]+=1 → [0,1,1,0,8,8]
i=0 (1), j=1 (5): mul=5. p1=1,p2=2. sum=5+1=6. pos[2]=6, pos[1]+=0 → [0,1,6,0,8,8]
i=0 (1), j=0 (4): mul=4. p1=0,p2=1. sum=4+1=5. pos[1]=5, pos[0]+=0 → [0,5,6,0,8,8]

pos=[0,5,6,0,8,8]. Skip leading 0. Result: "56088" ✓
```

---

## Key Takeaways

- **`p2 = i+j+1`, `p1 = i+j`** — p2 is the units digit, p1 receives the carry. This mapping is the entire algorithm.
- **No separate carry pass** — carries propagate naturally because `pos[p1] += sum/10` at each step; the final `pos` values can be multi-digit but that's fine since they represent carry accumulation before we emit digits.
- **Short-circuit on "0"** — handle `num1=="0" || num2=="0"` first; otherwise the loop runs but the leading-zero stripping at the end would return "".
- **`strings.Builder` for O(n)** — prefer it over `string +=` to avoid O(n²) string allocations.

---

## Related Problems

- LeetCode #2 — Add Two Numbers (addition on linked list digits)
- LeetCode #66 — Plus One (increment a digit array)
- LeetCode #415 — Add Strings (same idea but addition, not multiplication)
