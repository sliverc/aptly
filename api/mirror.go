package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smira/aptly/deb"
	"github.com/smira/aptly/query"
	"github.com/smira/aptly/task"
	"github.com/smira/aptly/utils"
	"sort"
	"strings"
)

func getVerifier(ignoreSignatures bool, keyRings []string) (utils.Verifier, error) {
	if ignoreSignatures {
		return nil, nil
	}

	verifier := &utils.GpgVerifier{}
	for _, keyRing := range keyRings {
		verifier.AddKeyring(keyRing)
	}

	err := verifier.InitKeyring()
	if err != nil {
		return nil, err
	}

	return verifier, nil
}

// GET /api/mirrors
func apiMirrorsList(c *gin.Context) {
	collectionFactory := context.NewCollectionFactory()
	collection := collectionFactory.RemoteRepoCollection()

	result := []*deb.RemoteRepo{}
	collection.ForEach(func(repo *deb.RemoteRepo) error {
		result = append(result, repo)
		return nil
	})

	c.JSON(200, result)
}

// POST /api/mirrors
func apiMirrorsCreate(c *gin.Context) {
	var err error
	var b struct {
		Name               string `binding:"required"`
		ArchiveURL         string `binding:"required"`
		Distribution       string
		Components         []string
		Architectures      []string
		DownloadSources    bool
		DownloadUdebs      bool
		Filter             string
		FilterWithDeps     bool
		SkipComponentCheck bool
		IgnoreSignatures   bool
		Keyrings           []string
	}

	b.DownloadSources = context.Config().DownloadSourcePackages
	b.IgnoreSignatures = context.Config().GpgDisableVerify
	b.Architectures = context.ArchitecturesList()

	if !c.Bind(&b) {
		return
	}

	collectionFactory := context.NewCollectionFactory()
	collection := collectionFactory.RemoteRepoCollection()

	if strings.HasPrefix(b.ArchiveURL, "ppa:") {
		b.ArchiveURL, b.Distribution, b.Components, err = deb.ParsePPA(b.ArchiveURL, context.Config())
		if err != nil {
			c.Fail(400, err)
			return
		}
	}

	if b.Filter != "" {
		_, err = query.Parse(b.Filter)
		if err != nil {
			c.Fail(400, fmt.Errorf("unable to create mirror: %s", err))
			return
		}
	}

	repo, err := deb.NewRemoteRepo(b.Name, b.ArchiveURL, b.Distribution, b.Components, b.Architectures,
		b.DownloadSources, b.DownloadUdebs)

	if err != nil {
		c.Fail(400, fmt.Errorf("unable to create mirror: %s", err))
		return
	}

	repo.Filter = b.Filter
	repo.FilterWithDeps = b.FilterWithDeps
	repo.SkipComponentCheck = b.SkipComponentCheck
	repo.DownloadSources = b.DownloadSources
	repo.DownloadUdebs = b.DownloadUdebs

	verifier, err := getVerifier(b.IgnoreSignatures, b.Keyrings)
	if err != nil {
		c.Fail(400, fmt.Errorf("unable to initialize GPG verifier: %s", err))
		return
	}

	downloader := context.NewDownloader(nil)
	err = repo.Fetch(downloader, verifier)
	if err != nil {
		c.Fail(400, fmt.Errorf("unable to fetch mirror: %s", err))
		return
	}

	err = collection.Add(repo)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to add mirror: %s", err))
		return
	}

	c.JSON(201, repo)
}

// DELETE /api/mirrors/:name
func apiMirrorsDrop(c *gin.Context) {
	name := c.Params.ByName("name")
	force := c.Request.URL.Query().Get("force") == "1"

	collectionFactory := context.NewCollectionFactory()
	mirrorCollection := collectionFactory.RemoteRepoCollection()
	snapshotCollection := collectionFactory.SnapshotCollection()

	repo, err := mirrorCollection.ByName(name)
	if err != nil {
		c.Fail(404, fmt.Errorf("unable to drop: %s", err))
		return
	}

	resources := []string{string(repo.Key())}
	taskName := fmt.Sprintf("Delete mirror %s", name)
	task, conflictErr := runTaskInBackground(taskName, resources, func(out *task.Output, detail *task.Detail) error {
		err := repo.CheckLock()
		if err != nil {
			return fmt.Errorf("unable to drop: %s", err)
		}

		if !force {
			snapshots := snapshotCollection.ByRemoteRepoSource(repo)

			if len(snapshots) > 0 {
				return fmt.Errorf("won't delete mirror with snapshots, use 'force=1' to override")
			}
		}

		return mirrorCollection.Drop(repo)
	})

	if conflictErr != nil {
		c.Error(conflictErr, conflictErr.Tasks)
		c.AbortWithStatus(409)
		return
	}

	c.JSON(202, task)
}

