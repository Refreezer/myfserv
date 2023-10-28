package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

func main() {
	port := flag.String("port", "3000", "set listen port")
	root := flag.String("root", "C:\\", "set Server root dir")
	flag.Parse()

	PrintLocalIps()
	log.Printf("Listening port %s\nRoot dir is %s\n", *port, *root)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	s := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", *port),
		Handler: http.FileServer(http.Dir(*root)),
	}

	go handleSignal(s, sig)

	err := s.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln(err)
	}
}

func handleSignal(serv *http.Server, sig chan os.Signal) {
	s := <-sig
	log.Printf("Signal %v arrived", s)
	err := serv.Shutdown(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}
}

func PrintLocalIps() /*net.IP*/ {
	ipRegexp, _ := regexp.Compile(`^(\d{1,4}\.){3}\d{1,4}`)
	networkInterfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, networkInterface := range networkInterfaces {
		if networkInterface.Flags&net.FlagUp == 0 || networkInterface.Flags&net.FlagRunning == 0 {
			continue
		}
		addrs, err := networkInterface.Addrs()
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println(networkInterface.Name)
		for _, addr := range addrs {
			ip := ipRegexp.Find([]byte(addr.String()))
			if ip == nil {
				continue
			}

			log.Println(string(ip))
		}
	}
}
