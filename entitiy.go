package main

type ContentsAPIResponse struct {
	Contents   []any `json:"contents"`
	TotalCount int   `json:"totalCount"`
	Offset     int   `json:"offset"`
	Limit      int   `json:"limit"`
}

type ManagementAPIMediaResponse struct {
	Media      []Media `json:"media"`
	TotalCount int     `json:"totalCount"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
}

type Media struct {
	Id     string `json:"id"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Config struct {
	Target    string   `json:"target"`
	ServiceID string   `json:"serviceId"`
	APIKey    string   `json:"apiKey"`
	Endpoints []string `json:"endpoints"`
}
