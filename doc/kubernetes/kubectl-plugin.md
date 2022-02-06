# kubectl plugin 源码分析

### kubectl 概述

`kubectl` 作为 [kubernetes](https://github.com/kubernetes/kubernetes) 官方提供的命令行工具，基于 [cobra](https://github.com/spf13/cobra) 实现，用于对 `kubernetes` 集群进行管理

- 本文仅针对 `kubectl plugin` 的源码进行分析
- 使用请移步 [Extend kubectl with plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)
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
- 构造 `KubectlOptions` 结构体， 其中 `PluginHandler` 接口实现了 `Lookup` 和 `Execute` 方法，分别对 `plugin` 的 `查找` 和 `执行`；先按下不表，用到时在详细分析
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
  - [NewKubectlCommand](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/kubectl/pkg/cmd/cmd.go#L250) 会完成全部原生 `kubectl` 命令的构造；本文仅需关注子命令 `plugin`，用于获取 plugin 列表，后续展开分析。

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
- 通过 `o.Arguments` (原始os.Args) 判断是否执行 `plugin`， 如果 `是` 则直接执行 `plugin`，否则返回 cmds。判断逻辑：
  - 存在 `o.Arguments`
  - `command` 未在 cmds 中注册
  - `command` 的名称不为空
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

kubectl plugin 支持两种功能：`获取列表` 和 `执行`， 接下来将逐一分析

- 获取 plugin 列表 - plugin.NewCmdPlugin
- 执行 plugin - HandlePluginCommand
  ``` go
    func HandlePluginCommand(pluginHandler PluginHandler, cmdArgs []string) error {
        ...
        // attempt to find binary, starting at longest possible name with given cmdArgs
        for len(remainingArgs) > 0 {
            path, found := pluginHandler.Lookup(strings.Join(remainingArgs, "-"))
            ...
            foundBinaryPath = path
            break
        }

        if err := pluginHandler.Execute(foundBinaryPath, cmdArgs[len(remainingArgs):], os.Environ()); err != nil {
            return err
        }

        return nil
    }
  ```

### 获取 plugin 列表
获取 plugin 列表的接口，由原生 `kubectl` 提供，在 `plugin.NewCmdPlugin` 中实现，代码如下：
``` go
func NewCmdPlugin(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "plugin [flags]",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Provides utilities for interacting with plugins"),
		Long:                  pluginLong,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.DefaultSubCommandRun(streams.ErrOut)(cmd, args)
		},
	}

	cmd.AddCommand(NewCmdPluginList(streams))
	return cmd
}
```
`NewCmdPlugin` 会新建一个 plugin 的 cmd，追加 `NewCmdPluginList` 子命令； 获取 plugin 列表的功能就在 `NewCmdPluginList` 中实现

- NewCmdPluginList
  ``` go
    func NewCmdPluginList(streams genericclioptions.IOStreams) *cobra.Command {
        o := &PluginListOptions{ // 构造 `PluginListOptions`
            IOStreams: streams,
        }

        cmd := &cobra.Command{
            Use:     "list",
            Short:   i18n.T("List all visible plugin executables on a user's PATH"),
            Example: pluginExample,
            Long:    pluginListLong,
            Run: func(cmd *cobra.Command, args []string) {
                cmdutil.CheckErr(o.Complete(cmd))
                cmdutil.CheckErr(o.Run())
            },
        }

        cmd.Flags().BoolVar(&o.NameOnly, "name-only", o.NameOnly, "If true, display only the binary name of each plugin, rather than its full path")
        return cmd
    }
  ```
  - 构造 [PluginListOptions](https://github.com/kubernetes/kubernetes/blob/fbdd0d7b4165bc5a677d45e4dc693e3260297bfa/staging/src/k8s.io/kubectl/pkg/cmd/plugin/plugin.go#L77)，
    - `PluginListOptions`  实现 `Complete` 和 `Run` 方法，它们提供 `plugin list` 的实现
  - 执行 `o.Complete`
    ``` go
    func (o *PluginListOptions) Complete(cmd *cobra.Command) error {
        o.Verifier = &CommandOverrideVerifier{
            root:        cmd.Root(),
            seenPlugins: make(map[string]string),
        }

        o.PluginPaths = filepath.SplitList(os.Getenv("PATH"))
        return nil
    }
    ```
    `o.Complete` 很简洁，用于完成 `PluginListOptions` 的初始化, 主要初始化 `Verifier` 和 `PluginPaths`
    - Verifier 主要检验：
      - 是否为可执行文件
      - 是否可能被其他 `plugin` 覆盖
      - 是否被原生命令行覆盖
    - PluginPaths
      - 执行路径，run 函数会遍历 PluginPaths，去寻找符合 plugin 要求的文件

  - 执行 `o.Run`
    ``` go
    func (o *PluginListOptions) Run() error {
        ...
        for _, dir := range uniquePathsList(o.PluginPaths) {
            ...
            files, err := ioutil.ReadDir(dir)
            ...

            for _, f := range files {
                if f.IsDir() {
                    continue
                }
                if !hasValidPrefix(f.Name(), ValidPluginFilenamePrefixes) {
                    continue
                }

                if isFirstFile {
                    fmt.Fprintf(o.Out, "The following compatible plugins are available:\n\n")
                    pluginsFound = true
                    isFirstFile = false
                }

                pluginPath := f.Name()
                if !o.NameOnly {
                    pluginPath = filepath.Join(dir, pluginPath)
                }

                fmt.Fprintf(o.Out, "%s\n", pluginPath)
                if errs := o.Verifier.Verify(filepath.Join(dir, f.Name())); len(errs) != 0 {
                    for _, err := range errs {
                        fmt.Fprintf(o.ErrOut, "  - %s\n", err)
                        pluginWarnings++
                    }
                }
            }
        }
        ...
        return nil
    }
    ```
    [o.Run](https://github.com/kubernetes/kubernetes/blob/fbdd0d7b4165bc5a677d45e4dc693e3260297bfa/staging/src/k8s.io/kubectl/pkg/cmd/plugin/plugin.go#L117) 主要实现：
    - 遍历 `o.PluginPaths`，读目录下全部文件
    - 判断是否为文件
    - 判断是否为 `plugin` 文件，以 `kubectl-` 开头
    - 找到第一个 `plugin` 文件时，写入 `The following compatible plugins are available`
    - 将 `plugin` 文件写入标准输出
    - 判断是否可执行，如果不是则将信息写入错误输出

  - 执行效果
    ``` shell
    # kubectl plugin list
    The following compatible plugins are available:

    /usr/local/bin/kubectl-test
    ```

### 执行 plugin