package cmd

import (
	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-bin/pkg/actions/tools/git"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:     "switch",
	Short:   "Switch branches",
	Run:     func(cmd *cobra.Command, args []string) {},
	GroupID: groups[group_main].ID,
}

func init() {
	carapace.Gen(switchCmd).Standalone()

	switchCmd.Flags().String("conflict", "", "conflict style (merge or diff3)")
	switchCmd.Flags().BoolP("create", "c", false, "create and switch to a new branch")
	switchCmd.Flags().BoolP("detach", "d", false, "detach HEAD at named commit")
	switchCmd.Flags().Bool("discard-changes", false, "throw away local modifications")
	switchCmd.Flags().BoolP("force", "f", false, "force checkout (throw away local modifications)")
	switchCmd.Flags().BoolP("force-create", "C", false, "create/reset and switch to a branch")
	switchCmd.Flags().Bool("guess", false, "second guess 'git switch <no-such-branch>'")
	switchCmd.Flags().Bool("ignore-other-worktrees", false, "do not check if another worktree is holding the given ref")
	switchCmd.Flags().BoolP("merge", "m", false, "perform a 3-way merge with the new branch")
	switchCmd.Flags().Bool("no-guess", false, "do not second guess 'git switch <no-such-branch>'")
	switchCmd.Flags().Bool("no-progress", false, "do not force progress reporting")
	switchCmd.Flags().Bool("no-track", false, "do not set upstream info for new branch")
	switchCmd.Flags().String("orphan", "", "new unparented branch")
	switchCmd.Flags().Bool("overwrite-ignore", false, "update ignored files (default)")
	switchCmd.Flags().Bool("progress", false, "force progress reporting")
	switchCmd.Flags().BoolP("quiet", "q", false, "suppress progress reporting")
	switchCmd.Flags().String("recurse-submodules", "", "control recursive updating of submodules")
	switchCmd.Flags().BoolP("track", "t", false, "set upstream info for new branch")
	rootCmd.AddCommand(switchCmd)

	switchCmd.Flag("recurse-submodules").NoOptDefVal = " "

	carapace.Gen(switchCmd).FlagCompletion(carapace.ActionMap{
		"conflict": carapace.ActionValues("merge", "diff3"),
		"orphan":   git.ActionRefs(git.RefOption{LocalBranches: true}),
	})

	carapace.Gen(switchCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			// from `man git-switch`:
			// git switch [<options>] (-c|-C) <new-branch> [<start-point>]
			if switchCmd.Flag("create").Changed || switchCmd.Flag("force-create").Changed {
				switch len(c.Args) {
				// no completions for <new-branch>
				case 0:
					return carapace.ActionValues()

				// complete local and remote branches for <start-point>
				case 1:
					return carapace.Batch(
						git.ActionRemoteBranches(""),
						git.ActionLocalBranches(),
					).ToA()
				}
			}

			if len(c.Args) == 0 {
				return carapace.Batch(
					git.ActionRemoteBranchNames(""),
					git.ActionRefs(git.RefOption{LocalBranches: true}),
				).ToA()
			}

			return carapace.ActionValues()
		}),
	)

	carapace.Gen(switchCmd).DashAnyCompletion(
		carapace.ActionPositional(switchCmd),
	)
}
