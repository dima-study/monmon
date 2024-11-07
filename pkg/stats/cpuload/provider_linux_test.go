//go:build linux

package cpuload

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDataProvider(t *testing.T) {
	p := NewDataProvider()

	require.NotEqual(t, 0, p.prev.total, "total must be not zero")
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    rawValue
		wantErr bool
	}{
		{
			name: "valid, no value",
			s:    "cpu  364963904 12712 107692855 1353061185 2762543 11910010 4125836 0 10 20",
			want: rawValue{
				user:   364963904 + 12712,
				system: 107692855,
				idle:   1353061185,
				total:  1844529045,
			},
			wantErr: false,
		},
		{
			name:    "invalid",
			s:       "cpu0 91781072 3183 26513368 336992676 675572 3894406 1193658 0 0 0",
			want:    rawValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.s)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Log(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
