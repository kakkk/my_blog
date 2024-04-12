package impl

func MustInitCachex() {
	initArticleCachex()
	initArticleSlugCachex()
	initArticleMetaCachex()
	initArticlePostIDsCachex()
	initArticleCategoriesCachex()
	initArticleTagsCachex()
	initCategoryCachex()
	initCategoryListCachex()
	initCategoryArticleIDsCachex()
	initTagArticleIDsCachex()
	initTagListCachex()
	initCommentCachex()
	initArticleCommentsCachex()
}
