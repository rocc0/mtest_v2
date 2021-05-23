package handlers

import (
	datapkg "mtest.com.ua/db/dataprocessor"
	hashpkg "mtest.com.ua/db/hasher"
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