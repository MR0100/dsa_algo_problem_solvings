# 0071 — Simplify Path

> LeetCode #71 · Difficulty: Medium
> **Categories:** String, Stack

---

## Problem Statement

Given an absolute path for a Unix-style file system, which begins with a slash `'/'`, transform this path into its **simplified canonical path**.

In the Unix-style file system context, a single period `'.'` signifies the current directory, a double period `".."` denotes moving up one level in the directory hierarchy, and any other names denote directories or files.

The simplified canonical path should follow these rules:
- The path must start with a single slash `'/'`.
- Any two directories within the path must be separated by exactly one slash `'/'`.
- The path must not end with a trailing `'/'`, unless it is the root directory `"/"`.
- The path must not have any single or double periods `'.'` or `".."` used to denote current or parent directories.

**Example 1**
```
Input:  path = "/home/"
Output: "/home"
```

**Example 2**
```
Input:  path = "/home//foo/"
Output: "/home/foo"
```

**Example 3**
```
Input:  path = "/.../a/../b/c/../d/./"
Output: "/.../b/d"
Explanation: "..." is a valid directory name, not a special token.
```

**Constraints**
- `1 <= path.length <= 3000`
- `path` consists of English letters, digits, period `'.'`, slash `'/'` or `'_'`.
- `path` is a valid absolute Unix path.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Stack** — push valid directory names; pop on `..`.
- **String Split** — split by `/` to get individual path components.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Stack-Based Path Processing ✅ | O(n) | O(n) | The canonical solution |

---

## Approach 1 — Stack-Based Path Processing (Recommended ✅)

### Intuition
Split the path by `'/'`. This gives us path components, some of which may be empty (from consecutive slashes). For each component:
- `""` or `"."`: skip (no-op).
- `".."`: pop the stack (go up), if non-empty (can't go above root).
- anything else: push onto the stack (enter directory).

Join the stack with `'/'` and prepend `'/'`. If the stack is empty, return `"/"`.

### Algorithm
```
parts = split(path, "/")
stack = []
for part in parts:
  if part == "" or part == ".": skip
  elif part == "..": if stack non-empty: pop
  else: push(part)
return "/" + join(stack, "/")
```

### Complexity
- **Time:** O(n) — split + iterate all components.
- **Space:** O(n) — stack stores at most n/2 components.

### Code
```go
func simplifyPath(path string) string {
    parts := strings.Split(path, "/")
    stack := []string{}
    for _, part := range parts {
        switch part {
        case "", ".":   // skip
        case "..": if len(stack) > 0 { stack = stack[:len(stack)-1] }
        default: stack = append(stack, part)
        }
    }
    return "/" + strings.Join(stack, "/")
}
```

### Dry Run — `path = "/a/./b/../../c/"`
```
Split: ["", "a", ".", "b", "..", "..", "c", ""]

"": skip. stack=[]
"a": push. stack=["a"]
".": skip. stack=["a"]
"b": push. stack=["a","b"]
"..": pop. stack=["a"]
"..": pop. stack=[]
"c": push. stack=["c"]
"": skip.

Result: "/" + "c" = "/c" ✓
```

---

## Key Takeaways

- **Three-way switch** — empty/dot (skip), `..` (pop), other (push). All other cases are just valid directory names, even `"..."`.
- **`".."` at root doesn't error** — can't go above root; just ignore (guard with `len(stack) > 0`).
- **`strings.Split(path, "/")` handles double slashes** — `"/a//b"` splits to `["","a","","b"]`; the empty strings are naturally skipped.
- **`"..."` is valid** — three or more dots are a valid directory name, NOT a special token.

---

## Related Problems

- LeetCode #150 — Evaluate Reverse Polish Notation (stack-based expression evaluation)
- LeetCode #20 — Valid Parentheses (stack-based matching)
- LeetCode #388 — Longest Absolute File Path (directory traversal with stack)
