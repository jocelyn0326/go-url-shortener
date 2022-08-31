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

func InitShortenUtilities(mapData map[string]interface{}, isResponse200 bool) (*httptest.ResponseRecorder, models.UrlElement, models.ErrorTemplate) {
	r := routers.InitRouter()
	data, _ := json.Marshal(mapData)
	var jsonStr = []byte(string(data))
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var result200 models.UrlElement
	var result400 models.ErrorTemplate
	if isResponse200 {
		json.Unmarshal(w.Body.Bytes(), &result200)
	} else {
		json.Unmarshal(w.Body.Bytes(), &result400)
	}

	return w, result200, result400
}

func TestPostShortenUrlHandlerShortenNormalData(t *testing.T) {
	faker := faker.NewFaker()
	longUrl := faker.RandomUrl()
	mapData := map[string]interface{}{
		"longUrl": longUrl,
	}
	w, result200, _ := InitShortenUtilities(mapData, true)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, longUrl, result200.LongUrl)
	assert.True(t, strings.HasPrefix(result200.ShortUrl, constants.BaseUrl))
}

func TestPostShortenUrlHandlerShortenReturnSameShortUrl(t *testing.T) {

	faker := faker.NewFaker()
	longUrl := faker.RandomUrl()
	mapData := map[string]interface{}{
		"longUrl": longUrl,
	}

	firstW, first200Result, _ := InitShortenUtilities(mapData, true)

	assert.Equal(t, http.StatusOK, firstW.Code)
	assert.Equal(t, longUrl, first200Result.LongUrl)
	assert.True(t, strings.HasPrefix(first200Result.ShortUrl, constants.BaseUrl))

	secondW, second200Result, _ := InitShortenUtilities(mapData, true)

	assert.Equal(t, http.StatusOK, secondW.Code)
	assert.Equal(t, first200Result.ShortUrl, second200Result.ShortUrl)
}

func TestPostShortenUrlHandlerShortenIncorrectRequestBody(t *testing.T) {

	faker := faker.NewFaker()
	longUrl := faker.RandomUrl()
	mapData := map[string]interface{}{
		"wrongKey": longUrl,
	}
	w, _, result400 := InitShortenUtilities(mapData, false)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, result400.Error, "Invalid request body.")
}

func TestPostShortenUrlHandlerShortenInvalidUrl(t *testing.T) {
	mapData := map[string]interface{}{
		"longURL": "INVALID URL",
	}
	w, _, result400 := InitShortenUtilities(mapData, false)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, result400.Error, "Invalid url format.")
}
