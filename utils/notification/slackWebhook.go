package notification

import (
	"fmt"
	"strings"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
)

func slackSend(w webhook) {
	slack.Send(w.WebhookUrl, "", w.Data.(*slack.Payload))
}

func TestSlackWebhook() {
	payload := &slack.Payload{
		Attachments: []*slack.Attachment{
			{
				Color:     "#00FF8B",
				Title:     "Successfully Checked Out!!!",
				TitleLink: "https://twitter.com/SICKO_AIO",
				Text:      fmt.Sprintf("*%s*\n%s", "Sicko Aio Test Product", "Product Description"),
				Fields: []*slack.Field{
					{
						Title: "Category",
						Value: "NIKE",
						Short: true,
					},
					{
						Title: "Region",
						Value: "US",
						Short: true,
					},
					{
						Title: "Product SKU",
						Value: "000000-000",
						Short: true,
					},
					{
						Title: "Size",
						Value: "10.5",
						Short: true,
					},
					{
						Title: "Quantity",
						Value: "1",
						Short: true,
					},

					{
						Title: "Order Number",
						Value: "C0000000000",
						Short: true,
					},
					{
						Title: "Profile",
						Value: "test",
						Short: true,
					},
					{
						Title: "Price",
						Value: "N/A",
						Short: true,
					},
					{
						Title: "Time",
						Value: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Short: false,
					},
					{
						Title: "Task Type",
						Value: "webhook:test",
						Short: true,
					},
				},
				Footer:       "Sicko AIO - 2.0",
				FooterIcon:   "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
				ThumbnailUrl: "https://secure-images.nike.com/is/image/DotCom/000000_000",
			},
		},
	}
	slackSend(webhook{
		WebhookUrl: communicator.Config.Settings.SlackWebhok[models.NIKE],
		Data:       payload,
	})
}

func sendSlackNikeACOCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	giftCards := []string{}
	for _, v := range worker.TaskInfo.GiftCardGroup {
		giftCards = append(giftCards, v.CardNumber)
	}
	payload := &slack.Payload{
		Attachments: []*slack.Attachment{
			{
				Title:     "Successfully Checked Out!!!",
				TitleLink: worker.TaskInfo.RedirectURL,
				Text:      fmt.Sprintf("*%s*\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*slack.Field{
					{
						Title: "Category",
						Value: string(taskGroupSetting.Category),
						Short: true,
					},
					{
						Title: "Region",
						Value: taskGroupSetting.Country,
						Short: true,
					},
					{
						Title: "Product SKU",
						Value: worker.Product.StyleColor,
						Short: true,
					},
					{
						Title: "Size",
						Value: worker.Product.Size,
						Short: true,
					},
					{
						Title: "Quantity",
						Value: fmt.Sprintf(`%v`, worker.Quantity),
						Short: true,
					},
					{
						Title: "Price",
						Value: func() string {
							if len(worker.Product.Price) == 0 {
								worker.Product.Price = "N/A"
							}
							return worker.Product.Price
						}(),
						Short: false,
					},
					{
						Title: "Order Number",
						Value: func() string {
							if len(worker.TaskInfo.OrderID) == 0 {
								worker.TaskInfo.OrderID = "N/A"
							}
							return worker.TaskInfo.OrderID
						}(),
						Short: true,
					},
					{
						Title: "Email",
						Value: func() string {
							if len(worker.TaskInfo.Email) == 0 {
								worker.TaskInfo.Email = "N/A"
							}
							return worker.TaskInfo.Email
						}(),
						Short: true,
					},
					{
						Title: "Profile",
						Value: func() string {
							switch {
							case worker.TaskInfo.Profile == nil:
								return "N/A"
							case len(worker.TaskInfo.Profile.ProfileName) == 0:
								worker.TaskInfo.Profile.ProfileName = "N/A"
							}
							return worker.TaskInfo.Profile.ProfileName
						}(),
						Short: true,
					},
					{
						Title: "Account",
						Value: func() string {
							if worker.TaskInfo.Account == nil {
								return "N/A"
							}
							return fmt.Sprintf(`%s:%s`, worker.TaskInfo.Account.Email, worker.TaskInfo.Account.Password)
						}(),
						Short: true,
					},
					{
						Title: "Discount Code",
						Value: func() string {
							if len(worker.TaskInfo.Discount) == 0 {
								worker.TaskInfo.Discount = "N/A"
							}
							return worker.TaskInfo.Discount
						}(),
						Short: true,
					},
					{
						Title: "Time",
						Value: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Short: false,
					},
					{
						Title: "Task Type",
						Value: taskGroupSetting.TaskType,
						Short: true,
					},
				},
				ThumbnailUrl: fmt.Sprintf("https://secure-images.nike.com/is/image/DotCom/%s", strings.Replace(worker.Product.StyleColor, "-", "_", 1)),
				Footer:       "Sicko AIO - 2.0",
				FooterIcon:   "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
			},
		},
	}
	logger.Info(fmt.Sprintf(`StyleColor: %s, OrderNumber: %s, Email: %s, Size: %v, URL: %s`, worker.Product.StyleColor, worker.TaskInfo.OrderID, worker.TaskInfo.Email, worker.Product.Size, worker.TaskInfo.RedirectURL))
	slackSend(webhook{
		WebhookUrl: communicator.Config.Settings.SlackWebhok[models.NIKE],
		Data:       payload,
	})
}

func sendSlackNikeLegacyCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	payload := &slack.Payload{
		Attachments: []*slack.Attachment{
			{
				Color:     "#ffb6c1",
				Title:     "Successfully Checked Out!!!",
				TitleLink: worker.TaskInfo.RedirectURL,
				Text:      fmt.Sprintf("*%s*\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*slack.Field{
					{
						Title: "Site",
						Value: fmt.Sprintf("NIKE_%s", taskGroupSetting.Country),
						Short: false,
					},
					{
						Title: "Product SKU",
						Value: worker.Product.StyleColor,
						Short: true,
					},
					{
						Title: "Size",
						Value: worker.Product.Size,
						Short: true,
					},
					{
						Title: "Quantity",
						Value: fmt.Sprintf(`%v`, worker.Quantity),
						Short: true,
					},
					// {
					// 	Title: "Profile",
					// 	Value: taskConfig.TaskInfo.Profile.ProfileName,
					// 	Short: true,
					// },
					// {
					// 	Title: "Account",
					// 	Value: fmt.Sprintf(`%s:%s`, taskConfig.TaskInfo.Account.Email, taskConfig.TaskInfo.Account.Password),
					// 	Short: true,
					// },
					// {
					// 	Title: "Discount Code",
					// 	Value: taskConfig.TaskInfo.Discount,
					// 	Short: true,
					// },
					{
						Title: "Task Type",
						Value: taskGroupSetting.TaskType,
						Short: true,
					},
				},
				ThumbnailUrl: fmt.Sprintf("https://secure-images.nike.com/is/image/DotCom/%s", strings.Replace(worker.Product.StyleColor, "-", "_", 1)),
				Footer:       "Sicko AIO - 2.0",
				FooterIcon:   "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
			},
		},
	}
	slackSend(webhook{
		WebhookUrl: communicator.Config.Settings.SlackWebhok[models.NIKE],
		Data:       payload,
	})
}

func sendSlackLuisaviaromaCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	payload := &slack.Payload{
		Attachments: []*slack.Attachment{
			{
				Color:     "#ffb6c1",
				Title:     "Successfully Checked Out!!!",
				TitleLink: worker.TaskInfo.RedirectURL,
				Text:      fmt.Sprintf("*%s*\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*slack.Field{

					{
						Title: "Category",
						Value: string(taskGroupSetting.Category),
						Short: true,
					},
					{
						Title: "Region",
						Value: taskGroupSetting.Country,
						Short: true,
					},
					{
						Title: "Product SKU",
						Value: worker.Product.StyleColor,
						Short: true,
					},
					{
						Title: "Size",
						Value: worker.Product.Size,
						Short: true,
					},
					{
						Title: "Quantity",
						Value: fmt.Sprintf(`%v`, worker.Quantity),
						Short: true,
					},
					{
						Title: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Short: true,
					},
					{
						Title: "Order Number",
						Value: func() string {
							if len(worker.TaskInfo.OrderID) == 0 {
								worker.TaskInfo.OrderID = "N/A"
							}
							return worker.TaskInfo.OrderID
						}(),
						Short: true,
					},
					{
						Title: "Email",
						Value: func() string {
							if len(worker.TaskInfo.Email) == 0 {
								worker.TaskInfo.Email = "N/A"
							}
							return worker.TaskInfo.Email
						}(),
						Short: true,
					},
					{
						Title: "Profile",
						Value: func() string {
							switch {
							case worker.TaskInfo.Profile == nil:
								return "N/A"
							case len(worker.TaskInfo.Profile.ProfileName) == 0:
								worker.TaskInfo.Profile.ProfileName = "N/A"
							}
							return worker.TaskInfo.Profile.ProfileName
						}(),
						Short: true,
					},
					{
						Title: "Account",
						Value: func() string {
							if worker.TaskInfo.Account == nil {
								return "N/A"
							}
							return fmt.Sprintf(`%s:%s`, worker.TaskInfo.Account.Email, worker.TaskInfo.Account.Password)
						}(),
						Short: true,
					},
					{
						Title: "Discount Code",
						Value: func() string {
							if len(worker.TaskInfo.Discount) == 0 {
								worker.TaskInfo.Discount = "N/A"
							}
							return worker.TaskInfo.Discount
						}(),
						Short: true,
					},
					{
						Title: "Time",
						Value: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Short: false,
					},
					{
						Title: "Task Type",
						Value: taskGroupSetting.TaskType,
						Short: true,
					},
				},
				ThumbnailUrl: fmt.Sprintf("https://cache.mrporter.com/variants/images/%s/in/w2000.jpg", worker.Product.StyleColor),
				Footer:       "Sicko AIO - 2.0",
				FooterIcon:   "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
			},
		},
	}
	slackSend(webhook{
		WebhookUrl: communicator.Config.Settings.SlackWebhok[models.MRPORTER],
		Data:       payload,
	})
}

func sendSlackMrPorterCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	payload := &slack.Payload{
		Attachments: []*slack.Attachment{
			{
				Color:     "#ffb6c1",
				Title:     "Successfully Checked Out!!!",
				TitleLink: worker.TaskInfo.RedirectURL,
				Text:      fmt.Sprintf("*%s*\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*slack.Field{

					{
						Title: "Category",
						Value: string(taskGroupSetting.Category),
						Short: true,
					},
					{
						Title: "Region",
						Value: taskGroupSetting.Country,
						Short: true,
					},
					{
						Title: "Product SKU",
						Value: worker.Product.StyleColor,
						Short: true,
					},
					{
						Title: "Size",
						Value: worker.Product.Size,
						Short: true,
					},
					{
						Title: "Quantity",
						Value: fmt.Sprintf(`%v`, worker.Quantity),
						Short: true,
					},
					{
						Title: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Short: true,
					},
					{
						Title: "Time",
						Value: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Short: false,
					},
					{
						Title: "Task Type",
						Value: taskGroupSetting.TaskType,
						Short: true,
					},
				},
				ThumbnailUrl: fmt.Sprintf("https://cache.mrporter.com/variants/images/%s/in/w2000.jpg", worker.Product.StyleColor),
				Footer:       "Sicko AIO - 2.0",
				FooterIcon:   "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
			},
		},
	}
	slackSend(webhook{
		WebhookUrl: communicator.Config.Settings.SlackWebhok[models.MRPORTER],
		Data:       payload,
	})
}

func sendSlackSsenseCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	payload := &slack.Payload{
		Attachments: []*slack.Attachment{
			{
				Color:     "#ffb6c1",
				Title:     "Successfully Checked Out!!!",
				TitleLink: worker.TaskInfo.RedirectURL,
				Text:      fmt.Sprintf("*%s*\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*slack.Field{

					{
						Title: "Category",
						Value: string(taskGroupSetting.Category),
						Short: true,
					},
					{
						Title: "Region",
						Value: taskGroupSetting.Country,
						Short: true,
					},
					{
						Title: "Product SKU",
						Value: worker.Product.StyleColor,
						Short: true,
					},
					{
						Title: "Size",
						Value: worker.Product.Size,
						Short: true,
					},
					{
						Title: "Quantity",
						Value: fmt.Sprintf(`%v`, worker.Quantity),
						Short: true,
					},
					{
						Title: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Short: true,
					},
					{
						Title: "Time",
						Value: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Short: false,
					},
					{
						Title: "Task Type",
						Value: taskGroupSetting.TaskType,
						Short: true,
					},
				},
				ThumbnailUrl: worker.Product.ImageURL,
				Footer:       "Sicko AIO - 2.0",
				FooterIcon:   "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
			},
		},
	}
	slackSend(webhook{
		WebhookUrl: communicator.Config.Settings.SlackWebhok[models.SSENSE],
		Data:       payload,
	})
}

func sendSlackPacsunCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	payload := &slack.Payload{
		Attachments: []*slack.Attachment{
			{
				Color:     "#ffb6c1",
				Title:     "Successfully Checked Out!!!",
				TitleLink: worker.TaskInfo.RedirectURL,
				Text:      fmt.Sprintf("*%s*\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*slack.Field{

					{
						Title: "Category",
						Value: string(taskGroupSetting.Category),
						Short: true,
					},
					{
						Title: "Region",
						Value: taskGroupSetting.Country,
						Short: true,
					},
					{
						Title: "Product SKU",
						Value: worker.Product.StyleColor,
						Short: true,
					},
					{
						Title: "Size",
						Value: worker.Product.Size,
						Short: true,
					},
					{
						Title: "Quantity",
						Value: fmt.Sprintf(`%v`, worker.Quantity),
						Short: true,
					},
					{
						Title: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Short: true,
					},
					{
						Title: "Time",
						Value: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Short: false,
					},
					{
						Title: "Task Type",
						Value: taskGroupSetting.TaskType,
						Short: true,
					},
				},
				ThumbnailUrl: worker.Product.ImageURL,
				Footer:       "Sicko AIO - 2.0",
				FooterIcon:   "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
			},
		},
	}
	slackSend(webhook{
		WebhookUrl: communicator.Config.Settings.SlackWebhok[models.PACSUN],
		Data:       payload,
	})
}

func sendSlackNewBalanceCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	payload := &slack.Payload{
		Attachments: []*slack.Attachment{
			{
				Color:     "#ffb6c1",
				Title:     "Successfully Checked Out!!!",
				TitleLink: "https://twitter.com/Sicko_AIO",
				Text:      fmt.Sprintf("*%s*\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*slack.Field{

					{
						Title: "Category",
						Value: string(taskGroupSetting.Category),
						Short: true,
					},
					{
						Title: "Region",
						Value: taskGroupSetting.Country,
						Short: true,
					},
					{
						Title: "Product SKU",
						Value: worker.Product.StyleColor,
						Short: true,
					},
					{
						Title: "Size",
						Value: worker.Product.Size,
						Short: true,
					},
					{
						Title: "Quantity",
						Value: fmt.Sprintf(`%v`, worker.Quantity),
						Short: true,
					},
					{
						Title: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Short: true,
					},
					{
						Title: "Time",
						Value: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Short: false,
					},
					{
						Title: "Task Type",
						Value: taskGroupSetting.TaskType,
						Short: true,
					},
				},
				ThumbnailUrl: worker.Product.ImageURL,
				Footer:       "Sicko AIO - 2.0",
				FooterIcon:   "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
			},
		},
	}
	slackSend(webhook{
		WebhookUrl: communicator.Config.Settings.SlackWebhok[models.NEWBALANCE],
		Data:       payload,
	})
}

func sendSlackTafCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	payload := &slack.Payload{
		Attachments: []*slack.Attachment{
			{
				Color:     "#ffb6c1",
				Title:     "Successfully Checked Out!!!",
				TitleLink: worker.TaskInfo.RedirectURL,
				Text:      fmt.Sprintf("*%s*\n%s", worker.Product.ProductName, worker.Product.ProductDescription),
				Fields: []*slack.Field{

					{
						Title: "Category",
						Value: string(taskGroupSetting.Category),
						Short: true,
					},
					{
						Title: "Region",
						Value: taskGroupSetting.Country,
						Short: true,
					},
					{
						Title: "Product SKU",
						Value: worker.Product.StyleColor,
						Short: true,
					},
					{
						Title: "Size",
						Value: worker.Product.Size,
						Short: true,
					},
					{
						Title: "Quantity",
						Value: fmt.Sprintf(`%v`, worker.Quantity),
						Short: true,
					},
					{
						Title: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Short: true,
					},
					{
						Title: "Order Number",
						Value: func() string {
							if len(worker.TaskInfo.OrderID) == 0 {
								worker.TaskInfo.OrderID = "N/A"
							}
							return worker.TaskInfo.OrderID
						}(),
						Short: true,
					},
					{
						Title: "Email",
						Value: func() string {
							if len(worker.TaskInfo.Email) == 0 {
								worker.TaskInfo.Email = "N/A"
							}
							return worker.TaskInfo.Email
						}(),
						Short: true,
					},
					{
						Title: "Profile",
						Value: func() string {
							switch {
							case worker.TaskInfo.Profile == nil:
								return "N/A"
							case len(worker.TaskInfo.Profile.ProfileName) == 0:
								worker.TaskInfo.Profile.ProfileName = "N/A"
							}
							return worker.TaskInfo.Profile.ProfileName
						}(),
						Short: true,
					},
					{
						Title: "Account",
						Value: func() string {
							if worker.TaskInfo.Account == nil {
								return "N/A"
							}
							return fmt.Sprintf(`%s:%s`, worker.TaskInfo.Account.Email, worker.TaskInfo.Account.Password)
						}(),
						Short: false,
					},
					{
						Title: "Time",
						Value: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Short: false,
					},
					{
						Title: "Task Type",
						Value: taskGroupSetting.TaskType,
						Short: true,
					},
				},
				ThumbnailUrl: worker.Product.ImageURL,
				Footer:       "Sicko AIO - 2.0",
				FooterIcon:   "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
			},
		},
	}
	slackSend(webhook{
		WebhookUrl: communicator.Config.Settings.SlackWebhok[models.TAF],
		Data:       payload,
	})
}

func sendSlackSneakerboyCheckoutWebhook(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	payload := &slack.Payload{
		Attachments: []*slack.Attachment{
			{
				Color:     "#ffb6c1",
				Title:     "Successfully Checked Out!!!",
				TitleLink: worker.TaskInfo.RedirectURL,
				Text:      fmt.Sprintf("*%s*", worker.Product.ProductName),
				Fields: []*slack.Field{

					{
						Title: "Category",
						Value: string(taskGroupSetting.Category),
						Short: true,
					},
					{
						Title: "Region",
						Value: taskGroupSetting.Country,
						Short: true,
					},
					{
						Title: "Product SKU",
						Value: worker.Product.StyleColor,
						Short: true,
					},
					{
						Title: "Size",
						Value: worker.Product.Size,
						Short: true,
					},
					{
						Title: "Quantity",
						Value: fmt.Sprintf(`%v`, worker.Quantity),
						Short: true,
					},
					{
						Title: "Price",
						Value: func() (price string) {
							price = worker.Product.Price
							if len(price) == 0 {
								price = "N/A"
							}
							return price
						}(),
						Short: true,
					},
					{
						Title: "Email",
						Value: worker.TaskInfo.Account.Email,
						Short: true,
					},
					{
						Title: "Password",
						Value: worker.TaskInfo.Account.Password,
						Short: true,
					},
					{
						Title: "Time",
						Value: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
						Short: false,
					},
					{
						Title: "Task Type",
						Value: taskGroupSetting.TaskType,
						Short: true,
					},
				},
				ThumbnailUrl: worker.Product.ImageURL,
				Footer:       "Sicko AIO - 2.0",
				FooterIcon:   "https://pbs.twimg.com/profile_images/1122681028210905088/2cZIhvv-_400x400.png",
			},
		},
	}
	slackSend(webhook{
		WebhookUrl: communicator.Config.Settings.SlackWebhok[models.SNEAKERBOY],
		Data:       payload,
	})
}
