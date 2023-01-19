package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
)

type InfluxClient struct {
	client influxdb2.Client
	api    api.WriteAPIBlocking

	chConnUpdates chan<- ConnectionUpdate
	chStop        chan bool

	opts ConnectOpts
}

type ConnectOpts struct {
	Host            string
	Port            int
	Token           string
	Org             string
	Bucket          string
	RefreshInterval time.Duration
}

type ConnectionUpdate struct {
	Err     error
	Name    string
	Version string
	Commit  string
}

func NewInfluxClient(opts ConnectOpts) (*InfluxClient, <-chan ConnectionUpdate) {
	influxURL := "http://" + opts.Host + ":" + strconv.Itoa(opts.Port)
	client := influxdb2.NewClientWithOptions(
		influxURL,
		opts.Token,
		influxdb2.DefaultOptions().
			SetApplicationName(constants.APP_ID).
			SetBatchSize(constants.INFLUX_BATCH_SIZE).
			SetPrecision(1*time.Second),
	)

	api := client.WriteAPIBlocking(opts.Org, opts.Bucket)
	api.EnableBatching()

	chConnUpdates := make(chan ConnectionUpdate, 1)
	return &InfluxClient{
		client:        client,
		api:           api,
		chConnUpdates: chConnUpdates,
		opts:          opts,
	}, chConnUpdates
}

func (c *InfluxClient) Start() {
	if c.chStop != nil {
		panic("trying to start InfluxClient while it's still running")
	}
	log.Println("starting InfluxClient")
	c.chStop = make(chan bool, 1)
	go c.work()
}

func (c *InfluxClient) Close() {
	if c.chStop == nil {
		panic("trying to Stop InfluxClient that's not running")
	}
	log.Println("stopping InfluxClient")
	close(c.chStop)
	close(c.chConnUpdates)
	c.chStop = nil
}

func (c *InfluxClient) work() {
	log.Printf("starting worker with refresh interval=%fs", c.opts.RefreshInterval.Seconds())
	retryInterval := 5 * time.Second
	maxInterval := time.Duration(math.Max(float64(5*time.Minute), float64(c.opts.RefreshInterval)))
	refreshTicker := time.NewTicker(c.opts.RefreshInterval)
	connected := false

	for {
		var update *ConnectionUpdate
		if connected {
			log.Println("pinging InfluxDB")
			// if no error on ping, everything stays the same, no need for update
			if err := c.ping(c.opts.RefreshInterval); err != nil {
				update = &ConnectionUpdate{Err: err}
			}
		} else {
			log.Println("attempting connect")
			update = c.connect(retryInterval)
		}

		// if update == nil, there is no change
		if update != nil {
			c.chConnUpdates <- *update
			if update.Err != nil {
				if !connected && retryInterval < maxInterval {
					retryInterval *= 2
					refreshTicker.Reset(retryInterval)
				}
				log.Printf("unsuccessful, retrying in %fs (%v)", retryInterval.Seconds(), update.Err)
				connected = false
			} else {
				log.Printf("connection successful, checking again in %fs", c.opts.RefreshInterval.Seconds())
				refreshTicker.Reset(c.opts.RefreshInterval)
				retryInterval = c.opts.RefreshInterval
				connected = true
			}
		} else {
			log.Println("ping successful")
		}

		select {
		case <-refreshTicker.C:
			continue
		case <-c.chStop:
			log.Println("received stop signal, stopping InfluxDB work")
			return
		}
	}
}

func (c *InfluxClient) connect(timeout time.Duration) (update *ConnectionUpdate) {
	update = new(ConnectionUpdate)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	health, healthErr := c.client.Health(ctx)
	if healthErr != nil {
		update.Err = fmt.Errorf("could not validate InfluxDB server health: %w", healthErr)
		return
	} else if health.Status != domain.HealthCheckStatusPass {
		update.Err = fmt.Errorf("InfluxDB server unhealthy (%s), please check your server logs", health.Status)
		return
	}

	if err := c.validateCanWrite(c.opts.Org, c.opts.Bucket); err != nil {
		update.Err = err
		return
	}
	update.Name = health.Name
	update.Version = *health.Version
	update.Commit = *health.Commit
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

func (c *InfluxClient) ping(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ok, err := c.client.Ping(ctx)
	if err != nil {
		return err
	} else if !ok {
		return errors.New("host unreachable")
	} else {
		return nil
	}
}
