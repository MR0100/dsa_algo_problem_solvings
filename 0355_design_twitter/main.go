package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// Design a simplified Twitter:
//   - postTweet(userId, tweetId)          user posts a tweet.
//   - getNewsFeed(userId) -> []tweetId    the 10 most-recent tweet ids from the
//                                         user and everyone they follow, newest
//                                         first.
//   - follow(followerId, followeeId)      start following.
//   - unfollow(followerId, followeeId)    stop following.
//
// The only ordering signal is time: we tag each tweet with a global, monotonically
// increasing timestamp so we can compare recency across users.
//
// Two designs:
//   1. bruteForceTwitter — collect all relevant tweets, sort by time, take 10.
//   2. Twitter (optimal) — k-way merge the per-user tweet lists with a max-heap,
//      pulling only the 10 newest.

// tweet is one post with a global timestamp for recency ordering.
type tweet struct {
	id   int
	time int
}

// ── Approach 1: Brute Force (Collect All + Sort) ─────────────────────────────
//
// bruteForceTwitter keeps each user's tweets and follow-set, and builds a feed by
// gathering every tweet from the user + followees, sorting by timestamp, and
// slicing the newest 10.
//
// Intuition:
//
//	The feed is "the 10 newest tweets among a set of users". The simplest correct
//	thing: dump all those users' tweets into one list, sort by time descending,
//	take the first 10.
//
// Algorithm:
//
//	postTweet:  append {id, ++clock} to the user's tweet slice.
//	follow:     add followee to the follower's set.
//	unfollow:   remove followee from the set.
//	getNewsFeed: union {self} ∪ followees; concatenate their tweets; sort by time
//	             desc; return the first up-to-10 ids.
//
// Time:  getNewsFeed O(T log T) where T = total tweets among relevant users.
// Space: O(U + total tweets) for users, follows, and tweets.
type bruteForceTwitter struct {
	clock   int                  // global timestamp, increments per post
	tweets  map[int][]tweet      // userId -> their tweets in post order
	follows map[int]map[int]bool // userId -> set of followees
}

// newBruteForceTwitter builds an empty Twitter.
func newBruteForceTwitter() *bruteForceTwitter {
	return &bruteForceTwitter{
		tweets:  make(map[int][]tweet),
		follows: make(map[int]map[int]bool),
	}
}

// PostTweet records a tweet by userId with the next global timestamp.
func (t *bruteForceTwitter) PostTweet(userId, tweetId int) {
	t.clock++ // newer tweets get strictly larger timestamps
	t.tweets[userId] = append(t.tweets[userId], tweet{id: tweetId, time: t.clock})
}

// GetNewsFeed returns up to 10 most-recent tweet ids for userId's feed.
func (t *bruteForceTwitter) GetNewsFeed(userId int) []int {
	// Set of users whose tweets appear: the user themself + everyone they follow.
	sources := map[int]bool{userId: true}
	for f := range t.follows[userId] {
		sources[f] = true
	}

	var all []tweet
	for u := range sources {
		all = append(all, t.tweets[u]...) // gather every relevant tweet
	}
	// Sort by timestamp descending (newest first).
	sort.Slice(all, func(i, j int) bool { return all[i].time > all[j].time })

	res := []int{}
	for i := 0; i < len(all) && i < 10; i++ {
		res = append(res, all[i].id)
	}
	return res
}

// Follow makes followerId follow followeeId.
func (t *bruteForceTwitter) Follow(followerId, followeeId int) {
	if t.follows[followerId] == nil {
		t.follows[followerId] = make(map[int]bool)
	}
	t.follows[followerId][followeeId] = true
}

// Unfollow makes followerId stop following followeeId.
func (t *bruteForceTwitter) Unfollow(followerId, followeeId int) {
	if s := t.follows[followerId]; s != nil {
		delete(s, followeeId)
	}
}

