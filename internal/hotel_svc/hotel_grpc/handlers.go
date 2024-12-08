package hotel_grpc

import (
	"context"

	"project/internal/hotel_svc"
	pb "project/internal/hotel_svc/hotel_grpc/proto/hotel_grpc"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	pb.UnimplementedHotelServiceServer
	HotelDB *hotel_svc.HotelDB
}

func RegisterHotelServiceServer(s grpc.ServiceRegistrar, srv pb.HotelServiceServer) {
	pb.RegisterHotelServiceServer(s, srv)
}

func (server *GrpcServer) GetRoomsByHotel(_ context.Context, in *pb.GetRoomsByHotelRequest) (*pb.GetRoomsByHotelResponse, error) {
	var response pb.GetRoomsByHotelResponse
	rows, err := server.HotelDB.Db.Query("SELECT id, hotel_id, price, available FROM Rooms WHERE hotel_id = $1 AND available = TRUE", in.HotelId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var room pb.RoomInfo
		if err := rows.Scan(&room.Id, &room.HotelId, &room.Price, &room.Available); err != nil {
			return nil, err
		}
		response.Rooms = append(response.Rooms, &room)
	}

	return &response, nil
}
