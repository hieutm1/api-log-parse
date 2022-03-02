// database/mysql.go
package database

import (
	"database/sql"
	"time"
	// "log"
	"fmt"
	// "github.com/gocolly/colly"
	_ "github.com/go-sql-driver/mysql"
	"path/filepath"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"flag"
)

type Config struct {
    Mysql  	struct {
		Username  string `yaml:"username"`
		Pass	  string `yaml:"pass"`
		Dbname    string `yaml:"dbname"`
		Host  	  string `yaml:"host"`
		Port	  string `yaml:"port"`
	}
	Redis struct {
		Host  	  string `yaml:"host"`
		Port	  string `yaml:"port"`
	}	
}

var (
	cfg   *string
	port   *string
)

const (
	USERNAME = "sbd_parser_user"
	PASS = "SBD@123!@#"
	DB_NAME = "sdb_log_db"
	DB_DRIVER = "mysql"
	CHARSET = "utf8"
)
var db *sql.DB

func init() {
    cfg 	= flag.String("cfg", "", "read file config")
}

func DBcon() (db *sql.DB){
	cnf := ConfigFile()
	// dsn := fmt.Sprintf("%s:%s@tcp(172.40.4.85:3306)/%s?charset=%s", USERNAME, PASS, DB_NAME, CHARSET)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", cnf.Mysql.Username, cnf.Mysql.Pass, cnf.Mysql.Host, cnf.Mysql.Port, cnf.Mysql.Dbname, CHARSET)
	db, err := sql.Open(DB_DRIVER, dsn)
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxOpenConns(20)
    db.SetMaxIdleConns(20)
    db.SetConnMaxLifetime(time.Minute * 5)
    return db 
}

func CloseDatabase(db *sql.DB){
	defer db.Close()
}

func TEST(e string){
	fmt.Println(e)
}

func ConfigFile() *Config{
	flag.Parse()
	var config Config
	filename, _ := filepath.Abs(*cfg)
	err := doesFileExist(*cfg)
	if err != nil {
        log.Error(err)
    }else{
		yamlFile, err := ioutil.ReadFile(filename)
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Error(err)
		}
	}
	return &config
}

// function to check if file exists
func doesFileExist(fileName string) (err error){
	_ , error := os.Stat(fileName)
	  return error
}
