package main

type BarrelWithCords struct {
	Id int `json:"id"`
	X  int `json:"x"`
	Y  int `json:"y"`
	Z  int `json:"z"`
}

type SellerShort struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	UUID         string `json:"uuid"`
	BarrelsCount int    `json:"barrelsCount"`
}

type ItemWithCount struct {
	Id    int32  `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}