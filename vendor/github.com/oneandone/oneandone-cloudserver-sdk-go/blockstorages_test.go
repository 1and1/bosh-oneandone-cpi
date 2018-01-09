package oneandone

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	set_bs         sync.Once
	set_bs_server  sync.Once
	test_bs_name   string
	test_bs_desc   string
	test_bs        *BlockStorage
	test_bs_server *Server
)

func setup_bs_server() {
	fmt.Println("Deploying a test server...")
	rand.Seed(time.Now().UnixNano())
	server_name = fmt.Sprintf("TestServer_%d", rand.Intn(1000000))
	fmt.Printf("Creating test server '%s'...\n", server_name)

	sap := get_random_appliance(hdd_size)
	ser_app_id = sap.Id
	mp := get_default_mon_policy()

	req := ServerRequest{
		Name:               server_name,
		Description:        server_name + " description",
		ApplianceId:        ser_app_id,
		MonitoringPolicyId: mp.Id,
		PowerOn:            true,
		DatacenterId:       "908DC2072407C94C8054610AD5A53B8C",
		Hardware: Hardware{
			Vcores:            v_cores,
			CoresPerProcessor: c_per_pr,
			Ram:               ram,
			Hdds: []Hdd{
				Hdd{
					Size:   hdd_size,
					IsMain: true,
				},
			},
		},
	}
	_, srv, err := api.CreateServer(&req)

	err = api.WaitForState(srv, "POWERED_ON", 10, 90)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	test_bs_server = srv
}

func create_block_storage() *BlockStorage {
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(999)
	test_bs_name = fmt.Sprintf("BlockStorage_%d", rint)
	test_bs_desc = fmt.Sprintf("BlockStorage_%d description", rint)
	req := BlockStorageRequest{
		Name:         test_bs_name,
		Description:  test_bs_desc,
		Size:         Int2Pointer(20),
		DatacenterId: "908DC2072407C94C8054610AD5A53B8C",
	}
	fmt.Printf("Creating new block storage '%s'...\n", test_bs_name)
	bs_id, bs, err := api.CreateBlockStorage(&req)

	if err != nil {
		fmt.Printf("Unable to create a block storage. Error: %s", err.Error())
		return nil
	}

	if bs_id == "" || bs.Id == "" {
		fmt.Printf("Unable to create block storage '%s'.", test_bs_name)
		return nil
	}

	api.WaitForState(bs, "POWERED_ON", 10, 30)
	return bs
}

func set_block_storage() {
	test_bs = create_block_storage()
}

func TestCreateBlockStorage(t *testing.T) {
	set_bs_server.Do(setup_bs_server)
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(999)
	test_bs_name = fmt.Sprintf("BlockStorage_%d", rint)
	test_bs_desc = fmt.Sprintf("BlockStorage_%d description", rint)
	req := BlockStorageRequest{
		Name:         test_bs_name,
		Description:  test_bs_desc,
		Size:         Int2Pointer(20),
		DatacenterId: "5091F6D8CBFEF9C26ACE957C652D5D49",
		ServerId:     test_bs_server.Id,
	}

	fmt.Printf("Creating new block storage '%s'...\n", test_bs_name)
	bs_id, bs, err := api.CreateBlockStorage(&req)
	if err != nil {
		t.Errorf("Unable to create a block storage. Error: %s", err.Error())
		t.Fail()
	}

	api.WaitForState(bs, "POWERED_ON", 10, 30)

	bs, _ = api.GetBlockStorage(bs_id)

	if bs.Server.Id != test_bs_server.Id {
		t.Errorf("Error while attaching a server to the block storage")
	}

	api.DeleteServer(test_bs_server.Id, false)
}

func TestListBlockStorages(t *testing.T) {
	set_bs.Do(set_block_storage)

	res, err := api.ListBlockStorages()
	if err != nil {
		t.Errorf("ListBlockStorages failed. Error: " + err.Error())
	}

	if len(res) == 0 {
		t.Errorf("No block storage found.")
	}
}

func TestGetBlockStorage(t *testing.T) {
	// set_bs_server.Do(setup_bs_server)
	set_bs.Do(set_block_storage)

	bs, err := api.GetBlockStorage(test_bs.Id)

	if err != nil {
		t.Errorf(err.Error())
	}

	test_bs = bs
}

func TestAddBlockStorageServer(t *testing.T) {
	set_bs.Do(set_block_storage)
	set_bs_server.Do(setup_bs_server)
	fmt.Printf("Adding server to block storage '%s'...\n", test_bs.Name)

	bs, err := api.AddBlockStorageServer(test_bs.Id, test_bs_server.Id)

	if err != nil {
		t.Errorf("AddBlockStorageServer failed. Error: " + err.Error())
		return
	}
	api.WaitForState(bs, "POWERED_ON", 10, 30)

	bs, _ = api.GetBlockStorage(bs.Id)

	fmt.Println(bs.Server)
	if bs.Server == nil {
		t.Errorf("Found no server added to the shared storage.")
	}

	bs, err = api.RemoveBlockStorageServer(bs.Id, test_bs_server.Id)

	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println("****")
	fmt.Println(bs)
	fmt.Println("****")

	if bs.Server != nil {
		t.Errorf("Server not removed from the shared storage.")
	}
	api.WaitForState(bs, "POWERED_ON", 10, 30)

	test_bs = bs
}

func TestDeleteBlockStorage(t *testing.T) {
	set_bs.Do(set_block_storage)
	// set_bs_server.Do(setup_bs_server)

	fmt.Printf("Deleting block storage '%s'...\n", test_bs.Name)
	bs, err := api.DeleteBlockStorage(test_bs.Id)
	if err != nil {
		t.Errorf("DeleteSharedStorage failed. Error: " + err.Error())
		return
	} else {
		api.WaitUntilDeleted(bs)
	}

	bs, err = api.GetBlockStorage(bs.Id)

	if bs != nil {
		t.Errorf("Unable to delete the block storage.")
	} else {
		test_bs = nil
	}

	if test_bs_server != nil {
		api.DeleteServer(test_bs_server.Id, false)
	}
}
