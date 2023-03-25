package rest_test

import (
	"context"
	"testing"

	"github.com/bdragon300/terminusgo/rest"
	"github.com/stretchr/testify/suite"
)

type OrganizationSuite struct {
	suite.Suite
	dbBaseURL  string
	testClient *rest.Client
}

func TestOrganizationSuite(t *testing.T) {
	suite.Run(t, new(OrganizationSuite))
}

func (s *OrganizationSuite) SetupSuite() {
	s.dbBaseURL = getDatabaseURL()
	s.testClient = rest.NewClient(nil, s.dbBaseURL).WithBasicAuth(DefaultUser, DefaultPassword)
}

func (s *OrganizationSuite) TearDownTest() {
	cleanUpTest(s.testClient)
}

func (s *OrganizationSuite) TestCRUD() {
	testOrg := rest.Organization{
		ID:          "Organization/" + testOrganizationName,
		Type:        "Organization",
		Name:        testOrganizationName,
		DatabaseIDs: []string{"SystemDatabase/system"},
	}
	s.assertHaveOrganizations(nil)

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := s.testClient.Organizations().WithContext(ctx).Create(testOrganizationName)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.assertHaveOrganizations(nil)

	createOrganization(s.Require(), s.testClient)
	s.assertHaveOrganizations([]rest.Organization{testOrg})

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	resp, err = s.testClient.Organizations().WithContext(ctx).Delete(testOrganizationName)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.assertHaveOrganizations([]rest.Organization{testOrg})

	resp, err = s.testClient.Organizations().Delete(testOrganizationName)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.assertHaveOrganizations(nil)
}

func (s *OrganizationSuite) TestOrganizationPath() {
	testPath := rest.OrganizationPath{Organization: "foo|bar"}

	s.Require().Equal("foo%7Cbar", testPath.String())
	s.Require().Equal("test_action/foo%7Cbar", testPath.GetURL("test_action"))

	buf := rest.OrganizationPath{}
	buf.FromString("foo%7Cbar")
	s.Require().Equal(testPath, buf)
}

func (s *OrganizationSuite) assertHaveOrganizations(orgs []rest.Organization) {
	checkOrgs := []rest.Organization{{
		ID:          "Organization/admin",
		Type:        "Organization",
		Name:        "admin",
		DatabaseIDs: []string{"SystemDatabase/system"},
	}}
	checkOrgs = append(checkOrgs, orgs...)

	listBuf := make([]rest.Organization, 0)
	resp, err := s.testClient.Organizations().ListAll(&listBuf)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Equal(len(checkOrgs), len(listBuf))

	for i := 0; i < len(checkOrgs); i++ {
		assertSameTerminusObjects(s.Require(), checkOrgs[i], listBuf[i], nil)

		elemBuf := rest.Organization{}
		resp, err = s.testClient.Organizations().Get(checkOrgs[i].Name, &elemBuf)
		s.Require().NoError(err)
		s.Require().True(resp.IsOK())
		assertSameTerminusObjects(s.Require(), checkOrgs[i], elemBuf, nil)
	}
}
