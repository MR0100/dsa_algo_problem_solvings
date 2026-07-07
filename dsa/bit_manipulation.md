# Bit Manipulation

> Operating directly on the binary representation of integers using the bitwise
> operators `&` (AND), `|` (OR), `^` (XOR), `&^` (AND NOT, Go-specific),
> `<<` (left shift), and `>>` (right shift). Bit tricks turn set operations,
> parity checks, doubling/halving, and state encoding into single CPU
> instructions — O(1) time, O(1) space.

---

## What it is

Every integer is a fixed-width vector of bits. Bit manipulation treats that
vector as:

1. **A set** — bit `i` set ⇔ element `i` is in the set. Union is `|`,
   intersection is `&`, difference is `&^`, symmetric difference is `^`,
   membership test is `(mask>>i)&1`. A 32-bit int is a set over a universe of
   up to 32 elements with O(1) operations.
2. **A number in base 2** — `x << 1` doubles, `x >> 1` halves (floor),
   `x & 1` tests odd/even, `x % (1<<k)` is `x & ((1<<k)-1)`.
3. **An algebraic object under XOR** — XOR is associative, commutative,
   `a ^ a = 0`, `a ^ 0 = a`. This makes duplicates cancel in any order,
   which powers "find the element that appears once" style problems.

### Truth-table refresher

| a | b | a & b | a \| b | a ^ b | a &^ b |
|---|---|-------|--------|-------|--------|
| 0 | 0 | 0     | 0      | 0     | 0      |
| 0 | 1 | 0     | 1      | 1     | 0      |
| 1 | 0 | 0     | 1      | 1     | 1      |
| 1 | 1 | 1     | 1      | 0     | 0      |

Note: in Go, unary `^x` is bitwise NOT (there is no `~` operator), and binary
`a &^ b` ("AND NOT" / bit clear) equals `a & (^b)`.

---

## How to recognise a bit-manipulation problem

Signals in the problem statement:

- **"Without using `+`, `-`, `*`, `/`, or `%`"** — arithmetic must be rebuilt
  from shifts and XOR/AND (e.g. Add Binary follow-up, Divide Two Integers,
  Sum of Two Integers).
- **"Every element appears twice except one"** (or *k* times except one) —
  XOR cancellation, or per-bit counting mod *k*.
- **"O(1) extra space" on a counting/duplicate problem** — often a hint that
  bits of the answer can be computed independently.
- **Small n (n ≤ ~20) and "all subsets / all combinations / assignments"** —
  enumerate `2^n` bitmasks; each mask *is* a subset. Also the doorway to
  **bitmask DP** (`dp[mask]` = best answer using the set of items in `mask`).
- **Binary strings, powers of two, "number of 1 bits", parity, Gray code** —
  the problem is literally about bits.
- **Need a tiny, hashable, copy-cheap "visited"/state set** — e.g. columns and
  diagonals under attack in N-Queens, 26 letters seen in a word, digits used
  in a Sudoku row.
- **"Divide/multiply fast"** — exponentiation by squaring and long division
  both walk the bits of the exponent/quotient (Pow(x, n), Divide Two
  Integers).

---

## Core toolbox (Go)

```go
// ---- single-bit operations on mask, bit index i (0 = least significant) ----
mask |= 1 << i          // SET bit i               (add i to the set)
mask &^= 1 << i         // CLEAR bit i             (remove i; Go's AND NOT)
mask ^= 1 << i          // TOGGLE bit i            (flip membership)
on := mask&(1<<i) != 0  // TEST bit i              (is i in the set?)

// ---- whole-mask tricks -----------------------------------------------------
low := mask & -mask     // isolate LOWEST set bit  (e.g. 0b10110 -> 0b00010)
mask &= mask - 1        // DROP lowest set bit     (Kernighan's trick)
isPow2 := mask > 0 && mask&(mask-1) == 0 // power of two has exactly one bit
full := (1 << n) - 1    // mask with the n lowest bits all set (the full set)
comp := full &^ mask    // complement within an n-element universe

// ---- counting and inspecting (math/bits package, all O(1)) ------------------
bits.OnesCount(uint(mask))      // popcount: number of set bits
bits.TrailingZeros(uint(mask))  // index of lowest set bit
bits.Len(uint(mask))            // position of highest set bit + 1
```

