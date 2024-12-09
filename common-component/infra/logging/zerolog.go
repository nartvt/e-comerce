package logging

import (
	"github.com/rs/zerolog"
)

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimestampFieldName
}
