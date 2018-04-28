/*
Package gap implements a genetic algorithm framework for problem solving
using advanced crossover and mutation techniques.
*/
package gap

import (
	"fmt"
	"runtime"
	"time"

	"github.com/stiganik/gap/combination"
	"github.com/stiganik/gap/selector"
	"github.com/stiganik/gap/solution"

	// Statically import all selection and combination algorithms to make
	// them register themselves at runtime.
	_ "github.com/stiganik/gap/combination/all"
	_ "github.com/stiganik/gap/selector/all"
)

// FitnessFn defines a fitness function which takes in a solution in the form of
// a byte array calculates the fitness of the the solution and expresses it as a
// unsigned integer value from less to more fit, zero being totally unsuitable.
type FitnessFn func(s []byte) uint

// Algorithm defines a problem and the genetic algorithm used to solve the
// problem.
type Algorithm struct {
	// The fitness function used to evaluate solutions.
	FFn FitnessFn

	// SolutionBitSize is the bit size of the solution slice.
	SolutionBitSize uint

	// SolutionPoolSize is the amount of solutions generated in each
	// iteration of the alogirthm.
	SolutionPoolSize uint

	// Elitism is the percentage of values that pass on to the next
	// generation without the selection and combination process. The value
	// will be clipped between 0 and 100.
	Elitism uint

	// SelectionAlgorithm is the algorithm used to select solutions from the
	// solution pool for crossover.
	SelectionAlgorithm selector.Algorithm

	// CombinationAlgorithms is a slice of algortihms used to combine the
	// selected solutions into the solution candidates. The combination
	// algorithms are applied sequentially.
	CombinationAlgorithms []combination.Algorithm

	// ThreadCount sets the amount of threads used by the genetic algorithm.
	// By default it is set to the number of physical cores the processor
	// has.
	// Note: Currently not in use.
	ThreadCount uint
}

func (a *Algorithm) check() error {
	switch {
	case a.FFn == nil:
		return fmt.Errorf("Fitness function missing")
	case a.SolutionBitSize == 0:
		return fmt.Errorf("Solution size 0")
	case a.SolutionPoolSize == 0:
		return fmt.Errorf("Solution pool size 0")
	case a.SelectionAlgorithm == "":
		return fmt.Errorf("No selection algorithm selected")
	case a.ThreadCount == 0:
		return fmt.Errorf("Thread count set to 0")
	default:
		if a.Elitism > 100 {
			a.Elitism = 100
		}
		return nil
	}
}

// Result contains the result of a genetic algorithm and also additional
// information about the running of the algorithm.
type Result struct {
	ElapsedTime time.Duration
	Generation  uint
	Solution    solution.Specimen
}

// New creates a new default genetic algorithm for solving the problem described
// by the fitness function fn. More customization can be achieved by editing the
// exported fields of the Algorithm object. sbl is the Solution Bit Length value
// which determines how many bits the solution must include.
func New(fn FitnessFn, sbl uint) *Algorithm {
	return &Algorithm{
		FFn:                fn,
		SolutionBitSize:    sbl,
		SolutionPoolSize:   1000,
		SelectionAlgorithm: selector.SCX,
		Elitism:            3,
		CombinationAlgorithms: []combination.Algorithm{
			combination.CROSSOVER_SINGLE_POINT,
		},
		ThreadCount: uint(runtime.NumCPU()),
	}
}

// Run runs the genetic algorithm and retrieves the correctest answer once the
// goal of the algorithm is reached.
func (a *Algorithm) Run(g Goal) (ret Result, err error) {
	if err = a.check(); err != nil {
		return
	}

	if err = g.init(); err != nil {
		return
	}
	defer g.finalize()

	poolA := solution.NewPool(a.SolutionPoolSize, a.SolutionBitSize)
	poolB := solution.NewPool(a.SolutionPoolSize, a.SolutionBitSize)

	if err = poolA.Seed(); err != nil {
		return
	}

	sel, err := selector.New(a.SelectionAlgorithm, a.Elitism)
	if err != nil {
		return
	}

	var combiners []combination.Combiner
	for _, comb := range a.CombinationAlgorithms {
		var c combination.Combiner
		if c, err = combination.New(comb, a.Elitism); err != nil {
			return
		}
		combiners = append(combiners, c)
	}

	curPool := &poolA
	otherPool := &poolB

	var best solution.Specimen
	best.Copy((*curPool).Specimens[0])

	generation := uint(0)
	start := time.Now()
	for {
		if g.checkGen(generation) || g.checkTime() {
			break
		}
		// XXX: Generate initial workers
		// XXX: Distribute work to workers
		for i := range curPool.Specimens {
			fitness := a.FFn((*curPool).Specimens[i].Buf)
			(*curPool).Specimens[i].Fitness = fitness
			if g.checkFitness(fitness) {
				ret.ElapsedTime = time.Since(start)
				ret.Generation = generation
				ret.Solution.Copy((*curPool).Specimens[i])
				return
			}
		}

		(*curPool).Specimens.SortDesc()
		best.Copy((*curPool).Specimens[0])

		if g.checkTime() {
			break
		}

		if err = sel.Select((*curPool), (*otherPool)); err != nil {
			return
		}

		if g.checkTime() {
			break
		}

		for _, combiner := range combiners {
			if err = combiner.Combine((*otherPool)); err != nil {
				return
			}
		}

		if g.checkTime() {
			break
		}

		// Switch buffers and repeat genetic algorithm
		if curPool == &poolA {
			curPool = &poolB
			otherPool = &poolA
		} else {
			curPool = &poolA
			otherPool = &poolB
		}
		generation++
	}

	ret.ElapsedTime = time.Since(start)
	ret.Generation = generation
	ret.Solution.Copy(best)

	return
}
