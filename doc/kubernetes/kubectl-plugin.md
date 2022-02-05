# kubectl plugin 源码分析

### kubectl 概述

`kubectl` 作为 [kubernetes](https://github.com/kubernetes/kubernetes) 官方提供的命令行工具，基于 [cobra](https://github.com/spf13/cobra) 实现，用于对 `kubernetes` 集群进行管理

- 本文仅针对 `kubectl plugin` 的源码进行分析，如何使用请移步 [Extend kubectl with plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)
- `cobra` demo 请移步 [pixiuctl](https://github.com/caoyingjunz/go-learning/tree/master/practise/cobra-practise)

### kubectl 版本
- master

### 入口函数 main
``` go
// cmd/kubectl/kubectl.go

func main() {
	command := cmd.NewDefaultKubectlCommand() // 主干函数
	if err := cli.RunNoErrOutput(command); err != nil {
		// Pretty-print the error and exit with an error.
		util.CheckErr(err)
	}
}
```
`kubectl` 的 `main` 函数非常简洁，新建 command 后直接执行；command 由核心函数 `NewDefaultKubectlCommand` 返回，让我们一起看看它的真面目。

### NewDefaultKubectlCommand
``` go
// NewDefaultKubectlCommand creates the `kubectl` command with default arguments
func NewDefaultKubectlCommand() *cobra.Command {
	return NewDefaultKubectlCommandWithArgs(KubectlOptions{
		PluginHandler: NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes),
		Arguments:     os.Args,
		ConfigFlags:   defaultConfigFlags,
		IOStreams:     genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
}
```

`NewDefaultKubectlCommand` 主要作用：
- 构造 `KubectlOptions` 结构体， 其中 `PluginHandler` 接口实现了 `Lookup` 和 `Execute` 方法，分别对 `plugin` 的查找和执行；先按下不表，用到时在详细分析
    ``` go
    type PluginHandler interface {
	    Lookup(filename string) (string, bool)
        Execute(executablePath string, cmdArgs, environment []string) error
    }
    ```
- 初始化 `Arguments`, `ConfigFlags`, `IOStreams` 字段
- 通过 `NewDefaultKubectlCommandWithArgs` 方法构造 `*cobra.Command`

### NewDefaultKubectlCommandWithArgs
``` go
func NewDefaultKubectlCommandWithArgs(o KubectlOptions) *cobra.Command {
	cmd := NewKubectlCommand(o) // 1. 构造原生 kubectl 命令行

	if o.PluginHandler == nil {
		return cmd
	}

	if len(o.Arguments) > 1 {
		cmdPathPieces := o.Arguments[1:]

		if _, _, err := cmd.Find(cmdPathPieces); err != nil {
			var cmdName string // first "non-flag" arguments
			for _, arg := range cmdPathPieces {  // 2. 判断是否为 plugin，如果是，执行处理
				if !strings.HasPrefix(arg, "-") {
					cmdName = arg
					break
				}
			}

			switch cmdName {
			case "help", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
				// Don't search for a plugin
			default:
				if err := HandlePluginCommand(o.PluginHandler, cmdPathPieces); err != nil {
					fmt.Fprintf(o.IOStreams.ErrOut, "Error: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}

	return cmd
}
```

`NewDefaultKubectlCommandWithArgs` 是 `kubectl` 的核心方法, 主要完成两件事:
- 通过 `NewKubectlCommand` 方法完成原生 `kubectl` 命令行的构建
  - [NewKubectlCommand](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/kubectl/pkg/cmd/cmd.go#L250) 会完成全部原生 `kubectl` 命令的构造；本文仅关注子命令 `plugins`

      ``` go
      func NewKubectlCommand(o KubectlOptions) *cobra.Command {
        ...
        cmds := &cobra.Command{
                Use:   "kubectl",
                Short: i18n.T("kubectl controls the Kubernetes cluster manager"),
                Long: templates.LongDesc(`
            kubectl controls the Kubernetes cluster manager.`
            ...
            }

        // 增加 plugin 子命令
        cmds.AddCommand(plugin.NewCmdPlugin(o.IOStreams))

        return cmds
        }
     ```
- 通过 `o.Arguments` (os.Args) 判断是否执行 `plugin`， 如果是则直接执行 `plugin`，否则返回 cmds。判断逻辑：
  - 存在 `o.Arguments`
  - `command` 未在 cmds 中注册
  - `command` 的名称为 `o.Arguments` 中的第一个字符串不是以 `-` 开头
  - `command` 的名称不是 `help`, `__complete`, `__completeNoDesc`
  ``` go
  if len(o.Arguments) > 1 {
		cmdPathPieces := o.Arguments[1:]

		if _, _, err := cmd.Find(cmdPathPieces); err != nil {
			var cmdName string // first "non-flag" arguments
			for _, arg := range cmdPathPieces {  // 2. 判断是否为 plugin，如果是，执行处理
				if !strings.HasPrefix(arg, "-") {
					cmdName = arg
					break
				}
			}

			switch cmdName {
			case "help", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
				// Don't search for a plugin
			default:
				if err := HandlePluginCommand(o.PluginHandler, cmdPathPieces); err != nil {
					fmt.Fprintf(o.IOStreams.ErrOut, "Error: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}
  ```

- 如果最终判定为执行 `plugin` ，则调用 `HandlePluginCommand` 进行下一步处理

### 阶段性总结

