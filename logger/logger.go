package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm/module/apmzerolog/v2"
	"go.elastic.co/apm/v2"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var once sync.Once
var log zerolog.Logger
var es *elasticsearch.Client

type Fields struct {
	Ctx      context.Context `json:"context"`
	Id       string          `json:"id"`
	RemoteIp string          `json:"remote_ip"`
	Host     string          `json:"host"`
	Method   string          `json:"method"`
	Header   string          `json:"header"`
	Request  string          `json:"request"`
	Response string          `json:"response"`
	Protocol string          `json:"protocol"`
	StartAt  time.Time       `json:"start_at"`
	Path     string          `json:"path"`
	Error    error           `json:"error"`
}

type CustomLogger struct {
	zerolog.Logger
}

func (l *CustomLogger) LogRoundTrip(
	req *http.Request,
	res *http.Response,
	err error,
	start time.Time,
	dur time.Duration,
) error {
	var (
	//e *zerolog.Event
	//nReq int64
	//nRes int64
	)

	// Set error level.
	//
	//switch {
	//case err != nil:
	//	e = l.Error()
	//case res != nil && res.StatusCode > 0 && res.StatusCode < 300:
	//	e = l.Info()
	//case res != nil && res.StatusCode > 299 && res.StatusCode < 500:
	//	e = l.Warn()
	//case res != nil && res.StatusCode > 499:
	//	e = l.Error()
	//default:
	//	e = l.Error()
	//}

	// Count number of bytes in request and response.
	//
	//if req != nil && req.Body != nil && req.Body != http.NoBody {
	//	nReq, _ = io.Copy(ioutil.Discard, req.Body)
	//}
	//if res != nil && res.Body != nil && res.Body != http.NoBody {
	//	nRes, _ = io.Copy(ioutil.Discard, res.Body)
	//}

	return nil
}

// RequestBodyEnabled makes the client pass request body to logger
func (l *CustomLogger) RequestBodyEnabled() bool { return true }

// RequestBodyEnabled makes the client pass response body to logger
func (l *CustomLogger) ResponseBodyEnabled() bool { return true }

