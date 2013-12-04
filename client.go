package main

import (
  "net"
  "fmt"
  "time"
)

/**
 * Establish a connection with the server, write an address to the stream
 * and await for a response.
 */
func fetchData(lockChan chan bool) {
  tcpAddr, _ := net.ResolveTCPAddr("tcp4", "localhost:1201")
  conn, _ := net.DialTCP("tcp", nil, tcpAddr)

  // Write the address to the stream
  conn.Write([]byte("mfgiXnSzJF6mb37FDorWJeeqeP3tFTERpo"))
  var buf[255]byte

  // Read back the response :-)
  n, err := conn.Read(buf[0:])
  if err!=nil {
    fmt.Println("Error reading data")
    return
  }

  fmt.Printf("Received: %s\n", string(buf[:(n-5)]))

  // 'Unlock' the lock chan (thread join)
  lockChan<-true

  conn.Close()
}

func main() {
  concurrentConnections := 100

  lockChan := make(chan bool, concurrentConnections)

  start := time.Now()

  /*
   * Create 100 concurrent connections and for each connection
   * and create a locking channel, works kinda like joining threads
   * in other languages
   */
  for i:=0; i<concurrentConnections; i++{
    go fetchData(lockChan)
  }
  for i:=0; i<int(concurrentConnections); i++{
    <-lockChan
  }

  elapsed := time.Since(start).Seconds()
  fmt.Println(concurrentConnections, "concurrent connections took me", elapsed, "seconds to complete.")
}
