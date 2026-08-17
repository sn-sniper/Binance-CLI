package cmd

import (
	"bytes"
	"testing"
)

func TestRootCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("rootCmd.Execute() error = %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("binance-cli [command]")) {
		t.Errorf("expected help output to contain usage, got:\n%s", output)
	}
}

func TestSubcommandsRegistered(t *testing.T) {
	expectedCommands := []string{"balance", "price", "symbols"}

	cmdMap := make(map[string]bool)
	for _, cmd := range rootCmd.Commands() {
		cmdMap[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !cmdMap[expected] {
			t.Errorf("expected %q subcommand to be registered in rootCmd", expected)
		}
	}
}
