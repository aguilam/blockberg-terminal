package main

import (
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func getBarrelsWithCords(conn *sqlite.Conn) ([]BarrelWithCords,error) {
	query := "SELECT id, x, y, z FROM barrel"
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

func getTopSellers(conn *sqlite.Conn) ([]DatedData[SellerShort],error) {
	query := `WITH daily AS (
		SELECT 
			seller.id,
			seller.username,
			seller.minecraft_uuid, 
			COUNT(*) AS barrels_count,
			date(b_i.record_date) AS day
		FROM barrel_item as b_i
		JOIN seller ON seller.id = b_i.seller_id
		GROUP BY date(b_i.record_date), seller_id
		),
		ranked AS (
		    SELECT *,
		        ROW_NUMBER() OVER (
		            PARTITION BY day
		            ORDER BY barrels_count DESC, id
		        ) AS position
		    FROM daily
		)
		SELECT *
		FROM ranked
		WHERE position <= 5
		ORDER BY day, position
		`
	var result []DatedData[SellerShort]
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			id := stmt.ColumnInt(0)
			username := stmt.ColumnText(1)
			uuid := stmt.ColumnText(2)
			barrelsCount := stmt.ColumnInt(3)
			day := stmt.ColumnText(4)
			seller := SellerShort{Id: id,Name: username,UUID: uuid,BarrelsCount: barrelsCount}
			if len(result) == 0 || result[len(result)-1].Day != day {
				result = append(result, DatedData[SellerShort]{
					Day:   day,
					Data: []SellerShort{},
				})
			}
			result[len(result)-1].Data = append(
				result[len(result)-1].Data,
				seller,
			)
			return nil
		},
	})
	if err != nil {
		return nil,err
	}
	return result,nil
}

func getItemMaxPrices(conn *sqlite.Conn) ([]DatedData[BarrelItem],error) {
	query := `WITH daily AS (
	    SELECT 
	        b_i.id, b_i.quantity, 
			b_i.price, b_i.benefit_ratio, b_i.record_date,
			item.id, item.name,
	        date(b_i.record_date) as day
	    FROM barrel_item as b_i
	    JOIN item ON b_i.item_id = item.id
		GROUP BY item.id
	),
	ranked AS (
	    SELECT *,
	        ROW_NUMBER() OVER (
	            PARTITION BY day 
	            ORDER BY price DESC, id
	        ) as position
	    FROM daily
	    WHERE benefit_ratio > 0
	)
	SELECT *
	FROM ranked
	WHERE position <= 5 
	ORDER BY day, position`
	var result []DatedData[BarrelItem]
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			item := BarrelItem{
				Id: stmt.ColumnInt32(0),
				Item: Item{
					Id: stmt.ColumnInt32(5),
					Name: stmt.ColumnText(6),
				},
				Quantity: stmt.ColumnInt32(1),
				Price: stmt.ColumnInt32(2),
				BenefitRatio: float32(stmt.ColumnFloat(3)),
				RecordDate: stmt.ColumnText(4),
			}
			day := stmt.ColumnText(7)
			if len(result) == 0 || result[len(result)-1].Day != day {
				result = append(result, DatedData[BarrelItem]{
					Day:   day,
					Data: []BarrelItem{},
				})
			}
			result[len(result)-1].Data = append(
				result[len(result)-1].Data,
				item,
			)
			return nil
		},
	})
	if err != nil {
		return nil,err
	}
	return result,nil
}

func getBarrelsMaxPrices(conn *sqlite.Conn) ([]DatedData[BarrelItem],error) {
	query := `WITH daily AS (
	    SELECT 
	        b_i.id, b_i.quantity, 
			b_i.price, b_i.benefit_ratio, b_i.record_date,
			item.id, item.name,
	        date(b_i.record_date) as day
	    FROM barrel_item as b_i
	    JOIN item ON b_i.item_id = item.id
	),
	ranked AS (
	    SELECT *,
	        ROW_NUMBER() OVER (
	            PARTITION BY day 
	            ORDER BY price DESC, id
	        ) as position
	    FROM daily
	    WHERE benefit_ratio > 0
	)
	SELECT *
	FROM ranked
	WHERE position <= 5 
	ORDER BY day, position`
	var result []DatedData[BarrelItem]
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			item := BarrelItem{
				Id: stmt.ColumnInt32(0),
				Item: Item{
					Id: stmt.ColumnInt32(5),
					Name: stmt.ColumnText(6),
				},
				Quantity: stmt.ColumnInt32(1),
				Price: stmt.ColumnInt32(2),
				BenefitRatio: float32(stmt.ColumnFloat(3)),
				RecordDate: stmt.ColumnText(4),
			}
			day := stmt.ColumnText(7)
			if len(result) == 0 || result[len(result)-1].Day != day {
				result = append(result, DatedData[BarrelItem]{
					Day:   day,
					Data: []BarrelItem{},
				})
			}
			result[len(result)-1].Data = append(
				result[len(result)-1].Data,
				item,
			)
			return nil
		},
	})
	if err != nil {
		return nil,err
	}
	return result,nil
}

