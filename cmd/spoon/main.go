package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/vkuznecovas/mouthful/cmd/spoon/command"
)

func main() {
	app := cli.NewApp()
	app.Name = "Spoon"
	app.Usage = "the helpful mouthful helper - spoon will export and import mouthful data, migrate between db engines and import comments from other commenting engines"
	app.Version = "v0.0.0"
	// app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "config, c",
					Value:  "",
					Usage:  "path to mouthful config file",
					EnvVar: "MOUTHFUL_CONFIG",
				},
			},
			Name:    "export",
			Aliases: []string{"e"},
			Usage:   "exports the database pointed by the config provided to a mouthful.dmp file",
			Action: func(c *cli.Context) error {
				configPath := c.String("config")
				return command.ExportCommandRun(configPath)
			},
		},
		{
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "config, c",
					Value:  "",
					Usage:  "path to mouthful config file",
					EnvVar: "MOUTHFUL_CONFIG",
				},
				cli.StringFlag{
					Name:   "dump, d",
					Value:  "",
					Usage:  "path to mouthful dump file",
					EnvVar: "MOUTHFUL_DUMP",
				},
			},
			Name:    "import",
			Aliases: []string{"i"},
			Usage:   "imports the  given dump to the database pointed by the config provided",
			Action: func(c *cli.Context) error {
				configPath := c.String("config")
				dumpPath := c.String("dump")
				return command.ImportCommandRun(configPath, dumpPath)
			},
		},
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "migrations between different commenting engines and supported mouthful database types",
			Subcommands: cli.Commands{
				{
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "dump, d",
							Value:  "",
							Usage:  "path to disqus dump file",
							EnvVar: "DISQUS_DUMP",
						},
					},
					Name:  "disqus",
					Usage: "imports the given dump to a mouthful sqlite instance",
					Action: func(c *cli.Context) error {
						dumpPath := c.String("dump")
						return command.DisqusMigrationRun(dumpPath)
					},
				},
				{
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "isso",
							Value:  "",
							Usage:  "path to isso sqlite file",
							EnvVar: "ISSO_FILE",
						},
					},
					Name:  "isso",
					Usage: "imports the given isso sqlite to a mouthful sqlite instance",
					Action: func(c *cli.Context) error {
						issoPath := c.String("isso")
						return command.IssoCommandRun(issoPath)
					},
				},
				{
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "sqlite, s",
							Value:  "",
							Usage:  "path to sqlite file to export from",
							EnvVar: "SQLITE_FILE",
						},
						cli.StringFlag{
							Name:   "config, c",
							Value:  "",
							Usage:  "path to mouthful config file",
							EnvVar: "MOUTHFUL_CONFIG",
						},
					},
					Name:  "dynamo",
					Usage: "imports the given sqlite database to the dynamodb database pointed by the config provided",
					Action: func(c *cli.Context) error {
						sqlitePath := c.String("sqlite")
						configPath := c.String("config")
						return command.DynamoCommandRun(sqlitePath, configPath)
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
