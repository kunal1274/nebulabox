# Test ID Coverage Report

**Generated:** $(date)
**Status:** In Progress (5/25 pages completed - 20%)

## ✅ Pages WITH Test IDs (5 pages)

| # | Page | File | Status | Test IDs Added |
|---|------|------|--------|----------------|
| 1 | Containers | `Containers.tsx` | ✅ Complete | refresh, create, list, card, stop, start, logs |
| 2 | Images | `Images.tsx` | ✅ Complete | refresh, pull, pullInput, build, buildTagInput, scan, push, list, card |
| 3 | Build Spec | `BuildSpec.tsx` | ✅ Complete | loadExample, tagInput, specEditor, validate, convert, build, dockerfileTab, validationTab, logsTab |
| 4 | Create Container | `CreateContainer.tsx` | ✅ Complete | imageInput, nameInput, createButton, cancelButton |
| 5 | Registry | `Registry.tsx` | ✅ Complete | loadCatalog, list, card, deleteTag |

## ❌ Pages WITHOUT Test IDs (20 pages remaining)

### Core Feature Pages (10)

| # | Page | File | Priority | Interactive Elements |
|---|------|------|----------|----------------------|
| 6 | Security | `Security.tsx` | High | generateKey, signImage, verifySignature, scanImage, keyList, signatureList |
| 7 | Orchestrator | `Orchestrator.tsx` | High | registerNode, createDeployment, scaleDeployment, restartDeployment, deleteDeployment |
| 8 | Runtime | `Runtime.tsx` | High | createContainer, pullImage, startContainer, stopContainer, deleteContainer |
| 9 | AIOps | `AIOps.tsx` | High | recordMetric, predictUsage, getScaling, setPolicy, chatCommand, sendCommand |
| 10 | Container Groups | `ContainerGroups.tsx` | High | createGroup, addContainer, removeContainer, startGroup, stopGroup |
| 11 | Composition | `Composition.tsx` | High | specName, addSource, previewComposition, saveSpec, composeContainer, deleteSpec |
| 12 | Templates | `Templates.tsx` | High | deployTemplate, saveTemplate, deleteTemplate, templateCard |
| 13 | Shared Runtime | `SharedRuntime.tsx` | High | createWorkspace, addMember, createInvite, acceptInvite, createSession, createTunnel, shareWorkspace, joinWorkspace |
| 14 | Snapshots | `Snapshots.tsx` | High | createSnapshot, snapshotName, resourceType, restoreSnapshot, deleteSnapshot |
| 15 | Ephemeral Runtime | `EphemeralRuntime.tsx` | High | provisionRuntime, runtimeName, instanceType, sleepRuntime, wakeRuntime, terminateRuntime |

### Management Pages (7)

| # | Page | File | Priority | Interactive Elements |
|---|------|------|----------|----------------------|
| 16 | Dashboard | `Dashboard.tsx` | Medium | viewContainers, createContainer, manageImages (button links) |
| 17 | Monitor | `Monitor.tsx` | Medium | refresh, filterLogs, metricCard |
| 18 | Networks | `Networks.tsx` | Medium | createNetwork, networkName, networkDriver, connectContainer |
| 19 | Services | `Services.tsx` | Medium | createService, serviceName, registerService |
| 20 | Teams | `Teams.tsx` | Medium | createTeam, teamName, addMember |
| 21 | Tenants | `Tenants.tsx` | Medium | createTenant, tenantName, setQuota |
| 22 | Settings | `Settings.tsx` | Medium | saveSettings, apiUrl |

### Sub-pages (4)

| # | Page | File | Priority | Interactive Elements |
|---|------|------|----------|----------------------|
| 23 | Container Logs | `ContainerLogs.tsx` | Low | (view only, minimal interactions) |
| 24 | Container Env | `ContainerEnv.tsx` | Low | (environment variable management) |
| 25 | Logs | `Logs.tsx` | Low | (log viewing) |
| 26 | Performance | `Performance.tsx` | Low | (performance metrics viewing) |

## ✅ Completed Infrastructure

### Test ID System
- ✅ `test-ids.ts` - Complete constants for all pages
- ✅ Pattern: `feature.action.element`
- ✅ Centralized constants for maintainability

### UI Components Support
- ✅ `Button` - Supports `data-test-id`
- ✅ `Input` - Supports `data-test-id`
- ✅ `Textarea` - Supports `data-test-id`
- ✅ `TabsTrigger` - Supports `data-test-id` (via Radix props)

### API Test-Aware Integration
- ✅ `api-test-aware.ts` - API call tracking
- ✅ Integrated into `ApiClient.request()`
- ✅ Auto-initialized in `main.tsx`
- ✅ Tracks: method, endpoint, request, response, errors, duration
- ✅ Exposes `window.__API_TRACKER__` in test mode

## 📊 Statistics

- **Total Pages:** 25 main pages + 4 sub-pages = 29 total
- **Completed:** 5 pages (17%)
- **Remaining:** 24 pages (83%)
- **Test ID Constants:** ✅ All defined
- **Infrastructure:** ✅ 100% Complete

## 🎯 Next Steps

### Priority Order:
1. **High Priority** (Core Features): Security, Orchestrator, Runtime, AIOps, Groups, Composition, Templates, SharedRuntime, Snapshots, EphemeralRuntime
2. **Medium Priority** (Management): Dashboard, Monitor, Networks, Services, Teams, Tenants, Settings
3. **Low Priority** (Sub-pages): ContainerLogs, ContainerEnv, Logs, Performance

### Implementation Pattern:
```typescript
// 1. Import TEST_IDS
import { TEST_IDS } from '@/lib/test-ids'

// 2. Add to buttons
<Button data-test-id={TEST_IDS.security.generateKey}>Generate Key</Button>

// 3. Add to inputs
<Input data-test-id={TEST_IDS.security.keyName} />

// 4. Add to lists/containers
<div data-test-id={TEST_IDS.security.keyList}>
  {keys.map(...)}
</div>
```

## 📝 Notes

- All test ID constants are already defined in `test-ids.ts`
- UI components (Button, Input, Textarea) already support `data-test-id`
- API tracking is fully integrated
- Remaining work is purely adding attributes to existing components

