package repositories

import (
	"github.com/sergera/star-notary-backend/internal/domain"
)

type WalletRepositoryInterface interface {
	CreateWallet(m *domain.WalletModel) error
}

type WalletRepository struct {
	conn *DBConnection
}

func NewWalletRepository(conn *DBConnection) *WalletRepository {
	return &WalletRepository{conn}
}

func (wr *WalletRepository) CreateWallet(m *domain.WalletModel) error {
	err := wr.conn.Session.QueryRow(
		`
		WITH inserted_wallet AS (
			INSERT INTO wallets (address)
			VALUES ($1)
			ON CONFLICT (address) DO NOTHING
			RETURNING id
		)

		SELECT * FROM inserted_wallet
		UNION
		SELECT id FROM wallets WHERE address=$1
		`,
		m.Address,
	).Scan(&m.Id)

	if err != nil {
		return err
	}

	return nil
}
