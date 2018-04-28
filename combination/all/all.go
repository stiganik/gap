/*
Package all is a convenience package for importing all combiner algorithms
implemented in this project.
*/
package all

import (
	// Blank importing all algorithms forces the algorithms to register
	// themselves at runtime saving the trouble of having to import all
	// algorithms one by one.
	_ "github.com/stiganik/gap/combination/crossover/singlepoint"
	_ "github.com/stiganik/gap/combination/crossover/twopoint"
	_ "github.com/stiganik/gap/combination/mutation/bitstring"
	_ "github.com/stiganik/gap/combination/mutation/flipbit"
)
