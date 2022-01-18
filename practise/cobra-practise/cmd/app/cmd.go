package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"go-learning/practise/cobra-practise/cmd/plugin"
)

// PluginHandler is capable of parsing command line arguments
// and performing executable filename lookups to search
// for valid plugin files, and execute found plugins.
type PluginHandler interface {
	// exists at the given filename, or a boolean false.
	// Lookup will iterate over a list of given prefixes
	// in order to recognize valid plugin filenames.
	// The first filepath to match a prefix is returned.
	Lookup(filename string) (string, bool)
	// Execute receives an executable's filepath, a slice
	// of arguments, and a slice of environment variables
	// to relay to the executable.
	Execute(executablePath string, cmdArgs, environment []string) error
}
type PixiuOptions struct {
	PluginHandler PluginHandler
	Arguments     []string
	ConfigFlags   *genericclioptions.ConfigFlags

	genericclioptions.IOStreams
}

// DefaultPluginHandler implements PluginHandler
type DefaultPluginHandler struct {
	ValidPrefixes []string
}

func NewDefaultPluginHandler(validPrefixes []string) *DefaultPluginHandler {
	return &DefaultPluginHandler{
		ValidPrefixes: validPrefixes,
	}
}

func (h *DefaultPluginHandler) Lookup(filename string) (string, bool) {
	// TODO
	return "", false
}

func (h *DefaultPluginHandler) Execute(executablePath string, cmdArgs, environment []string) error {
	return nil
}

func NewDefaultPixiuCommand() *cobra.Command {
	return NewDefaultPixiuctlCommandWithArgs(PixiuOptions{
		PluginHandler: NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes),
		Arguments:     os.Args,
		ConfigFlags:   genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag(),
		IOStreams:     genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
}

// NewDefaultPixiuCommand creates the `pixiuctl` command with default arguments
func NewDefaultPixiuctlCommandWithArgs(o PixiuOptions) *cobra.Command {
	cmd := NewPixiuCommand(o)

	if o.PluginHandler == nil {
		return cmd
	}

	// 预留命令行插件，暂时不做实现
	if len(o.Arguments) > 1 {
		cmdPathPieces := o.Arguments[1:]
		fmt.Println(cmdPathPieces)
	}

	return cmd
}

func NewPixiuCommand(o PixiuOptions) *cobra.Command {
	cmds := &cobra.Command{}

	return cmds
}
