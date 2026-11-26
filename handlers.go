package main

import (
	"blockberg-terminal/telemetry"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"zombiezen.com/go/sqlite/sqlitex"
)

var hTraces = otel.Tracer("blockberg-terminal/handlers")

func ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "ok",
		"version": "1.0.0",
	})
}

func getItems() (func(*gin.Context)) {
	return func (ctx *gin.Context)  {
		ctx.JSON(200,gin.H{
			"items": "items",
		})
	}
}

func getSnapshots(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := hTraces.Start(ctx.Request.Context(),"snapshots.search.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
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
		response, err := getItemsSnapshotByQuery(conn,query,intPage,intPageSize,ctx)
		if err != nil{
		  ctx.Error(err)
		  ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
		  return
		}
		ctx.JSON(200,response)
	}
}

func getAllTypes(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func (ctx *gin.Context)  {
		_,span := hTraces.Start(ctx.Request.Context(),"types.all.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		items_types, err := getTypes(conn,ctx)
		if err != nil{
		  ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
		  return
		}
		ctx.JSON(200,gin.H{
		  "types": items_types,
		})
	}
}

func getType(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := hTraces.Start(ctx.Request.Context(),"types.one.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		typeId := ctx.Param("id")
		intType,err := strconv.Atoi(typeId)
		if err != nil {
		  ctx.JSON(http.StatusBadRequest,gin.H{"error": "invalid type id"})
		  return
		}
		items, err := getItemsByType(conn,intType,ctx)
		if err != nil{
		  ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
		  return
		}
		ctx.JSON(200,gin.H{
		  "items": items,
		})
	}
}

func getBarrels(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := hTraces.Start(ctx.Request.Context(),"barrels.search.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
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
		response, err := getBarrelsByQuery(conn,query,intPage,intPageSize,ctx)
		if err != nil{
		  ctx.Error(err)
		  ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
		  return
		}
		ctx.JSON(200,response)
	}
}

func getBarrel(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := hTraces.Start(ctx.Request.Context(),"barrels.one.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		id := ctx.Param("id")
		intId, err := strconv.Atoi(id)
		if err != nil {
		  ctx.Error(err)
		  ctx.JSON(http.StatusBadRequest,gin.H{"error": "invalid barrel id"})
		  return
		}
		info, err := getBarrelInfo(conn,intId,ctx)
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
	}
}

func getBarrelItems(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
	    ctx.JSON(200,gin.H{
	      "items": "items",
	    })
	}
}

func getSeller(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := hTraces.Start(ctx.Request.Context(),"sellers.one.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		sellerId := ctx.Param("id")
		intId,err := strconv.Atoi(sellerId)
		if err != nil {
		  ctx.JSON(http.StatusBadRequest,gin.H{"error": "invalid seller id"})
		  return
		}
		items, err := getBarrelsBySellerId(conn,int32(intId),ctx)
		if err != nil{
		  ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
		  return
		}
		ctx.JSON(200,gin.H{
		  "items": items,
		})
	}
}

func postBarrelItems(pool *sqlitex.Pool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := hTraces.Start(ctx.Request.Context(),"snapshots.post.handlers")
		defer span.End()
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		var request ItemInBarrelPost
		if err := ctx.ShouldBindJSON(&request);err != nil{
		  ctx.Error(err)
		  ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			  return
		}
		err = postItemsInBarrel(conn,request,ctx)
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
	}
}

