package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	lab "github.com/zaquestion/lab/internal/gitlab"
)

// ciLintCmd represents the lint command
var ciLintCmd = &cobra.Command{
	Use:              "lint",
	Short:            "Validate .gitlab-ci.yml against GitLab",
	Long:             ``,
	PersistentPreRun: LabPersistentPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		path := ".gitlab-ci.yml"
		if len(args) == 1 {
			path = args[0]
		}
		b, err := ioutil.ReadFile(path)
		if !os.IsNotExist(err) && err != nil {
			log.Fatal(err)
		}
		ok, err := lab.Lint(string(b))
		if !ok || err != nil {
			log.Fatal(errors.Wrap(err, "ci yaml invalid"))
		}
		fmt.Println("Valid!")
	},
}

func init() {
	ciCmd.AddCommand(ciLintCmd)
	carapace.Gen(ciLintCmd).PositionalCompletion(
		carapace.ActionFiles(".gitlab-ci.yml"),
	)
}
