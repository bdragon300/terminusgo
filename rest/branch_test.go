package rest_test

import (
	"context"
	"strings"
	"testing"

	"github.com/bdragon300/terminusgo/rest"
	"github.com/bdragon300/terminusgo/srverror"
	"github.com/stretchr/testify/suite"
)

var excludeAssertionBranchFields = []string{"Head"}

type BranchSuite struct {
	suite.Suite
	dbBaseURL  string
	testClient *rest.Client
}

func TestBranchSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *BranchSuite) SetupSuite() {
	s.dbBaseURL = getDatabaseURL()
	s.testClient = rest.NewClient(nil, s.dbBaseURL).WithBasicAuth(DefaultUser, DefaultPassword)
}

func (s *BranchSuite) TearDownTest() {
	cleanUpTest(s.testClient)
}

func (s *BranchSuite) TestCRUD() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)

	testBranch := rest.Branch{
		ID:   "Branch/" + testBranchName,
		Type: "Branch",
		Name: testBranchName,
		Head: "",
	}
	s.assertHaveBranches(nil)

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := s.testClient.Branches().OnRepo(rest.RepoPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
	}).WithContext(ctx).Create(testBranchName, nil)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.assertHaveBranches(nil)

	createBranch(s.Require(), s.testClient)
	s.assertHaveBranches([]rest.Branch{testBranch})

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	resp, err = s.testClient.Branches().OnRepo(rest.RepoPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
	}).WithContext(ctx).Delete(testBranchName)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.assertHaveBranches([]rest.Branch{testBranch})

	resp, err = s.testClient.Branches().OnRepo(rest.RepoPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
	}).Delete(testBranchName)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.assertHaveBranches(nil)
}

// TODO: test Create options

func (s *BranchSuite) TestSquashOptions() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)
	createBranch(s.Require(), s.testClient)

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := s.testClient.Branches().OnRepo(rest.RepoPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
	}).WithContext(ctx).Squash(testBranchName, &rest.BranchSquashOptions{
		Author:  testCommitAuthor,
		Message: "test_message",
	})
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)

	resp, err = s.testClient.Branches().OnRepo(rest.RepoPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
	}).Squash(testBranchName, &rest.BranchSquashOptions{
		Author:  testCommitAuthor,
		Message: "test_message",
	})
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	commitPath := strings.Split(resp.(srverror.TerminusOkResponse).APIFields["commit"].(string), "/")
	commitID := commitPath[4]

	// Check if squashed commit actually appeared
	buf := make(rest.GenericDocument)
	resp, err = s.testClient.GenericDocuments().OnBranch(rest.BranchPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
		Branch:       rest.BranchCommits,
	}).Get(commitID, &buf, &rest.DocumentGetOptions{
		CompressIDs: true,
		Type:        "ValidCommit",
		Unfold:      false,
		GraphType:   rest.GraphTypeInstance,
		Prefixed:    false,
	})
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Equal(commitID, buf["identifier"])
	s.Require().Equal(testCommitAuthor, buf["author"])
	s.Require().Equal("test_message", buf["message"])
}

func (s *BranchSuite) TestReset() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)
	createBranch(s.Require(), s.testClient)
	setupSchema(s.Require(), s.testClient)

	doc1 := rest.GenericDocument{
		"@type":      "child1",
		"enumfield1": "id1",
		"num1":       "67",
		"str1":       "thery",
	}
	insertDocument(s.Require(), s.testClient, doc1)
	doc2 := rest.GenericDocument{
		"@type":      "child1",
		"enumfield1": "id2",
		"num1":       "463",
		"str1":       "something",
	}
	insertDocument(s.Require(), s.testClient, doc2)

	commits := make([]rest.Commit, 0)
	resp, err := s.testClient.Branches().OnRepo(rest.RepoPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
	}).CommitLog(testBranchName, &commits, nil)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Len(commits, 2)
	targetCommitID := commits[0].Identifier

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err = s.testClient.Branches().OnRepo(rest.RepoPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
	}).WithContext(ctx).Reset(testBranchName, commits[0])
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)

	resp, err = s.testClient.Branches().OnRepo(rest.RepoPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
	}).Reset(testBranchName, commits[0])
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())

	commits = make([]rest.Commit, 0)
	resp, err = s.testClient.Branches().OnRepo(rest.RepoPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
	}).CommitLog(testBranchName, &commits, nil)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Len(commits, 1)
	s.Require().Equal(targetCommitID, commits[0].Identifier)
}

func (s *BranchSuite) TestApply() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)
	createBranch(s.Require(), s.testClient)
	setupSchema(s.Require(), s.testClient)

	doc1 := rest.GenericDocument{
		"@type":      "child1",
		"enumfield1": "id1",
		"num1":       "67",
		"str1":       "thery",
	}
	insertDocument(s.Require(), s.testClient, doc1)
	doc2 := rest.GenericDocument{
		"@type":      "child1",
		"enumfield1": "id2",
		"num1":       "463",
		"str1":       "something",
	}
	insertDocument(s.Require(), s.testClient, doc2)

	commits := make([]rest.Commit, 0)
	resp, err := s.testClient.Branches().OnRepo(rest.RepoPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
	}).CommitLog(testBranchName, &commits, nil)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Len(commits, 2)

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err = s.testClient.Branches().OnRepo(rest.RepoPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
	}).WithContext(ctx).Apply(testBranchName, commits[1].Identifier, commits[0].Identifier, nil)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	// TODO: create a test case for Apply method
}

func (s *BranchSuite) assertHaveBranches(branches []rest.Branch) {
	repo := rest.RepoPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
	}
	checkBranches := []rest.Branch{{
		ID:   "Branch/main",
		Type: "Branch",
		Name: "main",
		Head: "",
	}}
	checkBranches = append(checkBranches, branches...)

	listBuf := make([]rest.Branch, 0)
	resp, err := s.testClient.Branches().OnRepo(repo).ListAll(&listBuf)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Equal(len(checkBranches), len(listBuf))

	for i := 0; i < len(checkBranches); i++ {
		assertSameTerminusObjects(s.Require(), checkBranches[i], listBuf[i], excludeAssertionBranchFields)

		elemBuf := rest.Branch{}
		resp, err = s.testClient.Branches().OnRepo(repo).Get(checkBranches[i].ID, &elemBuf)
		s.Require().NoError(err)
		s.Require().True(resp.IsOK())
		assertSameTerminusObjects(s.Require(), checkBranches[i], elemBuf, excludeAssertionBranchFields)
	}
}

