package main

import "fmt"

// Iterator is the given underlying iterator (as LeetCode provides it): it
// exposes HasNext and Next only — no peek. We model it over a slice.
type Iterator struct {
	data []int // backing sequence
	pos  int   // index of the next element to return
}

// NewIterator builds the base iterator over a slice.
func NewIterator(data []int) *Iterator {
	return &Iterator{data: data}
}

// HasNext reports whether the base iterator has more elements.
func (it *Iterator) HasNext() bool { return it.pos < len(it.data) }

// Next returns and consumes the next element of the base iterator.
func (it *Iterator) Next() int {
	v := it.data[it.pos]
	it.pos++
	return v
}

// peekingADT is the interface all approaches expose so main() can drive the
// same official operation sequence through each.
type peekingADT interface {
	Peek() int
	Next() int
	HasNext() bool
}

// ── Approach 1: Cache One Element (Eager Lookahead, Optimal) ─────────────────
//
// CachePeeking wraps the base iterator and always keeps the upcoming element
// pre-fetched in a field, so Peek returns it without consuming the base.
//
// Intuition:
//
//	The base iterator can only move forward and can't un-Next. To support peek,
//	we pull one element ahead into a buffer at construction. Peek just reads
//	the buffer; Next returns the buffer then refills it from the base (if any).
//	HasNext is true whenever we still hold a buffered element. One element of
//	extra state, all ops O(1).
//
// Algorithm:
//
//	Constructor: if base.HasNext, next = base.Next(), hasPeeked = true.
//	Peek:        return next (the cached element).
//	Next:        v = next; refill: if base.HasNext, next = base.Next() else
//	             hasPeeked = false; return v.
//	HasNext:     return hasPeeked.
//
// Time:  Peek O(1), Next O(1), HasNext O(1).
// Space: O(1) — one cached value plus a flag.
type CachePeeking struct {
	base      *Iterator // underlying forward-only iterator
	next      int       // pre-fetched upcoming element
	hasPeeked bool      // true while `next` holds a valid element
}

// NewCachePeeking builds the wrapper, eagerly fetching the first element.
func NewCachePeeking(it *Iterator) *CachePeeking {
	p := &CachePeeking{base: it}
	if it.HasNext() {
		p.next = it.Next() // pull the first element into the buffer
		p.hasPeeked = true
	}
	return p
}

// Peek returns the next element without advancing.
func (p *CachePeeking) Peek() int { return p.next }

// Next returns the buffered element and refills the buffer from the base.
func (p *CachePeeking) Next() int {
	v := p.next // element to hand out
	if p.base.HasNext() {
		p.next = p.base.Next() // refill lookahead
	} else {
		p.hasPeeked = false // base exhausted → buffer now empty
	}
	return v
}

// HasNext is true while a buffered element remains.
func (p *CachePeeking) HasNext() bool { return p.hasPeeked }

// ── Approach 2: Lazy Peek (Fetch Only When Asked) ────────────────────────────
//
// LazyPeeking defers the lookahead: it only pulls the next element from the
// base the first time Peek is called, tracking whether a peeked value is held.
//
// Intuition:
//
//	Instead of always buffering one ahead, buffer lazily. Keep a flag
//	`peeked`. Peek: if not yet peeked, consume one from the base and stash it,
//	set peeked=true; return the stash. Next: if peeked, return the stash and
//	clear the flag; otherwise delegate to base.Next(). HasNext: true if we
//	hold a peeked value OR the base still has one. Same O(1) costs, avoids the
//	constructor-time fetch.
//
// Algorithm:
//
//	Peek:    if !peeked { cache = base.Next(); peeked = true }; return cache.
//	Next:    if peeked { peeked = false; return cache }; return base.Next().
//	HasNext: return peeked || base.HasNext().
//
// Time:  Peek O(1), Next O(1), HasNext O(1).
// Space: O(1) — one cached value plus a flag.
type LazyPeeking struct {
	base   *Iterator // underlying iterator
	cache  int       // stashed element (valid only when peeked == true)
	peeked bool      // true if a value has been peeked but not yet consumed
}

// NewLazyPeeking builds the lazy wrapper without pre-fetching.
func NewLazyPeeking(it *Iterator) *LazyPeeking {
	return &LazyPeeking{base: it}
}

// Peek fetches (once) and returns the upcoming element without consuming it.
func (p *LazyPeeking) Peek() int {
	if !p.peeked {
		p.cache = p.base.Next() // consume from base but stash it
		p.peeked = true
	}
	return p.cache
}

// Next returns the stashed peek if present, else advances the base directly.
func (p *LazyPeeking) Next() int {
	if p.peeked {
		p.peeked = false // consume the stash
		return p.cache
	}
	return p.base.Next()
}

// HasNext is true if a value is stashed or the base still has elements.
func (p *LazyPeeking) HasNext() bool {
	return p.peeked || p.base.HasNext()
}

// runExample drives the official operation sequence and returns the outputs.
//
// Ops:  PeekingIterator([1,2,3]); next()→1; peek()→2; next()→2; next()→3; hasNext()→false
func runExample(newIt func(*Iterator) peekingADT) []string {
	p := newIt(NewIterator([]int{1, 2, 3}))
	out := []string{}
	out = append(out, fmt.Sprintf("%d", p.Next()))    // 1
	out = append(out, fmt.Sprintf("%d", p.Peek()))    // 2
	out = append(out, fmt.Sprintf("%d", p.Next()))    // 2
	out = append(out, fmt.Sprintf("%d", p.Next()))    // 3
	out = append(out, fmt.Sprintf("%t", p.HasNext())) // false
	return out
}

func main() {
	fmt.Println("=== Approach 1: Cache One Element (Optimal) ===")
	fmt.Println(runExample(func(it *Iterator) peekingADT { return NewCachePeeking(it) })) // [1 2 2 3 false]

	fmt.Println("=== Approach 2: Lazy Peek ===")
	fmt.Println(runExample(func(it *Iterator) peekingADT { return NewLazyPeeking(it) })) // [1 2 2 3 false]
}
