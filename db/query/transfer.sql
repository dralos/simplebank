-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id, to_account_id, amount)
VALUES ($1, $2, $3)
RETURNING id, from_account_id, to_account_id, amount, created_at;

-- name: GetTransfer :one
SELECT id, from_account_id, to_account_id, amount, created_at
FROM transfers
WHERE id = $1;

-- name: ListTransfers :many
SELECT id, from_account_id, to_account_id, amount, created_at
FROM transfers
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateTransferAmount :one
UPDATE transfers
SET amount = $2
WHERE id = $1
RETURNING id, from_account_id, to_account_id, amount, created_at;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;
