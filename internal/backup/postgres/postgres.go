package postgres

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/JCoupalK/go-pgdump"
	"github.com/Naumovets/Backuper/internal/backup"
	"github.com/Naumovets/Backuper/internal/config"
	"github.com/Naumovets/Backuper/internal/logger"
	"go.uber.org/zap"
)

type postgresBackup struct {
	dumper      *pgdump.Dumper
	pgConf      config.PostgresConfig
	storageConf config.StorageConfig
	name        string
}

func NewBackuper(pgConf config.PostgresConfig, storageConfig config.StorageConfig, name string) *postgresBackup {
	sslmode := "disable"
	if pgConf.SSLMode != nil {
		sslmode = *pgConf.SSLMode
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pgConf.Host, pgConf.Port, pgConf.User, pgConf.Password, pgConf.DBName, sslmode)

	if pgConf.ApplicationName != nil {
		psqlInfo += fmt.Sprintf(" application_name=%s", *pgConf.ApplicationName)
	}

	if pgConf.SSLCert != nil {
		psqlInfo += fmt.Sprintf(" sslcert=%s", *pgConf.SSLCert)
	}

	if pgConf.SSLKey != nil {
		psqlInfo += fmt.Sprintf(" sslkey=%s", *pgConf.SSLKey)
	}

	if pgConf.SSLRootCert != nil {
		psqlInfo += fmt.Sprintf(" sslrootcert=%s", *pgConf.SSLRootCert)
	}

	dumper := pgdump.NewDumper(psqlInfo, pgConf.Threads)
	return &postgresBackup{
		name:        name,
		dumper:      dumper,
		pgConf:      pgConf,
		storageConf: storageConfig,
	}
}

func (p *postgresBackup) MakeBackup(ctx context.Context) error {
	log := logger.GetLogger(ctx)

	log.Info("Start make backup database")

	targetDir := filepath.Join(p.storageConf.Path, p.name)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	dumpFilename := filepath.Join(targetDir, p.getFilename())

	opts := &pgdump.TableOptions{}
	if p.pgConf.TableSuffix != nil {
		opts.TableSuffix = *p.pgConf.TableSuffix
	}

	if p.pgConf.TablePrefix != nil {
		opts.TablePrefix = *p.pgConf.TablePrefix
	}

	if p.pgConf.TableSchema != nil {
		opts.Schema = *p.pgConf.TableSchema
	}

	if err := p.dumper.DumpDatabase(dumpFilename, opts); err != nil {
		return fmt.Errorf("error dumping database: %w", err)
	}

	if p.pgConf.Compress {
		log.Info("Compressing backup")
		compressedFile, err := p.compressFile(dumpFilename)
		if err != nil {
			log.Error("Failed to compress backup", zap.Error(err))
		} else {
			os.Remove(dumpFilename)
			dumpFilename = compressedFile
		}
	}

	log.Info("Database is dumped", zap.String("file", dumpFilename))

	// Run retention policy
	if err := p.runRetention(ctx, targetDir); err != nil {
		log.Error("Retention policy failed", zap.Error(err))
	}

	return nil
}

func (p *postgresBackup) compressFile(filename string) (string, error) {
	newFilename := filename + ".gz"
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	out, err := os.Create(newFilename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Use specified compression level if valid, otherwise default
	level := p.pgConf.CompressLevel
	if level < 0 || level > 9 {
		level = gzip.DefaultCompression
	}

	gw, err := gzip.NewWriterLevel(out, level)
	if err != nil {
		return "", err
	}
	defer gw.Close()

	if _, err := io.Copy(gw, f); err != nil {
		return "", err
	}

	return newFilename, nil
}

func (p *postgresBackup) runRetention(ctx context.Context, targetDir string) error {
	log := logger.GetLogger(ctx)

	if p.pgConf.MaxCount <= 0 {
		return nil
	}

	files, err := os.ReadDir(targetDir)
	if err != nil {
		return err
	}

	var backupFiles []os.DirEntry
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), p.pgConf.PrefixFilename) {
			backupFiles = append(backupFiles, f)
		}
	}

	if len(backupFiles) <= p.pgConf.MaxCount {
		return nil
	}

	// Sort by info (modification time) - oldest first
	sort.Slice(backupFiles, func(i, j int) bool {
		infoI, _ := backupFiles[i].Info()
		infoJ, _ := backupFiles[j].Info()
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	toDelete := len(backupFiles) - p.pgConf.MaxCount
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

func (p *postgresBackup) getFilename() string {
	return fmt.Sprintf("%s_%s_%s.sql", p.pgConf.PrefixFilename, p.pgConf.DBName, p.getTimeInFormat())
}

func (p *postgresBackup) getTimeInFormat() string {
	now := time.Now()
	switch p.pgConf.TimeFormat {
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

func (p *postgresBackup) GetInterval() time.Duration {
	return time.Second * time.Duration(p.pgConf.Interval)
}

func (p *postgresBackup) GetName() string {
	return p.name
}
