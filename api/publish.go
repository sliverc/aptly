package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smira/aptly/deb"
	"github.com/smira/aptly/utils"
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
	var sources []interface{}
	collectionFactory := context.NewCollectionFactory()

	if b.SourceKind == "snapshot" {
		var snapshot *deb.Snapshot

		snapshotCollection := collectionFactory.SnapshotCollection()

		for _, source := range b.Sources {
			components = append(components, source.Component)

			snapshot, err = snapshotCollection.ByName(source.Name)
			if err != nil {
				c.Fail(404, fmt.Errorf("unable to publish: %s", err))
				return
			}

			err = snapshotCollection.LoadComplete(snapshot)
			if err != nil {
				c.Fail(500, fmt.Errorf("unable to publish: %s", err))
				return
			}

			sources = append(sources, snapshot)
		}
	} else if b.SourceKind == deb.SourceLocalRepo {
		var localRepo *deb.LocalRepo

		localCollection := collectionFactory.LocalRepoCollection()

		for _, source := range b.Sources {
			components = append(components, source.Component)

			localRepo, err = localCollection.ByName(source.Name)
			if err != nil {
				c.Fail(404, fmt.Errorf("unable to publish: %s", err))
				return
			}

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

	collection := collectionFactory.PublishedRepoCollection()

	published, err := deb.NewPublishedRepo(storage, prefix, b.Distribution, b.Architectures, components, sources, collectionFactory)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to publish: %s", err))
		return
	}
	published.Origin = b.Origin
	published.Label = b.Label

	published.SkipContents = context.Config().SkipContentsPublishing
	if b.SkipContents != nil {
		published.SkipContents = *b.SkipContents
	}

	duplicate := collection.CheckDuplicate(published)
	if duplicate != nil {
		collectionFactory.PublishedRepoCollection().LoadComplete(duplicate, collectionFactory)
		c.Fail(400, fmt.Errorf("prefix/distribution already used by another published repo: %s", duplicate))
		return
	}

	err = published.Publish(context.PackagePool(), context, collectionFactory, signer, nil, b.ForceOverwrite)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to publish: %s", err))
		return
	}

	err = collection.Add(published)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to save to DB: %s", err))
		return
	}

	c.JSON(201, published)
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

	if published.SourceKind == deb.SourceLocalRepo {
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
			snapshot, err2 := snapshotCollection.ByName(snapshotInfo.Name)
			if err != nil {
				c.Fail(404, err2)
				return
			}

			err2 = snapshotCollection.LoadComplete(snapshot)
			if err2 != nil {
				c.Fail(500, err2)
				return
			}

			published.UpdateSnapshot(snapshotInfo.Component, snapshot)
			updatedComponents = append(updatedComponents, snapshotInfo.Component)
		}
	} else {
		c.Fail(500, fmt.Errorf("unknown published repository type"))
		return
	}

	if b.SkipContents != nil {
		published.SkipContents = *b.SkipContents
	}

	err = published.Publish(context.PackagePool(), context, collectionFactory, signer, nil, b.ForceOverwrite)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to update: %s", err))
		return
	}

	err = collection.Update(published)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to save to DB: %s", err))
		return
	}

	err = collection.CleanupPrefixComponentFiles(published.Prefix, updatedComponents,
		context.GetPublishedStorage(storage), collectionFactory, nil)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to update: %s", err))
		return
	}

	c.JSON(200, published)
}

// DELETE /publish/:prefix/:distribution
func apiPublishDrop(c *gin.Context) {
	force := c.Request.URL.Query().Get("force") == "1"

	param := parseEscapedPath(c.Params.ByName("prefix"))
	storage, prefix := deb.ParsePrefix(param)
	distribution := c.Params.ByName("distribution")

	collectionFactory := context.NewCollectionFactory()
	collection := collectionFactory.PublishedRepoCollection()

	err := collection.Remove(context, storage, prefix, distribution,
		collectionFactory, context.Progress(), force)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to drop: %s", err))
		return
	}

	c.JSON(200, gin.H{})
}
