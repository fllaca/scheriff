package utils_test

import (
	"testing"

	"github.com/fllaca/okay/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestJoinNotEmptyStrings(t *testing.T) {

	assert.Equal(t, "a", utils.JoinNotEmptyStrings("/", "", "a"))

}
