package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

var serverAddr = flag.String("server", "localhost:8080", "Server address")

type Client struct {
	connection net.Conn
	client     *jrpc2.Client
	id         int
	ctx        context.Context
	closeCtx   context.CancelFunc
}

func (c Client) close() {
	c.client.Close()
	c.closeCtx()
}

func newClient(id int) Client {
	conn, err := net.Dial(jrpc2.Network(*serverAddr))
	if err != nil {
		log.Fatalf("Dial %q: %v", *serverAddr, err)
	}
	log.Printf("Connected to %v as ID:%d", conn.RemoteAddr(), id)
	cli := jrpc2.NewClient(channel.Line(conn, conn), &jrpc2.ClientOptions{
		OnNotify: func(req *jrpc2.Request) {
			var params json.RawMessage
			req.UnmarshalParams(&params)
			log.Printf("[server push] Method %q params %#q", req.Method(), string(params))
		},
	})
	ctx, close := context.WithCancel(context.Background())
	return Client{
		connection: conn,
		client:     cli,
		id:         id,
		ctx:        ctx,
		closeCtx:   close,
	}
}

func countString(ctx context.Context, cli *jrpc2.Client, msg []string) (result int, err error) {
	err = cli.CallResult(ctx, "StringOperations.CountString", msg, &result)
	return
}

func createRandomString() string {
	n := rand.Intn(100)
	b := make([]byte, n)
	return base64.URLEncoding.EncodeToString(b)[:n]
}

func concurrencyStressTest(wg *sync.WaitGroup) {
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			cli := newClient(i)
			randomInput := createRandomString()
			clientSendMessage(cli.ctx, cli.client, i, wg, randomInput)

		}()
		wg.Wait()
		log.Println("Concurrencytest finished")
	}
}

func clientSendMessage(ctx context.Context, cli *jrpc2.Client, id int, wg *sync.WaitGroup, input string) {
	input = strings.TrimSpace(input)
	log.Print("\n-- Sending some individual requests...")
	if strLen, err := countString(ctx, cli, []string{input}); err != nil {
		log.Fatalln("StringOperations:", err)

	} else {
		log.Printf("StringCount result=%d as ID:%d", strLen, id)
	}
	wg.Done()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	wg := sync.WaitGroup{}
	id := 9999
	flag.Parse()

	if *serverAddr == "" {
		log.Fatal("You must provide -server address to connect to")
	}

	consoleClient := newClient(id)
	defer consoleClient.close()

	// Checking input
	for {
		fmt.Print("Enter a string to send to the server (or 'exit' to quit): ")
		if !scanner.Scan() {
			log.Fatalf("Failed to read input: %v", scanner.Err())
		}
		input := scanner.Text()
		switch input {
		case "-exit":
			fmt.Println("Client exited.")
			return
		case "-concurrency":
			concurrencyStressTest(&wg)
		default:
			wg.Add(1)
			clientSendMessage(consoleClient.ctx, consoleClient.client, id, &wg, input)
			wg.Wait()
		}
	}
}
