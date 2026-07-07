package main

import (
	"fmt"
	"strings"
)

// Node is the N-ary tree node used by LeetCode #431: a value plus a slice of
// child pointers (a node may have anywhere from 0 to N children).
type Node struct {
	Val      int
	Children []*Node
}

// TreeNode is the ordinary binary tree node we encode INTO: a value and exactly
// two child pointers, Left and Right.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// Codec is the interface every approach implements so main() can round-trip the
// same N-ary tree through each encode/decode strategy uniformly.
//
// The API is stateless by contract (LeetCode forbids storing tree state in the
// struct): encode takes an N-ary root and returns a binary root; decode does the
// reverse. The struct itself carries no per-tree fields.
type Codec interface {
	encode(root *Node) *TreeNode
	decode(root *TreeNode) *Node
}

// ── Approach 1: Left-Child / Right-Sibling (Optimal) ─────────────────────────
//
// LCRSCodec encodes an N-ary tree into a binary tree with the classic
// "left-child, right-sibling" (LCRS) mapping — the textbook, minimal-overhead
// answer.
//
// Intuition:
//
//	A binary node has two pointers; an N-ary node needs one pointer to its first
//	child and one to its "next sibling". So repurpose them:
//	  • binary Left  = the N-ary node's FIRST child
//	  • binary Right = the N-ary node's NEXT sibling
//	All the children of one N-ary node become a right-linked chain hanging off
//	the Left pointer. This is a bijection, so decoding is exact: follow Left to
//	recover the first child, then walk the Right chain to recover the rest.
//
// Algorithm:
//
//	encode(nary):
//	  make a binary node b with the same value.
//	  if nary has children:
//	    b.Left = encode(child[0]).
//	    walk cur = b.Left; for each remaining child i≥1:
//	      cur.Right = encode(child[i]); cur = cur.Right.
//	  return b.
//	decode(bin):
//	  make an N-ary node n with the same value.
//	  walk child = bin.Left along Right pointers:
//	    append decode(child) to n.Children; child = child.Right.
//	  return n.
//
// Time:  O(V) — encode and decode each visit every node exactly once.
// Space: O(H) recursion (H = N-ary tree height) plus O(V) for the output tree.
type LCRSCodec struct{}

// encode builds the binary tree using the left-child/right-sibling mapping.
func (LCRSCodec) encode(root *Node) *TreeNode {
	if root == nil {
		return nil // empty N-ary tree ↔ empty binary tree
	}
	b := &TreeNode{Val: root.Val} // binary node carrying the same value
	if len(root.Children) > 0 {
		// The FIRST child hangs off the Left pointer.
		b.Left = (LCRSCodec{}).encode(root.Children[0])
		cur := b.Left // cur walks the right-linked sibling chain
		// Every subsequent child is chained via Right (sibling links).
		for i := 1; i < len(root.Children); i++ {
			cur.Right = (LCRSCodec{}).encode(root.Children[i])
			cur = cur.Right // advance to the newly attached sibling
		}
	}
	return b
}

// decode rebuilds the N-ary tree from the left-child/right-sibling binary tree.
func (LCRSCodec) decode(root *TreeNode) *Node {
	if root == nil {
		return nil // empty binary tree ↔ empty N-ary tree
	}
	n := &Node{Val: root.Val, Children: []*Node{}} // N-ary node, same value
	child := root.Left                             // Left points at the first child
	// Walk the sibling chain: each Right hop is the next child of `n`.
	for child != nil {
		n.Children = append(n.Children, (LCRSCodec{}).decode(child))
		child = child.Right // move to the next sibling in the chain
	}
	return n
}

// ── Approach 2: BFS Level-Order String Serialization ─────────────────────────
//
// SerializeCodec ignores structural cleverness and instead serialises the N-ary
// tree into a flat string, stores that string inside a chain of binary nodes
// (one character code per node), then parses it back on decode.
//
// Intuition:
//
//	"Encode to a binary tree" only requires that the binary tree hold enough
//	information to rebuild the original. A serialised string already does that,
//	so we can smuggle the string through the binary tree: encode each token as a
//	binary node whose Val is the token and whose Left points to the next token.
//	This is deliberately NOT the intended trick — it shows the problem is really
//	just "serialise then deserialise", and highlights why the LCRS mapping
//	(Approach 1) is preferable: no string parsing, no delimiter bookkeeping.
//
// Serialisation format (per node, BFS/level order):
//
//	"<val> <childCount> " repeated. Root first, then each node's children in
//	order, so a queue reconstructs parents before children.
//
// Algorithm:
//
//	encode: BFS the N-ary tree, appending "<val> <#children> " per node;
//	        thread the whitespace-split tokens through Left pointers of binary
//	        nodes (Val = token's integer where numeric; counts stored too).
//	decode: read the tokens back off the Left chain, then BFS-rebuild: pop a
//	        (value,count) pair, create the node, enqueue it expecting `count`
//	        children which are the next popped nodes.
//
// Time:  O(V) — each node produces O(1) tokens; parsing is linear.
// Space: O(V) — the token chain and the queues.
type SerializeCodec struct{}

