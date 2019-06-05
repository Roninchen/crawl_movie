package movie

import (
	"context"
	"crawl_movie/models"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
)

type MovieService struct {

}

func (o *MovieService) GetResult(ctx context.Context,in *MovieRequest)(*MovieResult,error) {
	result :=&MovieResult{}
	maps :=make(map[string]string)
	params := strings.Split(in.Params, "||")
	maps["movie_name"] = params[0]
	movieInfo := models.GetMovie(maps)
	if movieInfo ==nil || len(movieInfo)<1{
		logs.Info("未查到相应数据",*in)
		result = &MovieResult{Code:200,Message:"ok",Data:[]byte("无数据,尝试其他查询策略")}
		return result,nil
	}
	logs.Info("========",movieInfo[0])
	//评分更新
	int, err := strconv.ParseFloat(movieInfo[0].Movie_grade,10)
	logs.Info(err)
	logs.Info("当前评分",int)
	//if int < 1 {
	//	grade := UpdateGrade(&movieInfo[0])
	//	movieInfo[0].Movie_grade = grade.Movie_grade
	//}
	//maps2 := MakeReturn(in.Params, movieInfo[0])
	bytes, err := json.Marshal(movieInfo[0:len(movieInfo)])
	//bytes, err := json.Marshal(maps2)
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

func MakeReturn(params string,movieInfo models.MovieInfo) map[string]string {
	maps := make(map[string]string)
	maps["电影"] = movieInfo.Movie_name
	maps["豆瓣评分"] = movieInfo.Movie_grade

	if strings.Contains(params,"||简介") {
		maps["简介"] = movieInfo.Movie_summary
	}
	if strings.Contains(params, "||演员") {
		maps["演员"] = movieInfo.Movie_main_character
	}
	if strings.Contains(params, "||热评") {
		maps["热评"] = movieInfo.Movie_hot_comment
	}
	if strings.Contains(params, "||导演") {
		maps["导演"] = movieInfo.Movie_director
	}
	if strings.Contains(params, "||上线日期") {
		maps["上线日期"] = movieInfo.Movie_on_time
	}
	return maps
}
func UpdateGrade(v *models.MovieInfo) (*models.MovieInfo) {
	movieId := strconv.FormatInt(v.Movie_id,10)
	url := "https://movie.douban.com/subject/"+movieId +"/"
	logs.Info("url:",url)
	movieInfo, err, _ := models.ReturnMovieInfoByUrl(url)
	if err != nil {
		logs.Info(err)
		//c.Data["json"] = map[string]interface{}{"success": 1, "message": "err","data":err}
		//c.ServeJSON()
	}
	if movieInfo == nil {
		logs.Info("movieInfo is null")
		return v
		//c.Data["json"] = map[string]interface{}{"success": 1, "message": "movieInfo is null","data":err}
		//c.ServeJSON()
	}
	movieInfo.Id = v.Id
	err = models.UpdateMovieGrade(movieInfo)
	if err != nil {
		logs.Info(err)
		return v
		//c.Data["json"] = map[string]interface{}{"success": 1, "message": "err","data":err}
		//c.ServeJSON()
	}
	return movieInfo
}