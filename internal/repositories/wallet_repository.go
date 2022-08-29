package repositories

import (
	"github.com/sergera/star-notary-backend/internal/domain"
)

type WalletRepository struct {
	conn *DBConnection
}

func NewWalletRepository(conn *DBConnection) *WalletRepository {
	return &WalletRepository{conn}
}

func (wr *WalletRepository) InsertWalletIfAbsent(m *domain.WalletModel) error {
	tx, err := wr.conn.Session.Begin()
	if err != nil {
		return err
	}

	tx.QueryRow(
		`
		SELECT id
		FROM wallets
		WHERE address=$1
		`,
		m.Address,
	).Scan(&m.Id)

	if m.Id == "" {
		if err := tx.QueryRow(
			`
			INSERT INTO wallets(address)
			VALUES ($1)
			RETURNING id
			`,
			m.Address,
		).Scan(&m.Id); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
