package archiver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

const archiverHTTPAddress = "https://testapi.qubic.org"

func RunHTTPExample() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	httpClient := http.DefaultClient

	// Get Status
	status, err := FetchStatusHTTP(ctx, httpClient, archiverHTTPAddress)
	if err != nil {
		return errors.Wrap(err, "fetching archiver status")
	}

	fmt.Printf("Last processed tick: %d\n", status.LastProcessedTick.TickNumber)
	fmt.Printf("Current epoch: %d\n", status.LastProcessedTick.Epoch)

	// Get transactions in the latest tick

	transactions, err := GetTickTransactionsHTTP(ctx, httpClient, archiverHTTPAddress, status.LastProcessedTick.TickNumber)
	if err != nil {
		return errors.Wrap(err, "fetching tick transactions")
	}

	fmt.Printf("Found %d transactions\n", len(transactions.Transactions))

	for _, transaction := range transactions.Transactions {
		fmt.Printf("  Transaction ID: %s\n", transaction.Transaction.TxId)
		fmt.Printf("  Transfered amount: %s\n", transaction.Transaction.Amount)
		fmt.Printf("  From: %s\n", transaction.Transaction.SourceId)
		fmt.Printf("  To: %s\n", transaction.Transaction.DestId)
		println("-------------------------------")

	}

	return nil

}

func FetchStatusHTTP(ctx context.Context, httpClient *http.Client, baseURL string) (*Status, error) {

	statusRequest, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/v1/status", nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating status request")
	}

	response, err := httpClient.Do(statusRequest)
	if err != nil {
		return nil, errors.Wrap(err, "performing status request")
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading status response")
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Status not 200! Info: %s", data))
	}

	var status Status

	err = json.Unmarshal(data, &status)
	if err != nil {
		return nil, errors.Wrap(err, "un-marshalling status response")
	}

	return &status, nil
}

func GetTickTransactionsHTTP(ctx context.Context, httpClient *http.Client, baseURL string, tickNumber uint32) (*TickTransactions, error) {

	finalURL := fmt.Sprintf(baseURL+"/v2/ticks/%d/transactions?transfers=true", tickNumber)

	tickTransactionsRequest, err := http.NewRequestWithContext(ctx, "GET", finalURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating tick transactions request")
	}

	response, err := httpClient.Do(tickTransactionsRequest)
	if err != nil {
		return nil, errors.Wrap(err, "performing tick transactions request")
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading tick transactions response")
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Status not 200! Info: %s", data))
	}

	var transactions TickTransactions

	err = json.Unmarshal(data, &transactions)
	if err != nil {
		return nil, errors.Wrap(err, "un-marshalling tick transactions response")
	}

	return &transactions, nil

}

type Status struct {
	LastProcessedTick struct {
		TickNumber uint32 `json:"tickNumber"`
		Epoch      uint32 `json:"epoch"`
	} `json:"lastProcessedTick"`
	LastProcessedTicksPerEpoch map[string]uint32 `json:"lastProcessedTicksPerEpoch"`
	SkippedTicks               []struct {
		StartTick uint32 `json:"startTick"`
		EndTick   uint32 `json:"endTick"`
	} `json:"skippedTicks"`
	ProcessedTickIntervalsPerEpoch []struct {
		Epoch     uint32 `json:"epoch"`
		Intervals []struct {
			InitialProcessedTick uint32 `json:"initialProcessedTick"`
			LastProcessedTick    uint32 `json:"lastProcessedTick"`
		} `json:"intervals"`
	}
	EmptyTicksPerEpoch map[string]uint32 `json:"emptyTicksPerEpoch"`
}

type TickTransactions struct {
	Transactions []struct {
		Transaction struct {
			SourceId     string `json:"sourceId"`
			DestId       string `json:"destId"`
			Amount       string `json:"amount"`
			TickNumber   uint32 `json:"tickNumber"`
			InputType    uint32 `json:"inputType"`
			InputSize    uint32 `json:"inputSize"`
			InputHex     string `json:"inputHex"`
			SignatureHex string `json:"signatureHex"`
			TxId         string `json:"txId"`
		} `json:"transaction"`
		Timestamp string `json:"timestamp"`
		MoneyFlew bool   `json:"moneyFlew"`
	} `json:"transactions"`
}
