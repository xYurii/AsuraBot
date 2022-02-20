package entities

import (
	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ItemType int

const (
	_ ItemType = iota
	LootboxType
	NormalType
	CosmeticType
)

type MissionType int

const (
	Win MissionType = iota
	Fight
	WinGalo
	FightGalo
)

type Mission struct {
	bun.BaseModel `bun:"table:mission,alias:mission"`

	ID       uuid.UUID         `bun:"id,pk"`
	UserID   disgord.Snowflake `bun:"userid"`
	Type     MissionType       `bun:"type"`
	Level    int               `bun:"level"`
	Progress int               `bun:"progress"`
	Adv      int               `bun:"adv"`
}

type Item struct {
	bun.BaseModel `bun:"table:item,alias:item"`

	ID       uuid.UUID         `bun:"id,pk"`
	UserID   disgord.Snowflake `bun:"userid"`
	Quantity int               `bun:"quantity"`
	ItemID   int               `bun:"itemid"`
	Equip    bool              `bun:"equip"`
	Type     ItemType          `bun:"type"`
}
type Rooster struct {
	bun.BaseModel `bun:"table:rooster,alias:galo"`

	ID       uuid.UUID         `bun:"id,pk"`
	UserID   disgord.Snowflake `bun:"userid"`
	Name     string            `bun:"name"`
	Resets   int               `bun:"resets"`
	Equip    bool              `bun:"equip"`
	Xp       int               `bun:"xp"`
	Type     int               `bun:"type"`
	Equipped []int             `bun:"equipped,array"`
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID              disgord.Snowflake `bun:"id,pk"`
	UserXp          int               `bun:"xp"`
	Galos           []*Rooster        `bun:"rel:has-many,join:id=userid"`
	Items           []*Item           `bun:"rel:has-many,join:id=userid"`
	Upgrades        []int             `bun:"upgrades,array"`
	Win             int               `bun:"win"`
	Lose            int               `bun:"lose"`
	Money           int               `bun:"money"`
	Dungeon         int               `bun:"dungeon"`
	DungeonReset    int               `bun:"dungeonreset"`
	TradeMission    uint64            `bun:"trademission"`
	LastMission     uint64            `bun:"lastmission"`
	Missions        []*Mission        `bun:"rel:has-many,join:id=userid"`
	Vip             uint64            `bun:"vip"`
	VipBackground   string            `bun:"vipbackground"`
	TrainLimit      int               `bun:"trainlimit"`
	TrainLimitReset uint64            `bun:"trainlimitreset"`
	AsuraCoin       int               `bun:"asuracoin"`
	ArenaActive     bool              `bun:"arenaactive"`
	ArenaWin        int               `bun:"arenawin"`
	ArenaLose       int               `bun:"arenalose"`
	ArenaLastFight  disgord.Snowflake `bun:"arenalastfight"`
	Rank            int               `bun:"rank"`
	TradeItem       uint64            `bun:"tradeitem"`
	Daily           uint64            `bun:"daily"`
	DailyStrikes    int               `bun:"dailystrikes"`
	Pity            int               `bun:"pity"`
}
