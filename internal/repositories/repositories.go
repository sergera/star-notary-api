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

func (sr *StarRepository) Create(m models.StarModel) error {
	fail := func(err error) error {
		return fmt.Errorf("Create Star: %v", err)
	}

	tx, err := sr.db.Begin()
	if err != nil {
		return fail(err)
	}

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
		return fail(err)
	}

	_, err = tx.Exec(
		`
		INSERT INTO stars (id, name, coordinates, is_for_sale, price_ether, date_created, owner_id)
		SELECT $1, $2, $3, $4, $5, $6, id
		FROM wallets 
		WHERE address=$7
		`,
		m.TokenId, m.Name, m.Coordinates, false, nil, m.Date, m.Owner,
	)
	if err != nil {
		log.Println(err)
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		log.Println(err)
		return fail(err)
	}

	return nil
}
