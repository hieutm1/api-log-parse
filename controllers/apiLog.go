package controllers

import(
        "github.com/gin-gonic/gin"
		"github.com/go-redis/redis/v8"
		database "logsquery/database"
		"code.cloudfoundry.org/bytefmt"
        "log"
		_"database/sql"
		
        "time"
        // "net"
        // "bytes"
        "encoding/json"
        // "strings"
        // "crypto/tls"
        "fmt"
		"context"
        "net/http"
        // "reflect"
        "strconv" 
		"database/sql"
		"math"
		// _ "github.com/go-sql-driver/mysql"
)

type HOST struct  {
	Http_Host []string `json:"http_host"`
}

type DomainHost struct  {
	Http_Host string `json:"http_host"`
	User string `json:"user"`
}

type Search struct {
	Hostname string `json:"hostname"`
	FromTime int `json:"fromtime"`
	ToTime int `json:"totime"`
}
type Detail struct {
	Http_Host string `json:"http_host"`
	Time string `json:"time"`
	Byte int `json:"byte"`
	Requests int `json:"requests"`
	Bytevn int `json:"byteVN"`
	Requestvn int `json:"requestVN"`
	Cache_hit int `json:"cache_hit"`
	Cache_miss int `json:"cache_miss"`
	Status_2x int `json:"status_2x"`
	Status_3x int `json:"status_3x"`
	Status_4x int `json:"status_4x"`
	Status_5x int `json:"status_5x"`
	Macintosh int `json:"macintosh"`
	Windows int `json:"windows"`
	Ios int `json:"ios"`
	Android int `json:"android"`
	Smartv int `json:"smartv"`
	Other int `json:"other"`
}

type TotalBytes struct {
	Byte string `json:"byte"`
	Time string `json:"time"`
}

type Bandwidth struct {
	Speed float64 `json:"speed"`
	Time string `json:"time"`
	Unit string `json:"unit"`
}

type BandwidthAll struct {
	Total float64 `json:"total"`
	VN float64 `json:"vn"`
	INT float64 `json:"int"`
	Time string `json:"time"`
	Unit string `json:"unit"`
}

type ByteRequestDetail struct {
	TotalByte float64 `json:"totalByte"`
	VnByte float64 `json:"vnByte"`
	INTByte float64 `json:"intByte"`
	TotalRequest float64 `json:"totalRequest"`
	VnRequest float64 `json:"vnRequest"`
	INTRequest float64 `json:"intRequest"`
	Time string `json:"time"`
	Unit string `json:"unit"`
}

type CountryISO struct {
	Http_Host string `json:"http_host"`
	Time string `json:"time"`
	Country string `json:"country"`
	Count string `json:"count"`
}

type TimeResponse struct {
	Http_Host string `json:"http_host"`
	Time string `json:"time"`
	Under4s string `json:"under4s"`
	Over4s string `json:"over4s"`
	Over6s string `json:"over6s"`
}

type ISPs struct {
	Http_Host string `json:"http_host"`
	Time string `json:"time"`
	Isp string `json:"isp"`
	Count string `json:"count"`
}

type City struct {
	Http_Host string `json:"http_host"`
	Time string `json:"time"`
	City string `json:"city"`
	Count string `json:"count"`
}

type Referer struct {
	Http_Host string `json:"http_host"`
	Time string `json:"time"`
	Referer string `json:"referer"`
	Count string `json:"count"`
}

type Ccu struct {
	Http_Host string `json:"http_host"`
	Time string `json:"time"`
	Ccu string `json:"ccu"`
}

type Avgspeeds struct {
	Time string `json:"time"`
	Speed string `json:"speed"`
}

type AvgspeedDetail struct {
	Time string `json:"time"`
	TotalSpeed float64 `json:"totalSpeed"`
	VnSpeed float64 `json:"vnSpeed"`
	INTSpeed float64 `json:"intSpeed"`
	Unit string `json:"unit"`
}

type TopHttpCodeError struct {
	Http_Host string `json:"http_host"`
	Time string `json:"time"`
	HttpCodeError string `json:"http_code_error"`
	Count string `json:"count"`
}

