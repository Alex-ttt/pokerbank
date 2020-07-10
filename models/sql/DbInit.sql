create schema poker;

create table poker."Games"
(
    "Id"   serial       not null
        constraint games_pk
            primary key,
    "Name" varchar(100) not null,
    "Date" date
);
INSERT INTO poker."Games" ("Id", "Name", "Date") VALUES (1, 'Долги до начала ведения учёта', null);
INSERT INTO poker."Games" ("Id", "Name", "Date") VALUES (2, 'Игра у Жени Перельмана', '2020-07-03');
INSERT INTO poker."Games" ("Id", "Name", "Date") VALUES (3, 'Игра у Егора', '2020-07-06');
INSERT INTO poker."Games" ("Id", "Name", "Date") VALUES (4, 'Игра у Жени', '2020-07-10');

alter table poker."Games" owner to postgres;

create table poker."Players"
(
    "Id"      serial      not null
        constraint players_pk
            primary key,
    "Name"    varchar(20) not null,
    "Surname" varchar(30)
);
INSERT INTO poker."Players" ("Id", "Name", "Surname") VALUES (1, 'Егор', 'Смеловский');
INSERT INTO poker."Players" ("Id", "Name", "Surname") VALUES (2, 'Евгений', 'Перельман');
INSERT INTO poker."Players" ("Id", "Name", "Surname") VALUES (3, 'Матвей', 'Минеев');
INSERT INTO poker."Players" ("Id", "Name", "Surname") VALUES (4, 'Александр', 'Тарасенко');
INSERT INTO poker."Players" ("Id", "Name", "Surname") VALUES (5, 'Роман', 'Огурешнов');
INSERT INTO poker."Players" ("Id", "Name", "Surname") VALUES (6, 'Алексей', 'Степанов');


alter table poker."Players"
    owner to postgres;

create table poker."Debts"
(
    "Id"       serial  not null
        constraint debts_pk
            primary key,
    "WinnerId" integer not null
        constraint debts_players_id_fk
            references poker."Players",
    "LoserId"  integer not null
        constraint debts_players_id_fk_2
            references poker."Players",
    "Amount"   integer not null,
    "GameId"   integer
        constraint debts_games_id_fk
            references poker."Games"
);
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

alter table poker."Debts"
    owner to postgres;

create table poker."DebtPayments"
(
    "Id"          serial  not null
        constraint debtpayments_pk
            primary key,
    "PayerId"     integer not null
        constraint debtpayments_players_id_fk
            references poker."Players",
    "RecipientId" integer not null
        constraint debtpayments_players_id_fk_2
            references poker."Players",
    "InsertStamp" timestamp default now(),
    "Amount"      integer not null
);
INSERT INTO poker."DebtPayments" ("Id", "PayerId", "RecipientId", "InsertStamp", "Amount") VALUES (1, 6, 5, '2020-07-09 07:55:55.458594', 2000);
INSERT INTO poker."DebtPayments" ("Id", "PayerId", "RecipientId", "InsertStamp", "Amount") VALUES (2, 3, 2, '2020-07-09 07:55:55.458594', 100);

alter table poker."DebtPayments"
    owner to postgres;

create function poker.playersdebts()
    returns TABLE
            (
                winnername      character varying,
                losername       character varying,
                playerwin       integer,
                commonplayerwin integer
            )
    language plpgsql
as
$$
BEGIN
    CREATE temporary table TempDebts AS
    select win."Id"  AS WinnerId,
           lose."Id" AS LoserId,
           sum(d."Amount")
               - coalesce(sum(invert_d."Amount"), 0)
                     as DebtAmount
    from poker."Debts" as d
             join poker."Players" as p on p."Id" = d."LoserId"
             left join poker."Debts" as invert_d ON
                invert_d."WinnerId" = d."LoserId"
            and invert_d."LoserId" = d."WinnerId"
             left join poker."Players" as win on win."Id" = d."WinnerId"
             left join poker."Players" as lose on lose."Id" = d."LoserId"
    group by win."Id", lose."Id"
    having sum(d."Amount")
               - coalesce(sum(invert_d."Amount"), 0)
               > 0;

    update TempDebts AS temp_debts
    set DebtAmount = temp_debts.DebtAmount - dp."Amount"
    from poker."DebtPayments" as dp
    where dp."PayerId" = temp_debts.LoserId
      and dp."RecipientId" = temp_debts.WinnerId;


    return query
        select w."Name"                       AS WinnerName,
               l."Name"                       AS LoserName,
               cast(source.DebtAmount AS int) AS PlayerWin,
               cast(
                               sum(source.DebtAmount) OVER (PARTITION BY w."Name")
                   AS int
                   )
                                              AS CommonPlayerWin
        from (
                 select debts.WinnerId,
                        debts.LoserId,
                        debts.DebtAmount
                 from TempDebts AS debts
                 union
                 select debts.WinnerId,
                        losers."Id",
                        0
                 from TempDebts as debts
                          join poker."Players" as losers on losers."Id" not in (
                     select td.LoserId
                     from TempDebts as td
                     where td.WinnerId = debts.WinnerId
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

