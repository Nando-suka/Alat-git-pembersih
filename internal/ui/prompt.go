package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func ConfirmDeletionPrompt(count int) (bool, error) {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Hapus %d cabang yang terdaftar", count),
		IsConfirm: true,
		Default:   "n",
	}

	result, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	if result == "y" || result == "Y" {
		return true, nil
	}
	return false, nil
}
