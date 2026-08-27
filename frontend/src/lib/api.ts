const API_BASE = '/api/v1'

interface ApiResponse<T> {
  data: T
}

interface PaginatedResponse<T> {
  data: {
    products?: T[]
    suppliers?: T[]
    total: number
    limit: number
    offset: number
  }
}

export interface Product {
  id: string
  sku: string
  name: string
  description: string
  status: 'active' | 'inactive' | 'archived'
  sellingPrice: number
  realCost: number
  margin: number
  variants: { sku: string; name: string; cost: number; stock: number; position: number }[]
  createdAt: string
  updatedAt: string
}

export interface Supplier {
  id: string
  name: string
  country: string
  integrationMethod: string
  contactEmail: string
  createdAt: string
  updatedAt: string
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => null)
    throw new Error(err?.error?.message || `Request failed: ${res.status}`)
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

export const api = {
  products: {
    list: (params?: { sku?: string; status?: string; limit?: number; offset?: number }) => {
      const qs = new URLSearchParams()
      if (params?.sku) qs.set('sku', params.sku)
      if (params?.status) qs.set('status', params.status)
      if (params?.limit) qs.set('limit', String(params.limit))
      if (params?.offset) qs.set('offset', String(params.offset))
      return request<PaginatedResponse<Product>>(`/products?${qs}`)
    },
    getByID: (id: string) => request<ApiResponse<Product>>(`/products/${id}`),
    getBySKU: (sku: string) => request<ApiResponse<Product>>(`/products/sku/${sku}`),
    create: (data: { sku: string; name: string; description?: string; sellingPrice: number }) =>
      request<ApiResponse<Product>>('/products', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: Partial<Pick<Product, 'name' | 'description' | 'status' | 'sellingPrice'>>) =>
      request<ApiResponse<Product>>(`/products/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id: string) =>
      request<void>(`/products/${id}`, { method: 'DELETE' }),
  },
  suppliers: {
    list: (params?: { name?: string; limit?: number; offset?: number }) => {
      const qs = new URLSearchParams()
      if (params?.name) qs.set('name', params.name)
      if (params?.limit) qs.set('limit', String(params.limit))
      if (params?.offset) qs.set('offset', String(params.offset))
      return request<PaginatedResponse<Supplier>>(`/suppliers?${qs}`)
    },
    getByID: (id: string) => request<ApiResponse<Supplier>>(`/suppliers/${id}`),
    create: (data: { name: string; country?: string; integrationMethod?: string; contactEmail?: string }) =>
      request<ApiResponse<Supplier>>('/suppliers', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: Partial<Pick<Supplier, 'name' | 'country' | 'integrationMethod' | 'contactEmail'>>) =>
      request<ApiResponse<Supplier>>(`/suppliers/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id: string) =>
      request<void>(`/suppliers/${id}`, { method: 'DELETE' }),
  },
}
