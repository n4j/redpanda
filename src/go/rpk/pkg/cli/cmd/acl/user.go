// Copyright 2021 Vectorized, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package acl

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vectorizedio/redpanda/src/go/rpk/pkg/api/admin"
	"github.com/vectorizedio/redpanda/src/go/rpk/pkg/config"
)

const (
	newUserFlag     = "new-username"
	newPasswordFlag = "new-password"

	deleteUsernameFlag = "delete-username"
)

func NewUserCommand(tls func() (*config.TLS, error)) *cobra.Command {
	var apiUrl string

	command := &cobra.Command{
		Use:          "user",
		Short:        "Manage users",
		SilenceUsage: true,
	}
	command.PersistentFlags().StringVar(
		&apiUrl,
		"api-url",
		fmt.Sprintf("localhost:%d", config.DefaultAdminPort),
		"The Admin API URL",
	)

	adminApi := buildAdminAPI(&apiUrl, tls)

	command.AddCommand(NewCreateUserCommand(adminApi))
	command.AddCommand(NewDeleteUserCommand(adminApi))
	return command
}

func NewCreateUserCommand(
	adminApi func() (admin.AdminAPI, error),
) *cobra.Command {
	var (
		newUser     string
		newPassword string
	)
	command := &cobra.Command{
		Use:          "create",
		Short:        "Create users",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			adminApi, err := adminApi()
			if err != nil {
				return err
			}
			err = adminApi.CreateUser(newUser, newPassword)
			if err != nil {
				return err
			}

			log.Infof("Created user '%s'", newUser)

			return nil
		},
	}

	command.Flags().StringVar(
		&newUser,
		newUserFlag,
		"",
		"The user to be created",
	)
	command.MarkFlagRequired(newUserFlag)
	command.Flags().StringVar(
		&newPassword,
		newPasswordFlag,
		"",
		"The new user's password",
	)
	command.MarkFlagRequired(newPasswordFlag)

	return command
}

func NewDeleteUserCommand(
	adminApi func() (admin.AdminAPI, error),
) *cobra.Command {
	var (
		username string
	)
	command := &cobra.Command{
		Use:          "delete",
		Short:        "Delete users",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			adminApi, err := adminApi()
			if err != nil {
				return err
			}
			err = adminApi.DeleteUser(username)
			if err != nil {
				return err
			}

			log.Infof("Deleted user '%s'", username)

			return nil
		},
	}

	command.Flags().StringVar(
		&username,
		deleteUsernameFlag,
		"",
		"The user to be deleted",
	)
	command.MarkFlagRequired(deleteUsernameFlag)

	return command
}

func buildAdminAPI(
	apiUrl *string, tls func() (*config.TLS, error),
) func() (admin.AdminAPI, error) {
	return func() (admin.AdminAPI, error) {
		tlsConfig, err := tls()
		if err != nil {
			return nil, err
		}
		return admin.NewAdminAPI(*apiUrl, tlsConfig)
	}
}
