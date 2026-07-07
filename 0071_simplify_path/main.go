package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Stack-Based Path Processing ───────────────────────────────────
//
// simplifyPath solves Simplify Path using a stack to process each component.
//
// Intuition:
//   Split the path by '/'. For each component:
//   - "" or ".": skip (empty segment or current directory).
//   - "..": pop the stack (go up one level), if non-empty.
//   - anything else: push onto the stack (a directory name).
//   Join the stack with '/' and prepend '/'.
//
// Algorithm:
//   parts = split(path, "/")
//   stack = []
//   for part in parts:
//     if part == "" or part == ".": continue
//     elif part == "..": if stack non-empty: pop
//     else: push part
//   return "/" + join(stack, "/")
//
// Time:  O(n) — n = len(path).
// Space: O(n) — stack size at most n/2 components.
func simplifyPath(path string) string {
	parts := strings.Split(path, "/")
	stack := []string{}

	for _, part := range parts {
		switch part {
		case "", ".":
			// skip: empty segment (double slash) or current directory
		case "..":
			if len(stack) > 0 {
				stack = stack[:len(stack)-1] // go up one level
			}
		default:
			stack = append(stack, part) // valid directory name
		}
	}

	return "/" + strings.Join(stack, "/")
}

func main() {
	fmt.Println("=== Simplify Path ===")
	fmt.Printf("%q  got=%q  expected %q\n", "/home/", simplifyPath("/home/"), "/home")
	fmt.Printf("%q  got=%q  expected %q\n", "/home//foo/", simplifyPath("/home//foo/"), "/home/foo")
	fmt.Printf("%q  got=%q  expected %q\n", "/.../a/../b/c/../d/./", simplifyPath("/.../a/../b/c/../d/./"), "/.../b/d")
	fmt.Printf("%q  got=%q  expected %q\n", "/a/./b/../../c/", simplifyPath("/a/./b/../../c/"), "/c")
	fmt.Printf("%q  got=%q  expected %q\n", "/../", simplifyPath("/../"), "/")
	fmt.Printf("%q  got=%q  expected %q\n", "/a//b////c/d//././/..", simplifyPath("/a//b////c/d//././/.."), "/a/b/c")
}
