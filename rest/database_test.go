package rest_test

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/bdragon300/terminusgo/rest"
	"github.com/bdragon300/tusgo"
	"github.com/stretchr/testify/suite"
)

var excludeAssertionDatabaseFields = []string{"ID", "CreationDate"}

type DatabaseSuite struct {
	suite.Suite
	dbBaseURL  string
	testClient *rest.Client
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

func (s *DatabaseSuite) SetupSuite() {
	s.dbBaseURL = getDatabaseURL()
	s.testClient = rest.NewClient(nil, s.dbBaseURL).WithBasicAuth(DefaultUser, DefaultPassword)
}

func (s *DatabaseSuite) TearDownTest() {
	cleanUpTest(s.testClient)
}

func (s *DatabaseSuite) TestListAll() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)

	// Check if canceled context matters and context doesn't get to client
	buf := make([]rest.Database, 0)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := s.testClient.Databases().WithContext(ctx).ListAll(&buf)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.Require().Empty(buf)

	buf = make([]rest.Database, 0)
	resp, err = s.testClient.Databases().ListAll(&buf)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Len(buf, 1)
	assertSameTerminusObjects(s.Require(), rest.Database{
		ID:           "", // Skip
		Type:         "UserDatabase",
		Name:         testDatabaseName,
		Comment:      "",
		CreationDate: time.Time{}, // Skip
		Label:        testDatabaseLabel,
		State:        "finalized",
		Path:         testOrganizationName + "/" + testDatabaseName,
		Branches:     []string{"main"},
	}, buf[0], excludeAssertionDatabaseFields)

	resp, err = s.testClient.Databases().OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).Delete(testDatabaseName, nil)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())

	buf = make([]rest.Database, 0)
	resp, err = s.testClient.Databases().ListAll(&buf)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Empty(buf)
}

func (s *DatabaseSuite) TestListAllOwned() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)

	// Check if canceled context matters and context doesn't get to client
	buf := make([]rest.Database, 0)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := s.testClient.Databases().WithContext(ctx).ListAllOwned(&buf)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.Require().Empty(buf)

	buf = make([]rest.Database, 0)
	resp, err = s.testClient.Databases().ListAllOwned(&buf)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Empty(buf)

	dbOwnerClient := rest.NewClient(nil, s.dbBaseURL).WithBasicAuth(testUserName, DefaultPassword)

	buf = make([]rest.Database, 0)
	resp, err = dbOwnerClient.Databases().ListAllOwned(&buf)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Len(buf, 1)
	assertSameTerminusObjects(s.Require(), rest.Database{
		ID:           "", // Skip
		Type:         "UserDatabase",
		Name:         testDatabaseName,
		Comment:      "",
		CreationDate: time.Time{}, // Skip
		Label:        testDatabaseLabel,
		State:        "finalized",
		Path:         testOrganizationName + "/" + testDatabaseName,
		Branches:     []string{"main"},
	}, buf[0], excludeAssertionDatabaseFields)

	resp, err = s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		Delete(testDatabaseName, nil)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())

	buf = make([]rest.Database, 0)
	resp, err = dbOwnerClient.Databases().ListAllOwned(&buf)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Empty(buf)
}

func (s *DatabaseSuite) TestCRUD() {
	const anotherDbLabel = "second_label"
	const anotherDbComment = "test comment"

	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)

	testDB := rest.Database{
		ID:           "",
		Type:         "Database",
		Name:         testDatabaseName,
		Comment:      "",
		CreationDate: time.Time{},
		Label:        testDatabaseLabel,
		State:        "finalized",
		Path:         testOrganizationName + "/" + testDatabaseName,
		Branches:     []string{"main"},
	}
	s.assertHaveDatabases([]rest.Database{})
	s.assertDatabaseExists(false)

	// CREATE/GET
	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		WithContext(ctx).
		Create(testDatabaseName, testDatabaseLabel, nil)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.assertHaveDatabases([]rest.Database{})
	s.assertDatabaseExists(false)

	createDatabase(s.Require(), s.testClient)
	s.assertHaveDatabases([]rest.Database{testDB})
	s.assertDatabaseExists(true)

	// UPDATE/GET
	// Check if canceled context matters and context doesn't get to client
	// TODO: add testing other Update's options -- seems that testing needed the acquiring documents from db
	changedDb := testDB
	changedDb.Label = anotherDbLabel
	changedDb.Comment = anotherDbComment
	updateOpts := rest.DatabaseUpdateOptions{
		Schema:   true,
		Public:   false,
		Label:    anotherDbLabel,
		Comment:  anotherDbComment,
		Prefixes: nil,
	}
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	resp, err = s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		WithContext(ctx).
		Update(testDatabaseName, &updateOpts)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.assertHaveDatabases([]rest.Database{testDB})
	s.assertDatabaseExists(true)

	resp, err = s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		Update(testDatabaseName, &updateOpts)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.assertHaveDatabases([]rest.Database{changedDb})
	s.assertDatabaseExists(true)

	// DELETE/GET
	// Check if canceled context matters and context doesn't get to client
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	resp, err = s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		WithContext(ctx).
		Delete(testDatabaseName, nil)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.assertHaveDatabases([]rest.Database{testDB})
	s.assertDatabaseExists(true)

	resp, err = s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		Delete(testDatabaseName, nil)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.assertHaveDatabases([]rest.Database{})
	s.assertDatabaseExists(false)
}

