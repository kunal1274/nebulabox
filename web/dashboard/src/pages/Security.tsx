import { useEffect, useState } from 'react'
import { Shield, Key, CheckCircle, XCircle, Lock, Scan, Plus, Copy } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { api, SigningKeyInfo, ImageSignature } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function Security() {
  const [keys, setKeys] = useState<SigningKeyInfo[]>([])
  const [loading, setLoading] = useState(true)
  const [newKeyId, setNewKeyId] = useState('')
  const [generating, setGenerating] = useState(false)
  
  // Signing state
  const [signImage, setSignImage] = useState('')
  const [signTag, setSignTag] = useState('latest')
  const [signDigest, setSignDigest] = useState('')
  const [signKeyId, setSignKeyId] = useState('')
  const [signing, setSigning] = useState(false)
  const [signature, setSignature] = useState<ImageSignature | null>(null)
  
  // Verification state
  const [verifyImage, setVerifyImage] = useState('')
  const [verifyTag, setVerifyTag] = useState('latest')
  const [verifyDigest, setVerifyDigest] = useState('')
  const [verifySignature, setVerifySignature] = useState<ImageSignature | null>(null)
  const [verifying, setVerifying] = useState(false)
  const [verifyResult, setVerifyResult] = useState<{ valid: boolean; message?: string } | null>(null)

  useEffect(() => {
    loadKeys()
  }, [])

  const loadKeys = async () => {
    setLoading(true)
    try {
      const data = await api.listSigningKeys()
      setKeys(data.keys)
      if (data.keys.length > 0 && !signKeyId) {
        setSignKeyId(data.keys[0].keyId)
      }
    } catch (error) {
      console.error('Failed to load keys:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleGenerateKey = async () => {
    if (!newKeyId.trim()) {
      alert('Please enter a key ID')
      return
    }
    
    setGenerating(true)
    try {
      const key = await api.generateSigningKey(newKeyId)
      setKeys([...keys, key])
      setNewKeyId('')
      setSignKeyId(key.keyId)
    } catch (error: any) {
      console.error('Failed to generate key:', error)
      alert(`Failed to generate key: ${error.message}`)
    } finally {
      setGenerating(false)
    }
  }

  const handleSign = async () => {
    if (!signImage.trim() || !signDigest.trim() || !signKeyId) {
      alert('Please fill in all required fields')
      return
    }
    
    setSigning(true)
    try {
      const sig = await api.signImage(signImage, signTag, signDigest, signKeyId)
      setSignature(sig)
    } catch (error: any) {
      console.error('Failed to sign image:', error)
      alert(`Failed to sign image: ${error.message}`)
    } finally {
      setSigning(false)
    }
  }

  const handleVerify = async () => {
    if (!verifyImage.trim() || !verifyDigest.trim() || !verifySignature) {
      alert('Please fill in all required fields and provide a signature')
      return
    }
    
    setVerifying(true)
    try {
      const result = await api.verifySignature(verifyImage, verifyTag, verifyDigest, verifySignature)
      setVerifyResult({
        valid: result.valid,
        message: result.valid ? 'Signature is valid ✓' : 'Signature verification failed ✗',
      })
    } catch (error: any) {
      console.error('Failed to verify signature:', error)
      setVerifyResult({
        valid: false,
        message: `Verification error: ${error.message}`,
      })
    } finally {
      setVerifying(false)
    }
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    alert('Copied to clipboard')
  }

  const formatDate = (dateStr: string) => {
    try {
      return new Date(dateStr).toLocaleString()
    } catch {
      return dateStr
    }
  }

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Security</h1>
        <p className="text-muted-foreground">
          Manage signing keys, sign images, and verify signatures
        </p>
      </div>

      <Tabs defaultValue="keys" className="space-y-6">
        <TabsList>
          <TabsTrigger value="keys" data-test-id={TEST_IDS.security.keyList}>
            <Key className="h-4 w-4 mr-2" />
            Signing Keys
          </TabsTrigger>
          <TabsTrigger value="sign" data-test-id={TEST_IDS.security.signImage}>
            <Lock className="h-4 w-4 mr-2" />
            Sign Image
          </TabsTrigger>
          <TabsTrigger value="verify" data-test-id={TEST_IDS.security.verifySignature}>
            <CheckCircle className="h-4 w-4 mr-2" />
            Verify Signature
          </TabsTrigger>
          <TabsTrigger value="scan" data-test-id={TEST_IDS.security.scanImage}>
            <Scan className="h-4 w-4 mr-2" />
            Vulnerability Scanning
          </TabsTrigger>
        </TabsList>

        {/* Signing Keys Tab */}
        <TabsContent value="keys" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Key className="h-5 w-5" />
                Generate New Key
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex gap-2">
                <div className="flex-1">
                  <Input
                    placeholder="Key ID (e.g., production-key-2024)"
                    value={newKeyId}
                    onChange={(e) => setNewKeyId(e.target.value)}
                    data-test-id={TEST_IDS.security.keyName}
                  />
                </div>
                <Button onClick={handleGenerateKey} disabled={generating} data-test-id={TEST_IDS.security.generateKey}>
                  <Plus className="h-4 w-4 mr-2" />
                  {generating ? 'Generating...' : 'Generate Key'}
                </Button>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Existing Keys</CardTitle>
            </CardHeader>
            <CardContent>
              {loading ? (
                <div className="text-muted-foreground">Loading keys...</div>
              ) : keys.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <Key className="h-12 w-12 mx-auto mb-2 opacity-50" />
                  <p>No signing keys found. Generate one to get started.</p>
                </div>
              ) : (
                <div className="space-y-4" data-test-id={TEST_IDS.security.keyList}>
                  {keys.map((key) => (
                    <div key={key.keyId} className="p-4 border rounded-lg">
                      <div className="flex items-start justify-between mb-2">
                        <div>
                          <div className="font-semibold text-lg">{key.keyId}</div>
                          <div className="text-sm text-muted-foreground">
                            Created by {key.createdBy} on {formatDate(key.createdAt)}
                          </div>
                        </div>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => copyToClipboard(key.publicKey)}
                        >
                          <Copy className="h-4 w-4 mr-2" />
                          Copy Public Key
                        </Button>
                      </div>
                      <details className="mt-2">
                        <summary className="cursor-pointer text-sm text-muted-foreground">
                          View Public Key (PEM)
                        </summary>
                        <pre className="mt-2 p-2 bg-muted rounded text-xs font-mono overflow-auto">
                          {key.publicKey}
                        </pre>
                      </details>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Sign Image Tab */}
        <TabsContent value="sign" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Lock className="h-5 w-5" />
                Sign Container Image
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label>Image Name</Label>
                <Input
                  value={signImage}
                  onChange={(e) => setSignImage(e.target.value)}
                  placeholder="myregistry/myapp"
                  data-test-id={TEST_IDS.security.imageRef}
                />
              </div>
              <div>
                <Label>Tag</Label>
                <Input
                  value={signTag}
                  onChange={(e) => setSignTag(e.target.value)}
                  placeholder="latest"
                />
              </div>
              <div>
                <Label>Digest</Label>
                <Input
                  value={signDigest}
                  onChange={(e) => setSignDigest(e.target.value)}
                  placeholder="sha256:abc123..."
                />
              </div>
              <div>
                <Label>Signing Key</Label>
                <select
                  className="w-full p-2 border rounded"
                  value={signKeyId}
                  onChange={(e) => setSignKeyId(e.target.value)}
                >
                  <option value="">Select a key...</option>
                  {keys.map((key) => (
                    <option key={key.keyId} value={key.keyId}>
                      {key.keyId}
                    </option>
                  ))}
                </select>
              </div>
              <Button onClick={handleSign} disabled={signing || !signKeyId} data-test-id={TEST_IDS.security.signImage}>
                <Shield className="h-4 w-4 mr-2" />
                {signing ? 'Signing...' : 'Sign Image'}
              </Button>

              {signature && (
                <div className="mt-4 p-4 bg-green-50 dark:bg-green-950 border border-green-200 dark:border-green-800 rounded">
                  <div className="flex items-center gap-2 mb-2">
                    <CheckCircle className="h-5 w-5 text-green-600 dark:text-green-400" />
                    <span className="font-semibold">Image Signed Successfully</span>
                  </div>
                  <div className="space-y-2 text-sm">
                    <div>
                      <strong>Signed by:</strong> {signature.signedBy}
                    </div>
                    <div>
                      <strong>Signed at:</strong> {formatDate(signature.signedAt)}
                    </div>
                    <div>
                      <strong>Algorithm:</strong> {signature.algorithm}
                    </div>
                    <div>
                      <strong>Key ID:</strong> {signature.keyId}
                    </div>
                    <details className="mt-2">
                      <summary className="cursor-pointer">View Signature</summary>
                      <pre className="mt-2 p-2 bg-white dark:bg-black rounded text-xs font-mono overflow-auto">
                        {JSON.stringify(signature, null, 2)}
                      </pre>
                    </details>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Verify Signature Tab */}
        <TabsContent value="verify" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <CheckCircle className="h-5 w-5" />
                Verify Image Signature
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label>Image Name</Label>
                <Input
                  value={verifyImage}
                  onChange={(e) => setVerifyImage(e.target.value)}
                  placeholder="myregistry/myapp"
                  data-test-id={TEST_IDS.security.imageRef}
                />
              </div>
              <div>
                <Label>Tag</Label>
                <Input
                  value={verifyTag}
                  onChange={(e) => setVerifyTag(e.target.value)}
                  placeholder="latest"
                />
              </div>
              <div>
                <Label>Digest</Label>
                <Input
                  value={verifyDigest}
                  onChange={(e) => setVerifyDigest(e.target.value)}
                  placeholder="sha256:abc123..."
                />
              </div>
              <div>
                <Label>Signature (JSON)</Label>
                <textarea
                  className="w-full p-2 border rounded font-mono text-sm"
                  rows={10}
                  placeholder='{"image": "...", "signature": "...", ...}'
                  value={verifySignature ? JSON.stringify(verifySignature, null, 2) : ''}
                  onChange={(e) => {
                    try {
                      const parsed = JSON.parse(e.target.value)
                      setVerifySignature(parsed)
                    } catch {
                      // Invalid JSON, keep as is
                    }
                  }}
                />
              </div>
              <Button onClick={handleVerify} disabled={verifying || !verifySignature} data-test-id={TEST_IDS.security.verifySignature}>
                <Shield className="h-4 w-4 mr-2" />
                {verifying ? 'Verifying...' : 'Verify Signature'}
              </Button>

              {verifyResult && (
                <div className={`mt-4 p-4 rounded border ${
                  verifyResult.valid
                    ? 'bg-green-50 dark:bg-green-950 border-green-200 dark:border-green-800'
                    : 'bg-red-50 dark:bg-red-950 border-red-200 dark:border-red-800'
                }`}>
                  <div className="flex items-center gap-2">
                    {verifyResult.valid ? (
                      <CheckCircle className="h-5 w-5 text-green-600 dark:text-green-400" />
                    ) : (
                      <XCircle className="h-5 w-5 text-red-600 dark:text-red-400" />
                    )}
                    <span className="font-semibold">
                      {verifyResult.valid ? 'Valid' : 'Invalid'}
                    </span>
                  </div>
                  {verifyResult.message && (
                    <p className="mt-2 text-sm">{verifyResult.message}</p>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Scanning Tab */}
        <TabsContent value="scan" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Scan className="h-5 w-5" />
                Vulnerability Scanning
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <p className="text-sm text-muted-foreground mb-4">
                    Enhanced vulnerability scanning is now available. Use the Images page to scan images.
                    The scanner uses an enhanced vulnerability database with support for:
                  </p>
                  <ul className="list-disc list-inside space-y-1 text-sm text-muted-foreground">
                    <li>Extended CVE database with multiple vulnerability entries</li>
                    <li>Automatic package detection based on image type</li>
                    <li>Severity-based sorting and reporting</li>
                    <li>Integration with registry version metadata</li>
                  </ul>
                </div>
                <Button variant="outline" onClick={() => window.location.href = '/images'} data-test-id={TEST_IDS.security.scanImage}>
                  <Scan className="h-4 w-4 mr-2" />
                  Go to Images Page
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

