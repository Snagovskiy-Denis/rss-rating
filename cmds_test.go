package main

import "testing"

func TestScoreSubcommandNewArticle(t *testing.T) {
	dbs, shutdown := fixtureDbs(t)
	defer shutdown()

	var err error
	articleUrl, scoreInt := fixtureArticles[1].ArticleURL, 1

	if _, err = GetArticleScore(dbs.Zettl, articleUrl); err == nil {
		t.Error("Invalid test setup")
	}

	scoreSubcommand(dbs, articleUrl, scoreInt)

	item, err := GetArticleScore(dbs.Zettl, articleUrl)
	if err != nil {
		t.Error(err)
	} else if item.Article.ArticleURL != articleUrl || item.Score != scoreInt {
		t.Errorf("invalid item:\n%v\n%v %v", item, articleUrl, scoreInt)
	}
}

func TestScoreSubcommandExistingArticle(t *testing.T) {
	dbs, shutdown := fixtureDbs(t)
	defer shutdown()

	var (
		before ArticleScore
		after  ArticleScore
		err    error
	)
	articleUrl, scoreInt := fixtureArticles[0].ArticleURL, 2

	before, err = GetArticleScore(dbs.Zettl, articleUrl)
	if err != nil {
		t.Error("Invalid test setup")
	}

	scoreSubcommand(dbs, articleUrl, scoreInt)

	after, err = GetArticleScore(dbs.Zettl, articleUrl)
	if err != nil {
		t.Error(err)
	} else if after.Score != scoreInt || before == after {
		t.Errorf("before and after are equal:\n%v\n%v", before, after)
	}
}

func TestReport(t *testing.T) {
	dbs, shutdown := fixtureDbs(t)
	defer shutdown()

	middleDate := "2023-01-01"
	var (
		defaultStart string
		defaultEnd   string
	)

	cases := []struct {
		wants Report
		start string
		end   string
	}{
		{
			Report{
				fixtureArticleScores[1].FeedURL,
				fixtureArticleScores[1].FeedTitle,
				float32(fixtureArticleScores[1].Score+fixtureArticleScores[2].Score) / 2,
			},
			defaultStart,
			defaultEnd,
		},
		{
			Report{
				fixtureArticleScores[2].FeedURL,
				fixtureArticleScores[2].FeedTitle,
				float32(fixtureArticleScores[2].Score),
			},

			middleDate,
			defaultEnd,
		},
		{
			Report{
				fixtureArticleScores[1].FeedURL,
				fixtureArticleScores[1].FeedTitle,
				float32(fixtureArticleScores[1].Score),
			},
			defaultStart,
			middleDate,
		},
		{
			Report{
				fixtureArticleScores[2].FeedURL,
				fixtureArticleScores[2].FeedTitle,
				float32(fixtureArticleScores[2].Score),
			},
			middleDate,
			"2030-01-01",
		},
	}

	for _, c := range cases {
		actual, _ := reportSubcommand(dbs, fixtureArticleScores[1].FeedURL, c.start, c.end)
		if c.wants.AvgScore != actual.AvgScore {
			t.Errorf(
				"Invalid report:\n%v != %v\nfrom %v to %v",
				actual, c.wants, c.start, c.end,
			)
		}
	}

	actual, _ := reportSubcommand(dbs, "non-existing-feed", defaultStart, defaultEnd)
	if actual.AvgScore != 0 {
		t.Errorf("Invalid report:\n%v != 0", actual)
	}
}
