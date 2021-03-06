package updatehub

import (
	"encoding/json"
	"fmt"

	"github.com/parnurzeal/gorequest"
	"github.com/updatehub/updatehub/metadata"
	"github.com/updatehub/updatehub/updatehub"
)

type Client struct {
}

type AgentInfo struct {
	Version  string                    `json:"version"`
	Config   updatehub.Settings        `json:"config"`
	Firmware metadata.FirmwareMetadata `json:"firmware"`
}

type LogEntry struct {
	Data    interface{} `json:"data"`
	Level   string      `json:"level"`
	Message string      `json:"message"`
	Time    string      `json:"time"`
}

type ProbeResponse struct {
	UpdateAvailable bool `json:"update-available"`
	TryAgainIn      int  `json:"try-again-in"`
}

// NewClient instantiates a new updatehub agent client
func NewClient() *Client {
	return &Client{}
}

// Probe default server address for update
func (c *Client) Probe() (*ProbeResponse, error) {
	return c.probe("")
}

// ProbeCustomServer probe custom server address for update
func (c *Client) ProbeCustomServer(serverAddress string) (*ProbeResponse, error) {
	return c.probe(serverAddress)
}

func (c *Client) probe(serverAddress string) (*ProbeResponse, error) {
	var probe ProbeResponse

	var req struct {
		ServerAddress string `json:"server-address"`
	}
	req.ServerAddress = serverAddress

	_, _, errs := gorequest.New().Post(buildURL("/probe")).Send(req).EndStruct(&probe)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return &probe, nil
}

// GetInfo get updatehub agent general information
func (c *Client) GetInfo() (*AgentInfo, error) {
	var info AgentInfo

	_, _, errs := gorequest.New().Get(buildURL("/info")).EndStruct(&info)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return &info, nil
}

// GetLogs get updatehub agent log entries
func (c *Client) GetLogs() ([]LogEntry, error) {
	_, body, errs := gorequest.New().Get(buildURL("/log")).End()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	var entries []LogEntry

	err := json.Unmarshal([]byte(body), &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func buildURL(path string) string {
	return fmt.Sprintf("http://localhost:8080/%s", path[1:])
}
