package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/local/go-zhihu-cli/internal/config"
)

type Client struct {
	httpClient *http.Client
	endpoints  config.Endpoints
	headers    map[string]string
	cookies    map[string]string
}

func New(endpoints config.Endpoints, cookies map[string]string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: 15 * time.Second,
	}

	c := &Client{
		httpClient: httpClient,
		endpoints:  endpoints,
		headers:    config.BrowserHeaders(endpoints.BaseURL),
		cookies:    cookies,
	}
	if len(cookies) > 0 {
		if err := applyCookies(endpoints.BaseURL, jar, cookies); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Client) HTTP() *http.Client {
	return c.httpClient
}

func (c *Client) Endpoints() config.Endpoints {
	return c.endpoints
}

func (c *Client) GetJSON(rawURL string, params url.Values, out any) error {
	if len(params) > 0 {
		u, err := url.Parse(rawURL)
		if err != nil {
			return err
		}
		q := u.Query()
		for key, values := range params {
			for _, value := range values {
				q.Add(key, value)
			}
		}
		u.RawQuery = q.Encode()
		rawURL = u.String()
	}

	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	c.applyHeaders(req)
	return c.doJSON(req, out)
}

func (c *Client) PostJSON(rawURL string, payload any, out any) error {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(http.MethodPost, rawURL, body)
	if err != nil {
		return err
	}
	c.applyHeaders(req)
	req.Header.Set("Content-Type", "application/json")
	return c.doJSON(req, out)
}

func (c *Client) applyHeaders(req *http.Request) {
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	if xsrf := c.cookies[config.RequiredCookieXSRF]; xsrf != "" {
		req.Header.Set("x-xsrftoken", xsrf)
	}
}

func (c *Client) doJSON(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("未登录或登录状态已过期")
	}
	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("知乎拒绝访问: %s", string(body))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("请求失败，状态码 %d: %s", resp.StatusCode, string(body))
	}
	if out == nil || len(body) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("解析 JSON 失败: %w", err)
	}
	return nil
}

func (c *Client) CurrentCookies() map[string]string {
	out := map[string]string{}
	base, err := url.Parse(c.endpoints.BaseURL)
	if err == nil {
		for _, cookie := range c.httpClient.Jar.Cookies(base) {
			out[cookie.Name] = cookie.Value
		}
	}
	for k, v := range c.cookies {
		if _, ok := out[k]; !ok {
			out[k] = v
		}
	}
	return out
}

func applyCookies(reqURL string, jar http.CookieJar, cookies map[string]string) error {
	u, err := url.Parse(reqURL)
	if err != nil {
		return err
	}
	var list []*http.Cookie
	for name, value := range cookies {
		list = append(list, &http.Cookie{Name: name, Value: value, Path: "/"})
	}
	jar.SetCookies(u, list)
	return nil
}
