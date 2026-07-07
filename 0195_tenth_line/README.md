# 0195 — Tenth Line

> LeetCode #195 · Difficulty: Easy
> **Categories:** Shell, String

---

## Problem Statement

Given a text file `file.txt`, print just the 10th line of the file.

**Example:**

Assume that `file.txt` has the following content:

```
Line 1
Line 2
Line 3
Line 4
Line 5
Line 6
Line 7
Line 8
Line 9
Line 10
```

Your script should output the tenth line:

```
Line 10
```

**Note:**

1. If the file contains less than 10 lines, what should you output?
2. There's at least three different solutions. Try to explore all possibilities.

> **Repo note:** #195 is one of LeetCode's four Shell problems (the official
> ask is a bash script; the note's answer to question 1 is "print nothing").
> Following this repo's Go-only convention, every approach below re-implements
> the logic in Go, treating the file's content as an in-memory string
> (`file.txt` → the `fileTxt` constant) and returning `""` for short files.
> `main()` also exercises a 3-line file to demonstrate the edge case. The
> three canonical bash answers appear in Key Takeaways.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★☆☆☆☆ Rare       | 2022          |
| Bloomberg  | ★☆☆☆☆ Rare       | 2022          |
| Uber       | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String Algorithms (line scanning / early exit)** — the entire task is locating the substring between the 9th and 10th newline; the three approaches trade memory for streaming and allocation-freedom → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (materialise every line) | O(N) | O(N) | Small files; shortest code |
| 2 | Line Stream with Early Exit | O(P) (first 10 lines only) | O(L) | Big files; stop reading the moment line 10 appears |
| 3 | Raw Byte Scan, Zero Allocation (Optimal) | O(P) | O(1) | Hot paths; no scanner, no line slice, no copies |

*N = total characters, P = characters in the first 10 lines, L = longest line length.*

---

## Approach 1 — Brute Force (Materialise Every Line)

### Intuition
The blunt instrument: load everything, then array-index. `lines[9]` is the 10th line (0-based indexing), and a length check answers the note's question 1 — a file with fewer than 10 lines outputs nothing (`""`). This is the spirit of `head -n 10 file.txt | tail -n +10`: process the whole prefix, pick the last piece.

### Algorithm
1. Split the whole content on `'\n'` into a slice of lines.
2. If the slice has fewer than 10 entries, return `""` (print nothing).
3. Otherwise return `lines[9]`.

### Complexity
- **Time:** O(N) — the split walks every character of the file, even the ones after line 10.
- **Space:** O(N) — every line is materialised in the slice although only one is needed.

### Code
```go
func bruteForce(fileContent string) string {
	lines := strings.Split(fileContent, "\n") // materialise ALL lines up front
	if len(lines) < 10 {                      // fewer than 10 lines → nothing to print
		return ""
	}
	return lines[9] // 10th line lives at 0-based index 9
}
```

### Dry Run
Example: the 10-line file.

| step | variable state | result |
|------|----------------|--------|
| 1 | `lines = [Line 1, Line 2, Line 3, Line 4, Line 5, Line 6, Line 7, Line 8, Line 9, Line 10]`, `len = 10` | |
| 2 | `len(lines) < 10` → `10 < 10` = false | keep going |
| 3 | return `lines[9]` | `Line 10` ✓ |

Edge case (3-line file): `lines = [Line 1, Line 2, Line 3]`, `len = 3 < 10` → return `""` (print nothing). ✓

---

## Approach 2 — Line Stream with Early Exit

### Intuition
Only one line matters, so reading past it is pure waste. A streaming scanner keeps just the current line in memory and abandons the input the moment the target is produced — exactly why `sed -n '10p;10q'` beats plain `sed -n '10p'` on a gigantic file: the `10q` quits instead of scanning to EOF.

### Algorithm
1. Wrap the content in a `bufio.Scanner` (default line-splitting mode).
2. Scan lines, incrementing a counter per line.
3. When the counter hits 10, return that line immediately — the rest of the input is never read.
4. If the input ends first, return `""`.

### Complexity
- **Time:** O(P) — only the characters of the first 10 lines are consumed (early exit); O(N) worst case when the file is shorter than 10 lines.
- **Space:** O(L) — the scanner buffers a single line at a time.

### Code
```go
func streamEarlyExit(fileContent string) string {
	scanner := bufio.NewScanner(strings.NewReader(fileContent)) // stream, don't slurp
	lineNo := 0
	for scanner.Scan() { // pulls ONE line per iteration
		lineNo++
		if lineNo == 10 {
			return scanner.Text() // found it — stop reading, skip the rest of the file
		}
	}
	return "" // ran out of lines before reaching 10
}
```

### Dry Run
Example: the 10-line file.

