// Time a single round trip to the server.
async function pingOnce() {
  const start = performance.now()
  await fetch(`/api/ping?t=${Date.now()}`, { cache: 'no-store' })
  return performance.now() - start
}

// Latency = average of several pings. Jitter = how much they vary.
export async function measurePing(count = 10) {
  const times = []
  for (let i = 0; i < count; i++) {
    times.push(await pingOnce())
  }
  const avg = times.reduce((a, b) => a + b, 0) / times.length
  const jitter =
    times.slice(1).reduce((sum, t, i) => sum + Math.abs(t - times[i]), 0) /
    (times.length - 1)
  return { latency: avg, jitter }
}

// Download: open several parallel streams, count bytes as they arrive.
export async function measureDownload({ streams = 4, sizeMB = 25, onProgress } = {}) {
  const size = sizeMB * 1024 * 1024
  let received = 0
  const start = performance.now()

  async function worker() {
    const res = await fetch(`/api/download?size=${size}&t=${Date.now()}`, {
      cache: 'no-store',
    })
    const reader = res.body.getReader()
    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      received += value.length
      if (onProgress) {
        const seconds = (performance.now() - start) / 1000
        onProgress((received * 8) / seconds / 1e6) // live Mbps
      }
    }
  }

  await Promise.all(Array.from({ length: streams }, worker))
  const seconds = (performance.now() - start) / 1000
  return (received * 8) / seconds / 1e6 // megabits per second
}

// Upload: POST random data in parallel and time it.
export async function measureUpload({ streams = 3, sizeMB = 8 } = {}) {
  const size = sizeMB * 1024 * 1024
  const data = new Uint8Array(size)
  // getRandomValues has a 64 KB limit per call, so fill in pieces
  for (let i = 0; i < size; i += 65536) {
    crypto.getRandomValues(data.subarray(i, Math.min(i + 65536, size)))
  }
  const start = performance.now()
  await Promise.all(
    Array.from({ length: streams }, () =>
      fetch('/api/upload', { method: 'POST', body: data })
    )
  )
  const seconds = (performance.now() - start) / 1000
  return (size * streams * 8) / seconds / 1e6
}