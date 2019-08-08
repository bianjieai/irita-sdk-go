package rpc

import (
	"github.com/tendermint/tendermint/p2p"
)

type (
	ResultStatus struct {
		NodeInfo p2p.DefaultNodeInfo `json:"node_info"`
	}
)

func (c *client) GetStatus() (ResultStatus, error) {
	var (
		res ResultStatus
	)
	status, err := c.rpc.Status()
	if err != nil {
		return res, nil
	} else {
		res.NodeInfo = status.NodeInfo
		return res, nil
	}
}
