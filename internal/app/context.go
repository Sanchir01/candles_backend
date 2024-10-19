package app

type ContextKey int
type AccessContextKey string
type EnumContextStatus string

const (
	DataLoadersContextKey ContextKey       = iota
	AccessTokenContextKey AccessContextKey = "user"
)

const (
	EnumContextProcessingStatus EnumContextStatus = "in processing"
)

const (
	RefreshToken = "refreshToken"
	AccessToken  = "accessToken"
)
