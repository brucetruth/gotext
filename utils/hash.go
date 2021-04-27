/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package utils

import (
	"crypto/md5"
	"encoding/base64"
	"reflect"
	"unsafe"
)

//Hash returns the md5 hash token from any object
func Hash(i interface{}) string {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		if !v.CanAddr(){
			return ""
		}
		v = v.Addr()
	}

	size := unsafe.Sizeof(v.Interface())
	b    := (*[1 << 10]uint8)(unsafe.Pointer(v.Pointer()))[:size:size]

	h    := md5.New()
	return base64.StdEncoding.EncodeToString(h.Sum(b))
}
