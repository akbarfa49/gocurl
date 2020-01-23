package gocurl

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
