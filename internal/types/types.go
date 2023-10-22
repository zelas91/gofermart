package types

import (
	"fmt"
	"time"
)

type ContextKey string

const (
	UserIDKey = ContextKey("userID")
	Logger    = ContextKey("logger")
)

type RFC3339Time time.Time

func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	ts := time.Time(t).Format(time.RFC3339)
	stamp := fmt.Sprintf("\"%s\"", ts)
	return []byte(stamp), nil
}

type OrderStatus string

const (
	NEW        = OrderStatus("NEW")
	PROCESSING = OrderStatus("PROCESSING")
	INVALID    = OrderStatus("INVALID")
	PROCESSED  = OrderStatus("PROCESSED")
)
