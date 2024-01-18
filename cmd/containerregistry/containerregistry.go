// Package containerregistry provides functionality for the CLI to control
// container registries
package containerregistry

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	long    = `Access information about container registries on the account and perform CRUD operations`
	example = `
	# Full example
	vultr-cli container-registry
	`
	createLong    = `Create a new container registry with specified options`
	createExample = `
	# Full example
	vultr-cli container-registry create --region="sjc" --name="my-registry" --public=true --plan="start_up"

	all flags are required

	# Shortened example with aliases
	vultr-cli cr c -i="sjc" -n="my-registry" -p=true -l="start_up"
	`

	getLong    = `Display information for a specific VPC`
	getExample = `
	# Full example
	vultr-cli container-registry get e8ba183d-df3b-487a-acbf-f6c06aa32468 

	# Shortened example with aliases
	vultr-cli cr g e8ba183d-df3b-487a-acbf-f6c06aa32468
	`
	updateLong    = `Update an existing container registry`
	updateExample = `
	# Full example
	vultr-cli container-registry update 835fd402-e0eb-47aa-a5a9-a9885feea1cf --plan="premium" --public="true" 

	# Shortened example with aliases
	vultr-cli cr u 835fd402-e0eb-47aa-a5a9-a9885feea1cf -p="premium" -b="true"
	`
	deleteLong    = `Delete a container registry`
	deleteExample = `
	#Full example
	vultr-cli container-registry delete b20fa61e-4abb-46c5-92c3-8700150e1f9a

	#Shortened example with aliases
	vultr-cli cr d b20fa61e-4abb-46c5-92c3-8700150e1f9a 
	`
	listLong    = `List all container registries on the account`
	listExample = `
	# Full example
	vultr-cli container-registry list

	# Shortened example with aliases
	vultr-cli cr l
	`
	credentialsLong    = `Commands for accessing the credentials on registries`
	credentialsExample = `
	# Full example
	vultr-cli container-registry credentials
	`
	credentialsDockerLong    = `Create the credential string used by docker`
	credentialsDockerExample = `
	# Full example
	vultr-cli container-registry credentials docker d24cfdcc-0534-4700-bf88-8ee48f20064e 
	`
	repoLong    = `Access commands for individual repositories on a container registry`
	repoExample = `
	# Full example
	vultr-cli container-registry repository

	# Shortened example with aliases
	vultr-cli cr r
	`
	repoUpdateLong    = `Update the details of registry's repository`
	repoUpdateExample = `
	# Full example
	vultr-cli container-registry repository update 4dcdc52e-9c63-401e-8c5f-1582490fe09c --image-name="my-thing" --description="new description"

	# Shortened example with aliases
	vultr-cli cr r u 4dcdc52e-9c63-401e-8c5f-1582490fe09c -i="my-thing" -d="new description"
	`
	repoDeleteLong    = `Delete a repository in a registry`
	repoDeleteExample = `
	# Full example
	vultr-cli container-registry repository delete 4dcdc52e-9c63-401e-8c5f-1582490fe09c --image-name="my-thing"

	# Shortened example with aliases
	vultr-cli cr r d 4dcdc52e-9c63-401e-8c5f-1582490fe09c -i="my-thing"
	`
	plansLong    = `Retrieve the current plan details for container registry`
	plansExample = `
	# Full example
	vultr-cli container-registry plans

	# Shortened example with aliases
	vultr-cli cr p
	`
	regionsLong    = `Retrieve the available regions for container registries`
	regionsExample = `
	# Full example
	vultr-cli container-registry regions

	# Shortened example with aliases
	vultr-cli cr r 
	`
)

type Options struct {
	Base                 *cli.Base
	CreateReq            *govultr.ContainerRegistryReq
	UpdateReq            *govultr.ContainerRegistryUpdateReq
	RepoName             string
	RepoUpdateReq        *govultr.ContainerRegistryRepoUpdateReq
	CredentialsDockerReq *govultr.DockerCredentialsOpt
}

