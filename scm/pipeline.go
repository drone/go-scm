package scm

type (
	ObjectAttributes struct {
		ID         int
		IID        int
		Name       string
		Ref        string
		Tag        bool
		SHA        string
		BeforeSHA  string
		Source     string
		Status     string
		Stages     []string
		CreatedAt  string
		FinishedAt string
		Duration   int
		Variables  []Variable
		URL        string
	}

	Variable struct {
		Key   string
		Value string
	}

	MergeRequest struct {
		ID                  int
		IID                 int
		Title               string
		SourceBranch        string
		SourceProjectID     int
		TargetBranch        string
		TargetProjectID     int
		State               string
		MergeStatus         string
		DetailedMergeStatus string
		URL                 string
	}

	Project struct {
		ID                int
		Name              string
		Description       string
		WebURL            string
		AvatarURL         *string
		GitSSHURL         string
		GitHTTPURL        string
		Namespace         string
		VisibilityLevel   int
		PathWithNamespace string
		DefaultBranch     string
	}

	SourcePipeline struct {
		Project    SourceProject
		PipelineID int
		JobID      int
	}

	SourceProject struct {
		ID                int
		WebURL            string
		PathWithNamespace string
	}

	Build struct {
		ID             int
		Stage          string
		Name           string
		Status         string
		CreatedAt      string
		StartedAt      *string
		FinishedAt     *string
		Duration       *float64
		QueuedDuration *float64
		FailureReason  *string
		When           string
		Manual         bool
		AllowFailure   bool
		User           User
		Runner         *Runner
		ArtifactsFile  Artifacts
		Environment    *Environment
	}

	Runner struct {
		ID          int
		Description string
		Active      bool
		RunnerType  string
		IsShared    bool
		Tags        []string
	}

	Artifacts struct {
		Filename *string
		Size     *int
	}

	Environment struct {
		Name           string
		Action         string
		DeploymentTier string
	}
)