| iteration | scanner.Text() | lineNo after | lineNo == 10? |
|-----------|----------------|--------------|----------------|
| 1 | `Line 1` | 1 | no |
| 2 | `Line 2` | 2 | no |
| 3–8 | `Line 3` … `Line 8` | 3–8 | no |
| 9 | `Line 9` | 9 | no |
| 10 | `Line 10` | 10 | **yes → return `Line 10`** ✓ |

Edge case (3-line file): iterations 1–3 consume all lines, `Scan()` returns false, fall through → return `""`. ✓

---

## Approach 3 — Raw Byte Scan, Zero Allocation (Optimal)

### Intuition
Lines are just the gaps between newline bytes: the 10th line starts right after the 9th `'\n'` and ends right before the 10th `'\n'` (or at end-of-input). So track only two integers — newlines seen, and where the 10th line starts — and slice the original string once. Go string slicing shares the underlying bytes, so this performs zero allocations: no scanner, no line slice, no copies.

### Algorithm
1. Walk the bytes with a `newlines` counter and a `start` marker.
2. On the 9th `'\n'`, record `start = i + 1` (the 10th line begins here).
3. On the 10th `'\n'`, return `content[start:i]` and stop — early exit.
4. If the walk ends with exactly 9 newlines seen, the 10th line is the final unterminated line: return `content[start:]`.
5. Otherwise (fewer than 9 newlines) return `""`.

### Complexity
- **Time:** O(P) — bytes up to the 10th newline only; O(N) worst case for short files.
- **Space:** O(1) — two integers; the returned slice aliases the input's memory (no allocation).

### Code
```go
func byteScan(fileContent string) string {
	newlines := 0 // how many '\n' bytes seen so far
	start := 0    // index just past the 9th newline = start of the 10th line
	for i := 0; i < len(fileContent); i++ {
		if fileContent[i] != '\n' {
			continue // only newline bytes drive the state machine
		}
		newlines++
		if newlines == 9 {
			start = i + 1 // the 10th line begins immediately after the 9th '\n'
		}
		if newlines == 10 {
			return fileContent[start:i] // ends just before the 10th '\n'; early exit
		}
	}
	if newlines == 9 {
		return fileContent[start:] // 10th line is the last line, no trailing '\n'
	}
	return "" // fewer than 10 lines in the file
}
```

### Dry Run
Example: the 10-line file (`"Line 1\nLine 2\n...\nLine 10"`, each `Line k\n` is 7 bytes; total length 69, no trailing newline).

| event | i (byte index of `'\n'`) | newlines after | start | action |
|-------|--------------------------|----------------|-------|--------|
| 1st `\n` (after `Line 1`) | 6 | 1 | 0 | keep walking |
| 2nd–8th `\n` | 13, 20, 27, 34, 41, 48 | 2–7 | 0 | keep walking |
| 8th `\n` (after `Line 8`) | 55 | 8 | 0 | keep walking |
| 9th `\n` (after `Line 9`) | 62 | 9 | **63** | 10th line starts at 63 |
| loop ends (no 10th `\n`) | — | 9 | 63 | `newlines == 9` → return `content[63:]` |

`content[63:]` = `Line 10`. ✓

Edge case (3-line file): only 2 newlines ever seen → falls to the final `return ""`. ✓

---

## Key Takeaways

- **k-th line = the slice between newline k−1 and newline k** — thinking in delimiters instead of "lines" removes all allocation from the problem.
- **Early exit is a real complexity win on streams**: O(first-k-lines) vs O(whole-file) is the difference between `sed -n '10p;10q'` and `sed -n '10p'` — say the `;10q` part in interviews.
- **Always answer the short-input question before coding** (here: < 10 lines → print nothing); index `[9]` without the length guard is the classic panic.
- **Go string slices share memory** — `s[a:b]` is O(1) and allocation-free, ideal for extract-a-substring tasks; copy only if the giant source must be garbage-collected.
- The three canonical bash answers the note asks you to explore, verbatim:
  ```bash
  sed -n '10p;10q' file.txt          # stream editor: print line 10, then quit
  awk 'NR == 10 {print; exit}' file.txt   # awk: row number 10, then exit
  tail -n +10 file.txt | head -n 1   # tail from line 10, keep the first
  ```
  All three print nothing when the file has fewer than 10 lines.

---

## Related Problems

- LeetCode #192 — Word Frequency (Shell problem set)
- LeetCode #193 — Valid Phone Numbers (Shell problem set)
- LeetCode #194 — Transpose File (Shell problem set)
- LeetCode #19 — Remove Nth Node From End of List (locate the k-th element with a single bounded pass)
- LeetCode #876 — Middle of the Linked List (positional element extraction from a stream you cannot index)
