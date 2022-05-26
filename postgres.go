package main

import (
	"github.com/law-a-1/product-service/ent"
	_ "github.com/lib/pq"
)

func NewPersistent() (*ent.Client, error) {
	client, err := ent.Open(
		"postgres",
		"host=localhost port=5432 user=postgres dbname=productdb password=root sslmode=disable")
	if err != nil {
		return nil, err
	}
	return client, nil
}
