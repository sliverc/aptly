package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smira/aptly/deb"
	"github.com/smira/aptly/utils"
	"github.com/smira/aptly/task"
	"strings"
)

// SigningOptions is a shared between publish API GPG options structure
type SigningOptions struct {
	Skip           bool
	Batch          bool
	GpgKey         string
	Keyring        string
	SecretKeyring  string
	Passphrase     string
	PassphraseFile string
}

func getSigner(options *SigningOptions) (utils.Signer, error) {
	if options.Skip {
		return nil, nil
	}

	signer := &utils.GpgSigner{}
	signer.SetKey(options.GpgKey)
	signer.SetKeyRing(options.Keyring, options.SecretKeyring)
	signer.SetPassphrase(options.Passphrase, options.PassphraseFile)
	signer.SetBatch(options.Batch)

	err := signer.Init()
	if err != nil {
		return nil, err
	}

	return signer, nil
}

// Replace '_' with '/' and double '__' with single '_'
func parseEscapedPath(path string) string {
	result := strings.Replace(strings.Replace(path, "_", "/", -1), "//", "_", -1)
	if result == "" {
		result = "."
	}
	return result
}

// GET /publish
func apiPublishList(c *gin.Context) {
	collectionFactory := context.NewCollectionFactory()
	collection := collectionFactory.PublishedRepoCollection()

	result := make([]*deb.PublishedRepo, 0, collection.Len())

	err := collection.ForEach(func(repo *deb.PublishedRepo) error {
		err := collection.LoadComplete(repo, collectionFactory)
		if err != nil {
			return err
		}

		result = append(result, repo)

		return nil
	})

	if err != nil {
		c.Fail(500, err)
		return
	}

	c.JSON(200, result)
}

// POST /publish/:prefix
func apiPublishRepoOrSnapshot(c *gin.Context) {
	param := parseEscapedPath(c.Params.ByName("prefix"))
	storage, prefix := deb.ParsePrefix(param)

	var b struct {
		SourceKind string `binding:"required"`
		Sources    []struct {
			Component string
			Name      string `binding:"required"`
		} `binding:"required"`
		Distribution   string
		Label          string
		Origin         string
		ForceOverwrite bool
		SkipContents   *bool
		Architectures  []string
		Signing        SigningOptions
	}

	if !c.Bind(&b) {
		return
	}

	signer, err := getSigner(&b.Signing)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to initialize GPG signer: %s", err))
		return
	}

	if len(b.Sources) == 0 {
		c.Fail(400, fmt.Errorf("unable to publish: soures are empty"))
		return
	}

	var components []string
	var names []string
	var sources []interface{}
	var resources []string
	collectionFactory := context.NewCollectionFactory()

	if b.SourceKind == "snapshot" {
		var snapshot *deb.Snapshot

		snapshotCollection := collectionFactory.SnapshotCollection()

		for _, source := range b.Sources {
			components = append(components, source.Component)
			names = append(names, source.Name)

			snapshot, err = snapshotCollection.ByName(source.Name)
			if err != nil {
				c.Fail(404, fmt.Errorf("unable to publish: %s", err))
				return
			}

			resources = append(resources, string(snapshot.ResourceKey()))
			err = snapshotCollection.LoadComplete(snapshot)
			if err != nil {
				c.Fail(500, fmt.Errorf("unable to publish: %s", err))
				return
			}

			sources = append(sources, snapshot)
		}
	} else if b.SourceKind == "local" {
		var localRepo *deb.LocalRepo

		localCollection := collectionFactory.LocalRepoCollection()

		for _, source := range b.Sources {
			components = append(components, source.Component)
			names = append(names, source.Name)

			localRepo, err = localCollection.ByName(source.Name)
			if err != nil {
				c.Fail(404, fmt.Errorf("unable to publish: %s", err))
				return
			}

			resources = append(resources, string(localRepo.Key()))
			err = localCollection.LoadComplete(localRepo)
			if err != nil {
				c.Fail(500, fmt.Errorf("unable to publish: %s", err))
			}

			sources = append(sources, localRepo)
		}
	} else {
		c.Fail(400, fmt.Errorf("unknown SourceKind"))
		return
	}

	published, err := deb.NewPublishedRepo(storage, prefix, b.Distribution, b.Architectures, components, sources, collectionFactory)
	if err != nil {
		c.Fail(400, fmt.Errorf("unable to publish: %s", err))
		return
	}

	resources = append(resources, string(published.Key()))
	collection := collectionFactory.PublishedRepoCollection()

	taskName := fmt.Sprintf("Publish %s: %s", b.SourceKind, strings.Join(names, ", "))
	task, conflictErr := runTaskInBackground(taskName, resources, func(out *task.Output) error {
		published.Origin = b.Origin
		published.Label = b.Label

		published.SkipContents = context.Config().SkipContentsPublishing
		if b.SkipContents != nil {
			published.SkipContents = *b.SkipContents
		}

		duplicate := collection.CheckDuplicate(published)
		if duplicate != nil {
			collectionFactory.PublishedRepoCollection().LoadComplete(duplicate, collectionFactory)
			return fmt.Errorf("prefix/distribution already used by another published repo: %s", duplicate)
		}

		err := published.Publish(context.PackagePool(), context, collectionFactory, signer, out, b.ForceOverwrite)
		if err != nil {
			return fmt.Errorf("unable to publish: %s", err)
		}

		err = collection.Add(published)
		if err != nil {
			return fmt.Errorf("unable to save to DB: %s", err)
		}

		return nil
	})

	if conflictErr != nil {
		c.Error(conflictErr, conflictErr.Tasks)
		c.AbortWithStatus(412)
		return
	}

	c.JSON(202, task)
}

