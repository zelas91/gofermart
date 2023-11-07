package types

import (
	"fmt"
	"time"
)

type ContextKey string

const (
	UserIDKey = ContextKey("userID")
	Logger    = ContextKey("logger")
	TxKey     = ContextKey("tx")
)

type RFC3339Time time.Time

func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	ts := time.Time(t).Format(time.RFC3339)
	stamp := fmt.Sprintf("\"%s\"", ts)
	return []byte(stamp), nil
}

const (
	PROCESSED = "PROCESSED"
)
