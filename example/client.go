package example

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

var mu sync.Mutex

type Client struct {
	Boxes []Box `json:Boxes`
}

type Box struct {
	Id string `json:Id`
	Bundle string `json:Bundle`
}

func NewClient() (*Client, error) {
	f, err := os.Create("/tmp/client_remote_data.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := f.WriteString("{}"); err != nil {
		return nil, err
	}

	return &Client{}, nil
}

func (c *Client) CreateBox(t *Box) error {
	mu.Lock()
	defer mu.Unlock()

	content, err := ioutil.ReadFile("/tmp/client_remote_data.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &c); err != nil {
		return err
	}

	c.Boxes = append(c.Boxes, *t)

	content, err = json.Marshal(c)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("/tmp/client_remote_data.json", content, 0); err != nil {
		return err
	}

	return nil
}

func (c *Client) ReadBox(id string) (*Box, error) {
	mu.Lock()
	defer mu.Unlock()

	content, err := ioutil.ReadFile("/tmp/client_remote_data.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(content, &c); err != nil {
		return nil, err
	}

	for _, b := range c.Boxes {
		if b.Id == id {
			return &b, nil
		}
	}

	return nil, fmt.Errorf("Box (%s) not found", id)
}

func (c *Client) UpdateBox(t *Box) error {
	mu.Lock()
	defer mu.Unlock()

	content, err := ioutil.ReadFile("/tmp/client_remote_data.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &c); err != nil {
		return err
	}

	var tmp []Box
	for _, b := range c.Boxes {
		if b.Id != t.Id {
			tmp = append(tmp, b)
		}
	}
	tmp = append(tmp, *t)
	c.Boxes = tmp

	content, err = json.Marshal(c)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("/tmp/client_remote_data.json", content, 0); err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteBox(id string) error {
	mu.Lock()
	defer mu.Unlock()

	content, err := ioutil.ReadFile("/tmp/client_remote_data.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &c); err != nil {
		return err
	}

	var tmp []Box
	for _, b := range c.Boxes {
		if b.Id != id {
			tmp = append(tmp, b)
		}
	}
	c.Boxes = tmp

	content, err = json.Marshal(c)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("/tmp/client_remote_data.json", content, 0); err != nil {
		return err
	}

	return nil
}
