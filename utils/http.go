package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"log"
	"net/url"
)

type RequestQueries map[string]string

type FastRequestParams struct {
	Url         string
	ContentType string
	Method      string
	Body        []byte
}

func FastDo(in FastRequestParams, out interface{}) error {
	var (
		err      error
		httpReq  *fasthttp.Request
		httpResp *fasthttp.Response
	)

	httpReq = fasthttp.AcquireRequest()
	httpResp = fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(httpReq)
		fasthttp.ReleaseResponse(httpResp)
	}()

	switch in.Method {
	case "GET":
	case "POST":
		if in.ContentType != "" {
			httpReq.Header.SetContentType(in.ContentType)
		}
	default:
		return errors.New("http method must 'GET' or 'POST'")
	}

	httpReq.Header.SetMethod(in.Method)
	httpReq.SetRequestURI(in.Url)
	httpReq.SetBody(in.Body)

	if err = fasthttp.Do(httpReq, httpResp); err != nil {
		return nil
	}

	log.Println("resp:", string(httpResp.Body()))
	return json.NewDecoder(bytes.NewReader(httpResp.Body())).Decode(out)
}

func EncodeURL(baseurl string, params RequestQueries) (string, error) {
	var (
		err     error
		_url    *url.URL
		_values url.Values
	)

	if _url, err = url.Parse(baseurl); err != nil {
		return "", err
	}

	_values = _url.Query()
	for k, v := range params {
		_values.Set(k, v)
	}

	_url.RawQuery = _values.Encode()

	return _url.String(), nil
}

func Get(api string, response interface{}) error {
	var (
		err  error
		resp *fasthttp.Response
		req  *fasthttp.Request
	)

	req = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("GET")
	req.SetRequestURI(api)

	resp = fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err = fasthttp.Do(req, resp); err != nil {
		return nil
	}

	log.Println("get resp:", string(resp.Body()))
	return json.NewDecoder(bytes.NewReader(resp.Body())).Decode(response)
}

func PostJSON(api string, params, response interface{}) error {
	var (
		err     error
		resp    *fasthttp.Response
		req     *fasthttp.Request
		reqBody []byte
	)

	if params != nil {
		if reqBody, err = json.Marshal(params); err != nil {
			return nil
		}
	}

	req = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")
	req.SetRequestURI(api)
	req.SetBody(reqBody)

	resp = fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err = fasthttp.Do(req, resp); err != nil {
		return err
	}

	log.Println("post resp:", string(resp.Body()))
	return json.NewDecoder(bytes.NewReader(resp.Body())).Decode(response)
}
