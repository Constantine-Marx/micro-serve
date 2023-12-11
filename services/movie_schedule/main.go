// Package movie_schedule services/movie_schedule/main.go
package movie_schedule

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
)

type MovieSchedule struct {
	ID           int    `json:"id"`
	MovieID      int    `json:"movie_id"`
	CinemaName   string `json:"cinema_name"`
	City         string `json:"city"`
	ScheduleDate string `json:"schedule_date"`
	ScheduleTime string `json:"schedule_time"`
}

type SeatsSchedule struct {
	ID int `json:"id"`
}

type MovieScheduleService interface {
	AddMovieSchedule(ctx context.Context, schedule *MovieSchedule) error
	GetMovieSchedule(ctx context.Context, schedule *MovieSchedule, reply *[]MovieSchedule) error
	GetMovieSeats(ctx context.Context, args *SeatsSchedule, reply *[][]int) error
}

type movieScheduleServiceImpl struct {
	db *sql.DB
}

func (s *movieScheduleServiceImpl) AddMovieSchedule(ctx context.Context, schedule *MovieSchedule) error {
	// 初始化座位
	seats := [10][10]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	// 插入电影场次到数据库
	result, err := s.db.Exec("INSERT INTO movie_schedules (movie_id, cinema_name, city, schedule_date, schedule_time,seats) VALUES (?, ?,?, ?, ?, ?)", schedule.MovieID, schedule.CinemaName, schedule.City, schedule.ScheduleDate, schedule.ScheduleTime, seats)
	if err != nil {
		return err
	}

	// 获取自动生成的 ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	schedule.ID = int(lastInsertID)

	return nil
}

func (s *movieScheduleServiceImpl) GetMovieSchedule(ctx context.Context, schedule *MovieSchedule, reply *[]MovieSchedule) error {
	rows, err := s.db.Query("SELECT id, cinema_name, city, schedule_date, schedule_time FROM movie_schedules WHERE movie_id = ?", schedule.MovieID)
	if err != nil {
		log.Println(err)
		return err
	}
	if rows == nil {
		log.Println("rows is nil")
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var retSchedules []MovieSchedule
	for rows.Next() {
		var schedules MovieSchedule
		err := rows.Scan(&schedules.ID, &schedules.CinemaName, &schedules.City, &schedules.ScheduleDate, &schedules.ScheduleTime)
		if err != nil {
			log.Println(err)
			return err
		}

		//schedules.ScheduleDate, err = time.Parse("2006-01-02", scheduleDate)
		//if err != nil {
		//	log.Printf("Error parsing schedule date: %v", err)
		//	return err
		//}
		//
		//schedules.ScheduleTime, err = time.Parse("15:04:05", scheduleTime)
		//if err != nil {
		//	log.Printf("Error parsing schedule time: %v", err)
		//	return err
		//}
		//schedules.ScheduleDate = schedules.ScheduleDate.UTC()
		//schedules.ScheduleTime = schedules.ScheduleTime.UTC()
		retSchedules = append(retSchedules, schedules)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return err
	}
	*reply = retSchedules
	return nil
}

func (s *movieScheduleServiceImpl) GetMovieSeats(ctx context.Context, args *SeatsSchedule, reply *[][]int) error {
	rows, err := s.db.Query("SELECT seats FROM movie_schedules WHERE id = ?", args.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	if rows == nil {
		log.Println("rows is nil")
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)
	var retSchedules [][]int
	for rows.Next() {
		var seats string
		err := rows.Scan(&seats)
		if err != nil {
			log.Println(err)
			return err
		}
		var seat [][]int
		err = json.Unmarshal([]byte(seats), &seat)
		if err != nil {
			log.Println(err)
			return err
		}
		retSchedules = append(retSchedules, seat...)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return err
	}

	*reply = retSchedules
	return nil
}

func NewMovieScheduleService(db *sql.DB) MovieScheduleService {
	return &movieScheduleServiceImpl{
		db: db,
	}
}