### Template 1 — enumerate all subsets of n items

```go
// For each of the 2^n masks, bit i of mask says whether item i is chosen.
for mask := 0; mask < 1<<n; mask++ {         // 2^n iterations
    subset := []int{}
    for i := 0; i < n; i++ {                 // decode the mask
        if mask&(1<<i) != 0 {                // bit i set -> nums[i] is in
            subset = append(subset, nums[i])
        }
    }
    // process subset (Subsets, LC #78, uses exactly this)
}
```

### Template 2 — iterate over set bits only (Kernighan)

```go
// Visits each set bit once: O(popcount) instead of O(width).
count := 0
for x != 0 {
    x &= x - 1 // clears the lowest set bit each pass
    count++    // Number of 1 Bits (LC #191) in a few lines
}
```

### Template 3 — enumerate all submasks of a mask

```go
// Classic bitmask-DP inner loop; total cost over all masks is O(3^n).
for sub := mask; sub > 0; sub = (sub - 1) & mask {
    // sub is a non-empty submask of mask
}
// (append the empty submask 0 manually if needed)
```

### Template 4 — addition from XOR + carry (no `+` operator)

```go
// XOR = sum without carries; AND<<1 = the carries. Repeat until no carry.
func add(a, b int) int {
    for b != 0 {
        carry := (a & b) << 1 // positions where both bits are 1 carry left
        a = a ^ b             // bitwise sum ignoring carries
        b = carry             // add the carries in the next round
    }
    return a
}
```

### Template 5 — bitmask as a visited set in backtracking

```go
// N-Queens II style: three masks encode every attacked line in O(1).
var solve func(row, cols, diag1, diag2 int) int
solve = func(row, cols, diag1, diag2 int) int {
    if row == n {
        return 1
    }
    count := 0
    free := ((1 << n) - 1) &^ (cols | diag1 | diag2) // safe columns this row
    for free != 0 {
        bit := free & -free      // pick lowest safe column
        free &^= bit             // remove it from the candidates
        count += solve(row+1, cols|bit, (diag1|bit)<<1, (diag2|bit)>>1)
    }
    return count
}
```

---

## Worked example — Single Number (LC #136) via XOR cancellation

**Problem.** Every element in `nums` appears exactly twice except one, which
appears once. Find it in O(n) time, O(1) space.

**Insight.** `a ^ a = 0` and `a ^ 0 = a`, and XOR is commutative — so XOR-ing
the whole array makes every pair annihilate itself, leaving the unique value.

```go
// singleNumber finds the element appearing once via XOR cancellation.
// Time: O(n) — one pass.  Space: O(1) — a single accumulator.
func singleNumber(nums []int) int {
    acc := 0                // identity element of XOR
    for _, x := range nums {
        acc ^= x            // pairs cancel to 0, unique value survives
    }
    return acc
}
```

**Step-by-step trace** on `nums = [4, 1, 2, 1, 2]` (bits shown 3-wide):

| step | x | x (binary) | acc before | acc ^ x (binary) | acc after |
|------|---|------------|------------|------------------|-----------|
| init | — | —          | —          | —                | 0 (`000`) |
| 1    | 4 | `100`      | `000`      | `000 ^ 100 = 100`| 4 (`100`) |
| 2    | 1 | `001`      | `100`      | `100 ^ 001 = 101`| 5 (`101`) |
| 3    | 2 | `010`      | `101`      | `101 ^ 010 = 111`| 7 (`111`) |
| 4    | 1 | `001`      | `111`      | `111 ^ 001 = 110`| 6 (`110`) |
| 5    | 2 | `010`      | `110`      | `110 ^ 010 = 100`| 4 (`100`) |