// GET /api/mirrors/:name
func apiMirrorsShow(c *gin.Context) {
	collectionFactory := context.NewCollectionFactory()
	collection := collectionFactory.RemoteRepoCollection()

	name := c.Params.ByName("name")
	repo, err := collection.ByName(name)
	if err != nil {
		c.Fail(404, fmt.Errorf("unable to show: %s", err))
		return
	}

	err = collection.LoadComplete(repo)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to show: %s", err))
	}

	c.JSON(200, repo)
}

// GET /api/mirrors/:name/packages
func apiMirrorsPackages(c *gin.Context) {
	collectionFactory := context.NewCollectionFactory()
	collection := collectionFactory.RemoteRepoCollection()

	name := c.Params.ByName("name")
	repo, err := collection.ByName(name)
	if err != nil {
		c.Fail(404, fmt.Errorf("unable to show: %s", err))
		return
	}

	err = collection.LoadComplete(repo)
	if err != nil {
		c.Fail(500, fmt.Errorf("unable to show: %s", err))
	}

	if repo.LastDownloadDate.IsZero() {
		c.Fail(404, fmt.Errorf("Unable to show package list, mirror hasn't been downloaded yet."))
		return
	}

	reflist := repo.RefList()
	result := []*deb.Package{}

	list, err := deb.NewPackageListFromRefList(reflist, collectionFactory.PackageCollection(), nil)
	if err != nil {
		c.Fail(404, err)
		return
	}

	queryS := c.Request.URL.Query().Get("q")
	if queryS != "" {
		q, err := query.Parse(c.Request.URL.Query().Get("q"))
		if err != nil {
			c.Fail(400, err)
			return
		}

		withDeps := c.Request.URL.Query().Get("withDeps") == "1"
		architecturesList := []string{}

		if withDeps {
			if len(context.ArchitecturesList()) > 0 {
				architecturesList = context.ArchitecturesList()
			} else {
				architecturesList = list.Architectures(false)
			}

			sort.Strings(architecturesList)

			if len(architecturesList) == 0 {
				c.Fail(400, fmt.Errorf("unable to determine list of architectures, please specify explicitly"))
				return
			}
		}

		list.PrepareIndex()

		list, err = list.Filter([]deb.PackageQuery{q}, withDeps,
			nil, context.DependencyOptions(), architecturesList)
		if err != nil {
			c.Fail(500, fmt.Errorf("unable to search: %s", err))
		}
	}

	if c.Request.URL.Query().Get("format") == "details" {
		list.ForEach(func(p *deb.Package) error {
			result = append(result, p)
			return nil
		})

		c.JSON(200, result)
	} else {
		c.JSON(200, list.Strings())
	}
}

