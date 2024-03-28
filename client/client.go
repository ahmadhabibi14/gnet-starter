package main

import (
	"fmt"
	"net"
)

func main() {
  for i := 0; i < 10; i++ {
    sendData(i)
  }
}

func sendData(idx int) {
  // Connect to the server
  conn, err := net.Dial("tcp", "localhost:9000")
  if err != nil {
    fmt.Println(err)
    return
  }
  // Close the connection
  defer conn.Close()

  // Send some data to the server
  _, err = conn.Write([]byte(
    fmt.Sprintf("%d. Hello, server!", idx+1),
  ))
  if err != nil {
    fmt.Println(err)
    return
  }
}