Result: `acc = 4`. Note how the second `1` (step 4) flipped bit 0 back off and
the second `2` (step 5) flipped bit 1 back off — each pair undoes itself
regardless of position, because XOR is order-independent.

---

## Common pitfalls (and Go-specific gotchas)

1. **Go has no `~`.** Bitwise NOT is unary `^x`; "clear these bits" is
   `a &^ b`. Writing `~x` is a compile error; porting C tricks needs this
   translation.
2. **Operator precedence:** in Go, `&`, `|`, `^`, `<<`, `>>` bind *tighter*
   than in C relative to comparison, but `mask & 1<<i != 0` still parses in
   surprising ways across languages — always parenthesise:
   `mask&(1<<i) != 0`.
3. **Shifting signed negatives:** `>>` on a signed int is an *arithmetic*
   shift (sign bit copies in), so `-1 >> 1 == -1` forever — a loop like
   `for x != 0 { x >>= 1 }` never terminates for negative `x`. Convert to
   `uint`/`uint32` first, or use `bits.OnesCount32(uint32(x))`.
4. **Shift overflow:** `1 << n` with `n ≥ 63` overflows `int`. For full
   subsets keep `n ≤ ~20–25` anyway (`2^20 ≈ 10^6` masks); for single bits on
   wide types use unsigned.
5. **`x & -x` on the minimum value:** works fine in Go two's complement, but
   negating `math.MinInt32`-style values in intermediate arithmetic (e.g.
   Divide Two Integers, Reverse Integer) overflows — handle the
   `-2^31` edge case explicitly or work in wider/negative space.
6. **XOR-swap and self-aliasing:** `a ^= b; b ^= a; a ^= b` zeroes both when
   `a` and `b` alias the same memory. Just use `a, b = b, a` in Go.
7. **Forgetting XOR only cancels *pairs*:** "appears 3 times except one"
   needs per-bit counting mod 3 (or the two-mask automaton), not a plain XOR.
8. **Off-by-one on the full mask:** the set of n elements is `(1<<n) - 1`,
   not `1<<n`. Similarly loop `mask < 1<<n`, not `<=`.
9. **Mixing up `^` (XOR) with exponentiation:** `2 ^ 10` is 8, not 1024.
   Powers of two are `1 << 10`.
10. **Assuming bit tricks always beat clarity:** for interview code, write the
    readable version first, then mention the bit-level optimisation
    (e.g. N-Queens with maps → N-Queens with three masks).

---

## Problems in this repo

Problems 0131+ are being added; this list covers what exists today.

- [0029 — Divide Two Integers](../0029_divide_two_integers/README.md) —
  long division by doubling the divisor with `<< 1`; no `*`, `/`, `%` allowed.
- [0050 — Pow(x, n)](../0050_powx_n/README.md) — exponentiation by squaring
  walks the bits of `n` (`n&1` for odd, `n >>= 1` to halve).
- [0052 — N-Queens II](../0052_n_queens_ii/README.md) — columns and both
  diagonals encoded as bitmasks for O(1) attack checks (Template 5 above).
- [0067 — Add Binary](../0067_add_binary/README.md) — binary addition with
  carry; follow-up is the XOR/AND-carry loop (Template 4 above).
- [0078 — Subsets](../0078_subsets/README.md) — enumerate all `2^n` bitmasks;
  bit `i` decides membership of `nums[i]` (Template 1 above).
- [0079 — Word Search](../0079_word_search/README.md) — in-place visited
  marking by XOR-ing cell bytes with 255 (self-inverse toggle).
- [0089 — Gray Code](../0089_gray_code/README.md) — direct formula
  `gray(i) = i ^ (i >> 1)`; consecutive codes differ by exactly one bit.
- [0090 — Subsets II](../0090_subsets_ii/README.md) — bitmask enumeration
  with duplicate-skipping on sorted input.
