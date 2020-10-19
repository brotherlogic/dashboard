package main

import (
	"context"

	pb "github.com/brotherlogic/dashboard/proto"
)

//GetData for the dashboard display
func (s *Server) GetData(ctx context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	return &pb.GetDataResponse{Display: "api response"}, nil
}
