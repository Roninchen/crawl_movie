package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	db orm.Ormer
)

type  MovieInfo struct {
	Id int64
	Movie_id int64
	Movie_name string
	Movie_pic string
	Movie_director string
	Movie_writer string
	Movie_country string
	Movie_language string
	Movie_main_character string
	Movie_type string
	Movie_on_time  string
	Movie_span string
	Movie_grade string
	Remark string
	_Create_time string
	_Modify_time string
	_Status int64
}

func init() {
	orm.Debug = true // 是否开启调试模式 调试模式下会打印出sql语句
	orm.RegisterDataBase("default","mysql","fuck_hacker:690383f40a51d07c@tcp(111.231.84.238:3306)/movie?charset=utf8&loc=Local",30)
	orm.RegisterModel(new(MovieInfo))
	db = orm.NewOrm()
}

func AddMovie(movie_info *MovieInfo)(int64,error){

		movie :=new(MovieInfo)
		movie.Movie_id = movie_info.Movie_id
		err1 := db.Read(movie)
		if err1 != nil {
			logs.Error(0,err1)
		}
		if movie.Movie_name!="" {
			return 0,err1
		}

		movie_info.Id = 0
		id,err := db.Insert(movie_info)
		logs.Error(id,err)
		logs.Info(movie.Movie_name+" movie insert success!")
		return id,err
}

//获得电影名字
func GetMovieName(movieHtml string)string{
	if movieHtml == ""{
		return ""
	}
	reg := regexp.MustCompile(`<span\s*property="v:itemreviewed">(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	return string(result[0][1])
}
//获得电影导演
func GetMovieDirector(movieHtml string) string{
	if movieHtml == ""{
		return ""
	}
	reg := regexp.MustCompile(`<a.*?rel="v:directedBy">(.*?)</a>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	return string(result[0][1])
}
//电影编剧
func GetMovieWriter(movieHtml string) string {
	reg := regexp.MustCompile(`"<a.*?href="/celebrity/*?/">(.*?)</a>"`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	movieWriters := ""
	for _,v := range result{
		movieWriters += v[1] + "/"
	}
	return strings.Trim(movieWriters,"/")
}
//电影主演
func GetMovieMainCharacters(movieHtml string)string{
	reg := regexp.MustCompile(`<a.*?rel="v:starring">(.*?)</a>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	mainCharacters := ""
	for _,v := range result{
		mainCharacters += v[1] + "/"
	}
	return strings.Trim(mainCharacters,"/")
}
//豆瓣评分
func GetMovieGrade(movieHtml string)string{
	reg := regexp.MustCompile(`<strong.*?property="v:average">(.*?)</strong>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	return string(result[0][1])
}
//电影类型
func GetMovieGenre(movieHtml string)string{
	reg := regexp.MustCompile(`<span.*?property="v:genre">(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	movieGenre := ""
	for _,v := range result{
		movieGenre += v[1] + "/"
	}
	return strings.Trim(movieGenre,"/")
}
//<span class="pl">制片国家/地区:</span> 中国大陆 / 香港<br/>
func GetMovieCountry(movieHtml string) string {
	reg := regexp.MustCompile(`<span.*?class="pl">制片国家/地区:</span>(.*?)<br/>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	movieCountry := ""
	for _,v := range result{
		movieCountry += v[1] + "/"
	}
	return strings.Trim(movieCountry,"/")
}
//电影语言
//<span class="pl">语言:</span> 英语<br/>
func GetMovieLanguage(movieHtml string) string {
	reg := regexp.MustCompile(`<span.*?class="pl">语言:</span>(.*?)<br/>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	movieLanguage := ""
	for _,v := range result{
		movieLanguage += v[1] + "/"
	}
	return strings.Trim(movieLanguage,"/")
}
//电影上线时间
func GetMovieOnTime(movieHtml string) string{
	reg := regexp.MustCompile(`<span.*?property="v:initialReleaseDate".*?>(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	movieOnTime :=""
	for _,v := range result{
		movieOnTime += v[1] + "/"
	}
	return strings.Trim(movieOnTime,"/")
}
//电影时长
func GetMovieRunningTime(movieHtml string) string{
	reg := regexp.MustCompile(`<span.*?property="v:runtime".*?>(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	return string(result[0][1])
}
//电影海报
//<img class="media" src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2522069454.webp" />
func GetMoviePhoto(movieHtml string) string {
	reg := regexp.MustCompile(`<img.*?class="media".*?src="(.*?)".*?/>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	return strings.Trim(string(result[0][1]),"")
}
//电影id
//share-id="26416062"
func GetMovieId(movieHtml string) int64 {
	reg := regexp.MustCompile(`share-id="(.*?)"`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return 0
	}
	s,_ :=strconv.Atoi(result[0][1])
	return int64(s)
}


func GetMovieUrls(movieHtml string)[]string{
	reg := regexp.MustCompile(`<a.*?href="(https://movie.douban.com/.*?)"`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	var movieSets []string
	for _,v := range result{
		movieSets = append(movieSets, v[1])
	}

	return movieSets
}

func Run(sUrl string)  {
	logs.Info("Run！！！")
	var movieInfo MovieInfo
	rsp := GetRep(sUrl)
	if rsp !=nil {

		logs.Info("begin sleep 20")
		//defer rsp.Body.Close()
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			logs.Info(err)
		}
		sMovieHtml := string(body)
		logs.Info("打印爬取信息")
		movieInfo.Movie_name = GetMovieName(sMovieHtml)
		//记录电影信息
		if movieInfo.Movie_name != "" {
			movieInfo.Movie_id = GetMovieId(sMovieHtml)
			movieInfo.Movie_name = GetMovieName(sMovieHtml)
			movieInfo.Movie_director = GetMovieDirector(sMovieHtml)
			movieInfo.Movie_writer = GetMovieWriter(sMovieHtml)
			movieInfo.Movie_main_character = GetMovieMainCharacters(sMovieHtml)
			movieInfo.Movie_grade = GetMovieGrade(sMovieHtml)
			movieInfo.Movie_type = GetMovieGenre(sMovieHtml)
			movieInfo.Movie_on_time = GetMovieOnTime(sMovieHtml)
			movieInfo.Movie_span = GetMovieRunningTime(sMovieHtml)
			movieInfo.Movie_language = GetMovieLanguage(sMovieHtml)
			movieInfo.Movie_pic = GetMoviePhoto(sMovieHtml)
			movieInfo.Movie_country = GetMovieCountry(sMovieHtml)

			AddMovie(&movieInfo)
		}
		logs.Info("提取该页面的所有连接")
		//提取该页面的所有连接
		urls := GetMovieUrls(sMovieHtml)

		for _,url := range urls{
			logs.Info(url)
			if strings.Contains(url, "subject") {
				PutinQueue(url)
			}
			//PutinQueue(url)
			//c.Ctx.WriteString("<br>" + url + "</br>")
		}

		//sUrl 应当记录到 访问set中
		AddToSet(sUrl)
		logs.Info("链接收集完毕")
		time.Sleep(time.Second*1)
	}
}

func PutUrl(sMovieHtml string,sUrl string)  {
	//提取该页面的所有连接
	urls := GetMovieUrls(sMovieHtml)

	for _,url := range urls{
		logs.Info(url)
		PutinQueue(url)
		//c.Ctx.WriteString("<br>" + url + "</br>")
	}

	//sUrl 应当记录到 访问set中
	AddToSet(sUrl)
	logs.Info("链接收集完毕")
}

