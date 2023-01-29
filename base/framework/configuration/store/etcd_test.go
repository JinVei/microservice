package store_test

import (
	"testing"

	"github.com/jinvei/microservice/base/framework/configuration/store"
)

func TestEtcdStoreGet(t *testing.T) {
	sto, err := store.NewEtcdStore("e30K")
	if err != nil {
		t.Fatal(err)
	}
	val, err := sto.Get("/microserice/configuration/test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("val:", val)
}

func TestEtcdStoreGetObj(t *testing.T) {
	sto, err := store.NewEtcdStore("e30K")
	if err != nil {
		t.Fatal(err)
	}
	obj := struct {
		Key    string `json:"key"`
		Status string `json:"status"`
	}{}
	err = sto.GetJson("/microserice/configuration/test/obj", &obj)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("val:", obj)
}
