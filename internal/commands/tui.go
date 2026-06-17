package commands

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/local/go-zhihu-cli/internal/auth"
	"github.com/local/go-zhihu-cli/internal/client"
	"github.com/local/go-zhihu-cli/internal/output"
	"github.com/spf13/cobra"

	"github.com/charmbracelet/bubbles/list"
	// "github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var tuiCmd = &cobra.Command{
	Use:   "tui [选项]",
	Short: "使用tui界面",
	Args:  cobra.NoArgs, 
	RunE: func(cmd *cobra.Command, args []string) error {
		cookies, err := auth.LoadCookieMap()
		if err != nil {
			return fmt.Errorf("未登录，请先执行 `zhihu login --cookies \"...\"`")
		}
		zhihu, err := client.New(state.endpoints, cookies)
		if err != nil {
			return err
		}
		p := tea.NewProgram(initialModel(zhihu), tea.WithAltScreen())
		_, err = p.Run()
		if err != nil {
			fmt.Printf("程序运行出错: %v\n", err) 
		}
		return err
	},
}

// 页面状态枚举
type page int

const (
	pageList page = iota
	pageDetail
	commentDetail
)

// 列表条目数据
type item struct {
	id          string // 唯一标识，用来匹配对应的详情内容
	title       string // 列表标题
	desc        string // 列表小字描述
	articletype string // 内容类型（如文章/问题/回答）
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func getFeedListItems(zhihu *client.Client) ([]list.Item, error) {
	feed, err := zhihu.Feed(feedLimit)
	if err != nil {
		return []list.Item{}, err
	}
	if feedJSON {
		data, _ := json.MarshalIndent(feed, "", "  ")
		fmt.Println(string(data))
		return nil, nil
	}
	if len(feed.Data) == 0 {
		fmt.Println("没有推荐内容。")
		return nil, nil
	}
	// fmt.Printf("%+v\n",feed.Data)
	feedItem := []list.Item{}
	for _, feedEntry := range feed.Data {
		// fmt.Printf("%+v\n", feedEntry)
		target := feedEntry.Target
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
		feedItem = append(feedItem, item{id: id, title: title, desc: fmt.Sprintf("[%s] %s", target.Type, author), articletype: target.Type})
	}
	return feedItem, nil
}

func getAnswerDetail(zhihu *client.Client, answerID string) (detailData, error) {
	answer, err := zhihu.Answer(answerID)
	if err != nil {
		return detailData{}, err
	}
	if readJSON {
		data, _ := json.MarshalIndent(answer, "", "  ")
		return detailData{text: string(data)}, nil
	}

	var title, author, text string
	if answer.Question != nil && answer.Question.Title != "" {
		title = answer.Question.Title
	}
	if answer.Author != nil && answer.Author.Name != "" {
		author = answer.Author.Name
	}
	if answer.Content != "" {
		text = output.StripHTMLPreserveLines(answer.Content)
	}

	return detailData{title: title, author: author, text: text, thumbsUpNum: answer.VoteupCount, commentNum: answer.CommentCount}, nil
}

type errMsg struct {
	err error
}

// 刷新列表完成的自定义消息
type refreshListMsg struct {
	newItems []list.Item
}

type loadDetailMsg struct {
	data detailData
}

type loadCommentsMsg struct {
	comment string
}

type detailData struct {
	title       string // 标题信息
	author      string // 作者信息
	text        string // 主体详情大文本
	thumbsUpNum int    // 点赞数
	commentNum  int    // 评论数
	comment     string // 评论内容，展开时展示

}

type model struct {
	zhihu *client.Client

	page         page        // 当前页面：列表/详情
	list         list.Model  // 官方默认列表组件实例
	selectedItem *item       // 当前选中的轻量列表条目
	detail       *detailData // 懒加载出来的详情（nil=未加载）
	loading      bool        // 是否正在加载（true显示加载提示）
	listLoading  bool        // 列表刷新加载状态
	viewport     viewport.Model
}

// 模拟刷新列表数据（可替换为真实网络/文件读取）
func refreshListData(zhihu *client.Client) tea.Cmd {
	return func() tea.Msg {
		items, err := getFeedListItems(zhihu)
		if err != nil {
			return refreshListMsg{newItems: []list.Item{item{id: "error", title: "刷新失败: " + err.Error(), desc: "", articletype: ""}}}
		}
		return refreshListMsg{newItems: items}
	}
}

var (
	listRefreshKey = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "刷新"),
	)
	listOpenKey = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "查看详情"),
	)
)

