package schemamd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SchemaOrgSuite struct {
	suite.Suite
	metaData *MetaData
}

func (suite *SchemaOrgSuite) SetupSuite() {
	suite.metaData = NewMetaData()
	for _, issue := range suite.metaData.issuesFound {
		fmt.Printf("[%s] %s: %s\n", issue.IssueCode(), issue.IssueContext(), issue.Issue())
	}
}

func (suite *SchemaOrgSuite) TearDownSuite() {
}

func (suite *SchemaOrgSuite) TestMetaData() {
	suite.Equal(2301, len(suite.metaData.Graph))
	suite.Equal(7, len(suite.metaData.dataTypes))
	suite.Equal(1053, len(suite.metaData.classes))
	suite.Equal(1248, len(suite.metaData.properties))
}

func (suite *SchemaOrgSuite) TestPerson() {
	person, ok := suite.metaData.classes["http://schema.org/Person"]
	suite.True(ok, "Node must be found")
	suite.NotNil(person, "Node must not be nil")
	suite.Equal(person.ClassName(), "Person", "ClassName should match")
	suite.Equal(person.PropertyName(), "person", "PropertyName should match")

	members := person.ClassMembers(suite.metaData)
	suite.NotNil(members, "Members should be found")
	suite.Equal(59, len(members), "Members should be present")
	for _, node := range members {
		fmt.Printf("%s.%s %+v\n", person.ClassName(), node.PropertyName(), node.rangeRefs)
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SchemaOrgSuite))
}
