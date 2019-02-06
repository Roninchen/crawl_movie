package main

import (
	"crawl_movie/movie"
	_ "crawl_movie/routers"
	"github.com/astaxie/beego"
	"google.golang.org/grpc"
	"net"

	//"github.com/astaxie/beego/logs"
)

func main() {

	go grpcgo()
	beego.Run()
	//	log := logs.NewLogger(10000)  // 创建一个日志记录器，参数为缓冲区的大小
	//	     log.SetLogger("console", "")  // 设置日志记录方式：控制台记录
	//	     log.SetLevel(logs.LevelDebug) // 设置日志写入缓冲区的等级：Debug级别（最低级别，所以所有log都会输入到缓冲区）
	//	     log.EnableFuncCallDepth(true) // 输出log时能显示输出文件名和行号（非必须）

	//	     log.Emergency("Emergency")
	//	     log.Alert("Alert")
	//	     log.Critical("Critical")
	//	     log.Error("Error")
	//	     log.Warning("Warning")
	//	     log.Notice("Notice")
	//	     log.Informational("Informational")
	//	     log.Debug("Debug")

	//	     log.Close

}
func grpcgo()  {
		serviceAddress := ":50052"
		movieservice := new(movie.MovieService)

		ls, _ := net.Listen("tcp", serviceAddress)

		gs := grpc.NewServer()

		movie.RegisterMovieServiceServer(gs, movieservice)

		gs.Serve(ls)
}
