# 0157 — Read N Characters Given Read4

> LeetCode #157 · Difficulty: Easy (Premium)
> **Categories:** Array, Simulation, Interactive, String

---

## Problem Statement

Given a `file` and assume that you can only read the file using a given method `read4`, implement a method to read `n` characters.

**Method `read4`:**

The API `read4` reads **four consecutive characters** from `file`, then writes those characters into the buffer array `buf4`.

The return value is the number of actual characters read.

Note that `read4()` has its own file pointer, much like `FILE *fp` in C.

**Definition of `read4`:**
```
    Parameter:  char[] buf4
    Returns:    int

buf4[] is a destination, not a source. The results from read4 will be copied to buf4[].
```

Below is a high-level example of how `read4` works:
```
File file("abcde");           // File is "abcde", initially file pointer (fp) points to 'a'
char[] buf4 = new char[4];    // Create buffer with enough space to store characters
read4(buf4);                  // read4 returns 4. Now buf4 = "abcd", fp points to 'e'
read4(buf4);                  // read4 returns 1. Now buf4 = "e", fp points to end of file
read4(buf4);                  // read4 returns 0. Now buf4 = "", fp points to end of file
```

**Method `read`:**

By using the `read4` method, implement the method `read` that reads `n` characters from `file` and store it in the buffer array `buf`. Consider that you **cannot** manipulate `file` directly.

The return value is the number of actual characters read.

**Definition of `read`:**
```
    Parameters: char[] buf, int n
    Returns:    int

buf[] is a destination, not a source. You will need to write the results to buf[].
```

**Note:**
- Consider that you cannot manipulate the file directly. The file is only accessible for `read4` but not for `read`.
- The `read` function will only be called **once** for each test case.
- You may assume the destination buffer array, `buf`, is guaranteed to have enough space for storing `n` characters.

**Example 1:**
```
Input: file = "abc", n = 4
Output: 3
Explanation: After calling your read method, buf should contain "abc". We read a total of 3 characters from the file, so return 3.
Note that "abc" is the file's content, not buf. buf is the destination buffer that you will have to write the results to.
```

**Example 2:**
```
Input: file = "abcde", n = 5
Output: 5
Explanation: After calling your read method, buf should contain "abcde". We read a total of 5 characters from the file, so return 5.
```

**Constraints:**
- `1 <= file.length <= 500`
- `file` consist of English letters and digits.
- `1 <= n <= 1000`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Facebook  | ★★★★★ Very High | 2024          |
| Google    | ★★★☆☆ Medium    | 2023          |
| Microsoft | ★★☆☆☆ Low       | 2023          |
| Amazon    | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **API Simulation / Buffered I/O** — you adapt one fixed-size read API into an arbitrary-size read API; the internal scratch buffer acts as a tiny queue between producer (`read4`) and consumer (`buf`) → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Two Pointers (read/write offsets)** — a write offset `total` into `buf` and the chunk cursor into `buf4` coordinate the copy → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Read Whole File (Brute Force) | O(F) — F = file length | O(F) | Trivially correct; wasteful when n ≪ F |
| 2 | Direct Copy (Optimal) | O(n) | O(1) | Always — stops as soon as n chars are delivered |

---

## Approach 1 — Read Whole File (Brute Force)

### Intuition
Ignore `n` while reading: drain the entire file 4 characters at a time into one big internal slice. Correctness is obvious because we now hold the full file in memory — the answer is just the first `min(n, fileLength)` characters of it.

### Algorithm
1. Create an empty internal slice `all` and a 4-byte scratch buffer `buf4`.
2. Loop: `cnt = read4(buf4)`; append `buf4[:cnt]` to `all`.
3. Stop when `cnt < 4` (read4 returned short → EOF reached).
4. `total = min(n, len(all))`.
5. Copy `all[:total]` into `buf`; return `total`.

### Complexity
- **Time:** O(F) where F is the file length — we always read to EOF even if `n = 1`.
- **Space:** O(F) — the internal slice stores the whole file.

### Code
```go
func bruteForceReadAll(read4 func([]byte) int, buf []byte, n int) int {
    all := []byte{}
    buf4 := make([]byte, 4)
    for {
        cnt := read4(buf4)
        all = append(all, buf4[:cnt]...)
        if cnt < 4 {
            break // EOF
        }
    }
    total := n
    if len(all) < n {
        total = len(all)
    }
    copy(buf, all[:total])
    return total
}
```

