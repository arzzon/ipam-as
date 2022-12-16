package main

import (
	"flag"
	"fmt"
	"github/arzzon/ipam-as/pkg/backend"
	"github/arzzon/ipam-as/pkg/manager/infoblox"
	"github/arzzon/ipam-as/pkg/types"
	"log"
	"net"

	pb "github/arzzon/ipam-as/api"
	"google.golang.org/grpc"
)

const (
	DefaultProvider = "infoblox"
)

var (
	port     *int
	provider *string

	// Infoblox
	ibLabelMap *string
	ibHost     *string
	ibVersion  *string
	ibPort     *string
	ibUsername *string
	ibPassword *string
	ibNetView  *string
)

func init() {
	port = flag.Int("port", 50051, "The port on which IPAM server runs")
	provider = flag.String("ipam-provider", DefaultProvider,
		"Required, the IPAM system that the controller will interface with.")
	ibLabelMap = flag.String("infoblox-labels", "",
		"Required for mapping the infoblox's dnsview and cidr to IPAM labels")
	ibHost = flag.String("infoblox-grid-host", "",
		"Required for infoblox, the grid manager host IP.")
	ibVersion = flag.String("infoblox-wapi-version", "",
		"Required for infoblox, the Web API version.")
	ibPort = flag.String("infoblox-wapi-port", "443",
		"Optional for infoblox, the Web API port.")
	ibUsername = flag.String("infoblox-username", "",
		"Required for infoblox, the login username.")
	ibPassword = flag.String("infoblox-password", "",
		"Required for infoblox, the login password.")
	ibNetView = flag.String("infoblox-netview", "",
		"Required for allocation of IP addresses")
	//sslInsecure = ibFlags.Bool("insecure", false,
	//	"Optional, when set to true, enable insecure SSL communication to Infoblox.")
}

func main() {
	// PARSE FLAGS
	flag.Parse()
	// SETUP CONNECTION WITH INFOBLOX
	switch *provider {
	case DefaultProvider:
		params := types.Params{
			Host:       *ibHost,
			Version:    *ibVersion,
			Port:       *ibPort,
			Username:   *ibUsername,
			Password:   *ibPassword,
			IbLabelMap: *ibLabelMap,
			NetView:    *ibNetView,
		}
		_, err := infoblox.NewInfobloxManager(params)
		if err != nil {
			log.Fatalf("Failed to setup Infoblox: %v", err)
		}
	default:
		log.Fatalf("Failed to setup Provider")
	}
	// SETUP AND START GRPC SERVER
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterIPManagementServer(s, &backend.Server{})
	log.Printf("IPAM gRPC Server Listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
