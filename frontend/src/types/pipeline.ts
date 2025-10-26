export interface Job {
  id: string;
  name: string;
  status: 'success' | 'failure' | 'pending' | 'cancelled';
  startedAt: string;
  completedAt: string;
  duration: string;
  steps: JobStep[];
}

export interface JobStep {
  name: string;
  status: 'success' | 'failure' | 'pending' | 'cancelled';
  duration: string;
  log?: string;
}

export interface PipelineRun {
  runId: string;
  status: 'success' | 'failure' | 'pending' | 'cancelled';
  branch: string;
  commitSha: string;
  startedAt: string;
  duration: string;
  commitMessage: string;
  author: {
    name: string;
    email: string;
    avatarUrl: string;
  };
  repositoryUrl: string;
  jobs: Job[];
}