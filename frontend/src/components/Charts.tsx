import {
  BarChart, Bar, LineChart, Line,
  XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend,
} from 'recharts'
import type { Product, PriceHistory } from '../lib/api'

const tooltipStyle = {
  contentStyle: { background: '#1e293b', border: '1px solid #334155', borderRadius: 8 },
  labelStyle: { color: '#94a3b8' },
}

interface ChartProps {
  products: Product[]
}

export function MarginChart({ products }: ChartProps) {
  const data = products.slice(0, 10).map((p) => ({
    name: p.sku,
    margin: Number(p.margin.toFixed(1)),
    price: p.sellingPrice,
  }))

  return (
    <div style={{
      background: '#1e293b',
      border: '1px solid #334155',
      borderRadius: 12,
      padding: '1.5rem',
    }}>
      <h3 style={{ margin: '0 0 1rem', fontSize: 16, color: '#f1f5f9' }}>Margin by Product (Top 10)</h3>
      <ResponsiveContainer width="100%" height={280}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
          <XAxis dataKey="name" stroke="#64748b" fontSize={11} angle={-45} textAnchor="end" height={60} />
          <YAxis stroke="#64748b" fontSize={12} unit="%" />
          <Tooltip {...tooltipStyle} />
          <Bar dataKey="margin" fill="#22c55e" radius={[6, 6, 0, 0]} name="Margin %" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}

export function PriceChart({ products }: ChartProps) {
  const data = products.slice(0, 10).map((p) => ({
    name: p.sku,
    price: p.sellingPrice,
    cost: p.realCost,
  }))

  return (
    <div style={{
      background: '#1e293b',
      border: '1px solid #334155',
      borderRadius: 12,
      padding: '1.5rem',
    }}>
      <h3 style={{ margin: '0 0 1rem', fontSize: 16, color: '#f1f5f9' }}>Price vs Cost (Top 10)</h3>
      <ResponsiveContainer width="100%" height={280}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
          <XAxis dataKey="name" stroke="#64748b" fontSize={11} angle={-45} textAnchor="end" height={60} />
          <YAxis stroke="#64748b" fontSize={12} />
          <Tooltip {...tooltipStyle} />
          <Bar dataKey="price" fill="#3b82f6" radius={[6, 6, 0, 0]} name="Price" />
          <Bar dataKey="cost" fill="#ef4444" radius={[6, 6, 0, 0]} name="Cost" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}

export function StockChart({ products }: ChartProps) {
  const statusCounts = products.reduce(
    (acc, p) => {
      acc[p.status] = (acc[p.status] || 0) + 1
      return acc
    },
    {} as Record<string, number>,
  )

  const data = Object.entries(statusCounts).map(([status, count]) => ({
    status,
    count,
  }))

  return (
    <div style={{
      background: '#1e293b',
      border: '1px solid #334155',
      borderRadius: 12,
      padding: '1.5rem',
    }}>
      <h3 style={{ margin: '0 0 1rem', fontSize: 16, color: '#f1f5f9' }}>Products by Status</h3>
      <ResponsiveContainer width="100%" height={280}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
          <XAxis dataKey="status" stroke="#64748b" fontSize={12} />
          <YAxis stroke="#64748b" fontSize={12} />
          <Tooltip {...tooltipStyle} />
          <Bar dataKey="count" fill="#8b5cf6" radius={[6, 6, 0, 0]} name="Products" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}

interface PriceHistoryChartProps {
  history: PriceHistory[]
  productName: string
}

export function PriceHistoryChart({ history, productName }: PriceHistoryChartProps) {
  const data = [...history]
    .reverse()
    .map((h) => ({
      time: new Date(h.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      price: h.price,
      shipping: h.shippingCost,
      total: Number((h.price + h.shippingCost).toFixed(2)),
    }))

  return (
    <div style={{
      background: '#1e293b',
      border: '1px solid #334155',
      borderRadius: 12,
      padding: '1.5rem',
    }}>
      <h3 style={{ margin: '0 0 1rem', fontSize: 16, color: '#f1f5f9' }}>
        Price History: {productName}
      </h3>
      {data.length === 0 ? (
        <p style={{ color: '#64748b', fontSize: 14, textAlign: 'center', padding: '2rem' }}>
          No price history yet. Changes will appear here when the sync worker detects price updates.
        </p>
      ) : (
        <ResponsiveContainer width="100%" height={280}>
          <LineChart data={data}>
            <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
            <XAxis dataKey="time" stroke="#64748b" fontSize={11} />
            <YAxis stroke="#64748b" fontSize={12} />
            <Tooltip {...tooltipStyle} />
            <Legend />
            <Line type="monotone" dataKey="price" stroke="#3b82f6" strokeWidth={2} dot={{ r: 4 }} name="Cost" />
            <Line type="monotone" dataKey="shipping" stroke="#f59e0b" strokeWidth={2} dot={{ r: 4 }} name="Shipping" />
            <Line type="monotone" dataKey="total" stroke="#22c55e" strokeWidth={2} dot={{ r: 4 }} name="Total" />
          </LineChart>
        </ResponsiveContainer>
      )}
    </div>
  )
}
