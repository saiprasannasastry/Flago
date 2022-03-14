package service

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/hashicorp/go-multierror"
	"github.com/segmentio/feature-flag/pkg/proto"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type PercentageReq struct {
	Ref        string `json:"ref,omitempty"`
	Feature    string `json:"feature,omitempty"`
	Percentage string `json:"percentage,omitempty"`
}

func UnmarshalandStore(ctx context.Context, f *Flago, flagData []byte, flagFamily string) error {

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
			log.WithError(err).Errorf("Invalid data passed for flag type %v", proto.CreateFlagReq_UNIVERSAL_ENABLED.String())
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
		f.Manager.TaskChan = make(chan FlagReq, 100)
		f.Manager.ErrorChan = make(chan error, 100)

		//	newCtx := context.Background()
		//if at all we get multiple customers for one feature, we can run all of them in background
		for i := 0; i < f.Manager.WorkerCount; i++ {
			f.Manager.Wg.Add(1)
			go f.Manager.Work(ctx, f.redisPool.AddToSetOfcustomers, i)
		}

		for _, request := range requests {
			f.Manager.TaskChan <- request
		}
		close(f.Manager.TaskChan)

		var errors error
		errDone := make(chan struct{})
		go func() {
			for {
				select {
				case err, ok := <-f.Manager.ErrorChan:
					if ok {
						log.Errorf("got error %v", err)
						errors = multierror.Append(errors, err)
						log.Errorf("got erro1r %v", errors)
					} else {
						log.Info("returning")
						close(errDone)
						return
					}

				}
			}
		}()
		f.Manager.Wg.Wait()
		close(f.Manager.ErrorChan)
		//time.Sleep(5*time.Second)
		<-errDone
		defer log.Errorf("blhgsgh   %v %v", len(f.Manager.ErrorChan), errors)
		return errors

	case proto.CreateFlagReq_REFERENCE_TYPE.String():
		type ref struct {
			Ref     string `json:"ref,omitempty"`
			Feature string `json:"feature,omitempty"`
		}
		var data ref
		err := json.Unmarshal(flagData, &data)
		if err != nil {
			log.WithError(err).Errorf("Invalid data passed for flag type %v", proto.CreateFlagReq_REFERENCE_TYPE.String())
			return err
		}
		return f.redisPool.AddToRef(data.Ref, data.Feature)

	case proto.CreateFlagReq_PERCENTAGE_OF.String():

		var data PercentageReq
		err := json.Unmarshal(flagData, &data)
		if err != nil {
			log.WithError(err).Errorf("Invalid data passed for flag type %v", proto.CreateFlagReq_PERCENTAGE_OF.String())
			return err
		}
		getAllCustomers, err := f.redisPool.GetAllCustomers(data.Ref)
		if err != nil {
			log.WithError(err).Errorf("Failed to get customers with Ref %v", data.Ref)
			return err
		}
		if len(getAllCustomers) == 0 {
			log.Warn("no customers to add the feature")
			return nil
		}
		percent, err := strconv.Atoi(data.Percentage)
		if err != nil {
			log.WithError(err).Error("Can't convert string to int")
		}

		randCustomersLength := int(math.Ceil(float64(percent*len(getAllCustomers)) / 100))
		rand.Seed(time.Now().UnixNano())
		randomIntegers := rand.Perm(len(getAllCustomers))
		var req []FlagReq

		for i := 0; i < randCustomersLength; i++ {
			customerDetails := strings.Split(getAllCustomers[randomIntegers[i]], "::")
			customerName, customerId := customerDetails[0], customerDetails[1]
			req = append(req, FlagReq{customerId, customerName, data.Feature})
		}
		log.Infof("Adding flag to customers %v", req)
		marshalReq, err := json.Marshal(req)
		if err != nil {
			log.WithError(err).Error("can't marshal request for flag percentage OF")
			return err
		}
		// once we have the random customers it'll jsut be a setOf Request flag
		return UnmarshalandStore(ctx, f, marshalReq, proto.CreateFlagReq_SET_OF.String())

	case proto.CreateFlagReq_COMBINATION_OF.String():

		type req struct {
			ReqType string      `json:"req_type"`
			Value   interface{} `json:"value"`
		}
		var unmarshalProto []req
		err := json.Unmarshal(flagData, &unmarshalProto)
		if err != nil {
			log.WithError(err).Errorf("Invalid data passed for flag type %v", proto.CreateFlagReq_COMBINATION_OF.String())
			return err
		}

		var reteErr error
		for _, req := range unmarshalProto {

			marshalData, err := json.Marshal(req.Value)
			if err != nil {
				log.WithError(err).Errorf("failed to marshal in combination of")
			}
			resp := UnmarshalandStore(ctx, f, marshalData, req.ReqType)
			if resp != nil {
				reteErr = multierror.Append(reteErr, resp)
			}

		}
		return reteErr

	}
	return nil
}
