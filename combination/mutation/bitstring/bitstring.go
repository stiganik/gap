/*
Package bitstring implements the bit string mutation technique for altering
genetic algorithm solutions.
*/
package bitstring

import (
	"math"
	"math/rand"
	"time"

	"github.com/stiganik/gap/combination"
	"github.com/stiganik/gap/solution"
)

var mask = []byte{128, 64, 32, 16, 8, 4, 2, 1}

func init() {
	combination.Register(combination.MUTATION_BIT_STRING, New)
}

type bitstring struct {
	elitism uint
	rnd     *rand.Rand
}

// New creates an instance of the bit string mutation technique.
func New(elitism uint) (combination.Combiner, error) {
	return &bitstring{
		elitism: elitism,
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func round(x float64) float64 {
	t := math.Trunc(x)
	if math.Abs(x-t) >= 0.5 {
		return t + math.Copysign(1, x)
	}
	return t
}

// Combine mutates one solution at a time by flipping a bit with probability
// 1/bitlen(solution)
//
// Solution A: 00000000
// OutA: 00000100
//
// The problem of choosing n independant values with a certain probability
// creates a binomial distibution. Fortunately the binomial distribution can be
// approximated using the normal distribution according to the formula N(np,
// np(1-p)). Although it is said that for this approximation to be more accurate
// a bit length of more than 20 is required, for our purposes it will do with
// any bit length.
func (b *bitstring) Combine(pool solution.Pool) error {
	specimens := pool.Specimens
	if len(specimens) == 0 {
		return nil
	}

	// XXX: Log debug feature.
	stdDeviance := math.Sqrt(1.0 - (1.0 / float64(pool.SpecimenBitSize)))

	elite := uint((float64(b.elitism) / float64(100)) * float64(len(specimens)))
	for i := elite; i < uint(len(specimens)); i++ {
		mutated := b.rnd.NormFloat64()*stdDeviance + 1.0
		var mutatedUint uint
		switch {
		case mutated < 0:
			mutatedUint = 0
		case mutated >= float64(pool.SpecimenBitSize):
			mutatedUint = pool.SpecimenBitSize - 1
		default:
			mutatedUint = uint(round(mutated))
		}

		for j := uint(0); j < mutatedUint; j++ {
			target := b.rnd.Intn(int(pool.SpecimenBitSize))
			specimens[i].Buf[target/8] ^= mask[target%8]
		}
	}

	return nil
}
