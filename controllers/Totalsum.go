package controllers

import(
        "github.com/gin-gonic/gin"
		// "github.com/go-redis/redis/v8"
		database "logsquery/database"
		// "code.cloudfoundry.org/bytefmt"
        "log"
		_"database/sql"
        "time"
        // "net"
        // "bytes"
        "encoding/json"
        "strings"
        // "crypto/tls"
        "fmt"
		"context"
        // "net/http"
        // "reflect"
        "strconv" 
		// "database/sql"
		// "math"
		// _ "github.com/go-sql-driver/mysql"
)


// @BasePath /api/v1

// log godoc
// @Summary total min
// @Schemes
// @Description do total min
// @Tags Total
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/min/{total}/{fromtime}/{totime} [get]
func Mintotal(c * gin.Context){

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
	// _= rdb.Del(ctx, "mintotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "mintotal"+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT time,sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE   UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY time ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")

		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"mintotal"+fromtime+totime , json , 1000000000 * 180).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("min total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary total daily
// @Schemes
// @Description do total daily
// @Tags Total
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/daily/{total}/{fromtime}/{totime} [get]
func Dailytotal(c * gin.Context){

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
	// _= rdb.Del(ctx, "daily_total"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "daily_total"+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  daily_"+s+ " WHERE  UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY date(time) ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"daily_total"+fromtime+totime , json , 1000000000 * 3000).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("daily total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary total monly
// @Schemes
// @Description do total monly
// @Tags Total
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/monly/{total}/{fromtime}/{totime} [get]
func Monthlyhtotal(c * gin.Context){

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
	// _= rdb.Del(ctx, "monly_total"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "monly_total"+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  monly_"+strconv.Itoa(year)+" WHERE   UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY date(time) ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"monly_total"+fromtime+totime , json , 1000000000 * 3000).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly total hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary total today
// @Schemes
// @Description do total today
// @Tags Total
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Router /log/total/today/{total} [get]
func Todaytotal(c * gin.Context){


	current_time := time.Now()
	
	today := fmt.Sprintf("%d-%02d-%02d", current_time.Year(), current_time.Month(), int(current_time.Day())) // strconv.Itoa(current_time.Year())+"-"+strconv.Itoa(int(current_time.Month()))+"-"+strconv.Itoa(current_time.Day())
	fmt.Println(today)

    s := fmt.Sprintf("%d%02d", current_time.Year(), int(current_time.Month()))

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "todaytotal" ).Err()

    val, err := rdb.Get(ctx, "todaytotal").Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE  date(time)='"+today+"'" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"todaytotal" , json , 1000000000 * 150).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary total thismonth
// @Schemes
// @Description do total thismonth
// @Tags Total
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Router /log/total/thismonth/{total} [get]
func ThisMonthTotal(c * gin.Context){

	current_time := time.Now()
	
	thismonth := fmt.Sprintf("%d-%02d", current_time.Year(), int(current_time.Month()))

    fmt.Println(thismonth)

	s := fmt.Sprintf("%d%02d", current_time.Year(), int(current_time.Month()))
	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "thismonthtotal" ).Err()

    val, err := rdb.Get(ctx, "thismonthtotal" ).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE date(time) LIKE  '"+thismonth+"%'" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"thismonthtotal" , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary time Respone
// @Schemes
// @Description count time Respone
// @Tags Total
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.TimeResponse
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/responseTime/{fromtime}/{totime} [get]
func ResponseTimeTotal(c * gin.Context){

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
	// _= rdb.Del(ctx, "ResponseTimeTotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ResponseTimeTotal"+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT time,sum(under4s),sum(over4s),sum(over6s) FROM  time_response_"+s+ " WHERE time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY time ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		
		// var detail Detail
		details := make([]*TimeResponse,0)
		for rows.Next(){
			
			detail := new(TimeResponse)
			if err := rows.Scan(  &detail.Time,  &detail.Under4s, &detail.Over4s, &detail.Over6s ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "ResponseTimeTotal"+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("ResponseTimeTotal hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary country
// @Schemes
// @Description count country
// @Tags Total
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.CountryISO
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/country/{fromtime}/{totime} [get]
func CountryTotal(c * gin.Context){

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
	// _= rdb.Del(ctx, "contryTotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "contryTotal"+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT date(time),name,sum(data) FROM  geoip_country_"+s+ " WHERE   time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		// fmt.Println("SELECT time,http_host,name,data FROM  geoip_"+strconv.Itoa(year)+ strconv.Itoa(month)+ " WHERE http_host='"+hostname+"' AND types='country' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name,time  ORDER BY time ASC" )
		// var detail Detail
		details := make([]*CountryISO,0)
		for rows.Next(){
			
			detail := new(CountryISO)
			if err := rows.Scan( &detail.Time,  &detail.Country, &detail.Count ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "contryTotal"+fromtime+totime , json , 1000000000 * 600).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println(" Country total hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary ISP
// @Schemes
// @Description count ISP
// @Tags Total
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.ISPs
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/isp/{fromtime}/{totime} [get]
func ISPTotal(c * gin.Context){

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
	// _= rdb.Del(ctx, "ISPTotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ISPTotal"+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT date(time),name,sum(data) FROM  geoip_isp_"+s+ " WHERE   time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		// fmt.Println("SELECT time,http_host,name,data FROM  geoip_"+strconv.Itoa(year)+ strconv.Itoa(month)+ " WHERE http_host='"+hostname+"' AND types='country' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name,time  ORDER BY time ASC" )
		// var detail Detail
		details := make([]*ISPs,0)
		for rows.Next(){
			
			detail := new(ISPs)
			if err := rows.Scan( &detail.Time,  &detail.Isp, &detail.Count ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "ISPTotal"+fromtime+totime , json , 1000000000 * 600).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("ISP total hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary City
// @Schemes
// @Description count City
// @Tags Total
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.City
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/city/{fromtime}/{totime} [get]
func CityTotal(c * gin.Context){

	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	// year 	:= tm.Year()
	// month 	:= int(tm.Month())
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "CityTotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "CityTotal"+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT date(time),name,sum(data) FROM  geoip_city_"+ s + " WHERE   time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")

		details := make([]*City,0)
		for rows.Next(){
			
			detail := new(City)
			if err := rows.Scan( &detail.Time,  &detail.City, &detail.Count ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "CityTotal"+fromtime+totime , json , 1000000000 * 600).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("ISP total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary sumary bandwidth Total
// @Schemes
// @Description do sumary bandwidth Total
// @Tags Total
// @Accept json
// @Produce json
// @success 200 {object} []controllers.Bandwidth
// @Param unit path string true "Kb,Mb,Gb,Tb"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/bandwidth//{unit}/{fromtime}/{totime} [get]
func BandwidthTotal(c * gin.Context){

	current_time := time.Now()
	var tunix int64
	unit := c.Param("unit")
	fromtime := c.Param("fromtime")
	totime	 := c.Param("totime")
	if fromtime == "now" {
		tunix = (current_time.Unix() - 1800)
		fromtime = strconv.FormatInt(tunix, 10)
		totime   = strconv.FormatInt(current_time.Unix(), 10)
	}

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	// year 	:= tm.Year()
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	
	// fmt.Println(s)
	var u float64
	// thismonth := strconv.Itoa(current_time.Year())+"-"+strconv.Itoa(int(current_time.Month()))
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
    
	// day 	:= tm.Day() 
	
	ctx := context.Background()
	rdb = database.RDb() 
	// _= rdb.Del(ctx, unit+"bwtotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, unit+"bwtotal"+fromtime+totime ).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection	

		rows, err := db.Query("SELECT sum(byte),time FROM  access_"+ s +" WHERE   time >= FROM_UNIXTIME('"+fromtime+"') AND  time <= FROM_UNIXTIME('"+totime+"')   GROUP BY time ASC" )
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
		errs := rdb.Set(ctx, unit+"bwtotal"+fromtime+totime , json , 1000000000 * 180).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, bws)
    }else{
		c.String(200, val)
		fmt.Println("BandwidthTotal  hit cache")
	}
	
}

// ISPnode
// @BasePath /api/v1

// log godoc
// @Summary total minn ode
// @Schemes
// @Description do total min node
// @Tags Total Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/min/{node}/{total}/{fromtime}/{totime} [get]
func MintotalISP(c * gin.Context){

	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 
	node :=c.Param("node")

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "mintotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "mintotal"+node+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT time,sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE %s UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY time ORDER BY time ASC" 
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")

		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"mintotal"+node+fromtime+totime , json , 1000000000 * 180).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("min total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary total daily 
// @Schemes
// @Description do total daily isp
// @Tags Total Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/daily/isp/{node}/{total}/{fromtime}/{totime} [get]
func DailytotalISP(c * gin.Context){

	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 
	node :=c.Param("node") 

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "daily_total"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "daily_total"+node+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  daily_"+s+ " WHERE %s UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY date(time) ORDER BY time ASC" 
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf(queryinsrt)
		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"daily_total"+node+fromtime+totime , json , 1000000000 * 3000).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("daily total hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary total monly isp
// @Schemes
// @Description do total monly isp
// @Tags Total Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/monly/isp/{node}/{total}/{fromtime}/{totime} [get]
func MonthlyhtotalISP(c * gin.Context){

	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 
	node :=c.Param("node")

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
	// _= rdb.Del(ctx, "monly_total"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "monly_total"+node+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  monly_"+strconv.Itoa(year)+" WHERE %s UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY date(time) ORDER BY time ASC" 
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"monly_total"+node+fromtime+totime , json , 1000000000 * 3000).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary total today isp
// @Schemes
// @Description do total today isp
// @Tags Total Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param node path string true "viettel"
// @Router /log/total/today/{node} [get]
func TodaytotalISP(c * gin.Context){

	node :=c.Param("node")
	current_time := time.Now()
	
	today := fmt.Sprintf("%d-%02d-%02d", current_time.Year(), current_time.Month(), int(current_time.Day())) // strconv.Itoa(current_time.Year())+"-"+strconv.Itoa(int(current_time.Month()))+"-"+strconv.Itoa(current_time.Day())
	fmt.Println(today)

    s := fmt.Sprintf("%d%02d", current_time.Year(), int(current_time.Month()))

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "todaytotal" ).Err()

    val, err := rdb.Get(ctx, "todaytotal"+node).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE %s  date(time)='"+today+"'"
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		rows, err := db.Query(queryinsrt )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"todaytotal"+node , json , 1000000000 * 150).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary total thismonth isp
// @Schemes
// @Description do total thismonth isp
// @Tags Total Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param node path string true "viettel"
// @Router /log/total/thismonth/{node} [get]
func ThisMonthTotalISP(c * gin.Context){

	node :=c.Param("node")
	current_time := time.Now()
	
	thismonth := fmt.Sprintf("%d-%02d", current_time.Year(), int(current_time.Month()))

    fmt.Println(thismonth)

	s := fmt.Sprintf("%d%02d", current_time.Year(), int(current_time.Month()))
	ctx := context.Background()
	rdb = database.RDb()
	_= rdb.Del(ctx, "thismonthtotal" ).Err()

    val, err := rdb.Get(ctx, "thismonthtotal"+node ).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE stringreplace date(time) LIKE  '"+thismonth+"%'" 
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
		// queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"thismonthtotal"+node , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc 
// @Summary time Respone isp
// @Schemes
// @Description count time Respone isp
// @Tags Total Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.TimeResponse
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/responseTime/isp/{node}/{fromtime}/{totime} [get]
func ResponseTimeTotalISP(c * gin.Context){

	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 
	node :=c.Param("node")

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "ResponseTimeTotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ResponseTimeTotal"+node+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT time,sum(under4s),sum(over4s),sum(over6s) FROM  time_response_"+s+ " WHERE %s time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY time ORDER BY time ASC"
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		rows, err := db.Query( queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		
		// var detail Detail
		details := make([]*TimeResponse,0)
		for rows.Next(){
			
			detail := new(TimeResponse)
			if err := rows.Scan(  &detail.Time,  &detail.Under4s, &detail.Over4s, &detail.Over6s ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "ResponseTimeTotal"+node+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("ResponseTimeTotal hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary country isp
// @Schemes
// @Description count country isp
// @Tags Total Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.CountryISO
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Param node path string true "viettel"
// @Router /log/total/country/isp/{node}/{fromtime}/{totime} [get]
func CountryTotalISP(c * gin.Context){

	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 
	node :=c.Param("node")

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "contryTotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "contryTotal"+node+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT date(time),name,sum(data) FROM  %s WHERE %s  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC"  
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		table := fmt.Sprintf("geoip_country_"+s)
		if checkTime(totime,fromtime) == true {
			table = "daily_geoip_country_"+s
		}
		queryinsrt = fmt.Sprintf(queryinsrt , table , nodestr)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		// fmt.Println("SELECT time,http_host,name,data FROM  geoip_"+strconv.Itoa(year)+ strconv.Itoa(month)+ " WHERE http_host='"+hostname+"' AND types='country' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name,time  ORDER BY time ASC" )
		// var detail Detail
		details := make([]*CountryISO,0)
		for rows.Next(){
			
			detail := new(CountryISO)
			if err := rows.Scan( &detail.Time,  &detail.Country, &detail.Count ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "contryTotal"+node+fromtime+totime , json , 1000000000 * 86400).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println(" Country total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary ISP isp
// @Schemes
// @Description count ISP isp
// @Tags Total Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.ISPs
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/isp/node/{node}/{fromtime}/{totime} [get]
func ISPTotalISP(c * gin.Context){
	node :=c.Param("node")
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
	// _= rdb.Del(ctx, "ISPTotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ISPTotal"+node+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT date(time),name,sum(data) FROM  %s WHERE %s time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC" 
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		table := fmt.Sprintf("geoip_isp_"+s)
		if checkTime(totime,fromtime) == true {
			table = "daily_geoip_isp_"+s
		}
		queryinsrt = fmt.Sprintf(queryinsrt , table , nodestr)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		details := make([]*ISPs,0)
		for rows.Next(){
			
			detail := new(ISPs)
			if err := rows.Scan( &detail.Time,  &detail.Isp, &detail.Count ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "ISPTotal"+node+fromtime+totime , json , 1000000000 * 86400).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("ISP total hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary City isp
// @Schemes
// @Description count City isp
// @Tags Total Node 
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.City
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Param node path string true "viettel"
// @Router /log/total/city/isp/{node}/{fromtime}/{totime} [get]
func CityTotalISP(c * gin.Context){

	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 
	node :=c.Param("node")
	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	// year 	:= tm.Year()
	// month 	:= int(tm.Month())
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "CityTotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "CityTotal"+node+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT date(time),name,sum(data) FROM  %s WHERE %s  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC"
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}

		table := fmt.Sprintf("geoip_city_"+s)
		if checkTime(totime,fromtime) == true {
			table = "daily_geoip_city_"+s
		}
		queryinsrt = fmt.Sprintf(queryinsrt , table , nodestr)
		rows, err := db.Query(queryinsrt )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")

		details := make([]*City,0)
		for rows.Next(){
			
			detail := new(City)
			if err := rows.Scan( &detail.Time,  &detail.City, &detail.Count ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "CityTotal"+node+fromtime+totime , json , 1000000000 * 86400).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("ISP total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary sumary bandwidth Total isp
// @Schemes
// @Description do sumary bandwidth Total isp
// @Tags Total Node
// @Accept json
// @Produce json
// @success 200 {object} []controllers.Bandwidth
// @Param unit path string true "Kb,Mb,Gb,Tb"
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/bandwidth/{node}/{unit}/{fromtime}/{totime} [get]
func BandwidthTotalISP(c * gin.Context){
	node :=c.Param("node")
	current_time := time.Now()
	var tunix int64
	unit := c.Param("unit")
	fromtime := c.Param("fromtime")
	totime	 := c.Param("totime")
	if fromtime == "now" {
		tunix = (current_time.Unix() - 1800)
		fromtime = strconv.FormatInt(tunix, 10)
		totime   = strconv.FormatInt(current_time.Unix(), 10)
	}

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	// year 	:= tm.Year()
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	
	// fmt.Println(s)
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
    
	// day 	:= tm.Day() 
	
	ctx := context.Background()
	rdb = database.RDb() 
	// _= rdb.Del(ctx, unit+"bwtotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, unit+"bwtotal"+node+fromtime+totime ).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection	
		queryinsrt := "SELECT sum(byte),sum(byteVN),(sum(byte)-sum(byteVN)),time FROM  access_"+ s +" WHERE  %s time >= FROM_UNIXTIME('"+fromtime+"') AND  time <= FROM_UNIXTIME('"+totime+"')   GROUP BY time ASC" 
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		bws := make([]*BandwidthAll,0)
		for rows.Next(){
			
			bw := new(BandwidthAll)
			if err := rows.Scan(  &bw.Total , &bw.VN , &bw.INT , &bw.Time ); err != nil {
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
		errs := rdb.Set(ctx, unit+"bwtotal"+node+fromtime+totime , json , 1000000000 * 180).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, bws)
    }else{
		c.String(200, val)
		fmt.Println("BandwidthTotal  hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary total ccu
// @Schemes
// @Description do total ccu
// @Tags Total Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Ccu
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/total/ccu/isp/{node}/{total}/{fromtime}/{totime} [get]
func CCuTotal(c * gin.Context){

	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 
	node :=c.Param("node")

	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	// month 	:= int(tm.Month())
	// day 	:= tm.Day()

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "monly_total"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ccu_total"+node+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT time,ROUND(sum(requests)/300) FROM  access_"+s+ " WHERE  %s  UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY time ORDER BY time ASC"
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*Ccu,0)
		for rows.Next(){
			
			detail := new(Ccu)
			if err := rows.Scan( &detail.Time, &detail.Ccu); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"ccu_total"+node+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary Referer isp
// @Schemes
// @Description count Referer isp
// @Tags Total Node 
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Referer
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Param node path string true "viettel"
// @Router /log/total/referer/isp/{node}/{fromtime}/{totime} [get]
func RefererTotalISP(c * gin.Context){

	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 
	node :=c.Param("node")
	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	// year 	:= tm.Year()
	// month 	:= int(tm.Month())
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "RefererTotal"+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "RefererTotalISP"+node+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT date(time),name,sum(data) FROM  %s WHERE %s name!='-' AND  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY sum(data) DESC LIMIT 50"
		var nodestr string
		var table string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		if checkTime(totime,fromtime) == false {
			table = "daily_referer_"+s
		}else{
			table = "daily_referer_"+s
		}
		queryinsrt = fmt.Sprintf(queryinsrt , table , nodestr)
		rows, err := db.Query(queryinsrt )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")

		details := make([]*Referer,0)
		for rows.Next(){
			
			detail := new(Referer)
			if err := rows.Scan(&detail.Time,  &detail.Referer, &detail.Count ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "RefererTotalISP"+node+fromtime+totime , json , 1000000000 * 86400).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("Referer total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary total summary isp
// @Schemes
// @Description do total summary isp
// @Tags  Total Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param node path string true "viettel"
// @Router /log/total/sumary/{node}/{fromtime}/{totime} [get]
func SummaryTotalISP(c * gin.Context){

	node :=c.Param("node")
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime")
	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "todaytotal" ).Err()

    val, err := rdb.Get(ctx, "SummaryTotalISP"+node+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM   %s WHERE %s  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"')"
		var nodestr string
		var table string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		if checkTime(totime,fromtime) == false {
			table = "access_"+s
		}else{
			table = "daily_"+s
		}
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
		queryinsrt = fmt.Sprintf(queryinsrt , table , nodestr)
		rows, err := db.Query(queryinsrt )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*Detail,0)
		for rows.Next(){
			
			detail := new(Detail)
			if err := rows.Scan( &detail.Time, &detail.Byte, &detail.Requests, &detail.Bytevn, &detail.Requestvn, &detail.Cache_hit, &detail.Cache_miss, &detail.Status_2x, &detail.Status_3x, &detail.Status_4x, &detail.Status_5x, &detail.Macintosh, &detail.Windows , &detail.Ios , &detail.Android , &detail.Smartv , &detail.Other ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"SummaryTotalISP"+node+fromtime+totime , json , 1000000000 * 3600).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Avgspeed total Avgspeed isp
// @Schemes
// @Description do total Avgspeed isp
// @Tags  Total Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Avgspeeds
// @Param node path string true "viettel"
// @Router /log/total/avgspeed/{node}/{unit}/{fromtime}/{totime} [get]
func TotalAvgspeed(c * gin.Context){

	node :=c.Param("node")
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

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "todaytotal" ).Err()

    val, err := rdb.Get(ctx, "TotalAvgspeed"+node+fromtime+totime+unit).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT time,sum(byte),sum(byteVN),(sum(byte)-sum(byteVN)),sum(if(byte>0,requests,0)),sum(if(byte>0,requestVN,0)),(sum(if(byte>0,requests,0))-sum(if(byte>0,requestVN,0))) FROM  access_"+s+ " WHERE %s  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY time ORDER BY time ASC"
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		rows, err := db.Query(queryinsrt )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*AvgspeedDetail,0)
		for rows.Next(){
			
			detail := new(ByteRequestDetail)
			if err := rows.Scan( &detail.Time, &detail.TotalByte, &detail.VnByte, &detail.INTByte, &detail.TotalRequest, &detail.VnRequest, &detail.INTRequest); err != nil {
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
		errs := rdb.Set(ctx,"TotalAvgspeed"+node+fromtime+totime+unit , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary total TopHttpCodeError isp
// @Schemes
// @Description do total TopHttpCodeError isp
// @Tags  Total Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.TopHttpCodeError
// @Param node path string true "viettel"
// @Router /log/total/codeerr/{node}/{fromtime}/{totime} [get]
func TopHttpCodeErrorTotalISP(c * gin.Context){

	node :=c.Param("node")
	code :=c.Param("code")
	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime")
	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "TopHttpCodeErrorTotalISP"+node+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "TopHttpCodeErrorTotalISP"+node+fromtime+totime+code).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection
		queryinsrt := "SELECT http_host,date(time),sum("+code+") as count FROM  %s WHERE %s time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY http_host HAVING count > 0 ORDER BY sum("+code+") DESC LIMIT 20"
		var nodestr string
		var table string
		fmt.Println(queryinsrt)
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		if checkTime(totime,fromtime) == false {
			table = "http_code_err_"+s
		}else{
			table = "daily_http_code_err_"+s
		}
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
		queryinsrt = fmt.Sprintf(queryinsrt , table , nodestr)
		rows, err := db.Query(queryinsrt )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*TopHttpCodeError,0)
		for rows.Next(){
			
			detail := new(TopHttpCodeError)
			if err := rows.Scan( &detail.Http_Host, &detail.Time,  &detail.Count ); err != nil {
				fmt.Println(err)
			}  
			detail.HttpCodeError = code     //http_code_err   
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"TopHttpCodeErrorTotalISP"+node+fromtime+totime+code , json , 1000000000 * 600).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly total hit cache")
	}
	
}