package resource

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/shurcooL/githubv4"
)

// Source represents the configuration for the resource.
type Source struct {
	Repository              string                      `json:"repository"`
	AccessToken             string                      `json:"access_token"`
	V3Endpoint              string                      `json:"v3_endpoint"`
	V4Endpoint              string                      `json:"v4_endpoint"`
	Paths                   []string                    `json:"paths"`
	IgnorePaths             []string                    `json:"ignore_paths"`
	DisableCISkip           bool                        `json:"disable_ci_skip"`
	DisableGitLFS           bool                        `json:"disable_git_lfs"`
	SkipSSLVerification     bool                        `json:"skip_ssl_verification"`
	DisableForks            bool                        `json:"disable_forks"`
	IgnoreDrafts            bool                        `json:"ignore_drafts"`
	GitCryptKey             string                      `json:"git_crypt_key"`
	BaseBranch              string                      `json:"base_branch"`
	RequiredReviewApprovals int                         `json:"required_review_approvals"`
	Labels                  []string                    `json:"labels"`
	States                  []githubv4.PullRequestState `json:"states"`
}

// Validate the source configuration.
func (s *Source) Validate() error {
	if s.AccessToken == "" {
		return errors.New("access_token must be set")
	}
	if s.Repository == "" {
		return errors.New("repository must be set")
	}
	if s.V3Endpoint != "" && s.V4Endpoint == "" {
		return errors.New("v4_endpoint must be set together with v3_endpoint")
	}
	if s.V4Endpoint != "" && s.V3Endpoint == "" {
		return errors.New("v3_endpoint must be set together with v4_endpoint")
	}
	for _, state := range s.States {
		switch state {
		case githubv4.PullRequestStateOpen:
		case githubv4.PullRequestStateClosed:
		case githubv4.PullRequestStateMerged:
		default:
			return errors.New(fmt.Sprintf("states value \"%s\" must be one of: OPEN, MERGED, CLOSED", state))
		}
	}
	return nil
}

// Metadata output from get/put steps.
type Metadata []*MetadataField

// Add a MetadataField to the Metadata.
func (m *Metadata) Add(name, value string) {
	*m = append(*m, &MetadataField{Name: name, Value: value})
}

// MetadataField ...
type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Version communicated with Concourse.
type Version struct {
	PR                  string                    `json:"pr"`
	Commit              string                    `json:"commit"`
	UpdatedTime         time.Time                 `json:"updated,omitempty"`
	ApprovedReviewCount string                    `json:"approved_review_count"`
	State               githubv4.PullRequestState `json:"state"`
}

// NewVersion constructs a new Version.
func NewVersion(p *PullRequest) Version {
	return Version{
		PR:                  strconv.Itoa(p.Number),
		Commit:              p.Tip.OID,
		UpdatedTime:         p.UpdatedAt.Time,
		ApprovedReviewCount: strconv.Itoa(p.ApprovedReviewCount),
		State:               p.State,
	}
}

// PullRequest represents a pull request and includes the tip (commit).
type PullRequest struct {
	PullRequestObject
	Tip                 CommitObject
	ApprovedReviewCount int
	Labels              []LabelObject
}

// PullRequestObject represents the GraphQL commit node.
// https://developer.github.com/v4/object/pullrequest/
type PullRequestObject struct {
	ID          string
	Number      int
	Title       string
	URL         string
	BaseRefName string
	HeadRefName string
	Repository  struct {
		URL string
	}
	IsCrossRepository bool
	IsDraft           bool
	State             githubv4.PullRequestState
	UpdatedAt         githubv4.DateTime
}

// CommitObject represents the GraphQL commit node.
// https://developer.github.com/v4/object/commit/
type CommitObject struct {
	ID            string
	OID           string
	CommittedDate githubv4.DateTime
	Message       string
	Author        struct {
		User struct {
			Login string
		}
		Email string
	}
}

// ChangedFileObject represents the GraphQL FilesChanged node.
// https://developer.github.com/v4/object/pullrequestchangedfile/
type ChangedFileObject struct {
	Path string
}

// LabelObject represents the GraphQL label node.
// https://developer.github.com/v4/object/label
type LabelObject struct {
	Name string
}