type TopHttpCodeErrorDomain struct {
	Http_Host string `json:"http_host"`
	Time string 	`json:"time"`
	Code403 string `json:"code403"`
	Code404 string `json:"code404"`
	Code503 string `json:"code503"`
	Code505 string `json:"code505"`
}
// 
var table string
var db *sql.DB
var rdb *redis.Client



// @BasePath /api/v1

// log godoc
// @Summary all domain
// @Schemes
// @Description get all domain
// @Tags domain
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: "Basic Authen"
// @success 200 {object} []controllers.HOST
// @Router /log/domains [get]
func Domains(c * gin.Context){

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "http_host").Err()

    val, err := rdb.Get(ctx, "http_host").Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT http_host FROM http_host ORDER BY http_host")
		if err != nil {
			fmt.Println(err.Error())
		}

		defer db.Close()

		log.Printf("Successfully connected to database")
		var host string
		var http_host HOST
		for rows.Next(){
			
			err = rows.Scan(&host)
			if err != nil {
				panic(err.Error())
			}
			http_host.Http_Host = append(http_host.Http_Host , host)

		}
		jshttp , err := json.Marshal(http_host)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(jshttp)
		errs := rdb.Set(ctx, "http_host", jshttp , 1000000000 * 60).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, http_host)
    }else{
		c.String(200, val)
		fmt.Println("http_host cache")
	}
	
}
// @BasePath /api/v1

// log godoc
// @Summary domain
// @Schemes
// @Description do domain
// @Tags domain
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param domain path string true "vovlive.1cdn.vn"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/domain/{domain}/{fromtime}/{totime} [get]
func DomainDetail(c * gin.Context){

	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, hostname+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, hostname+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT time,http_host,byte,requests,byteVN,requestVN,cache_hit,Cache_miss,status_2x,status_3x,status_4x,status_5x,macintosh,windows,ios,android,smartv,other FROM  access_"+s+ " WHERE http_host='"+hostname+"' AND  UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")

		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time,&detail.Http_Host, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, hostname+fromtime+totime , json , 1000000000 * 60).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("domain detail hit cache")
	}
	
}
// @BasePath /api/v1

// log godoc
// @Summary daily
// @Schemes
// @Description do daily
// @Tags domain
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param domain path string true "vovlive.1cdn.vn"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/daily/{domain}/{fromtime}/{totime} [get]
func DailyDetail(c * gin.Context){

	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "daily"+hostname+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "daily"+hostname+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT date(time),http_host,byte,requests,byteVN,requestVN,cache_hit,Cache_miss,status_2x,status_3x,status_4x,status_5x,macintosh,windows,ios,android,smartv,other FROM  daily_"+s+ " WHERE http_host='"+hostname+"' AND  date(time) >=date(FROM_UNIXTIME('"+fromtime+"')) AND  date(time) <=date(FROM_UNIXTIME('"+totime+"')) ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time,&detail.Http_Host, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "daily"+hostname+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("domain detail daily hit cache")
	}
	
}
// @BasePath /api/v1

