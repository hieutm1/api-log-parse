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
// @Summary all domain of user
// @Schemes
// @Description get all domain of user
// @Tags User
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: "Basic Authen"
// @Param user path string true "vtvcab"
// @success 200 {object} []controllers.DomainHost
// @Router /log/user/domains/{user} [get]
func DomainUser(c * gin.Context){

	user := c.Param("user") 
	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "user").Err()

    val, err := rdb.Get(ctx, "user").Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rows, err := db.Query("SELECT http_host,customerName FROM http_host WHERE customerName='"+user+"' ORDER BY http_host")
		if err != nil {
			fmt.Println(err.Error())
		}

		defer db.Close()

		log.Printf("Successfully connected to database")
		// var detail Detail
		details := make([]*DomainHost,0)
		for rows.Next(){
			
			detail := new(DomainHost)
			if err := rows.Scan( &detail.Http_Host,&detail.User); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

			}
		jshttp , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(jshttp)
		errs := rdb.Set(ctx, "user", jshttp , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("user domain cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary user min
// @Schemes
// @Description do user min
// @Tags User
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param user path string true "vtvcab"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/min/{user}/{fromtime}/{totime} [get]
func MinUser(c * gin.Context){

	user := c.Param("user") 
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
	// _= rdb.Del(ctx, "min"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "min"+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		rows, err := db.Query("SELECT time,sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE http_host IN ("+search+") AND  UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY time ORDER BY time ASC" )
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
		errs := rdb.Set(ctx,"min"+user+fromtime+totime , json , 1000000000 * 60).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("min user hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary user daily
// @Schemes
// @Description do user daily
// @Tags User
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param user path string true "vtvcab"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/daily/{user}/{fromtime}/{totime} [get]
func DailyUser(c * gin.Context){

	user := c.Param("user") 
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
	// _= rdb.Del(ctx, "daily_"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "daily_"+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		rows, err := db.Query("SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  daily_"+s+ " WHERE http_host IN ("+search+") AND  UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY date(time) ORDER BY time ASC" )
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
		errs := rdb.Set(ctx,"daily_"+user+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("daily user hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary user monly
// @Schemes
// @Description do user monly
// @Tags User
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param user path string true "vtvcab"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/monly/{user}/{fromtime}/{totime} [get]
func MonthlyhUser(c * gin.Context){

	user := c.Param("user") 
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
	// _= rdb.Del(ctx, "monly_"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "monly_"+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		rows, err := db.Query("SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  monly_"+strconv.Itoa(year)+" WHERE http_host IN ("+search+") AND  UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY date(time) ORDER BY time ASC" )
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
		errs := rdb.Set(ctx,"monly_"+user+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly user hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary user today
// @Schemes
// @Description do user today
// @Tags User
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param user path string true "vtvcab"
// @Router /log/user/today/{user} [get]
func TodayUser(c * gin.Context){

	user := c.Param("user") 

	current_time := time.Now()
	
	today := fmt.Sprintf("%d-%02d-%02d", current_time.Year(), current_time.Month(), int(current_time.Day()))

    fmt.Println(current_time)
	// year 	:= current_time.Year()
	// month 	:= int(current_time.Month())
	s := fmt.Sprintf("%d%02d", current_time.Year(), int(current_time.Month()))
	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "today"+user ).Err()

    val, err := rdb.Get(ctx, "today"+user).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		rows, err := db.Query("SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE http_host IN ("+search+") AND  date(time)='"+today+"'" )
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
		errs := rdb.Set(ctx,"today"+user , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly user hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary user thismonth
// @Schemes
// @Description do user thismonth
// @Tags User
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param user path string true "vtvcab"
// @Router /log/user/thismonth/{user} [get]
func ThisMonthUser(c * gin.Context){

	user := c.Param("user") 

	current_time := time.Now()
	
	thismonth := fmt.Sprintf("%d-%02d", current_time.Year(), int(current_time.Month()))

    fmt.Println(current_time)
	s := fmt.Sprintf("%d%02d", current_time.Year(), int(current_time.Month()))

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "thismonth"+user ).Err()

    val, err := rdb.Get(ctx, "thismonth"+user).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		rows, err := db.Query("SELECT date(time),sum(byte),sum(requests),COALESCE(sum(byteVN),0),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE http_host IN ("+search+") AND  date(time) LIKE  '"+thismonth+"%'" )
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
		errs := rdb.Set(ctx,"thismonth"+user , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly user hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary country
// @Schemes
// @Description count country
// @Tags User
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.CountryISO
// @Param user path string true "vtvcab"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/country/{user}/{fromtime}/{totime} [get]
func CountryUser(c * gin.Context){

	user :=c.Param("user")
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
	// _= rdb.Del(ctx, "CountryUser"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "CountryUser"+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")

		rows, err := db.Query("SELECT date(time),name,sum(data) FROM  geoip_country_"+s+ " WHERE http_host IN  ("+search+") AND  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		// fmt.Println("SELECT time,http_host,name,data FROM  geoip_"+s+ " WHERE http_host='"+hostname+"' AND types='country' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name,time  ORDER BY time ASC" )
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
		errs := rdb.Set(ctx, "CountryUser"+user+fromtime+totime , json , 1000000000 * 300).Err()
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
// @Tags User
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.ISPs
// @Param user path string true "vtvcab"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/isp/{user}/{fromtime}/{totime} [get]
func ISPUser(c * gin.Context){

	user :=c.Param("user")
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
	// _= rdb.Del(ctx, "ISPTotal"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ISPTotal"+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		rows, err := db.Query("SELECT date(time),name,sum(data) FROM  geoip_isp_"+s+ " WHERE http_host IN  ("+search+") AND  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		// fmt.Println("SELECT time,http_host,name,data FROM  geoip_"+s+ " WHERE http_host='"+hostname+"' AND types='country' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name,time  ORDER BY time ASC" )
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
		errs := rdb.Set(ctx, "ISPTotal"+user+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("ISPUser total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary time Respone
// @Schemes
// @Description count time Respone
// @Tags User
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.TimeResponse
// @Param user path string true "vtvcab"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/responseTime/{user}/{fromtime}/{totime} [get]
func ResponseTimeUser(c * gin.Context){

	user :=c.Param("user")
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
	// _= rdb.Del(ctx, "ResponseTimeUser"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ResponseTimeUser"+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		
		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		rows, err := db.Query("SELECT time,sum(under4s),sum(over4s),sum(over6s) FROM  time_response_"+s+ " WHERE http_host IN  ("+search+") AND  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY time ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		
		// var detail Detail
		details := make([]*TimeResponse,0)
		for rows.Next(){
			
			detail := new(TimeResponse)
			if err := rows.Scan( &detail.Time, &detail.Under4s, &detail.Over4s, &detail.Over6s  ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "ResponseTimeUser"+user+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("ResponseTimeUser hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary City
// @Schemes
// @Description count City
// @Tags User
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.City
// @Param user path string true "vtvcab"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/city/{user}/{fromtime}/{totime} [get]
func CityUser(c * gin.Context){

	user :=c.Param("user")
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
	// _= rdb.Del(ctx, "City"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "City"+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt := "SELECT date(time),name,sum(data) FROM %s WHERE http_host IN  ("+search+") AND   time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC" 
		var table string
		if checkTime(totime,fromtime) == false {
			table = "geoip_city_"+s
		}else{
			table = "daily_geoip_city_"+s
		}
		query := fmt.Sprintf(queryinsrt, table)
		rows, err := db.Query(query)
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
		errs := rdb.Set(ctx, "City"+user+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("City User total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary sumary bandwidth of domain 
// @Schemes
// @Description do sumary bandwidth of domain 
// @Tags User
// @Accept json
// @Produce json
// @success 200 {object} []controllers.Bandwidth
// @Param user path string true "vtvcab"
// @Param unit path string true "Kb,Mb,Gb,Tb"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/bandwidth/{user}/{unit}/{fromtime}/{totime} [get]
func BandwidthUser(c * gin.Context){

	current_time := time.Now()
	var tunix int64
	unit := c.Param("unit")
	user :=c.Param("user")
	fromtime := c.Param("fromtime")
	totime	 := c.Param("totime")
	if fromtime == "now" {
		tunix = (current_time.Unix() - 1800)
		fromtime = strconv.FormatInt(tunix, 10)
		totime   = strconv.FormatInt(current_time.Unix(), 10)
	}
	fmt.Println(tunix)
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
    
	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	
	// day 	:= tm.Day()
	ctx := context.Background()
	rdb = database.RDb() 
	// _= rdb.Del(ctx, unit+"bwtotal"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, unit+"bwtotal"+user+fromtime+totime ).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection	

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		
		rows, err := db.Query("SELECT sum(byte),time FROM  access_"+s+ " WHERE http_host IN  ("+search+") AND time >= FROM_UNIXTIME('"+fromtime+"') AND  time <= FROM_UNIXTIME('"+totime+"')   GROUP BY time ASC" )
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
		errs := rdb.Set(ctx, unit+"bwtotal"+user+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, bws)
    }else{
		c.String(200, val)
		fmt.Println("BandwidthTotal  hit cache")
	}
	
}

//ISP

// @BasePath /api/v1

// log godoc
// @Summary user min isp
// @Schemes 
// @Description do user min isp
// @Tags  User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/min/isp/{node}/{user}/{fromtime}/{totime} [get]
func MinUserISP(c * gin.Context){
	node := c.Param("node")
	user := c.Param("user") 
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
	// _= rdb.Del(ctx, "min"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "min"+node+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE  customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt := "SELECT time,sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE %s http_host IN ("+search+") AND  UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY time ORDER BY time ASC" 
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
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
		errs := rdb.Set(ctx,"min"+node+user+fromtime+totime , json , 1000000000 * 60).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("min user hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary user daily isp
// @Schemes
// @Description do user daily isp
// @Tags  User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/daily/isp/{node}/{user}/{fromtime}/{totime} [get]
func DailyUserISP(c * gin.Context){
	node := c.Param("node")
	user := c.Param("user") 
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
	// _= rdb.Del(ctx, "daily_"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "daily_"+node+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt :="SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  daily_"+s+ " WHERE %s http_host IN ("+search+") AND  UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY date(time) ORDER BY time ASC" 
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
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
		errs := rdb.Set(ctx,"daily_"+node+user+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("daily user hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary user monly isp
// @Schemes
// @Description do user monly isp
// @Tags  User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param user path string true "vtvcab"
// @Param fromtime path int true "1636633800"
// @Param user path string true "vnpt"
// @Param totime path int true "1637810200"
// @Router /log/user/monly/isp/{node}/{user}/{fromtime}/{totime} [get]
func MonthlyhUserISP(c * gin.Context){
	node :=c.Param("node")
	user := c.Param("user") 
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
	// _= rdb.Del(ctx, "monly_"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "monly_"+node+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt := "SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  monly_"+strconv.Itoa(year)+" WHERE %s  http_host IN ("+search+") AND  UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY date(time) ORDER BY time ASC" 
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
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
		errs := rdb.Set(ctx,"monly_"+node+user+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly user hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary user today isp
// @Schemes
// @Description do user today isp
// @Tags  User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Router /log/user/today/isp/{node}/{user} [get]
func TodayUserISP(c * gin.Context){
	node :=c.Param("node")
	user := c.Param("user") 

	current_time := time.Now()
	
	today := fmt.Sprintf("%d-%02d-%02d", current_time.Year(), current_time.Month(), int(current_time.Day()))

    fmt.Println(current_time)
	// year 	:= current_time.Year()
	// month 	:= int(current_time.Month())
	s := fmt.Sprintf("%d%02d", current_time.Year(), int(current_time.Month()))
	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "today"+user ).Err()

    val, err := rdb.Get(ctx, "today"+node+user).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt := "SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE  %s http_host IN ("+search+") AND  date(time)='"+today+"'"  
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
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
		errs := rdb.Set(ctx,"today"+node+user , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly user hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary user thismonth isp
// @Schemes
// @Description do user thismonth isp
// @Tags  User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Router /log/user/thismonth/{node}/{user} [get]
func ThisMonthUserISP(c * gin.Context){

	user := c.Param("user") 
	node :=c.Param("node")
	current_time := time.Now()
	
	thismonth := fmt.Sprintf("%d-%02d", current_time.Year(), int(current_time.Month()))

    fmt.Println(current_time)
	s := fmt.Sprintf("%d%02d", current_time.Year(), int(current_time.Month()))

	ctx := context.Background()
	rdb = database.RDb()
	// _= rdb.Del(ctx, "thismonth"+user ).Err()

    val, err := rdb.Get(ctx, "thismonth"+node+user).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt :="SELECT date(time),sum(byte),sum(requests),COALESCE(sum(byteVN),0),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE stringreplace http_host IN ("+search+") AND  date(time) LIKE  '"+thismonth+"%'" 
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		// queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
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
		errs := rdb.Set(ctx,"thismonth"+node+user , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly user hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary country isp
// @Schemes
// @Description count country isp
// @Tags  User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.CountryISO
// @Param user path string true "vtvcab"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Param node path string true "viettel"
// @Router /log//user/country/isp/{node}/{user}/{fromtime}/{totime} [get]
func CountryUserISP(c * gin.Context){
	node :=c.Param("node")
	user :=c.Param("user")
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
	// _= rdb.Del(ctx, "CountryUser"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "CountryUser"+node+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt := "SELECT date(time),name,sum(data) FROM  %s WHERE %s  http_host IN  ("+search+")  AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC"
		var nodestr  string
		var table string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		if checkTime(totime,fromtime) == false {
			table = "geoip_country_"+s 
		}else{
			table = "daily_geoip_country_"+s
		}
		queryinsrt = fmt.Sprintf(queryinsrt , table ,nodestr)
		rows, err := db.Query(queryinsrt)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		// fmt.Println("SELECT time,http_host,name,data FROM  geoip_"+s+ " WHERE http_host='"+hostname+"' AND types='country' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name,time  ORDER BY time ASC" )
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
		errs := rdb.Set(ctx, "CountryUser"+node+user+fromtime+totime , json , 1000000000 * 86400).Err()
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
// @Tags  User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.ISPs
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/isp/isp/node/{node}/{user}/{fromtime}/{totime} [get]
func ISPUserISP(c * gin.Context){

	node :=c.Param("node")
	user :=c.Param("user")
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
	// _= rdb.Del(ctx, "ISPTotal"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ISPTotal"+node+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt := "SELECT date(time),name,sum(data) FROM  %s  WHERE %s  http_host IN  ("+search+")  AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC"
		var nodestr string
		var table string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		if checkTime(totime,fromtime) == false {
			table = "geoip_isp_"+s
		}else{
			table = "daily_geoip_isp_"+s
		}
		queryinsrt = fmt.Sprintf(queryinsrt , table , nodestr)
		rows, err := db.Query( queryinsrt )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		// fmt.Println("SELECT time,http_host,name,data FROM  geoip_"+s+ " WHERE http_host='"+hostname+"' AND types='country' AND time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name,time  ORDER BY time ASC" )
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
		errs := rdb.Set(ctx, "ISPTotal"+node+user+fromtime+totime , json , 1000000000 * 86400).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("ISPUser total hit cache")
	}
	
}


// @BasePath /api/v1

// log godoc
// @Summary City isp
// @Schemes
// @Description count City isp
// @Tags  User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.City
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/city/isp/{node}/{user}/{fromtime}/{totime} [get]
func CityUserISP(c * gin.Context){
	node :=c.Param("node")
	user :=c.Param("user")
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
	// _= rdb.Del(ctx, "City"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "City"+node+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt := "SELECT date(time),name,sum(data) FROM  %s WHERE %s http_host IN  ("+search+") AND  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY time ASC" 
		var nodestr string
		var table string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		if checkTime(totime,fromtime) == false {
			table = "geoip_city_"+s
		}else{
			table = "daily_geoip_city_"+s
		}
		
		queryinsrt = fmt.Sprintf(queryinsrt , table ,nodestr)
		fmt.Println(checkTime(totime,fromtime))
		rows, err := db.Query(queryinsrt)
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
		errs := rdb.Set(ctx, "City"+node+user+fromtime+totime , json , 1000000000 * 86400).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("City User total hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary time Respone isp
// @Schemes
// @Description count time Respone isp
// @Tags  User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.TimeResponse
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/responseTime/isp/{node}/{user}/{fromtime}/{totime} [get]
func ResponseTimeUserISP(c * gin.Context){
	node :=c.Param("node")
	user :=c.Param("user")
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
	// _= rdb.Del(ctx, "ResponseTimeUser"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ResponseTimeUser"+node+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		
		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		rows, err := db.Query("SELECT time,sum(under4s),sum(over4s),sum(over6s) FROM  time_response_"+s+ " WHERE nodeISP='"+node+"' AND http_host IN  ("+search+") AND  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY time ORDER BY time ASC" )
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		log.Printf("Successfully connected to database")
		
		
		// var detail Detail
		details := make([]*TimeResponse,0)
		for rows.Next(){
			
			detail := new(TimeResponse)
			if err := rows.Scan( &detail.Time, &detail.Under4s, &detail.Over4s, &detail.Over6s  ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "ResponseTimeUser"+node+user+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("ResponseTimeUser hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary sumary bandwidth of domain  isp
// @Schemes
// @Description do sumary bandwidth of domain  isp
// @Tags User Node
// @Accept json
// @Produce json
// @success 200 {object} []controllers.Bandwidth
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Param unit path string true "Kb,Mb,Gb,Tb"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/bandwidth/isp/{node}/{user}/{unit}/{fromtime}/{totime} [get]
func BandwidthUserISP(c * gin.Context){
	node :=c.Param("node")
	current_time := time.Now()
	var tunix int64
	unit := c.Param("unit")
	user :=c.Param("user")
	fromtime := c.Param("fromtime")
	totime	 := c.Param("totime")
	if fromtime == "now" {
		tunix = (current_time.Unix() - 1800)
		fromtime = strconv.FormatInt(tunix, 10)
		totime   = strconv.FormatInt(current_time.Unix(), 10)
	}
	fmt.Println(tunix)
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
    
	i, err := strconv.ParseInt(fromtime, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    tm := time.Unix(i, 0)
    
	s := fmt.Sprintf("%d%02d", tm.Year(), int(tm.Month()))
	
	// day 	:= tm.Day()
	ctx := context.Background()
	rdb = database.RDb() 
	// _= rdb.Del(ctx, unit+"bwtotal"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, unit+"bwtotal"+node+user+fromtime+totime ).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection	

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt :="SELECT sum(byte),sum(byteVN),(sum(byte)-sum(byteVN)),time FROM  access_"+s+ " WHERE %s http_host IN  ("+search+") AND time >= FROM_UNIXTIME('"+fromtime+"') AND  time <= FROM_UNIXTIME('"+totime+"')   GROUP BY time ASC" 
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
			if err := rows.Scan( &bw.Total , &bw.VN , &bw.INT , &bw.Time ); err != nil {
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
		errs := rdb.Set(ctx, unit+"bwtotal"+node+user+fromtime+totime , json , 1000000000 * 300).Err()
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
// @Summary user ccu 
// @Schemes
// @Description do user ccu 
// @Tags User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Ccu
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/ccu/{node}/{user}/{fromtime}/{totime} [get]
func CCuUser(c * gin.Context){
	node := c.Param("node")
	user := c.Param("user") 
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
	// _= rdb.Del(ctx, "daily_"+user+fromtime+totime ).Err()

    val, err := rdb.Get(ctx, "ccu"+node+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt := "SELECT time,ROUND(sum(requests)/300) FROM  access_"+s+ " WHERE  %s  http_host IN ("+search+") AND  UNIX_TIMESTAMP(time) >='"+fromtime+"' AND  UNIX_TIMESTAMP(time) <='"+totime+"' GROUP BY time ORDER BY time ASC"
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		// fmt.Println(queryinsrt)
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
			if err := rows.Scan( &detail.Time, &detail.Ccu ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx,"ccu"+node+user+fromtime+totime , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("daily user hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary time referer isp
// @Schemes
// @Description count time referer isp
// @Tags User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Referer
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/referer/{node}/{user}/{fromtime}/{totime} [get]
func RefererUserISP(c * gin.Context){
	node :=c.Param("node")
	user :=c.Param("user")
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
	// _= rdb.Del(ctx, "refererUser"+user+fromtime+totime+node ).Err()

    val, err := rdb.Get(ctx, "refererUser"+user+fromtime+totime+node).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		
		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt := "SELECT date(time),name,sum(data) FROM  %s WHERE %s name!='-' AND  http_host IN  ("+search+") AND  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY name ORDER BY sum(data)  DESC LIMIT 100"
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
		
		
		// var detail Detail
		details := make([]*Referer,0)
		for rows.Next(){
			
			detail := new(Referer)
			if err := rows.Scan( &detail.Time, &detail.Referer, &detail.Count ); err != nil {
				fmt.Println(err)
			}       
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "refererUser"+user+fromtime+totime+node , json , 1000000000 * 86400).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("refererUser hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary user summary data isp 
// @Schemes
// @Description do user summary data isp
// @Tags User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Detail
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Router /log/user/today/isp/{node}/{user}/{fromtime}/{totime} [get]
func SummaryUserISP(c * gin.Context){
	node :=c.Param("node")
	user := c.Param("user") 

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
	// _= rdb.Del(ctx, "today"+user ).Err()

    val, err := rdb.Get(ctx, "SummaryUserISP"+node+user+fromtime+totime).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt := "SELECT date(time),sum(byte),sum(requests),sum(byteVN),COALESCE(sum(requestVN),0),sum(cache_hit),sum(Cache_miss),sum(status_2x),sum(status_3x),sum(status_4x),sum(status_5x),sum(macintosh),sum(windows),sum(ios),sum(android),sum(smartv),sum(other) FROM  access_"+s+ " WHERE  %s http_host IN ("+search+") AND  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"')"  
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
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
		errs := rdb.Set(ctx,"SummaryUserISP"+node+user+fromtime+totime , json , 1000000000 * 86400).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly user hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary user Avgspeed isp 
// @Schemes
// @Description Avgspeed isp
// @Tags User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.Avgspeeds
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Router /log/user/avgspeed/isp/{node}/{user}/{unit}/{fromtime}/{totime} [get]
func AvgspeedUserISP(c * gin.Context){
	node :=c.Param("node")
	user := c.Param("user") 
	unit := c.Param("unit")

	fromtime :=c.Param("fromtime") 
	totime :=c.Param("totime") 
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
	// _= rdb.Del(ctx, "today"+user ).Err()

    val, err := rdb.Get(ctx, "AvgspeedUserISP"+node+user+fromtime+totime+unit).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt := "SELECT time,sum(byte),sum(byteVN),(sum(byte)-sum(byteVN)),sum(if(byte>0,requests,0)),sum(if(byte>0,requestVN,0)),(sum(if(byte>0,requests,0))-sum(if(byte>0,requestVN,0))) FROM access_"+s+ " WHERE  %s http_host IN ("+search+") AND  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"')  GROUP BY time ORDER BY time ASC"  
		var nodestr string
		if node == "total" {
			nodestr="nodeISP !='' AND"
		}else{
			nodestr="nodeISP='"+node+"' AND"
		}
		queryinsrt = fmt.Sprintf(queryinsrt , nodestr)
		// queryinsrt = strings.ReplaceAll(queryinsrt, "stringreplace", nodestr)
		rows, err := db.Query(queryinsrt)
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
		errs := rdb.Set(ctx,"AvgspeedUserISP"+node+user+fromtime+totime+unit , json , 1000000000 * 300).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("monly user hit cache")
	}
	
}

// @BasePath /api/v1

// log godoc
// @Summary user http error code isp
// @Schemes
// @Description user http error code
// @Tags User Node
// @Accept json
// @Produce json
// @Header 200 {string} Authorization: Basic 
// @success 200 {object} []controllers.TopHttpCodeError
// @Param user path string true "vtvcab"
// @Param node path string true "viettel"
// @Param code path string true "404x"
// @Param fromtime path int true "1636633800"
// @Param totime path int true "1637810200"
// @Router /log/user/codeerr/{node}/{user}/{code}/{fromtime}/{totime} [get]
func TopHttpCodeErrorUserISP(c * gin.Context){
	node :=c.Param("node")
	user :=c.Param("user")
	code :=c.Param("code")
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
	// _= rdb.Del(ctx, "TopHttpCodeError"+user+fromtime+totime+node+code ).Err()

    val, err := rdb.Get(ctx, "TopHttpCodeError"+user+fromtime+totime+node+code).Result()
    if err != nil {
		fmt.Println(err)
		db := database.DBcon() //Khởi tạo biến conection

		
		rowsDomain,err := db.Query( "SELECT http_host FROM http_host WHERE customerName='"+user+"'")
		if err != nil {
			fmt.Println(err.Error())
		}
		var text string
		for rowsDomain.Next(){
			var a string
			if err := rowsDomain.Scan( &a ); err != nil {
				fmt.Println(err)
			}    
			text += "'"+a+"'," 
		}
		search := strings.TrimRight(text, ",")
		queryinsrt := "SELECT http_host,date(time),sum("+code+") as count FROM  %s WHERE %s http_host IN  ("+search+") AND  time >=FROM_UNIXTIME('"+fromtime+"') AND time <=FROM_UNIXTIME('"+totime+"') GROUP BY http_host HAVING count > 0  ORDER BY sum("+code+") DESC  LIMIT 20"
		var nodestr string
		var table string
		// fmt.Println(queryinsrt)
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
			detail.HttpCodeError = code     //http_code_err_
			details = append(details, detail)

		}

		json , err := json.Marshal(details)
		if err != nil {
			fmt.Println(err)
		}
		errs := rdb.Set(ctx, "TopHttpCodeError"+user+fromtime+totime+node+code , json , 1000000000 * 600).Err()
		if errs != nil {
			fmt.Println(err)
		}

		c.JSON(200, details)
    }else{
		c.String(200, val)
		fmt.Println("TopHttpCodeError hit cache")
	}
	
}