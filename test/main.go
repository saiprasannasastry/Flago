package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/segmentio/feature-flag/pkg/proto"
	"google.golang.org/grpc"
	"time"
)

func main() {

	requestCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	//requestCtx := context.Background()
	conn, err := grpc.Dial("localhost:5566", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to establish client connection to 0.0.0.0:5566: %v", err)
	}
	client := proto.NewFlagoServiceClient(conn)

	fmt.Println("***")

	fData := []map[string]string{
		{"customer_id": "1", "customer_name": "Twilio", "feature": "feature1"},
		{"customer_id": "1", "customer_name": "Twilio", "feature": "feature2"},
		{"customer_id": "2", "customer_name": "Segment", "feature": "feature1"},
		{"customer_id": "2", "customer_name": "Segment", "feature": "feature2"},
		{"customer_id": "1", "customer_name": "Twilio", "feature": "feature2"},
		{"customer_id": "2", "customer_name": "Segment", "feature": "feature1"},
		{"customer_id": "2", "customer_name": "Segment", "feature": "feature2"},
		{"customer_id": "1", "customer_name": "Twilio", "feature": "feature1"},
	}
	v, _ := json.Marshal(fData)
	_, err = client.CreateFlag(requestCtx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_SET_OF,
		FlagData: v,
	})
	log.Infof("%v,", err)

	flagData := "feature55"
	marsalData, _ := json.Marshal(flagData)
	_, err = client.CreateFlag(requestCtx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_UNIVERSAL_ENABLED, FlagData: marsalData})
	_, err = client.CreateFlag(requestCtx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_UNIVERSAL_DISABLED, FlagData: marsalData})

	var dat []map[string]interface{}
	dat = append(dat, map[string]interface{}{"req_type": proto.CreateFlagReq_REFERENCE_TYPE.String(), "value": map[string]string{"ref": "alpha", "feature": "feature0"}})
	dat = append(dat, map[string]interface{}{"req_type": proto.CreateFlagReq_PERCENTAGE_OF.String(), "value": map[string]string{"ref": "beta", "feature": "feature15", "percentage": "15"}})
	d, _ := json.Marshal(dat)
	_, err = client.CreateFlag(requestCtx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_COMBINATION_OF,
		FlagData: d,
	})
	if err != nil {
		log.Errorf("error during combinatio nof %v", err)
	}


	resp, err := client.GetFlag(requestCtx, &proto.FlagReq{Feature: "feature1", CustomerId: "1", CustomerName: "Twilio"})

	log.Infof("get flag %v", resp.Enabled)
	if err != nil {
		log.Errorf("%v", err)
	}

	stream, err := client.GetFlags(requestCtx, &proto.FlagReq{Feature: "feature1s", CustomerId: "1", CustomerName: "Twilio"})
	if err != nil {
		log.Errorf("%v", err)
	}

	log.Infof("GetFlags %v", stream.Flags)

	on, err := client.OnFlag(requestCtx, &proto.FlagReq{Feature: "feature1s", CustomerId: "1", CustomerName: "Twilio"})
	if err != nil {
		log.Errorf("%v", err)
	}

	log.Infof("dsi %v", on.Enabled)
	off, err := client.OffFlag(requestCtx, &proto.FlagReq{Feature: "feature1s", CustomerId: "1", CustomerName: "Twilio"})
	if err != nil {
		log.Errorf("%v", err)
	}

	log.Infof("dsi %v", off.Enabled)
}
