package app

type ContextKey int
type AccessContextKey string

const (
	DataLoadersContextKey ContextKey       = iota
	AccessTokenContextKey AccessContextKey = "user"
)
