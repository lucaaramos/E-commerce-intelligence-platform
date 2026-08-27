import {
  AreaChart, Area, BarChart, Bar,
  XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
} from 'recharts'

const revenueData = [
  { month: 'Jan', revenue: 12400, cost: 8200 },
  { month: 'Feb', revenue: 15800, cost: 9600 },
  { month: 'Mar', revenue: 14200, cost: 8900 },
  { month: 'Apr', revenue: 18600, cost: 11200 },
  { month: 'May', revenue: 21300, cost: 12800 },
  { month: 'Jun', revenue: 19800, cost: 11900 },
]

const stockData = [
  { category: 'Electronics', stock: 342 },
  { category: 'Clothing', stock: 567 },
  { category: 'Home', stock: 234 },
  { category: 'Sports', stock: 189 },
  { category: 'Books', stock: 421 },
]

const tooltipStyle = {
  contentStyle: { background: '#1e293b', border: '1px solid #334155', borderRadius: 8 },
  labelStyle: { color: '#94a3b8' },
}

export function RevenueChart() {
  return (
    <div style={{
      background: '#1e293b',
      border: '1px solid #334155',
      borderRadius: 12,
      padding: '1.5rem',
    }}>
      <h3 style={{ margin: '0 0 1rem', fontSize: 16, color: '#f1f5f9' }}>Revenue vs Cost</h3>
      <ResponsiveContainer width="100%" height={280}>
        <AreaChart data={revenueData}>
          <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
          <XAxis dataKey="month" stroke="#64748b" fontSize={12} />
          <YAxis stroke="#64748b" fontSize={12} />
          <Tooltip {...tooltipStyle} />
          <Area type="monotone" dataKey="revenue" stroke="#3b82f6" fill="#3b82f630" name="Revenue" />
          <Area type="monotone" dataKey="cost" stroke="#ef4444" fill="#ef444430" name="Cost" />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  )
}

export function StockChart() {
  return (
    <div style={{
      background: '#1e293b',
      border: '1px solid #334155',
      borderRadius: 12,
      padding: '1.5rem',
    }}>
      <h3 style={{ margin: '0 0 1rem', fontSize: 16, color: '#f1f5f9' }}>Stock by Category</h3>
      <ResponsiveContainer width="100%" height={280}>
        <BarChart data={stockData}>
          <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
          <XAxis dataKey="category" stroke="#64748b" fontSize={12} />
          <YAxis stroke="#64748b" fontSize={12} />
          <Tooltip {...tooltipStyle} />
          <Bar dataKey="stock" fill="#22c55e" radius={[6, 6, 0, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
