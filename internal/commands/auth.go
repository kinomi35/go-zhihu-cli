package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/local/go-zhihu-cli/internal/auth"
	"github.com/local/go-zhihu-cli/internal/client"
	"github.com/local/go-zhihu-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	loginCookies string
)

var loginCmd = &cobra.Command{
	Use:   "login [选项]",
	Short: "使用知乎浏览器 Cookie 登录",
	Args:  noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cookies, err := auth.ParseCookieString(loginCookies)
		if err != nil {
			return err
		}
		if err := auth.SaveCookieMap(cookies); err != nil {
			return err
		}
		fmt.Printf("登录成功，已保存 %d 个 Cookie 到 %s/%s\n", len(cookies), config.AppDirName, config.CookieFileName)
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout [选项]",
	Short: "清除已保存的登录 Cookie",
	Args:  noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := auth.ClearCookies(); err != nil {
			return err
		}
		fmt.Println("已退出登录。")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status [选项]",
	Short: "检查已保存的登录状态",
	Args:  noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cookies, err := auth.LoadCookieMap()
		if err != nil {
			return fmt.Errorf("未登录，请先执行 `zhihu login --cookies \"...\"`")
		}
		zhihu, err := client.New(state.endpoints, cookies)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(cmd.Context(), 15*time.Second)
		defer cancel()
		_ = ctx
		me, err := zhihu.Me()
		if err != nil {
			return err
		}
		if me.Name == "" {
			fmt.Println("已登录。")
			return nil
		}
		fmt.Printf("已登录为 %s", me.Name)
		if me.URLToken != "" {
			fmt.Printf(" (@%s)", me.URLToken)
		}
		fmt.Println()
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVar(&loginCookies, "cookies", "", "知乎浏览器 Cookie 请求头内容")
}
