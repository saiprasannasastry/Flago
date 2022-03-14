package redis

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"sync"
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

		feature1, feature2 := "Feature0", "Feature1"
		err := client.EnableAllCustomers(feature1)
		//already enabled above
		require.Error(t, err)

		err = client.DisableAllCustomers(feature2)
		require.NoError(t, err)

		disCount, enCount := disablekeysCount(t, redisClient), enableCount(t, redisClient)
		assert.Equal(t, 1, disCount)
		assert.Equal(t, 1, enCount)

		customerName := "test"
		customerId := "123"

		err = client.AddToSetOfcustomers(customerName, customerId, feature1)
		require.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			err = client.AddToSetOfcustomers(customerName, customerId, "feature3")
			require.NoError(t, err)
		}()

		err = client.AddToSetOfcustomers(customerName, customerId, feature1)
		require.Error(t, err)

		err = client.AddToSetOfcustomers(customerName, customerId, feature2)
		require.NoError(t, err)
		wg.Wait()
		disCount, enCount = disablekeysCount(t, redisClient), enableCount(t, redisClient)
		assert.Equal(t, 0, disCount)
		assert.Equal(t, 1, enCount)
	})

	t.Run("enable flag for a reference Type of customers (with and without errors)", func(t *testing.T) {
		feature0, feature1 := "feature0", "feature1"

		err := client.AddToRef("bets", feature0)
		require.Error(t, err)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			err = client.AddToRef("alpha", feature0)
			require.NoError(t, err)
		}()
		wg.Wait()
		err = client.AddToRef("alpha", feature1)
		require.NoError(t, err)
	})

	t.Run("Get all customers", func(t *testing.T) {

		customers, err := client.GetAllCustomers("alpha")
		require.NoError(t, err)
		assert.Equal(t, len(customers), 2)
	})

	t.Run("Get Flag for customer", func(t *testing.T) {

		flag, err := client.GetFlagForCustomer("test::1", "feature0")
		require.NoError(t, err)
		assert.Equal(t, flag, true)

		flag, err = client.GetFlagForCustomer("test::153", "feature0")
		require.NoError(t, err)
		assert.Equal(t, flag, false)
	})

}

func teardown(client *redis.Client) {
	client.FlushAll()
}

func setup(client *redis.Client) {
	client.SAdd(DISABLE_ALL_KEY)
	client.SAdd(ENABLE_ALL_KEY)
	client.SAdd("alpha", "test::1", "test::2")

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
