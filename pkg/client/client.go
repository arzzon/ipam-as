package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github/arzzon/ipam-as/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "test")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewIPManagementClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 360*time.Second)
	defer cancel()
	// Allocate
	allocateIP(c, ctx, "Test", "mytest1.com")

	// Release
	releaseIP(c, ctx, "10.9.0.2")
}

func allocateIP(c pb.IPManagementClient, ctx context.Context, label string, hostname string) {
	r, err := c.AllocateIP(ctx, &pb.AllocateIPRequest{Label: label, Hostname: hostname})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetIp())
}
func releaseIP(c pb.IPManagementClient, ctx context.Context, ip string) {
	r1, err := c.ReleaseIP(ctx, &pb.ReleaseIPRequest{Ip: ip})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("IP Released: %v", r1.Error)
}
