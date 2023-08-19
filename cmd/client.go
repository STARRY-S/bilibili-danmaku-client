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
		Use:   "bilibili-danmaku-client",
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
	cc.cmd.PersistentFlags().Bool("debug", false, "Debug mode")

	cc.cmd.Flags().Bool("voice", true, "Enable voice output")
	cc.cmd.Flags().String("voiceAPI", "sougou", "Voice API (available: sougou)")
	cc.cmd.Flags().IntP("roomID", "r", 0, "直播间 ID (Room ID)")
	cc.cmd.Flags().IntP("uid", "u", 0, "正在观看的用户名 ID (User ID)")

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
	if config.GetInt("uid") == 0 {
		logrus.Warnf("UID is not provided, some function may not working properly")
	}
	if config.GetInt("uid") < 0 {
		return fmt.Errorf("invalid uid")
	}
	if config.GetBool("voice") {
		logrus.Infof("Voice output enabled")
		switch config.GetString("voiceAPI") {
		case "sougou":
			logrus.Debugf("Voice API set to sougou")
		}
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
