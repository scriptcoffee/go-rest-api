package main

import (
	"testing"
	"net/http/httptest"
	"github.com/golang/mock/gomock"
	"net/http"
	"errors"
	"net/url"
	"strings"
)

func TestGetPeopleReturnsEmptyList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store = NewMockStore(mockCtrl)
	mockStore := store.(*MockStore)

	mockStore.EXPECT().getPeople().Return([]Person{}, nil).Times(1)

	req, err := http.NewRequest("GET", "/people", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPeople)

	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusNoContent
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestGetPeopleReturnsList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store = NewMockStore(mockCtrl)
	mockStore := store.(*MockStore)

	pList := []Person{}
	pList = append(pList, Person{Id:1, Name:"Paul", PhoneNr:"643265776357948984"})
	pList = append(pList, Person{Id:3, Name:"Peter", PhoneNr:"24525345626"})
	pList = append(pList, Person{Id:4, Name:"Alice", PhoneNr:"12343463462345243"})

	mockStore.EXPECT().getPeople().Return(pList, nil).Times(1)

	req, err := http.NewRequest("GET", "/people", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPeople)

	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := `[{"id":1,"name":"Paul","phoneNr":"643265776357948984"},{"id":3,"name":"Peter","phoneNr":"24525345626"},{"id":4,"name":"Alice","phoneNr":"12343463462345243"}]`
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestGetPeopleReturnsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store = NewMockStore(mockCtrl)
	mockStore := store.(*MockStore)

	mockStore.EXPECT().getPeople().Return([]Person{}, errors.New("getPeopleError")).Times(1)

	req, err := http.NewRequest("GET", "/people", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPeople)

	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusInternalServerError
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestGetPersonReturnsPerson(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store = NewMockStore(mockCtrl)
	mockStore := store.(*MockStore)

	mockStore.EXPECT().getPerson(3).Return(Person{Id:3, Name:"Peter", PhoneNr:"24525345626"}, nil).Times(1)

	req, err := http.NewRequest("GET", "/people/3", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := `{"id":3,"name":"Peter","phoneNr":"24525345626"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestGetPersonReturnsNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store = NewMockStore(mockCtrl)
	mockStore := store.(*MockStore)

	mockStore.EXPECT().getPerson(300).Return(Person{}, errors.New("Person not found")).Times(1)

	req, err := http.NewRequest("GET", "/people/300", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusNotFound
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestGetPersonReturnsBadRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/people/a", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusBadRequest
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestCreatePersonReturnsErrorBadRequest(t *testing.T) {
	v := url.Values{}
	v.Set("name", "Peter")

	req, err := http.NewRequest("POST", "/people", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusBadRequest
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestCreatePersonReturnsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store = NewMockStore(mockCtrl)
	mockStore := store.(*MockStore)

	p := Person{Name:"Peter", PhoneNr:"56468465613275"}
	mockStore.EXPECT().createPerson(p).Return(errors.New("createPersonError")).Times(1)

	v := url.Values{}
	v.Set("name", "Peter")
	v.Add("phoneNr", "56468465613275")

	req, err := http.NewRequest("POST", "/people", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusInternalServerError
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestCreatePersonReturnsOk(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store = NewMockStore(mockCtrl)
	mockStore := store.(*MockStore)

	p := Person{Name:"Peter", PhoneNr:"56468465613275"}
	mockStore.EXPECT().createPerson(p).Return(nil).Times(1)

	v := url.Values{}
	v.Set("name", "Peter")
	v.Add("phoneNr", "56468465613275")

	req, err := http.NewRequest("POST", "/people", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusCreated
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestUpdatePersonReturnsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store = NewMockStore(mockCtrl)
	mockStore := store.(*MockStore)

	p := Person{Id:1, Name:"Peter", PhoneNr:"56468465613275"}
	mockStore.EXPECT().updatePerson(p).Return(errors.New("Error in updatePerson")).Times(1)

	v := url.Values{}
	v.Set("name", "Peter")
	v.Add("phoneNr", "56468465613275")
	v.Add("id", "1")

	req, err := http.NewRequest("PUT", "/people/1", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusInternalServerError
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestUpdatePersonReturnsErrorIncompletePerson(t *testing.T) {
	v := url.Values{}
	v.Set("name", "Peter")

	req, err := http.NewRequest("PUT", "/people/1", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusBadRequest
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestUpdatePersonReturnsErrorMissingId(t *testing.T) {
	v := url.Values{}
	v.Set("name", "Peter")

	req, err := http.NewRequest("PUT", "/people", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusMethodNotAllowed
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestUpdatePersonReturnsErrorBadId(t *testing.T) {
	v := url.Values{}
	v.Set("name", "Peter")

	req, err := http.NewRequest("PUT", "/people/a", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusBadRequest
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestUpdatePersonReturnsOk(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store = NewMockStore(mockCtrl)
	mockStore := store.(*MockStore)

	p := Person{Id:1, Name:"Peter", PhoneNr:"56468465613275"}
	mockStore.EXPECT().updatePerson(p).Return(nil).Times(1)

	v := url.Values{}
	v.Set("name", "Peter")
	v.Add("phoneNr", "56468465613275")
	v.Add("id", "1")

	req, err := http.NewRequest("PUT", "/people/1", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestDeletePersonReturnsBadRequest(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/people/a", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusBadRequest
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestDeletePersonReturnsInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store = NewMockStore(mockCtrl)
	mockStore := store.(*MockStore)

	mockStore.EXPECT().deletePerson(1).Return(errors.New("Error in deletePerson")).Times(1)

	req, err := http.NewRequest("DELETE", "/people/1", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusInternalServerError
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestDeletePersonReturnsOk(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store = NewMockStore(mockCtrl)
	mockStore := store.(*MockStore)

	mockStore.EXPECT().deletePerson(1).Return(nil).Times(1)

	req, err := http.NewRequest("DELETE", "/people/1", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := createRouter()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	expectedBody := ``
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}