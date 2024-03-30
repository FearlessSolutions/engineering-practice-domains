package middleware

import (
	"example.com/sample/commonlib/auth"
	"example.com/sample/commonlib/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// LoggingMiddleware constructs an echo middleware which logs incoming requests on the global logger - logger.Log
func LoggingMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:       true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,

		LogValuesFunc: func(ctx echo.Context, logValues middleware.RequestLoggerValues) error {
			fieldsToLog := []zap.Field{
				zap.Duration("latency", logValues.Latency),
				zap.String("remoteIP", logValues.RemoteIP),
				zap.String("host", logValues.Host),
				zap.String("method", logValues.Method),
				zap.String("uri", logValues.URI),
				zap.String("userAgent", logValues.UserAgent),
				zap.String("contentLength", logValues.ContentLength),
				zap.Int("responseCode", ctx.Response().Status),
			}

			if logValues.Error != nil {
				fieldsToLog = append(fieldsToLog, zap.Error(logValues.Error))
			}
			if claims, claimsPresent := auth.RetrieveAuthClaims(ctx); claimsPresent {
				fieldsToLog = append(fieldsToLog, zap.String("requesterUsername", claims.PreferredUsername))
			}

			logger.Log.Info("request info", fieldsToLog...)
			return nil
		},
	})
}