func newList(items []list.Item) list.Model {
	keyMap := list.DefaultKeyMap()
	keyMap.CursorUp.SetHelp("↑/k", "上移")
	keyMap.CursorDown.SetHelp("↓/j", "下移")
	keyMap.PrevPage.SetHelp("←/h/pgup", "上一页")
	keyMap.NextPage.SetHelp("→/l/pgdn", "下一页")
	keyMap.GoToStart.SetHelp("g/home", "到开头")
	keyMap.GoToEnd.SetHelp("G/end", "到末尾")
	keyMap.Filter.SetHelp("/", "筛选")
	keyMap.ClearFilter.SetHelp("esc", "清除筛选")
	keyMap.CancelWhileFiltering.SetHelp("esc", "取消")
	keyMap.AcceptWhileFiltering.SetHelp("enter", "应用筛选")
	keyMap.ShowFullHelp.SetHelp("?", "更多")
	keyMap.CloseFullHelp.SetHelp("?", "收起帮助")
	keyMap.Quit.SetHelp("q", "退出")

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.KeyMap = keyMap
	l.Title = "知乎 TUI Demo"
	l.SetSize(60, 15) // 临时尺寸，窗口消息会自适应覆盖
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{listOpenKey, listRefreshKey}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{listOpenKey, listRefreshKey}
	}
	return l
}

func loadDetailContent(zhihu *client.Client, itemID string) tea.Cmd {
	return func() tea.Msg {
		detail, err := getAnswerDetail(zhihu, itemID)
		if err != nil {
			return errMsg{err: err}
		}
		return loadDetailMsg{data: detail}
	}
}

func loadComments(zhihu *client.Client, answerID string) tea.Cmd {
	return func() tea.Msg {
		comments, err := zhihu.AnswerComments(answerID, 100)
		if err != nil {
			return errMsg{err: err}
		}
		var commentText string
		for _, comment := range comments.Data {
			author := ""
			if comment.Author != nil {
				author = comment.Author.Name
			}
			commentText += fmt.Sprintf("%s: %s\n", author, output.StripHTMLPreserveLines(comment.Content))
		}
		return loadCommentsMsg{comment: commentText}
	}
}

func initialModel(zhihu *client.Client) model {
	// 仅初始化轻量item，无任何detailData

	l := newList([]list.Item{})
	l.SetSize(60, 15) // 临时尺寸，窗口消息会自适应覆盖
	vp := viewport.New(60, 15)

	return model{
		page:     pageList, // 启动默认列表页
		list:     l,
		zhihu:    zhihu,
		viewport: vp,
		// selectedItem/detail/loading/showComment 默认零值：nil/false
	}
}

