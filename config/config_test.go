package config

import (
	"testing"
)

func TestConfigClass_Test(t *testing.T) {

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
