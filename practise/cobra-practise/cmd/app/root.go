package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	author  string
	Verbose bool
	Source  string
	Region  string
)

var (
	rootCmd = &cobra.Command{
		Use:   "rootCmd [OPTIONS] [COMMANDS]",
		Short: "A short rootCmd demo",
		Long:  `A Long rootCmd demo`,

		// 不用也不能跟任何 positive args
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},

		// 以如下顺序执行: PersistentPreRun, PreRun, Run, PostRun, PersistentPostRun
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Run rootCmd PersistentPreRun with args: %v\n", args)
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("run rootCmd PreRun with args: %v\n", args)
		},

		// rootCmd 可以不实现，子命令来实现
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Verbose is:", Verbose)
			fmt.Println("author is:", author)
		},

		PostRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Run rootCmd PostRun with args: %v\n", args)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Run rootCmd PersistentPostRun with args: %v\n", args)
		},
	}
)

func Execute() error {
	//rootCmd.SuggestionsMinimumDistance = 1
	return rootCmd.Execute()
}

func init() {
	// 在 init() 函数中定义 flags 和 处理配置文件
	cobra.OnInitialize(initConfig)

	// Persistent Flags：全局 flag，指定命令和它的子命令均可用
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	// Local Flags: 只能指定给特定的 command
	rootCmd.Flags().StringVarP(&Source, "source", "s", "", "source directory to read from")

	// Bind Flags with Config（vipers）
	rootCmd.PersistentFlags().StringVarP(&author, "author", "a", "Caoyingjun", "Author name for the demo")
	viper.BindPFlag("auther", rootCmd.PersistentFlags().Lookup("auther"))

	// 全局必须的flag
	//rootCmd.PersistentFlags().StringVarP(&Region, "region", "r", "", "global region (required)")
	//rootCmd.MarkPersistentFlagRequired("region")

	// local 必须的flag
	//rootCmd.Flags().StringVarP(&Region, "region", "r", "", "local region (required)")
	//rootCmd.MarkFlagRequired("region")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
