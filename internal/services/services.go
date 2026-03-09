package services

import (
	"context"

	"github.com/Naumovets/Backuper/internal/config"
	"github.com/Naumovets/Backuper/internal/domain"
)

type service struct {
	cfg *config.Config
}

type Servicer interface {
	GetPostgresBackupers(ctx context.Context) ([]domain.PostgresBackuper, error)
	GetLoggingConfig(ctx context.Context) (domain.Logging, error)
	GetStorageConfig(ctx context.Context) (domain.Storage, error)
}

func NewService(cfg *config.Config) *service {
	return &service{
		cfg: cfg,
	}
}

func (s *service) GetPostgresBackupers(ctx context.Context) ([]domain.PostgresBackuper, error) {
	backupers := make([]domain.PostgresBackuper, 0)
	for name, backuper := range s.cfg.Backuper.Postgres {
		backupers = append(backupers, domain.PostgresBackuper{
			Name:            name,
			Threads:         backuper.Threads,
			Host:            backuper.Host,
			Port:            backuper.Port,
			DBName:          backuper.DBName,
			User:            backuper.User,
			Password:        backuper.Password,
			SSLKey:          backuper.SSLKey,
			SSLCert:         backuper.SSLCert,
			SSLRootCert:     backuper.SSLRootCert,
			SSLMode:         backuper.SSLMode,
			ApplicationName: backuper.ApplicationName,
			Compress:        backuper.Compress,
			CompressLevel:   backuper.CompressLevel,
			Interval:        backuper.Interval,
			MaxCount:        backuper.MaxCount,
			TimeFormat:      backuper.TimeFormat,
			PrefixFilename:  backuper.PrefixFilename,
			TableSchema:     backuper.TableSchema,
			TablePrefix:     backuper.TablePrefix,
			TableSuffix:     backuper.TableSuffix,
		})
	}

	return backupers, nil
}

func (s *service) GetLoggingConfig(ctx context.Context) (domain.Logging, error) {
	var rotationSettings *domain.RotationSettings
	if s.cfg.Backuper.Logging.RotationSettings != nil {
		rotationSettings = &domain.RotationSettings{
			MaxSize:  s.cfg.Backuper.Logging.RotationSettings.MaxSize,
			MaxCount: s.cfg.Backuper.Logging.RotationSettings.MaxCount,
			MaxAge:   s.cfg.Backuper.Logging.RotationSettings.MaxAge,
			Compress: s.cfg.Backuper.Logging.RotationSettings.Compress,
		}
	}

	return domain.Logging{
		Level:            s.cfg.Backuper.Logging.Level,
		Encoding:         s.cfg.Backuper.Logging.Encoding,
		OutputPaths:      s.cfg.Backuper.Logging.OutputPaths,
		ErrorOutputPaths: s.cfg.Backuper.Logging.ErrorOutputPaths,
		RotationSettings: rotationSettings,
	}, nil
}

func (s *service) GetStorageConfig(ctx context.Context) (domain.Storage, error) {
	return domain.Storage{Path: s.cfg.Backuper.StorageConfig.Path}, nil
}
