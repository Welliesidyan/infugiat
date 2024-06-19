import (
	"context"
	"fmt"
	"io"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func addVectorType(ctx context.Context, w io.Writer, db string) error {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	op, err := adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			"CREATE TYPE Singers AS STRUCT<FirstName STRING(1024), LastName STRING(1024), SingerInfo BYTES(MAX)>",
			"CREATE TYPE Venues AS STRUCT<VenueId INT64, VenueName STRING(1024), VenueInfo BYTES(MAX)>",
			"CREATE TABLE Performances (SingerId INT64, VenueId INT64, EventDate DATE, Revenue NUMERIC, LastUpdateTime TIMESTAMP OPTIONS (allow_commit_timestamp=true), VenueDetails Venues, SingerDetails Singers) PRIMARY KEY (SingerId, VenueId, EventDate), INTERLEAVE IN PARENT Singers ON DELETE CASCADE",
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Added vector type\n")
	return nil
}
  
