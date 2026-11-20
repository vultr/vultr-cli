package instance

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestOsChangeRequest(t *testing.T) {
	ignition := `{"ignition":{"version":"3.4.0"}}`
	encoded := base64.StdEncoding.EncodeToString([]byte(ignition))

	file := filepath.Join(t.TempDir(), "config.ign")
	if err := os.WriteFile(file, []byte(ignition), 0o600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		userData     string
		userDataFile string
		wantUserData *string
	}{
		{name: "without user data", wantUserData: nil},
		{name: "with user data text", userData: ignition, wantUserData: &encoded},
		{name: "with user data file", userDataFile: file, wantUserData: &encoded},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := osChangeRequest(391, test.userData, test.userDataFile)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if req.OsID != 391 {
				t.Errorf("OsID = %d, want 391", req.OsID)
			}
			switch {
			case test.wantUserData == nil && req.UserData != nil:
				t.Errorf("UserData = %q, want none", *req.UserData)
			case test.wantUserData != nil && (req.UserData == nil || *req.UserData != *test.wantUserData):
				t.Errorf("UserData = %v, want %q", req.UserData, *test.wantUserData)
			}
		})
	}

	t.Run("with a missing user data file", func(t *testing.T) {
		if _, err := osChangeRequest(391, "", filepath.Join(t.TempDir(), "missing.ign")); err == nil {
			t.Error("expected an error for a missing file")
		}
	})
}

func TestOsChangeHasUserDataFlags(t *testing.T) {
	osCmd, _, err := NewCmdInstance(nil).Find([]string{"os", "change"})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"os", "userdata", "userdata-file"} {
		if osCmd.Flags().Lookup(name) == nil {
			t.Errorf("instance os change has no --%s flag", name)
		}
	}
}
