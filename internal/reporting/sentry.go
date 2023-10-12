package reporting

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
)

func InitSentry(dsn string, env string) {
	if dsn == "" {
		log.Info().Msg("Sentry DSN is not specified. Sentry will not be initialized.")
		return	
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
		// TODO configure this
		Environment: env,
	})
	if err != nil {
		log.Error().Err(err).Msg("Could not initialize sentry.")
	} else {
		log.Info().Msg("Sentry initialized.")
	}
}

func FlushSentry() {
	sentry.Flush(2 * time.Second)
}

func CaptureErrorWithContext(err error, contextKey string, contextValue map[string]interface{}) {
	sentry.WithScope(func(scope *sentry.Scope) { 
		scope.SetContext(contextKey, contextValue)
		sentry.CaptureException(err)
	})
}
