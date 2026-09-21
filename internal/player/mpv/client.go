package mpv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Client struct {
	conn net.Conn
}

type response struct {
	Data  string `json:"data"`
	Error string `json:"error"`
}

type floatResponse struct {
	Data  float64 `json:"data"`
	Error string  `json:"error"`
}

type boolResponse struct {
	Data  bool   `json:"data"`
	Error string `json:"error"`
}

func Connect(socket string) (*Client, error) {

	conn, err := net.Dial("unix", socket)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil

}

func (c *Client) Close() error {

	return c.conn.Close()

}

func (c *Client) command(property string) ([]byte, error) {

	command := fmt.Sprintf(
		`{ "command": ["get_property", "%s"] }`+"\n",
		property,
	)

	_, err := c.conn.Write([]byte(command))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(c.conn)

	response, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	return response, nil

}

func (c *Client) loadFile(path string) error {

	command := struct {
		Command []string `json:"command"`
	}{
		Command: []string{
			"loadfile",
			path,
			"replace",
		},
	}

	data, err := json.Marshal(command)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	_, err = c.conn.Write(data)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(c.conn)

	response, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	var result struct {
		Error string `json:"error"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return err
	}

	if result.Error != "success" {
		return fmt.Errorf(
			"mpv loadfile failed: %s",
			result.Error,
		)
	}

	return nil

}

func (c *Client) stringProperty(
	property string,
) (string, error) {
	data, err := c.command(property)
	if err != nil {
		return "", err
	}

	var response response

	err = json.Unmarshal(data, &response)
	if err != nil {
		return "", err
	}

	return response.Data, nil
}

func (c *Client) floatProperty(
	property string,
) (float64, error) {

	data, err := c.command(property)
	if err != nil {
		return 0, err
	}

	var response floatResponse

	err = json.Unmarshal(data, &response)
	if err != nil {
		return 0, err
	}

	return response.Data, nil

}

func (c *Client) boolProperty(
	property string,
) (bool, error) {

	data, err := c.command(property)
	if err != nil {
		return false, err
	}

	var response boolResponse

	err = json.Unmarshal(data, &response)
	if err != nil {
		return false, err
	}

	return response.Data, nil

}

func (c *Client) DemuxerViaNetwork() (bool, error) {

	return c.boolProperty(
		"demuxer-via-network",
	)

}

func (c *Client) CoreIdle() (bool, error) {

	return c.boolProperty(
		"core-idle",
	)

}

func (c *Client) Load(path string) error {

	return c.loadFile(path)

}

func (c *Client) MediaTitle() (string, error) {
	return c.stringProperty("media-title")
}

func (c *Client) Artist() (string, error) {
	return c.stringProperty("metadata/by-key/artist")
}

func (c *Client) Title() (string, error) {
	return c.stringProperty("metadata/by-key/title")
}

func (c *Client) Album() (string, error) {
	return c.stringProperty("metadata/by-key/album")
}

func (c *Client) Filename() (string, error) {
	return c.stringProperty("filename")
}

func (c *Client) Path() (string, error) {
	return c.stringProperty("path")
}

func (c *Client) PlaybackTime() (time.Duration, error) {

	value, err := c.floatProperty("playback-time")
	if err != nil {
		return 0, err
	}

	return time.Duration(value * float64(time.Second)), nil

}

func (c *Client) Duration() (time.Duration, error) {
	value, err := c.floatProperty("duration")
	if err != nil {
		return 0, err
	}

	return time.Duration(value * float64(time.Second)), nil
}
