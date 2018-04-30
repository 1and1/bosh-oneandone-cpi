package fakes

import (
	sdk "github.com/1and1/oneandone-cloudserver-sdk-go"
)

type FakeDiskFinder struct {
	FindStorageCalled bool
	FindStorageID     string
	FindStorageResult *sdk.BlockStorage
	FindStorageError  error

	FindAllAttachedStoragesCalled bool
	FindAllAttachedInstanceID    string
	FindAllAttachedResult        []sdk.BlockStorage
	FindAllAttachedError         error
}

func (f *FakeDiskFinder) FindStorage(storageId string) (*sdk.BlockStorage, error) {
	f.FindStorageCalled = true
	f.FindStorageID = storageId
	return f.FindStorageResult, f.FindStorageError
}

func (f *FakeDiskFinder) FindAllAttachedStorages(instanceID string) ([]sdk.BlockStorage, error) {
	f.FindAllAttachedStoragesCalled = true
	f.FindAllAttachedInstanceID = instanceID

	return f.FindAllAttachedResult, f.FindAllAttachedError
}
