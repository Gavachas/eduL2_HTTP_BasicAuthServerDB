package sessiontoken

import (
	"time"

	"github.com/go-redis/redis"
)

type Session struct {
	Driver *redis.Client
	Name   string
	TTL    int64 // seconds
}

func (s *Session) Set(key string, value interface{}) error {
	exp := time.Duration(time.Duration(s.TTL) * time.Second) // 10 minutes
	return s.Driver.Set(key, value, exp).Err()
}

// MSet sets multiple key values
func (s *Session) MSet(keys []string, values []interface{}) error {
	exp := time.Duration(time.Duration(s.TTL) * time.Second) // 10 minutes
	var ifaces []interface{}
	pipe := s.Driver.TxPipeline()
	for i := range keys {
		ifaces = append(ifaces, keys[i], values[i])
		pipe.Expire(keys[i], exp)
	}

	if err := s.Driver.MSet(ifaces...).Err(); err != nil {
		return err
	}
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	return nil
}

// Get gets value
func (s *Session) Get(key string) (interface{}, error) {
	val, err := s.Driver.Get(key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

// MGet get multiple value
func (s *Session) MGet(keys []string) ([]interface{}, error) {
	return s.Driver.MGet(keys...).Result()
}

/*func (this *Session) Set(key string, value interface{}) error {
	var h = this.Driver
	var bytes []byte
	var err error
	var content string
	m := make(map[string]interface{})

	if h.Exists(this.Name).Val() == 1 {
		content = h.Get(this.Name).Val()
	} else {
		content = "{}"
	}
	fmt.Printf("%s => %s", this.Name, content)

	err = json.Unmarshal([]byte(content), &m)
	if err != nil {
		return err
	}

	var keys = strings.Split(key, ".")
	var depth = len(keys)

	if depth < 2 {
		m[key] = value
	} else {
		m = setSliceMap(m, keys, value)
	}

	bytes, _ = json.Marshal(m)
	return h.Set(this.Name, string(bytes), time.Duration(this.TTL)*time.Second).Err()
}

// s.Get
func (this *Session) Get(key string) (interface{}, error) {
	var h = this.Driver
	var m map[string]interface{}
	var content string
	var err error
	if err != nil {
		return nil, err
	}
	content = h.Get(this.Name).Val()
	err = json.Unmarshal([]byte(content), &m)

	var keys = strings.Split(key, ".")
	var n = len(keys)
	if n < 2 {
		return m[key], nil
	}
	return getSliceMap(m, keys)
}
func setSliceMap(m map[string]interface{}, keys []string, value interface{}) map[string]interface{} {
	var itMap = m
	var i int
	var limit = len(keys) - 1

	for i = 0; i < limit; i++ {
		_, ok := itMap[keys[i]]
		if ok {
			// fmt.Printf("%s yes\n", keys[i])
		} else {
			// fmt.Printf("%s no\n", keys[i])
			itMap[keys[i]] = make(map[string]interface{})
		}
		itMap = itMap[keys[i]].(map[string]interface{})
	}
	itMap[keys[limit]] = value

	return m
}

func getSliceMap(m map[string]interface{}, keys []string) (interface{}, error) {
	var itMap = m
	var i int
	var limit = len(keys) - 1
	var v interface{}
	var ok bool

	for i = 0; i < limit; i++ {
		v, ok = itMap[keys[i]]
		if !ok {
			break
		}
		itMap = v.(map[string]interface{})
	}
	v, ok = itMap[keys[i]]
	if !ok {
		return nil, fmt.Errorf("Err getSliceMap")
	}
	return v, nil
}*/
