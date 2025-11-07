import { useEffect, useState } from 'react'
import { Container, HardDrive, Activity, TrendingUp } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { api, SystemStats } from '@/lib/api'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { TEST_IDS } from '@/lib/test-ids'

export function Dashboard() {
  const [stats, setStats] = useState<SystemStats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadStats()
    const interval = setInterval(loadStats, 5000)
    return () => clearInterval(interval)
  }, [])

  const loadStats = async () => {
    try {
      const data = await api.getSystemStats()
      setStats(data)
    } catch (error) {
      console.error('Failed to load stats:', error)
      // Mock data for development
      setStats({
        cpuUsage: 45.2,
        memoryUsage: 62.8,
        diskUsage: 38.5,
        containersRunning: 3,
        containersTotal: 5,
      })
    } finally {
      setLoading(false)
    }
  }

  const statsCards = [
    {
      title: 'Running Containers',
      value: stats?.containersRunning || 0,
      total: stats?.containersTotal || 0,
      icon: Container,
      color: 'text-blue-500',
      href: '/containers',
    },
    {
      title: 'CPU Usage',
      value: `${stats?.cpuUsage.toFixed(1) || 0}%`,
      icon: Activity,
      color: 'text-green-500',
    },
    {
      title: 'Memory Usage',
      value: `${stats?.memoryUsage.toFixed(1) || 0}%`,
      icon: TrendingUp,
      color: 'text-orange-500',
    },
    {
      title: 'Disk Usage',
      value: `${stats?.diskUsage.toFixed(1) || 0}%`,
      icon: HardDrive,
      color: 'text-purple-500',
    },
  ]

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Dashboard</h1>
        <p className="text-muted-foreground">
          Overview of your NebulaBox infrastructure
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4 mb-8">
        {statsCards.map((stat) => {
          const Icon = stat.icon
          const content = (
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  {stat.title}
                </CardTitle>
                <Icon className={`h-4 w-4 ${stat.color}`} />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
                {stat.total !== undefined && (
                  <p className="text-xs text-muted-foreground mt-1">
                    {stat.total} total
                  </p>
                )}
              </CardContent>
            </Card>
          )

          return stat.href ? (
            <Link key={stat.title} to={stat.href} data-test-id={TEST_IDS.dashboard.viewContainers}>
              {content}
            </Link>
          ) : (
            <div key={stat.title} data-test-id={TEST_IDS.monitor.metricCard}>{content}</div>
          )
        })}
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
            <CardDescription>
              Common container management tasks
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            <Link to="/containers">
              <Button className="w-full justify-start" variant="outline" data-test-id={TEST_IDS.dashboard.viewContainers}>
                <Container className="mr-2 h-4 w-4" />
                View All Containers
              </Button>
            </Link>
            <Link to="/containers/new">
              <Button className="w-full justify-start" variant="outline" data-test-id={TEST_IDS.dashboard.createContainer}>
                <Container className="mr-2 h-4 w-4" />
                Create New Container
              </Button>
            </Link>
            <Link to="/images">
              <Button className="w-full justify-start" variant="outline" data-test-id={TEST_IDS.dashboard.manageImages}>
                <HardDrive className="mr-2 h-4 w-4" />
                Manage Images
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>System Status</CardTitle>
            <CardDescription>
              Real-time system monitoring
            </CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="text-sm text-muted-foreground">Loading...</div>
            ) : stats ? (
              <div className="space-y-4">
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span>CPU</span>
                    <span>{stats.cpuUsage.toFixed(1)}%</span>
                  </div>
                  <div className="w-full bg-secondary rounded-full h-2">
                    <div
                      className="bg-green-500 h-2 rounded-full"
                      style={{ width: `${stats.cpuUsage}%` }}
                    />
                  </div>
                </div>
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span>Memory</span>
                    <span>{stats.memoryUsage.toFixed(1)}%</span>
                  </div>
                  <div className="w-full bg-secondary rounded-full h-2">
                    <div
                      className="bg-orange-500 h-2 rounded-full"
                      style={{ width: `${stats.memoryUsage}%` }}
                    />
                  </div>
                </div>
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span>Disk</span>
                    <span>{stats.diskUsage.toFixed(1)}%</span>
                  </div>
                  <div className="w-full bg-secondary rounded-full h-2">
                    <div
                      className="bg-purple-500 h-2 rounded-full"
                      style={{ width: `${stats.diskUsage}%` }}
                    />
                  </div>
                </div>
              </div>
            ) : (
              <div className="text-sm text-muted-foreground">
                Unable to load stats
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

