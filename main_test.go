package main

import (
	"FunNow/url-shortener/constants"
	"FunNow/url-shortener/models"
	"FunNow/url-shortener/routers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ddosify/go-faker/faker"

	"github.com/stretchr/testify/assert"
)

func TestPostShortenUrlHandlerShortenNormalData(t *testing.T) {
	r := routers.InitRouter()
	faker := faker.NewFaker()
	longUrl := faker.RandomUrl()
	mapData := map[string]interface{}{
		"longUrl": longUrl,
	}
	data, _ := json.Marshal(mapData)
	var jsonStr = []byte(string(data))
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var result models.UrlElement
	json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, longUrl, result.LongUrl)
	assert.True(t, strings.HasPrefix(result.ShortUrl, constants.BaseUrl))
}

func TestPostShortenUrlHandlerShortenReturnSameShortUrl(t *testing.T) {
	r := routers.InitRouter()
	faker := faker.NewFaker()
	longUrl := faker.RandomUrl()
	mapData := map[string]interface{}{
		"longUrl": longUrl,
	}
	data, _ := json.Marshal(mapData)
	var jsonStr = []byte(string(data))
	firstReq, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonStr))
	firstW := httptest.NewRecorder()
	r.ServeHTTP(firstW, firstReq)

	var firstResult models.UrlElement
	json.Unmarshal(firstW.Body.Bytes(), &firstResult)

	assert.Equal(t, http.StatusOK, firstW.Code)
	assert.Equal(t, longUrl, firstResult.LongUrl)
	assert.True(t, strings.HasPrefix(firstResult.ShortUrl, constants.BaseUrl))

	expectedShortUrl := firstResult.ShortUrl
	secondReq, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonStr))
	secondW := httptest.NewRecorder()
	r.ServeHTTP(secondW, secondReq)

	var secondResult models.UrlElement
	json.Unmarshal(firstW.Body.Bytes(), &secondResult)

	assert.Equal(t, http.StatusOK, secondW.Code)
	assert.Equal(t, expectedShortUrl, secondResult.ShortUrl)
}

func TestPostShortenUrlHandlerShortenIncorrectRequestBody(t *testing.T) {
	r := routers.InitRouter()
	faker := faker.NewFaker()
	longUrl := faker.RandomUrl()
	mapData := map[string]interface{}{
		"wrongKey": longUrl,
	}
	data, _ := json.Marshal(mapData)
	var jsonStr = []byte(string(data))
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var errorResult models.ErrorTemplate
	json.Unmarshal(w.Body.Bytes(), &errorResult)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, errorResult.Error, "Invalid request body.")
}

func TestPostShortenUrlHandlerShortenInvalidUrl(t *testing.T) {
	r := routers.InitRouter()
	mapData := map[string]interface{}{
		"longURL": "INVALID URL",
	}
	data, _ := json.Marshal(mapData)
	var jsonStr = []byte(string(data))
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var errorResult models.ErrorTemplate
	json.Unmarshal(w.Body.Bytes(), &errorResult)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, errorResult.Error, "Invalid url format.")
}
