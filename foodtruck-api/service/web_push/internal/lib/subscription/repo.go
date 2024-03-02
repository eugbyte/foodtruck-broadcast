package subscription

import (
	"errors"
	debug "foodtruck/pkg/logger"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/samber/lo"
)

var logger = debug.Logger

// Subscription Repository to access Subscription DB
type subRepo struct {
	db        *dynamo.DB
	tableName string
	region    string
	endpoint  *string
}

// Create new subscription repository service to access Subscription DB.
// `endpoint` is optional argument.
func New(tableName string, region string, endpoint *string) *subRepo {
	return &subRepo{
		tableName: tableName,
		region:    region,
		endpoint:  endpoint,
	}
}

func (s *subRepo) Open() {
	sess := session.Must(session.NewSession())
	var awsConfig = aws.NewConfig().
		WithRegion(s.region)

	if s.endpoint != nil && *s.endpoint != "" {
		awsConfig = awsConfig.WithEndpoint(*s.endpoint)
	}

	s.db = dynamo.New(sess, awsConfig)
}

// Get all users within a location.
func (s *subRepo) GetAll(geohash string) ([]Subscription, error) {
	var subscriptions = make([]Subscription, 0)
	if s.db == nil {
		return subscriptions, errors.New("db not initialized. Initialize with Open()")
	}

	table := s.db.Table(s.tableName)
	err := table.Get("geohash", geohash).
		Index("geohashIndex").
		All(&subscriptions)
	return subscriptions, err
}

// Get all users within a location before a certain time.
func (s *subRepo) GetAllBefore(geohash string, before time.Time) ([]Subscription, error) {
	var subscriptions = make([]Subscription, 0)
	if s.db == nil {
		return subscriptions, errors.New("db not initialized. Initialize with Open()")
	}

	table := s.db.Table(s.tableName)
	err := table.Get("geohash", geohash).
		Index("geohashIndex").
		Range("lastSend", dynamo.Less, before.Unix()).
		All(&subscriptions)
	return subscriptions, err
}

func (s *subRepo) Get(endpointID string) (Subscription, error) {
	var subscription Subscription
	table := s.db.Table(s.tableName)
	err := table.Get("endpoint", endpointID).One(&subscription)
	return subscription, err
}

func (s *subRepo) Put(sub Subscription) error {
	if s.db == nil {
		return errors.New("db not initialized. Initialize with Open()")
	}
	table := s.db.Table(s.tableName)
	return table.Put(sub).Run()
}

func (s *subRepo) BatchPut(subs []Subscription) error {
	if s.db == nil {
		return errors.New("db not initialized. Initialize with Open()")
	}

	table := s.db.Table(s.tableName)

	// aggregate error
	var errs []string = make([]string, 0)

	for i := 0; i < len(subs); i += 25 {
		var subs []Subscription = lo.Slice(subs, i, i+25)
		var items []any = make([]any, 0)
		for _, sub := range subs {
			items = append(items, sub)
		}

		if wroteCount, err := table.Batch().Write().Put(items...).Run(); err != nil {
			errs = append(errs, err.Error())
		} else {
			logger.Info("wroteCount: ", wroteCount)
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ". "))
	}

	return nil
}
