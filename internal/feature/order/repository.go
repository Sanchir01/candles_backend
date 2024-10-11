package order

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		primaryDB: db,
	}
}

func (r *Repository) AllOrders(ctx context.Context) ([]model.Orders, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query, _, err := sq.Select("id, status,user_id,total_amount,created_at,updated_at, version").From("public.orders").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := make([]model.Orders, 0)

	for rows.Next() {
		var order model.Orders
		if err := rows.Scan(&order.ID, &order.Status, &order.UserID, &order.TotalAmount, &order.CreatedAt, &order.UpdatedAt, &order.Version); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, err
}

func (r *Repository) OrderByUserId(ctx context.Context, id uuid.UUID) (*model.Orders, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query, arg, err := sq.Select("id, status,user_id,total_amount,crated_at,updated_at, version").From("public.orders").Where(sq.Eq{"user_id": id}).ToSql()
	var order DBOrders
	if err := conn.QueryRow(ctx, query, arg...).
		Scan(&order.ID, &order.Status, &order.UserID, &order.TotalAmount, &order.CreatedAt, &order.UpdatedAt, &order.Version); err != nil {
		return nil, err
	}
	return (*model.Orders)(&order), err
}

func (r *Repository) OrderById(ctx context.Context, id uuid.UUID) (*model.Orders, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query, arg, err := sq.Select("id, status,user_id,total_amount,crated_at,updated_at, version").From("public.orders").Where(sq.Eq{"id": id}).ToSql()
	var order DBOrders
	if err := conn.QueryRow(ctx, query, arg...).
		Scan(&order.ID, &order.Status, &order.UserID, &order.TotalAmount, &order.CreatedAt, &order.UpdatedAt, &order.Version); err != nil {
		return nil, err
	}
	return (*model.Orders)(&order), err
}

func (r *Repository) CreateOrder(ctx context.Context, id uuid.UUID, status string, totalAmount int, tx pgx.Tx) (uuid.UUID, error) {
	query, args, err := sq.
		Insert("orders").
		Columns("status", "user_id", "total_amount").
		Values(status, id, totalAmount).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	var orderID uuid.UUID
	if err = tx.QueryRow(ctx, query, args...).Scan(&orderID); err != nil {
		return uuid.Nil, err
	}

	return orderID, nil
}

func (r *Repository) CreateOrderItem(
	ctx context.Context, orderID uuid.UUID, productIDs []uuid.UUID, quantity []int, price []int, tx pgx.Tx,
) ([]uuid.UUID, error) {

	queryBuilder := sq.Insert("order_items").Columns("order_id", "product_id", "quantity", "price")

	for i := 0; i < len(productIDs); i++ {
		queryBuilder = queryBuilder.Values(orderID, productIDs[i], quantity[i], price[i])
	}
	query, args, err := queryBuilder.Suffix("RETURNING id").PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}
	var ids []uuid.UUID

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err.Error())
		return nil, err
	}
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	log.Printf("ids items", ids)
	return ids, err
}
