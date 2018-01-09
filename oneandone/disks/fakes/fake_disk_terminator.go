package fakes

type FakeDiskTerminator struct {
	DeleteStorageCalled bool
	DeleteStorageError  error
	DeleteStorageID     string
}

func (f *FakeDiskTerminator) DeleteStorage(storageId string) error {
	f.DeleteStorageCalled = true
	f.DeleteStorageID = storageId
	return f.DeleteStorageError
}
