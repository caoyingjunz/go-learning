/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cache

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// helper function to reduce stuttering
func testPop(f *DeltaFIFO) testFifoObject {
	return Pop(f).(Deltas).Newest().Object.(testFifoObject)
}

// literalListerGetter is a KeyListerGetter that is based on a
// function that returns a slice of objects to list and get.
// The function must list the same objects every time.
type literalListerGetter func() []testFifoObject

var _ KeyListerGetter = literalListerGetter(nil)

// ListKeys just calls kl.
func (kl literalListerGetter) ListKeys() []string {
	result := []string{}
	for _, fifoObj := range kl() {
		result = append(result, fifoObj.name)
	}
	return result
}

// GetByKey returns the key if it exists in the list returned by kl.
func (kl literalListerGetter) GetByKey(key string) (interface{}, bool, error) {
	for _, v := range kl() {
		if v.name == key {
			return v, true, nil
		}
	}
	return nil, false, nil
}

func TestDeltaFIFO_basic(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{KeyFunction: testFifoObjectKeyFunc})
	const amount = 500
	go func() {
		for i := 0; i < amount; i++ {
			f.Add(mkFifoObj(string([]rune{'a', rune(i)}), i+1))
		}
	}()
	go func() {
		for u := uint64(0); u < amount; u++ {
			f.Add(mkFifoObj(string([]rune{'b', rune(u)}), u+1))
		}
	}()

	lastInt := int(0)
	lastUint := uint64(0)
	for i := 0; i < amount*2; i++ {
		switch obj := testPop(f).val.(type) {
		case int:
			if obj <= lastInt {
				t.Errorf("got %v (int) out of order, last was %v", obj, lastInt)
			}
			lastInt = obj
		case uint64:
			if obj <= lastUint {
				t.Errorf("got %v (uint) out of order, last was %v", obj, lastUint)
			} else {
				lastUint = obj
			}
		default:
			t.Fatalf("unexpected type %#v", obj)
		}
	}
}

// TestDeltaFIFO_replaceWithDeleteDeltaIn tests that a `Sync` delta for an
// object `O` with ID `X` is added when .Replace is called and `O` is among the
// replacement objects even if the DeltaFIFO already stores in terminal position
// a delta of type `Delete` for ID `X`. Not adding the `Sync` delta causes
// SharedIndexInformers to miss `O`'s create notification, see https://github.com/kubernetes/kubernetes/issues/83810
// for more details.
func TestDeltaFIFO_replaceWithDeleteDeltaIn(t *testing.T) {
	oldObj := mkFifoObj("foo", 1)
	newObj := mkFifoObj("foo", 2)

	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: testFifoObjectKeyFunc,
		KnownObjects: literalListerGetter(func() []testFifoObject {
			return []testFifoObject{oldObj}
		}),
	})

	f.Delete(oldObj)
	f.Replace([]interface{}{newObj}, "")

	actualDeltas := Pop(f)
	expectedDeltas := Deltas{
		Delta{Type: Deleted, Object: oldObj},
		Delta{Type: Sync, Object: newObj},
	}
	if !reflect.DeepEqual(expectedDeltas, actualDeltas) {
		t.Errorf("expected %#v, got %#v", expectedDeltas, actualDeltas)
	}
}

