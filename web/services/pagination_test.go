package services

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"gorm.io/gorm"
)

type TestEntity struct {
	ID uint
}

type PaginationTestSuite struct {
	suite.Suite
	db *gorm.DB
	tx *gorm.DB
}

func TestPaginationTestSuite(t *testing.T) {
	suite.Run(t, new(PaginationTestSuite))
}

func (suite *PaginationTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(TestEntity{})
}

func (suite *PaginationTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(TestEntity{})
}

func (suite *PaginationTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
}

func (suite *PaginationTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func (suite *PaginationTestSuite) TestPagination() {
	var entities []TestEntity
	for i := 1; i < 11; i++ {
		entity := TestEntity{ID: uint(i)}
		entities = append(entities, entity)
	}

	err := suite.tx.Create(&entities).Error
	suite.NoError(err)

	var nonPaginatedEntities []TestEntity
	suite.tx.Scopes(Paginate(nil)).Find(&nonPaginatedEntities)
	suite.Equal(10, len(nonPaginatedEntities))

	var paginatedEntities []TestEntity
	suite.tx.Scopes(Paginate(&Page{
		Number: 2,
		Size:   5,
	})).Find(&paginatedEntities)

	suite.Equal(5, len(paginatedEntities))
	suite.Equal(uint(6), paginatedEntities[0].ID)
}
