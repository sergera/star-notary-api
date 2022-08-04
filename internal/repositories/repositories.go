package repositories

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/sergera/star-notary-backend/internal/models"
)

type StarRepository struct {
	psqlConfig string
	db         *sql.DB
}

func NewStarRepository(host string, port string, dbname string, user string, password string, sslmode bool) *StarRepository {
	var sslconfig string
	if sslmode {
		sslconfig = "enable"
	} else {
		sslconfig = "disable"
	}

	return &StarRepository{
		fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
			host, port, dbname, user, password, sslconfig),
		nil,
	}
}

func (sr *StarRepository) Open() {
	db, err := sql.Open("postgres", sr.psqlConfig)
	if err != nil {
		panic(err)
	}

	sr.db = db
}

func (sr *StarRepository) Close() {
	sr.db.Close()
}

func (sr *StarRepository) InsertWalletIfAbsent(m models.StarModel) error {
	tx, err := sr.db.Begin()
	if err != nil {
		return err
	}

	var ownerId string
	tx.QueryRow(
		`
		SELECT id
		FROM wallets
		WHERE address=$1
		`,
		m.Owner,
	).Scan(&ownerId)

	if ownerId == "" {
		_, err = tx.Exec(
			`
			INSERT INTO wallets(address)
			VALUES ($1)
			ON CONFLICT DO NOTHING
			`,
			m.Owner,
		)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (sr *StarRepository) CreateStar(m models.StarModel) error {
	tx, err := sr.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`
		INSERT INTO stars (id, name, coordinates, is_for_sale, price_ether, date_created, owner_id)
		SELECT $1, $2, $3, $4, $5, $6, id
		FROM wallets
		WHERE address=$7
		`,
		m.TokenId, m.Name, m.Coordinates, false, "0", m.Date, m.Owner,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (sr *StarRepository) GetStars(m models.StarRangeModel) ([]models.StarModel, error) {
	rows, err := sr.db.Query(
		`
		SELECT stars.id, stars.name, stars.coordinates, stars.is_for_sale, stars.price_ether, stars.date_created, wallets.address
		FROM stars, wallets
		WHERE stars.owner_id = wallets.id
		AND stars.id >= $1
		AND stars.id <= $2
		ORDER BY stars.id ASC
		`,
		m.FirstId, m.LastId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stars []models.StarModel
	for rows.Next() {
		var st models.StarModel
		err := rows.Scan(&st.TokenId, &st.Name, &st.Coordinates, &st.IsForSale, &st.Price, &st.Date, &st.Owner)
		if err != nil {
			return nil, err
		}
		stars = append(stars, st)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return stars, nil
}
