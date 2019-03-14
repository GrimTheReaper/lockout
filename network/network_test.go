package network

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/grimthereaper/lockout/pb"
	"google.golang.org/grpc"
)

var flagAPIPort int
var flagAPIHost string

var flagGRPCPort int
var flagGRPCHost string

func init() {
	// Using "flag" to slim down the external imports,
	//  for this test, the less imports we have the easier it will be to run..
	flag.IntVar(&flagAPIPort, "api-port", 8080, "Port to bind the network to")
	flag.StringVar(&flagAPIHost, "api-host", "", "Host to bind the network to")

	flag.IntVar(&flagGRPCPort, "grpc-port", 8082, "Port to bind the network to")
	flag.StringVar(&flagGRPCHost, "grpc-host", "", "Host to bind the network to")

	flag.Parse()
}

func TestGRPC(t *testing.T) {
	go func() {
		_, err := Serve(flagGRPCHost, flagGRPCPort, true)
		if err != nil {
			panic(err)
		}
	}()

	// Give it a few seconds to initialize.
	time.Sleep(100 * time.Millisecond)

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", flagGRPCHost, flagGRPCPort), grpc.WithInsecure())
	ok(t, err)

	client := pb.NewWhitelistCheckerClient(conn)

	response, err := client.CheckIP(context.TODO(), &pb.IPCheckRequest{
		Ip:        "8.8.8.8",
		Countries: []string{"US"},
	})
	ok(t, err)
	equals(t, true, response.GetWhitelisted())

	ok(t, conn.Close())
}

func TestAPI(t *testing.T) {
	go func() {
		_, err := Serve(flagAPIHost, flagAPIPort, false)
		if err != nil {
			panic(err)
		}
	}()

	// Give it a few seconds to initialize.
	time.Sleep(100 * time.Millisecond)

	byts, err := json.Marshal(pb.IPCheckRequest{
		Ip:        "8.8.8.8",
		Countries: []string{"US"},
	})
	ok(t, err)
	equals(t, true, len(byts) != 0)

	// TODO: Figure out why "http://%v" is invalid.
	request, err := http.NewRequest("POST", "http://"+fmt.Sprintf("%v:%v/api/v0/ip/whitelist", flagAPIHost, flagAPIPort), bytes.NewReader(byts))
	ok(t, err)

	http.DefaultClient.Timeout = 4 * time.Second
	response, err := http.DefaultClient.Do(request)
	ok(t, err)
	defer response.Body.Close()

	// The expected result is small, but its always recommended to use a Decoder for http responses.
	decoder := json.NewDecoder(response.Body)

	var resp pb.IPCheckResponse
	ok(t, decoder.Decode(&resp))

	equals(t, true, resp.GetWhitelisted())
}
