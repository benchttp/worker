// Package firestoreconv offers a way to convert a Firestore protobuf document
// received on a Firestore Trigger event to a defined Go struct.
//
// The functionnality is limited to a simple conversion into internal.Benchmark.
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

// ToBenchmark converts a Firestore event payload to a usable internal.RunOutput.
func ToBenchmark(v *firestore.Value) (internal.Benchmark, error) {
	rs := []internal.Record{}

	// {"Value": {"Fields": {"report": ...} } }
	reportField, ok := v.Fields["report"]
	if !ok {
		return internal.Benchmark{}, fmt.Errorf(`%w: "%s"`, ErrMapValueField, "report")
	}

	// ...: {"MapValue": {"Fields": {"records": ...} } }
	recordField, ok := reportField.MapValue.Fields["records"]
	if !ok {
		return internal.Benchmark{}, fmt.Errorf(`%w: "%s"`, ErrMapValueField, "records")
	}

	// ...: {"ArrayValue": {"Values": [...] } }
	for _, v := range recordField.ArrayValue.Values {
		// ...: {"MapValue": {"Fields": {"time": {"IntegerValue": ... } } } }
		timeField, ok := v.MapValue.Fields["time"]
		if !ok {
			return internal.Benchmark{}, fmt.Errorf(`%w: "%s"`, ErrMapValueField, "time")
		}
		rs = append(rs, internal.Record{Time: time.Duration(*timeField.IntegerValue)})
	}

	var b internal.Benchmark
	b.Report.Records = rs

	return b, nil
}
