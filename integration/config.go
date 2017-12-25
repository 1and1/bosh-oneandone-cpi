package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bosh-oneandone-cpi/action"
	boshapi "github.com/bosh-oneandone-cpi/api"
	boshdisp "github.com/bosh-oneandone-cpi/api/dispatcher"
	"github.com/bosh-oneandone-cpi/api/transport"
	boshcfg "github.com/bosh-oneandone-cpi/config"
	client "github.com/bosh-oneandone-cpi/oneandone/client"
	boshlogger "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/bosh-utils/uuid"
	"path/filepath"
)

var (
	// A stemcell that will be created in integration_suite_test.go
	existingStemcell string

	// Configurable defaults
	stemcellFile         = envOrDefault("STEMCELL_FILE", "")
	stemcellVersion      = envOrDefault("STEMCELL_VERSION", "https://s3.amazonaws.com/bosh-core-stemcells/google/bosh-stemcell-3468.15-google-kvm-ubuntu-trusty-go_agent.tgz")
	networkName          = envOrDefault("NETWORK_NAME", "cfintegration")
	customNetworkName    = envOrDefault("CUSTOM_NETWORK_NAME", "cfintegration-custom")
	customSubnetworkName = envOrDefault("CUSTOM_SUBNETWORK_NAME", "cfintegration-custom-us-central1")
	ipAddrs              = strings.Split(envOrDefault("PRIVATE_IP", "192.168.100.102,192.168.100.103,192.168.100.104"), ",")
	imageURL             = envOrDefault("IMAGE_URL", "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-1404-trusty-v20161213")

	// Channel that will be used to retrieve IPs to use
	ips chan string
	//apikeyPath = fakeAPIKeyPath()
	apikeyPath = "C:/gopath/src/github.com/bosh-oneandone-cpi/integration/test/assets/fake_api_key.pem"

	cfgContent = fmt.Sprintf(`{
	  "cloud": {
		"plugin": "oneandone",
          "properties": {
		  "oao": {
			"apikeyfile": "%v",
			"cpikeyfile": "sss",
			"cpiuser": "sss",
			"authorized_keys":{
				"cpi":"xxx",
				"user":"xxx"

			}
		  },
		  "agent": {
			"mbus": "http://127.0.0.1",
			"blobstore": {
			  "provider": "local"
			}
		  },
		  "registry": {
			"protocol": "http",
			"host": "fake-host",
			"port": 5555,
			"username": "fake-username",
			"password": "fake-username"
		  }
		}
      }
	}`, apikeyPath)
)

func execCPI(request string) (boshdisp.Response, error) {
	var err error
	var cfg boshcfg.Config
	var in, out, errOut, errOutLog bytes.Buffer
	var oaoClient client.Connector
	var boshResponse boshdisp.Response

	if cfg, err = boshcfg.NewConfigFromString(cfgContent); err != nil {
		return boshResponse, err
	}

	oaoClient = client.NewConnector(cfg.Cloud, boshlogger.NewLogger(boshlogger.LevelWarn))
	multiWriter := io.MultiWriter(&errOut, &errOutLog)
	logger := boshlogger.NewWriterLogger(boshlogger.LevelDebug, multiWriter)
	multiLogger := boshapi.MultiLogger{Logger: logger, LogBuff: &errOutLog}
	uuidGen := uuid.NewGenerator()

	actionFactory := action.NewConcreteFactory(
		oaoClient,
		uuidGen,
		cfg,
		multiLogger,
	)

	caller := boshdisp.NewJSONCaller()
	dispatcher := boshdisp.NewJSON(actionFactory, caller, multiLogger)

	in.WriteString(request)
	cli := transport.NewCLI(&in, &out, dispatcher, multiLogger)

	var response []byte

	if err = cli.ServeOnce(); err != nil {
		return boshResponse, err
	}

	if response, err = ioutil.ReadAll(&out); err != nil {
		return boshResponse, err
	}

	if err = json.Unmarshal(response, &boshResponse); err != nil {
		return boshResponse, err
	}
	return boshResponse, nil
}

func envRequired(key string) (val string) {
	if val = os.Getenv(key); val == "" {
		panic(fmt.Sprintf("Could not find required environment variable '%s'", key))
	}
	return
}

func envOrDefault(key, defaultVal string) (val string) {
	if val = os.Getenv(key); val == "" {
		val = defaultVal
	}
	return
}

func assetsDir() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "test/assets")
}

func fakeAPIKeyPath() string {
	return filepath.Join(assetsDir(), "fake_api_key.pem")
}