CREATE SEQUENCE wallets_id_seq;
CREATE FUNCTION wallets_next_id()
	RETURNS bigint
	LANGUAGE 'plpgsql'
AS $BODY$
DECLARE
	current_epoch bigint := 1314220021721;
	current_milliseconds bigint;
	shard_id int := 1;
	sequence_id bigint;
	result bigint := 0;
BEGIN
	SELECT nextval('wallets_id_seq') % 1024 INTO sequence_id;

	SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO current_milliseconds;
	result := (current_milliseconds - current_epoch) << 23;
	result := result | (shard_id << 10);
	result := result | (sequence_id);
	return result;
END;
$BODY$;

CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS wallets (
  id BIGINT PRIMARY KEY DEFAULT wallets_next_id(),
	address CITEXT NOT NULL,
	CONSTRAINT wallet_address UNIQUE (address),
	CONSTRAINT wallet_address_length CHECK (LENGTH(address) = 42)
);

CREATE INDEX wallet_address_idx ON wallets(address);

CREATE TABLE IF NOT EXISTS stars (
	id BIGINT PRIMARY KEY,
	name CITEXT NOT NULL,
	coordinates CHAR(19) UNIQUE NOT NULL,
	is_for_sale BOOLEAN NOT NULL,
	price_ether NUMERIC(77,18),
	date_created TIMESTAMP NOT NULL,
	owner_id BIGINT REFERENCES wallets(id) NOT NULL,
	CONSTRAINT star_name UNIQUE (name),
	CONSTRAINT star_name_length CHECK (LENGTH(name) >= 4 AND LENGTH(name) <= 32),
	CONSTRAINT star_coordinates UNIQUE (coordinates)
);

CREATE INDEX star_name_idx ON stars(name);
CREATE INDEX star_coordinates_idx ON stars(coordinates);

CREATE TABLE IF NOT EXISTS sales_history (
	date_sold TIMESTAMP,
	price_ether NUMERIC(30,18) NOT NULL,
	star_id BIGINT REFERENCES stars(id) NOT NULL,
	seller_id BIGINT REFERENCES wallets(id) NOT NULL,
	buyer_id BIGINT REFERENCES wallets(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS names_history (
	date_set TIMESTAMP NOT NULL,
	name CITEXT NOT NULL,
	star_id BIGINT REFERENCES stars(id) NOT NULL,
	owner_id INTEGER REFERENCES wallets(id) NOT NULL,
	CONSTRAINT star_name_history UNIQUE (name),
	CONSTRAINT star_name_history_length CHECK (LENGTH(name) >= 4 AND LENGTH(name) <= 32)
);
