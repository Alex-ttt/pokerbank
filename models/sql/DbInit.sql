create schema poker;

create table poker."Games" (
   "Id" serial not null constraint games_pk primary key,
   "Name" varchar(100) not null,
   "Date" date
);

alter table poker."Games" owner to postgres;

INSERT INTO poker."Games" ("Name", "Date") VALUES ('Банк до начала ведения учёта', null);
INSERT INTO poker."Games" ("Name", "Date") VALUES ('Игра у Жени', '2020-07-03');
INSERT INTO poker."Games" ("Name", "Date") VALUES ('Игра у Егора', '2020-07-06');
INSERT INTO poker."Games" ("Name", "Date") VALUES ('Игра у Жени', '2020-07-10');
INSERT INTO poker."Games" ("Name", "Date") VALUES ('Игра у Егора', '2020-07-12');

create table poker."Players" (
    "Id" serial not null constraint players_pk primary key,
    "Name" varchar(20) not null,
    "Surname" varchar(30)
);

alter table poker."Players" owner to postgres;

INSERT INTO poker."Players" ("Name", "Surname") VALUES ('Егор', 'Смеловский');
INSERT INTO poker."Players" ("Name", "Surname") VALUES ('Евгений', 'Перельман');
INSERT INTO poker."Players" ("Name", "Surname") VALUES ('Матвей', 'Минеев');
INSERT INTO poker."Players" ("Name", "Surname") VALUES ('Александр', 'Тарасенко');
INSERT INTO poker."Players" ("Name", "Surname") VALUES ('Роман', 'Огурешнов');
INSERT INTO poker."Players" ("Name", "Surname") VALUES ('Алексей', 'Степанов');


create table poker."Debts" (
   "Id" serial not null constraint debts_pk primary key,
   "WinnerId" integer not null constraint debts_players_id_fk references poker."Players",
   "LoserId" integer not null constraint debts_players_id_fk_2 references poker."Players",
   "Amount" integer not null, "GameId" integer constraint debts_games_id_fk references poker."Games"
);

alter table poker."Debts" owner to postgres;

INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (4, 1, 6000, 1);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (4, 3, 1200, 1);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (2, 5, 8350, 1);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (2, 1, 3250, 1);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (3, 1, 3000, 1);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (5, 6, 2000, 2);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (5, 1, 2600, 2);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (4, 1, 3400, 2);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (5, 2, 1000, 2);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (2, 1, 3750, 3);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (2, 3, 1250, 3);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (4, 3, 4500, 4);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (5, 1, 1150, 4);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (2, 6, 550, 4);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (5, 6, 450, 4);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (3, 5, 4250, 5);
INSERT INTO poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId") VALUES (2, 5, 750, 5);

create table poker."DebtPayments" (
    "Id" serial not null constraint debtpayments_pk primary key,
    "PayerId" integer not null constraint debtpayments_players_id_fk references poker."Players",
    "RecipientId" integer not null constraint debtpayments_players_id_fk_2 references poker."Players",
    "Amount" integer not null,
    "InsertStamp" timestamp default now()
);

alter table poker."DebtPayments" owner to postgres;

INSERT INTO poker."DebtPayments" ("PayerId", "RecipientId", "InsertStamp", "Amount") VALUES (6, 5, now(), 2000);
INSERT INTO poker."DebtPayments" ("PayerId", "RecipientId", "InsertStamp", "Amount") VALUES (3, 2, now(), 100);
INSERT INTO poker."DebtPayments" ("PayerId", "RecipientId", "InsertStamp", "Amount") VALUES (5, 2, now(), 5350);
INSERT INTO poker."DebtPayments" ("PayerId", "RecipientId", "InsertStamp", "Amount") VALUES (3, 2, now(), 200);

create or replace function poker.playersdebts()
    returns TABLE(winnerId int, winnername character varying, losername character varying, playerwin integer, commonplayerwin integer)
    language plpgsql
