package config

type Config struct {
	JWTKey string `env:"JWT_KEY" envDefault:"supersecret"`
	Logger Logger
}

type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

type HTTP struct {
	AppPort string `env:"APP_PORT" envDefault:"8000"`
}

type Database struct {
	DBHost         string `env:"MONGO_HOST" envDefault:"studentsdb"`
	DBPort         string `env:"MONGO_PORT" envDefault:"27017"`
	DBName         string `env:"DATABASE_NAME" envDefault:"studentsdb"`
	CollectionName string `env:"DATABASE_NAME" envDefault:"students"`
}
