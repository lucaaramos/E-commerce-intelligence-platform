import { useState } from 'react'

interface Product {
  id: string
  sku: string
  name: string
  status: 'active' | 'inactive' | 'archived'
  sellingPrice: number
  margin: number
}

const mockProducts: Product[] = [
  { id: '1', sku: 'SKU-1001', name: 'Wireless Headphones', status: 'active', sellingPrice: 89.99, margin: 42.5 },
  { id: '2', sku: 'SKU-1002', name: 'Bluetooth Speaker', status: 'active', sellingPrice: 59.99, margin: 38.2 },
  { id: '3', sku: 'SKU-1003', name: 'Smart Watch', status: 'active', sellingPrice: 279.00, margin: 31.8 },
  { id: '4', sku: 'SKU-1004', name: 'USB-C Hub', status: 'inactive', sellingPrice: 45.00, margin: 55.1 },
  { id: '5', sku: 'SKU-1005', name: 'Running Shoes', status: 'active', sellingPrice: 129.99, margin: 28.4 },
]

const statusColors: Record<string, string> = {
  active: '#22c55e',
  inactive: '#f59e0b',
  archived: '#64748b',
}

export default function Products() {
  const [search, setSearch] = useState('')

  const filtered = mockProducts.filter(
    (p) => p.name.toLowerCase().includes(search.toLowerCase()) || p.sku.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0, fontSize: 24, color: '#f1f5f9' }}>Products</h1>
        <button style={{
          background: '#3b82f6',
          color: '#fff',
          border: 'none',
          borderRadius: 8,
          padding: '0.5rem 1rem',
          fontSize: 14,
          cursor: 'pointer',
        }}>
          + Add Product
        </button>
      </div>

      <input
        type="text"
        placeholder="Search by name or SKU..."
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
              {['SKU', 'Name', 'Status', 'Price', 'Margin', 'Actions'].map((h) => (
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
            {filtered.map((p) => (
              <tr key={p.id} style={{ borderBottom: '1px solid #334155' }}>
                <td style={{ padding: '0.75rem 1rem', fontSize: 14, color: '#94a3b8', fontFamily: 'monospace' }}>{p.sku}</td>
                <td style={{ padding: '0.75rem 1rem', fontSize: 14, color: '#f1f5f9' }}>{p.name}</td>
                <td style={{ padding: '0.75rem 1rem' }}>
                  <span style={{
                    padding: '0.25rem 0.5rem',
                    borderRadius: 4,
                    fontSize: 12,
                    color: statusColors[p.status],
                    background: `${statusColors[p.status]}20`,
                  }}>
                    {p.status}
                  </span>
                </td>
                <td style={{ padding: '0.75rem 1rem', fontSize: 14, color: '#f1f5f9' }}>${p.sellingPrice.toFixed(2)}</td>
                <td style={{ padding: '0.75rem 1rem', fontSize: 14, color: '#22c55e' }}>{p.margin}%</td>
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
