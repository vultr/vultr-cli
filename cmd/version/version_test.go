package version

import (
	"reflect"
	"testing"
)

func TestNewCmdVersion(t *testing.T) {
	cmd := NewCmdVersion()

	if cmd.Short != "vultr-cli version" {
		t.Errorf("invalid short name")
	}

	if cmd.Use != "version" {
		t.Errorf("invalid use")
	}

	if !reflect.DeepEqual(cmd.Aliases, []string{"v"}) {
		t.Errorf("expected alias %v got %v", []string{"v"}, cmd.Aliases)
	}

}

func TestNewVersionOptions(t *testing.T) {
	version := NewVersionOptions()

	ref := reflect.TypeOf(version)
	if _, ok := ref.MethodByName("Get"); !ok {
		t.Errorf("missing Get")
	}

	vInterface := reflect.TypeOf(new(Interface)).Elem()
	if !ref.Implements(vInterface) {
		t.Errorf("VersionOptions does not implement Interface")
	}
}

func TestOptions_Get(t *testing.T) {
	version := NewVersionOptions()

	if reflect.TypeOf(version.Get()).Kind() != reflect.String {
		t.Errorf("response from get should be string")
	}
}
