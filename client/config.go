package client

import (
	"encoding/json"
	"os"
)

func (c *Client) LoadConfig() error {
	// デフォルト値を設定
	c.Config.RequestUnit = 10

	f, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	d.DisallowUnknownFields()
	d.Decode(c.Config)
	return nil
}
