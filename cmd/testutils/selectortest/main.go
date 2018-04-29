package main

import (
	"fmt"
	"os"

	"github.com/stiganik/gap/selection"
	"github.com/stiganik/gap/solution"

	_ "github.com/stiganik/gap/selection/all"
)

const (
	poolSize     = 20
	solutionSize = 10
	elitism      = 10
)

var tests = []selection.Algorithm{selection.SCX}

func main() {
	poolA := solution.NewPool(poolSize, solutionSize)
	poolB := solution.NewPool(poolSize, solutionSize)

	err := poolA.Seed()
	if err != nil {
		fmt.Println("Failed to seed pool:", err)
		os.Exit(1)
	}

	for i := range poolA.Specimens {
		poolA.Specimens[i].Fitness = uint(i)
	}

	poolA.Specimens.SortDesc()
	before := make([]uint, poolSize, poolSize)

	for i := range poolA.Specimens {
		before[i] = poolA.Specimens[i].Fitness
	}

	fmt.Println("PoolA:", before)
	for _, test := range tests {
		fmt.Println("\nAlgorithm:", string(test))

		sel, err := selection.New(test, elitism)
		if err != nil {
			fmt.Println("Failed to create select algorithm:", err)
			os.Exit(1)
		}
		if err = sel.Select(poolA, poolB); err != nil {
			fmt.Println("Failed to select:", string(test))
			os.Exit(1)
		}

		poolB.Specimens.SortDesc()

		after := make([]uint, poolSize, poolSize)
		for i := range poolB.Specimens {
			after[i] = poolB.Specimens[i].Fitness
		}

		fmt.Println("PoolB:", after)
	}
}
