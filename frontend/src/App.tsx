import { Container, Typography, Box } from '@mui/material'
import PipelineTable from './components/PipelineTable'

function App() {
  return (
    <Container maxWidth="lg">
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          GitHub Actions Pipeline Runs
        </Typography>
        <PipelineTable />
      </Box>
    </Container>
  )
}

export default App