func TestDeltaFIFO_requeueOnPop(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{KeyFunction: testFifoObjectKeyFunc})

	f.Add(mkFifoObj("foo", 10))
	_, err := f.Pop(func(obj interface{}) error {
		if obj.(Deltas)[0].Object.(testFifoObject).name != "foo" {
			t.Fatalf("unexpected object: %#v", obj)
		}
		return ErrRequeue{Err: nil}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok, err := f.GetByKey("foo"); !ok || err != nil {
		t.Fatalf("object should have been requeued: %t %v", ok, err)
	}

	_, err = f.Pop(func(obj interface{}) error {
		if obj.(Deltas)[0].Object.(testFifoObject).name != "foo" {
			t.Fatalf("unexpected object: %#v", obj)
		}
		return ErrRequeue{Err: fmt.Errorf("test error")}
	})
	if err == nil || err.Error() != "test error" {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok, err := f.GetByKey("foo"); !ok || err != nil {
		t.Fatalf("object should have been requeued: %t %v", ok, err)
	}

	_, err = f.Pop(func(obj interface{}) error {
		if obj.(Deltas)[0].Object.(testFifoObject).name != "foo" {
			t.Fatalf("unexpected object: %#v", obj)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok, err := f.GetByKey("foo"); ok || err != nil {
		t.Fatalf("object should have been removed: %t %v", ok, err)
	}
}

func TestDeltaFIFO_addUpdate(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{KeyFunction: testFifoObjectKeyFunc})
	f.Add(mkFifoObj("foo", 10))
	f.Update(mkFifoObj("foo", 12))
	f.Delete(mkFifoObj("foo", 15))

	if e, a := []interface{}{mkFifoObj("foo", 15)}, f.List(); !reflect.DeepEqual(e, a) {
		t.Errorf("Expected %+v, got %+v", e, a)
	}
	if e, a := []string{"foo"}, f.ListKeys(); !reflect.DeepEqual(e, a) {
		t.Errorf("Expected %+v, got %+v", e, a)
	}

	got := make(chan testFifoObject, 2)
	go func() {
		for {
			obj := testPop(f)
			t.Logf("got a thing %#v", obj)
			t.Logf("D len: %v", len(f.queue))
			got <- obj
		}
	}()

	first := <-got
	if e, a := 15, first.val; e != a {
		t.Errorf("Didn't get updated value (%v), got %v", e, a)
	}
	select {
	case unexpected := <-got:
		t.Errorf("Got second value %v", unexpected.val)
	case <-time.After(50 * time.Millisecond):
	}
	_, exists, _ := f.Get(mkFifoObj("foo", ""))
	if exists {
		t.Errorf("item did not get removed")
	}
}

func TestDeltaFIFO_enqueueingNoLister(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{KeyFunction: testFifoObjectKeyFunc})
	f.Add(mkFifoObj("foo", 10))
	f.Update(mkFifoObj("bar", 15))
	f.Add(mkFifoObj("qux", 17))
	f.Delete(mkFifoObj("qux", 18))

	// This delete does not enqueue anything because baz doesn't exist.
	f.Delete(mkFifoObj("baz", 20))

	expectList := []int{10, 15, 18}
	for _, expect := range expectList {
		if e, a := expect, testPop(f).val; e != a {
			t.Errorf("Didn't get updated value (%v), got %v", e, a)
		}
	}
	if e, a := 0, len(f.items); e != a {
		t.Errorf("queue unexpectedly not empty: %v != %v\n%#v", e, a, f.items)
	}
}

func TestDeltaFIFO_enqueueingWithLister(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: testFifoObjectKeyFunc,
		KnownObjects: literalListerGetter(func() []testFifoObject {
			return []testFifoObject{mkFifoObj("foo", 5), mkFifoObj("bar", 6), mkFifoObj("baz", 7)}
		}),
	})
	f.Add(mkFifoObj("foo", 10))
	f.Update(mkFifoObj("bar", 15))

	// This delete does enqueue the deletion, because "baz" is in the key lister.
	f.Delete(mkFifoObj("baz", 20))

	expectList := []int{10, 15, 20}
	for _, expect := range expectList {
		if e, a := expect, testPop(f).val; e != a {
			t.Errorf("Didn't get updated value (%v), got %v", e, a)
		}
	}
	if e, a := 0, len(f.items); e != a {
		t.Errorf("queue unexpectedly not empty: %v != %v", e, a)
	}
}

func TestDeltaFIFO_addReplace(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{KeyFunction: testFifoObjectKeyFunc})
	f.Add(mkFifoObj("foo", 10))
	f.Replace([]interface{}{mkFifoObj("foo", 15)}, "0")
	got := make(chan testFifoObject, 2)
	go func() {
		for {
			got <- testPop(f)
		}
	}()

	first := <-got
	if e, a := 15, first.val; e != a {
		t.Errorf("Didn't get updated value (%v), got %v", e, a)
	}
	select {
	case unexpected := <-got:
		t.Errorf("Got second value %v", unexpected.val)
	case <-time.After(50 * time.Millisecond):
	}
	_, exists, _ := f.Get(mkFifoObj("foo", ""))
	if exists {
		t.Errorf("item did not get removed")
	}
}

