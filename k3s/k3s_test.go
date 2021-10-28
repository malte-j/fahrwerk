package k3s

import (
	"strings"
	"testing"
)

func TestGenerateMasterScript(t *testing.T) {
	type args struct {
		config MasterScriptConfig
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test K3s version",
			args: args{config: MasterScriptConfig{
				K3sVersion: "1.0",
				K3sToken:   "alksnd1n010s92",
			}},
			want: "INSTALL_K3S_VERSION=\"1.0\"",
		},
		{
			name: "Test K3s Token",
			args: args{config: MasterScriptConfig{
				K3sVersion: "1.0",
				K3sToken:   "alksnd1n010s92",
			}},
			want: "K3S_TOKEN=\"alksnd1n010s92\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateMasterScript(tt.args.config); !strings.Contains(got, tt.want) {
				t.Errorf("GenerateMasterScript() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateK3sToken(t *testing.T) {
	t.Run("k3s token length", func(t *testing.T) {
		if got := GenerateK3sToken(); len(got) != 20 {
			t.Errorf("GenerateK3sToken() = %v, want %v", got, "length of 20")
		}
	})
}
