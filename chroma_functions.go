package project05loganb73

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	chroma "github.com/amikos-tech/chroma-go"
	chromaOpenai "github.com/amikos-tech/chroma-go/openai"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/gocarina/gocsv"
)

type Class struct {
	Subj                       string `csv:"SUBJ" json:"SUBJ"`
	CourseNum                  string `csv:"CRSE NUM" json:"CRSE NUM"`
	Sec                        string `csv:"SEC" json:"SEC"`
	Crn                        string `csv:"CRN" json:"CRN"`
	SchedType                  string `csv:"Schedule Type Code" json:"Schedule Type Code"`
	CampusCode                 string `csv:"Campus Code" json:"Campus Code"`
	TitleShortDesc             string `csv:"Title Short Desc" json:"Title Short Desc"`
	InstructionModeDesc        string `csv:"Instruction Mode Desc" json:"Instruction Mode Desc"`
	MeetingTypeCodes           string `csv:"Meeting Type Codes" json:"Meeting Type Codes"`
	MeetDays                   string `csv:"Meet Days" json:"Meet Days"`
	BeginTime                  string `csv:"Begin Time" json:"Begin Time"`
	EndTime                    string `csv:"End Time" json:"End Time"`
	MeetStart                  string `csv:"Meet Start" json:"Meet Start"`
	MeetEnd                    string `csv:"Meet End" json:"Meet End"`
	Bldg                       string `csv:"BLDG" json:"BLDG"`
	Rm                         string `csv:"RM" json:"RM"`
	ActualEnrollment           string `csv:"Actual Enrollment" json:"Actual Enrollment"`
	PrimaryInstructorFirstName string `csv:"Primary Instructor First Name" json:"Primary Instructor First Name"`
	PrimaryInstructorLastName  string `csv:"Primary Instructor Last Name" json:"Primary Instructor Last Name"`
	PrimaryInstructorEmail     string `csv:"Primary Instructor Email" json:"Primary Instructor Email"`
	College                    string `csv:"College" json:"College"`
}

func SetupDb() chroma.Client {

	//chroma setup
	classesFile, err := os.OpenFile("Fall 2024 Class Schedule 08082024.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("OpenFile error: %v\n", err)
	}
	defer classesFile.Close()

	classes := []*Class{}
	err = gocsv.UnmarshalFile(classesFile, &classes)
	if err != nil {
		fmt.Printf("UnmarshalFile error: %v\n", err)
	}

	//set up client and embedding function
	client, err := chroma.NewClient("http://localhost:8000")
	if err != nil {
		log.Fatalf("Error creating client: %s \n", err)
	}

	openaiEf, err := chromaOpenai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatalf("Error creating OpenAI embedding function: %s \n", err)
	}

	//create collections
	collection, err := client.CreateCollection(context.TODO(), "full-collection", map[string]any{"professor": "name", "department": "dname", "location": "bname"}, true, openaiEf, types.L2)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	profCollection, err := client.CreateCollection(context.TODO(), "professor-collection", map[string]any{}, true, openaiEf, types.L2)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	departmentCollection, err := client.CreateCollection(context.TODO(), "class-collection", map[string]any{}, true, openaiEf, types.L2)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	locationCollection, err := client.CreateCollection(context.TODO(), "location-collection", map[string]any{}, true, openaiEf, types.L2)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	//create record sets
	rs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(collection.EmbeddingFunction), // we pass the embedding function from the collection
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		log.Fatalf("Error creating record set: %s \n", err)
	}

	professorsRs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(profCollection.EmbeddingFunction), // we pass the embedding function from the collection
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		log.Fatalf("Error creating record set: %s \n", err)
	}

	departmentsRs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(collection.EmbeddingFunction), // we pass the embedding function from the collection
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		log.Fatalf("Error creating record set: %s \n", err)
	}

	locationsRs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(collection.EmbeddingFunction), // we pass the embedding function from the collection
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		log.Fatalf("Error creating record set: %s \n", err)
	}

	//make maps to avoid duplicate inserts
	profMap := make(map[string]bool)
	departmentMap := make(map[string]bool)
	locationMap := make(map[string]bool)

	//loop over classes and add data to record sets
	i := 0
	for _, class := range classes {
		fmt.Printf("on loop %d\n", i)
		classJson, err := json.Marshal(class)
		if err != nil {
			fmt.Printf("error marshalling struct: %s\n", err.Error())
		}

		//update all helper record sets
		professorFullName := class.PrimaryInstructorFirstName + " " + class.PrimaryInstructorLastName
		professorsRs = UpdateRecordset(professorsRs, professorFullName, profMap, "professor")
		departmentsRs = UpdateRecordset(departmentsRs, class.Subj, departmentMap, "department")
		locationsRs = UpdateRecordset(locationsRs, class.Bldg, locationMap, "location")

		//update full record set with all info
		//rs.WithRecord(types.WithDocument(string(classJson)), types.WithMetadata("professor", professorFullName), types.WithMetadata("subject", class.Subj), types.WithMetadata("location", class.Bldg))

		//test approach
		metadata := map[string]interface{}{}
		metadata["professor"] = professorFullName
		metadata["department"] = class.Subj
		metadata["location"] = class.Bldg
		rs.WithRecord(types.WithDocument(string(classJson)), types.WithMetadatas(metadata))

		//build in batches
		if ((i % 500) == 0) && i > 0 {
			_, err = rs.BuildAndValidate(context.TODO())
			if err != nil {
				log.Fatalf("Error validating record set full: %s \n", err)
			}
			fmt.Printf("inserted %d docs\n", i)
		}
		i++
	}

	//insert stragglers (not caught by last batch insert)
	_, err = rs.BuildAndValidate(context.TODO())
	if err != nil {
		log.Fatalf("Error validating record set full post loop: %s \n", err)
	}

	// Add the records to the collection
	_, err = collection.AddRecords(context.Background(), rs)
	if err != nil {
		log.Fatalf("Error adding documents: %s \n", err)
	}

	_, err = profCollection.AddRecords(context.Background(), professorsRs)
	if err != nil {
		log.Fatalf("Error adding documents: %s \n", err)
	}

	_, err = departmentCollection.AddRecords(context.Background(), departmentsRs)
	if err != nil {
		log.Fatalf("Error adding documents: %s \n", err)
	}

	_, err = locationCollection.AddRecords(context.Background(), locationsRs)
	if err != nil {
		log.Fatalf("Error adding documents: %s \n", err)
	}

	// Count the number of documents in the collection
	countDocs, qrerr := collection.Count(context.TODO())
	if qrerr != nil {
		log.Fatalf("Error counting documents: %s \n", qrerr)
	}

	// Query the collection
	fmt.Printf("countDocs: %v\n", countDocs)
	qr, qrerr := profCollection.Query(context.TODO(), []string{"phil peterson"}, 1, nil, nil, nil)
	if qrerr != nil {
		log.Fatalf("Error querying documents: %s \n", qrerr)
	}
	fmt.Printf("qr: %v\n", qr.Documents[0])

	return *client
}

