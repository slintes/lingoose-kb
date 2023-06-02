package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/henomis/lingoose/chat"
	"github.com/henomis/lingoose/prompt"
	"github.com/henomis/lingoose/types"
	"os"
	"strings"

	openaiembedder "github.com/henomis/lingoose/embedder/openai"
	"github.com/henomis/lingoose/index"
	"github.com/henomis/lingoose/llm/openai"
	"github.com/henomis/lingoose/textsplitter"
)

func main() {

	if len(os.Args) != 2 {
		panic("no title provided!")
	}

	title := os.Args[1]

	openaiEmbedder := openaiembedder.New(openaiembedder.AdaEmbeddingV2)

	docsVectorIndex := index.NewSimpleVectorIndex("db", ".", openaiEmbedder)
	indexIsEmpty, _ := docsVectorIndex.IsEmpty()

	if indexIsEmpty {
		err := ingestData(docsVectorIndex, title)
		if err != nil {
			panic(err)
		}
	}

	llmOpenAI := openai.NewChat()

	fmt.Println("Enter a query to search the knowledge base. Type 'quit' to exit.")
	query := ""
	for query != "quit" {

		fmt.Printf("> ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		query = strings.TrimSpace(input)

		if query == "quit" {
			break
		}

		similarities, err := docsVectorIndex.SimilaritySearch(context.Background(), query, index.WithTopK(3))
		if err != nil {
			panic(err)
		}

		content := ""

		for _, similarity := range similarities {
			//fmt.Printf("Similarity: %f\n", similarity.Score)
			//fmt.Printf("Document: %s\n", similarity.Document.Content)
			//fmt.Println("Metadata: ", similarity.Document.Metadata)
			//fmt.Println("----------")
			content += similarity.Document.Content + "\n"
		}

		systemPrompt := prompt.New("You are an helpful assistant. Answer to the questions using only " +
			"the provided context. Don't add any information that is not in the context. " +
			"If you don't know the answer, just say 'I don't know'.",
		)
		userPrompt := prompt.NewPromptTemplate(
			"Based on the following context answer to the question.\n\nContext:\n{{.context}}\n\nQuestion: {{.query}}").WithInputs(
			types.M{
				"query":   query,
				"context": content,
			},
		)

		chat := chat.New(
			chat.PromptMessage{
				Type:   chat.MessageTypeSystem,
				Prompt: systemPrompt,
			},
			chat.PromptMessage{
				Type:   chat.MessageTypeUser,
				Prompt: userPrompt,
			},
		)

		response, err := llmOpenAI.Chat(context.Background(), chat)
		if err != nil {
			panic(err)
		}

		fmt.Println(response)

	}

}

func ingestData(docsVectorIndex *index.SimpleVectorIndex, title string) error {

	fmt.Printf("Learning Wiki page... ")

	loader := NewWikiLoader(title)

	documents, err := loader.Load(context.Background())
	if err != nil {
		return err
	}

	textSplitter := textsplitter.NewRecursiveCharacterTextSplitter(2000, 200)

	documentChunks := textSplitter.SplitDocuments(documents)

	err = docsVectorIndex.LoadFromDocuments(context.Background(), documentChunks)
	if err != nil {
		return err
	}

	fmt.Printf("Done\n")

	return nil
}
