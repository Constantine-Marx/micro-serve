// Package movie services/movie/movie.go
package movie

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Movie struct {
	ID           int     `json:"id"`
	OriginalName string  `json:"original_name"`
	ImdbVotes    int     `json:"imdb_votes"`
	ImdbRating   float64 `json:"imdb_rating"`
	RottenRating string  `json:"rotten_rating"`
	RottenVotes  int     `json:"rotten_votes"`
	Year         int     `json:"year"`
	ImdbId       string  `json:"imdb_id"`
	Alias        string  `json:"alias"`
	DoubanId     int     `json:"douban_id"`
	Type         string  `json:"type"`
	DoubanRating float64 `json:"douban_rating"`
	DoubanVotes  int     `json:"douban_votes"`
	Duration     int     `json:"duration"`
	DateReleased string  `json:"date_released"`
	Poster       string  `json:"poster"`
	Name         string  `json:"name"`
	Genre        string  `json:"genre"`
	Description  string  `json:"description"`
	Language     string  `json:"language"`
	Country      string  `json:"country"`
	Lang         string  `json:"lang"`
	ShareImage   string  `json:"share_image"`
}

type MovieService interface {
	GetMovieByID(ctx context.Context, args *Movie, reply *Movie) error
	GetMoviesByPage(ctx context.Context, args *PageRequest, reply *[]Movie) error
}

type PageRequest struct {
	Page int `json:"page"`
}

type movieServiceImpl struct {
	db     *sql.DB
	movies map[int]*Movie
}

func (s *movieServiceImpl) GetMovieByID(ctx context.Context, args *Movie, reply *Movie) error {
	row := s.db.QueryRow("SELECT id, original_name, imdb_votes, imdb_rating, rotten_rating, rotten_votes, year, imdb_id, alias, douban_id, type, douban_rating, douban_votes, duration, date_released, poster, name, genre, description, language, country, lang, share_image FROM movies WHERE id = ?", args.ID)
	err := row.Scan(&reply.ID, &reply.OriginalName, &reply.ImdbVotes, &reply.ImdbRating, &reply.RottenRating, &reply.RottenVotes, &reply.Year, &reply.ImdbId, &reply.Alias, &reply.DoubanId, &reply.Type, &reply.DoubanRating, &reply.DoubanVotes, &reply.Duration, &reply.DateReleased, &reply.Poster, &reply.Name, &reply.Genre, &reply.Description, &reply.Language, &reply.Country, &reply.Lang, &reply.ShareImage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("movie not found")
		}
		return err
	}
	return nil
}

func (s *movieServiceImpl) GetMoviesByPage(ctx context.Context, args *PageRequest, reply *[]Movie) error {
	offset := (args.Page - 1) * 10
	rows, err := s.db.Query("SELECT id, original_name, imdb_votes, imdb_rating, rotten_rating, rotten_votes, year, imdb_id, alias, douban_id, type, douban_rating, douban_votes, duration, date_released, poster, name, genre, description, language, country, lang, share_image FROM movies LIMIT 10 OFFSET ?", offset)
	if err != nil {
		log.Printf("Error fetching movies: %v", err) // 添加日志记录
		return err
	}
	if rows == nil {
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing movie rows: %v", err) // 添加日志记录
		}
	}(rows)

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.ID, &movie.OriginalName, &movie.ImdbVotes, &movie.ImdbRating, &movie.RottenRating, &movie.RottenVotes, &movie.Year, &movie.ImdbId, &movie.Alias, &movie.DoubanId, &movie.Type, &movie.DoubanRating, &movie.DoubanVotes, &movie.Duration, &movie.DateReleased, &movie.Poster, &movie.Name, &movie.Genre, &movie.Description, &movie.Language, &movie.Country, &movie.Lang, &movie.ShareImage)
		if err != nil {
			log.Printf("Error scanning movie row: %v", err) // 添加日志记录
			return err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error in movie rows: %v", err) // 添加日志记录
		return err
	}

	*reply = movies
	return nil
}

func NewMovieService(db *sql.DB) MovieService {
	return &movieServiceImpl{
		db: db,
	}
}
