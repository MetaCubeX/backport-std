// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !goexperiment.swissmap

package abi

// See comment in map_select_swiss.go.
type mapType = OldMapType

const SwissMap = false
const SwissMapInt = 0
