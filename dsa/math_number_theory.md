# Math & Number Theory

> **Category:** Math / Number Theory
> **Typical complexity:** O(log n) per operation (GCD, fast power, digit loops), O(n log log n) for a sieve, O(1) for closed-form formulas

---

## What it is

"Math / Number Theory" problems are solved not by a data structure but by a
**property of the numbers themselves**: divisibility, primes, GCD/LCM, modular
arithmetic, digit manipulation, base conversion, combinatorics, or a closed-form
formula that replaces an O(n) simulation with O(1) or O(log n) arithmetic.

The toolbox splits into a few recurring families:

| Family | Core operations | Typical complexity |
|--------|-----------------|--------------------|
| **Digit manipulation** | pop/push digits with `% 10` and `/ 10`, reverse, palindrome checks | O(number of digits) = O(log₁₀ n) |
| **GCD / LCM / divisibility** | Euclid's algorithm, `lcm(a,b) = a/gcd(a,b)*b`, counting divisors | O(log min(a,b)) |
| **Primes** | trial division to √n, Sieve of Eratosthenes, prime factorisation | O(√n) / O(n log log n) |
| **Modular arithmetic** | `(a+b) % m`, `(a*b) % m`, fast power, Fermat inverse | O(log e) for powers |
| **Fast exponentiation** | binary (square-and-multiply) power, matrix power for linear recurrences | O(log n) |
| **Combinatorics** | nCr, factorials, Catalan numbers, permutation ranking (factorial number system) | O(n) precompute, O(1) query |
| **Base conversion / positional systems** | binary, Roman numerals, Excel columns, factorial base | O(log n) |
| **Overflow-safe arithmetic** | check *before* the operation that would overflow | O(1) per check |

The unifying skill: **turn a simulation into arithmetic**. Instead of building
the string / iterating all numbers / recursing over every branch, find the
formula or number-theoretic invariant that jumps straight to the answer.

---

## How to recognise it — signals in the problem statement

| Signal | Example phrasing | Reach for |
|--------|------------------|-----------|
| The input is a *number*, not an array | "Given an integer `x`…" | digit manipulation |
| "Reverse the digits", "is it a palindrome" | LC 7, LC 9 | pop/push digit loop |
| "Without using multiplication / division / mod / built-in pow / sqrt" | LC 29, LC 50, LC 69 | bit shifts, exponentiation by squaring, binary search on answer |
| "Answer may be very large, return it modulo 10⁹+7" | many hard counting problems | modular arithmetic + fast power / Fermat inverse |
| "Numbers given as strings, too big for int64" | LC 43 Multiply Strings, LC 66/67 Plus One / Add Binary | grade-school column arithmetic with carry |
| "How many ways…" with no weights/order constraints beyond counting | LC 62 Unique Paths, LC 96 Unique BSTs, LC 70 Climbing Stairs | nCr, Catalan, Fibonacci closed forms |
| "k-th permutation / k-th element without enumerating" | LC 60 Permutation Sequence | factorial number system |
| "Assume the environment stores 32-bit integers; return 0 on overflow" | LC 7, LC 8 | pre-multiplication overflow checks |
| Cycles / periodicity in an index pattern | LC 6 Zigzag (`period = 2*(rows-1)`) | modular index arithmetic |
| "Prime", "divisor", "GCD", "coprime", "divisible by" | LC 204 Count Primes, LC 1071 GCD of Strings | sieve, Euclid |
| Roman numerals, Excel columns, base-k digits | LC 12, 13, 168, 171 | positional/value-symbol decomposition |

**Rule of thumb:** if a brute-force simulation is obviously O(n) or worse over a
*value* (not a collection) — e.g. "count up to n", "build the whole sequence" —
suspect there's an O(log n) or O(1) mathematical shortcut.

---

## General templates (Go)

### 1. Digit manipulation — pop & push

```go
// reverseDigits reverses the decimal digits of x (sign handled by caller).
// Pattern: pop the last digit with %10, push it with *10 + digit.
// Time: O(log10 x) — one iteration per digit. Space: O(1).
func reverseDigits(x int) int {
    rev := 0
    for x != 0 {
        digit := x % 10 // pop the last digit (in Go, keeps sign of x)
        x /= 10         // drop the last digit (truncates toward zero)
        rev = rev*10 + digit // push digit onto the reversed number
    }
    return rev
}
```

For 32-bit overflow safety, check **before** pushing:

