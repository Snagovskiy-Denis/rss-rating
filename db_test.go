package main

import (
	"testing"
)

func TestGetFeedByArticleUrl(t *testing.T) {
	db, shutdown := fixtureNewsboatDb(t)
	defer shutdown()

	cases := []struct {
		in   string
		want Article
	}{
		{
			fixtureArticles[0].ArticleURL,
			fixtureArticles[0],
		},
		{
			fixtureArticles[1].ArticleURL,
			fixtureArticles[1],
		},
		{
			fixtureArticles[2].ArticleURL,
			fixtureArticles[2],
		},
	}

	for _, tc := range cases {
		item, err := GetFeedByArticleUrl(db, tc.in)
		if err != nil || item != tc.want {
			t.Errorf("GetFeedByArticleUrl\nexpected: %v, %v\n actual: %v %v", tc.want, nil, item, err)
		}
	}
}

func TestGetFeedByArticleUrlThatDoesNotExist(t *testing.T) {
	db, shutdown := fixtureNewsboatDb(t)
	defer shutdown()

	_, err := GetFeedByArticleUrl(db, "https://go.dev/non/existing/path")
	if err == nil {
		t.Errorf("Excpects error, but got %v", err)
	}
}

func TestGetExistingArticleScore(t *testing.T) {
	dbs, shutdown := fixtureDbs(t)
	defer shutdown()

	article := fixtureArticles[0]
	itemZ, err := GetArticleScore(dbs.Zettl, article.ArticleURL)
	if err != nil || itemZ.ArticleURL != article.ArticleURL {
		t.Errorf("GetArticleScore\nexpected: %v, %v\n actual: %v %v", article, nil, itemZ, err)
	}
}

func TestGetNonExistingArticleScore(t *testing.T) {
	dbs, shutdown := fixtureDbs(t)
	defer shutdown()

	_, err := GetArticleScore(dbs.Zettl, "https://go.dev/non/existing/path")
	if err == nil {
		t.Errorf("Excpects error, but got %v", err)
	}
}
