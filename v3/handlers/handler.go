package handlers

import (
	datapkg "mtest.com.ua/v3/db/dataprocessor"
	hashpkg "mtest.com.ua/v3/db/hasher"
)

func NewService(data *datapkg.Service, hash *hashpkg.HashHandler, search indexUpdater) *Handlers {
	return &Handlers{
		mtestDataProcessor:    data,
		executorDataProcessor: data,
		regionDataProcessor:   data,
		userDataProcessor:     data,
		hasher:                hash,
		indexUpdater:          search,
	}
}
