package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/local/go-zhihu-cli/internal/config"
)

type cookieStore struct {
	Cookies map[string]string `json:"cookies"`
}

func LoadCookieMap() (map[string]string, error) {
	path, err := config.CookieFile()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var store cookieStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	if !HasRequiredCookies(store.Cookies) {
		return nil, errors.New("已保存的 Cookie 不完整")
	}
	return store.Cookies, nil
}

func SaveCookieMap(cookies map[string]string) error {
	if !HasRequiredCookies(cookies) {
		return errors.New("Cookie 缺少 z_c0、_xsrf 或 d_c0")
	}
	dir, err := config.ConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cookieStore{Cookies: cookies}, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(dir, config.CookieFileName)
	return os.WriteFile(path, data, 0o600)
}

func ParseCookieString(cookieHeader string) (map[string]string, error) {
	cookies := map[string]string{}
	for _, part := range strings.Split(cookieHeader, ";") {
		part = strings.TrimSpace(part)
		if part == "" || !strings.Contains(part, "=") {
			continue
		}
		name, value, _ := strings.Cut(part, "=")
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		if name != "" {
			cookies[name] = value
		}
	}
	if !HasRequiredCookies(cookies) {
		return nil, errors.New("Cookie 必须包含 z_c0、_xsrf 和 d_c0")
	}
	return cookies, nil
}

func ClearCookies() error {
	path, err := config.CookieFile()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func HasRequiredCookies(cookies map[string]string) bool {
	if cookies == nil {
		return false
	}
	return cookies[config.RequiredCookieZC0] != "" &&
		cookies[config.RequiredCookieXSRF] != "" &&
		cookies[config.RequiredCookieDC0] != ""
}

func CookieHeader(cookies map[string]string) string {
	parts := make([]string, 0, len(cookies))
	for name, value := range cookies {
		parts = append(parts, name+"="+value)
	}
	return strings.Join(parts, "; ")
}