as
$$
BEGIN
    CREATE temporary table TempDebts AS
    select
        win."WinnerId" as WinnerId,
        win."LoserId" as LoserId,
        coalesce(win.win, 0) - coalesce(lose.lose, 0) as DebtAmount
    from poker."Players" as p
             join (
        select d."WinnerId", d."LoserId", sum("Amount") as win
        from poker."Debts" as d
        group by d."WinnerId", d."LoserId"
    ) as win on win."WinnerId" = p."Id"
             left  join (
        select d."WinnerId", d."LoserId", sum("Amount") as lose
        from poker."Debts" AS d
        group by d."WinnerId", d."LoserId"
    ) as lose on lose."WinnerId" = win."LoserId" and lose."LoserId" = p."Id"
    where coalesce(win.win, 0) - coalesce(lose.lose, 0) >= 0;

    update TempDebts AS temp_debts
    set DebtAmount = temp_debts.DebtAmount - all_payments.Amount
    from (
             select
                 dp."PayerId",
                 dp."RecipientId",
                 sum(dp."Amount") as Amount
             from poker."DebtPayments" as dp
             group by  dp."PayerId", dp."RecipientId"
         ) as all_payments
    where
            all_payments."PayerId" = temp_debts.LoserId
      and all_payments."RecipientId" = temp_debts.WinnerId;

    return query
        select
            w."Id",
            w."Name" AS WinnerName,
            l."Name" AS LoserName,
            cast(source.DebtAmount AS int) AS PlayerWin,
            cast(
                            sum(source.DebtAmount) OVER (PARTITION BY w."Name")
                AS int
                )
                AS CommonPlayerWin
        from
            (
                select
                    debts.WinnerId,
                    debts.LoserId,
                    debts.DebtAmount
                from TempDebts AS debts
                union
                select
                    debts.WinnerId,
                    losers."Id",
                    0
                from TempDebts as debts
                         join poker."Players" as losers on losers."Id" not in (
                    select td.LoserId from TempDebts as td
                    where td.WinnerId = debts.WinnerId
                ) and exists(
                   select * from TempDebts AS td_ where td_.LoserId = losers."Id"
               )
            ) as source
            join poker."Players" as w on w."Id" = source.WinnerId
            join poker."Players" as l on l."Id" = source.LoserId
        order by w."Name", l."Name";

    drop table TempDebts;
    return;

END
$$;

alter function poker.playersdebts() owner to postgres;

create or replace function poker.gamesinfo()
    returns TABLE
            (
                gameid int,
                gamename character varying,
                playername character varying,
                playeramount integer,
                gamedate date,
                commongamebank int
            )
    language plpgsql
as
$$
BEGIN
    return query
        select
            games."Id",
            games."Name",
            players."Name",
            cast(games_info.Amount as int) as Amount,
            games."Date",
            coalesce(game_bank.Bank, 0)
        from (
                 select debts."GameId"      as GameId,
                        debts."WinnerId"    as PlayerId,
                        sum(debts."Amount") as Amount
                 from poker."Debts" as debts
                 group by debts."GameId", debts."WinnerId"
                 union
                 select debts."GameId"      as GameId,
                        debts."LoserId"     as PlayerId,
                        -sum(debts."Amount") as Amount
                 from poker."Debts" as debts
                 group by debts."GameId", debts."LoserId"
             ) as games_info
                 join poker."Games" as games ON games."Id" = games_info.GameId
                 left join (
            select
                g."GameId" as GameId,
                cast(sum(g."Amount") as int) as Bank
            from poker."Debts" as g
            group by g."GameId"
        ) as game_bank ON game_bank.GameId = games."Id"
                 join poker."Players" as players ON players."Id" = games_info.PlayerId
        order by
            games."Id" desc,
            games."Date" desc nulls last,
            games_info.Amount desc;
END
$$;

alter function poker.gamesinfo() owner to postgres;

create type poker.PlayerGameResult AS
(
    WinnerId int,
    LoserId int,
    Amount int
);

create function poker.InsertGameResult
(
    gameName varchar(40),
    gameDate date,
    results poker.PlayerGameResult[]
)
    returns void
    language plpgsql
AS $$
declare
    newGameId int;
begin
    insert into poker."Games"("Name", "Date")
    values (gameName, gameDate)
    returning "Id" into newGameId;

    insert into poker."Debts" ("WinnerId", "LoserId", "Amount", "GameId")
    select
        (game_results.res).WinnerId,
        (game_results.res).LoserId,
        (game_results.res).Amount,
        newGameId
    from (
             select unnest(results) as res
         ) as game_results;
end;
$$;

alter function poker.InsertGameResult(gameName varchar, gameDate date, results poker.PlayerGameResult[]) owner to postgres;