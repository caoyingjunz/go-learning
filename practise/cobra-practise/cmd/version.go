package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var show bool

func init() {
	versionCmd.Flags().BoolVarP(&show, "show", "s", false, "verbose output")
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the demo version information",
	Long:  `A Long versionCmd demo`,

	// PersistentPreRun(global)： 子命令会继承 rootCmd
	// PreRun(local): 子命令不会继承 rootCmd
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//fmt.Printf("Run versionCmd PersistentPreRun with args: %v\n", args)
		fmt.Printf("")
	},
	//PreRun: func(cmd *cobra.Command, args []string) {
	//	fmt.Printf("run versionCmd PreRun with args: %v\n", args)
	//},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(show)
		fmt.Println("Cobra-demo version is v1.0.0")
	},
	// 定义自己的 PersistentPostRun 否则会运行 rootCmd 的 PersistentPostRun
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("")
	},
}
