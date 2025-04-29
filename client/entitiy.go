package client

import (
	"encoding/json"
)

type ContentsAPIResponse struct {
	Contents   json.RawMessage `json:"contents"`
	TotalCount int             `json:"totalCount"`
	Offset     int             `json:"offset"`
	Limit      int             `json:"limit"`
}

type ManagementAPIMediaResponse struct {
	Media      []Media `json:"media"`
	TotalCount int     `json:"totalCount"`
	Token      string  `json:"token"`
}

type Media struct {
	Id     string `json:"id"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// ContentsConfig はコンテンツバックアップの設定を保持する構造体
type ContentsConfig struct {
	// 公開コンテンツを取得するためのAPIキー（classifyByStatusがfalseの場合はこれのみ必要）
	GetPublishContentsAPIKey string `json:"getPublishContentsAPIKey"`
	// 以下のフィールドはclassifyByStatusがtrueの場合のみ必要
	GetAllStatusContentsAPIKey string   `json:"getAllStatusContentsAPIKey,omitempty"`
	GetContentsMetaDataAPIKey  string   `json:"getContentsMetaDataAPIKey,omitempty"`
	Endpoints                  []string `json:"endpoints"`
	RequestUnit                int      `json:"requestUnit"`
	ClassifyByStatus           bool     `json:"classifyByStatus"`
}

// MediaConfig はメディアバックアップの設定を保持する構造体
type MediaConfig struct {
	APIKey      string `json:"apiKey"`
	RequestUnit int    `json:"requestUnit"`
}

type Config struct {
	Target    string         `json:"target"`
	ServiceID string         `json:"serviceId"`
	Contents  ContentsConfig `json:"contents,omitempty"`
	Media     MediaConfig    `json:"media,omitempty"`
}

type Client struct {
	Config *Config
}
