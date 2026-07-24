package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var ErrNotFound = errors.New("Not found")

func getBarrelHistory(conn *sqlite.Conn,x int,y int, z int) {
	return
}

func getTypes(conn *sqlite.Conn)([]AllTypes,error){
	query := `
		SELECT COUNT(item.id),type.id,type.name 
		FROM item_type as type 
		LEFT JOIN item ON type_id = type.id
		GROUP BY type.id;`
	allTypes := []AllTypes{}
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			item := AllTypes{
				Count : stmt.ColumnInt32(0),
				ID : stmt.ColumnInt32(1),
				Name: stmt.ColumnText(2),
			}
			allTypes = append(allTypes, item)
			return nil
		},
	})
	if err != nil {
		return nil,fmt.Errorf("Get all types error: %w",err)
	}
	return allTypes,nil
}

func getItemsByType(conn *sqlite.Conn, typeId int) ([]Item,error){
	query := `
		SELECT id, name, description, normalized_name 
		FROM item
		WHERE item.type_id = ?;`
	var items []Item
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			var desc *string

    		if stmt.ColumnType(2) != sqlite.TypeNull {
    		    text := stmt.ColumnText(2)
    		    desc = &text 
    		}
			item := Item{
				Id: stmt.ColumnInt32(0),
				Name: stmt.ColumnText(1),
				Description: desc,
				NormalizedName: stmt.ColumnText(3),
			}
			items = append(items,item)
			return nil
		},
		Args: []any{typeId},
	})
	if err != nil {
		return nil,fmt.Errorf("Get items by type error: %w",err)
	}
	return items,nil
}

func getBarrelsBySellerId(conn *sqlite.Conn,sellerId int32) ([]BarrelItem, error){
	query := `
		SELECT 
			b_item.id, b_item.item_id, b_item.barrel_id, b_item.seller_id, b_item.quantity, 
    		b_item.price, b_item.benefit_ratio, b_item.record_date,
			item.id, item.name, item.description, item.normalized_name, item.type_id
		FROM barrel_item as b_item
		JOIN seller ON b_item.seller_id = seller.id
		JOIN item ON b_item.item_id = item.id
		WHERE seller.id = ?;`
	var sellerBarrels []BarrelItem
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
			sellerBarrels = append(sellerBarrels, item)
			return nil
		},
		Args: []any{sellerId},
	})
	if err != nil {
		return nil,fmt.Errorf("Get barrel by seller error: %w",err)
	}
	return sellerBarrels,nil
}

func getBarrelByCords(conn *sqlite.Conn, x int, y int, z int) (*int,error){
	query := `
		SELECT id FROM barrel
		WHERE x = ? AND y = ? AND z = ?`
	var id int
	var found bool
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			id = stmt.ColumnInt(0)
			found = true
			return nil
		},
		Args: []any{x,y,z},
	})
	if err != nil {
        return nil, err
    }
	if !found {
		return nil,nil
	}
	return &id,nil
}
func createItemInBarrel(conn *sqlite.Conn, item DBItemInBarrel) error{
	query := `
	INSERT INTO item_in_barrel (item_id, snapshot_id, quantity)
	VALUES (?,?,?)
	`
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		Args: []any{item.ItemId,item.SnapshotId,item.Quantity},
	})
	if err != nil{
		return err
	}
	return nil
}

func createItemsSnapshot(conn *sqlite.Conn, barrelId int) (*int, error) {
	var snapshotId int
	query := `
		INSERT INTO storage_snapshot (barrel_id)
		VALUES (?)
		RETURNING id
	`
	err := sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			snapshotId = stmt.ColumnInt(0)
			return nil
		},
		Args: []any{barrelId},
	})
	if err != nil {
		return nil, err
	}
	return &snapshotId, err
}

func getOrCreateItem(conn *sqlite.Conn, itemName string) (*int32, error) {
    normalizedName := strings.ToLower(strings.TrimSpace(itemName))
    var id int32
    var found bool

    selectQuery := `
		SELECT id, name, normalized_name 
		FROM item 
		WHERE normalized_name = ?;`
    err := sqlitex.Execute(conn, selectQuery, &sqlitex.ExecOptions{
        Args: []any{normalizedName},
        ResultFunc: func(stmt *sqlite.Stmt) error {
        	id = stmt.ColumnInt32(0)
            found = true
            return nil
        },
    })
	if found {
        return &id, nil
    }

    if err != nil {
        return nil, err
    }

    insertQuery := `
		INSERT INTO item (name, normalized_name) 
		VALUES (?, ?) 
		RETURNING id, name, normalized_name;`
    err = sqlitex.Execute(conn, insertQuery, &sqlitex.ExecOptions{
        Args: []any{itemName, normalizedName},
        ResultFunc: func(stmt *sqlite.Stmt) error {
            id = stmt.ColumnInt32(0)
            return nil
        },
    })
    if err != nil {
        return nil, err
    }

    return &id, nil
}

