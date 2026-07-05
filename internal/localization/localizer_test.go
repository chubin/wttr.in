package localization

import "testing"

func TestL10n_IsRTL(t *testing.T) {
	cases := []struct {
		lang string
		want bool
	}{
		{"he", true},
		{"ar", true},
		{"fa", true},
		{"en", false},
		{"de", false},
		{"ru", false},
		{"", false},
	}

	for _, tc := range cases {
		l10n := L10n{Lang: tc.lang}
		if got := l10n.IsRTL(); got != tc.want {
			t.Errorf("L10n{Lang: %q}.IsRTL() = %v, want %v", tc.lang, got, tc.want)
		}
	}
}
