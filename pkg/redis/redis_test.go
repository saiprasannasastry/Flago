package redis

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"testing"
)

func TestRedis(t *testing.T) {
	redisInstance := miniredis.NewMiniRedis()
	defer redisInstance.Close()
	redisInstance.Start()

	redisClient := redis.NewClient(&redis.Options{Addr: redisInstance.Addr()})
	store, errNew := GetRedisClient(redisInstance.Addr(), 50, 5)
	require.NoError(t, errNew)
	client := NewPool(store)

	setup(redisClient)

	t.Run("Create disable all flag with feature 0", func(t *testing.T) {
		feature := "Feature0"
		err := client.DisableAllCustomers(feature)
		require.NoError(t, err)

		disableCount, enableCount := disablekeysCount(t, redisClient), enableCount(t, redisClient)
		assert.Equal(t, 1, disableCount)
		assert.Equal(t, 0, enableCount)

	})
	t.Run("Create disable all flag with feature 0 require Error and feature 1 with no error", func(t *testing.T) {
		feature := "Feature0"
		err := client.DisableAllCustomers(feature)
		require.Error(t, err)
		disCount, enCount := disablekeysCount(t, redisClient), enableCount(t, redisClient)
		assert.Equal(t, 1, disCount)
		assert.Equal(t, 0, enCount)

		feature = "Feature1"
		err = client.DisableAllCustomers(feature)
		require.NoError(t, err)
		disCount, enCount = disablekeysCount(t, redisClient), enableCount(t, redisClient)
		assert.Equal(t, 2, disCount)
		assert.Equal(t, 0, enCount)

	})

	t.Run("Create enable all flag with feature 0", func(t *testing.T) {
		feature := "Feature0"
		err := client.EnableAllCustomers(feature)
		require.NoError(t, err)
		disCount, enCount := disablekeysCount(t, redisClient), enableCount(t, redisClient)
		assert.Equal(t, 1, disCount)
		assert.Equal(t, 1, enCount)

	})
	t.Run("Create enable all flag with feature 0 require Error and feature 1 with no error", func(t *testing.T) {

		feature := "Feature0"
		err := client.EnableAllCustomers(feature)
		require.Error(t, err)

		disCount, enCount := disablekeysCount(t, redisClient), enableCount(t, redisClient)
		assert.Equal(t, 1, disCount)
		assert.Equal(t, 1, enCount)

		feature = "Feature1"
		err = client.EnableAllCustomers(feature)
		require.NoError(t, err)

		disCount, enCount = disablekeysCount(t, redisClient), enableCount(t, redisClient)
		assert.Equal(t, 0, disCount)
		assert.Equal(t, 2, enCount)
	})
	t.Run("enable flag for set of customers", func(t *testing.T) {
		teardown(redisClient)
		feature1, feature2 := "Feature0", "Feature1"
		err := client.EnableAllCustomers(feature1)
		require.NoError(t, err)

		err = client.DisableAllCustomers(feature2)
		require.NoError(t, err)

		disCount, enCount := disablekeysCount(t, redisClient), enableCount(t, redisClient)
		assert.Equal(t, 1, disCount)
		assert.Equal(t, 1, enCount)

		customerName := "test"
		customerId := "123"

		err = client.AddToSetOfcustomers(customerName, customerId, feature1)
		require.NoError(t, err)

		err = client.AddToSetOfcustomers(customerName, customerId, feature1)
		require.Error(t, err)

		err = client.AddToSetOfcustomers(customerName, customerId, feature2)
		require.NoError(t, err)

		disCount, enCount = disablekeysCount(t, redisClient), enableCount(t, redisClient)
		assert.Equal(t, 0, disCount)
		assert.Equal(t, 0, enCount)

	})

}

func teardown(client *redis.Client) {
	client.Del(DISABLE_ALL_KEY)
	client.Del(ENABLE_ALL_KEY)
}

func setup(client *redis.Client) {
	client.SAdd(DISABLE_ALL_KEY)
	client.SAdd(ENABLE_ALL_KEY)
}

func disablekeysCount(t *testing.T, client *redis.Client) int {
	ctx := context.Background()
	disableCount, err := client.WithContext(ctx).SMembers(DISABLE_ALL_KEY).Result()
	require.NoError(t, err)
	return len(disableCount)
}

func enableCount(t *testing.T, client *redis.Client) int {
	ctx := context.Background()
	enableCount, err := client.WithContext(ctx).SMembers(ENABLE_ALL_KEY).Result()
	require.NoError(t, err)
	return len(enableCount)
}
