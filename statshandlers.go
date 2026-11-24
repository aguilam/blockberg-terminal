package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"zombiezen.com/go/sqlite/sqlitex"
)

func GetAllBarrels(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		barrels, err := getBarrelsWithCords(conn)
		if err != nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"total": len(barrels), "barrels": barrels})
	}
}

func GetTopSellers(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		sellers,err := getTopSellers(conn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"sellers": sellers}) 
	}
}

func GetBarrelsMaxPrices(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		barrels,err := getBarrelsMaxPrices(conn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"barrels": barrels}) 
	}
}

func GetBarrelsMaxBenefit(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		barrels,err := getBarrelsMaxBenefit(conn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"barrels": barrels}) 
	}
}

func GetTopSellingBarrelItems(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		items,err := topSellingBarrelItems(conn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"items": items}) 
	}
}

func GetTopStoredItems(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		items,err := topStoredItems(conn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"items": items}) 
	}
}