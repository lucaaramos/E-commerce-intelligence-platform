import { useEffect, useState } from 'react'
import { api, type Product } from '../lib/api'

const statusColors: Record<string, string> = {
  active: '#22c55e',
  inactive: '#f59e0b',
  archived: '#64748b',
}

export default function Products() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [form, setForm] = useState({ sku: '', name: '', description: '', sellingPrice: '' })
  const [error, setError] = useState<string | null>(null)

  const load = () => {
    setLoading(true)
    api.products.list({ limit: 100 })
      .then((res) => setProducts(res.data.products || []))
      .catch((err) => setError(String(err)))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const filtered = products.filter(
    (p) => p.name.toLowerCase().includes(search.toLowerCase()) || p.sku.toLowerCase().includes(search.toLowerCase())
  )

  const openCreate = () => {
    setEditingProduct(null)
    setForm({ sku: '', name: '', description: '', sellingPrice: '' })
    setShowForm(true)
    setError(null)
  }

  const openEdit = (p: Product) => {
    setEditingProduct(p)
    setForm({ sku: p.sku, name: p.name, description: p.description, sellingPrice: String(p.sellingPrice) })
    setShowForm(true)
    setError(null)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    try {
      if (editingProduct) {
        await api.products.update(editingProduct.id, {
          name: form.name,
          description: form.description,
          sellingPrice: Number(form.sellingPrice),
        })
      } else {
        await api.products.create({
          sku: form.sku,
          name: form.name,
          description: form.description,
          sellingPrice: Number(form.sellingPrice),
        })
      }
      setShowForm(false)
      load()
    } catch (err: any) {
      setError(err.message || 'Failed to save product')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this product?')) return
    try {
      await api.products.delete(id)
      load()
    } catch (err: any) {
      setError(err.message)
    }
  }

  const handleStatusToggle = async (p: Product) => {
    const next = p.status === 'active' ? 'inactive' : 'active'
    try {
      await api.products.update(p.id, { status: next })
      load()
    } catch (err: any) {
      setError(err.message)
    }
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0, fontSize: 24, color: '#f1f5f9' }}>Products</h1>
        <button onClick={openCreate} style={{
          background: '#3b82f6', color: '#fff', border: 'none',
          borderRadius: 8, padding: '0.5rem 1rem', fontSize: 14, cursor: 'pointer',
        }}>+ Add Product</button>
      </div>

      {error && <p style={{ color: '#ef4444', margin: 0 }}>{error}</p>}

      {showForm && (
        <form onSubmit={handleSubmit} style={{
          background: '#1e293b', border: '1px solid #334155', borderRadius: 12,
          padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '1rem',
        }}>
          <h3 style={{ margin: 0, color: '#f1f5f9' }}>{editingProduct ? 'Edit Product' : 'New Product'}</h3>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
            <input placeholder="SKU" value={form.sku} disabled={!!editingProduct}
              onChange={(e) => setForm({ ...form, sku: e.target.value })} style={inputStyle} />
            <input placeholder="Name" value={form.name} required
              onChange={(e) => setForm({ ...form, name: e.target.value })} style={inputStyle} />
          </div>
          <input placeholder="Description" value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })} style={inputStyle} />
          <input placeholder="Selling Price" type="number" step="0.01" min="0.01" required value={form.sellingPrice}
            onChange={(e) => setForm({ ...form, sellingPrice: e.target.value })} style={inputStyle} />
          <div style={{ display: 'flex', gap: '0.75rem', justifyContent: 'flex-end' }}>
            <button type="button" onClick={() => setShowForm(false)} style={btnSecondaryStyle}>Cancel</button>
            <button type="submit" style={btnPrimaryStyle}>{editingProduct ? 'Update' : 'Create'}</button>
          </div>
        </form>
      )}

      <input type="text" placeholder="Search by name or SKU..." value={search}
        onChange={(e) => setSearch(e.target.value)} style={inputStyle} />

      {loading ? <p style={{ color: '#94a3b8' }}>Loading...</p> : (
        <div style={{ background: '#1e293b', border: '1px solid #334155', borderRadius: 12, overflow: 'hidden' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid #334155' }}>
                {['SKU', 'Name', 'Status', 'Price', 'Cost', 'Margin', 'Actions'].map((h) => (
                  <th key={h} style={thStyle}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {filtered.map((p) => (
                <tr key={p.id} style={{ borderBottom: '1px solid #334155' }}>
                  <td style={tdStyleMono}>{p.sku}</td>
                  <td style={tdStyle}>{p.name}</td>
                  <td style={{ padding: '0.75rem 1rem' }}>
                    <button onClick={() => handleStatusToggle(p)} style={{
                      padding: '0.25rem 0.5rem', borderRadius: 4, fontSize: 12, cursor: 'pointer',
                      border: 'none', color: statusColors[p.status], background: `${statusColors[p.status]}20`,
                    }}>{p.status}</button>
                  </td>
                  <td style={tdStyle}>${p.sellingPrice.toFixed(2)}</td>
                  <td style={tdStyle}>${p.realCost.toFixed(2)}</td>
                  <td style={{ padding: '0.75rem 1rem', fontSize: 14, color: p.margin >= 0 ? '#22c55e' : '#ef4444' }}>
                    {p.margin.toFixed(1)}%
                  </td>
                  <td style={{ padding: '0.75rem 1rem', display: 'flex', gap: '0.5rem' }}>
                    <button onClick={() => openEdit(p)} style={btnSmall}>Edit</button>
                    <button onClick={() => handleDelete(p.id)} style={{ ...btnSmall, color: '#ef4444' }}>Delete</button>
                  </td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr><td colSpan={7} style={{ ...tdStyle, textAlign: 'center', color: '#64748b' }}>No products found</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

const inputStyle: React.CSSProperties = {
  background: '#0f172a', border: '1px solid #334155', borderRadius: 8,
  padding: '0.75rem 1rem', color: '#f1f5f9', fontSize: 14, outline: 'none',
}
const thStyle: React.CSSProperties = {
  padding: '0.75rem 1rem', textAlign: 'left', fontSize: 12,
  color: '#94a3b8', fontWeight: 600, textTransform: 'uppercase' as const,
}
const tdStyle: React.CSSProperties = { padding: '0.75rem 1rem', fontSize: 14, color: '#f1f5f9' }
const tdStyleMono: React.CSSProperties = { ...tdStyle, fontFamily: 'monospace', color: '#94a3b8' }
const btnPrimaryStyle: React.CSSProperties = {
  background: '#3b82f6', color: '#fff', border: 'none',
  borderRadius: 8, padding: '0.5rem 1rem', fontSize: 14, cursor: 'pointer',
}
const btnSecondaryStyle: React.CSSProperties = {
  background: 'transparent', color: '#94a3b8', border: '1px solid #334155',
  borderRadius: 8, padding: '0.5rem 1rem', fontSize: 14, cursor: 'pointer',
}
const btnSmall: React.CSSProperties = {
  background: 'transparent', border: '1px solid #334155', borderRadius: 4,
  padding: '0.25rem 0.5rem', color: '#94a3b8', fontSize: 12, cursor: 'pointer',
}
