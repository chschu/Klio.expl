package main

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source expldb/expldb.go        -destination generated/expldb_mocks/expldb.go        -package expldb_mocks
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source types/entry_stringer.go -destination generated/types_mocks/entry_stringer.go -package types_mocks
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source webhook/add_handler.go  -destination generated/webhook_mocks/add_handler.go  -package webhook_mocks
