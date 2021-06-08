package external_service

import (
	"github.com/pefish/go-test-assert"
	"testing"
)

func TestExternalServiceDriver_Register(t *testing.T) {
	ExternalServiceDriverInstance.Register("test", &TestExternalServiceInstance)

	ExternalServiceDriverInstance.Startup()

	svcs := ExternalServiceDriverInstance.ExternalServices()
	_, ok := svcs["test"]
	test.Equal(t, true, ok)
}