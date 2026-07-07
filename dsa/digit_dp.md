# Digit DP (Digit Dynamic Programming)

> **Also known as:** digit-by-digit counting, combinatorial digit counting.
> **Core question it answers:** "How many integers in `[0, N]` (or `[L, R]`) satisfy some
> property of their decimal digits?" — without iterating over all N numbers.
> **Core trick:** build the number one digit at a time from the most-significant end,
> carrying a **`tight`** flag that says "so far my prefix equals N's prefix, so my next
> digit is capped".

---

## What it is

Digit DP is a counting technique. Instead of testing every number up to `N` one by one
(O(N), impossible when N ≈ 10⁹ or 10¹⁸), you **fix the digits from left to right** and
count how many completions of the remaining suffix are valid.

The whole method rests on one observation about the range `[0, N]`. Look at the numbers
digit-position by digit-position. At each position you either:

1. place a digit **strictly less** than N's digit there — and then *every* choice for all
   remaining positions is allowed (the number is already guaranteed `< N`), so the suffix
   is "free"; or
2. place a digit **equal** to N's digit — you stay "tight" against the bound and must keep
   respecting it at the next position; or
3. place a digit **greater** than N's digit — forbidden, it would exceed `N`.

That single **tight → free** transition is the heart of every digit-DP solution. Once a
prefix goes free, the number of completions is a plain combinatorial count (often a power,
a factorial, or a nested DP that ignores the bound entirely).

There are two flavours you will meet on LeetCode:

