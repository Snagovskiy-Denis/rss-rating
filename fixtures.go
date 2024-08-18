package main

import (
	"database/sql"
	"testing"
)

var fixtureArticles = [3]Article{
	{
		"https://github.com/eradman/entr/releases.atom",
		"https://github.com/eradman/entr/releases/tag/5.6",
		"Release notes from entr",
		"entr-5.6",
		1719857071,
	},
	{
		"https://github.com/eradman/entr/releases.atom",
		"https://github.com/eradman/entr/releases/tag/3.3",
		"Release notes from entr",
		"3.3",
		1611930981,
	},
	{
		"https://go.dev/blog/feed.atom",
		"https://go.dev/blog/go1.22",
		"The Go Blog",
		"Go 1.22 is released!",
		1611878400,
	},
}

var fixtureArticleScores = [3]ArticleScore{
	{
		1,
		fixtureArticles[0],
	},
	{
		1,
		fixtureArticles[2],
	},
	{
		2,
		Article{
			"https://go.dev/blog/feed.atom",
			"https://go.dev/blog/routing-enhancements",
			"The Go Blog",
			"Routing Enhancements for Go 1.22",
			1719792000,
		},
	},
}

func fixtureNewsboatDb(t *testing.T) (*sql.DB, func()) {
	db, _ := sql.Open("sqlite", ":memory:")

	db.Exec(`
        CREATE TABLE rss_item (
            id INTEGER PRIMARY KEY,
            url VARCHAR(1024) NOT NULL,
            feedurl VARCHAR(1024) NOT NULL,
            pubDate INTEGER NOT NULL,
            title VARCHAR(1024) NOT NULL
        )
    `)

	db.Exec(`
        CREATE TABLE rss_feed (
            rssurl VARCHAR(1024) PRIMARY KEY NOT NULL,
            title VARCHAR(1024) NOT NULL
        )
    `)

	for _, item := range fixtureArticles {
		db.Exec(
			`INSERT INTO rss_item (url, feedurl, pubDate, title) VALUES (?, ?, ?, ?)`,
			item.ArticleURL, item.FeedURL, item.PubDate, item.ArticleTitle,
		)
		db.Exec(
			`INSERT INTO rss_feed (rssurl, title) VALUES (?, ?)`,
			item.FeedURL, item.FeedTitle,
		)
	}

	return db, func() {
		db.Close()
	}
}

func fixtureZettelkastenDb(t *testing.T) (*sql.DB, func()) {
	db, _ := sql.Open("sqlite", ":memory:")

	db.Exec(`
        CREATE TABLE rss_scores (
            -- id INTEGER PRIMARY KEY AUTOINCREMENT,
            feed_url TEXT NOT NULL,
            feed_title TEXT NOT NULL,
            article_url TEXT NOT NULL,
            article_title TEXT NOT NULL,
            pub_date INTEGER NOT NULL,
            score INTEGER NOT NULL,
            scored_at INTEGER DEFAULT (unixepoch()),
            -- UNIQUE(feed_url, article_url)
            PRIMARY KEY (feed_url, article_url)
        )
    `)

	for _, articleScore := range fixtureArticleScores {
		db.Exec(
			`INSERT INTO rss_scores (
                feed_url,
                feed_title,
                article_url,
                article_title,
                pub_date,
                score
            )
            VALUES (?, ?, ?, ?, ?, ?)`,
			articleScore.FeedURL,
			articleScore.FeedTitle,
			articleScore.ArticleURL,
			articleScore.ArticleTitle,
			articleScore.PubDate,
			articleScore.Score,
		)
	}

	return db, func() {
		db.Close()
	}
}

func fixtureDbs(t *testing.T) (*Databases, func()) {
	dbN, closeDbN := fixtureNewsboatDb(t)
	dbZ, closeDbZ := fixtureZettelkastenDb(t)
	dbs := &Databases{dbN, dbZ}
	return dbs, func() {
		closeDbN()
		closeDbZ()
	}
}
