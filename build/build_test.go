package build

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestDistributionNodeIndexOutOfBounds(t *testing.T) {
	g := NewGomegaWithT(t)

	s := InputDistributionsSubset{
		NodeTotal: 2,
		NodeIndex: 2,
		TestFiles: []InputDistributionsSubsetTestFile{},
	}
	_, err := s.Distribution()
	g.Expect(err).To(HaveOccurred())
}

func TestDistributionSimpleSplitByNodes(t *testing.T) {
	g := NewGomegaWithT(t)

	tf1 := InputDistributionsSubsetTestFile{
		Path: "foo",
	}
	tf2 := InputDistributionsSubsetTestFile{
		Path: "bar",
	}
	tf3 := InputDistributionsSubsetTestFile{
		Path: "baz",
	}
	tf := []InputDistributionsSubsetTestFile{tf1, tf2, tf3}

	s := InputDistributionsSubset{
		NodeTotal: 2,
		NodeIndex: 0,
		TestFiles: tf,
	}
	subset, err := s.Distribution()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(subset.NodeIndex).To(Equal(0))
	g.Expect(subset.TestFiles).To(ContainElements(DistributionTestFile{Path: "foo"}, DistributionTestFile{Path: "bar"}))
	g.Expect(subset.TestFiles).NotTo(ContainElements(DistributionTestFile{Path: "baz"}))

	s = InputDistributionsSubset{
		NodeTotal: 2,
		NodeIndex: 1,
		TestFiles: tf,
	}
	subset, err = s.Distribution()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(subset.NodeIndex).To(Equal(1))
	g.Expect(subset.TestFiles).NotTo(ContainElements(DistributionTestFile{Path: "foo"}, DistributionTestFile{Path: "bar"}))
	g.Expect(subset.TestFiles).To(ContainElements(DistributionTestFile{Path: "baz"}))

	tf4 := InputDistributionsSubsetTestFile{
		Path: "boo",
	}
	tf = append(tf, tf4)
	s = InputDistributionsSubset{
		NodeTotal: 2,
		NodeIndex: 1,
		TestFiles: tf,
	}
	subset, err = s.Distribution()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(subset.NodeIndex).To(Equal(1))
	g.Expect(subset.TestFiles).NotTo(ContainElements(DistributionTestFile{Path: "foo"}, DistributionTestFile{Path: "bar"}))
	g.Expect(subset.TestFiles).To(ContainElements(DistributionTestFile{Path: "baz"}, DistributionTestFile{Path: "boo"}))
}
