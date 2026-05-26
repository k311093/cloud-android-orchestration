// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package instances

import (
	"testing"
	"time"
)

type mockHostReadinessChecker struct {
	isHostReadyFunc func() (bool, error)
}

func (c *mockHostReadinessChecker) IsHostReady() (bool, error) {
	return c.isHostReadyFunc()
}

func TestWaitForHostReadySucceeds(t *testing.T) {
	checker := &mockHostReadinessChecker{
		isHostReadyFunc: func() (bool, error) {
			return true, nil
		},
	}

	err := WaitForHostReady(checker, 100*time.Millisecond, 1*time.Millisecond)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWaitForHostReadyRetries(t *testing.T) {
	calls := 0
	checker := &mockHostReadinessChecker{
		isHostReadyFunc: func() (bool, error) {
			calls++
			if calls == 1 {
				return false, nil
			}
			return true, nil
		},
	}

	err := WaitForHostReady(checker, 100*time.Millisecond, 1*time.Millisecond)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}
