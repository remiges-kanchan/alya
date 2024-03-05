// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package batchsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	FetchBatchRowsData(ctx context.Context, batch pgtype.UUID) ([]FetchBatchRowsDataRow, error)
	FetchBlockOfRows(ctx context.Context, arg FetchBlockOfRowsParams) ([]FetchBlockOfRowsRow, error)
	GetBatchStatus(ctx context.Context, id pgtype.UUID) (StatusEnum, error)
	InsertIntoBatchRows(ctx context.Context, arg InsertIntoBatchRowsParams) error
	InsertIntoBatches(ctx context.Context, arg InsertIntoBatchesParams) (pgtype.UUID, error)
	UpdateBatchOutputFiles(ctx context.Context, arg UpdateBatchOutputFilesParams) error
	UpdateBatchRowsBatchJob(ctx context.Context, arg UpdateBatchRowsBatchJobParams) error
	UpdateBatchRowsSlowQuery(ctx context.Context, arg UpdateBatchRowsSlowQueryParams) error
	UpdateBatchRowsStatus(ctx context.Context, arg UpdateBatchRowsStatusParams) error
}

var _ Querier = (*Queries)(nil)
