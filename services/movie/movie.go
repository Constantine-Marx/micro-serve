package movie

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Movie struct {
	ID     int     `json:"ID"`
	Title  string  `json:"Title"`
	Rating float64 `json:"Rating"`
}

type MovieService interface {
	GetMovieByID(ctx context.Context, args *Movie, reply *Movie) error
	CreateMovie(ctx context.Context, movie *Movie, reply *struct{}) error
}

type movieServiceImpl struct {
	db     *sql.DB
	movies map[int]*Movie
}

func (s *movieServiceImpl) GetMovieByID(ctx context.Context, args *Movie, reply *Movie) error {
	row := s.db.QueryRow("SELECT id, title, rating FROM movies WHERE id = ?", args.ID)
	err := row.Scan(&reply.ID, &reply.Title, &reply.Rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("movie not found")
		}
		return err
	}
	return nil
}

func (s *movieServiceImpl) CreateMovie(ctx context.Context, movie *Movie, reply *struct{}) error {
	result, err := s.db.Exec("INSERT INTO movies (id, title, rating) VALUES (?, ?, ?)", movie.ID, movie.Title, movie.Rating)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("movie already exists")
	}
	return nil
}

func NewMovieService(db *sql.DB) MovieService {
	return &movieServiceImpl{
		db: db,
	}
}