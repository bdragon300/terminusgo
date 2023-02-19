package tusc

import (
	"context"
	"net/http"
	"net/url"

	"github.com/bdragon300/tusgo"
)

type Client struct {
	*tusgo.Client
}

func New(client *http.Client, baseURL string) *Client {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}
	return &Client{tusgo.NewClient(client, u)}
}

func (c *Client) CreateFile(filename string, remoteSize int64, metadata map[string]string) (*tusgo.Upload, error) {
	u := tusgo.Upload{}
	meta := make(map[string]string)
	for k, v := range metadata {
		meta[k] = v
	}
	meta["filename"] = filename
	_, err := c.CreateUpload(&u, remoteSize, false, meta)

	return &u, err
}

func (c *Client) WithContext(ctx context.Context) *Client {
	res := c.Client.WithContext(ctx)
	return &Client{res}
}
