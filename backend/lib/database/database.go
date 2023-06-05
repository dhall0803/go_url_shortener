package database

import (
	"database/sql"
	"os"
)

type ShortUrl struct {
	Id       string `json:"id"`
	UserId   string `json:"user_id"`
	LongUrl  string `json:"long_url"`
	ShortUrl string `json:"short_url"`
}

// queryWrapper provides the code needed to connect to the database and execute a query. The third argument is a call back function that determines what to do with the returned rows
func queryWrapper(connStr, query string, parameters any, callback func(*sql.Rows) ([]interface{}, error)) ([]interface{}, error) {
	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Perform a query
	rows, err := db.Query(query, parameters)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := callback(rows)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetShortUrl(userId, longUrl string) ([]ShortUrl, error) {
	query := "SELECT * FROM short_urls WHERE user_id = $1 AND long_url = $2"
	parameters := []interface{}{userId, longUrl}

	rowsToShortUrls := func(rows *sql.Rows) ([]interface{}, error) {
		var shortUrls []interface{}
		for rows.Next() {
			var shortUrl ShortUrl
			err := rows.Scan(&shortUrl.Id, &shortUrl.UserId, &shortUrl.LongUrl, &shortUrl.ShortUrl)
			if err != nil {
				return nil, err
			}
			shortUrls = append(shortUrls, shortUrl)
		}
		return shortUrls, nil
	}

	result, err := queryWrapper(os.Getenv("CONNECTION_STRING"), query, parameters, rowsToShortUrls)

	if err != nil {
		return nil, err
	}

	finalResult := make([]ShortUrl, len(result))

	for i, shortUrl := range result {
		finalResult[i] = shortUrl.(ShortUrl)
	}
	return finalResult, nil
}

func CreateShortUrl(newShortUrl ShortUrl) error {
	query := "INSERT INTO short_urls (id, user_id, long_url, short_url) VALUES ($1, $2, $3)"
	parameters := []interface{}{newShortUrl.Id, newShortUrl.UserId, newShortUrl.LongUrl, newShortUrl.ShortUrl}

	_, err := queryWrapper(os.Getenv("CONNECTION_STRING"), query, parameters, func(rows *sql.Rows) ([]interface{}, error) {
		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}
