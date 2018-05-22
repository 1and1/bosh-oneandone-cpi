package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bosh-oneandone-cpi/action"
	boshapi "github.com/bosh-oneandone-cpi/api"
	boshdisp "github.com/bosh-oneandone-cpi/api/dispatcher"
	"github.com/bosh-oneandone-cpi/api/transport"
	boshcfg "github.com/bosh-oneandone-cpi/config"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	boshlogger "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/bosh-utils/uuid"
	"io"
	"io/ioutil"
	"os"
)

var oaoClient client.Connector

var (
	// A stemcell that will be created in integration_suite_test.go
	existingStemcell string
	token            = os.Getenv("ONEANDONE_TOKEN")
	sshKey           = os.Getenv("RSA_KEY")
	// registry
	registryUser     = envOrDefault("CPI_REGISTRY_USER", "admin")
	registryPassword = envOrDefault("CPI_REGISTRY_PASSWORD", "admin-password")
	registryHost     = envOrDefault("CPI_REGISTRY_ADDRESS", "10.4.92.139")
	registryPort     = envOrDefault("CPI_REGISTRY_PORT", "25777")

	cfgContent = fmt.Sprintf(`{
		  "cloud": {
			"plugin": "oneandone",
			  "properties": {
			  "oao": {
				"token": "%v"
			  },
			  "agent": {
				"mbus": "http://127.0.0.1",
				"blobstore": {
				  "provider": "local"
				}
			  },
			  "registry": {
				"protocol": "http",
				"host": "%v",
				"port": %v,
				"username": "%v",
				"password": "%v"
			  }
			}
		  }
		}`, token, registryHost, registryPort, registryUser, registryPassword)
)

func initAPI() error {
	var err error
	var cfg boshcfg.Config
	if cfg, err = boshcfg.NewConfigFromString(cfgContent); err != nil {
		return err
	}
	oaoClient = client.NewConnector(cfg.Cloud, boshlogger.NewLogger(boshlogger.LevelWarn))
	oaoClient.Connect()
	return nil
}

func execCPI(request string) (boshdisp.Response, error) {
	var err error
	var boshResponse boshdisp.Response
	var cfg boshcfg.Config

	var in, out, errOut, errOutLog bytes.Buffer
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

func envOrDefault(key, defaultVal string) (val string) {
	if val = os.Getenv(key); val == "" {
		val = defaultVal
	}
	return
}
