package middleware

import (
	"strings"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/constants"
)

func CaselessMatcherMiddleware(ctx iris.Context) {
	ctx.Request().URL.Path = strings.ToLower(ctx.Path())
	ctx.Next()
}

func RequestIDMiddleware(ctx iris.Context) {
	requestId := ctx.Request().Header.Get(constants.CTXRequestIdKey)
	if len(requestId) == 0 {
		requestId = uuid.NewString()
	}
	ctx.Values().Set(constants.CTXRequestIdKey, requestId)
	ctx.Next()
}

func MakeAccessLog() *accesslog.AccessLog {
	// Initialize a new access log middleware.
	var ac = accesslog.File("./access.log")
	// Remove this line to disable logging to console:
	//ac.AddOutput(os.Stdout)

	// The default configuration:
	ac.Delim = '|'
	ac.TimeFormat = "2006-01-02 15:04:05"
	ac.Async = false
	ac.IP = true
	ac.BytesReceivedBody = true
	ac.BytesSentBody = true
	ac.BytesReceived = false
	ac.BytesSent = false
	ac.BodyMinify = true
	ac.RequestBody = true
	ac.ResponseBody = false
	ac.KeepMultiLineError = true
	ac.PanicLog = accesslog.LogHandler

	// Default line format if formatter is missing:
	// Time|Latency|Code|Method|Path|IP|Path Params Query Fields|Bytes Received|Bytes Sent|Request|Response|
	//
	// Set Custom Formatter:
	ac.SetFormatter(&accesslog.JSON{
		Indent:    "  ",
		HumanTime: true,
	})
	// ac.SetFormatter(&accesslog.CSV{})
	// ac.SetFormatter(&accesslog.Template{Text: "{{.Code}}"})

	return ac
}
