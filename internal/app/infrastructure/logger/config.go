package logger

type Config struct {
	// LogToConsole enables console logging.
	LogToConsole bool `yaml:"log_to_console"`

	// EncodeLogsAsJson makes the log framework use JSON format.
	EncodeLogsAsJson bool `yaml:"encode_logs_as_json"`

	// LogToFile makes the framework log to a file.
	// The fields below can be skipped if this value is false.
	LogToFile bool `yaml:"log_to_file"`

	// Directory is the name of the log file output directory.
	Directory string `yaml:"directory"`

	// Filename is the name of the log file which will be placed inside the directory.
	Filename string `yaml:"filename"`

	// MaxSize is the max size in MB of the log file before it's rolled.
	MaxSize int `yaml:"max_size"`

	// MaxBackups is the max number of rolled files to keep.
	MaxBackups int `yaml:"max_backups"`

	// MaxAge is the max age in days to keep a log file.
	MaxAge int `yaml:"max_age"`
}
