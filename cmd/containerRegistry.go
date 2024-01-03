// Copyright © 2023 The Vultr-cli Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

var (
	crLong    = `Access information about container registries on the account and perform CRUD operations`
	crExample = `
	# Full example
	vultr-cli container-registry
	`
	crCreateLong    = `Create a new container registry with specified options`
	crCreateExample = `
	# Full example
	vultr-cli container-registry create --region="sjc" --name="my-registry" --public=true --plan="start_up"

	all flags are required

	# Shortened example with aliases
	vultr-cli cr c -i="sjc" -n="my-registry" -p=true -l="start_up"
	`

	crGetLong    = `Display information for a specific VPC`
	crGetExample = `
	# Full example
	vultr-cli container-registry get e8ba183d-df3b-487a-acbf-f6c06aa32468 

	# Shortened example with aliases
	vultr-cli cr g e8ba183d-df3b-487a-acbf-f6c06aa32468
	`
	crUpdateLong    = `Update an existing container registry`
	crUpdateExample = `
	# Full example
	vultr-cli container-registry update 835fd402-e0eb-47aa-a5a9-a9885feea1cf --plan="premium" --public="true" 

	# Shortened example with aliases
	vultr-cli cr u 835fd402-e0eb-47aa-a5a9-a9885feea1cf -p="premium" -b="true"
	`
	crDeleteLong    = `Delete a container registry`
	crDeleteExample = `
	#Full example
	vultr-cli container-registry delete b20fa61e-4abb-46c5-92c3-8700150e1f9a

	#Shortened example with aliases
	vultr-cli cr d b20fa61e-4abb-46c5-92c3-8700150e1f9a 
	`
	crListLong    = `List all container registries on the account`
	crListExample = `
	# Full example
	vultr-cli container-registry list

	# Shortened example with aliases
	vultr-cli cr l
	`
	crCredentialsLong = `Commands for accessing the credentials on registries`
	//nolint: gosec
	crCredentialsExample = `
	# Full example
	vultr-cli container-registry credentials
	`
	crCredentialsDockerLong = `Create the credential string used by docker`
	//nolint: gosec
	crCredentialsDockerExample = `
	# Full example
	vultr-cli container-registry credentials docker d24cfdcc-0534-4700-bf88-8ee48f20064e 
	`
	crRepoLong    = `Access commands for individual repositories on a container registry`
	crRepoExample = `
	# Full example
	vultr-cli container-registry repository

	# Shortened example with aliases
	vultr-cli cr r
	`
	crRepoUpdateLong    = `Update the details of registry's repository`
	crRepoUpdateExample = `
	# Full example
	vultr-cli container-registry repository update 4dcdc52e-9c63-401e-8c5f-1582490fe09c --image-name="my-thing" --description="new description"

	# Shortened example with aliases
	vultr-cli cr r u 4dcdc52e-9c63-401e-8c5f-1582490fe09c -i="my-thing" -d="new description"
	`
	crRepoDeleteLong    = `Delete a repository in a registry`
	crRepoDeleteExample = `
	# Full example
	vultr-cli container-registry repository delete 4dcdc52e-9c63-401e-8c5f-1582490fe09c --image-name="my-thing"

	# Shortened example with aliases
	vultr-cli cr r d 4dcdc52e-9c63-401e-8c5f-1582490fe09c -i="my-thing"
	`
	crPlansLong    = `Retrieve the current plan details for container registry`
	crPlansExample = `
	# Full example
	vultr-cli container-registry plans

	# Shortened example with aliases
	vultr-cli cr p
	`
	crRegionsLong    = `Retrieve the available regions for container registries`
	crRegionsExample = `
	# Full example
	vultr-cli container-registry regions

	# Shortened example with aliases
	vultr-cli cr r 
	`
)

