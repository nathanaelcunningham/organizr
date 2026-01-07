package fileutil

import "testing"

func TestParseTemplate(t *testing.T) {
	type args struct {
		template string
		vars     map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No series",
			args: args{
				template: "{author}/{title}",
				vars: map[string]string{
					"author": "author",
					"title":  "title",
				},
			},
			want: "author/title",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseTemplate(tt.args.template, tt.args.vars); got != tt.want {
				t.Errorf("ParseTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
