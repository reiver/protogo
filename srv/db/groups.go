package dbsrv

import (
	"database/sql"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
	"codeberg.org/reiver/go-log"
)

type GroupRow struct {
	ID       int64
	Name     string
	Favorite bool
	Members  []string
}

func LoadGroups(logger log.Logger, db *sql.DB) ([]GroupRow, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return nil, erorr.Wrap(erorr.Error("nil db"), "failed to load groups")
	}

	rows, err := db.Query(`SELECT id, name, favorite FROM groups ORDER BY id`)
	if nil != err {
		log.Error(field.S("failed to load groups"), field.E(err))
		return nil, erorr.Wrap(err, "failed to load groups")
	}
	defer rows.Close()

	var groups []GroupRow
	for rows.Next() {
		var g GroupRow
		err := rows.Scan(&g.ID, &g.Name, &g.Favorite)
		if nil != err {
			log.Error(field.S("failed to scan group row"), field.E(err))
			return nil, erorr.Wrap(err, "failed to scan group row")
		}
		groups = append(groups, g)
	}
	if err := rows.Err(); nil != err {
		return nil, erorr.Wrap(err, "failed to iterate group rows")
	}

	for i := range groups {
		members, err := loadGroupMembers(logger, db, groups[i].ID)
		if nil != err {
			return nil, err
		}
		groups[i].Members = members
	}

	return groups, nil
}

func loadGroupMembers(logger log.Logger, db *sql.DB, groupID int64) ([]string, error) {
	log := logger.Begin()
	defer log.End()

	rows, err := db.Query(`SELECT person_name FROM group_members WHERE group_id = ? ORDER BY id`, groupID)
	if nil != err {
		log.Error(field.S("failed to load group members"), field.E(err), field.Int64("group_id", groupID))
		return nil, erorr.Wrap(err, "failed to load group members")
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if nil != err {
			return nil, erorr.Wrap(err, "failed to scan group member row")
		}
		members = append(members, name)
	}
	if err := rows.Err(); nil != err {
		return nil, erorr.Wrap(err, "failed to iterate group member rows")
	}

	return members, nil
}

func LoadMessagesForGroup(logger log.Logger, db *sql.DB, groupID int64) ([]MessageRow, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return nil, erorr.Wrap(erorr.Error("nil db"), "failed to load messages for group")
	}

	rows, err := db.Query(`SELECT id, person_id, group_id, from_me, sender, text, timestamp FROM messages WHERE group_id = ? ORDER BY id`, groupID)
	if nil != err {
		log.Error(field.S("failed to load messages for group"), field.E(err), field.Int64("group_id", groupID))
		return nil, erorr.Wrap(err, "failed to load messages for group")
	}
	defer rows.Close()

	var messages []MessageRow
	for rows.Next() {
		var m MessageRow
		err := rows.Scan(&m.ID, &m.PersonID, &m.GroupID, &m.FromMe, &m.Sender, &m.Text, &m.Timestamp)
		if nil != err {
			return nil, erorr.Wrap(err, "failed to scan message row")
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); nil != err {
		return nil, erorr.Wrap(err, "failed to iterate message rows")
	}

	return messages, nil
}

func InsertGroup(logger log.Logger, db *sql.DB, g GroupRow) (int64, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return 0, erorr.Wrap(erorr.Error("nil db"), "failed to insert group")
	}

	tx, err := db.Begin()
	if nil != err {
		log.Error(field.S("failed to start transaction"), field.E(err))
		return 0, erorr.Wrap(err, "failed to start transaction for group insert")
	}
	defer tx.Rollback()

	result, err := tx.Exec(`INSERT INTO groups (name, favorite) VALUES (?, ?)`, g.Name, g.Favorite)
	if nil != err {
		log.Error(field.S("failed to insert group"), field.E(err), field.String("name", g.Name))
		return 0, erorr.Wrap(err, "failed to insert group")
	}

	groupID, err := result.LastInsertId()
	if nil != err {
		return 0, erorr.Wrap(err, "failed to get last insert id for group")
	}

	for _, member := range g.Members {
		_, err := tx.Exec(`INSERT INTO group_members (group_id, person_name) VALUES (?, ?)`, groupID, member)
		if nil != err {
			log.Error(field.S("failed to insert group member"), field.E(err), field.String("member", member))
			return 0, erorr.Wrap(err, "failed to insert group member")
		}
	}

	if err := tx.Commit(); nil != err {
		log.Error(field.S("failed to commit transaction"), field.E(err))
		return 0, erorr.Wrap(err, "failed to commit group insert transaction")
	}

	return groupID, nil
}

func UpdateGroupFavorite(logger log.Logger, db *sql.DB, groupID int64, favorite bool) error {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return erorr.Wrap(erorr.Error("nil db"), "failed to update group favorite")
	}

	_, err := db.Exec(`UPDATE groups SET favorite = ? WHERE id = ?`, favorite, groupID)
	if nil != err {
		log.Error(field.S("failed to update group favorite"), field.E(err), field.Int64("group_id", groupID))
		return erorr.Wrap(err, "failed to update group favorite")
	}

	return nil
}
