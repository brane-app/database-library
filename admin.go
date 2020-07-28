package monkebase

import (
	"github.com/imonke/monketype"
	"github.com/jmoiron/sqlx"

	"database/sql"
	"time"
)

func WriteBan(ban map[string]interface{}) (err error) {
	var statement string
	var values []interface{}
	statement, values = makeSQLInsertable(BAN_TABLE, ban)

	_, err = database.Exec(statement, values...)
	return
}

func readSingleBanKey(key, query string) (ban monketype.Ban, exists bool, err error) {
	var statement string = "SELECT * FROM " + BAN_TABLE + " WHERE " + key + "=? LIMIT 1"
	if err = database.QueryRowx(statement, query).StructScan(&ban); err != nil {
		exists = false

		if err == sql.ErrNoRows {
			err = nil
		}

		return
	}

	exists = true
	return
}

/**
 * Read a single ban of id `ID`
 * Done in one query
 */
func ReadSingleBan(ID string) (ban monketype.Ban, exists bool, err error) {
	ban, exists, err = readSingleBanKey("id", ID)
	return
}

/**
 * Read a slice of bans of a user
 * Done in one query
 */
func ReadBansOfUser(ID string, offset, count int) (bans []monketype.Ban, size int, err error) {
	var statement string = "SELECT * FROM " + BAN_TABLE + " WHERE banned=? ORDER BY created DESC LIMIT ?, ?"
	var rows *sqlx.Rows
	if rows, err = database.Queryx(statement, ID, size, offset); err != nil || rows == nil {
		return
	}

	defer rows.Close()

	bans = make([]monketype.Ban, count)
	size = 0
	for rows.Next() {
		rows.StructScan(&bans[size])
	}

	return
}

/**
 * Get whether or not a user is banned, either by a permanent ban, or an expirable ban
 * Done in one query
 */
func IsBanned(ID string) (banned bool, err error) {
	var count int

	var now int64 = time.Now().Unix()
	var statement string = "SELECT COUNT(*) FROM " + BAN_TABLE + " WHERE (banned=? AND forever) OR (banned=? AND expires>?)"
	if err = database.QueryRowx(statement, ID, ID, now).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		banned = false
		return
	}

	banned = count != 0
	return
}

func WriteReport(report map[string]interface{}) (err error) {
	var statement string
	var values []interface{}
	statement, values = makeSQLInsertable(REPORT_TABLE, report)

	_, err = database.Exec(statement, values...)
	return
}

func ReadManyUnresolvedReport(offset, count int) (reports []monketype.Report, size int, err error) {
	var statement string = "SELECT * FROM " + REPORT_TABLE + " WHERE resolved!=0 ORDER BY created DESC LIMIT ?, ?"
	var rows *sqlx.Rows
	if rows, err = database.Queryx(statement, offset, count); err != nil || rows == nil {
		return
	}

	defer rows.Close()
	reports = make([]monketype.Report, count)
	size = 0

	for rows.Next() {
		rows.StructScan(&reports[size])
		size++
	}

	reports = reports[0:size]
	return
}

func ReadSingleReport(ID string) (report monketype.Report, exists bool, err error) {
	var statement string = "SELECT * FROM " + REPORT_TABLE + " WHERE id=?"
	if err = database.QueryRowx(statement, ID).StructScan(&report); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}

	exists = true
	return
}