func postItemsInBarrel(conn *sqlite.Conn,items ItemInBarrelPost) (error) {
	barrelId, err := getBarrelByCords(conn,*items.X,*items.Y,*items.Z)
	if err != nil{
		return err;
	}
	if barrelId == nil{
		return ErrNotFound;
	}

	snapshotId, err := createItemsSnapshot(conn, *barrelId)
	if err != nil {
		return nil
	}

	for _,v := range items.Items {
		id, err := getOrCreateItem(conn,v.Name)
		if err != nil {
			continue
		}
		err = createItemInBarrel(conn, DBItemInBarrel{ItemId: *id, SnapshotId: *snapshotId, Quantity: v.Quantity})
		if err != nil {
			return err
		}
	}
	return nil;
}

func getSellerByName(conn *sqlite.Conn,sellerName string) (*int32, error) {
	normalizedName := strings.ToLower(strings.TrimSpace(sellerName))
	found := false
	var id int32
	stmt := `
		SELECT id FROM seller
		WHERE normalized_username = ?`
	err := sqlitex.Execute(conn,stmt,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			id = stmt.ColumnInt32(0)
			found = true
			return nil
		},
		Args: []any {normalizedName},
	})
	if err != nil {
		return nil, err
	}
	if !found {
		return nil,nil
	}
	return &id, nil
}

func createSeller(conn *sqlite.Conn, name string, minecraft_id *string) (*int32, error) {
	normalized_username := strings.ToLower(strings.TrimSpace(name))
	var id int32
	var uuid any = nil
	if minecraft_id != nil {
	    uuid = *minecraft_id
	}

	stmt := `
		INSERT INTO seller (username,normalized_username,minecraft_uuid)
		VALUES (?, ?, ?)
		RETURNING id
	`
	err := sqlitex.Execute(conn,stmt,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			id = stmt.ColumnInt32(0)
			return nil
		},
		Args: []any{name,normalized_username,uuid},
	})
	if err != nil{
		return nil,err
	}
	return &id,nil
}

func createBarrel(conn *sqlite.Conn, x int, y int, z int ) (*int,error) {
	var id int
	stmt := `
		INSERT INTO barrel (x,y,z)
		VALUES (?, ?, ?)
		RETURNING id
	`
	err := sqlitex.Execute(conn,stmt,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			id = stmt.ColumnInt(0)
			return nil
		},
		Args: []any{x,y,z},
	})
	if err != nil {
		return nil,err
	}
	return &id,nil
}

