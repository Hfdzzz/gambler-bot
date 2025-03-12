package entity

type User struct {
	Username	string `bson:"username"`
	Money 		int `bson:"money"`
}