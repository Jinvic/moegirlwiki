package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const apiURL = "https://zh.moegirl.org.cn/api.php"

// APIResponse 表示API响应的基本结构
type APIResponse struct {
	Query *QueryResult `json:"query,omitempty"`
}

// QueryResult 表示查询结果
type QueryResult struct {
	SearchInfo *SearchInfo `json:"searchinfo,omitempty"`
	Search     []SearchResult `json:"search,omitempty"`
	Pages      map[string]Page `json:"pages,omitempty"`
}

// SearchInfo 搜索信息
type SearchInfo struct {
	TotalHits int `json:"totalhits"`
}

// SearchResult 搜索结果
type SearchResult struct {
	NS        int    `json:"ns"`
	Title     string `json:"title"`
	PageID    int    `json:"pageid"`
	Size      int    `json:"size"`
	WordCount int    `json:"wordcount"`
	Snippet   string `json:"snippet"`
	Timestamp string `json:"timestamp"`
}

// Page 页面信息
type Page struct {
	PageID    int    `json:"pageid"`
	NS        int    `json:"ns"`
	Title     string `json:"title"`
	Revisions []Revision `json:"revisions,omitempty"`
}

// Revision 页面修订版本
type Revision struct {
	ContentFormat string `json:"contentformat"`
	ContentModel  string `json:"contentmodel"`
	Content       string `json:"*"`
}

// MoegirlClient 萌娘百科API客户端
type MoegirlClient struct {
	client *http.Client
}

// NewMoegirlClient 创建新的萌娘百科客户端
func NewMoegirlClient() *MoegirlClient {
	return &MoegirlClient{
		client: &http.Client{},
	}
}

