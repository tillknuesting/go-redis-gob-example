package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type Data struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	dataMarshal := Data{
		Name: "Peter",
		Age:  55,
	}

	// Struct to JSON
	var jsonDataMarshal []byte

	jsonDataMarshal, err := json.Marshal(dataMarshal)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("json data: ", string(jsonDataMarshal))

	// JSON to Struct
	jsonDataUnmarshal := []byte(`{"name":"Peter","age":55}`)

	var dataUnmarshal Data
	err = json.Unmarshal(jsonDataUnmarshal, &dataUnmarshal)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("data unmarshalled from json", dataUnmarshal)

	dataEncode := Data{
		Name: "Peter",
		Age:  55,
	}

	// struct to Gob
	bufEn := &bytes.Buffer{}
	if err := gob.NewEncoder(bufEn).Encode(dataEncode); err != nil {
		panic(err)
	}
	// buf.Bytes()

	BufEnString := bufEn.String()

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	nameKey := dataEncode.Name + "key"

	err = rdb.Set(ctx, nameKey, BufEnString, 0).Err()
	if err != nil {
		log.Println(err)
	}

	val, err := rdb.Get(ctx, nameKey).Result()
	if err != nil {
		log.Println(err)
	}

	// Gob to Struct
	bufDe := &bytes.Buffer{}

	bufDe.WriteString(val)

	var dataDecode Data
	if err := gob.NewDecoder(bufDe).Decode(&dataDecode); err != nil {
		log.Println(err)
	}
	fmt.Println("data decoded from gob:", dataDecode)
}
