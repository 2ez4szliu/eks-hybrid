package init

import (
	"github.com/integrii/flaggy"
	"go.uber.org/zap"
	"k8s.io/utils/strings/slices"

	"github.com/aws/eks-hybrid/internal/api"
	"github.com/aws/eks-hybrid/internal/cli"
	"github.com/aws/eks-hybrid/internal/configenricher"
	"github.com/aws/eks-hybrid/internal/configprovider"
	"github.com/aws/eks-hybrid/internal/containerd"
	"github.com/aws/eks-hybrid/internal/daemon"
	"github.com/aws/eks-hybrid/internal/kubelet"
	"github.com/aws/eks-hybrid/internal/ssm"
	"github.com/aws/eks-hybrid/internal/system"
)

const (
	configPhase = "config"
	runPhase    = "run"
)

func NewInitCommand() cli.Command {
	init := initCmd{}
	init.cmd = flaggy.NewSubcommand("init")
	init.cmd.StringSlice(&init.daemons, "d", "daemon", "specify one or more of `containerd` and `kubelet`. This is intended for testing and should not be used in a production environment.")
	init.cmd.StringSlice(&init.skipPhases, "s", "skip", "phases of the bootstrap you want to skip")
	init.cmd.Description = "Initialize this instance as a node in an EKS cluster"
	return &init
}

type initCmd struct {
	cmd        *flaggy.Subcommand
	skipPhases []string
	daemons    []string
}

func (c *initCmd) Flaggy() *flaggy.Subcommand {
	return c.cmd
}

func (c *initCmd) Run(log *zap.Logger, opts *cli.GlobalOptions) error {
	log.Info("Checking user is root..")
	root, err := cli.IsRunningAsRoot()
	if err != nil {
		return err
	} else if !root {
		return cli.ErrMustRunAsRoot
	}

	log.Info("Loading configuration..", zap.String("configSource", opts.ConfigSource))
	provider, err := configprovider.BuildConfigProvider(opts.ConfigSource)
	if err != nil {
		return err
	}
	nodeConfig, err := provider.Provide()
	if err != nil {
		return err
	}
	log.Info("Loaded configuration", zap.Reflect("config", nodeConfig))

	zap.L().Info("Validating configuration..")
	v := api.NewValidator(nodeConfig)
	if err := v.Validate(nodeConfig); err != nil {
		return err
	}

	log.Info("Enriching configuration..")
	enricher := configenricher.New(log, nodeConfig)
	if err := enricher.Enrich(nodeConfig); err != nil {
		return err
	}

	log.Info("Creating daemon manager..")
	daemonManager, err := daemon.NewDaemonManager()
	if err != nil {
		return err
	}
	defer daemonManager.Close()

	aspects := []system.SystemAspect{
		system.NewLocalDiskAspect(),
		system.NewNetworkingAspect(),
	}

	var daemons []daemon.Daemon
	// If Hybrid w/ SSM is enabled, we need to make sure SSM daemon is configured first
	// in order to register the instance first. This will provide us both aws credentials
	// and managed instance ID which will override hostname in both kubelet configs & provider-id
	if nodeConfig.IsSSM() {
		daemons = append(daemons, ssm.NewSsmDaemon(daemonManager))
	}

	daemons = append(daemons,
		containerd.NewContainerdDaemon(daemonManager),
		kubelet.NewKubeletDaemon(daemonManager),
	)

	if !slices.Contains(c.skipPhases, configPhase) {
		log.Info("Configuring daemons...")
		for _, daemon := range daemons {
			if len(c.daemons) > 0 && !slices.Contains(c.daemons, daemon.Name()) {
				continue
			}
			nameField := zap.String("name", daemon.Name())

			log.Info("Configuring daemon...", nameField)
			if err := daemon.Configure(nodeConfig); err != nil {
				return err
			}
			log.Info("Configured daemon", nameField)

			// Check if SSM daemon and set node name
			if daemon.Name() == ssm.SsmDaemonName {
				registeredNodeName, err := ssm.GetManagedHybridInstanceId()
				if err != nil {
					return err
				}
				nodeNameField := zap.String("registered instance-id", registeredNodeName)
				log.Info("Re-setting node name with registered managed instance id", nodeNameField)
				nodeConfig.Spec.Hybrid.NodeName = registeredNodeName
			}
		}
	}

	if !slices.Contains(c.skipPhases, runPhase) {
		// Aspects are not required for hybrid nodes
		// Setting up aspects fall under user responsibility for hybrid nodes
		if !nodeConfig.IsHybridNode() {
			log.Info("Setting up system aspects...")
			for _, aspect := range aspects {
				nameField := zap.String("name", aspect.Name())
				log.Info("Setting up system aspect..", nameField)
				if err := aspect.Setup(nodeConfig); err != nil {
					return err
				}
				log.Info("Set up system aspect", nameField)
			}
		}
		for _, daemon := range daemons {
			if len(c.daemons) > 0 && !slices.Contains(c.daemons, daemon.Name()) {
				continue
			}

			nameField := zap.String("name", daemon.Name())

			log.Info("Ensuring daemon is running..", nameField)
			if err := daemon.EnsureRunning(); err != nil {
				return err
			}
			log.Info("Daemon is running", nameField)

			log.Info("Running post-launch tasks..", nameField)
			if err := daemon.PostLaunch(nodeConfig); err != nil {
				return err
			}
			log.Info("Finished post-launch tasks", nameField)
		}
	}

	return nil
}
