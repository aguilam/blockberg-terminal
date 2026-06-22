package main

type NewBarrelPost struct {
	X       string `json:"x" binding:"required"`
	Y       string `json:"y" binding:"required"`
	Z       string `json:"z" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type ItemInBarrel struct {
	Name     string `json:"name" binding:"required"`
	Quantity string `json:"quantity" binding:"required"`
}

type ItemInBarrelPost struct {
	Name     string         `json:"name" binding:"required"`
	Quantity string         `json:"quantity" binding:"required"`
	Items    []ItemInBarrel `json:"items" binding:"required"`
}