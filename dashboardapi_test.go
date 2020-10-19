package main

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/dashboard/proto"
)

func TestDashboard(t *testing.T) {
	s := InitTest()

	resp, err := s.GetData(context.Background(), &pb.GetDataRequest{})
	if err != nil {
		t.Errorf("Error in getting data: %v", err)
	}

	if len(resp.GetDisplay()) == 0 {
		t.Errorf("Bad response: %v", resp)
	}
}
