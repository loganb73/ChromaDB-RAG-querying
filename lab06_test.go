package project05loganb73

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	chroma "github.com/amikos-tech/chroma-go"
	chromaOpenai "github.com/amikos-tech/chroma-go/openai"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/gocarina/gocsv"
	"github.com/sashabaranov/go-openai"
)

func TestVectorQuery(t *testing.T) {

	//openai setup
	openaiKey := os.Getenv("openai_key_temp")

	aiClient := openai.NewClient(openaiKey)
	resp, err := aiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Logan!",
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)

	//chroma setup
	classesFile, err := os.OpenFile("Fall 2024 Class Schedule 08082024.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("OpenFile error: %v\n", err)
		return
	}
	defer classesFile.Close()

	classes := []*Class{}
	err = gocsv.UnmarshalFile(classesFile, &classes)
	if err != nil {
		fmt.Printf("UnmarshalFile error: %v\n", err)
	}

	jsonStrings := []string{}
	for _, class := range classes {
		classJson, err := json.Marshal(class)
		if err != nil {
			fmt.Printf("error marshalling struct: %s\n", err.Error())
		}
		jsonStrings = append(jsonStrings, string(classJson))
	}

	client, err := chroma.NewClient("http://localhost:8000")
	if err != nil {
		log.Fatalf("Error creating client: %s \n", err)
		return
	}

	openaiEf, err := chromaOpenai.NewOpenAIEmbeddingFunction(os.Getenv("openai_key_temp"))
	if err != nil {
		log.Fatalf("Error creating OpenAI embedding function: %s \n", err)
	}

	collection, err := client.CreateCollection(context.TODO(), "my-collection", map[string]interface{}{"key1": "value1"}, true, openaiEf, types.L2)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}
	if err != nil {
		log.Fatalf("Error creating collection: %s \n", err)
	}

	rs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(collection.EmbeddingFunction), // we pass the embedding function from the collection
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		log.Fatalf("Error creating record set: %s \n", err)
	}

	i := 0
	for _, classString := range jsonStrings {
		rs.WithRecord(types.WithDocument(classString))

		if (i % 500) == 0 {
			//buildandvalidate
			_, err = rs.BuildAndValidate(context.TODO())
			if err != nil {
				log.Fatalf("Error validating record set: %s \n", err)
			}
			fmt.Printf("inserted %d docs\n", i)
		}
		i++
	}

	//insert stragglers (not caught by last batch insert)
	_, err = rs.BuildAndValidate(context.TODO())
	if err != nil {
		log.Fatalf("Error validating record set: %s \n", err)
	}

	// Add the records to the collection
	_, err = collection.AddRecords(context.Background(), rs)
	if err != nil {
		log.Fatalf("Error adding documents: %s \n", err)
	}

	// Count the number of documents in the collection
	countDocs, qrerr := collection.Count(context.TODO())
	if qrerr != nil {
		log.Fatalf("Error counting documents: %s \n", qrerr)
	}

	// Query the collection
	fmt.Printf("countDocs: %v\n", countDocs) //this should result in 2
	qr, qrerr := collection.Query(context.TODO(), []string{"Logan"}, 1, nil, nil, nil)
	if qrerr != nil {
		log.Fatalf("Error querying documents: %s \n", qrerr)
	}
	fmt.Printf("qr: %v\n", qr.Documents[0]) //this should result in the document about
}
