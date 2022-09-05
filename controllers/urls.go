package controllers

import (
	"FunNow/url-shortener/constants"
	"FunNow/url-shortener/db"
	"FunNow/url-shortener/models"
	"FunNow/url-shortener/utils"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
	"unicode/utf8"

	log "github.com/cihub/seelog"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
)

func SaveUrl(urlElement models.UrlElement) (models.UrlElement, error) {
	var urlCode string
	var resultUrlElement models.UrlElement
	redisClient := db.DatabaseClient.RedisClient

	// Step 1. Check if the origin long url exists
	result, err := redisClient.Get(constants.Ctx, urlElement.LongUrl).Result()

	// Step 2. If the long url exists in cache:
	if err == nil {
		// Return the data from cache directly
		json.Unmarshal([]byte(result), &resultUrlElement)
		return resultUrlElement, nil
	}

	// Step 3. If the long url does not exist in cache:
	// Generate a random url code until it's unique by get all the keys in.
	// Guarantee unique by using redis transactions optimistic locking.
	// Transactional function.
	txf := func(tx *redis.Tx) error {
		// It returns redis.Nil error when key does not exist.
		getKeyError := tx.Get(constants.Ctx, urlCode).Err()
		if getKeyError != redis.Nil {
			// urlCode exists in cache
			if getKeyError == nil {
				return errors.New("This Key already exists.")
			}

			// Raise the original error raised by tx.Get().Err()
			return getKeyError
		}

		// Finalize the urlElement and save it to cache.
		urlElement.ShortUrl = constants.BaseUrl + urlCode
		urlElement.CreateDateTime = time.Now()
		urlJsonValue, _ := json.Marshal(urlElement)

		// Operation is committed only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(constants.Ctx, func(pipe redis.Pipeliner) error {

			// First, set the url code into the cache for checking duplication
			pipe.Set(constants.Ctx, urlCode, urlJsonValue, 0)

			// Second, save the long url and the data
			pipe.Set(constants.Ctx, urlElement.LongUrl, urlJsonValue, 0)

			return nil
		})

		return err
	}

	const maxRetries = 1000
	var generateShortUrlSuccesses bool = false
	// Retry if the key exists, do the loop with a limit of 1000 times.
	for i := 0; i < maxRetries; i++ {

		// Generate the random string as an url code.
		urlCode = utils.RandString(constants.ShortUrlCodeLength)

		// Watch the key and check if it exists or the transaction is interrupted.
		err := redisClient.Watch(constants.Ctx, txf, urlCode)

		// The transaction is successful, break the loop.
		if err == nil {
			generateShortUrlSuccesses = true
			break
		}

	}

	if generateShortUrlSuccesses {
		return urlElement, nil
	}

	return urlElement, errors.New("Failed to generate a unique key.")

}

func PostShortenUrlHandler(c *gin.Context) {
	var urlElement models.UrlElement
	var result models.UrlElement

	// Check request body validation
	if bindError := c.BindJSON(&urlElement); bindError != nil {
		log.Info("Received invalid request body.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body."})
		return
	}

	// Check url character length. Stick to the lowest common denominator and limit their URLs to 2,083 characters in length.
	if utf8.RuneCountInString(urlElement.LongUrl) > 2083 {
		log.Info("Received invalid url. Url characters can not greater than 2083.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid url.Url characters can not greater than 2083."})
	}

	//Check url validation.
	_, parseError := url.ParseRequestURI(urlElement.LongUrl)
	if parseError != nil {
		log.Info("Received invalid url request, url: ", urlElement.LongUrl, parseError.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid url format."})
		return
	}

	// Save url into cache
	result, saveError := SaveUrl(urlElement)
	if saveError != nil {
		log.Error("Unexpected error raised: ", saveError.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown error."})
		return
	}
	c.JSON(http.StatusOK, result)
}