// NewCmdContainerRegistry provides the CLI command functionality for container registry
func NewCmdContainerRegistry(base *cli.Base) *cobra.Command {
	o := &Options{Base: base}

	cmd := &cobra.Command{
		Use:     "container-registry",
		Aliases: []string{"cr"},
		Short:   "commands to interact with container registries",
		Long:    long,
		Example: example,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}
			return nil
		},
	}

	// List
	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "list all container registries",
		Long:    listLong,
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Base.Options = utils.GetPaging(cmd)
			regs, meta, err := o.List()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving container registry list : %v", err))
				os.Exit(1)
			}

			data := &ContainerRegistriesPrinter{Registries: regs, Meta: meta}
			o.Base.Printer.Display(data, nil)
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	list.Flags().IntP("per-page", "p", utils.PerPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	// Get
	get := &cobra.Command{
		Use:     "get <Registry ID>",
		Aliases: []string{"g"},
		Short:   "get a container registry",
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			reg, err := o.Get()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving container registry info : %v", err))
				os.Exit(1)
			}

			data := &ContainerRegistryPrinter{Registry: reg}
			o.Base.Printer.Display(data, nil)
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "create a container registry",
		Long:    createLong,
		Example: createExample,
		Run: func(cmd *cobra.Command, args []string) {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				printer.Error(fmt.Errorf("error parsing 'name' flag for container registry create : %v", err))
				os.Exit(1)
			}

			region, err := cmd.Flags().GetString("region")
			if err != nil {
				printer.Error(fmt.Errorf("error parsing 'region' flag for container registry create : %v", err))
				os.Exit(1)
			}

			public, err := cmd.Flags().GetBool("public")
			if err != nil {
				printer.Error(fmt.Errorf("error parsing 'public' flag for container registry create : %v", err))
				os.Exit(1)
			}

			plan, err := cmd.Flags().GetString("plan")
			if err != nil {
				printer.Error(fmt.Errorf("error parsing 'plan' flag for container registry create : %v", err))
				os.Exit(1)
			}

			o.CreateReq = &govultr.ContainerRegistryReq{
				Name:   name,
				Region: region,
				Public: public,
				Plan:   plan,
			}

			reg, err := o.Create()
			if err != nil {
				printer.Error(fmt.Errorf("error creating container registry : %v", err))
				os.Exit(1)
			}

			data := &ContainerRegistryPrinter{Registry: reg}
			o.Base.Printer.Display(data, nil)
		},
	}

	create.Flags().StringP("name", "n", "", "The name to use for the container registry")
	create.Flags().StringP("region", "i", "", "The ID of the region in which to create the container registry")
	create.Flags().BoolP("public", "p", false, "If the registry is publicly available. Should be true | false (default is false)")
	create.Flags().StringP("plan", "l", "", "The type of plan to use for the container registry")
	if err := create.MarkFlagRequired("name"); err != nil {
		printer.Error(fmt.Errorf("error marking container registry create 'name' flag required: %v\n", err))
		os.Exit(1)
	}

	if err := create.MarkFlagRequired("region"); err != nil {
		printer.Error(fmt.Errorf("error marking container registry create 'region' flag required: %v\n", err))
		os.Exit(1)
	}

	if err := create.MarkFlagRequired("public"); err != nil {
		printer.Error(fmt.Errorf("error marking container registry create 'public' flag required: %v\n", err))
		os.Exit(1)
	}

	if err := create.MarkFlagRequired("plan"); err != nil {
		printer.Error(fmt.Errorf("error marking container registry create 'plan' flag required: %v\n", err))
		os.Exit(1)
	}

	// Update
	update := &cobra.Command{
		Use:     "update <Registry ID>",
		Aliases: []string{"u"},
		Short:   "update a container registry",
		Long:    updateLong,
		Example: updateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			public, errPu := cmd.Flags().GetBool("public")
			if errPu != nil {
				printer.Error(fmt.Errorf("error parsing 'public' flag for container registry update : %v", errPu))
				os.Exit(1)
			}

			plan, errPl := cmd.Flags().GetString("plan")
			if errPl != nil {
				printer.Error(fmt.Errorf("error parsing 'plan' flag for container registry update : %v", errPl))
				os.Exit(1)
			}

			o.UpdateReq = &govultr.ContainerRegistryUpdateReq{
				Plan: govultr.StringToStringPtr(plan),
			}

			if cmd.Flags().Changed("public") {
				o.UpdateReq.Public = govultr.BoolToBoolPtr(public)
			}

			if err := o.Update(); err != nil {
				printer.Error(fmt.Errorf("error updating container registry : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("container registry has been updated"), nil)
		},
	}

	update.Flags().StringP("plan", "p", "", "Name of the plan used for the container registry")
	update.Flags().BoolP("public", "b", false, "The container registry availability status")

	// Delete
	del := &cobra.Command{
		Use:     "delete <Registry ID>",
		Aliases: []string{"destroy", "d"},
		Short:   "delete a container registry",
		Long:    deleteLong,
		Example: deleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.Delete(); err != nil {
				printer.Error(fmt.Errorf("error deleting container registry : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("container registry has been deleted"), nil)
		},
	}

	// Plans
	plans := &cobra.Command{
		Use:     "plans",
		Aliases: []string{"p"},
		Short:   "list the plan names for container registry",
		Long:    plansLong,
		Example: plansExample,
		Run: func(cmd *cobra.Command, args []string) {
			plans, err := o.Plans()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving plans for container registry : %v", err))
				os.Exit(1)
			}

			data := &ContainerRegistryPlansPrinter{Plans: plans.Plans}
			o.Base.Printer.Display(data, nil)
		},
	}

	// Regions
	regions := &cobra.Command{
		Use:     "regions",
		Aliases: []string{"i"},
		Short:   "list the available regions for container registry",
		Long:    regionsLong,
		Example: regionsExample,
		Run: func(cmd *cobra.Command, args []string) {
			regs, meta, err := o.Regions()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving regions for container registry : %v", err))
				os.Exit(1)
			}

			data := &ContainerRegistryRegionsPrinter{Regions: regs, Meta: meta}
			o.Base.Printer.Display(data, nil)
		},
	}

	// Repository
	repository := &cobra.Command{
		Use:     "repository",
		Aliases: []string{"r", "repo"},
		Short:   "interact with container registry repositories",
		Long:    repoLong,
		Example: repoExample,
	}

	// Repository List
	repoList := &cobra.Command{
		Use:     "list <Registry ID>",
		Aliases: []string{"l"},
		Short:   "list all container registries",
		Long:    listLong,
		Example: listExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Base.Options = utils.GetPaging(cmd)
			repos, meta, err := o.RepositoryList()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving repositories for container registry : %v", err))
				os.Exit(1)
			}

			data := &ContainerRegistryRepositoriesPrinter{Repositories: repos, Meta: meta}
			o.Base.Printer.Display(data, nil)
		},
	}

	// Repository Get
	repoGet := &cobra.Command{
		Use:     "get <Registry ID>",
		Aliases: []string{"g"},
		Short:   "get a container registry repository",
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name, errIm := cmd.Flags().GetString("image-name")
			if errIm != nil {
				printer.Error(fmt.Errorf("error parsing 'image-name' flag for container registry repository get : %v", errIm))
				os.Exit(1)
			}

			o.RepoName = name

			repo, err := o.RepositoryGet()
			if err != nil {
				printer.Error(fmt.Errorf("error getting container registry repository : %v", err))
				os.Exit(1)
			}

			data := &ContainerRegistryRepositoryPrinter{Repository: repo}
			o.Base.Printer.Display(data, nil)
		},
	}

	repoGet.Flags().StringP("image-name", "i", "", "The name of the image/repo")
	if err := repoGet.MarkFlagRequired("image-name"); err != nil {
		printer.Error(fmt.Errorf("error marking get container registry repository 'image-name' flag required: %v\n", err))
		os.Exit(1)
	}

	// Repository Update
	repoUpdate := &cobra.Command{
		Use:     "update <registry ID>",
		Aliases: []string{"u"},
		Short:   "update a container registry repository",
		Long:    repoUpdateLong,
		Example: repoUpdateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name, errIm := cmd.Flags().GetString("image-name")
			if errIm != nil {
				printer.Error(fmt.Errorf("error parsing 'image-name' flag for container registry repository update : %v", errIm))
				os.Exit(1)
			}

			description, errDe := cmd.Flags().GetString("description")
			if errDe != nil {
				printer.Error(fmt.Errorf("error parsing 'description' flag for container registry repository update : %v", errDe))
				os.Exit(1)
			}

			o.RepoName = name
			o.RepoUpdateReq = &govultr.ContainerRegistryRepoUpdateReq{
				Description: description,
			}

			if err := o.RepositoryUpdate(); err != nil {
				printer.Error(fmt.Errorf("error updating container registry repository : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("container registry repository has been updated"), nil)
		},
	}

	repoUpdate.Flags().StringP("image-name", "i", "", "The name of the image/repo")
	if err := repoUpdate.MarkFlagRequired("image-name"); err != nil {
		printer.Error(fmt.Errorf("error marking update container registry repository 'image-name' flag required: %v\n", err))
		os.Exit(1)
	}

	repoUpdate.Flags().StringP("description", "d", "", "The description of the image/repo")
	if err := repoUpdate.MarkFlagRequired("description"); err != nil {
		printer.Error(fmt.Errorf("error marking update container registry repository 'description' flag required: %v\n", err))
		os.Exit(1)
	}

	// Repository Delete
	repoDelete := &cobra.Command{
		Use:     "delete <registry ID>",
		Aliases: []string{"destroy", "d"},
		Short:   "delete a container registry repository",
		Long:    repoDeleteLong,
		Example: repoDeleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name, errIm := cmd.Flags().GetString("image-name")
			if errIm != nil {
				printer.Error(fmt.Errorf("error parsing 'image-name' flag for container registry repository delete : %v", errIm))
				os.Exit(1)
			}

			o.RepoName = name

			if err := o.RepositoryDelete(); err != nil {
				printer.Error(fmt.Errorf("error deleting container registry repository : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("container registry repository has been deleted"), nil)
		},
	}

	repoDelete.Flags().StringP("image-name", "i", "", "The name of the image/repo")
	if err := repoDelete.MarkFlagRequired("image-name"); err != nil {
		printer.Error(fmt.Errorf("error marking delete container registry repository 'image-name' flag required: %v\n", err))
		os.Exit(1)
	}

	repository.AddCommand(repoGet, repoList, repoUpdate, repoDelete)

	// Credentials
	credentials := &cobra.Command{
		Use:     "credentials",
		Aliases: []string{""},
		Short:   "Commands for container registry credentials",
		Long:    credentialsLong,
		Example: credentialsExample,
	}

	// Credentials Docker
	credentialsDocker := &cobra.Command{
		Use:     "docker <Registry ID>",
		Aliases: []string{"d"},
		Short:   "create Docker credentials for a container registry",
		Long:    credentialsDockerLong,
		Example: credentialsDockerExample,
		Run: func(cmd *cobra.Command, args []string) {
			expiry, errEx := cmd.Flags().GetInt("expiry-seconds")
			if errEx != nil {
				printer.Error(fmt.Errorf("error parsing 'expiry-seconds' flag for container registry docker creds : %v", errEx))
				os.Exit(1)
			}

			access, errAc := cmd.Flags().GetBool("write-access")
			if errAc != nil {
				printer.Error(fmt.Errorf("error parsing 'write-access' flag for container registry docker creds : %v", errAc))
				os.Exit(1)
			}

			o.CredentialsDockerReq = &govultr.DockerCredentialsOpt{
				ExpirySeconds: govultr.IntToIntPtr(expiry),
				WriteAccess:   govultr.BoolToBoolPtr(access),
			}

			cred, err := o.CredentialsDocker()
			if err != nil {
				printer.Error(fmt.Errorf("error generating container registry repository docker credentials : %v", err))
				os.Exit(1)
			}

			data := &ContainerRegistryCredentialDockerPrinter{Credential: cred}
			o.Base.Printer.Display(data, nil)
		},
	}

	credentialsDocker.Flags().IntP(
		"expiry-seconds", 
		"e", 
		0, 
		"(optional) The seconds until these credentials expire.  Default is 0, never",
	)

	credentialsDocker.Flags().BoolP(
		"write-access", 
		"w", 
		false, 
		"(optional) Whether or not these credentials have write access.  Should be true or false.  Default is false",
	)

	credentials.AddCommand(credentialsDocker)

	cmd.AddCommand(
		list,
		get,
		create,
		update,
		del,
		plans,
		regions,
		repository,
		credentials,
	)

	return cmd
}

// List ...
func (o *Options) List() ([]govultr.ContainerRegistry, *govultr.Meta, error) {
	cr, meta, _, err := o.Base.Client.ContainerRegistry.List(o.Base.Context, o.Base.Options)
	return cr, meta, err
}

// Get ...
func (o *Options) Get() (*govultr.ContainerRegistry, error) {
	cr, _, err := o.Base.Client.ContainerRegistry.Get(o.Base.Context, o.Base.Args[0])
	return cr, err
}

// Create ...
func (o *Options) Create() (*govultr.ContainerRegistry, error) {
	cr, _, err := o.Base.Client.ContainerRegistry.Create(o.Base.Context, o.CreateReq)
	return cr, err
}

// Update ...
func (o *Options) Update() error {
	_, _, err := o.Base.Client.ContainerRegistry.Update(o.Base.Context, o.Base.Args[0], o.UpdateReq)
	return err
}

// Delete ...
func (o *Options) Delete() error {
	return o.Base.Client.ContainerRegistry.Delete(o.Base.Context, o.Base.Args[0])
}

// Plans list ...
func (o *Options) Plans() (*govultr.ContainerRegistryPlans, error) {
	plans, _, err := o.Base.Client.ContainerRegistry.ListPlans(o.Base.Context)
	return plans, err
}

// Regions list ...
func (o *Options) Regions() ([]govultr.ContainerRegistryRegion, *govultr.Meta, error) {
	regions, meta, _, err := o.Base.Client.ContainerRegistry.ListRegions(o.Base.Context)
	return regions, meta, err
}

// RepositoryList ...
func (o *Options) RepositoryList() ([]govultr.ContainerRegistryRepo, *govultr.Meta, error) {
	repos, meta, _, err := o.Base.Client.ContainerRegistry.ListRepositories(o.Base.Context, o.Base.Args[0], o.Base.Options)
	return repos, meta, err
}

// RepositoryGet ...
func (o *Options) RepositoryGet() (*govultr.ContainerRegistryRepo, error) {
	repo, _, err := o.Base.Client.ContainerRegistry.GetRepository(o.Base.Context, o.Base.Args[0], o.RepoName)
	return repo, err
}

// RepositoryUpdate ...
func (o *Options) RepositoryUpdate() error {
	_, _, err := o.Base.Client.ContainerRegistry.UpdateRepository(o.Base.Context, o.Base.Args[0], o.RepoName, o.RepoUpdateReq)
	return err
}

// RepositoryDelete ...
func (o *Options) RepositoryDelete() error {
	return o.Base.Client.ContainerRegistry.DeleteRepository(o.Base.Context, o.Base.Args[0], o.RepoName)
}

// CredentialsDocker ...
func (o *Options) CredentialsDocker() (*govultr.ContainerRegistryDockerCredentials, error) {
	cred, _, err := o.Base.Client.ContainerRegistry.CreateDockerCredentials(o.Base.Context, o.Base.Args[0], o.CredentialsDockerReq)
	return cred, err
}
