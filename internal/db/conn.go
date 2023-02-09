package db

import (
	"context"
	"fmt"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
)

type InfluxClient struct {
	URL string

	client influxdb2.Client
	api    api.WriteAPIBlocking

	opts ConnectOpts
}

type ConnectOpts struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

type ConnectionDetails struct {
	Name    string
	Version string
	Commit  string
}

func NewInfluxClient(opts ConnectOpts) *InfluxClient {
	client := influxdb2.NewClientWithOptions(
		opts.URL,
		opts.Token,
		influxdb2.DefaultOptions().
			SetApplicationName(constants.APP_ID).
			SetBatchSize(constants.INFLUX_BATCH_SIZE).
			SetPrecision(1*time.Second),
	)

	api := client.WriteAPIBlocking(opts.Org, opts.Bucket)
	api.EnableBatching()

	return &InfluxClient{
		URL:    opts.URL,
		client: client,
		api:    api,
		opts:   opts,
	}
}

func (c *InfluxClient) Close() {
	log.Println("closing connection")
	c.client.Close()
}

func (c *InfluxClient) ValidateConn(timeout time.Duration) (update *ConnectionDetails, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	health, healthErr := c.client.Health(ctx)
	if healthErr != nil {
		err = fmt.Errorf("could not validate InfluxDB server health: %w", healthErr)
		return
	} else if health.Status != domain.HealthCheckStatusPass {
		err = fmt.Errorf("InfluxDB server unhealthy (%s), please check your server logs", health.Status)
		return
	}

	if err = c.validateCanWrite(c.opts.Org, c.opts.Bucket); err != nil {
		return
	}
	update = &ConnectionDetails{
		Name:    health.Name,
		Version: *health.Version,
		Commit:  *health.Commit,
	}
	return
}

const (
	testMeasurement = constants.APP_ID + "_test_measurement"
	testField       = constants.APP_ID + "_test_field"
)

func (c *InfluxClient) validateCanWrite(org string, bucket string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testTime := time.Now()
	if err = c.api.WritePoint(
		ctx,
		write.NewPoint(
			testMeasurement,
			nil,
			map[string]interface{}{
				testField: 0,
			},
			testTime,
		),
	); err != nil {
		return
	}
	if err = c.api.Flush(ctx); err != nil {
		return
	}

	deleteAPI := c.client.DeleteAPI()
	err = deleteAPI.DeleteWithName(
		ctx,
		org,
		bucket,
		testTime.Add(-30*time.Second),
		testTime.Add(30*time.Second),
		`_measurement="`+testMeasurement+`"`,
	)
	return
}
