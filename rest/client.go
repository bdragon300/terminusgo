package rest

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"github.com/bdragon300/tusgo"

	"github.com/bdragon300/terminusgo/srverror"
	"github.com/bdragon300/terminusgo/woql/bare"

	"github.com/dghubble/sling"
	"github.com/hashicorp/go-cleanhttp"
)

// TODO: move from rest to another place
// TODO: api_init.pl paths and filenames and headers
type Client struct {
	C          *sling.Sling
	baseAPIURL string
	implClient *http.Client
}

func NewClient(hostPath string) *Client {
	cl := &Client{C: sling.New()}
	cl.implClient = cleanhttp.DefaultPooledClient()
	cl.C.Client(cl.implClient) // TODO: passing context (Timeout for instance) by user to a request
	cl.baseAPIURL = hostPath + "/api"
	cl.C.Base(cl.baseAPIURL)

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
	return &DatabaseIntroducer{BaseIntroducer: BaseIntroducer{client: c}}
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

func (c *Client) Diffs() *DiffRequester {
	return &DiffRequester{Client: c}
}

func (c *Client) Files() *FilesIntroducer {
	u := path.Join(c.baseAPIURL, "files")
	if filesURL, err := url.Parse(u); err != nil {
		panic(err)
	} else {
		tusClient := tusgo.NewClient(c.implClient, filesURL)
		return &FilesIntroducer{BaseIntroducer: BaseIntroducer{client: c}, tusClient: tusClient}
	}
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

type ClientWOQLOptions struct {
	CommitAuthor  string
	CommitMessage string
	AllWitnesses  bool
}

// Query with empty context
func (c *Client) WOQL(ctx context.Context, buf *srverror.WOQLResponse, query bare.RawQuery, options *ClientWOQLOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	type commitInfo struct {
		Author  string `json:"author"`
		Message string `json:"message"`
	}
	body := struct {
		AllWitnesses bool          `json:"all_witnesses,omitempty"`
		CommitInfo   commitInfo    `json:"commit_info"`
		Query        bare.RawQuery `json:"query"`
	}{options.AllWitnesses, commitInfo{options.CommitAuthor, options.CommitMessage}, query}
	sl := c.C.BodyJSON(body).Post("woql")
	return doRequest(ctx, sl, buf)
}
