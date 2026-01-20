package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/urfave/cli/v3"
)

type Champ struct {
	Name       string
	WinCount   float64
	MatchCount float64
	Counters   []ChampCounter
	AllMatches float64 // all the matches analyzed by u.gg for pick rate calculation
	Winrate    float64
	Pickrate   float64
}

type ChampCounter struct {
	Name             string
	LostMatchCounter float64
	MatchCount       float64
	LoseRate         float64
}

var lolzyServerHost string = "lolzy.bozsik-services.me"

func Commands() *cli.Command {
	return &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "meta",
				Usage: "meta <role_name> <options>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "top",
						Aliases: []string{"t"},
						Value:   "0",
						Usage:   "Sets how many champs do you want to get. If not specified, all the champs will be queried.",
					},
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Value:   false,
						Usage:   "If true, champs will be included in the query with less than 0.5 percent pick rate",
					},
					&cli.StringFlag{
						Name:    "rank",
						Aliases: []string{"r"},
						Value:   "overall",
						Usage:   "Filters the results to include only matches from the specified rank. ",
					},
				},
				UseShortOptionHandling: true,
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.NArg() == 0 {
						return fmt.Errorf("Champ name not specified")
					}
					serverURL := url.URL{
						Host:   lolzyServerHost,
						Scheme: "http",
						Path:   fmt.Sprintf("/api/%s/meta", c.Args().First()),
					}
					q := serverURL.Query()
					var boolStr string
					if c.Bool("all") {
						boolStr = "true"
					} else {
						boolStr = "false"
					}
					q.Add("top", c.String("top"))
					q.Add("all", boolStr)
					q.Add("rank", c.String("rank"))
					serverURL.RawQuery = q.Encode()
					fmt.Println("serverurl: ", serverURL.String())
					req, err := http.NewRequest("GET", serverURL.String(), nil)
					if err != nil {
						return nil
					}
					client := http.Client{}
					res, err := client.Do(req)
					if err != nil {
						return nil
					}
					defer res.Body.Close()
					var result []Champ
					err = json.NewDecoder(res.Body).Decode(&result)
					if err != nil {
						return nil
					}
					RenderAllChamps(result)
					return nil
				},
			},
			{
				Name:  "counter",
				Usage: "counter <champname> <options>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "role",
						Aliases: []string{"r"},
						Value:   "",
						Usage:   "You can specify which role are you intersted in. ",
					},
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Value:   true,
						Usage:   "If true, champs will be included in the query with less than 10 match data",
					},
					&cli.StringFlag{
						Name:    "rank",
						Aliases: []string{"rk"},
						Value:   "overall",
						Usage:   "Filters the results to include only matches from the specified rank. ",
					},
				},
				UseShortOptionHandling: true,
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.NArg() == 0 {
						return fmt.Errorf("Champoin name not specified ")
					}
					serverURL := url.URL{
						Host:   lolzyServerHost,
						Scheme: "http",
						Path:   fmt.Sprintf("/api/%s/counter", c.Args().First()),
					}
					q := serverURL.Query()
					var boolStr string
					if c.Bool("all") {
						boolStr = "true"
					} else {
						boolStr = "false"
					}
					q.Add("role", c.String("role"))
					q.Add("rank", c.String("rank"))
					q.Add("all", boolStr)
					serverURL.RawQuery = q.Encode()
					req, err := http.NewRequest("GET", serverURL.String(), nil)
					if err != nil {
						return nil
					}
					client := http.Client{}
					res, err := client.Do(req)
					if err != nil {
						return nil
					}
					defer res.Body.Close()
					var result map[string][]ChampCounter
					err = json.NewDecoder(res.Body).Decode(&result)
					if err != nil {
						fmt.Println("There was a problem with the parsing of the result: ", err)
						return err
					}
					RenderCounterMap(result)
					return nil
				},
			},
		},
	}
}
