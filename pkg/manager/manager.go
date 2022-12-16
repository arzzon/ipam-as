package manager

//import (
//	"fmt"
//	"github/arzzon/ipam-as/pkg/manager/infoblox"
//	. "github/arzzon/ipam-as/pkg/types"
//	"log"
//)
//
//type Manager interface {
//	// Creates an A record
//	//CreateARecord(req IPAMRequest) bool
//	// Deletes an A record and releases the IP address
//	//DeleteARecord(req IPAMRequest)
//	// Gets IP Address associated with hostname/key
//	GetIPAddress(req IPAMRequest) string
//	// Gets and reserves the next available IP address
//	AllocateNextIPAddress(req IPAMRequest) string
//	// Releases an IP address
//	ReleaseIPAddress(req IPAMRequest)
//}
//
//const InfobloxProvider = "infoblox"
//
//type Params struct {
//	Provider       string
//	ProviderParams infoblox.InfobloxParams
//}
//
//func NewManager(params Params) (Manager, error) {
//	switch params.Provider {
//	case InfobloxProvider:
//		log.Printf("[DEBUG] Creating Manager with Provider: %v", InfobloxProvider)
//		ibxParams := params.ProviderParams
//		return infoblox.NewInfobloxManager(ibxParams), nil
//	default:
//		log.Printf("[DEBUG] Unknown Provider: %v", params.Provider)
//	}
//	return nil, fmt.Errorf("manager cannot be initialized")
//}
