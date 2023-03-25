package rest_test

import (
	"context"
	"testing"
	"time"

	"github.com/bdragon300/terminusgo/rest"
	"github.com/bdragon300/terminusgo/srverror"
	"github.com/stretchr/testify/suite"
)

var excludeAssertionUserFields = []string{"Capabilities"}

type UserSuite struct {
	suite.Suite
	dbBaseURL  string
	testClient *rest.Client
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) SetupSuite() {
	s.dbBaseURL = getDatabaseURL()
	s.testClient = rest.NewClient(nil, s.dbBaseURL).WithBasicAuth(DefaultUser, DefaultPassword)
}

func (s *UserSuite) TearDownTest() {
	cleanUpTest(s.testClient)
}

func (s *UserSuite) TestCRUDOnServer() {
	testUser := rest.User{
		ID:           "User/" + testUserName,
		Type:         "User",
		Name:         testUserName,
		Capabilities: nil,
	}
	s.assertHaveUsers([]rest.User{})

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := s.testClient.Users().OnServer().WithContext(ctx).Create(testUserName, "test_password")
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.assertHaveUsers([]rest.User{})

	resp, err = s.testClient.Users().OnServer().Create(testUserName, "test_password")
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.assertHaveUsers([]rest.User{testUser})

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	resp, err = s.testClient.Users().OnServer().WithContext(ctx).Delete(testUserName)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.assertHaveUsers([]rest.User{testUser})

	resp, err = s.testClient.Users().OnServer().WithContext(ctx).Delete(testUserName)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.assertHaveUsers([]rest.User{})
}

func (s *UserSuite) TestGetOptions() {
	s.createCapableUser()

	buf := rest.User{}
	resp, err := s.testClient.Users().OnServer().Get(testUserName, &buf, &rest.UserGetOptions{Capability: false})
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.assertUserHasSimpleCapabilities(buf)

	buf = rest.User{}
	resp, err = s.testClient.Users().OnServer().Get(testUserName, &buf, &rest.UserGetOptions{Capability: true})
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.assertUserHasObjectCapabilities(buf)
}

func (s *UserSuite) TestListAllOptions() {
	s.createCapableUser()

	buf := make([]rest.User, 0)
	resp, err := s.testClient.Users().OnServer().ListAll(&buf, &rest.UserListAllOptions{Capability: false})
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	basicUsers := []rest.User{{
		ID:   "User/admin",
		Type: "User",
		Capabilities: []srverror.Union[rest.UserCapability, []string]{{
			V0:       rest.UserCapability{},
			V1:       []string{"Capability/server_access"},
			Selector: 1,
		}},
		Name: "admin",
	}, {
		ID:           "User/anonymous",
		Type:         "User",
		Capabilities: nil,
		Name:         "anonymous",
	}}
	s.Require().Len(buf, 3)
	s.Require().EqualValues(basicUsers, buf[:2])
	s.assertUserHasSimpleCapabilities(buf[2])

	buf = make([]rest.User, 0)
	resp, err = s.testClient.Users().OnServer().ListAll(&buf, &rest.UserListAllOptions{Capability: true})
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Len(buf, 3)
	s.Require().EqualValues(basicUsers, buf[:2])
	s.assertUserHasObjectCapabilities(buf[2])
}

func (s *UserSuite) TestUpdateCapabilitiesOptions() {
	roles := []rest.Role{{
		ID:     "Role/" + testRoleName,
		Type:   "Role",
		Name:   testRoleName,
		Action: nil,
	}}
	capOptsDB := rest.UserUpdateCapabilitiesOptions{
		Scope: rest.Database{
			ID:           "",
			Type:         "UserDatabase",
			Name:         testDatabaseName,
			Comment:      "",
			CreationDate: time.Time{},
			Label:        testDatabaseLabel,
			State:        "finalized",
			Path:         testOrganizationName + "/" + testDatabaseName,
			Branches:     []string{"main"},
		},
		Roles:     roles,
		Operation: rest.UserCapabilitiesGrant,
	}
	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := s.testClient.Users().OnServer().WithContext(ctx).UpdateCapabilities(testUserName, &capOptsDB)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)

	s.createCapableUser()
	createDatabase(s.Require(), s.testClient)

	db := rest.Database{}
	resp, err = s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		Get(testDatabaseName, &db)
	s.Require().Error(err)
	s.Require().Equal("HTTP 404: Database does not exist, or you do not have permission", resp.String())

	resp, err = s.testClient.Users().OnServer().UpdateCapabilities(testUserName, &capOptsDB)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())

	u := rest.User{}
	resp, err = s.testClient.Users().OnServer().Get(testUserName, &u, &rest.UserGetOptions{Capability: true})
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Equal(testOrganizationName, extractField[string](u.Capabilities[0].V0.Scope, "Name"))
	s.Require().Equal(testDatabaseName, extractField[string](u.Capabilities[1].V0.Scope, "Name"))

	resp, err = s.testClient.Users().OnServer().UpdateCapabilities(testUserName, &rest.UserUpdateCapabilitiesOptions{
		Scope: rest.Organization{
			ID:          "Organization/" + testOrganizationName,
			Type:        "Organization",
			Name:        testOrganizationName,
			DatabaseIDs: nil,
		},
		Roles:     roles,
		Operation: rest.UserCapabilitiesRevoke,
	})
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())

	resp, err = s.testClient.Users().OnServer().UpdateCapabilities(testUserName, &rest.UserUpdateCapabilitiesOptions{
		Scope:     db,
		Roles:     roles,
		Operation: rest.UserCapabilitiesRevoke,
	})
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())

	resp, err = s.testClient.Users().OnServer().Get(testUserName, &u, &rest.UserGetOptions{Capability: true})
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Empty(u.Capabilities)
}

