package helper

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/faizauthar12/doku/app/models"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Options :
type Options struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Body        []byte            `json:"body"`
	Headers     map[string]string `json:"headers"`
	Timeout     time.Duration
	ContentType string                 `json:"content_type"`
	QueryParams map[string]interface{} `json:"query_params"`
	IsPrintCurl bool                   `json:"is_print_curl"`
}

// Response :
type Response struct {
	StatusCode int    `json:"status-code"`
	Body       []byte `json:"body"`
	Error      error  `json:"error"`
}

func GET(opt *Options) Response {
	res := <-request(opt, "GET")
	return res
}

func POST(opt *Options) Response {
	res := <-request(opt, "POST")
	return res
}

func PUT(opt *Options) Response {
	res := <-request(opt, "PUT")
	return res
}

func DELETE(opt *Options) Response {
	res := <-request(opt, "DELETE")
	return res
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Request :
func request(opt *Options, method string) <-chan Response {
	res := make(chan Response)
	go func() {
		defer Recover("rest http")
		defer close(res)
		var rsp *http.Response
		var e error
		c := http.Client{
			Timeout:   opt.Timeout,
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		}
		logrus.Debugf("http request body : %s", opt.Body)
		clientReqHeader := http.Header{}
		for k, v := range opt.Headers {
			clientReqHeader.Add(k, v)
		}
		if opt.ContentType == "" {
			clientReqHeader.Add("Content-Type", "application/json")
		} else {
			clientReqHeader.Add("Content-Type", opt.ContentType)
		}

		reqObj := http.Request{}
		reqObj.Method = method
		reqObj.URL, _ = url.Parse(opt.URL)
		reqObj.Header = clientReqHeader
		reqObj.Body = io.NopCloser(bytes.NewBuffer(opt.Body))
		queryParams := reqObj.URL.Query()

		for queryKey, queryVal := range opt.QueryParams {
			myType := reflect.TypeOf(queryVal)

			if k := myType.Kind(); k == reflect.Int {
				queryValInt := queryVal.(int)
				queryValStr := strconv.Itoa(queryValInt)
				queryParams.Add(queryKey, queryValStr)
			} else if k == reflect.String {
				queryParams.Add(queryKey, queryVal.(string))
			}

		}
		reqObj.URL.RawQuery = queryParams.Encode()
		if opt.IsPrintCurl {
			PrintCurl(&reqObj, opt.Body)
		}
		rsp, e = c.Do(&reqObj)
		if e != nil {
			logrus.Debugf("error when creating http request %s", e.Error())
			res <- Response{Error: fmt.Errorf("failed to create new request")}
			return
		}
		defer rsp.Body.Close()
		body, _ := io.ReadAll(rsp.Body)
		res <- Response{StatusCode: rsp.StatusCode, Body: body}
	}()
	return res
}

func GetStatusCode(err error, statusCodeDefault int) int {
	var statusCode int
	if strings.Contains(err.Error(), "not found") {
		statusCode = http.StatusNotFound
	} else if strings.Contains(err.Error(), "already") {
		statusCode = http.StatusConflict
	} else if strings.Contains(err.Error(), "expired") {
		statusCode = http.StatusGone
	} else {
		statusCode = http.StatusInternalServerError
	}

	if statusCodeDefault > 0 {
		statusCode = statusCodeDefault
	}

	return statusCode
}
func GenerateResultByError(err error, statusCode int, systemMessage string) models.Response {
	errLog := WriteLog(err, statusCode, err.Error())
	if systemMessage == "" {
		systemMessage = err.Error()
	}
	errLog.SystemMessage = systemMessage
	return models.Response{
		StatusCode: statusCode,
		Error:      errLog,
	}
}

func GenerateResultByErrorWithMessage(err error, statusCode int, message interface{}) models.Response {
	return models.Response{
		StatusCode: statusCode,
		Error:      WriteLog(err, statusCode, message),
	}
}

func GenerateResultByErrorLog(err *models.ErrorLog) models.Response {
	return models.Response{
		StatusCode: err.StatusCode,
		Error:      err,
	}
}

func PrintCurl(httpReq *http.Request, payload []byte) {
	fmt.Printf("curl -X %s '%s'", httpReq.Method, httpReq.URL.String())
	for name, values := range httpReq.Header {
		for _, value := range values {
			fmt.Printf(" -H '%s: %s'", name, value)
		}
	}
	if len(payload) > 0 {
		fmt.Printf(" -d '%s'", string(payload))
	}
	fmt.Println()
}