// ContainerRegistry represents the container-registry command
func ContainerRegistry() *cobra.Command {
	crCmd := &cobra.Command{
		Use:     "container-registry",
		Aliases: []string{"cr"},
		Short:   "commands to interact with container registries",
		Long:    crLong,
		Example: crExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Context().Value(ctxAuthKey{}).(bool) {
				return errors.New(apiKeyError)
			}
			return nil
		},
	}

	crCmd.AddCommand(crCreate, crGet, crList, crDelete, crUpdate, crListPlans, crListRegions)
	crCreate.Flags().StringP("name", "n", "", "The name to use for the container registry")
	crCreate.Flags().StringP("region", "i", "", "The ID of the region in which to create the container registry")
	crCreate.Flags().BoolP("public", "p", false, "If the registry is publicly available. Should be true | false (default is false)")
	crCreate.Flags().StringP("plan", "l", "", "The type of plan to use for the container registry")
	if err := crCreate.MarkFlagRequired("name"); err != nil {
		fmt.Printf("error container registry create 'name' flag required: %v\n", err)
		os.Exit(1)
	}

	if err := crCreate.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error container registry create 'region' flag required: %v\n", err)
		os.Exit(1)
	}

	if err := crCreate.MarkFlagRequired("public"); err != nil {
		fmt.Printf("error container registry create 'public' flag required: %v\n", err)
		os.Exit(1)
	}

	if err := crCreate.MarkFlagRequired("plan"); err != nil {
		fmt.Printf("error container registry create 'plan' flag required: %v\n", err)
		os.Exit(1)
	}

	crList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	crList.Flags().IntP("per-page", "p", perPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	crUpdate.Flags().StringP("plan", "p", "", "Name of the plan used for the container registry")
	crUpdate.Flags().BoolP("public", "b", false, "The container registry availability status")

	crCredentialsCmd := &cobra.Command{
		Use:     "credentials",
		Aliases: []string{""},
		Short:   "Commands for container registry credentials",
		Long:    crCredentialsLong,
		Example: crCredentialsExample,
	}

	crCredentialsCmd.AddCommand(crCredentialsDocker)
	crCredentialsDocker.Flags().IntP("expiry-seconds", "e", 0, "(optional) The seconds until these credentials expire.  Default is 0, never")
	crCredentialsDocker.Flags().BoolP("read-write", "r", false, "(optional) Whether or not these credentials have write access.  Should be true or false.  Default is false") //nolint:lll

	crCmd.AddCommand(crCredentialsCmd)

	crRepoCmd := &cobra.Command{
		Use:     "repository",
		Aliases: []string{"r", "repo"},
		Short:   "Interact with container registry repositories",
		Long:    crRepoLong,
		Example: crRepoExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Context().Value(ctxAuthKey{}).(bool) {
				return errors.New(apiKeyError)
			}
			return nil
		},
	}

	crRepoCmd.AddCommand(crRepoGet, crRepoList, crRepoUpdate, crRepoDelete)

	crRepoGet.Flags().StringP("image-name", "i", "", "The name of the image/repo")
	if err := crRepoGet.MarkFlagRequired("image-name"); err != nil {
		fmt.Printf("error marking get container registry repository 'image-name' flag required: %v\n", err)
		os.Exit(1)
	}

	crRepoUpdate.Flags().StringP("image-name", "i", "", "The name of the image/repo")
	if err := crRepoUpdate.MarkFlagRequired("image-name"); err != nil {
		fmt.Printf("error marking update container registry repository 'image-name' flag required: %v\n", err)
		os.Exit(1)
	}

	crRepoUpdate.Flags().StringP("description", "d", "", "The description of the image/repo")

	crRepoDelete.Flags().StringP("image-name", "i", "", "The name of the image/repo")
	if err := crRepoDelete.MarkFlagRequired("image-name"); err != nil {
		fmt.Printf("error marking delete container registry repository 'image-name' flag required: %v\n", err)
		os.Exit(1)
	}

	crRepoList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	crRepoList.Flags().IntP("per-page", "p", perPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	crCmd.AddCommand(crRepoCmd)

	return crCmd
}

var crCreate = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "create a container registry",
	Long:    crCreateLong,
	Example: crCreateExample,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		region, _ := cmd.Flags().GetString("region")
		public, _ := cmd.Flags().GetBool("public")
		plan, _ := cmd.Flags().GetString("plan")

		options := &govultr.ContainerRegistryReq{
			Name:   name,
			Region: region,
			Public: public,
			Plan:   plan,
		}

		cr, _, err := client.ContainerRegistry.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("error creating container registry: %v\n", err)
			os.Exit(1)
		}

		printer.ContainerRegistry(cr)
	},
}

var crGet = &cobra.Command{
	Use:     "get <Registry ID>",
	Aliases: []string{"g"},
	Short:   "get a container registry",
	Long:    crGetLong,
	Example: crGetExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a container registry ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		cr, _, err := client.ContainerRegistry.Get(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting container registry : %v\n", err)
			os.Exit(1)
		}

		printer.ContainerRegistry(cr)
	},
}

var crList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "list all container registries",
	Long:    crListLong,
	Example: crListExample,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		cr, meta, _, err := client.ContainerRegistry.List(context.Background(), options)
		if err != nil {
			fmt.Printf("error getting container registry list : %v\n", err)
			os.Exit(1)
		}

		printer.ContainerRegistryList(cr, meta)
	},
}

var crUpdate = &cobra.Command{
	Use:     "update <Registry ID>",
	Aliases: []string{"u"},
	Short:   "update a container registry",
	Long:    crUpdateLong,
	Example: crUpdateExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a container registry ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		public, _ := cmd.Flags().GetBool("public")
		plan, _ := cmd.Flags().GetString("plan")

		options := &govultr.ContainerRegistryUpdateReq{
			Plan: govultr.StringToStringPtr(plan),
		}

		if cmd.Flags().Changed("public") {
			options.Public = govultr.BoolToBoolPtr(public)
		}

		cr, _, err := client.ContainerRegistry.Update(context.Background(), id, options)
		if err != nil {
			fmt.Printf("error updating container registry : %v\n", err)
			os.Exit(1)
		}

		printer.ContainerRegistry(cr)
	},
}