```go
// Before rev = rev*10 + digit, guard against int32 overflow:
if rev > math.MaxInt32/10 || (rev == math.MaxInt32/10 && digit > 7) {
    return 0 // would overflow 2147483647 (…7 is its last digit)
}
if rev < math.MinInt32/10 || (rev == math.MinInt32/10 && digit < -8) {
    return 0 // would underflow -2147483648 (…8 is its last digit)
}
```

### 2. Euclid's GCD and LCM

```go
// gcd computes the greatest common divisor by Euclid's algorithm.
// Invariant: gcd(a, b) == gcd(b, a mod b); the pair shrinks every step.
// Time: O(log min(a,b)). Space: O(1).
func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b // replace (a,b) with (b, a mod b)
    }
    return a
}

// lcm — divide FIRST to avoid overflow of a*b.
func lcm(a, b int) int {
    return a / gcd(a, b) * b
}
```

### 3. Sieve of Eratosthenes

```go
// sieve returns isComposite[0..n]; primes are the indices left false (>=2).
// Time: O(n log log n). Space: O(n).
func sieve(n int) []bool {
    isComposite := make([]bool, n+1)
    for p := 2; p*p <= n; p++ { // only need factors up to sqrt(n)
        if !isComposite[p] {
            for multiple := p * p; multiple <= n; multiple += p {
                isComposite[multiple] = true // p*p is the first NEW composite
            }
        }
    }
    return isComposite
}
```

### 4. Fast exponentiation (binary power / square-and-multiply)

```go
// fastPow computes x^n in O(log n) multiplications.
// Idea: x^n = (x^2)^(n/2) when n even; peel one x off when n odd.
// Each loop iteration halves n -> log2(n) iterations.
func fastPow(x float64, n int) float64 {
    if n < 0 {
        x, n = 1/x, -n // x^-n == (1/x)^n; careful: negate AFTER copying
    }
    result := 1.0
    for n > 0 {
        if n&1 == 1 {   // current lowest bit of n is set
            result *= x // multiply in the current power of x
        }
        x *= x  // square: x, x^2, x^4, x^8, ...
        n >>= 1 // shift to the next bit of the exponent
    }
    return result
}
```

Modular variant (the workhorse for "answer mod 10⁹+7"):

```go
const MOD = 1_000_000_007

// modPow computes (base^exp) % MOD in O(log exp).
func modPow(base, exp int) int {
    base %= MOD
    result := 1
    for exp > 0 {
        if exp&1 == 1 {
            result = result * base % MOD // reduce after EVERY multiply
        }
        base = base * base % MOD
        exp >>= 1
    }
    return result
}

// modInverse via Fermat's little theorem: a^(MOD-2) ≡ a^-1 (mod MOD),
// valid because MOD is prime and a is not a multiple of MOD.
func modInverse(a int) int { return modPow(a, MOD-2) }
```

### 5. Combinatorics — nCr without overflow

```go
// nCr computes C(n, r) = n! / (r! (n-r)!) iteratively.
// Multiply-then-divide in an order that keeps every intermediate an integer:
// after i steps the partial product is C(n, i), always integral.
// Time: O(r). Space: O(1).
func nCr(n, r int) int {
    if r > n-r {
        r = n - r // symmetry: fewer iterations, smaller intermediates
    }
    result := 1
    for i := 1; i <= r; i++ {
        result = result * (n - r + i) / i // divide immediately each step
    }
    return result
}
```

Catalan numbers (Unique BSTs, valid parentheses counts, ballot problems):
`Cat(n) = C(2n, n) / (n+1)`, or the recurrence `Cat(n) = Σ Cat(i)·Cat(n-1-i)`.

### 6. Factorial number system — k-th permutation

```go
// kthPermutation returns the k-th (1-indexed) permutation of 1..n.
// Key fact: with i symbols left, each choice of the leading symbol
// accounts for (i-1)! permutations — so the leading symbol's index
// is (k-1) / (i-1)!, then recurse on the remainder.
// Time: O(n^2) (slice removal). Space: O(n).
func kthPermutation(n, k int) []int {
    fact := 1
    digits := []int{}
    for i := 1; i <= n; i++ {
        fact *= i                    // fact = n!
        digits = append(digits, i)   // available symbols, ascending
    }
    k--          // convert to 0-indexed rank
    result := []int{}
    for i := n; i >= 1; i-- {
        fact /= i                    // fact = (i-1)!
        idx := k / fact              // which of the remaining symbols leads
        k %= fact                    // rank within that block
        result = append(result, digits[idx])
        digits = append(digits[:idx], digits[idx+1:]...) // remove used symbol
    }
    return result
}
```

