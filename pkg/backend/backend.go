package backend

import (
	"context"
	"errors"
	pb "github/arzzon/ipam-as/api"
	infblx "github/arzzon/ipam-as/pkg/manager/infoblox"
	"github/arzzon/ipam-as/pkg/types"
	"log"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	//pb.IPManagementServer
	pb.UnimplementedIPManagementServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) AllocateIP(ctx context.Context, in *pb.AllocateIPRequest) (*pb.AllocateIPResponse, error) {
	log.Printf("Received: %v", in.GetLabel())
	// CHECK IF IP ALREADY ASSIGNED
	infblxMngr := infblx.GetInfobloxManager()
	if infblxMngr == nil {
		return &pb.AllocateIPResponse{Ip: "", Error: "Infoblox not initialized"}, errors.New("Infoblox not initialized")
	}
	req := types.IPAMRequest{
		HostName:  in.Hostname,
		IPAMLabel: in.Label,
	}
	newIP := infblxMngr.AllocateIP(req)
	return &pb.AllocateIPResponse{Ip: newIP, Error: ""}, nil
}

// SayHello implements helloworld.GreeterServer
func (s *Server) ReleaseIP(ctx context.Context, in *pb.ReleaseIPRequest) (*pb.ReleaseIPResponse, error) {
	log.Printf("Received: %v", in.GetIp())
	// CHECK AND RELEASE IP
	infblxMngr := infblx.GetInfobloxManager()
	if infblxMngr == nil {
		return &pb.ReleaseIPResponse{Error: "Infoblox not initialized"}, errors.New("Infoblox not initialized")
	}
	req := types.IPAMRequest{
		IPAddr: in.Ip,
	}
	err := infblxMngr.ReleaseIP(req)
	if err != nil {
		log.Printf("[ERROR] Unable to Release Address: %+v", req)
	}
	return &pb.ReleaseIPResponse{Error: ""}, nil
}
