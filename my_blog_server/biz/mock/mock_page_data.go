package mock

import (
	"bytes"
	"io/ioutil"

	"my_blog/biz/common/config"
	"my_blog/biz/model/blog/page"

	"github.com/russross/blackfriday/v2"
)

func PostListPageRespMocker(pageType string, name string, pre string, next string, slug string) *page.PostListPageResp {
	return &page.PostListPageResp{
		Meta: &page.PageMeta{
			Title:       "kakkk's blog",
			Description: "this is kakkk's blog",
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			PageType:    pageType,
		},
		Name:     name,
		PrevPage: pre,
		NextPage: next,
		Slug:     slug,
		PostList: []*page.PostItem{
			{
				ID:       "1",
				Title:    "测试文章",
				Abstract: "这是一篇测试文章测试一下测试测试测试测试测试测试测试测试测试测试测试测试测试测试测试",
				Info:     "November 25, 2022 · kakkk",
			},
			{
				ID:       "2",
				Title:    "测试文章",
				Abstract: "这是一篇测试文章测试一下测试测试测试测试测试测试测试测试测试测试测试测试测试测试测试",
				Info:     "November 25, 2022 · 1 min · kakkk",
			},
			{
				ID:       "3",
				Title:    "测试文章",
				Abstract: "这是一篇测试文章测试一下测试测试测试测试测试测试测试测试测试测试测试测试测试测试测试",
				Info:     "November 25, 2022 · 1 min · kakkk",
			},
			{
				ID:       "4",
				Title:    "测试文章",
				Abstract: "这是一篇测试文章测试一下测试测试测试测试测试测试测试测试测试测试测试测试测试测试测试",
				Info:     "November 25, 2022 · 1 min · kakkk",
			},
			{
				ID:       "5",
				Title:    "测试文章",
				Abstract: "这是一篇测试文章测试一下测试测试测试测试测试测试测试测试测试测试测试测试测试测试测试",
				Info:     "November 25, 2022 · 1 min · kakkk",
			},
		},
	}
}

func ArchivesPageRespMocker() *page.ArchivesPageResp {
	return &page.ArchivesPageResp{
		Meta: &page.PageMeta{
			Title:       "kakkk's blog",
			Description: "this is kakkk's blog",
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			PageType:    page.PageTypeArchives,
		},
		PostArchives: []*page.ArchiveByYear{
			{
				Year: "2022",
				Archives: []*page.ArchiveByMonth{
					{
						Month: "November",
						Count: "1",
						Posts: []*page.PostItem{
							{
								ID:    "1",
								Title: "测试文章1",
								Info:  "2022.11.1",
							},
						},
					},
					{
						Month: "January",
						Count: "2",
						Posts: []*page.PostItem{
							{
								ID:    "2",
								Title: "测试文章2",
								Info:  "2022.1.25",
							},
							{
								ID:    "3",
								Title: "测试文章3",
								Info:  "2022.1.2",
							},
						},
					},
				},
			},
			{
				Year: "2021",
				Archives: []*page.ArchiveByMonth{
					{
						Month: "December",
						Count: "1",
						Posts: []*page.PostItem{
							{
								ID:    "4",
								Title: "测试文章4",
								Info:  "2021.12.1",
							},
						},
					},
				},
			},
		},
	}
}

func TagsMocker() *page.TermsPageResp {
	return &page.TermsPageResp{
		Meta: &page.PageMeta{
			Title:       "kakkk's blog",
			Description: "test",
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			PageType:    page.PageTypeTagList,
		},
		List: []*page.TermListItem{
			{
				Name:  "标签",
				Count: "12",
				Slug:  "标签",
			},
			{
				Name:  "这是一个标签",
				Count: "6",
				Slug:  "这是一个标签",
			},
			{
				Name:  "Golang",
				Count: "34",
				Slug:  "Golang",
			},
			{
				Name:  "Java",
				Count: "4",
				Slug:  "Java",
			},
			{
				Name:  "Linux",
				Count: "1",
				Slug:  "Linux",
			},
			{
				Name:  "踩过的坑",
				Count: "6",
				Slug:  "踩过的坑",
			},
			{
				Name:  "测试",
				Count: "3",
				Slug:  "测试",
			},
		},
	}
}

func CategoriesMocker() *page.TermsPageResp {
	return &page.TermsPageResp{
		Meta: &page.PageMeta{
			Title:       "kakkk's blog",
			Description: "this is kakkk's blog",
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			PageType:    page.PageTypeCategoryList,
		},
		List: []*page.TermListItem{
			{
				Name:  "测试分类",
				Count: "12",
				Slug:  "test_category",
			},
			{
				Name:  "这是一个分类",
				Count: "6",
				Slug:  "this_is_a_category",
			},
			{
				Name:  "Golang",
				Count: "34",
				Slug:  "golag",
			},
			{
				Name:  "Java",
				Count: "4",
				Slug:  "java",
			},
		},
	}
}

func SearchMocker() *page.BasicPageResp {
	return &page.BasicPageResp{
		Meta: &page.PageMeta{
			Title:       "Search",
			Description: "search",
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			PageType:    page.PageTypeSearch,
		},
	}
}

func ErrorPageMocker(code string) *page.BasicPageResp {
	return &page.BasicPageResp{
		Meta: &page.PageMeta{
			Title:       "error",
			Description: "error",
			SiteDomain:  "http://127.0.0.1:8888",
			PageType:    page.PageTypeError,
		},
	}
}

func PostPageMocker(id string) *page.PostPageResponse {
	file, err := ioutil.ReadFile("../mock_post.md")
	// mock测试，错误直接panic
	if err != nil {
		panic(err)
	}
	file = bytes.Replace(file, []byte("\r"), nil, -1)
	content := string(blackfriday.Run(
		file,
		blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak),
	))
	return &page.PostPageResponse{
		Title: "测试文章" + id,
		Info: &page.PostInfo{
			Author:    "kakkk",
			PublishAt: "November 25, 2022",
			UV:        "20",
			WordCount: "1024",
			CategoryList: []*page.TermListItem{
				{
					Name: "测试分类",
					Slug: "test",
				},
				{
					Name: "Golang",
					Slug: "golang",
				},
			},
		},
		Content: content,
		Tags:    []string{"标签", "测试", "Golang"},
		PrevPage: &page.PostNav{
			Title: "测试文章1",
			ID:    "1",
		},
		NextPage: &page.PostNav{
			Title: "测试文章3",
			ID:    "3",
		},
		Meta: &page.PageMeta{
			Title:       "post page",
			Description: "post page",
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			PageType:    page.PageTypePost,
		},
	}
}