var crDelete = &cobra.Command{
	Use:     "delete <Registry ID>",
	Aliases: []string{"destroy", "d"},
	Short:   "delete a container registry",
	Long:    crDeleteLong,
	Example: crDeleteExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a container registry ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.ContainerRegistry.Delete(context.Background(), id); err != nil {
			fmt.Printf("error deleting container registry : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted container registry")
	},
}

var crListPlans = &cobra.Command{
	Use:     "plans",
	Aliases: []string{"p"},
	Short:   "list the plan names for container registry",
	Long:    crPlansLong,
	Example: crPlansExample,
	Run: func(cmd *cobra.Command, args []string) {
		plans, _, err := client.ContainerRegistry.ListPlans(context.Background())
		if err != nil {
			fmt.Printf("error getting container registry plans : %v\n", err)
			os.Exit(1)
		}

		printer.ContainerRegistryPlans(&plans.Plans)
	},
}

var crListRegions = &cobra.Command{
	Use:     "regions",
	Aliases: []string{"i"},
	Short:   "list the available regions for container registry",
	Long:    crRegionsLong,
	Example: crRegionsExample,
	Run: func(cmd *cobra.Command, args []string) {
		regions, _, _, err := client.ContainerRegistry.ListRegions(context.Background())
		if err != nil {
			fmt.Printf("error getting container registry regions : %v\n", err)
			os.Exit(1)
		}

		printer.ContainerRegistryRegions(regions)
	},
}

var crCredentialsDocker = &cobra.Command{
	Use:     "docker <Registry ID>",
	Aliases: []string{"d"},
	Short:   "create Docker credentials for a container registry",
	Long:    crCredentialsDockerLong,
	Example: crCredentialsDockerExample,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		expiry, _ := cmd.Flags().GetInt("expiry-seconds")
		access, _ := cmd.Flags().GetBool("write-access")

		options := &govultr.DockerCredentialsOpt{
			ExpirySeconds: govultr.IntToIntPtr(expiry),
			WriteAccess:   govultr.BoolToBoolPtr(access),
		}

		if access {
			options.WriteAccess = govultr.BoolToBoolPtr(access)
		}

		creds, _, err := client.ContainerRegistry.CreateDockerCredentials(context.Background(), id, options)
		if err != nil {
			fmt.Printf("error getting container registry docker credentials : %v\n", err)
			os.Exit(1)
		}

		printer.ContainerRegistryDockerCredentials(creds)
	},
}

var crRepoGet = &cobra.Command{
	Use:     "get <Registry ID>",
	Aliases: []string{"g"},
	Short:   "get a container registry repository",
	Long:    crGetLong,
	Example: crGetExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a container registry ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		name, _ := cmd.Flags().GetString("image-name")
		cr, _, err := client.ContainerRegistry.GetRepository(context.Background(), id, name)
		if err != nil {
			fmt.Printf("error getting container registry repository : %v\n", err)
			os.Exit(1)
		}

		printer.ContainerRegistryRepository(cr)
	},
}

var crRepoList = &cobra.Command{
	Use:     "list <Registry ID>",
	Aliases: []string{"l"},
	Short:   "list all container registries",
	Long:    crListLong,
	Example: crListExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a container registry ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		options := getPaging(cmd)
		cr, meta, _, err := client.ContainerRegistry.ListRepositories(context.Background(), id, options)
		if err != nil {
			fmt.Printf("error getting container registry repository list : %v\n", err)
			os.Exit(1)
		}

		printer.ContainerRegistryRepositoryList(cr, meta)
	},
}

var crRepoUpdate = &cobra.Command{
	Use:     "update <registry ID>",
	Aliases: []string{"u"},
	Short:   "update a container registry repository",
	Long:    crRepoUpdateLong,
	Example: crRepoUpdateExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a container registry ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		name, _ := cmd.Flags().GetString("image-name")
		description, _ := cmd.Flags().GetString("description")

		options := &govultr.ContainerRegistryRepoUpdateReq{
			Description: description,
		}

		cr, _, err := client.ContainerRegistry.UpdateRepository(context.Background(), id, name, options)
		if err != nil {
			fmt.Printf("error updating container registry repository : %v\n", err)
			os.Exit(1)
		}

		printer.ContainerRegistryRepository(cr)
	},
}

var crRepoDelete = &cobra.Command{
	Use:     "delete <registry ID>",
	Aliases: []string{"destroy", "d"},
	Short:   "delete a container registry repository",
	Long:    crRepoDeleteLong,
	Example: crRepoDeleteExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a container registry ID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		name, _ := cmd.Flags().GetString("image-name")
		if err := client.ContainerRegistry.DeleteRepository(context.Background(), id, name); err != nil {
			fmt.Printf("error deleting container registry repository : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted container registry repository")
	},
}
