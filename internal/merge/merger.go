// Package merge provides functionality to merge multiple sorted log streams
// into a single chronologically ordered output.
package merge

import (
	"container/heap"
	"time"
)

// Entry represents a single log line with its associated timestamp and source.
type Entry struct {
	Line      string
	Timestamp time.Time
	Source    string
}

// slot pairs an Entry with the index of the channel it came from.
type slot struct {
	Entry
	src int
}

// slotHeap implements heap.Interface ordered by timestamp (min-heap).
type slotHeap []slot

func (h slotHeap) Len() int            { return len(h) }
func (h slotHeap) Less(i, j int) bool  { return h[i].Timestamp.Before(h[j].Timestamp) }
func (h slotHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *slotHeap) Push(x interface{}) { *h = append(*h, x.(slot)) }
func (h *slotHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Merger merges multiple channels of log entries into one ordered stream.
type Merger struct {
	inputs []<-chan Entry
}

// New creates a Merger that will merge the provided input channels.
func New(inputs []<-chan Entry) *Merger {
	return &Merger{inputs: inputs}
}

// Merge reads from all input channels and emits entries in timestamp order.
// The returned channel is closed when all inputs are exhausted.
func (m *Merger) Merge() <-chan Entry {
	out := make(chan Entry, 64)
	go func() {
		defer close(out)

		h := make(slotHeap, 0, len(m.inputs))
		heap.Init(&h)

		// Seed the heap with the first entry from each input.
		for i, ch := range m.inputs {
			if e, ok := <-ch; ok {
				heap.Push(&h, slot{e, i})
			}
		}

		for h.Len() > 0 {
			min := heap.Pop(&h).(slot)
			out <- min.Entry
			if next, ok := <-m.inputs[min.src]; ok {
				heap.Push(&h, slot{next, min.src})
			}
		}
	}()
	return out
}