// encode serialises the N-ary tree (BFS) and threads the integer tokens through
// a Left-linked chain of binary nodes. Even indices hold a value, odd indices
// hold that value's child-count — a simple self-describing stream.
func (SerializeCodec) encode(root *Node) *TreeNode {
	if root == nil {
		return nil
	}
	tokens := []int{}      // flat stream: val, count, val, count, ...
	queue := []*Node{root} // BFS frontier over the N-ary tree
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]                                     // dequeue
		tokens = append(tokens, node.Val, len(node.Children)) // emit (val,count)
		queue = append(queue, node.Children...)               // children explored next
	}
	// Thread tokens through a Left-only binary chain (each node = one token).
	dummy := &TreeNode{}
	cur := dummy
	for _, t := range tokens {
		cur.Left = &TreeNode{Val: t} // store token in a binary node's value
		cur = cur.Left               // extend the chain downward via Left
	}
	return dummy.Left // first real token node is the encoded root
}

// decode reads the token chain back and BFS-rebuilds the N-ary tree.
func (SerializeCodec) decode(root *TreeNode) *Node {
	if root == nil {
		return nil
	}
	tokens := []int{} // recover the flat (val,count) stream
	for cur := root; cur != nil; cur = cur.Left {
		tokens = append(tokens, cur.Val)
	}
	idx := 0 // read cursor into tokens
	// Pop the root's (val,count); create it; queue it with how many children
	// it still expects. Each queued node consumes the next `count` popped nodes.
	rootVal, rootCnt := tokens[idx], tokens[idx+1]
	idx += 2
	nRoot := &Node{Val: rootVal, Children: []*Node{}}
	type item struct {
		node      *Node
		remaining int // children still to attach to this node
	}
	queue := []item{{nRoot, rootCnt}}
	for len(queue) > 0 && idx < len(tokens) {
		parent := &queue[0]
		if parent.remaining == 0 {
			queue = queue[1:] // this parent is satisfied; move on
			continue
		}
		// Next token pair is one child of the current parent.
		val, cnt := tokens[idx], tokens[idx+1]
		idx += 2
		child := &Node{Val: val, Children: []*Node{}}
		parent.node.Children = append(parent.node.Children, child)
		parent.remaining--                      // one fewer child to place for this parent
		queue = append(queue, item{child, cnt}) // the child may have its own kids
	}
	return nRoot
}

// naryToString renders an N-ary tree as a canonical nested string so main() can
// prove encode∘decode is the identity (input string == round-tripped string).
func naryToString(n *Node) string {
	if n == nil {
		return "nil"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d", n.Val))
	if len(n.Children) > 0 {
		sb.WriteString("[")
		for i, c := range n.Children {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(naryToString(c)) // recurse into each child
		}
		sb.WriteString("]")
	}
	return sb.String()
}

// buildSampleTree constructs the canonical LeetCode N-ary example tree:
//
//	      1
//	   /  |  \
//	  3   2   4
//	 / \
//	5   6
//
// serialised by LeetCode as [1,null,3,2,4,null,5,6].
func buildSampleTree() *Node {
	five := &Node{Val: 5, Children: []*Node{}}
	six := &Node{Val: 6, Children: []*Node{}}
	three := &Node{Val: 3, Children: []*Node{five, six}}
	two := &Node{Val: 2, Children: []*Node{}}
	four := &Node{Val: 4, Children: []*Node{}}
	return &Node{Val: 1, Children: []*Node{three, two, four}}
}

// roundTrip encodes then decodes a tree through a codec and returns the decoded
// tree's canonical string — equal to the input's string iff the codec is a
// faithful bijection.
func roundTrip(c Codec, root *Node) string {
	encoded := c.encode(root)    // N-ary → binary
	decoded := c.decode(encoded) // binary → N-ary
	return naryToString(decoded)
}

func main() {
	tree := buildSampleTree()
	original := naryToString(tree)
	fmt.Println("Original N-ary tree:", original) // 1[3[5,6],2,4]

	fmt.Println("=== Approach 1: Left-Child / Right-Sibling (Optimal) ===")
	fmt.Println(roundTrip(LCRSCodec{}, tree)) // 1[3[5,6],2,4]

	fmt.Println("=== Approach 2: BFS Level-Order String Serialization ===")
	fmt.Println(roundTrip(SerializeCodec{}, tree)) // 1[3[5,6],2,4]

	// Edge case: a single-node tree round-trips through both codecs.
	single := &Node{Val: 7, Children: []*Node{}}
	fmt.Println("=== Edge: single node (both codecs) ===")
	fmt.Println(roundTrip(LCRSCodec{}, single))      // 7
	fmt.Println(roundTrip(SerializeCodec{}, single)) // 7
}
