import { NavLink, Outlet } from 'react-router-dom'

const links = [
  { to: '/', label: 'Dashboard' },
  { to: '/products', label: 'Products' },
  { to: '/suppliers', label: 'Suppliers' },
]

export default function Layout() {
  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      <aside style={{
        width: 240,
        background: '#1e293b',
        borderRight: '1px solid #334155',
        padding: '1.5rem 0',
        display: 'flex',
        flexDirection: 'column',
      }}>
        <h2 style={{ padding: '0 1.5rem', margin: 0, fontSize: 18, color: '#f1f5f9' }}>
          E-Commerce Intel
        </h2>
        <nav style={{ marginTop: '2rem', display: 'flex', flexDirection: 'column', gap: 2 }}>
          {links.map(({ to, label }) => (
            <NavLink
              key={to}
              to={to}
              style={({ isActive }) => ({
                display: 'block',
                padding: '0.75rem 1.5rem',
                color: isActive ? '#f1f5f9' : '#94a3b8',
                background: isActive ? '#334155' : 'transparent',
                textDecoration: 'none',
                fontSize: 14,
              })}
            >
              {label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <main style={{ flex: 1, padding: '2rem', overflow: 'auto' }}>
        <Outlet />
      </main>
    </div>
  )
}
