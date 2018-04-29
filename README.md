Genetic Algorithm Project. (GAP v0.1)

This project aims to implement a Go based genetic algorithm package. The interface is designed to be simple to use for both beginners (default values) and advanced users (customizable).

The project is being developed using the git branching model described in this post:

http://nvie.com/posts/a-successful-git-branching-model/

This means that the "master" branch will contain stable released versions of the project, but may not always contain all the latest commits from "develop".

## Quick start

A somewhat trivial example will generate a 1 byte solution where having a higher byte value is favored by the fitness function. It will stop after 1 minute or when a fitness of 255 is reached. Usually this will not get past generation 0 since a terminating solution (255) is generated with high probability while seeding the solution pool with random values.

```go
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/stiganik/gap"
)

func Fitness(s []byte) uint {
    return uint(s[0])
}

func main() {
    alg := gap.Algorithm {
        FFn: Fitness,
        SolutionBitSize: uint(8),
    }

    goal := gap.Goal {
        Goals: gap.TIME | gap.FITNESS,
        TimeN: 1 * time.Minute,
        FitN:  255,
    }

    res, err := alg.Run(goal)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to run algorithm: %s", err.Error())
        os.Exit(1)
    }

    fmt.Println("Algorithm completed")
    fmt.Println("Elapsed time: ", res.ElapsedTime)
    fmt.Println("Generations: ", res.Generation)
    fmt.Println("Fitness: ", res.Solution.Fitness)
}
```

## Algorithms

Currently implemented algorithms.

### Selection

- SCX - Fitness proportionate selection algorithm. Selects solutions with
	probability proportianate to their fitness compared to the total
	fitness of the solution pool.	https://en.wikipedia.org/wiki/Fitness_proportionate_selection

### Combination

- Single Point Crossover - The single point crossover algorithm creates two
  output solutions by combining two input solutions via one pivot point.
	https://en.wikipedia.org/wiki/Crossover_(genetic_algorithm)#Single-point
- Two Point Crossover - The two point crossover algorithm creates two output
  solutions by combining two input solutions via two pivot points.
  https://en.wikipedia.org/wiki/Crossover_(genetic_algorithm)#Two-point

- Bit String Mutation - The bit string mutation algorithm mutates every bit
  of a solution with probability P = 1/solution length. This gives on average
  1 mutation per solution.
  https://en.wikipedia.org/wiki/Mutation_(genetic_algorithm)
- Flip bit mutation - The flip bit string mutation algorithm flips all the
  bits in the bitstring without looking into it further.
  https://en.wikipedia.org/wiki/Mutation_(genetic_algorithm)

## Constructor

While the `Algorithm` structure is usually initialized manually, the package does provide a constructor in case you need a function to create a new algorithm for whatever reason. The function has the signature:

```go
func New(fn FitnessFn, sbl uint) *Algorithm
```

where `sbl` is the solution bit size.

## Customization

The algorithm does fine with its default settings up to a point, but at some point a little customization is required.

### Goals

Goals can be customized to contain one or more terminating conditions. The possible conditions are:

- TIME
- GENERATION
- FITNESS

They can be combined in any sequence and amount using the bitwise OR operator '|':

```go
var goal gap.Goal
goal.Goals = gap.TIME | gap.GENERATION
```

Depending on which terminating conditions you have chosen the corresponding fields in the Goal structure need to be evaluated:

- TimeN
- GenN
- FitN

An example utilizing all three terminating conditions would look like this:

```go
goal := gap.Goal {
    Goals: gap.TIME | gap.GENERATION | gap.FITNESS,
    TimeN: 30 * time.Second, // Terminate after 30 seconds
    GenN:  20,               // Terminate after 20 generations
    FitN:  300,              // Terminate after fitness reaches 300
}

```

All of the conditions are valid and checked at once and the first one to become true will stop the algorithm.

### Algorithm

The algorithm itself has 6 customizable features:

- Fitness function
- Solution bit size
- Solution pool size
- Elitism
- Selection algorithm
- Combination algorithms

#### Fitness function

The fitness function is a function pointer defined thusly:

```go
type FitnessFn func(s []byte) uint
```

Every time a solution needs to be evaluated this function is called and the solution byte array is passed into the function. The fitness function can be changed by accessing the structure field "FFn":

```go
gap.Algorithm{
    FFn: func(s []byte) uint { return 0 },
}
```

#### Solution bit size

The solution bit size determines how many bits the solution must have. Since bits come in bunches of 8 (a.k.a bytes) then only the guarantee is made that the solution will contain at least the solution bit size amount of bits.

```go
gap.Algorithm{
    SolutionBitSize: 34, // 4 bytes and 2 bits will be represented as 5 bytes
}
```

#### Solution pool size

The solution poolsize determines how many solutions are generated into the gene pool. By default this value is set to `1000`.

```go
gap.Algorithm{
    SolutionPoolSize: 10000, // Use 10 000 Solution candidates
}
```

#### Elitism

Elitism is a pointer to a value between 0 and 100 that determines which percentage of the best solutions of each generation pass on to the next generation unaltered. By default this value is set to `3`.

```go
elite := uint(10)
gap.Algorithm{
    Elitism: &elite, // 10% of the elite pass on without selection and altering
}
```

#### Selection algorithm

The selection algorithm determines which algorithm is used to choose solutions from the possible solutions into the next generation. By default this value is set to `selection.SCX`

```go
gap.Algorithm{
    SelectionAlgorithm: selection.SCX, // Use fitness proportionate selection algorithm
}
```

#### Combination algorithms

The combination algorithms determine which algorithms are chosen to mutate and/or combine the selected solutions to form the next generation. This value can contain multiple algorithms which are applied sequentially one after the other. The end result is the next generation of solutions. By default this value is set to `[]combination.Algorithm{combination.CROSSOVER_SINGLE_POINT}`

```go
gap.Algorithm{
    CombinationAlgorithms: []combination.Algorithm{
        combination.CROSSOVER_SINGLE_POINT,
    },
}
```

#### Example

An example of customizing a genetic algorithm:

```go
package main

import (
    "github.com/stiganik/gap"
    "github.com/stiganik/gap/combination"
)

func Fitness(s []byte) uint {
    return uint(s[0])
}

func main() {
    elite := uint(10)
    alg := gap.Algorithm{
        FFn:              Fitness,
        SolutionBitSize:  uint(8),
        SolutionPoolSize: 10000,
        Elitism:          &elite,
        CombinationAlgorithms: []combination.Algorithm{
            combination.CROSSOVER_SINGLE_POINT,
            combination.MUTATION_BIT_STRING,
        },
    }
}
```
