package monkebase

import (
	"github.com/google/uuid"
	"github.com/imonke/monketype"

	"testing"
)

type banSet struct {
	Ban  monketype.Ban
	Want bool
}

func Test_WriteBan(test *testing.T) {
	var banner string = uuid.New().String()
	var banned string = uuid.New().String()
	var duration int64 = 60 * 60 * 24 * 7
	var ban monketype.Ban = monketype.NewBan(banner, banned, "", duration, false)

	var err error
	if err = WriteBan(ban.Map()); err != nil {
		test.Fatal(err)
	}
}

func Test_ReadSingleBan(test *testing.T) {
	var banned string = uuid.New().String()
	var ban monketype.Ban = monketype.NewBan("", banned, "", 0, false)
	WriteBan(ban.Map())

	var fetched monketype.Ban
	var exists bool
	var err error
	if fetched, exists, err = ReadSingleBan(ban.ID); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("ban %s does not exist", ban.ID)
	}

	if fetched.Banned != ban.Banned {
		test.Errorf("banned mismatch! have: %s, want: %s", fetched.Banned, ban.Banned)
	}

	if _, exists, err = ReadSingleBan(uuid.New().String()); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("random uuid got some ban")
	}
}

func Test_ReadSingleBan_nobody(test *testing.T) {
	var exists bool
	var err error
	if _, exists, err = ReadSingleBan(uuid.New().String()); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("random uuid got some ban")
	}
}

func Test_IsBanned_OldAndForever(test *testing.T) {
	var set banSet
	var sets []banSet = []banSet{
		banSet{
			Ban:  monketype.NewBan("", uuid.New().String(), "", -100, true),
			Want: true,
		},
		banSet{
			Ban:  monketype.NewBan("", uuid.New().String(), "", -100, false),
			Want: false,
		},
		banSet{
			Ban:  monketype.NewBan("", uuid.New().String(), "", 60*60*360, false),
			Want: true,
		},
		banSet{
			Ban:  monketype.NewBan("", uuid.New().String(), "", 60*60*360, true),
			Want: true,
		},
	}

	var err error
	var banned bool

	for _, set = range sets {
		if err = WriteBan(set.Ban.Map()); err != nil {
			test.Fatal(err)
		}

		if banned, err = IsBanned(set.Ban.Banned); err != nil {
			test.Fatal(err)
		}

		if banned != set.Want {
			test.Errorf("Ban state is %t!\n%#v", banned, set.Ban)
		}
	}
}

func Test_ReadBansOfUser(test *testing.T) {
	EmptyTable(BAN_TABLE)
	var banned string = uuid.New().String()
	var count, index int = 20, 0
	for count != index {
		WriteBan(monketype.NewBan(uuid.New().String(), uuid.New().String(), "", 0, true).Map())
		WriteBan(monketype.NewBan(uuid.New().String(), banned, "", 0, true).Map())
		WriteBan(monketype.NewBan(uuid.New().String(), uuid.New().String(), "", 0, true).Map())
		index++
	}

	count = 10
	var bans []monketype.Ban
	var size int
	var err error
	if bans, size, err = ReadBansOfUser(banned, "", count); err != nil {
		test.Fatal(err)
	}

	if count != size {
		test.Errorf("size expected mismatch! have: %d, want: %d", size, count)
	}

	if len(bans) != size {
		test.Errorf("size actual mismatch! have: %d, want: %d", size, len(bans))
	}

	var ban monketype.Ban
	for _, ban = range bans {
		if ban.Banned != banned {
			test.Errorf("Ban banned mismatch! have: %s, want: %s", ban.Banned, banned)
		}
	}
}

func Test_ReadBansOfUser_after(test *testing.T) {
	EmptyTable(BAN_TABLE)
	var banned string = uuid.New().String()
	var count, index int = 20, 0
	for count != index {
		WriteBan(monketype.NewBan(uuid.New().String(), uuid.New().String(), "", 0, true).Map())
		WriteBan(monketype.NewBan(uuid.New().String(), banned, "", 0, true).Map())
		WriteBan(monketype.NewBan(uuid.New().String(), uuid.New().String(), "", 0, true).Map())
		index++
	}

	count = 10
	var offset int = 5
	var first, second []monketype.Ban
	var err error
	if first, _, err = ReadBansOfUser(banned, "", count); err != nil {
		test.Fatal(err)
	}

	if second, _, err = ReadBansOfUser(banned, first[offset].ID, count); err != nil {
		test.Fatal(err)
	}

	var single monketype.Ban
	for index, single = range first[offset+1:] {
		if single.ID != second[index].ID {
		}
	}
}

func Test_WriteReport(test *testing.T) {
	var reporter string = uuid.New().String()
	var reported string = uuid.New().String()
	var report monketype.Report = monketype.NewReport(reporter, reported, "", "happens to be that smelly && stinky == True")

	var err error
	if err = WriteReport(report.Map()); err != nil {
		test.Fatal(err)
	}

	var exists bool
	if _, exists, err = ReadSingleReport(report.ID); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("report %s does not exist", report.ID)
	}
}

func Test_ReadReport(test *testing.T) {
	var reporter string = uuid.New().String()
	var reported string = uuid.New().String()
	var report monketype.Report = monketype.NewReport(reporter, reported, "", "Called me the J word (javascript developer)")
	WriteReport(report.Map())

	var fetched monketype.Report
	var exists bool
	var err error
	if fetched, exists, err = ReadSingleReport(report.ID); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("report %s does not exist", report.ID)
	}

	var mapped map[string]interface{} = fetched.Map()

	var key string
	var value interface{}
	for key, value = range report.Map() {
		if mapped[key] != value {
			test.Errorf("mismatch at %s! have: %#v, want: %#v", key, value, mapped[key])
		}
	}
}

func Test_ReadReport_notExists(test *testing.T) {
	var id string = uuid.New().String()

	var fetched monketype.Report
	var exists bool
	var err error
	if _, exists, err = ReadSingleReport(id); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("id %s references a report %#v", id, fetched)
	}
}

func Test_ReadManyUnresolvedReport(test *testing.T) {
	EmptyTable(REPORT_TABLE)

	var report monketype.Report
	var index, limit int = 0, 20
	for index != limit {
		report = monketype.NewReport(uuid.New().String(), uuid.New().String(), "user", "")
		report.Resolved = true
		WriteReport(report.Map())

		WriteReport(monketype.NewReport(uuid.New().String(), uuid.New().String(), "user", "").Map())

		report = monketype.NewReport(uuid.New().String(), uuid.New().String(), "user", "")
		report.Resolved = true
		WriteReport(report.Map())
		index++
	}

	var count, offset int = 10, 5
	var first, second []monketype.Report
	var err error
	if first, _, err = ReadManyUnresolvedReport("", count); err != nil {
		test.Fatal(err)
	}

	if second, _, err = ReadManyUnresolvedReport(first[offset].ID, count); err != nil {
		test.Fatal(err)
	}

	var single monketype.Report
	for index, single = range first[offset+1:] {
		if single.ID != second[index].ID {
			test.Errorf("IDs not aligned! have: %s, want: %s", second[index].ID, single.ID)
		}
	}
}
