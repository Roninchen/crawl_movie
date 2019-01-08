package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"regexp"
	"strings"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"time"
	"strconv"
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
	var movieInfo MovieInfo
	rsp := GetRep(sUrl)
	defer rsp.Body.Close()
	body, _ := ioutil.ReadAll(rsp.Body)
	sMovieHtml := string(body)
	//rsp := httplib.Get(sUrl)
	////设置User-agent以及cookie是为了防止  豆瓣网的 403
	//rsp.Header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:50.0) Gecko/20100101 Firefox/50.0")
	//rsp.Header("Cookie", `bid=gFP9qSgGTfA; __utma=30149280.1124851270.1482153600.1483055851.1483064193.8; __utmz=30149280.1482971588.4.2.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; ll="118221"; _pk_ref.100001.4cf6=%5B%22%22%2C%22%22%2C1483064193%2C%22https%3A%2F%2Fwww.douban.com%2F%22%5D; _pk_id.100001.4cf6=5afcf5e5496eab22.1482413017.7.1483066280.1483057909.; __utma=223695111.1636117731.1482413017.1483055857.1483064193.7; __utmz=223695111.1483055857.6.5.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; _vwo_uuid_v2=BDC2DBEDF8958EC838F9D9394CC5D9A0|2cc6ef7952be8c2d5408cb7c8cce2684; ap=1; viewed="1006073"; gr_user_id=e5c932fc-2af6-4861-8a4f-5d696f34570b; __utmc=30149280; __utmc=223695111; _pk_ses.100001.4cf6=*; __utmb=30149280.0.10.1483064193; __utmb=223695111.0.10.1483064193`)
	//sMovieHtml,err := rsp.String()
	//if err != nil{
	//	//panic(err)
	//	fmt.Println(err)
	//	return
	//}

	movieInfo.Movie_name            = GetMovieName(sMovieHtml)
	//记录电影信息
	if movieInfo.Movie_name !="" {
		movieInfo.Movie_id					=  GetMovieId(sMovieHtml)
		movieInfo.Movie_name 				=  GetMovieName(sMovieHtml)
		movieInfo.Movie_director 			=  GetMovieDirector(sMovieHtml)
		movieInfo.Movie_writer				=  GetMovieWriter(sMovieHtml)
		movieInfo.Movie_main_character 		=  GetMovieMainCharacters(sMovieHtml)
		movieInfo.Movie_grade 				=  GetMovieGrade(sMovieHtml)
		movieInfo.Movie_type 				=  GetMovieGenre(sMovieHtml)
		movieInfo.Movie_on_time 			=  GetMovieOnTime(sMovieHtml)
		movieInfo.Movie_span 				=  GetMovieRunningTime(sMovieHtml)
		movieInfo.Movie_language 			=  GetMovieLanguage(sMovieHtml)
		movieInfo.Movie_pic  				=  GetMoviePhoto(sMovieHtml)
		movieInfo.Movie_country				=  GetMovieCountry(sMovieHtml)



		AddMovie(&movieInfo)
	}

	//提取该页面的所有连接
	urls := GetMovieUrls(sMovieHtml)

	for _,url := range urls{
		PutinQueue(url)
		//c.Ctx.WriteString("<br>" + url + "</br>")
	}

	//sUrl 应当记录到 访问set中
	AddToSet(sUrl)

	time.Sleep(time.Second)
}

