package main

type ContentsAPIResponse struct {
	Contents   []any `json:"contents"`
	TotalCount int   `json:"totalCount"`
	Offset     int   `json:"offset"`
	Limit      int   `json:"limit"`
}

type ManagementAPIMediaResponse struct {
	Media      []Media
	TotalCount int
	Limit      int
	Offset     int
}

type Media struct {
	Id     string
	Url    string
	Width  int
	Height int
}

type Config struct {
	Target    string   `json:"target"`
	ServiceID string   `json:"serviceId"`
	APIKey    string   `json:"apiKey"`
	Endpoints []string `json:"endpoints"`
}
