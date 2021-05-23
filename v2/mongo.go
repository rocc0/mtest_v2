package v2

//
//import (
//	"log"
//
//	"gopkg.in/mgo.v2"
//)
//
//func uploadDataToMongo() error {
//	return nil
//}
//
//
//func getAdministrativeActions(col string) (interface{}, error) {
//	session, err := mgo.Dial("mongodb://adder:password@192.168.99.100:27017")
//	if err != nil {
//		panic(err)
//	}
//	defer session.Close()
//
//	var result interface{}
//
//	c := session.DB("questions").C("i"+col)
//
//	err = c.Find(nil).All(&result)
//
//	if err != nil {
//		log.Print(err.Error(),"data not found")
//	}
//
//	return result,  nil
//}
//
//func saveUserMtestsToMongo(id int, u UserMtest) error {
//	session, err := mgo.Dial("mongodb://adder:password@192.168.99.100:27017")
//	if err != nil {
//		panic(err)
//	}
//	defer session.Close()
//
//
//
//
//
//}
