package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestRootCommand_DisplayNameRendersAsGhExtension(t *testing.T) {
	root := newRootCommand()

	require.Equal(t, "gh pr-review", root.Annotations[cobra.CommandDisplayNameAnnotation],
		"root must set CommandDisplayNameAnnotation so help text reads 'gh pr-review'")
	require.Equal(t, "gh pr-review", root.DisplayName(),
		"DisplayName() must return the gh-extension invocation form, not the binary name")
}

func TestRootCommand_HelpOutputUsesSpacedForm(t *testing.T) {
	root := newRootCommand()

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"--help"})

	require.NoError(t, root.Execute())

	out := buf.String()
	require.Contains(t, out, "gh pr-review",
		"help output must show the user-facing invocation form")
	require.NotContains(t, out, "gh-pr-review",
		"help output must not leak the on-disk binary name (with hyphen)")
}

func TestRootCommand_SubcommandHelpUsesSpacedForm(t *testing.T) {
	for _, sub := range []string{"comments", "review", "threads"} {
		t.Run(sub, func(t *testing.T) {
			root := newRootCommand()

			var buf bytes.Buffer
			root.SetOut(&buf)
			root.SetErr(&buf)
			root.SetArgs([]string{sub, "--help"})

			require.NoError(t, root.Execute())

			out := buf.String()
			require.True(t, strings.Contains(out, "gh pr-review "+sub),
				"subcommand help must render parent path as 'gh pr-review %s', got:\n%s", sub, out)
			require.NotContains(t, out, "gh-pr-review",
				"subcommand help must not leak the on-disk binary name")
		})
	}
}
