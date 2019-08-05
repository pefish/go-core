package config

import (
	"fmt"
	"testing"
)

func TestConfigClass_LoadYamlConfig(t *testing.T) {
	Config.LoadYamlConfig(Configuration{
		ConfigFilepath: `/Users/joy/Work/backend/go-core/_example/config/local.yaml`,
		SecretFilepath: `/Users/joy/Work/backend/go-core/_example/secret/local1.yaml`,
	})
	a := struct {
		Host string `json:"host"`
	}{}
	Config.GetStruct(`mysql`, &a)
	fmt.Println(a)
}

func TestConfigClass_GetString(t *testing.T) {
	type fields struct {
		configs map[string]interface{}
	}
	type args struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: `test GetString`,
			fields: fields{
				map[string]interface{}{
					`test`: `haha`,
				},
			},
			args: args{
				`test`,
			},
			want: `haha`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ConfigClass{
				configs: tt.fields.configs,
			}
			if got := this.GetString(tt.args.str); got != tt.want {
				t.Errorf("ConfigClass.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}
