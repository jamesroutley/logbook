package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	logbookDir string
	gitRepo    string
)

func init() {
	homedir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logbookDir = filepath.Join(homedir, "logbook")

	initCmd.Flags().StringVarP(&gitRepo, "git", "g", "", "Git repo to clone from")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises a new logbook",
	Long: `init initialises a new logbook directory, at ~/logbook. If you
have an existing logbook directory stored in git, use the -git flag to clone it.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		if _, err := os.Stat(logbookDir); !os.IsNotExist(err) {
			fmt.Printf("%s already exists\n", logbookDir)
			os.Exit(1)
		}
		if gitRepo != "" {
			cmd := exec.Command("git", "clone", gitRepo, logbookDir)
			output, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(output)
			return
		}
		os.Mkdir(logbookDir, os.ModePerm)
	},
}
