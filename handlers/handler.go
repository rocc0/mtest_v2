package handlers

import (
	datapkg "mtest.com.ua/db/dataprocessor"
	hashpkg "mtest.com.ua/db/hasher"
	"mtest.com.ua/mail"
)

func NewService(data *datapkg.Service, hash *hashpkg.HashHandler, search indexUpdater, auth mail.Auth) *Handlers {
	return &Handlers{
		mtestDataProcessor:      data,
		executorDataProcessor:   data,
		regionDataProcessor:     data,
		userDataProcessor:       data,
		admActionsProcessor:     data,
		governmentDataProcessor: data,
		regActUpdater:           data,
		BusinessDataProcessor:   data,
		hasher:                  hash,
		indexUpdater:            search,
		SynonymsProcessor:       data,
		auth:                    auth,
	}
}
