package main

type NewBarrelPost struct {
	X       string `json:"x" binding:"required"`
	Y       string `json:"y" binding:"required"`
	Z       string `json:"z" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type ItemInBarrelPost struct {
	Name     string         `json:"name" binding:"required"`
	Quantity string         `json:"quantity" binding:"required"`
	Items    []ItemInBarrel `json:"items" binding:"required"`
}

type ItemInBarrel struct {
	Name     string `json:"name" binding:"required"`
	Quantity string `json:"quantity" binding:"required"`
}

type AllTypes struct {
	Id    int32
	Name  string
	Count int32
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