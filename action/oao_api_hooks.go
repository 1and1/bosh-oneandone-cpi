package action

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/stemcell"
	"github.com/bosh-oneandone-cpi/oneandone/vm"
	//"github.com/bosh-oneandone-cpi/oneandone/disks"
	//"github.com/bosh-oneandone-cpi/oneandone/resource"
)

var (
	stemcellFinderFactory    stemcell.FinderFactory    = stemcell.NewFinder
	stemcellCreatorFactory   stemcell.CreatorFactory   = stemcell.NewCreator
	stemcellDestroyerFactory stemcell.DestroyerFactory = stemcell.NewDestroyer

	vmFinderFactory     vm.FinderFactory     = vm.NewFinder
	vmCreatorFactory    vm.CreatorFactory    = vm.NewCreator
	vmTerminatorFactory vm.TerminatorFactory = vm.NewTerminator
	vmUpdaterFactory    vm.UpdaterFactory    = vm.NewUpdater

	//diskCreatorFactory      disks.CreatorFactory                  = disks.NewCreator
	//diskFinderFactory       disks.FinderFactory                   = disks.NewFinder
	//diskTerminatorFactory   disks.TerminatorFactory               = disks.NewTerminator
	//attacherDetacherFactory disks.InstanceAttacherDetacherFactory = disks.NewAttacherDetacherForInstance
)

func newStemcellFinder(c client.Connector, l boshlog.Logger) stemcell.Finder {
	return stemcellFinderFactory(c, l)
}

func newStemcellCreator(c client.Connector, l boshlog.Logger) stemcell.Creator {
	return stemcellCreatorFactory(c, l)
}

func newStemcellDestroyer(c client.Connector, l boshlog.Logger) stemcell.Destroyer {
	return stemcellDestroyerFactory(c, l)
}

func newVMFinder(c client.Connector, l boshlog.Logger) vm.Finder {
	return vmFinderFactory(c, l)
}

func newVMCreator(c client.Connector, l boshlog.Logger) vm.Creator {
	return vmCreatorFactory(c, l)
}

func newVMTerminator(c client.Connector, l boshlog.Logger) vm.Terminator {
	return vmTerminatorFactory(c, l)
}
func newVMUpdater(c client.Connector, l boshlog.Logger) vm.Updater {
	return vmUpdaterFactory(c, l)
}

//func newDiskCreator(c client.Connector, l boshlog.Logger) disks.Creator {
//	return diskCreatorFactory(c, l)
//}
//
//func newDiskFinder(c client.Connector, l boshlog.Logger) disks.Finder {
//	return diskFinderFactory(c, l)
//}
//
//func newDiskTerminator(c client.Connector, l boshlog.Logger) disks.Terminator {
//	return diskTerminatorFactory(c, l)
//}
//
//func newAttacherDetacherForInstance(in *resource.Instance, c client.Connector, l boshlog.Logger) (disks.AttacherDetacher, error) {
//	return attacherDetacherFactory(in, c, l)
//}

func installStemcellCreatorFactory(fac stemcell.CreatorFactory) {
	stemcellCreatorFactory = fac
}
func installStemcellDestroyerFactory(fac stemcell.DestroyerFactory) {
	stemcellDestroyerFactory = fac
}
func installStemcellFinderFactory(fac stemcell.FinderFactory) {
	stemcellFinderFactory = fac
}

func installVMCreatorFactory(fac vm.CreatorFactory) {
	vmCreatorFactory = fac
}
func installVMTerminatorFactory(fac vm.TerminatorFactory) {
	vmTerminatorFactory = fac
}
func installVMFinderFactory(fac vm.FinderFactory) {
	vmFinderFactory = fac
}
func installVMUpdaterFactory(fac vm.UpdaterFactory) {
	vmUpdaterFactory = fac
}

//func installDiskCreatorFactory(fac disks.CreatorFactory) {
//	diskCreatorFactory = fac
//}
//func installDiskFinderFactory(fac disks.FinderFactory) {
//	diskFinderFactory = fac
//}
//func installDiskTerminatorFactory(fac disks.TerminatorFactory) {
//	diskTerminatorFactory = fac
//}
//
//func installInstanceAttacherDetacherFactory(fac disks.InstanceAttacherDetacherFactory) {
//	attacherDetacherFactory = fac
//}

func resetAllFactories() {
	stemcellCreatorFactory = stemcell.NewCreator
	stemcellDestroyerFactory = stemcell.NewDestroyer
	stemcellFinderFactory = stemcell.NewFinder

	vmFinderFactory = vm.NewFinder
	vmCreatorFactory = vm.NewCreator
	vmTerminatorFactory = vm.NewTerminator
	vmUpdaterFactory = vm.NewUpdater

	//diskCreatorFactory = disks.NewCreator
	//diskFinderFactory = disks.NewFinder
	//diskTerminatorFactory = disks.NewTerminator
	//attacherDetacherFactory = disks.NewAttacherDetacherForInstance

}
