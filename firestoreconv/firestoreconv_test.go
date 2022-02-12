package firestoreconv_test

import (
	"reflect"
	"testing"

	"github.com/googleapis/google-cloudevents-go/cloud/firestore/v1"

	"github.com/benchttp/worker/firestoreconv"
	"github.com/benchttp/worker/internal"
)

func TestToReport(t *testing.T) {
	e := firestore.DocumentEventData{
		Value: &firestore.Value{
			Fields: map[string]firestore.OldValueField{
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
	}

	want := internal.Report{
		Records: []internal.Record{
			{Time: 100},
			{Time: 200},
		},
	}

	got, err := firestoreconv.ToReport(e.Value)
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
