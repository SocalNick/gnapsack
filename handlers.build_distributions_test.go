package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

func TestSplitsSuiteIntoSubset(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	initializeRoutes(router)

	postBody := `
{
	"commit_hash": "foo",
	"branch": "master",
	"node_total": 2,
	"node_index": 0,
	"test_files":[
	  {"path":"foo"},
		{"path":"bar"},
		{"path":"baz"}
	]
}
`

	req, _ := http.NewRequest("POST", "/v1/build_distributions/subset", strings.NewReader(postBody))

	router.ServeHTTP(w, req)

	response := w.Result()

	g.Expect(response.StatusCode).To(Equal(http.StatusOK))

	actual, err := ioutil.ReadAll(response.Body)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(string(actual)).Should(ContainSubstring("foo"))
	g.Expect(string(actual)).Should(ContainSubstring("bar"))
	g.Expect(string(actual)).ShouldNot(ContainSubstring("baz"))
}

func TestCommitHashIsRequired(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	initializeRoutes(router)

	req, _ := http.NewRequest("POST", "/v1/build_distributions/subset", strings.NewReader("{}"))

	router.ServeHTTP(w, req)

	response := w.Result()

	g.Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))

	actual, err := ioutil.ReadAll(response.Body)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(string(actual)).Should(ContainSubstring("Key: 'InputDistributionsSubset.CommitHash' Error:Field validation for 'CommitHash' failed on the 'required' tag"))
}

func TestBranchIsRequired(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	initializeRoutes(router)

	req, _ := http.NewRequest("POST", "/v1/build_distributions/subset", strings.NewReader(`{}`))

	router.ServeHTTP(w, req)

	response := w.Result()

	g.Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))

	actual, err := ioutil.ReadAll(response.Body)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(string(actual)).Should(ContainSubstring("Key: 'InputDistributionsSubset.Branch' Error:Field validation for 'Branch' failed on the 'required'"))
}

func TestNodeTotalMustBeGreaterThanOrEqualTo2(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	initializeRoutes(router)

	req, _ := http.NewRequest("POST", "/v1/build_distributions/subset", strings.NewReader(`{}`))

	router.ServeHTTP(w, req)

	response := w.Result()

	g.Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))

	actual, err := ioutil.ReadAll(response.Body)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(string(actual)).Should(ContainSubstring("Key: 'InputDistributionsSubset.NodeTotal' Error:Field validation for 'NodeTotal' failed on the 'gte' tag"))
}
