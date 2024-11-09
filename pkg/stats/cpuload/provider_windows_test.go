//go:build windows

package cpuload

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDataProvider(t *testing.T) {
	p := NewDataProvider()

	if err := p.Available(); err != nil {
		t.Skip("not available")
	}

	require.NotEqual(t, 0, p.prev.total, "total must be not zero")
}
