// Copyright 2015 The go-mxt Authors
// This file is part of the go-mxt library.
//
// The go-mxt library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-mxt library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-mxt library. If not, see <http://www.gnu.org/licenses/>.

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/mxt/go-mxt/metrics"
)

var (
	headerInMeter      = metrics.NewRegisteredMeter("mxt/downloader/headers/in", nil)
	headerReqTimer     = metrics.NewRegisteredTimer("mxt/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("mxt/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("mxt/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("mxt/downloader/bodies/in", nil)
	bodyReqTimer     = metrics.NewRegisteredTimer("mxt/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("mxt/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("mxt/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("mxt/downloader/receipts/in", nil)
	receiptReqTimer     = metrics.NewRegisteredTimer("mxt/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("mxt/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("mxt/downloader/receipts/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("mxt/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("mxt/downloader/states/drop", nil)

	throttleCounter = metrics.NewRegisteredCounter("mxt/downloader/throttle", nil)
)
