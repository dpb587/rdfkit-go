package rdfcanon

import "errors"

var ErrMaxIterationsReached = errors.New("maximum iterations reached")
var ErrMaxRecursionDepthReached = errors.New("maximum recursion depth reached")
