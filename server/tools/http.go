package tools

import (
	"bytes"
	"encoding/base64"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Create client
var client = &http.Client{
	Timeout: 30 * time.Second,
}

// Opt HTTP选项
type Opt interface {
	before(logger ILog, req *http.Request) error
	after(logger ILog, resp *http.Response) error
}

// HeaderOpt HTTP请求头选项
type HeaderOpt struct {
	Header http.Header
}

// NewHeaderOpt 创建HTTP请求头选项
func NewHeaderOpt(header http.Header) *HeaderOpt {
	return &HeaderOpt{Header: header}
}

func (h *HeaderOpt) before(logger ILog, req *http.Request) error {
	for k, v := range h.Header {
		req.Header[k] = v
	}
	return nil
}

func (h *HeaderOpt) after(logger ILog, resp *http.Response) error {
	return nil
}

// JSONDecodeOpt HTTP请求头选项
type JSONDecodeOpt struct {
	Object    interface{}
	BodyBytes []byte
}

// NewJSONDecodeOpt 创建响应JSON解析选项
func NewJSONDecodeOpt(object interface{}) *JSONDecodeOpt {
	return &JSONDecodeOpt{Object: object}
}

func (h *JSONDecodeOpt) before(logger ILog, req *http.Request) error {
	return nil
}

func (h *JSONDecodeOpt) after(logger ILog, resp *http.Response) error {
	var err error
	if resp.Body == nil {
		return nil
	}
	h.BodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	return jsoniter.Unmarshal(h.BodyBytes, h.Object)
}

// HTTPGet 发送HTTP Get请求
func HTTPGet(logger ILog, surl string, params map[string]string, options ...Opt) (resp *http.Response, err error) {
	if !strings.HasPrefix(strings.ToLower(surl), "http") {
		return nil, logger.Error("url failed")
	}

	// Create request
	req, err := http.NewRequest("GET", surl, nil)
	if err != nil {
		return nil, logger.Error(err)
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	return HTTPDo(logger, req, options...)
}

// HTTPDo 发送HTTP请求
func HTTPDo(logger ILog, req *http.Request, options ...Opt) (resp *http.Response, err error) {
	furl := req.URL.String()
	logger.Info("do http ", req.Method, "  request:", furl)
	var dumped []byte
	dumped, err = httputil.DumpRequest(req, true)
	if err == nil {
		logger.Infof("do http %s request:%s", req.Method, dumped)
	} else {
		logger.Info("do http %s request:", req.Method, furl)
	}

	for _, opt := range options {
		err := opt.before(logger, req)
		if err != nil {
			return nil, logger.Error("set option failed:", err)
		}
	}

	t := time.Now()
	resp, err = client.Do(req)
	costTime := time.Now().Sub(t)
	if err != nil {
		logger.Error("http ", req.Method, " ", furl, "failed:", err, " cost ", costTime)
		return
	}
	dumped, err = httputil.DumpResponse(resp, true)
	if err != nil {
		logger.Error("http ", req.Method, " ", furl, "failed:", err, " cost ", costTime)
		return
	}

	logger.Infof("http %s %v cost %v response: %s", req.Method, furl, costTime, dumped)

	if resp.StatusCode != 200 {
		return resp, logger.Error("http ", req.Method, " ", furl, "failed, status code:", resp.StatusCode)
	}

	for _, opt := range options {
		err := opt.after(logger, resp)
		if err != nil {
			return nil, logger.Error("commit option failed:", err)
		}
	}

	return
}

// HTTPPostForm 发送HTTP Post Form请求
func HTTPPostForm(logger ILog, surl string, params map[string]string, options ...Opt) (resp *http.Response, err error) {
	if !strings.HasPrefix(strings.ToLower(surl), "http") {
			return nil, logger.Error("url failed")
	}

	paramValues := url.Values{}
	for k, v := range params {
		paramValues.Set(k, v)
	}
	body := bytes.NewBufferString(paramValues.Encode())

	// Create request
	req, err := http.NewRequest("POST", surl, body)
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	return HTTPDo(logger, req, options...)
}

// HTTPPostJSON 发送HTTP Post JSON请求
func HTTPPostJSON(logger ILog, surl string, object interface{}, options ...Opt) (resp *http.Response, err error) {
	if !strings.HasPrefix(strings.ToLower(surl), "http") {
			return nil, logger.Error("url failed")
	}

	body, err := jsoniter.Marshal(object)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest("POST", surl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	return HTTPDo(logger, req, options...)
}

//ES加密请求
// HTTPPostJSON 发送HTTP Post JSON请求
func HTTPPostJSONESEncrypt(logger ILog, surl string, object interface{}, user, password string, options ...Opt) (resp *http.Response, err error) {
	if !strings.HasPrefix(strings.ToLower(surl), "http") {
		return nil, logger.Error("url failed")
	}

	body, err := jsoniter.Marshal(object)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest("POST", surl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	Base64Encrypt := base64.StdEncoding.EncodeToString([]byte(user + ":" + password))
	req.Header.Add("Authorization", "Basic ZGF0YW1vcmU6QWNUM1NEMVRiY1NJ")
	logger.Info(Base64Encrypt)
	logger.Info(req.Header)
	return HTTPDo(logger, req, options...)
}
