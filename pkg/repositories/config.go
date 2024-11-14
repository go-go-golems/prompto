package repositories

import "github.com/go-go-golems/prompto/pkg"

type RepositoryConfig struct {
	Repositories []string
}

func NewRepositoryConfig(repositories []string) *RepositoryConfig {
	return &RepositoryConfig{
		Repositories: repositories,
	}
}

func (rc *RepositoryConfig) LoadRepositories() ([]pkg.Repository, error) {
	var repositories []pkg.Repository
	for _, repoPath := range rc.Repositories {
		repo := pkg.NewRepository(repoPath)
		err := repo.LoadPromptos()
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, *repo)
	}
	return repositories, nil
}
