// cmd/clean.go
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Nando-suka/git-cleaner-pro/internal/gitutils"
	"github.com/Nando-suka/git-cleaner-pro/internal/ui"
)

var (
	targetBranch  string
	includeRemote bool
	yesFlag       bool
)

func init() {
	rootCMD.AddCommand(cleanCmd)
	cleanCmd.Flags().StringVarP(&targetBranch, "target", "t", "", "Cabang target untuk pengecekan merge (default: 'main' atau 'master')")
	cleanCmd.Flags().BoolVarP(&includeRemote, "remote", "r", false, "Sertakan cabang remote dalam pembersihan")
	cleanCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "Lewati konfirmasi dan hapus semua cabang yang ditemukan")
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Membersihkan cabang lokal (dan remote) yang sudah di-merge",
	Long: `Memindai repositori Git di direktori saat ini, mencari cabang yang sudah di-merge ke cabang target.
Jika tidak ada flag --target yang diberikan, alat akan mendeteksi 'main' atau 'master' secara otomatis.`,
	RunE: runClean,
}

func runClean(cmd *cobra.Command, args []string) error {
	// 1. Cari repositori Git
	repo, err := gitutils.OpenRepository(".")
	if err != nil {
		return fmt.Errorf("gagal membuka repositori: %w", err)
	}

	// 2. Tentukan cabang target
	baseBranch, err := gitutils.DetectBaseBranch(repo)
	if err != nil {
		return err
	}
	if targetBranch != "" {
		baseBranch = targetBranch
	}
	fmt.Printf("Memeriksa cabang yang sudah di-merge ke '%s'...\n", baseBranch)

	// 3. Cari cabang yang sudah di-merge
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mergedBranches, err := gitutils.ListMergedBranches(ctx, repo, baseBranch, includeRemote)
	if err != nil {
		return err
	}

	// 4. Tampilkan hasil
	if len(mergedBranches) == 0 {
		fmt.Println("Tidak ada cabang yang perlu dibersihkan. Repositori Anda sudah rapi!")
		return nil
	}

	fmt.Printf("\n📋 Ditemukan %d cabang yang sudah di-merge:\n", len(mergedBranches))
	for _, branch := range mergedBranches {
		branchType := "lokal"
		if branch.IsRemote {
			branchType = "remote"
		}
		fmt.Printf("  - %s (%s)\n", branch.Name, branchType)
	}

	// 5. Minta konfirmasi (jika tidak ada flag --yes)
	if !yesFlag {
		shouldProceed, err := ui.ConfirmDeletionPrompt(len(mergedBranches))
		if err != nil {
			return err
		}
		if !shouldProceed {
			fmt.Println("Operasi dibatalkan.")
			return nil
		}
	}

	// 6. Hapus cabang
	fmt.Println("\n🗑️  Menghapus cabang...")
	for _, branch := range mergedBranches {
		err := gitutils.DeleteBranch(repo, branch)
		if err != nil {
			fmt.Printf("  Gagal menghapus %s: %v\n", branch.Name, err)
		} else {
			fmt.Printf("  Berhasil menghapus %s\n", branch.Name)
		}
	}

	fmt.Println("\n Pembersihan selesai!")
	return nil
}
