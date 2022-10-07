package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"go-learning/practise/cobra-practise/demo-server/app/options"
)

func NewDemoCommand() *cobra.Command {
	o, err := options.NewOptions()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}

	cmd := &cobra.Command{
		Use:  "demo-server",
		Long: `The demo server controller is a daemon than embeds the core control loops shipped with demo.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err = o.Complete(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			if err := Run(o); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	o.BindFlags(cmd)

	return cmd
}

func Run(c *options.Options) error {
	// 打印测试
	fmt.Println(c.ComponentConfig.Mysql)
	// 测试工厂函数
	fmt.Println(c.DBFactory.Test().Get("dd"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	StartDemoServer(ctx)

	select {}
}

func StartDemoServer(ctx context.Context) {
	go func(ctx context.Context) {
		t := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("接到中断信号，退出!")
				return
			case <-t.C:
				fmt.Println("demo")
			}
		}
	}(ctx)
}
