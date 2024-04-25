package network

/**
 * @Author: lee
 * @Description:
 * @File: rest
 * @Date: 2021/9/7 3:48 下午
 */

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/juju/ratelimit"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type RestAgent struct {
	NetAgentBase
	Client   *resty.Client
	proxyUrl string
}

var _ HttpInterface = (*RestAgent)(nil)

type RestOptionFunc func(*RestAgent)

func NewRestClient(host string, port uint, isHttps bool, options ...RestOptionFunc) (*RestAgent, error) {
	trimHost := strings.TrimLeft(host, " ")

	if strings.HasPrefix(trimHost, "http") && strings.Contains(trimHost, "://") {
		trimHost = trimHost
	} else {
		if isHttps {
			trimHost = "https://" + trimHost
		} else {
			trimHost = "http://" + trimHost
		}
	}

	if 0 != port {
		trimHost += ":" + strconv.FormatUint((uint64)(port), 10)
	}

	hostUrl, err := url.Parse(trimHost)
	if nil != err {
		return nil, err
	}
	ret := &RestAgent{
		NetAgentBase: NetAgentBase{
			URL: hostUrl,
		},
	}

	for _, option := range options {
		if nil != option {
			option(ret)
		}
	}
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	hc := http.Client{
		Jar:     cookieJar,
		Timeout: 20 * time.Second,
	}

	if "" != ret.proxyUrl {
		proxyUrl, err := url.Parse(ret.proxyUrl)
		if nil == err {
			hc.Transport = &http.Transport{
				// 设置代理
				Proxy: http.ProxyURL(proxyUrl),
			}
		}

	}
	client := resty.NewWithClient(&hc)
	ret.Client = client

	return ret, nil
}

func WithRestProxy(url string) RestOptionFunc {
	return func(agent *RestAgent) {
		agent.proxyUrl = url
	}
}

type RestAgentBase struct {
	Bucket *ratelimit.Bucket
	Client *RestAgent
}

func (agent *RestAgentBase) Get(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) ([]byte, error) {
	agent.Bucket.Wait(1)
	ret, err := agent.Client.Get(path, params, headers, cookies)
	return []byte(ret), err
}

func (agent *RestAgentBase) Post(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) ([]byte, error) {
	agent.Bucket.Wait(1)
	ret, err := agent.Client.Post(path, reqBody, params, headers, cookies)
	return []byte(ret), err
}

func (agent *RestAgentBase) PostForm(path string, reqBody map[string]string, params map[string]string, headers map[string]string, cookies []*http.Cookie) ([]byte, error) {
	agent.Bucket.Wait(1)
	ret, err := agent.Client.PostForm(path, reqBody, params, headers, cookies)
	return []byte(ret), err
}

func (agent *RestAgentBase) Put(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) ([]byte, error) {
	agent.Bucket.Wait(1)
	ret, err := agent.Client.Put(path, reqBody, params, headers, cookies)
	return []byte(ret), err
}

func (agent *RestAgentBase) PutForm(path string, reqBody map[string]string, params map[string]string, headers map[string]string, cookies []*http.Cookie) ([]byte, error) {
	agent.Bucket.Wait(1)
	ret, err := agent.Client.PutForm(path, reqBody, params, headers, cookies)
	return []byte(ret), err
}

func (agent *RestAgentBase) Delete(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) ([]byte, error) {
	agent.Bucket.Wait(1)
	ret, err := agent.Client.Delete(path, reqBody, params, headers, cookies)
	return []byte(ret), err
}

func (h *RestAgent) SimpleGet(path string, params map[string]string) (string, error) {
	url := h.URL.String() + path
	if nil != params {

	}
	res, err := h.Client.R().SetQueryParams(params).Get(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != 200 {
		return string(res.Body()), fmt.Errorf("response err: %s", res.String())
	}

	return res.String(), nil
}

func (h *RestAgent) SimplePost(path string, reqBody string, params map[string]string) (string, error) {
	url := h.URL.String() + path
	r := h.Client.R()

	res, err := r.SetQueryParams(params).SetBody(reqBody).SetHeader("Content-Type", "application/json").Post(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}

	return res.String(), nil
}

func (h *RestAgent) Get(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	r := h.Client.R()
	res, err := r.SetQueryParams(params).SetHeaders(headers).SetCookies(cookies).Get(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}

	return res.String(), nil
}

func (h *RestAgent) Post(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	r := h.Client.R()
	res, err := r.SetQueryParams(params).SetBody(reqBody).SetHeaders(headers).SetCookies(cookies).Post(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}
	return res.String(), nil

}

func (h *RestAgent) PostForm(path string, reqBody map[string]string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	r := h.Client.R()
	res, err := r.SetFormData(reqBody).SetQueryParams(params).SetHeaders(headers).SetCookies(cookies).Post(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}
	return res.String(), nil
}

func (h *RestAgent) Put(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	r := h.Client.R()
	res, err := r.SetQueryParams(params).SetBody(reqBody).SetHeaders(headers).SetCookies(cookies).Put(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}
	return res.String(), nil

}

func (h *RestAgent) PutForm(path string, reqBody map[string]string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	r := h.Client.R()
	res, err := r.SetQueryParams(params).SetFormData(reqBody).SetHeaders(headers).SetCookies(cookies).Put(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}
	return res.String(), nil
}

func (h *RestAgent) Delete(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	r := h.Client.R()
	res, err := r.SetQueryParams(params).SetBody(reqBody).SetHeaders(headers).SetCookies(cookies).Delete(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}
	return res.String(), nil

}
