// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package storage

import (
	"io"
	"net/url"
	"os"
	"path/filepath"

	"code.gitea.io/gitea/modules/util"
)

var (
	_ ObjectStorage = &LocalStorage{}
)

// LocalStorage represents a local files storage
type LocalStorage struct {
	dir string
}

// NewLocalStorage returns a local files
func NewLocalStorage(bucket string) (*LocalStorage, error) {
	if err := os.MkdirAll(bucket, os.ModePerm); err != nil {
		return nil, err
	}

	return &LocalStorage{
		dir: bucket,
	}, nil
}

// Open a file
func (l *LocalStorage) Open(path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(l.dir, path))
}

// Save a file
func (l *LocalStorage) Save(path string, r io.Reader) (int64, error) {
	p := filepath.Join(l.dir, path)
	if err := os.MkdirAll(filepath.Dir(p), os.ModePerm); err != nil {
		return 0, err
	}

	// always override
	if err := util.Remove(p); err != nil {
		return 0, err
	}

	f, err := os.Create(p)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return io.Copy(f, r)
}

// Delete delete a file
func (l *LocalStorage) Delete(path string) error {
	p := filepath.Join(l.dir, path)
	return util.Remove(p)
}

// URL gets the redirect URL to a file
func (l *LocalStorage) URL(path, name string) (*url.URL, error) {
	return nil, ErrURLNotSupported
}
