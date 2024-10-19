package spacemesh_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: baseURL,
	}
}

// Common types
type Account struct {
	Address  string        `json:"address"`
	Current  AccountState  `json:"current"`
	Projected AccountState `json:"projected"`
	Template string        `json:"template"`
}

type AccountState struct {
	Balance string `json:"balance"`
	Counter string `json:"counter"`
	Layer   int64  `json:"layer"`
}

type Activation struct {
	Coinbase     string `json:"coinbase"`
	Height       string `json:"height"`
	ID           []byte `json:"id"`
	NumUnits     int64  `json:"numUnits"`
	PublishEpoch int64  `json:"publishEpoch"`
	SmesherID    []byte `json:"smesherId"`
	Weight       string `json:"weight"`
}

type Layer struct {
	Number              int64  `json:"number"`
	Status              string `json:"status"`
	Hash                []byte `json:"hash"`
	BlockID             []byte `json:"blockId"`
	StateHash           []byte `json:"stateHash"`
	CumulativeStateHash []byte `json:"cumulativeStateHash"`
}

type MalfeasanceProof struct {
	SmesherID  []byte            `json:"smesher"`
	Layer      int64             `json:"layer"`
	Kind       string            `json:"kind"`
	Proof      []byte            `json:"proof"`
	Properties map[string]string `json:"properties"`
}

type Reward struct {
	Layer       int64  `json:"layer"`
	Total       string `json:"total"`
	LayerReward string `json:"layerReward"`
	Coinbase    string `json:"coinbase"`
	Smesher     []byte `json:"smesher"`
}

type Transaction struct {
	ID       []byte               `json:"id"`
	Principal string              `json:"principal"`
	Template string               `json:"template"`
	Method   int64                `json:"method"`
	Nonce    Nonce                `json:"nonce"`
	Type     string               `json:"type"`
	GasPrice string               `json:"gasPrice"`
	MaxGas   string               `json:"maxGas"`
	MaxSpend string               `json:"maxSpend"`
	Raw      []byte               `json:"raw"`
}

type Nonce struct {
	Counter string `json:"counter"`
}

// AccountService endpoints
func (c *Client) AccountList(req AccountRequest) (*AccountList, error) {
	resp := &AccountList{}
	err := c.sendRequest("POST", "/spacemesh.v2alpha1.AccountService/List", req, resp)
	return resp, err
}

// ActivationService endpoints
func (c *Client) ActivationsList(req ActivationRequest) (*ActivationList, error) {
	resp := &ActivationList{}
	err := c.sendRequest("POST", "/spacemesh.v2alpha1.ActivationService/List", req, resp)
	return resp, err
}

func (c *Client) ActivationsCount(req ActivationsCountRequest) (*ActivationsCountResponse, error) {
	resp := &ActivationsCountResponse{}
	err := c.sendRequest("POST", "/spacemesh.v2alpha1.ActivationService/ActivationsCount", req, resp)
	return resp, err
}

// LayerService endpoints
func (c *Client) LayerList(req LayerRequest) (*LayerList, error) {
	resp := &LayerList{}
	err := c.sendRequest("POST", "/spacemesh.v2alpha1.LayerService/List", req, resp)
	return resp, err
}

// MalfeasanceService endpoints
func (c *Client) MalfeasanceList(req MalfeasanceRequest) (*MalfeasanceList, error) {
	resp := &MalfeasanceList{}
	err := c.sendRequest("POST", "/spacemesh.v2alpha1.MalfeasanceService/List", req, resp)
	return resp, err
}

// NetworkService endpoints
func (c *Client) NetworkInfo(req NetworkInfoRequest) (*NetworkInfoResponse, error) {
	resp := &NetworkInfoResponse{}
	err := c.sendRequest("POST", "/spacemesh.v2alpha1.NetworkService/Info", req, resp)
	return resp, err
}

// NodeService endpoints
func (c *Client) NodeStatus(req NodeStatusRequest) (*NodeStatusResponse, error) {
	resp := &NodeStatusResponse{}
	err := c.sendRequest("POST", "/spacemesh.v2alpha1.NodeService/Status", req, resp)
	return resp, err
}

// RewardService endpoints
func (c *Client) RewardList(req RewardRequest) (*RewardList, error) {
	resp := &RewardList{}
	err := c.sendRequest("POST", "/spacemesh.v2alpha1.RewardService/List", req, resp)
	return resp, err
}

// TransactionService endpoints
func (c *Client) TransactionList(req TransactionRequest) (*TransactionList, error) {
	resp := &TransactionList{}
	err := c.sendRequest("POST", "/spacemesh.v2alpha1.TransactionService/List", req, resp)
	return resp, err
}

func (c *Client) EstimateGas(req EstimateGasRequest) (*EstimateGasResponse, error) {
	resp := &EstimateGasResponse{}
	err := c.sendRequest("POST", "/spacemesh.v2alpha1.TransactionService/EstimateGas", req, resp)
	return resp, err
}

func (c *Client) ParseTransaction(req ParseTransactionRequest) (*ParseTransactionResponse, error) {
	resp := &ParseTransactionResponse{}
	err := c.sendRequest("POST", "/spacemesh.v2alpha1.TransactionService/ParseTransaction", req, resp)
	return resp, err
}

func (c *Client) SubmitTransaction(req SubmitTransactionRequest) (*SubmitTransactionResponse, error) {
	resp := &SubmitTransactionResponse{}
	err := c.sendRequest("POST", "/spacemesh.v2alpha1.TransactionService/SubmitTransaction", req, resp)
	return resp, err
}

