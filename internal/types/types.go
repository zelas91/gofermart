package types

type ContextKey string

const (
	UserIDKey = ContextKey("userID")
	Logger    = ContextKey("logger")
)
