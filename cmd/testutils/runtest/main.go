package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stiganik/gap"
)

const solutionSize = 100

func Fitness(s []byte) uint {
	var fitness uint
	counter := 0
	for _, b := range s {
		for i := 0; i < 8; i++ {
			if b&1 == 0 && counter%2 == 0 {
				fitness += 1
			}
			if b&1 == 1 && counter%2 != 0 {
				fitness += 1
			}
			b = (b >> 1)
			counter++
			if counter == solutionSize {
				return fitness
			}
		}
	}
	return fitness
}

func main() {
	alg := gap.New(Fitness, solutionSize)
	alg.SolutionPoolSize = 10000

	res, err := alg.Run(gap.Goal{
		Goals: gap.TIME | gap.FITNESS,
		TimeN: 10 * time.Minute,
		FitN:  100,
	})

	if err != nil {
		fmt.Println("Error running algorithm:", err)
		os.Exit(1)
	}

	fmt.Println("Algorithm completed")
	fmt.Println("Elapsed time: ", res.ElapsedTime)
	fmt.Println("Generations: ", res.Generation)
	fmt.Println("Fitness: ", res.Solution.Fitness)
	var solution string
	counter := 0
	for _, b := range res.Solution.Buf {
		for i := 0; i < 8; i++ {
			if b&1 == 0 {
				solution += "0"
			}
			if b&1 == 1 {
				solution += "1"
			}
			b = (b >> 1)
			counter++
			if counter == solutionSize {
				fmt.Println("The solution is:", solution)
				os.Exit(0)
			}
		}
	}
}
