package rest

import (
	"context"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/hashicorp/go-cleanhttp"
)

// TODO: move from rest to another place
// TODO: api_init.pl paths and filenames and headers
type Client struct {
	C          *sling.Sling
	implClient *http.Client
}

func NewClient(hostPath string) *Client {
	cl := &Client{C: sling.New()}
	cl.implClient = cleanhttp.DefaultPooledClient()
	cl.C.Client(cl.implClient) // TODO: passing context (Timeout for instance) by user to a request
	cl.C.Base(hostPath + "/api")

	return cl
}

func (c *Client) WithJWTAuth(jwtToken string) *Client {
	c.C.Set("Authorization", "Bearer "+jwtToken)
	return c
}

func (c *Client) WithAPITokenAuth(token string) *Client {
	c.C.Set("API_TOKEN", token) // TODO: check if it is actually used in db
	return c
}

func (c *Client) Organizations() *OrganizationRequester {
	return &OrganizationRequester{Client: c}
}

func (c *Client) Databases() *DatabaseIntroducer {
	return &DatabaseIntroducer{client: c}
}

func (c *Client) Repos() *RepoIntroducer {
	return &RepoIntroducer{client: c}
}

func (c *Client) Branches() *BranchIntroducer {
	return &BranchIntroducer{client: c}
}

func (c *Client) Commits() *CommitIntroducer {
	return &CommitIntroducer{client: c}
}

func (c *Client) Users() *UserIntroducer {
	return &UserIntroducer{client: c}
}

func (c *Client) GenericDocuments() *DocumentIntroducer[GenericDocument] {
	return &DocumentIntroducer[GenericDocument]{client: c}
}

func (c *Client) Remotes() *RemoteIntroducer {
	return &RemoteIntroducer{client: c}
}

func (c *Client) Roles() *RoleRequester {
	return &RoleRequester{Client: c}
}

func (c *Client) Ping(ctx context.Context) (response TerminusResponse, err error) {
	sl := c.C.Get("ok")
	return doRequest(ctx, sl, nil)
}

type VersionInfo struct {
	Version string `json:"version"`
	GitHash string `json:"git_hash,omitempty"`
}

type TerminusVersionInfo struct {
	Authority       string      `json:"authority"`
	Storage         VersionInfo `json:"storage"`
	TerminusDB      VersionInfo `json:"terminusdb"`
	TerminusDBStore VersionInfo `json:"terminusdb_store"`
}

func (c *Client) VersionInfo(ctx context.Context, buf *TerminusVersionInfo) (response TerminusResponse, err error) {
	var respBuf struct {
		Info *TerminusVersionInfo `json:"api:info"`
	}
	sl := c.C.Get("info")
	response, err = doRequest(ctx, sl, &respBuf)
	if err != nil {
		return
	}
	*buf = *respBuf.Info
	return
}
