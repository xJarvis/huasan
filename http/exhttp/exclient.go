package exhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/xjarvis/huashan/log/logger"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/**
 *  HTTP-客户端
 */


var HClient *http.Client

const CONTENT_TAG 	string = "Content-Type"
const COOKIE_TAG  	string = "Cookie"

const POST_ENCODED string = "application/x-www-form-urlencoded"
const POST_FORM 	string = "multipart/form-data"
const POST_JSON 	string = "application/json"
const POST_XML 	string = "text/xml"


func init() {
	HClient = &http.Client{
		Timeout:time.Duration(30 * time.Second),
	}
	//代理使用
	//proxyUrl, _ := url.Parse("http://localhost:1087")
	//HClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
}

/**
 * HTTP POST 请求
 */
func HttpPost(urlStr string, params interface{}, headers map[string]string, cookies map[string]string) ([]byte,error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	if nil == headers {
		headers = map[string]string{CONTENT_TAG: POST_ENCODED}
	} else if _, ok := headers[CONTENT_TAG]; !ok {
		headers[CONTENT_TAG] = POST_ENCODED
	}
	urlStr = u.String()
	var reader io.Reader
	switch headers[CONTENT_TAG] {
	case POST_ENCODED:
		q := u.Query()
		for k, v := range params.(map[string]string) {
			q.Set(k, v)
		}
		reader = strings.NewReader(q.Encode())
	case POST_JSON:
		reqByte, err := json.Marshal(params)
		if err != nil {
			logger.Error("json marshal data failed!", err)
			return nil, err
		}
		reader = bytes.NewBuffer(reqByte)
	default:
		if params != nil {
			reader = bytes.NewBuffer(params.([]byte))
		}
	}
	req,err := http.NewRequest("POST", urlStr, reader)
	if err != nil {
		return nil, err
	}

	//设置请求头
	for k,v := range headers {
		req.Header.Set(k, v)
	}

	//设置cookie
	for k,v := range cookies {
		req.Header.Set(COOKIE_TAG, k + "=" + v)
	}

	req.Close = true

	//执行请求
	resp, err := HClient.Do(req)
	if err != nil {
		return nil, errors.New("post request failed!" + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/**
 * HTTP GET 请求
 */
func HttpGet(urlStr string, params map[string]string, headers map[string]string, cookies map[string]string) ([]byte,error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	urlStr = u.String()

	var reader io.Reader
	req,err := http.NewRequest("GET", urlStr, reader)
	if err != nil {
		return nil, err
	}

	//设置请求头
	for k,v := range headers {
		req.Header.Set(k, v)
	}

	//设置cookie
	for k,v := range cookies {
		req.Header.Set(COOKIE_TAG, k + "=" + v)
	}

	req.Close = true

	//执行请求
	resp, err := HClient.Do(req)
	if err != nil {
		return nil, errors.New("get request failed!" + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/**
 * 简单GET请求
 */
func SimpleGet(urlStr string, params map[string]string) ([]byte,error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil,err
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	res, err := HClient.Get(u.String())
	if err != nil {
		return nil,errors.New("get request failed!" + err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil,err
	}
	return body,nil
}


/**
 * URL GET请求
 */
func Get(urlStr string) ([]byte,error) {
	res, err := HClient.Get(urlStr)
	if err != nil {
		return nil,errors.New("get request failed!" + err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil,err
	}
	return body,nil
}