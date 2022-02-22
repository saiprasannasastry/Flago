package service

import (
	"context"
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

type FlagReq struct {
	CustomerId   string `json:"customer_id,omitempty"`
	CustomerName string `json:"customer_name,omitempty"`
	Feature      string `json:"feature,omitempty"`
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

const (
	FlagExpressionTypeConstant FlagExpressionType = "constant"
	FlagExpressionTypePercent                     = "percent"
	FlagExpressionValueIn                         = "value_in"
	FlagExpressionNot                             = "not"
	FlagExpressionAllOf                           = "all_of"
	FlagExpressionAnyOf                           = "any_of"
	FlagExpressionRef                             = "ref"
)

type Manager struct {
	TaskChan    chan FlagReq
	WorkerCount int
	BusyCount   int64
	ErrorChan   chan error
	QuitChan    chan bool
	Wg          sync.WaitGroup
	RetryCount int
	RetryBackOff time.Duration
	RetryFactor int
}

func NewConcurrentManager(workerCount int,retryCount int,retryFactor int,retryBackoff time.Duration) *Manager {
	wg := sync.WaitGroup{}
	m := &Manager{
		TaskChan: make(chan FlagReq),
		WorkerCount: workerCount,
		ErrorChan: make(chan error),
		Wg: wg,
		QuitChan: make(chan bool),
		RetryCount: retryCount,
		RetryFactor: retryFactor,
		RetryBackOff: retryBackoff,
	}
	return m
}

func NewFlagoServer(redisInterface redis_pkg.PoolInterface, manager *Manager) proto.FlagoServiceServer {

	return &Flago{redisPool: redisInterface, Manager: manager}
}

func (f *Flago) CreateFlag(ctx context.Context, input *proto.CreateFlagReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), UnmarshalandStore(ctx,f, input.FlagData, input.FlagFamily.String())
}

func (f *Flago) GetFlag(ctx context.Context, input *proto.GetFlagReq) (*proto.FlagResp, error) {
	return nil, nil
}
func (f *Flago) OnFlag(ctx context.Context, input *proto.FlagReq) (*proto.FlagResp, error) {
	return nil, nil
}
func (f *Flago) OffFlag(ctx context.Context, input *proto.FlagReq) (*proto.FlagResp, error) {
	return nil, nil
}
