import { useEffect, useState } from 'react'
import { Brain, TrendingUp, TrendingDown, AlertTriangle, MessageSquare, Send, BarChart3, Settings, Zap } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { api, ResourcePrediction, ScalingRecommendation, ChatCommand, ScalingPolicy } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function AIOps() {
  const [activeTab, setActiveTab] = useState('predictions')

  // Predictions state
  const [containerId, setContainerId] = useState('')
  const [prediction, setPrediction] = useState<ResourcePrediction | null>(null)
  const [predicting, setPredicting] = useState(false)

  // Scaling state
  const [targetId, setTargetId] = useState('')
  const [currentReplicas, setCurrentReplicas] = useState(1)
  const [scalingRec, setScalingRec] = useState<ScalingRecommendation | null>(null)
  const [loadingScaling, setLoadingScaling] = useState(false)

  // Scaling policy state
  const [policy, setPolicy] = useState<ScalingPolicy>({
    targetId: '',
    type: 'container',
    minReplicas: 1,
    maxReplicas: 10,
    scaleUpThreshold: 80,
    scaleDownThreshold: 20,
    cooldownPeriod: '5m',
  })
  const [savingPolicy, setSavingPolicy] = useState(false)

  // ChatOps state
  const [chatCommand, setChatCommand] = useState('')
  const [chatHistory, setChatHistory] = useState<ChatCommand[]>([])
  const [processingChat, setProcessingChat] = useState(false)

  useEffect(() => {
    loadChatHistory()
  }, [])

  const loadChatHistory = async () => {
    try {
      const data = await api.getChatHistory(20)
      setChatHistory(data.commands)
    } catch (error) {
      console.error('Failed to load chat history:', error)
    }
  }

  const handlePredict = async () => {
    if (!containerId.trim()) {
      alert('Please enter a container ID')
      return
    }

    setPredicting(true)
    try {
      const pred = await api.predictResourceUsage(containerId, '30m')
      setPrediction(pred)
    } catch (error: any) {
      console.error('Failed to predict:', error)
      alert(`Failed to predict: ${error.message}`)
    } finally {
      setPredicting(false)
    }
  }

  const handleGetScalingRecommendation = async () => {
    if (!targetId.trim()) {
      alert('Please enter a target ID')
      return
    }

    setLoadingScaling(true)
    try {
      const rec = await api.getScalingRecommendation(targetId, currentReplicas)
      setScalingRec(rec)
    } catch (error: any) {
      console.error('Failed to get scaling recommendation:', error)
      alert(`Failed to get recommendation: ${error.message}`)
    } finally {
      setLoadingScaling(false)
    }
  }

  const handleSavePolicy = async () => {
    if (!policy.targetId.trim()) {
      alert('Please enter a target ID')
      return
    }

    setSavingPolicy(true)
    try {
      await api.setScalingPolicy(policy)
      alert('Scaling policy saved successfully')
      setPolicy({
        targetId: '',
        type: 'container',
        minReplicas: 1,
        maxReplicas: 10,
        scaleUpThreshold: 80,
        scaleDownThreshold: 20,
        cooldownPeriod: '5m',
      })
    } catch (error: any) {
      console.error('Failed to save policy:', error)
      alert(`Failed to save policy: ${error.message}`)
    } finally {
      setSavingPolicy(false)
    }
  }

  const handleChatCommand = async () => {
    if (!chatCommand.trim()) {
      return
    }

    setProcessingChat(true)
    try {
      const response = await api.processChatCommand(chatCommand)
      
      // Add command to history
      const newCommand: ChatCommand = {
        command: chatCommand,
        timestamp: new Date().toISOString(),
        processed: true,
        success: response.success,
        response: response.message,
      }
      setChatHistory([...chatHistory, newCommand])
      setChatCommand('')
    } catch (error: any) {
      console.error('Failed to process command:', error)
      const errorCommand: ChatCommand = {
        command: chatCommand,
        timestamp: new Date().toISOString(),
        processed: true,
        success: false,
        response: `Error: ${error.message}`,
      }
      setChatHistory([...chatHistory, errorCommand])
    } finally {
      setProcessingChat(false)
    }
  }

  const formatDate = (dateStr: string) => {
    try {
      return new Date(dateStr).toLocaleString()
    } catch {
      return dateStr
    }
  }

  const getTrendIcon = (trend: string) => {
    switch (trend.toLowerCase()) {
      case 'increasing':
        return <TrendingUp className="h-4 w-4 text-red-500" />
      case 'decreasing':
        return <TrendingDown className="h-4 w-4 text-green-500" />
      default:
        return <BarChart3 className="h-4 w-4 text-blue-500" />
    }
  }

  const getActionColor = (action: string) => {
    switch (action) {
      case 'scale_up':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300'
      case 'scale_down':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
      case 'none':
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
    }
  }

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">AI-Driven Operations</h1>
        <p className="text-muted-foreground">
          Predictive analytics, intelligent scaling, and ChatOps for automated container management
        </p>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList>
          <TabsTrigger value="predictions">
            <BarChart3 className="h-4 w-4 mr-2" />
            Predictions
          </TabsTrigger>
          <TabsTrigger value="scaling">
            <Zap className="h-4 w-4 mr-2" />
            Auto-Scaling
          </TabsTrigger>
          <TabsTrigger value="chatops">
            <MessageSquare className="h-4 w-4 mr-2" />
            ChatOps
          </TabsTrigger>
        </TabsList>

        {/* Predictions Tab */}
        <TabsContent value="predictions" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Brain className="h-5 w-5" />
                Resource Usage Prediction
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex gap-2 mb-4">
                <Input
                  placeholder="Container ID"
                  value={containerId}
                  onChange={(e) => setContainerId(e.target.value)}
                  data-test-id={TEST_IDS.aiops.containerId}
                />
                <Button onClick={handlePredict} disabled={predicting} data-test-id={TEST_IDS.aiops.predictUsage}>
                  {predicting ? 'Predicting...' : 'Predict Usage'}
                </Button>
              </div>

              {prediction && (
                <div className="space-y-4 mt-6">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="p-4 border rounded-lg">
                      <div className="flex items-center gap-2 mb-2">
                        <span className="font-semibold">CPU Prediction</span>
                        {getTrendIcon(prediction.cpuTrend)}
                      </div>
                      <div className="text-2xl font-bold">{prediction.predictedCpu.toFixed(1)}%</div>
                      <div className="text-sm text-muted-foreground">Trend: {prediction.cpuTrend}</div>
                    </div>
                    <div className="p-4 border rounded-lg">
                      <div className="flex items-center gap-2 mb-2">
                        <span className="font-semibold">Memory Prediction</span>
                        {getTrendIcon(prediction.memoryTrend)}
                      </div>
                      <div className="text-2xl font-bold">{prediction.predictedMemory.toFixed(1)}%</div>
                      <div className="text-sm text-muted-foreground">Trend: {prediction.memoryTrend}</div>
                    </div>
                  </div>

                  <div className="p-4 border rounded-lg">
                    <div className="font-semibold mb-2">Confidence Score</div>
                    <div className="w-full bg-gray-200 rounded-full h-2.5">
                      <div
                        className="bg-blue-600 h-2.5 rounded-full"
                        style={{ width: `${prediction.confidence * 100}%` }}
                      ></div>
                    </div>
                    <div className="text-sm text-muted-foreground mt-1">
                      {(prediction.confidence * 100).toFixed(0)}% confidence
                    </div>
                  </div>

                  {prediction.anomalies && prediction.anomalies.length > 0 && (
                    <div className="p-4 border rounded-lg border-yellow-200 bg-yellow-50 dark:bg-yellow-900/20">
                      <div className="flex items-center gap-2 mb-2">
                        <AlertTriangle className="h-4 w-4 text-yellow-600" />
                        <span className="font-semibold">Detected Anomalies</span>
                      </div>
                      <div className="space-y-2">
                        {prediction.anomalies.map((anomaly, idx) => (
                          <div key={idx} className="text-sm">
                            <Badge variant="outline" className="mr-2">{anomaly.severity}</Badge>
                            {anomaly.message}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {prediction.recommendations && prediction.recommendations.length > 0 && (
                    <div className="p-4 border rounded-lg">
                      <div className="font-semibold mb-2">Recommendations</div>
                      <div className="space-y-2">
                        {prediction.recommendations.map((rec, idx) => (
                          <div key={idx} className="p-3 bg-blue-50 dark:bg-blue-900/20 rounded">
                            <div className="flex items-center gap-2 mb-1">
                              <span className="font-semibold">{rec.title}</span>
                              <Badge>{rec.priority}</Badge>
                            </div>
                            <div className="text-sm text-muted-foreground">{rec.description}</div>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Auto-Scaling Tab */}
        <TabsContent value="scaling" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Zap className="h-5 w-5" />
                Scaling Recommendations
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <Label>Target ID</Label>
                  <Input
                    value={targetId}
                    onChange={(e) => setTargetId(e.target.value)}
                    placeholder="container-id or deployment-id"
                    data-test-id={TEST_IDS.aiops.containerId}
                  />
                </div>
                <div>
                  <Label>Current Replicas</Label>
                  <Input
                    type="number"
                    value={currentReplicas}
                    onChange={(e) => setCurrentReplicas(parseInt(e.target.value) || 1)}
                  />
                </div>
              </div>
              <Button onClick={handleGetScalingRecommendation} disabled={loadingScaling} className="mb-4" data-test-id={TEST_IDS.aiops.getScaling}>
                Get Scaling Recommendation
              </Button>

              {scalingRec && (
                <div className="mt-6 p-4 border rounded-lg">
                  <div className="flex items-center gap-2 mb-4">
                    <Badge className={getActionColor(scalingRec.action)}>
                      {scalingRec.action.replace('_', ' ').toUpperCase()}
                    </Badge>
                    <span className="text-sm text-muted-foreground">Confidence: {(scalingRec.confidence * 100).toFixed(0)}%</span>
                  </div>
                  
                  <div className="grid grid-cols-3 gap-4 mb-4">
                    <div>
                      <div className="text-sm text-muted-foreground">Current</div>
                      <div className="text-xl font-bold">{scalingRec.currentReplicas}</div>
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground">Recommended</div>
                      <div className="text-xl font-bold">{scalingRec.recommendedReplicas}</div>
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground">Reason</div>
                      <div className="text-sm">{scalingRec.reason}</div>
                    </div>
                  </div>

                  <div className="mb-4">
                    <div className="text-sm font-semibold mb-1">Predicted Usage</div>
                    <div className="grid grid-cols-2 gap-2 text-sm">
                      <div>CPU: {scalingRec.predictedCpu.toFixed(1)}%</div>
                      <div>Memory: {scalingRec.predictedMemory.toFixed(1)}%</div>
                    </div>
                  </div>

                  <div className="text-sm text-muted-foreground">{scalingRec.message}</div>
                </div>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Settings className="h-5 w-5" />
                Scaling Policy
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <Label>Target ID</Label>
                  <Input
                    value={policy.targetId}
                    onChange={(e) => setPolicy({ ...policy, targetId: e.target.value })}
                    placeholder="container-id or deployment-id"
                  />
                </div>
                <div>
                  <Label>Type</Label>
                  <select
                    value={policy.type}
                    onChange={(e) => setPolicy({ ...policy, type: e.target.value })}
                    className="w-full px-3 py-2 border rounded-md"
                  >
                    <option value="container">Container</option>
                    <option value="deployment">Deployment</option>
                  </select>
                </div>
                <div>
                  <Label>Min Replicas</Label>
                  <Input
                    type="number"
                    value={policy.minReplicas}
                    onChange={(e) => setPolicy({ ...policy, minReplicas: parseInt(e.target.value) || 1 })}
                  />
                </div>
                <div>
                  <Label>Max Replicas</Label>
                  <Input
                    type="number"
                    value={policy.maxReplicas}
                    onChange={(e) => setPolicy({ ...policy, maxReplicas: parseInt(e.target.value) || 10 })}
                  />
                </div>
                <div>
                  <Label>Scale Up Threshold (%)</Label>
                  <Input
                    type="number"
                    value={policy.scaleUpThreshold}
                    onChange={(e) => setPolicy({ ...policy, scaleUpThreshold: parseFloat(e.target.value) || 80 })}
                  />
                </div>
                <div>
                  <Label>Scale Down Threshold (%)</Label>
                  <Input
                    type="number"
                    value={policy.scaleDownThreshold}
                    onChange={(e) => setPolicy({ ...policy, scaleDownThreshold: parseFloat(e.target.value) || 20 })}
                  />
                </div>
              </div>
              <Button onClick={handleSavePolicy} disabled={savingPolicy} data-test-id={TEST_IDS.aiops.setPolicy}>
                {savingPolicy ? 'Saving...' : 'Save Policy'}
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        {/* ChatOps Tab */}
        <TabsContent value="chatops" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <MessageSquare className="h-5 w-5" />
                ChatOps Interface
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex gap-2">
                  <Input
                    placeholder="Ask about scaling, predictions, or status..."
                    value={chatCommand}
                    onChange={(e) => setChatCommand(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && !processingChat && handleChatCommand()}
                    data-test-id={TEST_IDS.aiops.chatCommand}
                  />
                  <Button onClick={handleChatCommand} disabled={processingChat || !chatCommand.trim()} data-test-id={TEST_IDS.aiops.sendCommand}>
                    <Send className="h-4 w-4" />
                  </Button>
                </div>

                <div className="text-sm text-muted-foreground">
                  Try commands like: "Show scaling recommendations", "Predict resource usage", "Get container status"
                </div>

                <div className="border rounded-lg p-4 h-96 overflow-y-auto space-y-3">
                  {chatHistory.length === 0 ? (
                    <div className="text-center text-muted-foreground py-8">
                      <MessageSquare className="h-12 w-12 mx-auto mb-2 opacity-50" />
                      <p>No chat history. Start a conversation!</p>
                    </div>
                  ) : (
                    chatHistory.map((cmd, idx) => (
                      <div key={idx} className="space-y-1">
                        <div className="flex items-center gap-2">
                          <span className="text-xs text-muted-foreground">
                            {formatDate(cmd.timestamp)}
                          </span>
                          <Badge variant={cmd.success ? "default" : "destructive"}>
                            {cmd.success ? 'Success' : 'Error'}
                          </Badge>
                        </div>
                        <div className="p-2 bg-gray-50 dark:bg-gray-900 rounded">
                          <div className="font-semibold text-sm mb-1">You:</div>
                          <div className="text-sm">{cmd.command}</div>
                        </div>
                        <div className="p-2 bg-blue-50 dark:bg-blue-900/20 rounded">
                          <div className="font-semibold text-sm mb-1">AI:</div>
                          <div className="text-sm">{cmd.response}</div>
                        </div>
                      </div>
                    ))
                  )}
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

