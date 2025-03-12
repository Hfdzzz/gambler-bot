package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Gagal memuat file .env")
	} else {
		fmt.Println("File .env berhasil dimuat")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
