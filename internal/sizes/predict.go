package sizes

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/ViaQ/logerr/kverrors"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
)

var (
	metricsClient Client

	// durationPredictSecs is the value of time series in seconds from now.
	// It is passed as second parameter to predict_linear.
	durationPredictSecs = 1 * 24 * 3600 // 1d = 1 * 24h * 3600s = 86400s
	// timeoutClient is the timeout duration for prometheus client.
	timeoutClient = 10 * time.Second

	// promURL is the URL of the prometheus thanos querier
	promURL string
	// promToken is the token to connect to prometheus thanos querier.
	promToken string

	cmd []byte
)

type client struct {
	api     v1.API
	timeout time.Duration
}

// Client is the interface which contains methods for querying and extracting metrics.
type Client interface {
	LogLoggedBytesReceivedTotal(duration model.Duration) (float64, error)
}

func newClient(url, token string) (*client, error) {
	httpConfig := config.HTTPClientConfig{
		BearerToken: config.Secret(token),
		TLSConfig: config.TLSConfig{
			InsecureSkipVerify: true,
		},
	}

	rt, rtErr := config.NewRoundTripperFromConfig(httpConfig, "size-calculator-metrics")

	if rtErr != nil {
		return nil, kverrors.Wrap(rtErr, "failed creating prometheus configuration")
	}

	pc, err := api.NewClient(api.Config{
		Address:      url,
		RoundTripper: rt,
	})
	if err != nil {
		return nil, kverrors.Wrap(err, "failed creating prometheus client")
	}

	return &client{
		api:     v1.NewAPI(pc),
		timeout: timeoutClient,
	}, nil
}

func (c *client) executeScalarQuery(query string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	res, _, err := c.api.Query(ctx, query, time.Now())
	if err != nil {
		return 0.0, kverrors.Wrap(err, "failed executing query",
			"query", query)
	}

	if res.Type() == model.ValScalar {
		value := res.(*model.Scalar)
		return float64(value.Value), nil
	}

	if res.Type() == model.ValVector {
		vec := res.(model.Vector)
		if vec.Len() == 0 {
			return 0.0, nil
		}

		return float64(vec[0].Value), nil
	}

	return 0.0, kverrors.Wrap(nil, "failed to parse result for query",
		"query", query)
}

func (c *client) LogLoggedBytesReceivedTotal(duration model.Duration) (float64, error) {
	query := fmt.Sprintf(
		`sum(predict_linear(log_logged_bytes_total[%s], %d))`,
		duration,
		durationPredictSecs,
	)

	return c.executeScalarQuery(query)
}

// PredictFor takes the default duration and predicts
// the amount of logs expected in 1 day
func PredictFor(duration model.Duration) (logsCollected float64, err error) {
	// execute the bash script to get the prometheus thanos querier URL and token.
	cmd, err = exec.Command("./internal/sizes/run.sh").Output()
	if err != nil {
		return 0, kverrors.Wrap(err, "Failed to execute the script")
	}

	promURL = strings.Split(string(cmd), ",")[0]
	promToken = strings.Split(string(cmd), ",")[1]

	// Create a client to collect metrics
	metricsClient, err = newClient(promURL, promToken)
	if err != nil {
		return 0, kverrors.Wrap(err, "Failed to create metrics client")
	}

	logsCollected, err = metricsClient.LogLoggedBytesReceivedTotal(duration)
	if err != nil {
		return 0, err
	}

	return logsCollected, nil
}
