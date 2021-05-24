package consul

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/mitchellh/mapstructure" //MIT license, is this a problem?
)

// The Trento Agent is periodically updating data structures in the Consul Key-Value Store
// to be consumed by the Trento Web Console
//
// At some point in time this structure needs to become a versioned protocol between the
// Trento Web Console and the Trento Agent

// For now, until we tighten the structure around versioning, this file is a start
// that records constants that need to be in sync or have compat behavior between
// Console and Agent

const (
	stringFlag uint64 = 0
	int32Flag  uint64 = 1
	boolFlag   uint64 = 2
	sliceFlag  uint64 = 3
	intFlag    uint64 = 4

	KvClustersPath       string = "trento/v0/clusters/"
	KvHostsPath          string = "trento/v0/hosts/"
	KvHostsMetadataPath  string = "trento/v0/hosts/%s/metadata/"
	KvHostsSAPSystemPath string = "trento/v0/hosts/%s/sapsystems/"
	KvEnvironmentsPath   string = "trento/v0/environments/"
)

type ClusterStonithType int

const (
	ClusterStonithNone ClusterStonithType = iota
	ClusterStonithSBD
	ClusterStonithUnknown
)

type KV interface {
	Get(key string, q *consulApi.QueryOptions) (*consulApi.KVPair, *consulApi.QueryMeta, error)
	List(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error)
	Keys(prefix, separator string, q *consulApi.QueryOptions) ([]string, *consulApi.QueryMeta, error)
	Put(p *consulApi.KVPair, q *consulApi.WriteOptions) (*consulApi.WriteMeta, error)
	DeleteTree(prefix string, w *consulApi.WriteOptions) (*consulApi.WriteMeta, error)
	ListMap(prefix, offset string) (map[string]interface{}, error)
	PutMap(prefix string, data map[string]interface{}) error
	PutTyped(prefix string, value interface{}) error
}

func newKV(wrapped *consulApi.KV) KV {
	return &kv{
		wrapped,
		wrapped.List,
		wrapped.Put,
	}
}

type kv struct {
	wrapped *consulApi.KV
	// we need this dedicated function fields because the some of the internal KV methods depends on them
	// and we want to mock it internally
	list func(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error)
	put  func(p *consulApi.KVPair, q *consulApi.WriteOptions) (*consulApi.WriteMeta, error)
}

func (k *kv) Get(key string, q *consulApi.QueryOptions) (*consulApi.KVPair, *consulApi.QueryMeta, error) {
	return k.wrapped.Get(key, q)
}

func (k *kv) Keys(prefix, separator string, q *consulApi.QueryOptions) ([]string, *consulApi.QueryMeta, error) {
	return k.wrapped.Keys(prefix, separator, q)
}

func (k *kv) Put(p *consulApi.KVPair, q *consulApi.WriteOptions) (*consulApi.WriteMeta, error) {
	return k.put(p, q)
}

func (k *kv) DeleteTree(prefix string, w *consulApi.WriteOptions) (*consulApi.WriteMeta, error) {
	return k.wrapped.DeleteTree(prefix, w)
}

func (k *kv) List(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error) {
	return k.list(prefix, q)
}

// Maps converts the KVs output to a Map
// prefix -> The KV prefix to get the data
// offset -> Remove the offset from the outgoing map
// 	         (used to remove initial key from the map like trento/v0/)
func (k *kv) ListMap(prefix, offset string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	entries, _, err := k.list(prefix, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for KV values")
	}

	if !strings.HasSuffix(offset, "/") {
		offset = fmt.Sprintf("%s/", offset)
	}

	currentItem := result

	for _, entry := range entries {
		modEntry := strings.TrimPrefix(strings.Trim(entry.Key, " "), offset)
		if len(modEntry) == 0 {
			continue
		}

		lastKey := ""
		sliceFound := false
		currentItem = result
		keys := strings.Split(modEntry, "/")
		for i, key := range keys {
			// Handle slice type
			if sliceFound {
				keyInt, _ := strconv.Atoi(key)
				value := reflect.ValueOf(currentItem[lastKey])
				if value.Len() == keyInt {
					newSliceItem := make(map[string]interface{})
					currentItem[lastKey] = append(currentItem[lastKey].([]interface{}), newSliceItem)
					currentItem = newSliceItem
				} else {
					currentItem = value.Index(keyInt).Interface().(map[string]interface{})
				}
				sliceFound = false
				continue
				// Last element slice
			} else if entry.Flags == sliceFlag && i == len(keys)-2 {
				value := getTypeByFlag(entry)
				currentItem[key] = value
				break
				// Other types
			} else if i == len(keys)-1 {
				value := getTypeByFlag(entry)
				if value != "" {
					currentItem[key] = value
				}
				break
			}

			value := reflect.ValueOf(currentItem[key])
			if value.Kind() == reflect.Slice {
				lastKey = key
				sliceFound = true
				continue
			}

			if _, ok := currentItem[key]; !ok {
				currentItem[key] = make(map[string]interface{})
			}

			item, _ := currentItem[key].(map[string]interface{})
			currentItem = item
		}
	}
	return result, nil
}

