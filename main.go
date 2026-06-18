package main

import "github.com/gin-gonic/gin"

func main() {
  router := gin.Default()
  router.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })
  router.GET("items",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "items": "items",
    })
  })
  router.GET("type",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "type": "type",
    })
  })
  router.GET("barrel",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "barrel": "barrel",
    })
  })
  router.GET("seller",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "seller": "seller",
    })
  })

  router.POST("item",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "item": "item",
    })
  })
  router.POST("type",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "type": "type",
    })
  })
  router.POST("barrel",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "barrel": "barrel",
    })
  })
  router.POST("seller",func(ctx *gin.Context) {
    ctx.JSON(200,gin.H{
      "seller": "seller",
    })
  })
  router.Run()
}