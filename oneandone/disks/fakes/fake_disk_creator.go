package fakes

import (
	sdk "github.com/oneandone/oneandone-cloudserver-sdk-go"
)

type FakeDiskCreator struct {
	CreateStorageCalled   bool
	CreateStorageResult   *sdk.BlockStorage
	CreateStorageSize     int
	CreateStorageError    error
}

func (f *FakeDiskCreator) CreateStorage(name string, sizeinGB int, dcId string) (*sdk.BlockStorage, error) {

	f.CreateStorageCalled = true
	f.CreateStorageSize = sizeinGB
	return f.CreateStorageResult, f.CreateStorageError

}
