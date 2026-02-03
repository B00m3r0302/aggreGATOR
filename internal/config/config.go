package config

import (
	"os"
)

type Config struct {
	DbUrl string `json:"db_url"`
}

func Read() Config {

}
