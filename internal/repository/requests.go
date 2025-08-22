package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"test-task/internal/model"
)

var ErrOrderIter = errors.New("error in orders iterations")
var ErrOrderNotFound = errors.New("error not found in DB")

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetLastNOrders(limit int) ([]model.Order, error) {
	mainQuery := fmt.Sprintf(`
		SELECT
			o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
			o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
			d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
			p.transaction_id, p.request_id, p.currency, p.provider, p.amount,
			p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders AS o
		JOIN deliveries AS d ON o.delivery_id = d.id
		JOIN payments AS p ON o.payment_transaction_id = p.transaction_id
		ORDER BY o.date_created DESC
		LIMIT %d;`, limit)

	rows, err := r.db.Query(mainQuery)
	if err != nil {
		return nil, fmt.Errorf("error in order requests: %w", err)
	}
	defer rows.Close()

	ordersMap := make(map[string]*model.Order)

	for rows.Next() {
		var o model.Order
		err := rows.Scan(
			&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSign,
			&o.CustomerID, &o.DeliveryService, &o.ShardKey, &o.SmID, &o.DateCreated, &o.OofShard,
			&o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip, &o.Delivery.City, &o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email,
			&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency, &o.Payment.Provider, &o.Payment.Amount,
			&o.Payment.PaymentDT, &o.Payment.Bank, &o.Payment.DeliveryCost, &o.Payment.GoodsTotal, &o.Payment.CustomFee,
		)
		if err != nil {
			log.Printf("Error in scan orders: %v", err)
			continue
		}

		o.Items = []model.Item{}
		ordersMap[o.OrderUID] = &o
	}
	if err = rows.Err(); err != nil {
		return nil, ErrOrderIter
	}

	itemsQuery := `
		SELECT
			i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale,
			i.size, i.total_price, i.nm_id, i.brand, i.status
		FROM items AS i
		JOIN order_items AS oi ON i.chrt_id = oi.chrt_id
		WHERE oi.order_uid = $1;`

	stmt, err := r.db.Prepare(itemsQuery)
	if err != nil {
		return nil, errors.New("error in ready orders")
	}
	defer stmt.Close()

	for uid, order := range ordersMap {
		itemRows, err := stmt.Query(uid)
		if err != nil {
			log.Printf("Erroe in request for orders %s: %v", uid, err)
			continue
		}

		for itemRows.Next() {
			var item model.Item
			err := itemRows.Scan(
				&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale,
				&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
			)
			if err != nil {
				log.Printf("Error in scan %s: %v", uid, err)
				continue
			}
			order.Items = append(order.Items, item)
		}
		itemRows.Close()
	}

	orders := make([]model.Order, 0, len(ordersMap))
	for _, order := range ordersMap {
		orders = append(orders, *order)
	}

	return orders, nil
}

func (r *OrderRepository) GetByUID(uid string) (model.Order, error) {
	var o model.Order

	mainQuery := `
		SELECT
			o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
			o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
			d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
			p.transaction_id, p.request_id, p.currency, p.provider, p.amount,
			p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders AS o
		JOIN deliveries AS d ON o.delivery_id = d.id
		JOIN payments AS p ON o.payment_transaction_id = p.transaction_id
		WHERE o.order_uid = $1;`

	err := r.db.QueryRow(mainQuery, uid).Scan(
		&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSign,
		&o.CustomerID, &o.DeliveryService, &o.ShardKey, &o.SmID, &o.DateCreated, &o.OofShard,
		&o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip, &o.Delivery.City, &o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email,
		&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency, &o.Payment.Provider, &o.Payment.Amount,
		&o.Payment.PaymentDT, &o.Payment.Bank, &o.Payment.DeliveryCost, &o.Payment.GoodsTotal, &o.Payment.CustomFee,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, ErrOrderNotFound
		}
		return model.Order{}, fmt.Errorf("error scanning order: %w", err)
	}

	itemsQuery := `
		SELECT
			i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale,
			i.size, i.total_price, i.nm_id, i.brand, i.status
		FROM items AS i
		JOIN order_items AS oi ON i.chrt_id = oi.chrt_id
		WHERE oi.order_uid = $1;`

	rows, err := r.db.Query(itemsQuery, uid)
	if err != nil {
		return model.Order{}, fmt.Errorf("error querying items for order %s: %w", uid, err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			log.Printf("Error scanning item for order %s: %v", uid, err)
			continue
		}
		o.Items = append(o.Items, item)
	}
	if err = rows.Err(); err != nil {
		return model.Order{}, fmt.Errorf("error iterating items for order %s: %w", uid, err)
	}

	return o, nil
}

func (r *OrderRepository) SaveOrder(order *model.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	deliveryQuery := `INSERT INTO deliveries (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var deliveryID int
	err = tx.QueryRow(deliveryQuery, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email).Scan(&deliveryID)
	if err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	paymentQuery := `INSERT INTO payments (transaction_id, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = tx.Exec(paymentQuery, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	orderQuery := `INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_transaction_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err = tx.Exec(orderQuery, order.OrderUID, order.TrackNumber, order.Entry, deliveryID, order.Payment.Transaction, order.Locale, order.InternalSign, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	itemQuery := `INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT (chrt_id) DO NOTHING`
	orderItemQuery := `INSERT INTO order_items (order_uid, chrt_id) VALUES ($1, $2)`

	for _, item := range order.Items {
		_, err = tx.Exec(itemQuery, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("failed to insert item with chrt_id %d: %w", item.ChrtID, err)
		}
		_, err = tx.Exec(orderItemQuery, order.OrderUID, item.ChrtID)
		if err != nil {
			return fmt.Errorf("failed to insert into order_items for order %s and chrt_id %d: %w", order.OrderUID, item.ChrtID, err)
		}
	}

	return tx.Commit()
}
