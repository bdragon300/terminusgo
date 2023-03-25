package rest_test

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/bdragon300/terminusgo/rest"
	"github.com/stretchr/testify/require"
)

const (
	DatabaseURLVariable = "DATABASE_URL"
	DefaultDatabaseURL  = "http://localhost:6363"
	DefaultUser         = "admin"
	DefaultPassword     = "admin"
)

const (
	testOrganizationName = "test_organization"
	testUserName         = "test_user"
	testRoleName         = "test_role"
	testDatabaseName     = "test_database"
	testDatabaseLabel    = "TestDatabase"
	testBranchName       = "test_branch"
	testCommitAuthor     = "terminusgo tests"
)

func createOrganization(r *require.Assertions, c *rest.Client) {
	resp, err := c.Organizations().Create(testOrganizationName)
	r.NoError(err)
	r.True(resp.IsOK())
}

func createUser(r *require.Assertions, c *rest.Client) {
	resp, err := c.Users().OnServer().Create(testUserName, DefaultPassword)
	r.NoError(err)
	r.True(resp.IsOK())
}

func createRole(r *require.Assertions, c *rest.Client) {
	resp, err := c.Roles().Create(
		testRoleName,
		&rest.RoleCreateOptions{Action: []rest.RoleAction{rest.RoleActionCreateDatabase, rest.RoleActionDeleteDatabase}},
	)
	r.NoError(err)
	r.True(resp.IsOK())
}

func createDatabase(r *require.Assertions, c *rest.Client) {
	resp, err := c.Databases().
		OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).
		Create(testDatabaseName, testDatabaseLabel, nil)
	r.NoError(err)
	r.True(resp.IsOK())
}

func createBranch(r *require.Assertions, c *rest.Client) {
	resp, err := c.Branches().
		OnRepo(rest.RepoPath{
			Organization: testOrganizationName,
			Database:     testDatabaseName,
			Repo:         "local",
		}).Create(testBranchName, nil)
	r.NoError(err)
	r.True(resp.IsOK())
}

func setupSchema(r *require.Assertions, c *rest.Client) {
	rawDoc := `
	{
		"@context": {
			"@base": "terminusdb:///data/",
			"@schema": "terminusdb:///schema#",
			"@type": "Context"
		},
		"child1": {
			"@inherits": [
				"doc1"
			],
			"@key": {
				"@type": "Random"
			},
			"@type": "Class",
			"enumfield1": {
				"@class": {
					"@id": "enum1",
					"@type": "Enum",
					"@values": [
						"id1",
						"id2",
						"id3"
					]
				},
				"@type": "Optional"
			},
			"link1": {
				"@class": {
					"@class": "sub1",
					"@subdocument": []
				},
				"@type": "Optional"
			},
			"num1": {
				"@class": "xsd:integer",
				"@type": "Optional"
			},
			"str1": {
				"@class": "xsd:string",
				"@type": "Optional"
			}
		},
		"child2": {
			"@inherits": [
				"child1",
				"doc1"
			],
			"@key": {
				"@type": "Random"
			},
			"@type": "Class",
			"enumfield1": {
				"@class": {
					"@id": "enum1",
					"@type": "Enum",
					"@values": [
						"id1",
						"id2",
						"id3"
					]
				},
				"@type": "Optional"
			},
			"link1": {
				"@class": {
					"@class": "sub1",
					"@subdocument": []
				},
				"@type": "Optional"
			},
			"num1": {
				"@class": "xsd:integer",
				"@type": "Optional"
			},
			"str1": {
				"@class": "xsd:string",
				"@type": "Optional"
			}
		},
		"doc1": {
			"@key": {
				"@type": "Random"
			},
			"@type": "Class",
			"link1": {
				"@class": {
					"@class": "sub1",
					"@subdocument": []
				},
				"@type": "Optional"
			},
			"num1": {
				"@class": "xsd:integer",
				"@type": "Optional"
			},
			"str1": {
				"@class": "xsd:string",
				"@type": "Optional"
			}
		},
		"enum1": {
			"@type": "Enum",
			"@values": [
				"id1",
				"id2",
				"id3"
			]
		},
		"sub1": {
			"@key": {
				"@type": "Random"
			},
			"@subdocument": [],
			"@type": "Class",
			"substring1": {
				"@class": "xsd:string",
				"@type": "Optional"
			}
		}
	}`
	doc := make(rest.GenericDocument)
	err := json.Unmarshal([]byte(rawDoc), &doc)
	r.NoError(err)

	resp, err := c.GenericDocuments().OnBranch(rest.BranchPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
		Branch:       testBranchName,
	}).Create(doc, &rest.DocumentCreateOptions{
		GraphType:   rest.GraphTypeSchema,
		Message:     "create schema",
		Author:      testCommitAuthor,
		RawJSON:     false,
		FullReplace: true,
	})
	r.NoError(err)
	r.True(resp.IsOK())
}

func insertDocument(r *require.Assertions, c *rest.Client, doc rest.GenericDocument) {
	resp, err := c.GenericDocuments().OnBranch(rest.BranchPath{
		Organization: testOrganizationName,
		Database:     testDatabaseName,
		Repo:         "local",
		Branch:       testBranchName,
	}).Create(doc, &rest.DocumentCreateOptions{
		GraphType:   rest.GraphTypeInstance,
		Message:     "insert document",
		Author:      testCommitAuthor,
		RawJSON:     false,
		FullReplace: true,
	})
	r.NoError(err)
	r.True(resp.IsOK())
}

func cleanUpTest(c *rest.Client) {
	_, _ = c.Databases().OnOrganization(rest.OrganizationPath{Organization: testOrganizationName}).Delete(testDatabaseName, nil)
	_, _ = c.Users().OnServer().Delete(testUserName)
	_, _ = c.Roles().Delete(testRoleName)
	_, _ = c.Organizations().Delete(testOrganizationName)
}

func getDatabaseURL() string {
	url := os.Getenv(DatabaseURLVariable)
	if url == "" {
		url = DefaultDatabaseURL
	}
	return url
}

func extractField[T any](obj any, fieldName string) T {
	val := reflect.Indirect(reflect.ValueOf(obj))
	if !val.IsValid() {
		panic("object is nil")
	}

	t := new(T)
	tTyp := reflect.TypeOf(*t)
	fld := val.FieldByName(fieldName)
	return fld.Convert(tTyp).Interface().(T)
}
