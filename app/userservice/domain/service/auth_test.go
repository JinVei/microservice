package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jinvei/microservice/base/framework/cache"
	"github.com/jinvei/microservice/base/framework/configuration"
)

func TestSessionCache(t *testing.T) {
	os.Setenv("MICROSERVICE_CONFIGURATION_TOKEN", "e30K")
	configuration.SetSystemID("10001")

	rediscli := cache.RedisClient(nil)

	if rediscli == nil {
		t.Fatal()
	}
	seesionsKey := fmt.Sprintf(sessionKeyFormat, "1", "sse1")
	skeys, err := rediscli.Keys(context.Background(), seesionsKey).Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("keys3", skeys)
}

func TestSessionCache1(t *testing.T) {
	os.Setenv("MICROSERVICE_CONFIGURATION_TOKEN", "e30K")
	configuration.SetSystemID("10001")

	rediscli := cache.RedisClient(nil)
	if rediscli == nil {
		t.Fatal()
	}
	us := UserSession{
		rdb:         rediscli,
		sessDura:    time.Minute * 5,
		maxSesssion: 3,
	}
	us.AddSession(context.TODO(), "11", "s33")
	us.AddSession(context.TODO(), "11", "s22")
	us.AddSession(context.TODO(), "11", "s11")

	us.AddSession(context.TODO(), "22", "s23")
	us.AddSession(context.TODO(), "22", "s11")
}
