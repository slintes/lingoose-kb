package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/henomis/lingoose/document"
	"github.com/henomis/lingoose/loader"
	"github.com/henomis/lingoose/types"
)

const wikiURL = "https://en.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=extracts&explaintext"

type wikiResult struct {
	BatchComplete string `json:"batchcomplete"`
	Query         struct {
		Normalized []struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"normalized"`
		Pages map[string]struct {
			PageID  int    `json:"pageid"`
			Ns      int    `json:"ns"`
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

type WikiLoader struct {
	splitter loader.TextSplitter
	title    string
}

func NewWikiLoader(title string) *WikiLoader {
	return &WikiLoader{
		title: title,
	}
}

func (w *WikiLoader) WithTextSplitter(textSplitter loader.TextSplitter) *WikiLoader {
	w.splitter = textSplitter
	return w
}

func (w *WikiLoader) Load(ctx context.Context) ([]document.Document, error) {

	documents := make([]document.Document, 1)
	var err error
	doc, err := w.load(ctx)
	if err != nil {
		return nil, err
	}
	documents[0] = *doc

	if w.splitter != nil {
		documents = w.splitter.SplitDocuments(documents)
	}

	return documents, nil
}

func (w *WikiLoader) load(ctx context.Context) (*document.Document, error) {

	url := fmt.Sprintf(wikiURL, w.title)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result wikiResult
	err = json.Unmarshal(jsonContent, &result)
	if err != nil {
		return nil, err
	}

	content := ""

	for _, page := range result.Query.Pages {
		content += page.Extract
	}

	return &document.Document{
		Content: content,
		Metadata: types.Meta{
			"source": url,
		},
	}, nil
}
