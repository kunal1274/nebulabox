import { useEffect, useState } from 'react'
import { Settings, CheckCircle2, AlertCircle, Zap } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { api } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

type Mode = 'mock' | 'test' | 'live'

interface ModeInfo {
  mode: Mode
  description: string
  builtImages: number
  totalImages: number
}

export function ModeSwitcher() {
  const [modeInfo, setModeInfo] = useState<ModeInfo | null>(null)
  const [loading, setLoading] = useState(true)
  const [changing, setChanging] = useState(false)

  useEffect(() => {
    loadMode()
    // Refresh mode info every 5 seconds
    const interval = setInterval(loadMode, 5000)
    return () => clearInterval(interval)
  }, [])

  const loadMode = async () => {
    try {
      const info = await api.getMode()
      setModeInfo({
        mode: info.mode as Mode,
        description: info.description,
        builtImages: info.builtImages,
        totalImages: info.totalImages,
      })
    } catch (error) {
      console.error('Failed to load mode:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleModeChange = async (newMode: Mode) => {
    if (modeInfo?.mode === newMode) return
    
    setChanging(true)
    try {
      await api.setMode(newMode)
      await loadMode() // Refresh mode info
      // Reload page to ensure all components use new mode
      window.location.reload()
    } catch (error) {
      console.error('Failed to change mode:', error)
      alert(`Failed to change mode: ${error}`)
    } finally {
      setChanging(false)
    }
  }

  const getModeIcon = (mode: Mode) => {
    switch (mode) {
      case 'mock':
        return <AlertCircle className="h-4 w-4" />
      case 'test':
        return <CheckCircle2 className="h-4 w-4" />
      case 'live':
        return <Zap className="h-4 w-4" />
    }
  }

  const getModeColor = (mode: Mode) => {
    switch (mode) {
      case 'mock':
        return 'bg-gray-500'
      case 'test':
        return 'bg-blue-500'
      case 'live':
        return 'bg-green-500'
    }
  }

  if (loading || !modeInfo) {
    return (
      <div className="flex items-center gap-2">
        <Button variant="outline" size="sm" disabled>
          <Settings className="mr-2 h-4 w-4" />
          Loading...
        </Button>
      </div>
    )
  }

  return (
    <div className="flex items-center gap-2">
      <div className="flex items-center gap-2 px-3 py-1.5 rounded-md border bg-card" data-test-id={TEST_IDS.modeSwitcher.badge}>
        {getModeIcon(modeInfo.mode)}
        <span className="text-sm font-medium">Mode:</span>
        <Badge className={`${getModeColor(modeInfo.mode)} text-white`}>
          {modeInfo.mode.toUpperCase()}
        </Badge>
      </div>
      <div className="flex gap-1">
        <Button
          variant={modeInfo.mode === 'mock' ? 'default' : 'outline'}
          size="sm"
          onClick={() => handleModeChange('mock')}
          disabled={changing}
          title="Mock Mode - Static mock data only"
          data-test-id={TEST_IDS.modeSwitcher.mock}
        >
          Mock
        </Button>
        <Button
          variant={modeInfo.mode === 'test' ? 'default' : 'outline'}
          size="sm"
          onClick={() => handleModeChange('test')}
          disabled={changing}
          title="Test Mode - UAT sandbox with in-memory storage"
          data-test-id={TEST_IDS.modeSwitcher.test}
        >
          Test
        </Button>
        <Button
          variant={modeInfo.mode === 'live' ? 'default' : 'outline'}
          size="sm"
          onClick={() => handleModeChange('live')}
          disabled={changing}
          title="Live Mode - Production with full persistence"
          data-test-id={TEST_IDS.modeSwitcher.live}
        >
          Live
        </Button>
      </div>
      <div className="text-xs text-muted-foreground">
        Images: {modeInfo.builtImages}/{modeInfo.totalImages}
      </div>
    </div>
  )
}

