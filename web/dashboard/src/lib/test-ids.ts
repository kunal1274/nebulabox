/**
 * Test ID Constants
 * Centralized test IDs for all interactive elements
 * Pattern: feature.action.element
 */

export const TEST_IDS = {
  // Containers Page
  containers: {
    refresh: 'containers.refresh.button',
    create: 'containers.create.button',
    list: 'containers.list.container',
    stop: 'containers.stop.button',
    start: 'containers.start.button',
    logs: 'containers.logs.button',
    card: 'containers.card.container',
  },

  // Images Page
  images: {
    refresh: 'images.refresh.button',
    pull: 'images.pull.button',
    pullInput: 'images.pull.input',
    build: 'images.build.button',
    buildTagInput: 'images.buildTag.input',
    dockerfileInput: 'images.dockerfile.textarea',
    scan: 'images.scan.button',
    push: 'images.push.button',
    list: 'images.list.container',
    card: 'images.card.container',
  },

  // Build Spec Page
  buildspec: {
    specEditor: 'buildspec.editor.textarea',
    tagInput: 'buildspec.tag.input',
    loadExample: 'buildspec.loadExample.button',
    validate: 'buildspec.validate.button',
    convert: 'buildspec.convert.button',
    build: 'buildspec.build.button',
    dockerfileTab: 'buildspec.dockerfile.tab',
    validationTab: 'buildspec.validation.tab',
    logsTab: 'buildspec.logs.tab',
  },

  // Create Container Page
  createContainer: {
    imageInput: 'createContainer.image.input',
    nameInput: 'createContainer.name.input',
    commandInput: 'createContainer.command.input',
    envInput: 'createContainer.env.textarea',
    portsInput: 'createContainer.ports.textarea',
    volumesInput: 'createContainer.volumes.textarea',
    createButton: 'createContainer.create.button',
    cancelButton: 'createContainer.cancel.button',
  },

  // Mode Switcher
  modeSwitcher: {
    mock: 'modeSwitcher.mock.button',
    test: 'modeSwitcher.test.button',
    live: 'modeSwitcher.live.button',
    badge: 'modeSwitcher.badge.container',
  },

  // Dashboard
  dashboard: {
    viewContainers: 'dashboard.viewContainers.button',
    createContainer: 'dashboard.createContainer.button',
    manageImages: 'dashboard.manageImages.button',
  },

  // Sidebar
  sidebar: {
    dashboard: 'sidebar.dashboard.link',
    containers: 'sidebar.containers.link',
    images: 'sidebar.images.link',
    registry: 'sidebar.registry.link',
    buildspec: 'sidebar.buildspec.link',
    security: 'sidebar.security.link',
    orchestrator: 'sidebar.orchestrator.link',
    runtime: 'sidebar.runtime.link',
    aiops: 'sidebar.aiops.link',
    groups: 'sidebar.groups.link',
    composition: 'sidebar.composition.link',
    templates: 'sidebar.templates.link',
    shareruntime: 'sidebar.shareruntime.link',
    snapshots: 'sidebar.snapshots.link',
    ephemeral: 'sidebar.ephemeral.link',
    monitor: 'sidebar.monitor.link',
    networks: 'sidebar.networks.link',
    services: 'sidebar.services.link',
    teams: 'sidebar.teams.link',
    tenants: 'sidebar.tenants.link',
    settings: 'sidebar.settings.link',
  },

  // Registry Page
  registry: {
    login: 'registry.login.button',
    logout: 'registry.logout.button',
    username: 'registry.username.input',
    password: 'registry.password.input',
    loadCatalog: 'registry.loadCatalog.button',
    repository: 'registry.repository.input',
    loadTags: 'registry.loadTags.button',
    retagFrom: 'registry.retagFrom.input',
    retagTo: 'registry.retagTo.input',
    retag: 'registry.retag.button',
    deleteTag: 'registry.deleteTag.button',
    card: 'registry.card.container',
    list: 'registry.list.container',
  },

  // Security Page
  security: {
    generateKey: 'security.generateKey.button',
    keyName: 'security.keyName.input',
    signImage: 'security.signImage.button',
    imageRef: 'security.imageRef.input',
    verifySignature: 'security.verifySignature.button',
    scanImage: 'security.scanImage.button',
    keyList: 'security.keyList.container',
    signatureList: 'security.signatureList.container',
  },

  // Orchestrator Page
  orchestrator: {
    registerNode: 'orchestrator.registerNode.button',
    nodeId: 'orchestrator.nodeId.input',
    nodeName: 'orchestrator.nodeName.input',
    nodeAddress: 'orchestrator.nodeAddress.input',
    createDeployment: 'orchestrator.createDeployment.button',
    deployName: 'orchestrator.deployName.input',
    deployImage: 'orchestrator.deployImage.input',
    deployReplicas: 'orchestrator.deployReplicas.input',
    scaleDeployment: 'orchestrator.scaleDeployment.button',
    restartDeployment: 'orchestrator.restartDeployment.button',
    deleteDeployment: 'orchestrator.deleteDeployment.button',
  },

  // Runtime Page
  runtime: {
    createContainer: 'runtime.createContainer.button',
    containerId: 'runtime.containerId.input',
    containerImage: 'runtime.containerImage.input',
    pullImage: 'runtime.pullImage.button',
    imageName: 'runtime.imageName.input',
    startContainer: 'runtime.startContainer.button',
    stopContainer: 'runtime.stopContainer.button',
    deleteContainer: 'runtime.deleteContainer.button',
  },

  // AIOps Page
  aiops: {
    recordMetric: 'aiops.recordMetric.button',
    containerId: 'aiops.containerId.input',
    predictUsage: 'aiops.predictUsage.button',
    getScaling: 'aiops.getScaling.button',
    setPolicy: 'aiops.setPolicy.button',
    chatCommand: 'aiops.chatCommand.input',
    sendCommand: 'aiops.sendCommand.button',
  },

  // Groups Page
  groups: {
    createGroup: 'groups.createGroup.button',
    groupName: 'groups.groupName.input',
    addContainer: 'groups.addContainer.button',
    removeContainer: 'groups.removeContainer.button',
    startGroup: 'groups.startGroup.button',
    stopGroup: 'groups.stopGroup.button',
  },

  // Composition Page
  composition: {
    specName: 'composition.specName.input',
    addSource: 'composition.addSource.button',
    previewComposition: 'composition.previewComposition.button',
    saveSpec: 'composition.saveSpec.button',
    composeContainer: 'composition.composeContainer.button',
    deleteSpec: 'composition.deleteSpec.button',
  },

  // Templates Page
  templates: {
    deployTemplate: 'templates.deployTemplate.button',
    saveTemplate: 'templates.saveTemplate.button',
    deleteTemplate: 'templates.deleteTemplate.button',
    templateCard: 'templates.templateCard.container',
  },

  // SharedRuntime Page
  shareruntime: {
    createWorkspace: 'shareruntime.createWorkspace.button',
    workspaceName: 'shareruntime.workspaceName.input',
    addMember: 'shareruntime.addMember.button',
    createInvite: 'shareruntime.createInvite.button',
    acceptInvite: 'shareruntime.acceptInvite.button',
    createSession: 'shareruntime.createSession.button',
    createTunnel: 'shareruntime.createTunnel.button',
    shareWorkspace: 'shareruntime.shareWorkspace.button',
    joinWorkspace: 'shareruntime.joinWorkspace.button',
  },

  // Snapshots Page
  snapshots: {
    createSnapshot: 'snapshots.createSnapshot.button',
    snapshotName: 'snapshots.snapshotName.input',
    resourceType: 'snapshots.resourceType.select',
    restoreSnapshot: 'snapshots.restoreSnapshot.button',
    deleteSnapshot: 'snapshots.deleteSnapshot.button',
  },

  // EphemeralRuntime Page
  ephemeral: {
    provisionRuntime: 'ephemeral.provisionRuntime.button',
    runtimeName: 'ephemeral.runtimeName.input',
    instanceType: 'ephemeral.instanceType.select',
    sleepRuntime: 'ephemeral.sleepRuntime.button',
    wakeRuntime: 'ephemeral.wakeRuntime.button',
    terminateRuntime: 'ephemeral.terminateRuntime.button',
  },

  // Monitor Page
  monitor: {
    refresh: 'monitor.refresh.button',
    filterLogs: 'monitor.filterLogs.input',
    metricCard: 'monitor.metricCard.container',
  },

  // Networks Page
  networks: {
    createNetwork: 'networks.createNetwork.button',
    networkName: 'networks.networkName.input',
    networkDriver: 'networks.networkDriver.select',
    connectContainer: 'networks.connectContainer.button',
  },

  // Services Page
  services: {
    createService: 'services.createService.button',
    serviceName: 'services.serviceName.input',
    registerService: 'services.registerService.button',
  },

  // Teams Page
  teams: {
    createTeam: 'teams.createTeam.button',
    teamName: 'teams.teamName.input',
    addMember: 'teams.addMember.button',
  },

  // Tenants Page
  tenants: {
    createTenant: 'tenants.createTenant.button',
    tenantName: 'tenants.tenantName.input',
    setQuota: 'tenants.setQuota.button',
  },

  // Settings Page
  settings: {
    saveSettings: 'settings.saveSettings.button',
    apiUrl: 'settings.apiUrl.input',
  },
} as const

/**
 * Generate test ID helper
 */
export function getTestId(path: string): string {
  const parts = path.split('.')
  let current: any = TEST_IDS
  for (const part of parts) {
    if (current[part]) {
      current = current[part]
    } else {
      return path
    }
  }
  return typeof current === 'string' ? current : path
}

