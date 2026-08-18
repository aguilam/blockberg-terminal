package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
)
const unsetCoord = 65531
func main() {
  apiKey := flag.String("api-key","","Key used for add new storages. If not set, anyone can add new storages")
  port := flag.Int("port",8080,"Server port")
  aiKey := flag.String("ai-key","","AI provider API key")
  aiUrl := flag.String("ai-url","","AI provider base url")
  aiPrompt := flag.String("ai-prompt","","Replace base prompt for ai recognition")
  aiModel := flag.String("ai-model","","Selected AI model name")
  xMin := flag.Int("x-min",unsetCoord,"")
  xMax := flag.Int("x-max",unsetCoord,"")
  yMin := flag.Int("y-min",unsetCoord,"")
  yMax := flag.Int("y-max",unsetCoord,"")
  zMin := flag.Int("z-min",unsetCoord,"")
  zMax := flag.Int("z-max",unsetCoord,"")

  flag.Parse()
  prompt := defaultPrompt
  if *aiPrompt != "" {
    prompt = *aiPrompt
  }

  xMin = nilCordFlag(xMin)
  xMax = nilCordFlag(xMax)
  yMin = nilCordFlag(yMin)
  yMax = nilCordFlag(yMax)
  zMin = nilCordFlag(zMin)
  zMax = nilCordFlag(zMax)

  router := gin.New()

  router.Use(gin.Logger())
  router.Use(gin.Recovery())
  conn,err := InitDB("sqlite.db")
  if err != nil {
    panic("failed to connect database: " + err.Error())
  }

  router.GET("/ping", ping)

  router.GET("/items",getItems())

  router.GET("/snapshots",getSnapshots(conn))

  router.GET("/types",getAllTypes(conn))

  router.GET("/types/:id",getType(conn))

  router.GET("/barrels",getBarrels(conn))

  router.GET("/barrels/:id",getBarrel(conn))

  router.GET("/barrels/:id/items",getBarrelItems(conn))

  router.GET("/seller/:id",getSeller(conn))

  router.POST("/barrels/items",postBarrelItems(conn))

  router.POST("/barrels",postBarrels(conn,apiKey,xMin,xMax,yMin,yMax,zMin,zMax,prompt,aiKey,aiModel,aiUrl))

  router.Use(func(ctx *gin.Context) {
    ctx.Next()

    for _, err := range ctx.Errors {
        log.Println("ERROR:", err.Err)
    }
  })

  addr := fmt.Sprintf("127.0.0.1:%d", *port)
  listener, err := net.Listen("tcp", addr)
  if err != nil {
		log.Fatalf("Failed to bind port: %v", err)
	}

  actualAddr := listener.Addr().String()
  fmt.Printf("SERVER_READY:%s\n", actualAddr)
  if err := router.RunListener(listener); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}