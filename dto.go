package main

import (
	"time"
)

type NewBarrelPost struct {
	X       *int    `json:"x" binding:"required"`
	Y       *int    `json:"y" binding:"required"`
	Z       *int    `json:"z" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type ItemInBarrelPost struct {
	X     *int            `json:"x" binding:"required"`
	Y     *int            `json:"y" binding:"required"`
	Z     *int            `json:"z" binding:"required"`
	Items []ItemInBarrel `json:"items" binding:"required"`
}

type ItemInBarrel struct {
	Name     string `json:"name" binding:"required"`
	Quantity int32  `json:"quantity" binding:"required"`
}

type DBItemInBarrel struct {
	Id         *int32
	ItemId     int32
	SnapshotId int
	Quantity   int32
}
type AllTypes struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Count int32  `json:"count"`
}

type Item struct {
	Id             int32
	Name           string
	Description    *string
	NormalizedName string
}

type BarrelItem struct {
	Id           int32
	Item         Item
	Quantity     int32
	Price        int32
	BenefitRatio float32
	RecordDate   string
}

type BarrelItemFull struct {
	Id           int  `json:"id"`
	Name         string  `json:"name"`
	Seller       string  `json:"seller"`
	Price        int     `json:"price"`
	Quantity     int     `json:"quantity"`
	X            int     `json:"x"`
	Y            int     `json:"y"`
	Z            int     `json:"z"`
	BenefitRatio float32 `json:"benefitRatio"`
	SnapshotsCount int `json:"snapshotsCount"`
	RecordDate   time.Time  `json:"recordDate"`
}

type ItemsSnapshotSearchResponse struct {
	Id           int  `json:"id"`
	X            int     `json:"x"`
	Y            int     `json:"y"`
	Z            int     `json:"z"`
	Items []ItemInBarrel `json:"items"`
	RecordDate   time.Time  `json:"recordDate"`
}

type ItemsSnapshot struct {
	Items []ItemInBarrel `json:"items"`
	RecordDate   time.Time  `json:"recordDate"`
}


type BarrelItemResponse struct {
	Barrels []BarrelItemFull `json:"barrels"`
	Total   int              `json:"total"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
}

type PaginatedSnapshotSearchResponse struct {
	Items []ItemsSnapshotSearchResponse `json:"items"`
	Total   int              `json:"total"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
}


type BarrelInfo struct {
	Id           int  `json:"id"`
	BarrelText	 string `json:"barrelText"`
	Name         string  `json:"name"`
	Seller       string  `json:"seller"`
	Price        int     `json:"price"`
	Quantity     int     `json:"quantity"`
	BenefitRatio float32 `json:"benefitRatio"`
	X            int     `json:"x"`
	Y            int     `json:"y"`
	Z            int     `json:"z"`
	RecordDate   time.Time  `json:"recordDate"`
	ItemsSnapshot ItemsSnapshot `json:"barrelItems"`
}
type AiRecognizedBarrelItem struct {
	SellerName string `json:"seller_name" jsonschema_description:"Minecraft username of the seller. If no seller name is present, return 'None'"`
    ItemName string `json:"item_name" jsonschema_description:"Normalized Minecraft item name without emojis, extra symbols, or quantity information"`
    Quantity int32 `json:"quantity" jsonschema_description:"Total number of items. Convert stacks to items (1 stack = 64 items)"`
    Price int32 `json:"price" jsonschema_description:"Price of items. If input contains blocks, convert them to single items (1 block = 9 diamonds)"`
	TypeName string `json:"type_name" jsonschema_description:"Item category" jsonschema:"enum=building blocks, enum=colored blocks, enum=natural blocks, enum=functional blocks, enum=redstone blocks, enum=tools & utilities, enum=combat, enum=food & drinks, enum=ingredients, enum=spawn eggs,enum=other"`
}
type RecognizedBarrelItem struct {
	SellerName string 
    ItemName string
    BenefitRatio float64
    Quantity float64
    Price float64 
	TypeName string
}