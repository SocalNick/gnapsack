package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gusto/gnapsack/db"

	"github.com/gusto/gnapsack/build"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var _ = Describe("Build", func() {

	var (
		db   *gorm.DB
		bids build.InputDistributionsSubset
	)

	BeforeEach(func() {
		var err error
		db, err = gorm.Open("sqlite3", ":memory:")
		Expect(err).ToNot(HaveOccurred())

		db.AutoMigrate(
			&build.Build{},
			&build.Distribution{},
			&build.DistributionTestFile{},
			&build.Subset{},
		)

		bids = build.InputDistributionsSubset{
			CommitHash: "foo",
			Branch:     "master",
			NodeTotal:  2,
		}
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("Finding or creating a build", func() {

		Context("When there is NOT an existing build", func() {

			It("should create a build", func() {
				build, err := FindOrCreateBuildForDistributionSubset(db, bids)
				Expect(err).ToNot(HaveOccurred())

				Expect(build.CommitHash).To(Equal("foo"))
			})
		})

	})

})
