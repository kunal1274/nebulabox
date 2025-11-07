import { Link, useLocation } from 'react-router-dom'
import { 
  LayoutDashboard, 
  Container, 
  HardDrive, 
  Activity,
  Settings,
  Network,
  Users,
  Package,
  Wrench,
  Shield,
  Zap,
  Cpu,
  Brain,
  Layers,
  Shuffle,
  FileStack,
  Globe,
  Camera,
  Cloud
} from 'lucide-react'
import { cn } from '@/lib/utils'

const navigation = [
  { name: 'Dashboard', href: '/', icon: LayoutDashboard },
  { name: 'Containers', href: '/containers', icon: Container },
  { name: 'Images', href: '/images', icon: HardDrive },
  { name: 'Registry', href: '/registry', icon: Package },
  { name: 'Build Spec', href: '/buildspec', icon: Wrench },
  { name: 'Security', href: '/security', icon: Shield },
  { name: 'Orchestrator', href: '/orchestrator', icon: Zap },
  { name: 'Runtime', href: '/runtime', icon: Cpu },
  { name: 'AI Ops', href: '/aiops', icon: Brain },
  { name: 'Groups', href: '/groups', icon: Layers },
  { name: 'Composition', href: '/composition', icon: Shuffle },
  { name: 'Templates', href: '/templates', icon: FileStack },
  { name: 'Shared Runtime', href: '/shareruntime', icon: Globe },
  { name: 'Snapshots', href: '/snapshots', icon: Camera },
  { name: 'Ephemeral Runtime', href: '/ephemeral', icon: Cloud },
  { name: 'Monitor', href: '/monitor', icon: Activity },
  { name: 'Networks', href: '/networks', icon: Network },
  { name: 'Services', href: '/services', icon: Activity },
  { name: 'Teams', href: '/teams', icon: Users },
  { name: 'Tenants', href: '/tenants', icon: Users },
  { name: 'Settings', href: '/settings', icon: Settings },
]

export function Sidebar() {
  const location = useLocation()

  return (
    <div className="flex h-screen w-64 flex-col bg-card border-r">
      <div className="flex h-16 items-center border-b px-6">
        <h1 className="text-xl font-bold bg-gradient-to-r from-primary to-purple-600 bg-clip-text text-transparent">
          🚀 NebulaBox
        </h1>
      </div>
      <nav className="flex-1 space-y-1 px-3 py-4">
        {navigation.map((item) => {
          const isActive = location.pathname === item.href
          return (
            <Link
              key={item.name}
              to={item.href}
              className={cn(
                'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                isActive
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
              )}
            >
              <item.icon className="h-5 w-5" />
              {item.name}
            </Link>
          )
        })}
      </nav>
      <div className="border-t p-4">
        <div className="text-xs text-muted-foreground">
          <div>Version 0.1.0-alpha</div>
          <div className="mt-1">Phase 1 Development</div>
        </div>
      </div>
    </div>
  )
}

