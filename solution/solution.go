/*
Package solution wraps all complicated solution pool operations into one
package.
*/
package solution

import (
	"math/rand"
	"sort"
	"time"
)

// Specimen is a single solution for the genetic algorithm.
type Specimen struct {
	Fitness uint
	Buf     []byte
}

// Copy copies the values from the argument Specimen to the current specimen
func (s *Specimen) Copy(s2 Specimen) {
	s.Fitness = s2.Fitness
	if s.Buf == nil {
		s.Buf = make([]byte, len(s2.Buf))
	}
	copy(s.Buf, s2.Buf)
}

// Specimens is an array of specimen. It is a convenience wrapper for a Specimen
// slice with some additional functions.
type Specimens []Specimen

// Pool is a wrapper for the Specimens array with some additional functionality
// and problem information.
type Pool struct {
	SpecimenBitSize  uint
	SpecimenByteSize uint
	Specimens        Specimens
}

// Len is the number of elements in the collection.
func (s Specimens) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s Specimens) Less(i, j int) bool {
	if s[i].Fitness < s[j].Fitness {
		return true
	}
	return false
}

// Swap swaps the elements with indexes i and j.
func (s Specimens) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

// SortDesc sorts the elements in the collection in a descending order as
// defined by Less(j, i int).
func (s Specimens) SortDesc() {
	sort.Sort(sort.Reverse(s))
}

// SortAsc sorts the elements in the collection in an ascending order as defined
// by Less(i, j int).
func (s Specimens) SortAsc() {
	sort.Sort(s)
}

func byteSize(bitSize uint) uint {
	return uint((bitSize + 7) / 8)
}

// New creates a new solution pool object that contains poolSize solutions with
// length at least bitSize and returns a pointer to it.
func NewPool(poolSize, bitSize uint) Pool {
	p := Pool{
		SpecimenBitSize:  bitSize,
		SpecimenByteSize: byteSize(bitSize),
		Specimens:        make([]Specimen, poolSize, poolSize),
	}

	for i := range p.Specimens {
		p.Specimens[i].Buf = make([]byte, p.SpecimenByteSize, p.SpecimenByteSize)
	}
	return p
}

// Seed seeds the pool with random values.
func (p Pool) Seed() error {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	for i := range p.Specimens {
		if _, err := r.Read(p.Specimens[i].Buf); err != nil {
			return err
		}
	}

	return nil
}
