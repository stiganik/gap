/*
Package twopoint implements the two point point crossover technique for
combining genetic algorithm solutions.
*/
package twopoint

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/stiganik/gap/combination"
	"github.com/stiganik/gap/solution"
)

func init() {
	combination.Register(combination.CROSSOVER_TWO_POINT, New)
}

type twopoint struct {
	elitism uint
	rnd     *rand.Rand
}

// New creates an instance of the two point crossover technique.
func New(elitism uint) (combination.Combiner, error) {
	return &twopoint{
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

// Combine combines two solution by selecting two random pivoting points (in
// bits) and splicing the solutions together into two output solutions.
//
// Solution A: --------------------
// Solution B: ////////////////////
//
// Pivot 1 = 38 (bits)
// Pivot 2 = 82 (bits)
//
// OutA: ----?/////?---------
// OutB: ////?-----?/////////
func (t *twopoint) Combine(pool solution.Pool) error {
	specimens := pool.Specimens
	if len(specimens)%2 != 0 {
		return fmt.Errorf("Poolsize must be divisible by 2")
	}

	if len(specimens) == 0 {
		return nil
	}

	elite := uint((float64(t.elitism) / float64(100)) * float64(len(specimens)))
	if elite%2 != 0 {
		elite++
	}

	tmp := make([]byte, len(specimens[0].Buf))
	for i := elite; i < uint(len(specimens)); i += 2 {
		a := specimens[i].Buf
		b := specimens[i+1].Buf

		r1 := t.rnd.Intn(len(a) * 8)
		r2 := t.rnd.Intn(len(b) * 8)
		if r1 > r2 {
			tmp := r2
			r2 = r1
			r1 = tmp
		}

		p1Bytes, p1Bits := byteValue(r1)
		p2Bytes, p2Bits := byteValue(r2)

		// only copy bytes in between, if there are any full bytes.
		if p2Bytes-p1Bytes > 1 {
			copy(tmp, a[(p1Bytes+1):p2Bytes])
			copy(a[(p1Bytes+1):p2Bytes], b[(p1Bytes+1):p2Bytes])
			copy(b[(p1Bytes+1):p2Bytes], tmp[0:p2Bytes-(p1Bytes+1)])
		}

		a[p1Bytes], b[p1Bytes] = spliceByte(p1Bits, a[p1Bytes], b[p1Bytes])
		a[p2Bytes], b[p2Bytes] = spliceByte(p2Bits, a[p2Bytes], b[p2Bytes])
	}

	return nil
}
