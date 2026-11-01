import { useState } from 'react'
import { measurePing, measureDownload, measureUpload } from './lib/speedtest'

export default function App() {
  const [result, setResult] = useState(null)
  const [running, setRunning] = useState(false)

  async function run() {
    setRunning(true)
    setResult(null)
    const ping = await measurePing()
    const download = await measureDownload()
    const upload = await measureUpload()
    setResult({ ...ping, download, upload })
    setRunning(false)
  }

  return (
    <div style={{ padding: 24, fontFamily: 'sans-serif' }}>
      <h1>tibeb-meter</h1>
      <button onClick={run} disabled={running}>
        {running ? 'Testing...' : 'Start test'}
      </button>
      {result && (
        <pre>{JSON.stringify(result, null, 2)}</pre>
      )}
    </div>
  )
}