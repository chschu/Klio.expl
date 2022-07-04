package main

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source webhook/add_handler.go          -destination generated/webhook_mocks/add_handler.go          -package webhook_mocks
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source webhook/expl_handler.go         -destination generated/webhook_mocks/expl_handler.go         -package webhook_mocks
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source webhook/find_handler.go         -destination generated/webhook_mocks/find_handler.go         -package webhook_mocks
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source webhook/interfaces.go           -destination generated/webhook_mocks/interfaces.go           -package webhook_mocks
