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
	gov_id int,
	gov_name varchar(300) not null,
	constraint govs_pk
		primary key (id)
);
`
	initGovernments = `INSERT INTO govs (id, gov_id, gov_name) VALUES
(1, 1, 'Міністерство аграрної політики та продовольства України'),
(2, 2, 'Державна ветеринарна та фітосанітарна служба України'),
(3, 3, 'Державне агентство рибного господарства України'),
(4, 4, 'Міністерство внутрішніх справ України '),
(5, 5, 'Державна служба України з надзвичайних ситуацій'),
(6, 6, 'Міністерство екології та природних ресурсів України'),
(7, 7, 'Державна служба геології та надр України'),
(8, 8, 'Міністерство економічного розвитку і торгівлі України '),
(9, 9, 'Міністерство енергетики та вугільної промисловості України '),
(10, 10, 'Міністерство інфраструктури України '),
(11, 11, 'Міністерство культури України'),
(12, 12, 'Міністерство молоді та спорту України'),
(13, 13, 'Міністерство оборони України '),
(14, 14, 'Міністерство освіти і науки України'),
(15, 15, 'Міністерство охорони здоров\'\'я України'),
(16, 16, 'Міністерство регіонального розвитку, будівництва та житлово-комунального господарства України'),
(17, 17, 'Державне агентство земельних ресурсів України (NULL,2014)  / Державна служба України з питань геодезії, картографії та кадастру (NULL,2015)'),
(18, 18, 'Державне агентство з енергоефективності та енергозбереження України'),
(19, 19, 'Міністерство соціальної політики України '),
(20, 20, 'Міністерство фінансів України '),
(21, 21, 'Державна фіскальна служба України'),
(22, 22, 'Міністерство юстиції України'),
(23, 23, 'Антимонопольний комітет України'),
(24, 24, 'Державний комітет телебачення і радіомовлення України'),
(25, 25, 'Фонд державного майна України'),
(26, 26, 'Державна служба спеціального зв\'язку та захисту інформації України'),
(27, 27, 'Державна інспекція ядерного регулювання України'),
(28, 28, 'Національна комісія, що здійснює державне регулювання у сферах енергетики та комунальних послуг'),
(29, 29, 'Державна комісія з регулювання ринків фінансових послуг України'),
(30, 30, 'Національна комісія з цінних паперів та фондового ринку'),
(31, 31, 'Національна комісія, що здійснює державне регулювання у сфері зв’язку та інформатизації'),
(32, 32, 'Державна інспекція сільського господарства України'),
(33, 33, 'Державна архітектурно-будівельна інспекція України'),
(34, 34, 'Державна служба інтелектуальної власності України'),
(35, 35, 'Державна фінансова інспекція України'),
(36, 36, 'Національна гвардія України'),
(37, 37, 'Пенсійний фонд України'),
(38, 38, 'Державна служба статистики України'),
(39, 39, 'Державна прикордонна служба України'),
(40, 40, 'Державна аудиторська служба України'),
(41, 41, 'Державна авіаційна служба України'),
(42, 42, 'Державна служба України з безпеки на транспорті'),
(43, 43, 'Державна служба гірничого нагляду та промислової безпеки України'),
(44, 44, 'Державної служба України з лікарських засобів та контролю за наркотиками'),
(45, 45, 'Державна служба України з питань праці'),
(46, 46, 'Служба безпеки України');`

	initAdmActionsTable = `create table if not exists adm_actions(
	id bigint auto_increment,
	act_id int,
	act_name varchar(300) not null,
	constraint businesses_pk
		primary key (id)
);
`
	initAdmActions = `INSERT INTO adm_actions (id, act_id, act_name) VALUES
(1, 1, 'Ознайомлення з інформаційними вимогами'),
(2, 2, 'Навчання та інструктаж керівників і працівників щодо інформаційних вимог, підготовка внутрішніх процедур (розподіл обов’язків, підготовка інструкцій тощо)'),
(3, 3, 'Пошук та отримання потрібної інформації з вже існуючих документів (внутрішні процедури)'),
(4, 4, 'Підготовка нових даних/документів (коригування існуючих) та отримання документів (дозволів, погоджень тощо) від третіх сторін (зовнішні процедури)'),
(5, 5, 'Розробка інформаційних матеріалів (наприклад, інформаційних листків тощо)'),
(6, 6, 'Заповнення бланків, форм, таблиць (у тому числі облікової інформації)'),
(7, 7, 'Проведення нарад, зустрічей (внутрішніх/зовнішніх з аудитором, юристом тощо)'),
(8, 8, 'Зовнішні перевірки та контрольні заходи'),
(9, 9, 'Копіювання (звітів тощо) документів, реєстрація.'),
(10, 10, 'Подання інформації до відповідних органів (наприклад, надсилання або доставка відомостей до органів влади)'),
(11, 11, 'Зберігання інформації (архівування тощо)'),
(12, 12, 'Придбання обладнання та витратних матеріалів, що використовуються саме для виконання інформаційних вимог'),
(13, 15, 'Інше');`

	initRegionsTable = `create table if not exists regions(
	id bigint auto_increment,
	reg_id int,
	reg_name varchar(300) not null,
	constraint businesses_pk
		primary key (id)
);
`

	initRegions = `INSERT INTO regions (id, reg_id, reg_name) VALUES
