
/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package clock

import (
	"time"
)

type Clock struct {
	start, last time.Time
}

func New() *Clock {
	n := time.Now()
	return &Clock{
		start: n,
		last:  n,
	}
}

func (c *Clock) AllElapsed() time.Duration {
	return time.Now().Sub(c.start)
}
