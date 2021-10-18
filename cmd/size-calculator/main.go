package main

import (
	"fmt"
	"math"
	"os"

	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"

	"github.com/ViaQ/logerr/log"
	"github.com/ViaQ/loki-operator/internal/sizes"
	"github.com/prometheus/common/model"
)

const (
	// defaultDuration is the default time duration to consider for metric scraping
	defaultDuration string = "1h"
	// range1xSmall defines the range (in GB)
	// of t-shirt size 1x.small i.e., 0 <= 1x.small <= 500
	range1xSmall int = 500
)

func init() {
	log.Init("size-calculator")
}

func main() {
	duration, parseErr := model.ParseDuration(defaultDuration)
	if parseErr != nil {
		log.Error(parseErr, "failed to parse duration")
		os.Exit(1)
	}

	logsCollected, err := sizes.PredictFor(duration)
	if err != nil {
		log.Error(err, "Failed to collect metrics data")
		os.Exit(1)
	}

	logsCollectedInGB := int(math.Ceil(logsCollected / math.Pow(1024, 3)))
	log.Info(fmt.Sprintf("Amount of logs expected in 24 hours is %f Bytes or %dGB", logsCollected, logsCollectedInGB))

	if logsCollectedInGB <= range1xSmall {
		log.Info(fmt.Sprintf("Recommended t-shirt size for %dGB is %s", logsCollectedInGB, lokiv1beta1.SizeOneXSmall))
	} else {
		log.Info(fmt.Sprintf("Recommended t-shirt size for %dGB is %s", logsCollectedInGB, lokiv1beta1.SizeOneXMedium))
	}
}
