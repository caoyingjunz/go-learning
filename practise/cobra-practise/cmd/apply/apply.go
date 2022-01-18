package apply

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/dynamic"
	"k8s.io/kubectl/pkg/cmd/delete"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/openapi"
	"k8s.io/kubectl/pkg/util/templates"
	"k8s.io/kubectl/pkg/validation"
)

var (
	applyLong = templates.LongDesc(i18n.T(`
		Apply a configuration to a resource by file name or stdin.
		The resource name must be specified. This resource will be created if it doesn't exist yet.
		To use 'apply', always create the resource initially with either 'apply' or 'create --save-config'.

		JSON and YAML formats are accepted.

		Alpha Disclaimer: the --prune functionality is not yet complete. Do not use unless you are aware of what the current state is. See https://issues.k8s.io/34274.`))

	applyExample = templates.Examples(i18n.T(`
		# Apply the configuration in pod.json to a pod
		pixiuctl apply -f ./pod.json`))
)

// NewCmdApply create the `apply` command
func NewCmdApply(baseName string, f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	flags := NewApplyFlags(f, ioStreams)

	cmd := &cobra.Command{
		Use:                   "apply (-f FILENAME | -k DIRECTORY)",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Apply a configuration to a resource by file name or stdin"),
		Long:                  applyLong,
		Example:               applyExample,
		Run: func(cmd *cobra.Command, args []string) {
			o, err := flags.ToOptions(cmd, baseName, args)
			cmdutil.CheckErr(err)
			cmdutil.CheckErr(o.Run())
		},
	}

	return cmd
}

// ApplyFlags directly reflect the information that CLI is gathering via flags.  They will be converted to Options, which
// reflect the runtime requirements for the command.  This structure reduces the transformation to wiring and makes
// the logic itself easy to unit test
type ApplyFlags struct {
	Factory cmdutil.Factory

	RecordFlags *genericclioptions.RecordFlags
	PrintFlags  *genericclioptions.PrintFlags

	DeleteFlags *delete.DeleteFlags

	FieldManager   string
	Selector       string
	Prune          bool
	PruneResources []pruneResource
	All            bool
	Overwrite      bool
	OpenAPIPatch   bool
	PruneWhitelist []string

	genericclioptions.IOStreams
}

type pruneResource struct {
	group      string
	version    string
	kind       string
	namespaced bool
}

// NewApplyFlags returns a default ApplyFlags
func NewApplyFlags(f cmdutil.Factory, streams genericclioptions.IOStreams) *ApplyFlags {
	return &ApplyFlags{
		Factory:     f,
		RecordFlags: genericclioptions.NewRecordFlags(),
		DeleteFlags: delete.NewDeleteFlags("that contains the configuration to apply"),
		PrintFlags:  genericclioptions.NewPrintFlags("created").WithTypeSetter(scheme.Scheme),

		Overwrite:    true,
		OpenAPIPatch: true,

		IOStreams: streams,
	}
}

// ApplyOptions defines flags and other configuration parameters for the `apply` command
type ApplyOptions struct {
	Kubeconfig string

	Recorder genericclioptions.Recorder

	PrintFlags *genericclioptions.PrintFlags

	DeleteOptions *delete.DeleteOptions

	ServerSideApply bool
	ForceConflicts  bool
	FieldManager    string
	Selector        string
	DryRunStrategy  cmdutil.DryRunStrategy
	DryRunVerifier  *resource.DryRunVerifier
	Prune           bool
	PruneResources  []pruneResource
	cmdBaseName     string
	All             bool
	Overwrite       bool
	OpenAPIPatch    bool
	PruneWhitelist  []string

	Validator     validation.Schema
	Builder       *resource.Builder
	Mapper        meta.RESTMapper
	DynamicClient dynamic.Interface
	OpenAPISchema openapi.Resources

	Namespace        string
	EnforceNamespace bool

	genericclioptions.IOStreams

	// Objects (and some denormalized data) which are to be
	// applied. The standard way to fill in this structure
	// is by calling "GetObjects()", which will use the
	// resource builder if "objectsCached" is false. The other
	// way to set this field is to use "SetObjects()".
	// Subsequent calls to "GetObjects()" after setting would
	// not call the resource builder; only return the set objects.
	objects       []*resource.Info
	objectsCached bool

	// Stores visited objects/namespaces for later use
	// calculating the set of objects to prune.
	VisitedUids       sets.String
	VisitedNamespaces sets.String

	// Function run after the objects are generated and
	// stored in the "objects" field, but before the
	// apply is run on these objects.
	PreProcessorFn func() error
	// Function run after all objects have been applied.
	// The standard PostProcessorFn is "PrintAndPrunePostProcessor()".
	PostProcessorFn func() error
}

// ToOptions converts from CLI inputs to runtime inputs
func (flags *ApplyFlags) ToOptions(cmd *cobra.Command, baseName string, args []string) (*ApplyOptions, error) {
	s, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		panic(err)
	}

	o := &ApplyOptions{
		cmdBaseName: baseName,
		Kubeconfig:  s,
		PrintFlags:  flags.PrintFlags,

		Prune:          flags.Prune,
		PruneResources: flags.PruneResources,
		All:            flags.All,
		Overwrite:      flags.Overwrite,
		OpenAPIPatch:   flags.OpenAPIPatch,
		PruneWhitelist: flags.PruneWhitelist,

		IOStreams: flags.IOStreams,

		objects:       []*resource.Info{},
		objectsCached: false,

		VisitedUids:       sets.NewString(),
		VisitedNamespaces: sets.NewString(),
	}

	return o, nil
}

// Run executes the `apply` command.
func (o *ApplyOptions) Run() error {
	fmt.Println("TEST Kubeconfig", o.Kubeconfig)
	return nil
}
