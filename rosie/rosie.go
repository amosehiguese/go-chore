package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	pb "github.com/amosehiguese/go-chore/proto/v1"
	"google.golang.org/grpc"
)

type Rosie struct {
	mu sync.Mutex
	chores []*pb.Chore
	pb.UnimplementedRobotMaidServer
}

func (r *Rosie) Add(_ context.Context, chores *pb.Chores) (*pb.Response, error) {
	r.mu.Lock()
	r.chores = append(r.chores, chores.Chores...)
	r.mu.Unlock()

	return &pb.Response{Message: "ok"}, nil
}

func (r * Rosie) Complete(_ context.Context, req *pb.CompleteRequest) (*pb.Response, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.chores == nil || req.ChoreNumber < 1 || int(req.ChoreNumber) > len(r.chores) {
		return nil, fmt.Errorf("chore %d not found", req.ChoreNumber)
	}

	r.chores[req.ChoreNumber-1].Complete = true
	return &pb.Response{Message: "ok"},nil
}

func (r *Rosie) List(_ context.Context, _ *pb.Empty) (*pb.Chores, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.chores == nil {
		r.chores = make([]*pb.Chore, 0)
	}

	return &pb.Chores{Chores: r.chores}, nil
}

var addr, certFn, keyFn string

func init() {
	flag.StringVar(&addr, "address", "localhost:34444", "listen address")
	flag.StringVar(&certFn, "cert", "cert.pem", "certificate file")
	flag.StringVar(&keyFn, "key", "key.pem", "private key file")
}

func main() {
	flag.Parse()

	server := grpc.NewServer()
	rosie := new(Rosie)

	pb.RegisterRobotMaidServer(server, rosie)

	cert, err := tls.LoadX509KeyPair(certFn, keyFn)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	go func ()  {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		for {
			if <-c == os.Interrupt {
				server.GracefulStop()
				log.Println("Server shutting down")
				return
			}
		}		
	}()

	fmt.Printf("listening for TLS connections on %s ...", addr)
	log.Fatal(
		server.Serve(tls.NewListener(listener, &tls.Config{
			Certificates: []tls.Certificate{cert},
			CurvePreferences: []tls.CurveID{tls.CurveP256},
			MinVersion: tls.VersionTLS12,
			PreferServerCipherSuites: true,
		})),
	)

}
