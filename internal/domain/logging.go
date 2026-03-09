package domain

type Logging struct {
	Level            string   `json:"level"`
	Encoding         string   `json:"encoding"`
	OutputPaths      []string `json:"output_paths"`
	ErrorOutputPaths []string `json:"error_output_paths"`

	RotationSettings *RotationSettings `json:"rotation_settings,omitempty"`
}

type RotationSettings struct {
	MaxSize  int  `json:"max_size"`
	MaxCount int  `json:"max_count"`
	MaxAge   int  `json:"max_age"`
	Compress bool `json:"compress"`
}
