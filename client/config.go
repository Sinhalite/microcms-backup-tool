package client

import (
	"encoding/json"
	"os"
)

func (c *Client) LoadConfig(configPath string) error {
	// デフォルト値を設定
	c.Config.Contents.RequestUnit = 10

	f, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	d.DisallowUnknownFields()
	if err := d.Decode(c.Config); err != nil {
		return err
	}
	return nil
}