func getTypeByFlag(entry *consulApi.KVPair) interface{} {
	var value interface{}

	switch entry.Flags {
	case stringFlag:
		value = string(entry.Value)
	case int32Flag:
		i, err := strconv.ParseInt(string(entry.Value), 10, 32)
		if err != nil {
			return err
		}
		value = int32(i)
	case intFlag:
		i, err := strconv.Atoi(string(entry.Value))
		if err != nil {
			return err
		}
		value = i
	case boolFlag:
		b, err := strconv.ParseBool(string(entry.Value))
		if err != nil {
			return err
		}
		value = b
	case sliceFlag:
		value = make([]interface{}, 0)
	}

	return value
}

// Store a map[string]interface data in KV storage under the prefix key
func (k *kv) PutMap(prefix string, data map[string]interface{}) error {
	if !strings.HasSuffix(prefix, "/") {
		prefix = fmt.Sprintf("%s/", prefix)
	}

	// Empty KV directories
	if len(data) == 0 {
		err := k.PutTyped(fmt.Sprintf("%s", prefix), "")
		if err != nil {
			return err
		}
		return nil
	}

	for key, value := range data {
		switch reflect.ValueOf(value).Kind() {
		case reflect.Map, reflect.Struct, reflect.Ptr:
			mapInterface := make(map[string]interface{})
			mapstructure.Decode(value, &mapInterface)
			err := k.PutMap(fmt.Sprintf("%s%s", prefix, key), mapInterface)
			if err != nil {
				return err
			}
		case reflect.Slice:
			// Store the slice with slice flag, to be able to reload as list in the ListMap funciton
			err := k.PutTyped(fmt.Sprintf("%s%s/", prefix, key), []string{})
			if err != nil {
				return err
			}

			// Store slice elements
			sliceValue := reflect.ValueOf(value)
			for i := 0; i < sliceValue.Len(); i++ {
				mapInterface := make(map[string]interface{})
				mapstructure.Decode(sliceValue.Index(i).Interface(), &mapInterface)
				// Index is composed by 4 digits to keep correct numbers order in KV storage
				index := fmt.Sprintf("%04d", i)
				err = k.PutMap(fmt.Sprintf("%s%s/%s", prefix, key, index), mapInterface)
				if err != nil {
					return err
				}
			}
		default:
			err := k.PutTyped(fmt.Sprintf("%s%s", prefix, key), value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k *kv) PutTyped(prefix string, value interface{}) error {
	var flag uint64 = stringFlag // By default string flag is used

	switch reflect.ValueOf(value).Kind() {
	case reflect.Int32:
		flag = int32Flag
	case reflect.Bool:
		flag = boolFlag
	case reflect.Slice:
		flag = sliceFlag
		value = ""
	case reflect.Int:
		flag = intFlag
	}

	_, err := k.put(&consulApi.KVPair{
		Key:   prefix,
		Value: []byte(fmt.Sprintf("%v", value)),
		Flags: flag,
	}, nil)

	if err != nil {
		return errors.Wrap(err, "Error storing a new value in the KV storage")
	}

	// TO-DO make this a debug log statement when we introduce levelled logging
	//log.Printf("Value %s properly stored at %s", value, prefix)
	return nil
}
