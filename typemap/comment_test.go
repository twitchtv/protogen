package typemap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceComments(t *testing.T) {
	files := loadTestPb(t)
	file := protoFile(files, "service.proto")
	service := service(file, "ServiceWithManyMethods")

	comments, err := ServiceComments(file, service)
	require.NoError(t, err, "unable to load service comments")
	assert.Equal(t, " ServiceWithManyMethods leading\n", comments.Leading)
}

func TestMethodComments(t *testing.T) {
	files := loadTestPb(t)
	file := protoFile(files, "service.proto")
	service := service(file, "ServiceWithManyMethods")
	method1 := method(service, "Method1")

	comments, err := MethodComments(file, service, method1)
	require.NoError(t, err, "unable to load method comments")
	assert.Equal(t, " Method1 leading\n", comments.Leading)
}
