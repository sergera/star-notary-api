package repositories

import (
	"database/sql"

	"github.com/sergera/star-notary-backend/internal/domain"
)

type StarRepository struct {
	conn *DBConnection
}

func NewStarRepository(conn *DBConnection) *StarRepository {
	return &StarRepository{conn}
}

func (sr *StarRepository) CreateStar(m domain.StarModel) error {
	tx, err := sr.conn.Session.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(
		`
		INSERT INTO stars (id, name, coordinates, is_for_sale, price_ether, date_created, owner_wallet_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
		m.TokenId, m.Name, m.Coordinates, m.IsForSale, m.Price, m.Date, m.Wallet.Id,
	); err != nil {
		return err
	}

	if _, err := tx.Exec(
		`
		INSERT INTO names_history (star_id, name, date_set, owner_wallet_id)
		VALUES ($1, $2, $3, $4)
		`,
		m.TokenId, m.Name, m.Date, m.Wallet.Id,
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (sr *StarRepository) SetPrice(m domain.StarModel) error {
	if _, err := sr.conn.Session.Exec(
		`
		UPDATE stars
		SET is_for_sale = $2, price_ether = $3
		WHERE id = $1
		`,
		m.TokenId, m.IsForSale, m.Price,
	); err != nil {
		return err
	}

	return nil
}

func (sr *StarRepository) SetName(m domain.StarModel) error {
	tx, err := sr.conn.Session.Begin()
	if err != nil {
		return err
	}

	if err := tx.QueryRow(
		`
		UPDATE stars
		SET name = $2
		WHERE id = $1
		RETURNING owner_wallet_id
		`,
		m.TokenId, m.Name,
	).Scan(&m.Wallet.Id); err != nil {
		return err
	}

	if _, err := tx.Exec(
		`
		INSERT INTO names_history (star_id, name, date_set, owner_wallet_id)
		VALUES ($1, $2, $3, $4)
		`,
		m.TokenId, m.Name, m.Date, m.Wallet.Id,
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (sr *StarRepository) GetStarRange(m domain.StarRangeModel) ([]domain.StarModel, error) {
	var rows *sql.Rows
	var err error
	if m.OldestFirst {
		rows, err = sr.conn.Session.Query(
			`
			SELECT stars.id, stars.name, stars.coordinates, stars.is_for_sale, stars.price_ether, stars.date_created, wallets.address
			FROM stars, wallets
			WHERE stars.owner_wallet_id = wallets.id
			AND stars.id >= $1
			AND stars.id <= $2
			ORDER BY stars.id ASC
			`,
			m.Start, m.End,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	} else {
		rows, err = sr.conn.Session.Query(
			`
			WITH last_star AS (
				SELECT id
				FROM stars
				ORDER BY id DESC
				LIMIT 1
			)
			
			SELECT stars.id, stars.name, stars.coordinates, stars.is_for_sale, stars.price_ether, stars.date_created, wallets.address
			FROM stars, wallets
			WHERE stars.owner_wallet_id = wallets.id
			AND stars.id <= ((SELECT id FROM last_star) - ($1 - 1))
			AND stars.id >= ((SELECT id FROM last_star) - ($2 - 1))
			ORDER BY stars.id DESC
			`,
			m.Start, m.End,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	}

	var stars []domain.StarModel
	for rows.Next() {
		st := domain.StarModel{
			Wallet: &domain.WalletModel{},
		}
		err := rows.Scan(&st.TokenId, &st.Name, &st.Coordinates, &st.IsForSale, &st.Price, &st.Date, &st.Wallet.Address)
		if err != nil {
			return nil, err
		}
		stars = append(stars, st)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stars, nil
}
