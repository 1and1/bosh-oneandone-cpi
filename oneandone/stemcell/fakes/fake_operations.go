package fakes

type FakeFinder struct {
	FindStemcellCalled       bool
	FindStemcellCalledWithID string

	FindStemcellResult string
	FindStemcellError  error
}

type FakeCreator struct {
	FindStemcellCalled       bool
	FindStemcellCalledWithID string

	FindStemcellResult string
	FindStemcellError  error
}

type FakeDestroyer struct {
	DestroyStemcellCalled bool
	DestroyStemcellError  error
}

func (f *FakeCreator) CreateStemcell(imageId string) (stemcellId string, err error) {
	f.FindStemcellCalled = true
	f.FindStemcellCalledWithID = imageId
	f.FindStemcellResult = imageId
	return f.FindStemcellResult, f.FindStemcellError
}

func (f *FakeDestroyer) DeleteStemcell(stemcellId string) (err error) {
	f.DestroyStemcellCalled = true
	return f.DestroyStemcellError
}

func (f *FakeFinder) FindStemcell(imageOCID string) (stemcellId string, err error) {
	f.FindStemcellCalled = true
	f.FindStemcellCalledWithID = imageOCID
	f.FindStemcellResult = imageOCID
	return f.FindStemcellResult, f.FindStemcellError
}
