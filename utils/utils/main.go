package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
)

type MovieData struct {
	Data []struct {
		CreatedAt   int64  `json:"createdAt"`
		UpdatedAt   int64  `json:"updatedAt"`
		ID          string `json:"id"`
		Poster      string `json:"poster"`
		Name        string `json:"name"`
		Genre       string `json:"genre"`
		Description string `json:"description"`
		Language    string `json:"language"`
		Country     string `json:"country"`
		Lang        string `json:"lang"`
		ShareImage  string `json:"shareImage"`
		Movie       string `json:"movie"`
	} `json:"data"`
	CreatedAt    int64  `json:"createdAt"`
	UpdatedAt    int64  `json:"updatedAt"`
	ID           string `json:"id"`
	OriginalName string `json:"originalName"`
	ImdbVotes    int    `json:"imdbVotes"`
	ImdbRating   string `json:"imdbRating"`
	RottenRating string `json:"rottenRating"`
	RottenVotes  int    `json:"rottenVotes"`
	Year         string `json:"year"`
	ImdbId       string `json:"imdbId"`
	Alias        string `json:"alias"`
	DoubanId     string `json:"doubanId"`
	Type         string `json:"type"`
	DoubanRating string `json:"doubanRating"`
	DoubanVotes  int    `json:"doubanVotes"`
	Duration     int    `json:"duration"`
	DateReleased string `json:"dateReleased"`
}

type UtilService interface {
	LoadMoviesFromFile(ctx context.Context, filename string, reply *struct{}) error
}

type movieServiceImpl struct {
	db *sql.DB
}

func (s *movieServiceImpl) LoadMoviesFromFile(ctx context.Context, filename string, reply *struct{}) error {
	// 读取文件内容
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// 解析 JSON 数据
	var moviesData []MovieData
	err = json.Unmarshal(data, &moviesData)
	if err != nil {
		return fmt.Errorf("error parsing JSON data: %w", err)
	}

	// 遍历电影数据并插入到 movie 表中
	for _, movieData := range moviesData {
		// 检查并转换 "null" 值
		if movieData.RottenRating == "null" {
			movieData.RottenRating = "0"
		}
		// 插入数据
		_, err := s.db.Exec("INSERT INTO movies (original_name, imdb_votes, imdb_rating, rotten_rating, rotten_votes, year, imdb_id, alias, douban_id, type, douban_rating, douban_votes, duration, date_released, poster, name, genre, description, language, country, lang, share_image) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			movieData.OriginalName, movieData.ImdbVotes, movieData.ImdbRating, movieData.RottenRating, movieData.RottenVotes,
			movieData.Year, movieData.ImdbId, movieData.Alias, movieData.DoubanId, movieData.Type, movieData.DoubanRating,
			movieData.DoubanVotes, movieData.Duration, movieData.DateReleased, movieData.Data[0].Poster, movieData.Data[0].Name,
			movieData.Data[0].Genre, movieData.Data[0].Description, movieData.Data[0].Language, movieData.Data[0].Country,
			movieData.Data[0].Lang, movieData.Data[0].ShareImage)
		if err != nil {
			return fmt.Errorf("error inserting movie data: %w", err)
		}
	}

	reply = &struct{}{}

	return nil
}
func NewExtractService(db *sql.DB) UtilService {
	return &movieServiceImpl{
		db: db,
	}
}
