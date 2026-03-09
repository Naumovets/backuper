package domain

type PostgresBackuper struct {
	Name            string  `json:"name"`
	Threads         int     `json:"threads"`
	Host            string  `json:"host"`
	Port            int     `json:"port"`
	DBName          string  `json:"dbname"`
	User            string  `json:"user"`
	Password        string  `json:"password"`
	SSLKey          *string `json:"sslkey,omitempty"`
	SSLCert         *string `json:"sslcert,omitempty"`
	SSLRootCert     *string `json:"sslrootcert,omitempty"`
	SSLMode         *string `json:"sslmode,omitempty"`
	ApplicationName *string `json:"application_name,omitempty"`
	Compress        bool    `json:"compress"`
	CompressLevel   int     `json:"compress_level"`
	Interval        int     `json:"internal"`
	MaxCount        int     `json:"max_count"`
	TimeFormat      string  `json:"timeformat"`
	PrefixFilename  string  `json:"prefix_filename"`
	TableSchema     *string `json:"table_schema,omitempty"`
	TablePrefix     *string `json:"table_prefix,omitempty"`
	TableSuffix     *string `json:"table_suffix,omitempty"`
}
