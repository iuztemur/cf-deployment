package git

import (
	"fmt"
  "gopkg.in/src-d/go-git.v4/plumbing"
  "gopkg.in/src-d/go-git.v4/plumbing/storer"

  gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

//go:generate counterfeiter . Repository

type Repository interface {
	CreateRemote(*config.RemoteConfig) (*gogit.Remote, error)

	Log(*gogit.LogOptions) (object.CommitIter, error)

	Tags() (storer.ReferenceIter, error)
}

type Helper struct {
	TagsMap map[plumbing.Hash]*plumbing.Reference
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


// Code for emulating git-describe taken from here:
// https://github.com/edupo/semver-cli/blob/master/gitWrapper/git.go
// e.g.
// >> git describe --tags
// >> v7.6.0-13-g74262de9

func (helper Helper) getTagMap() error {
	tags, err := helper.r.Tags()
	if err != nil {
		return err
	}

	err = tags.ForEach(func(t *plumbing.Reference) error{
		helper.TagsMap[t.Hash()] = t
		return nil
	})

	return err
}


func (helper Helper) describe(reference *plumbing.Reference) (string, error) {
	r := helper.r
	log, err := r.Log(&gogit.LogOptions{
		From:  reference.Hash(),
		Order: gogit.LogOrderCommitterTime,
	})

	err = helper.getTagMap()
	if err != nil {
		return "", err
	}

	var tag *plumbing.Reference
	var count int
	err = log.ForEach(func(c *object.Commit) error{
	  if t, ok := helper.TagsMap[c.Hash]; ok {
	    tag = t
    }
	  if tag != nil {
	    return storer.ErrStop
    }
	  count++
	  return nil
  })
	if count == 0 {
    return fmt.Sprint(tag.Name().Short()), nil
  } else {
    return fmt.Sprintf("%v-%v-%v",
                        tag.Name().Short(),
                        count,
                        tag.Hash().String()[0:8],
                        ), nil
  }
}


func (helper Helper) GetMajorVersion(reference *plumbing.Reference) (string, error) {
  _, err:= helper.describe(reference)

  if err != nil {
    return "", err
  }

  // verify version is numbered


  // pull first 3 characters: v7.


  //append 0.0

  //return string


  return "v1.0.0", nil
}

func (helper Helper) CheckoutSubDirectory(subdir, version string) error { return nil }
