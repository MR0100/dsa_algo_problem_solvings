package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Node is the N-ary tree node used by LeetCode #428: a value plus a slice of
// children (arbitrary fan-out).
type Node struct {
	Val      int
	Children []*Node
}

// ── Approach 1: Preorder with Explicit Child Counts ──────────────────────────
//
// Codec1 serializes an N-ary tree as a preorder stream where each node is
// written as "val,childCount", so deserialization knows exactly how many
// children to consume without any sentinel.
//
// Intuition:
//
//	The reason serializing an N-ary tree is harder than a binary tree is that the
//	number of children is not fixed, so a reader cannot tell where one node's
//	children end. Fix that by writing the child count next to each value. A
//	preorder walk then reconstructs uniquely: read a value and its count k, then
//	recursively read exactly k children.
//
// Algorithm (serialize):
//  1. Preorder: emit "Val,len(Children)" for the node.
//  2. Recurse into each child in order.
//  3. Join all tokens with commas.
//
// Algorithm (deserialize):
//  1. Split into tokens; walk them left to right with a shared cursor.
//  2. Read Val, then k = child count; build the node.
//  3. Recurse k times to build children; attach them.
//
// Time:  O(n) serialize and deserialize — each node emitted/consumed once.
// Space: O(n) for the string plus O(h) recursion depth.
type Codec1 struct{}

// serialize turns the tree into "v0,c0,v1,c1,..." preorder token stream.
func (Codec1) serialize(root *Node) string {
	if root == nil {
		return "" // empty tree ⇒ empty string
	}
	tokens := []string{}
	var pre func(n *Node)
	pre = func(n *Node) {
		// Each node contributes its value AND its child count so the reader
		// knows how many children to pull next.
		tokens = append(tokens, strconv.Itoa(n.Val), strconv.Itoa(len(n.Children)))
		for _, c := range n.Children {
			pre(c)
		}
	}
	pre(root)
	return strings.Join(tokens, ",")
}

// deserialize rebuilds the tree from the "v,c,v,c,..." token stream.
func (Codec1) deserialize(data string) *Node {
	if data == "" {
		return nil
	}
	tokens := strings.Split(data, ",")
	pos := 0 // shared cursor into tokens
	var build func() *Node
	build = func() *Node {
		val, _ := strconv.Atoi(tokens[pos])     // current node's value
		count, _ := strconv.Atoi(tokens[pos+1]) // its child count
		pos += 2                                // consume the (val,count) pair
		node := &Node{Val: val, Children: []*Node{}}
		for i := 0; i < count; i++ {
			node.Children = append(node.Children, build()) // pull exactly `count` children
		}
		return node
	}
	return build()
}

// ── Approach 2: Bracketed / Parenthesized Encoding ───────────────────────────
//
// Codec2 serializes using explicit brackets to delimit each node's children,
// e.g. "1[3[5 6]2 4]". Structure is carried by the brackets rather than by a
// count, mirroring how nested data (JSON/S-expressions) is written.
//
// Intuition:
//
//	Another way to make an unbounded child list unambiguous is to bracket it: a
//	node prints its value, then "[" its children "]" if it has any. The matching
//	"]" tells the reader where this node's children end, no counting required.
//	This is the classic S-expression / DFS-with-delimiters idea.
//
// Algorithm (serialize):
//  1. Emit the node's value.
//  2. If it has children, emit "[", recurse into each child (space-separated),
//     then "]".
//
// Algorithm (deserialize):
//  1. Tokenize into numbers, "[", "]".
//  2. Read a number → new node. If the next token is "[", recurse to read
//     children until the matching "]".
//
// Time:  O(n) serialize and deserialize.
// Space: O(n) string plus O(h) recursion depth.
type Codec2 struct{}

// serialize produces "val" or "val[child child ...]" recursively.
func (Codec2) serialize(root *Node) string {
	if root == nil {
		return ""
	}
	var sb strings.Builder
	var enc func(n *Node)
	enc = func(n *Node) {
		sb.WriteString(strconv.Itoa(n.Val)) // node value
		if len(n.Children) > 0 {
			sb.WriteByte('[') // open the child group
			for i, c := range n.Children {
				if i > 0 {
					sb.WriteByte(' ') // separate siblings
				}
				enc(c)
			}
			sb.WriteByte(']') // close the child group
		}
	}
	enc(root)
	return sb.String()
}

