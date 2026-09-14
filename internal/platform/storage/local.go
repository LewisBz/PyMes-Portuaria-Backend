package storage

import (
	"context"
	"os"
	"path/filepath"
)

type Local struct {
	Root string
}

func (l *Local) path(key string) string {
	return filepath.Join(l.Root, filepath.FromSlash(key))
}

func (l *Local) Put(_ context.Context, key string, data []byte) error {
	p := l.path(key)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

func (l *Local) Get(_ context.Context, key string) ([]byte, error) {
	return os.ReadFile(l.path(key))
}

func (l *Local) Delete(_ context.Context, key string) error {
	return os.Remove(l.path(key))
}
