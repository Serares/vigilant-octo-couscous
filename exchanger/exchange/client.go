package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type ExchangeRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type ExchangeConvertResponse struct {
	Success bool    `json:"success"`
	Result  float64 `json:"result"`
}

type ExchangeClient struct {
	Log    *slog.Logger
	Client *http.Client
}

func NewExchangeClient(log *slog.Logger) *ExchangeClient {
	return &ExchangeClient{
		Log: log,
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (c *ExchangeClient) GetConvertAmount(data ExchangeRequest) (ExchangeConvertResponse, error) {
	accessKey := os.Getenv("EXCHANGE_API_KEY")
	baseUrl := os.Getenv("EXCHANGE_API_URL")
	url := fmt.Sprintf(
		"%sconvert?access_key=%s&from=%s&to=%s&amount=%f&format=1",
		baseUrl,
		accessKey,
		data.From,
		data.To,
		data.Amount,
	)
	r, err := c.Client.Get(url)
	if err != nil {
		c.Log.Error("Error trying to query the api", "error", err)
		return ExchangeConvertResponse{
			Success: false,
			Result:  0,
		}, err
	}
	defer r.Body.Close()
	var resp ExchangeConvertResponse

	if r.StatusCode == 200 {
		if err != nil {
			c.Log.Error("Error trying to close the body", "error", err)
			return ExchangeConvertResponse{
				Success: false,
				Result:  0,
			}, err
		}
	}
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		c.Log.Error("Error trying to read the body", "error", err)
		return ExchangeConvertResponse{
			Success: false,
			Result:  0,
		}, err
	}

	//c.Log.Info("Response from the api", "response", string(resp))
	err = json.Unmarshal(msg, &resp)
	if err != nil {
		c.Log.Error("Error trying to unmarshal the body", "error", err)
		return ExchangeConvertResponse{
			Success: false,
			Result:  0,
		}, err
	}

	//c.Log.Info("Response from the api", "response", resp)
	return resp, nil
}
