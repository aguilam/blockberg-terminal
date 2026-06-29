package main

type NewBarrelPost struct {
	X       int    `json:"x" binding:"required"`
	Y       int    `json:"y" binding:"required"`
	Z       int    `json:"z" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type ItemInBarrelPost struct {
	X     int            `json:"x" binding:"required"`
	Y     int            `json:"y" binding:"required"`
	Z     int            `json:"z" binding:"required"`
	Items []ItemInBarrel `json:"items" binding:"required"`
}

type ItemInBarrel struct {
	Name     string `json:"name" binding:"required"`
	Quantity int32  `json:"quantity" binding:"required"`
}

type DBItemInBarrel struct {
	Id         *int32
	ItemId     int32
	BarrelId   int32
	Quantity   int32
	RecordDate *string
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