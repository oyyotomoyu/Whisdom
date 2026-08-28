package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
)

var unsafeFilenameChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

// SanitizeFilename strips path separators and unusual characters so a
// client-supplied filename is safe to use as a single path segment.
func SanitizeFilename(name string) string {
	base := filepath.Base(name)
	cleaned := unsafeFilenameChars.ReplaceAllString(base, "_")
	if cleaned == "" || cleaned == "." || cleaned == ".." {
		return "upload"
	}
	return cleaned
}

// MaterialsDir returns the directory where a given material's original file
// is stored, under dataDir/materials/<materialID>/.
func MaterialsDir(dataDir, materialID string) string {
	return filepath.Join(dataDir, "materials", materialID)
}

// SaveMaterialFile persists an uploaded file under dataDir/materials/<id>/
// and returns the absolute path it was written to.
func SaveMaterialFile(dataDir, materialID, filename string, src multipart.File) (string, error) {
	dir := MaterialsDir(dataDir, materialID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create material directory: %w", err)
	}

	dest := filepath.Join(dir, SanitizeFilename(filename))
	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return "", fmt.Errorf("create material file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", fmt.Errorf("write material file: %w", err)
	}

	return dest, nil
}

// DeleteMaterialDir removes the on-disk directory for a material, including
// its original file.
func DeleteMaterialDir(dataDir, materialID string) error {
	return os.RemoveAll(MaterialsDir(dataDir, materialID))
}
