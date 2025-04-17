package handler

import (
	"net/http"

	"github.com/gogf/gf/net/ghttp"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks/handler/mrporter"
)

type syncAccountPayload struct {
	Category         models.CATEGORY         `json:"category"`
	AccountGroupName models.AccountGroupName `json:"accountGroupName"`
	Email            string                  `json:"email"`
	Method           string                  `json:"method"`
}

func SyncAccount(r *ghttp.Request) {

	var payload []*syncAccountPayload
	if err := r.Parse(&payload); err != nil {
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteJsonExit(map[string]bool{"success": false})
		return
	}

	for _, v := range payload {

		go func() {
			var index int
			var err *models.Error
			var oldAccount *models.Account
			var syncedAccount *models.Account

			for k, account := range communicator.Config.Accounts[v.AccountGroupName] {
				if v.Email == account.Email {
					index = k
					oldAccount = account
					break
				}
			}

			switch v.Category {
			case models.NIKE:

			case models.MRPORTER:

				syncedAccount, err = mrporter.Sync(oldAccount, "")

			case models.SSENSE:

			}

			if syncedAccount != nil && err == nil {
				communicator.Config.Accounts[v.AccountGroupName][index] = syncedAccount
			} else {
				oldAccount.Status = "error"
				communicator.Config.Accounts[v.AccountGroupName][index] = oldAccount
			}
		}()
	}
	r.Response.WriteJsonExit(map[string]bool{"success": true})
}
