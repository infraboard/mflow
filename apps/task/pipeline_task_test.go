package task_test

import (
	"testing"

	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/test/tools"
)

func TestNewPipelineTask(t *testing.T) {
	p := pipeline.NewDefaultPipeline()
	tools.MustReadYamlFile("impl/test/pipeline.yml", p)
	req := pipeline.NewRunPipelineRequest("")
	pt := task.NewPipelineTask(p, req)

	tasks := pt.JobTasks()
	for i := range tasks.Items {
		task := tasks.Items[i]
		t.Log(task.Spec.TaskId, task.Spec.JobName, task.Status.Stage)
	}

	// 即将运行的tasks
	nexts, err := pt.NextRun()
	if err != nil {
		t.Fatal(err)
	}

	for i := range nexts.Items {
		nexts, err = pt.NextRun()
		t.Log(nexts, err)

		next := nexts.Items[i]
		next.Status.MarkedSuccess()
	}
}
