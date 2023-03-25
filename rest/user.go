package rest

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/bdragon300/terminusgo/srverror"
)

type UserCapability struct {
	ID    string         `json:"@id"`
	Type  string         `json:"@type"`
	Role  []Role         `json:"role"`
	Scope TerminusObject `json:"scope"`
}

type User struct {
	ID           string                                     `json:"@id"`
	Type         string                                     `json:"@type"`
	Name         string                                     `json:"name"`
	Capabilities []srverror.Union[UserCapability, []string] `json:"capability"`
}

type UserIntroducer BaseIntroducer

func (ui *UserIntroducer) OnOrganization(path OrganizationPath) *UserRequester {
	return &UserRequester{Client: ui.client, path: path}
}

func (ui *UserIntroducer) OnServer() *UserRequester {
	return &UserRequester{Client: ui.client, path: nil}
}

type UserRequester BaseRequester

func (ur *UserRequester) WithContext(ctx context.Context) *UserRequester {
	r := *ur
	r.ctx = ctx
	return &r
}

type UserListAllOptions struct {
	Capability bool `url:"capability" default:"false"` // True to expand `capability` into "Capability" schema item list, False to get them as string list
}

func (ur *UserRequester) ListAll(buf *[]User, options *UserListAllOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := ur.Client.C.QueryStruct(options).Get(ur.getURL(""))
	return doRequest(ur.ctx, sl, buf)
}

type UserGetOptions struct {
	Capability bool `url:"capability" default:"false"` // True to expand `capability` into "Capability" schema item list, False to get them as string list
}

func (ur *UserRequester) Get(name string, buf *User, options *UserGetOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := ur.Client.C.QueryStruct(options).Get(ur.getURL(name))
	return doRequest(ur.ctx, sl, buf)
}

func (ur *UserRequester) Create(name, password string) (response TerminusResponse, err error) {
	body := struct {
		Name     string `json:"name"`
		Password string `json:"password,omitempty"`
	}{name, password}
	sl := ur.Client.C.BodyJSON(body).Post("users")
	return doRequest(ur.ctx, sl, nil)
}

func (ur *UserRequester) UpdatePassword(name, password string) (response TerminusResponse, err error) {
	body := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{name, password}
	sl := ur.Client.C.BodyJSON(body).Put("users")
	return doRequest(ur.ctx, sl, nil)
}

type UserUpdateCapabilitiesOptions struct {
	Scope     TerminusObject
	Roles     []Role
	Operation UserCapabilitiesOperation `default:"revoke"`
}

type UserCapabilitiesOperation string

const (
	UserCapabilitiesGrant  UserCapabilitiesOperation = "grant"
	UserCapabilitiesRevoke UserCapabilitiesOperation = "revoke"
)

func (ur *UserRequester) UpdateCapabilities(name string, options *UserUpdateCapabilitiesOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	roles := make([]string, 0)
	for _, v := range options.Roles {
		roles = append(roles, v.Name)
	}
	scopeType := extractField[string](options.Scope, "Type")
	if strings.HasSuffix(scopeType, "Database") {
		scopeType = "database"
	}
	body := struct {
		Scope     string   `json:"scope"`
		ScopeType string   `json:"scope_type"`
		Roles     []string `json:"roles"`
		User      string   `json:"user"`
		Operation string   `json:"operation"`
	}{extractField[string](options.Scope, "Name"), scopeType, roles, name, string(options.Operation)}
	sl := ur.Client.C.BodyJSON(body).Post("capabilities")
	return doRequest(ur.ctx, sl, nil)
}

func (ur *UserRequester) Delete(name string) (response TerminusResponse, err error) {
	sl := ur.Client.C.Delete(ur.getURL(name))
	return doRequest(ur.ctx, sl, nil)
}

func (ur *UserRequester) getURL(objectID string) string {
	org := ""
	action := "users"
	if ur.path != nil {
		path := ur.path.(OrganizationPath)
		org = path.Organization
		action = "organizations"
	}
	return UserPath{
		Organization: org,
		User:         objectID,
	}.GetURL(action)
}

type UserPath struct {
	Organization, User string
}

func (up UserPath) GetURL(action string) string {
	return fmt.Sprintf("%s/%s", action, up.String())
}

func (up UserPath) String() string {
	if up.Organization == "" {
		return url.PathEscape(up.User)
	}
	return fmt.Sprintf("%s/users/%s", url.PathEscape(up.Organization), url.PathEscape(up.User))
}

func (up UserPath) FromString(s string) UserPath {
	parts := strings.SplitN(s, "/", 3)
	if len(parts) < 3 {
		panic(fmt.Sprintf("too short path %q", s))
	}
	parts = append(parts[:1], parts[2:]...) // Cut "users" part
	res := UserPath{}
	fillUnescapedStringFields(parts, &res)
	return res
}
