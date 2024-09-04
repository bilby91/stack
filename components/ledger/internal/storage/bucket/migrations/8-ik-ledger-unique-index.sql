drop index "{{.Bucket}}_logs_idempotency_key";

create unique index "{{.Bucket}}_logs_idempotency_key" on "{{.Bucket}}".logs (ledger, idempotency_key);