package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/tiktoken-go/tokenizer"
	"os"
)

func main() {

	_ = godotenv.Load()
	// Initialize openai client
	openapi_key := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(openapi_key), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	// Initialize tokenizer client

	codec, err := tokenizer.ForModel(tokenizer.GPT5)
	if err != nil {
		panic(err)
	}

	prompt := "Give me 2 reasons to love you "
	// this should print a list of token ids
	ids, tokens, _ := codec.Encode(prompt)
	fmt.Println(ids)

	fmt.Println(tokens)
	// this should print the original string back
	text, _ := codec.Decode(ids)
	fmt.Println(text)

	// Vector representation of tokens from input string
	params := openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: tokens,
		},
		Model: openai.EmbeddingModelTextEmbedding3Large,
	}
	embeddings, err := client.Embeddings.New(context.TODO(), params)
	if err != nil {
		panic(err.Error())
	}

	// printing first 10 embeddings
	fmt.Println(len(embeddings.Data[0].Embedding))
	fmt.Printf("%v\n", embeddings.Data[0].Embedding[:10])

	completion, err := client.Completions.New(context.TODO(), openai.CompletionNewParams{
		Model: openai.CompletionNewParamsModelGPT3_5TurboInstruct,
		Prompt: openai.CompletionNewParamsPromptUnion{
			OfString: openai.String(prompt),
		},
		Temperature: openai.Float(0.0),
	})
	if err != nil {
		panic(err)
	}

	choice := completion.Choices[0]

	// ✅ This is the final output that OpenAI selected
	fmt.Printf("\nBest completion selected:\n%s\n", choice.Text)

}
