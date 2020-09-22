package integration_test

//import (
//	"fmt"
//	"time"
//
//	"github.com/stretchr/testify/require"
//	"github.com/bianjieai/irita-sdk-go/modules/ibc"
//	"github.com/bianjieai/irita-sdk-go/types"
//)
//
//func (s IntegrationTestSuite) TestIBC() {
//	ibcClient := s.Module(ibc.ModuleName).(ibc.IBCI)
//
//	baseTx := types.BaseTx{
//		From:     s.Account().Name,
//		Gas:      200000,
//		Memo:     "test",
//		Mode:     types.Commit,
//		Password: s.Account().Password,
//	}
//
//	recordPacket := ibc.RecordPacket{
//		ID:        "record1",
//		Timestamp: uint64(time.Now().Unix()),
//		Height:    0,
//		TxHash:    "",
//		Contents:  nil,
//		Creator:   "",
//	}
//
//	packetBZ := s.Codec().MustMarshalJSON(recordPacket)
//	request := ibc.SendPacketRequest{
//		ClientID:    "client1",
//		Module:      "ibcrecord",
//		Proof:       []byte("proof"),
//		ProofHeight: 10,
//		ProofPath:   []string{"1"},
//		ProofData:   []byte("proof"),
//		Packet: ibc.Packet{
//			Data: packetBZ,
//		},
//	}
//
//	res, err := ibcClient.SendPacket(request, baseTx)
//	require.NoError(s.T(), err)
//	require.NotEmpty(s.T(), res.Hash)
//
//	fmt.Println(res.Hash)
//}
//
//func (s IntegrationTestSuite) TestQueryClientState(){
//	ibcClient := s.Module(ibc.ModuleName).(ibc.IBCI)
//	clientState,err := ibcClient.QueryClientState("chaintest")
//	require.NoError(s.T(), err)
//	fmt.Println(clientState)
//}
