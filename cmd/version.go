package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Get version number",
		Long:  `Get version number`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", Version)
		},
	}
	Version = "development"
)
