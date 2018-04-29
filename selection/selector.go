/*
Package selection is the interface package for all selection algorithms.
*/
package selection

import (
	"fmt"
	"sync"

	"github.com/stiganik/gap/solution"
)

// Algorithm defines a set of supported solution selection algorithms.
type Algorithm string

const (
	// Fitness proportionate selection algorithm. Selects solutions with
	// probability proportianate to their fitness compared to the total
	// fitness of the solution pool.
	// https://en.wikipedia.org/wiki/Fitness_proportionate_selection
	SCX Algorithm = "scx"
)

var syncMutex sync.RWMutex
var algorithms map[Algorithm]NewFunc

// NewFunc creates a new instance of the algorithm implementation this function
// belongs to. Elitism is the percetage of solutions that should be considered
// "elite" and selected implicitly.
type NewFunc func(elitism uint) (Selector, error)

// Register registers a new selection algorithm for use through the Selector
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
func New(alg Algorithm, elitism uint) (Selector, error) {
	syncMutex.RLock()
	defer syncMutex.RUnlock()

	newFn, ok := algorithms[alg]
	if !ok {
		return nil, fmt.Errorf("Algorithm not linked: %s", alg)
	}

	return newFn(elitism)
}

// Selector is the interface for all selection algorithms in this project.
// Selection algorithms should not be used directly, only through this
// interface.
type Selector interface {
	// Select selects solutions from poolA and deposits them in poolB based on
	// the selection algorithm chosen.
	//
	// The selection algorithm MAY change the existing values of both
	// solution pools. The selection algorithm MUST NOT change the length
	// and/or capacity of the solution pools or the solutions.
	Select(poolA, poolB solution.Pool) error
}
