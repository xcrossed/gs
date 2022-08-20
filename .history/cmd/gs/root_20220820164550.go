package gs

import (
	"fmt"
	"os"

	"github.com/go-spring/gs/internal"
	"github.com/spf13/cobra"
)

// 配置文件
var projectXml internal.ProjectXml

var rootCmd = &cobra.Command{
	Use:   "gs",
	Short: "gs - a simple CLI to transform and inspect strings",
	Long: `gs is a super fancy CLI (kidding)
   
One can use gs to add or modfiy go spring project from the terminal`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
