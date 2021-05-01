package config

import (
	"reflect"
	"testing"
)

func TestFromProjectTOML(t *testing.T) {
	for _, tt := range []struct {
		name string
		data string
		want Config
	}{
		{
			name: "Full sample",
			data: `# syntax = erichripko/cnbp
			[_]
			api = "0.2"
			id = "my-app"
			name = "My App"
			
			[io.buildpacks]
			api = "0.1"
			
			[io.buildpacks.build]
			builder = "some-builder"
			
			[io.buildpacks.ext.buildkit]
			previous-image = "prev-app-image"
			`,
			want: &config{
				builder:   "some-builder",
				prevImage: "prev-app-image",
			},
		}, {
			name: "Minimal sample",
			data: `# syntax = erichripko/cnbp
			[io.buildpacks.build]
			builder = "some-builder"
			`,
			want: &config{
				builder: "some-builder",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := FromProjectTOML(tt.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromProjectTOML() = %+v, want %+v", got, tt.want)
			}
		})
	}

	for _, tt := range []struct {
		name    string
		data    string
		wantErr string
	}{
		{
			name: "Builder is required",
			data: `# syntax = erichripko/cnbp
			`,
			wantErr: "no builder provided",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromProjectTOML(tt.data)
			if err.Error() != tt.wantErr {
				t.Errorf("FromProjectTOML() = %+v error = %v, wantErr = %v", got, err, tt.wantErr)
				return
			}
		})
	}
}
