package models

import (
	"crawl_movie/conf"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	db orm.Ormer
	gdb *gorm.DB
)

type  MovieInfo struct {
	Id 						int64  `json:"id"`
	Movie_id 				int64  `json:"movie_id"`
	Movie_name 				string `json:"movie_name"`
	Movie_pic 				string `json:"movie_pic"`
	Movie_director 			string `json:"movie_director"`
	Movie_writer 			string `json:"movie_writer"`
	Movie_country 			string `json:"movie_country"`
	Movie_language 			string `json:"movie_language"`
	Movie_main_character 	string `json:"movie_main_character"`
	Movie_type 				string `json:"movie_type"`
	Movie_on_time  			string `json:"movie_on_time"`
	Movie_span 				string `json:"movie_span"`
	Movie_grade 			string `json:"movie_grade"`
	Remark 					string `json:"remark"`
	Movie_summary 			string `json:"movie_summary"`
	Movie_hot_comment 		string `json:"movie_hot_comment"`
	Episode 				string `json:"episode"`
	Season 					string `json:"season"`
	_Create_time 			string `json:"_create_time"`
	_Modify_time 			string `json:"_modify_time"`
	_Status 				int64  `json:"_status"`
}

func init() {
	//初始化配置文件
	conf.Init("")
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	name := viper.GetString("mysql.name")
	password := viper.GetString("mysql.password")
	database := viper.GetString("mysql.database")
	orm.Debug = true // 是否开启调试模式 调试模式下会打印出sql语句
	orm.RegisterDataBase("default","mysql",name+":"+password+"@tcp("+host+":"+port+")/"+database+"?charset=utf8&loc=Local",30)
	orm.RegisterModel(new(MovieInfo))
	db = orm.NewOrm()

	//grom
	gdb,_=gorm.Open("mysql",name+":"+password+"@tcp("+host+":"+port+")/"+database+"?charset=utf8&loc=Local")
	gdb.LogMode(true) // open debug

}

