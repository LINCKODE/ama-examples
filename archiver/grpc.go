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

	// GRPC client creation

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	connection, err := grpc.NewClient(archiverGRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.Wrap(err, "creating grpc connection")
	}
	defer connection.Close()

	archiverClient := protobuff.NewArchiveServiceClient(connection)

	// Get status
	status, err := FetchStatus(ctx, archiverClient)
	if err != nil {
		return errors.Wrap(err, "fetching status")
	}

	fmt.Printf("Last processed tick: %d\n", status.LastProcessedTick.TickNumber)
	fmt.Printf("Current epoch: %d\n", status.LastProcessedTick.Epoch)

	return nil
}

func FetchStatus(ctx context.Context, archiverClient protobuff.ArchiveServiceClient) (*protobuff.GetStatusResponse, error) {

	result, err := archiverClient.GetStatus(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting archiver status")
	}

	return result, nil
}

func GetCurrentTickTransactions(ctx context.Context, archiverClient protobuff.ArchiveServiceClient) {

	result, err := archiverClient.GetTickTransactions

}