- **Direct combinatorial counting** — the property is simple enough that, once free, you can
  write a closed-form count for the suffix (e.g. "how many `k`-digit strings with distinct
  digits" = a falling-factorial product). No memo table needed; you just sum contributions
  position by position. Problems: #357, #400, #233.
- **General memoised digit DP** — the property needs state carried between digits (a bitmask
  of used digits, a running remainder mod m, a "started yet?" leading-zero flag, previous
  digit, count of some digit so far…). You recurse over
  `(position, state, tight, started)` and memoise the `tight == false` sub-results. This is
  the fully general hammer for "count numbers with property P" (e.g. #902, #1012, #1067, #788).

`[L, R]` ranges are handled by the standard subtraction:
`count(L..R) = f(R) − f(L−1)`, where `f(X)` counts valid numbers in `[0, X]`.

### Mental model

Think of the decimal expansion of `N` as a **fence** you walk beside. As long as you hug
the fence exactly (tight), your next step is limited to at most N's digit. The moment you
step *inside* the fence (choose a smaller digit), you are free to roam the entire field
behind it — and the size of that field is pure combinatorics.

---

## When to recognise it

| Signal in the problem | Why digit DP fits |
|-----------------------|-------------------|
| "Count numbers in `[1, N]` / `[L, R]` such that …" with **N up to 10⁹–10¹⁸** | Range is astronomically large but the *number of digits* is tiny (≤ 18). Iterate over ~18 positions, not N values. |
| The predicate is a **property of the digits** (distinct digits, no `4`, digit-sum divisible by k, contains "13", monotone digits, occurrences of digit `d`) | State carried between positions is small; the answer decomposes digit by digit. |
| "How many times does digit `d` appear across all numbers `1..N`" (#233) | Fix each digit-position as the "d-slot"; count high/low completions combinatorially. |
| "How many k-digit numbers have all-distinct digits" (#357) | Once the first free digit is placed, remaining slots are a falling factorial `9 · 9 · 8 · …`. Pure combinatorics per length. |
| "Find the Nth digit of the infinite string 1234567891011…" (#400) | Numbers are grouped into **digit-length blocks** (1-digit: 9 numbers, 2-digit: 90, …). Count block by block — the same "count by digit-length" arithmetic as digit DP's free case. |
| Counting strobogrammatic / palindromic numbers of a given length within a bound (#248) | Build from both ends by length; within a length, respect the low/high bound like a tight flag. |

**When *not* to use it:** you need the numbers themselves, not just a count (enumerate /
backtrack instead); or the property is about the *value* (`n % 7 == 0`) rather than its
digits and N is small enough to loop. Digit DP counts; it does not list.

---

## General templates (Go)

### Template A — general memoised digit DP

The reusable skeleton. State: index into the digit array, whatever problem state you carry
(here a bitmask of used digits as an example), `tight` (prefix still equals the bound), and
`started` (have we placed a non-leading-zero digit yet — needed so leading zeros don't count
as "used digits" or as a started number).

```go
// countLE returns how many integers in [0, n] satisfy the digit predicate.
// This example counts numbers whose decimal digits are all DISTINCT.
func countLE(n int) int {
    if n < 0 {
        return 0
    }
    digits := toDigits(n) // most-significant first, e.g. 325 -> [3,2,5]
    L := len(digits)

    // memo[pos][mask] is valid ONLY when tight==false && started==true,
    // because those are the states independent of the specific prefix.
    // -1 = not computed yet.
    memo := make([][]int, L)
    for i := range memo {
        memo[i] = make([]int, 1<<10)
        for j := range memo[i] {
            memo[i][j] = -1
        }
    }

    var dp func(pos, mask int, tight, started bool) int
    dp = func(pos, mask int, tight, started bool) int {
        if pos == L {
            // A full number was built. Count it (started guards the empty/zero case
            // if the problem excludes 0 — adjust per problem).
            return 1
        }
        if !tight && started && memo[pos][mask] != -1 {
            return memo[pos][mask] // reuse the bound-independent sub-answer
        }

        limit := 9
        if tight {
            limit = digits[pos] // capped by the bound's digit at this position
        }

        total := 0
        for d := 0; d <= limit; d++ {
            nextStarted := started || d > 0
            if !nextStarted {
                // still in the leading-zero prefix: this position contributes nothing
                // to the mask and the number "hasn't started".
                total += dp(pos+1, mask, false, false)
                continue
            }
            if mask&(1<<d) != 0 {
                continue // predicate-specific prune: digit d already used → skip
            }
            total += dp(pos+1, mask|(1<<d), tight && d == limit, true)
        }

        if !tight && started {
            memo[pos][mask] = total // safe to cache: independent of the prefix
        }
        return total
    }

    return dp(0, 0, true, false)
}

func toDigits(n int) []int {
    if n == 0 {
        return []int{0}
    }
    var d []int
    for n > 0 {
        d = append(d, n%10)
        n /= 10
    }
    // reverse to most-significant-first
    for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
        d[i], d[j] = d[j], d[i]
    }
    return d
}
```

Key rules for the memo:
- **Only cache states with `tight == false`.** A tight state's answer depends on the exact
  bound prefix, so it is not reusable across different `N`s or positions.
- Include every piece of carried state in the memo key. Forgetting `started` (leading-zero
  flag) is the classic bug — it conflates "the number 07" with "the number 7".

### Template B — direct combinatorial counting (no memo)

When the suffix count has a closed form, skip the table and sum contributions. Example: count
of positive integers `< 10^k` with all-distinct digits (the shape of #357).

```go
// countUniqueDigits returns how many integers in [0, 10^k) have all-distinct digits.
func countUniqueDigits(k int) int {
    if k == 0 {
        return 1 // only the number 0
    }
    count := 1 // the number 0 itself
    // For each length len = 1..k, count numbers with exactly `len` digits and no repeats:
    //   first digit: 9 choices (1..9, no leading zero)
    //   second:      9 choices (0..9 minus the one used)
    //   third:       8, then 7, ...  -> falling factorial
    uniqueForLen := 9 // running product for the first digit
    available := 9    // choices remaining for the next position
    for length := 1; length <= k; length++ {
        count += uniqueForLen
        uniqueForLen *= available // extend to one more position
        available--
    }
    return count
}
```

### Template C — digit-block counting (the #400 "Nth digit" family)

Not a bound-tight DP, but the same "group numbers by digit-length and count each block"
arithmetic that powers digit DP's free case.

```go
// findNthDigit returns the n-th digit (1-indexed) of the string 1,2,3,...,10,11,12,...
func findNthDigit(n int) int {
    // Block of `length`-digit numbers has `count` numbers, contributing count*length digits.
    length := 1               // current digit-length block (1-digit, then 2-digit, ...)
    count := 9                // how many numbers have this many digits (9, 90, 900, ...)
    start := 1                // first number with this many digits (1, 10, 100, ...)

    // 1. Skip whole blocks until n lands inside the current block.
    for n > length*count {
        n -= length * count
        length++
        count *= 10
        start *= 10
    }

    // 2. Locate the exact number containing the n-th digit.
    number := start + (n-1)/length // which number in the block
    // 3. Pick the right digit within that number.
    indexInNumber := (n - 1) % length
    s := strconv.Itoa(number)
    return int(s[indexInNumber] - '0')
}
```

---

## Worked example — counting the digit `1` from `1` to `N = 13` (LeetCode #233)

We want the total number of times the digit `1` appears when writing out `1, 2, …, 13`.
The digit-DP / combinatorial view fixes each **position** (units, tens, …) and asks: "in how
many of these numbers is *this* position a `1`?"

For a position with place value `p`, split `N` into `high`, `cur`, `low`:

- `high = N / (p*10)` — the digits above this position
- `cur  = (N / p) % 10` — the digit *at* this position
- `low  = N % p` — the digits below this position

The count of `1`s contributed by this position is:

| Case on `cur` | Ones contributed at this position |
|---------------|-----------------------------------|
| `cur == 0` | `high * p` |
| `cur == 1` | `high * p + (low + 1)` |
| `cur >= 2` | `(high + 1) * p` |

Trace for `N = 13`:

**Units position, `p = 1`:** `high = 13/10 = 1`, `cur = (13/1)%10 = 3`, `low = 13%1 = 0`.
`cur = 3 ≥ 2` → contributes `(high+1)*p = (1+1)*1 = 2`.
(Sanity: units digit is `1` in `1` and `11` → 2 ones. ✓)

**Tens position, `p = 10`:** `high = 13/100 = 0`, `cur = (13/10)%10 = 1`, `low = 13%10 = 3`.
`cur = 1` → contributes `high*p + (low+1) = 0*10 + (3+1) = 4`.
(Sanity: tens digit is `1` in `10, 11, 12, 13` → 4 ones. ✓)

**Higher positions:** `high` becomes 0 and `cur` 0 → contribute 0.

**Total = 2 + 4 = 6.** Enumerating by hand: `1,10,11,12,13` give ones counts
`1,1,2,1,1` = **6**. ✓ — computed in O(number of digits) = O(log N), never touching most of
the 13 numbers individually (and identical work for N = 10¹⁸).

```go
// countDigitOne counts occurrences of the digit 1 in all numbers from 1 to n.
// Time O(log n), Space O(1).
func countDigitOne(n int) int {
    count := 0
    for p := 1; p <= n; p *= 10 { // p = 1, 10, 100, ... place values
        high := n / (p * 10)
        cur := (n / p) % 10
        low := n % p
        switch {
        case cur == 0:
            count += high * p
        case cur == 1:
            count += high*p + low + 1
        default: // cur >= 2
            count += (high + 1) * p
        }
    }
    return count
}
```

---

## Complexity

Let `D` = number of digits of `N` (≈ `log₁₀ N`, at most ~19 for 64-bit), and `S` = size of
the carried state space (e.g. `2¹⁰` for a used-digit mask, `m` for "remainder mod m").

| Variant | Time | Space | Reason |
|---------|------|-------|--------|
| Direct combinatorial (#233, #357) | **O(D)** | **O(1)** | One pass over the digit positions; each contributes a closed-form count. |
| Digit-block counting (#400) | **O(D)** | **O(1)** | Skip at most D length-blocks, then O(1) index arithmetic. |
| General memoised digit DP | **O(D · S · 10)** | **O(D · S)** | Each `(pos, state)` computed once (for the free branch), trying 10 digits. The tight branch is a single O(D) spine on top. |

The decisive win: complexity depends on the **number of digits**, not the magnitude of `N`.
Counting a property up to 10¹⁸ costs the same ~19-position pass as up to 10.

---

## Common pitfalls

1. **Caching tight states.** Only `tight == false` sub-results are prefix-independent and
   therefore reusable. Memoising tight states (or keying the memo without the tight flag)
   corrupts counts. Standard fix: never write to `memo` unless `!tight`.

2. **Dropping the leading-zero (`started`) flag.** Without it, the prefix `0…0` counts as
   real digits — a `mask`-based "distinct digits" DP will think many zeros were "used", and
   short numbers get miscounted. Carry `started` and treat leading zeros as contributing
   nothing until the number "starts".

3. **Off-by-one on the range.** `f(X)` usually counts `[0, X]`. If the problem wants `[1, N]`,
   decide explicitly whether `0` is included and subtract it; for `[L, R]` use
   `f(R) − f(L−1)`, and mind that `L−1` can underflow if `L = 0`.

4. **`cur == 1` boundary in #233-style counting.** When the current digit *equals* the target
   digit, the contribution is `high*p + (low + 1)` — the `+1` (for `low = 0`) and the inclusion
   of the low part are the two most-missed pieces. Derive it from a tiny example each time.

5. **Integer overflow.** Products like `high * p` or falling factorials grow fast; for
   N near 10¹⁸ use `int` on 64-bit Go (which is 64-bit) or `int64`, and watch intermediate
   multiplications.

6. **Digit order.** `n % 10` peels the **least**-significant digit first; digit DP walks
   **most**-significant first. Reverse the slice (or index from the top) — mixing the two
   silently caps the wrong position.

7. **Confusing "value" properties with "digit" properties.** Digit DP only helps when the
   predicate reads the decimal digits. "Divisible by 7" is a value property (though
   *digit-sum* mod k *is* expressible as digit-DP state). Check what the predicate actually
   inspects before reaching for this tool.

8. **#400 block arithmetic off-by-one.** The digit index within the located number is
   `(n-1) % length`, and the number is `start + (n-1)/length`. Both use `n-1` because the
   string is 1-indexed — a very common fencepost slip.

---

## Problems in this repo that use it

- [0233 — Number of Digit One](/0233_number_of_digit_one/README.md) — count occurrences of digit `1` in `1..N` by fixing each place value; the `high/cur/low` combinatorial split (O(log N)).
- [0357 — Count Numbers with Unique Digits](/0357_count_numbers_with_unique_digits/README.md) — count all-distinct-digit numbers in `[0, 10^n)`; falling-factorial per digit-length (Template B).
- [0400 — Nth Digit](/0400_nth_digit/README.md) — digit-block counting: group `1,2,3,…` into 1-digit/2-digit/… blocks and index into the right number (Template C).
- [0248 — Strobogrammatic Number III](/0248_strobogrammatic_number_iii/README.md) — count strobogrammatic numbers in `[low, high]` by building from both ends per length, respecting the bounds like tight flags.

### Related classics to know (not yet in repo)

- LeetCode #902 — Numbers At Most N Given Digit Set (canonical tight-flag digit DP)
- LeetCode #788 — Rotated Digits (per-digit state: is the number "good"?)
- LeetCode #1012 — Numbers With Repeated Digits (complement: total − distinct-digit count)
- LeetCode #600 — Non-negative Integers without Consecutive Ones (digit DP in **binary**)
- LeetCode #1067 — Digit Count in Range (the `[L, R]` generalisation of #233)
