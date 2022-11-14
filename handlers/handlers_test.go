package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ivanstepachev/tg_refferal/utils"
	"github.com/stretchr/testify/require"
)

func TestStartMessageLanding(t *testing.T) {
	tgID := utils.RandomTgID()
	user, err := startMessage(&testH, tgID, "/start landing")
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, tgID, user.TelegramId)
	require.Equal(t, "landing", user.BeneficiaryId)
	require.Equal(t, true, user.IsPaid)

	require.NotZero(t, user.Id)
	require.NotZero(t, user.CreatedAt)

	// _ = testH.DB.Delete(&user)
}

func TestStartMessageDirectly(t *testing.T) {
	tgID := utils.RandomTgID()
	user, err := startMessage(&testH, tgID, "/start")
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, tgID, user.TelegramId)
	require.Equal(t, "directly", user.BeneficiaryId)
	require.Equal(t, false, user.IsPaid)

	require.NotZero(t, user.Id)
	require.NotZero(t, user.CreatedAt)

	// _ = testH.DB.Delete(&user)
}

func TestStartMessageRefLink(t *testing.T) {
	tgID := utils.RandomTgID()
	startMes := fmt.Sprintf("/start %v", utils.RandomTgID())
	user, err := startMessage(&testH, tgID, startMes)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, tgID, user.TelegramId)
	require.Equal(t, startMes[7:], user.BeneficiaryId)
	require.Equal(t, false, user.IsPaid)

	require.NotZero(t, user.Id)
	require.NotZero(t, user.CreatedAt)

	// _ = testH.DB.Delete(&user)
}

type PaymentResponse struct {
	PaymentId string `json:"payment_id"`
	Status string `json:"status"`
}

func TestRefferalSystem(t *testing.T) {
	tgID1 := utils.RandomTgID()
	tgID2 := utils.RandomTgID()
	tgID3 := utils.RandomTgID()

	startMes2 := fmt.Sprintf("/start %v", tgID1)
	startMes3 := fmt.Sprintf("/start %v", tgID2)

	_, err := startMessage(&testH, tgID1, utils.RandomStartMessage())
	require.NoError(t, err)
	_, err = startMessage(&testH, tgID2, startMes2)
	require.NoError(t, err)
	_, err = startMessage(&testH, tgID3, startMes3)
	require.NoError(t, err)

	payment := PaymentResponse{
		PaymentId: tgID3,
		Status: "success",
	}
	jsonData, _ := json.Marshal(payment)

	req, err := http.NewRequest("POST", "/payment", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testH.PaymentApiHandler)
	handler.ServeHTTP(rr, req)
	
	user1, err := testH.selectUser(tgID1)
	require.NoError(t, err)
	require.Equal(t, 50, user1.Balance)
	user2, err := testH.selectUser(tgID2)
	require.NoError(t, err)
	require.Equal(t, 250, user2.Balance)

	user3, err := testH.selectUser(tgID3)
	require.NoError(t, err)
	require.Equal(t, true, user3.IsPaid)
}