package middlewares

import (
	"bytes"
	"chatproxy/appconfigs"
	"chatproxy/util"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

func RestyOnBeforeRequest(c *resty.Client, req *resty.Request) error {

	entry := log.WithFields(log.Fields{
		"time":    time.Now().Format(time.RFC3339),
		"method":  req.Method,
		"uri":     req.URL,
		"headers": req.Header,
		"body":    req,
	})

	entry.Info("Resty OnBeforeRequest")

	// Explore trace info
	/*
		fmt.Println("Request Trace Info:")
		ti := req.TraceInfo()
		fmt.Println("  DNSLookup     :", ti.DNSLookup)
		fmt.Println("  ConnTime      :", ti.ConnTime)
		fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
		fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
		fmt.Println("  ServerTime    :", ti.ServerTime)
		fmt.Println("  ResponseTime  :", ti.ResponseTime)
		fmt.Println("  TotalTime     :", ti.TotalTime)
		fmt.Println("  IsConnReused  :", ti.IsConnReused)
		fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
		fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
		fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	*/
	return nil // if its success otherwise return error
}

func RestyOnAfterResponse(c *resty.Client, resp *resty.Response) error {
	// Explore response object
	entry := log.WithFields(log.Fields{
		"status code": resp.StatusCode(),
		"time":        time.Now().Format(time.RFC3339),
		"Status":      resp.Status(),
		"Proto":       resp.Proto(),
		"headers":     resp.Header(),
		"body":        resp,
	})

	entry.Info("Resty OnBeforeResponse")
	// Explore trace info
	/*
		fmt.Println("Request Trace Info:")
		ti := resp.Request.TraceInfo()
		fmt.Println("  DNSLookup     :", ti.DNSLookup)
		fmt.Println("  ConnTime      :", ti.ConnTime)
		fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
		fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
		fmt.Println("  ServerTime    :", ti.ServerTime)
		fmt.Println("  ResponseTime  :", ti.ResponseTime)
		fmt.Println("  TotalTime     :", ti.TotalTime)
		fmt.Println("  IsConnReused  :", ti.IsConnReused)
		fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
		fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
		fmt.Println("  RequestAttempt:", ti.RequestAttempt)
		if ti.RemoteAddr != nil {
			fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
		}*/
	return nil // if its success otherwise return error
}

// ServerHeader middleware
// - Called before request handling
// - Sets "x" response header
// - Calls next handler in chain
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(util.HttpResponseHeaderRecipientKey, util.HttpResponseHeaderRecipient)
		return next(c)
	}
}

// LoggingMiddleware middleware
// - Called before request handling
// - Logs request method and URI
// - Uses logrus structured logging
// - Logs to configured outputs (file, stdout)
// - Calls next handler in chain
func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract fields
		method := c.Request().Method
		uri := c.Request().RequestURI
		ip := c.RealIP()

		// Extract other useful fields
		headers := c.Request().Header
		// Request
		reqBody := []byte{}
		if c.Request().Body != nil { // Read
			reqBody, _ = io.ReadAll(c.Request().Body)
		}
		c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset
		// Log request
		entry := log.WithFields(log.Fields{
			"time":    time.Now().Format(time.RFC3339),
			"method":  method,
			"uri":     uri,
			"ip":      ip,
			"headers": headers,
			"body":    string(reqBody),
		})

		entry.Info("Handled request")
		return next(c)
	}
}

func ValidateSignatureMiddleware(signingKey []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			//url := fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, req.URL.RequestURI())
			url := fmt.Sprintf("%s", req.URL.RequestURI())
			expectedSignature := req.Header.Get(util.HttpSignatureHeader)
			validator := hmac.New(sha1.New, []byte(signingKey))
			validator.Write([]byte(url))
			actualSignature := base64.URLEncoding.EncodeToString(validator.Sum(nil))
			if actualSignature != expectedSignature {
				return echo.ErrForbidden
			}
			return next(c)
		}
	}
}

func AppContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ac := appconfigs.GetContextByHost(c.Request().Host)
		c.Set(util.EchoAppContext, ac)
		return next(c)
	}
}
