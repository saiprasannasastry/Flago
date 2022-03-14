package service

import (
	"context"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/segmentio/feature-flag/pkg/proto"
	red "github.com/segmentio/feature-flag/pkg/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	DISABLE_ALL_KEY = "disable:all"
	ENABLE_ALL_KEY  = "enable:all"
)

func TestCreateServer(t *testing.T) {
	redisInstance := miniredis.NewMiniRedis()
	defer redisInstance.Close()
	redisInstance.Start()
	//setup redis
	redisClient := redis.NewClient(&redis.Options{Addr: redisInstance.Addr()})
	store, errNew := red.GetRedisClient(redisInstance.Addr(), 50, 5)
	require.NoError(t, errNew)
	client := red.NewPool(store)
	setup(redisClient)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//setup server details
	manager := NewConcurrentManager(5, 2, 3, 1*time.Minute)
	flagReq := NewFlagoServer(client, manager)

	t.Run("test Disable All Flag", func(t *testing.T) {
		flagData := map[string]string{"feature": "feature0"}
		marsalData, _ := json.Marshal(flagData)
		_, err := flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_UNIVERSAL_DISABLED, FlagData: marsalData})

		require.Error(t, err, "bad marshal data")
		marsalData, _ = json.Marshal("feature1")
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_UNIVERSAL_DISABLED, FlagData: marsalData})

		require.NoError(t, err, "good marshal data")
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_UNIVERSAL_DISABLED, FlagData: marsalData})

		require.Error(t, err, "flag already enabled")
		disCount, enCount := disablekeysCount(t, redisClient), enableCount(t, redisClient, ENABLE_ALL_KEY)
		assert.Equal(t, disCount, 1)
		assert.Equal(t, enCount, 0)
	})

	t.Run("test Enable All Flag", func(t *testing.T) {
		flagData := map[string]string{"feature": "feature0"}
		marsalData, _ := json.Marshal(flagData)
		_, err := flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_UNIVERSAL_ENABLED, FlagData: marsalData})

		require.Error(t, err, "bad marshal data")
		marsalData, _ = json.Marshal("feature1")
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_UNIVERSAL_ENABLED, FlagData: marsalData})

		require.NoError(t, err, "good marshal data")
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_UNIVERSAL_ENABLED, FlagData: marsalData})

		require.Error(t, err, "flag already enabled")
		disCount, enCount := disablekeysCount(t, redisClient), enableCount(t, redisClient, ENABLE_ALL_KEY)
		assert.Equal(t, disCount, 0)
		assert.Equal(t, enCount, 1)
		marsalData, _ = json.Marshal("feature2")
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_UNIVERSAL_ENABLED, FlagData: marsalData})

		require.NoError(t, err, "adding 2nd feature")
		disCount, enCount = disablekeysCount(t, redisClient), enableCount(t, redisClient, ENABLE_ALL_KEY)
		assert.Equal(t, disCount, 0)
		assert.Equal(t, enCount, 2)
	})

	t.Run("create flag req SET_OF", func(t *testing.T) {
		marsalData, _ := json.Marshal("feature1")
		_, err := flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_UNIVERSAL_DISABLED, FlagData: marsalData})

		require.NoError(t, err, "good marshal data")
		disCount, enCount := disablekeysCount(t, redisClient), enableCount(t, redisClient, ENABLE_ALL_KEY)
		assert.Equal(t, disCount, 1)
		//because in previous test case we'd enabled the same feature
		assert.Equal(t, enCount, 1)

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
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_SET_OF,
			FlagData: v,
		})
		disCount, enCount = disablekeysCount(t, redisClient), enableCount(t, redisClient, ENABLE_ALL_KEY)
		assert.Equal(t, disCount, 0)
		require.Error(t, err)

		twCount, segCount := enableCount(t, redisClient, "Twilio::1"), enableCount(t, redisClient, "Segment::2")
		assert.Equal(t, twCount, 2)
		assert.Equal(t, segCount, 2)

		fData = []map[string]string{
			{"customer_id": "5", "customer_name": "NewCustomer", "feature": "feature1"}}
		v, _ = json.Marshal(fData)
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_SET_OF,
			FlagData: v,
		})
		newCount := enableCount(t, redisClient, "NewCustomer::5")
		assert.Equal(t, newCount, 1)
		require.NoError(t, err)

		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_SET_OF,
			FlagData: v,
		})
		require.Error(t, err, "same data")

		badfData := "test"
		v, _ = json.Marshal(badfData)
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_SET_OF,
			FlagData: v,
		})
		require.Error(t, err, "bad marshal error")
	})
	t.Run("create flag req REF_TYPE", func(t *testing.T) {
		d, _ := json.Marshal(map[string]string{"ref": "alpha", "feature": "feature0"})
		go func() {
			_, err := flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_REFERENCE_TYPE,
				FlagData: d,
			})
			require.NoError(t, err)
		}()
		time.Sleep(1 * time.Second)
		test1, test2 := enableCount(t, redisClient, "test::1"), enableCount(t, redisClient, "test::2")
		assert.Equal(t, 1, test1)
		assert.Equal(t, 1, test2)

		d, _ = json.Marshal(map[string]string{"ref": "alpha", "feature": "feature1"})
		_, err := flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_REFERENCE_TYPE,
			FlagData: d,
		})
		require.NoError(t, err)
		test1, test2 = enableCount(t, redisClient, "test::1"), enableCount(t, redisClient, "test::2")
		assert.Equal(t, 2, test1)
		assert.Equal(t, 2, test2)

		d, _ = json.Marshal(map[string]string{"ref": "alpha", "feature": "feature1"})
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_REFERENCE_TYPE,
			FlagData: d,
		})
		require.Error(t, err, "value already present")

		d, _ = json.Marshal(map[string]string{"ref": "blag", "feature": "feature1"})
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_REFERENCE_TYPE,
			FlagData: d,
		})
		require.Error(t, err, "ref does not exist in Redis")

		d, _ = json.Marshal("bad data")
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_REFERENCE_TYPE,
			FlagData: d,
		})
		require.Error(t, err, "bad marshal data")

	})

	t.Run("create flag req PERCENTAGE_OF", func(t *testing.T) {
		d, _ := json.Marshal(map[string]string{"ref": "beta", "feature": "feature15", "percentage": "100"})

		_, err := flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_PERCENTAGE_OF,
			FlagData: d,
		})
		require.NoError(t, err, "random Data")
	})

	t.Run("create flag req COMBINATION_OF", func(t *testing.T) {

		var dat []map[string]interface{}
		dat = append(dat, map[string]interface{}{"req_type": proto.CreateFlagReq_REFERENCE_TYPE.String(), "value": map[string]string{"ref": "alpha", "feature": "feature0"}})
		d, _ := json.Marshal(dat)
		_, err := flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_COMBINATION_OF,
			FlagData: d,
		})
		require.Error(t, err, "flags already enabled for same feature above")

		var newDat []map[string]interface{}
		newDat = append(newDat, map[string]interface{}{"req_type": proto.CreateFlagReq_REFERENCE_TYPE.String(), "value": map[string]string{"ref": "alpha", "feature": "feature73"}})
		newD, _ := json.Marshal(newDat)
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_COMBINATION_OF,
			FlagData: newD,
		})

		var newData []map[string]interface{}
		newData = append(newData, map[string]interface{}{"req_type": proto.CreateFlagReq_PERCENTAGE_OF.String(), "value": map[string]string{"ref": "beta", "feature": "feature15", "percentage": "15"}})

		d, _ = json.Marshal(newData)
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_COMBINATION_OF,
			FlagData: d,
		})

		require.Error(t, err, "we enabled 100 % of flags above")

		var newData1 []map[string]interface{}
		fData := []map[string]string{
			{"customer_id": "1", "customer_name": "Twilio", "feature": "feature55"}}
		newData1 = append(newData1, map[string]interface{}{"req_type": proto.CreateFlagReq_SET_OF.String(), "value": fData})
		newMarshal, _ := json.Marshal(newData1)
		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_COMBINATION_OF,
			FlagData: newMarshal,
		})
		require.NoError(t, err)

		_, err = flagReq.CreateFlag(ctx, &proto.CreateFlagReq{FlagFamily: proto.CreateFlagReq_COMBINATION_OF,
			FlagData: newMarshal,
		})
		require.Error(t, err, "same flag as above")

	})

	tearDown(redisClient)
}

func setup(client *redis.Client) {
	client.SAdd(DISABLE_ALL_KEY)
	client.SAdd(ENABLE_ALL_KEY)
	client.SAdd("alpha", "test::1", "test::2")
	client.SAdd("beta", "google::65", "adobe::88", "rudder::12", "sap::15", "cisco::18", "microsoft::10", "uber::11")
}
func tearDown(client *redis.Client) {
	client.FlushAll()
}
func enableCount(t *testing.T, client *redis.Client, key string) int {
	ctx := context.Background()
	enableCount, err := client.WithContext(ctx).SMembers(key).Result()
	require.NoError(t, err)
	return len(enableCount)
}
func disablekeysCount(t *testing.T, client *redis.Client) int {
	ctx := context.Background()
	disableCount, err := client.WithContext(ctx).SMembers(DISABLE_ALL_KEY).Result()
	require.NoError(t, err)
	return len(disableCount)
}
