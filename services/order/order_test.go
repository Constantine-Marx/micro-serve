package order

import (
	"database/sql"
	"rpcx/utils/storage"
	"sync"
	"testing"
)

func Test_orderServiceImpl_updateSeats(t *testing.T) {
	db, _ := storage.ConnectDB("root", "228809", "localhost:3306", "movie_ticket_service")
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	type fields struct {
		orders     map[int]*Order
		orderMutex sync.Mutex
		db         *sql.DB
	}
	type args struct {
		scheduleID int
		seatRow    int
		seatColumn int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"1", fields{db: db}, args{scheduleID: 1, seatRow: 1, seatColumn: 2}, true},
		{"2", fields{db: db}, args{scheduleID: 1, seatRow: 0, seatColumn: 0}, false},
		{"3", fields{db: db}, args{scheduleID: 1, seatRow: 2, seatColumn: 2}, true}, // 座位坐标超出范围
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &orderServiceImpl{
				orders:     tt.fields.orders,
				orderMutex: tt.fields.orderMutex,
				db:         tt.fields.db,
			}
			if err := s.updateSeats(tt.args.scheduleID, tt.args.seatRow, tt.args.seatColumn); (err != nil) != tt.wantErr {
				t.Errorf("updateSeats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
