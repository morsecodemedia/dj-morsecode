package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {

	conn, err := net.Dial(
		"unix",
		"/tmp/dj-morsecode.sock",
	)

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	fmt.Println("Connected!")

	command := `{ "command": ["get_property", "media-title"] }` + "\n"

	_, err = conn.Write([]byte(command))
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(conn)

	response, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	fmt.Println(response)

}
