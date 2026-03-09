package backup

import (
	"context"
	"time"
)

type Backuper interface {
	MakeBackup(ctx context.Context) error
	GetInterval() time.Duration
	GetName() string
}
