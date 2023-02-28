package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/frain-dev/convoy/server/models"
	"github.com/frain-dev/convoy/services"

	"github.com/frain-dev/convoy/config"
	"github.com/frain-dev/convoy/database/postgres"
	"github.com/frain-dev/convoy/datastore"
	"github.com/frain-dev/convoy/internal/pkg/migrator"
	"github.com/frain-dev/convoy/pkg/log"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
)

func addMigrateCommand(a *app) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Convoy migrations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			err = config.LoadConfig(cfgPath)
			if err != nil {
				return err
			}

			_, err = config.Get()
			if err != nil {
				return err
			}

			// Override with CLI Flags
			cliConfig, err := buildCliConfiguration(cmd)
			if err != nil {
				return err
			}

			if err = config.Override(cliConfig); err != nil {
				return err
			}

			return nil

		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(addUpCommand())
	cmd.AddCommand(addDownCommand())
    cmd.AddCommand(addRunCommand())
	return cmd
}

func addUpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "up",
		Aliases: []string{"migrate-up"},
		Short:   "Run all pending migrations",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Get()
			if err != nil {
				log.WithError(err).Fatalf("Error fetching the config.")
			}

			db, err := postgres.NewDB(cfg)
			if err != nil {
				log.Fatal(err)
			}

			defer db.GetDB().Close()

			m := migrator.New(db)
			err = m.Up()
			if err != nil {
				log.Fatalf("migration up failed with error: %+v", err)
			}
		},
	}

	return cmd
}

func addDownCommand() *cobra.Command {
	var migrationID string

	cmd := &cobra.Command{
		Use:     "down",
		Aliases: []string{"migrate-down"},
		Short:   "Rollback migrations",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Get()
			if err != nil {
				log.WithError(err).Fatalf("Error fetching the config.")
			}

			db, err := postgres.NewDB(cfg)
			if err != nil {
				log.Fatal(err)
			}

			defer db.GetDB().Close()

			m := migrator.New(db)
			err = m.Down()
			if err != nil {
				log.Fatalf("migration down failed with error: %+v", err)
			}
		},
	}

	cmd.Flags().StringVar(&migrationID, "id", "", "Migration ID")

	return cmd
}

func addRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"migrate-run"},
		Short:   "Run arbitrary SQL queries",
		Run: func(cmd *cobra.Command, args []string) {
			db, err := postgres.NewDB(config.Configuration{Database: config.DatabaseConfiguration{Dsn: "postgres://admin:password@localhost:5432/convoy?sslmode=disable"}})
			if err != nil {
				log.Fatal(err)
			}

			psw := datastore.Password{
				Plaintext: "12345678",
				Hash:      nil,
			}

			err = psw.GenerateHash()
			if err != nil {
				log.Fatal(err)
			}

			ur := postgres.NewUserRepo(db)
			user := &datastore.User{
				UID:                    ulid.Make().String(),
				FirstName:              "Daniel",
				LastName:               "O.J",
				Email:                  "danvixent@gmail.com",
				EmailVerified:          true,
				Password:               string(psw.Hash),
				ResetPasswordToken:     "bfuyudy",
				EmailVerificationToken: "vvfedfef",
				CreatedAt:              time.Now(),
				UpdatedAt:              time.Now(),
			}

			err = ur.CreateUser(cmd.Context(), user)
			if err != nil {
				log.Fatal("create user ", err)
			}

			orgserv := services.NewOrganisationService(
				postgres.NewOrgRepo(db),
				postgres.NewOrgMemberRepo(db),
			)

			org, err := orgserv.CreateOrganisation(context.Background(), &models.Organisation{
				Name: "big-org-name",
			}, user)
			if err != nil {
				log.Fatal("create org", err)
			}
			// orgs, _, err := o.LoadOrganisationsPaged(cmd.Context(), datastore.Pageable{
			// 	Page:    1,
			// 	PerPage: 10,
			// })

			// if err != nil {
			// 	return
			// }

			p := postgres.NewProjectRepo(db)
			proj := &datastore.Project{
				UID:            ulid.Make().String(),
				Name:           "mob psycho",
				Type:           datastore.OutgoingProject,
				OrganisationID: org.UID,
				Config:         &datastore.DefaultProjectConfig,
			}

			err = p.CreateProject(cmd.Context(), proj)
			if err != nil {
				fmt.Printf("CreateProject: %+v", err)
				return
			}

			endpoint := &datastore.Endpoint{
				UID:       ulid.Make().String(),
				ProjectID: proj.UID,
				OwnerID:   "owner1",
				TargetURL: "http://localhost",
				Title:     "test_endpoint",
				Secrets: []datastore.Secret{
					{
						UID:       ulid.Make().String(),
						Value:     "secret1",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				},
				Description:       "testing",
				HttpTimeout:       "10s",
				RateLimit:         100,
				Status:            datastore.ActiveEndpointStatus,
				RateLimitDuration: "3s",
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}

			endpointRepo := postgres.NewEndpointRepo(db)
			err = endpointRepo.CreateEndpoint(context.Background(), endpoint, proj.UID)
			if err != nil {
				fmt.Printf("CreateEndpoint: %+v", err)
				return
			}

			s := postgres.NewSubscriptionRepo(db)
			sub := &datastore.Subscription{
				UID:         ulid.Make().String(),
				Name:        "test_sub",
				Type:        datastore.SubscriptionTypeAPI,
				ProjectID:   proj.UID,
				EndpointID:  endpoint.UID,
				AlertConfig: &datastore.DefaultAlertConfig,
				RetryConfig: &datastore.DefaultRetryConfig,
				FilterConfig: &datastore.FilterConfiguration{
					EventTypes: []string{"*"},
					Filter: datastore.FilterSchema{
						Headers: datastore.M{},
						Body:    datastore.M{},
					},
				},
				RateLimitConfig: &datastore.DefaultRateLimitConfig,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}

			err = s.CreateSubscription(context.Background(), proj.UID, sub)
			if err != nil {
				fmt.Printf("CreateSubscription: %+v", err)
				return
			}

			e := postgres.NewEventRepo(db)
			event := &datastore.Event{
				UID:       ulid.Make().String(),
				EventType: "*",
				ProjectID: proj.UID,
				Endpoints: []string{endpoint.UID},
				Data:      []byte(`{"ref":"terf"}`),
				Raw:       `{"ref":"terf"}`,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err = e.CreateEvent(context.Background(), event)
			if err != nil {
				fmt.Printf("CreateEvent: %+v", err)
				return
			}

			now := time.Now()
			edRepo := postgres.NewEventDeliveryRepo(db)

			for i := 0; i < 500; i++ {

				n := rand.Intn(200)
				now = now.Add(-time.Hour * 24)

				for j := 0; j < n; j++ {

					delivery := &datastore.EventDelivery{
						UID:            ulid.Make().String(),
						ProjectID:      proj.UID,
						EventID:        event.UID,
						EndpointID:     endpoint.UID,
						SubscriptionID: sub.UID,
						Status:         datastore.SuccessEventStatus,
						Metadata: &datastore.Metadata{
							Data:            event.Data,
							Raw:             event.Raw,
							Strategy:        sub.RetryConfig.Type,
							NextSendTime:    time.Now(),
							NumTrials:       1,
							IntervalSeconds: 2,
							RetryLimit:      2,
						},
						Description: "ccc",
						CreatedAt:   now,
						UpdatedAt:   now,
					}

					fmt.Println("now", now.Format(time.RFC3339))

					err = edRepo.CreateEventDelivery(context.Background(), delivery)
					if err != nil {
						fmt.Printf("CreateEventDelivery: %+v", err)
						return
					}
				}
			}

			params := datastore.SearchParams{
				CreatedAtStart: time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC).Unix(),
				CreatedAtEnd:   time.Now().Add(time.Hour).Unix(),
			}

			intervals, err := edRepo.LoadEventDeliveriesIntervals(context.Background(), "01GSJ4033W1V3Y4SREXTYREXRY", params, datastore.Yearly, 1)
			if err != nil {
				fmt.Printf("LoadEventDeliveriesIntervals: %+v", err)
				return
			}

			fmt.Printf("intervals %+v", intervals)
		},
	}

	return cmd
}
