package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, "go-project-template", NAME)
	assert.Equal(t, "v1.0.0", SUMMARY)
}
