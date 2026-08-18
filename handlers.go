package main

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"zombiezen.com/go/sqlite"
)

func ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func getItems() (func(*gin.Context)) {
	return func (ctx *gin.Context)  {
		ctx.JSON(200,gin.H{
			"items": "items",
		})
	}
}

func getSnapshots(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
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

func getAllTypes(conn *sqlite.Conn) (func(*gin.Context)) {
	return func (ctx *gin.Context)  {
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

func getType(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
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

func getBarrels(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
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

func getBarrel(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
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

func getBarrelItems(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
	    ctx.JSON(200,gin.H{
	      "items": "items",
	    })
	}
}

func getSeller(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
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

func postBarrelItems(conn *sqlite.Conn) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
		var request ItemInBarrelPost
		if err := ctx.ShouldBindJSON(&request);err != nil{
		  ctx.Error(err)
		  ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			  return
		}
		err := postItemsInBarrel(conn,request)
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

func postBarrels(conn *sqlite.Conn,apiKey *string, xMin *int, xMax *int, yMin *int, yMax * int, zMin *int, zMax *int, prompt string, aiKey *string, aiModel *string, aiUrl *string) (func(*gin.Context)) {
	return func(ctx *gin.Context) {
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
		  var recognizedBarrel *RecognizedBarrelItem
	
		  splittedMessage := strings.Split(barrel.Message,"\n")
		  if len(splittedMessage) == 3{
			recognizedBarrel = &RecognizedBarrelItem{}
			recognizedBarrel.ItemName = splittedMessage[0]
			benefitParts := strings.Split(strings.ReplaceAll(splittedMessage[1]," ",""), "-")
			recognizedBarrel.SellerName = splittedMessage[2]
			re := regexp.MustCompile("[^0-9]")
			quantity, err := strconv.ParseFloat(re.ReplaceAllString(benefitParts[0],""),32)
			recognizedBarrel.Quantity = quantity;
			recognizedBarrel.Price, err = strconv.ParseFloat(re.ReplaceAllString(benefitParts[1],""),32)
			if err != nil {
				ctx.Error(err)
				return
			}
			recognizedBarrel.BenefitRatio = recognizedBarrel.Quantity / recognizedBarrel.Price
		  } else if *aiKey != "" && *aiUrl != "" && *aiModel != "" {
			barrel, err := aiBarrelRecognition(*aiKey,*aiUrl,*aiModel,barrel.Message,prompt)
			recognizedBarrel = barrel
			if err != nil {
			  ctx.Error(err)
			  return
			}
		  }
		  if recognizedBarrel == nil {
			continue
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
	}
}