### 7. Grade-school big-number arithmetic (numbers as strings/arrays)

```go
// addStrings adds two non-negative decimal numbers given as strings.
// Pattern: walk both from the right, keep a carry, prepend digits.
// Works for ANY base b: sum%b is the digit, sum/b is the carry.
// Time: O(max(len(a), len(b))). Space: O(result length).
func addStrings(a, b string) string {
    i, j, carry := len(a)-1, len(b)-1, 0
    out := []byte{}
    for i >= 0 || j >= 0 || carry > 0 { // keep going while ANY input remains
        sum := carry
        if i >= 0 { sum += int(a[i] - '0'); i-- }
        if j >= 0 { sum += int(b[j] - '0'); j-- }
        out = append(out, byte(sum%10)+'0') // current digit
        carry = sum / 10                    // propagate carry
    }
    // digits were collected least-significant first — reverse them
    for l, r := 0, len(out)-1; l < r; l, r = l+1, r-1 {
        out[l], out[r] = out[r], out[l]
    }
    return string(out)
}
```

---

## Worked example — Pow(x, n) by exponentiation by squaring (LC 50)

Compute `fastPow(2.0, 13)`. Binary of 13 is `1101` — read low bit to high:

| Iter | `n` (binary) | low bit `n&1` | `result` after step | `x` after squaring |
|------|--------------|---------------|---------------------|--------------------|
| start | 1101 (13) | — | 1 | 2 |
| 1 | 1101 | **1** → multiply | 1 × 2 = **2** | 2² = 4 |
| 2 | 110 (6) | 0 → skip | 2 | 4² = 16 |
| 3 | 11 (3) | **1** → multiply | 2 × 16 = **32** | 16² = 256 |
| 4 | 1 (1) | **1** → multiply | 32 × 256 = **8192** | 256² = 65536 |
| — | 0 | loop ends | **8192** | — |

Check: 2¹³ = 8192. ✓ Only **4 iterations** (⌈log₂ 13⌉) and 3 multiplies into
`result` — one per set bit of 13 = 8+4+1, i.e. `2¹³ = 2⁸ · 2⁴ · 2¹`. The
squaring line manufactures exactly those powers `2¹, 2², 2⁴, 2⁸`, and the bit
test picks which ones to keep. That is the whole trick: the exponent's binary
representation *is* the multiplication plan.

Same skeleton, swapping the operation, gives:
- **modPow** — multiply mod m (counting problems mod 10⁹+7)
- **matrix power** — O(log n) Fibonacci / linear recurrences (LC 70 follow-up)
- **fast doubling** — F(2k) = F(k)·(2F(k+1) − F(k))

---

## Common pitfalls (and how to avoid them)

1. **Go's `%` keeps the sign of the dividend.** `-7 % 3 == -1` in Go (unlike
   Python's `2`). For a guaranteed non-negative residue use
   `((a % m) + m) % m`. Bites hard in modular arithmetic with subtraction:
   always write `(a - b + m) % m`, never `(a - b) % m`.

2. **Go's integer division truncates toward zero.** `-7 / 2 == -3`, not `-4`.
   This is exactly what LC 29 (Divide Two Integers) requires — but it is NOT
   floor division; don't assume floor semantics from other languages.

3. **Overflow: check *before* the operation, not after.** `rev*10 + d` may
   already have wrapped by the time you compare it. Guard with
   `rev > math.MaxInt32/10` style pre-checks (see template 1). On LeetCode,
   Go's `int` is 64-bit, so 32-bit problems won't wrap "for free" — you must
   compare against `math.MaxInt32` / `math.MinInt32` explicitly.

4. **`-MinInt` overflows.** `-(-2147483648)` doesn't fit in int32. In
   `fastPow`, negating `n = math.MinInt64` (or int32 min in 32-bit settings)
   is undefined territory — handle the minimum value as a special case or
   widen the type first.

5. **`lcm(a,b) = a*b/gcd(a,b)` overflows.** Divide first:
   `a / gcd(a,b) * b`.

6. **Reduce mod m after *every* multiplication.** `a*b % m` is safe for
   m ≈ 10⁹ in 64-bit ints (product < 10¹⁸ < 2⁶³), but chaining two multiplies
   before reducing is not. Never `a*b*c % m`.

7. **Division doesn't distribute over mod.** `(a/b) % m ≠ (a%m)/(b%m)`. Use
   the modular inverse (`modInverse`, template 4) — requires m prime (10⁹+7
   is) and `b` not a multiple of m.