func TestDeltaFIFO_ResyncNonExisting(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: testFifoObjectKeyFunc,
		KnownObjects: literalListerGetter(func() []testFifoObject {
			return []testFifoObject{mkFifoObj("foo", 5)}
		}),
	})
	f.Delete(mkFifoObj("foo", 10))
	f.Resync()

	deltas := f.items["foo"]
	if len(deltas) != 1 {
		t.Fatalf("unexpected deltas length: %v", deltas)
	}
	if deltas[0].Type != Deleted {
		t.Errorf("unexpected delta: %v", deltas[0])
	}
}

func TestDeltaFIFO_Resync(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: testFifoObjectKeyFunc,
		KnownObjects: literalListerGetter(func() []testFifoObject {
			return []testFifoObject{mkFifoObj("foo", 5)}
		}),
	})
	f.Resync()

	deltas := f.items["foo"]
	if len(deltas) != 1 {
		t.Fatalf("unexpected deltas length: %v", deltas)
	}
	if deltas[0].Type != Sync {
		t.Errorf("unexpected delta: %v", deltas[0])
	}
}

func TestDeltaFIFO_DeleteExistingNonPropagated(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: testFifoObjectKeyFunc,
		KnownObjects: literalListerGetter(func() []testFifoObject {
			return []testFifoObject{}
		}),
	})
	f.Add(mkFifoObj("foo", 5))
	f.Delete(mkFifoObj("foo", 6))

	deltas := f.items["foo"]
	if len(deltas) != 2 {
		t.Fatalf("unexpected deltas length: %v", deltas)
	}
	if deltas[len(deltas)-1].Type != Deleted {
		t.Errorf("unexpected delta: %v", deltas[len(deltas)-1])
	}
}

func TestDeltaFIFO_ReplaceMakesDeletions(t *testing.T) {
	// We test with only one pre-existing object because there is no
	// promise about how their deletes are ordered.

	// Try it with a pre-existing Delete
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: testFifoObjectKeyFunc,
		KnownObjects: literalListerGetter(func() []testFifoObject {
			return []testFifoObject{mkFifoObj("foo", 5), mkFifoObj("bar", 6), mkFifoObj("baz", 7)}
		}),
	})
	f.Delete(mkFifoObj("baz", 10))
	f.Replace([]interface{}{mkFifoObj("foo", 5)}, "0")

	expectedList := []Deltas{
		{{Deleted, mkFifoObj("baz", 10)}},
		{{Sync, mkFifoObj("foo", 5)}},
		// Since "bar" didn't have a delete event and wasn't in the Replace list
		// it should get a tombstone key with the right Obj.
		{{Deleted, DeletedFinalStateUnknown{Key: "bar", Obj: mkFifoObj("bar", 6)}}},
	}

	for _, expected := range expectedList {
		cur := Pop(f).(Deltas)
		if e, a := expected, cur; !reflect.DeepEqual(e, a) {
			t.Errorf("Expected %#v, got %#v", e, a)
		}
	}

	// Now try starting with an Add instead of a Delete
	f = NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: testFifoObjectKeyFunc,
		KnownObjects: literalListerGetter(func() []testFifoObject {
			return []testFifoObject{mkFifoObj("foo", 5), mkFifoObj("bar", 6), mkFifoObj("baz", 7)}
		}),
	})
	f.Add(mkFifoObj("baz", 10))
	f.Replace([]interface{}{mkFifoObj("foo", 5)}, "0")

	expectedList = []Deltas{
		{{Added, mkFifoObj("baz", 10)},
			{Deleted, DeletedFinalStateUnknown{Key: "baz", Obj: mkFifoObj("baz", 7)}}},
		{{Sync, mkFifoObj("foo", 5)}},
		// Since "bar" didn't have a delete event and wasn't in the Replace list
		// it should get a tombstone key with the right Obj.
		{{Deleted, DeletedFinalStateUnknown{Key: "bar", Obj: mkFifoObj("bar", 6)}}},
	}

	for _, expected := range expectedList {
		cur := Pop(f).(Deltas)
		if e, a := expected, cur; !reflect.DeepEqual(e, a) {
			t.Errorf("Expected %#v, got %#v", e, a)
		}
	}

	// Now try starting without an explicit KeyListerGetter
	f = NewDeltaFIFOWithOptions(DeltaFIFOOptions{KeyFunction: testFifoObjectKeyFunc})
	f.Add(mkFifoObj("baz", 10))
	f.Replace([]interface{}{mkFifoObj("foo", 5)}, "0")

	expectedList = []Deltas{
		{{Added, mkFifoObj("baz", 10)},
			{Deleted, DeletedFinalStateUnknown{Key: "baz", Obj: mkFifoObj("baz", 10)}}},
		{{Sync, mkFifoObj("foo", 5)}},
	}

	for _, expected := range expectedList {
		cur := Pop(f).(Deltas)
		if e, a := expected, cur; !reflect.DeepEqual(e, a) {
			t.Errorf("Expected %#v, got %#v", e, a)
		}
	}
}

