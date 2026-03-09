package scheduler

import (
	"context"
	"sync"

	"github.com/Naumovets/Backuper/internal/backup"
	"github.com/Naumovets/Backuper/internal/backup/postgres"
	"github.com/Naumovets/Backuper/internal/config"
)

type scheduler struct {
	backupers []backup.Backuper
	stopChs   map[string]chan struct{}
	mx        sync.Mutex
}

type Scheduler interface {
	Run(ctx context.Context)
	Stop(ctx context.Context, name string) error
	Add(ctx context.Context, backuper backup.Backuper) error
}

func NewScheduler(conf *config.Config) *scheduler {
	backupers := make([]backup.Backuper, 0)
	stopChs := make(map[string]chan struct{})

	for name, pgConf := range conf.Backuper.Postgres {
		pgbackuper := postgres.NewBackuper(pgConf, conf.Backuper.StorageConfig, name)
		backupers = append(backupers, pgbackuper)
		stopChs[pgbackuper.GetName()] = make(chan struct{})
	}

	return &scheduler{
		backupers: backupers,
		mx:        sync.Mutex{},
		stopChs:   stopChs,
	}
}
