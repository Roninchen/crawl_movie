package movie

import (
	"context"
	"crawl_movie/models"
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

type MovieService struct {

}

func (o *MovieService) GetResult(ctx context.Context,in *MovieRequest)(*MovieResult,error) {
	result :=&MovieResult{}
	maps :=make(map[string]string)
	maps["movie_name"] = in.Params
	movieInfo := models.GetMovie(maps)
	if movieInfo ==nil || len(movieInfo)<1{
		logs.Info("未查到相应数据",*in)
		result = &MovieResult{Code:200,Message:"ok",Data:[]byte("无数据,尝试其他查询策略")}
		return result,nil
	}
	logs.Info("========",movieInfo[0])

	maps2 :=make(map[string]string)
	maps2["电影"] = movieInfo[0].Movie_name
	maps2["豆瓣评分"] = movieInfo[0].Movie_grade
	//bytes, err := json.Marshal(movieInfo[0])
	bytes, err := json.Marshal(maps2)
	if err != nil {
		logs.Info(err)
		result = &MovieResult{Code:200,Message:"ok",Data:[]byte("无数据,尝试其他查询策略")}
		return result,nil
	}
	result.Data = bytes
	result.Code =200
	result.Message = "ok"
	logs.Info("ok",bytes)
	return result,nil
}
