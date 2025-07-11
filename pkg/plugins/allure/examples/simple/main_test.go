//go:build example

package main

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/metafates/testo"
	"github.com/metafates/testo/pkg/plugins/allure"
)

func Test(t *testing.T) {
	testo.RunSuite[*Suite, T](t)
}

type T = struct {
	*testo.T
	*allure.Allure
}

type Suite struct {
	client *http.Client
}

// BeforeAll hook is executed before running any tests.
func (s *Suite) BeforeAll(t T) {
	s.client = &http.Client{Timeout: 10 * time.Second}
}

// CasesURL provides URLs.
func (s *Suite) CasesURL() []string {
	return []string{
		"https://example.com",
		"https://example.net",
		"https://example.org",
	}
}

// TestExample is executed for each URL from [Suite.CasesURL] output.
func (s *Suite) TestExample(t T, params struct{ URL string }) {
	t.Parallel()

	t.Title("Send request to " + params.URL)
	t.Tags("http", "html", "example")
	t.Description("This test sends a simple request to the given URL and inspects its response.")
	t.Severity(allure.SeverityCritical)
	t.Links(allure.NewLink("https://github.com").Issue().Named("github repo"))

	var res *http.Response

	allure.Step(t, "send request", func(t T) {
		t.Parameters(
			allure.NewParameter("method", http.MethodGet),
			allure.NewParameter("token", "super secret here").Masked(),
		)

		var err error
		res, err = s.client.Get(params.URL)

		is := allure.Require(t)
		is.NoError(err, "successful round trip")
		is.Equal(http.StatusOK, res.StatusCode, "status is OK")
	})

	allure.Step(t, "inspect response", func(t T) {
		allure.Step(t, "check headers", func(t T) {
			contentType := res.Header.Get("Content-Type")

			allure.Assert(t).Equal("text/html", contentType, "Content-Type is HTML")
		})

		allure.Step(t, "read body", func(t T) {
			body, err := io.ReadAll(res.Body)

			allure.Require(t).NoError(err)

			t.Attach("body", allure.NewAttachmentString(string(body)).As(allure.TextHTML))

			allure.Assert(t).Greater(len(body), 0, "not empty")
		})
	})
}
