# 0388 — Longest Absolute File Path

> LeetCode #388 · Difficulty: Medium
> **Categories:** String, Stack, Depth-First Search

---

## Problem Statement

Suppose we have a file system that stores both files and directories. An example of one system is represented in the following picture:

```
dir
    subdir1
        file1.ext
        subsubdir1
    subdir2
        subsubdir2
            file2.ext
```

Here, we have `dir` as the only directory in the root. `dir` contains two subdirectories, `subdir1` and `subdir2`. `subdir1` contains a file `file1.ext` and an empty second-level subdirectory `subsubdir1`. `subdir2` contains a second-level subdirectory `subsubdir2` containing a file `file2.ext`.

In text form, it looks like this (with `⟶` representing the tab character):

```
dir
⟶ subdir1
⟶ ⟶ file1.ext
⟶ ⟶ subsubdir1
⟶ subdir2
⟶ ⟶ subsubdir2
⟶ ⟶ ⟶ file2.ext
```

If we were to write this representation in code, it will look like this: `"dir\n\tsubdir1\n\t\tfile1.ext\n\t\tsubsubdir1\n\tsubdir2\n\t\tsubsubdir2\n\t\t\tfile2.ext"`. Note that the `'\n'` and `'\t'` are the new-line and tab characters.

Every file and directory has a unique absolute path in the file system, which is the order of directories that must be opened to reach the file/directory itself, all concatenated by `'/'`s. Using the above example, the absolute path to `file2.ext` is `"dir/subdir2/subsubdir2/file2.ext"`. Each directory name consists of letters, digits, and/or spaces. Each file name is of the form `name.extension`, where `name` and `extension` consist of letters, digits, and/or spaces.

Given a string `input` representing the file system in the explained format, return the length of the longest absolute path to a file in the abstracted file system. If there is no file in the system, return `0`.

Note that the testcases are generated such that the file system is valid and no file or directory name has length `0`.

**Example 1:**

```
Input: input = "dir\n\tsubdir1\n\tsubdir2\n\t\tfile.ext"
Output: 20
Explanation: We have only one file, and the absolute path is "dir/subdir2/file.ext" of length 20.
```

**Example 2:**

```
Input: input = "dir\n\tsubdir1\n\t\tfile1.ext\n\t\tsubsubdir1\n\tsubdir2\n\t\tsubsubdir2\n\t\t\tfile2.ext"
Output: 32
Explanation: We have two files:
"dir/subdir1/file1.ext" of length 21
"dir/subdir2/subsubdir2/file2.ext" of length 32.
We return 32 since it is the longest absolute path to a file.
```

**Example 3:**

```
Input: input = "a"
Output: 0
Explanation: We do not have any files, just a single directory named "a".
```

**Constraints:**

- `1 <= input.length <= 10^4`
- `input` may contain lowercase or uppercase English letters, a new line character `'\n'`, a tab character `'\t'`, a dot `'.'`, a space `' '`, and digits.
- All file and directory names have positive length.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Stack** — the tab depth defines a nesting; a stack (or array indexed by depth) holds the running path length of each ancestor so a parent's length is an O(1) lookup → see [`/dsa/stack.md`](/dsa/stack.md)
- **String Processing** — split on `'\n'`, count leading `'\t'`, detect files by a `'.'` in the name → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Depth-First Traversal (implicit)** — the input is a preorder listing of a tree; we reconstruct depths without ever building the tree → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Stack of Path Lengths | O(N) | O(D) | Clean array-as-stack indexed by depth |
| 2 | Map depth → running length (Optimal) | O(N) | O(D) | Same accounting; explicit parent lookup |

(N = total characters, D = maximum directory depth.)

---

## Approach 1 — Stack of Path Lengths

### Intuition

The number of leading `'\t'` tabs of a line is its depth. The absolute path length of an entry at depth `d` equals its own name length plus its parent's accumulated length (at depth `d-1`) plus 1 for the `'/'`. A stack indexed by depth gives the parent length in O(1). A line containing a `'.'` is a file — update the best answer there.

### Algorithm

1. Split input on `'\n'` into lines.
2. For each line: `depth` = number of leading `'\t'`; `name` = line without those tabs.
3. `curLen = depth==0 ? len(name) : stack[depth-1] + 1 + len(name)`.
4. Set `stack[depth] = curLen`. If `name` contains `'.'`, it is a file → update `longest`.
5. Return `longest` (0 if no file was seen).

### Complexity

- **Time:** O(N) — each character is scanned a constant number of times.
- **Space:** O(D) — the stack holds one length per active depth.

### Code

