/*
 * Copyright (c) 2022, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package global

import (
	"fmt"

	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/constraints"
)

const (
	// DefaultRetryTimes is 3 times
	DefaultRetryTimes = 3
	// DefaultRetryInterval is 5 seconds
	DefaultRetryInterval = time.Second * 5
	// DefaultTimeFormat is "2006-01-02 15:04:05 Z0700"
	DefaultTimeFormat = "2006-01-02 15:04:05 Z0700"
	// DefaultTimeZone is "UTC"
	DefaultTimeZone = "UTC"
	// DefaultProbeInterval is 1 minutes
	DefaultProbeInterval = time.Second * 60
	// DefaultTimeOut is 30 seconds
	DefaultTimeOut = time.Second * 30
)

// Retry is the settings of retry
type Retry struct {
	Times    int           `yaml:"times" json:"times,omitempty" jsonschema:"title=Retry Times,description=how many times need to retry,minimum=1"`
	Interval time.Duration `yaml:"interval" json:"interval,omitempty" jsonschema:"type=string,format=duration,title=Retry Interval,description=the interval between each retry"`
}

// The normalize() function logic as below:
// - if both global and local are not set, then return the _default.
// - if set the global, but not the local, then return the global
// - if set the local, but not the global, then return the local
// - if both global and local are set, then return the local
func normalize[T constraints.Ordered](global, local, valid, _default T) T {
	// if the val is invalid, then assign the default value
	if local <= valid {
		local = _default
		//if the global configuration is validated, assign the global
		if global > valid {
			local = global
		}
	}
	return local
}

// ErrNoRetry is the error need not retry
type ErrNoRetry struct {
	Message string
}

func (e *ErrNoRetry) Error() string {
	return e.Message
}

// DoRetry is a help function to retry the function if it returns error
func DoRetry(kind, name, tag string, r Retry, fn func() error) error {
	var err error
	for i := 0; i < r.Times; i++ {
		err = fn()
		_, ok := err.(*ErrNoRetry)
		if err == nil || ok {
			return err
		}
		log.Warnf("[%s / %s / %s] Retried to send %d/%d - %v", kind, name, tag, i+1, r.Times, err)

		// last time no need to sleep
		if i < r.Times-1 {
			time.Sleep(r.Interval)
		}
	}
	return fmt.Errorf("[%s / %s / %s] failed after %d retries - %v", kind, name, tag, r.Times, err)
}
