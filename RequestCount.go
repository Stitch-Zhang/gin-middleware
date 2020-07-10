package RequestCount

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
)

//RCJSON is the struct of log file
type RCJSON struct {
	TotalTimes int            `json:"totalTimes"`
	MaxReqIP   string         `json:"maxReqIP"`
	Details    map[string]int `json:"details"`
}

//Loger is collect website log to file
var Loger = func() gin.HandlerFunc {
	return func(c *gin.Context) {
		logs := struct {
			Time   time.Time `json:"time"`
			IP     string    `json:"ip"`
			Method string    `json:"method"`
			URL    string    `json:"url"`
			UA     string    `json:"ua"`
		}{
			Time:   time.Now(),
			IP:     c.ClientIP(),
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			UA:     c.Request.UserAgent(),
		}
		//Add more details to url of logs
		if c.Request.URL.RawQuery != "" {
			logs.URL += "?" + c.Request.URL.RawQuery
		}
		logerJson, err := json.Marshal(logs)
		if err != nil {
			fmt.Println("logs to Json wrong!\n", err)
			return
		}
		logsFile, err := os.OpenFile("logs.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		_, err = logsFile.Write(logerJson)
		logsFile.WriteString("\n")
		if err != nil {
			fmt.Println("Writing logs wrong!", err)
			return
		}
		fmt.Println(string(logerJson))
		defer func() {
			logsFile.Close()
		}()
	}
}

//RCMAP is for count requests times
type RCMAP map[string]int

func ExitAction() {
	stopTime := time.Now().Format("01-02-15:04:05")
	log.Print("Ready to exit\n")
	err := os.Rename("RC.log", "Request-Times---"+stopTime+".log")
	if err != nil {
		fmt.Println("Renamed RC.log wrong!", err)
	}
	err = os.Rename("logs.json", "Logs---"+stopTime+".log")
	if err != nil {
		fmt.Println("Renamed logs.json wrong!", err)
	}
	log.Print("Bye!")
	os.Exit(0)
}

func GetMapMaxKey(targetMap map[string]int) string {
	maxOne := struct {
		key   string
		value int
	}{}
	for k, v := range targetMap {
		if maxOne.value < v {
			maxOne.key = k
			maxOne.value = v
		}
	}
	return maxOne.key
}

/* CORE FUNCTION
RC := func(c *gin.Context) {
	rcJSON.TotalTimes += 1
	rcJSON.MaxReqIP = c.ClientIP()
	rcMap[c.ClientIP()] += 1
	rcJSON.Details = rcMap
	func() {
		rcFile, err := os.Create("RC.log")
		rcJSONer, err := json.Marshal(rcJSON)
		if err != nil {
			fmt.Println("Created RC.log failed!", err)
			return
		}
		_, err = rcFile.Write(rcJSONer)
		if err != nil {
			fmt.Println("Wrote rclog wrong!", err)
			return
		}
		defer rcFile.Close()
	}()
}*/
