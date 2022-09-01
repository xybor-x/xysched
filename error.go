package xysched

import (
	"github.com/xybor-x/xyerror"
)

var egen = xyerror.Register("xysched", 300000)

// Errors of package xysched.
var (
	CallError = egen.NewClass("CallError")
)
