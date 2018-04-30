package fakes

import (
	"github.com/bosh-oneandone-cpi/oneandone/resource"
	sdk "github.com/1and1/oneandone-cloudserver-sdk-go"
)

type FakeAttacherDetacher struct {
	AttachVolumeCalled bool
	AttachedVolume     *sdk.BlockStorage
	AttachedInstance   *resource.Instance
	AttachmentError    error

	DetachVolumeCalled bool
	DetachedVolume     *sdk.BlockStorage
	DetachmentError    error
}

func (f *FakeAttacherDetacher) AttachInstanceToStorage(v *sdk.BlockStorage, in *resource.Instance) (string, error) {
	f.AttachVolumeCalled = true
	f.AttachedVolume = v
	f.AttachedInstance = in
	return "/dev/fake-path", f.AttachmentError
}

func (f *FakeAttacherDetacher) DetachInstanceFromStorage(v *sdk.BlockStorage, in *resource.Instance) error {
	f.DetachVolumeCalled = true
	f.DetachedVolume = v
	return f.DetachmentError
}
