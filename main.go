package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
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
  router.GET("/types",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "type": "type",
    })
  })
  router.GET("/types/:id",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "type": "type",
    })
  })
  router.GET("/barrels",func(ctx *gin.Context) {
    query := ctx.Query("query")
    ctx.JSON(200,gin.H{
      "barrel": "barrel",
    })
  })
  router.GET("/barrels/:id",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "barrel": "barrel",
    })
  })
  router.GET("/seller/:id",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "seller": "seller",
    })
  })

  router.POST("/barrels/items",func(ctx *gin.Context) {
    var request ItemInBarrelPost
    if err := ctx.ShouldBindJSON(&request);err != nil{
      ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		  return
    }
    ctx.JSON(200,gin.H{
      "item": "item",
    })
  })
  router.POST("/barrels",func(ctx *gin.Context) {
    var request NewBarrelPost
    if err := ctx.ShouldBindJSON(&request);err != nil{
      ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		  return
    }
    ctx.JSON(200,gin.H{
      "barrel": "barrel",
    })
  })
  router.Run()
}