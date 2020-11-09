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

	var size_bans, index int = 20, 0
	var bans []monketype.Ban = make([]monketype.Ban, size_bans)
	for ; index != size_bans; index++ {
		bans[index] = monketype.NewBan(uuid.New().String(), banned, "", 0, true)
		bans[index].Created = bans[index].Created + int64(100+size_bans-index)
		WriteBan(bans[index].Map())
		WriteBan(monketype.NewBan(uuid.New().String(), uuid.New().String(), "", 0, true).Map())
	}

	var fetched []monketype.Ban
	var offset, count, size int = 7, 9, 0
	var err error
	if fetched, size, err = ReadBansOfUser(banned, offset, count); err != nil {
		test.Fatal(err)
	}

	if size != count {
		test.Errorf("read %d bans, expected %d\n%#v", size, count, fetched)
	}

	if len(fetched) != size {
		test.Errorf("actual size %d does not match size %d", len(fetched), size)
	}

	var now, last int64 = 0, fetched[0].Created
	var ban monketype.Ban
	for index, ban = range fetched[1:] {
		if ban.ID != bans[1+index+offset].ID {
			test.Errorf("ban mismatch: \nhave: %s, \nwant: %s", ban.ID, bans[1+index+offset].ID)
		}

		now = ban.Created
		if now > last {
			test.Errorf("Fail at %d: %d < %d", index, now, last)
		}

		last = now
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
	var index int
	var size_resolved, size_unresolved int = 10, 20
	for index = 0; index != size_resolved; index++ {
		report = monketype.NewReport(uuid.New().String(), uuid.New().String(), "user", "")
		report.Resolved = true
		WriteReport(report.Map())
	}

	var unresolved []monketype.Report = make([]monketype.Report, size_unresolved)
	for index = 0; index != size_unresolved; index++ {
		unresolved[index] = monketype.NewReport(uuid.New().String(), uuid.New().String(), "user", "")
		unresolved[index].Resolved = false
		unresolved[index].Created = unresolved[index].Created + int64(100+size_unresolved-index)
		WriteReport(unresolved[index].Map())
	}

	var fetched []monketype.Report
	var offset, count, size int = 3, 7, 0
	var err error
	if fetched, size, err = ReadManyUnresolvedReport(offset, count); err != nil {
		test.Fatal(err)
	}

	if size != count {
		test.Errorf("read %d reports, expected %d\n%#v", size, count, fetched)
	}

	if len(fetched) != size {
		test.Errorf("actual size %d does not match size %d", len(fetched), size)
	}

	var now, last int64 = 0, fetched[0].Created
	for index, report = range fetched[1:] {
		if report.ID != unresolved[1+index+offset].ID {
			test.Errorf("report mismatch: \nhave: %#v, \nwant: %#v", report, unresolved[1+index+offset])
		}

		now = report.Created
		if now > last {
			test.Errorf("Fail at %d: %d < %d", index, now, last)
		}

		last = now
	}
}
