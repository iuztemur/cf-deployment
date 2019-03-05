package git

import (
	"fmt"

	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

//go:generate counterfeiter . Repository

type Repository interface {
	CreateRemote(*config.RemoteConfig) (*gogit.Remote, error)

	Log(*gogit.LogOptions) (object.CommitIter, error)
}

type Helper struct {
	r Repository
}

func NewHelper(r Repository) *Helper {
	return &Helper{r: r}
}

func (helper Helper) SetupRepository() error {
	config := &config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/cloudfoundry/cf-deployment.git"},
	}

	if _, err := helper.r.CreateRemote(config); err != nil {
		return fmt.Errorf("error adding remote: %v", err)
	}

	return nil
}

func (helper Helper) GetMajorVersion() (string, error) { return "", nil }

func (helper Helper) CheckoutSubDirectory(subdir, version string) error { return nil }
