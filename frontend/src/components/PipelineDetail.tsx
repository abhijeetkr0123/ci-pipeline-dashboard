import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  IconButton,
  Typography,
  Box,
  Avatar,
  Link,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Chip,
  Grid,
  Divider,
  CircularProgress,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import GitHubIcon from '@mui/icons-material/GitHub';
import type { PipelineRun, Job } from '../types/pipeline';
import { format, parseISO } from 'date-fns';
import { getPipelineRunById } from '../services/pipelineService';

interface PipelineDetailProps {
  open: boolean;
  onClose: () => void;
  runId: string | null;
}

const getStatusColor = (status?: string) => {
  switch (status) {
    case 'success':
      return '#2E7D32';
    case 'failure':
      return '#D32F2F';
    case 'pending':
      return '#ED6C02';
    default:
      return '#757575';
  }
};

const JobStatus: React.FC<{ job: Job }> = ({ job }) => {
  const [expanded, setExpanded] = useState(false);

  return (
    <Accordion expanded={expanded} onChange={() => setExpanded(!expanded)} sx={{ mb: 1 }}>
      <AccordionSummary expandIcon={<ExpandMoreIcon />}>
        <Box sx={{ display: 'flex', alignItems: 'center', width: '100%' }}>
          <Typography sx={{ flexGrow: 1 }}>{job.name || 'Unnamed Job'}</Typography>
          <Chip
            label={(job.status || 'unknown').toUpperCase()}
            size="small"
            sx={{ backgroundColor: getStatusColor(job.status), color: 'white', ml: 2 }}
          />
          <Typography variant="body2" color="text.secondary" sx={{ ml: 2 }}>
            {job.duration || '-'}
          </Typography>
        </Box>
      </AccordionSummary>
      <AccordionDetails>
        <Box>
          {job.steps?.length ? (
            job.steps.map((step, idx) => (
              <Box key={idx} sx={{ mb: 1, pl: 2 }}>
                <Grid container spacing={2} alignItems="center">
                  <Grid item xs>
                    <Typography variant="body2">{step.name || '-'}</Typography>
                  </Grid>
                  <Grid item>
                    <Chip
                      label={(step.status || 'unknown').toUpperCase()}
                      size="small"
                      sx={{ backgroundColor: getStatusColor(step.status), color: 'white' }}
                    />
                  </Grid>
                  <Grid item>
                    <Typography variant="body2" color="text.secondary">
                      {step.duration || '-'}
                    </Typography>
                  </Grid>
                </Grid>
              </Box>
            ))
          ) : (
            <Typography variant="body2">No steps available</Typography>
          )}
        </Box>
      </AccordionDetails>
    </Accordion>
  );
};

export default function PipelineDetail({ open, onClose, runId }: PipelineDetailProps) {
  const [pipeline, setPipeline] = useState<PipelineRun | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!runId) return;

    const fetchPipeline = async () => {
      setLoading(true);
      try {
        const data = await getPipelineRunById(runId);
        setPipeline(data);
      } catch (error) {
        console.error('Error fetching pipeline detail:', error);
        setPipeline(null);
      } finally {
        setLoading(false);
      }
    };

    fetchPipeline();
  }, [runId]);

  if (!open) return null;

  const author = pipeline?.author || { name: 'Unknown', email: 'N/A', avatarUrl: '' };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        <Box display="flex" alignItems="center" justifyContent="space-between">
          <Typography variant="h6">Pipeline Run: {pipeline?.runId || '-'}</Typography>
          <IconButton onClick={onClose}>
            <CloseIcon />
          </IconButton>
        </Box>
      </DialogTitle>

      <DialogContent>
        {loading ? (
          <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
            <CircularProgress />
          </Box>
        ) : !pipeline ? (
          <Typography>No data available</Typography>
        ) : (
          <>
            {/* Author & Commit Info */}
            <Box sx={{ mb: 3 }}>
              <Grid container spacing={2}>
                <Grid item xs={12}>
                  <Box display="flex" alignItems="center" gap={2} mb={2}>
                    <Avatar src={author.avatarUrl} alt={author.name} />
                    <Box>
                      <Typography variant="subtitle1">{author.name}</Typography>
                      <Typography variant="body2" color="text.secondary">{author.email}</Typography>
                    </Box>
                  </Box>
                </Grid>

                <Grid item xs={12}>
                  <Typography variant="body1" sx={{ mb: 1 }}>
                    {pipeline.commitMessage || 'No commit message'}
                  </Typography>
                  <Box display="flex" gap={1} alignItems="center">
                    <GitHubIcon fontSize="small" />
                    <Link
                      href={`${pipeline.repositoryUrl || '#'}${pipeline.commitSha ? `/commit/${pipeline.commitSha}` : ''}`}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      {pipeline.commitSha || '-'}
                    </Link>
                  </Box>
                </Grid>
              </Grid>
            </Box>

            {/* Pipeline Info */}
            <Divider sx={{ my: 2 }} />
            <Box sx={{ mb: 2 }}>
              <Typography variant="h6" gutterBottom>Pipeline Information</Typography>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">Branch</Typography>
                  <Typography variant="body1">{pipeline.branch || '-'}</Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">Started At</Typography>
                  <Typography variant="body1">
                    {pipeline.startedAt ? format(parseISO(pipeline.startedAt), 'PPpp') : '-'}
                  </Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">Status</Typography>
                  <Chip
                    label={(pipeline.status || 'unknown').toUpperCase()}
                    size="small"
                    sx={{ backgroundColor: getStatusColor(pipeline.status), color: 'white', mt: 0.5 }}
                  />
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">Duration</Typography>
                  <Typography variant="body1">{pipeline.duration || '-'}</Typography>
                </Grid>
              </Grid>
            </Box>

            {/* Jobs */}
            <Divider sx={{ my: 2 }} />
            <Box>
              <Typography variant="h6" gutterBottom>Jobs</Typography>
              {pipeline.jobs?.length ? (
                pipeline.jobs.map((job) => <JobStatus key={job.id} job={job} />)
              ) : (
                <Typography variant="body2">No jobs available</Typography>
              )}
            </Box>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}
