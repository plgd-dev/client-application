// ************************************************************************
// Copyright (C) 2022 plgd.dev, s.r.o.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ************************************************************************

package tar

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func makeTarFileFunc(src string, tw *tar.Writer) func(file string, fi os.FileInfo, err error) error {
	return func(file string, fi os.FileInfo, err error) error {
		// return on any error
		if err != nil {
			return err
		}
		// return on non-regular files (thanks to [kumo](https://medium.com/@komuw/just-like-you-did-fbdd7df829d3) for this suggested update)
		if !fi.Mode().IsRegular() {
			return nil
		}
		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.ReplaceAll(file, src, ""), string(filepath.Separator))
		// write the header
		if err = tw.WriteHeader(header); err != nil {
			return err
		}
		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			_ = f.Close()
			return err
		}
		return f.Close()
	}
}

// Tar takes a source and variable writers and walks 'source' writing each file
// found to the tar writer; the purpose for accepting multiple writers is to allow
// for multiple outputs (for example a file, or md5 hash)
func Tar(src string, writers ...io.Writer) (err error) {
	// ensure the src actually exists before trying to tar it
	if _, err = os.Stat(src); err != nil {
		return fmt.Errorf("unable to tar files: %w", err)
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer func() {
		errClose := gzw.Close()
		if err == nil {
			err = errClose
		}
	}()

	tw := tar.NewWriter(gzw)
	defer func() {
		errClose := tw.Close()
		if err == nil {
			err = errClose
		}
	}()

	// walk path
	err = filepath.Walk(src, makeTarFileFunc(src, tw))
	return
}

func copyFile(target string, perm fs.FileMode, tr io.Reader) error {
	f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, perm)
	if err != nil {
		return err
	}
	// copy over contents
	for {
		_, err := io.CopyN(f, tr, 1024)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			_ = f.Close()
			return err
		}
	}
	return f.Close()
}

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer func() {
		_ = gzr.Close()
	}()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		switch {
		// if no more files are found return
		case errors.Is(err, io.EOF):
			return nil
		// return any other error
		case err != nil:
			return err
		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}
		// Good: Check that path does not contain ".." before using it
		// If this check is ignored then there's a possibility of
		// arbitrary file write during zip extraction — "zip slip".
		// https://deepsource.io/gh/friendsofshopware/shopware-cli/issue/GSC-G305/occurrences
		if strings.Contains(header.Name, "..") {
			return fmt.Errorf("%s: target contains '..': security issue G305", header.Name)
		}
		target, err := filepath.Abs(filepath.Join(dst, header.Name)) //nolint:gosec
		if err != nil {
			return err
		}

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()
		// check the file type
		switch header.Typeflag {
		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0o755); err != nil {
					return err
				}
			}
		// if it's a file create it
		case tar.TypeReg:
			err := copyFile(target, os.FileMode(header.Mode), tr)
			if err != nil {
				return err
			}
		}
	}
}
