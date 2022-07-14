package connections_test

import (
	"log"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUpdateFileSQL(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("../../gallery.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = db.Debug()
	db.DryRun = true
	db.Session(&gorm.Session{
		SkipDefaultTransaction: true,
	})
	db.Raw(`BEGIN TRANSACTION;`)
	db.Raw(`
			INSERT INTO
			file_entries 
			(
				created_at,
				updated_at,
				name,
				relative_path,
				file_size,
				thumbnail,
				machine_id,
				is_invalid,
				parent_directory_id,
				thumbnail_height,
				thumbnail_width
			)
			VALUES
			(
				CURRENT_TIMESTAMP,
				CURRENT_TIMESTAMP,
				?,
				?,
				?,
				?,
				?,
				0,
				(SELECT id FROM directories WHERE relative_path=? and machine_id=? LIMIT 1),
				?,
				?
			)
			ON CONFLICT(relative_path, machine_id) DO
			UPDATE
			SET updated_at = EXCLUDED.updated_at,
				name = EXCLUDED.name,
				file_size = EXCLUDED.file_size,
				thumbnail = EXCLUDED.thumbnail,
				is_invalid = 0,
				parent_directory_id = EXCLUDED.parent_directory_id,
				thumbnail_height = EXCLUDED.thumbnail_height,
				thumbnail_width = EXCLUDED.thumbnail_width
			;
			`, "a.jpg", "IMG_8654.jpg", 1234, "data:image/png;fijewaj",
		1, "ps吧", 1, 1024, 1024)

	db.Raw(`
			INSERT INTO
			file_entries 
			(
				created_at,
				updated_at,
				name,
				relative_path,
				file_size,
				thumbnail,
				machine_id,
				is_invalid,
				parent_directory_id,
				thumbnail_height,
				thumbnail_width
			)
			VALUES
			(
				CURRENT_TIMESTAMP,
				CURRENT_TIMESTAMP,
				?,
				?,
				?,
				?,
				?,
				0,
				(SELECT id FROM directories WHERE relative_path=? and machine_id=? LIMIT 1),
				?,
				?
			)
			ON CONFLICT(relative_path, machine_id) DO
			UPDATE
			SET updated_at = EXCLUDED.updated_at,
				name = EXCLUDED.name,
				file_size = EXCLUDED.file_size,
				thumbnail = EXCLUDED.thumbnail,
				is_invalid = 0,
				parent_directory_id = EXCLUDED.parent_directory_id,
				thumbnail_height = EXCLUDED.thumbnail_height,
				thumbnail_width = EXCLUDED.thumbnail_width
			;
			`, "a.jpg", "IMG_8654.jpg", 1234, "data:image/png;fijewaj",
		1, "ps吧", 1, 1024, 1024)

	if db.RowsAffected != 0 {
		t.Logf("rows Affected is %d\n", db.RowsAffected)
		t.Fail()
	}
	db.Raw(`COMMIT;`)
	log.Fatalf("%s\n", db.Statement.SQL.String())
	var exists bool
	db.Row().Scan(&exists)
	if !exists {
		t.Fail()
	}
	if db.Error != nil {
		t.Log(db.Error)
		t.Fail()
	}
	if db.RowsAffected != 1 {
		t.Logf("rows Affected is %d\n", db.RowsAffected)
		t.Fail()
	}
}