func getBarrelsMaxBenefit(conn *sqlite.Conn) ([]DatedData[BarrelItem],error) {
	query := `WITH daily AS (
	    SELECT 
	        b_i.id, b_i.quantity, 
			b_i.price, b_i.benefit_ratio, b_i.record_date,
			item.id, item.name,
	        date(b_i.record_date) as day
	    FROM barrel_item as b_i
	    JOIN item ON b_i.item_id = item.id
	),
	ranked AS (
	    SELECT *,
	        ROW_NUMBER() OVER (
	            PARTITION BY day 
	            ORDER BY benefit_ratio DESC, id
	        ) as position
	    FROM daily
	    WHERE benefit_ratio > 0
	)
	SELECT *
	FROM ranked
	WHERE position <= 5 
	ORDER BY day, position`
	var result []DatedData[BarrelItem]
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			item := BarrelItem{
				Id: stmt.ColumnInt32(0),
				Item: Item{
					Id: stmt.ColumnInt32(5),
					Name: stmt.ColumnText(6),
				},
				Quantity: stmt.ColumnInt32(1),
				Price: stmt.ColumnInt32(2),
				BenefitRatio: float32(stmt.ColumnFloat(3)),
				RecordDate: stmt.ColumnText(4),
			}
			day := stmt.ColumnText(7)
			if len(result) == 0 || result[len(result)-1].Day != day {
				result = append(result, DatedData[BarrelItem]{
					Day:   day,
					Data: []BarrelItem{},
				})
			}
			result[len(result)-1].Data = append(
				result[len(result)-1].Data,
				item,
			)
		return nil
	},
	})
	if err != nil {
	return nil,err
	}
	return result,nil
}

func topSellingBarrelItems(conn *sqlite.Conn) ([]DatedData[ItemWithCount],error) {
	query := `WITH daily AS (
		    SELECT 
		        item.id,
				item.name,
		        COUNT(barrel_item.id) AS barrels_count,
		        date(barrel_item.record_date) as day
		    FROM barrel_item 
		    JOIN item ON barrel_item.item_id = item.id
		    GROUP BY date(barrel_item.record_date), item.id
		),
		ranked AS (
		    SELECT *,
		        ROW_NUMBER() OVER (
		            PARTITION BY day 
		            ORDER BY barrels_count DESC,id
		        ) as position
		    FROM daily
		)
		SELECT *
		FROM ranked
		WHERE position <= 5
		ORDER BY day, position
		`
	var result []DatedData[ItemWithCount]
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			item := ItemWithCount{
				Id: stmt.ColumnInt32(0),
				Name: stmt.ColumnText(1),
				Count: stmt.ColumnInt(2),
			}
			day := stmt.ColumnText(3)
			if len(result) == 0 || result[len(result)-1].Day != day {
				result = append(result, DatedData[ItemWithCount]{
					Day:day,
					Data: []ItemWithCount{},
				})
			}
			result[len(result)-1].Data = append(result[len(result)-1].Data, item)
			return nil
		},
	})
	if err != nil {
	return nil,err
	}
	return result,nil
}

func topStoredItems(conn *sqlite.Conn) ([]DatedData[ItemWithCount],error) {
	query := `WITH daily AS (
		    SELECT 
		        item.id,
				item.name,
		        SUM(item_in_barrel.quantity) AS item_count,
		        date(storage_snapshot.record_date) AS day
		    FROM item_in_barrel 
		    JOIN item ON item_in_barrel.item_id = item.id
		    LEFT JOIN storage_snapshot ON storage_snapshot.id = item_in_barrel.snapshot_id
		    GROUP BY date(storage_snapshot.record_date),item.id
		),
		ranked AS (
		    SELECT *,
		        ROW_NUMBER() OVER (
		            PARTITION BY day 
		            ORDER BY item_count DESC,id
		        ) AS position
		    FROM daily
		)
		SELECT *
		FROM ranked
		WHERE position <= 5
		ORDER BY day, position
		`
	var result []DatedData[ItemWithCount]
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			item := ItemWithCount{
				Id: stmt.ColumnInt32(0),
				Name: stmt.ColumnText(1),
				Count: stmt.ColumnInt(2),
			}
			day := stmt.ColumnText(3)
			if len(result) == 0 || result[len(result)-1].Day != day {
				result = append(result, DatedData[ItemWithCount]{
					Day:day,
					Data: []ItemWithCount{},
				})
			}
			result[len(result)-1].Data = append(result[len(result)-1].Data, item)
			return nil
		},
	})
	if err != nil {
	return nil,err
	}
	return result,nil
}