package db

import (
	"github.com/gusto/gnapsack/build"
	"github.com/jinzhu/gorm"
)

func FindOrCreateBuildForDistributionSubset(db *gorm.DB, bids build.InputDistributionsSubset) (build.Build, error) {
	var b build.Build
	if err := db.Where("commit_hash = ? AND branch = ? AND node_total = ?", bids.CommitHash, bids.Branch, bids.NodeTotal).Find(&b).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return b, err
		}

		b = build.Build{
			CommitHash: bids.CommitHash,
			Branch:     bids.Branch,
			NodeTotal:  bids.NodeTotal,
		}

		db.Create(&b)
	}
	return b, nil
}

func CreateBuildSubset(db *gorm.DB, bs build.Subset) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&bs).Error; err != nil {
			return err
		}
		return nil
	})
}

func ListBuildSubsets(db *gorm.DB) ([]build.Subset, error) {
	var buildSubsets []build.Subset
	if err := db.Preload("TestFiles").Find(&buildSubsets).Error; err != nil {
		return nil, err
	}
	return buildSubsets, nil
}