```go
func stackLengths(input string) int {
	lines := strings.Split(input, "\n") // each token is one dir/file entry
	// stack[d] = cumulative path length of the current chain at depth d.
	stack := make([]int, len(lines)+1)
	longest := 0

	for _, line := range lines {
		// Count leading tabs → this is the entry's nesting depth.
		depth := 0
		for depth < len(line) && line[depth] == '\t' {
			depth++
		}
		name := line[depth:]      // strip the leading tabs to get the raw name
		nameLen := len(name)      // characters in this dir/file name

		curLen := nameLen // path length if this were a top-level entry
		if depth > 0 {
			// parent chain length + '/' + this name
			curLen = stack[depth-1] + 1 + nameLen
		}
		stack[depth] = curLen // record cumulative length at this depth

		// A '.' in the name marks a file (e.g. "file.ext"); measure it.
		if strings.Contains(name, ".") {
			if curLen > longest {
				longest = curLen
			}
		}
	}
	return longest
}
```

### Dry Run

Input `"dir\n\tsubdir1\n\tsubdir2\n\t\tfile.ext"` → lines: `dir`, `\tsubdir1`, `\tsubdir2`, `\t\tfile.ext`.

| line | depth | name | curLen | stack after | file? |
|------|-------|------|--------|-------------|-------|
| `dir` | 0 | `dir` | 3 | `[3]` | no |
| `\tsubdir1` | 1 | `subdir1` | 3+1+7 = 11 | `[3,11]` | no |
| `\tsubdir2` | 1 | `subdir2` | 3+1+7 = 11 | `[3,11]` | no |
| `\t\tfile.ext` | 2 | `file.ext` | 11+1+8 = 20 | `[3,11,20]` | yes → longest = 20 |

Answer: `20`.

---

## Approach 2 — Map depth → running length (Optimal)

### Intuition

Identical accounting, but expressed as a map `lengths[depth]` seeded with a virtual root `lengths[0] = 0`. Directories store their length *with* a trailing `'/'` so a child just adds its name; files add their name with no trailing slash and update the answer. Reads cleanly as "parent length + my name (+ slash)."

### Algorithm

1. Split on `'\n'`. Seed `lengths[0] = 0` (virtual root).
2. For each line: `depth` = leading tabs; `level = depth + 1`; `name` = stripped line.
3. If file (`'.'` in name): `answer = max(answer, lengths[level-1] + len(name))`.
4. Else (directory): `lengths[level] = lengths[level-1] + len(name) + 1`.
5. Return `answer`.

### Complexity

- **Time:** O(N) — a single pass over all characters.
- **Space:** O(D) — one map entry per active depth.

### Code

```go
func mapDepthLength(input string) int {
	lines := strings.Split(input, "\n")
	// lengths[d] = path length (including trailing '/') of the current dir at depth d.
	// lengths[0] = 0 acts as the virtual filesystem root.
	lengths := map[int]int{0: 0}
	answer := 0

	for _, line := range lines {
		depth := 0
		for depth < len(line) && line[depth] == '\t' {
			depth++
		}
		name := line[depth:]
		level := depth + 1 // shift so top-level entries sit at level 1

		if strings.Contains(name, ".") {
			// File: full path = parent dir path + file name (no trailing slash).
			pathLen := lengths[level-1] + len(name)
			if pathLen > answer {
				answer = pathLen
			}
		} else {
			// Directory: store its path length WITH a trailing '/' for children.
			lengths[level] = lengths[level-1] + len(name) + 1
		}
	}
	return answer
}
```

### Dry Run

Input `"dir\n\tsubdir1\n\tsubdir2\n\t\tfile.ext"` (`lengths[0]=0`):

| line | depth | level | name | file? | update |
|------|-------|-------|------|-------|--------|
| `dir` | 0 | 1 | `dir` | no | `lengths[1] = 0 + 3 + 1 = 4` |
| `\tsubdir1` | 1 | 2 | `subdir1` | no | `lengths[2] = 4 + 7 + 1 = 12` |
| `\tsubdir2` | 1 | 2 | `subdir2` | no | `lengths[2] = 4 + 7 + 1 = 12` |
| `\t\tfile.ext` | 2 | 3 | `file.ext` | yes | `answer = lengths[2] + 8 = 12 + 8 = 20` |

Answer: `20`.

---

## Key Takeaways

- **Tab count = tree depth.** An indented preorder listing encodes a tree; the indent level replaces explicit parent pointers.
- **Stack/array indexed by depth gives O(1) parent lookup.** Writing `stack[depth]` automatically "pops" deeper stale entries — they get overwritten when the next entry at that depth arrives.
- **Store directory lengths with the trailing `/`** so children add only their own name; files omit the slash and are the only entries that update the answer.
- No file ⇒ answer stays `0`, which handles the single-directory edge case for free.

---

## Related Problems

- LeetCode #71 — Simplify Path (path parsing with a stack)
- LeetCode #1592 — Rearrange Spaces Between Words (string tokenizing)
- LeetCode #394 — Decode String (depth via a stack)
