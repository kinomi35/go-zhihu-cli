package commands

import (
	"encoding/json"
	"fmt"

	"github.com/local/go-zhihu-cli/internal/auth"
	"github.com/local/go-zhihu-cli/internal/client"
	"github.com/local/go-zhihu-cli/internal/output"
	"github.com/spf13/cobra"
)

var feedLimit int
var feedJSON bool

var feedCmd = &cobra.Command{
	Use:   "feed [选项]",
	Short: "查看知乎推荐流",
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
		feed, err := zhihu.Feed(feedLimit)
		if err != nil {
			return err
		}
		if feedJSON {
			data, _ := json.MarshalIndent(feed, "", "  ")
			fmt.Println(string(data))
			return nil
		}
		if len(feed.Data) == 0 {
			fmt.Println("没有推荐内容。")
			return nil
		}
		for i, item := range feed.Data {
			target := item.Target
			id := output.AnyID(target.ID)
			title := target.Title
			if title == "" && target.Question != nil {
				title = target.Question.Title
			}
			if title == "" {
				title = output.StripHTML(target.Excerpt)
			}
			author := ""
			if target.Author != nil {
				author = target.Author.Name
			}
			fmt.Printf("%2d. [%s] %s\n", i+1, target.Type, output.Truncate(output.StripHTML(title), 100))
			fmt.Printf("    ID: %s", id)
			if author != "" {
				fmt.Printf("  作者: %s", author)
			}
			fmt.Println()
			if target.Type == "answer" && id != "" {
				fmt.Printf("    阅读: zhihu read %s\n", id)
			}
		}
		return nil
	},
}

func init() {
	feedCmd.Flags().IntVarP(&feedLimit, "limit", "l", 10, "推荐内容数量")
	feedCmd.Flags().BoolVar(&feedJSON, "json", false, "输出原始 JSON")
}
