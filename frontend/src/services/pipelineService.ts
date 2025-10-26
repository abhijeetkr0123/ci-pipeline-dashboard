import axios from 'axios';
import type { PipelineRun } from '../types/pipeline';


// List view
export const getPipelineRuns = async (): Promise<PipelineRun[]> => {
  try {
    const response = await axios.get(`/api/pipelines`);
    return response.data;
  } catch (error) {
    console.error('Error fetching pipeline runs:', error);
    throw error;
  }
};

// Detail view
export const getPipelineRunById = async (runId: string): Promise<PipelineRun> => {
  try {
    const response = await axios.get(`/api/pipelines/details?id=${runId}`);
    return response.data; // Assumes backend returns object matching PipelineRun type
  } catch (error) {
    console.error('Error fetching pipeline run detail:', error);
    throw error;
  }
};
