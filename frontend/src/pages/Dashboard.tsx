import { useEffect, useState } from 'react'
import KPICard from '../components/KPICard'
import { MarginChart, PriceChart, StockChart, PriceHistoryChart } from '../components/Charts'
import { api, type Product, type Supplier, type PriceHistory } from '../lib/api'

export default function Dashboard() {
  const [products, setProducts] = useState<Product[]>([])
  const [suppliers, setSuppliers] = useState<Supplier[]>([])
  const [priceHistory, setPriceHistory] = useState<PriceHistory[]>([])
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    Promise.all([
      api.products.list({ limit: 100 }),
      api.suppliers.list({ limit: 100 }),
    ])
      .then(([prodRes, suppRes]) => {
        const prods = prodRes.data.products || []
        setProducts(prods)
        setSuppliers(suppRes.data.suppliers || [])
        if (prods.length > 0) {
          setSelectedProduct(prods[0])
        }
      })
      .catch((err) => setError(String(err)))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    if (!selectedProduct) return
    api.priceHistory.byProduct(selectedProduct.id)
      .then((res) => setPriceHistory(res.data || []))
      .catch(() => setPriceHistory([]))
  }, [selectedProduct])

  const activeProducts = products.filter((p) => p.status === 'active').length
  const avgMargin = products.length > 0
    ? (products.reduce((sum, p) => sum + p.margin, 0) / products.length).toFixed(1)
    : '0'

  if (loading) return <p style={{ color: '#94a3b8' }}>Loading...</p>
  if (error) return <p style={{ color: '#ef4444' }}>Error: {error}</p>

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
      <h1 style={{ margin: 0, fontSize: 24, color: '#f1f5f9' }}>Dashboard</h1>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '1rem' }}>
        <KPICard title="Total Products" value={products.length} subtitle={`${activeProducts} active`} color="#3b82f6" />
        <KPICard title="Suppliers" value={suppliers.length} color="#8b5cf6" />
        <KPICard title="Avg Margin" value={`${avgMargin}%`} color="#22c55e" />
        <KPICard title="Products w/ Variants" value={products.filter((p) => p.variants && p.variants.length > 0).length} color="#f59e0b" />
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
        <MarginChart products={products} />
        <PriceChart products={products} />
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
        <StockChart products={products} />
        <div style={{
          background: '#1e293b',
          border: '1px solid #334155',
          borderRadius: 12,
          padding: '1.5rem',
        }}>
          <h3 style={{ margin: '0 0 1rem', fontSize: 16, color: '#f1f5f9' }}>Price Distribution</h3>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            {[
              { range: '$0 - $50', count: products.filter((p) => p.sellingPrice <= 50).length },
              { range: '$50 - $100', count: products.filter((p) => p.sellingPrice > 50 && p.sellingPrice <= 100).length },
              { range: '$100 - $200', count: products.filter((p) => p.sellingPrice > 100 && p.sellingPrice <= 200).length },
              { range: '$200+', count: products.filter((p) => p.sellingPrice > 200).length },
            ].map(({ range, count }) => (
              <div key={range} style={{ display: 'flex', justifyContent: 'space-between', padding: '0.5rem 0', borderBottom: '1px solid #334155' }}>
                <span style={{ fontSize: 14, color: '#94a3b8' }}>{range}</span>
                <span style={{ fontSize: 14, color: '#f1f5f9', fontWeight: 600 }}>{count}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '200px 1fr', gap: '1rem' }}>
        <div style={{
          background: '#1e293b',
          border: '1px solid #334155',
          borderRadius: 12,
          padding: '1.5rem',
        }}>
          <h3 style={{ margin: '0 0 1rem', fontSize: 14, color: '#94a3b8' }}>Select Product</h3>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            {products.map((p) => (
              <button
                key={p.id}
                onClick={() => setSelectedProduct(p)}
                style={{
                  textAlign: 'left',
                  padding: '0.5rem',
                  borderRadius: 6,
                  border: selectedProduct?.id === p.id ? '1px solid #3b82f6' : '1px solid transparent',
                  background: selectedProduct?.id === p.id ? '#3b82f620' : 'transparent',
                  color: '#f1f5f9',
                  fontSize: 13,
                  cursor: 'pointer',
                }}
              >
                {p.sku}
              </button>
            ))}
          </div>
        </div>
        <PriceHistoryChart
          history={priceHistory}
          productName={selectedProduct ? `${selectedProduct.sku} - ${selectedProduct.name}` : ''}
        />
      </div>

      <div style={{
        background: '#1e293b',
        border: '1px solid #334155',
        borderRadius: 12,
        padding: '1.5rem',
      }}>
        <h3 style={{ margin: '0 0 1rem', fontSize: 16, color: '#f1f5f9' }}>Recent Products</h3>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
          {products.slice(0, 5).map((p) => (
            <div key={p.id} style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              padding: '0.75rem 1rem',
              background: '#0f172a',
              borderRadius: 8,
            }}>
              <div>
                <span style={{ fontSize: 14, color: '#f1f5f9' }}>{p.name}</span>
                <span style={{ fontSize: 14, color: '#64748b', marginLeft: '0.75rem', fontFamily: 'monospace' }}>{p.sku}</span>
              </div>
              <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
                <span style={{ fontSize: 14, color: '#f1f5f9' }}>${p.sellingPrice.toFixed(2)}</span>
                <span style={{ fontSize: 12, color: '#22c55e' }}>{p.margin.toFixed(1)}%</span>
              </div>
            </div>
          ))}
          {products.length === 0 && (
            <p style={{ color: '#64748b', fontSize: 14 }}>No products yet. Add one from the Products page.</p>
          )}
        </div>
      </div>
    </div>
  )
}
