package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/local/go-zhihu-cli/configs"
)

type Endpoints struct {
	BaseURL                   string `json:"base_url"`
	APIV3                     string `json:"api_v3"`
	APIV4                     string `json:"api_v4"`
	MeURL                     string `json:"me_url"`
	FeedRecommendURL          string `json:"feed_recommend_url"`
	AnswerURLTemplate         string `json:"answer_url_template"`
	AnswerCommentsURLTemplate string `json:"answer_comments_url_template"`
}

func DefaultEndpoints() Endpoints {
	endpoints, err := defaultEndpoints()
	if err != nil {
		panic("加载内置接口配置失败: " + err.Error())
	}
	return endpoints
}

func LoadEndpoints(path string) (Endpoints, error) {
	endpoints, err := defaultEndpoints()
	if err != nil {
		return Endpoints{}, err
	}

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return Endpoints{}, err
		}
		if err := json.Unmarshal(data, &endpoints); err != nil {
			return Endpoints{}, err
		}
	}

	if endpoints.BaseURL == "" || endpoints.FeedRecommendURL == "" {
		return Endpoints{}, errors.New("接口配置缺少必要 URL")
	}
	return endpoints, nil
}

func defaultEndpoints() (Endpoints, error) {
	var endpoints Endpoints
	if err := json.Unmarshal(configs.DefaultEndpointsJSON, &endpoints); err != nil {
		return Endpoints{}, err
	}
	return endpoints, nil
}
