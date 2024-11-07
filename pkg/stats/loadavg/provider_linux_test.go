//go:build linux

package loadavg

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    Value
		wantErr bool
	}{
		{
			name: "valid",
			s:    "1.67 1.69 1.82 2/3708 3162107",
			want: Value{
				One:     167,
				Five:    169,
				Fifteen: 182,
			},
			wantErr: false,
		},
		{
			name: "valid 2",
			s:    "1.67 1.69 1.82 2/3708ABCDE3162107",
			want: Value{
				One:     167,
				Five:    169,
				Fifteen: 182,
			},
			wantErr: false,
		},
		{
			name:    "invalid 1",
			s:       "1.67abc 1.69 1.82 2/3708 3162107",
			want:    Value{},
			wantErr: true,
		},
		{
			name:    "invalid 5",
			s:       "1.67 1.69abc 1.82 2/3708 3162107",
			want:    Value{},
			wantErr: true,
		},
		{
			name:    "invalid 15",
			s:       "1.67abc 1.69 1.82abc 2/3708 3162107",
			want:    Value{},
			wantErr: true,
		},
		{
			name:    "invalid num fields",
			s:       "1 2",
			want:    Value{},
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

func TestNewProvider(t *testing.T) {
	p := DataProvider{}
	if err := p.Available(); err != nil {
		t.Skip("not available")
	}

	_, err := p.Data()
	require.NoError(t, err, "must not have error")
}
