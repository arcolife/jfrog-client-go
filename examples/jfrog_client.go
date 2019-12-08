package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/arcolife/jfrog-client-go/bintray"
	"github.com/arcolife/jfrog-client-go/bintray/auth"
	"github.com/arcolife/jfrog-client-go/bintray/services"
	"github.com/arcolife/jfrog-client-go/bintray/services/packages"
	"github.com/arcolife/jfrog-client-go/bintray/services/repositories"
	"github.com/arcolife/jfrog-client-go/bintray/services/versions"
	"github.com/arcolife/jfrog-client-go/utils/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	// "github.com/jfrog/jfrog-client-go/utils/io/fileutils"
)

var file *os.File

type (
	// Package - struct fields must be public in order for unmarshal to
	// correctly populate the data.
	Package struct {
		Package string                `yaml:"package"`
		Config  packages.Params       `yaml:"config"`
		Upload  services.UploadParams `yaml:"upload"`
	}

	// Repo - struct fields must be public in order for unmarshal to
	// correctly populate the data.
	Repo struct {
		Name     string              `yaml:"name"`
		Subject  string              `yaml:"subject"`
		Config   repositories.Config `yaml:"config"`
		Packages []Package           `yaml:"packages"`
	}

	// Bintray - struct fields must be public in order for unmarshal to
	// correctly populate the data.
	Bintray struct {
		Repos []Repo `yaml:"repos"`
	}

	// BintrayConfig -> entrypoint for all configs
	BintrayConfig struct {
		Threads int  `yaml:"threads"`
		Cleanup bool `yaml:"cleanup"`
		Bintray `yaml:"bintray"`
	}
)

// PrettyPrint to print a struct like JSON
func PrettyPrint(v interface{}, label string) (err error) {
	fmt.Println(label, "=>")
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func (config *BintrayConfig) initConfig() *bintray.ServicesManager {

	configdata, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Print(errors.Wrap(err, "Error reading file"))
		os.Exit(1)
	}

	err = yaml.Unmarshal(configdata, &config)
	if err != nil {
		fmt.Print(errors.Wrap(err, "Error Unmarshalling file"))
		os.Exit(1)
	}

	btDetails := auth.NewBintrayDetails()
	btDetails.SetUser(os.Getenv("BINTRAY_USER"))
	btDetails.SetKey(os.Getenv("BINTRAY_KEY"))
	btDetails.SetDefPackageLicense("Apache 2.0")

	serviceConfig := bintray.NewConfigBuilder().
		SetBintrayDetails(btDetails).
		SetDryRun(false).
		SetThreads(config.Threads).
		Build()

	btManager, err := bintray.New(serviceConfig)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	return btManager
}

func (config *BintrayConfig) deleteRepo(btManager *bintray.ServicesManager, repo *Repo) error {
	fmt.Printf("\nDeleting Repo.. [%s]\n", repo.Name)
	repoPath := repositories.Path{Subject: repo.Subject, Repo: repo.Name}
	return errors.Wrap(btManager.ExecDeleteRepoRest(&repoPath), "Repo non-existent")
}

func (config *BintrayConfig) deletePackage(btManager *bintray.ServicesManager, repo *Repo, pack *Package) error {
	pkg := fmt.Sprintf("%s/%s/%s", repo.Subject, repo.Name, pack.Package)
	fmt.Printf("\nDeleting Package.. [%s]\n", pkg)
	packagePath, _ := packages.CreatePath(pkg)
	return errors.Wrap(btManager.DeletePackage(packagePath), "Package non-existent")
}

func (config *BintrayConfig) createRepo(btManager *bintray.ServicesManager, repo *Repo) (bool, error) {
	repoPath := repositories.Path{Subject: repo.Subject, Repo: repo.Name}
	return btManager.CreateReposIfNeeded(&repoPath, &repo.Config, repo.Config.RepoConfigFilePath)
}

func (config *BintrayConfig) createPackage(btManager *bintray.ServicesManager, repo *Repo, pack *Package) error {
	pkg := fmt.Sprintf("%s/%s/%s", repo.Subject, repo.Name, pack.Package)
	pack.Config.Path, _ = packages.CreatePath(pkg)
	return btManager.CreatePackage(&pack.Config)
}

func (config *BintrayConfig) publishPackage(btManager *bintray.ServicesManager, repo *Repo, pack *Package) {
	pkg := fmt.Sprintf("%s/%s/%s", repo.Subject, repo.Name, pack.Package)
	fmt.Printf("\nPublishing GPG Signatures.. [%s]", pkg)
	versionPathString := fmt.Sprintf("%s/%s", pkg, pack.Upload.Version)
	versionPath, _ := versions.CreatePath(versionPathString)
	err := btManager.PublishVersion(versionPath)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}

