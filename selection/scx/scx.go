/*
Package scx is the fitness proportionate selection algorithm implementation.
*/
package scx

import (
	"math/rand"
	"time"

	"github.com/stiganik/gap/selection"
	"github.com/stiganik/gap/solution"
)

func init() {
	selection.Register(selection.SCX, New)
}

type scx struct {
	elitism uint
	rnd     *rand.Rand
}

// New creates an instance of the fitness proportionate selection algorithm.
func New(elitism uint) (selection.Selector, error) {
	scx := &scx{
		elitism: elitism,
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return scx, nil
}

// Select selects a solution from poolA with probability
// P(solution.fitness / poolA.totalFitness) and deposits the solution in poolB.
// This process is repeated until poolB is full.
func (s *scx) Select(poolA, poolB solution.Pool) error {
	specimensA := poolA.Specimens
	specimensB := poolB.Specimens

	var totalFitness uint
	for _, sol := range specimensA {
		totalFitness += sol.Fitness
	}

	elite := uint((float64(s.elitism) / float64(100)) * float64(len(specimensA)))
	for i := range specimensB {
		if uint(i) < elite {
			specimensB[i].Fitness = specimensA[i].Fitness
			copy(specimensB[i].Buf, specimensA[i].Buf)
			continue
		}

		r := uint(s.rnd.Float32() * float32(totalFitness)) // Generate value between 0 and totalFitness
		var sum uint
		for _, el := range specimensA {
			sum += el.Fitness
			if r <= sum {
				specimensB[i].Fitness = el.Fitness
				copy(specimensB[i].Buf, el.Buf)
				break
			}
		}
	}

	return nil
}
