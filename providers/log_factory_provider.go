package providers

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"

	"github.com/google/wire"
)

func LoggerOutputProvider() zerolog.ConsoleWriter {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("**** %s ****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	return output
}

func LoggerProvider(output zerolog.ConsoleWriter) zerolog.Logger {
	return zerolog.New(output).With().Timestamp().Caller().Logger()
}

var LoggerFactoryProviderSet = wire.NewSet(LoggerOutputProvider, LoggerProvider)
