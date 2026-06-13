package commands

import (
	"encoding/json"
	"fmt"

	"github.com/local/go-zhihu-cli/internal/auth"
	"github.com/local/go-zhihu-cli/internal/client"
	"github.com/local/go-zhihu-cli/internal/output"
	"github.com/spf13/cobra"
)

var readJSON bool
var readComments int

var readCmd = &cobra.Command{
	Use:   "read <回答ID> [选项]",
	Short: "按回答 ID 阅读推荐回答",
	Args:  exactArgs(1, "回答ID"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cookies, err := auth.LoadCookieMap()
		if err != nil {
			return fmt.Errorf("未登录，请先执行 `zhihu login --cookies \"...\"`")
		}
		zhihu, err := client.New(state.endpoints, cookies)
		if err != nil {
			return err
		}
		answer, err := zhihu.Answer(args[0])
		if err != nil {
			return err
		}
		if readJSON {
			data, _ := json.MarshalIndent(answer, "", "  ")
			fmt.Println(string(data))
			return nil
		}
		if answer.Question != nil && answer.Question.Title != "" {
			fmt.Println(output.StripHTML(answer.Question.Title))
			fmt.Println()
		}
		if answer.Author != nil && answer.Author.Name != "" {
			fmt.Printf("作者：%s\n\n", answer.Author.Name)
		}
		fmt.Println(output.StripHTML(answer.Content))
		fmt.Printf("\n赞同：%d  评论：%d\n", answer.VoteupCount, answer.CommentCount)

		if readComments > 0 {
			comments, err := zhihu.AnswerComments(args[0], readComments)
			if err != nil {
				return err
			}
			if len(comments.Data) > 0 {
				fmt.Println("\n评论：")
				for i, comment := range comments.Data {
					author := ""
					if comment.Author != nil {
						author = comment.Author.Name
					}
					fmt.Printf("%2d. ", i+1)
					if author != "" {
						fmt.Printf("%s: ", author)
					}
					fmt.Printf("%s", output.Truncate(output.StripHTML(comment.Content), 160))
					if comment.VoteCount > 0 {
						fmt.Printf("  (%d 赞)", comment.VoteCount)
					}
					fmt.Println()
				}
			}
		}
		return nil
	},
}

func init() {
	readCmd.Flags().BoolVar(&readJSON, "json", false, "输出原始 JSON")
	readCmd.Flags().IntVarP(&readComments, "comments", "c", 0, "显示的评论数量")
}
