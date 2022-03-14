package service

import (
	"context"
	"github.com/apex/log"
	"sync"
	"time"

	//"encoding/json"
	redis_pkg "github.com/segmentio/feature-flag/pkg/redis"

	"github.com/segmentio/feature-flag/pkg/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Flago struct {
	redisPool redis_pkg.PoolInterface
	Manager   *Manager
	proto.UnimplementedFlagoServiceServer
}
type FlagExpressionType string

type Manager struct {
	TaskChan     chan FlagReq
	WorkerCount  int
	BusyCount    int64
	ErrorChan    chan error
	QuitChan     chan bool
	Wg           *sync.WaitGroup
	RetryCount   int
	RetryBackOff time.Duration
	RetryFactor  int
}

type FlagExpression struct {
	Type     FlagExpressionType
	Constant bool
	Percent  float64
	ValueIn  []interface{}
	Not      *FlagExpression
	AllOf    []FlagExpression
	AnyOf    []FlagExpression
	Ref      string
}

type FlagReq struct {
	CustomerId   string `json:"customer_id,omitempty"`
	CustomerName string `json:"customer_name,omitempty"`
	Feature      string `json:"feature,omitempty"`
}

func NewConcurrentManager(workerCount int, retryCount int, retryFactor int, retryBackoff time.Duration) *Manager {
	wg := new(sync.WaitGroup)
	m := &Manager{
		WorkerCount:  workerCount,
		Wg:           wg,
		QuitChan:     make(chan bool),
		RetryCount:   retryCount,
		RetryFactor:  retryFactor,
		RetryBackOff: retryBackoff,
	}
	return m
}

func NewFlagoServer(redisInterface redis_pkg.PoolInterface, manager *Manager) proto.FlagoServiceServer {

	return &Flago{redisPool: redisInterface, Manager: manager}
}

func (f *Flago) CreateFlag(ctx context.Context, input *proto.CreateFlagReq) (*emptypb.Empty, error) {

	return new(emptypb.Empty), UnmarshalandStore(ctx, f, input.FlagData, input.FlagFamily.String())
}

func (f *Flago) GetFlag(ctx context.Context, input *proto.FlagReq) (*proto.FlagResp, error) {
	flagDetails, err := f.redisPool.GetFlagForCustomer(input.CustomerName+"::"+input.CustomerId, input.Feature)
	if err != nil {
		log.WithError(err).Errorf("failed to get flag details for customer %v", input.CustomerName)
		return nil, err
	}
	return &proto.FlagResp{Enabled: flagDetails}, nil
}
func (f *Flago) GetFlags(ctx context.Context, input *proto.FlagReq) (*proto.GetFlagResp, error) {
	flagDetails, err := f.redisPool.GetAllCustomers(input.CustomerName + "::" + input.CustomerId)
	if err != nil {
		log.WithError(err).Errorf("failed to get flag details for customer %v", input.CustomerName)
		return nil, err
	}
	log.Infof("%v", flagDetails)
	return &proto.GetFlagResp{Flags: flagDetails}, nil

}

func (f *Flago) OnFlag(ctx context.Context, input *proto.FlagReq) (*proto.FlagResp, error) {
	return nil, nil
}
func (f *Flago) OffFlag(ctx context.Context, input *proto.FlagReq) (*proto.FlagResp, error) {
	return nil, nil
}

//work spins up concurrent workers to add data to redis
func (m *Manager) Work(ctx context.Context, fn func(string, string, string) error, workerNumber int) {
	log.Infof("spawnning worker %v", workerNumber)
	defer m.Wg.Done()
	defer log.Info("done working")
	for {
		select {
		case t, ok := <-m.TaskChan:
			if ok {
				err := fn(t.CustomerName, t.CustomerId, t.Feature)
				if err != nil {
					log.Infof("pushing error from %v", workerNumber)
					m.ErrorChan <- err
				}
			} else {
				return
			}
		case <-ctx.Done():
			log.Infof("closing channel %v", ctx.Err())
			return
		}
	}
}
