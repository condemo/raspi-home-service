package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// TODO: Podría implementar un método en cada tipo de configuración
// para guardar cambios realizados en al config
// Serian métodos con pointers para hacer persistentes los cambios
// y para que esos cambios se vean reflejados en todo el programa
// en tiempo real

var (
	Envs    = initConfig()
	APIConf = initApiConf()
)

type apiConfig struct {
	InfoTick time.Ticker
}

func initApiConf() apiConfig {
	return apiConfig{InfoTick: *time.NewTicker(2 * time.Second)}
}

type envConfig struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

// TODO: Mejorar la estructura y plantear añadir fallbacks
func initConfig() envConfig {
	godotenv.Load()
	return envConfig{
		DBHost: os.Getenv("DB_HOST"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
		DBPort: os.Getenv("DB_PORT"),
	}
}
