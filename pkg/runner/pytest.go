package runner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/envs"
	"github.com/kubeshop/testkube/pkg/executor"
	"github.com/kubeshop/testkube/pkg/executor/content"
	"github.com/kubeshop/testkube/pkg/executor/env"
	"github.com/kubeshop/testkube/pkg/executor/output"
	"github.com/kubeshop/testkube/pkg/executor/runner"
	"github.com/kubeshop/testkube/pkg/executor/scraper"
	"github.com/kubeshop/testkube/pkg/ui"
)

func NewPytestRunner(dependency string) (*PytestRunner error) {
    output.PrintLog(fmt.Sprintf("%s Preparing test runner", ui.IconTruck))
    params, err := envs.LoadTestkubeVariables()
    if err != nil {
        return nil, fmt.Errorf("could not initialize PytestRunner runner variables: %w", err)
    }
	return &PytestRunner{
		Params:  params,
		Fetcher: content.NewFetcher(""),
		Scraper: scraper.NewMinioScraper(
			params.Endpoint,
			params.AccessKeyID,
			params.SecretAccessKey,
			params.Location,
			params.Token,
			params.Bucket,
			params.Ssl,
		),
		dependency: dependency,
	}, nil
	}
}

// PytestRunner
type PytestRunner struct {
	Params     envs.Params
	Fetcher    content.ContentFetcher
	Scraper    scraper.Scraper
	dependency string
}

func (r *PytestRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {

    output.PrintLog(fmt.Sprintf("%s Preparing for test run", ui.IconTruck))

    runPath := filepath.Join(r.Params.DataDir, "repo", execution.Content.Repository.Path)
    if execution.Content.Repository != nil && execution.Content.Repository.WorkingDir != "" {
        runPath = filepath.Join(r.Params.DataDir, "repo", execution.Content.Repository.WorkingDir)
    }

    if _, err := os.Stat(filepath.Join(runPath, "requirement.txt")); err == nil {
        out, err := executor.Run(runPath, r.dependency, nil, "install")
        if err != nil {
            output.PrintLog(fmt.Sprintf("%s Dependency installation error %s", ui.IconCross, r.dependency))
            return result, fmt.Errorf("%s install error: %w\n\n%s", r.dependency, err, out)
        }
        output.PrintLog(fmt.Sprintf("%s Dependencies successfully installed", ui.IconBox))
    }

    	var runner string
    	var args []string

    	if r.dependency == "pip" {
    		runner = "pip"
    		args = []string{"pytest"}
    	}

	// use `execution.Variables` for variables passed from Test/Execution
	// variables of type "secret" will be automatically decoded
	envManager := env.NewManagerWithVars(execution.Variables)
	env.NewManager().GetReferenceVars(envManager.Variables)

	output.PrintEvent("Running", runPath, "pytest", args)

	out, runErr := executor.Run(runPath, runner, envManager, args...)

	out = envManager.ObfuscateSecrets(out)

    if runErr != nil {
        output.PrintLog(fmt.Sprintf("%s Test run failed", ui.IconCross))
        result = testkube.ExecutionResult{
            Status:     testkube.ExecutionStatusFailed,
            OutputType: "text/plain",
            Output:     fmt.Sprintf("pytest error: %s\n\n%s", runErr.Error(), out),
        }
    } else {
        result = testkube.ExecutionResult{
            Status:     testkube.ExecutionStatusPassed,
            OutputType: "text/plain",
            Output:     string(out),
        }
    }

    if runErr == nil {
        output.PrintLog(fmt.Sprintf("%s Test run successful", ui.IconCheckMark))
    }
    return result, runErr

// 	path, err := r.Fetcher.Fetch(execution.Content)
// 	if err != nil {
// 		return result, err
// 	}
//
// 	output.PrintEvent("created content path", path)

// 	fileInfo, err := os.Stat(path)
// 	if err != nil {
// 		return result, err
// 	}
//
// 	if !fileInfo.IsDir() {
// 		output.PrintEvent("using file", execution)
// 		// TODO implement file based test content for string, git-file, file-uri, git
// 		//      or remove if not used
// 	}

// 	if fileInfo.IsDir() {
// 		output.PrintEvent("using dir", execution)
// 		// TODO implement file based test content for git-dir, git
// 		//      or remove if not used
// 	}

	// TODO run executor here

	// error result should be returned if something is not ok
	// return result.Err(fmt.Errorf("some test execution related error occured"))

	// TODO return ExecutionResult
// 	return testkube.ExecutionResult{
// 		Status: testkube.ExecutionStatusPassed,
// 		Output: "exmaple test output",
// 	}, nil
}

// GetType returns runner type
func (r *PytestRunner) GetType() runner.Type {
	return runner.TypeMain
}
