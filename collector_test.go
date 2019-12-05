package monitoring

import (
	"reflect"
	"testing"
)

func TestNewCollector(t *testing.T) {
	fn := func() (tags map[string]string, fields map[string]interface{}) { return nil, nil }

	collector := NewCollector("name", fn)
	name := collector.GetName()

	if name != "name" {
		t.Fatalf("expected to match. want: %v, got: %v", "name", name)
	}
}

func TestCollector_GetName(t *testing.T) {
	collector := NewCollector("", nil)
	name := collector.GetName()
	if name != "" {
		t.Fatal("expected empty string")
	}

	collector = NewCollector("Name", nil)
	name = collector.GetName()

	if name != "Name" {
		t.Fatalf("expected: %v, got: %v", "Name", name)
	}
}

func TestCollector_Collect(t *testing.T) {
	var expectedTags map[string]string
	var expectedFields map[string]interface{}
	fn := func() (tags map[string]string, fields map[string]interface{}) { return expectedTags, expectedFields }
	collector := NewCollector("name", fn)
	name := collector.GetName()

	tags, fields := collector.Collect()
	if name != "name" {
		t.Fatalf("expected to match. want: %v, got: %v", "name", name)
	}
	if !reflect.DeepEqual(tags, expectedTags) {
		t.Fatalf("expected to match. want: %v, got: %v", expectedTags, tags)
	}
	if !reflect.DeepEqual(fields, expectedFields) {
		t.Fatalf("expected to match. want: %v, got: %v", expectedFields, fields)
	}
}
