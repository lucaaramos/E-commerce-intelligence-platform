import { useEffect, useState } from 'react'
import { api, type Supplier } from '../lib/api'

export default function Suppliers() {
  const [suppliers, setSuppliers] = useState<Supplier[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [editingSupplier, setEditingSupplier] = useState<Supplier | null>(null)
  const [form, setForm] = useState({ name: '', country: '', integrationMethod: '', contactEmail: '' })
  const [error, setError] = useState<string | null>(null)

  const load = () => {
    setLoading(true)
    api.suppliers.list({ limit: 100 })
      .then((res) => setSuppliers(res.data.suppliers || []))
      .catch((err) => setError(String(err)))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const filtered = suppliers.filter(
    (s) => s.name.toLowerCase().includes(search.toLowerCase()) || s.country.toLowerCase().includes(search.toLowerCase())
  )

  const openCreate = () => {
    setEditingSupplier(null)
    setForm({ name: '', country: '', integrationMethod: '', contactEmail: '' })
    setShowForm(true)
    setError(null)
  }

  const openEdit = (s: Supplier) => {
    setEditingSupplier(s)
    setForm({ name: s.name, country: s.country, integrationMethod: s.integrationMethod, contactEmail: s.contactEmail })
    setShowForm(true)
    setError(null)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    try {
      if (editingSupplier) {
        await api.suppliers.update(editingSupplier.id, form)
      } else {
        await api.suppliers.create(form)
      }
      setShowForm(false)
      load()
    } catch (err: any) {
      setError(err.message || 'Failed to save supplier')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this supplier?')) return
    try {
      await api.suppliers.delete(id)
      load()
    } catch (err: any) {
      setError(err.message)
    }
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0, fontSize: 24, color: '#f1f5f9' }}>Suppliers</h1>
        <button onClick={openCreate} style={{
          background: '#3b82f6', color: '#fff', border: 'none',
          borderRadius: 8, padding: '0.5rem 1rem', fontSize: 14, cursor: 'pointer',
        }}>+ Add Supplier</button>
      </div>

      {error && <p style={{ color: '#ef4444', margin: 0 }}>{error}</p>}

      {showForm && (
        <form onSubmit={handleSubmit} style={{
          background: '#1e293b', border: '1px solid #334155', borderRadius: 12,
          padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '1rem',
        }}>
          <h3 style={{ margin: 0, color: '#f1f5f9' }}>{editingSupplier ? 'Edit Supplier' : 'New Supplier'}</h3>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
            <input placeholder="Name" value={form.name} required
              onChange={(e) => setForm({ ...form, name: e.target.value })} style={inputStyle} />
            <input placeholder="Country" value={form.country}
              onChange={(e) => setForm({ ...form, country: e.target.value })} style={inputStyle} />
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
            <input placeholder="Integration Method (API, CSV, Manual)" value={form.integrationMethod}
              onChange={(e) => setForm({ ...form, integrationMethod: e.target.value })} style={inputStyle} />
            <input placeholder="Contact Email" type="email" value={form.contactEmail}
              onChange={(e) => setForm({ ...form, contactEmail: e.target.value })} style={inputStyle} />
          </div>
          <div style={{ display: 'flex', gap: '0.75rem', justifyContent: 'flex-end' }}>
            <button type="button" onClick={() => setShowForm(false)} style={btnSecondaryStyle}>Cancel</button>
            <button type="submit" style={btnPrimaryStyle}>{editingSupplier ? 'Update' : 'Create'}</button>
          </div>
        </form>
      )}

      <input type="text" placeholder="Search by name or country..." value={search}
        onChange={(e) => setSearch(e.target.value)} style={inputStyle} />

      {loading ? <p style={{ color: '#94a3b8' }}>Loading...</p> : (
        <div style={{ background: '#1e293b', border: '1px solid #334155', borderRadius: 12, overflow: 'hidden' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid #334155' }}>
                {['Name', 'Country', 'Integration', 'Contact', 'Actions'].map((h) => (
                  <th key={h} style={thStyle}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {filtered.map((s) => (
                <tr key={s.id} style={{ borderBottom: '1px solid #334155' }}>
                  <td style={{ ...tdStyle, fontWeight: 500 }}>{s.name}</td>
                  <td style={tdStyle}>{s.country || '—'}</td>
                  <td style={{ padding: '0.75rem 1rem' }}>
                    <span style={{
                      padding: '0.25rem 0.5rem', borderRadius: 4, fontSize: 12,
                      color: s.integrationMethod === 'API' ? '#3b82f6' : '#f59e0b',
                      background: s.integrationMethod === 'API' ? '#3b82f620' : '#f59e0b20',
                    }}>{s.integrationMethod || '—'}</span>
                  </td>
                  <td style={{ ...tdStyle, color: '#94a3b8' }}>{s.contactEmail || '—'}</td>
                  <td style={{ padding: '0.75rem 1rem', display: 'flex', gap: '0.5rem' }}>
                    <button onClick={() => openEdit(s)} style={btnSmall}>Edit</button>
                    <button onClick={() => handleDelete(s.id)} style={{ ...btnSmall, color: '#ef4444' }}>Delete</button>
                  </td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr><td colSpan={5} style={{ ...tdStyle, textAlign: 'center', color: '#64748b' }}>No suppliers found</td></tr>
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