8. **Float precision.** `math.Pow`, `math.Sqrt` and float accumulation can be
   off by 1 ulp — fatal when you cast to int (`int(math.Sqrt(2147395600))`
   can land on the wrong side). For integer answers (LC 69 Sqrt(x)) use
   binary search or Newton's method on integers, and verify with `mid*mid <= x`.

9. **Sieve inner loop must start at `p*p`, not `2*p`.** Smaller multiples were
   already crossed off by smaller primes; starting at `2*p` is correct but
   wastes work — and forgetting `p*p <= n` bounds check panics on overflow-ish
   large n. Also remember 0 and 1 are neither prime nor composite.

10. **nCr with factorials overflows immediately.** 21! > 2⁶³. Use the
    incremental multiply-then-divide loop (template 5), Pascal's triangle DP,
    or factorials under a modulus with Fermat inverses.

11. **Carry loop termination.** In string/array addition, the loop condition
    must include `carry > 0` — otherwise `99 + 1` drops the final `1`
    (LC 66 Plus One's classic trap: all-nines needs a brand-new leading digit).

12. **Don't string-convert when the problem forbids it or when digits
    suffice.** "Reverse half the number" (LC 9) avoids both string allocation
    and overflow, since you only rebuild half the digits.

---

## Problems in this repo

| Problem | What math it uses |
|---------|-------------------|
| [0002 — Add Two Numbers](../0002_add_two_numbers/README.md) | column addition with carry propagation (`% 10`, `/ 10`) |
| [0006 — Zigzag Conversion](../0006_zigzag_conversion/README.md) | periodic index formula, cycle length `2*(numRows-1)` |
| [0007 — Reverse Integer](../0007_reverse_integer/README.md) | digit pop/push, pre-multiplication 32-bit overflow guards |
| [0008 — String to Integer (atoi)](../0008_string_to_integer_atoi/README.md) | digit accumulation with overflow clamping |
| [0009 — Palindrome Number](../0009_palindrome_number/README.md) | digit reversal (reverse-half trick avoids overflow) |
| [0012 — Integer to Roman](../0012_integer_to_roman/README.md) | greedy value-symbol decomposition, digit-position lookup |
| [0013 — Roman to Integer](../0013_roman_to_integer/README.md) | positional value system with subtractive pairs |
| [0029 — Divide Two Integers](../0029_divide_two_integers/README.md) | division via repeated doubling / bit shifts, truncation toward zero, MinInt edge case |
| [0043 — Multiply Strings](../0043_multiply_strings/README.md) | grade-school multiplication, digit position `i+j` / `i+j+1` |
| [0048 — Rotate Image](../0048_rotate_image/README.md) | coordinate transform: transpose + reflect = 90° rotation |
| [0050 — Pow(x, n)](../0050_powx_n/README.md) | exponentiation by squaring, O(log n), negative-exponent handling |
| [0060 — Permutation Sequence](../0060_permutation_sequence/README.md) | factorial number system, k-th permutation without enumeration |
| [0062 — Unique Paths](../0062_unique_paths/README.md) | combinatorics: answer is C(m+n−2, m−1), overflow-safe nCr |
| [0066 — Plus One](../0066_plus_one/README.md) | carry propagation on a digit array, all-nines expansion |
| [0067 — Add Binary](../0067_add_binary/README.md) | base-2 column addition (same carry template, base swapped) |
| [0069 — Sqrt(x)](../0069_sqrt_x/README.md) | integer square root: binary search on answer / Newton's method |
| [0070 — Climbing Stairs](../0070_climbing_stairs/README.md) | Fibonacci; closed form / matrix power as O(log n) follow-ups |
| [0089 — Gray Code](../0089_gray_code/README.md) | formula `i ^ (i >> 1)`, reflect-and-prefix construction |
| [0096 — Unique Binary Search Trees](../0096_unique_binary_search_trees/README.md) | Catalan numbers: recurrence and C(2n,n)/(n+1) closed form |

*(Problems 0131–0400 are being added concurrently; a later pass will extend
this list.)*

---

## Related concepts in this library

- [`/dsa/binary_search.md`](binary_search.md) — binary search *on the answer* (integer sqrt, k-th value problems)
- Bit manipulation — shifts as multiply/divide by 2, `n & 1` parity, `i ^ (i >> 1)` Gray code
- Dynamic programming — when the counting recurrence has no clean closed form
