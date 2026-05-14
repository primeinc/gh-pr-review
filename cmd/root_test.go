package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// Tests follow spf13/cobra's canonical plugin pattern (TestPlugin /
// TestPluginWithSubCommands in cobra command_test.go) where the binary is
// installed as `kubectl-plugin` but invoked as `kubectl plugin`. Our setup
// is identical: `gh-pr-review` binary, `gh pr-review` invocation.

const wantDisplayName = "gh pr-review"

func TestRootCommand_DisplayName(t *testing.T) {
	root := newRootCommand()

	require.Equal(t, wantDisplayName, root.Annotations[cobra.CommandDisplayNameAnnotation])
	require.Equal(t, wantDisplayName, root.DisplayName())
}

func TestRootCommand_HelpOutput(t *testing.T) {
	root := newRootCommand()
	out := executeHelp(t, root, nil)

	require.Contains(t, out, wantDisplayName+" [command]")
	require.Contains(t, out, "help for "+wantDisplayName)
	require.Contains(t, out, `Use "`+wantDisplayName+` [command] --help"`)
}

func TestAllCommandsInTree_CommandPathStartsWithDisplayName(t *testing.T) {
	walkCommands(newRootCommand(), func(c *cobra.Command) {
		require.True(t, strings.HasPrefix(c.CommandPath(), wantDisplayName),
			"CommandPath %q must start with %q", c.CommandPath(), wantDisplayName)
	})
}

func TestAllCommandsInTree_HelpOutput(t *testing.T) {
	var paths [][]string
	walkCommands(newRootCommand(), func(c *cobra.Command) {
		if c.HasParent() {
			paths = append(paths, argvPath(c))
		}
	})

	for _, p := range paths {
		t.Run(strings.Join(p, " "), func(t *testing.T) {
			out := executeHelp(t, newRootCommand(), p)
			require.Contains(t, out, wantDisplayName,
				"help output for %q missing %q:\n%s", strings.Join(p, " "), wantDisplayName, out)
		})
	}
}

func walkCommands(c *cobra.Command, fn func(*cobra.Command)) {
	fn(c)
	for _, child := range c.Commands() {
		walkCommands(child, fn)
	}
}

func argvPath(c *cobra.Command) []string {
	var path []string
	for cur := c; cur.HasParent(); cur = cur.Parent() {
		path = append([]string{cur.Name()}, path...)
	}
	return path
}

func executeHelp(t *testing.T, root *cobra.Command, subPath []string) string {
	t.Helper()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(append(append([]string{}, subPath...), "--help"))
	require.NoError(t, root.Execute())
	return buf.String()
}
