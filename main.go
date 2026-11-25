package main

import (
	"blockberg-terminal/telemetry"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
  isOnlyAI := flag.Bool("only-ai",false,"")
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

  ctx := context.Background()

  shutdown, err := telemetry.Setup(
      ctx,
      "barrel-api",
      "localhost:4317",
  )
  if err != nil {
      log.Fatal(err)
  }

  defer func() {
      ctx, cancel := context.WithTimeout(
          context.Background(),
          5*time.Second,
      )
      defer cancel()

      if err := shutdown(ctx); err != nil {
          log.Printf("Telemetry shutdown error: %v", err)
      }
  }()
  router := gin.New()

  router.Use(gin.Logger())
  router.Use(gin.Recovery())
  router.Use(func(ctx *gin.Context) {
    ctx.Next()

    for _, err := range ctx.Errors {
        log.Println("ERROR:", err.Err)
    }
  })
  router.Use(otelgin.Middleware("barrel-api"))
  pool,err := InitDB("sqlite.db")
  if err != nil {
    panic("failed to connect database: " + err.Error())
  }

  router.GET("/ping", ping)

  router.GET("/items",getItems())

  router.GET("/snapshots",getSnapshots(pool))

  router.GET("/types",getAllTypes(pool))

  router.GET("/types/:id",getType(pool))

  router.GET("/barrels",getBarrels(pool))

  router.GET("/barrels/:id",getBarrel(pool))

  router.GET("/barrels/:id/items",getBarrelItems(pool))

  router.GET("/seller/:id",getSeller(pool))

  router.POST("/barrels/items",postBarrelItems(pool))

  router.POST("/barrels",postBarrels(pool,apiKey,xMin,xMax,yMin,yMax,zMin,zMax,prompt,aiKey,aiModel,aiUrl,*isOnlyAI))

  router.GET("/stats/barrels",GetAllBarrels(pool))

  router.GET("/stats/sellers/top",GetTopSellers(pool))

  router.GET("/stats/items/stored",GetTopStoredItems(pool))

  router.GET("/stats/items/barrels",GetTopSellingBarrelItems(pool))

  router.GET("/stats/items/price",GetBarrelsMaxPrices(pool))

  router.GET("/stats/items/benefit",GetBarrelsMaxBenefit(pool))
  
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