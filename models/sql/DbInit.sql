create schema poker;

create table poker."Games" (
   "Id" serial not null constraint games_pk primary key,
   "Name" varchar(100) not null,
   "Date" date
);

alter table poker."Games" owner to postgres;

INSERT INTO poker."Games" ("Id", "Name", "Date") VALUES (1, 'Банк до начала ведения учёта', null);
INSERT INTO poker."Games" ("Id", "Name", "Date") VALUES (2, 'Игра у Жени', '2020-07-03');
INSERT INTO poker."Games" ("Id", "Name", "Date") VALUES (3, 'Игра у Егора', '2020-07-06');
INSERT INTO poker."Games" ("Id", "Name", "Date") VALUES (4, 'Игра у Жени', '2020-07-10');
INSERT INTO poker."Games" ("Id", "Name", "Date") VALUES (5, 'Игра у Егора', '2020-07-12');

create table poker."Players" (
    "Id" serial not null constraint players_pk primary key,
    "Name" varchar(20) not null,
    "Surname" varchar(30)
);

alter table poker."Players" owner to postgres;

INSERT INTO poker."Players" ("Id", "Name", "Surname") VALUES (1, 'Егор', 'Смеловский');
INSERT INTO poker."Players" ("Id", "Name", "Surname") VALUES (2, 'Евгений', 'Перельман');
INSERT INTO poker."Players" ("Id", "Name", "Surname") VALUES (3, 'Матвей', 'Минеев');
INSERT INTO poker."Players" ("Id", "Name", "Surname") VALUES (4, 'Александр', 'Тарасенко');
INSERT INTO poker."Players" ("Id", "Name", "Surname") VALUES (5, 'Роман', 'Огурешнов');
INSERT INTO poker."Players" ("Id", "Name", "Surname") VALUES (6, 'Алексей', 'Степанов');


create table poker."Debts" (
   "Id" serial not null constraint debts_pk primary key,
   "WinnerId" integer not null constraint debts_players_id_fk references poker."Players",
   "LoserId" integer not null constraint debts_players_id_fk_2 references poker."Players",
   "Amount" integer not null, "GameId" integer constraint debts_games_id_fk references poker."Games"
);

alter table poker."Debts" owner to postgres;

INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (1, 4, 1, 6000, 1);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (2, 4, 3, 1200, 1);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (3, 2, 5, 8350, 1);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (4, 2, 1, 3250, 1);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (5, 3, 1, 3000, 1);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (6, 5, 6, 2000, 2);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (7, 5, 1, 2600, 2);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (8, 4, 1, 3400, 2);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (9, 5, 2, 1000, 2);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (10, 2, 1, 3750, 3);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (11, 2, 3, 1250, 3);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (12, 4, 3, 4500, 4);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (13, 5, 1, 1150, 4);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (14, 2, 6, 550, 4);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (15, 5, 6, 450, 4);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (16, 3, 5, 4250, 5);
INSERT INTO poker."Debts" ("Id", "WinnerId", "LoserId", "Amount", "GameId") VALUES (17, 2, 5, 750, 5);

create table poker."DebtPayments" (
    "Id" serial not null constraint debtpayments_pk primary key,
    "PayerId" integer not null constraint debtpayments_players_id_fk references poker."Players",
    "RecipientId" integer not null constraint debtpayments_players_id_fk_2 references poker."Players",
    "Amount" integer not null,
    "InsertStamp" timestamp default now()
);

alter table poker."DebtPayments" owner to postgres;

INSERT INTO poker."DebtPayments" ("Id", "PayerId", "RecipientId", "InsertStamp", "Amount") VALUES (1, 6, 5, now(), 2000);
INSERT INTO poker."DebtPayments" ("Id", "PayerId", "RecipientId", "InsertStamp", "Amount") VALUES (2, 3, 2, now(), 100);
INSERT INTO poker."DebtPayments" ("Id", "PayerId", "RecipientId", "InsertStamp", "Amount") VALUES (3, 5, 2, now(), 5350);
INSERT INTO poker."DebtPayments" ("Id", "PayerId", "RecipientId", "InsertStamp", "Amount") VALUES (4, 3, 2, now(), 200);

create function poker.playersdebts()
returns TABLE
    (
        winnername      character varying,
        losername       character varying,
        playerwin       integer,
        commonplayerwin integer
    )
language plpgsql
as $$
BEGIN
    CREATE temporary table TempDebts AS
    select
       winner."Id" as WinnerId,
       loser."Id" as LoserId,
       sum(debts."Amount") - coalesce(sum(invert_debts."Amount"), 0) as DebtAmount
    from poker."Debts" as debts
      inner join poker."Players" as losers on losers."Id" = debts."LoserId"
      left join poker."Debts" as invert_debts ON
        invert_debts."WinnerId" = debts."LoserId"
        and invert_debts."LoserId" = debts."WinnerId"
      inner join poker."Players" as winner on winner."Id" = debts."WinnerId"
      inner join poker."Players" as loser on loser."Id" = debts."LoserId"
    group by winner."Id", loser."Id"
    having sum(debts."Amount") - coalesce(sum(invert_debts."Amount"), 0) > 0;

    update TempDebts AS temp_debts
        set DebtAmount = temp_debts.DebtAmount - debt_pay."Amount"
        from poker."DebtPayments" as debt_pay
        where
            debt_pay."PayerId" = temp_debts.LoserId
            and debt_pay."RecipientId" = temp_debts.WinnerId;


    return query
        select
           w."Name" AS WinnerName,
           l."Name" AS LoserName,
           cast(debts.DebtAmount AS int) AS PlayerWin,
           cast(sum(debts.DebtAmount) OVER (PARTITION BY w."Name") AS int) AS CommonPlayerWin
        from (
             select
                debts.WinnerId,
                debts.LoserId,
                debts.DebtAmount
             from TempDebts as debts
             union
             select
                debts.WinnerId,
                losers."Id",
                0
             from TempDebts as debts
              join poker."Players" as losers on
                  losers."Id" not in (
                     select td.LoserId
                         from TempDebts as td
                         where td.WinnerId = debts.WinnerId
                  ) and exists
                  (
                      select * from TempDebts AS td_ where td_.LoserId = losers."Id"
                  )
         ) as debts
         join poker."Players" as w on w."Id" = debts.WinnerId
         join poker."Players" as l on l."Id" = debts.LoserId
         order by w."Name", l."Name";

    drop table TempDebts;
    return;
END
$$;

alter function poker.playersdebts() owner to postgres;

create or replace function poker.gamesinfo()
    returns TABLE
            (
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
            games."Date" desc nulls last,
            games_info.Amount desc;
END
$$;

alter function poker.gamesinfo() owner to postgres;
