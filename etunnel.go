/**
 * etunnel TLS <=> Plain or Plain <=> TLS
 */

package main

import "os"
import "flag"
import "fmt"
import "strings"
import "time"
import "math/rand"
import "path/filepath"
import "sync"

type Tunnel struct {
	id string
	l0 *socket
	t1 *socket
	pump chan []byte
}


// I like these as IDs
func generateULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	u, _ := ulid.New(ulid.Timestamp(t), entropy)
	return u.String()
	// Output: 0000XSNJG0MQJHBF4QX1EFD6Y3
}

// var bridge_list
var tunnel_list = make(map[*Tunnel]int64)


/**
 * Main
 */
func main() {

	cmd, _ := os.Executable()

	dir, _ := filepath.Abs(filepath.Dir(cmd))
	fmt.Println(dir)

	err50 := os.Chdir(dir)
	if err50 != nil {
		panic(err50)
	}

	// Read Options
	listenAddr := flag.String("listen", "127.0.0.1:1236379", "Listen Address and Port")
	listenCert := flag.String("listen-tls", "/path/to/file.pem", "Certificate Chain for Listen Socket")
	tunnelAddr := flag.String("tunnel", "localhost:1234", "Tunnel Address and Port")
	tunnelCert := flag.String("tunnel-tls", "tunnel.pem", "Certificate Chain for Tunnel Socket")

	flag.Parse()

	fmt.Println("Listen: ", *listenAddr)
	if (listenCert.length) {
		fmt.Println("Listen Certificate: ", *listenAddr)
	}
	fmt.Println("Tunnel:", *tunnelAddr)
	if (tunnelCert.length) {
		fmt.Println("Tunnel Certificate: ", *listenAddr)
	}

	// Now do the thing
}
