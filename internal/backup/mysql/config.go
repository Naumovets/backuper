package mysql

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/Naumovets/backuper/internal/config"
	"github.com/go-sql-driver/mysql"
)

func ToMysqlConfig(c *config.MysqlConfig) mysql.Config {
	cfg := mysql.NewConfig()
	cfg.User = c.User
	cfg.Passwd = c.Password
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%d", c.Host, c.Port)
	cfg.DBName = c.DBName

	cfg.Params = make(map[string]string)
	if c.ApplicationName != nil {
		cfg.Params["program_name"] = *c.ApplicationName
	}
	cfg.Params["charset"] = "utf8mb4"
	cfg.Params["parseTime"] = "true"

	// TLS (по SSLMode)
	if c.SSLMode != nil {
		switch *c.SSLMode {
		case "disable":
			cfg.AllowFallbackToPlaintext = true
		case "require":
			// TLS без CA
		case "verify-ca", "verify-full":
			tlsCfg := &tls.Config{}
			if c.SSLCert != nil {
				tlsCfg.Certificates = make([]tls.Certificate, 1)
				// tlsCfg.LoadX509KeyPair(*c.SSLCert, *c.SSLKey) — если файлы
			}
			if c.SSLRootCert != nil {
				// tlsCfg.RootCAs
			}
			cfg.TLS = tlsCfg
		}
	}

	cfg.InterpolateParams = true
	cfg.ParseTime = true
	cfg.Loc = time.UTC

	return *cfg
}
