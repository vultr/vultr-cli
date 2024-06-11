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
	credentialsLong = `Commands for accessing the credentials on registries`
	//nolint:gosec
	credentialsExample = `
	# Full example
	vultr-cli container-registry credentials
	`
	credentialsDockerLong = `Create the credential string used by docker`
	//nolint:gosec
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
	vultr-cli container-registry repository update 4dcdc52e-9c63-401e-8c5f-1582490fe09c --image-name="my-thing" 
	--description="new description"

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

// NewCmdContainerRegistry provides the CLI command functionality for container registry
func NewCmdContainerRegistry(base *cli.Base) *cobra.Command { //nolint:funlen,gocyclo
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "container-registry",
		Short:   "Commands to interact with container registries",
		Aliases: []string{"cr"},
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
		Short:   "List all container registries",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)
			regs, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error retrieving container registry list : %v", err)
			}

			data := &ContainerRegistriesPrinter{Registries: regs, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	list.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
	)

	// Get
	get := &cobra.Command{
		Use:     "get <Registry ID>",
		Short:   "Get a container registry",
		Aliases: []string{"g"},
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, err := o.get()
			if err != nil {
				return fmt.Errorf("error retrieving container registry info : %v", err)
			}

			data := &ContainerRegistryPrinter{Registry: reg}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create a container registry",
		Aliases: []string{"c"},
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing 'name' flag for container registry create : %v", errNa)
			}

			region, errRe := cmd.Flags().GetString("region")
			if errRe != nil {
				return fmt.Errorf("error parsing 'region' flag for container registry create : %v", errRe)
			}

			public, errPu := cmd.Flags().GetBool("public")
			if errPu != nil {
				return fmt.Errorf("error parsing 'public' flag for container registry create : %v", errPu)
			}

			plan, errPl := cmd.Flags().GetString("plan")
			if errPl != nil {
				return fmt.Errorf("error parsing 'plan' flag for container registry create : %v", errPl)
			}

			o.CreateReq = &govultr.ContainerRegistryReq{
				Name:   name,
				Region: region,
				Public: public,
				Plan:   plan,
			}

			reg, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating container registry : %v", err)
			}

			data := &ContainerRegistryPrinter{Registry: reg}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	create.Flags().StringP("name", "n", "", "The name to use for the container registry")
	if err := create.MarkFlagRequired("name"); err != nil {
		fmt.Printf("error marking container registry create 'name' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("region", "i", "", "The ID of the region in which to create the container registry")
	if err := create.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking container registry create 'region' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().BoolP(
		"public",
		"p",
		false,
		"If the registry is publicly available. Should be true | false (default is false)",
	)
	if err := create.MarkFlagRequired("public"); err != nil {
		fmt.Printf("error marking container registry create 'public' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP("plan", "l", "", "The type of plan to use for the container registry")
	if err := create.MarkFlagRequired("plan"); err != nil {
		fmt.Printf("error marking container registry create 'plan' flag required: %v", err)
		os.Exit(1)
	}

	// Update
	update := &cobra.Command{
		Use:     "update <Registry ID>",
		Short:   "Update a container registry",
		Aliases: []string{"u"},
		Long:    updateLong,
		Example: updateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			public, errPu := cmd.Flags().GetBool("public")
			if errPu != nil {
				return fmt.Errorf("error parsing 'public' flag for container registry update : %v", errPu)
			}

			plan, errPl := cmd.Flags().GetString("plan")
			if errPl != nil {
				return fmt.Errorf("error parsing 'plan' flag for container registry update : %v", errPl)
			}

			o.UpdateReq = &govultr.ContainerRegistryUpdateReq{
				Plan: govultr.StringToStringPtr(plan),
			}

			if cmd.Flags().Changed("public") {
				o.UpdateReq.Public = govultr.BoolToBoolPtr(public)
			}

			if err := o.update(); err != nil {
				return fmt.Errorf("error updating container registry : %v", err)
			}

			o.Base.Printer.Display(printer.Info("container registry has been updated"), nil)

			return nil
		},
	}

	update.Flags().StringP("plan", "p", "", "Name of the plan used for the container registry")
	update.Flags().BoolP("public", "b", false, "The container registry availability status")

	// Delete
	del := &cobra.Command{
		Use:     "delete <Registry ID>",
		Short:   "Delete a container registry",
		Aliases: []string{"destroy", "d"},
		Long:    deleteLong,
		Example: deleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting container registry : %v", err)
			}

			o.Base.Printer.Display(printer.Info("container registry has been deleted"), nil)

			return nil
		},
	}

	// Plans
	plans := &cobra.Command{
		Use:     "plans",
		Short:   "List the plan names for a container registry",
		Aliases: []string{"p"},
		Long:    plansLong,
		Example: plansExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			plans, err := o.plans()
			if err != nil {
				return fmt.Errorf("error retrieving plans for container registry : %v", err)
			}

			data := &ContainerRegistryPlansPrinter{Plans: plans.Plans}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Regions
	regions := &cobra.Command{
		Use:     "regions",
		Short:   "List the available regions for a container registry",
		Aliases: []string{"i"},
		Long:    regionsLong,
		Example: regionsExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			regs, meta, err := o.regions()
			if err != nil {
				return fmt.Errorf("error retrieving regions for container registry : %v", err)
			}

			data := &ContainerRegistryRegionsPrinter{Regions: regs, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Repository
	repository := &cobra.Command{
		Use:     "repository",
		Short:   "Interact with container registry repositories",
		Aliases: []string{"r", "repo"},
		Long:    repoLong,
		Example: repoExample,
	}

	// Repository List
	repoList := &cobra.Command{
		Use:     "list <Registry ID>",
		Short:   "List all container registries",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)
			repos, meta, err := o.repositoryList()
			if err != nil {
				return fmt.Errorf("error retrieving repositories for container registry : %v", err)
			}

			data := &ContainerRegistryRepositoriesPrinter{Repositories: repos, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Repository Get
	repoGet := &cobra.Command{
		Use:     "get <Registry ID>",
		Short:   "Get a container registry repository",
		Aliases: []string{"g"},
		Long:    getLong,
		Example: getExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errIm := cmd.Flags().GetString("image-name")
			if errIm != nil {
				return fmt.Errorf("error parsing 'image-name' flag for container registry repository get : %v", errIm)
			}

			o.RepoName = name

			repo, err := o.repositoryGet()
			if err != nil {
				return fmt.Errorf("error getting container registry repository : %v", err)
			}

			data := &ContainerRegistryRepositoryPrinter{Repository: repo}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	repoGet.Flags().StringP("image-name", "i", "", "The name of the image/repo")
	if err := repoGet.MarkFlagRequired("image-name"); err != nil {
		fmt.Printf("error marking get container registry repository 'image-name' flag required: %v", err)
		os.Exit(1)
	}

	// Repository Update
	repoUpdate := &cobra.Command{
		Use:     "update <Registry ID>",
		Short:   "Update a container registry repository",
		Aliases: []string{"u"},
		Long:    repoUpdateLong,
		Example: repoUpdateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errIm := cmd.Flags().GetString("image-name")
			if errIm != nil {
				return fmt.Errorf(
					"error parsing 'image-name' flag for container registry repository update : %v",
					errIm,
				)
			}

			description, errDe := cmd.Flags().GetString("description")
			if errDe != nil {
				return fmt.Errorf(
					"error parsing 'description' flag for container registry repository update : %v",
					errDe,
				)
			}

			o.RepoName = name
			o.RepoUpdateReq = &govultr.ContainerRegistryRepoUpdateReq{
				Description: description,
			}

			if err := o.repositoryUpdate(); err != nil {
				return fmt.Errorf("error updating container registry repository : %v", err)
			}

			o.Base.Printer.Display(printer.Info("container registry repository has been updated"), nil)

			return nil
		},
	}

	repoUpdate.Flags().StringP("image-name", "i", "", "The name of the image/repo")
	if err := repoUpdate.MarkFlagRequired("image-name"); err != nil {
		fmt.Printf("error marking update container registry repository 'image-name' flag required: %v", err)
		os.Exit(1)
	}

	repoUpdate.Flags().StringP("description", "d", "", "The description of the image/repo")
	if err := repoUpdate.MarkFlagRequired("description"); err != nil {
		fmt.Printf("error marking update container registry repository 'description' flag required: %v", err)
		os.Exit(1)
	}

	// Repository Delete
	repoDelete := &cobra.Command{
		Use:     "delete <Registry ID>",
		Short:   "Delete a container registry repository",
		Aliases: []string{"destroy", "d"},
		Long:    repoDeleteLong,
		Example: repoDeleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a container registry ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errIm := cmd.Flags().GetString("image-name")
			if errIm != nil {
				return fmt.Errorf(
					"error parsing 'image-name' flag for container registry repository delete : %v",
					errIm,
				)
			}

			o.RepoName = name

			if err := o.repositoryDelete(); err != nil {
				return fmt.Errorf("error deleting container registry repository : %v", err)
			}

			o.Base.Printer.Display(printer.Info("container registry repository has been deleted"), nil)

			return nil
		},
	}

	repoDelete.Flags().StringP("image-name", "i", "", "The name of the image/repo")
	if err := repoDelete.MarkFlagRequired("image-name"); err != nil {
		fmt.Printf("error marking delete container registry repository 'image-name' flag required: %v", err)
		os.Exit(1)
	}

	repository.AddCommand(
		repoGet,
		repoList,
		repoUpdate,
		repoDelete,
	)

	// Credentials
	credentials := &cobra.Command{
		Use:     "credentials",
		Short:   "Commands for container registry credentials",
		Aliases: []string{""},
		Long:    credentialsLong,
		Example: credentialsExample,
	}

	// Credentials Docker
	credentialsDocker := &cobra.Command{
		Use:     "docker <Registry ID>",
		Short:   "Create Docker credentials for a container registry",
		Aliases: []string{"d"},
		Long:    credentialsDockerLong,
		Example: credentialsDockerExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			expiry, errEx := cmd.Flags().GetInt("expiry-seconds")
			if errEx != nil {
				return fmt.Errorf(
					"error parsing 'expiry-seconds' flag for container registry docker creds : %v",
					errEx,
				)
			}

			access, errAc := cmd.Flags().GetBool("read-write")
			if errAc != nil {
				return fmt.Errorf("error parsing 'read-write' flag for container registry docker creds : %v", errAc)
			}

			o.CredentialsDockerReq = &govultr.DockerCredentialsOpt{
				ExpirySeconds: govultr.IntToIntPtr(expiry),
				WriteAccess:   govultr.BoolToBoolPtr(access),
			}

			cred, err := o.credentialsDocker()
			if err != nil {
				return fmt.Errorf("error generating container registry repository docker credentials : %v", err)
			}

			data := &ContainerRegistryCredentialDockerPrinter{Credential: cred}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	credentialsDocker.Flags().IntP(
		"expiry-seconds",
		"e",
		0,
		"(optional) The seconds until these credentials expire.  Default is 0, never",
	)

	credentialsDocker.Flags().BoolP(
		"read-write",
		"w",
		false,
		"(optional) Whether or not these credentials have write access.  Should be true or false.  Default is false",
	)

	credentials.AddCommand(
		credentialsDocker,
	)

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

type options struct {
	Base                 *cli.Base
	CreateReq            *govultr.ContainerRegistryReq
	UpdateReq            *govultr.ContainerRegistryUpdateReq
	RepoName             string
	RepoUpdateReq        *govultr.ContainerRegistryRepoUpdateReq
	CredentialsDockerReq *govultr.DockerCredentialsOpt
}

func (o *options) list() ([]govultr.ContainerRegistry, *govultr.Meta, error) {
	cr, meta, _, err := o.Base.Client.ContainerRegistry.List(o.Base.Context, o.Base.Options)
	return cr, meta, err
}

func (o *options) get() (*govultr.ContainerRegistry, error) {
	cr, _, err := o.Base.Client.ContainerRegistry.Get(o.Base.Context, o.Base.Args[0])
	return cr, err
}

func (o *options) create() (*govultr.ContainerRegistry, error) {
	cr, _, err := o.Base.Client.ContainerRegistry.Create(o.Base.Context, o.CreateReq)
	return cr, err
}

func (o *options) update() error {
	_, _, err := o.Base.Client.ContainerRegistry.Update(o.Base.Context, o.Base.Args[0], o.UpdateReq)
	return err
}

func (o *options) del() error {
	return o.Base.Client.ContainerRegistry.Delete(o.Base.Context, o.Base.Args[0])
}

func (o *options) plans() (*govultr.ContainerRegistryPlans, error) {
	plans, _, err := o.Base.Client.ContainerRegistry.ListPlans(o.Base.Context)
	return plans, err
}

func (o *options) regions() ([]govultr.ContainerRegistryRegion, *govultr.Meta, error) {
	regions, meta, _, err := o.Base.Client.ContainerRegistry.ListRegions(o.Base.Context)
	return regions, meta, err
}

func (o *options) repositoryList() ([]govultr.ContainerRegistryRepo, *govultr.Meta, error) {
	repos, meta, _, err := o.Base.Client.ContainerRegistry.ListRepositories(o.Base.Context, o.Base.Args[0], o.Base.Options) //nolint:lll
	return repos, meta, err
}

func (o *options) repositoryGet() (*govultr.ContainerRegistryRepo, error) {
	repo, _, err := o.Base.Client.ContainerRegistry.GetRepository(o.Base.Context, o.Base.Args[0], o.RepoName)
	return repo, err
}

func (o *options) repositoryUpdate() error {
	_, _, err := o.Base.Client.ContainerRegistry.UpdateRepository(
		o.Base.Context,
		o.Base.Args[0],
		o.RepoName,
		o.RepoUpdateReq,
	)

	return err
}

func (o *options) repositoryDelete() error {
	return o.Base.Client.ContainerRegistry.DeleteRepository(o.Base.Context, o.Base.Args[0], o.RepoName)
}

func (o *options) credentialsDocker() (*govultr.ContainerRegistryDockerCredentials, error) {
	cred, _, err := o.Base.Client.ContainerRegistry.CreateDockerCredentials(
		o.Base.Context,
		o.Base.Args[0],
		o.CredentialsDockerReq,
	)
	return cred, err
}