func postBarrels(pool *sqlitex.Pool,apiKey *string, xMin *int, xMax *int, yMin *int, yMax * int, zMin *int, zMax *int, prompt string, aiKey *string, aiModel *string, aiUrl *string, isOnlyAI bool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		_,span := hTraces.Start(ctx.Request.Context(),"barrels.post.handlers")
		defer span.End()
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
		telemetry.RecordsReceived.Add(ctx,int64(len(request)))
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		var barrels []AddBarrelStruct
		var success,aiRecognize,failed int64
		for _, barrel := range request {
			if (xMin != nil && *barrel.X < *xMin) || (xMax != nil && *barrel.X > *xMax) {
				continue
			}
			if (yMin != nil && *barrel.Y < *yMin) || (yMax != nil && *barrel.Y > *yMax) {
				continue
			}
			if (zMin != nil && *barrel.Z < *zMin) || (zMax != nil && *barrel.Z > *zMax) {
				continue
			}
			recognizedBarrel := &RecognizedBarrelItem{}
			splittedMessage := strings.Split(barrel.Message,"\n")
		  	if (isOnlyAI) {
				recBarrel, err := aiFallbackRecognizeBarrel(barrel.Message,*aiKey,*aiUrl,*aiModel,prompt,nil)
				if err != nil {
					ctx.Error(err)
					failed++
					continue
				}
				recognizedBarrel = recBarrel
				aiRecognize++
				goto adding
		  	}
		  	if len(splittedMessage) == 3{
				recognizedBarrel.ItemName = splittedMessage[0]
				benefitParts := strings.Split(strings.ReplaceAll(splittedMessage[1]," ",""), "-")
				recognizedBarrel.SellerName = splittedMessage[2]
				re := regexp.MustCompile("[^0-9]")
				quantity, err := strconv.ParseInt(re.ReplaceAllString(benefitParts[0],""),10,32)
				if err != nil {
					recognizedBarrel, err = aiFallbackRecognizeBarrel(barrel.Message,*aiKey,*aiUrl,*aiModel,prompt,err)
					if err != nil {
						ctx.Error(err)
						failed++
						continue
					}
					aiRecognize++
					goto adding
				}
				recognizedBarrel.Quantity = int(quantity);
				price,err := strconv.ParseInt(re.ReplaceAllString(benefitParts[1],""),10,32)
				if err != nil {
					recognizedBarrel, err = aiFallbackRecognizeBarrel(barrel.Message,*aiKey,*aiUrl,*aiModel,prompt,err)
					if err != nil {
						ctx.Error(err)
						failed++
						continue
					}
					aiRecognize++
					goto adding
				}
				recognizedBarrel.Price = int(price)
				recognizedBarrel.BenefitRatio = float64(recognizedBarrel.Quantity) / float64(recognizedBarrel.Price)
		  	} else if len(splittedMessage) == 4 {
				recognizedBarrel.ItemName = splittedMessage[0] + " " + splittedMessage[1]
				benefitParts := strings.Split(strings.ReplaceAll(splittedMessage[2]," ",""), "-")
				if len(benefitParts) != 2 {
					newBarrel, err := aiFallbackRecognizeBarrel(barrel.Message,*aiKey,*aiUrl,*aiModel,prompt,nil)
					if newBarrel == nil {
						failed++
						continue
					}
					if err != nil {
						ctx.Error(err)
						failed++
						continue
					}
					recognizedBarrel = newBarrel
					aiRecognize++
					goto adding				
				}
				recognizedBarrel.SellerName = splittedMessage[3]
				re := regexp.MustCompile("[^0-9]")
				quantity, err := strconv.ParseInt(re.ReplaceAllString(benefitParts[0],""),10,32)
				if err != nil {
					recognizedBarrel, err = aiFallbackRecognizeBarrel(barrel.Message,*aiKey,*aiUrl,*aiModel,prompt,err)
					if err != nil {
						ctx.Error(err)
						failed++
						continue
					}
					aiRecognize++
					goto adding
				}
				recognizedBarrel.Quantity = int(quantity);
				price,err := strconv.ParseInt(re.ReplaceAllString(benefitParts[1],""),10,32)
				if err != nil {
					recognizedBarrel, err = aiFallbackRecognizeBarrel(barrel.Message,*aiKey,*aiUrl,*aiModel,prompt,err)
					if err != nil {
						ctx.Error(err)
						failed++
						continue
					}
					aiRecognize++
					goto adding
				}
				recognizedBarrel.Price = int(price)
				recognizedBarrel.BenefitRatio = float64(recognizedBarrel.Quantity) / float64(recognizedBarrel.Price)
		  	} else if *aiKey != "" && *aiUrl != "" && *aiModel != "" {
				barrel, err := aiBarrelRecognition(*aiKey,*aiUrl,*aiModel,barrel.Message,prompt)
				if err != nil {
				  ctx.Error(err)
				  failed++
				  continue
				}
				recognizedBarrel = barrel
				aiRecognize++
		  	}
			adding:
		  	if recognizedBarrel == nil {
				failed++
				continue
		  	}
		  	sellerId, err := getOrCreateSeller(conn,recognizedBarrel.SellerName,ctx)
		  	barrelId, err := getOrCreateBarrel(conn,*barrel.X,*barrel.Y,*barrel.Z,ctx)
		  	itemId,err := getOrCreateItem(conn,recognizedBarrel.ItemName,ctx)
			if err != nil {
				ctx.Error(err)
			}
			barrels = append(barrels, AddBarrelStruct{ItemId: itemId,SellerId: sellerId,BarrelId: barrelId,BarrelText: barrel.Message,Quantity: recognizedBarrel.Quantity,Price: recognizedBarrel.Price,BenefitRatio: recognizedBarrel.BenefitRatio})
			success++
		}
		err = createBarrelItems(conn,barrels,ctx)
		if err != nil {
			ctx.Error(err)
		}
		telemetry.RecordsProcessed.Add(ctx, int64(success),
    	    metric.WithAttributes(attribute.String("result", "success")),
    	)

    	telemetry.RecordsProcessed.Add(ctx, int64(aiRecognize),
    	    metric.WithAttributes(attribute.String("result", "aiRecognize")),
    	)

    	telemetry.RecordsProcessed.Add(ctx, int64(failed),
    	    metric.WithAttributes(attribute.String("result", "failed")),
    	)
		ctx.JSON(200,gin.H{
		  "status": "success",
		})
	}
}