// Search 搜索功能
func (m *MoegirlClient) Search(query string, limit int) ([]SearchResult, error) {
	params := url.Values{}
	params.Set("action", "query")
	params.Set("format", "json")
	params.Set("list", "search")
	params.Set("srsearch", query)
	params.Set("srlimit", fmt.Sprintf("%d", limit))

	resp, err := m.makeRequest(params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	if apiResp.Query == nil || apiResp.Query.Search == nil {
		return nil, fmt.Errorf("未找到搜索结果")
	}

	return apiResp.Query.Search, nil
}

// GetPageByTitle 通过标题获取页面内容
func (m *MoegirlClient) GetPageByTitle(title string) (*Page, error) {
	params := url.Values{}
	params.Set("action", "query")
	params.Set("format", "json")
	params.Set("titles", title)
	params.Set("prop", "revisions")
	params.Set("rvprop", "content")

	resp, err := m.makeRequest(params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	if apiResp.Query == nil || apiResp.Query.Pages == nil {
		return nil, fmt.Errorf("未找到页面")
	}

	// 获取第一个页面（通常只有一个）
	for _, page := range apiResp.Query.Pages {
		return &page, nil
	}

	return nil, fmt.Errorf("未找到页面内容")
}

// GetPageByID 通过页面ID获取页面内容
func (m *MoegirlClient) GetPageByID(pageID int) (*Page, error) {
	params := url.Values{}
	params.Set("action", "query")
	params.Set("format", "json")
	params.Set("pageids", fmt.Sprintf("%d", pageID))
	params.Set("prop", "revisions")
	params.Set("rvprop", "content")

	resp, err := m.makeRequest(params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	if apiResp.Query == nil || apiResp.Query.Pages == nil {
		return nil, fmt.Errorf("未找到页面")
	}

	// 获取第一个页面（通常只有一个）
	for _, page := range apiResp.Query.Pages {
		return &page, nil
	}

	return nil, fmt.Errorf("未找到页面内容")
}

// makeRequest 发起API请求
func (m *MoegirlClient) makeRequest(params url.Values) (*http.Response, error) {
	reqURL := apiURL + "?" + params.Encode()
	resp, err := m.client.Get(reqURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	return resp, nil
}

// printSearchResults 打印搜索结果
func printSearchResults(results []SearchResult) {
	if len(results) == 0 {
		fmt.Println("没有找到匹配的结果")
		return
	}

	fmt.Printf("找到 %d 条结果:\n", len(results))
	fmt.Println(strings.Repeat("-", 80))
	
	for i, result := range results {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   页面ID: %d\n", result.PageID)
		fmt.Printf("   字符数: %d\n", result.Size)
		if result.Snippet != "" {
			fmt.Printf("   摘要: %s\n", stripHTML(result.Snippet))
		}
		fmt.Println(strings.Repeat("-", 80))
	}
}

// printPageContent 打印页面内容
func printPageContent(page *Page) {
	if page.Title == "" {
		fmt.Println("未找到页面")
		return
	}

	fmt.Printf("页面标题: %s\n", page.Title)
	fmt.Printf("页面ID: %d\n", page.PageID)
	fmt.Println(strings.Repeat("=", 80))

	if len(page.Revisions) > 0 && page.Revisions[0].Content != "" {
		content := page.Revisions[0].Content
		// 限制输出长度，避免输出过长的内容
		if len(content) > 2000 {
			content = content[:2000] + "\n\n... (内容太长，已截断) ..."
		}
		fmt.Println(content)
	} else {
		fmt.Println("页面内容为空")
	}
}

// stripHTML 简单的HTML标签移除函数
func stripHTML(text string) string {
	result := text
	// 简单地移除HTML标签（更完整的实现可能需要使用HTML解析库）
	for {
		start := strings.Index(result, "<")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ">")
		if end == -1 {
			break
		}
		end = start + end + 1
		result = result[:start] + result[end:]
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法:")
		fmt.Println("  搜索: go run main.go search <关键词> [数量]")
		fmt.Println("  查看(按标题): go run main.go view <页面标题>")
		fmt.Println("  查看(按ID): go run main.go viewid <页面ID>")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  go run main.go search 萌娘 5")
		fmt.Println("  go run main.go view 萌娘百科:首页")
		fmt.Println("  go run main.go viewid 23528")
		return
	}

	client := NewMoegirlClient()
	command := os.Args[1]

	switch command {
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("错误: 请提供搜索关键词")
			return
		}
		
		query := os.Args[2]
		limit := 10 // 默认搜索结果数量
		if len(os.Args) > 3 {
			if l, err := strconv.Atoi(os.Args[3]); err == nil {
				limit = l
			}
		}
		
		fmt.Printf("正在搜索: %s\n", query)
		results, err := client.Search(query, limit)
		if err != nil {
			fmt.Printf("搜索失败: %v\n", err)
			return
		}
		
		printSearchResults(results)

	case "view":
		if len(os.Args) < 3 {
			fmt.Println("错误: 请提供页面标题")
			return
		}
		
		title := strings.Join(os.Args[2:], " ")  // 支持带空格的标题
		fmt.Printf("正在获取页面: %s\n", title)
		page, err := client.GetPageByTitle(title)
		if err != nil {
			fmt.Printf("获取页面失败: %v\n", err)
			return
		}
		
		printPageContent(page)

	case "viewid":
		if len(os.Args) < 3 {
			fmt.Println("错误: 请提供页面ID")
			return
		}
		
		pageID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("页面ID必须是数字: %v\n", err)
			return
		}
		
		fmt.Printf("正在获取页面ID: %d\n", pageID)
		page, err := client.GetPageByID(pageID)
		if err != nil {
			fmt.Printf("获取页面失败: %v\n", err)
			return
		}
		
		printPageContent(page)

	default:
		fmt.Printf("未知命令: %s\n", command)
		fmt.Println("用法:")
		fmt.Println("  搜索: go run main.go search <关键词> [数量]")
		fmt.Println("  查看(按标题): go run main.go view <页面标题>")
		fmt.Println("  查看(按ID): go run main.go viewid <页面ID>")
	}
}