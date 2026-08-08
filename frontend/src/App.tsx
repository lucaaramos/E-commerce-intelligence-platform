import { useEffect, useState } from 'react'

type ReadyState = { status: string; checks?: Record<string, { status: string }> }

function App() {
  const [ready, setReady] = useState<ReadyState | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetch('/readyz')
      .then((res) => res.json())
      .then(setReady)
      .catch((err) => setError(String(err)))
  }, [])

  return (
    <main style={{ fontFamily: 'system-ui, sans-serif', maxWidth: 640, margin: '4rem auto', padding: '0 1rem' }}>
      <h1>E-Commerce Intelligence Platform</h1>
      <p>Bootstrap complete. Dashboard and product management ship in later iterations.</p>
      {error && <p style={{ color: 'red' }}>Cannot reach backend: {error}</p>}
      {ready && (
        <pre>{JSON.stringify(ready, null, 2)}</pre>
      )}
    </main>
  )
}

export default App
