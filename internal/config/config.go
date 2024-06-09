package config


import (
	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
	HostDB   	string `env:"DB_HOST" env-default:"localhost"`
    PortDB   	string `env:"DB_PORT" env-default:"5432"`
	UserDB		string `env:"DB_USER" env-default:"admin"`
	PasswordDB 	string `env:"DB_PASSWORD" env-default:"admin"`
	NameDB 		string `env:"DB_NAME" env-default:"posts"`

	RepoType    string `env:"REPO_TYPE" env-default:"in-memory"`

	MaxComLen   int    `env:"MAX_COM_LEN" env-default:2000`

	ServerPort  string `env:"PORT" env-default:"6055"`
}


func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
