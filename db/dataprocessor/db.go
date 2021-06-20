package dataprocessor

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func NewService(db *sql.DB) (*Service, error) {
	return &Service{db: db}, nil
}

func ConnectToSQL(address string) (*sql.DB, error) {
	if address == "" {
		address = "root:password@tcp(localhost:3306)/mtest"
	}

	db, err := sql.Open("mysql", address)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

const (
	disableGroupBy   = `SET GLOBAL sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));`
	initRegActsTable = `create table if not exists reg_acts
(
	id bigint auto_increment,
	mid varchar(300), 
	doc_id varchar(300),
	doc_text text, 
	doc_name varchar(300),
	doc_type varchar(300),
	constraint businesses_pk
		primary key (id)
);
`
	initBusinessesTable = `create table if not exists businesses
(
	id bigint auto_increment,
	name varchar(300) not null,
	constraint businesses_pk
		primary key (id)
);
`
	initGovernmentsTable = `create table if not exists govs(
	id bigint auto_increment,
	gov_id int auto_increment,
	gov_name varchar(300) not null,
	constraint businesses_pk
		primary key (id)
);
`
	initAdmActionsTable = `create table if not exists adm_actions(
	id bigint auto_increment,
	act_id int auto_increment,
	act_name varchar(300) not null,
	constraint businesses_pk
		primary key (id)
);
`
	initRegionsTable = `create table if not exists regions(
	id bigint auto_increment,
	reg_id int auto_increment,
	reg_name varchar(300) not null,
	constraint businesses_pk
		primary key (id)
);
`
)

func (mt Service) Init() error {
	qs := []string{disableGroupBy, initRegActsTable, initBusinessesTable, initRegionsTable, initGovernmentsTable, initAdmActionsTable, initSynonyms}

	for _, q := range qs {
		stmt, err := mt.db.Prepare(q)
		if err != nil {
			return err
		}

		defer func() {
			if err := stmt.Close(); err != nil {
				log.Error(err)
			}
		}()
		if _, err = stmt.Exec(); err != nil {
			return err
		}
	}
	return nil
}
