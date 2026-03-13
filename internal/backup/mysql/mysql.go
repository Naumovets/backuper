package mysql

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Naumovets/backuper/internal/backup"
	"github.com/Naumovets/backuper/internal/config"
	"github.com/Naumovets/backuper/internal/logger"
	"github.com/Naumovets/backuper/pkg/compressor"
	"github.com/Naumovets/go-mysqldump"
	"go.uber.org/zap"
)

type mysqlBackup struct {
	mysqlConf   config.MysqlConfig
	storageConf config.StorageConfig
	dumper      *mysqldump.Data
	name        string
}

func NewBackuper(mysqlConf config.MysqlConfig, storageConfig config.StorageConfig, name string) *mysqlBackup {
	config := ToMysqlConfig(&mysqlConf)

	opts := mysqldump.TableOptions{}

	if mysqlConf.TablePrefix != nil {
		opts.TablePrefix = *mysqlConf.TablePrefix
	}

	if mysqlConf.TableSuffix != nil {
		opts.TableSuffix = *mysqlConf.TableSuffix
	}

	dumper, _ := mysqldump.Init(config, opts)

	return &mysqlBackup{
		name:        name,
		mysqlConf:   mysqlConf,
		storageConf: storageConfig,
		dumper:      dumper,
	}
}

func (m *mysqlBackup) MakeBackup(ctx context.Context) error {
	log := logger.GetLogger(ctx)

	log.Info("Start make backup mysql database")

	targetDir := filepath.Join(m.storageConf.Path, m.name)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	dumpFilename := filepath.Join(targetDir, m.getFilename())
	file, err := os.Create(dumpFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Dump database to file
	if err := m.dumper.MakeDump(file); err != nil {
		os.Remove(dumpFilename)

		return fmt.Errorf("error dumping: %w", err)
	}

	if m.mysqlConf.Compress {
		log.Info("Compressing backup")
		compressedFile, err := compressor.CompressFile(dumpFilename, m.mysqlConf.CompressLevel)
		if err != nil {
			log.Error("Failed to compress backup", zap.Error(err))
		} else {
			os.Remove(dumpFilename)
			dumpFilename = compressedFile
		}
	}

	log.Info("Database is dumped", zap.String("file", dumpFilename))

	// Run retention policy
	if err := m.runRetention(ctx, targetDir); err != nil {
		log.Error("Retention policy failed", zap.Error(err))
	}

	return nil
}

func (m *mysqlBackup) runRetention(ctx context.Context, targetDir string) error {
	log := logger.GetLogger(ctx)

	if m.mysqlConf.MaxCount <= 0 {
		return nil
	}

	files, err := os.ReadDir(targetDir)
	if err != nil {
		return err
	}

	var backupFiles []os.DirEntry
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), m.mysqlConf.PrefixFilename) {
			backupFiles = append(backupFiles, f)
		}
	}

	if len(backupFiles) <= m.mysqlConf.MaxCount {
		return nil
	}

	// Sort by info (modification time) - oldest first
	sort.Slice(backupFiles, func(i, j int) bool {
		infoI, _ := backupFiles[i].Info()
		infoJ, _ := backupFiles[j].Info()
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	toDelete := len(backupFiles) - m.mysqlConf.MaxCount
	for i := 0; i < toDelete; i++ {
		path := filepath.Join(targetDir, backupFiles[i].Name())
		if err := os.Remove(path); err != nil {
			log.Error("Failed to delete old backup", zap.String("path", path), zap.Error(err))
		} else {
			log.Info("Deleted old backup", zap.String("path", path))
		}
	}

	return nil
}

func (m *mysqlBackup) getFilename() string {
	return fmt.Sprintf("%s_%s_%s.sql", m.mysqlConf.PrefixFilename, m.mysqlConf.DBName, m.getTimeInFormat())
}

func (m *mysqlBackup) getTimeInFormat() string {
	now := time.Now()
	switch m.mysqlConf.TimeFormat {
	case backup.TimestampFormat:
		return strconv.FormatInt(now.Unix(), 10)
	case backup.DateOnlyFormat:
		return now.Format(time.DateOnly)
	case backup.DateTimeFormat:
		return now.Format(time.DateTime)
	default:
		return strconv.FormatInt(now.Unix(), 10)
	}
}

func (m *mysqlBackup) GetInterval() time.Duration {
	return time.Second * time.Duration(m.mysqlConf.Interval)
}

func (m *mysqlBackup) GetName() string {
	return m.name
}
