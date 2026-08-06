package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
  apiKey := flag.String("api-key","","Key used for add new storages. If not set, anyone can add new storages")
  port := flag.Int("port",8080,"Server port")
  aiKey := flag.String("ai-key","","Key")
  aiUrl := flag.String("ai-url","","URL")
  aiPrompt := flag.String("ai-prompt","","")
  aiModel := flag.String("ai-model","","")
  prompt := defaultPrompt
  if *aiPrompt != "" {
    prompt = *aiPrompt
  }
  flag.Parse()

  router := gin.New()

  router.Use(gin.Logger())
  router.Use(gin.Recovery())
  conn,err := InitDB("sqlite.db")
  if err != nil {
    panic("failed to connect database: " + err.Error())
  }

  router.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })

  router.GET("/items",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "items": "items",
    })
  })
  router.GET("/items/search",func (ctx *gin.Context)  {
    query := ctx.Query("query")
    page := ctx.Query("page")
    pageSize := ctx.DefaultQuery("page_size","10")
    if strings.TrimSpace(query) == "" {
      ctx.JSON(400, gin.H{"error": "query is required"})
      return
    }
    intPageSize, err := strconv.Atoi(pageSize)
    if err != nil || intPageSize < 1 {
      intPageSize = 10
    }
    intPage, err := strconv.Atoi(page)
    if err != nil || intPage < 1 {
      intPage = 1
    }
    intPageSize = min(intPageSize,100)
    response, err := getItemsSnapshotByQuery(conn,query,intPage,intPageSize)
    if err != nil{
      ctx.Error(err)
      ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
      return
    }
    ctx.JSON(200,response)
  })
  router.GET("/types",func(ctx *gin.Context) {
    items_types, err := getTypes(conn)
    if err != nil{
      ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
      return
    }
    ctx.JSON(200,gin.H{
      "types": items_types,
    })
  })

  router.GET("/types/:id",func(ctx *gin.Context) {
    typeId := ctx.Param("id")
    intType,err := strconv.Atoi(typeId)
    if err != nil {
      ctx.JSON(http.StatusBadRequest,gin.H{"error": "invalid type id"})
      return
    }
    items, err := getItemsByType(conn,intType)
    if err != nil{
      ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
      return
    }
    ctx.JSON(200,gin.H{
      "items": items,
    })
  })

  router.GET("/barrels",func(ctx *gin.Context) {
    query := ctx.Query("query")
    page := ctx.Query("page")
    pageSize := ctx.DefaultQuery("page_size","10")
    if strings.TrimSpace(query) == "" {
      ctx.JSON(400, gin.H{"error": "query is required"})
      return
    }
    intPageSize, err := strconv.Atoi(pageSize)
    if err != nil || intPageSize < 1 {
      intPageSize = 10
    }
    intPage, err := strconv.Atoi(page)
    if err != nil || intPage < 1 {
      intPage = 1
    }
    intPageSize = min(intPageSize,100)
    response, err := getBarrelsByQuery(conn,query,intPage,intPageSize)
    if err != nil{
      ctx.Error(err)
      ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
      return
    }
    ctx.JSON(200,response)
  })

  router.GET("/barrels/:id",func(ctx *gin.Context) {
    id := ctx.Param("id")
    intId, err := strconv.Atoi(id)
    if err != nil {
      ctx.Error(err)
      ctx.JSON(http.StatusBadRequest,gin.H{"error": "invalid barrel id"})
      return
    }
    info, err := getBarrelInfo(conn,intId)
    if err != nil {
      ctx.Error(err)
      ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
      return
    }
    if info == nil {
      ctx.Error(err)
      ctx.JSON(http.StatusNotFound,gin.H{"error": "barrel not found"})
      return
    }
    ctx.JSON(http.StatusOK,info)
  })

  router.GET("/barrels/:id/items",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "items": "items",
    })
  })

  router.GET("/seller/:id",func(ctx *gin.Context) {
    sellerId := ctx.Param("id")
    intId,err := strconv.Atoi(sellerId)
    if err != nil {
      ctx.JSON(http.StatusBadRequest,gin.H{"error": "invalid seller id"})
      return
    }
    items, err := getBarrelsBySellerId(conn,int32(intId))
    if err != nil{
      ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
      return
    }
    ctx.JSON(200,gin.H{
      "items": items,
    })
  })

  router.POST("/barrels/items",func(ctx *gin.Context) {
    var request ItemInBarrelPost
    if err := ctx.ShouldBindJSON(&request);err != nil{
      ctx.Error(err)
      ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		  return
    }
    err = postItemsInBarrel(conn,request)
    if err != nil {
      ctx.Error(err)
      if errors.Is(err,ErrNotFound){
        ctx.JSON(http.StatusNotFound,gin.H{"error": "Barrel not found"})
        return
      }
      ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
      return
    }
    ctx.JSON(http.StatusCreated,gin.H{
      "status": "success",
    })
  })

  router.POST("/barrels",func(ctx *gin.Context) {
    auth := ctx.GetHeader("Authorization")
    expected := "Bearer " + *apiKey
    if *apiKey != "" && auth != expected {
      ctx.JSON(http.StatusUnauthorized, gin.H{
          "error": "Invalid API key",
      })
      return
    }
    var request []NewBarrelPost
    if err := ctx.ShouldBindJSON(&request); err != nil {
      ctx.Error(err)
      ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		  return
    }
    var recognizedBarrel *RecognizedBarrelItem
    for _, barrel := range request {
      splittedMessage := strings.Split(barrel.Message,"\n")
      if len(splittedMessage) == 3{
        recognizedBarrel.ItemName = splittedMessage[0]
        benefitParts := strings.Split(splittedMessage[1],"-")
        recognizedBarrel.SellerName = splittedMessage[2]
        re := regexp.MustCompile("[^0-9]")
        recognizedBarrel.Quantity, err = strconv.ParseFloat(re.ReplaceAllString(benefitParts[0],""),32)
        recognizedBarrel.Price, err = strconv.ParseFloat(re.ReplaceAllString(benefitParts[1],""),32)
        recognizedBarrel.BenefitRatio = recognizedBarrel.Quantity / recognizedBarrel.Price
      } else if *aiKey != "" && *aiUrl != "" && *aiModel != "" {
        recognizedBarrel, err = aiBarrelRecognition(*aiKey,*aiUrl,*aiModel,barrel.Message,prompt)
        if err != nil {
          ctx.Error(err)
          return
        }
      }
      sellerId, err := getSellerByName(conn,recognizedBarrel.SellerName)
      if sellerId == nil {
        id, err := getUserMinecraftUUID(recognizedBarrel.SellerName)
        if err != nil {
          ctx.Error(err)
          return
        }
        sellerId, err = createSeller(conn,recognizedBarrel.SellerName,&id)
      }
      barrelId, err := getBarrelByCords(conn,*barrel.X,*barrel.Y,*barrel.Z)
      if barrelId == nil {
        barrelId, err = createBarrel(conn,*barrel.X,*barrel.Y,*barrel.Z)
        if err != nil {
          ctx.Error(err)
          return
        }
      }
      itemId,err := getOrCreateItem(conn,recognizedBarrel.ItemName)
      _, err = createBarrelItem(conn,*itemId,*barrelId,*sellerId,*recognizedBarrel,barrel.Message)
      if err != nil {
        ctx.Error(err)
      }
    }
    ctx.JSON(200,gin.H{
      "status": "success",
    })
  })
  addr := fmt.Sprintf(":%d",*port)
  router.Run(addr)
}