package types

type NodeFilterRequest struct {
	Nodes []string `json:"nodes"`
}

type AccounGroupRequest struct {
	Accounts []string `json:"accounts"`
}