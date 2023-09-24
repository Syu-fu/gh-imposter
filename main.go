package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/shurcooL/githubv4"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gh imposter",
		Usage: "gh-imposter allows you to configure repositories based on a settings file.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "~/.config/gh-imposter/config.yaml",
				Usage:   "configuration file path",
			},
			&cli.PathFlag{
				Name:    "repository",
				Aliases: []string{"r"},
				Usage:   "specify the repository name in the format (default: current directory name)",
			},
			&cli.BoolFlag{
				Name:  "time",
				Usage: "display overtime in hh:mm format",
			},
		},
		Action: func(c *cli.Context) error {
			Imposter(c)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Imposter(c *cli.Context) {
	config := Config{}
	if err := ReadConfig(c.Path("config"), &config); err != nil {
		fmt.Println(err)
		return
	}
	gqlClient, err := api.DefaultGraphQLClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	owner, name, err := SelectUpdateRepository(c.String("repository"))
	fmt.Printf("owner: %v\n", owner)
	fmt.Printf("name: %v\n", name)
	if err != nil {
		fmt.Println(err)
		return
	}
	repoId, patterns, err := GetRepositoryInfo(gqlClient, owner, name)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := RegisterBranchProtectionRule(gqlClient, repoId, patterns, &config); err != nil {
		fmt.Println(err)
		return
	}
}

func SelectUpdateRepository(repo string) (string, string, error) {
	// if commandArgs exist, use it
	if repo != "" {
		repoArr := strings.Split(repo, "/")
		if len(repoArr) != 2 {
			return "", "", fmt.Errorf("repository name is not specified")
		}
		return repoArr[0], repoArr[1], nil
	}
	// if commandArgs not exist, use current directory info
	cr, err := repository.Current()
	if err != nil {
		// if both not exist, error
		return "", "", fmt.Errorf("repository name is not specified")
	}
	return cr.Owner, cr.Name, nil
}

func GetRepositoryInfo(gqlClient *api.GraphQLClient, owner string, name string) (string, []struct {
	ID      string
	Pattern string
}, error,
) {
	var query struct {
		Repository struct {
			ID                    string
			BranchProtectionRules struct {
				Nodes []struct {
					ID      string
					Pattern string
				}
			} `graphql:"branchProtectionRules(first: 100)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}
	if err := gqlClient.Query("GetRepoID", &query, variables); err != nil {
		return "", nil, err
	}
	var branchProtectionRules []struct {
		ID      string
		Pattern string
	}

	branchProtectionRules = append(branchProtectionRules, query.Repository.BranchProtectionRules.Nodes...)
	return query.Repository.ID, branchProtectionRules, nil
}

func RegisterBranchProtectionRule(gqlClient *api.GraphQLClient, repoId string, branchProtectionRules []struct {
	ID      string
	Pattern string
}, config *Config,
) error {
	for _, rule := range config.BranchProtectionRules {
		// if pattern is empty, skip
		if rule.Pattern == "" {
			continue
		}
		// if pattern is already exist, update
		// if pattern is not exist, create
		id := containsRule(branchProtectionRules, string(rule.Pattern))
		if id != "" {
			fmt.Printf("update %s\n", rule.Pattern)
			if err := UpdateBranchProtectionRule(gqlClient, id, rule); err != nil {
				return err
			}
		} else {
			fmt.Printf("create %s\n", rule.Pattern)
			if err := CreateBranchProtectionRule(gqlClient, repoId, rule); err != nil {
				return err
			}
		}
	}
	return nil
}

func UpdateBranchProtectionRule(gqlClient *api.GraphQLClient, ruleId string, rule BranchProtectionRule) error {
	var mutation struct {
		UpdateBranchProtectionRule struct {
			ClientMutationId githubv4.String
		} `graphql:"updateBranchProtectionRule(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": githubv4.UpdateBranchProtectionRuleInput{
			BranchProtectionRuleID:         ruleId,
			Pattern:                        githubv4.NewString(rule.Pattern),
			RequiresApprovingReviews:       rule.RequiresApprovingReviews,
			RequiredApprovingReviewCount:   rule.RequiredApprovingReviewCount,
			RequiresCommitSignatures:       rule.RequiresCommitSignatures,
			RequiresLinearHistory:          rule.RequiresLinearHistory,
			BlocksCreations:                rule.BlocksCreations,
			AllowsForcePushes:              rule.AllowsForcePushes,
			AllowsDeletions:                rule.AllowsDeletions,
			IsAdminEnforced:                rule.IsAdminEnforced,
			RequiresStatusChecks:           rule.RequiresStatusChecks,
			RequiresStrictStatusChecks:     rule.RequiresStrictStatusChecks,
			RequiresCodeOwnerReviews:       rule.RequiresCodeOwnerReviews,
			DismissesStaleReviews:          rule.DismissesStaleReviews,
			RestrictsReviewDismissals:      rule.RestrictsReviewDismissals,
			ReviewDismissalActorIDs:        rule.ReviewDismissalActorIDs,
			BypassPullRequestActorIDs:      rule.BypassPullRequestActorIDs,
			BypassForcePushActorIDs:        rule.BypassForcePushActorIDs,
			RestrictsPushes:                rule.RestrictsPushes,
			PushActorIDs:                   rule.PushActorIDs,
			RequiredStatusCheckContexts:    rule.RequiredStatusCheckContexts,
			RequiredStatusChecks:           rule.RequiredStatusChecks,
			RequiresDeployments:            rule.RequiresDeployments,
			RequiredDeploymentEnvironments: rule.RequiredDeploymentEnvironments,
			RequiresConversationResolution: rule.RequiresConversationResolution,
			RequireLastPushApproval:        rule.RequireLastPushApproval,
			LockBranch:                     rule.LockBranch,
			LockAllowsFetchAndMerge:        rule.LockAllowsFetchAndMerge,
			ClientMutationID:               rule.ClientMutationID,
		},
	}
	if err := gqlClient.Mutate("UpdateBranchProtectionRule", &mutation, variables); err != nil {
		return err
	}
	return nil
}

func CreateBranchProtectionRule(gqlClient *api.GraphQLClient, repoId string, rule BranchProtectionRule) error {
	var mutation struct {
		CreateBranchProtectionRuleInput struct {
			ClientMutationId githubv4.String
		} `graphql:"createBranchProtectionRule(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": githubv4.CreateBranchProtectionRuleInput{
			RepositoryID:                 repoId,
			Pattern:                      rule.Pattern,
			RequiresApprovingReviews:     rule.RequiresApprovingReviews,
			RequiredApprovingReviewCount: rule.RequiredApprovingReviewCount,
			RequiresCommitSignatures:     rule.RequiresCommitSignatures,
			RequiresLinearHistory:        rule.RequiresLinearHistory,
			BlocksCreations:              rule.BlocksCreations,
			AllowsForcePushes:            rule.AllowsForcePushes,
			AllowsDeletions:              rule.AllowsDeletions,

			IsAdminEnforced:      rule.IsAdminEnforced,
			RequiresStatusChecks: rule.RequiresStatusChecks,

			RequiresStrictStatusChecks:     rule.RequiresStrictStatusChecks,
			RequiresCodeOwnerReviews:       rule.RequiresCodeOwnerReviews,
			DismissesStaleReviews:          rule.DismissesStaleReviews,
			RestrictsReviewDismissals:      rule.RestrictsReviewDismissals,
			ReviewDismissalActorIDs:        rule.ReviewDismissalActorIDs,
			BypassPullRequestActorIDs:      rule.BypassPullRequestActorIDs,
			BypassForcePushActorIDs:        rule.BypassForcePushActorIDs,
			RestrictsPushes:                rule.RestrictsPushes,
			PushActorIDs:                   rule.PushActorIDs,
			RequiredStatusCheckContexts:    rule.RequiredStatusCheckContexts,
			RequiredStatusChecks:           rule.RequiredStatusChecks,
			RequiresDeployments:            rule.RequiresDeployments,
			RequiredDeploymentEnvironments: rule.RequiredDeploymentEnvironments,
			RequiresConversationResolution: rule.RequiresConversationResolution,
			RequireLastPushApproval:        rule.RequireLastPushApproval,
			LockBranch:                     rule.LockBranch,

			LockAllowsFetchAndMerge: rule.LockAllowsFetchAndMerge,
			ClientMutationID:        rule.ClientMutationID,
		},
	}
	if err := gqlClient.Mutate("CreateBranchProtectionRule", &mutation, variables); err != nil {
		return err
	}
	return nil
}

func containsRule(rules []struct {
	ID      string
	Pattern string
}, pattern string,
) string {
	for _, rule := range rules {
		if rule.Pattern == pattern {
			return rule.ID
		}
	}
	return ""
}
