package build

import (
	"errors"

	. "github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type InputDistributionsSubset struct {
	NodeTotal  int                                `json:"node_total" binding:"gte=2"`
	NodeIndex  int                                `json:"node_index" binding:"lt=node_total`
	CommitHash string                             `json:"commit_hash" binding:"required"`
	Branch     string                             `json:"branch" binding:"required"`
	TestFiles  []InputDistributionsSubsetTestFile `json:"test_files" binding:"required"`
}

type InputDistributionsSubsetTestFile struct {
	Path string `json:"path" binding:"required"`
}

type Build struct {
	gorm.Model
	Uuid          UUID `sql:"index"`
	CommitHash    string
	Branch        string
	NodeTotal     int
	Distributions []Distribution
	Subsets       []Subset
}

type Distribution struct {
	gorm.Model
	DistributionId UUID                   `json:"build_distribution_id"`
	NodeIndex      int                    `json:"node_index"`
	TestFiles      []DistributionTestFile `json:"test_files"`
}

type DistributionTestFile struct {
	gorm.Model
	SubsetID      int
	Path          string  `json:"path" binding:"required"`
	TimeExecution float64 `json:"time_execution"`
}

type Subset struct {
	gorm.Model
	NodeTotal  int                    `json:"node_total" binding:"gte=2"`
	NodeIndex  int                    `json:"node_index" binding:"lt=node_total`
	CommitHash string                 `json:"commit_hash" binding:"required"`
	Branch     string                 `json:"branch" binding:"required"`
	TestFiles  []DistributionTestFile `json:"test_files" binding:"required"`
}

func (ds InputDistributionsSubset) Distribution() (distribution Distribution, err error) {
	if ds.NodeIndex >= ds.NodeTotal {
		return Distribution{}, errors.New("NodeIndex must be less than NodeTotal")
	}

	numTestFiles := len(ds.TestFiles)

	testFilesPerNode := numTestFiles / ds.NodeTotal
	remainder := numTestFiles % ds.NodeTotal

	firstIndex := 0
	lastIndex := testFilesPerNode + remainder
	if ds.NodeIndex > 0 {
		firstIndex = ds.NodeIndex*testFilesPerNode + remainder
		lastIndex = firstIndex + testFilesPerNode
	}

	chunk := ds.TestFiles[firstIndex:lastIndex]
	distributionTestFiles := []DistributionTestFile{}

	for _, v := range chunk {
		distributionTestFile := DistributionTestFile{Path: v.Path}
		distributionTestFiles = append(distributionTestFiles, distributionTestFile)
	}

	distribution = Distribution{
		NodeIndex: ds.NodeIndex,
		TestFiles: distributionTestFiles,
	}
	return distribution, nil
}
