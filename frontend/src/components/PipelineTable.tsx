import { useState, useEffect } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  TableSortLabel,
  Box,
  CircularProgress,
} from '@mui/material';
import { format, parseISO } from 'date-fns';
import { getPipelineRuns } from '../services/pipelineService';
import type { PipelineRun } from '../types/pipeline';
import PipelineDetail from './PipelineDetail';

type Order = 'asc' | 'desc';

interface HeadCell {
  id: keyof PipelineRun;
  label: string;
  sortable: boolean;
}

const headCells: HeadCell[] = [
  { id: 'runId', label: 'Run ID', sortable: false },
  { id: 'status', label: 'Status', sortable: true },
  { id: 'branch', label: 'Branch', sortable: false },
  { id: 'commitSha', label: 'Commit SHA', sortable: false },
  { id: 'startedAt', label: 'Started At', sortable: true },
  { id: 'duration', label: 'Duration', sortable: false },
];

export default function PipelineTable() {
  const [pipelineRuns, setPipelineRuns] = useState<PipelineRun[]>([]);
  const [loading, setLoading] = useState(true);
  const [order, setOrder] = useState<Order>('desc');
  const [orderBy, setOrderBy] = useState<keyof PipelineRun>('startedAt');
  const [selectedPipeline, setSelectedPipeline] = useState<PipelineRun | null>(null);

  useEffect(() => {
    fetchPipelineRuns();
  }, []);

  const fetchPipelineRuns = async () => {
    try {
      const data = await getPipelineRuns();
      setPipelineRuns(data);
    } catch (error) {
      console.error('Error fetching pipeline runs:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleRequestSort = (property: keyof PipelineRun) => {
    const isAsc = orderBy === property && order === 'asc';
    setOrder(isAsc ? 'desc' : 'asc');
    setOrderBy(property);
  };

  const getStatusColor = (status: string) => {
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

  const sortedPipelineRuns = [...pipelineRuns].sort((a, b) => {
    if (orderBy === 'startedAt') {
      const dateA = parseISO(a.startedAt);
      const dateB = parseISO(b.startedAt);
      return order === 'asc' ? dateA.getTime() - dateB.getTime() : dateB.getTime() - dateA.getTime();
    }
    
    if (orderBy === 'status') {
      return order === 'asc'
        ? a.status.localeCompare(b.status)
        : b.status.localeCompare(a.status);
    }
    
    return 0;
  });

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <>
      <TableContainer component={Paper}>
        <Table>
        <TableHead>
          <TableRow>
            {headCells.map((headCell) => (
              <TableCell key={headCell.id}>
                {headCell.sortable ? (
                  <TableSortLabel
                    active={orderBy === headCell.id}
                    direction={orderBy === headCell.id ? order : 'asc'}
                    onClick={() => handleRequestSort(headCell.id)}
                  >
                    {headCell.label}
                  </TableSortLabel>
                ) : (
                  headCell.label
                )}
              </TableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {sortedPipelineRuns.slice(0, 25).map((run) => (
            <TableRow 
              key={run.runId}
              hover
              onClick={() => setSelectedPipeline(run)}
              sx={{ cursor: 'pointer' }}
            >
              <TableCell>{run.runId}</TableCell>
              <TableCell>
                <Box
                  sx={{
                    backgroundColor: getStatusColor(run.status),
                    color: 'white',
                    padding: '4px 8px',
                    borderRadius: '4px',
                    display: 'inline-block',
                  }}
                >
                  {run.status.toUpperCase()}
                </Box>
              </TableCell>
              <TableCell>{run.branch}</TableCell>
              <TableCell>{run.commitSha}</TableCell>
              <TableCell>{format(parseISO(run.startedAt), 'MMM d, yyyy HH:mm:ss')}</TableCell>
              <TableCell>{run.duration}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
      <PipelineDetail
        open={Boolean(selectedPipeline)}
        onClose={() => setSelectedPipeline(null)}
        runId={selectedPipeline?.runId || null}
      />
    </>
  );
}