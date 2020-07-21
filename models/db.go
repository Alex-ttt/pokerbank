package models

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

type DatabaseSettings struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func InitDb(connectionSettings DatabaseSettings) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		connectionSettings.Host,
		connectionSettings.Port,
		connectionSettings.User,
		connectionSettings.Password,
		connectionSettings.DbName,
	)
	var err error
	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic(err)
	}

	CreateDatabaseStructure(Db)

	if err = Db.Ping(); err != nil {
		log.Panic(err)
	}
}
func CreateDatabaseStructure(db *sql.DB) {
	schemaRow := db.QueryRow(
		`select exists(
				SELECT schema_name
				FROM information_schema.schemata
				WHERE schema_name = 'poker')`)

	var isSchemaExists bool
	err := schemaRow.Scan(&isSchemaExists)
	if err != nil {
		panic(err)
	}

	if isSchemaExists {
		return
	}

	//pathWd, _ := os.Getwd()
	//absolutePath := pathWd + "\\models\\sql\\DbInit.sql"
	//query, _ := ioutil.ReadFile(absolutePath)

	_, err = db.Exec("create schema poker;\n\ncreate table poker.\"Games\" (\n   \"Id\" serial not null constraint games_pk primary key,\n   \"Name\" varchar(100) not null,\n   \"Date\" date\n);\n\nalter table poker.\"Games\" owner to postgres;\n\nINSERT INTO poker.\"Games\" (\"Name\", \"Date\") VALUES ('Банк до начала ведения учёта', null);\nINSERT INTO poker.\"Games\" (\"Name\", \"Date\") VALUES ('Игра у Жени', '2020-07-03');\nINSERT INTO poker.\"Games\" (\"Name\", \"Date\") VALUES ('Игра у Егора', '2020-07-06');\nINSERT INTO poker.\"Games\" (\"Name\", \"Date\") VALUES ('Игра у Жени', '2020-07-10');\nINSERT INTO poker.\"Games\" (\"Name\", \"Date\") VALUES ('Игра у Егора', '2020-07-12');\n\ncreate table poker.\"Players\" (\n    \"Id\" serial not null constraint players_pk primary key,\n    \"Name\" varchar(20) not null,\n    \"Surname\" varchar(30)\n);\n\nalter table poker.\"Players\" owner to postgres;\n\nINSERT INTO poker.\"Players\" (\"Name\", \"Surname\") VALUES ('Егор', 'Смеловский');\nINSERT INTO poker.\"Players\" (\"Name\", \"Surname\") VALUES ('Евгений', 'Перельман');\nINSERT INTO poker.\"Players\" (\"Name\", \"Surname\") VALUES ('Матвей', 'Минеев');\nINSERT INTO poker.\"Players\" (\"Name\", \"Surname\") VALUES ('Александр', 'Тарасенко');\nINSERT INTO poker.\"Players\" (\"Name\", \"Surname\") VALUES ('Роман', 'Огурешнов');\nINSERT INTO poker.\"Players\" (\"Name\", \"Surname\") VALUES ('Алексей', 'Степанов');\n\n\ncreate table poker.\"Debts\" (\n   \"Id\" serial not null constraint debts_pk primary key,\n   \"WinnerId\" integer not null constraint debts_players_id_fk references poker.\"Players\",\n   \"LoserId\" integer not null constraint debts_players_id_fk_2 references poker.\"Players\",\n   \"Amount\" integer not null, \"GameId\" integer constraint debts_games_id_fk references poker.\"Games\"\n);\n\nalter table poker.\"Debts\" owner to postgres;\n\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (4, 1, 6000, 1);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (4, 3, 1200, 1);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (2, 5, 8350, 1);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (2, 1, 3250, 1);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (3, 1, 3000, 1);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (5, 6, 2000, 2);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (5, 1, 2600, 2);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (4, 1, 3400, 2);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (5, 2, 1000, 2);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (2, 1, 3750, 3);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (2, 3, 1250, 3);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (4, 3, 4500, 4);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (5, 1, 1150, 4);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (2, 6, 550, 4);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (5, 6, 450, 4);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (3, 5, 4250, 5);\nINSERT INTO poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\") VALUES (2, 5, 750, 5);\n\ncreate table poker.\"DebtPayments\" (\n    \"Id\" serial not null constraint debtpayments_pk primary key,\n    \"PayerId\" integer not null constraint debtpayments_players_id_fk references poker.\"Players\",\n    \"RecipientId\" integer not null constraint debtpayments_players_id_fk_2 references poker.\"Players\",\n    \"Amount\" integer not null,\n    \"InsertStamp\" timestamp default now()\n);\n\nalter table poker.\"DebtPayments\" owner to postgres;\n\nINSERT INTO poker.\"DebtPayments\" (\"PayerId\", \"RecipientId\", \"InsertStamp\", \"Amount\") VALUES (6, 5, now(), 2000);\nINSERT INTO poker.\"DebtPayments\" (\"PayerId\", \"RecipientId\", \"InsertStamp\", \"Amount\") VALUES (3, 2, now(), 100);\nINSERT INTO poker.\"DebtPayments\" (\"PayerId\", \"RecipientId\", \"InsertStamp\", \"Amount\") VALUES (5, 2, now(), 5350);\nINSERT INTO poker.\"DebtPayments\" (\"PayerId\", \"RecipientId\", \"InsertStamp\", \"Amount\") VALUES (3, 2, now(), 200);\n\ncreate or replace function poker.playersdebts()\n    returns TABLE(winnerId int, winnername character varying, losername character varying, playerwin integer, commonplayerwin integer)\n    language plpgsql\nas\n$$\nBEGIN\n    CREATE temporary table TempDebts AS\n    select\n        win.\"WinnerId\" as WinnerId,\n        win.\"LoserId\" as LoserId,\n        coalesce(win.win, 0) - coalesce(lose.lose, 0) as DebtAmount\n    from poker.\"Players\" as p\n             join (\n        select d.\"WinnerId\", d.\"LoserId\", sum(\"Amount\") as win\n        from poker.\"Debts\" as d\n        group by d.\"WinnerId\", d.\"LoserId\"\n    ) as win on win.\"WinnerId\" = p.\"Id\"\n             left  join (\n        select d.\"WinnerId\", d.\"LoserId\", sum(\"Amount\") as lose\n        from poker.\"Debts\" AS d\n        group by d.\"WinnerId\", d.\"LoserId\"\n    ) as lose on lose.\"WinnerId\" = win.\"LoserId\" and lose.\"LoserId\" = p.\"Id\"\n    where coalesce(win.win, 0) - coalesce(lose.lose, 0) >= 0;\n\n    update TempDebts AS temp_debts\n    set DebtAmount = temp_debts.DebtAmount - all_payments.Amount\n    from (\n             select\n                 dp.\"PayerId\",\n                 dp.\"RecipientId\",\n                 sum(dp.\"Amount\") as Amount\n             from poker.\"DebtPayments\" as dp\n             group by  dp.\"PayerId\", dp.\"RecipientId\"\n         ) as all_payments\n    where\n            all_payments.\"PayerId\" = temp_debts.LoserId\n      and all_payments.\"RecipientId\" = temp_debts.WinnerId;\n\n    return query\n        select\n            w.\"Id\",\n            w.\"Name\" AS WinnerName,\n            l.\"Name\" AS LoserName,\n            cast(source.DebtAmount AS int) AS PlayerWin,\n            cast(\n                            sum(source.DebtAmount) OVER (PARTITION BY w.\"Name\")\n                AS int\n                )\n                AS CommonPlayerWin\n        from\n            (\n                select\n                    debts.WinnerId,\n                    debts.LoserId,\n                    debts.DebtAmount\n                from TempDebts AS debts\n                union\n                select\n                    debts.WinnerId,\n                    losers.\"Id\",\n                    0\n                from TempDebts as debts\n                         join poker.\"Players\" as losers on losers.\"Id\" not in (\n                    select td.LoserId from TempDebts as td\n                    where td.WinnerId = debts.WinnerId\n                ) and exists(\n                   select * from TempDebts AS td_ where td_.LoserId = losers.\"Id\"\n               )\n            ) as source\n            join poker.\"Players\" as w on w.\"Id\" = source.WinnerId\n            join poker.\"Players\" as l on l.\"Id\" = source.LoserId\n        order by w.\"Name\", l.\"Name\";\n\n    drop table TempDebts;\n    return;\n\nEND\n$$;\n\nalter function poker.playersdebts() owner to postgres;\n\ncreate or replace function poker.gamesinfo()\n    returns TABLE\n            (\n                gameid int,\n                gamename character varying,\n                playername character varying,\n                playeramount integer,\n                gamedate date,\n                commongamebank int\n            )\n    language plpgsql\nas\n$$\nBEGIN\n    return query\n        select\n            games.\"Id\",\n            games.\"Name\",\n            players.\"Name\",\n            cast(games_info.Amount as int) as Amount,\n            games.\"Date\",\n            coalesce(game_bank.Bank, 0)\n        from (\n                 select debts.\"GameId\"      as GameId,\n                        debts.\"WinnerId\"    as PlayerId,\n                        sum(debts.\"Amount\") as Amount\n                 from poker.\"Debts\" as debts\n                 group by debts.\"GameId\", debts.\"WinnerId\"\n                 union\n                 select debts.\"GameId\"      as GameId,\n                        debts.\"LoserId\"     as PlayerId,\n                        -sum(debts.\"Amount\") as Amount\n                 from poker.\"Debts\" as debts\n                 group by debts.\"GameId\", debts.\"LoserId\"\n             ) as games_info\n                 join poker.\"Games\" as games ON games.\"Id\" = games_info.GameId\n                 left join (\n            select\n                g.\"GameId\" as GameId,\n                cast(sum(g.\"Amount\") as int) as Bank\n            from poker.\"Debts\" as g\n            group by g.\"GameId\"\n        ) as game_bank ON game_bank.GameId = games.\"Id\"\n                 join poker.\"Players\" as players ON players.\"Id\" = games_info.PlayerId\n        order by\n            games.\"Id\" desc,\n            games.\"Date\" desc nulls last,\n            games_info.Amount desc;\nEND\n$$;\n\nalter function poker.gamesinfo() owner to postgres;\n\ncreate type poker.PlayerGameResult AS\n(\n    WinnerId int,\n    LoserId int,\n    Amount int\n);\n\ncreate function poker.InsertGameResult\n(\n    gameName varchar(40),\n    gameDate date,\n    results poker.PlayerGameResult[]\n)\n    returns void\n    language plpgsql\nAS $$\ndeclare\n    newGameId int;\nbegin\n    insert into poker.\"Games\"(\"Name\", \"Date\")\n    values (gameName, gameDate)\n    returning \"Id\" into newGameId;\n\n    insert into poker.\"Debts\" (\"WinnerId\", \"LoserId\", \"Amount\", \"GameId\")\n    select\n        (game_results.res).WinnerId,\n        (game_results.res).LoserId,\n        (game_results.res).Amount,\n        newGameId\n    from (\n             select unnest(results) as res\n         ) as game_results;\nend;\n$$;\n\nalter function poker.InsertGameResult(gameName varchar, gameDate date, results poker.PlayerGameResult[]) owner to postgres;")
	if err != nil {
		panic(err)
	}
}
