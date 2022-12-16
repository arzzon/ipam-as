package backend

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/arzzon/ipam-as/api"
	infblx "github.com/arzzon/ipam-as/pkg/manager/infoblox"
	"github.com/arzzon/ipam-as/pkg/types"
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
	if len(newIP) == 0 {
		return &pb.AllocateIPResponse{IP: newIP, Error: ""}, nil
	}
	return &pb.AllocateIPResponse{IP: newIP, Error: "Failed to allocated IP"}, errors.New("Failed to allocated IP")
}

// ReleaseIP implements IPManagementServer
func (s *Server) ReleaseIP(ctx context.Context, in *pb.ReleaseIPRequest) (*pb.ReleaseIPResponse, error) {
	// CHECK AND RELEASE IP
	infblxMngr := infblx.GetInfobloxManager()
	if infblxMngr == nil {
		return &pb.ReleaseIPResponse{Error: "Infoblox not initialized"}, errors.New("Infoblox not initialized")
	}
	req := types.IPAMRequest{
		HostName:  in.Hostname,
		IPAMLabel: in.Label,
	}
	err := infblxMngr.ReleaseIP(req)
	if err != nil {
		log.Printf("[ERROR] Unable to Release Address: %+v", req)
		return &pb.ReleaseIPResponse{Error: fmt.Sprint(err)}, err
	}
	log.Printf("[INFO] Released for Host: %v", in.GetHostname())
	return &pb.ReleaseIPResponse{Error: ""}, nil
}
