package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"

	cp "github.com/otiai10/copy"
)

func TestRun(t *testing.T) {

	t.Run("runner should run test based on execution data", func(t *testing.T) {
		// given
		// runner := NewPytestRunner()
		// execution := testkube.NewQueuedExecution()
		// execution.Content = testkube.NewStringTestContent("hello I'm test content")

		// // when
		// result, err := runner.Run(*execution)

		// // then
		// assert.NoError(t, err)
		// assert.Equal(t, result.Status, testkube.ExecutionStatusPassed)

		runner := NewPytestRunner()


		// setup
		tempDir, _ := os.MkdirTemp("", "*")
		os.Setenv("RUNNER_DATADIR", tempDir)
		repoDir := filepath.Join(tempDir, "repo")
		os.Mkdir(repoDir, 0755)
		_ = cp.Copy("../../examples", repoDir)

		runner, err := NewPytestRunner("pip")
		if err != nil {
			t.Fail()
		}

		result, err := runner.Run(testkube.Execution{
			Content: &testkube.TestContent{
				Type_: string(testkube.TestContentTypeGitDir),
				Repository: &testkube.Repository{
					Type_:  "git",
					Uri:    "",
					Branch: "master",
					Path:   "",
				},
			},
		})

		fmt.Printf("RESULT: %+v\n", result)
		fmt.Printf("ERROR:  %+v\n", err)

		t.Fail()

	})

}
