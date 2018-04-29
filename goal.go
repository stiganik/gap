package gap

import (
	"context"
	"fmt"
	"time"
)

// GoalFlag is a uint value type, that is used to specify different genetic
// algorithm goals of cancellation.
type GoalFlag uint

const (
	// TIME is a goal type flag. If this flag is set a counter is started
	// using the TimeN field in the goal structure which requests a graceful
	// stop from the algorithm after the time elapses. NOTE: The graceful
	// shutdown may take a while after the timer elapses.
	TIME GoalFlag = 1 << iota

	// GENERATION is a goal type flag. It starts counting generations and
	// ends the algorithm after the amount of generations specified in the
	// goal structure field GenN has been completed.
	GENERATION

	// FITNESS is a goal type flag. It monitors the fitness of all
	// solutions and stops the algorithm after the fitness specified in the
	// goal structure field FitN has been achieved.
	FITNESS
)

// Goal is a structure that contains information about the goals of the
// algorithm being run. It can be customized to use one or more end conditions.
type Goal struct {
	// A value that can be constructed from combining one or more goal flags
	// via bitwise OR.
	Goals GoalFlag

	// If TIME is set - the duration the algorithm will run before
	// cancellation.
	TimeN time.Duration

	// If GENERATION is set - the amount of generations the algorithm will
	// run before cancellation.
	GenN uint

	// If FITNESS is set - the fitness after which the algorithm will be
	// cancelled.
	FitN uint

	ctx    context.Context
	cancel context.CancelFunc
	term   bool
}

func (g *Goal) init() error {
	g.ctx = nil
	g.cancel = nil
	g.term = false
	if g.Goals&(TIME|FITNESS|GENERATION) == 0 {
		return fmt.Errorf("no goal set for algorithm")
	}
	if g.Goals&TIME != 0 {
		g.ctx, g.cancel = context.WithTimeout(context.Background(), g.TimeN)
	}
	return nil
}

func (g *Goal) checkTime() bool {
	if g.Goals&TIME != 0 {
		select {
		case <-g.ctx.Done():
			g.term = true
		default:
		}
	}
	return g.term
}

func (g *Goal) checkGen(gen uint) bool {
	if g.Goals&GENERATION != 0 && gen >= g.GenN {
		g.term = true
	}
	return g.term
}

func (g *Goal) checkFitness(fitness uint) bool {
	if g.Goals&FITNESS != 0 && fitness >= g.FitN {
		g.term = true
	}
	return g.term
}

func (g *Goal) finalize() {
	if g.cancel != nil {
		g.cancel()
	}
}
