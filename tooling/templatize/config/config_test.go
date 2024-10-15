package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigProvider(t *testing.T) {
	region := "uksouth"
	user := "chiac"

	configProvider := NewConfigProvider("../testdata/config.yaml", region, user)
	ctx := context.Background()

	variables, err := configProvider.GetVariables(ctx, "public", "int")
	assert.NoError(t, err)
	assert.NotNil(t, variables)

	// key is not in the config file
	assert.Equal(t, "", variables["svc_resourcegroup"])

	// key is in the config file, region constant value
	assert.Equal(t, "uksouth", variables["test"])

	// key is in the config file, default in INT, constant value
	assert.Equal(t, "aro-hcp-int.azurecr.io/maestro-server:the-stable-one", variables["maestro_image"])

	// key is in the config file, default, varaible value
	assert.Equal(t, fmt.Sprintf("hcp-underlay-%s", region), variables["region_resourcegroup"])
}
