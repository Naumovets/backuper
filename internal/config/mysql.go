package config

type MysqlConfig struct {
	Threads         int     `yaml:"threads"`
	Host            string  `yaml:"host"`
	Port            int     `yaml:"port"`
	DBName          string  `yaml:"dbname"`
	User            string  `yaml:"user"`
	Password        string  `yaml:"password"`
	SSLKey          *string `yaml:"sslkey"`
	SSLCert         *string `yaml:"sslcert"`
	SSLRootCert     *string `yaml:"sslrootcert"`
	SSLMode         *string `yaml:"sslmode"`
	ApplicationName *string `yaml:"application_name"`
	Compress        bool    `yaml:"compress"`
	CompressLevel   int     `yaml:"compress_level"`
	Interval        int     `yaml:"internal"`
	MaxCount        int     `yaml:"max_count"`
	TimeFormat      string  `yaml:"timeformat"`
	PrefixFilename  string  `yaml:"prefix_filename"`
	TablePrefix     *string `yaml:"table_prefix"`
	TableSuffix     *string `yaml:"table_suffix"`
}
