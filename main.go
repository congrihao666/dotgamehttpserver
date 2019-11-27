package main

import (
	"github.com/gin-gonic/gin"
	////"fmt"
	"net/http"
	"encoding/json"
	////"strings"
)

func cors() gin.HandlerFunc {
    return func(c *gin.Context) {
        method := c.Request.Method

        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
        c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
        c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
        c.Header("Access-Control-Allow-Credentials", "true")
        if method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
        }
        c.Next()
    }
}

func handleDotGame(c *gin.Context) {
	raw_str,_ := c.GetRawData()
	var param Param
	json.Unmarshal(raw_str,&param)
	//////fmt.Println("gameID,gameName",param.GameID,param.GameName)

	if param.Data != "" && param.Channel != "" {
		insert_dot_log(&param)
	} 
	
	/*
	func_list := strings.Split(param.Data,"&")
	
	for _,funcData := range func_list {
		var unFunc FuncItem
		json.Unmarshal([]byte(funcData),&unFunc)
		fmt.Println("src====",funcData)
		////if unFunc.Param1 != "" {
		/////	fmt.Println("funcName=========",unFunc.FuncName,unFunc.TimeStamp,unFunc.Param1)
		////}

	}
	*/

}

func main() {
	/////TestMongo()
	/////go init_mongo()
	r := gin.Default()
	r.Use(cors())
	r.POST("/",handleDotGame)
	r.Run(":8723")
}