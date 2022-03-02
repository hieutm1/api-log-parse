package main

import (
        "github.com/gin-gonic/gin"
        "logsquery/controllers"
		swaggerfiles "github.com/swaggo/files"
		ginSwagger "github.com/swaggo/gin-swagger"
		// "logsquery/db"
		// "log"
		// "context"
		// rdb "github.com/go-redis/redis/v8"
		// "database/sql"
		// _ "github.com/go-sql-driver/mysql"
		"fmt"
		"time"
		docs "logsquery/docs"
		"net/http"
		"github.com/gin-contrib/cors"
		"flag"
)

var (
	port		*string
	host		*string
	keyaccess	*string
	keysecret	*string
)

func init() {
	port 		= flag.String("port","8081", "read port")
	host 		= flag.String("host","", "read host")
	keyaccess 	= flag.String("keyaccess","9LAHFgjhKQLTUYfsMYfh", "read keyaccess")
	keysecret 	= flag.String("keysecret","UthRjdorucYw0uTIinLdSQbP09Sk1O", "read keysecret")
}

// @securityDefinitions.basic BasicAuth
// @name Authorization
func setupRouter() *gin.Engine {
        r := gin.Default()
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"PUT", "PATCH","GET"},
			AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return origin == "*"
			},
			MaxAge: 12 * time.Hour,
			}))
		r.Use(CORSMiddleware())
		docs.SwaggerInfo.Title = "LOG API"
		docs.SwaggerInfo.Description = "This is a server API LOGS server."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = "113.164.27.23:8081"
		docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
        r.Static("/public", "./public")
        r.LoadHTMLGlob("templates/*")

        v1 := r.Group("/api/v1") 
        {
			keyaccess := fmt.Sprintf("%s", *keyaccess)
			keysecret := fmt.Sprintf("%s", *keysecret)
			log := v1.Group("/log",gin.BasicAuth(gin.Accounts{
				// "9LAHFgjhKQLTUYfsMYfh": "UthRjdorucYw0uTIinLdSQbP09Sk1O",
				keyaccess : keysecret,
			}))
			{
				log.GET("/domains", controllers.Domains)
				log.GET("/domain/:hostname/:fromtime/:totime", controllers.DomainDetail)
				log.GET("/daily/:hostname/:fromtime/:totime", controllers.DailyDetail)
				log.GET("/city/:hostname/:fromtime/:totime", controllers.DailyCity)
				log.GET("/monthly/:hostname/:fromtime/:totime", controllers.MonthlyDetail)
				log.GET("/thismonth/:hostname", controllers.ThisMonth)
				log.GET("/today/:hostname", controllers.TodaySum)
				log.GET("/summarydata/:hostname/:fromtime/:totime", controllers.SumaryData)
				log.GET("/avgspeed/:hostname/:unit/:fromtime/:totime", controllers.AvgSpeed)
				log.GET("/totalByte", controllers.TotalByte)
				log.GET("/bandwidth/:domain/:unit/:fromtime/:totime", controllers.BandwidthDomainAll)
				log.GET("/bandwidthAll/:domain/:unit/:fromtime/:totime", controllers.BandwidthDomainAll)
				log.GET("/bytefmt/:byte", controllers.Byteto)
				log.GET("/country/:hostname/:fromtime/:totime", controllers.Country)
				log.GET("/isp/:hostname/:fromtime/:totime", controllers.ISP)
				log.GET("/responseTime/:hostname/:fromtime/:totime", controllers.ResponseTime)
				log.GET("/responseTimeDaily/:hostname/:fromtime/:totime", controllers.ResponseTimeDaily)
				log.GET("/ccu/:hostname/:fromtime/:totime", controllers.CCUdomain)
				log.GET("/referer/:hostname/:fromtime/:totime", controllers.RefererDomain)
				log.GET("/codeerr/:hostname/:fromtime/:totime", controllers.HttpCodeErr)
				// user
				log.GET("/user/domains/:user", controllers.DomainUser)
				log.GET("/user/today/:user", controllers.TodayUser)
				log.GET("/user/thismonth/:user", controllers.ThisMonthUser)
				log.GET("/user/min/:user/:fromtime/:totime", controllers.MinUser)
				log.GET("/user/daily/:user/:fromtime/:totime", controllers.DailyUser)
				log.GET("/user/monthly/:user", controllers.MonthlyhUser)
				log.GET("/user/country/:user/:fromtime/:totime", controllers.CountryUser)
				log.GET("/user/isp/:user/:fromtime/:totime", controllers.ISPUser)
				log.GET("/user/city/:user/:fromtime/:totime", controllers.CityUser)
				log.GET("/user/responseTime/:user/:fromtime/:totime", controllers.ResponseTimeUser)
				log.GET("/user/bandwidth/:user/:unit/:fromtime/:totime", controllers.BandwidthUser)
				log.GET("/user/ccu/:node/:user/:fromtime/:totime", controllers.CCuUser)
				//user isp
				log.GET("/user/today/isp/:node/:user", controllers.TodayUserISP)
				log.GET("/user/thismonth/isp/:node/:user", controllers.ThisMonthUserISP)
				log.GET("/user/min/isp/:node/:user/:fromtime/:totime", controllers.MinUserISP)
				log.GET("/user/daily/isp/:node/:user/:fromtime/:totime", controllers.DailyUserISP)
				log.GET("/user/monthly/isp/:node/:user", controllers.MonthlyhUserISP)
				log.GET("/user/country/isp/:node/:user/:fromtime/:totime", controllers.CountryUserISP)
				log.GET("/user/isp/node/:node/:user/:fromtime/:totime", controllers.ISPUserISP)
				log.GET("/user/city/isp/:node/:user/:fromtime/:totime", controllers.CityUserISP)
				log.GET("/user/responseTime/isp/:node/:user/:fromtime/:totime", controllers.ResponseTimeUserISP)
				log.GET("/user/bandwidth/isp/:node/:user/:unit/:fromtime/:totime", controllers.BandwidthUserISP)
				log.GET("/user/referer/isp/:node/:user/:fromtime/:totime", controllers.RefererUserISP)
				log.GET("/user/summary/isp/:node/:user/:fromtime/:totime", controllers.SummaryUserISP)
				log.GET("/user/avgspeed/isp/:node/:user/:unit/:fromtime/:totime", controllers.AvgspeedUserISP)
				log.GET("/user/codeerr/isp/:code/:node/:user/:fromtime/:totime", controllers.TopHttpCodeErrorUserISP)
				// log.PUT("/user/uploadUserDomain", controllers.EditUserDomain)
				log.GET("/total/today", controllers.Todaytotal)
				log.GET("/total/thismonth", controllers.ThisMonthTotal)
				log.GET("/total/min/:fromtime/:totime", controllers.Mintotal)
				log.GET("/total/daily/:fromtime/:totime", controllers.Dailytotal)
				log.GET("/total/monthly", controllers.Monthlyhtotal)
				log.GET("/total/responseTime/:fromtime/:totime", controllers.ResponseTimeTotal)
				log.GET("/total/country/:fromtime/:totime", controllers.CountryTotal)
				log.GET("/total/isp/:fromtime/:totime", controllers.ISPTotal)
				log.GET("/total/city/:fromtime/:totime", controllers.CityTotal)
				log.GET("/total/bandwidth/:unit/:fromtime/:totime", controllers.BandwidthTotal)
				// // total iSP
				log.GET("/total/today/:node", controllers.TodaytotalISP)
				log.GET("/total/thismonth/:node", controllers.ThisMonthTotalISP)
				log.GET("/total/min/isp/:node/:fromtime/:totime", controllers.MintotalISP)
				log.GET("/total/daily/isp/:node/:fromtime/:totime", controllers.DailytotalISP)
				log.GET("/total/monthly/isp/:node", controllers.MonthlyhtotalISP)
				log.GET("/total/responseTime/isp/:node/:fromtime/:totime", controllers.ResponseTimeTotalISP)
				log.GET("/total/country/isp/:node/:fromtime/:totime", controllers.CountryTotalISP)
				log.GET("/total/isp/node/:node/:fromtime/:totime", controllers.ISPTotalISP)
				log.GET("/total/city/isp/:node/:fromtime/:totime", controllers.CityTotalISP)
				log.GET("/total/bandwidth/isp/:node/:unit/:fromtime/:totime", controllers.BandwidthTotalISP)
				log.GET("/total/ccu/isp/:node/:fromtime/:totime", controllers.CCuTotal)
				log.GET("/total/referer/isp/:node/:fromtime/:totime", controllers.RefererTotalISP)
				log.GET("/total/summary/:node/:fromtime/:totime", controllers.SummaryTotalISP)
				log.GET("/total/avgspeed/:node/:unit/:fromtime/:totime", controllers.TotalAvgspeed)
				log.GET("/total/codeerr/:code/:node/:fromtime/:totime", controllers.TopHttpCodeErrorTotalISP)
			}
			v1.GET("/secret", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"secret": "The secret ingredient to the BBQ sauce is stiring it in an old whiskey barrel.",
				})
			})
			// client.POST("/backend/push", controllers.Push)
			// client.PATCH("/story/update/:id", controllers.Update)
			// client.DELETE("/story/:id", controllers.Delete)
			//user
        }
		
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
        return r
}


func main() {
		// programmatically set swagger info
        r := setupRouter()
		listen := fmt.Sprintf("%s:%s", *host, *port)
        r.Run(listen) // Ứng dụng chạy tại cổng 8081
			
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}