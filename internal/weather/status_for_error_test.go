package weather

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/chubin/wttr.in/internal/query"
)

func TestStatusForError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "unknown option", err: query.ErrUnknownOption, want: http.StatusBadRequest},
		{name: "wrapped unknown option", err: fmt.Errorf("%w: a", query.ErrUnknownOption), want: http.StatusBadRequest},
		{name: "invalid query string", err: query.ErrInvalidQueryString, want: http.StatusBadRequest},
		{name: "invalid integer", err: query.ErrInvalidInteger, want: http.StatusBadRequest},
		{name: "server error stays 500", err: ErrDataSource, want: http.StatusInternalServerError},
		{name: "wrapped server error stays 500", err: fmt.Errorf("location not found: %w", ErrDataSource), want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := statusForError(tt.err); got != tt.want {
				t.Fatalf("statusForError() = %d, want %d", got, tt.want)
			}
		})
	}
}