// PUT /publish/:prefix/:distribution
func apiPublishUpdateSwitch(c *gin.Context) {
	param := parseEscapedPath(c.Params.ByName("prefix"))
	storage, prefix := deb.ParsePrefix(param)
	distribution := c.Params.ByName("distribution")

	var b struct {
		ForceOverwrite bool
		Signing        SigningOptions
		SkipContents   *bool
		Snapshots      []struct {
			Component string `binding:"required"`
			Name      string `binding:"required"`
		}
	}

	if !c.Bind(&b) {
		return
	}

	signer, err := getSigner(&b.Signing)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to initialize GPG signer: %s", err))
		return
	}

	collectionFactory := context.NewCollectionFactory()
	collection := collectionFactory.PublishedRepoCollection()

	published, err := collection.ByStoragePrefixDistribution(storage, prefix, distribution)
	if err != nil {
		c.Fail(404, fmt.Errorf("unable to update: %s", err))
		return
	}
	err = collection.LoadComplete(published, collectionFactory)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to update: %s", err))
		return
	}

	var updatedComponents []string
	var updatedSnapshots []string
	var resources []string

	if published.SourceKind == "local" {
		if len(b.Snapshots) > 0 {
			c.Fail(400, fmt.Errorf("snapshots shouldn't be given when updating local repo"))
			return
		}
		updatedComponents = published.Components()
		for _, component := range updatedComponents {
			published.UpdateLocalRepo(component)
		}
	} else if published.SourceKind == "snapshot" {
		publishedComponents := published.Components()
		for _, snapshotInfo := range b.Snapshots {
			if !utils.StrSliceHasItem(publishedComponents, snapshotInfo.Component) {
				c.Fail(404, fmt.Errorf("component %s is not in published repository", snapshotInfo.Component))
				return
			}

			snapshotCollection := collectionFactory.SnapshotCollection()
			snapshot, err := snapshotCollection.ByName(snapshotInfo.Name)
			if err != nil {
				c.Fail(404, err)
				return
			}

			err = snapshotCollection.LoadComplete(snapshot)
			if err != nil {
				c.Fail(500, err)
				return
			}

			published.UpdateSnapshot(snapshotInfo.Component, snapshot)
			updatedComponents = append(updatedComponents, snapshotInfo.Component)
			updatedSnapshots = append(updatedSnapshots, snapshot.Name)
		}
	} else {
		c.Fail(500, fmt.Errorf("unknown published repository type"))
		return
	}

	if b.SkipContents != nil {
		published.SkipContents = *b.SkipContents
	}

	resources = append(resources, string(published.Key()))
	taskName := fmt.Sprintf("Update published %s (%s): %s", published.SourceKind, strings.Join(updatedComponents, " "), strings.Join(updatedSnapshots, ", "))
	task, conflictErr := runTaskInBackground(taskName, resources, func(out *task.Output) error {
		err := published.Publish(context.PackagePool(), context, collectionFactory, signer, out, b.ForceOverwrite)
		if err != nil {
			return fmt.Errorf("unable to update: %s", err)
		}

		err = collection.Update(published)
		if err != nil {
			return fmt.Errorf("unable to save to DB: %s", err)
		}

		err = collection.CleanupPrefixComponentFiles(published.Prefix, updatedComponents,
			context.GetPublishedStorage(storage), collectionFactory, out)
		if err != nil {
			return fmt.Errorf("unable to update: %s", err)
		}

		return nil
	})

	if conflictErr != nil {
		c.Error(conflictErr, conflictErr.Tasks)
		c.AbortWithStatus(412)
		return
	}


	c.JSON(202, task)
}

// DELETE /publish/:prefix/:distribution
func apiPublishDrop(c *gin.Context) {
	force := c.Request.URL.Query().Get("force") == "1"

	param := parseEscapedPath(c.Params.ByName("prefix"))
	storage, prefix := deb.ParsePrefix(param)
	distribution := c.Params.ByName("distribution")

	collectionFactory := context.NewCollectionFactory()
	collection := collectionFactory.PublishedRepoCollection()

	published, err := collection.ByStoragePrefixDistribution(storage, prefix, distribution)
	if err != nil {
		c.Fail(400, fmt.Errorf("unable to remove: %s", err))
		return
	}

	resources := []string{string(published.Key())}

	taskName := fmt.Sprintf("Delete published %s (%s)", prefix, distribution)
	task, conflictErr := runTaskInBackground(taskName, resources, func(out *task.Output) error {
		err := collection.Remove(context, storage, prefix, distribution,
			collectionFactory, out, force)
		if err != nil {
			return fmt.Errorf("unable to drop: %s", err)
		}

		return nil
	})

	if conflictErr != nil {
		c.Error(conflictErr, conflictErr.Tasks)
		c.AbortWithStatus(412)
		return
	}

	c.JSON(202, task)
}
