package bpe

import (
	"container/heap"
	"errors"
	"math/rand"
)

const DefaultCacheCapacity int = 10000

type Merge struct {
	Pos   int
	Rank  int
	NewId int
}

type mergeHeap []Merge

func (h mergeHeap) Len() int { return len(h) }

func (h mergeHeap) Less(i, j int) bool {
	if h[i].Rank != h[j].Rank {
		return h[i].Rank < h[j].Rank
	}
	return h[i].Pos < h[j].Pos
}

func (h mergeHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *mergeHeap) Push(x interface{}) {
	*h = append(*h, x.(Merge))
}

func (h *mergeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

type Symbol struct {
	C    int
	Prev int
	Next int
	Len  int
}

// Some slice methods to manipulate slice struct Symbol
type Symbols []Symbol

// Insert inserts a symbol to the slice at `i` index point
func (ss *Symbols) Insert(s Symbol, i int) error {
	var err error
	if i < 0 || i > len(*ss) {
		err = errors.New("`i` index is out of bound.")
		return err
	}
	*ss = append((*ss)[:i], append([]Symbol{s}, (*ss)[i:]...)...)
	return nil
}

// Remove removes a symbol from the slice at `i` index point
func (ss *Symbols) Remove(i int) error {
	var err error
	if i < 0 || i > len(*ss)-1 {
		err = errors.New("`i` index is out of bound.")
		return err
	}
	*ss = append((*ss)[:i], (*ss)[i+1:]...)
	return nil
}

func (s *Symbol) MergeWith(other *Symbol, newC int) {
	s.C = newC
	s.Len += other.Len
	s.Next = other.Next
}

type Word struct {
	Symbols Symbols
}

func NewWord() *Word {
	return &Word{
		// Symbols: Symbols{},
		Symbols: []Symbol{},
	}
}

func (w *Word) Add(c int, byteLen int) {

	var symbols []Symbol

	symLen := len(w.Symbols)

	if symLen == 0 {
		newSym := Symbol{
			C:    c,
			Prev: -1,
			Next: -1,
			Len:  byteLen,
		}
		symbols = append(symbols, newSym)
	} else {
		for i, s := range w.Symbols {
			// first
			if i == 0 {
				sym := &w.Symbols[i]
				sym.Next = 1
				sym.Prev = -1
				symbols = append(symbols, *sym)
			} else if i == symLen-1 { // last
				sym := &w.Symbols[i]
				sym.Next = symLen
				sym.Prev = symLen - 2
				symbols = append(symbols, *sym)
			} else {
				symbols = append(symbols, s)
			}
		}

		newSym := Symbol{
			C:    c,
			Prev: symLen - 1,
			Next: -1,
			Len:  byteLen,
		}
		symbols = append(symbols, newSym)
	}

	w.Symbols = symbols
}

type Pair struct {
	C1 int
	C2 int
}

// PairVal holds pair's rank and NewId
type PairVal struct {
	Rank  int
	NewId int
}

type WChange struct {
	C1     int
	C2     int
	Change int
}

// Merge finds any pairs of (c1, c2) and removes in place. It also maps changes depending
// on the position of the pair in word.
func (w *Word) Merge(c1, c2, replacement int) ([]WChange, error) {
	// fmt.Printf("before merge word symbols: %v\n", w.Symbols)
	// fmt.Printf("c1: %v - c2: %v- replacement: %v\n", c1, c2, replacement)
	var changes []WChange

	i := 0

	for {
		if i >= len(w.Symbols) {
			break
		}
		// found a pair
		if w.Symbols[i].C == c1 && (i+1) < len(w.Symbols) && w.Symbols[i+1].C == c2 {
			first := w.Symbols[i]
			second := w.Symbols[i+1]

			// If there's other characters before the pair
			if i > 0 {
				changes = append(changes, WChange{
					C1:     w.Symbols[i-1].C,
					C2:     first.C,
					Change: -1,
				})
				changes = append(changes, WChange{
					C1:     w.Symbols[i-1].C,
					C2:     replacement,
					Change: 1,
				})
			}

			// Remove in place
			newS := Symbol{
				C:    replacement,
				Prev: first.Prev,
				Next: second.Next,
				Len:  first.Len + second.Len,
			}

			// Insert replacement before first `char` of pair
			err := w.Symbols.Insert(newS, i)
			if err != nil {
				return nil, err
			}

			// Remove first `char` of pair
			err = w.Symbols.Remove(i + 1)
			if err != nil {
				return nil, err
			}
			// And then the second
			err = w.Symbols.Remove(i + 1)
			if err != nil {
				return nil, err
			}

			// If there are other `chars` after the pair
			if i > 0 && i < len(w.Symbols)-1 {
				// fmt.Println("Yes, there some char after the pair")
				changes = append(changes, WChange{
					C1:     second.C,
					C2:     w.Symbols[i+1].C,
					Change: -1,
				})
				changes = append(changes, WChange{
					C1:     replacement,
					C2:     w.Symbols[i+1].C,
					Change: 1,
				})
			}
		}

		i++

	} // End of `for` loop

	// fmt.Printf("After merge word symbols: %v\n", w.Symbols)

	// fmt.Printf("Num of changes: %v\n", len(changes))
	// fmt.Printf("They are: %v\n", changes)

	return changes, nil
}

func (w *Word) MergeAll(merges map[Pair]PairVal, dropoutOpt ...float32) {
	var dropout float32 = 0.0
	if dropoutOpt != nil {
		dropout = dropoutOpt[0]
	}

	pq := &mergeHeap{}
	heap.Init(pq)

	for i := 0; i < len(w.Symbols)-1; i++ {
		pair := Pair{
			C1: w.Symbols[i].C,
			C2: w.Symbols[i+1].C,
		}
		if m, ok := merges[pair]; ok {
			heap.Push(pq, Merge{Pos: i, Rank: m.Rank, NewId: m.NewId})
		}
	}

	skip := make([]Merge, 0)
	r := rand.New(rand.NewSource(99))

	for pq.Len() > 0 {
		top := heap.Pop(pq).(Merge)

		if dropout > 0.0 && r.Float32() < dropout {
			skip = append(skip, top)
			continue
		}

		for _, s := range skip {
			heap.Push(pq, s)
		}
		skip = skip[:0]

		if top.Pos < 0 || top.Pos >= len(w.Symbols) {
			continue
		}
		current := &w.Symbols[top.Pos]
		if current.Len == 0 || current.Next == -1 {
			continue
		}

		nextPos := current.Next
		if nextPos < 0 || nextPos >= len(w.Symbols) {
			continue
		}
		right := w.Symbols[nextPos]
		if right.Len == 0 {
			continue
		}

		targetPair := Pair{C1: current.C, C2: right.C}
		m, ok := merges[targetPair]
		if !ok || m.NewId != top.NewId {
			continue
		}

		current.MergeWith(&right, top.NewId)
		w.Symbols[nextPos].Len = 0

		if right.Next > -1 && right.Next < len(w.Symbols) {
			w.Symbols[right.Next].Prev = top.Pos
		}

		if current.Prev >= 0 {
			prev := current.Prev
			prevSymbol := w.Symbols[prev]
			if prevSymbol.Len != 0 {
				newPair := Pair{C1: prevSymbol.C, C2: current.C}
				if m, ok := merges[newPair]; ok {
					heap.Push(pq, Merge{Pos: prev, Rank: m.Rank, NewId: m.NewId})
				}
			}
		}

		next := current.Next
		if next >= 0 && next < len(w.Symbols) {
			nextSymbol := w.Symbols[next]
			if nextSymbol.Len != 0 {
				newPair := Pair{C1: current.C, C2: nextSymbol.C}
				if m, ok := merges[newPair]; ok {
					heap.Push(pq, Merge{Pos: top.Pos, Rank: m.Rank, NewId: m.NewId})
				}
			}
		}
	}

	w.removeSymbols()
}

// removeSymbols removes all symbols with lenth == 0
func (w *Word) removeSymbols() {
	var filtered []Symbol
	for _, s := range w.Symbols {
		if s.Len != 0 {
			filtered = append(filtered, s)
		}
	}
	w.Symbols = filtered
}

func (w *Word) GetChars() []int {
	var res []int
	for _, s := range w.Symbols {
		res = append(res, s.C)
	}
	return res
}

func (w *Word) GetOffsets() [][]int {
	var offsets [][]int

	var pos int = 0
	for _, s := range w.Symbols {
		end := pos + s.Len
		offsets = append(offsets, []int{pos, end})
		pos += s.Len
	}

	return offsets
}
