package types

type IPAMRequest struct {
	//Metadata  interface{}
	//Operation string
	HostName  string
	IPAddr    string
	Key       string
	IPAMLabel string
}

type IPAMResponse struct {
	Request IPAMRequest
	IPAddr  string
	Status  bool
}

type Params struct {
	Host       string
	Version    string
	Port       string
	Username   string
	Password   string
	SslVerify  string
	NetView    string
	IbLabelMap string
}
