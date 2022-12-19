package infoblox

import (
	"encoding/json"
	"errors"
	"github.com/arzzon/ipam-as/pkg/types"
	ibxclient "github.com/infobloxopen/infoblox-go-client"
	"log"
)

const (
	EAKey          = "F5IPAM"
	EAVal          = "managed"
	netview string = "default"
)

type InfobloxManager struct {
	connector   *ibxclient.Connector
	objMgr      *ibxclient.ObjectManager
	ea          ibxclient.EA
	NetView     string
	IBLabels    map[string]IBConfig
	HostIPCache map[string]string
}

type IBConfig struct {
	DNSView string `json:"dnsView,omitempty"`
	CIDR    string `json:"cidr"`
}

var InfblxMngr *InfobloxManager

func GetInfobloxManager() *InfobloxManager {
	return InfblxMngr
}
func NewInfobloxManager(params types.Params) (*InfobloxManager, error) {
	if InfblxMngr != nil {
		return InfblxMngr, nil
	}
	// CREATE NEW INFOBLOX MANAGER
	hostConfig := ibxclient.HostConfig{
		Host:     params.Host,
		Version:  params.Version,
		Port:     params.Port,
		Username: params.Username,
		Password: params.Password,
	}

	labels, err := ParseLabels(params.IbLabelMap)
	if err != nil {
		return nil, err
	}

	// TransportConfig params: sslVerify, httpRequestsTimeout, httpPoolConnections
	// These are the common values
	transportConfig := ibxclient.NewTransportConfig("false", 20, 10)
	requestBuilder := &ibxclient.WapiRequestBuilder{}
	requestor := &ibxclient.WapiHttpRequestor{}
	connector, err := ibxclient.NewConnector(hostConfig, transportConfig, requestBuilder, requestor)
	if err != nil {
		return nil, err
	}
	objMgr := ibxclient.NewObjectManager(connector, "F5IPAM", "0")

	objMgr.OmitCloudAttrs = true

	// Create an Extensible Attribute for resource tracking
	if eaDef, _ := objMgr.GetEADefinition(EAKey); eaDef == nil {
		eaDef := ibxclient.EADefinition{
			Name:    EAKey,
			Type:    "STRING",
			Comment: "Managed by the F5 IPAM Controller",
		}
		_, err = objMgr.CreateEADefinition(eaDef)
		if err != nil {
			return nil, err
		}
	}
	ibMgr := &InfobloxManager{
		connector:   connector,
		objMgr:      objMgr,
		ea:          ibxclient.EA{EAKey: EAVal},
		IBLabels:    labels,
		NetView:     params.NetView,
		HostIPCache: make(map[string]string),
	}

	// Validating that dnsView, CIDR exist on infoblox Server
	for _, parameter := range labels {
		result, err := ibMgr.validateIPAMLabels(parameter.DNSView, parameter.CIDR)
		if !result {
			return nil, err
		}
	}
	InfblxMngr = ibMgr

	return ibMgr, nil
}

func (infMgr *InfobloxManager) validateIPAMLabels(dnsView, cidr string) (bool, error) {
	_, err := infMgr.objMgr.GetNetwork(infMgr.NetView, cidr, infMgr.ea)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (infMgr *InfobloxManager) AllocateIP(req types.IPAMRequest) string {
	// Check if IP is already assigned for the hostname
	ipAddr := infMgr.getIPAddressFromName(req)
	if ipAddr != "" {
		return ipAddr
	}
	label, ok := infMgr.IBLabels[req.IPAMLabel]
	if !ok {
		return ""
	}
	name := req.HostName
	if req.Key != "" {
		name = req.Key
	}
	//fixedAddr, err := infMgr.objMgr.AllocateIP(infMgr.NetView, label.CIDR, "", "", name, infMgr.ea)
	fixedAddr, err := infMgr.objMgr.AllocateIP(infMgr.NetView, label.CIDR, "", "", name, infMgr.ea)
	if err != nil {
		log.Printf("[ERROR] Unable to get a new IP address: %+v", req)
		return ""
	}
	infMgr.HostIPCache[req.HostName] = fixedAddr.IPAddress
	log.Printf("[DEBUG] Allocated new IP address for the request: %+v, IP:%s", req, fixedAddr.IPAddress)
	return fixedAddr.IPAddress
}

func (infMgr *InfobloxManager) ReleaseIP(req types.IPAMRequest) error {
	if req.IPAMLabel == "" {
		return errors.New("IPAM Label is required for releasing an IP")
	}
	label, ok := infMgr.IBLabels[req.IPAMLabel]
	if !ok {
		return errors.New("IPAM Label is not found")
	}
	ip := ""
	if req.IPAddr != "" {
		ip = req.IPAddr
	} else if req.HostName != "" {
		if mapped_ip, ok := infMgr.HostIPCache[req.HostName]; ok {
			ip = mapped_ip
		} else {
			return errors.New("No IP assigned to Host")
		}
	}
	_, err := infMgr.objMgr.ReleaseIP(infMgr.NetView, label.CIDR, ip, "")
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] IP associated with host key %v has been released.", req.HostName)
	delete(infMgr.HostIPCache, req.HostName)
	return nil
}

func (mgr *InfobloxManager) GetNetworkView(name string) ([]ibxclient.Network, error) {
	var res []ibxclient.Network
	network := ibxclient.NewNetwork(ibxclient.Network{NetviewName: netview})
	mgr.connector.GetObject(network, "", &res)
	return res, nil
}

func (infMgr *InfobloxManager) getIPAddressFromName(req types.IPAMRequest) (ip string) {
	var returnFixedAddresses []ibxclient.FixedAddress

	label, ok := infMgr.IBLabels[req.IPAMLabel]
	if !ok {
		return ""
	}

	name := req.HostName
	if req.Key != "" {
		name = req.Key
	}

	fixedAddr := ibxclient.NewFixedAddress(ibxclient.FixedAddress{
		NetviewName: infMgr.NetView,
		Cidr:        label.CIDR,
	})

	err := infMgr.connector.GetObject(fixedAddr, "", &returnFixedAddresses)

	if err != nil || returnFixedAddresses == nil || len(returnFixedAddresses) == 0 {
		log.Printf("[ERROR] IP not available, %+v", req)
		return ""
	}

	for _, fixedAddress := range returnFixedAddresses {
		if fixedAddress.Name == name {
			return fixedAddress.IPAddress
		}
	}
	return ""
}

func ParseLabels(params string) (map[string]IBConfig, error) {
	ibLabelMap := make(map[string]IBConfig)
	err := json.Unmarshal([]byte(params), &ibLabelMap)
	if err != nil {
		return nil, err
	}
	for label, ibParam := range ibLabelMap {
		// DNSView is being disabled
		// The below line can be removed when DNSView support is enabled
		ibParam.DNSView = ""

		ibLabelMap[label] = ibParam
	}
	return ibLabelMap, nil
}
