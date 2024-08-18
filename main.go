package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// TODO:
//      - auto-mark video as watched in newsboat cache at score cmd

func main() {
	dbPathNewsboat := filepath.Join(os.Getenv("XDG_DATA_HOME"), "newsboat/cache.db")
	dbPathZettelkasten, ok := os.LookupEnv("ZETTELKASTEN_DB")
	if !ok {
		log.Fatal("Missing ZETTELKASTEN_DB")
	}

	dbs, err := NewDatabases(dbPathNewsboat, dbPathZettelkasten)
	if err != nil {
		log.Fatalf("Cannot connect to db: %v", err)
	}
	defer dbs.Close()

	scoreCmd := flag.NewFlagSet("score", flag.ExitOnError)
	scoreInt := scoreCmd.Int("score", -777, "negative or positive integer")
	articleUrl1 := scoreCmd.String(
		"article-url",
		"",
		"expects '%u' value from newsboat macro on article-detail page",
	)

	reportCmd := flag.NewFlagSet("report", flag.ExitOnError)
	feedUrl := reportCmd.String(
		"feed-url",
		"",
		"expects '%u' value from newsboat macro on feed-list page",
	)
	// limit := reportCmd.Int("limit", -1, "max num of articles to include in report")
	dateStart := reportCmd.String("start-date", "", "ge pubDate filter (format: YYYY-MM-DD)")
	dateEnd := reportCmd.String("end-date", "", "le pubDate filter (format: YYYY-MM-DD)")

	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
	articleUrl2 := checkCmd.String(
		"article-url",
		"",
		"expects '%u' value from newsboat macro on article-detail page",
	)

	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
	cmdName := helpCmd.String("cmd", "help", "subcommand name to show help to")

	if len(os.Args) < 2 {
		help(nil)
		os.Exit(1)
	}

	switch os.Args[1] {
	default:
		help(nil)
	case "help":
		helpCmd.Parse(os.Args[2:])
		cmd, _ := map[string]*flag.FlagSet{
			"report": reportCmd,
			"score":  scoreCmd,
			"check":  checkCmd,
			"help":   helpCmd,
		}[*cmdName]
		help(cmd)
	case "score":
		scoreCmd.Parse(os.Args[2:])
		if *articleUrl1 == "" {
			log.Fatal("Missing article-url option")
		}
		if *scoreInt == -777 {
			log.Fatal("Missing score option")
		}
		if err := scoreSubcommand(dbs, *articleUrl1, *scoreInt); err != nil {
			log.Fatal(err)
		}
	case "report":
		reportCmd.Parse(os.Args[2:])
		if *feedUrl == "" {
			log.Fatal("Missing feed-url option")
		}
        report, err := reportSubcommand(dbs, *feedUrl, *dateStart, *dateEnd)
        if err != nil {
            os.Exit(1)
        }
		fmt.Println(report)
	case "check":
		checkCmd.Parse(os.Args[2:])
		if *articleUrl2 == "" {
			log.Fatal("Missing article-url option")
		}
		as, _ := GetArticleScore(dbs.Zettl, *articleUrl2)
		fmt.Println(as.Score)
	}
}

func help(cmd *flag.FlagSet) {
	fmt.Print("Newsboat plugin that adds personal likes and dislikes to articles\n\n")
	fmt.Println("Usage: rss-rating (score | report | check | help) OPTION... [ARGUMENT]...")
	if cmd != nil {
		fmt.Printf("\n%v options:\n", cmd.Name())
		cmd.PrintDefaults()
	}
}
