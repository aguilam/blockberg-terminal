package main

import (
	"time"
)

type NewBarrelPost struct {
	X       int    `json:"x" binding:"required"`
	Y       int    `json:"y" binding:"required"`
	Z       int    `json:"z" binding:"required"`
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
	ItemsSnapshot ItemsSnapshot `json:"barrel_items"`
}