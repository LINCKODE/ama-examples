package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/qubic/go-node-connector/types"
	"github.com/qubic/go-schnorrq"
	"io"
	"net/http"
	"time"
)

const ApiURL = "https://testapi.qubic.org"

const SenderSeed = ""
const SenderID = ""
const TransactionAmount = 5
const DestinationID = ""

func CreateAndBroadcastTransaction() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	subSeed, err := types.GetSubSeed(SenderSeed)
	if err != nil {
		return errors.Wrap(err, "getting sub-seed from seed")
	}

	latestTick, err := GetLatestTick(ctx, ApiURL)
	if err != nil {
		return errors.Wrap(err, "getting latest tick")
	}

	targetTick := latestTick + 10

	transaction, err := CreateTransaction(SenderID, DestinationID, TransactionAmount, targetTick, subSeed)
	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	err = BroadCastTransaction(ctx, transaction, ApiURL)
	if err != nil {
		return errors.Wrap(err, "broadcasting transaction")
	}

	return nil
}

func CreateTransaction(sourceID, destinationID string, amount int64, targetTick uint32, subSeed [32]byte) (*types.Transaction, error) {
	transaction, err := types.NewSimpleTransferTransaction(sourceID, destinationID, amount, targetTick)
	if err != nil {
		return nil, errors.Wrap(err, "creating transaction")
	}

	unsignedDigest, err := transaction.GetUnsignedDigest()
	if err != nil {
		return nil, errors.Wrap(err, "getting transaction unsigned digest")
	}

	signature, err := schnorrq.Sign(subSeed, transaction.SourcePublicKey, unsignedDigest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign transaction")
	}
	transaction.Signature = signature

	return &transaction, nil

}

func GetLatestTick(ctx context.Context, apiUrl string) (uint32, error) {

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl+"/v1/latestTick", nil)
	if err != nil {
		return 0, errors.Wrap(err, "creating latest tick request")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, errors.Wrap(err, "performing latest tick request")
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, errors.Wrap(err, "reading latest tick response")
	}

	if response.StatusCode != http.StatusOK {
		return 0, errors.New(fmt.Sprintf("Status not 200! Info: %s", data))
	}

	type LatestTickStruct struct {
		LatestTick uint32 `json:"latestTick"`
	}

	var latestTick LatestTickStruct
	err = json.Unmarshal(data, &latestTick)
	if err != nil {
		return 0, errors.Wrap(err, "un-marshalling latest tick response")
	}

	return latestTick.LatestTick, nil
}

func BroadCastTransaction(ctx context.Context, transaction *types.Transaction, apiUrl string) error {

	transactionID, err := transaction.ID()
	if err != nil {
		return errors.Wrap(err, "getting transaction ID")
	}

	finalUrl := apiUrl + "/v1/broadcast-transaction"

	// Encode transaction
	encodedTransaction, err := transaction.EncodeToBase64()
	if err != nil {
		return errors.Wrap(err, "encoding transaction")
	}

	// Create request payload from encoded transaction
	requestPayload := struct {
		EncodedTransaction string `json:"encodedTransaction"`
	}{
		EncodedTransaction: encodedTransaction,
	}
	buffer := new(bytes.Buffer)

	err = json.NewEncoder(buffer).Encode(requestPayload)
	if err != nil {
		return errors.Wrap(err, "encoding request payload")
	}

	//Create and send request
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, finalUrl, buffer)
	if err != nil {
		return errors.Wrap(err, "creating broadcast request")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "performing broadcast request")
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "reading broadcast response")
	}

	if response.StatusCode != http.StatusOK {
		fmt.Printf("DEBUG: %s\n", finalUrl)
		return errors.New(fmt.Sprintf("Status not 200! Info: %s", data))
	}

	type ResponseInfo struct {
		PeersBroadcasted   uint32 `json:"peersBroadcasted"`
		EncodedTransaction string `json:"encodedTransaction"`
	}

	var info ResponseInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		return errors.Wrap(err, "un-marshalling broadcast response")
	}

	fmt.Printf("Broadcasted to %d peers.\n", info.PeersBroadcasted)
	fmt.Printf("Transaction ID: %s\n", transactionID)
	fmt.Printf("Target tick: %d\n", transaction.Tick)

	return nil
}
