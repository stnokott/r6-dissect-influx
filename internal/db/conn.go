// Package db contains functions for communication with an InfluxDB host.
package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"github.com/stnokott/r6-dissect-influx/internal/game"
)

// InfluxClient is used to send R6 round data to an InfluxDB host.
type InfluxClient struct {
	URL string

	client influxdb2.Client
	// use blocking writes since all calls from frontend are already asynchronous
	api api.WriteAPIBlocking

	loopCancel context.CancelFunc

	opts ConnectOpts
}

// ConnectOpts allows defining options for a new client.
type ConnectOpts struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

// ConnectionDetails contains details about the connection to an InfluxDB host.
type ConnectionDetails struct {
	Name    string
	Version string
	Commit  string
}

// NewInfluxClient creates a new client with the provided options.
func NewInfluxClient(opts ConnectOpts) *InfluxClient {
	client := influxdb2.NewClientWithOptions(
		opts.URL,
		opts.Token,
		influxdb2.DefaultOptions().
			SetApplicationName(constants.APP_ID).
			SetPrecision(1*time.Second),
	)

	api := client.WriteAPIBlocking(opts.Org, opts.Bucket)

	return &InfluxClient{
		URL:    opts.URL,
		client: client,
		api:    api,
		opts:   opts,
	}
}

// ValidateConn validates that the InfluxDB host can be reached, is healthy and that the client has write permissions.
func (c *InfluxClient) ValidateConn(timeout time.Duration) (details *ConnectionDetails, err error) {
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
	details = &ConnectionDetails{
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

var (
	roundBuffer      []game.RoundInfo = []game.RoundInfo{}
	roundBufferMutex sync.Mutex
)

// AddRound adds a game round to the buffer to be pushed to the InfluxDB when possible.
func AddRound(r game.RoundInfo) {
	roundBufferMutex.Lock()
	roundBuffer = append(roundBuffer, r)
	roundBufferMutex.Unlock()
}

// InfluxEvent contains data sent from the asynchronous loop started by LoopAsync.
// It contains the MatchID and RoundIndex for the related event and an Err field indicating success if nil.
type InfluxEvent struct {
	Err        error
	MatchID    string
	RoundIndex int
}

// LoopAsync starts a goroutine which pushes the data from the queue iteratively (FIFO) to the InfluxDB.
// It returns a channel used for providing event information.
func (c *InfluxClient) LoopAsync() <-chan InfluxEvent {
	if c.loopCancel != nil {
		log.Println("WARNING: loop already running, cancelling previous loop")
		c.loopCancel()
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	ticker := time.NewTicker(5 * time.Second)
	c.loopCancel = cancelFunc
	eventChan := make(chan InfluxEvent)

	go func() {
		defer ticker.Stop()
		defer close(eventChan)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				roundBufferMutex.Lock()
				if len(roundBuffer) > 0 {
					var r game.RoundInfo
					r, roundBuffer = roundBuffer[0], roundBuffer[1:]
					eventChan <- InfluxEvent{
						MatchID:    r.MatchID,
						RoundIndex: r.RoundIndex,
						Err:        c.writeRound(r),
					}
				}
				roundBufferMutex.Unlock()
			}
		}

	}()
	return eventChan
}

const (
	measurement string = "r6-dissect"

	tagAppVersion string = "appversion"
	tagMatchID    string = "matchid"
	tagSeasonSlug string = "season"
	tagMatchType  string = "matchtype"
	tagGameMode   string = "gamemode"
	tagMap        string = "map"
	tagSite       string = "site"

	fieldWon string = "won"
)

func (c *InfluxClient) writeRound(r game.RoundInfo) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = c.api.WritePoint(ctx, write.NewPoint(
		measurement,
		map[string]string{
			tagAppVersion: constants.Version,
			tagMatchID:    r.MatchID,
			tagSeasonSlug: r.SeasonSlug,
			tagMatchType:  r.MatchType,
			tagGameMode:   r.GameMode,
			tagMap:        r.MapName,
			tagSite:       r.Site,
		},
		map[string]interface{}{
			fieldWon: r.Won,
		},
		r.Time,
	))
	if err != nil {
		return
	}
	err = c.api.Flush(ctx)
	return
}

// Close calls the client's underlying Close function, finishing all asynchronous writes.
func (c *InfluxClient) Close() {
	if c.loopCancel != nil {
		log.Println("stopping loop")
		c.loopCancel()
	}
	log.Println("closing connection")
	c.client.Close()
}