// TestDeltaFIFO_ReplaceMakesDeletionsReplaced is the same as the above test, but
// ensures that a Replaced DeltaType is emitted.
func TestDeltaFIFO_ReplaceMakesDeletionsReplaced(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: testFifoObjectKeyFunc,
		KnownObjects: literalListerGetter(func() []testFifoObject {
			return []testFifoObject{mkFifoObj("foo", 5), mkFifoObj("bar", 6), mkFifoObj("baz", 7)}
		}),
		EmitDeltaTypeReplaced: true,
	})

	f.Delete(mkFifoObj("baz", 10))
	f.Replace([]interface{}{mkFifoObj("foo", 6)}, "0")

	expectedList := []Deltas{
		{{Deleted, mkFifoObj("baz", 10)}},
		{{Replaced, mkFifoObj("foo", 6)}},
		// Since "bar" didn't have a delete event and wasn't in the Replace list
		// it should get a tombstone key with the right Obj.
		{{Deleted, DeletedFinalStateUnknown{Key: "bar", Obj: mkFifoObj("bar", 6)}}},
	}

	for _, expected := range expectedList {
		cur := Pop(f).(Deltas)
		if e, a := expected, cur; !reflect.DeepEqual(e, a) {
			t.Errorf("Expected %#v, got %#v", e, a)
		}
	}
}

// TestDeltaFIFO_ReplaceDeltaType checks that passing EmitDeltaTypeReplaced
// means that Replaced is correctly emitted.
func TestDeltaFIFO_ReplaceDeltaType(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: testFifoObjectKeyFunc,
		KnownObjects: literalListerGetter(func() []testFifoObject {
			return []testFifoObject{mkFifoObj("foo", 5)}
		}),
		EmitDeltaTypeReplaced: true,
	})
	f.Replace([]interface{}{mkFifoObj("foo", 5)}, "0")

	expectedList := []Deltas{
		{{Replaced, mkFifoObj("foo", 5)}},
	}

	for _, expected := range expectedList {
		cur := Pop(f).(Deltas)
		if e, a := expected, cur; !reflect.DeepEqual(e, a) {
			t.Errorf("Expected %#v, got %#v", e, a)
		}
	}
}

func TestDeltaFIFO_UpdateResyncRace(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: testFifoObjectKeyFunc,
		KnownObjects: literalListerGetter(func() []testFifoObject {
			return []testFifoObject{mkFifoObj("foo", 5)}
		}),
	})
	f.Update(mkFifoObj("foo", 6))
	f.Resync()

	expectedList := []Deltas{
		{{Updated, mkFifoObj("foo", 6)}},
	}

	for _, expected := range expectedList {
		cur := Pop(f).(Deltas)
		if e, a := expected, cur; !reflect.DeepEqual(e, a) {
			t.Errorf("Expected %#v, got %#v", e, a)
		}
	}
}

