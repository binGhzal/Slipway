// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package plugins

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// CatalogEntry defines plugin metadata

type CatalogEntry struct {
	Source struct {
		OCI string `yaml:"oci"`
	} `yaml:"source"`
	Provides []struct {
		Capability string            `yaml:"capability"`
		Details    map[string]string `yaml:"details"`
	} `yaml:"provides"`
}

func LoadCatalog(path string) (map[string]CatalogEntry, error) {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	m := map[string]CatalogEntry{}
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}
