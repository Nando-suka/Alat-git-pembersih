package gitutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type BranchInfo struct {
	Name     string
	IsRemote bool
	Ref      *plumbing.Reference
}

func OpenRepository(path string) (*git.Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("direktori ini bukan termasuk repositori git: %w", err)
	}
	return repo, nil
}

func DetectBaseBranch(repo *git.Repository) (string, error) {
	headRef, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("gagal mendapatkan HEAD: %w", err)
	}
	branchName := headRef.Name().Short()

	possibleBases := []string{"main", "master"}
	for _, base := range possibleBases {
		if _, err := repo.Branch(base); err == nil {
			return base, nil
		}
	}

	return branchName, nil
}

func ListMergedBranches(ctx context.Context, repo *git.Repository, baseBranch string, includeRemote bool) ([]BranchInfo, error) {
	var mergedBranches []BranchInfo

	// Dapatkan referensi ke cabang dasar
	baseRef, err := repo.Reference(plumbing.NewBranchReferenceName(baseBranch), true)
	if err != nil {
		return nil, fmt.Errorf("cabang dasar '%s' tidak ditemukan: %w", baseBranch, err)
	}
	baseCommit, err := repo.CommitObject(baseRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan commit untuk cabang dasar: %w", err)
	}

	// Iterasi semua cabang lokal
	branches, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("gagal membaca daftar cabang: %w", err)
	}
	defer branches.Close()

	err = branches.ForEach(func(ref *plumbing.Reference) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			branchName := ref.Name().Short()
			// Jangan sertakan cabang dasar itu sendiri
			if branchName == baseBranch {
				return nil
			}
			commit, err := repo.CommitObject(ref.Hash())
			if err != nil {
				return nil // Lewati jika gagal
			}
			isAncestor, err := commit.IsAncestor(baseCommit)
			if err != nil {
				return nil
			}
			if isAncestor {
				mergedBranches = append(mergedBranches, BranchInfo{
					Name:     branchName,
					IsRemote: false,
					Ref:      ref,
				})
			}
			return nil
		}
	})
	if err != nil {
		return nil, err
	}

	if includeRemote {
		remoteBranches, err := listRemoteBranches(ctx, repo, baseCommit)
		if err != nil {
			return nil, err
		}
		mergedBranches = append(mergedBranches, remoteBranches...)
	}
	return mergedBranches, nil
}

func listRemoteBranches(ctx context.Context, repo *git.Repository, baseCommit *object.Commit) ([]BranchInfo, error) {
	var mergedBranches []BranchInfo

	remotes, err := repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan daftar remote: %w", err)
	}
	if len(remotes) == 0 {
		return nil, nil
	}
	remote := remotes[0] // Gunakan remote pertama (biasanya 'origin')

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("gagal listing referensi remote: %w", err)
	}

	for _, ref := range refs {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if !ref.Name().IsRemote() {
				continue
			}

			if ref.Name().Short() == "HEAD" {
				continue
			}
			commit, err := repo.CommitObject(ref.Hash())
			if err != nil {
				continue
			}
			isAncestor, err := commit.IsAncestor(baseCommit)
			if err != nil {
				continue
			}
			if isAncestor {
				mergedBranches = append(mergedBranches, BranchInfo{
					Name:     ref.Name().Short(),
					IsRemote: true,
					Ref:      ref,
				})
			}
		}

	}
	return mergedBranches, nil
}

func DeleteBranch(repo *git.Repository, branch BranchInfo) error {
	if branch.IsRemote {
		// Hapus cabang remote
		remoteName := strings.Split(branch.Name, "/")[0]
		branchName := strings.Join(strings.Split(branch.Name, "/")[1:], "/")
		remote, err := repo.Remote(remoteName)
		if err != nil {
			return fmt.Errorf("gagal mendapatkan remote '%s': %w", remoteName, err)
		}
		refSpec := config.RefSpec(fmt.Sprintf(":refs/heads/%s", branchName))
		return remote.Push(&git.PushOptions{
			RemoteName: remoteName,
			RefSpecs:   []config.RefSpec{refSpec},
		})
	} else {
		// Hapus cabang lokal
		return repo.Storer.RemoveReference(branch.Ref.Name())
	}
}
