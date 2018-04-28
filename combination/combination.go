/*
Package combination is the interface package for all crossover and mutation algorithms.
*/
package combination

import (
	"fmt"
	"sync"

	"github.com/stiganik/gap/solution"
)

// Algorithm defines a set of supported solution combination and mutation
// algorithms.
type Algorithm string

const (
	// The single point crossover algorithm creates two output solutions by
	// combining two input solutions via one pivot point.
	// https://en.wikipedia.org/wiki/Crossover_(genetic_algorithm)#Single-point
	CROSSOVER_SINGLE_POINT Algorithm = "crossover_single_point"

	// The two point crossover algorithm creates two output solutions by
	// combining two input solutions via two pivot points.
	// https://en.wikipedia.org/wiki/Crossover_(genetic_algorithm)#Two-point
	CROSSOVER_TWO_POINT Algorithm = "crossover_two_point"

	// The bit string mutation algorithm mutates every bit of a solution
	// with probability 1/bitlen(solution). This gives on average 1 mutation
	// per solution.
	// https://en.wikipedia.org/wiki/Mutation_(genetic_algorithm)
	MUTATION_BIT_STRING Algorithm = "mutation_bit_string"

	// The flip bit string mutation algorithm flips all the bits in the
	// bitstring without looking into it further.
	// https://en.wikipedia.org/wiki/Mutation_(genetic_algorithm)
	MUTATION_FLIP_BIT Algorithm = "mutation_flip_bit"
)

var syncMutex sync.RWMutex
var algorithms map[Algorithm]NewFunc

// NewFunc creates a new instance of the algorithm implementation this function
// belongs to. Elitism is the percetage of solutions that should be considered
// "elite" and left unaltered.
type NewFunc func(elitism uint) (Combiner, error)

// Register registers a new combination algorithm for use through the Combiner
// interface.
func Register(alg Algorithm, new NewFunc) {
	syncMutex.Lock()
	defer syncMutex.Unlock()

	if algorithms == nil {
		algorithms = make(map[Algorithm]NewFunc)
	}
	algorithms[alg] = new
}

// New creates a new instance of the selection algorithm defined by alg.
func New(alg Algorithm, elitism uint) (Combiner, error) {
	syncMutex.RLock()
	defer syncMutex.RUnlock()

	newFn, ok := algorithms[alg]
	if !ok {
		return nil, fmt.Errorf("Algorithm not linked: %s", alg)
	}

	return newFn(elitism)
}

// Combiner is the interface for all combination algorithms in this project.
// Combination algorithms should not be used directly, only through this
// interface.
type Combiner interface {
	// Combine combines or mutates solutions in the same pool to create new
	// solutions.
	//
	// The combination algorithm MAY change the existing byte values in the
	// solution pool. The combination algorithm MUST NOT change the length
	// and/or capacity of the solution pool or the solutions.
	Combine(pool solution.Pool) error
}
