package sesha3

import (
	"github.com/moul/gotty-client"
	"github.com/pkg/errors"
)

type SeshaClientInput struct {
	URL string
}

type sesha3Client struct {
	client *gottyclient.Client
}

func NewClient(in *SeshaClientInput) (*sesha3Client, error) {
	var err error
	if len(in.URL) < 1 {
		err = errors.Wrap(err, "url should not be empty")
		return nil, err
	}

	client, err := gottyclient.NewClient(in.URL)
	if err != nil {
		err = errors.Wrap(err, "sesha3 client creation failed")
		return nil, err
	}

	return &sesha3Client{client: client}, err
}

func (c *sesha3Client) Run() error {
	err := c.client.Loop()
	if err != nil {
		err = errors.Wrap(err, "sesha3 run failed")
	}

	return err
}
