// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package pkg

import "runtime"

// SupportMessageShort shows how to support arkade
var SupportMessageShort = supportMessage()

func supportMessage() string {
	if runtime.GOOS == "darwin" {
		return "slicervm.com: boot Linux microVMs directly on your Mac in <1s"
	}
	if runtime.GOOS == "linux" {
		return "slicervm.com: boot Linux microVMs instantly for AI agents, dev and e2e testing"
	}
	if runtime.GOOS == "windows" {
		return "slicervm.com: launch Linux microVMs in WSL2 for isolated dev and testing"
	}
	return ""
}