func UpdateRecordset(rs *types.RecordSet, insertString string, existanceMap map[string]bool, metadataTitle string) (updatedRs *types.RecordSet) {

	//don't add empty strings (space is for professor names concatenated with a space in between)
	if insertString == "" || insertString == " " {
		return rs
	}

	//insert if does not already exist in record set
	_, exists := existanceMap[insertString]
	if exists {
		return rs
	} else {
		updatedRs = rs.WithRecord(types.WithDocument(string(insertString)), types.WithMetadata(metadataTitle, insertString))
		existanceMap[insertString] = true
		_, err := rs.BuildAndValidate(context.TODO())
		if err != nil {
			log.Fatalf("Error validating record set profs post loop: %s \n", err)
		}
		return updatedRs
	}
}

func QueryDb(queryString string, collectionName string) (resp []string) {
	ctx := context.Background()
	client, err := chroma.NewClient("http://localhost:8000") //connects to localhost:8000

	if err != nil {
		fmt.Printf("Failed to create client: %v", err)
	}

	openaiEf, err := chromaOpenai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatalf("Error creating OpenAI embedding function: %s \n", err)
	}

	collection, err := client.GetCollection(ctx, "full-collection", openaiEf)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	data, err := collection.Query(context.TODO(), []string{queryString}, 1, nil, nil, nil)
	if err != nil {
		log.Fatalf("Error querying documents: %v\n", err)
	}

	return data.Documents[0]
}

func QueryWithMetadata(queryString string, collectionName string) (resp []string) {
	ctx := context.Background()
	client, err := chroma.NewClient("http://localhost:8000") //connects to localhost:8000

	if err != nil {
		fmt.Printf("Failed to create client: %v", err)
	}

	openaiEf, err := chromaOpenai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatalf("Error creating OpenAI embedding function: %s \n", err)
	}

	collection, err := client.GetCollection(ctx, "full-collection", openaiEf)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	metadata := make(map[string]interface{})
	metadata["instructor"] = "phil peterson"

	data, err := collection.Query(context.TODO(), []string{queryString}, 1, metadata, nil, nil)
	if err != nil {
		log.Fatalf("Error querying documents: %v\n", err)
	}

	return data.Documents[0]
}
