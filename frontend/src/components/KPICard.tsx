interface Props {
  title: string
  value: string | number
  subtitle?: string
  trend?: { value: number; positive: boolean }
  color?: string
}

export default function KPICard({ title, value, subtitle, trend, color = '#3b82f6' }: Props) {
  return (
    <div style={{
      background: '#1e293b',
      border: '1px solid #334155',
      borderRadius: 12,
      padding: '1.5rem',
      display: 'flex',
      flexDirection: 'column',
      gap: '0.5rem',
    }}>
      <span style={{ fontSize: 13, color: '#94a3b8', textTransform: 'uppercase', letterSpacing: 0.5 }}>
        {title}
      </span>
      <div style={{ display: 'flex', alignItems: 'baseline', gap: '0.75rem' }}>
        <span style={{ fontSize: 28, fontWeight: 700, color: '#f1f5f9' }}>{value}</span>
        {trend && (
          <span style={{
            fontSize: 13,
            color: trend.positive ? '#22c55e' : '#ef4444',
          }}>
            {trend.positive ? '↑' : '↓'} {Math.abs(trend.value)}%
          </span>
        )}
      </div>
      {subtitle && <span style={{ fontSize: 12, color: '#64748b' }}>{subtitle}</span>}
      <div style={{ height: 3, borderRadius: 2, background: '#334155', marginTop: '0.5rem' }}>
        <div style={{ height: '100%', width: '65%', borderRadius: 2, background: color }} />
      </div>
    </div>
  )
}
