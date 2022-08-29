package repositories

import (
	"database/sql"
	"strconv"

	"github.com/sergera/star-notary-backend/internal/domain"
)

type StarRepository struct {
	conn *DBConnection
}

func NewStarRepository(conn *DBConnection) *StarRepository {
	return &StarRepository{conn}
}

func (sr *StarRepository) CreateStar(m domain.StarModel) error {
	if _, err := sr.conn.Session.Exec(
		`
		INSERT INTO stars (id, name, coordinates, is_for_sale, price_ether, date_created, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
		m.TokenId, m.Name, m.Coordinates, m.IsForSale, m.Price, m.Date, m.Wallet.Id,
	); err != nil {
		return err
	}

	return nil
}

func (sr *StarRepository) GetStarRange(m domain.StarRangeModel) ([]domain.StarModel, error) {
	tx, err := sr.conn.Session.Begin()
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	if m.OldestFirst {
		rows, err = tx.Query(
			`
			SELECT stars.id, stars.name, stars.coordinates, stars.is_for_sale, stars.price_ether, stars.date_created, wallets.address, wallets.id
			FROM stars, wallets
			WHERE stars.owner_id = wallets.id
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
		var maxIdString string
		tx.QueryRow(
			`
			SELECT id
			FROM stars
			ORDER BY id DESC
			LIMIT 1
			`,
		).Scan(&maxIdString)

		if maxIdString == "" {
			var emptyStarSlice []domain.StarModel
			return emptyStarSlice, nil
		}

		maxId, err := strconv.ParseInt(maxIdString, 10, 64)
		if err != nil {
			return nil, err
		}

		start, err := strconv.ParseInt(m.Start, 10, 64)
		if err != nil {
			return nil, err
		}

		end, err := strconv.ParseInt(m.End, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = tx.Query(
			`
			SELECT stars.id, stars.name, stars.coordinates, stars.is_for_sale, stars.price_ether, stars.date_created, wallets.address, wallets.id
			FROM stars, wallets
			WHERE stars.owner_id = wallets.id
			AND stars.id >= $1
			AND stars.id <= $2
			ORDER BY stars.id DESC
			`,
			maxId-(end-1), maxId-(start-1),
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
		err := rows.Scan(&st.TokenId, &st.Name, &st.Coordinates, &st.IsForSale, &st.Price, &st.Date, &st.Wallet.Address, &st.Wallet.Id)
		if err != nil {
			return nil, err
		}
		stars = append(stars, st)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return stars, nil
}
