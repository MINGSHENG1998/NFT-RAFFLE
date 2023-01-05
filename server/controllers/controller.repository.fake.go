package controllers

import (
	"nft-raffle/models"
	"sync"
)

var (
	fakeControllerRepository     *fakeControllerRepositoryStruct
	fakeControllerRepositoryOnce sync.Once
)

type IFakeControllerRepository interface {
	MockFindUserByEmailSuccess(email string) (models.User, error)
	MockFindUserByIdSuccess(userId string) (models.User, error)
}

type fakeControllerRepositoryStruct struct{}

func GetFakeControllerRepository() *fakeControllerRepositoryStruct {
	if fakeControllerRepository == nil {
		fakeControllerRepositoryOnce.Do(func() {
			fakeControllerRepository = &fakeControllerRepositoryStruct{}
		})
	}
	return fakeControllerRepository
}

func (fake *fakeControllerRepositoryStruct) MockFindUserByEmailSuccess(email string) (models.User, error) {
	return models.User{
		Email:    "testingaaa@gmail.com",
		Password: "$2a$14$X7pxIBiQtS/SFhyOHo1aIO6PFTEY5.w2xHR84e.0nOi.kqwdiTylm",
	}, nil
}

func (fake *fakeControllerRepositoryStruct) MockFindUserByIdSuccess(userId string) (models.User, error) {
	return models.User{
		Email:      "testingaaa@gmail.com",
		Password:   "$2a$14$X7pxIBiQtS/SFhyOHo1aIO6PFTEY5.w2xHR84e.0nOi.kqwdiTylm",
		First_name: "aaa",
		Last_name:  "bbb",
	}, nil
}
