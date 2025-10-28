import axios from 'axios';
import type { PipelineRun } from '../types/pipeline';

// Use the deployed backend URL
const API_BASE_URL = 'https://ci-pipeline-dashboard.onrender.com';

const api = axios.create({
  baseURL: API_BASE_URL,
});

// List view
export const getPipelineRuns = async (): Promise<PipelineRun[]> => {
  try {
    const response = await api.get('/api/pipelines');
    return response.data;
  } catch (error) {
    console.error('Error fetching pipeline runs:', error);
    throw error;
  }
};

// Detail view
export const getPipelineRunById = async (pipelineId: string): Promise<PipelineRun> => {
  try {
    const response = await api.get(`/api/pipelines/details?id=${pipelineId}`);
    return response.data; // Assumes backend returns object matching PipelineRun type
  } catch (error) {
    console.error('Error fetching pipeline run detail:', error);
    throw error;
  }
};
