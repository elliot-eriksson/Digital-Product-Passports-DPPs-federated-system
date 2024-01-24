package main

import (
	"fmt"
	//"time"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	//"github.com/graphql-go/graphql"
)

/* MongoDB object to interact with the PassPort struct
type PassPortMongo struct {
	Collection *mgo.Collection
}*/

func createPassport(ItemID int, ItemName string, Origin string, IsNew bool, LinkMadeFrom []string, LinkMakes []string, CreationDate string) *PassPort {
	return &PassPort{
		ItemID:       ItemID,
		ItemName:     ItemName,
		Origin:       Origin,
		IsNew:        IsNew,
		LinkMadeFrom: LinkMadeFrom,
		LinkMakes:    LinkMakes,
		CreationDate: CreationDate,
	}
}

func (value *PassPort) UpdatePassport(IsNewTemp bool, LinkMadeFromTemp []string, LinkMakesTemp []string) {
	value.IsNew = IsNewTemp
	value.LinkMadeFrom = LinkMadeFromTemp
	value.LinkMakes = LinkMakesTemp
}

func (value *PassPort) UpdateIsNew(IsNewTemp bool) {
	value.IsNew = IsNewTemp

}
func (value *PassPort) UpdateLinkMadeFrom(LinkMadeFromTemp []string) {
	value.LinkMadeFrom = LinkMadeFromTemp
}
func (value *PassPort) UpdateLinkMakes(LinkMakesTemp []string) {
	value.LinkMakes = LinkMakesTemp
}

//func (value *PassPort) UpdatePassport(IsNewTemp bool, LinkMadeFromTemp []string, LinkMakesTemp []string) {
//value.IsNew = (value *PassPort) UpdateIsNew(IsNewTemp)
//ny.UpdatePassport(false, []string{"123"}, []string{"123"})
//value.isNewTemp = updateIsNew()

//}

func main() {

	ny := createPassport(1, "Sten", "malm", true, []string{}, []string{}, "2")
	//fmt.Println("ItemID:", ny.ItemID)

	fmt.Println(ny)
	ny.UpdateIsNew(false)
	ny.UpdateLinkMadeFrom([]string{"123456789"})
	ny.UpdateLinkMakes([]string{"12345678"})
	//ny.UpdatePassport(false, []string{"123"}, []string{"123"})
	fmt.Println(ny)

	//fmt.Println(ny)

}
