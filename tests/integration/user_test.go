package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"webnote/internal/adapters/postgres"
	usermodels "webnote/internal/user/models"
	userrepo "webnote/internal/user/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testDBPool *pgxpool.Pool
	userRepo   *userrepo.PostgresUserRepository
)

const testPassHash = "passhash"

func TestMain(m *testing.M) {
	fmt.Println(os.Environ())
	cfg := postgres.MustNewConfig()
	pool := postgres.MustNewPool(context.Background(), cfg)

	testDBPool = pool
	r := userrepo.NewPostgresUserRepository(testDBPool)
	userRepo = r
	code := m.Run()

	testDBPool.Close()

	os.Exit(code)
}

func TestUserGetByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("Context created")

	u1Email := "test1@test.loc"
	u1 := usermodels.User{
		Email:    u1Email,
		Password: testPassHash,
	}

	id, err := userRepo.CreateUser(ctx, u1)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatalf("user create error %v", err)
	}

	res, err := userRepo.GetUserByID(ctx, id)
	if err != nil {
		t.Fatalf("getUserByID err %v", err)
	}

	if res.ID != id {
		t.Fatalf("expected %d, got %d", id, res.ID)
	}
}

func TestUserGetByEmail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u2Email := "test2@test.loc"
	u2 := usermodels.User{
		Email:    u2Email,
		Password: testPassHash,
	}

	id, err := userRepo.CreateUser(ctx, u2)
	if err != nil {
		t.Fatalf("user create error %v", err)
	}

	user, err := userRepo.GetUserByID(ctx, id)
	if err != nil {
		t.Fatalf("getUserByID err %v", err)
	}

	res, err := userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("getUserByEmail err %v", err)
	}

	if res.Email != u2Email {
		t.Fatalf("expected %s, got %s", u2Email, res.Email)
	}
}

func TestUserCreateSession(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u3Email := "test3@test.loc"
	u3 := usermodels.User{
		Email:    u3Email,
		Password: testPassHash,
	}

	id, err := userRepo.CreateUser(ctx, u3)
	if err != nil {
		t.Fatalf("user create error: %v", err)
	}

	var (
		sessToken = "sessToken3"
		csrfToken = "csrfToken3" // #nosec G101
		maxAge    = 3600
		expireAt  = time.Now().Add(time.Hour * 1)
	)
	tmpSess := usermodels.Session{
		UserID:       id,
		SessionToken: sessToken,
		CSRFToken:    csrfToken,
		MaxAge:       maxAge,
		ExpireAt:     expireAt,
	}
	u3Sess, err := userRepo.CreateSession(ctx, tmpSess)
	if err != nil {
		t.Fatalf("err while creating session: %s", err)
	}

	if u3Sess.UserID != tmpSess.UserID {
		t.Fatalf("expected %d, got %d", tmpSess.UserID, u3Sess.UserID)
	}

	if u3Sess.SessionToken != tmpSess.SessionToken {
		t.Fatalf("expected %s, got %s", tmpSess.SessionToken, u3Sess.SessionToken)
	}

	if u3Sess.CSRFToken != tmpSess.CSRFToken {
		t.Fatalf("expected %s, got %s", tmpSess.CSRFToken, u3Sess.CSRFToken)
	}

	if u3Sess.MaxAge != tmpSess.MaxAge {
		t.Fatalf("expected %d, got %d", tmpSess.MaxAge, u3Sess.MaxAge)
	}
}

func TestUserGetSessionByToken(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u4Email := "test4@test.loc"
	u4 := usermodels.User{
		Email:    u4Email,
		Password: testPassHash,
	}

	id, err := userRepo.CreateUser(ctx, u4)
	if err != nil {
		t.Fatalf("user create error: %v", err)
	}

	var (
		sessToken = "sessToken4"
		csrfToken = "csrfToken4"
		maxAge    = 3600
		expireAt  = time.Now().Add(time.Hour * 1)
	)

	tmpSess := usermodels.Session{
		UserID:       id,
		SessionToken: sessToken,
		CSRFToken:    csrfToken,
		MaxAge:       maxAge,
		ExpireAt:     expireAt,
	}

	u4Sess, err := userRepo.CreateSession(ctx, tmpSess)
	if err != nil {
		t.Fatalf("err while creating session: %s", err)
	}

	res, err := userRepo.GetSessionByToken(ctx, u4Sess.SessionToken)
	if err != nil {
		t.Fatalf("")
	}

	if u4Sess.UserID != res.UserID {
		t.Fatalf("expected %s, got %s", res.SessionToken, u4Sess.SessionToken)
	}

	if u4Sess.SessionToken != res.SessionToken {
		t.Fatalf("expected %s, got %s", res.SessionToken, u4Sess.SessionToken)
	}

	if u4Sess.CSRFToken != res.CSRFToken {
		t.Fatalf("expected %s, got %s", res.CSRFToken, u4Sess.CSRFToken)
	}

	if u4Sess.MaxAge != res.MaxAge {
		t.Fatalf("expected %d, got %d", res.MaxAge, u4Sess.MaxAge)
	}
}

func TestUserDeleteSessionByUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u5Email := "test5@test.loc"
	u5 := usermodels.User{
		Email:    u5Email,
		Password: testPassHash,
	}

	id, err := userRepo.CreateUser(ctx, u5)
	if err != nil {
		t.Fatalf("user create error: %v", err)
	}

	var (
		sessToken = "sessToken5"
		csrfToken = "csrfToken5"
		maxAge    = 3600
		expireAt  = time.Now().Add(time.Hour * 1)
	)

	tmpSess := usermodels.Session{
		UserID:       id,
		SessionToken: sessToken,
		CSRFToken:    csrfToken,
		MaxAge:       maxAge,
		ExpireAt:     expireAt,
	}

	u5Sess, err := userRepo.CreateSession(ctx, tmpSess)
	if err != nil {
		t.Fatalf("err while creating session: %s", err)
	}

	id, err = userRepo.DeleteSessionByUserID(ctx, u5Sess.UserID)
	if err != nil {
		t.Fatalf("err while deleting session: %s", err)
	}

	if u5Sess.ID != id {
		t.Fatalf("expected %d, got %d", u5Sess.ID, id)
	}
}
