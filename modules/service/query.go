package service

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// queryRequestContextByTxQuery will query for a single request context via a direct txs tags query.
func (s serviceClient) queryRequestContextByTxQuery(reqCtxID string) (RequestContext, error) {
	txHash, msgIndex, err := splitRequestContextID(reqCtxID)
	if err != nil {
		return RequestContext{}, err
	}

	txInfo, err := s.QueryTx(hex.EncodeToString(txHash))
	if err != nil {
		return RequestContext{}, err
	}

	if int64(len(txInfo.Tx.Body.Msgs)) > msgIndex {
		msg := txInfo.Tx.Body.Msgs[msgIndex]
		if msg, ok := msg.(*MsgCallService); ok {
			return RequestContext{
				ServiceName:        msg.ServiceName,
				Providers:          msg.Providers,
				Consumer:           msg.Consumer,
				Input:              msg.Input,
				ServiceFeeCap:      msg.ServiceFeeCap,
				Timeout:            msg.Timeout,
				Repeated:           msg.Repeated,
				RepeatedFrequency:  msg.RepeatedFrequency,
				RepeatedTotal:      msg.RepeatedTotal,
				BatchCounter:       uint64(msg.RepeatedTotal),
				BatchRequestCount:  0,
				BatchResponseCount: 0,
				BatchState:         BATCHCOMPLETED,
				State:              COMPLETED,
				ResponseThreshold:  0,
				ModuleName:         "",
			}, nil
		}
	}
	return RequestContext{}, fmt.Errorf("invalid reqCtxID:%s", reqCtxID)
}

// queryRequestByTxQuery will query for a single request via a direct txs tags query.
func (s serviceClient) queryRequestByTxQuery(requestID string) (Request, error) {
	reqCtxID, _, requestHeight, batchRequestIndex, err := splitRequestID(requestID)
	if err != nil {
		return Request{}, err
	}

	// query request context
	reqCtx, err := s.QueryRequestContext(hex.EncodeToString(reqCtxID))
	if err != nil {
		return Request{}, err
	}

	blockResult, err := s.BlockResults(context.Background(), &requestHeight)
	if err != nil {
		return Request{}, err
	}

	for _, event := range blockResult.EndBlockEvents {
		if event.Type == eventTypeNewBatchRequest {
			var found bool
			var requestsBz []byte
			for _, attribute := range event.Attributes {
				if string(attribute.Key) == attributeKeyRequests {
					requestsBz = attribute.GetValue()
				}
				if string(attribute.Key) == attributeKeyRequestContextID &&
					string(attribute.GetValue()) == reqCtxID.String() {
					found = true
				}
			}

			if found {
				var requests []CompactRequest
				if err := json.Unmarshal(requestsBz, &requests); err != nil {
					return Request{}, err
				}
				if len(requests) > int(batchRequestIndex) {
					compactRequest := requests[batchRequestIndex]
					return Request{
						Id:                         sdk.MustHexBytesFrom(requestID).String(),
						ServiceName:                reqCtx.ServiceName,
						Provider:                   compactRequest.Provider,
						Consumer:                   reqCtx.Consumer,
						Input:                      reqCtx.Input,
						ServiceFee:                 compactRequest.ServiceFee,
						SuperMode:                  reqCtx.SuperMode,
						RequestHeight:              compactRequest.RequestHeight,
						ExpirationHeight:           compactRequest.RequestHeight + reqCtx.Timeout,
						RequestContextId:           compactRequest.RequestContextId,
						RequestContextBatchCounter: compactRequest.RequestContextBatchCounter,
					}, nil
				}
			}
		}
	}
	return Request{}, fmt.Errorf("invalid requestID:%s", requestID)
}

// queryResponseByTxQuery will query for a single request via a direct txs tags query.
func (s serviceClient) queryResponseByTxQuery(requestID string) (Response, error) {
	builder := sdk.NewEventQueryBuilder().AddCondition(
		sdk.NewCond(
			sdk.EventTypeMessage,
			sdk.AttributeKeyAction,
		).EQ("respond_service"),
	).AddCondition(
		sdk.NewCond(
			sdk.EventTypeMessage,
			sdk.AttributeKeyAction,
		).EQ(attributeKeyRequestID),
	)

	result, err := s.QueryTxs(builder, 1, 1)
	if err != nil {
		return Response{}, err
	}

	if len(result.Txs) == 0 {
		return Response{}, fmt.Errorf("unknown response: %s", requestID)
	}

	reqCtxID, batchCounter, _, _, err := splitRequestID(requestID)
	if err != nil {
		return Response{}, err
	}

	// query request context
	reqCtx, err := s.QueryRequestContext(hex.EncodeToString(reqCtxID))
	if err != nil {
		return Response{}, err
	}

	for _, msg := range result.Txs[0].Tx.Body.Msgs {
		if responseMsg, ok := msg.(*MsgRespondService); ok {
			if responseMsg.RequestId != requestID {
				continue
			}
			return Response{
				Provider:                   responseMsg.Provider,
				Consumer:                   reqCtx.Consumer,
				Output:                     responseMsg.Output,
				Result:                     responseMsg.Result,
				RequestContextId:           sdk.HexStringFrom(reqCtxID),
				RequestContextBatchCounter: batchCounter,
			}, nil
		}
	}

	return Response{}, nil
}

// SplitRequestContextID splits the given contextID to txHash and msgIndex
func splitRequestContextID(reqCtxID string) (sdk.HexBytes, int64, error) {
	contextID, err := hex.DecodeString(reqCtxID)
	if err != nil {
		return nil, 0, errors.New("invalid request context id")
	}

	if len(contextID) != contextIDLen {
		return nil, 0, fmt.Errorf("invalid request context id:%s", reqCtxID)
	}

	txHash := contextID[0:32]
	msgIndex := int64(binary.BigEndian.Uint64(contextID[32:40]))

	return txHash, msgIndex, nil
}

// SplitRequestID splits the given contextID to contextID, batchCounter, requestHeight, batchRequestIndex
func splitRequestID(reqID string) (sdk.HexBytes, uint64, int64, int16, error) {
	requestID, err := hex.DecodeString(reqID)
	if err != nil {
		return nil, 0, 0, 0, errors.New("invalid request id")
	}

	if len(requestID) != requestIDLen {
		return nil, 0, 0, 0, errors.New("invalid request id")
	}

	reqCtxID := requestID[0:40]
	batchCounter := binary.BigEndian.Uint64(requestID[40:48])
	requestHeight := int64(binary.BigEndian.Uint64(requestID[48:56]))
	batchRequestIndex := int16(binary.BigEndian.Uint16(requestID[56:]))

	return reqCtxID, batchCounter, requestHeight, batchRequestIndex, nil
}
