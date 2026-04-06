package constants

import "atomicptr.dev/deeperr"

const (
	ErrWorktreeNotFound deeperr.Code = iota + 100
	ErrUnstagedChanges
)
