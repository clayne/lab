package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/pkg/errors"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	"github.com/zaquestion/lab/internal/action"
	"github.com/zaquestion/lab/internal/git"
	lab "github.com/zaquestion/lab/internal/gitlab"
)

// ciCreateCmd represents the run command
var ciCreateCmd = &cobra.Command{
	Use:     "create [branch]",
	Aliases: []string{"run"},
	Short:   "Create a CI pipeline",
	Long: heredoc.Doc(`
		Run the CI pipeline for the given or current branch if none provided.
		This API uses your GitLab token to create CI pipelines

		Project will be inferred from branch if not provided

		Note: "lab ci create" differs from "lab ci trigger" which is a
		different API`),
	Example: heredoc.Doc(`
		lab ci create feature_branch
		lab ci create -p engineering/integration_tests master`),
	PersistentPreRun: labPersistentPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		pid, branch, err := getCIRunOptions(cmd, args)
		if err != nil {
			log.Fatal(err)
		}
		pipeline, err := lab.CICreate(pid, &gitlab.CreatePipelineOptions{Ref: &branch})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(pipeline.WebURL)
	},
}

var ciTriggerCmd = &cobra.Command{
	Use:   "trigger [branch]",
	Short: "Trigger a CI pipeline",
	Long: heredoc.Doc(`
		Runs a trigger for a CI pipeline on the given or current branch if none provided.
		This API supports variables and must be called with a trigger token or from within GitLab CI.

		Project will be inferred from branch if not provided

		Note: "lab ci trigger" differs from "lab ci create" which is a different API`),
	Example: heredoc.Doc(`
		lab ci trigger feature_branch
		lab ci trigger -p engineering/integration_tests master
		lab ci trigger -p engineering/integration_tests -v foo=bar master`),
	PersistentPreRun: labPersistentPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		pid, branch, err := getCIRunOptions(cmd, args)
		if err != nil {
			log.Fatal(err)
		}
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			log.Fatal(err)
		}
		vars, err := cmd.Flags().GetStringSlice("variable")
		if err != nil {
			log.Fatal(err)
		}
		ciVars, err := parseCIVariables(vars)
		if err != nil {
			log.Fatal(err)
		}
		pipeline, err := lab.CITrigger(pid, gitlab.RunPipelineTriggerOptions{
			Ref:       &branch,
			Token:     &token,
			Variables: ciVars,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(pipeline.WebURL)
	},
}

func getCIRunOptions(cmd *cobra.Command, args []string) (string, string, error) {
	branch, err := git.CurrentBranch()
	if err != nil {
		return "", "", err
	}
	if len(args) > 0 {
		branch = args[0]
	}

	var pid string

	remote, err := determineSourceRemote(branch)
	if err != nil {
		return "", "", err
	}
	rn, err := git.PathWithNamespace(remote)
	if err != nil {
		return "", "", err
	}
	pid = rn

	project, err := cmd.Flags().GetString("project")
	if err != nil {
		return "", "", err
	}
	if project != "" {
		_, err := lab.FindProject(project)
		if err != nil {
			return "", "", err
		}
		pid = project
	}
	return pid, branch, nil
}

func parseCIVariables(vars []string) (map[string]string, error) {
	variables := make(map[string]string)
	for _, v := range vars {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) < 2 {
			return nil, errors.Errorf("Invalid Variable: \"%s\", Variables must be in the format key=value", v)
		}
		variables[parts[0]] = parts[1]

	}
	return variables, nil
}

func init() {
	ciCreateCmd.Flags().StringP("project", "p", "", "project to create pipeline on")
	ciCmd.AddCommand(ciCreateCmd)
	carapace.Gen(ciCreateCmd).PositionalCompletion(
		action.Remotes(),
	)

	ciTriggerCmd.Flags().StringP("project", "p", "", "project to run pipeline trigger on")
	ciTriggerCmd.Flags().StringP("token", "t", os.Getenv("CI_JOB_TOKEN"), "pipeline trigger token, optional if run within GitLabCI")
	ciTriggerCmd.Flags().StringSliceP("variable", "v", []string{}, "variables to pass to pipeline")

	ciCmd.AddCommand(ciTriggerCmd)
	carapace.Gen(ciTriggerCmd).PositionalCompletion(
		action.RemoteBranches(-1),
	)
}
