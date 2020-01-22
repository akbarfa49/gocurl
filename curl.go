package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
)

const (
	MethodPost    = "POST"
	MethodGet     = "GET"
	MethodOptions = "Options"
	MethodPut     = "Put"
	MethodDelete  = "Delete"
)

//from NET/HTTP
const (
	StatusContinue           = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols = 101 // RFC 7231, 6.2.2
	StatusProcessing         = 102 // RFC 2518, 10.1
	StatusEarlyHints         = 103 // RFC 8297

	StatusOK                   = 200 // RFC 7231, 6.3.1
	StatusCreated              = 201 // RFC 7231, 6.3.2
	StatusAccepted             = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
	StatusNoContent            = 204 // RFC 7231, 6.3.5
	StatusResetContent         = 205 // RFC 7231, 6.3.6
	StatusPartialContent       = 206 // RFC 7233, 4.1
	StatusMultiStatus          = 207 // RFC 4918, 11.1
	StatusAlreadyReported      = 208 // RFC 5842, 7.1
	StatusIMUsed               = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices  = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently = 301 // RFC 7231, 6.4.2
	StatusFound            = 302 // RFC 7231, 6.4.3
	StatusSeeOther         = 303 // RFC 7231, 6.4.4
	StatusNotModified      = 304 // RFC 7232, 4.1
	StatusUseProxy         = 305 // RFC 7231, 6.4.5

	StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect = 308 // RFC 7538, 3

	StatusBadRequest                   = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                 = 401 // RFC 7235, 3.1
	StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
	StatusForbidden                    = 403 // RFC 7231, 6.5.3
	StatusNotFound                     = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
	StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
	StatusConflict                     = 409 // RFC 7231, 6.5.8
	StatusGone                         = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable = 416 // RFC 7233, 4.4
	StatusExpectationFailed            = 417 // RFC 7231, 6.5.14
	StatusTeapot                       = 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest           = 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity          = 422 // RFC 4918, 11.2
	StatusLocked                       = 423 // RFC 4918, 11.3
	StatusFailedDependency             = 424 // RFC 4918, 11.4
	StatusTooEarly                     = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         = 428 // RFC 6585, 3
	StatusTooManyRequests              = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

	StatusInternalServerError           = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           = 507 // RFC 4918, 11.5
	StatusLoopDetected                  = 508 // RFC 5842, 7.2
	StatusNotExtended                   = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)

type Client struct {
	Method string
	Url    string
	Body   interface{}
	//http version default is 1.0
	Header  []string
	Version string
}
type Response struct {
	Status int
	Header map[string]string
	Body   interface{}
}

var crlf = []byte("\r\n")

func New() *Client {
	return &Client{
		Version: "HTTP/1.0",
	}
}
func (c *Client) setHeader(key, value string) {
	c.Header = append(c.Header, fmt.Sprintf("%v:%v", key, value))
}

func (c *Client) Curl(method, url string, body interface{}) {
	c.Method = method
	c.Url = url
	if body != nil {
		c.Body = body
	}
}

func (c *Client) Exec() (*Response, error) {

	var wrap [][]byte
	var status []byte
	var interf interface{}
	var head []map[string]string
	//var tmp map[string]string
	d := new(net.Dialer)
	r := new(Response)
	r.Header = map[string]string{}
	u, _ := url.Parse(c.Url)

	conn, err := d.Dial("tcp", u.Host)
	if err != nil {
		return nil, err
	}
	if u.Path == "" {
		u.Path = "/"
	}
	bayt := []byte(fmt.Sprintf("%v %v %v\r\n", c.Method, u.Path, c.Version))
	for i := range c.Header {
		head := []byte(c.Header[i])
		for j := range head {
			bayt = append(bayt, head[j])
		}
		if i != len(c.Header)-1 {
			for j := range crlf {
				bayt = append(bayt, crlf[j])
			}
		} else {
			for j := 0; j < 2; j++ {
				for k := range crlf {
					bayt = append(bayt, crlf[k])
				}

			}
		}
	}
	if c.Body != nil {
		body, _ := json.Marshal(c.Body)
		for i := range body {
			bayt = append(bayt, body[i])
		}
	}
	conn.Write(bayt)
	bufr := bufio.NewReader(conn)
	bufs := bufio.NewScanner(bufr)
	//bufs.Split()
	for bufs.Scan() {
		wrap = append(wrap, bufs.Bytes())
	}
	for i := range wrap[0] {
		if wrap[0][i] == 32 {
			status = append(status, wrap[0][i+1])
			status = append(status, wrap[0][i+2])
			status = append(status, wrap[0][i+3])
			break
		}
	}
	for i := 1; i < len(wrap); i++ {
		if i == len(wrap)-1 {
			json.Unmarshal(wrap[i], &interf)
		} else if i == len(wrap)-2 {
			continue
		} else {
			head = append(head, toHeader(wrap[i]))
		}
	}
	for i := 0; i < len(head); i++ {
		j := head[i]
		for k, l := range j {
			r.Header[k] = l
		}
	}
	stats, _ := strconv.Atoi(string(status))
	r.Status = stats
	r.Body = interf
	return r, nil
}

func toHeader(b []byte) map[string]string {
	var key []byte
	var value []byte
	for i := range b {
		if b[i] == 58 {
			for j := i + 2; j < len(b); j++ {
				value = append(value, b[j])
			}
			break
		}
		key = append(key, b[i])
	}
	a := map[string]string{string(key): string(value)}
	return a
}
