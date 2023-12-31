// © 2023 Khulnasoft Limited All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

// ErrFailedToMarshalManifest is returned when we were unable to marshal the
// manifest to JSON
var ErrFailedToMarshalManifest = errors.New("failed to marshal manifest")

// ErrFailedToUnmarshalManifest is returned when we were unable to unmarshal the
// manifest from JSON
var ErrFailedToUnmarshalManifest = errors.New("failed to unmarshal manifest")

// Manifest contains metadata about the custom rules project.
type Manifest struct {
	Name string         `json:"name"`
	Push []ManifestPush `json:"push,omitempty"`
}

// ManifestPush contains metadata about where this rule bundle should be pushed
// to.  Currently this will always be the cloud API service.
type ManifestPush struct {
	CustomRulesID  string `json:"custom_rules_id,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
}

// copy creates a copy of the manifest so we don't accidentally modify the
// original.
func (m Manifest) copy() Manifest {
	cpy := m
	if m.Push != nil {
		cpy.Push = make([]ManifestPush, len(m.Push))
		copy(cpy.Push, m.Push)
	}
	return cpy
}

type manifestFile struct {
	*File
	manifest Manifest
}

func (m *manifestFile) WriteChanges(fsys afero.Fs) error {
	// This implementation is simpler if we just always update the manifest
	// file.
	b, err := json.MarshalIndent(m.manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToMarshalManifest, err)
	}
	m.File.UpdateContents(b)
	if err := m.File.WriteChanges(fsys); err != nil {
		return err
	}

	return nil
}

func (m *manifestFile) UpdateContents(manifest Manifest) {
	m.manifest = manifest
}

func manifestFromDir(fsys afero.Fs, root string) (*manifestFile, error) {
	path := filepath.Join(root, "manifest.json")
	file, err := FileFromPath(fsys, path)
	if err != nil {
		return nil, err
	}
	if !file.Exists() {
		return &manifestFile{File: file}, nil
	}
	b, err := afero.ReadFile(fsys, path)
	if err != nil {
		return nil, readPathError(path, err)
	}
	manifest := Manifest{}
	if err := json.Unmarshal(b, &manifest); err != nil {
		return nil, pathError(path, ErrFailedToUnmarshalManifest, err)
	}
	m := &manifestFile{
		File:     file,
		manifest: manifest,
	}
	return m, nil
}
