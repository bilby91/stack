create sequence "{{.Bucket}}"."{{.Ledger}}_transaction_id" owned by "{{.Bucket}}".transactions.id;
select setval('"{{.Bucket}}"."{{.Ledger}}_transaction_id"', coalesce((
    select max(id) + 1
    from "{{.Bucket}}".transactions
    where ledger = '{{ .Ledger }}'
), 1)::bigint, false);