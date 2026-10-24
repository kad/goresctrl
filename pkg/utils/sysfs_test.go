// Copyright 2026 Intel Corporation. All Rights Reserved.
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

package utils

import (
	"os"
	"path/filepath"
	"testing"

	goresctrlpath "github.com/intel/goresctrl/pkg/path"
)

func setupUncoreDie(t *testing.T, vals map[string]string) {
	t.Helper()
	dir := t.TempDir()
	goresctrlpath.SetPrefix(dir)
	t.Cleanup(func() { goresctrlpath.SetPrefix("") })
	dieDir := filepath.Join(dir, SysfsUncoreBasepath, "package_00_die_00")
	if err := os.MkdirAll(dieDir, 0755); err != nil {
		t.Fatal(err)
	}
	for name, val := range vals {
		if err := os.WriteFile(filepath.Join(dieDir, name), []byte(val), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestUncoreFreqGettersAndReset(t *testing.T) {
	setupUncoreDie(t, map[string]string{
		"initial_min_freq_khz": "800000",
		"initial_max_freq_khz": "2500000",
		"min_freq_khz":         "2200000",
		"max_freq_khz":         "2300000",
	})

	if v, err := GetUncoreMinFreq(0, 0); err != nil || v != 2200000 {
		t.Errorf("GetUncoreMinFreq = %d, %v; want 2200000", v, err)
	}
	if v, err := GetUncoreMaxFreq(0, 0); err != nil || v != 2300000 {
		t.Errorf("GetUncoreMaxFreq = %d, %v; want 2300000", v, err)
	}
	if v, err := GetUncoreInitialMinFreq(0, 0); err != nil || v != 800000 {
		t.Errorf("GetUncoreInitialMinFreq = %d, %v; want 800000", v, err)
	}
	if v, err := GetUncoreInitialMaxFreq(0, 0); err != nil || v != 2500000 {
		t.Errorf("GetUncoreInitialMaxFreq = %d, %v; want 2500000", v, err)
	}

	if err := ResetUncoreFreq(0, 0); err != nil {
		t.Fatalf("ResetUncoreFreq: %v", err)
	}
	if v, _ := GetUncoreMinFreq(0, 0); v != 800000 {
		t.Errorf("min after reset = %d, want 800000", v)
	}
	if v, _ := GetUncoreMaxFreq(0, 0); v != 2500000 {
		t.Errorf("max after reset = %d, want 2500000", v)
	}
}

func TestUncoreFreqResetMissingDie(t *testing.T) {
	setupUncoreDie(t, nil)
	if err := ResetUncoreFreq(0, 0); err == nil {
		t.Error("ResetUncoreFreq on a die without attributes should fail")
	}
}
