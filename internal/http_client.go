package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/jefferyjob/go-easy-utils/anyUtil"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type HttpType string

var (
	HttpJson           HttpType = "application/json"
	HttpFormUrlencoded HttpType = "application/x-www-form-urlencoded"
	HttpSsml           HttpType = "application/ssml+xml"
)

type HTTPClient struct {
	ctx              context.Context
	client           *http.Client
	contentType      HttpType
	header           map[string]any
	requestLogSwitch bool
	body             any
	httpReq          *http.Request
	httpRep          *http.Response
}

const timeout = 5 * time.Second

type Option func(*HTTPClient)

func NewHTTPClient(ctx context.Context, opts ...Option) *HTTPClient {
	h := &HTTPClient{
		ctx: ctx,
		client: &http.Client{
			Timeout: timeout,
		},
		//requestLogSwitch: true,
	}

	for _, o := range opts {
		o(h)
	}

	return h
}

func WithTimeout(timeout time.Duration) Option {
	return func(h *HTTPClient) {
		h.client.Timeout = timeout
	}
}

func WithContentType(conType HttpType) Option {
	return func(h *HTTPClient) {
		h.contentType = conType
	}
}

func WithHeader(header map[string]any) Option {
	return func(h *HTTPClient) {
		h.header = header
	}
}

func WithRequestLogSwitch(switchLog bool) Option {
	return func(h *HTTPClient) {
		h.requestLogSwitch = switchLog
	}
}

func (hc *HTTPClient) SendRequest(method, url string, body map[string]any) (*http.Response, func(), error) {
	hc.body = body
	defer hc.requestLogPrint()
	var reqBody []byte

	switch method {
	case http.MethodGet:
		url = hc.setGet(url, body)
		break
	case http.MethodPost:
		reqBody = hc.setPost(body)
		break
	case http.MethodDelete:
		reqBody = hc.setPost(body)
		break
	default:
		return nil, func() {}, errors.New("not define method: " + method)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	hc.httpReq = req
	if err != nil {
		return nil, func() {}, err
	}

	hc.setHeader(req)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(reqBody)))

	resp, err := hc.client.Do(req)
	hc.httpRep = resp
	if err != nil {
		return nil, func() {}, err
	}

	//if resp.StatusCode != http.StatusOK {
	//	return nil, func() { resp.Body.Close() }, errors.New("http response code failed, code: " + resp.Status)
	//}

	return resp, func() { resp.Body.Close() }, nil

	//respBody, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return nil, err
	//}
	//return respBody, nil
}

// 设置 Http Header
func (hc *HTTPClient) setHeader(req *http.Request) {
	// 自定义 Header
	if len(hc.header) > 0 {
		for key, value := range hc.header {
			req.Header.Set(key, anyUtil.AnyToStr(value))
		}
	}

	// 设置 Header Content-Type
	if hc.contentType != "" {
		req.Header.Set("Content-Type", string(hc.contentType))
	}
}

func (hc *HTTPClient) setPost(body map[string]any) []byte {
	// json 方式请求
	if hc.contentType == HttpJson {
		jsonData, ok := body["json"]
		if !ok {
			log.Println("找不到对应的json的key")
			return nil
		}
		return []byte(anyUtil.AnyToStr(jsonData))
	}

	// ssml+xml 方式请求
	if hc.contentType == HttpSsml {
		xmlData, ok := body["xml"]
		if !ok {
			return nil
		}
		xmlBytes, _ := xml.MarshalIndent(xmlData, "", "    ")
		return xmlBytes
	}

	// 其他方式请求
	values := url.Values{}
	for k, v := range body {
		values.Set(k, fmt.Sprintf("%v", v))
	}

	params := values.Encode()
	return []byte(params)
}

func (hc *HTTPClient) setGet(targetUrl string, body map[string]any) string {
	values := url.Values{}
	for k, v := range body {
		values.Set(k, fmt.Sprintf("%v", v))
	}
	params := values.Encode()
	if params == "" {
		return targetUrl
	}
	return targetUrl + "?" + values.Encode()
}

func (hc *HTTPClient) requestLogPrint() {
	if !hc.requestLogSwitch {
		return
	}

	requestLog := struct {
		Url      string      `json:"url"`
		Method   string      `json:"method"`
		Header   http.Header `json:"header"`
		Body     any         `json:"body"`
		Req      string      `json:"req"`
		RespCode int         `json:"resp_code"`
		Resp     string      `json:"resp"`
	}{
		Url:      hc.httpReq.URL.String(),
		Method:   hc.httpReq.Method,
		Header:   hc.httpReq.Header,
		Body:     hc.body,
		RespCode: hc.httpRep.StatusCode,
	}

	if hc.httpReq.Body != nil {
		bodyBytes, err := io.ReadAll(hc.httpReq.Body)
		if err == nil {
			requestLog.Req = string(bodyBytes)
		}
	}

	if hc.httpReq.Body != nil {
		bodyBytes, err := io.ReadAll(io.LimitReader(hc.httpReq.Body, 4096)) // 限制响应正文内容以避免过度日志记录
		if err == nil {
			requestLog.Resp = string(bodyBytes)
		}
	}

	jsonData, err := json.Marshal(&requestLog)
	if err != nil {
		log.Printf("日志转json报错, err:%v \n", err)
	}
	log.Printf("http client request log: %v", string(jsonData))
}
