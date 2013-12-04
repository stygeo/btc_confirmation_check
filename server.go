package main

import (
  "net"
  "runtime"
  "fmt"
  "net/http"
  "io/ioutil"
  "strings"
  "bytes"
)

/**
 * JSON RPC Wrapper method
 * Takes an required method argument and optional **args** argumens as parameters for the json rpc call
 */
func jsonRpcPost(method string, args ...string) ([]byte, error) {
  // Create a buffer which will be used to create the JSON array.
  // I'm sure there's a marshalling function using maps but I
  // decided to take a string concat approach.
  var buffer bytes.Buffer

  buffer.WriteString("[")

  // Loop over the optional args and 'fill' the array.
  for i, value := range args {
    buffer.WriteString(fmt.Sprintf("\"%s\"", value))

    // Add a comma if it isn't the last argument.
    if i != len(args) - 1 {
      buffer.WriteString(",")
    }
  }

  buffer.WriteString("]")

  data := fmt.Sprintf("{\"method\":\"%s\", \"params\":%s}", method, buffer.String())

  // Create a new req object which will be used to build the POST request.
  req, err := http.NewRequest("POST", "http://localhost:8332", strings.NewReader(data))
  if err!=nil { return nil, err }
  req.Header.Set("Content-Type", "application/json")
  // Bitcoind uses basic auth.
  req.SetBasicAuth("bit", "kojn")

  client := &http.Client{}
  resp, err := client.Do(req)

  if err!=nil { return nil, err }

  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)

  return body, nil
}

func handleClient(conn net.Conn) {
  defer conn.Close()

  // Create a buffer to store the bitcoin address (I know BTC addresses are only 34 characters)
  var buf [255]byte
  for {
    // Read the btc address
    n, err := conn.Read(buf[0:])
    if err!=nil {return}

    // Get the amount received by the specified address
    body, err := jsonRpcPost("getreceivedbyaddress", string(buf[:n]))

    if err!=nil {
      fmt.Println("JSON RPC Error", err)
      return
    }

    // Write back data to the client
    data := fmt.Sprintf("%s\r\n\r\n", body)
    _, err = conn.Write([]byte(data))
    if err!=nil {
      fmt.Println("Error writing to connection")
      return
    }
  }
}

func main() {
  runtime.GOMAXPROCS(4)

  tcpAddr, _ := net.ResolveTCPAddr("tcp4", ":1201")
  listener, _ := net.ListenTCP("tcp", tcpAddr)

  fmt.Println("Kojn Address Checker accepting connections..")
  for {
    conn, err := listener.Accept()
    if err!=nil {
      fmt.Println(err)
      continue
    }

    go handleClient(conn)
  }
}
