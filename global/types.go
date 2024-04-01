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

// "github.com/megaease/easeprobe/global"
// "github.com/megaease/easeprobe/probe"
)

// Format is the format of text
type Format int

// The format types
const (
	Unknown        Format = iota
	MarkdownSocial        // *text* is bold
	Markdown              // **text** is bold
	HTML
	JSON
	Text
	Log
	Slack
	Discord
	Lark
	SMS
	Shell
)
