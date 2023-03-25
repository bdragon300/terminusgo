package rest_test

import (
	"context"
	"testing"

	"github.com/bdragon300/terminusgo/rest"
	"github.com/stretchr/testify/suite"
)

type RepoSuite struct {
	suite.Suite
	dbBaseURL  string
	testClient *rest.Client
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}

func (s *RepoSuite) SetupSuite() {
	s.dbBaseURL = getDatabaseURL()
	s.testClient = rest.NewClient(nil, s.dbBaseURL).WithBasicAuth(DefaultUser, DefaultPassword)
}

func (s *RepoSuite) TearDownTest() {
	cleanUpTest(s.testClient)
}

func (s *RepoSuite) TestOptimize() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := s.testClient.Repos().OnDatabase(rest.DatabasePath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
	}).WithContext(ctx).Optimize("local")
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)

	resp, err = s.testClient.Repos().OnDatabase(rest.DatabasePath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
	}).Optimize("local")
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
}

// TODO: make test for Schema* methods
