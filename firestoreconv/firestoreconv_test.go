package firestoreconv_test

import (
	"reflect"
	"testing"

	"github.com/googleapis/google-cloudevents-go/cloud/firestore/v1"

	"github.com/benchttp/worker/benchttp"
	"github.com/benchttp/worker/firestoreconv"
)

func TestToBenchmark(t *testing.T) {
	e := firestore.DocumentEventData{
		Value: &firestore.Value{
			Fields: map[string]firestore.OldValueField{
				"benchmark": {
					MapValue: &firestore.MapValue{
						Fields: map[string]firestore.MapValueField{
							"records": {
								ArrayValue: &firestore.ArrayValue{
									Values: []firestore.ValueElement{
										{
											MapValue: &firestore.MapValue{
												Fields: map[string]firestore.MapValueField{
													"time": {
														IntegerValue: newInt64(100),
													},
												},
											},
										},
										{
											MapValue: &firestore.MapValue{
												Fields: map[string]firestore.MapValueField{
													"time": {
														IntegerValue: newInt64(200),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	want := benchttp.Report{
		Benchmark: benchttp.Benchmark{
			Records: []benchttp.Record{
				{Time: 100},
				{Time: 200},
			},
		},
	}

	got, err := firestoreconv.Report(e.Value)
	if err != nil {
		t.Errorf("want nil error, got %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("incorrect conversion of firestore.DocumentEventData want %v, got %v", want, got)
	}
}

// newInt64 returns a pointer to the given int64 value.
func newInt64(x int64) *int64 {
	return &x
}
