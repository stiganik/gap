package main

import (
	"fmt"
	"os"

	"github.com/stiganik/gap/combination"
	"github.com/stiganik/gap/solution"

	_ "github.com/stiganik/gap/combination/all"
)

const (
	poolSize     = 2
	solutionSize = 80
)

var algos = []combination.Algorithm{
	combination.CROSSOVER_SINGLE_POINT,
	combination.CROSSOVER_TWO_POINT,
	combination.MUTATION_BIT_STRING,
	combination.MUTATION_FLIP_BIT,
}

func reset(orig, new solution.Pool) {
	for i := range orig.Specimens {
		copy(new.Specimens[i].Buf, orig.Specimens[i].Buf)
	}
}

func main() {
	pool := solution.NewPool(poolSize, solutionSize)
	poolResult := solution.NewPool(poolSize, solutionSize)
	pool.Seed()

	fmt.Println("Before:")
	fmt.Println("A:", pool.Specimens[0].Buf)
	fmt.Println("B:", pool.Specimens[1].Buf)
	fmt.Println()

	for _, algo := range algos {
		reset(pool, poolResult)

		comb, err := combination.New(algo, 0)
		if err != nil {
			fmt.Println("Failed to create combiner:", err)
			os.Exit(1)
		}

		if err = comb.Combine(poolResult); err != nil {
			fmt.Println("Failed combine values:", err)
			os.Exit(1)
		}

		fmt.Println("Algorithm:", string(algo))
		fmt.Println("A:", poolResult.Specimens[0].Buf)
		fmt.Println("B:", poolResult.Specimens[1].Buf)
		fmt.Println()
	}
}
