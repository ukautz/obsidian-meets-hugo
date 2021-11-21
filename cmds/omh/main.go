package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "omh"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "obsidian-root",
			Aliases:  []string{"O"},
			Required: true,
			Usage:    "Path to root of Obsidian Vault",
		},
		&cli.StringFlag{
			Name:     "hugo-root",
			Aliases:  []string{"H"},
			Required: true,
			Usage:    "Path to root of Hugo setup",
		},
		&cli.StringFlag{
			Name:    "sub-path",
			Aliases: []string{"p"},
			Usage:   "Sub-path used in Hugo setup below content and static",
			Value:   "obsidian",
		},
		&cli.StringSliceFlag{
			Name:    "include-tag",
			Aliases: []string{"i"},
			Usage:   "Tag to include (accept list - accepts all, if unset)",
		},
		&cli.StringSliceFlag{
			Name:    "exclude-tags",
			Aliases: []string{"e"},
			Usage:   "Tag to exclude (reject list - reject none, if unset)",
		},
	}
	app.Action = func(c *cli.Context) error {
		todo()
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func todo() {
	fmt.Println("Obsidian Meets Hugo")
	fmt.Println("  Command line tool to export (partial) Obsidian Vault to Hugo")
	fmt.Println("Input:")
	fmt.Println("  - Obsidian Directory: Path to root of Obsidian Vault")
	fmt.Println("  - Hugo Directory: Path to root of Hugo setup")
	fmt.Println("    - Sub-Path, default `obsidian`, used in `content/<sub-path>` and `static/<sub-path>`")
	fmt.Println("  - Optional Tag include list")
	fmt.Println("  - Optional Tag exclude list")
	fmt.Println("Execution:")
	fmt.Println("  - Find all Markdown files in Obsidian Directory and Subdirectories")
	fmt.Println("    - Copy and Transform from Obsidian Note into Hugo Page in `<hugo-root>/content/<sub-path>`")
	fmt.Println("      - Make file name snake-case")
	fmt.Println("      - Replace all internal links, so that they work in Hugo (point to snake case, respective sub-path in content)")
	fmt.Println("      - Replace all internal references to non-Markdown files with appropriate Markdown")
	fmt.Println("  - Find all none-Markdown files and copy them to `<hugo-root>/static/<sub-path>")
}
