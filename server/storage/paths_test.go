package storage

import "testing"

func TestValidateDestinationPath(t *testing.T) {
	cases := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid", "/rag/sop/", false},
		{"valid no trailing slash", "/training/company-default", false},
		{"empty", "", true},
		{"blank", "   ", true},
		{"missing leading slash", "rag/sop/", true},
		{"parent traversal", "/rag/../secret/", true},
		{"null byte", "/rag/\x00sop/", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateDestinationPath(tc.path)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidateDestinationPath(%q) error = %v, wantErr %v", tc.path, err, tc.wantErr)
			}
		})
	}
}

func TestNormalizeDestinationPath(t *testing.T) {
	cases := map[string]string{
		"/rag/sop/":   "/rag/sop/",
		"/rag/sop":    "/rag/sop",
		"/rag//sop/":  "/rag/sop/",
		"/rag/./sop/": "/rag/sop/",
		"/":           "/",
	}
	for in, want := range cases {
		if got := NormalizeDestinationPath(in); got != want {
			t.Errorf("NormalizeDestinationPath(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSanitizeFilename(t *testing.T) {
	cases := map[string]string{
		"policy.pdf":            "policy.pdf",
		"../../etc/passwd":      "passwd",
		"my report (final).pdf": "my_report_final_.pdf",
		"..":                    "upload",
		"":                      "upload",
	}
	for in, want := range cases {
		if got := SanitizeFilename(in); got != want {
			t.Errorf("SanitizeFilename(%q) = %q, want %q", in, got, want)
		}
	}
}