func (m model) Init() tea.Cmd {

	return refreshListData(m.zhihu) // 启动时自动加载列表数据
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case errMsg:
		m.loading = false
		errorText := fmt.Sprintf("加载出错: %v", msg.err)
		m.detail = &detailData{text: errorText}
		m.syncDetailViewport()
		return m, nil
	// 加载完成回调消息
	case loadDetailMsg:
		m.loading = false

		m.detail = &msg.data
		m.syncDetailViewport()
		return m, nil

	case loadCommentsMsg:
		if m.detail != nil {
			m.detail.comment = msg.comment
		}
		m.syncCommentViewport()
		return m, nil

	// 刷新列表完成回调消息
	case refreshListMsg:
		m.listLoading = false
		cmd := m.list.SetItems(msg.newItems)
		m.list.Select(0)
		return m, cmd

	case tea.KeyMsg:

		switch msg.String() {

		case "q":
			return m, tea.Quit
		case "esc":
			if m.page == pageDetail {
				m.page = pageList
				m.selectedItem = nil
				m.detail = nil
				m.loading = false
	
				m.syncDetailViewport()
				return m, nil
			}

		// z 返回列表，清空详情、加载状态、评论开关
		case "z":
			if m.page == pageDetail {
				m.page = pageList
				m.selectedItem = nil
				m.detail = nil
				m.loading = false
				m.syncDetailViewport()
			}
			if m.page == commentDetail {
				m.page = pageDetail
				m.syncDetailViewport()
			}
			return m, nil

			// c 切换评论（仅加载完成后生效）
		case "c":
			if m.page == pageDetail && !m.loading && m.detail != nil {
				m.page = commentDetail
				m.syncCommentViewport()
				if  m.detail.comment == "" && m.selectedItem != nil {
					return m, loadComments(m.zhihu, m.selectedItem.id)
				}
			}
			return m, nil
		// 列表界面按 r 刷新
		case "r":
			// 只有在列表页、不在刷新中才能触发
			if m.page == pageList && !m.listLoading {
				m.listLoading = true
				return m, refreshListData(m.zhihu)
			}
			return m, nil

		// 回车：选中条目，触发异步加载详情
		case "enter":
			if m.page == pageList {
				selRaw := m.list.SelectedItem()
				sel, ok := selRaw.(item)
				if ok {
					m.selectedItem = &sel
					m.page = pageDetail
					m.loading = true
					// 发起加载指令
					return m, loadDetailContent(m.zhihu, m.selectedItem.id)
				}
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width-4, msg.Height-6)
		m.setDetailViewportSize(msg.Width, msg.Height)
		m.syncDetailViewport()
		m.syncCommentViewport()
		return m, nil
	}

	if m.page == pageDetail && !m.loading && m.detail != nil {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	// 列表页面交给list自身处理上下移动
	if m.page == pageList {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}
func (m model) View() string {
	switch m.page {
	case pageList:
		style := lipgloss.NewStyle().Padding(2)
		// 如果正在刷新列表，覆盖展示加载提示
		if m.listLoading {
			return style.Render("正在刷新列表，请稍候...")
		}
		return style.Render(m.list.View())

	case pageDetail:
		style := lipgloss.NewStyle().Padding(2)
		if m.loading {
			return style.Render("正在加载详情，请稍候...\n\n[z] 返回列表")
		}
		if m.detail == nil {
			return style.Render("详情加载失败\n\n[z] 返回列表")
		}
		tips := "\n[c] 展开评论 | [↑/↓] 滚动 | [z] 返回列表 | [q] 退出"
		return style.Render(m.viewport.View() + tips)

	case commentDetail:
		style := lipgloss.NewStyle().Padding(2)
		if m.detail == nil {
			return style.Render("评论加载失败\n\n[z] 返回详情")
		}
		tips := "\n[z] 返回详情 | [q] 退出"
		return style.Render(m.viewport.View() + tips)
	default:
		return ""
	}
}

func (m *model) syncDetailViewport() {
	if m.detail == nil {
		m.viewport.SetContent("")
		return
	}
	m.viewport.SetContent(m.renderDetailContent())
}

func (m *model) syncCommentViewport() {
	if m.detail == nil {
		m.viewport.SetContent("")
		return
	}
	m.viewport.SetContent(m.renderCommentContent())
}

func (m *model) setDetailViewportSize(width, height int) {
	m.viewport.Width = max(1, width-4)
	m.viewport.Height = max(1, height-7)
}

func (m model) renderDetailContent() string {
	if m.detail == nil {
		return ""
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Underline(true).Foreground(lipgloss.Color("#FFA500"))
	authorStyle := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#00BFFF"))
	thumbsUpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4500"))
	commentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#32CD32"))

	content := titleStyle.Render(m.detail.title) +
		"\n" + authorStyle.Render("作者: "+m.detail.author) +
		"\n\n" + m.detail.text +
		"\n\n" + thumbsUpStyle.Render(fmt.Sprintf("赞同: %d", m.detail.thumbsUpNum)) +
		"  " + commentStyle.Render(fmt.Sprintf("评论: %d", m.detail.commentNum))

	return content
}
func (m model) renderCommentContent() string {
	if m.detail.comment == "" {
		return "暂无评论"
	}
	var content string
		comment := m.detail.comment
		if comment == "" {
			comment = "正在加载评论..."
		}
		content = lipgloss.NewStyle().Foreground(lipgloss.Color("#ADFF2F")).Render("评论内容:\n"+comment)

	return content
}