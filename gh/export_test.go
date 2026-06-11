package gh

// HiddenMarker exposes the unexported constant for use in external test packages.
var HiddenMarker = hiddenMarker

// NewTestClient constructs a Client with injected service mocks.
// Only available during test compilation.
func NewTestClient(issues issuesInterface, pr pullRequestsInterface, repos repositoriesInterface) *Client {
	return &Client{
		issuesService:       issues,
		pullRequestsService: pr,
		repositoriesService: repos,
	}
}