func GetMovie(maps map[string]string) (movieInfo []MovieInfo) {
	//err :=gdb.Model(&MovieInfo{}).Where(maps).Find(&movieInfo).Order("_modify_time DESC", true).Error
	_,err :=db.QueryTable("movie_info").Filter("movie_name__contains",maps["movie_name"]).All(&movieInfo)
	if err != nil {
		glog.Info("数据库查询错误")
		return nil
	}
	return movieInfo
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
//获得电影简介
func GetMovieSummary(movieHtml string) string {
	if movieHtml == ""{
		return ""
	}
	movieHtml = strings.Replace(movieHtml,"\n"," ",-1)
	reg := regexp.MustCompile(`<span.*?property="v:summary"(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	movieSummary :=""
	for _,v := range result{
		movieSummary += v[1] + "\n"
	}
	logs.Info("简介:"+movieSummary)
	movieSummary = strings.Replace(movieSummary,"class=\"\">","",-1)
	movieSummary = strings.Replace(movieSummary,">","",-1)
	return strings.Trim(movieSummary," \n")
}
//获得电影热评
func GetMovieHotComment(movieHtml string) string {
	if movieHtml == ""{
		return ""
	}
	movieHtml = strings.Replace(movieHtml,"\n"," ",-1)
	reg := regexp.MustCompile(`<span.*?class="short">(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	hotComment :=""
	var i = int64(1)
	for _,v := range result{
		if strings.Contains(v[1], "substring") || strings.Contains(v[1],"summary") {
			continue
		}
		hotComment +="热评" + strconv.FormatInt(i,10) +":" + v[1] +" \n"
		i++
	}
	return strings.Trim(hotComment," \n")
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
//获取剧集数
func GetEpisode(movieHtml string)string{
	reg := regexp.MustCompile(`<span.*?class="pl">集数:</span>(.*?)<br/>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	return string(result[0][1])
}

//获取季数
func GetSeason(movieHtml string)string{
	reg := regexp.MustCompile(`<span.*?class="pl">季数:</span>(.*?)<br/>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result)==0 {
		return ""
	}
	return string(result[0][1])
}

func GetMovieUrls(movieHtml string)[]string{
	reg := regexp.MustCompile(`<a.*?href="(https://movie.douban.com/subject/.*?)"`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	var movieSets []string
	for _,v := range result{
		movieSets = append(movieSets, v[1])
	}

	return movieSets
}
func GetMovieGradeIsZero() (movieInfo []MovieInfo,err error) {
	_, err = db.Raw("SELECT * FROM movie_info where movie_grade < ? and movie_on_time > ? and _modify_time < ? ORDER BY _modify_time ASC", 1, "2019-00-00 00:00:00","2019-02-07 00:00:00").QueryRows(&movieInfo)
	//_, err = db.QueryTable("movie_info").Filter("movie_grade < ?", 1).Filter("movie_on_time > ? ", "'2019-00-00 00:00:00'").All(&movieInfo)
	if err != nil {
		logs.Info(err)
		return nil,err
	}
	return movieInfo,nil
}
func UpdateMovieGrade(info *MovieInfo) error {
	_, err := db.Update(info)
	if err != nil {
		logs.Info(err)
		return err
	}
	return nil
}
func Run(sUrl string) {
	//logs.Info("Run！！！")
	//var movieInfo MovieInfo
	//rsp := GetRep(sUrl)
	//if rsp != nil {
	//
	//	logs.Info("begin sleep 20")
	//	//defer rsp.Body.Close()
	//	body, err := ioutil.ReadAll(rsp.Body)
	//	if err != nil {
	//		logs.Info(err)
	//	}
	//	sMovieHtml := string(body)
	//	logs.Info("打印爬取信息")
	//	movieInfo.Movie_name = GetMovieName(sMovieHtml)
	//	//记录电影信息
	//	if movieInfo.Movie_name != "" {
	//		movieInfo.Movie_id = GetMovieId(sMovieHtml)
	//		movieInfo.Movie_name = GetMovieName(sMovieHtml)
	//		movieInfo.Movie_director = GetMovieDirector(sMovieHtml)
	//		movieInfo.Movie_writer = GetMovieWriter(sMovieHtml)
	//		movieInfo.Movie_main_character = GetMovieMainCharacters(sMovieHtml)
	//		movieInfo.Movie_grade = GetMovieGrade(sMovieHtml)
	//		movieInfo.Movie_type = GetMovieGenre(sMovieHtml)
	//		movieInfo.Movie_on_time = GetMovieOnTime(sMovieHtml)
	//		movieInfo.Movie_span = GetMovieRunningTime(sMovieHtml)
	//		movieInfo.Movie_language = GetMovieLanguage(sMovieHtml)
	//		movieInfo.Movie_pic = GetMoviePhoto(sMovieHtml)
	//		movieInfo.Movie_country = GetMovieCountry(sMovieHtml)
	//		movieInfo.Movie_summary = GetMovieSummary(sMovieHtml)
	//		movieInfo.Movie_hot_comment = GetMovieHotComment(sMovieHtml)
	//		movieInfo.Episode = GetEpisode(sMovieHtml)
	//		movieInfo.Season = GetSeason(sMovieHtml)
	//
	//		AddMovie(&movieInfo)
	//	}

		movieInfo, err,sMovieHtml := ReturnMovieInfoByUrl(sUrl)
		if err==nil && movieInfo!=nil {
			AddMovie(movieInfo)
		}
		if sMovieHtml == "" {
			return
		}
		logs.Info("提取该页面的所有连接")
		//如果redis数据小于100万继续提取链接
		queueLen := GetQueueLength()
		logs.Info("queueLen:", queueLen)
		if queueLen <= 1000000 {
			//提取该页面的所有连接
			urls := GetMovieUrls(sMovieHtml)
			for _, url := range urls {
				logs.Info(url)
				if strings.Contains(url, "subject") {
					PutinQueue(url)
				}
			}
		}

		//sUrl 应当记录到 访问set中
		AddToSet(sUrl)
		logs.Info("链接收集完毕")
		time.Sleep(time.Second * 1)
}

func PutUrl(sMovieHtml string,sUrl string)  {
	//提取该页面的所有连接
	urls := GetMovieUrls(sMovieHtml)

	for _,url := range urls{
		logs.Info(url)
		PutinQueue(url)
	}

	//sUrl 应当记录到 访问set中
	AddToSet(sUrl)
	logs.Info("链接收集完毕")
}

func ReturnMovieInfoByUrl(url string) (*MovieInfo,error,string) {
	var movieInfo MovieInfo
	rsp := GetRep(url)
	if rsp != nil {
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			logs.Info(err)
			return nil,err,""
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
			movieInfo.Movie_summary = GetMovieSummary(sMovieHtml)
			movieInfo.Movie_hot_comment = GetMovieHotComment(sMovieHtml)
			movieInfo.Episode = GetEpisode(sMovieHtml)
			movieInfo.Season = GetSeason(sMovieHtml)
			return &movieInfo,nil,sMovieHtml
		}
		return nil,nil,sMovieHtml
	}
	return nil,nil,""
}