package main

import (
	"go-web/router"
	"go-web/swagger"
	"go-web/utils"
)

func main() {
	utils.InitMysql()
	utils.InitRedis()
	utils.InitConfig()
	utils.InitWorkerPool()
	r := router.Router()
	swagger.SwaggerInit(r)

	r.Run(":8081")
}