### Dry Run
Example 1: `file = "abc"`, `n = 4`.

| Step | read4 returns | buf4 (valid part) | all after append | Loop continues? |
|------|---------------|-------------------|------------------|-----------------|
| 1 | 3 | `"abc"` | `"abc"` | no (3 < 4 → EOF) |

Then `total = min(4, 3) = 3`, copy `"abc"` into `buf`, return **3** ✓

---

## Approach 2 — Direct Copy (Optimal)

### Intuition
We never need more than `n` characters, so stream chunks of ≤ 4 directly into the destination at the current write offset. Two natural stop conditions:
1. `total == n` — the request is satisfied.
2. `read4` returns 0 — the file is exhausted.

One subtlety: the last `read4` may deliver more characters than we still need (e.g. `n = 5` with the file longer — chunk 2 delivers 4 chars but only 1 fits). Clamp `cnt` to `n - total` before copying. Since `read` is called only once per test case, silently discarding the surplus characters is safe (contrast with LeetCode #158 where they must be remembered).

### Algorithm
1. `total = 0`; allocate the fixed scratch buffer `buf4[4]`.
2. While `total < n`:
   1. `cnt = read4(buf4)`.
   2. If `cnt == 0` → EOF, break.
   3. If `total + cnt > n` → clamp `cnt = n - total`.
   4. `copy(buf[total:], buf4[:cnt])`; `total += cnt`.
3. Return `total`.

### Complexity
- **Time:** O(n) — at most ⌈n/4⌉ + 1 calls to `read4`, each doing O(1) work per delivered character.
- **Space:** O(1) — one fixed 4-byte scratch buffer, independent of both `n` and file size.

### Code
```go
func directCopy(read4 func([]byte) int, buf []byte, n int) int {
    buf4 := make([]byte, 4)
    total := 0
    for total < n {
        cnt := read4(buf4)
        if cnt == 0 {
            break // EOF
        }
        if total+cnt > n {
            cnt = n - total // clamp to the request size
        }
        copy(buf[total:], buf4[:cnt])
        total += cnt
    }
    return total
}
```

### Dry Run
Example 1: `file = "abc"`, `n = 4`.

| Iter | total (before) | read4 returns | clamp? | copied into buf | total (after) | loop check |
|------|----------------|---------------|--------|------------------|----------------|------------|
| 1 | 0 | 3 (`"abc"`) | no (0+3 ≤ 4) | `buf[0:3] = "abc"` | 3 | 3 < 4 → continue |
| 2 | 3 | 0 (EOF) | — | — | 3 | break |

Return **3**, `buf = "abc"` ✓

Example 2: `file = "abcde"`, `n = 5`.

| Iter | total (before) | read4 returns | clamp? | copied into buf | total (after) |
|------|----------------|---------------|--------|------------------|----------------|
| 1 | 0 | 4 (`"abcd"`) | no | `buf[0:4] = "abcd"` | 4 |
| 2 | 4 | 1 (`"e"`) | no (4+1 ≤ 5) | `buf[4:5] = "e"` | 5 |

Loop exits (`total == n`), return **5**, `buf = "abcde"` ✓

---

## Key Takeaways
- **Adapter pattern for I/O APIs**: converting a fixed-chunk read into an arbitrary-size read is a real systems pattern (`fread` over raw block reads); the interview version is exactly this problem.
- Two independent termination conditions — *request satisfied* (`total == n`) and *source exhausted* (`read4` returns 0) — must both be handled; forgetting EOF causes an infinite loop.
- **Clamp the final chunk** (`cnt = n - total`): the API delivers up to 4 chars whether you need them or not.
- The "call once" guarantee is what makes discarding surplus characters legal. The moment `read` can be called multiple times (LeetCode #158), you must buffer leftovers between calls — that is the entire difficulty jump from Easy to Hard.

---

## Related Problems
- LeetCode #158 — Read N Characters Given Read4 II — Call Multiple Times (same API, must persist leftovers)
- LeetCode #251 — Flatten 2D Vector (iterator over chunked data)
- LeetCode #604 — Design Compressed String Iterator (streaming consumption of a source)
