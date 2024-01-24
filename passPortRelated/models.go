package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Definerar en struct för en produkt. ID är mongoDBs _id

type PassPort struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"` //ObjectID som skapas i mongoDB och är unikt.
	ItemID       int                `bson:"ItemID"`        //ItemID som ökar med ett för varje tillagt pass
	ItemName     string             `bson:"ItemName"`      //Namn för produkten som man får själv lägga till vid skapningen
	Origin       string             `bson:"Origin"`        //Samma för namn kanske ska ändras så det är länk etc för nu matas endast ett namn in
	IsNew        bool               `bson:"IsNew"`         //IsNew ska uppdateras senare med IPFS
	LinkMadeFrom []string           `bson:"LinkMadeFrom"`  //LinkMadeFrom ska vara länkar till exempelvis stål och trä som används i en spade
	LinkMakes    []string           `bson:"LinkMakes"`     //LinkMakes blir LinkMadeFrom fast för andra hållet alltså att stålet går till spaden
	LinkEvents   []string           `bson:"LinkEvents"`    //LinkEvents är eventen som kan hända exempelvis att en bil repareras etc
	Sensitive    []int              `bson:"Sensitive"`     //Sensitive är en array som kommer ändras beroende på om data ska kunna visas av privatpersoner etc
	CreationDate string             `bson:"CreationDate"`  //CreationDate skapas vid skapningen av produkten alltså dagens datum.
}
