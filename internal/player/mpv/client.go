package mpv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
}

type response struct {
	Data  string `json:"data"`
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

func (c *Client) stringProperty(
	property string,
) (string, error) {
	data, err := c.command("media-title")
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

func (c *Client) MediaTitle() (string, error) {
	return c.stringProperty("media-title")
}

func (c *Client) Filename() (string, error) {
	return c.stringProperty("filename")
}

// func (c *Client) MediaTitle() (string, error) {

// 	data, err := c.command("media-title")
// 	if err != nil {
// 		return "", err
// 	}

// 	var response response

// 	err = json.Unmarshal(data, &response)
// 	if err != nil {
// 		return "", err
// 	}

// 	return response.Data, nil

// }
