package benchttp

import (
	"time"

	"github.com/benchttp/worker/stats"
)

// StatsDescriptor contains a computed stats group description information
type StatsDescriptor struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userID"`
	Tag        string    `json:"tag"`
	FinishedAt time.Time `json:"finishedAt"`
}

// Stats contains StatsDescriptor, Codestats and stats.Stats of a given computed stats group
type Stats struct {
	Descriptor StatsDescriptor          `json:"descriptor"`
	Code       stats.StatusDistribution `json:"code"`
	Time       stats.Common             `json:"time"`
}

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	IdleConn int
	MaxConn  int
}
