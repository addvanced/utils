package slugger

import (
	"testing"
)

func TestSlugger_Slug(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		separator          string
		withEmoji          bool
		substitutions      map[string]string
		clearSubstitutions bool

		substitutionsChange       map[string]string
		addSubstitutionsChange    bool
		removeSubstitutionsChange bool

		expected string
	}{
		{
			name:      "Basic slug generation",
			input:     "Hello World",
			separator: "-",
			expected:  "hello-world",
		},
		{
			name:      "Basic slug generation using the default separator",
			input:     "Hello World",
			separator: "",
			expected:  "hello-world",
		},
		{
			name:      "Custom separator",
			input:     "Hello World",
			separator: "_",
			expected:  "hello_world",
		},
		{
			name:          "With substitutions",
			input:         "10% or 5€",
			separator:     "-",
			substitutions: map[string]string{"%": "percent", "€": "euro"},
			expected:      "10-percent-or-5-euro",
		},
		{
			name:      "With emoji replacement",
			input:     "Hello 🌍",
			separator: "-",
			withEmoji: true,
			expected:  "hello-globe-showing-europe-africa",
		},
		{
			name:      "Normalize to safe ASCII",
			input:     "Wôrķšpáçè ~~sèťtïñğš~~",
			separator: "-",
			expected:  "workspace-settings",
		},
		{
			name:                "Changing substitutions",
			input:               "10% or 5€",
			separator:           "-",
			substitutions:       map[string]string{"%": "percent", "€": "euro"},
			substitutionsChange: map[string]string{"%": "pct", "€": "eur"},
			expected:            "10-pct-or-5-eur",
		},
		{
			name:                   "Adding a new substitution",
			input:                  "Hello & World #HelloWorld",
			separator:              "-",
			substitutions:          map[string]string{"&": "and"},
			substitutionsChange:    map[string]string{"#": "hashtag"},
			addSubstitutionsChange: true,
			expected:               "hello-and-world-hashtaghelloworld",
		},
		{
			name:                      "Removing a substitution",
			input:                     "Hello & World #HelloWorld",
			separator:                 "-",
			substitutions:             map[string]string{"&": "and", "#": "hashtag"},
			substitutionsChange:       map[string]string{"#": "hashtag"},
			removeSubstitutionsChange: true,
			expected:                  "hello-and-world-helloworld",
		},
		{
			name:               "Clear all substitutions",
			input:              "Hello & World #HelloWorld",
			separator:          "-",
			substitutions:      map[string]string{"&": "and", "#": "hashtag"},
			clearSubstitutions: true,
			expected:           "hello-world-helloworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sl := New(tt.substitutions, tt.withEmoji)

			if tt.clearSubstitutions {
				sl.SetSubstitutions(nil)
			} else if tt.substitutionsChange != nil {
				if tt.addSubstitutionsChange {
					for key, value := range tt.substitutionsChange {
						sl.AddSubstitution(key, value)
					}
				} else if tt.removeSubstitutionsChange {
					for key := range tt.substitutionsChange {
						sl.RemoveSubstitution(key)
					}
				} else {
					sl.SetSubstitutions(tt.substitutionsChange)
				}
			}

			got := sl.Slug(tt.input, tt.separator)
			if got != tt.expected {
				t.Errorf("input: %q - expected %q, got %q", tt.input, tt.expected, got)
			}
		})
	}
}

func TestSlugger_Slug_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		separator     string
		substitutions map[string]string
		withEmoji     bool
		expected      string
	}{
		{
			name:      "Empty input",
			input:     "",
			separator: "-",
			expected:  "",
		},
		{
			name:      "Only whitespace collapses to empty",
			input:     " \t \n ",
			separator: "-",
			expected:  "",
		},
		{
			name:      "Tabs and newlines become single separators",
			input:     "Hello\tWorld\nTest",
			separator: "-",
			expected:  "hello-world-test",
		},
		{
			name:      "Keep ASCII safe set - _ .",
			input:     "A_b-C.d",
			separator: "-",
			expected:  "a_b-c.d",
		},
		{
			name:      "NFKD diacritics stripped including special characters",
			input:     "æ ø å Æ Ø Å ä ö Ä Ö ß",
			separator: " ",
			expected:  "ae oe a ae oe a a o a o ss",
		},
		{
			name:          "Substitutions are case-insensitive",
			input:         "10 % OR 5 €",
			separator:     "-",
			substitutions: map[string]string{"%": "percent", "€": "euro"},
			expected:      "10-percent-or-5-euro",
		},
		{
			name:          "Overlapping substitutions prefer longest (&& before &)",
			input:         "A && B & C",
			separator:     "-",
			substitutions: map[string]string{"&&": "andand", "&": "and"},
			expected:      "a-andand-b-and-c", // this will FAIL if shorter key is applied first
		},
		{
			name:      "Emoji ignored when WithEmoji=false",
			input:     "Hello 🌍",
			separator: "-",
			withEmoji: false,
			expected:  "hello",
		},
		{
			name:      "Custom multi-char separator",
			input:     "Hello   World",
			separator: "__",
			expected:  "hello__world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sl := New(tt.substitutions, tt.withEmoji)
			got := sl.Slug(tt.input, tt.separator)
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func BenchmarkSlugger_Slug(b *testing.B) {
	sl := New(map[string]string{"&": "and"}, false)

	for b.Loop() {
		sl.Slug("Wôrķšpáçè ~~sèťtïñğš~~", "")
	}
}

func BenchmarkSlugger_Slug_WithEmoji(b *testing.B) {
	sl := New(map[string]string{"&": "and"}, true)

	for b.Loop() {
		sl.Slug("a 😺, 🐈‍⬛, and a 🦁 go to 🏞️", "")
	}
}

func BenchmarkSlugger_Slug_CustomSeparator(b *testing.B) {
	sl := New(map[string]string{"&": "and"}, false)

	for b.Loop() {
		sl.Slug("Wôrķšpáçè ~~sèťtïñğš~~", "|")
	}
}