func Log(ctx context.Context) zerolog.Logger {
	once.Do(func() {
		var isDebugMode bool
		isDebugMode, _ = strconv.ParseBool(os.Getenv("DEBUG_MODE"))
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if isDebugMode {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		var gitRevision string
		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		zerolog.ErrorStackMarshaler = apmzerolog.MarshalErrorStack
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			pwd, _ := os.Getwd()
			file = strings.ReplaceAll(file, pwd+"/", "")
			return file + ":" + strconv.Itoa(line)
		}

		runLogFile, _ := os.OpenFile(
			"myapp.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		multi := zerolog.MultiLevelWriter(os.Stdout, runLogFile, &apmzerolog.Writer{})

		log = zerolog.New(multi).
			With().
			Stack().
			Timestamp().
			Caller().
			Str("git_revision", gitRevision).
			Str("go_version", buildInfo.GoVersion).
			Logger()
		es, _ = elasticsearch.NewClient(elasticsearch.Config{
			Logger: &CustomLogger{log},
		})
	})
	log.Hook(apmzerolog.TraceContextHook(ctx))

	return log
}

func LogError(ctx *fiber.Ctx, err error) {
	var logs *zerolog.Event
	var l = Log(ctx.Context())
	logs = l.Error()
	tx := apm.TransactionFromContext(ctx.Context())

	logs.
		Str("path", ctx.Path()).
		Str("trace_id", tx.TraceContext().Trace.String()).
		Err(err)
}

func LogInfo(ctx *fiber.Ctx, msg string) {
	var logs *zerolog.Event
	var l = Log(ctx.Context())
	logs = l.Info()

	tx := apm.TransactionFromContext(ctx.Context())

	logs.
		Str("path", ctx.Path()).
		Str("trace_id", tx.TraceContext().Trace.String()).
		Msg(msg)
}

func LogDebug(ctx *fiber.Ctx, msg string) {
	var logs *zerolog.Event
	var l = Log(ctx.Context())
	logs = l.Debug()

	ctx.Request().Header.VisitAll(func(key, val []byte) {
		k := bytes.NewBuffer(key).String()
		if k == "Elastic-Apm-Traceparent" || k == "Traceparent" {
			logs.Str(k, bytes.NewBuffer(val).String())
		}
	})

	logs.
		Str("path", ctx.Path()).
		Msg(msg)
}

//func LogOutboundRequest(path string, header []byte, request []byte, response []byte, statusCode int, fields Fields, errs []error) {
//	defer catch()
//	var isDebugMode bool
//	isDebugMode, _ = strconv.ParseBool(os.Getenv("DEBUG_MODE"))
//
//	var l = Log(ctx.)
//	var logs *zerolog.Event
//
//	switch strconv.Itoa(statusCode)[0] {
//	case '1', '2', '3':
//		logs = l.Info()
//	case '4', '5':
//		l.Error()
//		logs = l.Error()
//	default:
//		logs = l.Panic()
//	}
//
//	if len(errs) != 0 {
//		logs.Err(errors.Wrap(errs[0], "Reason"))
//	}
//
//	if !isDebugMode {
//		header = []byte("{}")
//	}
//
//	if len(request) != 0 {
//		logs.RawJSON("request", request)
//	}
//
//	if json.Valid(response) {
//		logs.RawJSON("response", response)
//	} else {
//		logs.Str("response", string(response))
//	}
//
//	logs.
//		Str("path", path).
//		RawJSON("header", header).
//		Int("status_code", statusCode).
//		Str("id", fields.Id).
//		Str("method", fields.Method).
//		Str("protocol", fields.Protocol).
//		Float64("latency", time.Since(fields.StartAt).Seconds()).
//		Msg(fmt.Sprintf("%v %v callback %v", statusCode, fields.Method, path))
//
//	headers := `{
//      "Content-Type": "application/json",
//      "X-Real-Ip": "200",
//      "X-Forwarded-Proto": "https",
//      "X-Forwarded-For": "erer",
//      "X-Forwarded-Prefix": "/api/rest/",
//      "X-Forwarded-Host": "staging.id",
//      "Accept-Encoding": "gzip,deflate",
//      "X-Forwarded-Path": "/tst,
//      "Token": "tst",
//      "Connection": "keep-alive",
//      "X-Forwarded-Port": "443",
//      "Host": "10.188.0.7:2345",
//      "User-Agent": "Apache-HttpClient/4.5.6 (Java/1.8.0_342)",
//      "Content-Length": "263",      "Accept": "application/xml, text/xml, application/json, application/*+xml, application/*+json"
//    }
//     `
//	requests := `
//{
//      "virtualAccountNumber": "8902110000000535",
//      "paymentReff": "TEST-inquiry4-8902110000000535",
//      "paymentDate": "13/03/2023 09:57:33",
//      "transactionId": "2023010900009",
//      "paymentMethod": "VIRTUAL_ACCOUNT_SINARMAS",
//      "paymentAmount": 400000,
//      "settlementDate": "13/03/2023 09:57:33"
//    }`
//
//	responses := `
//{
//      "responseMessage": "DATA_NOT_FOUND",
//      "responseCode": "752"
//    }
//`
//
//	fmt.Println("test")
//	fieldss := Fields{
//		Id:       "test-id",
//		RemoteIp: "test-remote-ip",
//		Host:     "localhost",
//		Method:   "POST",
//		Protocol: "HTTP",
//		StartAt:  time.Now(),
//		Request:  requests,
//		Response: responses,
//		Header:   headers,
//		Path:     "test-path",
//		Error:    nil,
//	}
//
//	result, _ := json.Marshal(fieldss)
//
//	logs.
//		Int("status_code", 200).
//		Str("method", fields.Method).
//		Str("protocol", fields.Protocol).
//		Float64("latency", time.Since(fields.StartAt).Seconds()).
//		Msg(fmt.Sprintf("%v %v callback %v", statusCode, fields.Method))
//	es.Index("mantap", strings.NewReader(string(result)), es.Index.WithRefresh("true"))
//
//}

func LogMiddleware(
	ctx context.Context,
	path string,
	header []byte,
	queryParam []byte,
	request []byte,
	response []byte,
	statusCode int,
	fields Fields,
	err error,
) {
	defer catch(ctx)
	var l = Log(ctx)
	var logs *zerolog.Event
	switch strconv.Itoa(statusCode)[0] {
	case '1', '2', '3':
		logs = l.Info()
	case '4', '5':
		l.Error()
		logs = l.Error()
	default:
		logs = l.Panic()
	}

	tx := apm.TransactionFromContext(ctx)

	if err != nil {
		logs.Err(errors.Wrap(err, "Reason"))
	}

	if len(queryParam) != 0 {
		logs.Str("query_param", string(queryParam))
	}

	if len(request) != 0 {
		dst := &bytes.Buffer{}
		if err := json.Compact(dst, request); err != nil {
			panic(err)
		}

		logs.RawJSON("request", dst.Bytes())
	}

	fields.Path = path
	fields.Response = string(response)
	fields.Request = string(request)
	fields.Header = string(header)
	fields.Error = err

	result, _ := json.Marshal(fields)

	logs.
		Str("path", "test").
		Str("trace_id", tx.TraceContext().Trace.String()).
		RawJSON("header", header).
		RawJSON("response", response).
		Interface("error", err).
		Int("status_code", statusCode).
		Str("id", fields.Id).
		Str("remote_ip", fields.RemoteIp).
		Str("host", fields.Host).
		Str("method", fields.Method).
		Str("protocol", fields.Protocol).
		Float64("latency", time.Since(fields.StartAt).Seconds()).
		Msg(fmt.Sprintf("%v %v %v", statusCode, fields.Method, path))
	es.Index("mantap", strings.NewReader(string(result)), es.Index.WithRefresh("true"))

}

func catch(ctx context.Context) {
	rvr := recover()
	if rvr != nil {
		var ok bool
		l := Log(ctx)

		err, ok := rvr.(error)
		if !ok {
			err = fmt.Errorf("%v", rvr)
		}

		l.Error().Err(err).Msg(err.Error())
	}
}
