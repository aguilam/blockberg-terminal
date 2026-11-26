package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"zombiezen.com/go/sqlite/sqlitex"
)
var handlerTracer = otel.Tracer("blockberg-terminal/stats/handlers")

func GetAllBarrels(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := handlerTracer.Start(ctx.Request.Context(),"stats.barrels.all.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		barrels, err := getBarrelsWithCords(conn,ctx)
		if err != nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"total": len(barrels), "barrels": barrels})
	}
}

func GetTopSellers(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := handlerTracer.Start(ctx.Request.Context(),"stats.sellers.top.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		sellers,err := getTopSellers(conn,ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"sellers": sellers}) 
	}
}

func GetBarrelsMaxPrices(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := handlerTracer.Start(ctx.Request.Context(),"stats.barrels.prices.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		barrels,err := getBarrelsMaxPrices(conn,ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"barrels": barrels}) 
	}
}

func GetBarrelsMaxBenefit(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := handlerTracer.Start(ctx.Request.Context(),"stats.barrels.benefits.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		barrels,err := getBarrelsMaxBenefit(conn,ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"barrels": barrels}) 
	}
}

func GetTopSellingBarrelItems(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := handlerTracer.Start(ctx.Request.Context(),"stats.barrels.items.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		items,err := topSellingBarrelItems(conn,ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"items": items}) 
	}
}

func GetTopStoredItems(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := handlerTracer.Start(ctx.Request.Context(),"stats.items.stored.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		items,err := topStoredItems(conn,ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		ctx.JSON(200,gin.H{"items": items}) 
	}
}