package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Post struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

func main() {

	for i := 0; i < 30000; i++ {
		//Nodo 1
		//request := getRequest("http://localhost:8081/rest/api/2/issue/")
		//CreateTicket(request)

		//Nodo2
		request := getRequest("http://localhost:8082/rest/api/2/issue/")
		CreateTicket(request)
	}

}

func CreateTicket(request *http.Request) {

	client := &http.Client{}

	res, err := client.Do(request)

	if err != nil {
		fmt.Println("Problema en la creacion  del request:", err)
		panic(err)
	}

	defer res.Body.Close()

	post := &Post{}
	derr := json.NewDecoder(res.Body).Decode(post)

	if derr != nil {
		panic(derr)
	}

	if res.StatusCode != http.StatusCreated {
		fmt.Println("No se ha podido crear el ticket")
		panic(res.Status)
	}

	//fmt.Pringln("Issue Creado")
	//fmt.Println("KEY:", post.Key)
	//fmt.Println("")

}

func getRequest(url string) *http.Request {

	posturl := url

	body := []byte(`{
		"fields": {
		   "project":
		   {
			  "id": "10000"
		   },
		   "summary": "No REST for the Wicked.",
		   "description": "Creating of an issue using IDs for projects and issue types using the REST API",
		   "issuetype": {
			  "id": "10002"
		   }
	   }
	}`)

	request, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))

	if err != nil {
		fmt.Println("Problema en el setup del request")
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Bearer <USER APIKEY>")

	return request
}
