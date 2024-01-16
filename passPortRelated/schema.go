package main

//"time"

//"go.mongodb.org/mongo-driver/bson"
//"go.mongodb.org/mongo-driver/bson/primitive"
//"go.mongodb.org/mongo-driver/mongo"
//"go.mongodb.org/mongo-driver/mongo/options"

func Createpassport() {
	Passport := PassPort{
		ItemName:     "Your Item Name",
		Origin:       "Your Origin",
		IsNew:        true,
		LinkMadeFrom: []string{"Link1", "Link2"},
		LinkMakes:    []string{"Link3", "Link4"},
		CreationDate: "2024-01-16",
	}

	/*bsonData, err := bson.Marshal(Passport)
	if err != nil {
		log.Fatal(err)
	}

	var jsonData map[string]interface{}
	err = bson.Unmarshal(bsonData, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("BSON representation:", bsonData)

	jsonString, _ := json.MarshalIndent(jsonData, "", "  ")
	fmt.Printf("JSON data (for better readability):\n%s\n", jsonString)*/

}

/*
func UpdatePassport(IsNew bool,LinkMadeFrom []string, LinkMakes []string){

	IsNew: IsNew,
	linkMadefromLinkMadeFrom: LinkMadeFrom,

}
*/

/*var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createPassport": &graphql.Field{
			Type:        createPassport,
			Description: "Create a new passport",
			Args:        graphql.FieldConfigArgument{
				// Define any arguments needed for creating a passport
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Implement logic to create a passport and return its data
				// Access arguments using p.Args
				// Example:
				// itemID, _ := p.Args["ItemID"].(int)
				// itemName, _ := p.Args["ItemName"].(string)
				// ... (similarly for other fields)

				// Perform the actual creation logic (not implemented in this example)

				// Return the created passport data
				return map[string]interface{}{
					"ItemID":       1,
					"ItemName":     "Example",
					"Origin":       "ExampleOrigin",
					"IsNew":        true,
					"LinkMadeFrom": []string{"Link1", "Link2"},
					"LinkMakes":    []string{"Link3", "Link4"},
					"CreationDate": "2022-01-16",
				}, nil
			},
		},
	},
})

func main() {

	//ny := createPassport(1, "Sten", "malm", true, []string{}, []string{}, "2")

	//fmt.Println(ny)
	//fmt.Println(ny)
	//fmt.Println(ny)

}*/
