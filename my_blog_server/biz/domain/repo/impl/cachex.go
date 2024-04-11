package impl

func MustInitCachex() {
	initArticleCachex()
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
