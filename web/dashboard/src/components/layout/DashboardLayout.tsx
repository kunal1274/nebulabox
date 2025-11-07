import { Outlet } from 'react-router-dom'
import { Sidebar } from './Sidebar'
import { ModeSwitcher } from '@/components/ModeSwitcher'

export function DashboardLayout() {
  return (
    <div className="flex h-screen bg-background">
      <Sidebar />
      <main className="flex-1 flex flex-col overflow-hidden">
        <div className="border-b bg-card px-6 py-3 flex items-center justify-end shrink-0">
          <ModeSwitcher />
        </div>
        <div className="flex-1 overflow-y-auto">
          <Outlet />
        </div>
      </main>
    </div>
  )
}

