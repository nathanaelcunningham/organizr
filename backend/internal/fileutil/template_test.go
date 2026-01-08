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

func TestValidateTemplate(t *testing.T) {
	allowedVars := []string{"author", "series", "title"}

	tests := []struct {
		name      string
		template  string
		allowed   []string
		wantError bool
	}{
		{
			name:      "Valid template with all allowed placeholders",
			template:  "{author}/{series}/{title}",
			allowed:   allowedVars,
			wantError: false,
		},
		{
			name:      "Valid template with subset of allowed placeholders",
			template:  "{author}/{title}",
			allowed:   allowedVars,
			wantError: false,
		},
		{
			name:      "Invalid placeholder",
			template:  "{author}/{invalid}/{title}",
			allowed:   allowedVars,
			wantError: true,
		},
		{
			name:      "Empty template",
			template:  "",
			allowed:   allowedVars,
			wantError: false,
		},
		{
			name:      "Template with no placeholders",
			template:  "audiobooks/folder",
			allowed:   allowedVars,
			wantError: false,
		},
		{
			name:      "Mixed valid and invalid placeholders",
			template:  "{author}/{bad}/{title}",
			allowed:   allowedVars,
			wantError: true,
		},
		{
			name:      "Multiple invalid placeholders",
			template:  "{invalid1}/{invalid2}",
			allowed:   allowedVars,
			wantError: true,
		},
		{
			name:      "Valid single placeholder",
			template:  "{author}",
			allowed:   allowedVars,
			wantError: false,
		},
		{
			name:      "Valid template with series_number",
			template:  "{author}/{series}/{series_number} - {title}",
			allowed:   []string{"author", "series", "series_number", "title"},
			wantError: false,
		},
		{
			name:      "Invalid template with series_number but not in allowed",
			template:  "{author}/{invalid}/{series_number}",
			allowed:   allowedVars,
			wantError: true,
		},
		{
			name:      "Template with only series_number",
			template:  "{series_number}",
			allowed:   []string{"series_number"},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplate(tt.template, tt.allowed)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateTemplate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
