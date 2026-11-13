package instance

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

type vpcCapture struct {
	path    string
	body    string
	hits    int
	readErr error
}

func newVPCTestCommand(t *testing.T, cap *vpcCapture) *cobra.Command {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, errRead := io.ReadAll(r.Body)
		cap.readErr = errRead
		cap.hits++
		cap.path = r.URL.Path
		cap.body = string(body)
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)

	base := cli.NewCLIBase("vultr-cli/test")
	if err := base.Client.SetBaseURL(srv.URL); err != nil {
		t.Fatalf("set base url: %v", err)
	}
	base.Client.SetRetryLimit(0)

	cmd := NewCmdInstance(base)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	// Skip the API-key check so flag parsing and the VPC call can be observed locally.
	cmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		utils.SetOptions(base, c, args)
		base.Printer.Output = "text"
		return nil
	}

	return cmd
}

func TestInstanceVPCAttachDetach(t *testing.T) {
	cases := []struct {
		name string
		args []string
		path string
		body string
	}{
		{
			name: "attach flag",
			args: []string{"vpc", "attach", "inst-1", "--vpc-id=vpc-9"},
			path: "/v2/instances/inst-1/vpcs/attach",
			body: `"vpc_id":"vpc-9"`,
		},
		{
			name: "detach flag",
			args: []string{"vpc", "detach", "inst-1", "--vpc-id=vpc-9"},
			path: "/v2/instances/inst-1/vpcs/detach",
			body: `"vpc_id":"vpc-9"`,
		},
		{
			name: "attach positional",
			args: []string{"vpc", "attach", "inst-1", "vpc-9"},
			path: "/v2/instances/inst-1/vpcs/attach",
			body: `"vpc_id":"vpc-9"`,
		},
		{
			name: "detach positional",
			args: []string{"vpc", "detach", "inst-1", "vpc-9"},
			path: "/v2/instances/inst-1/vpcs/detach",
			body: `"vpc_id":"vpc-9"`,
		},
		{
			name: "flag overrides positional",
			args: []string{"vpc", "attach", "inst-1", "vpc-a", "--vpc-id=vpc-9"},
			path: "/v2/instances/inst-1/vpcs/attach",
			body: `"vpc_id":"vpc-9"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cap := &vpcCapture{}
			cmd := newVPCTestCommand(t, cap)
			cmd.SetArgs(tc.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
			if cap.readErr != nil {
				t.Fatalf("read body: %v", cap.readErr)
			}
			if cap.hits != 1 {
				t.Fatalf("api calls = %d, want 1", cap.hits)
			}
			if cap.path != tc.path {
				t.Fatalf("path = %q, want %q", cap.path, tc.path)
			}
			if !strings.Contains(cap.body, tc.body) {
				t.Fatalf("body = %q, want substring %q", cap.body, tc.body)
			}
		})
	}

	t.Run("missing vpc id", func(t *testing.T) {
		cap := &vpcCapture{}
		cmd := newVPCTestCommand(t, cap)
		cmd.SetArgs([]string{"vpc", "attach", "inst-1"})

		if err := cmd.Execute(); err == nil {
			t.Fatal("expected an error when no VPC ID is given")
		}
		if cap.hits != 0 {
			t.Fatalf("api calls = %d, want 0", cap.hits)
		}
	})
}

// The reported command passes --vpc-id and no API key. Flag parsing must succeed
// and stop at the existing API-key check, with no request.
func TestInstanceVPCDocumentedFlagParsesWithoutAPIKey(t *testing.T) {
	base := cli.NewCLIBase("vultr-cli/test")

	for _, action := range []string{"attach", "detach"} {
		t.Run(action, func(t *testing.T) {
			cmd := NewCmdInstance(base)
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true
			cmd.SetArgs([]string{"vpc", action, "demo-instance", "--vpc-id=demo-vpc"})

			err := cmd.Execute()
			if err == nil {
				t.Fatal("expected the API key check")
			}
			if strings.Contains(err.Error(), "unknown flag") {
				t.Fatalf("documented --vpc-id was rejected: %v", err)
			}
			if !strings.Contains(err.Error(), "VULTR_API_KEY") {
				t.Fatalf("error = %v, want the API key check", err)
			}
		})
	}
}