(1, 1, 'АР Крим'),
(2, 2, 'Вінницька область'),
(3, 3, 'Волинська область'),
(4, 4, 'Дніпропетровська область'),
(5, 5, 'Донецька область'),
(6, 6, 'Житомирська область'),
(7, 7, 'Закарпатська область'),
(8, 8, 'Запорізька область'),
(9, 9, 'Івано-Франківська область'),
(10, 10, 'Київська область'),
(11, 11, 'Кіровоградська область'),
(12, 12, 'Луганська область'),
(13, 13, 'Львівська область'),
(14, 14, 'м. Київ'),
(15, 15, 'Миколаївська область'),
(16, 16, 'Одеська область'),
(17, 17, 'Полтавська область'),
(18, 18, 'Рівненська область'),
(19, 19, 'Сумська область'),
(20, 20, 'Тернопільська область'),
(21, 21, 'Харківська область'),
(22, 22, 'Херсонська область'),
(23, 23, 'Хмельницька область'),
(24, 24, 'Черкаська область'),
(25, 25, 'Чернівецька область'),
(26, 26, 'Чернігівська область');
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

	initDefaultMTESTS = `
INSERT INTO mtests (id, mid, name, region, govern, calculations, calc_type, calc_data, executors, developer, dev_mid, author, business) VALUES
(109, '287251ad-8125-4a5d-80da-c9c21fdc8367', 'Тестова назва', 10, 7, '{\"1\":[{\"type\":\"container\",\"id\":3,\"columns\":[[{\"type\":\"itemplus\",\"id\":3,\"columns\":[[{\"type\":\"item\",\"id\":3,\"name\":\"Додати дію\",\"subsum\":null},{\"type\":\"item\",\"id\":6,\"name\":\"Додати дію\",\"subsum\":null}]],\"name\":\"Додати складову інф. вимоги\"}]],\"name\":\"Додати інф. вимогу\",\"contsub\":0,\"awgsub\":null},{\"type\":\"container\",\"id\":null,\"columns\":[[{\"type\":\"itemplus\",\"id\":4,\"columns\":[[{\"type\":\"item\",\"id\":3,\"name\":\"Додати дію\",\"subsum\":null},{\"type\":\"item\",\"id\":4,\"name\":\"Додати дію\",\"subsum\":null}]],\"name\":\"Додати складову інф. вимоги\"}]],\"name\":\"Додати інф. вимогу\",\"contsub\":0,\"awgsub\":null}]}', 1, '{}', '{\"d53467b6-5859-4e6c-aa00-611d08166425\":{\"email\":\"vlad.kotlyarenko@gmail.com\",\"mid\":\"d53467b6-5859-4e6c-aa00-611d08166425\",\"checked\":true}}', NULL, NULL, 'vk@clc.com.ua', 0),
(117, '7751f388-e0ff-40f9-b7e7-75f05cb8fb15', 'Вирубка корупціонерів', 10, 9, '{\"1\":[{\"type\":\"container\",\"id\":3,\"columns\":[[{\"type\":\"itemplus\",\"id\":3,\n                    \"columns\":[[{\"type\":\"item\",\"id\":3,\"name\":\"Додати дію\",\"subsum\":0},{\"type\":\"item\",\"id\":6,\"name\":\"Додати дію\",\"subsum\":0}]],\n                    \"name\":\"Додати складову інф. вимоги\"}]],\"name\":\"Додати інф. вимогу\",\"contsub\":0},\n                {\"type\":\"container\",\"id\":null,\"columns\":[[{\"type\":\"itemplus\",\"id\":4,\"columns\":[[{\"type\":\"item\",\"id\":3,\"name\":\"Додати дію\",\"subsum\":0},\n                            {\"type\":\"item\",\"id\":4,\"name\":\"Додати дію\",\"subsum\":0}]],\"name\":\"Додати складову інф. вимоги\"}]],\"name\":\"Додати інф. вимогу\",\"contsub\":0}]}', 0, '{}', '{}', NULL, NULL, 'vk@clc.com.ua', 0),
(125, 'd53467b6-5859-4e6c-aa00-611d08166425', 'Тестова назва', 10, 7, '{\"1\":[{\"type\":\"container\",\"id\":3,\"columns\":[[{\"type\":\"itemplus\",\"id\":3,\n                    \"columns\":[[{\"type\":\"item\",\"id\":3,\"name\":\"Додати дію\",\"subsum\":0},{\"type\":\"item\",\"id\":6,\"name\":\"Додати дію\",\"subsum\":0}]],\n                    \"name\":\"Додати складову інф. вимоги\"}]],\"name\":\"Додати інф. вимогу\",\"contsub\":0},\n                {\"type\":\"container\",\"id\":null,\"columns\":[[{\"type\":\"itemplus\",\"id\":4,\"columns\":[[{\"type\":\"item\",\"id\":3,\"name\":\"Додати дію\",\"subsum\":0},\n                            {\"type\":\"item\",\"id\":4,\"name\":\"Додати дію\",\"subsum\":0}]],\"name\":\"Додати складову інф. вимоги\"}]],\"name\":\"Додати інф. вимогу\",\"contsub\":0}]}', 3, '{}', '{}', 'vk@clc.com.ua', '287251ad-8125-4a5d-80da-c9c21fdc8367', 'vlad.kotlyarenko@gmail.com', 0);
`

	initUsersTable = `CREATE TABLE  if not exists users (
  id bigint UNSIGNED NOT NULL,
  name varchar(100) NOT NULL,
  surename varchar(20) NOT NULL,
  email varchar(100) DEFAULT NULL,
  password varchar(100) NOT NULL,
  rights varchar(100) NOT NULL DEFAULT '1',
  activated int NOT NULL DEFAULT '0',
  records text
);`

	initUsers = `INSERT INTO users (id, name, surename, email, password, rights, activated, records) VALUES
(1, 'Владислав', 'Котляренко', 'vk@clc.com.ua', '$2a$10$8rFbqXKPyKo2kj9IfbapiuwWPgoGeHb.TpnVK0uAxWqKDjxqodtte', '1', 1, '{\"287251ad-8125-4a5d-80da-c9c21fdc8367\":{\"id\":\"287251ad-8125-4a5d-80da-c9c21fdc8367\",\"name\":\"Тестова назва\",\"region\":10,\"govern\":7,\"calc_type\":1,\"developer\":\"\",\"dev_mid\":\"\",\"executors\":{\"d53467b6-5859-4e6c-aa00-611d08166425\":{\"email\":\"vlad.kotlyarenko@gmail.com\",\"mid\":\"d53467b6-5859-4e6c-aa00-611d08166425\",\"checked\":true}}},\"7751f388-e0ff-40f9-b7e7-75f05cb8fb15\":{\"id\":\"7751f388-e0ff-40f9-b7e7-75f05cb8fb15\",\"name\":\"Вирубка корупціонерів\",\"region\":10,\"govern\":9,\"calc_type\":0,\"developer\":\"\",\"dev_mid\":\"\",\"executors\":{}}}'),
(33, 'Владислав', 'Котляренко', 'vlad.kotlyarenko@gmail.com', '$2a$10$bNBC1L3mxl6B76yGjcRxVu2MXxh8m9CSDDYZyIyPgW2jCDvvGKilG', '1', 1, '{\"d53467b6-5859-4e6c-aa00-611d08166425\":{\"id\":\"d53467b6-5859-4e6c-aa00-611d08166425\",\"name\":\"Тестова назва\",\"region\":10,\"govern\":7,\"calc_type\":3,\"developer\":\"vk@clc.com.ua\",\"dev_mid\":\"287251ad-8125-4a5d-80da-c9c21fdc8367\",\"executors\":null}}');
`
)

func (mt Service) Init() error {
	qs := []string{
		disableGroupBy, initRegActsTable,
		initBusinessesTable, initGlobalSynonyms,
		initRegionsTable, initGovernmentsTable, initAdmActionsTable,
		initSynonyms, initMTESTsTable, initUsersTable}

	for _, q := range qs {
		aff, err := exec(q, mt.db)
		if err != nil {
			return err
		}
		if aff != 0 {
			switch {
			case q == initGovernmentsTable:
				if _, err := exec(initGovernments, mt.db); err != nil {
					return err
				}
			case q == initAdmActionsTable:
				if _, err := exec(initAdmActions, mt.db); err != nil {
					return err
				}
			case q == initRegionsTable:
				if _, err := exec(initRegions, mt.db); err != nil {
					return err
				}
			case q == initMTESTsTable:
				if _, err := exec(initDefaultMTESTS, mt.db); err != nil {
					return err
				}
			case q == initUsersTable:
				if _, err := exec(initUsers, mt.db); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func exec(q string, client *sql.DB) (int64, error) {
	stmt, err := client.Prepare(q)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()
	res, err := stmt.Exec()
	if err != nil {
		log.Error(err, q)
		return 0, err
	}
	return res.RowsAffected()
}
