package gs

import (
	"github.com/spf13/cobra"
)

var backpupCmd = &cobra.Command{
	Use:     "backup",
	Aliases: []string{"bak"},
	Short:   "backup a project",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// res := stringer.Reverse(args[0])
		// fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(backpupCmd)
}
