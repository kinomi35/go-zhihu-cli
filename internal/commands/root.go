package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/local/go-zhihu-cli/internal/config"
	"github.com/spf13/cobra"
)

type appState struct {
	endpointsPath string
	endpoints     config.Endpoints
}

var state appState

var rootCmd = &cobra.Command{
	Use:                   "zhihu",
	Short:                 "知乎命令行客户端",
	SilenceUsage:          true,
	SilenceErrors:         true,
	DisableFlagsInUseLine: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		endpoints, err := config.LoadEndpoints(state.endpointsPath)
		if err != nil {
			return fmt.Errorf("加载接口配置失败: %w", err)
		}
		state.endpoints = endpoints
		return nil
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&state.endpointsPath, "endpoints", "", "知乎 Web API 接口配置 JSON 文件路径")
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(feedCmd)
	rootCmd.AddCommand(readCmd)
	localizeCobra(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误：%s\n", localizeError(err))
		fmt.Fprintln(os.Stderr, "执行 `zhihu --help` 查看帮助。")
		os.Exit(1)
	}
}

const chineseUsageTemplate = `用法:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [命令]{{end}}{{if gt (len .Aliases) 0}}

别名:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

示例:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

可用命令:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

其他命令:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

选项:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

全局选项:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

其他帮助主题:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

使用 "{{.CommandPath}} [命令] --help" 查看某个命令的帮助。{{end}}
`

const chineseHelpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

func localizeCobra(cmd *cobra.Command) {
	disableFlagsInUseLine(cmd)
	cmd.SetUsageTemplate(chineseUsageTemplate)
	cmd.SetHelpTemplate(chineseHelpTemplate)
	cmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return errors.New(localizeError(err))
	})
	cmd.SetHelpCommand(newHelpCommand())
	addChineseHelpFlag(cmd)
	for _, child := range cmd.Commands() {
		addChineseHelpFlag(child)
	}
}

func disableFlagsInUseLine(cmd *cobra.Command) {
	cmd.DisableFlagsInUseLine = true
	for _, child := range cmd.Commands() {
		disableFlagsInUseLine(child)
	}
}

func newHelpCommand() *cobra.Command {
	helpCmd := &cobra.Command{
		Use:   "help [命令]",
		Short: "显示命令帮助",
		Long:  "显示指定命令的帮助信息。",
		Run: func(cmd *cobra.Command, args []string) {
			target := cmd.Root()
			if len(args) > 0 {
				found, _, err := cmd.Root().Find(args)
				if err != nil || found == nil {
					cmd.Printf("未知帮助主题：%s\n", strings.Join(args, " "))
					_ = cmd.Root().Usage()
					return
				}
				target = found
			}
			_ = target.Help()
		},
	}
	addChineseHelpFlag(helpCmd)
	return helpCmd
}

func addChineseHelpFlag(cmd *cobra.Command) {
	if cmd.Flags().Lookup("help") == nil {
		cmd.Flags().BoolP("help", "h", false, "显示帮助")
	}
}

func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return fmt.Errorf("%s 不接受参数: %s", cmd.CommandPath(), strings.Join(args, " "))
}

func exactArgs(count int, name string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == count {
			return nil
		}
		return fmt.Errorf("%s 需要 %d 个参数（%s），实际收到 %d 个", cmd.CommandPath(), count, name, len(args))
	}
}

func localizeError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if strings.HasPrefix(msg, "unknown command ") {
		msg = strings.Replace(msg, "unknown command", "未知命令", 1)
		msg = strings.Replace(msg, " for ", "，所属命令 ", 1)
	}
	replacer := strings.NewReplacer(
		"unknown command", "未知命令",
		"unknown flag:", "未知选项:",
		"unknown shorthand flag:", "未知短选项:",
		" in -", "，来自 -",
		"Did you mean this?", "你是不是想执行：",
		"Run", "执行",
		"for usage.", "查看用法。",
	)
	return replacer.Replace(msg)
}
