package utils_test

import (
	"encoder/framework/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUtils_IsJson(t *testing.T) {

	json := `{
		"id": "1",
		"name": "Zelda",
		"email": "zelda@hyrule.com"
	}`

	err := utils.IsJson(json)
	require.Nil(t, err)

	json = "not a json"
	err = utils.IsJson(json)
	require.Error(t, err)
}
