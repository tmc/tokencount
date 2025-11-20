package main

import (
	"context"
	"os/exec"
	"path"
	"testing"

	"rsc.io/script"
	"rsc.io/script/scripttest"
)

func TestTokencount(t *testing.T) {
	// Build current version to temp directory to ensure we're testing the right code
	tokencountPath := path.Join(t.TempDir(), "tokencount")
	if err := exec.Command("go", "build", "-o", tokencountPath).Run(); err != nil {
		t.Fatalf("failed to build tokencount: %v", err)
	}
	// Create custom commands without exec
	cmds := scripttest.DefaultCmds()
	delete(cmds, "exec") // Remove exec command
	cmds["tokencount"] = script.Program(tokencountPath, nil, 0)
	cmds["stdin-tokencount"] = stdinTokencountCmd(tokencountPath)
	engine := &script.Engine{
		Cmds:  cmds,
		Conds: scripttest.DefaultConds(),
	}
	scripttest.Test(t, context.Background(), engine, []string{}, "testdata/*.txt")
}

// stdinTokencountCmd runs tokencount with the first arg redirected to stdin, rest as normal args
func stdinTokencountCmd(tokencountPath string) script.Cmd {
	return script.Command(
		script.CmdUsage{
			Summary: "run tokencount with file as stdin",
			Args:    "file [args...]",
		},
		func(s *script.State, args ...string) (script.WaitFunc, error) {
			if len(args) < 1 {
				return nil, script.ErrUsage
			}
			stdinFile := s.Path(args[0])
			// Build command: tokencount [extra args] < file
			cmd := tokencountPath
			if len(args) > 1 {
				for _, arg := range args[1:] {
					cmd += " " + arg
				}
			}
			cmd += " < " + stdinFile
			// Use sh to redirect stdin with full path to tokencount
			return script.Program("sh", nil, 0).Run(s, "-c", cmd)
		},
	)
}
