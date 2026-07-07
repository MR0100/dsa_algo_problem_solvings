# 0158 — Read N Characters Given Read4 II — Call Multiple Times

> LeetCode #158 · Difficulty: Hard (Premium)
> **Categories:** Array, Simulation, Interactive, Design, String

---

## Problem Statement

Given a `file` and assume that you can only read the file using a given method `read4`, implement a method `read` to read `n` characters. Your method `read` **may be called multiple times**.

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
- Consider that you cannot manipulate `file` directly. The file is only accessible for `read4` but not for `read`.
- The `read` function may be called **multiple times**.
- Please remember to **RESET** your class variables declared in Solution, as static/class variables are persisted across multiple test cases. Please see [here](https://leetcode.com/faq/) for more details.
- You may assume the destination buffer array, `buf`, is guaranteed to have enough space for storing `n` characters.
- It is guaranteed that in a given test case the same buffer `buf` is called by `read`.

**Example 1:**
```
Input: file = "abc", queries = [1,2,1]
Output: [1,2,0]
Explanation: The test case represents the following scenario:
File file("abc");
Solution sol;
sol.read(buf, 1); // After calling your read method, buf should contain "a". We read a total of 1 character from the file, so return 1.
sol.read(buf, 2); // Now buf should contain "bc". We read a total of 2 characters from the file, so return 2.
sol.read(buf, 1); // We have reached the end of file, no more characters can be read. So return 0.
Assume buf is allocated and guaranteed to have enough space for storing all characters from the file.
```

**Example 2:**
```
Input: file = "abc", queries = [4,1]
Output: [3,0]
Explanation: The test case represents the following scenario:
File file("abc");
Solution sol;
sol.read(buf, 4); // After calling your read method, buf should contain "abc". We read a total of 3 characters from the file, so return 3.
sol.read(buf, 1); // We have reached the end of file, no more characters can be read. So return 0.
```

**Constraints:**
- `1 <= file.length <= 500`
- `file` consist of English letters and digits.
- `1 <= queries.length <= 10`
- `1 <= queries[i] <= 500`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Facebook  | ★★★★★ Very High | 2024          |
| Google    | ★★★☆☆ Medium    | 2023          |
| Amazon    | ★★☆☆☆ Low       | 2023          |
| Bloomberg | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Queue / Buffered I/O** — leftover characters form a FIFO staging area between the chunked producer (`read4`) and the arbitrary-size consumer (`read`); Approach 1 uses a literal queue, Approach 2 compresses it into a ring-like fixed buffer with head/size cursors → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Two Pointers (cursor management)** — `i4` (consume cursor) and `n4` (fill level) walk the internal buffer, while `total` walks the destination → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Leftover Queue (Brute Force) | O(n) per call | O(n) queue | Easiest to get right under pressure |
| 2 | Persistent Buffer + Pointers (Optimal) | O(n) per call | O(1) — fixed 4 bytes + 2 ints | The expected interview answer |

---

## Approach 1 — Leftover Queue (Brute Force)

### Intuition
The entire difficulty over #157 is that a `read4` chunk can **straddle two `read()` calls**: with file `"abc"`, `read(1)` forces us to pull the full chunk `"abc"` from the file but deliver only `"a"` — if `"bc"` is thrown away, the next call returns garbage. The bluntest fix: never throw anything away. Push every fetched character into a persistent queue; every `read()` pops from the front. The queue *is* the state between calls.

### Algorithm
1. (State) `queue []byte` lives on the reader object, surviving across calls.
2. On `read(buf, n)`: while `len(queue) < n`, call `cnt = read4(buf4)`; if `cnt == 0` break (EOF); else append `buf4[:cnt]` to the queue.
3. `total = min(n, len(queue))`.
4. Copy `queue[:total]` into `buf`.
5. `queue = queue[total:]` — delivered chars are popped; surplus stays for next time.
6. Return `total`.

### Complexity
- **Time:** O(n) per call — every character is enqueued once and dequeued once across the reader's lifetime.
- **Space:** O(n) — after over-reading, the queue can hold up to `n + 3` characters (worst case just under the request size plus one chunk).

### Code
```go
type queueReader struct {
    read4 func([]byte) int
    queue []byte // characters fetched but not yet delivered
}

func (r *queueReader) read(buf []byte, n int) int {
    buf4 := make([]byte, 4)
    for len(r.queue) < n {
        cnt := r.read4(buf4)
        if cnt == 0 {
            break // EOF
        }
        r.queue = append(r.queue, buf4[:cnt]...)
    }
    total := n
    if len(r.queue) < n {
        total = len(r.queue)
    }
    copy(buf, r.queue[:total])
    r.queue = r.queue[total:] // leftovers persist to the next call
    return total
}
```

### Dry Run
Example 1: `file = "abc"`, queries `[1, 2, 1]`.

| Call | queue before | read4 activity | queue after fill | total = min(n, len) | delivered | queue after pop | returns |
|------|--------------|----------------|------------------|---------------------|-----------|-----------------|---------|
| `read(1)` | `""` | returns 3 → `"abc"` | `"abc"` (3 ≥ 1, stop) | min(1,3) = 1 | `"a"` | `"bc"` | **1** ✓ |
| `read(2)` | `"bc"` | none (2 ≥ 2 already) | `"bc"` | min(2,2) = 2 | `"bc"` | `""` | **2** ✓ |
| `read(1)` | `""` | returns 0 → EOF | `""` | min(1,0) = 0 | `""` | `""` | **0** ✓ |

Output `[1,2,0]` ✓

---

## Approach 2 — Persistent Buffer + Pointers (Optimal)

### Intuition
Observe how much can actually be left over: `read4` hands us at most 4 characters, so after serving any request the undelivered surplus is **at most 3 characters** — always from the *most recent* chunk. So a growable queue is overkill. Keep exactly three pieces of persistent state:

- `buf4 [4]byte` — the last chunk fetched from the file,
- `i4` — index of the next unconsumed character inside `buf4`,
- `n4` — how many characters in `buf4` are valid.

Every `read()` first drains `buf4[i4:n4]`, and only calls `read4` again when `i4 == n4`. Nothing is ever discarded, and space is constant.

### Algorithm
1. (State) `buf4`, `i4`, `n4` persist on the reader across calls; all start at zero.
2. On `read(buf, n)`: set `total = 0`.
3. While `total < n`:
   1. If `i4 == n4` (internal buffer drained): `n4 = read4(buf4)`, `i4 = 0`; if `n4 == 0` → EOF, break.
   2. Inner copy: while `total < n` **and** `i4 < n4`: `buf[total] = buf4[i4]`; `total++`; `i4++`.
4. Return `total`. (Whatever remains in `buf4[i4:n4]` is automatically the leftover for the next call.)

### Complexity
- **Time:** O(n) per call — each character moves file → `buf4` → `buf` exactly once over the reader's lifetime; at most ⌈n/4⌉ + 1 `read4` calls per request.
- **Space:** O(1) — a fixed 4-byte array and two integers, independent of `n`, query count, and file size.

### Code
```go
type pointerReader struct {
    read4 func([]byte) int
    buf4  [4]byte // last fetched chunk (persists across calls)
    i4    int     // next unconsumed index in buf4
    n4    int     // number of valid chars in buf4
}

func (r *pointerReader) read(buf []byte, n int) int {
    total := 0
    for total < n {
        if r.i4 == r.n4 { // buffer drained → refill
            r.n4 = r.read4(r.buf4[:])
            r.i4 = 0
            if r.n4 == 0 {
                break // EOF
            }
        }
        for total < n && r.i4 < r.n4 {
            buf[total] = r.buf4[r.i4]
            total++
            r.i4++
        }
    }
    return total
}
```

### Dry Run
Example 1: `file = "abc"`, queries `[1, 2, 1]`.

| Call | Step | i4 | n4 | buf4 (valid) | total | buf (this call) | Notes |
|------|------|----|----|---------------|-------|------------------|-------|
| `read(1)` | enter | 0 | 0 | — | 0 | — | i4 == n4 → refill |
| | refill | 0 | 3 | `"abc"` | 0 | — | read4 returned 3 |
| | copy | 1 | 3 | `"abc"` | 1 | `"a"` | inner loop stops: total == n |
| | return | 1 | 3 | `"bc"` left | — | — | **returns 1** ✓ |
| `read(2)` | enter | 1 | 3 | `"bc"` left | 0 | — | i4 ≠ n4 → no refill |
| | copy | 3 | 3 | — | 2 | `"bc"` | drained leftover exactly |
| | return | 3 | 3 | none left | — | — | **returns 2** ✓ |
| `read(1)` | enter | 3 | 3 | — | 0 | — | i4 == n4 → refill |
| | refill | 0 | 0 | — | 0 | — | read4 returned 0 → EOF, break |
| | return | 0 | 0 | — | — | — | **returns 0** ✓ |

Output `[1,2,0]` ✓

---

## Key Takeaways
- **The state between calls is the whole problem.** #157 lets you discard surplus characters; #158 makes discarding a correctness bug. Whenever an API says "may be called multiple times", ask: *what must survive between calls?*
- The surplus is bounded by the chunk size (≤ 3 chars here) — bounded leftovers mean you can replace a growable queue with a **fixed buffer + cursors**, dropping space to O(1).
- The `i4 == n4` guard ("drain before refill") is the invariant that guarantees no character is skipped or duplicated.
- Reset semantics matter: on LeetCode the Solution object may be reused, so leftover state must live on the instance (fields), not in globals.
- This is the textbook model of real buffered I/O — `fread`/`BufferedReader` wrap block-device reads exactly this way.

---

## Related Problems
- LeetCode #157 — Read N Characters Given Read4 (single-call version, no persistence needed)
- LeetCode #251 — Flatten 2D Vector (iterator with internal cursors across chunks)
- LeetCode #604 — Design Compressed String Iterator (stateful character streaming)
- LeetCode #900 — RLE Iterator (consume-k semantics with persistent leftover state)
