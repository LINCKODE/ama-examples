package archiver

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/qubic/go-archiver/protobuff"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const archiverGRPCAddress = "213.170.135.5:8003"

func RunGRPCExample() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// GRPC client creation
	connection, err := grpc.NewClient(archiverGRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.Wrap(err, "creating grpc connection")
	}
	defer connection.Close()

	archiverClient := protobuff.NewArchiveServiceClient(connection)

	// Get status
	status, err := FetchStatusGRPC(ctx, archiverClient)
	if err != nil {
		return errors.Wrap(err, "fetching status")
	}

	fmt.Printf("Last processed tick: %d\n", status.LastProcessedTick.TickNumber)
	fmt.Printf("Current epoch: %d\n", status.LastProcessedTick.Epoch)

	// Get transactions in the latest tick

	transactions, err := GetTickTransactionsGRPC(ctx, archiverClient, status.LastProcessedTick.TickNumber, true)
	if err != nil {
		return errors.Wrap(err, "fetching tick transactions")
	}

	fmt.Printf("Found %d transactions\n", len(transactions.Transactions))

	for _, transaction := range transactions.Transactions {
		fmt.Printf("  Transaction ID: %s\n", transaction.Transaction.TxId)
		fmt.Printf("  Transfered amount: %d\n", transaction.Transaction.Amount)
		fmt.Printf("  From: %s\n", transaction.Transaction.SourceId)
		fmt.Printf("  To: %s\n", transaction.Transaction.DestId)
		println("-------------------------------")

	}

	return nil
}

func FetchStatusGRPC(ctx context.Context, archiverClient protobuff.ArchiveServiceClient) (*protobuff.GetStatusResponse, error) {

	result, err := archiverClient.GetStatus(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "performing request")
	}

	return result, nil
}

func GetTickTransactionsGRPC(ctx context.Context, archiverClient protobuff.ArchiveServiceClient, tickNumber uint32, transfersOnly bool) (*protobuff.GetTickTransactionsResponseV2, error) {

	request := protobuff.GetTickTransactionsRequestV2{
		TickNumber: tickNumber,
		// Do not filter, we want to see all transactions
		Approved:  false,
		Transfers: transfersOnly,
	}

	result, err := archiverClient.GetTickTransactionsV2(ctx, &request)
	if err != nil {
		return nil, errors.Wrap(err, "performing request")
	}

	return result, nil

}
