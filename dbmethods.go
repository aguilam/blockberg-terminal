package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var ErrNotFound = errors.New("Not found")
var dTracer = otel.Tracer("blockberg-terminal/db")

func getBarrelHistory(conn *sqlite.Conn,x int,y int, z int) {
	return
}

func getTypes(conn *sqlite.Conn,c *gin.Context)([]AllTypes,error){
	_,span := dTracer.Start(c.Request.Context(),"types.all.db")
	defer span.End()
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

func getItemsByType(conn *sqlite.Conn, typeId int, c *gin.Context) ([]Item,error){
	_,span := dTracer.Start(c.Request.Context(),"type.items.db")
	defer span.End()
	query := `
		SELECT id, name, normalized_name 
		FROM item
		WHERE item.type_id = ?;`
	var items []Item
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			item := Item{
				Id: stmt.ColumnInt32(0),
				Name: stmt.ColumnText(1),
				NormalizedName: stmt.ColumnText(2),
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

func getBarrelsBySellerId(conn *sqlite.Conn,sellerId int32, c *gin.Context) ([]BarrelItem, error){
	_,span := dTracer.Start(c.Request.Context(),"sellers.barrels.db")
	defer span.End()
	query := `
		SELECT 
			b_item.id, b_item.item_id, b_item.barrel_id, b_item.seller_id, b_item.quantity, 
    		b_item.price, b_item.benefit_ratio, b_item.record_date,
			item.id, item.name, item.normalized_name, item.type_id
		FROM barrel_item as b_item
		JOIN seller ON b_item.seller_id = seller.id
		JOIN item ON b_item.item_id = item.id
		WHERE seller.id = ?;`
	var sellerBarrels []BarrelItem
	err := sqlitex.Execute(conn,query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			item := BarrelItem{
				Id: stmt.ColumnInt32(0),
				Item: Item{
					Id: stmt.ColumnInt32(8),
					Name: stmt.ColumnText(9),
					NormalizedName: stmt.ColumnText(10),
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

func getBarrelByCords(conn *sqlite.Conn, x int, y int, z int,c *gin.Context) (*int,error){
	_,span := dTracer.Start(c.Request.Context(),"barrel.get.coords.db")
	defer span.End()
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

func getOrCreateBarrel(conn *sqlite.Conn, x int, y int, z int, c *gin.Context) (*int,error) {
	_,span := dTracer.Start(c.Request.Context(),"barrel.get_or_create.db")
	defer span.End()
	barrelId, err := getBarrelByCords(conn,x,y,z,c)
	if barrelId == nil {
	  barrelId, err = createBarrel(conn,x,y,z,c)
	  if err != nil {
		return nil,err
	  }
	}
	return barrelId,nil
}

func createItemInBarrel(conn *sqlite.Conn, item DBItemInBarrel, c *gin.Context) error{
	_,span := dTracer.Start(c.Request.Context(),"item_in_barrel.create.db")
	defer span.End()
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

func createItemsSnapshot(conn *sqlite.Conn, barrelId int, c *gin.Context) (*int, error) {
	_,span := dTracer.Start(c.Request.Context(),"snapshots.create.db")
	defer span.End()
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

func getOrCreateItem(conn *sqlite.Conn, itemName string, c *gin.Context) (*int32, error) {
	_,span := dTracer.Start(c.Request.Context(),"item.get_or_create.db")
	defer span.End()
    normalizedName := normalizeString(itemName)
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
	formatedName := formatString(itemName)
    err = sqlitex.Execute(conn, insertQuery, &sqlitex.ExecOptions{
        Args: []any{formatedName, normalizedName},
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

func postItemsInBarrel(conn *sqlite.Conn,items ItemInBarrelPost, c *gin.Context) (error) {
	_,span := dTracer.Start(c.Request.Context(),"item_barrel.post.db")
	defer span.End()
	barrelId, err := getBarrelByCords(conn,*items.X,*items.Y,*items.Z,c)
	if err != nil{
		return err;
	}
	if barrelId == nil{
		return ErrNotFound;
	}

	snapshotId, err := createItemsSnapshot(conn, *barrelId,c)
	if err != nil {
		return nil
	}

	for _,v := range items.Items {
		id, err := getOrCreateItem(conn,v.Name,c)
		if err != nil {
			continue
		}
		err = createItemInBarrel(conn, DBItemInBarrel{ItemId: *id, SnapshotId: *snapshotId, Quantity: v.Quantity},c)
		if err != nil {
			return err
		}
	}
	return nil;
}

func getOrCreateSeller(conn *sqlite.Conn,sellerName string,c *gin.Context) (*int32, error) {
	_,span := dTracer.Start(c.Request.Context(),"seller.get_or_create.db")
	defer span.End()
	sellerId, err := getSellerByName(conn,sellerName,c)
	if sellerId == nil || err != nil {
	  id, err := getUserMinecraftUUID(sellerName)
	  sellerId, err = createSeller(conn,sellerName,&id,c)
	  if err != nil {
		return nil,err
	  }
	}
	return sellerId,nil
}

func getSellerByName(conn *sqlite.Conn,sellerName string, c *gin.Context) (*int32, error) {
	_,span := dTracer.Start(c.Request.Context(),"seller.search.name.db")
	defer span.End()
	normalizedName := normalizeStringWithNumbers(sellerName)
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

func createSeller(conn *sqlite.Conn, name string, minecraft_id *string, c *gin.Context) (*int32, error) {
	_,span := dTracer.Start(c.Request.Context(),"seller.post.db")
	defer span.End()
	normalized_username := normalizeStringWithNumbers(name)
	name = formatStringWithNumbers(name)
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

func createBarrel(conn *sqlite.Conn, x int, y int, z int, c *gin.Context) (*int,error) {
	_,span := dTracer.Start(c.Request.Context(),"barrel.post.db")
	defer span.End()
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
func createBarrelItems(conn *sqlite.Conn,barrels []AddBarrelStruct,c *gin.Context) (error){
	_,span := dTracer.Start(c.Request.Context(),"barrels.post_batch.db")
	defer span.End()
	if len(barrels) == 0 {
        return nil
    }
	var query strings.Builder
	query.WriteString(`
		INSERT INTO barrel_item (
			item_id, barrel_id, seller_id,
			quantity, price, benefit_ratio, barrel_text
		)
		VALUES
	`)
	args := make([]any, 0, len(barrels)*7)

	for i, barrel := range barrels {
		if i > 0 {
			query.WriteByte(',')
		}

		query.WriteString("(?, ?, ?, ?, ?, ?, ?)")

		args = append(args,
			*barrel.ItemId,
			*barrel.BarrelId,
			*barrel.SellerId,
			barrel.Quantity,
			barrel.Price,
			barrel.BenefitRatio,
			barrel.BarrelText,
		)
	}
	err := sqlitex.Execute(conn,query.String(),&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
//			id = stmt.ColumnInt32(0)
			return nil
		},
		Args: args,
	})
	if err != nil {
		return err
	}
	return nil
}

func createBarrelItem(conn *sqlite.Conn,itemId int32, barrelId int, sellerId int32,recognizedBarrel RecognizedBarrelItem,barrelText string,c *gin.Context) (*int32, error){
	_,span := dTracer.Start(c.Request.Context(),"barrels.post.db")
	defer span.End()
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

func getItemsSnapshotByQuery(conn *sqlite.Conn, query string, page int, pageSize int,c *gin.Context) (*PaginatedSnapshotSearchResponse,error){
	_,span := dTracer.Start(c.Request.Context(),"snapshots.search.db")
	defer span.End()
	normalized_query := normalizeString(query)
	var snapshots []ItemsSnapshotSearchResponse;
	offset := (page - 1) * pageSize
	var total int;
	sqlite_query := `
		SELECT id,x,y,z, items, record_date FROM (
			SELECT 
				snapshot.id, 
				barrel.x, 
				barrel.y, 
				barrel.z, 
				snapshot.record_date, 
				ROW_NUMBER() OVER (
					PARTITION by snapshot.barrel_id 
					ORDER BY snapshot.record_date DESC, snapshot.id DESC
				) AS snapshot_pos,
				(
                	SELECT json_group_array(
                	    json_object('name', sub_i.name, 'quantity', sub_ib.quantity)
                	)
                	FROM item_in_barrel AS sub_ib
                	JOIN item AS sub_i ON sub_i.id = sub_ib.item_id
                	WHERE sub_ib.snapshot_id = snapshot.id
            	) AS items
			FROM item
			JOIN item_in_barrel AS ib ON ib.item_id = item.id
			JOIN storage_snapshot AS snapshot ON snapshot.id = ib.snapshot_id
			JOIN barrel ON barrel.id = snapshot.barrel_id
			WHERE item.normalized_name LIKE '%' || ? || '%'
			GROUP BY snapshot.id
		)
		WHERE snapshot_pos = 1 
		LIMIT ?
		OFFSET ?
	`
	err := sqlitex.Execute(conn, sqlite_query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			record_date, err := time.Parse("2006-01-02 15:04:05",stmt.ColumnText(5))
			if err != nil {
				return err
			}
			itemsRaw := stmt.ColumnText(4)
			var items []ItemInBarrel
			if itemsRaw != "" {
				if err := json.Unmarshal([]byte(itemsRaw), &items); err != nil {
					return err
				}
			}
			snapshot := ItemsSnapshotSearchResponse{
				Id: stmt.ColumnInt(0),
				X: stmt.ColumnInt(1),
				Y: stmt.ColumnInt(2),
				Z: stmt.ColumnInt(3),
				Items: items,
				RecordDate: record_date,
			}
			snapshots = append(snapshots,snapshot)
			return nil
		},
		Args: []any{normalized_query,pageSize,offset},
	})
	if err != nil {
		return nil, err
	}

	sqlite_query = `
		SELECT COUNT(*) FROM (
			SELECT 
				snapshot.id, 
				ROW_NUMBER() OVER (
					PARTITION by snapshot.barrel_id 
					ORDER BY snapshot.record_date DESC, snapshot.id DESC
				) AS snapshot_pos
			FROM item
			JOIN item_in_barrel AS ib ON ib.item_id = item.id
			JOIN storage_snapshot AS snapshot ON snapshot.id = ib.snapshot_id
			JOIN barrel ON barrel.id = snapshot.barrel_id
			WHERE item.normalized_name LIKE '%' || ? || '%'
			GROUP BY snapshot.id
		)
		WHERE snapshot_pos = 1 
	`
	err = sqlitex.Execute(conn, sqlite_query,&sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			total = stmt.ColumnInt(0);
			return nil
		},
		Args: []any{normalized_query},
	})
	if err != nil {
		return nil, err
	}
	return &PaginatedSnapshotSearchResponse{Items: snapshots,Total: total,Page: page, Limit: pageSize}, nil
}

func getBarrelsByQuery(conn *sqlite.Conn, query string, page int, pageSize int,c *gin.Context) (*BarrelItemResponse, error) {
	_,span := dTracer.Start(c.Request.Context(),"barrels.search.db")
	defer span.End()
	normalized_query := normalizeString(query)
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
		ORDER BY benefit_ratio DESC, INSTR(normalized_name, ?) ASC, LENGTH(normalized_name) ASC
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

func getBarrelInfo(conn *sqlite.Conn, barrelId int,c *gin.Context) (*BarrelInfo,error) {
	_,span := dTracer.Start(c.Request.Context(),"barrels.info.db")
	defer span.End()
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