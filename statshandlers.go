package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"zombiezen.com/go/sqlite"
)

func GetAllBarrels(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		barrels, err := getBarrelsWithCords(conn)
		if err != nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"total": len(barrels), "barrels": barrels})
	}
}

func GetTopSellers(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		sellers,err := getTopSellers(conn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"sellers": sellers}) 
	}
}

func GetBarrelsMaxPrices(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		barrels,err := getBarrelsMaxPrices(conn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"barrels": barrels}) 
	}
}

func GetBarrelsMaxBenefit(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		barrels,err := getBarrelsMaxBenefit(conn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"barrels": barrels}) 
	}
}

func GetTopSellingBarrelItems(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		items,err := topSellingBarrelItems(conn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"items": items}) 
	}
}

func GetTopStoredItems(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		items,err := topStoredItems(conn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"items": items}) 
	}
}