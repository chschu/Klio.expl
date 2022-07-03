package main

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source expldb/expldb.go            -destination generated/expldb_mocks/expldb.go            -package expldb_mocks
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source security/jwt.go             -destination generated/security_mocks/jwt.go             -package security_mocks
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source types/entry_stringer.go     -destination generated/types_mocks/entry_stringer.go     -package types_mocks
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source types/index_spec.go         -destination generated/types_mocks/index_spec.go         -package types_mocks
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source types/index_spec_parser.go  -destination generated/types_mocks/index_spec_parser.go  -package types_mocks
