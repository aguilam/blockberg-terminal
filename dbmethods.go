package main

import (
	"errors"
	"fmt"
	"strings"

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
	var allTypes []AllTypes
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			item := AllTypes{
				Count : stmt.ColumnInt32(0),
				Id : stmt.ColumnInt32(1),
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

func getBarrelByCords(conn *sqlite.Conn, x int, y int, z int) (*int32,error){
	query := `
		SELECT id FROM barrel
		WHERE x = ? AND y = ? AND z = ?`
	var id int32
	var found bool
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			id = stmt.ColumnInt32(0)
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
	INSERT INTO item_in_barrel (item_id, barrel_id, quantity,record_date)
	VALUES (?,?,?,?)
	`
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		Args: []any{item.ItemId,item.BarrelId,item.Quantity,item.RecordDate},
	})
	if err != nil{
		return err
	}
	return nil
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

func postBarrelItems(conn *sqlite.Conn,items ItemInBarrelPost) (error){
	barrelId, err := getBarrelByCords(conn,items.X,items.Y,items.Z)
	if err != nil{
		return err;
	}
	if barrelId == nil{
		return ErrNotFound;
	}

	for _,v := range items.Items {
		id, err := getOrCreateItem(conn,v.Name)
		if err != nil {
			continue
		}
		err = createItemInBarrel(conn,DBItemInBarrel{ItemId: *id,BarrelId: *barrelId,Quantity: v.Quantity})
		if err != nil {
			return err
		}
	}
	return nil;
}