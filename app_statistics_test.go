package monitoring

import (
	"reflect"
	"testing"
)

func TestAppStatistics_GetName(t *testing.T) {
	appStat := AppStatistics{}
	name := appStat.GetName()
	if name != "" {
		t.Fatal("expected empty string")
	}

	appStat = AppStatistics{Name: "Name"}
	name = appStat.GetName()
	if name != "Name" {
		t.Fatalf("expected: %v, got: %v", "Name", name)
	}
}

func TestAppStatistics_Collect(t *testing.T) {
	appStat := AppStatistics{}
	tags, fields := appStat.Collect()
	tags2, fields2 := appStat.Collect()

	if !reflect.DeepEqual(tags, tags2) {
		t.Fatalf("expected to match. want: %v, got: %v", tags, tags2)
	}

	if reflect.DeepEqual(fields, fields2) {
		t.Fatalf("expected to not match. want: %v, got: %v", fields, fields2)
	}
}
