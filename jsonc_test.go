package jsonc_test

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unvurn/httpc/form"

	"github.com/unvurn/jsonc"
)

const httpbinEndpoint = "https://httpbin.org"

type httpbinGetResponse[T any] struct {
	Args    T                 `json:"args"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
	Url     string            `json:"url"`
}

type httpbinPostResponse struct {
	Args    map[string]string `json:"args"`
	Data    string            `json:"data"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
	Url     string            `json:"url"`
}

type httpbinPostFormResponse[T any] struct {
	httpbinPostResponse
	Form  T                 `json:"form"`
	Files map[string]string `json:"files"`
}

type httpbinPostJsonResponse[T any] struct {
	httpbinPostResponse
	Json T `json:"json"`
}

func TestHttpbin_Get(t *testing.T) {
	u, _ := url.JoinPath(httpbinEndpoint, "get")
	resp, err := jsonc.NewRequest[*httpbinGetResponse[any]]().Get(context.Background(), u)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

type params struct {
	Name         string  `schema:"name" json:"name"`
	Age          int     `schema:"age" json:"age"`
	Scores       []int   `schema:"scores" json:"scores"`
	Description  string  `schema:"description" json:"description"`
	Description2 string  `schema:"description_2,omitempty" json:"description_2,omitempty"`
	ExtraNote    *string `schema:"extra_note" json:"extra_note"`
	ExtraNote2   *string `schema:"extra_note_2,omitempty" json:"extra_note_2,omitempty"`
}

func (p *params) UnmarshalJSON(b []byte) error {
	var temp struct {
		Name         string   `json:"name"`
		Age          string   `json:"age"`
		Scores       []string `json:"scores"`
		Description  string   `json:"description"`
		Description2 string   `json:"description_2,omitempty"`
		ExtraNote    *string  `json:"extra_note"`
		ExtraNote2   *string  `json:"extra_note_2,omitempty"`
	}

	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	p.Name = temp.Name
	p.Age, err = strconv.Atoi(temp.Age)
	if err != nil {
		return err
	}
	for _, score := range temp.Scores {
		s, err := strconv.Atoi(score)
		if err != nil {
			return err
		}
		p.Scores = append(p.Scores, s)
	}
	p.Description = temp.Description
	p.Description2 = temp.Description2
	if temp.ExtraNote == nil || *temp.ExtraNote == "null" {
		p.ExtraNote = nil
	} else {
		p.ExtraNote = temp.ExtraNote
	}
	if temp.ExtraNote2 == nil || *temp.ExtraNote2 == "null" {
		p.ExtraNote2 = nil
	} else {
		p.ExtraNote2 = temp.ExtraNote2
	}
	return nil
}

func TestHttpbin_GetQuery(t *testing.T) {
	u, _ := url.JoinPath(httpbinEndpoint, "get")
	p := params{
		Name:   "John Doe",
		Age:    25,
		Scores: []int{100, 90, 80},
	}
	resp, err := jsonc.NewRequest[httpbinGetResponse[params]]().Get(context.Background(), u, p)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "John Doe", resp.Args.Name)
	assert.Equal(t, 25, resp.Args.Age)
	assert.Equal(t, []int{100, 90, 80}, resp.Args.Scores)
	assert.Empty(t, resp.Args.Description)
	assert.Empty(t, resp.Args.Description2)
	assert.Nil(t, resp.Args.ExtraNote)
	assert.Nil(t, resp.Args.ExtraNote2)
	assert.Len(t, resp.Args.Scores, 3)
	assert.NotEmpty(t, resp.Headers)
	assert.NotEmpty(t, resp.Origin)
	u2, _ := url.Parse(resp.Url)
	u2.RawQuery = ""
	assert.Equal(t, u, u2.String())
	assert.Contains(t, resp.Headers["User-Agent"], "Go-http-client/2.0")
}

