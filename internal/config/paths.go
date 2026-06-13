package config

import (
	"os"
	"path/filepath"
)

const (
	AppDirName         = ".go-zhihu-cli"
	CookieFileName     = "cookies.json"
	RequiredCookieZC0  = "z_c0"
	RequiredCookieXSRF = "_xsrf"
	RequiredCookieDC0  = "d_c0"
)

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, AppDirName), nil
}

func CookieFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, CookieFileName), nil
}
