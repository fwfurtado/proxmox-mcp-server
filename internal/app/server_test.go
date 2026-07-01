package app

import "testing"

func TestValidateTransport(t *testing.T) {
	tests := []struct {
		name      string
		transport string
		wantErr   bool
	}{
		{name: "stdio", transport: TransportStdio},
		{name: "streamable-http", transport: TransportStreamableHTTP},
		{name: "invalid", transport: "http", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTransport(tt.transport)
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}

func TestSafeURL(t *testing.T) {
	t.Run("removes credentials", func(t *testing.T) {
		got := safeURL("https://user:secret@example.com:8006/api2/json")
		want := "https://example.com:8006/api2/json"

		if got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})

	t.Run("keeps invalid url as-is", func(t *testing.T) {
		raw := "://bad url"
		if got := safeURL(raw); got != raw {
			t.Fatalf("expected %q, got %q", raw, got)
		}
	})
}