func TestHttpbin_GetQuery2(t *testing.T) {
	u, _ := url.JoinPath(httpbinEndpoint, "get")
	empty := ""
	p := params{
		Name:       "John Doe",
		Age:        25,
		Scores:     []int{100, 90, 80},
		ExtraNote:  &empty,
		ExtraNote2: &empty,
	}
	resp, err := jsonc.NewRequest[httpbinGetResponse[params]]().Get(context.Background(), u, p)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "John Doe", resp.Args.Name)
	assert.Equal(t, 25, resp.Args.Age)
	assert.Equal(t, []int{100, 90, 80}, resp.Args.Scores)
	assert.Empty(t, resp.Args.Description)
	assert.Empty(t, resp.Args.Description2)
	assert.Empty(t, resp.Args.ExtraNote)
	assert.Empty(t, resp.Args.ExtraNote2)
	assert.Len(t, resp.Args.Scores, 3)
	assert.NotEmpty(t, resp.Headers)
	assert.NotEmpty(t, resp.Origin)
	u2, _ := url.Parse(resp.Url)
	u2.RawQuery = ""
	assert.Equal(t, u, u2.String())
	assert.Contains(t, resp.Headers["User-Agent"], "Go-http-client/2.0")
}

func TestHttpbin_PostForm(t *testing.T) {
	u, _ := url.JoinPath(httpbinEndpoint, "post")
	p := params{
		Name:   "Jane Doe",
		Age:    25,
		Scores: []int{100, 90, 80},
	}
	resp, err := jsonc.NewRequest[httpbinPostFormResponse[params]]().Post(context.Background(), u, p)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Jane Doe", resp.Form.Name)
	assert.Equal(t, 25, resp.Form.Age)
	assert.Equal(t, []int{100, 90, 80}, resp.Form.Scores)
	assert.Len(t, resp.Form.Scores, 3)
	assert.NotEmpty(t, resp.Headers)
	assert.NotEmpty(t, resp.Origin)
	assert.Equal(t, u, resp.Url)
	assert.Contains(t, resp.Headers["User-Agent"], "Go-http-client/2.0")
}

func TestHttpbin_PostFileUpload(t *testing.T) {
	u, _ := url.JoinPath(httpbinEndpoint, "post")
	p := params{
		Name:   "Jane Doe",
		Age:    25,
		Scores: []int{100, 90, 80},
	}
	resp, err := jsonc.NewRequest[httpbinPostFormResponse[params]]().Post(context.Background(), u, p,
		form.Bytes("data1", "data1.txt", []byte("This is data1 content.")),
		form.File("data2", "samples/dummy.pdf"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Jane Doe", resp.Form.Name)
	assert.Equal(t, 25, resp.Form.Age)
	assert.Len(t, resp.Form.Scores, 3)
	assert.Equal(t, []int{100, 90, 80}, resp.Form.Scores)
	assert.Len(t, resp.Files, 2)
	assert.NotEmpty(t, resp.Headers)
	assert.NotEmpty(t, resp.Origin)
	assert.Equal(t, u, resp.Url)
	assert.Contains(t, resp.Headers["User-Agent"], "Go-http-client/2.0")
}

func TestHttpbin_PostJson(t *testing.T) {
	u, _ := url.JoinPath(httpbinEndpoint, "post")
	type params struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Scores []int  `json:"scores"`
	}
	p := params{
		Name:   "Jane Doe",
		Age:    25,
		Scores: []int{100, 90, 80},
	}
	resp, err := jsonc.NewRequest[httpbinPostJsonResponse[params]]().PostJson(context.Background(), u, p)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Jane Doe", resp.Json.Name)
	assert.Equal(t, 25, resp.Json.Age)
	assert.Len(t, resp.Json.Scores, 3)
	assert.Equal(t, []int{100, 90, 80}, resp.Json.Scores)
	assert.NotEmpty(t, resp.Headers)
	assert.NotEmpty(t, resp.Origin)
	assert.Equal(t, u, resp.Url)
	assert.Contains(t, resp.Headers["User-Agent"], "Go-http-client/2.0")
}
