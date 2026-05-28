package dbsrv

import (
	"database/sql"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
	"codeberg.org/reiver/go-log"
)

type PersonRow struct {
	ID       int64
	Name     string
	Title    string
	Company  string
	FediID   string
	Note     string
	Favorite bool
}

type MessageRow struct {
	ID        int64
	PersonID  *int64
	GroupID   *int64
	FromMe    bool
	Sender    string
	Text      string
	Timestamp string
}

func LoadPeople(logger log.Logger, db *sql.DB) ([]PersonRow, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return nil, erorr.Wrap(erorr.Error("nil db"), "failed to load people")
	}

	rows, err := db.Query(`SELECT id, name, title, company, fedi_id, note, favorite FROM people ORDER BY id`)
	if nil != err {
		log.Error(field.S("failed to load people"), field.E(err))
		return nil, erorr.Wrap(err, "failed to load people")
	}
	defer rows.Close()

	var people []PersonRow
	for rows.Next() {
		var p PersonRow
		err := rows.Scan(&p.ID, &p.Name, &p.Title, &p.Company, &p.FediID, &p.Note, &p.Favorite)
		if nil != err {
			log.Error(field.S("failed to scan person row"), field.E(err))
			return nil, erorr.Wrap(err, "failed to scan person row")
		}
		people = append(people, p)
	}
	if err := rows.Err(); nil != err {
		return nil, erorr.Wrap(err, "failed to iterate people rows")
	}

	return people, nil
}

func LoadMessagesForPerson(logger log.Logger, db *sql.DB, personID int64) ([]MessageRow, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return nil, erorr.Wrap(erorr.Error("nil db"), "failed to load messages for person")
	}

	rows, err := db.Query(`SELECT id, person_id, group_id, from_me, sender, text, timestamp FROM messages WHERE person_id = ? ORDER BY id`, personID)
	if nil != err {
		log.Error(field.S("failed to load messages for person"), field.E(err), field.Int64("person_id", personID))
		return nil, erorr.Wrap(err, "failed to load messages for person")
	}
	defer rows.Close()

	var messages []MessageRow
	for rows.Next() {
		var m MessageRow
		err := rows.Scan(&m.ID, &m.PersonID, &m.GroupID, &m.FromMe, &m.Sender, &m.Text, &m.Timestamp)
		if nil != err {
			log.Error(field.S("failed to scan message row"), field.E(err))
			return nil, erorr.Wrap(err, "failed to scan message row")
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); nil != err {
		return nil, erorr.Wrap(err, "failed to iterate message rows")
	}

	return messages, nil
}

func InsertPerson(logger log.Logger, db *sql.DB, p PersonRow) (int64, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return 0, erorr.Wrap(erorr.Error("nil db"), "failed to insert person")
	}

	result, err := db.Exec(
		`INSERT INTO people (name, title, company, fedi_id, note, favorite) VALUES (?, ?, ?, ?, ?, ?)`,
		p.Name, p.Title, p.Company, p.FediID, p.Note, p.Favorite,
	)
	if nil != err {
		log.Error(field.S("failed to insert person"), field.E(err), field.String("name", p.Name))
		return 0, erorr.Wrap(err, "failed to insert person")
	}

	id, err := result.LastInsertId()
	if nil != err {
		return 0, erorr.Wrap(err, "failed to get last insert id for person")
	}

	return id, nil
}

func UpdatePersonFavorite(logger log.Logger, db *sql.DB, personID int64, favorite bool) error {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return erorr.Wrap(erorr.Error("nil db"), "failed to update person favorite")
	}

	_, err := db.Exec(`UPDATE people SET favorite = ? WHERE id = ?`, favorite, personID)
	if nil != err {
		log.Error(field.S("failed to update person favorite"), field.E(err), field.Int64("person_id", personID))
		return erorr.Wrap(err, "failed to update person favorite")
	}

	return nil
}

func InsertMessage(logger log.Logger, db *sql.DB, m MessageRow) (int64, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return 0, erorr.Wrap(erorr.Error("nil db"), "failed to insert message")
	}

	result, err := db.Exec(
		`INSERT INTO messages (person_id, group_id, from_me, sender, text, timestamp) VALUES (?, ?, ?, ?, ?, ?)`,
		m.PersonID, m.GroupID, m.FromMe, m.Sender, m.Text, m.Timestamp,
	)
	if nil != err {
		log.Error(field.S("failed to insert message"), field.E(err))
		return 0, erorr.Wrap(err, "failed to insert message")
	}

	id, err := result.LastInsertId()
	if nil != err {
		return 0, erorr.Wrap(err, "failed to get last insert id for message")
	}

	return id, nil
}
