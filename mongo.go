package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/landdron/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
	"sync"
	"time"
)

var channelName sync.Map

const (
	URL string = "mongodb://localhost:27017"
)

type Param struct {
	UserID string `json:"userid" bson:"userid"`
	Data string `json:"data" bson:"data"`
	GameID string `json:"gameID" bson:"gameID"`
	GameName string `json:"gameName" bson:"gameName"`
	Channel  string `json:"channel" bson:"channel"`    //渠道
}

type FuncItem struct {
	FuncName string `json:"funcName" bson:"funcName"`
	Param1 string `json:"param1" bson:"param1"`
	Param2 string `json:"param2" bson:"param2"`
	Param3 string `json:"param3" bson:"param3"`
	TimeStamp int64 `json:"timestamp" bson:"timestamp"`
}

var log_channel chan *Param


func insert_dot_log(param *Param) {
	go func() {
		log_channel <- param
	}()
}

func loop_mongo() {
	for param := range log_channel {
		store_db(param)
	}
}

func store_db(param *Param) {
	channel_game_map,ok := channelName.Load(param.Channel)
	client := get_mongo_client()
	if client == nil {
		return
	}
	defer client.Disconnect(context.Background())

	if !ok {
		game_name_map := make(map[string]string)
		game_name_map[param.GameID] = param.GameName
		channelName.Store(param.Channel,game_name_map)
		store_channel_game(client,param.Channel,param.GameID,param.GameName)
	} else {
		game_name_map := reflect.ValueOf(channel_game_map).Interface().(map[string]string)
		_,ok = game_name_map[param.GameID]
		if !ok {
			store_channel_game(client,param.Channel,param.GameID,param.GameName)
			game_name_map[param.GameID] = param.GameName
		}
	}
	store_log(client,param)
}

func get_mongo_client()*mongo.Client {
	client,err := mongo.NewClient(options.Client().ApplyURI(URL))
	if err != nil {
		log.Error("new client error",err.Error())
		return nil
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Error("connect error",err.Error())
		return nil
	}
	return client
}

func store_channel_game(client *mongo.Client,channelName,gameID,gameName string) {
	collection := client.Database("dotgame").Collection("channelgame")

	collection.InsertOne(context.TODO(),bson.D{
		{"channel",channelName},
		{"gameID",gameID},
		{"gameName",gameName}})
}

func store_log(client *mongo.Client,param *Param) {
	///collection_name := time.Now().Format("2006_01_02")
	/////collection := client.Database("dotgame").Collection(collection_name)
	/*
	b := make(bsonrw.SliceWriter,0,100)
	wm,_ := bsonrw.NewBSONValueWriter(&b)
	encoder,_ := bson.NewEncoder(wm)
	encoder.Encode(param)
	fmt.Println("buf=================",string(b))
	 */


	game_data := bson.A{}

	func_list := strings.Split(param.Data,"&")

	for _,funcData := range func_list {
		var unFunc FuncItem
		json.Unmarshal([]byte(funcData),&unFunc)
	    if unFunc.FuncName == "" || len(unFunc.FuncName) == 0{
	    	continue
		}
		game_data = append(game_data,bson.D{
			{"funcName",unFunc.FuncName},
			{"param1",unFunc.Param1},
			{"param2",unFunc.Param2},
			{"param3", unFunc.Param3},
			{"timestamp",unFunc.TimeStamp}})
	}

	insert_data := bson.D{
		{"userid",param.UserID},
		{"gameID",param.GameID},
		{"gameName",param.GameName},
		{"channel",param.Channel},
		{"data",game_data},
	}

	ollection_name := time.Now().Format("log20060102")
	fmt.Println("collection_name",ollection_name)
	collection := client.Database("dotgame").Collection(ollection_name)
	collection.InsertOne(context.TODO(),insert_data)
}

func LoadDB() {

}

func TestMongo() {
	cred:= options.Credential{
		AuthMechanism: "PLAIN",
    	Username:      "root",
    	Password:      "KuaiYou2018",
	}
	ops := options.Client().ApplyURI("mongodb://root:KuaiZhiYou2018@dds-2ze13536a1a0f7841214-pub.mongodb.rds.aliyuncs.com:3717,dds-2ze13536a1a0f7842160-pub.mongodb.rds.aliyuncs.com:3717/admin?replicaSet=mgset-16845519").SetAuth(cred)
	client, err := mongo.Connect(context.TODO(), ops)
	if err != nil {
		fmt.Println("err=========== err",err.Error())
		return 
	}

	defer func() {
		if err = client.Disconnect(context.TODO());err != nil {
			fmt.Println("err w2",err.Error())
		}
	}()

	//////if err = client.Ping()
}

func init() {
	log_channel = make(chan *Param)
	go loop_mongo()
}