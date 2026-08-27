import {
  BarChart, Bar,
  XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
} from 'recharts'
import type { Product } from '../lib/api'

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
