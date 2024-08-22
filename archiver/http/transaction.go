package http

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func FetchTransactionHTTP(transactionID string) error {

	archiverAddress := fmt.Sprintf("https://testapi.qubic.org/v1/transactions/%s", transactionID)

	httpClient := http.DefaultClient

	request, err := http.NewRequest(http.MethodGet, archiverAddress, nil)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "executing request")
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "reading response body")
	}

	fmt.Printf("%s\n", string(data))

	return nil
}
