package archiver;

import com.google.gson.Gson;
import org.qubic.archiver.proto.Archive;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.Map;

public class HttpExample {

    private final String baseUrl;
    private final Gson gson;

    public HttpExample(String baseUrl) {
        this.baseUrl = baseUrl;
        this.gson = new Gson();

        ArchiverStatus status = fetchArchiverStatus();

        int latestTick = status.getLastProcessedTick().getTickNumber();
        int epoch = status.getLastProcessedTick().getEpoch();

        System.out.println("Latest tick: " + latestTick);
        System.out.println("Current epoch: " + epoch);

        Transaction[] latestTickTransactions = fetchTickTransactions(latestTick);

        for (Transaction transaction : latestTickTransactions) {
            System.out.println(" Transaction ID: " + transaction.getTxId());
            System.out.println(" Transferred amount: " + transaction.getAmount());
            System.out.println(" From: " + transaction.getSourceId());
            System.out.println(" To: " + transaction.getDestId());
            System.out.println("-------------------------------");
        }

    }

    private ArchiverStatus fetchArchiverStatus() {

        try {
            String response = performRequest(this.baseUrl + "/v1/status", "GET");
            return gson.fromJson(response, ArchiverStatus.class);

        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private Transaction[] fetchTickTransactions(int tickNumber) {
        try {
            String response = performRequest(this.baseUrl + "/v1/ticks/" + tickNumber + "/transactions", "GET");
            TransactionsResponse transactionsResponse = gson.fromJson(response, TransactionsResponse.class);

            return transactionsResponse.getTransactions();
        }
        catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private String performRequest(String address, String requestType) throws Exception {
        URL url = new URL(address);
        HttpURLConnection connection = (HttpURLConnection) url.openConnection();
        connection.setRequestMethod(requestType);

        int statusCode = connection.getResponseCode();
        if (statusCode != 200) {
            throw new Exception("Non 200 status code: " + statusCode);
        }


        // Read response
        StringBuilder result = new StringBuilder();
        try (BufferedReader reader = new BufferedReader(
                new InputStreamReader(connection.getInputStream()))) {
            for (String line; (line = reader.readLine()) != null; ) {
                result.append(line);
            }
        }
        connection.disconnect();
        return result.toString();

    }

    static class TransactionsResponse {
        private Transaction[] transactions;

        public TransactionsResponse(Transaction[] transactions) {
            this.transactions = transactions;
        }

        public Transaction[] getTransactions() {
            return transactions;
        }

        public void setTransactions(Transaction[] transactions) {
            this.transactions = transactions;
        }
    }

    static class Transaction {
        private String sourceId;
        private String destId;
        private int amount;
        private int tickNumber;
        private int inputType;
        private int inputSize;
        private String inputHex;
        private String signatureHex;
        private String txId;

        public Transaction(String sourceId, String destId, int amount, int tickNumber, int inputType, int inputSize, String inputHex, String signatureHex, String txId) {
            this.sourceId = sourceId;
            this.destId = destId;
            this.amount = amount;
            this.tickNumber = tickNumber;
            this.inputType = inputType;
            this.inputSize = inputSize;
            this.inputHex = inputHex;
            this.signatureHex = signatureHex;
            this.txId = txId;
        }

        public String getSourceId() {
            return sourceId;
        }

        public void setSourceId(String sourceId) {
            this.sourceId = sourceId;
        }

        public String getDestId() {
            return destId;
        }

        public void setDestId(String destId) {
            this.destId = destId;
        }

        public int getAmount() {
            return amount;
        }

        public void setAmount(int amount) {
            this.amount = amount;
        }

        public int getTickNumber() {
            return tickNumber;
        }

        public void setTickNumber(int tickNumber) {
            this.tickNumber = tickNumber;
        }

        public int getInputType() {
            return inputType;
        }

        public void setInputType(int inputType) {
            this.inputType = inputType;
        }

        public int getInputSize() {
            return inputSize;
        }

        public void setInputSize(int inputSize) {
            this.inputSize = inputSize;
        }

        public String getInputHex() {
            return inputHex;
        }

        public void setInputHex(String inputHex) {
            this.inputHex = inputHex;
        }

        public String getSignatureHex() {
            return signatureHex;
        }

        public void setSignatureHex(String signatureHex) {
            this.signatureHex = signatureHex;
        }

        public String getTxId() {
            return txId;
        }

        public void setTxId(String txId) {
            this.txId = txId;
        }
    }

    static class LastProcessedTick {
        private int tickNumber;
        private int epoch;

        public LastProcessedTick(int tickNumber, int epoch) {
            this.tickNumber = tickNumber;
            this.epoch = epoch;
        }

        public int getTickNumber() {
            return tickNumber;
        }

        public void setTickNumber(int tickNumber) {
            this.tickNumber = tickNumber;
        }

        public int getEpoch() {
            return epoch;
        }

        public void setEpoch(int epoch) {
            this.epoch = epoch;
        }
    }

    static class SkippedTicksInterval {
        private int startTick;
        private int endTick;

        public SkippedTicksInterval(int startTick, int endTick) {
            this.startTick = startTick;
            this.endTick = endTick;
        }

        public int getStartTick() {
            return startTick;
        }

        public void setStartTick(int startTick) {
            this.startTick = startTick;
        }

        public int getEndTick() {
            return endTick;
        }

        public void setEndTick(int endTick) {
            this.endTick = endTick;
        }
    }

    static class ProcessedTicksInterval {
        private int initialProcessedTick;
        private int lastProcessedTick;

        public ProcessedTicksInterval(int initialProcessedTick, int lastProcessedTick) {
            this.initialProcessedTick = initialProcessedTick;
            this.lastProcessedTick = lastProcessedTick;
        }

        public int getInitialProcessedTick() {
            return initialProcessedTick;
        }

        public void setInitialProcessedTick(int initialProcessedTick) {
            this.initialProcessedTick = initialProcessedTick;
        }

        public int getLastProcessedTick() {
            return lastProcessedTick;
        }

        public void setLastProcessedTick(int lastProcessedTick) {
            this.lastProcessedTick = lastProcessedTick;
        }
    }

    static class ProcessedTicksIntervalsPerEpoch {
        private int epoch;
        private ProcessedTicksInterval[] intervals;

        public ProcessedTicksIntervalsPerEpoch(int epoch, ProcessedTicksInterval[] intervals) {
            this.epoch = epoch;
            this.intervals = intervals;
        }

        public int getEpoch() {
            return epoch;
        }

        public void setEpoch(int epoch) {
            this.epoch = epoch;
        }

        public ProcessedTicksInterval[] getIntervals() {
            return intervals;
        }

        public void setIntervals(ProcessedTicksInterval[] intervals) {
            this.intervals = intervals;
        }
    }

    static class ArchiverStatus {

        private LastProcessedTick lastProcessedTick;
        private Map<String, Integer> lastProcessedTicksPerEpoch;
        private SkippedTicksInterval[] skippedTicks;
        private ProcessedTicksIntervalsPerEpoch[] processedTickIntervalsPerEpoch;
        private Map<String, Integer> emptyTicksPerEpoch;

        public ArchiverStatus(LastProcessedTick lastProcessedTick, Map<String, Integer> lastProcessedTicksPerEpoch, SkippedTicksInterval[] skippedTicks, ProcessedTicksIntervalsPerEpoch[] processedTickIntervalsPerEpoch, Map<String, Integer> emptyTicksPerEpoch) {
            this.lastProcessedTick = lastProcessedTick;
            this.lastProcessedTicksPerEpoch = lastProcessedTicksPerEpoch;
            this.skippedTicks = skippedTicks;
            this.processedTickIntervalsPerEpoch = processedTickIntervalsPerEpoch;
            this.emptyTicksPerEpoch = emptyTicksPerEpoch;
        }

        public LastProcessedTick getLastProcessedTick() {
            return lastProcessedTick;
        }

        public void setLastProcessedTick(LastProcessedTick lastProcessedTick) {
            this.lastProcessedTick = lastProcessedTick;
        }

        public Map<String, Integer> getLastProcessedTicksPerEpoch() {
            return lastProcessedTicksPerEpoch;
        }

        public void setLastProcessedTicksPerEpoch(Map<String, Integer> lastProcessedTicksPerEpoch) {
            this.lastProcessedTicksPerEpoch = lastProcessedTicksPerEpoch;
        }

        public SkippedTicksInterval[] getSkippedTicks() {
            return skippedTicks;
        }

        public void setSkippedTicks(SkippedTicksInterval[] skippedTicks) {
            this.skippedTicks = skippedTicks;
        }

        public ProcessedTicksIntervalsPerEpoch[] getProcessedTickIntervalsPerEpoch() {
            return processedTickIntervalsPerEpoch;
        }

        public void setProcessedTickIntervalsPerEpoch(ProcessedTicksIntervalsPerEpoch[] processedTickIntervalsPerEpoch) {
            this.processedTickIntervalsPerEpoch = processedTickIntervalsPerEpoch;
        }

        public Map<String, Integer> getEmptyTicksPerEpoch() {
            return emptyTicksPerEpoch;
        }

        public void setEmptyTicksPerEpoch(Map<String, Integer> emptyTicksPerEpoch) {
            this.emptyTicksPerEpoch = emptyTicksPerEpoch;
        }
    }

}
