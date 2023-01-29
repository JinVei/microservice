package store

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	etcdcli "go.etcd.io/etcd/client/v3"
)

type EtcdStore struct {
	cli *etcdcli.Client
}

func NewEtcdStore(token string) (*EtcdStore, error) {
	decdat, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	conf := config{}

	if err := json.Unmarshal(decdat, &conf); err != nil {
		return nil, err
	}

	if len(conf.Addr) == 0 {
		// todo log
		conf.Addr = append(conf.Addr, "localhost:2379")
	}

	cli, err := etcdcli.New(etcdcli.Config{
		Endpoints:   conf.Addr,
		DialTimeout: 5 * time.Second,
		Username:    conf.Username,
		Password:    conf.Password,
	})
	if err != nil {
		return nil, err
	}

	store := EtcdStore{
		cli: cli,
	}

	return &store, nil
}

func (store *EtcdStore) Get(path string) (string, error) {
	resp, err := store.cli.Get(context.Background(), path)
	if err != nil {
		return "", err
	}
	if resp.Count == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

func (store *EtcdStore) GetJson(path string, obj interface{}) error {
	resp, err := store.cli.Get(context.Background(), path)
	if err != nil {
		return err
	}
	if resp.Count == 0 {
		return fmt.Errorf("key:`%s` Not Value", path)
	}

	val := resp.Kvs[0].Value
	err = json.Unmarshal(val, obj)
	if err != nil {
		return fmt.Errorf("path='%s' UnmarshalError:'%s', val='%s'", path, err, val)
	}
	return nil
}
