SET statement_timeout = 0;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE ClanUpgrades AS (
    members INT,
    banks INT,
    missions INT
);
CREATE TYPE Vip as (
    vip TIMESTAMPTZ,
    background varchar(200)
);
CREATE TYPE Arena as(
    active BOOL,
    win INT,
    lose INT,
    lastFight BigInt
);

CREATE TYPE Daily AS (last TIMESTAMPTZ, strikes INT);

CREATE TYPE Cooldowns AS (trade TIMESTAMPTZ, daily Daily);
CREATE TYPE MissionProgress AS (
    type INT,
    level INT,
    progress INT,
    adv INT
);

CREATE TYPE Mission AS(
    trade TIMESTAMPTZ,
    last TIMESTAMPTZ,
    missions MissionProgress []
);

CREATE TYPE Dungeon AS (floor int, resets int);

CREATE TYPE Status AS (win int, lose int);

CREATE TABLE Clan(
    name VARCHAR(25) PRIMARY KEY,
    xp INT,
    createdAt TIMESTAMPTZ,
    background VARCHAR(300),
    money INT,
    upgrades ClanUpgrades,
    lastIncome TIMESTAMPTZ,
    mission TIMESTAMPTZ,
    missionProgress INT
);

CREATE TABLE Users(
    ID BigInt PRIMARY KEY,
    xp INT,
    upgrades INT [],
    status Status,
    money INT,
    clan VARCHAR(26) REFERENCES Clan(name),
    dungeon Dungeon,
    mission Mission,
    vip Vip,
    trainLimit INT,
    asuraCoin INT,
    arena Arena,
    rank INT,
    cooldowns Cooldowns,
    pity INT
);


CREATE TABLE Rooster(
    ID uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    userID BIGINT REFERENCES Users (ID) ON DELETE CASCADE NOT NULL,
    name VARCHAR(26),
    resets INT,
	equip BOOL,
    xp INT,
    type INT,
    equipped INT []
);



CREATE TABLE ClanMember(
    ID BigInt PRIMARY KEY,
    clan VARCHAR(26) REFERENCES Clan(name) ON DELETE CASCADE NOT NULL,
    role INT,
    xp INT
);

CREATE TABLE Item(
    ID uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    userID BIGINT REFERENCES Users (ID) ON DELETE CASCADE NOT NULL,
    quatity INT,
    itemID INT,
	equip BOOL,
    type INT
);

