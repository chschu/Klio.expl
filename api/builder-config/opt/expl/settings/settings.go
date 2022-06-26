package settings

import "time"

const (
	MaxUTF16LengthForKey    = 50
	MaxUTF16LengthForValue  = 500
	EntryToStringTimeZone   = "Europe/Berlin"
	EntryToStringTimeFormat = "02.01.2006 15:04"
	MaxExplCount            = 20
	ExplTokenValidity       = time.Minute * 5
	MaxFindCount            = 20
	FindTokenValidity       = time.Minute * 5
	TopExplCount            = 100
	HandlerTimeout          = time.Second * 2
)
