package integration

import (
	"fmt"

	"github.com/1and1/oneandone-cloudserver-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func assertSucceeds(request string) {
	response, err := execCPI(request)
	Expect(err).ToNot(HaveOccurred())
	Expect(response.Error).To(BeNil())
}

func assertFails(request string) error {
	response, _ := execCPI(request)
	Expect(response.Error).ToNot(BeNil())
	return response.Error
}

func assertSucceedsWithResult(request string) interface{} {
	response, err := execCPI(request)
	Expect(err).ToNot(HaveOccurred())
	Expect(response.Error).To(BeNil())
	Expect(response.Result).ToNot(BeNil())
	return response.Result
}

func toStringArray(raw []interface{}) []string {
	strings := make([]string, len(raw), len(raw))
	for i := range raw {
		strings[i] = raw[i].(string)
	}
	return strings
}

func assertValidVM(id string, valFunc func(server *oneandone.Server)) {

	server, err := oaoClient.Client().GetServer(id)
	if err != nil {
		//todo:throw error
	}
	valFunc(server)
	return
	Fail(fmt.Sprintf("Instance %q not found\n", id))
}