// ── Approach 2: k-Way Merge with a Max-Heap (Optimal) ────────────────────────
//
// Twitter stores each user's tweets in post order (so the last element is their
// newest) and builds the feed by seeding a max-heap with each source's newest
// tweet, then popping 10 times, pushing the next-older tweet from whichever list
// the popped tweet came from.
//
// Intuition:
//
//	Each user's own tweet list is already sorted by time (append order). We only
//	need the 10 newest across k such sorted lists — a classic k-way merge. A
//	max-heap keyed by timestamp yields the global newest in O(log k) per pop, and
//	we stop after 10 pops instead of sorting everything.
//
// Algorithm:
//  1. sources = {self} ∪ followees.
//  2. Push each source's newest tweet (with a pointer to its position) into a
//     max-heap keyed by timestamp.
//  3. Repeat up to 10 times: pop the newest; record its id; if that list has an
//     older tweet, push it.
//
// Time:  getNewsFeed O(k + 10 log k), k = number of sources.
// Space: O(k) heap.
type Twitter struct {
	clock   int
	tweets  map[int][]tweet
	follows map[int]map[int]bool
}

// heapItem is a tweet plus enough info to fetch the next-older one from its list.
type heapItem struct {
	t     tweet // the tweet currently in the heap
	owner int   // userId whose list it came from
	idx   int   // its index in that user's tweet slice
}

// maxHeap orders heapItems by timestamp, newest at the top.
type maxHeap []heapItem

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i].t.time > h[j].t.time } // max by time
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(heapItem)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// Constructor builds an empty Twitter.
func Constructor() Twitter {
	return Twitter{
		tweets:  make(map[int][]tweet),
		follows: make(map[int]map[int]bool),
	}
}

// PostTweet records a tweet with the next global timestamp.
func (t *Twitter) PostTweet(userId, tweetId int) {
	t.clock++
	t.tweets[userId] = append(t.tweets[userId], tweet{id: tweetId, time: t.clock})
}

// GetNewsFeed returns up to 10 newest tweet ids via a k-way merge.
func (t *Twitter) GetNewsFeed(userId int) []int {
	sources := map[int]bool{userId: true}
	for f := range t.follows[userId] {
		sources[f] = true
	}

	h := &maxHeap{}
	heap.Init(h)
	// Seed the heap with each source's NEWEST tweet (last element of its slice).
	for u := range sources {
		lst := t.tweets[u]
		if len(lst) > 0 {
			last := len(lst) - 1
			heap.Push(h, heapItem{t: lst[last], owner: u, idx: last})
		}
	}

	res := []int{}
	for h.Len() > 0 && len(res) < 10 {
		top := heap.Pop(h).(heapItem) // globally newest remaining tweet
		res = append(res, top.t.id)
		if top.idx > 0 { // there is an older tweet in the same user's list
			prev := top.idx - 1
			heap.Push(h, heapItem{t: t.tweets[top.owner][prev], owner: top.owner, idx: prev})
		}
	}
	return res
}

// Follow makes followerId follow followeeId.
func (t *Twitter) Follow(followerId, followeeId int) {
	if t.follows[followerId] == nil {
		t.follows[followerId] = make(map[int]bool)
	}
	t.follows[followerId][followeeId] = true
}

// Unfollow makes followerId stop following followeeId.
func (t *Twitter) Unfollow(followerId, followeeId int) {
	if s := t.follows[followerId]; s != nil {
		delete(s, followeeId)
	}
}

func main() {
	// Official example:
	//   postTweet(1,5)
	//   getNewsFeed(1)      -> [5]
	//   follow(1,2)
	//   postTweet(2,6)
	//   getNewsFeed(1)      -> [6, 5]
	//   unfollow(1,2)
	//   getNewsFeed(1)      -> [5]

	fmt.Println("=== Approach 1: Brute Force (Collect + Sort) ===")
	b := newBruteForceTwitter()
	b.PostTweet(1, 5)
	fmt.Println(b.GetNewsFeed(1)) // expected [5]
	b.Follow(1, 2)
	b.PostTweet(2, 6)
	fmt.Println(b.GetNewsFeed(1)) // expected [6 5]
	b.Unfollow(1, 2)
	fmt.Println(b.GetNewsFeed(1)) // expected [5]

	fmt.Println("=== Approach 2: k-Way Merge Max-Heap (Optimal) ===")
	tw := Constructor()
	tw.PostTweet(1, 5)
	fmt.Println(tw.GetNewsFeed(1)) // expected [5]
	tw.Follow(1, 2)
	tw.PostTweet(2, 6)
	fmt.Println(tw.GetNewsFeed(1)) // expected [6 5]
	tw.Unfollow(1, 2)
	fmt.Println(tw.GetNewsFeed(1)) // expected [5]
}
