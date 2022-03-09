// Package firestoreconv offers a way to convert a Firestore protobuf document
// received on a Firestore Trigger event to a defined Go struct.
//
// The functionnality is limited to a simple conversion into benchttp.Benchmark.
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

	"github.com/googleapis/google-cloudevents-go/cloud/firestore/v1"

	"github.com/benchttp/worker/benchttp"
)

// ErrMapValueField is return when trying to access a field on
// a protobuf MapValue and the field is not present in the map.
var ErrMapValueField = errors.New("key is not in protobuf map value")

// Report converts a Firestore event payload to a usable benchttp.Report.
func Report(v *firestore.Value) (benchttp.Report, error) {
	benchmark, ok := v.Fields["benchmark"]
	if !ok {
		return benchttp.Report{}, fmt.Errorf(`%w: "%s"`, ErrMapValueField, "benchmark")
	}

	recordField, ok := benchmark.MapValue.Fields["records"]
	if !ok {
		return benchttp.Report{}, fmt.Errorf(`%w: "%s"`, ErrMapValueField, "records")
	}

	records := make([]benchttp.Record, len(recordField.ArrayValue.Values))

	for i, v := range recordField.ArrayValue.Values {
		timeField, ok := v.MapValue.Fields["time"]
		if !ok {
			return benchttp.Report{}, fmt.Errorf(`%w: "%s"`, ErrMapValueField, "time")
		}

		codepb, ok := v.MapValue.Fields["code"]
		if !ok {
			return benchttp.Report{}, fmt.Errorf(`%w: "%s"`, ErrMapValueField, "code")
		}

		records[i] = benchttp.Record{
			Time: float64(*timeField.IntegerValue),
			Code: int(*codepb.IntegerValue),
		}
	}

	var report benchttp.Report
	report.Benchmark.Records = records

	return report, nil
}