// deserialize parses the bracketed string back into a tree.
func (Codec2) deserialize(data string) *Node {
	if data == "" {
		return nil
	}
	tokens := tokenizeBrackets(data)
	pos := 0
	var build func() *Node
	build = func() *Node {
		val, _ := strconv.Atoi(tokens[pos]) // a number token starts a node
		pos++
		node := &Node{Val: val, Children: []*Node{}}
		if pos < len(tokens) && tokens[pos] == "[" {
			pos++ // consume "["
			for tokens[pos] != "]" {
				node.Children = append(node.Children, build()) // read one child
			}
			pos++ // consume the matching "]"
		}
		return node
	}
	return build()
}

// tokenizeBrackets splits a bracketed string into number / "[" / "]" tokens.
func tokenizeBrackets(s string) []string {
	tokens := []string{}
	i := 0
	for i < len(s) {
		switch s[i] {
		case '[', ']':
			tokens = append(tokens, string(s[i])) // structural token
			i++
		case ' ':
			i++ // skip separators
		default:
			j := i
			for j < len(s) && s[j] != '[' && s[j] != ']' && s[j] != ' ' {
				j++ // extend over the whole number (supports multi-digit / negatives)
			}
			tokens = append(tokens, s[i:j])
			i = j
		}
	}
	return tokens
}

// ── verification helpers ─────────────────────────────────────────────────────

// levelOrder returns the BFS level groups of a tree so we can confirm a
// round-trip rebuilt the same structure (LeetCode's #429 output form).
func levelOrder(root *Node) [][]int {
	res := [][]int{}
	if root == nil {
		return res
	}
	q := []*Node{root}
	for len(q) > 0 {
		w := len(q)
		lvl := []int{}
		for i := 0; i < w; i++ {
			n := q[0]
			q = q[1:]
			lvl = append(lvl, n.Val)
			q = append(q, n.Children...)
		}
		res = append(res, lvl)
	}
	return res
}

// buildExample constructs the standard 3-ary example tree:
//
//	     1
//	   / | \
//	  3  2  4
//	 / \
//	5   6
func buildExample() *Node {
	n5 := &Node{Val: 5}
	n6 := &Node{Val: 6}
	n3 := &Node{Val: 3, Children: []*Node{n5, n6}}
	n2 := &Node{Val: 2}
	n4 := &Node{Val: 4}
	return &Node{Val: 1, Children: []*Node{n3, n2, n4}}
}

func main() {
	root := buildExample()
	want := levelOrder(root) // [[1] [3 2 4] [5 6]]

	fmt.Println("=== Approach 1: Preorder + Child Counts ===")
	c1 := Codec1{}
	s1 := c1.serialize(root)
	fmt.Println("serialized:", s1) // 1,3,3,2,5,0,6,0,2,0,4,0
	r1 := c1.deserialize(s1)
	fmt.Println("round-trip level order:", levelOrder(r1))                             // [[1] [3 2 4] [5 6]]
	fmt.Println("round-trip matches:", fmt.Sprint(levelOrder(r1)) == fmt.Sprint(want)) // true

	fmt.Println("=== Approach 2: Bracketed Encoding ===")
	c2 := Codec2{}
	s2 := c2.serialize(root)
	fmt.Println("serialized:", s2) // 1[3[5 6] 2 4]
	r2 := c2.deserialize(s2)
	fmt.Println("round-trip level order:", levelOrder(r2))                             // [[1] [3 2 4] [5 6]]
	fmt.Println("round-trip matches:", fmt.Sprint(levelOrder(r2)) == fmt.Sprint(want)) // true

	fmt.Println("=== Edge case: empty tree ===")
	fmt.Printf("Codec1 empty round-trips to nil: %v\n", c1.deserialize(c1.serialize(nil)) == nil) // true
	fmt.Printf("Codec2 empty round-trips to nil: %v\n", c2.deserialize(c2.serialize(nil)) == nil) // true

	fmt.Println("=== Edge case: single node (val 42) ===")
	single := &Node{Val: 42}
	fmt.Println("Codec1:", c1.serialize(single), "->", levelOrder(c1.deserialize(c1.serialize(single)))) // 42,0 -> [[42]]
	fmt.Println("Codec2:", c2.serialize(single), "->", levelOrder(c2.deserialize(c2.serialize(single)))) // 42 -> [[42]]
}
