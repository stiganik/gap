/*
Package flipbit implements the flip bit mutation technique for altering
genetic algorithm solutions.
*/
package flipbit

import (
	"github.com/stiganik/gap/combination"
	"github.com/stiganik/gap/solution"
)

func init() {
	combination.Register(combination.MUTATION_FLIP_BIT, New)
}

type flipbit struct {
	elitism uint
}

// New creates an instance of the flip bit mutation technique.
func New(elitism uint) (combination.Combiner, error) {
	return &flipbit{
		elitism: elitism,
	}, nil
}

// Combine mutates one solution at a time by flipping all bits to their opposite
//
// Solution A: 00000000
// OutA: 11111111
func (f *flipbit) Combine(pool solution.Pool) error {
	specimens := pool.Specimens
	if len(specimens) == 0 {
		return nil
	}

	elite := uint((float64(f.elitism) / float64(100)) * float64(len(specimens)))
	for i := elite; i < uint(len(specimens)); i++ {
		for j := range specimens[i].Buf {
			specimens[i].Buf[j] = ^specimens[i].Buf[j]
		}
	}

	return nil
}
