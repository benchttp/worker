// Package firestoreconv offers a way to convert a Firestore protobuf document
// received on a Firestore Trigger event to a defined Go struct.
//
// The functionnality is limited to a simple conversion into internal.Report.
// Only the strictly required fields for the usecase are supported.
//
// firestoreconv exists to work around googleapis/google-cloud-go not offering
// a native way of converting Firestore Trigger data.
// See related issue https://github.com/googleapis/google-cloud-go/issues/1438
// for future solutions.
//
package firestoreconv

import (
	"errors"
	"fmt"
	"time"

	"github.com/googleapis/google-cloudevents-go/cloud/firestore/v1"

	"github.com/benchttp/worker/internal"
)

// ErrMapValueField is return when trying to access a field on
// a protobuf MapValue and the field is not present in the map.
var ErrMapValueField = errors.New("key is not in protobuf map value")

// ToReport converts a Firestore event payload to a usable internal.Report.
func ToReport(v *firestore.Value) (internal.Report, error) {
	r := internal.Report{}

	recordspb, ok := v.Fields["records"]
	if !ok {
		return internal.Report{}, fmt.Errorf(`%w: "%s"`, ErrMapValueField, "records")
	}

	for _, v := range recordspb.ArrayValue.Values {
		timepb, ok := v.MapValue.Fields["time"]
		if !ok {
			return internal.Report{}, fmt.Errorf(`%w: "%s"`, ErrMapValueField, "time")
		}
		r.Records = append(r.Records, internal.Record{Time: time.Duration(*timepb.IntegerValue)})
	}

	return r, nil
}