func (s *UserSuite) createCapableUser() {
	createOrganization(s.Require(), s.testClient)
	createUser(s.Require(), s.testClient)
	createRole(s.Require(), s.testClient)
	resp, err := s.testClient.Users().OnServer().UpdateCapabilities(testUserName, &rest.UserUpdateCapabilitiesOptions{
		Scope: rest.Organization{
			ID:          "Organization/" + testOrganizationName,
			Type:        "Organization",
			Name:        testOrganizationName,
			DatabaseIDs: nil,
		},
		Roles: []rest.Role{{
			ID:     "Role/" + testRoleName,
			Type:   "Role",
			Name:   testRoleName,
			Action: nil,
		}},
		Operation: rest.UserCapabilitiesGrant,
	})
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
}

func (s *UserSuite) assertHaveUsers(users []rest.User) {
	checkUsers := []rest.User{{
		ID:   "User/admin",
		Type: "User",
		Capabilities: []srverror.Union[rest.UserCapability, []string]{{
			V0:       rest.UserCapability{},
			V1:       []string{"Capability/server_access"},
			Selector: 1,
		}},
		Name: "admin",
	}, {
		ID:           "User/anonymous",
		Type:         "User",
		Capabilities: nil,
		Name:         "anonymous",
	}}
	checkUsers = append(checkUsers, users...)

	listBuf := make([]rest.User, 0)
	resp, err := s.testClient.Users().OnServer().ListAll(&listBuf, nil)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Equal(len(checkUsers), len(listBuf))

	for i := 0; i < len(checkUsers); i++ {
		assertSameTerminusObjects(s.Require(), checkUsers[i], listBuf[i], excludeAssertionUserFields)

		elemBuf := rest.User{}
		resp, err = s.testClient.Users().OnServer().Get(checkUsers[i].Name, &elemBuf, nil)
		s.Require().NoError(err)
		s.Require().True(resp.IsOK())
		assertSameTerminusObjects(s.Require(), checkUsers[i], elemBuf, excludeAssertionUserFields)
	}
}

func (s *UserSuite) assertUserHasSimpleCapabilities(u rest.User) {
	assertSameTerminusObjects(s.Require(), rest.User{
		ID:           "User/" + testUserName,
		Type:         "User",
		Name:         testUserName,
		Capabilities: nil, // Skip
	}, u, excludeAssertionUserFields)
	s.Require().Equal(1, u.Capabilities[0].Selector)
	s.Require().Len(u.Capabilities[0].V1, 1)
	s.Require().Regexp("^Capability/[0-9a-f]+$", u.Capabilities[0].V1[0])
}

func (s *UserSuite) assertUserHasObjectCapabilities(u rest.User) {
	assertSameTerminusObjects(s.Require(), rest.User{
		ID:           "User/" + testUserName,
		Type:         "User",
		Name:         testUserName,
		Capabilities: nil, // Skip
	}, u, excludeAssertionUserFields)

	s.Require().Equal(0, u.Capabilities[0].Selector)
	assertSameTerminusObjects(s.Require(), rest.UserCapability{
		ID:   "",
		Type: "Capability",
		Role: []rest.Role{
			{
				ID:   "Role/" + testRoleName,
				Type: "Role",
				Name: testRoleName,
				Action: []rest.RoleAction{
					rest.RoleActionCreateDatabase,
					rest.RoleActionDeleteDatabase,
				},
			},
		},
		Scope: "Organization/" + testOrganizationName,
	}, u.Capabilities[0].V0, excludeAssertionFields)
	s.Require().Regexp("^Capability/[0-9a-f]+$", u.Capabilities[0].V0.ID)
}
