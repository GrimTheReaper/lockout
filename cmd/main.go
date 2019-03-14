package main

import (
	"flag"
	"os"

	"github.com/grimthereaper/lockout/network"
)

var flagAPIPort int
var flagAPIHost string

var flagGRPCPort int
var flagGRPCHost string

var exitChan = make(chan int)

func init() {
	// Using "flag" to slim down the external imports,
	//  for this test, the less imports we have the easier it will be to run..
	flag.IntVar(&flagAPIPort, "api-port", 8080, "Port to bind the network to")
	flag.StringVar(&flagAPIHost, "api-host", "", "Host to bind the network to")

	flag.IntVar(&flagGRPCPort, "grpc-port", 8082, "Port to bind the network to")
	flag.StringVar(&flagGRPCHost, "grpc-host", "", "Host to bind the network to")

	flag.Parse()

}

func main() {
	if flagAPIPort == flagGRPCPort {
		panic("The API and GRPC Can't share the same port!")
	}

	go func() {
		_, err := network.Serve(flagAPIHost, flagAPIPort, false)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		_, err := network.Serve(flagGRPCHost, flagGRPCPort, true)
		if err != nil {
			panic(err)
		}
	}()

	// Pretty much gonna wait forever.
	os.Exit(<-exitChan)
}
