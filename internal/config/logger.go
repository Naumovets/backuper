package config

type LoggingConfig struct {
	Level            string   `yaml:"level"`
	Encoding         string   `yaml:"encoding"`
	OutputPaths      []string `yaml:"output_paths"`
	ErrorOutputPaths []string `yaml:"error_output_paths"`

	RotationSettings *RotationSettingsConfig `yaml:"rotation_settings"`
}

type RotationSettingsConfig struct {
	MaxSize  int  `yaml:"max_size"`
	MaxCount int  `yaml:"max_count"`
	MaxAge   int  `yaml:"max_age"`
	Compress bool `yaml:"compress"`
}
