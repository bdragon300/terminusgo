package rest_test

import (
	"context"
	"testing"

	"github.com/bdragon300/terminusgo/rest"
	"github.com/bdragon300/terminusgo/srverror"
	"github.com/stretchr/testify/suite"
)

type ClientSuite struct {
	suite.Suite
	dbBaseURL  string
	testClient *rest.Client
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func (s *ClientSuite) SetupSuite() {
	s.dbBaseURL = getDatabaseURL()
	s.testClient = rest.NewClient(nil, s.dbBaseURL).WithBasicAuth(DefaultUser, DefaultPassword)
}

func (s *ClientSuite) TestPing() {
	ctx := context.Background()
	resp, err := s.testClient.Ping(ctx)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())

	// Check if canceled context matters
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err = s.testClient.Ping(ctx)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
}

func (s *ClientSuite) TestVersionInfo() {
	ctx := context.Background()
	buf := rest.TerminusVersionInfo{}
	resp, err := s.testClient.VersionInfo(ctx, &buf)
	s.Require().NoError(err)
	s.Require().True(resp.IsOK())
	s.Require().Equal(resp.(srverror.TerminusOkResponse).APIFields["success"], "api:success")

	// Check if canceled context matters
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err = s.testClient.VersionInfo(ctx, &buf)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Nil(resp)
}

// TODO: add WOQL tests
