package redis

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/go-redis/redis"
	"github.com/hashicorp/go-multierror"
	"strings"
)

//EnableAllCustomers enables features for all customers
// feature ready for pro
func (p Pool) EnableAllCustomers(feature string) error {
	found := p.RedisClient.SMembers(ENABLE_ALL_KEY)
	foundKeys, err := found.Result()

	if err != nil {
		log.WithError(err).Errorf("failed to check key %v ", ENABLE_ALL_KEY)
		return err
	}
	keys := convertSliceToMap(foundKeys)

	if keys[feature] {
		newError := errors.New("flag already enabled for feature")
		log.Infof("feture %v already enabled for all customers", feature)
		return newError
	}

	removed := p.removeKey(DISABLE_ALL_KEY, feature)

	if removed != nil {
		log.WithError(removed).Errorf("failed to remove flag for feature %v for all customers ", feature)
		return removed
	}
	log.Infof("no key found for enabling feature %v, proceeding to insert", feature)
	return p.addKey(ENABLE_ALL_KEY, feature)
}

// DisableAllCustomers disables all customers for a particular feature
// Not ready for a prod release
func (p Pool) DisableAllCustomers(feature string) error {

	found := p.RedisClient.SMembers(DISABLE_ALL_KEY)

	foundKeys, err := found.Result()

	if err != nil {
		log.WithError(err).Errorf("failed to check key %v ", DISABLE_ALL_KEY)
		return err
	}
	keys := convertSliceToMap(foundKeys)

	if keys[feature] {
		log.Infof("feture %v already disabled for all customers", feature)
		return errors.New("flag already enabled for feature")
	}

	if found.Err() != redis.Nil {
		log.Infof("no key found for Disable all feature %v, proceeding to insert", feature)
		err := p.removeKey(ENABLE_ALL_KEY, feature)
		if err != nil {
			log.WithError(err).Errorf("failed to add key")
		}
		return p.addKey(DISABLE_ALL_KEY, feature)
	}

	return nil
}

//AddToSetOfcustomers enables the flag for the reature to AddToSetOfCustomers
func (p Pool) AddToSetOfcustomers(customerName string, customerId string, feature string) error {

	key := fmt.Sprintf("%v::%v", customerName, customerId)

	found := p.RedisClient.SMembers(key)

	foundValues, err := found.Result()

	if err != nil {
		log.WithError(err).Errorf("failed to check key %v ", key)
		return err
	}
	keys := convertSliceToMap(foundValues)

	if keys[feature] {
		newError := fmt.Errorf("flag already enabled for feature %v for customer %v", feature, customerName)
		log.Errorf("%v ", newError)
		return newError
	}

	//since we are enabling for one or more customers we need to remove from disable all
	p.RedisClient.SRem(DISABLE_ALL_KEY, feature)
	//p.RedisClient.SRem(ENABLE_ALL_KEY, feature)
	log.Infof("enabling feature %v for customer %v ", feature, customerName)
	return p.addKey(key, feature)
}

//AddToRef enables flag for all reference type customers
func (p Pool) AddToRef(refType string, feature string) error {
	foundKeys, err := p.GetAllCustomers(refType)
	if err != nil {
		log.WithError(err).Errorf("failed to get result for refType %v ", refType)
		return err
	}

	if len(foundKeys) == 0 {
		log.Errorf("refType not found %v", refType)
		return errors.New("refType not found")
	}
	var error error

	for _, customers := range foundKeys {
		//we want to insert to each company
		customerDetails := strings.Split(customers, "::")
		customerName, iD := customerDetails[0], customerDetails[1]
		err = p.AddToSetOfcustomers(customerName, iD, feature)
		if err != nil {
			error = multierror.Append(error, err)
		}
	}
	return error
}

func (p Pool) AddFlag(refType string, feature string) (int64, error) {
	return p.RedisClient.SAdd(refType, feature).Result()
}
func (p Pool) DeleteFlag(refType string, feature string) (int64, error) {
	return p.RedisClient.SRem(refType, feature).Result()
}

//Get all customers returns us list of customers for given data
func (p Pool) GetAllCustomers(refType string) ([]string, error) {
	return p.RedisClient.SMembers(refType).Result()
}

//Get flag for customer returns if the flag is enabled for customer
func (p Pool) GetFlagForCustomer(customerDetails, feature string) (bool, error) {
	return p.RedisClient.SIsMember(customerDetails, feature).Result()
}
func (p Pool) removeKey(keyName string, feature string) error {

	removed := p.RedisClient.SRem(keyName, feature)
	if removed.Err() != nil {
		log.WithError(removed.Err()).Errorf("failed to insert key for feature %v", feature)
		return removed.Err()
	}
	return nil
}

func (p Pool) addKey(keyName string, feature string) error {
	inserted := p.RedisClient.SAdd(keyName, feature)

	if inserted.Err() != nil {
		log.WithError(inserted.Err()).Errorf("failed to insert key for feature %v", feature)
		return inserted.Err()
	}
	return nil

}
func convertSliceToMap(foundKeys []string) map[string]bool {
	keys := map[string]bool{}
	for _, foundkey := range foundKeys {
		keys[foundkey] = true
	}
	return keys
}
