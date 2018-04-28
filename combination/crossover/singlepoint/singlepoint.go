/*
Package singlepoint implements the single point crossover technique for
combining genetic algorithm solutions.
*/
package singlepoint

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/stiganik/gap/combination"
	"github.com/stiganik/gap/solution"
)

func init() {
	combination.Register(combination.CROSSOVER_SINGLE_POINT, New)
}

type singlePoint struct {
	elitism uint
	rnd     *rand.Rand
}

// New creates an instance of the single point crossover technique.
func New(elitism uint) (combination.Combiner, error) {
	return &singlePoint{
		elitism: elitism,
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func byteValue(pivot int) (int, int) {
	// pad to next byte, then convert to bytes and subtract 1
	excessbits := pivot % 8
	pad := 8 - excessbits
	return ((pivot + pad) / 8) - 1, excessbits
}

// Since it depends on the hardware how bits are kept in memory by Go and the
// machine just slice the byte at the proper spot and swap with the other byte
// regardless whether it's little or big endian.
func spliceByte(pivot int, a, b byte) (ra byte, rb byte) {
	rightMask := byte((1 << uint(pivot)) - 1)
	leftMask := byte(^rightMask)

	// XXX: Debug log feature ?
	//
	//fmt.Println("Pivot:", pivot)
	//fmt.Printf("Right mask: %08b\n", rightMask)
	//fmt.Printf("Left mask: %08b\n", leftMask)

	ra = (a & rightMask) | (b & leftMask)
	rb = (b & rightMask) | (a & leftMask)

	// XXX: Debug log feature ?
	//
	//fmt.Printf("a byte: %08b\n", a)
	//fmt.Printf("b byte: %08b\n", b)
	//fmt.Printf("ra byte: %08b\n", ra)
	//fmt.Printf("rb byte: %08b\n", rb)

	return
}

// Combine combines two solution by selecting a random pivoting point (in bits)
// and splicing the solutions together into two output solutions.
//
// Solution A: ----------
// Solution B: //////////
//
// Pivot = 38 (bits)
//
// OutA: ----?/////
// OutB: ////?-----
func (s *singlePoint) Combine(pool solution.Pool) error {
	specimens := pool.Specimens
	if len(specimens)%2 != 0 {
		return fmt.Errorf("Poolsize must be divisible by 2")
	}

	if len(specimens) == 0 {
		return nil
	}

	elite := uint((float64(s.elitism) / float64(100)) * float64(len(specimens)))
	if elite%2 != 0 {
		elite++
	}

	tmp := make([]byte, pool.SpecimenByteSize)
	for i := elite; i < uint(len(specimens)); i += 2 {
		a := specimens[i].Buf
		b := specimens[i+1].Buf

		r := s.rnd.Intn(len(a) * 8)
		pBytes, pBits := byteValue(r)

		copy(tmp, a[(pBytes+1):len(a)])
		copy(a[(pBytes+1):len(a)], b[(pBytes+1):len(b)])
		copy(b[(pBytes+1):len(b)], tmp[0:len(a)-(pBytes+1)])

		a[pBytes], b[pBytes] = spliceByte(pBits, a[pBytes], b[pBytes])
	}

	return nil
}
