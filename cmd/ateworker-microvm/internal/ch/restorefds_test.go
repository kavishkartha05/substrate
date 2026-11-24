// Copyright 2026 The Agent Substrate Authors
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

package ch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSnapshotNetDevices(t *testing.T) {
	const mac = "02:a8:1e:00:00:02"
	tests := []struct {
		name    string
		config  string
		want    []SnapshotNetDevice
		wantErr bool
	}{
		{
			name:   "one device with one queue pair",
			config: `{"net":[{"id":"_net1","num_queues":2,"mac":"` + mac + `"}]}`,
			want:   []SnapshotNetDevice{{ID: "_net1", QueuePairs: 1, MAC: mac}},
		},
		{
			name:    "two net devices",
			config:  `{"net":[{"id":"_net1","num_queues":2,"mac":"` + mac + `"},{"id":"_net2","num_queues":2,"mac":"` + mac + `"}]}`,
			wantErr: true,
		},
		{
			name:    "too many queues",
			config:  `{"net":[{"id":"_net1","num_queues":2048,"mac":"` + mac + `"}]}`,
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(tc.config), 0o600); err != nil {
				t.Fatal(err)
			}
			got, err := SnapshotNetDevices(dir)
			if (err != nil) != tc.wantErr {
				t.Fatalf("SnapshotNetDevices() err = %v, wantErr %t", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("SnapshotNetDevices() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
