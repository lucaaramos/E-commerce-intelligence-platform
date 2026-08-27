import { useState } from 'react'

interface Supplier {
  id: string
  name: string
  country: string
  integrationMethod: string
  contactEmail: string
}

const mockSuppliers: Supplier[] = [
  { id: '1', name: 'TechCorp International', country: 'China', integrationMethod: 'API', contactEmail: 'sales@techcorp.com' },
  { id: '2', name: 'GlobalTextile Ltd', country: 'Vietnam', integrationMethod: 'CSV', contactEmail: 'orders@globaltextile.com' },
  { id: '3', name: 'SportGear Direct', country: 'USA', integrationMethod: 'API', contactEmail: 'wholesale@sportgear.com' },
  { id: '4', name: 'HomeBasics Co', country: 'Germany', integrationMethod: 'Manual', contactEmail: 'supply@homebasics.de' },
  { id: '5', name: 'DigitalWorld', country: 'Japan', integrationMethod: 'API', contactEmail: 'b2b@digitalworld.jp' },
]

const countryFlags: Record<string, string> = {
  China: '🇨🇳', Vietnam: '🇻🇳', USA: '🇺🇸', Germany: '🇩🇪', Japan: '🇯🇵',
}

export default function Suppliers() {
  const [search, setSearch] = useState('')

  const filtered = mockSuppliers.filter(
    (s) => s.name.toLowerCase().includes(search.toLowerCase()) || s.country.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0, fontSize: 24, color: '#f1f5f9' }}>Suppliers</h1>
        <button style={{
          background: '#3b82f6',
          color: '#fff',
          border: 'none',
          borderRadius: 8,
          padding: '0.5rem 1rem',
          fontSize: 14,
          cursor: 'pointer',
        }}>
          + Add Supplier
        </button>
      </div>

      <input
        type="text"
        placeholder="Search by name or country..."
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        style={{
          background: '#1e293b',
          border: '1px solid #334155',
          borderRadius: 8,
          padding: '0.75rem 1rem',
          color: '#f1f5f9',
          fontSize: 14,
          outline: 'none',
        }}
      />

      <div style={{
        background: '#1e293b',
        border: '1px solid #334155',
        borderRadius: 12,
        overflow: 'hidden',
      }}>
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr style={{ borderBottom: '1px solid #334155' }}>
              {['Name', 'Country', 'Integration', 'Contact', 'Actions'].map((h) => (
                <th key={h} style={{
                  padding: '0.75rem 1rem',
                  textAlign: 'left',
                  fontSize: 12,
                  color: '#94a3b8',
                  fontWeight: 600,
                  textTransform: 'uppercase',
                }}>{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {filtered.map((s) => (
              <tr key={s.id} style={{ borderBottom: '1px solid #334155' }}>
                <td style={{ padding: '0.75rem 1rem', fontSize: 14, color: '#f1f5f9', fontWeight: 500 }}>{s.name}</td>
                <td style={{ padding: '0.75rem 1rem', fontSize: 14, color: '#f1f5f9' }}>
                  {countryFlags[s.country] || ''} {s.country}
                </td>
                <td style={{ padding: '0.75rem 1rem' }}>
                  <span style={{
                    padding: '0.25rem 0.5rem',
                    borderRadius: 4,
                    fontSize: 12,
                    color: s.integrationMethod === 'API' ? '#3b82f6' : '#f59e0b',
                    background: s.integrationMethod === 'API' ? '#3b82f620' : '#f59e0b20',
                  }}>
                    {s.integrationMethod}
                  </span>
                </td>
                <td style={{ padding: '0.75rem 1rem', fontSize: 14, color: '#94a3b8' }}>{s.contactEmail}</td>
                <td style={{ padding: '0.75rem 1rem' }}>
                  <button style={{
                    background: 'transparent',
                    border: '1px solid #334155',
                    borderRadius: 4,
                    padding: '0.25rem 0.5rem',
                    color: '#94a3b8',
                    fontSize: 12,
                    cursor: 'pointer',
                  }}>Edit</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
