package err

type Code uint

const (
	CodeWorktreeNotFound Code = iota + 100
	CodeUnstagedChanges
)
