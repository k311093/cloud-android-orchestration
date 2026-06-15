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
	"net/http"
	"time"

	"github.com/google/cloud-android-orchestration/pkg/app/errors"
)

func WaitForHostReady(checker HostReadinessChecker, maxWait time.Duration, retryDelay time.Duration) error {
	deadline := time.Now().Add(maxWait)

	for time.Now().Before(deadline) {
		ready, err := checker.IsHostReady()
		if err != nil {
			return err
		}
		if ready {
			return nil
		}
		time.Sleep(retryDelay)
	}
	return errors.NewInternalError("wait for host orchestrator readiness timed out", nil)
}

type BasicHostReadinessChecker struct {
	Client HostClient
}

func (c *BasicHostReadinessChecker) IsHostReady() (bool, error) {
	status, err := c.Client.Get("/", "", nil)
	if err != nil || status == http.StatusBadGateway || status == http.StatusServiceUnavailable {
		return false, nil
	}
	return true, nil
}
