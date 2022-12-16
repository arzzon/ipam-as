package backend

import (
	"context"
	"errors"
	"fmt"
	pb "github/arzzon/ipam-as/api"
	infblx "github/arzzon/ipam-as/pkg/manager/infoblox"
	"github/arzzon/ipam-as/pkg/types"
	"log"
)

// server is used to implement IPManagementServer
type Server struct {
	//pb.IPManagementServer
	pb.UnimplementedIPManagementServer
}

// AllocateIP implements IPManagementServer
func (s *Server) AllocateIP(ctx context.Context, in *pb.AllocateIPRequest) (*pb.AllocateIPResponse, error) {
	log.Printf("Received: %v", in.GetLabel())
	// CHECK IF IP ALREADY ASSIGNED
	infblxMngr := infblx.GetInfobloxManager()
	if infblxMngr == nil {
		return &pb.AllocateIPResponse{IP: "", Error: "Infoblox not initialized"}, errors.New("Infoblox not initialized")
	}
	req := types.IPAMRequest{
		HostName:  in.Hostname,
		IPAMLabel: in.Label,
	}
	newIP := infblxMngr.AllocateIP(req)
	return &pb.AllocateIPResponse{IP: newIP, Error: ""}, nil
}

// ReleaseIP implements IPManagementServer
func (s *Server) ReleaseIP(ctx context.Context, in *pb.ReleaseIPRequest) (*pb.ReleaseIPResponse, error) {
	// CHECK AND RELEASE IP
	infblxMngr := infblx.GetInfobloxManager()
	if infblxMngr == nil {
		return &pb.ReleaseIPResponse{Error: "Infoblox not initialized"}, errors.New("Infoblox not initialized")
	}
	req := types.IPAMRequest{
		IPAddr:    in.IP,
		IPAMLabel: in.Label,
	}
	err := infblxMngr.ReleaseIP(req)
	if err != nil {
		log.Printf("[ERROR] Unable to Release Address: %+v", req)
		return &pb.ReleaseIPResponse{Error: fmt.Sprint(err)}, err
	}
	log.Printf("[INFO] Released IP: %v", in.GetIP())
	return &pb.ReleaseIPResponse{Error: ""}, nil
}