func (config *BintrayConfig) uploadPackage(btManager *bintray.ServicesManager, repo *Repo, pack *Package) {
	pkg := fmt.Sprintf("%s/%s/%s", repo.Subject, repo.Name, pack.Package)
	fmt.Printf("\nUploading Files to Package [%s] with Publish: [%t]\n", pkg, pack.Upload.Publish)
	versionPath := fmt.Sprintf("%s/%s", pkg, pack.Upload.Version)
	pack.Upload.Path, _ = versions.CreatePath(versionPath)
	PrettyPrint(&pack, "Package")
	uploaded, failed, err := btManager.UploadFiles(&pack.Upload)
	fmt.Println("UPLOADED", uploaded)
	fmt.Println("FAILED: ", failed)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}

func (config *BintrayConfig) signPackageVersion(btManager *bintray.ServicesManager, repo *Repo, pack *Package) {
	pkg := fmt.Sprintf("%s/%s/%s", repo.Subject, repo.Name, pack.Package)
	fmt.Printf("\nSigning versioned Package files.. [%s]", pkg)
	versionPath := fmt.Sprintf("%s/%s", pkg, pack.Upload.Version)
	path, _ := versions.CreatePath(versionPath)
	err := btManager.GpgSignVersion(path, os.Getenv("BINTRAY_ADMIN_GPG_PASSPHRASE"))
	if err != nil {
		fmt.Println(err)
	}
}

func (config *BintrayConfig) calcMetadata(btManager *bintray.ServicesManager, repo *Repo, pack *Package) bool {
	pkg := fmt.Sprintf("%s/%s/%s", repo.Subject, repo.Name, pack.Package)
	fmt.Printf("\nScheduling metadata calculation.. [%s]\n", pkg)
	versionPath := fmt.Sprintf("%s/%s", pkg, pack.Upload.Version)
	path, _ := versions.CreatePath(versionPath)
	scheduledOk, err := btManager.CalcMetadata(path)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	return scheduledOk
}

func (config *BintrayConfig) showPackage(btManager *bintray.ServicesManager, repo *Repo, pack *Package) {
	pkg := fmt.Sprintf("%s/%s/%s", repo.Subject, repo.Name, pack.Package)
	fmt.Printf("\nPackage details.. [%s]", pkg)
	pkgPath, _ := packages.CreatePath(pkg)
	err := btManager.ShowPackage(pkgPath)
	if err != nil {
		fmt.Println(err)
	}
}

func (config *BintrayConfig) checkDetails(btManager *bintray.ServicesManager) {
	var err error
	for _, repo := range config.Bintray.Repos {
		for _, pack := range repo.Packages {
			pkg := fmt.Sprintf("%s/%s/%s", repo.Subject, repo.Name, pack.Package)
			fmt.Printf("\nChecking details.. [%s]", pkg)

			// Repository
			RepoExistsOk, _ := btManager.IsRepoExists(&repositories.Path{Subject: repo.Subject, Repo: repo.Name})
			repoPath := fmt.Sprintf("%s/%s", repo.Subject, repo.Name)
			if RepoExistsOk != true {
				fmt.Printf("\nRepo does not exist.. [%s]", repoPath)
				fmt.Printf("\nCreating Repo..")
				_, err = config.createRepo(btManager, &repo)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				fmt.Printf("\nRepo already exists.. [%s]", repoPath)
			}

			// Package
			pkgPath, _ := packages.CreatePath(pkg)
			PackageExistsOk, _ := btManager.IsPackageExists(pkgPath)
			if PackageExistsOk != true {
				fmt.Printf("\nPackage [%s] does not exist", pkg)
				fmt.Printf("\nCreating Package..\n")
				err = config.createPackage(btManager, &repo, &pack)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				fmt.Printf("\nPackage already exists.. [%s]\n", pkg)
			}
		}
	}
}

func (config *BintrayConfig) cleanup(btManager *bintray.ServicesManager) {
	var err error
	for _, repo := range config.Bintray.Repos {
		for _, pack := range repo.Packages {
			err = config.deletePackage(btManager, &repo, &pack)
			if err != nil {
				fmt.Print(err)
			} else {
				config.deleteRepo(btManager, &repo)
			}
		}
	}
}

func main() {
	log.SetLogger(log.NewLogger(log.INFO, file))

	config := BintrayConfig{}
	btManager := config.initConfig()
	if config.Cleanup {
		config.cleanup(btManager)
		config.checkDetails(btManager)
	}
	for _, repo := range config.Bintray.Repos {
		for _, pack := range repo.Packages {
			config.uploadPackage(btManager, &repo, &pack)
			config.signPackageVersion(btManager, &repo, &pack)
			if pack.Upload.Publish == true {
				config.publishPackage(btManager, &repo, &pack)
			}
			config.calcMetadata(btManager, &repo, &pack)
			config.showPackage(btManager, &repo, &pack)
		}
	}
}
