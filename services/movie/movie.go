package movie

import "fmt"

type Movie struct {
	ID     int
	Title  string
	Rating float64
}

type MovieService interface {
	GetMovieByID(id int) (*Movie, error)
	CreateMovie(movie *Movie) error
}

type movieServiceImpl struct {
	movies map[int]*Movie
}

func NewMovieService() MovieService {
	return &movieServiceImpl{
		movies: make(map[int]*Movie),
	}
}

func (s *movieServiceImpl) GetMovieByID(id int) (*Movie, error) {
	movie, ok := s.movies[id]
	if !ok {
		return nil, fmt.Errorf("movie not found")
	}
	return movie, nil
}

func (s *movieServiceImpl) CreateMovie(movie *Movie) error {
	if _, ok := s.movies[movie.ID]; ok {
		return fmt.Errorf("movie already exists")
	}
	s.movies[movie.ID] = movie
	return nil
}
