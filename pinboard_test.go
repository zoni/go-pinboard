package pinboard

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var p1 = Pinboard{
	User:     "drags",
	Password: "foobar",
}

var p2 = Pinboard{
	User:  "drags",
	Token: "AC1638B3E618FD194CA0",
}

var p3 = Pinboard{
	User: "drags",
}

func TestAuthQueryPassword(t *testing.T) {
	want := "http://drags:foobar@pinboard.example"
	got, _ := url.Parse("http://pinboard.example")
	err := p1.authQuery(got)
	if err != nil {
		t.Errorf("Error from authQuery: %v", err)
	}
	if got.String() != want {
		t.Errorf("Want %s, Got %s", want, got.String())
	}
}

func TestAuthQueryToken(t *testing.T) {
	want := "http://pinboard.example?auth_token=drags%3AAC1638B3E618FD194CA0"
	got, _ := url.Parse("http://pinboard.example")
	err := p2.authQuery(got)
	if err != nil {
		t.Errorf("Error from authQuery: %v", err)
	}
	if got.String() != want {
		t.Errorf("Want %s, Got %s", want, got.String())
	}
}

func TestAuthQueryInvalid(t *testing.T) {
	want := errors.New("Pinboard requires either a Password or Token for authentication")
	u, _ := url.Parse("http://pinboard.example")
	got := p3.authQuery(u)
	if got == nil {
		t.Error("Did not got an error from authQuery when expected")
	}
	if got.Error() != want.Error() {
		t.Errorf("Wanted %v, Got %v", want, got)
	}
}

var testResponse = `<books>
  <book>
    <author>Ray Robinson</author>
	<title>Sunny's Listening</title>
  </book>
  <book>
	<author>STS RJD2</author>
	<title>The Music of Now</title>
  </book>
</books>
  `

type testBooks struct {
	XMLName xml.Name
	Books   []testBook `xml:"book"`
}

type testBook struct {
	XMLName xml.Name
	Title   string `xml:"title"`
	Author  string `xml:"author"`
}

func TestParseResponse(t *testing.T) {
	want := testBooks{
		XMLName: xml.Name{Space: "", Local: "books"},
		Books: []testBook{
			{
				XMLName: xml.Name{Space: "", Local: "book"},
				Author:  "Ray Robinson",
				Title:   "Sunny's Listening",
			},
			{
				XMLName: xml.Name{Space: "", Local: "book"},
				Author:  "STS RJD2",
				Title:   "The Music of Now",
			},
		}}
	// httptest mock server
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testResponse)
	}))
	defer s.Close()

	// request from mock server
	c := s.Client()
	resp, err := c.Get(s.URL)
	if err != nil {
		t.Error("Got error from httptest server", err)
	}

	// testing pinboard.parseResponse
	tmp, err := parseResponse(resp, &testBooks{})
	if err != nil {
		t.Errorf("Error parsing TestBook response: %v", err)
	}
	got := tmp.(*testBooks)

	if !reflect.DeepEqual(want, *got) {
		t.Log("Want", want)
		t.Log("Gott", *got)
		t.Errorf("testBooks did not parse as expected")
	}
}
