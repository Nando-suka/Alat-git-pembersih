package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   "git-cleaner-pro",
	Short: "Git cleaner pro membersihkan cabang yang sudah di-merge dengan aman dan interaktif",
	Long: `Alat CLI untuk membantu developer menjaga repositori Git tetap rapi.
Alat ini akan memindai cabang lokal dan remote, mengidentifikasi cabang yang sudah di-merge ke cabang utama (main/master),
dan menghapusnya setelah mendapatkan konfirmasi dari Anda.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Jalankan 'git-cleaner-pro untuk memulai pembersihan")
	},
}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
