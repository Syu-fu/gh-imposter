package main

import "github.com/shurcooL/githubv4"

type Config struct {
	BranchProtectionRules []BranchProtectionRule `yaml:"branchProtectionRules"`
}

type BranchProtectionRule struct {
	Pattern                        githubv4.String                      `yaml:"pattern"`
	RequiresApprovingReviews       *githubv4.Boolean                    `yaml:"requiresApprovingReviews"`
	RequiredApprovingReviewCount   *githubv4.Int                        `yaml:"requiredApprovingReviewCount"`
	RequiresCommitSignatures       *githubv4.Boolean                    `yaml:"requiresCommitSignatures"`
	RequiresLinearHistory          *githubv4.Boolean                    `yaml:"requiresLinearHistory"`
	BlocksCreations                *githubv4.Boolean                    `yaml:"blocksCreations"`
	AllowsForcePushes              *githubv4.Boolean                    `yaml:"allowsForcePushes"`
	AllowsDeletions                *githubv4.Boolean                    `yaml:"allowsDeletions"`
	IsAdminEnforced                *githubv4.Boolean                    `yaml:"isAdminEnforced"`
	RequiresStatusChecks           *githubv4.Boolean                    `yaml:"requiresStatusChecks"`
	RequiresStrictStatusChecks     *githubv4.Boolean                    `yaml:"requiresStrictStatusChecks"`
	RequiresCodeOwnerReviews       *githubv4.Boolean                    `yaml:"requiresCodeOwnerReviews"`
	DismissesStaleReviews          *githubv4.Boolean                    `yaml:"dismissesStaleReviews"`
	RestrictsReviewDismissals      *githubv4.Boolean                    `yaml:"restrictsReviewDismissals"`
	ReviewDismissalActorIDs        *[]githubv4.ID                       `yaml:"reviewDismissalActorIDs"`
	BypassPullRequestActorIDs      *[]githubv4.ID                       `yaml:"bypassPullRequestActorIDs"`
	BypassForcePushActorIDs        *[]githubv4.ID                       `yaml:"bypassForcePushActorIDs"`
	RestrictsPushes                *githubv4.Boolean                    `yaml:"restrictsPushes"`
	PushActorIDs                   *[]githubv4.ID                       `yaml:"pushActorIDs"`
	RequiredStatusCheckContexts    *[]githubv4.String                   `yaml:"requiredStatusCheckContexts"`
	RequiredStatusChecks           *[]githubv4.RequiredStatusCheckInput `yaml:"requiredStatusChecks"`
	RequiresDeployments            *githubv4.Boolean                    `yaml:"requiresDeployments"`
	RequiredDeploymentEnvironments *[]githubv4.String                   `yaml:"requiredDeploymentEnvironments"`
	RequiresConversationResolution *githubv4.Boolean                    `yaml:"requiresConversationResolution"`
	RequireLastPushApproval        *githubv4.Boolean                    `yaml:"requireLastPushApproval"`
	LockBranch                     *githubv4.Boolean                    `yaml:"lockBranch"`
	LockAllowsFetchAndMerge        *githubv4.Boolean                    `yaml:"lockAllowsFetchAndMerge"`
	ClientMutationID               *githubv4.String                     `yaml:"clientMutationID"`
}
