/**
 * etunnel TLS <=> Plain or Plain <=> TLS
 */

package main

import "os"
import "io"
import "flag"
import "fmt"
import "log"
// import "strings"
// import "time"
// import "math/rand"
import "path/filepath"
// import "sync"

import "net"
import "crypto/tls"

type Tunnel struct {
	id string
	l0 *net.Listener
	t1 *net.Listener
	pump chan []byte
}


// I like these as IDs
// func generateULID() string {
// 	t := time.Now()
// 	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
// 	u, _ := ulid.New(ulid.Timestamp(t), entropy)
// 	return u.String()
// 	// Output: 0000XSNJG0MQJHBF4QX1EFD6Y3
// }

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
	listenAddr := flag.String("listen", "127.0.0.1:1234", "Listen Address and Port")
	listenCert := flag.String("listen-tls", "server.pem", "Certificate Chain for Listen Socket")
	tunnelAddr := flag.String("tunnel", "localhost:5678", "Tunnel Address and Port")
	tunnelCert := flag.String("tunnel-tls", "tunnel.pem", "Certificate Chain for Tunnel Socket")

	flag.Parse()

	fmt.Println("Listen: ", *listenAddr)
	if (len(*listenCert) > 0) {
		fmt.Println("Listen Certificate: ", *listenCert)
	}
	fmt.Println("Tunnel:", *tunnelAddr)
	if (len(*tunnelCert) > 0) {
		fmt.Println("Tunnel Certificate: ", *tunnelCert)
	}

	// Now do the thing
	// var serverCert = nil
	// var clientCert = nil

	serverCert, err := tls.LoadX509KeyPair(*listenCert, "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	clientCert, err := tls.LoadX509KeyPair(*tunnelCert, "tunnel.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: []tls.Certificate{serverCert},
	}

	server, err := tls.Listen("tcp", *listenAddr, config)
	if err != nil {
		log.Println(err)
		return
	}
	defer server.Close()

	for {
		T, err := server.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go tunnel(T, tunnelAddr, clientCert)
	}
}

/**
 * Process the Tunnel
 */
func tunnel(server net.Conn, tunnelAddr *string, clientCert tls.Certificate) {

	// defer server.Close()

	log.Println("Create Tunnel")

	client, err := net.Dial("tcp", *tunnelAddr)
	if err != nil {
		log.Println(err)
		return
	}


	// Stream from Server to Client
	go pumpServerToClient(server, client)

	// Stream from Client to Server
	go pumpClientToServer(server, client)

	// for {
		// // read in input from stdin
		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Text to send: ")
		// text, _ := reader.ReadString('\n')
		// // send to socket
		// fmt.Fprintf(conn, text + "\n")
		// // listen for reply
		// message, _ := bufio.NewReader(conn).ReadString('\n')
		// fmt.Print("Message from server: "+message)
		//

		// r := bufio.NewReader(conn)
		// for {
		// msg, err := r.ReadString('\n')
		// if err != nil {
		//     log.Println(err)
		//     return
		// }
		//
		// println(msg)

		// n, err := conn.Write([]byte("world\n"))
		// if err != nil {
		//     log.Println(n, err)
		//     return
		// }
		// ?

	// }
}

func pumpServerToClient(server net.Conn, client net.Conn) {

	// defer server.Close()
	// defer client.Close()

	for {

		fmt.Print(".")

		tmp := make([]byte, 256)

		br, err := server.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		bw, err := client.Write(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("send error:", err)
			}
			break
		}

		if (br != bw) {
			fmt.Println("Byte Count Not Matched:", br, bw)
			break;
		}

		// // read in input from stdin
		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Text to send: ")
		// text, _ := reader.ReadString('\n')
		// // send to socket
		// fmt.Fprintf(conn, text + "\n")
		// // listen for reply
		// message, _ := bufio.NewReader(conn).ReadString('\n')
		// fmt.Print("Message from server: "+message)
		//

		// r := bufio.NewReader(conn)
		// for {
		// msg, err := r.ReadString('\n')
		// if err != nil {
		//     log.Println(err)
		//     return
		// }
		//
		// println(msg)

		// n, err := client.Write([]byte("world\n"))
		// if err != nil {
		//     log.Println(n, err)
		//     return
		// }
		// ?

	}

	fmt.Println("pumpServerToClient - done")

}

func pumpClientToServer(server net.Conn, client net.Conn) {

	// defer server.Close()
	// defer client.Close()
	for {

		fmt.Print(".")

		tmp := make([]byte, 256)

		br, err := client.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("C->T read error:", err)
			}
			break
		}

		bw, err := server.Write(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("C->S send error:", err)
			}
			break
		}

		if (br != bw) {
			fmt.Println("Byte Count Not Matched:", br, bw)
			break;
		}

	}

	fmt.Println("pumpClientToServer - done")

}
