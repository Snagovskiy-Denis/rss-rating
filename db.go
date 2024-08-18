package main

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

type Databases struct {
	Newsboat *sql.DB
	Zettl    *sql.DB
}

type Article struct {
	FeedURL      string
	ArticleURL   string
	FeedTitle    string
	ArticleTitle string
	PubDate      int
}

type ArticleScore struct {
	Score int
	Article
}

func NewDatabases(pathNewsboat, pathZettl string) (*Databases, error) {
	var dbs Databases
	var err error
	dbs.Newsboat, err = sql.Open("sqlite", pathNewsboat)
	if err != nil {
		return &dbs, err
	}
	dbs.Zettl, err = sql.Open("sqlite", pathZettl)
	if err != nil {
		return &dbs, err
	}
	return &dbs, nil
}

func (dbs *Databases) Close() {
	dbs.Newsboat.Close()
	dbs.Zettl.Close()
}

func GetFeedByArticleUrl(dbNewsboat *sql.DB, url string) (Article, error) {
	row := dbNewsboat.QueryRow(
		`SELECT
            feed.rssurl,
            feed.title,
            item.url,
            item.title,
            item.pubDate
        FROM rss_item AS item
        JOIN rss_feed AS feed ON item.feedurl = feed.rssurl
        WHERE item.url = ?;`,
		url,
	)

	var a Article
	err := row.Scan(
		&a.FeedURL,
		&a.FeedTitle,
		&a.ArticleURL,
		&a.ArticleTitle,
		&a.PubDate,
	)
	return a, err
}

func GetArticleScore(dbZetl *sql.DB, articleURL string) (ArticleScore, error) {
	row := dbZetl.QueryRow(
		`SELECT
            feed_url,
            feed_title,
            article_url,
            article_title,
            pub_date,
            score
        FROM rss_scores
        WHERE article_url = ?`,
		articleURL,
	)

	var as ArticleScore
	err := row.Scan(
		&as.FeedURL,
		&as.FeedTitle,
		&as.ArticleURL,
		&as.ArticleTitle,
		&as.PubDate,
		&as.Score,
	)
    return as, err
}
