package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"zombiezen.com/go/sqlite/sqlitex"
)

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
		response, err := getItemsSnapshotByQuery(conn,query,intPage,intPageSize)
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
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		
		items_types, err := getTypes(conn)
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
		items, err := getItemsByType(conn,intType)
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
		response, err := getBarrelsByQuery(conn,query,intPage,intPageSize)
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
		items, err := getBarrelsBySellerId(conn,int32(intId))
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
	}
}

func postBarrels(pool *sqlitex.Pool,apiKey *string, xMin *int, xMax *int, yMin *int, yMax * int, zMin *int, zMax *int, prompt string, aiKey *string, aiModel *string, aiUrl *string, isOnlyAI bool) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		allStart := time.Now()
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
		conn, err := pool.Take(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": "Internal error"})
			return
		}
		defer pool.Put(conn)
		var barrels []AddBarrelStruct
		for _, barrel := range request {
			parseStart := time.Now()
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
					continue
				}
				recognizedBarrel = recBarrel
				goto adding
		  	}
		  	if len(splittedMessage) == 3{
				recognizedBarrel.ItemName = splittedMessage[0]
				benefitParts := strings.Split(strings.ReplaceAll(splittedMessage[1]," ",""), "-")
				recognizedBarrel.SellerName = splittedMessage[2]
				re := regexp.MustCompile("[^0-9]")
				quantity, err := strconv.ParseFloat(re.ReplaceAllString(benefitParts[0],""),32)
				if err != nil {
					recognizedBarrel, err = aiFallbackRecognizeBarrel(barrel.Message,*aiKey,*aiUrl,*aiModel,prompt,err)
					if err != nil {
						ctx.Error(err)
						continue
					}
					goto adding
				}
				recognizedBarrel.Quantity = quantity;
				recognizedBarrel.Price, err = strconv.ParseFloat(re.ReplaceAllString(benefitParts[1],""),32)
				if err != nil {
					recognizedBarrel, err = aiFallbackRecognizeBarrel(barrel.Message,*aiKey,*aiUrl,*aiModel,prompt,err)
					if err != nil {
						ctx.Error(err)
						continue
					}
					goto adding
				}
				recognizedBarrel.BenefitRatio = recognizedBarrel.Quantity / recognizedBarrel.Price	
		  	} else if len(splittedMessage) == 4 {
				recognizedBarrel.ItemName = splittedMessage[0] + " " + splittedMessage[1]
				benefitParts := strings.Split(strings.ReplaceAll(splittedMessage[2]," ",""), "-")
				if len(benefitParts) != 2 {
					newBarrel, err := aiFallbackRecognizeBarrel(barrel.Message,*aiKey,*aiUrl,*aiModel,prompt,nil)
					if err != nil {
						ctx.Error(err)
						continue
					}
					recognizedBarrel = newBarrel
					goto adding				
				}
				recognizedBarrel.SellerName = splittedMessage[3]
				re := regexp.MustCompile("[^0-9]")
				quantity, err := strconv.ParseFloat(re.ReplaceAllString(benefitParts[0],""),32)
				if err != nil {
					recognizedBarrel, err = aiFallbackRecognizeBarrel(barrel.Message,*aiKey,*aiUrl,*aiModel,prompt,err)
					if err != nil {
						ctx.Error(err)
						continue
					}
					goto adding
				} 
				recognizedBarrel.Quantity = quantity;
				recognizedBarrel.Price, err = strconv.ParseFloat(re.ReplaceAllString(benefitParts[1],""),32)
				if err != nil {
					recognizedBarrel, err = aiFallbackRecognizeBarrel(barrel.Message,*aiKey,*aiUrl,*aiModel,prompt,err)
					if err != nil {
						ctx.Error(err)
						continue
					}
					goto adding
				}
				recognizedBarrel.BenefitRatio = recognizedBarrel.Quantity / recognizedBarrel.Price
		  	} else if *aiKey != "" && *aiUrl != "" && *aiModel != "" {
				barrel, err := aiBarrelRecognition(*aiKey,*aiUrl,*aiModel,barrel.Message,prompt)
				if err != nil {
				  ctx.Error(err)
				  continue
				}
				recognizedBarrel = barrel
		  	}
			adding:
		  	if recognizedBarrel == nil {
				continue
		  	}
			parseElapsed := time.Since(parseStart)
			fmt.Println("Время парсинга:", parseElapsed)
			dbTime := time.Now()
		  	sellerId, err := getOrCreateSeller(conn,recognizedBarrel.SellerName)
		  	barrelId, err := getOrCreateBarrel(conn,*barrel.X,*barrel.Y,*barrel.Z)
		  	itemId,err := getOrCreateItem(conn,recognizedBarrel.ItemName)
			if err != nil {
				ctx.Error(err)
			}
			dbElapsed := time.Since(dbTime)
			fmt.Println("Время БД:", dbElapsed)
			barrels = append(barrels, AddBarrelStruct{ItemId: itemId,SellerId: sellerId,BarrelId: barrelId,BarrelText: barrel.Message,Quantity: recognizedBarrel.Quantity,Price: recognizedBarrel.Price,BenefitRatio: recognizedBarrel.BenefitRatio})
		}
		dbTimee := time.Now()
		err = createBarrelItems(conn,barrels)
		dbElapsed := time.Since(dbTimee)
		fmt.Println("Время БД:", dbElapsed)
		if err != nil {
			ctx.Error(err)
		}
		elapsed := time.Since(allStart)
		fmt.Println("Время выполнения:", elapsed)
		ctx.JSON(200,gin.H{
		  "status": "success",
		})
	}
}