func TestDeltaFIFO_HasSyncedCorrectOnDeletion(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: testFifoObjectKeyFunc,
		KnownObjects: literalListerGetter(func() []testFifoObject {
			return []testFifoObject{mkFifoObj("foo", 5), mkFifoObj("bar", 6), mkFifoObj("baz", 7)}
		}),
	})
	f.Replace([]interface{}{mkFifoObj("foo", 5)}, "0")

	expectedList := []Deltas{
		{{Sync, mkFifoObj("foo", 5)}},
		// Since "bar" didn't have a delete event and wasn't in the Replace list
		// it should get a tombstone key with the right Obj.
		{{Deleted, DeletedFinalStateUnknown{Key: "bar", Obj: mkFifoObj("bar", 6)}}},
		{{Deleted, DeletedFinalStateUnknown{Key: "baz", Obj: mkFifoObj("baz", 7)}}},
	}

	for _, expected := range expectedList {
		if f.HasSynced() {
			t.Errorf("Expected HasSynced to be false")
		}
		cur := Pop(f).(Deltas)
		if e, a := expected, cur; !reflect.DeepEqual(e, a) {
			t.Errorf("Expected %#v, got %#v", e, a)
		}
	}
	if !f.HasSynced() {
		t.Errorf("Expected HasSynced to be true")
	}
}

func TestDeltaFIFO_detectLineJumpers(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{KeyFunction: testFifoObjectKeyFunc})

	f.Add(mkFifoObj("foo", 10))
	f.Add(mkFifoObj("bar", 1))
	f.Add(mkFifoObj("foo", 11))
	f.Add(mkFifoObj("foo", 13))
	f.Add(mkFifoObj("zab", 30))

	if e, a := 13, testPop(f).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}

	f.Add(mkFifoObj("foo", 14)) // ensure foo doesn't jump back in line

	if e, a := 1, testPop(f).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}

	if e, a := 30, testPop(f).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}

	if e, a := 14, testPop(f).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
}

func TestDeltaFIFO_addIfNotPresent(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{KeyFunction: testFifoObjectKeyFunc})

	emptyDeltas := Deltas{}
	if err := f.AddIfNotPresent(emptyDeltas); err == nil || !errors.Is(err, ErrZeroLengthDeltasObject) {
		t.Errorf("Expected error '%v', got %v", ErrZeroLengthDeltasObject, err)
	}

	f.Add(mkFifoObj("b", 3))
	b3 := Pop(f)
	f.Add(mkFifoObj("c", 4))
	c4 := Pop(f)
	if e, a := 0, len(f.items); e != a {
		t.Fatalf("Expected %v, got %v items in queue", e, a)
	}

	f.Add(mkFifoObj("a", 1))
	f.Add(mkFifoObj("b", 2))
	f.AddIfNotPresent(b3)
	f.AddIfNotPresent(c4)

	if e, a := 3, len(f.items); a != e {
		t.Fatalf("expected queue length %d, got %d", e, a)
	}

	expectedValues := []int{1, 2, 4}
	for _, expected := range expectedValues {
		if actual := testPop(f).val; actual != expected {
			t.Fatalf("expected value %d, got %d", expected, actual)
		}
	}
}

func TestDeltaFIFO_KeyOf(t *testing.T) {
	f := DeltaFIFO{keyFunc: testFifoObjectKeyFunc}

	table := []struct {
		obj interface{}
		key string
	}{
		{obj: testFifoObject{name: "A"}, key: "A"},
		{obj: DeletedFinalStateUnknown{Key: "B", Obj: nil}, key: "B"},
		{obj: Deltas{{Object: testFifoObject{name: "C"}}}, key: "C"},
		{obj: Deltas{{Object: DeletedFinalStateUnknown{Key: "D", Obj: nil}}}, key: "D"},
	}

	for _, item := range table {
		got, err := f.KeyOf(item.obj)
		if err != nil {
			t.Errorf("Unexpected error for %q: %v", item.obj, err)
			continue
		}
		if e, a := item.key, got; e != a {
			t.Errorf("Expected %v, got %v", e, a)
		}
	}
}

