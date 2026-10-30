package transfer

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestTransferPreservesIdentityAndIsRepeatable(t *testing.T) {
	base := os.Getenv("TEST_DATABASE_URL")
	if base == "" {
		t.Skip("set TEST_DATABASE_URL")
	}
	ctx := context.Background()
	admin, err := pgx.Connect(ctx, base)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close(ctx)
	makeSchema := func() string {
		name := "transfer_" + strings.ReplaceAll(uuid.NewString(), "-", "")
		quoted := pgx.Identifier{name}.Sanitize()
		if _, err := admin.Exec(ctx, "CREATE SCHEMA "+quoted); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			c, e := pgx.Connect(ctx, base)
			if e != nil {
				t.Error(e)
				return
			}
			defer c.Close(ctx)
			if _, e = c.Exec(ctx, "DROP SCHEMA "+quoted+" CASCADE"); e != nil {
				t.Error(e)
			}
		})
		u, e := url.Parse(base)
		if e != nil {
			t.Fatal(e)
		}
		q := u.Query()
		q.Set("search_path", name)
		u.RawQuery = q.Encode()
		return u.String()
	}
	sourceDSN, targetDSN := makeSchema(), makeSchema()
	source, err := pgx.Connect(ctx, sourceDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer source.Close(ctx)
	_, err = source.Exec(ctx, `CREATE TABLE staff(id uuid,first_name text,last_name text,role text,deleted_at timestamptz);
 CREATE TABLE accounts(id uuid,username text,password_hash text,created_at timestamptz,updated_at timestamptz,is_active bool,staff_id uuid,is_admin bool,totp_secret text,last_totp_step bigint);
 CREATE TABLE refresh_tokens(token_hash text,account_id uuid,expires_at timestamptz);
 CREATE TABLE audit_events(event_id uuid,event_time timestamptz,actor text,message text,event_type text);
 CREATE TABLE security_settings(id bool,setup_secret text);`)
	if err != nil {
		t.Fatal(err)
	}
	adminID, staffID, workerID, eventID := uuid.NewString(), uuid.NewString(), uuid.NewString(), uuid.NewString()
	_, err = source.Exec(ctx, `INSERT INTO staff VALUES($1,'А','Б','doctor',NULL)`, staffID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = source.Exec(ctx, `INSERT INTO accounts VALUES($1,'admin','admin-hash',now(),now(),true,NULL,true,'test-totp-secret',123),($2,'doctor','worker-hash',now(),now(),true,$3,false,NULL,-1)`, adminID, workerID, staffID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = source.Exec(ctx, `INSERT INTO refresh_tokens VALUES('test-refresh-hash',$1,now()+interval '1 day')`, workerID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = source.Exec(ctx, `INSERT INTO audit_events VALUES($1,now(),$2,'test event','login')`, eventID, "account:"+adminID)
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if err = Migrate(ctx, sourceDSN, targetDSN); err != nil {
			t.Fatal(err)
		}
	}
	target, err := pgx.Connect(ctx, targetDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer target.Close(ctx)
	var matches bool
	err = target.QueryRow(ctx, `SELECT
 EXISTS(SELECT 1 FROM accounts WHERE id=$1 AND password_hash='admin-hash' AND totp_secret='test-totp-secret' AND last_totp_step=123)
 AND EXISTS(SELECT 1 FROM accounts WHERE id=$2 AND staff_id=$3 AND password_hash='worker-hash')
 AND EXISTS(SELECT 1 FROM refresh_tokens WHERE account_id=$2 AND token_hash='test-refresh-hash')
 AND EXISTS(SELECT 1 FROM audit_events WHERE event_id=$4 AND event_type='login')
 AND (SELECT count(*) FROM accounts)=2`, adminID, workerID, staffID, eventID).Scan(&matches)
	if err != nil || !matches {
		t.Fatalf("transfer did not preserve identity data: %v", err)
	}
	var marked bool
	if err = source.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM _identity_transfer)").Scan(&marked); err != nil || !marked {
		t.Fatalf("missing safe-transfer marker: %v", err)
	}
	if _, err = source.Exec(ctx, "UPDATE accounts SET username='changed' WHERE id=$1", workerID); err != nil {
		t.Fatal(err)
	}
	if err = Migrate(ctx, sourceDSN, targetDSN); err == nil {
		t.Fatal("changed legacy data silently discarded on retry")
	}

}