func (s *DatabaseSuite) TestCreateOptions() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)

	// TODO: seems that testing needed the acquiring documents from db
}

// TODO: Test WOQL

func (s *DatabaseSuite) TestPrefixes() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)
	buf := rest.Prefix{}

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		WithContext(ctx).
		Prefixes(testDatabaseName, &buf)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)

	resp, err = s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		Prefixes(testDatabaseName, &buf)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Equal(rest.Prefix{
		Base:   "terminusdb:///data1",
		Schema: "terminusdb:///schema",
		Type:   "Context",
	}, buf)
}

func (s *DatabaseSuite) TestCommitLogAndOptions() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)

	excludeAssertionCommitFields := []string{"ID", "Identifier", "Instance", "Parent", "Schema", "Timestamp"}
	commits := []rest.Commit{{Type: "InitialCommit", Author: "system", Message: "create initial schema"}}
	buf := make([]rest.Commit, 0)

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		WithContext(ctx).
		CommitLog(testDatabaseName, &buf, nil)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)

	examples := [][]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}} // {offset,limit}
	for _, v := range examples {
		offset, limit := v[0], v[1]
		expect := commits[offset:limit]
		resp, err = s.testClient.Databases().
			OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
			CommitLog(testDatabaseName, &buf, &rest.DatabaseCommitLogOptions{Count: limit, Start: offset})
		s.Require().ErrorIs(err, context.Canceled)
		s.Require().Nil(resp)
		s.Require().Equal(len(expect), len(buf))
		for i := 0; i < len(expect); i++ {
			assertSameTerminusObjects(s.Require(), expect[i], buf[i], excludeAssertionCommitFields)
		}
	}
}

func (s *DatabaseSuite) TestOptimize() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		WithContext(ctx).
		Optimize(testDatabaseName)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)

	resp, err = s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		Optimize(testDatabaseName)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
}

// TODO: add Schema methods tests

func (s *DatabaseSuite) TestPack() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)

	buf := bytes.NewBuffer(make([]byte, 0))

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	written, resp, err := s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		WithContext(ctx).
		Pack(testDatabaseName, buf, nil)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.Require().Equal(0, buf.Len())
	s.Require().Equal(0, written)

	written, resp, err = s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		Pack(testDatabaseName, buf, nil)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Positive(buf.Len())
	s.Require().Equal(buf.Len(), written)
}

// TODO: add test for Pack options

func (s *DatabaseSuite) TestUnpackUpload() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)

	buf := bytes.NewBuffer(make([]byte, 0))

	_, resp, err := s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		Pack(testDatabaseName, buf, nil)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Positive(buf.Len())

	bufLen := buf.Len()

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	read, resp, err := s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		WithContext(ctx).
		UnpackUpload(testDatabaseName, buf)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.Require().Equal(0, buf.Len())
	s.Require().Equal(0, read)

	read, resp, err = s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		WithContext(ctx).
		UnpackUpload(testDatabaseName, buf)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Equal(bufLen, read)
}

func (s *DatabaseSuite) TestUnpackTusResource() {
	createUser(s.Require(), s.testClient)
	createOrganization(s.Require(), s.testClient)
	createDatabase(s.Require(), s.testClient)

	buf := bytes.NewBuffer(make([]byte, 0))

	_, resp, err := s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		Pack(testDatabaseName, buf, nil)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Positive(buf.Len())

	bufLen := buf.Len()

	// TUS upload file
	u, err := url.Parse(s.dbBaseURL + "/files")
	s.Require().NoError(err)
	upload := tusgo.Upload{}
	tusClient := tusgo.NewClient(nil, u)
	_, err = tusClient.CreateUpload(&upload, int64(bufLen), false, nil)
	s.Require().NoError(err)
	tusStream := tusgo.NewUploadStream(tusClient, &upload)
	_, err = io.Copy(tusStream, buf)
	s.Require().NoError(err)

	// Check if canceled context matters and context doesn't get to client
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	read, resp, err := s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		WithContext(ctx).
		UnpackTusResource(testDatabaseName, upload.Location)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
	s.Require().Equal(0, read)

	read, resp, err = s.testClient.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		UnpackTusResource(testDatabaseName, upload.Location)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Equal(bufLen, read)
}

func (s *DatabaseSuite) assertHaveDatabases(checkDbs []rest.Database) {
	listBuf := make([]rest.Database, 0)
	resp, err := s.testClient.Databases().ListAll(&listBuf)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Equal(len(checkDbs), len(listBuf))

	for i := 0; i < len(checkDbs); i++ {
		assertSameTerminusObjects(s.Require(), checkDbs[i], listBuf[i], excludeAssertionDatabaseFields)

		elemBuf := rest.Database{}
		resp, err = s.testClient.Databases().
			OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
			Get(checkDbs[i].Name, nil)
		s.Require().NoError(err)
		s.Require().True(resp.IsOK())
		assertSameTerminusObjects(s.Require(), checkDbs[i], elemBuf, excludeAssertionDatabaseFields)
	}
}

func (s *DatabaseSuite) assertDatabaseExists(exists bool) {
	ok, resp, err := s.testClient.Databases().OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).IsExists(testDatabaseName)
	s.Require().NoError(err)
	s.Require().Equal(resp.IsOK(), ok)
	s.Require().Equal(exists, ok)
}