// PUT /api/mirrors/:name
func apiMirrorsUpdate(c *gin.Context) {
	var (
		err    error
		remote *deb.RemoteRepo
	)

	var b struct {
		Name               string
		Filter             string
		FilterWithDeps     bool
		ForceComponents    bool
		DownloadSources    bool
		DownloadUdebs      bool
		Architectures      []string
		Components         []string
		SkipComponentCheck bool
		MaxTries           int
		IgnoreSignatures   bool
		Keyrings           []string
		ForceUpdate        bool
		DownloadLimit      int64
	}

	collectionFactory := context.NewCollectionFactory()
	collection := collectionFactory.RemoteRepoCollection()

	remote, err = collection.ByName(c.Params.ByName("name"))
	if err != nil {
		c.Fail(404, err)
		return
	}

	b.Name = remote.Name
	b.DownloadUdebs = remote.DownloadUdebs
	b.DownloadSources = remote.DownloadSources
	b.SkipComponentCheck = remote.SkipComponentCheck
	b.FilterWithDeps = remote.FilterWithDeps
	b.Filter = remote.Filter
	b.Architectures = remote.Architectures
	b.Components = remote.Components

	if !c.Bind(&b) {
		return
	}

	if b.Name != remote.Name {
		_, err = collection.ByName(b.Name)
		if err == nil {
			c.Fail(409, fmt.Errorf("unable to rename: mirror %s already exists", b.Name))
			return
		}
	}

	if b.DownloadUdebs != remote.DownloadUdebs {
		if remote.IsFlat() && b.DownloadUdebs {
			c.Fail(400, fmt.Errorf("unable to update: flat mirrors don't support udebs"))
			return
		}
	}

	remote.Name = b.Name
	remote.DownloadUdebs = b.DownloadUdebs
	remote.DownloadSources = b.DownloadSources
	remote.SkipComponentCheck = b.SkipComponentCheck
	remote.FilterWithDeps = b.FilterWithDeps
	remote.Filter = b.Filter
	remote.Architectures = b.Architectures
	remote.Components = b.Components

	verifier, err := getVerifier(b.IgnoreSignatures, b.Keyrings)
	if err != nil {
		c.Fail(400, fmt.Errorf("unable to initialize GPG verifier: %s", err))
		return
	}

	resources := []string{string(remote.Key())}
	task, conflictErr := runTaskInBackground("Update mirror "+b.Name, resources, func(out *task.Output, detail *task.Detail) error {

		downloader := context.NewDownloader(out)
		err := remote.Fetch(downloader, verifier)
		if err != nil {
			return fmt.Errorf("unable to update: %s", err)
		}

		if !b.ForceUpdate {
			err = remote.CheckLock()
			if err != nil {
				return fmt.Errorf("unable to update: %s", err)
			}
		}

		if b.MaxTries <= 0 {
			b.MaxTries = 1
		}

		err = remote.DownloadPackageIndexes(out, downloader, collectionFactory, b.SkipComponentCheck, b.MaxTries)
		if err != nil {
			return fmt.Errorf("unable to update: %s", err)
		}

		if remote.Filter != "" {
			var filterQuery deb.PackageQuery

			filterQuery, err = query.Parse(remote.Filter)
			if err != nil {
				return fmt.Errorf("unable to update: %s", err)
			}

			_, _, err = remote.ApplyFilter(context.DependencyOptions(), filterQuery)
			if err != nil {
				return fmt.Errorf("unable to update: %s", err)
			}
		}

		queue, downloadSize, err := remote.BuildDownloadQueue(context.PackagePool())
		if err != nil {
			return fmt.Errorf("unable to update: %s", err)
		}

		defer func() {
			// on any interruption, unlock the mirror
			err := context.ReOpenDatabase()
			if err == nil {
				remote.MarkAsIdle()
				collection.Update(remote)
			}
		}()

		remote.MarkAsUpdating()
		err = collection.Update(remote)
		if err != nil {
			return fmt.Errorf("unable to update: %s", err)
		}

		count := len(queue)
		taskDetail := struct {
			TotalDownloadSize int64
			RemainingDownloadSize int64
			TotalNumberOfPackages int
			RemainingNumberOfPackages int
		}{
			downloadSize, downloadSize, count, count,
		}
		detail.Store(taskDetail)

		// In separate goroutine (to avoid blocking main), push queue to downloader
		ch := make(chan error, count)
		go func() {
			for _, task := range queue {
				downloader.DownloadWithChecksum(remote.PackageURL(task.RepoURI).String(), task.DestinationPath, ch, task.Checksums, b.SkipComponentCheck, b.MaxTries)

				taskDetail.RemainingDownloadSize -= task.Checksums.Size
				taskDetail.RemainingNumberOfPackages--
				detail.Store(taskDetail)
			}

			queue = nil
		}()

		// Wait for all downloads to finish
		var errors []string
		for count > 0 {
			select {
			case err = <-ch:
				if err != nil {
					errors = append(errors, err.Error())
				}
				count--
			}
		}

		remote.FinalizeDownload()
		return nil
	})

	if conflictErr != nil {
		c.Error(conflictErr, conflictErr.Tasks)
		c.AbortWithStatus(409)
		return
	}

	c.JSON(202, task)
}
