package configs

import (
	"context"
	"fmt"
	"gambling-bot/internal/entity"
	"log"
	

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() {
	url := "mongodb://localhost:27017"

	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Cannot accest mongodb: ", err)
	}

	fmt.Println("Connect to mongodb successfully")

	Client = client
	
}

func GetCollection(collectionName string) *mongo.Collection {
	db := Client.Database("Caligula")
	return db.Collection(collectionName)
}

func CheckUserData(s *discordgo.Session, m *discordgo.MessageCreate, username string){
	//db := Client.Database("Caligula")
	collection := GetCollection("user")

	filter := bson.M{"username": username}
	existingUser := entity.User{}

	err := collection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "You are not register, type 'bj' to register")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Your money: %d", existingUser.Money))

}

func AddUserData(username string, money int){
	//db := Client.Database("Caligula")
	collection := GetCollection("user")

	filter := bson.M{"username": username}

	existingUser := entity.User{}

	err := collection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err == nil {
		fmt.Println("User sudah terdaftar")
		return
	}

	data := entity.User{
		Username: username,
		Money: money,
	}

	_, err = collection.InsertOne(context.TODO(), data)
	
	
	if err != nil {
		fmt.Println("Gagal menambahkan user ke database:", err)
		return
	}

	log.Println("Data berhasil disimpan untuk user: ", username)

}

func UpdateUserDataWin(username string, money int){
	//db := Client.Database("Caligula")
	collection := GetCollection("user")

	user := entity.User{}
	filter := bson.M{"username": username}

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		println("No user.")
		return
	}

	updateMoney := user.Money + money

	update := bson.M{"$set": bson.M{"money": updateMoney}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal("Fail update data: ", err)
	}

	log.Println("Succes update data for user: : ", username)
}

func UpdateUserDataLose(username string, money int){
	//db := Client.Database("Caligula")
	collection := GetCollection("user")

	user := entity.User{}
	filter := bson.M{"username": username}

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		println("No user.")
		return
	}

	updateMoney := user.Money - money

	update := bson.M{"$set": bson.M{"money": updateMoney}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal("Fail update data: ", err)
	}

	log.Println("Succes update data for user: : ", username)
}