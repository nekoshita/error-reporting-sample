// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var stackTraceMsg = `panic: hoge
	goroutine 1 [running]:
	main.main()
			/github.com/nekoshita/error-reporting-sample/main.go:12 +0x2c
	exit status 2`

var stackTraceMsgWith = func(msg string, line int) string {
	return fmt.Sprintf(`panic: %s
	goroutine 1 [running]:
	main.main()
			/github.com/nekoshita/error-reporting-sample/main.go:%d +0x2c
	exit status 2`, msg, line)
}

var exceptionMsg = `java.lang.CloneNotSupportedException
	at sample.24.SmapleUserList.Loging(SampleUserList.java:27)`

var exceptionMsgWith = func(msg string) string {
	return fmt.Sprintf(`java.lang.CloneNotSupportedException.%s`, msg)
}

func main() {
	ctx := context.Background()
	log := NewLogger(ctx).WithOptions(zap.AddStacktrace(zap.DPanicLevel))

	log.Info("starting server...")
	http.HandleFunc("/", handler)
	http.HandleFunc("/panic", handlerPanic)
	http.HandleFunc("/panic2", handlerPanic2)
	http.HandleFunc("/error", handlerErrorMsg)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Info(fmt.Sprintf("defaulting to port %s", port))
	}

	// Start HTTP server.
	log.Info(fmt.Sprintf("listening on port %s", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

func handlerErrorMsg(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	if msg == "" {
		msg = "empty msg"
	}

	ctx := context.Background()
	log := NewLogger(ctx).WithOptions(zap.AddStacktrace(zap.DPanicLevel))

	// この3つは同じグループになる
	// なぜならstack_traceのファイルが同じだから？
	log.Error("error message with stack trace", []zap.Field{
		zap.Any("stack_trace", stackTraceMsgWith(msg, 89)),
	}...)
	log.Error("error message with stack trace in exception filed", []zap.Field{
		zap.Any("exception", stackTraceMsgWith(msg, 67)),
	}...)
	log.Error(stackTraceMsgWith(msg, 24), []zap.Field{}...)

	// これはメッセージの内容が異なると違うグループになる
	log.Error(msg, []zap.Field{
		zap.Any("@type", "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"),
	}...)

	log.Error("error message with stack trace in exception filed", []zap.Field{
		zap.Any("@type", "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"),
		zap.Any("exception", exceptionMsgWith(msg)),
	}...)

	log.Error("これは収集されない", []zap.Field{
		zap.Any("stacktrace", stackTraceMsgWith(msg, 89)),
	}...)
}

func handlerPanic(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	if msg == "" {
		msg = "empty msg"
	}

	panic(msg)
}

func handlerPanic2(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	if msg == "" {
		msg = "empty msg"
	}

	panic(msg)
}

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  "DEBUG",
	zapcore.InfoLevel:   "INFO",
	zapcore.WarnLevel:   "WARNING",
	zapcore.ErrorLevel:  "ERROR",
	zapcore.DPanicLevel: "CRITICAL",
	zapcore.PanicLevel:  "ALERT",
	zapcore.FatalLevel:  "EMERGENCY",
}

// https://www.wheatandcat.me/entry/2022/07/12/232723
func EncodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logLevelSeverity[l])
}
func newProductionEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()

	cfg.TimeKey = "time"
	cfg.LevelKey = "severity"
	cfg.MessageKey = "message"
	cfg.StacktraceKey = "stack_trace"
	cfg.EncodeLevel = EncodeLevel
	cfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	return cfg
}
func newProductConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(zap.DebugLevel)
	cfg.EncoderConfig = newProductionEncoderConfig()

	return cfg
}
func NewLogger(ctx context.Context) *zap.Logger {
	cfg := newProductConfig()
	logger, _ := cfg.Build()

	return logger
}
