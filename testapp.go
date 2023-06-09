package main

import (
	"errors"
	"fmt"
	"strings"
)

type Artifact struct {
	Value    string
	Size     int
	ObjectId string
}

type Request struct {
	Method   string
	UrlPath  string
	Artifact Artifact
}

type Response struct {
	statusCode int
	Value      string
	OID        string
	Size       int
}

/*
This is a rather naive function for figuring out the
repository and objectID from the URL.
*/
func urlMatch(url string) (repository string, objectID string) {
	fragments := strings.SplitN(url, "/", -1)
	repository = fragments[2]
	objectID = ""
	if len(fragments) > 3 {
		objectID = fragments[3]
	}
	return repository, objectID
}

func main() {
	// We'll store the data in memory in a map.
	storage := make(map[string]string)

	err := testPut(storage)
	if err != nil {
		panic(err)
	}

	err = testGet(storage)
	if err != nil {
		panic(err)
	}
	err = testDelete(storage)
	if err != nil {
		panic(err)
	}
}

func Backend(request Request, storage map[string]string) Response {
	response := Response{}
	repository, objectID := urlMatch(request.UrlPath)
	switch request.Method {
	/*
		Download an Object

		GET /data/{repository}/{objectID}

		Response

		Status: 200 OK
		{object data}
	*/
	case "GET":
		/*
			This implementation of GET is incomplete at this time and won't
			pass the tests, please improve it.
		*/
		response = Response{
			statusCode: 200,
			Value:      storage[objectID],
			OID:        objectID,
			Size:       len(storage[objectID]),
		}
	}

	fmt.Println(request.Method + " repository: " + repository + " objectID: " + objectID)
	return response
}

func testPut(storage map[string]string) error {

	a1 := Artifact{
		Value: "something",
	}

	req1 := Request{
		Method:   "PUT",
		UrlPath:  "/data/codingtest",
		Artifact: a1,
	}

	res1 := Backend(req1, storage)

	a2 := Artifact{
		Value: "other",
	}

	req2 := Request{
		Method:   "PUT",
		UrlPath:  "/data/codingtest",
		Artifact: a2,
	}

	res2 := Backend(req2, storage)

	if res1.OID == res2.OID {
		return errors.New("expected to have unique oid")
	}

	if res1.Size != len(a1.Value) {
		return errors.New(fmt.Sprintf("expected a size of %d, got %d", len(a1.Value), res1.Size))
	}

	if res2.Size != len(a2.Value) {
		return errors.New(fmt.Sprintf("expected a size of %d, got %d", len(a2.Value), res2.Size))
	}

	return nil
}

func testGet(storage map[string]string) error {
	a1 := Artifact{
		Value: "something2",
	}

	req1 := Request{
		Method:   "PUT",
		UrlPath:  "/data/codingtest",
		Artifact: a1,
	}

	res1 := Backend(req1, storage)
	fmt.Printf("put response: %v \n", res1)

	if res1.statusCode != 201 {
		return errors.New(fmt.Sprintf("expected status of 200, got %d", res1.statusCode))
	}

	getReq := Request{
		Method:  "GET",
		UrlPath: fmt.Sprintf("/data/codingtest/%s", res1.OID),
	}
	fmt.Printf("get resquest: %v \n", getReq)
	getResponse := Backend(getReq, storage)
	fmt.Printf("get response: %v \n", getResponse)

	if getResponse.statusCode != 200 {
		return errors.New(fmt.Sprintf("expected status of 200, got %d", getResponse.statusCode))
	}
	if getResponse.Value != a1.Value {
		return errors.New(fmt.Sprintf("expected content of %s, got %s", a1.Value, getResponse.Value))
	}

	return nil
}

func testDelete(storage map[string]string) error {
	a1 := Artifact{
		Value: "something3",
	}

	req1 := Request{
		Method:   "PUT",
		UrlPath:  "/data/codingtest",
		Artifact: a1,
	}

	res1 := Backend(req1, storage)
	fmt.Printf("put response: %v \n", res1)

	deleteRequest := Request{
		Method:  "DELETE",
		UrlPath: fmt.Sprintf("/data/codingtest/%s", res1.OID),
	}

	deleteResponse := Backend(deleteRequest, storage)
	fmt.Printf("put response: %v \n", deleteResponse)

	if deleteResponse.statusCode != 200 {
		return errors.New(fmt.Sprintf("expected status of 200, got %d", deleteResponse.statusCode))
	}

	getReq := Request{
		Method:  "GET",
		UrlPath: fmt.Sprintf("/data/codingtest/%s", res1.OID),
	}
	fmt.Printf("get resquest: %v \n", getReq)
	getResponse := Backend(getReq, storage)
	fmt.Printf("get response: %v \n", getResponse)

	if getResponse.statusCode != 404 {
		return errors.New(fmt.Sprintf("expected status of 404, got %d", getResponse.statusCode))
	}

	deleteRequest2 := Request{
		Method:  "DELETE",
		UrlPath: fmt.Sprintf("/data/codingtest/%s", res1.OID),
	}

	deleteResponse2 := Backend(deleteRequest2, storage)
	fmt.Printf("put response: %v \n", deleteResponse2)

	if deleteResponse2.statusCode != 404 {
		return errors.New(fmt.Sprintf("expected status of 200, got %d", deleteResponse2.statusCode))
	}

	return nil
}
