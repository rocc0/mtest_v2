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
	initMTESTsTable = `create table if not exists mtests (
	id bigint auto_increment,
	mid varchar(100) not null,
	name varchar(500) not null,
	region int not null,
	govern int not null,
	business int not null,
	calculations text not null,
	calc_type int default 0,
	calc_data varchar(1000) default '{}',
	executors varchar(500) default '{}',
	developer varchar(100),
	dev_mid varchar(100),
	pub_data date,
	author varchar(100),
	tags varchar(500) default '{}',
	math_result int default 0,
	corr_result int default 0,
	constraint mtests_pk
		primary key (id))`
)

func (mt Service) Init() error {
	qs := []string{disableGroupBy, initRegActsTable, initBusinessesTable,
		initRegionsTable, initGovernmentsTable, initAdmActionsTable, initSynonyms, initMTESTsTable}

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
