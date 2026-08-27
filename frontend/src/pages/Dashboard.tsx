import KPICard from '../components/KPICard'
import { RevenueChart, StockChart } from '../components/Charts'

const recentActivity = [
  { id: 1, action: 'Product created', detail: 'SKU-1001 Wireless Headphones', time: '2 min ago' },
  { id: 2, action: 'Supplier updated', detail: 'TechCorp International', time: '15 min ago' },
  { id: 3, action: 'Stock alert', detail: 'SKU-1045 Running Shoes - low stock (12)', time: '1 hr ago' },
  { id: 4, action: 'Price changed', detail: 'SKU-1023 Smart Watch $299 → $279', time: '2 hr ago' },
  { id: 5, action: 'Product archived', detail: 'SKU-1099 Old Model Speaker', time: '3 hr ago' },
]

export default function Dashboard() {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
      <h1 style={{ margin: 0, fontSize: 24, color: '#f1f5f9' }}>Dashboard</h1>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '1rem' }}>
        <KPICard title="Total Products" value={128} subtitle="12 active, 3 archived" color="#3b82f6" />
        <KPICard title="Suppliers" value={24} subtitle="8 countries" color="#8b5cf6" />
        <KPICard title="Avg Margin" value="34.2%" trend={{ value: 2.1, positive: true }} color="#22c55e" />
        <KPICard title="Revenue (MTD)" value="$21,300" trend={{ value: 8.5, positive: true }} color="#f59e0b" />
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '1rem' }}>
        <RevenueChart />
        <StockChart />
      </div>

      <div style={{
        background: '#1e293b',
        border: '1px solid #334155',
        borderRadius: 12,
        padding: '1.5rem',
      }}>
        <h3 style={{ margin: '0 0 1rem', fontSize: 16, color: '#f1f5f9' }}>Recent Activity</h3>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
          {recentActivity.map((item) => (
            <div key={item.id} style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              padding: '0.75rem 1rem',
              background: '#0f172a',
              borderRadius: 8,
            }}>
              <div>
                <span style={{ fontSize: 14, color: '#f1f5f9' }}>{item.action}</span>
                <span style={{ fontSize: 14, color: '#64748b', marginLeft: '0.75rem' }}>{item.detail}</span>
              </div>
              <span style={{ fontSize: 12, color: '#64748b', whiteSpace: 'nowrap' }}>{item.time}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
