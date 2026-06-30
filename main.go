package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
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
    err := postBarrelItems(conn,request)
    if err != nil {
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