// Helper function to send HTTP requests
func (c *Client) sendRequest(method, path string, req interface{}, resp interface{}) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := c.baseURL + path
	httpReq, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code: %d", httpResp.StatusCode)
	}

	return json.NewDecoder(httpResp.Body).Decode(resp)
}

// Request and response types for each endpoint
type AccountRequest struct {
	Addresses []string `json:"addresses,omitempty"`
	Limit     string   `json:"limit,omitempty"`
	Offset    string   `json:"offset,omitempty"`
}

type AccountList struct {
	Accounts []Account `json:"accounts"`
}

type ActivationRequest struct {
	Coinbase   string   `json:"coinbase,omitempty"`
	EndEpoch   int64    `json:"endEpoch,omitempty"`
	ID         [][]byte `json:"id,omitempty"`
	Limit      string   `json:"limit,omitempty"`
	Offset     string   `json:"offset,omitempty"`
	SmesherID  [][]byte `json:"smesherId,omitempty"`
	StartEpoch int64    `json:"startEpoch,omitempty"`
}

type ActivationList struct {
	Activations []Activation `json:"activations"`
}

type ActivationsCountRequest struct {
	Epoch int64 `json:"epoch"`
}

type ActivationsCountResponse struct {
	Count int64 `json:"count"`
}

type LayerRequest struct {
	EndLayer  int64  `json:"endLayer,omitempty"`
	Limit     string `json:"limit,omitempty"`
	Offset    string `json:"offset,omitempty"`
	SortOrder string `json:"sortOrder,omitempty"`
	StartLayer int64 `json:"startLayer,omitempty"`
}

type LayerList struct {
	Layers []Layer `json:"layers"`
}

type MalfeasanceRequest struct {
	Limit     string   `json:"limit,omitempty"`
	Offset    string   `json:"offset,omitempty"`
	SmesherID [][]byte `json:"smesherId,omitempty"`
}

type MalfeasanceList struct {
	Proofs []MalfeasanceProof `json:"proofs"`
}

type NetworkInfoRequest struct{}

type NetworkInfoResponse struct {
	EffectiveGenesisLayer int64     `json:"effectiveGenesisLayer"`
	GenesisID             []byte    `json:"genesisId"`
	GenesisTime           time.Time `json:"genesisTime"`
	HRP                   string    `json:"hrp"`
	LabelsPerUnit         string    `json:"labelsPerUnit"`
	LayerDuration         string    `json:"layerDuration"`
	LayersPerEpoch        int64     `json:"layersPerEpoch"`
}

type NodeStatusRequest struct{}

type NodeStatusResponse struct {
	AppliedLayer    int64  `json:"appliedLayer"`
	ConnectedPeers  string `json:"connectedPeers"`
	CurrentLayer    int64  `json:"currentLayer"`
	LatestLayer     int64  `json:"latestLayer"`
	ProcessedLayer  int64  `json:"processedLayer"`
	Status          string `json:"status"`
}

type RewardRequest struct {
	Coinbase   string `json:"coinbase,omitempty"`
	EndLayer   int64  `json:"endLayer,omitempty"`
	Limit      string `json:"limit,omitempty"`
	Offset     string `json:"offset,omitempty"`
	Smesher    []byte `json:"smesher,omitempty"`
	SortOrder  string `json:"sortOrder,omitempty"`
	StartLayer int64  `json:"startLayer,omitempty"`
}

type RewardList struct {
	Rewards []Reward `json:"rewards"`
}

type TransactionRequest struct {
	Address       string   `json:"address,omitempty"`
	EndLayer      int64    `json:"endLayer,omitempty"`
	IncludeResult bool     `json:"includeResult,omitempty"`
	IncludeState  bool     `json:"includeState,omitempty"`
	Limit         string   `json:"limit,omitempty"`
	Offset        string   `json:"offset,omitempty"`
	SortOrder     string   `json:"sortOrder,omitempty"`
	StartLayer    int64    `json:"startLayer,omitempty"`
	TxID          [][]byte `json:"txid,omitempty"`
}

type TransactionList struct {
	Transactions []TransactionResponse `json:"transactions"`
}

type TransactionResponse struct {
	Tx       Transaction        `json:"tx"`
	TxResult TransactionResult  `json:"txResult,omitempty"`
	TxState  string             `json:"txState,omitempty"`
}

type TransactionResult struct {
	Block            []byte   `json:"block"`
	Fee              string   `json:"fee"`
	GasConsumed      string   `json:"gasConsumed"`
	Layer            int64    `json:"layer"`
	Message          string   `json:"message"`
	Status           string   `json:"status"`
	TouchedAddresses []string `json:"touchedAddresses"`
}

type EstimateGasRequest struct {
	Transaction []byte `json:"transaction"`
}

type EstimateGasResponse struct {
	RecommendedMaxGas string `json:"recommendedMaxGas"`
	Status            Status `json:"status"`
}

type ParseTransactionRequest struct {
	Transaction []byte `json:"transaction"`
	Verify      bool   `json:"verify"`
}

type ParseTransactionResponse struct {
	Status Status      `json:"status"`
	Tx     Transaction `json:"tx"`
}

type SubmitTransactionRequest struct {
	Transaction []byte `json:"transaction"`
}

type SubmitTransactionResponse struct {
	Status Status `json:"status"`
	TxID   []byte `json:"txId"`
}

type Status struct {
	Code    int32             `json:"code"`
	Details []json.RawMessage `json:"details"`
	Message string            `json:"message"`
}