func TestDeltaFIFO_HasSynced(t *testing.T) {
	tests := []struct {
		actions        []func(f *DeltaFIFO)
		expectedSynced bool
	}{
		{
			actions:        []func(f *DeltaFIFO){},
			expectedSynced: false,
		},
		{
			actions: []func(f *DeltaFIFO){
				func(f *DeltaFIFO) { f.Add(mkFifoObj("a", 1)) },
			},
			expectedSynced: true,
		},
		{
			actions: []func(f *DeltaFIFO){
				func(f *DeltaFIFO) { f.Replace([]interface{}{}, "0") },
			},
			expectedSynced: true,
		},
		{
			actions: []func(f *DeltaFIFO){
				func(f *DeltaFIFO) { f.Replace([]interface{}{mkFifoObj("a", 1), mkFifoObj("b", 2)}, "0") },
			},
			expectedSynced: false,
		},
		{
			actions: []func(f *DeltaFIFO){
				func(f *DeltaFIFO) { f.Replace([]interface{}{mkFifoObj("a", 1), mkFifoObj("b", 2)}, "0") },
				func(f *DeltaFIFO) { Pop(f) },
			},
			expectedSynced: false,
		},
		{
			actions: []func(f *DeltaFIFO){
				func(f *DeltaFIFO) { f.Replace([]interface{}{mkFifoObj("a", 1), mkFifoObj("b", 2)}, "0") },
				func(f *DeltaFIFO) { Pop(f) },
				func(f *DeltaFIFO) { Pop(f) },
			},
			expectedSynced: true,
		},
		{
			// This test case won't happen in practice since a Reflector, the only producer for delta_fifo today, always passes a complete snapshot consistent in time;
			// there cannot be duplicate keys in the list or apiserver is broken.
			actions: []func(f *DeltaFIFO){
				func(f *DeltaFIFO) { f.Replace([]interface{}{mkFifoObj("a", 1), mkFifoObj("a", 2)}, "0") },
				func(f *DeltaFIFO) { Pop(f) },
			},
			expectedSynced: true,
		},
	}

	for i, test := range tests {
		f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{KeyFunction: testFifoObjectKeyFunc})

		for _, action := range test.actions {
			action(f)
		}
		if e, a := test.expectedSynced, f.HasSynced(); a != e {
			t.Errorf("test case %v failed, expected: %v , got %v", i, e, a)
		}
	}
}

// TestDeltaFIFO_PopShouldUnblockWhenClosed checks that any blocking Pop on an empty queue
// should unblock and return after Close is called.
func TestDeltaFIFO_PopShouldUnblockWhenClosed(t *testing.T) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KeyFunction: testFifoObjectKeyFunc,
		KnownObjects: literalListerGetter(func() []testFifoObject {
			return []testFifoObject{mkFifoObj("foo", 5)}
		}),
	})

	c := make(chan struct{})
	const jobs = 10
	for i := 0; i < jobs; i++ {
		go func() {
			f.Pop(func(obj interface{}) error {
				return nil
			})
			c <- struct{}{}
		}()
	}

	runtime.Gosched()
	f.Close()

	for i := 0; i < jobs; i++ {
		select {
		case <-c:
		case <-time.After(500 * time.Millisecond):
			t.Fatalf("timed out waiting for Pop to return after Close")
		}
	}
}

func BenchmarkDeltaFIFOListKeys(b *testing.B) {
	f := NewDeltaFIFOWithOptions(DeltaFIFOOptions{KeyFunction: testFifoObjectKeyFunc})
	const amount = 10000

	for i := 0; i < amount; i++ {
		f.Add(mkFifoObj(string([]rune{'a', rune(i)}), i+1))
	}
	for u := uint64(0); u < amount; u++ {
		f.Add(mkFifoObj(string([]rune{'b', rune(u)}), u+1))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = f.ListKeys()
		}
	})
	b.StopTimer()
}
