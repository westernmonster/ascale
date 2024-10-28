// Copyright 2017 Bo-Yi Wu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package json

import jsoniter "github.com/json-iterator/go"

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	// Marshal is exported by vin/json package.
	Marshal = json.Marshal
	// Unmarshal is exported by vin/json package.
	Unmarshal = json.Unmarshal
	// MarshalIndent is exported by vin/json package.
	MarshalIndent = json.MarshalIndent
	// NewDecoder is exported by vin/json package.
	NewDecoder = json.NewDecoder
	// NewEncoder is exported by vin/json package.
	NewEncoder = json.NewEncoder
)
