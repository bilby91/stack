drop trigger "{{.Bucket}}_insert_log" on "{{.Bucket}}".logs;

alter table "{{.Bucket}}".transactions
add column inserted_at timestamp without time zone
default now();

alter table "{{.Bucket}}".transactions
alter column timestamp
set default now();

alter table "{{.Bucket}}".transactions
alter column id
type bigint;

-- create function "{{.Bucket}}".set_effective_volumes()
--     returns trigger
--     security definer
--     language plpgsql
-- as
-- $$
-- begin
--     new.post_commit_effective_volumes = coalesce((
--         select (
--             (post_commit_effective_volumes).inputs + case when new.is_source then 0 else new.amount end,
--             (post_commit_effective_volumes).outputs + case when new.is_source then new.amount else 0 end
--         )
--         from "{{.Bucket}}".moves
--         where accounts_seq = new.accounts_seq
--             and asset = new.asset
--             and ledger = new.ledger
--             and (effective_date < new.effective_date or (effective_date = new.effective_date and seq < new.seq))
--         order by effective_date desc, seq desc
--         limit 1
--     ), (
--         case when new.is_source then 0 else new.amount end,
--         case when new.is_source then new.amount else 0 end
--     ));
--
--     return new;
-- end;
-- $$;
--
-- create trigger "{{.Bucket}}_set_effective_volumes"
-- before insert
-- on "{{.Bucket}}"."moves"
-- for each row
-- execute procedure "{{.Bucket}}".set_effective_volumes();
--
-- create function "{{.Bucket}}".update_effective_volumes()
--     returns trigger
--     security definer
--     language plpgsql
-- as
-- $$
-- begin
--     update "{{.Bucket}}".moves
--     set post_commit_effective_volumes =
--             (
--              (post_commit_effective_volumes).inputs + case when new.is_source then 0 else new.amount end,
--              (post_commit_effective_volumes).outputs + case when new.is_source then new.amount else 0 end
--                 )
--     where accounts_seq = new.accounts_seq
--         and asset = new.asset
--         and effective_date > new.effective_date
--         and ledger = new.ledger;
--
--     return new;
-- end;
-- $$;
--
-- create trigger "{{.Bucket}}_update_effective_volumes"
--     after insert
--     on "{{.Bucket}}"."moves"
--     for each row
-- execute procedure "{{.Bucket}}".update_effective_volumes();

-- todo: need to populate balances with existing data
create table "{{.Bucket}}".balances (
    ledger varchar,
    account varchar,
    asset varchar,
    balance numeric,

    primary key (ledger, account, asset)
);

alter table "{{.Bucket}}".moves
alter column post_commit_volumes
drop not null,
alter column post_commit_effective_volumes
drop not null;