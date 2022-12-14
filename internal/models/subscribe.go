package models

import (
	"time"

	"github.com/robfig/cron/v3"
)

// SubInfo Структура с информацией о приобретенной пользователем подпиской
type SubInfo struct {
	Username	string
	Exp			time.Duration
}

// CurrentSub Структура с информацией о текущей подписке пользователя
type CurrentSub struct {
	Exp		time.Duration
	RemId	[]cron.EntryID
}
