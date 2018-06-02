package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "database/sql"

	_ "github.com/go-sql-driver/mysql"
	// Importing the plugins enables collection of AWS resource information at runtime.
	// Every plugin should be imported after "github.com/aws/aws-xray-sdk-go/xray" library.

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
)

type person struct {
	name string
	age  int
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	SQL_URI := os.Getenv("MY_SQL_URI")
	if SQL_URI == "" {
		log.Println("No SQL_URI defined")
		return events.APIGatewayProxyResponse{Body: string("No SQL_URI defined"), StatusCode: 500}, nil
	}

	fmt.Println(os.Getenv("MY_SQL_URI"))
	db, err := xray.SQL("mysql", os.Getenv("MY_SQL_URI"))
	defer db.Close()

	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	age := r1.Intn(100)
	rows, err := db.Query(ctx, "SELECT name FROM customers WHERE age=?", age)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	defer rows.Close()

	var results = []person{}

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		}
		fmt.Printf("%s is %d\n", name, age)
		results = append(results, person{name: name, age: age})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(results)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	return events.APIGatewayProxyResponse{Body: string(data), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
