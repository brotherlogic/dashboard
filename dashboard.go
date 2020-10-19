package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	fcpb "github.com/brotherlogic/filecopier/proto"
	pbg "github.com/brotherlogic/goserver/proto"
	"github.com/brotherlogic/goserver/utils"
)

//Server main server type
type Server struct {
	*goserver.GoServer
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
	}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {

}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "magic", Value: int64(12345)},
	}
}

func (s *Server) buildDash() {
	// Kick off a refresh
	out, err := exec.Command("sudo", "/etc/init.d/lightdm", "restart").Output()
	time.Sleep(time.Second * 3)
	s.Log(fmt.Sprintf("%v and %v", out, err))
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("dashboard", false, true)
	if err != nil {
		return
	}

	//Upload the file to the remote
	ctx, cancel := utils.ManualContext("db-ul", "db-ul", time.Minute, true)
	conn, err := server.FDialServer(ctx, "filecopier")
	if err != nil {
		server.Log(fmt.Sprintf("Bad dial: %v", err))
		return
	}
	client := fcpb.NewFileCopierServiceClient(conn)

	//Make a local copy of the file
	data, _ := Asset("index.pb")
	ioutil.WriteFile("/tmp/index.html", data, 0644)
	_, err = client.Copy(ctx, &fcpb.CopyRequest{
		InputFile:    "/tmp/index.html",
		InputServer:  "stack1",
		OutputFile:   "/var/www/html/dashboard/index.html",
		OutputServer: "root@www.brotherlogic.com",
	})
	if err != nil {
		log.Fatalf("Bad copy: %v", err)
	}
	cancel()
	conn.Close()
	os.Remove("/tmp/index.html")

	server.buildDash()

	fmt.Printf("%v", server.Serve())
}