// log godoc
// @Summary monthly
// @Schemes
// @Description do monthly
// @Tags domain
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param domain path string true "vovlive.1cdn.vn"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/monthly/{domain}/{fromtime}/{totime} [get]
func MonthlyDetail(c * gin.Context){

	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	year 	:= tm.Year()
	// month 	:= int(tm.Month())
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "daily"+hostname+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "monthly"+hostname+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT time,http_host,byte,requests,byteVN,requestVN,cache_hit,Cache_miss,status_2x,status_3x,status_4x,status_5x,macintosh,windows,ios,android,smartv,other FROM  monly_"+strconv.Itoa(year)+ " WHERE http_host='"+hostname+"' AND  date(time) >=date(FROM_UNIXTIME('"+fromtime+"')) AND  date(time) <=date(FROM_UNIXTIME('"+totime+"')) ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time,&detail.Http_Host, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "monthly"+hostname+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("daily monthly hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary sumary today
// @Schemes
// @Description do sumary
// @Tags domain
// @Accept json
// @Produce json 
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param domain path string true "vovlive.1cdn.vn"
// @Router /log/today/{domain} [get]
func TodaySum(c * gin.Context){

	hostname := c.Param("hostname") 

	current_time := time.Now()
	
	today := fmt.Sprintf("%d-%02d-%02d", current_time.Year(), current_time.Month(), int(current_time.Day()))

    fmt.Println(current_time)
	s := fmt.Sprintf("%d%02d", current_time.Year(), int(current_time.Month()))
	
	// day 	:= tm.Day()
	
	ctx := context.Background()
	rdb = database.RDb() 
	// _= rdb.Del(ctx, "today"+hostname).Err()

    val, err := rdb.Get(ctx, "today"+hostname).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT time,http_host,sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE http_host='"+hostname+"' AND  date(time)='"+today+"'" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// fmt.Println("SELECT time,sum(http_host),sum(byte),sum(requests),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE http_host='"+hostname+"' AND  date(time)='"+today+"'" )
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time,&detail.Http_Host, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "today"+hostname , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("today detail hit cache")
	}
	
}
// @BasePath /api/v1

// log godoc
// @Summary sumary this month
// @Schemes
// @Description do sumary
// @Tags domain
// @Accept json
// @Produce json
// @success 200 {object} []controllers.Detail
// @Param domain path string true "vovlive.1cdn.vn"
// @Router /log/thismonth/{domain} [get]
func ThisMonth(c * gin.Context){

	hostname := c.Param("hostname") 

	current_time := time.Now()
	
	thismonth := fmt.Sprintf("%d-%02d", current_time.Year(), int(current_time.Month()))

    fmt.Println(current_time)
	s := fmt.Sprintf("%d%02d", current_time.Year(), int(current_time.Month()))
	
	// day 	:= tm.Day()
	
	ctx := context.Background()
	rdb = database.RDb() 
	// _= rdb.Del(ctx, "today"+hostname).Err()

    val, err := rdb.Get(ctx, "thismonth"+hostname).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT time,http_host,sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE http_host='"+hostname+"' AND  date(time) LIKE  '"+thismonth+"%'" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		fmt.Println("SELECT time,http_host,sum(byte),sum(requests),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE http_host='"+hostname+"' AND  date(time) LIKE  '"+thismonth+"%'" )
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time,&detail.Http_Host, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "thismonth"+hostname , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(http.StatusOK, details)
    }else{
		c.String(http.StatusOK, val)
		fmt.Println("thismonth detail hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary byte to KB,MB,GB,TB
// @Schemes
// @Description byte to KB,MB,GB,TB
// @Tags domain
// @Accept json
// @Produce json
// @success 200 {string} 1K
// @Param byte path int true "1024"
// @Router /log/bytefmt/{byte} [get]
func Byteto(c * gin.Context){
	byte := c.Param("byte") 
	id, err := strconv.ParseInt(byte, 10, 64)
	if err != nil {
		fmt.Println("error converting string to int")
		c.String(401, "error converting string to int")
	}else{
		to := bytefmt.ByteSize(uint64(id))+"B"
		c.String(http.StatusOK, to)
	}

}

// @BasePath /api/v1

// log godoc
// @Summary sumary total byte
// @Schemes
// @Description do sumary total byte
// @Tags log
// @Accept json
// @Produce json
// @success 200 {object} []controllers.Bandwidth
// @Router /log/totalbyte [get]
func TotalByte(c * gin.Context){

	current_time := time.Now()
	
	// thismonth := strconv.Itoa(current_time.Year())+"-"+strconv.Itoa(int(current_time.Month()))

    fmt.Println(current_time)
	s := fmt.Sprintf("%d%02d", current_time.Year(), int(current_time.Month()))
	
	// day 	:= tm.Day()
	
	ctx := context.Background()
	rdb = database.RDb() 
	// _= rdb.Del(ctx, "bytetotal" ).Err()

    val, err := rdb.Get(ctx, "bytetotal" ).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT sum(byte),time FROM  access_"+s+ " GROUP BY time DESC LIMIT 20 " )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		bws := make([]*TotalBytes,0)
		for rows.Next(){
			
			bw := new(TotalBytes)
			if err := rows.Scan( &bw.Byte,&bw.Time ); err != nil {
				fmt.Println(err)
			}       
			bws = append(bws, bw)

		}

		json , err := json.Marshal(bws)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "bytetotal" , json , 1000000000 * 250).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(http.StatusOK, bws)
    }else{
		c.String(http.StatusOK, val)
		fmt.Println("TotalByte  hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary sumary bandwidth of domain 
// @Schemes
// @Description do sumary bandwidth of domain 
// @Tags domain
// @Accept json
// @Produce json
// @success 200 {object} []controllers.Bandwidth
// @Param domain path string true "vovlive.1cdn.vn"
// @Param unit path string true "Kb,Mb,Gb,Tb"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/bandwidth/{domain}/{unit}/{fromtime}/{totime} [get]
func BandwidthDomain(c * gin.Context){

	current_time := time.Now()

	var tunix int64

	unit := c.Param("unit")
	
	domain := c.Param("domain")

	fromtime := c.Param("fromtime")

	totime	 := c.Param("totime")
	// now
	if fromtime == "now" {
		tunix = (current_time.Unix() - 1800)
		fromtime = strconv.FormatInt(tunix, 10)
		totime   = strconv.FormatInt(current_time.Unix(), 10)
	}

	var u float64

	if unit == "Kb" {
		u = 1000
	}
	if unit == "Mb" {
		u = 1000*1000
	}
	if unit == "Gb" {
		u = 1000*1000*1000
	}
	if unit == "Tb" {
		u = 1000*1000*1000*1000
	}
    
	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()
	
	ctx := context.Background()
	rdb = database.RDb() 
	_= rdb.Del(ctx, unit+"bwtotal"+domain+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, unit+"bwtotal"+domain+fromtime+totime ).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection	

		rows, err := db.Query("SELECT sum(byte),time FROM  access_"+s+ " WHERE http_host='"+domain+"' AND time >= FROM_UNIXTIME('"+fromtime+"') AND  time <= FROM_UNIXTIME('"+totime+"')   GROUP BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		bws := make([]*Bandwidth,0)
		for rows.Next(){
			
			bw := new(Bandwidth)
			if err := rows.Scan( &bw.Speed , &bw.Time ); err != nil {
				fmt.Println(err)
			}       
			bw.Speed =  RoundDown(bw.Speed*8/u/300,3)
			bw.Unit = unit+"ps"
 			bws = append(bws, bw)

		}

		json , err := json.Marshal(bws)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, unit+"bwtotal"+domain+fromtime+totime , json , 1000000000 * 150).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(http.StatusOK, bws)
    }else{
		c.String(http.StatusOK, val)
		fmt.Println("BandwidthTotal  hit cache")
	}
	
}

func RoundDown(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = round / pow
	return
}

// @BasePath /api/v1

// log godoc
// @Summary country
// @Schemes
// @Description count country
// @Tags domain
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.CountryISO
// @Param domain path string true "vovlive.1cdn.vn"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/contry/{domain}/{fromtime}/{totime} [get]
func Country(c * gin.Context){

	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "contry"+hostname+fromtime+totime ).Err()
	
    val, err := rdb.Get(ctx, "contry"+hostname+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		query := "SELECT date(time),http_host,name,sum(data) FROM  %s WHERE http_host='"+hostname+"'  AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC"
		if checkTime(totime,fromtime) == false { 
			table  = "geoip_country_"+s
		}else {
			table  = "daily_geoip_country_"+s
		}
		queryinsrt := fmt.Sprintf(query , table)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		details := make([]*CountryISO,0)
		for rows.Next(){
			
			detail := new(CountryISO)
			if err := rows.Scan( &detail.Time,&detail.Http_Host,  &detail.Country, &detail.Count ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "contry"+hostname+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println(" Country hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary time Respone
// @Schemes
// @Description count time Respone
// @Tags domain
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.TimeResponse
// @Param domain path string true "vovlive.1cdn.vn"
// @Param name path string true "under4sec,under6sec,over6sec"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/responseTimeDaily/{domain}/{name}/{fromtime}/{totime} [get]
func ResponseTimeDaily(c * gin.Context){

	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 
	 

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "ResponseTimeDaily"+hostname+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ResponseTimeDaily"+hostname+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT time,http_host,sum(under4s),sum(over4s),sum(over6s) FROM  time_response_"+s+ " WHERE http_host='"+hostname+"'  AND date(time) >=date(FROM_UNIXTIME('"+fromtime+"')) AND date(time) <=date(FROM_UNIXTIME('"+totime+"')) GROUP BY time ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")

		details := make([]*TimeResponse,0)
		for rows.Next(){
			
			detail := new(TimeResponse)
			if err := rows.Scan( &detail.Time,&detail.Http_Host,  &detail.Under4s, &detail.Over4s, &detail.Over6s ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "ResponseTimeDaily"+hostname+fromtime+totime , json , 1000000000 * 150).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("responseTimeDaily hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary time Respone
// @Schemes
// @Description count time Respone
// @Tags domain 
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.TimeResponse
// @Param domain path string true "vovlive.1cdn.vn"
// @Param name path string true "under4sec,under6sec,over6sec"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/responseTime/{domain}/{name}/{fromtime}/{totime} [get]
func ResponseTime(c * gin.Context){

	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 
	 

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()

	// _= rdb.Del(ctx, "ResponseTime"+hostname+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ResponseTime"+hostname+fromtime+totime ).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		query := "SELECT time,http_host,sum(under4s),sum(over4s),sum(over6s) FROM  time_response_"+s+ " WHERE http_host='"+hostname+"'  AND date(time) >=date(FROM_UNIXTIME('"+fromtime+"')) AND date(time) <=date(FROM_UNIXTIME('"+totime+"')) GROUP BY time ORDER BY time ASC"
	
		rows, err := db.Query(query)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
	
		details := make([]*TimeResponse,0)
		for rows.Next(){
			
			detail := new(TimeResponse)
			if err := rows.Scan(  &detail.Time,&detail.Http_Host,  &detail.Under4s, &detail.Over4s, &detail.Over6s  ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "ResponseTime"+hostname+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("daily reponse time hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary ISP
// @Schemes
// @Description count ISP
// @Tags domain
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.ISPs
// @Param domain path string true "vovlive.1cdn.vn"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/isp/{domain}/{fromtime}/{totime} [get]
func ISP(c * gin.Context){

	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "ISP"+hostname+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ISP"+hostname+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		query :="SELECT date(time),http_host,name,sum(data) FROM  %s WHERE http_host='"+hostname+"' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC" 
		if  checkTime(totime,fromtime) == false {
			table  = "geoip_isp_"+s
		}else {
			table  = "daily_geoip_isp_"+s
		}
		queryinsrt := fmt.Sprintf(query,table)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		details := make([]*ISPs,0)
		for rows.Next(){
			
			detail := new(ISPs)
			if err := rows.Scan( &detail.Time,&detail.Http_Host,  &detail.Isp, &detail.Count ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "ISP"+hostname+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("daily ISP hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary city domain
// @Schemes
// @Description count city domain
// @Tags domain
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.City
// @Param domain path string true "vovlive.1cdn.vn"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/city/{domain}/{fromtime}/{totime} [get]
func DailyCity(c * gin.Context){

	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "City"+hostname+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "City"+hostname+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		query := "SELECT date(time),http_host,name,sum(data) FROM  %s WHERE http_host='"+hostname+"' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC"
		if checkTime(totime,fromtime) == false {
			table  = "geoip_city_"+s
		}else {
			table  = "daily_geoip_city_"+s
		}
		queryinsrt := fmt.Sprintf(query,table)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		details := make([]*City,0)
		for rows.Next(){
			
			detail := new(City)
			if err := rows.Scan( &detail.Time,&detail.Http_Host,  &detail.City, &detail.Count ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "City"+hostname+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("City hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary ccu domain
// @Schemes
// @Description count ccu domain
// @Tags domain
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Ccu
// @Param domain path string true "vovlive.1cdn.vn"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/ccu/{domain}/{fromtime}/{totime} [get]
func CCUdomain(c * gin.Context){

	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "City"+hostname+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ccu"+hostname+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT time,http_host,ROUND(sum(requests)/300) FROM  access_"+s+ " WHERE http_host='"+hostname+"' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY time ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		details := make([]*Ccu,0)
		for rows.Next(){
			
			detail := new(Ccu)
			if err := rows.Scan( &detail.Time,&detail.Http_Host,  &detail.Ccu, ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "ccu"+hostname+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("ccu hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary sumary bandwidth of domain all
// @Schemes
// @Description do sumary bandwidth of domain all
// @Tags domain
// @Accept json
// @Produce json
// @success 200 {object} []controllers.Bandwidth
// @Param domain path string true "vovlive.1cdn.vn"
// @Param unit path string true "Kb,Mb,Gb,Tb"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/bandwidthAll/{domain}/{unit}/{fromtime}/{totime} [get]
func BandwidthDomainAll(c * gin.Context){

	current_time := time.Now()

	var tunix int64

	unit := c.Param("unit")
	
	domain := c.Param("domain")

	fromtime := c.Param("fromtime")

	totime	 := c.Param("totime")
	// now
	if fromtime == "now" {
		tunix = (current_time.Unix() - 1800)
		fromtime = strconv.FormatInt(tunix, 10)
		totime   = strconv.FormatInt(current_time.Unix(), 10)
	}

	var u float64

	if unit == "Kb" {
		u = 1000
	}
	if unit == "Mb" {
		u = 1000*1000
	}
	if unit == "Gb" {
		u = 1000*1000*1000
	}
	if unit == "Tb" {
		u = 1000*1000*1000*1000
	}
    
	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()
	
	ctx := context.Background()
	rdb = database.RDb() 
	// _= rdb.Del(ctx, unit+"bwtotalall"+domain+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, unit+"bwtotalall"+domain+fromtime+totime ).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection	

		rows, err := db.Query("SELECT sum(byte),sum(byteVN),(sum(byte)-sum(byteVN)),time FROM  access_"+s+ " WHERE http_host='"+domain+"' AND time >= FROM_UNIXTIME('"+fromtime+"') AND  time <= FROM_UNIXTIME('"+totime+"')   GROUP BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		bws := make([]*BandwidthAll,0)
		for rows.Next(){
			
			bw := new(BandwidthAll)
			if err := rows.Scan( &bw.Total ,&bw.VN  ,&bw.INT , &bw.Time ); err != nil {
				fmt.Println(err)
			}       
			bw.Total =  RoundDown(bw.Total*8/u/300,3)
			bw.VN =  RoundDown(bw.VN*8/u/300,3)
			bw.INT =  RoundDown(bw.INT*8/u/300,3)
			bw.Unit = unit+"ps"
 			bws = append(bws, bw)

		}

		json , err := json.Marshal(bws)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, unit+"bwtotalall"+domain+fromtime+totime , json , 1000000000 * 150).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(http.StatusOK, bws)
    }else{
		c.String(http.StatusOK, val)
		fmt.Println("BandwidthTotal  hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary referer
// @Schemes
// @Description count referer
// @Tags domain
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.ISPs
// @Param domain path string true "vovlive.1cdn.vn"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/referer/{domain}/{fromtime}/{totime} [get]
func RefererDomain(c * gin.Context){

	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "ISP"+hostname+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "referer"+hostname+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		query := "SELECT date(time),http_host,name,sum(data) FROM  %s WHERE name!='-' AND  http_host='"+hostname+"' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name  ORDER BY sum(data) DESC LIMIT 50" 
		if checkTime(totime,fromtime) == false {
			table  = "referer_"+s
		} else {
			table  = "daily_referer_"+s
		}
		queryinsrt := fmt.Sprintf(query,table)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		details := make([]*Referer,0)
		for rows.Next(){
			
			detail := new(Referer)
			if err := rows.Scan( &detail.Time,&detail.Http_Host,  &detail.Referer, &detail.Count ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "referer"+hostname+fromtime+totime , json , 1000000000 * 86400).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("daily referer hit cache")
	}
	
}

func checkTime(totime string,fromtime string) bool {
	to, _ := strconv.Atoi(totime)
	fro, _ := strconv.Atoi(fromtime)
	var status bool
	status = true 
	if to - fro < 86400 {
		status = false 
	}
	return status
}

// @BasePath /api/v1

// log godoc
// @Summary sumary data
// @Schemes
// @Description do sumary
// @Tags domain
// @Accept json
// @Produce json 
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param domain path string true "vovlive.1cdn.vn"
// @Router /log/sumarydata/{domain}/{fromtime}/{totime} [get]
func SumaryData(c * gin.Context){

	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 	
	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	
	// day 	:= tm.Day()
	
	ctx := context.Background()
	rdb = database.RDb() 
	// _= rdb.Del(ctx, "today"+hostname).Err()

    val, err := rdb.Get(ctx, "sumarydata"+hostname+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT time,http_host,sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE http_host='"+hostname+"' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"')" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time,&detail.Http_Host, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "sumarydata"+hostname+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("today detail hit cache")
	}
	
}



// @BasePath /api/v1

// log godoc
// @Summary Avgspeed data
// @Schemes
// @Description do Avgspeed
// @Tags domain
// @Accept json
// @Produce json 
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Avgspeeds
// @Param domain path string true "vovlive.1cdn.vn"
// @Router /log/avgspeed/{domain}/{unit}/{fromtime}/{totime} [get]
func AvgSpeed(c * gin.Context){

	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 
	unit := c.Param("unit")
	var u float64
	if unit == "Kb" {
		u = 1000
	}
	if unit == "Mb" {
		u = 1000*1000
	}
	if unit == "Gb" {
		u = 1000*1000*1000
	}
	if unit == "Tb" {
		u = 1000*1000*1000*1000
	}
	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	
	// day 	:= tm.Day()
	
	ctx := context.Background()
	rdb = database.RDb() 
	// _= rdb.Del(ctx, "today"+hostname).Err()

    val, err := rdb.Get(ctx, "Avgspeeds"+hostname+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT time,byte,byteVN,(byte-byteVN),requests,requestVN,(requests-requestVN)  FROM  access_"+s+ " WHERE http_host='"+hostname+"' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY time ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		details := make([]*AvgspeedDetail,0)
		for rows.Next(){
			
			detail := new(ByteRequestDetail)
			if err := rows.Scan( &detail.Time, &detail.TotalByte, &detail.VnByte, &detail.INTByte, &detail.TotalRequest, &detail.VnRequest, &detail.INTRequest ); err != nil {
				fmt.Println(err)
			}       
			avgdetail := new(AvgspeedDetail)
			avgdetail.Time = detail.Time
			avgdetail.TotalSpeed = 0
			if detail.TotalRequest > 0 {
				avgdetail.TotalSpeed = RoundDown(detail.TotalByte*8/u/300/(detail.TotalRequest/300),2)
			}
			avgdetail.VnSpeed = 0
			if detail.VnRequest > 0 {
				avgdetail.VnSpeed = RoundDown(detail.VnByte*8/u/300/(detail.VnRequest/300),2)
			}
			avgdetail.INTSpeed = 0
			if detail.INTRequest > 0 {
				avgdetail.INTSpeed = RoundDown(detail.INTByte*8/u/300/(detail.INTRequest/300),2)
			}
			avgdetail.Unit = unit
			details = append(details, avgdetail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "Avgspeeds"+hostname+fromtime+totime +unit, json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("today detail hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary top http code error data
// @Schemes
// @Description do top http code error
// @Tags domain
// @Accept json
// @Produce json 
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.TopHttpCodeError
// @Param domain path string true "vovlive.1cdn.vn"
// @Router /log/codeerr/{domain}/{fromtime}/{totime} [get]
func HttpCodeErr(c * gin.Context){
	hostname := c.Param("hostname") 
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 	
	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	
	// day 	:= tm.Day()
	
	ctx := context.Background()
	rdb = database.RDb() 
	// _= rdb.Del(ctx, "today"+hostname).Err()

    val, err := rdb.Get(ctx, "http_code_error"+hostname+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT http_host,time,sum(403x),sum(404x),sum(503x),sum(505x)  FROM  http_code_err_"+s+ " WHERE http_host='"+hostname+"' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY time ORDER BY time ASC")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		details := make([]*TopHttpCodeErrorDomain,0)
		for rows.Next(){		
			detail := new(TopHttpCodeErrorDomain)
			if err := rows.Scan( &detail.Http_Host,&detail.Time,&detail.Code403,&detail.Code404,&detail.Code503,&detail.Code505 ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)
		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "http_code_error"+hostname+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("today detail hit cache")
	}
	
}