func createBarrelItem(conn *sqlite.Conn,itemId int32, barrelId int, sellerId int32,recognizedBarrel RecognizedBarrelItem,barrelText string) (*int32, error){
	var id int32
	stmt := `
		INSERT INTO barrel_item (item_id, barrel_id, seller_id, quantity, price, benefit_ratio, barrel_text)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	err := sqlitex.Execute(conn,stmt,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			id = stmt.ColumnInt32(0)
			return nil
		},
		Args: []any{itemId,barrelId,sellerId,recognizedBarrel.Quantity,recognizedBarrel.Price,recognizedBarrel.BenefitRatio,barrelText},
	})
	if err != nil {
		return nil,err
	}
	return &id,nil
}

func getBarrelsByQuery(conn *sqlite.Conn, query string, page int, pageSize int) (*BarrelItemResponse, error) {
	normalized_query := strings.ToLower(strings.TrimSpace(query))
	var barrels []BarrelItemFull;
	offset := (page - 1) * pageSize
	var total int;

	sqlite_query := `
		SELECT id, name, username, price, quantity, x, y, z, benefit_ratio, record_date, snapshots_count FROM (
		    SELECT 
		        barrel.id, 
		        item.name, 
		        seller.username, 
		        barrel_item.price, 
		        barrel_item.quantity, 
		        barrel.x, 
		        barrel.y, 
		        barrel.z, 
		        barrel_item.benefit_ratio, 
		        barrel_item.record_date, 
		        item.normalized_name, 
		        ROW_NUMBER() OVER (
		            PARTITION BY barrel_item.barrel_id
		            ORDER BY barrel_item.record_date DESC, barrel_item.id DESC
		        ) AS barrel_pos,
		        COUNT(snapshot.id) OVER (PARTITION BY barrel_item.barrel_id) AS snapshots_count
		    FROM item
		    JOIN barrel_item ON barrel_item.item_id = item.id
		    JOIN barrel ON barrel.id = barrel_item.barrel_id
		    LEFT JOIN storage_snapshot AS snapshot ON snapshot.barrel_id = barrel.id
		    JOIN seller ON seller.id = barrel_item.seller_id
		)
		WHERE barrel_pos = 1 AND normalized_name LIKE '%' || ? || '%'
		ORDER BY INSTR(normalized_name, ?) ASC, LENGTH(normalized_name) ASC, benefit_ratio ASC
		LIMIT ?
		OFFSET ?
	`
	err := sqlitex.Execute(conn,sqlite_query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			record_date, _ := time.Parse("2006-01-02 15:04:05",stmt.ColumnText(9))
			barrel := BarrelItemFull{
				Id: stmt.ColumnInt(0),
				Name: stmt.ColumnText(1),
				Seller: stmt.ColumnText(2),
				Price: stmt.ColumnInt(3),
				Quantity: stmt.ColumnInt(4),
				X: stmt.ColumnInt(5),
				Y: stmt.ColumnInt(6),
				Z: stmt.ColumnInt(7),
				BenefitRatio: float32(stmt.ColumnFloat(8)),
				SnapshotsCount: stmt.ColumnInt(10),
				RecordDate: record_date,
			}
			barrels = append(barrels, barrel)
			return nil
		},
		Args: []any{normalized_query,normalized_query,pageSize,offset},
	})

	if err != nil {
		return nil,err
	}

	sqlite_query = `
		SELECT COUNT(*) FROM (
			SELECT barrel_item.id,
				ROW_NUMBER() OVER (
    				PARTITION BY barrel_item.barrel_id
    				ORDER BY barrel_item.record_date DESC, barrel_item.id DESC
				) AS barrel_pos, item.normalized_name
			FROM item
			JOIN barrel_item ON barrel_item.item_id = item.id
		)
		WHERE barrel_pos = 1 AND normalized_name LIKE '%' || ? || '%'
	`
	err = sqlitex.Execute(conn,sqlite_query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			total = int(stmt.ColumnInt32(0))
			return nil
		},
		Args: []any{normalized_query},
	})

	if err != nil {
		return nil,err
	}
	return &BarrelItemResponse{Barrels: barrels, Total: total, Page: page, Limit: pageSize}, nil
}

func getBarrelInfo(conn *sqlite.Conn, barrelId int) (*BarrelInfo,error) {
	var barrelInfo BarrelInfo
	var lastSnapshotId int
	var found bool
	
	query := `
		SELECT barrel.id, bi.barrel_text, item.name, seller.username, bi.price, bi.quantity, bi.benefit_ratio, barrel.x, barrel.y, barrel.z, bi.record_date, snapshot.id, snapshot.record_date
		FROM barrel
		JOIN barrel_item bi ON bi.id = (SELECT id FROM barrel_item WHERE barrel_id = ? ORDER BY record_date DESC, id DESC LIMIT 1)
		JOIN item ON item.id = bi.item_id
		JOIN seller ON seller.id = bi.seller_id
		LEFT JOIN storage_snapshot as snapshot ON snapshot.id = (SELECT id FROM storage_snapshot WHERE barrel_id = ? ORDER BY record_date DESC, id DESC LIMIT 1) 
		WHERE barrel.id = ?
	`
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			record_date, _ := time.Parse("2006-01-02 15:04:05",stmt.ColumnText(10))
			barrelInfo.Id = stmt.ColumnInt(0)
			barrelInfo.BarrelText = stmt.ColumnText(1)
			barrelInfo.Name = stmt.ColumnText(2)
			barrelInfo.Seller = stmt.ColumnText(3)
			barrelInfo.Price = stmt.ColumnInt(4)
			barrelInfo.Quantity = stmt.ColumnInt(5)
			barrelInfo.BenefitRatio = float32(stmt.ColumnFloat(6))
			barrelInfo.X = stmt.ColumnInt(7)
			barrelInfo.Y = stmt.ColumnInt(8)
			barrelInfo.Z = stmt.ColumnInt(9)
			barrelInfo.RecordDate = record_date
			lastSnapshotId = stmt.ColumnInt(11)
			snapshot_date, _ := time.Parse("2006-01-02 15:04:05",stmt.ColumnText(12))
			barrelInfo.ItemsSnapshot.RecordDate = snapshot_date
			found = true
			return nil
		},
		Args: []any{barrelId,barrelId,barrelId},
	})
	if err != nil {
		return nil, err
	}

	query = `
		SELECT item.name, quantity
		FROM item_in_barrel
		JOIN item ON item.id = item_in_barrel.item_id
		WHERE snapshot_id = ?
	`

	err = sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			item := ItemInBarrel{
				Name: stmt.ColumnText(0),
				Quantity: stmt.ColumnInt32(1),
			}
			barrelInfo.ItemsSnapshot.Items = append(barrelInfo.ItemsSnapshot.Items, item)
			return nil
		},
		Args: []any{lastSnapshotId},
	})
	if err != nil {
		return nil, err
	}

	if !found {
		return nil,nil
	}

	return &barrelInfo, nil
}