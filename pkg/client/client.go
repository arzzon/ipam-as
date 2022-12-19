package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/arzzon/ipam-as/api"
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
	//allocateIP(c, ctx, "Test", "mytest1.com")

	// Release
	releaseIP(c, ctx, "foo.com", "default")
}

func allocateIP(c pb.IPManagementClient, ctx context.Context, label string, hostname string) {
	r, err := c.AllocateIP(ctx, &pb.AllocateIPRequest{Label: label, Hostname: hostname})
	if err != nil {
		log.Fatalf("failed to allocated IP: %v", err)
	}
	if r.Error != "" {
		log.Fatalf("failed to allocated IP: %v", err)
	} else {
		log.Printf("Allocated IP: %s", r.GetIP())
	}
}
func releaseIP(c pb.IPManagementClient, ctx context.Context, host string, label string) {
	r1, err := c.ReleaseIP(ctx, &pb.ReleaseIPRequest{Hostname: host, Label: label})
	if err != nil {
		log.Printf("failed to release IP: %v", err)
	}
	if r1.Error != "" {
		log.Printf("failed to release IP: %v", r1.Error)
	} else {
		log.Printf("IP Released: %v", r1.Error)
	}
}
