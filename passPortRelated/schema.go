package main

//"time"

//"go.mongodb.org/mongo-driver/bson"
//"go.mongodb.org/mongo-driver/bson/primitive"
//"go.mongodb.org/mongo-driver/mongo"
//"go.mongodb.org/mongo-driver/mongo/options"

// Funktion för att skapa ett passport.
// TODO: Ändra så att funktionen tar query parametrar istället för hårdkodad data
func Createpassport() PassPort {
	Passport := PassPort{
		ItemName:     "Your Item Name",
		Origin:       "Your Origin",
		IsNew:        true,
		LinkMadeFrom: []string{"Link1", "Link2"},
		LinkMakes:    []string{"Link3", "Link4"},
		LinkEvents:   []string{"Link5", "Link6"},
		CreationDate: "2024-01-16",
	}
	return Passport
}

// Funktion för att lägga till remanufacture event till ett item
/*func RemanufactureEvents() {

	var NewEvent = fmt.Scan()

}*/
