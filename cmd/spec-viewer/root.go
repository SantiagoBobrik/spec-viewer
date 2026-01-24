package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	port   string
	folder string
)

var rootCmd = &cobra.Command{
	Use:   "spec-viewer",
	Short: "A live spec viewer for your markdown files",
	Long: `Spec Viewer is a CLI tool that serves your local markdown specifications
as a live-reloading website. It watches for changes in your folder
and updates the browser automatically.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
