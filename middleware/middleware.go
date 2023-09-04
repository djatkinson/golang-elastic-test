package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmhttp/v2"
	"net/http"
	"splunk-test/logger"
	"time"
)

func Config() fiber.Handler {

	return func(c *fiber.Ctx) error {
		rid := c.Get(fiber.HeaderXRequestID)

		if rid == "" {
			rid = uuid.New().String()
			c.Set(fiber.HeaderXRequestID, rid)
		}

		header := make(map[string]string)
		c.Request().Header.VisitAll(func(key, val []byte) {
			k := bytes.NewBuffer(key).String()
			if k == apmhttp.W3CTraceparentHeader || k == apmhttp.ElasticTraceparentHeader {
				fmt.Println("test-middleware", c.Path(), k, bytes.NewBuffer(val).String())
			}

			header[k] = bytes.NewBuffer(val).String()
		})

		startAt := time.Now()

		defer func() {
			var err error
			rvr := recover()
			if rvr != nil {
				var ok bool
				err, ok = rvr.(error)
				if !ok {
					err = fmt.Errorf("%v", rvr)
				}
				err = errors.New(fmt.Sprintf("%v", rvr))

				c.Status(http.StatusInternalServerError).JSON(`{"error": "Internal Server Error"}`)
			}

			headerByte, _ := json.Marshal(header)
			fields := logger.Fields{
				Id:       rid,
				RemoteIp: c.IP(),
				Method:   c.Method(),
				Host:     c.Hostname(),
				Protocol: c.Protocol(),
				StartAt:  startAt,
				Ctx:      c.Context(),
			}

			logger.LogMiddleware(
				c.Context(),
				c.Path(),
				headerByte,
				c.Request().URI().QueryString(),
				c.Body(),
				c.Response().Body(),
				c.Response().StatusCode(),
				fields,
				err)
		}()

		return c.Next()
	}
}
