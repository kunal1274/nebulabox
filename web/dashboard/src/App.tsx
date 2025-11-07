import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { DashboardLayout } from '@/components/layout/DashboardLayout'
import { Dashboard } from '@/pages/Dashboard'
import { Containers } from '@/pages/Containers'
import { CreateContainer } from '@/pages/CreateContainer'
import { ContainerLogs } from '@/pages/ContainerLogs'
import { Images } from '@/pages/Images'
import { Registry } from '@/pages/Registry'
import { BuildSpec } from '@/pages/BuildSpec'
import { Security } from '@/pages/Security'
import { Orchestrator } from '@/pages/Orchestrator'
import { Runtime } from '@/pages/Runtime'
import { AIOps } from '@/pages/AIOps'
import { ContainerGroups } from '@/pages/ContainerGroups'
import { Composition } from '@/pages/Composition'
import { Templates } from '@/pages/Templates'
import { SharedRuntime } from '@/pages/SharedRuntime'
import { Snapshots } from '@/pages/Snapshots'
import { EphemeralRuntime } from '@/pages/EphemeralRuntime'
import { Monitor } from '@/pages/Monitor'
import { Settings } from '@/pages/Settings'
import { Networks } from '@/pages/Networks'
import { Services } from '@/pages/Services'
import { Teams } from '@/pages/Teams'
import { Tenants } from '@/pages/Tenants'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<DashboardLayout />}>
          <Route index element={<Dashboard />} />
          <Route path="containers" element={<Containers />} />
          <Route path="containers/new" element={<CreateContainer />} />
          <Route path="containers/:id/logs" element={<ContainerLogs />} />
          <Route path="images" element={<Images />} />
          <Route path="registry" element={<Registry />} />
          <Route path="buildspec" element={<BuildSpec />} />
          <Route path="security" element={<Security />} />
          <Route path="orchestrator" element={<Orchestrator />} />
          <Route path="runtime" element={<Runtime />} />
          <Route path="aiops" element={<AIOps />} />
          <Route path="groups" element={<ContainerGroups />} />
          <Route path="composition" element={<Composition />} />
          <Route path="templates" element={<Templates />} />
          <Route path="shareruntime" element={<SharedRuntime />} />
          <Route path="snapshots" element={<Snapshots />} />
          <Route path="ephemeral" element={<EphemeralRuntime />} />
          <Route path="monitor" element={<Monitor />} />
          <Route path="networks" element={<Networks />} />
          <Route path="services" element={<Services />} />
          <Route path="teams" element={<Teams />} />
          <Route path="tenants" element={<Tenants />} />
          <Route path="settings" element={<Settings />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App

