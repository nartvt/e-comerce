package util

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found. Using exists environment variables")
	}
}

func SliceToJson(slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	resp, err := json.Marshal(slice)
	if err != nil {
		return ""
	}
	return string(resp)
}

func StructToJson(structAny any) string {
	resp, err := json.Marshal(structAny)
	if err != nil {
		return ""
	}
	return string(resp)
}

func UUID() string {
	id := uuid.New()
	return id.String()
}

func UUIDFunc() func() string {
	return func() string {
		return UUID()
	}
}
