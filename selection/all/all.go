/*
Package all is a convenience package for importing all selection algorithms
implemented in this project.
*/
package all

import (
	// Blank importing all algorithms forces the algorithms to register
	// themselves at runtime saving the trouble of having to import all
	// algorithms one by one.
	_ "github.com/stiganik/gap/selection/scx"
)
