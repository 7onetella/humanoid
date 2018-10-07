package main

import "strings"

// IsAllowed is the command allowed
func IsAllowed(cmd string) bool {
	for _, allowedCmd := range allowedCommmands {
		if strings.HasPrefix(cmd, allowedCmd) {
			return true
		}
	}
	return false
}

// IsApprovalRequired does the command require approval from peers
func IsApprovalRequired(cmd string) bool {
	for _, requiredCmd := range approvalRequired {
		if strings.HasPrefix(cmd, requiredCmd) {
			return true
		}
	}
	return false
}
