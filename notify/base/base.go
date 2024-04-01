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

// Package base is the base implementation of the notification.
package base

import (
	"time"

	"wozaizhao.com/wzzapi/global"
	// "github.com/megaease/easeprobe/probe"
	log "github.com/sirupsen/logrus"
)

// DefaultNotify is the base struct of the Notify
type DefaultNotify struct {
	NotifyKind   string        `json:"-"`
	NotifyFormat global.Format `json:"-"`
	// NotifySendFunc func(string, string) error `json:"-"`
	NotifyName   string        `json:"name" jsonschema:"required,title=Notification Name,description=The name of the notification"`
	Dry          bool          `json:"dry,omitempty" jsonschema:"title=Dry Run,description=If true the notification will not send the message"`
	Timeout      time.Duration `json:"timeout,omitempty" jsonschema:"format=duration,title=Timeout,description=The timeout of the notification"`
	global.Retry `json:"retry,omitempty" jsonschema:"title=Retry,description=The retry of the notification"`
}

// Kind returns the kind of the notification
func (c *DefaultNotify) Kind() string {
	return c.NotifyKind
}

// Config is the default configuration for notification
func (c *DefaultNotify) Config(gConf global.NotifySettings) error {
	mode := "Live"
	if c.Dry {
		mode = "Dry"
	}
	log.Infof("Notification [%s] - [%s] is running on %s mode!", c.NotifyKind, c.NotifyName, mode)
	c.Timeout = gConf.NormalizeTimeOut(c.Timeout)
	// c.Retry = gConf.NormalizeRetry(c.Retry)

	log.Infof("Notification [%s] - [%s] is configured!", c.NotifyKind, c.NotifyName)
	return nil
}

// Name returns the name of the notification
func (c *DefaultNotify) Name() string {
	return c.NotifyName
}
