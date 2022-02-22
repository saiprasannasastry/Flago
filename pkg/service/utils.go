package service

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/segmentio/feature-flag/pkg/proto"
)

func UnmarshalandStore(ctx context.Context,f *Flago, flagData []byte, flagFamily string) error {
	switch flagFamily {
	case proto.CreateFlagReq_UNIVERSAL_DISABLED.String():
		var feature string
		err := json.Unmarshal(flagData, &feature)
		if err != nil {
			log.WithError(err).Errorf("Invalid data passed for flag type %v", proto.CreateFlagReq_UNIVERSAL_DISABLED.String())
			return err
		}
		return f.redisPool.DisableAllCustomers(feature)

	case proto.CreateFlagReq_UNIVERSAL_ENABLED.String():
		var feature string
		err := json.Unmarshal(flagData, &feature)
		if err != nil {
			log.WithError(err).Errorf("Invalid data passed for flag type %v", proto.CreateFlagReq_UNIVERSAL_DISABLED.String())
			return err
		}
		return f.redisPool.EnableAllCustomers(feature)

	case proto.CreateFlagReq_SET_OF.String():
		var requests []FlagReq
		err := json.Unmarshal(flagData, &requests)
		if err != nil {
			log.WithError(err).Errorf("Invalid data passed for flag type %v", proto.CreateFlagReq_SET_OF.String())
			return err
		}
	//	newCtx := context.Background()
		//if at all we get multiple customers for one feature, we can run all of them in background
		for i := 0; i < f.Manager.WorkerCount; i++ {
			f.Manager.Wg.Add(1)
			go f.Manager.work(ctx,  f.redisPool.AddToSetOfcustomers,i)
		}

		for _, request := range requests {
			f.Manager.TaskChan <- request
		}

		go func (m *Manager){
			for {
				if len(m.TaskChan)==0{
					log.Debug("Task channel is empty")
					m.QuitChan<-true
				}
			}
		}(f.Manager)

		//var errors error
		//for {
		//	select{
		//	case err:= <- f.Manager.ErrorChan:
		//		errors= multierror.Append(errors,err)
		//	case <- f.Manager.QuitChan:
		//		return errors
		//	case <- ctx.Done():
		//		return errors
		//	default:
		//		log.Info("in default")
		//	}
		//}

		f.Manager.Wg.Wait()
		return nil

	case proto.CreateFlagReq_REFERENCE_TYPE.String():
		type ref struct {
			ref     string `json:"ref,omitempty"`
			feature string `json:"feature,omitempty"`
		}
		var data ref
		err := json.Unmarshal(flagData, &data)
		if err != nil {
			log.WithError(err).Errorf("Invalid data passed for flag type %v", proto.CreateFlagReq_SET_OF.String())
			return err
		}
		log.Infof("%v",data)
	}
	return nil

}

//work spins up concurrent workers to add data to redis
func (m *Manager) work(ctx context.Context, fn func( string, string, string) error, workerNumber int) {
	log.Infof("spawnning worker %v",workerNumber)
	defer m.Wg.Done()
	defer log.Info("returning ")
	for {
		select {
		case t := <-m.TaskChan:
			fn(t.CustomerName, t.CustomerId, t.Feature)
		case <-ctx.Done():
			log.Infof("closing channel %v",ctx.Err())
			return
		}
	}
}
