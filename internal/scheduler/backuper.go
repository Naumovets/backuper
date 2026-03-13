package scheduler

import (
	"context"
	"time"

	"github.com/Naumovets/backuper/internal/backup"
	"github.com/Naumovets/backuper/internal/logger"
	"go.uber.org/zap"
)

func (s *scheduler) Run(ctx context.Context) {
	log := logger.GetLogger(ctx)

	log.Info("scheduler run")

	s.mx.Lock()
	defer s.mx.Unlock()
	for _, backuper := range s.backupers {
		stopCh, ok := s.stopChs[backuper.GetName()]
		if !ok {
			stopCh = make(chan struct{})
			s.stopChs[backuper.GetName()] = stopCh
		}

		go start(ctx, backuper, stopCh)
	}
}

func (s *scheduler) Add(ctx context.Context, backuper backup.Backuper) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	log := logger.GetLogger(ctx).With(zap.String("name", backuper.GetName()))

	_, ok := s.stopChs[backuper.GetName()]
	if ok {
		log.Error("cannot add backuper", zap.Error(ErrCloseChanAlreadyExist))

		return ErrCloseChanAlreadyExist
	}

	stop := make(chan struct{})
	s.stopChs[backuper.GetName()] = stop

	s.backupers = append(s.backupers, backuper)

	go start(ctx, backuper, stop)

	log.Info("add backuper")

	return nil
}

func start(ctx context.Context, backuper backup.Backuper, stop chan struct{}) {
	log := logger.GetLogger(ctx).With(zap.String("name", backuper.GetName()))
	ctx = context.WithValue(ctx, logger.CtxKeyLogger, log)

	log.Info("start backuper")

	// Make initial backup immediately
	if err := backuper.MakeBackup(ctx); err != nil {
		log.Error("cannot make initial backup", zap.Error(err))
	}

	ticker := time.NewTicker(backuper.GetInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("backuper done, graceful shutdown")

			return
		case <-ticker.C:
			err := backuper.MakeBackup(ctx)
			if err != nil {
				log.Error("cannot make backup", zap.Error(err))
			}
		case <-stop:
			log.Info("backuper done, stop chan")

			return
		}
	}
}

func (s *scheduler) Stop(ctx context.Context, name string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	log := logger.GetLogger(ctx).With(zap.String("name", name))

	ch, ok := s.stopChs[name]
	if !ok {
		log.Error("cannot stop backuper", zap.Error(ErrCloseChanNotFound))

		return ErrCloseChanNotFound
	}

	close(ch)
	delete(s.stopChs, name)

	for i, b := range s.backupers {
		if b.GetName() == name {
			s.backupers = append(s.backupers[:i], s.backupers[i+1:]...)
			break
		}
	}

	log.Info("backuper stoped")

	return nil
}
