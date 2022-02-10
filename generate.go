package requester

import _ "github.com/golang/mock/mockgen/model"

//go:generate mockgen -destination=mocks/mock_hash_getter.go -package=mocks requester/getter HashGetter
