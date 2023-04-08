package cmd

import (
	"fmt"

	"github.com/STARRY-S/bilibili-danmaku-client/pkg/client"
	"github.com/STARRY-S/bilibili-danmaku-client/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Execute(args []string) error {
	clientCmd := newClientCmd()
	clientCmd.addCommands()
	clientCmd.cmd.SetArgs(args)
	_, err := clientCmd.cmd.ExecuteC()
	if err != nil {
		return err
	}
	return nil
}

type clientCmd struct {
	baseCmd
}

func newClientCmd() *clientCmd {
	cc := &clientCmd{}
	cc.baseCmd.cmd = &cobra.Command{
		Use:   "client",
		Short: "Bilibili Danmaku Client Go",
		Long:  "Bilibili Danmaku Client Go",
		RunE: func(cmd *cobra.Command, args []string) error {
			initializeFlagsConfig(cmd, config.DefaultProvider)

			if config.GetBool("debug") {
				logrus.SetLevel(logrus.DebugLevel)
			}
			if err := cc.setupFlags(); err != nil {
				return err
			}
			if err := cc.run(); err != nil {
				return err
			}

			return nil
		},
	}
	cc.cmd.CompletionOptions = cobra.CompletionOptions{
		// HiddenDefaultCmd: true,
	}
	cc.cmd.Version = getVersion()
	cc.cmd.SilenceUsage = true
	cc.cmd.PersistentFlags().Bool("debug", false, "debug mode")

	cc.cmd.Flags().IntP("roomID", "r", 0, "room id")

	return cc
}

func (cc clientCmd) setupFlags() error {
	logrus.Debugf("%v", getConfigJson(config.DefaultProvider))
	if config.GetInt("roomID") == 0 {
		return fmt.Errorf("room ID not specified, " +
			"use '--roomID' to specify the room ID.")
	}
	if config.GetInt("roomID") < 0 {
		return fmt.Errorf("invalid room ID")
	}

	return nil
}

func (cc clientCmd) run() error {
	client := client.NewClient(config.GetInt("roomID"))
	if err := client.Connect(); err != nil {
		return err
	}

	return nil
}

func (cc clientCmd) getCommand() *cobra.Command {
	return cc.cmd
}

func (cc *clientCmd) addCommands() {
	addCommands(
		cc.cmd,
		newVersionCmd(),
	)
}
