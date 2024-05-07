// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: batch.sql

package batchsqlc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const bulkInsertIntoBatchRows = `-- name: BulkInsertIntoBatchRows :execrows
INSERT INTO batchrows (batch, line, input, status, reqat) 
VALUES 
    (unnest($1::uuid[]), unnest($2::int[]), unnest($3::jsonb[]), 'queued', NOW())
`

type BulkInsertIntoBatchRowsParams struct {
	Batch []uuid.UUID `json:"batch"`
	Line  []int32     `json:"line"`
	Input [][]byte    `json:"input"`
}

func (q *Queries) BulkInsertIntoBatchRows(ctx context.Context, arg BulkInsertIntoBatchRowsParams) (int64, error) {
	result, err := q.db.Exec(ctx, bulkInsertIntoBatchRows, arg.Batch, arg.Line, arg.Input)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const countBatchRowsByBatchIDAndStatus = `-- name: CountBatchRowsByBatchIDAndStatus :one
SELECT COUNT(*)
FROM batchrows
WHERE batch = $1 AND status IN ($2, $3)
`

type CountBatchRowsByBatchIDAndStatusParams struct {
	Batch    uuid.UUID  `json:"batch"`
	Status   StatusEnum `json:"status"`
	Status_2 StatusEnum `json:"status_2"`
}

func (q *Queries) CountBatchRowsByBatchIDAndStatus(ctx context.Context, arg CountBatchRowsByBatchIDAndStatusParams) (int64, error) {
	row := q.db.QueryRow(ctx, countBatchRowsByBatchIDAndStatus, arg.Batch, arg.Status, arg.Status_2)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const fetchBatchRowsForBatchDone = `-- name: FetchBatchRowsForBatchDone :many
SELECT line, status, res, messages
FROM batchrows
WHERE batch = $1
`

type FetchBatchRowsForBatchDoneRow struct {
	Line     int32      `json:"line"`
	Status   StatusEnum `json:"status"`
	Res      []byte     `json:"res"`
	Messages []byte     `json:"messages"`
}

func (q *Queries) FetchBatchRowsForBatchDone(ctx context.Context, batch uuid.UUID) ([]FetchBatchRowsForBatchDoneRow, error) {
	rows, err := q.db.Query(ctx, fetchBatchRowsForBatchDone, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchBatchRowsForBatchDoneRow
	for rows.Next() {
		var i FetchBatchRowsForBatchDoneRow
		if err := rows.Scan(
			&i.Line,
			&i.Status,
			&i.Res,
			&i.Messages,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchBlockOfRows = `-- name: FetchBlockOfRows :many
SELECT batches.app, batches.status, batches.op, batches.context, batchrows.batch, batchrows.rowid, batchrows.line, batchrows.input
FROM batchrows
INNER JOIN batches ON batchrows.batch = batches.id
WHERE batchrows.status = $1 AND batches.status != 'wait'
LIMIT $2
FOR UPDATE OF batchrows SKIP LOCKED
`

type FetchBlockOfRowsParams struct {
	Status StatusEnum `json:"status"`
	Limit  int32      `json:"limit"`
}

type FetchBlockOfRowsRow struct {
	App     string     `json:"app"`
	Status  StatusEnum `json:"status"`
	Op      string     `json:"op"`
	Context []byte     `json:"context"`
	Batch   uuid.UUID  `json:"batch"`
	Rowid   int32      `json:"rowid"`
	Line    int32      `json:"line"`
	Input   []byte     `json:"input"`
}

func (q *Queries) FetchBlockOfRows(ctx context.Context, arg FetchBlockOfRowsParams) ([]FetchBlockOfRowsRow, error) {
	rows, err := q.db.Query(ctx, fetchBlockOfRows, arg.Status, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchBlockOfRowsRow
	for rows.Next() {
		var i FetchBlockOfRowsRow
		if err := rows.Scan(
			&i.App,
			&i.Status,
			&i.Op,
			&i.Context,
			&i.Batch,
			&i.Rowid,
			&i.Line,
			&i.Input,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBatchByID = `-- name: GetBatchByID :one
SELECT id, app, op, context, inputfile, status, reqat, doneat, outputfiles, nsuccess, nfailed, naborted
FROM batches
WHERE id = $1 
FOR UPDATE
`

func (q *Queries) GetBatchByID(ctx context.Context, id uuid.UUID) (Batch, error) {
	row := q.db.QueryRow(ctx, getBatchByID, id)
	var i Batch
	err := row.Scan(
		&i.ID,
		&i.App,
		&i.Op,
		&i.Context,
		&i.Inputfile,
		&i.Status,
		&i.Reqat,
		&i.Doneat,
		&i.Outputfiles,
		&i.Nsuccess,
		&i.Nfailed,
		&i.Naborted,
	)
	return i, err
}

const getBatchRowsByBatchID = `-- name: GetBatchRowsByBatchID :many
SELECT rowid, batch, line, input, status, reqat, doneat, res, blobrows, messages, doneby FROM batchrows WHERE batch = $1
`

func (q *Queries) GetBatchRowsByBatchID(ctx context.Context, batch uuid.UUID) ([]Batchrow, error) {
	rows, err := q.db.Query(ctx, getBatchRowsByBatchID, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Batchrow
	for rows.Next() {
		var i Batchrow
		if err := rows.Scan(
			&i.Rowid,
			&i.Batch,
			&i.Line,
			&i.Input,
			&i.Status,
			&i.Reqat,
			&i.Doneat,
			&i.Res,
			&i.Blobrows,
			&i.Messages,
			&i.Doneby,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBatchRowsByBatchIDSorted = `-- name: GetBatchRowsByBatchIDSorted :many
SELECT rowid, line, input, status, reqat, doneat, res, blobrows, messages, doneby
FROM batchrows
WHERE batch = $1
ORDER BY line
FOR UPDATE
`

type GetBatchRowsByBatchIDSortedRow struct {
	Rowid    int32            `json:"rowid"`
	Line     int32            `json:"line"`
	Input    []byte           `json:"input"`
	Status   StatusEnum       `json:"status"`
	Reqat    pgtype.Timestamp `json:"reqat"`
	Doneat   pgtype.Timestamp `json:"doneat"`
	Res      []byte           `json:"res"`
	Blobrows []byte           `json:"blobrows"`
	Messages []byte           `json:"messages"`
	Doneby   pgtype.Text      `json:"doneby"`
}

func (q *Queries) GetBatchRowsByBatchIDSorted(ctx context.Context, batch uuid.UUID) ([]GetBatchRowsByBatchIDSortedRow, error) {
	rows, err := q.db.Query(ctx, getBatchRowsByBatchIDSorted, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBatchRowsByBatchIDSortedRow
	for rows.Next() {
		var i GetBatchRowsByBatchIDSortedRow
		if err := rows.Scan(
			&i.Rowid,
			&i.Line,
			&i.Input,
			&i.Status,
			&i.Reqat,
			&i.Doneat,
			&i.Res,
			&i.Blobrows,
			&i.Messages,
			&i.Doneby,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBatchRowsCount = `-- name: GetBatchRowsCount :one
SELECT COUNT(*) FROM batchrows WHERE batch = $1
`

func (q *Queries) GetBatchRowsCount(ctx context.Context, batch uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, getBatchRowsCount, batch)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getBatchStatus = `-- name: GetBatchStatus :one
SELECT status
FROM batches
WHERE id = $1
`

func (q *Queries) GetBatchStatus(ctx context.Context, id uuid.UUID) (StatusEnum, error) {
	row := q.db.QueryRow(ctx, getBatchStatus, id)
	var status StatusEnum
	err := row.Scan(&status)
	return status, err
}

const getCompletedBatches = `-- name: GetCompletedBatches :many
SELECT id
FROM batches
WHERE status IN ('success', 'failed', 'aborted')
FOR UPDATE
`

func (q *Queries) GetCompletedBatches(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := q.db.Query(ctx, getCompletedBatches)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPendingBatchRows = `-- name: GetPendingBatchRows :many
SELECT rowid, line, input, status, reqat, doneat, res, blobrows, messages, doneby
FROM batchrows
WHERE batch = $1 AND status IN ('queued', 'inprog')
FOR UPDATE
`

type GetPendingBatchRowsRow struct {
	Rowid    int32            `json:"rowid"`
	Line     int32            `json:"line"`
	Input    []byte           `json:"input"`
	Status   StatusEnum       `json:"status"`
	Reqat    pgtype.Timestamp `json:"reqat"`
	Doneat   pgtype.Timestamp `json:"doneat"`
	Res      []byte           `json:"res"`
	Blobrows []byte           `json:"blobrows"`
	Messages []byte           `json:"messages"`
	Doneby   pgtype.Text      `json:"doneby"`
}

func (q *Queries) GetPendingBatchRows(ctx context.Context, batch uuid.UUID) ([]GetPendingBatchRowsRow, error) {
	rows, err := q.db.Query(ctx, getPendingBatchRows, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPendingBatchRowsRow
	for rows.Next() {
		var i GetPendingBatchRowsRow
		if err := rows.Scan(
			&i.Rowid,
			&i.Line,
			&i.Input,
			&i.Status,
			&i.Reqat,
			&i.Doneat,
			&i.Res,
			&i.Blobrows,
			&i.Messages,
			&i.Doneby,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProcessedBatchRowsByBatchIDSorted = `-- name: GetProcessedBatchRowsByBatchIDSorted :many
SELECT rowid, line, input, status, reqat, doneat, res, blobrows, messages, doneby
FROM batchrows
WHERE batch = $1 AND status IN ('success', 'failed')
ORDER BY line
FOR UPDATE
`

type GetProcessedBatchRowsByBatchIDSortedRow struct {
	Rowid    int32            `json:"rowid"`
	Line     int32            `json:"line"`
	Input    []byte           `json:"input"`
	Status   StatusEnum       `json:"status"`
	Reqat    pgtype.Timestamp `json:"reqat"`
	Doneat   pgtype.Timestamp `json:"doneat"`
	Res      []byte           `json:"res"`
	Blobrows []byte           `json:"blobrows"`
	Messages []byte           `json:"messages"`
	Doneby   pgtype.Text      `json:"doneby"`
}

func (q *Queries) GetProcessedBatchRowsByBatchIDSorted(ctx context.Context, batch uuid.UUID) ([]GetProcessedBatchRowsByBatchIDSortedRow, error) {
	rows, err := q.db.Query(ctx, getProcessedBatchRowsByBatchIDSorted, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProcessedBatchRowsByBatchIDSortedRow
	for rows.Next() {
		var i GetProcessedBatchRowsByBatchIDSortedRow
		if err := rows.Scan(
			&i.Rowid,
			&i.Line,
			&i.Input,
			&i.Status,
			&i.Reqat,
			&i.Doneat,
			&i.Res,
			&i.Blobrows,
			&i.Messages,
			&i.Doneby,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertIntoBatchRows = `-- name: InsertIntoBatchRows :exec
INSERT INTO batchrows (batch, line, input, status, reqat)
VALUES ($1, $2, $3, 'queued', NOW())
`

type InsertIntoBatchRowsParams struct {
	Batch uuid.UUID `json:"batch"`
	Line  int32     `json:"line"`
	Input []byte    `json:"input"`
}

func (q *Queries) InsertIntoBatchRows(ctx context.Context, arg InsertIntoBatchRowsParams) error {
	_, err := q.db.Exec(ctx, insertIntoBatchRows, arg.Batch, arg.Line, arg.Input)
	return err
}

const insertIntoBatches = `-- name: InsertIntoBatches :one
INSERT INTO batches (id, app, op, context, status, reqat)
VALUES ($1, $2, $3, $4, $5, NOW())
RETURNING id
`

type InsertIntoBatchesParams struct {
	ID      uuid.UUID  `json:"id"`
	App     string     `json:"app"`
	Op      string     `json:"op"`
	Context []byte     `json:"context"`
	Status  StatusEnum `json:"status"`
}

func (q *Queries) InsertIntoBatches(ctx context.Context, arg InsertIntoBatchesParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, insertIntoBatches,
		arg.ID,
		arg.App,
		arg.Op,
		arg.Context,
		arg.Status,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const updateBatchCounters = `-- name: UpdateBatchCounters :exec
UPDATE batches
SET nsuccess = COALESCE(nsuccess, 0) + $2,
    nfailed = COALESCE(nfailed, 0) + $3,
    naborted = COALESCE(naborted, 0) + $4
WHERE id = $1
`

type UpdateBatchCountersParams struct {
	ID       uuid.UUID   `json:"id"`
	Nsuccess pgtype.Int4 `json:"nsuccess"`
	Nfailed  pgtype.Int4 `json:"nfailed"`
	Naborted pgtype.Int4 `json:"naborted"`
}

func (q *Queries) UpdateBatchCounters(ctx context.Context, arg UpdateBatchCountersParams) error {
	_, err := q.db.Exec(ctx, updateBatchCounters,
		arg.ID,
		arg.Nsuccess,
		arg.Nfailed,
		arg.Naborted,
	)
	return err
}

const updateBatchOutputFiles = `-- name: UpdateBatchOutputFiles :exec
UPDATE batches
SET outputfiles = $2
WHERE id = $1
`

type UpdateBatchOutputFilesParams struct {
	ID          uuid.UUID `json:"id"`
	Outputfiles []byte    `json:"outputfiles"`
}

func (q *Queries) UpdateBatchOutputFiles(ctx context.Context, arg UpdateBatchOutputFilesParams) error {
	_, err := q.db.Exec(ctx, updateBatchOutputFiles, arg.ID, arg.Outputfiles)
	return err
}

const updateBatchRowStatus = `-- name: UpdateBatchRowStatus :exec
UPDATE batchrows
SET status = $2
WHERE rowid = $1
`

type UpdateBatchRowStatusParams struct {
	Rowid  int32      `json:"rowid"`
	Status StatusEnum `json:"status"`
}

func (q *Queries) UpdateBatchRowStatus(ctx context.Context, arg UpdateBatchRowStatusParams) error {
	_, err := q.db.Exec(ctx, updateBatchRowStatus, arg.Rowid, arg.Status)
	return err
}

const updateBatchRowsBatchJob = `-- name: UpdateBatchRowsBatchJob :exec
UPDATE batchrows
SET status = $2, doneat = $3, res = $4, blobrows = $5, messages = $6, doneby = $7
WHERE rowid = $1
`

type UpdateBatchRowsBatchJobParams struct {
	Rowid    int32            `json:"rowid"`
	Status   StatusEnum       `json:"status"`
	Doneat   pgtype.Timestamp `json:"doneat"`
	Res      []byte           `json:"res"`
	Blobrows []byte           `json:"blobrows"`
	Messages []byte           `json:"messages"`
	Doneby   pgtype.Text      `json:"doneby"`
}

func (q *Queries) UpdateBatchRowsBatchJob(ctx context.Context, arg UpdateBatchRowsBatchJobParams) error {
	_, err := q.db.Exec(ctx, updateBatchRowsBatchJob,
		arg.Rowid,
		arg.Status,
		arg.Doneat,
		arg.Res,
		arg.Blobrows,
		arg.Messages,
		arg.Doneby,
	)
	return err
}

const updateBatchRowsSlowQuery = `-- name: UpdateBatchRowsSlowQuery :exec
UPDATE batchrows
SET status = $2, doneat = $3, res = $4, messages = $5, doneby = $6
WHERE rowid = $1
`

type UpdateBatchRowsSlowQueryParams struct {
	Rowid    int32            `json:"rowid"`
	Status   StatusEnum       `json:"status"`
	Doneat   pgtype.Timestamp `json:"doneat"`
	Res      []byte           `json:"res"`
	Messages []byte           `json:"messages"`
	Doneby   pgtype.Text      `json:"doneby"`
}

func (q *Queries) UpdateBatchRowsSlowQuery(ctx context.Context, arg UpdateBatchRowsSlowQueryParams) error {
	_, err := q.db.Exec(ctx, updateBatchRowsSlowQuery,
		arg.Rowid,
		arg.Status,
		arg.Doneat,
		arg.Res,
		arg.Messages,
		arg.Doneby,
	)
	return err
}

const updateBatchRowsStatus = `-- name: UpdateBatchRowsStatus :exec
UPDATE batchrows
SET status = $1
WHERE rowid = ANY($2::int[])
`

type UpdateBatchRowsStatusParams struct {
	Status  StatusEnum `json:"status"`
	Column2 []int32    `json:"column_2"`
}

func (q *Queries) UpdateBatchRowsStatus(ctx context.Context, arg UpdateBatchRowsStatusParams) error {
	_, err := q.db.Exec(ctx, updateBatchRowsStatus, arg.Status, arg.Column2)
	return err
}

const updateBatchStatus = `-- name: UpdateBatchStatus :exec
UPDATE batches
SET status = $2, doneat = $3, outputfiles = $4, nsuccess = $5, nfailed = $6, naborted = $7
WHERE id = $1
`

type UpdateBatchStatusParams struct {
	ID          uuid.UUID        `json:"id"`
	Status      StatusEnum       `json:"status"`
	Doneat      pgtype.Timestamp `json:"doneat"`
	Outputfiles []byte           `json:"outputfiles"`
	Nsuccess    pgtype.Int4      `json:"nsuccess"`
	Nfailed     pgtype.Int4      `json:"nfailed"`
	Naborted    pgtype.Int4      `json:"naborted"`
}

func (q *Queries) UpdateBatchStatus(ctx context.Context, arg UpdateBatchStatusParams) error {
	_, err := q.db.Exec(ctx, updateBatchStatus,
		arg.ID,
		arg.Status,
		arg.Doneat,
		arg.Outputfiles,
		arg.Nsuccess,
		arg.Nfailed,
		arg.Naborted,
	)
	return err
}

const updateBatchSummary = `-- name: UpdateBatchSummary :exec
UPDATE batches
SET status = $2, doneat = $3, outputfiles = $4, nsuccess = $5, nfailed = $6, naborted = $7
WHERE id = $1
`

type UpdateBatchSummaryParams struct {
	ID          uuid.UUID        `json:"id"`
	Status      StatusEnum       `json:"status"`
	Doneat      pgtype.Timestamp `json:"doneat"`
	Outputfiles []byte           `json:"outputfiles"`
	Nsuccess    pgtype.Int4      `json:"nsuccess"`
	Nfailed     pgtype.Int4      `json:"nfailed"`
	Naborted    pgtype.Int4      `json:"naborted"`
}

func (q *Queries) UpdateBatchSummary(ctx context.Context, arg UpdateBatchSummaryParams) error {
	_, err := q.db.Exec(ctx, updateBatchSummary,
		arg.ID,
		arg.Status,
		arg.Doneat,
		arg.Outputfiles,
		arg.Nsuccess,
		arg.Nfailed,
		arg.Naborted,
	)
	return err
}

const updateBatchSummaryOnAbort = `-- name: UpdateBatchSummaryOnAbort :exec
UPDATE batches
SET status = $2, doneat = $3, naborted = $4
WHERE id = $1
`

type UpdateBatchSummaryOnAbortParams struct {
	ID       uuid.UUID        `json:"id"`
	Status   StatusEnum       `json:"status"`
	Doneat   pgtype.Timestamp `json:"doneat"`
	Naborted pgtype.Int4      `json:"naborted"`
}

func (q *Queries) UpdateBatchSummaryOnAbort(ctx context.Context, arg UpdateBatchSummaryOnAbortParams) error {
	_, err := q.db.Exec(ctx, updateBatchSummaryOnAbort,
		arg.ID,
		arg.Status,
		arg.Doneat,
		arg.Naborted,
	)
	return err
}
