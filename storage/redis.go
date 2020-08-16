package storage

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Open(url string) error {
	r.client = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := r.client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Insert(data Data) error {
	b, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	err = r.client.Set(data.IP, string(b), 0).Err()
	if err != nil {
		return err
	}
	fmt.Println(data.IP)
	return nil
}

func (r *Redis) InsertBulk(bulk []interface{}) error {
	return nil
}

func (r *Redis) Get(ip string) (*Data, error) {
	val, err := r.client.Get(ip).Result()
	if err == redis.Nil {
		fmt.Printf("%s is not exist\n", ip)
		return &Data{}, nil
	} else if err != nil {
		return &Data{}, err
	}
	b := []byte(val)
	data := &Data{}

	err = json.Unmarshal(b, data)
	if err != nil {
		return &Data{}, err
	}
	return data, nil
}

func (r *Redis) Exist(collName string) (bool, error) {
	return false, nil
}
