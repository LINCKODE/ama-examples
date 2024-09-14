package archiver;

import com.google.protobuf.Empty;
import io.grpc.Grpc;
import io.grpc.InsecureChannelCredentials;
import io.grpc.ManagedChannel;
import org.qubic.archiver.proto.Archive;
import org.qubic.archiver.proto.ArchiveServiceGrpc;

public class GrpcExample {

    private final ManagedChannel channel;
    private final ArchiveServiceGrpc.ArchiveServiceBlockingStub blockingStub;

    public GrpcExample(String host, int port) {
        // Create connection
        this.channel = createChannel(host, port);
        this.blockingStub = ArchiveServiceGrpc.newBlockingStub(channel);

        // Get archiver status
        Archive.GetStatusResponse status = blockingStub.getStatus(Empty.getDefaultInstance());

        // Get latest processed tick and current epoch from status
        int latestTick = status.getLastProcessedTick().getTickNumber();
        int currentEpoch = status.getLastProcessedTick().getEpoch();

        System.out.println("Latest tick: " + latestTick);
        System.out.println("Current epoch: " + currentEpoch);

        // Create tick transactions request
        Archive.GetTickTransactionsRequest tickTransactionsRequest  = Archive.GetTickTransactionsRequest.newBuilder()
                .setTickNumber(latestTick)
                .build();

        // Perform request and get list of transactions for the latest tick
        Archive.GetTickTransactionsResponse tickTransactionsResponse = blockingStub.getTickTransactions(tickTransactionsRequest);

        for (Archive.Transaction transaction : tickTransactionsResponse.getTransactionsList()) {
            System.out.println(" Transaction ID: " + transaction.getTxId());
            System.out.println(" Transferred amount: " + transaction.getAmount());
            System.out.println(" From: " + transaction.getSourceId());
            System.out.println(" To: " + transaction.getDestId());
            System.out.println("-------------------------------");

        }

        // Close the connection
        this.channel.shutdown();

    }

    private ManagedChannel createChannel(String host, int port) {
        return Grpc.newChannelBuilder(host + ":" + port, InsecureChannelCredentials.create()).build();
    }


}
