package main

import (
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func getBarrelsWithCords(conn *sqlite.Conn) ([]BarrelWithCords,error) {
	query := "SELECT COUNT(id), id, x, y, z FROM barrel"
	barrels := []BarrelWithCords{}
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			id := stmt.ColumnInt(0)
			x := stmt.ColumnInt(1)
			y := stmt.ColumnInt(2)
			z := stmt.ColumnInt(3)
			barrel := BarrelWithCords{Id: id,X: x,Y: y,Z: z}
			barrels = append(barrels,barrel)
			return nil
		},
	})
	if err != nil {
		return nil,err
	}
	return barrels,nil
}

func getTopSellers(conn *sqlite.Conn) ([]SellerShort,error) {
	query := `SELECT seller.id, seller.username, seller.minecraft_uuid, COUNT(barrel.id) 
				FROM seller 
				JOIN barrel_item AS barrel ON barrel.seller_id = seller.id 
				GROUP BY seller.id
				ORDER BY COUNT(barrel.id) DESC, seller.id DESC
				LIMIT 10`
	sellers := []SellerShort{}
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			id := stmt.ColumnInt(0)
			username := stmt.ColumnText(1)
			uuid := stmt.ColumnText(2)
			barrelsCount := stmt.ColumnInt(3)
			seller := SellerShort{Id: id,Name: username,UUID: uuid,BarrelsCount: barrelsCount}
			sellers = append(sellers, seller)
			return nil
		},
	})
	if err != nil {
		return nil,err
	}
	return sellers,nil
}

func getBarrelsMaxPrices(conn *sqlite.Conn) ([]BarrelItem,error) {
	query := `SELECT 
			b_item.id, b_item.item_id, b_item.barrel_id, b_item.seller_id, b_item.quantity, 
    		b_item.price, b_item.benefit_ratio, b_item.record_date,
			item.id, item.name, item.description, item.normalized_name
		FROM barrel_item as b_item
		LEFT JOIN seller ON b_item.seller_id = seller.id
		JOIN item ON b_item.item_id = item.id
		ORDER BY b_item.price DESC
		LIMIT 5`
	items := []BarrelItem{}
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			var desc *string

    		if stmt.ColumnType(10) != sqlite.TypeNull {
    		    text := stmt.ColumnText(10)
    		    desc = &text 
    		}
			item := BarrelItem{
				Id: stmt.ColumnInt32(0),
				Item: Item{
					Id: stmt.ColumnInt32(8),
					Name: stmt.ColumnText(9),
					Description: desc,
					NormalizedName: stmt.ColumnText(11),
				},
				Quantity: stmt.ColumnInt32(4),
				Price: stmt.ColumnInt32(5),
				BenefitRatio: float32(stmt.ColumnFloat(6)),
				RecordDate: stmt.ColumnText(7),
			}
			items = append(items, item)
			return nil
		},
	})
	if err != nil {
		return nil,err
	}
	return items,nil
}

func getBarrelsMaxBenefit(conn *sqlite.Conn) ([]BarrelItem,error) {
	query := `SELECT 
		b_item.id, b_item.item_id, b_item.barrel_id, b_item.seller_id, b_item.quantity, 
		b_item.price, b_item.benefit_ratio, b_item.record_date,
		item.id, item.name, item.description, item.normalized_name
		FROM barrel_item as b_item
		LEFT JOIN seller ON b_item.seller_id = seller.id
		JOIN item ON b_item.item_id = item.id
		ORDER BY b_item.benefit_ratio DESC
		LIMIT 5`
	items := []BarrelItem{}
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
	ResultFunc: func(stmt *sqlite.Stmt) error {
		var desc *string

		if stmt.ColumnType(10) != sqlite.TypeNull {
			text := stmt.ColumnText(10)
			desc = &text 
		}
		item := BarrelItem{
			Id: stmt.ColumnInt32(0),
			Item: Item{
				Id: stmt.ColumnInt32(8),
				Name: stmt.ColumnText(9),
				Description: desc,
				NormalizedName: stmt.ColumnText(11),
			},
			Quantity: stmt.ColumnInt32(4),
			Price: stmt.ColumnInt32(5),
			BenefitRatio: float32(stmt.ColumnFloat(6)),
			RecordDate: stmt.ColumnText(7),
		}
		items = append(items, item)
		return nil
	},
	})
	if err != nil {
	return nil,err
	}
	return items,nil
}

func topSellingBarrelItems(conn *sqlite.Conn) ([]ItemWithCount,error) {
	query := `SELECT item.id, item.name,COUNT(b_item.id)
		FROM item
		JOIN barrel_item as b_item ON b_item.item_id = item.id
		GROUP BY item.id
		ORDER BY COUNT(b_item.id) DESC
		LIMIT 5`
		items := []ItemWithCount{}
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			item := ItemWithCount{
				Id: stmt.ColumnInt32(0),
				Name: stmt.ColumnText(1),
				Count: stmt.ColumnInt(2),
			}
			items = append(items, item)
			return nil
		},
	})
	if err != nil {
	return nil,err
	}
	return items,nil
}

func topStoredItems(conn *sqlite.Conn) ([]ItemWithCount,error) {
	query := `SELECT item.id, item.name, SUM(item_in_barrel.quantity)
		FROM item
		JOIN item_in_barrel ON item_in_barrel.item_id = item.id
		GROUP BY item.id
		ORDER BY SUM(item_in_barrel.quantity) DESC
		LIMIT 5`
		items := []ItemWithCount{}
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			item := ItemWithCount{
				Id: stmt.ColumnInt32(0),
				Name: stmt.ColumnText(1),
				Count: stmt.ColumnInt(2),
			}
			items = append(items, item)
			return nil
		},
	})
	if err != nil {
	return nil,err
	}
	return items,nil
}