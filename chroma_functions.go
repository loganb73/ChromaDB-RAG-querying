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

	openaiEf, err := chromaOpenai.NewOpenAIEmbeddingFunction(os.Getenv("openai_key_temp"))
	if err != nil {
		log.Fatalf("Error creating OpenAI embedding function: %s \n", err)
	}

	//create collections
	collection, err := client.CreateCollection(context.TODO(), "full-collection", map[string]interface{}{}, true, openaiEf, types.L2)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	profCollection, err := client.CreateCollection(context.TODO(), "professor-collection", map[string]interface{}{}, true, openaiEf, types.L2)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	// classCollection, err := client.CreateCollection(context.TODO(), "class-collection", map[string]interface{}{"professor": "value1"}, true, openaiEf, types.L2)
	// if err != nil {
	// 	log.Fatalf("Failed to create collection: %v", err)
	// }

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

	//make maps to avoid duplicate inserts
	profMap := make(map[string]bool)

	//loop over classes and add data to record sets
	i := 0
	for _, class := range classes {
		classJson, err := json.Marshal(class)
		if err != nil {
			fmt.Printf("error marshalling struct: %s\n", err.Error())
		}
		rs.WithRecord(types.WithDocument(string(classJson)))

		professorFullName := class.PrimaryInstructorFirstName + class.PrimaryInstructorLastName

		if professorFullName != "" { //don't add for classes without a professor
			_, exists := profMap[professorFullName]
			if exists {
				continue
			} else {
				professorsRs.WithRecord(types.WithDocument(string(professorFullName)))
				if err != nil {
					fmt.Printf("error adding professor: %s\n", err.Error())
				}
				profMap[professorFullName] = true
			}
		}

		//build in batches
		if ((i % 500) == 0) && i > 0 {
			_, err = rs.BuildAndValidate(context.TODO())
			if err != nil {
				fmt.Printf("Error validating record set full: %s \n", err)
				log.Fatalf("Error validating record set full: %s \n", err)
			}

			_, err = professorsRs.BuildAndValidate(context.TODO())
			if err != nil {
				fmt.Printf("Error validating record set profs: %s \n", err)
				log.Fatalf("Error validating record set profs: %s \n", err)
			}
			fmt.Printf("inserted %d docs\n", i)
		}
		i++
	}

	//insert stragglers (not caught by last batch insert)
	_, err = rs.BuildAndValidate(context.TODO())
	if err != nil {
		fmt.Printf("Error validating record set full post loop: %s \n", err)
		log.Fatalf("Error validating record set full post loop: %s \n", err)
	}
	_, err = professorsRs.BuildAndValidate(context.TODO())
	if err != nil {
		fmt.Printf("Error validating record set profs post loop: %s \n", err)
		log.Fatalf("Error validating record set profs post loop: %s \n", err)
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

	// Count the number of documents in the collection
	countDocs, qrerr := profCollection.Count(context.TODO())
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
