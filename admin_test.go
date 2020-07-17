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
